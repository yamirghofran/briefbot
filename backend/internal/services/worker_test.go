package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
)

func TestNewWorkerService(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	tests := []struct {
		name           string
		config         WorkerConfig
		expectedConfig WorkerConfig
	}{
		{
			name: "valid config",
			config: WorkerConfig{
				WorkerCount:    5,
				PollInterval:   10 * time.Second,
				MaxRetries:     5,
				BatchSize:      20,
				EnablePodcasts: true,
			},
			expectedConfig: WorkerConfig{
				WorkerCount:    5,
				PollInterval:   10 * time.Second,
				MaxRetries:     5,
				BatchSize:      20,
				EnablePodcasts: true,
			},
		},
		{
			name: "zero values use defaults",
			config: WorkerConfig{
				WorkerCount:    0,
				PollInterval:   0,
				MaxRetries:     0,
				BatchSize:      0,
				EnablePodcasts: false,
			},
			expectedConfig: WorkerConfig{
				WorkerCount:  2,
				PollInterval: 5 * time.Second,
				MaxRetries:   3,
				BatchSize:    10,
			},
		},
		{
			name: "negative values use defaults",
			config: WorkerConfig{
				WorkerCount:  -1,
				PollInterval: -1 * time.Second,
				MaxRetries:   -1,
				BatchSize:    -1,
			},
			expectedConfig: WorkerConfig{
				WorkerCount:  2,
				PollInterval: 5 * time.Second,
				MaxRetries:   3,
				BatchSize:    10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, tt.config)

			assert.NotNil(t, service)
			ws := service.(*workerService)
			assert.Equal(t, tt.expectedConfig.WorkerCount, ws.workerCount)
			assert.Equal(t, tt.expectedConfig.PollInterval, ws.pollInterval)
			assert.Equal(t, tt.expectedConfig.MaxRetries, ws.maxRetries)
			assert.Equal(t, tt.expectedConfig.BatchSize, ws.batchSize)
		})
	}
}

func TestWorkerService_StartStop(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 100 * time.Millisecond,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config)

	// Should not be running initially
	assert.False(t, service.IsRunning())

	// Mock dequeue to return empty items (workers will poll)
	ctx := context.Background()
	mockJobQueue.On("DequeuePendingItems", mock.Anything, config.BatchSize).Return([]db.Item{}, nil).Maybe()

	// Start the service
	err := service.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, service.IsRunning())

	// Starting again should fail
	err = service.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")

	// Stop the service
	err = service.Stop()
	assert.NoError(t, err)
	assert.False(t, service.IsRunning())

	// Stopping again should fail
	err = service.Stop()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
}

