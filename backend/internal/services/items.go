package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/yamirghofran/briefbot/internal/db"
)

type ItemService interface {
	// Background processing methods
	CreateItemAsync(ctx context.Context, userID int32, url string) (*db.Item, error)
	ProcessURL(ctx context.Context, userID int32, url string) (*db.Item, error)

	// Traditional CRUD methods
	CreateItem(ctx context.Context, userID *int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string) (*db.Item, error)
	GetItem(ctx context.Context, id int32) (*db.Item, error)
	GetItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error)
	GetUnreadItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error)
	GetUnreadItemsFromPreviousDay(ctx context.Context) ([]db.Item, error)
	UpdateItem(ctx context.Context, id int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string, isRead *bool) error
	MarkItemAsRead(ctx context.Context, id int32) error
	ToggleItemReadStatus(ctx context.Context, id int32) (*db.Item, error)
	DeleteItem(ctx context.Context, id int32) error
	GetItemProcessingStatus(ctx context.Context, itemID int32) (*ItemStatus, error)
	GetItemsByProcessingStatus(ctx context.Context, status *string) ([]db.Item, error)
}

type itemService struct {
	querier         db.Querier
	aiService       AIService
	scrapingService ScrapingService
	jobQueueService JobQueueService
}

func NewItemService(querier db.Querier, aiService AIService, scrapingService ScrapingService, jobQueueService JobQueueService) ItemService {
	return &itemService{
		querier:         querier,
		aiService:       aiService,
		scrapingService: scrapingService,
		jobQueueService: jobQueueService,
	}
}

// CreateItemAsync creates an item asynchronously - just saves the URL and returns immediately
func (s *itemService) CreateItemAsync(ctx context.Context, userID int32, url string) (*db.Item, error) {
	// For async creation, we just save the URL with a placeholder title
	// The actual processing will happen in the background
	item, err := s.jobQueueService.EnqueueItem(ctx, userID, url, url)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue item: %w", err)
	}

	return item, nil
}

// ProcessURL processes a URL synchronously (for backward compatibility or manual processing)
func (s *itemService) ProcessURL(ctx context.Context, userID int32, url string) (*db.Item, error) {
	content, err := s.scrapingService.Scrape(url)
	if err != nil {
		return nil, err
	}
	extraction, err := s.aiService.ExtractContent(ctx, content)
	if err != nil {
		return nil, err
	}

	summary, err := s.aiService.SummarizeContent(ctx, content)
	if err != nil {
		return nil, err
	}

	concatenatedSummary := ConcatenateSummary(summary)

	return s.CreateItem(ctx, &userID, extraction.Title, &url, &content, &concatenatedSummary, &extraction.Type, &extraction.Platform,
		extraction.Tags, extraction.Authors)
}

func (s *itemService) CreateItem(ctx context.Context, userID *int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string) (*db.Item, error) {
	completedStatus := ProcessingStatusCompleted
	params := db.CreateItemParams{
		UserID:           userID,
		Title:            title,
		Url:              url,
		TextContent:      textContent,
		Summary:          summary,
		Type:             itemType,
		Tags:             tags,
		Platform:         platform,
		Authors:          authors,
		ProcessingStatus: &completedStatus, // Mark as completed since we're creating with all data
		ProcessingError:  nil,
	}
	item, err := s.querier.CreateItem(ctx, params)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *itemService) GetItem(ctx context.Context, id int32) (*db.Item, error) {
	item, err := s.querier.GetItem(ctx, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *itemService) GetItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	items, err := s.querier.GetItemsByUser(ctx, userID)
	if err != nil {
		return []db.Item{}, err
	}
	return items, nil
}

func (s *itemService) GetUnreadItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	items, err := s.querier.GetUnreadItemsByUser(ctx, userID)
	if err != nil {
		return []db.Item{}, err
	}
	return items, nil
}

func (s *itemService) UpdateItem(ctx context.Context, id int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string, isRead *bool) error {
	params := db.UpdateItemParams{
		ID:          id,
		Title:       title,
		Url:         url,
		IsRead:      isRead,
		TextContent: textContent,
		Summary:     summary,
		Type:        itemType,
		Tags:        tags,
		Platform:    platform,
		Authors:     authors,
	}
	return s.querier.UpdateItem(ctx, params)
}

func (s *itemService) MarkItemAsRead(ctx context.Context, id int32) error {
	return s.querier.MarkItemAsRead(ctx, id)
}

func (s *itemService) ToggleItemReadStatus(ctx context.Context, id int32) (*db.Item, error) {
	item, err := s.querier.ToggleItemReadStatus(ctx, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *itemService) DeleteItem(ctx context.Context, id int32) error {
	return s.querier.DeleteItem(ctx, id)
}

func (s *itemService) GetItemProcessingStatus(ctx context.Context, itemID int32) (*ItemStatus, error) {
	if s.jobQueueService == nil {
		return nil, fmt.Errorf("job queue service not available")
	}
	return s.jobQueueService.GetItemStatus(ctx, itemID)
}

func (s *itemService) GetItemsByProcessingStatus(ctx context.Context, status *string) ([]db.Item, error) {
	if s.jobQueueService == nil {
		return nil, fmt.Errorf("job queue service not available")
	}
	return s.jobQueueService.GetItemsByStatus(ctx, *status)
}

func (s *itemService) GetUnreadItemsFromPreviousDay(ctx context.Context) ([]db.Item, error) {
	items, err := s.querier.GetUnreadItemsFromPreviousDay(ctx)
	if err != nil {
		return []db.Item{}, err
	}
	return items, nil
}

func ConcatenateSummary(summary ItemSummary) string {
	return summary.Overview + " " + strings.Join(summary.KeyPoints, " ")
}
