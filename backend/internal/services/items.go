package services

import (
	"context"

	"github.com/yamirghofran/briefbot/internal/db"
)

type ItemService interface {
	CreateItem(ctx context.Context, userID *int32, url, fileKey, textContent, summary *string) (*db.Item, error)
	GetItem(ctx context.Context, id int32) (*db.Item, error)
	GetItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error)
	GetUnreadItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error)
	UpdateItem(ctx context.Context, id int32, url, fileKey, textContent, summary *string, isRead *bool) error
	MarkItemAsRead(ctx context.Context, id int32) error
	DeleteItem(ctx context.Context, id int32) error
}

type itemService struct {
	querier db.Querier
}

func NewItemService(querier db.Querier) ItemService {
	return &itemService{querier: querier}
}

func (s *itemService) CreateItem(ctx context.Context, userID *int32, url, fileKey, textContent, summary *string) (*db.Item, error) {
	params := db.CreateItemParams{
		UserID:      userID,
		Url:         url,
		FileKey:     fileKey,
		TextContent: textContent,
		Summary:     summary,
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

func (s *itemService) UpdateItem(ctx context.Context, id int32, url, fileKey, textContent, summary *string, isRead *bool) error {
	params := db.UpdateItemParams{
		ID:          id,
		Url:         url,
		IsRead:      isRead,
		FileKey:     fileKey,
		TextContent: textContent,
		Summary:     summary,
	}
	return s.querier.UpdateItem(ctx, params)
}

func (s *itemService) MarkItemAsRead(ctx context.Context, id int32) error {
	return s.querier.MarkItemAsRead(ctx, id)
}

func (s *itemService) DeleteItem(ctx context.Context, id int32) error {
	return s.querier.DeleteItem(ctx, id)
}
