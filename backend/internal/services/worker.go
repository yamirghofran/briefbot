package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/yamirghofran/briefbot/internal/db"
)

type WorkerService interface {
	Start(ctx context.Context) error
	Stop() error
	IsRunning() bool
}

type workerService struct {
	jobQueueService JobQueueService
	aiService       AIService
	scrapingService ScrapingService
	podcastService  PodcastService

	// Configuration
	workerCount    int
	pollInterval   time.Duration
	maxRetries     int
	batchSize      int32
	enablePodcasts bool

	// Runtime state
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	running   bool
	runningMu sync.Mutex
}

type WorkerConfig struct {
	WorkerCount    int
	PollInterval   time.Duration
	MaxRetries     int
	BatchSize      int32
	EnablePodcasts bool
}

func NewWorkerService(
	jobQueueService JobQueueService,
	aiService AIService,
	scrapingService ScrapingService,
	podcastService PodcastService,
	config WorkerConfig,
) WorkerService {
	if config.WorkerCount <= 0 {
		config.WorkerCount = 2
	}
	if config.PollInterval <= 0 {
		config.PollInterval = 5 * time.Second
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 10
	}

	return &workerService{
		jobQueueService: jobQueueService,
		aiService:       aiService,
		scrapingService: scrapingService,
		podcastService:  podcastService,
		workerCount:     config.WorkerCount,
		pollInterval:    config.PollInterval,
		maxRetries:      config.MaxRetries,
		batchSize:       config.BatchSize,
		enablePodcasts:  config.EnablePodcasts,
	}
}

func (s *workerService) Start(ctx context.Context) error {
	s.runningMu.Lock()
	if s.running {
		s.runningMu.Unlock()
		return fmt.Errorf("worker service is already running")
	}
	s.runningMu.Unlock()

	s.ctx, s.cancel = context.WithCancel(ctx)

	// Start worker goroutines
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker(i + 1)
	}

	s.runningMu.Lock()
	s.running = true
	s.runningMu.Unlock()

	log.Printf("Worker service started with %d workers", s.workerCount)
	return nil
}

func (s *workerService) Stop() error {
	s.runningMu.Lock()
	if !s.running {
		s.runningMu.Unlock()
		return fmt.Errorf("worker service is not running")
	}
	s.runningMu.Unlock()

	log.Println("Stopping worker service...")
	s.cancel()
	s.wg.Wait()

	s.runningMu.Lock()
	s.running = false
	s.runningMu.Unlock()

	log.Println("Worker service stopped")
	return nil
}

func (s *workerService) IsRunning() bool {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()
	return s.running
}

func (s *workerService) worker(id int) {
	defer s.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("Worker %d stopping", id)
			return
		default:
			// Process a batch of items
			if err := s.processBatch(); err != nil {
				log.Printf("Worker %d error processing batch: %v", id, err)
			}

			// Wait before next poll
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(s.pollInterval):
				// Continue to next iteration
			}
		}
	}
}

// Code smell: hardcoded worker service (OCP violation)
// Solution: Strategy pattern

// func (s *workerService) processBatch() error {
// 	// Process items first
// 	if err := s.processItemBatch(); err != nil {
// 		return fmt.Errorf("failed to process item batch: %w", err)
// 	}

// 	// Process podcasts if enabled
// 	if s.enablePodcasts && s.podcastService != nil {
// 		if err := s.processPodcastBatch(); err != nil {
// 			return fmt.Errorf("failed to process podcast batch: %w", err)
// 		}
// 	}

// 	return nil
// }

func (s *workerService) processBatch() error {
	// Process items first
	if err := s.processItemBatch(); err != nil {
		return fmt.Errorf("failed to process item batch: %w", err)
	}

	// Process podcasts if enabled
	if s.enablePodcasts && s.podcastService != nil {
		if err := s.processPodcastBatch(); err != nil {
			return fmt.Errorf("failed to process podcast batch: %w", err)
		}
	}

	return nil
}

func (s *workerService) processItemBatch() error {
	// Get pending items
	items, err := s.jobQueueService.DequeuePendingItems(s.ctx, s.batchSize)
	if err != nil {
		return fmt.Errorf("failed to dequeue items: %w", err)
	}

	if len(items) == 0 {
		// No items to process, sleep longer
		time.Sleep(s.pollInterval)
		return nil
	}

	log.Printf("Processing batch of %d items", len(items))

	// Process each item
	for _, item := range items {
		if err := s.processItem(s.ctx, item); err != nil {
			log.Printf("Failed to process item %d: %v", item.ID, err)
			// Continue with next item even if one fails
		}
	}

	return nil
}

