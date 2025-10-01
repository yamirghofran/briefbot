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

	// SSE integration
	SetSSEManager(sseManager *SSEManager)
}

type ItemStatus struct {
	Item            *db.Item
	IsProcessing    bool
	IsCompleted     bool
	IsFailed        bool
	ProcessingError *string
}

type jobQueueService struct {
	querier    db.Querier
	sseManager *SSEManager
}

func NewJobQueueService(querier db.Querier) JobQueueService {
	return &jobQueueService{
		querier:    querier,
		sseManager: nil, // Will be set later via SetSSEManager
	}
}

// SetSSEManager sets the SSE manager for the job queue service
func (s *jobQueueService) SetSSEManager(sseManager *SSEManager) {
	s.sseManager = sseManager
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

	// Notify SSE clients about new item
	if s.sseManager != nil && item.ProcessingStatus != nil {
		s.sseManager.NotifyItemUpdate(userID, item.ID, item.ProcessingStatus, "created")
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
	// Get item to find user ID
	item, err := s.querier.GetItem(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	err = s.querier.UpdateItemAsProcessing(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to mark item as processing: %w", err)
	}

	// Notify SSE clients about processing status
	if s.sseManager != nil && item.UserID != nil {
		processingStatus := ProcessingStatusProcessing
		s.sseManager.NotifyItemUpdate(*item.UserID, itemID, &processingStatus, "processing")
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
		Title:       title,       // Use AI-extracted title
		Url:         item.Url,    // Preserve existing URL
		IsRead:      item.IsRead, // Preserve existing is_read value
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

	// Notify SSE clients about completion
	if s.sseManager != nil && item.UserID != nil {
		s.sseManager.NotifyItemUpdate(*item.UserID, itemID, &completedStatus, "completed")
	}

	return nil
}

func (s *jobQueueService) FailItem(ctx context.Context, itemID int32, errorMsg string) error {
	// Get item to find user ID
	item, err := s.querier.GetItem(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	failedStatus := ProcessingStatusFailed
	statusParams := db.UpdateItemProcessingStatusParams{
		ID:               itemID,
		ProcessingStatus: &failedStatus,
		ProcessingError:  &errorMsg,
	}

	err = s.querier.UpdateItemProcessingStatus(ctx, statusParams)
	if err != nil {
		return fmt.Errorf("failed to mark item as failed: %w", err)
	}

	// Notify SSE clients about failure
	if s.sseManager != nil && item.UserID != nil {
		s.sseManager.NotifyItemUpdate(*item.UserID, itemID, &failedStatus, "failed")
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
