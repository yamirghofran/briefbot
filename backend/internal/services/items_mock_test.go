package services

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
)

type MockItemService struct {
	mock.Mock
}

func (m *MockItemService) CreateItemAsync(ctx context.Context, userID int32, url string) (*db.Item, error) {
	args := m.Called(ctx, userID, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) ProcessURL(ctx context.Context, userID int32, url string) (*db.Item, error) {
	args := m.Called(ctx, userID, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) CreateItem(ctx context.Context, userID *int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string) (*db.Item, error) {
	args := m.Called(ctx, userID, title, url, textContent, summary, itemType, platform, tags, authors)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) GetItem(ctx context.Context, id int32) (*db.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) GetItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockItemService) GetUnreadItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockItemService) GetUnreadItemsFromPreviousDay(ctx context.Context) ([]db.Item, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockItemService) UpdateItem(ctx context.Context, id int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string, isRead *bool) error {
	args := m.Called(ctx, id, title, url, textContent, summary, itemType, platform, tags, authors, isRead)
	return args.Error(0)
}

func (m *MockItemService) MarkItemAsRead(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemService) ToggleItemReadStatus(ctx context.Context, id int32) (*db.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) DeleteItem(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemService) GetItemProcessingStatus(ctx context.Context, itemID int32) (*ItemStatus, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ItemStatus), args.Error(1)
}

func (m *MockItemService) GetItemsByProcessingStatus(ctx context.Context, status *string) ([]db.Item, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Item), args.Error(1)
}