func (s *workerService) processItem(ctx context.Context, item db.Item) error {
	// Mark item as processing
	if err := s.jobQueueService.MarkItemAsProcessing(ctx, item.ID); err != nil {
		return fmt.Errorf("failed to mark item as processing: %w", err)
	}

	log.Printf("Processing item %d: %s", item.ID, *item.Url)

	// Process the URL with retry logic
	var textContent string
	var extraction ItemExtraction
	var summary string
	var err error

	// Retry processing up to maxRetries times
	for attempt := 1; attempt <= s.maxRetries; attempt++ {
		textContent, extraction, summary, err = s.processURL(ctx, *item.Url)
		if err == nil {
			break // Success!
		}

		log.Printf("Attempt %d failed for item %d: %v", attempt, item.ID, err)

		if attempt < s.maxRetries {
			// Wait before retry with exponential backoff
			backoffDuration := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoffDuration):
				// Continue to next attempt
			}
		}
	}

	if err != nil {
		// All retries failed, mark as failed
		errorMsg := fmt.Sprintf("Failed after %d attempts: %v", s.maxRetries, err)
		if failErr := s.jobQueueService.FailItem(ctx, item.ID, errorMsg); failErr != nil {
			log.Printf("Failed to mark item %d as failed: %v", item.ID, failErr)
		}
		return fmt.Errorf("failed to process URL after %d attempts: %w", s.maxRetries, err)
	}

	// Mark as completed with AI-extracted title
	if err := s.jobQueueService.CompleteItem(ctx, item.ID, extraction.Title, textContent, summary, extraction.Type, extraction.Platform, extraction.Tags, extraction.Authors); err != nil {
		return fmt.Errorf("failed to complete item: %w", err)
	}

	log.Printf("Successfully processed item %d", item.ID)
	return nil
}

func (s *workerService) processURL(ctx context.Context, url string) (string, ItemExtraction, string, error) {
	// Scrape content
	content, err := s.scrapingService.Scrape(url)
	if err != nil {
		return "", ItemExtraction{}, "", fmt.Errorf("failed to scrape URL: %w", err)
	}

	// Extract metadata
	extraction, err := s.aiService.ExtractContent(ctx, content)
	if err != nil {
		return "", ItemExtraction{}, "", fmt.Errorf("failed to extract content: %w", err)
	}

	// Summarize content
	summary, err := s.aiService.SummarizeContent(ctx, content)
	if err != nil {
		return "", ItemExtraction{}, "", fmt.Errorf("failed to summarize content: %w", err)
	}

	// Concatenate summary
	concatenatedSummary := ConcatenateSummary(summary)

	return content, extraction, concatenatedSummary, nil
}

func (s *workerService) processPodcastBatch() error {
	// Get pending podcasts with atomic acquisition to prevent multiple workers from processing the same podcast
	pendingPodcasts, err := s.podcastService.AcquirePendingPodcasts(s.ctx, s.batchSize)
	if err != nil {
		return fmt.Errorf("failed to acquire pending podcasts: %w", err)
	}

	if len(pendingPodcasts) == 0 {
		// No podcasts to process, sleep longer
		time.Sleep(s.pollInterval)
		return nil
	}

	log.Printf("Processing batch of %d podcasts", len(pendingPodcasts))

	// Process each podcast
	for _, podcast := range pendingPodcasts {
		if err := s.processPodcast(s.ctx, podcast); err != nil {
			log.Printf("Failed to process podcast %d: %v", podcast.ID, err)
			// Continue with next podcast even if one fails
		}
	}

	return nil
}

func (s *workerService) processPodcast(ctx context.Context, podcast db.Podcast) error {
	log.Printf("Processing podcast %d: %s", podcast.ID, podcast.Title)

	// Process the podcast with retry logic
	var err error

	// Retry processing up to maxRetries times
	for attempt := 1; attempt <= s.maxRetries; attempt++ {
		err = s.podcastService.ProcessPodcast(ctx, podcast.ID)
		if err == nil {
			break // Success!
		}

		log.Printf("Attempt %d failed for podcast %d: %v", attempt, podcast.ID, err)

		if attempt < s.maxRetries {
			// Wait before retry with exponential backoff
			backoffDuration := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoffDuration):
				// Continue to next attempt
			}
		}
	}

	if err != nil {
		// All retries failed, mark as failed
		if failErr := s.podcastService.UpdatePodcastStatus(ctx, podcast.ID, PodcastStatusFailed); failErr != nil {
			log.Printf("Failed to mark podcast %d as failed: %v", podcast.ID, failErr)
		}
		return fmt.Errorf("failed to process podcast after %d attempts: %w", s.maxRetries, err)
	}

	log.Printf("Successfully processed podcast %d", podcast.ID)
	return nil
}
