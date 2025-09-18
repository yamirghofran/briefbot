package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

	// Connect to database
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Test the connection
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	fmt.Println("Successfully connected to database")

	// Initialize querier
	querier := db.New(conn)

	oaiClient := openai.NewClient(
		option.WithBaseURL("https://api.groq.com/openai/v1"),
		option.WithAPIKey(os.Getenv("GROQ_API_KEY")),
	)

	// Initialize services
	aiService, err := services.NewAIService(&oaiClient)
	if err != nil {
		log.Fatal("Unable to start AI service")
	}
	scrapingService := services.NewScraper()
	userService := services.NewUserService(querier)
	jobQueueService := services.NewJobQueueService(querier)
	itemService := services.NewItemService(querier, aiService, scrapingService, jobQueueService)

	// Initialize email service
	emailService, err := services.NewEmailService()
	if err != nil {
		log.Printf("Warning: Email service not initialized: %v", err)
		emailService = nil
	}

	// Initialize daily digest service
	var dailyDigestService services.DailyDigestService
	if emailService != nil {
		dailyDigestService = services.NewDailyDigestService(querier, emailService)
	}

	// Initialize worker service
	workerConfig := services.WorkerConfig{
		WorkerCount:  2,               // Number of concurrent workers
		PollInterval: 5 * time.Second, // How often to check for new jobs
		MaxRetries:   3,               // Max retries for failed jobs
		BatchSize:    10,              // Number of items to process per batch
	}
	workerService := services.NewWorkerService(jobQueueService, aiService, scrapingService, workerConfig)

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
	handlers.SetupRoutes(router, userService, itemService, dailyDigestService)

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
