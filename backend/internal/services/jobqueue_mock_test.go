package services

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
)

type MockJobQueueService struct {
	mock.Mock
}

func (m *MockJobQueueService) EnqueueItem(ctx context.Context, userID int32, title string, url string) (*db.Item, error) {
	args := m.Called(ctx, userID, title, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockJobQueueService) DequeuePendingItems(ctx context.Context, limit int32) ([]db.Item, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockJobQueueService) MarkItemAsProcessing(ctx context.Context, itemID int32) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockJobQueueService) CompleteItem(ctx context.Context, itemID int32, title, textContent, summary, itemType, platform string, tags, authors []string) error {
	args := m.Called(ctx, itemID, title, textContent, summary, itemType, platform, tags, authors)
	return args.Error(0)
}

func (m *MockJobQueueService) FailItem(ctx context.Context, itemID int32, errorMsg string) error {
	args := m.Called(ctx, itemID, errorMsg)
	return args.Error(0)
}

func (m *MockJobQueueService) GetItemStatus(ctx context.Context, itemID int32) (*ItemStatus, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ItemStatus), args.Error(1)
}

func (m *MockJobQueueService) GetItemsByStatus(ctx context.Context, status string) ([]db.Item, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockJobQueueService) GetFailedItemsForRetry(ctx context.Context, limit int32) ([]db.Item, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockJobQueueService) RetryItem(ctx context.Context, itemID int32) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}
