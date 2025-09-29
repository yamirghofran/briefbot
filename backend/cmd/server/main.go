package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/yamirghofran/briefbot/internal/db"
	"github.com/yamirghofran/briefbot/internal/handlers"
	"github.com/yamirghofran/briefbot/internal/services"
)

func main() {
	// Load .env file into the environment
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Get database URL from environment or use default
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/briefbot?sslmode=disable"
	}

	// Create connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	// Configure connection pool for better concurrency handling
	poolConfig.MaxConns = 25                      // Maximum number of connections (increased for workers)
	poolConfig.MinConns = 5                       // Minimum number of connections to maintain
	poolConfig.MaxConnLifetime = time.Hour        // Maximum lifetime of a connection
	poolConfig.MaxConnIdleTime = 30 * time.Minute // Maximum idle time
	poolConfig.HealthCheckPeriod = time.Minute    // How often to check connection health

	// Additional settings to prevent connection issues
	poolConfig.ConnConfig.ConnectTimeout = 30 * time.Second // Connection timeout

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	fmt.Println("Successfully connected to database with connection pool")

	// Initialize querier
	querier := db.New(pool)

	oaiClient := openai.NewClient(
		option.WithBaseURL("https://api.groq.com/openai/v1"),
		option.WithAPIKey(os.Getenv("GROQ_API_KEY")),
	)

	// Initialize R2 service configuration
	r2Config := services.R2Config{
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		PublicHost:      os.Getenv("R2_PUBLIC_HOST"),
	}

	// Create R2 service if configuration is complete
	var r2Service *services.R2Service
	if r2Config.AccessKeyID != "" && r2Config.SecretAccessKey != "" && r2Config.AccountID != "" && r2Config.BucketName != "" {
		r2Service, err = services.NewR2Service(r2Config)
		if err != nil {
			log.Printf("Warning: R2 service not initialized: %v", err)
			r2Service = nil
		} else {
			log.Println("R2 service initialized successfully")
		}
	} else {
		log.Printf("Warning: R2 service not configured - missing environment variables")
		r2Service = nil
	}

	// Initialize services
	aiService, err := services.NewAIService(&oaiClient)
	if err != nil {
		log.Fatal("Unable to start AI service")
	}
	scrapingService := services.NewScraper()
	userService := services.NewUserService(querier)
	jobQueueService := services.NewJobQueueService(querier)
	itemService := services.NewItemService(querier, aiService, scrapingService, jobQueueService)

	// Initialize podcast service
	podcastConfig := services.DefaultPodcastConfig()
	podcastService := services.NewPodcastService(querier, aiService, nil, r2Service, podcastConfig)

	// Initialize email service
	emailService, err := services.NewEmailService()
	if err != nil {
		log.Printf("Warning: Email service not initialized: %v", err)
		emailService = nil
	}

	// Initialize unified digest service (replaces separate daily and integrated digest services)
	var digestService services.DigestService
	if emailService != nil {
		digestService = services.NewDigestService(querier, emailService, podcastService)
		log.Println("Unified digest service initialized successfully")
	} else {
		log.Printf("Warning: Digest service not initialized - email service not available")
	}

	// Initialize speech service for podcast audio generation
	var speechService services.SpeechService
	falAPIKey := os.Getenv("FAL_API_KEY")
	if falAPIKey != "" {
		falClient := services.NewFalClient(falAPIKey)
		// Configure for indefinite polling with concurrent processing
		// maxAttempts is now ignored - we'll poll until COMPLETED or FAILED
		// 3 second interval provides good balance between responsiveness and API load
		speechService = services.NewSpeechService(falClient, 0, 3*time.Second) // 0 = unlimited attempts
		log.Println("Speech service initialized with Fal client (unlimited attempts, 3s interval, max 5 concurrent)")
	} else {
		log.Printf("Warning: Speech service not configured - FAL_API_KEY environment variable not set")
		speechService = services.NewSpeechService(nil, 0, 3*time.Second)
	}

	// Get max concurrent requests from environment (optional)
	maxConcurrent := 5 // Default
	if maxConcurrentStr := os.Getenv("MAX_CONCURRENT_AUDIO_REQUESTS"); maxConcurrentStr != "" {
		if val, err := strconv.Atoi(maxConcurrentStr); err == nil && val > 0 {
			maxConcurrent = val
			log.Printf("Max concurrent audio requests set to %d from environment", maxConcurrent)
		}
	}

	// Update podcast service with speech service and max concurrent setting
	podcastConfig.MaxConcurrentAudio = int32(maxConcurrent)
	podcastService = services.NewPodcastService(querier, aiService, speechService, r2Service, podcastConfig)

	// Initialize worker service
	workerConfig := services.WorkerConfig{
		WorkerCount:    2,               // Number of concurrent workers
		PollInterval:   5 * time.Second, // How often to check for new jobs
		MaxRetries:     3,               // Max retries for failed jobs
		BatchSize:      10,              // Number of items to process per batch
		EnablePodcasts: true,            // Enable podcast processing
	}
	workerService := services.NewWorkerService(jobQueueService, aiService, scrapingService, podcastService, workerConfig)

	// Start worker service in background
	go func() {
		if err := workerService.Start(context.Background()); err != nil {
			log.Printf("Failed to start worker service: %v", err)
		}
	}()

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	handlers.SetupRoutes(router, userService, itemService, digestService, podcastService)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Background worker service is running")

	// Handle graceful shutdown
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	// Stop worker service
	if err := workerService.Stop(); err != nil {
		log.Printf("Error stopping worker service: %v", err)
	}

	// Shutdown server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
