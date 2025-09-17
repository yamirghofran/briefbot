package services

import (
	"context"
	"fmt"

	"github.com/yamirghofran/briefbot/internal/db"
)

// ProcessingStatus constants
const (
	ProcessingStatusPending    = "pending"
	ProcessingStatusProcessing = "processing"
	ProcessingStatusCompleted  = "completed"
	ProcessingStatusFailed     = "failed"
)

type JobQueueService interface {
	// Queue management
	EnqueueItem(ctx context.Context, userID int32, title string, url string) (*db.Item, error)
	DequeuePendingItems(ctx context.Context, limit int32) ([]db.Item, error)
	MarkItemAsProcessing(ctx context.Context, itemID int32) error

	// Status management
	CompleteItem(ctx context.Context, itemID int32, title, textContent, summary, itemType, platform string, tags, authors []string) error
	FailItem(ctx context.Context, itemID int32, errorMsg string) error
	GetItemStatus(ctx context.Context, itemID int32) (*ItemStatus, error)

	// Utility methods
	GetItemsByStatus(ctx context.Context, status string) ([]db.Item, error)
	GetFailedItemsForRetry(ctx context.Context, limit int32) ([]db.Item, error)
	RetryItem(ctx context.Context, itemID int32) error
}

type ItemStatus struct {
	Item            *db.Item
	IsProcessing    bool
	IsCompleted     bool
	IsFailed        bool
	ProcessingError *string
}

type jobQueueService struct {
	querier db.Querier
}

func NewJobQueueService(querier db.Querier) JobQueueService {
	return &jobQueueService{querier: querier}
}

func (s *jobQueueService) EnqueueItem(ctx context.Context, userID int32, title string, url string) (*db.Item, error) {
	params := db.CreatePendingItemParams{
		UserID: &userID,
		Title:  title,
		Url:    &url,
	}

	item, err := s.querier.CreatePendingItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue item: %w", err)
	}

	return &item, nil
}

func (s *jobQueueService) DequeuePendingItems(ctx context.Context, limit int32) ([]db.Item, error) {
	items, err := s.querier.GetPendingItems(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue pending items: %w", err)
	}

	return items, nil
}

func (s *jobQueueService) MarkItemAsProcessing(ctx context.Context, itemID int32) error {
	err := s.querier.UpdateItemAsProcessing(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to mark item as processing: %w", err)
	}

	return nil
}

func (s *jobQueueService) CompleteItem(ctx context.Context, itemID int32, title, textContent, summary, itemType, platform string, tags, authors []string) error {
	// Get the current item to preserve existing URL
	item, err := s.querier.GetItem(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to get item for completion: %w", err)
	}

	// Update the item with processed data, using the AI-extracted title and preserving URL
	params := db.UpdateItemParams{
		ID:          itemID,
		Title:       title,    // Use AI-extracted title
		Url:         item.Url, // Preserve existing URL
		IsRead:      nil,
		TextContent: &textContent,
		Summary:     &summary,
		Type:        &itemType,
		Tags:        tags,
		Platform:    &platform,
		Authors:     authors,
	}

	err = s.querier.UpdateItem(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update item with processed data: %w", err)
	}

	// Then mark as completed
	completedStatus := ProcessingStatusCompleted
	statusParams := db.UpdateItemProcessingStatusParams{
		ID:               itemID,
		ProcessingStatus: &completedStatus,
		ProcessingError:  nil,
	}

	err = s.querier.UpdateItemProcessingStatus(ctx, statusParams)
	if err != nil {
		return fmt.Errorf("failed to mark item as completed: %w", err)
	}

	return nil
}

func (s *jobQueueService) FailItem(ctx context.Context, itemID int32, errorMsg string) error {
	failedStatus := ProcessingStatusFailed
	statusParams := db.UpdateItemProcessingStatusParams{
		ID:               itemID,
		ProcessingStatus: &failedStatus,
		ProcessingError:  &errorMsg,
	}

	err := s.querier.UpdateItemProcessingStatus(ctx, statusParams)
	if err != nil {
		return fmt.Errorf("failed to mark item as failed: %w", err)
	}

	return nil
}

func (s *jobQueueService) GetItemStatus(ctx context.Context, itemID int32) (*ItemStatus, error) {
	item, err := s.querier.GetItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	status := &ItemStatus{
		Item:            &item,
		IsProcessing:    item.ProcessingStatus != nil && *item.ProcessingStatus == ProcessingStatusProcessing,
		IsCompleted:     item.ProcessingStatus != nil && *item.ProcessingStatus == ProcessingStatusCompleted,
		IsFailed:        item.ProcessingStatus != nil && *item.ProcessingStatus == ProcessingStatusFailed,
		ProcessingError: item.ProcessingError,
	}

	return status, nil
}

func (s *jobQueueService) GetItemsByStatus(ctx context.Context, status string) ([]db.Item, error) {
	statusPtr := &status
	items, err := s.querier.GetItemsByProcessingStatus(ctx, statusPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get items by status: %w", err)
	}

	return items, nil
}

func (s *jobQueueService) GetFailedItemsForRetry(ctx context.Context, limit int32) ([]db.Item, error) {
	// This would require a new SQL query, for now return empty slice
	// In a real implementation, you'd add a SQL query to get failed items from last 24h
	return []db.Item{}, nil
}

func (s *jobQueueService) RetryItem(ctx context.Context, itemID int32) error {
	// Reset the item status to pending for retry
	pendingStatus := ProcessingStatusPending
	statusParams := db.UpdateItemProcessingStatusParams{
		ID:               itemID,
		ProcessingStatus: &pendingStatus,
		ProcessingError:  nil,
	}

	err := s.querier.UpdateItemProcessingStatus(ctx, statusParams)
	if err != nil {
		return fmt.Errorf("failed to retry item: %w", err)
	}

	return nil
}
