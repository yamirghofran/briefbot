package services

import (
	"context"
	"strings"

	"github.com/yamirghofran/briefbot/internal/db"
)

type ItemService interface {
	ProcessURL(ctx context.Context, userID int32, url string) (*db.Item, error)
	CreateItem(ctx context.Context, userID *int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string) (*db.Item, error)
	GetItem(ctx context.Context, id int32) (*db.Item, error)
	GetItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error)
	GetUnreadItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error)
	UpdateItem(ctx context.Context, id int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string, isRead *bool) error
	MarkItemAsRead(ctx context.Context, id int32) error
	DeleteItem(ctx context.Context, id int32) error
}

type itemService struct {
	querier         db.Querier
	aiService       AIService
	scrapingService ScrapingService
}

func NewItemService(querier db.Querier, aiService AIService, scrapingService ScrapingService) ItemService {
	return &itemService{querier: querier, aiService: aiService, scrapingService: scrapingService}
}

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
	params := db.CreateItemParams{
		UserID:      userID,
		Title:       title,
		Url:         url,
		TextContent: textContent,
		Summary:     summary,
		Type:        itemType,
		Tags:        tags,
		Platform:    platform,
		Authors:     authors,
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

func (s *itemService) DeleteItem(ctx context.Context, id int32) error {
	return s.querier.DeleteItem(ctx, id)
}

func ConcatenateSummary(summary ItemSummary) string {
	return summary.Overview + " " + strings.Join(summary.KeyPoints, " ")
}