func TestWorkerService_ProcessItem_Success(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		MaxRetries:   3,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	url := "https://example.com/article"
	item := db.Item{
		ID:  1,
		Url: &url,
	}

	content := "Article content here"
	extraction := ItemExtraction{
		Title:    "Test Article",
		Type:     "article",
		Platform: "example.com",
		Tags:     []string{"tech", "news"},
		Authors:  []string{"John Doe"},
	}
	summary := ItemSummary{
		Overview:  "Brief overview about the test article",
		KeyPoints: []string{"Key point 1", "Key point 2"},
	}

	mockJobQueue.On("MarkItemAsProcessing", ctx, item.ID).Return(nil)
	mockScraping.On("Scrape", url).Return(content, nil)
	mockAI.On("ExtractContent", ctx, content).Return(extraction, nil)
	mockAI.On("SummarizeContent", ctx, content).Return(summary, nil)
	mockJobQueue.On("CompleteItem", ctx, item.ID, extraction.Title, content, mock.Anything, extraction.Type, extraction.Platform, extraction.Tags, extraction.Authors).Return(nil)

	err := service.processItem(ctx, item)

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
	mockScraping.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestWorkerService_ProcessItem_ScrapingFails(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		MaxRetries:   2,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	url := "https://example.com/article"
	item := db.Item{
		ID:  1,
		Url: &url,
	}

	scrapingError := errors.New("scraping failed")

	mockJobQueue.On("MarkItemAsProcessing", ctx, item.ID).Return(nil)
	mockScraping.On("Scrape", url).Return("", scrapingError).Times(2) // MaxRetries = 2
	mockJobQueue.On("FailItem", ctx, item.ID, mock.MatchedBy(func(msg string) bool {
		return true // Accept any error message
	})).Return(nil)

	err := service.processItem(ctx, item)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to process URL after")
	mockJobQueue.AssertExpectations(t)
	mockScraping.AssertExpectations(t)
}

func TestWorkerService_ProcessItem_ExtractionFails(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		MaxRetries:   2,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	url := "https://example.com/article"
	item := db.Item{
		ID:  1,
		Url: &url,
	}

	content := "Article content"
	extractionError := errors.New("extraction failed")

	mockJobQueue.On("MarkItemAsProcessing", ctx, item.ID).Return(nil)
	mockScraping.On("Scrape", url).Return(content, nil).Times(2)
	mockAI.On("ExtractContent", ctx, content).Return(ItemExtraction{}, extractionError).Times(2)
	mockJobQueue.On("FailItem", ctx, item.ID, mock.Anything).Return(nil)

	err := service.processItem(ctx, item)

	assert.Error(t, err)
	mockJobQueue.AssertExpectations(t)
	mockScraping.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestWorkerService_ProcessItem_SummarizationFails(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		MaxRetries:   2,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	url := "https://example.com/article"
	item := db.Item{
		ID:  1,
		Url: &url,
	}

	content := "Article content"
	extraction := ItemExtraction{
		Title:    "Test Article",
		Type:     "article",
		Platform: "example.com",
	}
	summarizationError := errors.New("summarization failed")

	mockJobQueue.On("MarkItemAsProcessing", ctx, item.ID).Return(nil)
	mockScraping.On("Scrape", url).Return(content, nil).Times(2)
	mockAI.On("ExtractContent", ctx, content).Return(extraction, nil).Times(2)
	mockAI.On("SummarizeContent", ctx, content).Return(ItemSummary{}, summarizationError).Times(2)
	mockJobQueue.On("FailItem", ctx, item.ID, mock.Anything).Return(nil)

	err := service.processItem(ctx, item)

	assert.Error(t, err)
	mockJobQueue.AssertExpectations(t)
	mockScraping.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestWorkerService_ProcessItem_RetrySuccess(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		MaxRetries:   3,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	url := "https://example.com/article"
	item := db.Item{
		ID:  1,
		Url: &url,
	}

	content := "Article content"
	extraction := ItemExtraction{
		Title:    "Test Article",
		Type:     "article",
		Platform: "example.com",
	}
	summary := ItemSummary{
		Overview:  "Brief overview",
		KeyPoints: []string{"Key point 1"},
	}

	mockJobQueue.On("MarkItemAsProcessing", ctx, item.ID).Return(nil)
	// Fail twice, succeed on third attempt
	mockScraping.On("Scrape", url).Return("", errors.New("scraping failed")).Times(2)
	mockScraping.On("Scrape", url).Return(content, nil).Once()
	mockAI.On("ExtractContent", ctx, content).Return(extraction, nil)
	mockAI.On("SummarizeContent", ctx, content).Return(summary, nil)
	mockJobQueue.On("CompleteItem", ctx, item.ID, extraction.Title, content, mock.Anything, extraction.Type, extraction.Platform, mock.Anything, mock.Anything).Return(nil)

	err := service.processItem(ctx, item)

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
	mockScraping.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestWorkerService_ProcessItemBatch_NoItems(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 10 * time.Millisecond,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx

	mockJobQueue.On("DequeuePendingItems", ctx, config.BatchSize).Return([]db.Item{}, nil)

	err := service.processItemBatch()

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
}

func TestWorkerService_ProcessItemBatch_DequeueError(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx
	dequeueError := errors.New("database error")

	mockJobQueue.On("DequeuePendingItems", ctx, config.BatchSize).Return(nil, dequeueError)

	err := service.processItemBatch()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to dequeue items")
	mockJobQueue.AssertExpectations(t)
}

func TestWorkerService_ProcessPodcast_Success(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   1 * time.Second,
		MaxRetries:     3,
		BatchSize:      5,
		EnablePodcasts: true,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	podcast := db.Podcast{
		ID:    1,
		Title: "Test Podcast",
	}

	mockPodcast.On("ProcessPodcast", ctx, podcast.ID).Return(nil)

	err := service.processPodcast(ctx, podcast)

	assert.NoError(t, err)
	mockPodcast.AssertExpectations(t)
}

func TestWorkerService_ProcessPodcast_Failure(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   1 * time.Second,
		MaxRetries:     2,
		BatchSize:      5,
		EnablePodcasts: true,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	podcast := db.Podcast{
		ID:    1,
		Title: "Test Podcast",
	}

	processingError := errors.New("processing failed")

	mockPodcast.On("ProcessPodcast", ctx, podcast.ID).Return(processingError).Times(2)
	mockPodcast.On("UpdatePodcastStatus", ctx, podcast.ID, PodcastStatusFailed).Return(nil)

	err := service.processPodcast(ctx, podcast)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to process podcast after")
	mockPodcast.AssertExpectations(t)
}

func TestWorkerService_ProcessPodcast_RetrySuccess(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   1 * time.Second,
		MaxRetries:     3,
		BatchSize:      5,
		EnablePodcasts: true,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	podcast := db.Podcast{
		ID:    1,
		Title: "Test Podcast",
	}

	// Fail twice, succeed on third attempt
	mockPodcast.On("ProcessPodcast", ctx, podcast.ID).Return(errors.New("temporary error")).Times(2)
	mockPodcast.On("ProcessPodcast", ctx, podcast.ID).Return(nil).Once()

	err := service.processPodcast(ctx, podcast)

	assert.NoError(t, err)
	mockPodcast.AssertExpectations(t)
}

func TestWorkerService_ProcessPodcastBatch_NoPodcasts(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   10 * time.Millisecond,
		BatchSize:      5,
		EnablePodcasts: true,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx
	mockPodcast.On("AcquirePendingPodcasts", ctx, config.BatchSize).Return([]db.Podcast{}, nil)

	err := service.processPodcastBatch()

	assert.NoError(t, err)
	mockPodcast.AssertExpectations(t)
}

func TestWorkerService_ProcessPodcastBatch_AcquireError(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   1 * time.Second,
		BatchSize:      5,
		EnablePodcasts: true,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx
	acquireError := errors.New("database error")

	mockPodcast.On("AcquirePendingPodcasts", ctx, config.BatchSize).Return(nil, acquireError)

	err := service.processPodcastBatch()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to acquire pending podcasts")
	mockPodcast.AssertExpectations(t)
}

func TestWorkerService_ProcessBatch_PodcastsDisabled(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   10 * time.Millisecond,
		BatchSize:      5,
		EnablePodcasts: false, // Podcasts disabled
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx
	mockJobQueue.On("DequeuePendingItems", ctx, config.BatchSize).Return([]db.Item{}, nil)

	err := service.processBatch()

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
	// Podcast service should not be called
	mockPodcast.AssertNotCalled(t, "AcquirePendingPodcasts")
}

func TestWorkerService_ProcessBatch_PodcastsEnabled(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   10 * time.Millisecond,
		BatchSize:      5,
		EnablePodcasts: true, // Podcasts enabled
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx
	mockJobQueue.On("DequeuePendingItems", ctx, config.BatchSize).Return([]db.Item{}, nil)
	mockPodcast.On("AcquirePendingPodcasts", ctx, config.BatchSize).Return([]db.Podcast{}, nil)

	err := service.processBatch()

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
	mockPodcast.AssertExpectations(t)
}

func TestWorkerService_ProcessBatch_PodcastServiceNil(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   10 * time.Millisecond,
		BatchSize:      5,
		EnablePodcasts: true,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, nil, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx
	mockJobQueue.On("DequeuePendingItems", ctx, config.BatchSize).Return([]db.Item{}, nil)

	err := service.processBatch()

	// Should not error when podcast service is nil even if podcasts enabled
	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
}

func TestWorkerService_ProcessURL_Success(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount: 1,
		BatchSize:   5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	url := "https://example.com/article"
	content := "Article content here"
	extraction := ItemExtraction{
		Title:    "Test Article",
		Type:     "article",
		Platform: "example.com",
		Tags:     []string{"tech"},
		Authors:  []string{"John Doe"},
	}
	summary := ItemSummary{
		Overview:  "Brief overview about the test article",
		KeyPoints: []string{"Key point 1"},
	}

	mockScraping.On("Scrape", url).Return(content, nil)
	mockAI.On("ExtractContent", ctx, content).Return(extraction, nil)
	mockAI.On("SummarizeContent", ctx, content).Return(summary, nil)

	resultContent, resultExtraction, resultSummary, err := service.processURL(ctx, url)

	assert.NoError(t, err)
	assert.Equal(t, content, resultContent)
	assert.Equal(t, extraction, resultExtraction)
	assert.NotEmpty(t, resultSummary)
	mockScraping.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestWorkerService_MarkItemAsProcessingError(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		MaxRetries:   1,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	url := "https://example.com/article"
	item := db.Item{
		ID:  1,
		Url: &url,
	}

	markError := errors.New("failed to mark as processing")

	mockJobQueue.On("MarkItemAsProcessing", ctx, item.ID).Return(markError)

	err := service.processItem(ctx, item)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to mark item as processing")
	mockJobQueue.AssertExpectations(t)
}

func TestWorkerService_CompleteItemError(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		MaxRetries:   1,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	ctx := context.Background()
	url := "https://example.com/article"
	item := db.Item{
		ID:  1,
		Url: &url,
	}

	content := "Article content"
	extraction := ItemExtraction{
		Title:    "Test Article",
		Type:     "article",
		Platform: "example.com",
	}
	summary := ItemSummary{
		Overview:  "Brief overview",
		KeyPoints: []string{"Key point 1"},
	}

	mockJobQueue.On("MarkItemAsProcessing", ctx, item.ID).Return(nil)
	mockScraping.On("Scrape", url).Return(content, nil)
	mockAI.On("ExtractContent", ctx, content).Return(extraction, nil)
	mockAI.On("SummarizeContent", ctx, content).Return(summary, nil)
	mockJobQueue.On("CompleteItem", ctx, item.ID, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("failed to complete item"))

	err := service.processItem(ctx, item)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to complete item")
	mockJobQueue.AssertExpectations(t)
	mockScraping.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestWorkerService_ContextCancellation(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 100 * time.Millisecond,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config)

	ctx, cancel := context.WithCancel(context.Background())

	// Mock dequeue to return empty items (workers will poll)
	mockJobQueue.On("DequeuePendingItems", mock.Anything, config.BatchSize).Return([]db.Item{}, nil).Maybe()

	// Start the service
	err := service.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, service.IsRunning())

	// Cancel the context
	cancel()

	// Give it time to stop
	time.Sleep(200 * time.Millisecond)

	// The service should still report as running since Stop() wasn't called
	// Context cancellation stops workers but doesn't update the running flag
	assert.True(t, service.IsRunning())

	// Stop should succeed
	err = service.Stop()
	assert.NoError(t, err)
	assert.False(t, service.IsRunning())
}

func TestWorkerService_ProcessItemBatch_MultipleItems(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:  1,
		PollInterval: 1 * time.Second,
		MaxRetries:   1,
		BatchSize:    5,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx
	url1 := "https://example.com/article1"
	url2 := "https://example.com/article2"

	items := []db.Item{
		{ID: 1, Url: &url1},
		{ID: 2, Url: &url2},
	}

	content := "Article content"
	extraction := ItemExtraction{
		Title:    "Test Article",
		Type:     "article",
		Platform: "example.com",
	}
	summary := ItemSummary{
		Overview:  "Brief overview",
		KeyPoints: []string{"Key point 1"},
	}

	mockJobQueue.On("DequeuePendingItems", ctx, config.BatchSize).Return(items, nil)

	// Setup expectations for both items
	for _, item := range items {
		mockJobQueue.On("MarkItemAsProcessing", ctx, item.ID).Return(nil)
		mockScraping.On("Scrape", *item.Url).Return(content, nil)
		mockAI.On("ExtractContent", ctx, content).Return(extraction, nil)
		mockAI.On("SummarizeContent", ctx, content).Return(summary, nil)
		mockJobQueue.On("CompleteItem", ctx, item.ID, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	}

	err := service.processItemBatch()

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
	mockScraping.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestWorkerService_ProcessPodcastBatch_MultiplePodcasts(t *testing.T) {
	mockJobQueue := new(MockJobQueueService)
	mockAI := new(MockAIService)
	mockScraping := new(MockScrapingService)
	mockPodcast := new(MockPodcastService)

	config := WorkerConfig{
		WorkerCount:    1,
		PollInterval:   1 * time.Second,
		MaxRetries:     1,
		BatchSize:      5,
		EnablePodcasts: true,
	}

	service := NewWorkerService(mockJobQueue, mockAI, mockScraping, mockPodcast, config).(*workerService)

	// Initialize context for the service
	ctx := context.Background()
	service.ctx = ctx
	podcasts := []db.Podcast{
		{ID: 1, Title: "Podcast 1"},
		{ID: 2, Title: "Podcast 2"},
	}

	mockPodcast.On("AcquirePendingPodcasts", ctx, config.BatchSize).Return(podcasts, nil)

	// Setup expectations for both podcasts
	for _, podcast := range podcasts {
		mockPodcast.On("ProcessPodcast", ctx, podcast.ID).Return(nil)
	}

	err := service.processPodcastBatch()

	assert.NoError(t, err)
	mockPodcast.AssertExpectations(t)
}
