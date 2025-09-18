package test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
)

// MockQuerier is a mock implementation of db.Querier
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateItem(ctx context.Context, arg db.CreateItemParams) (db.Item, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Item), args.Error(1)
}

func (m *MockQuerier) CreatePendingItem(ctx context.Context, arg db.CreatePendingItemParams) (db.Item, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Item), args.Error(1)
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockQuerier) DeleteItem(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) DeleteUser(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) GetFailedItemsForRetry(ctx context.Context, limit int32) ([]db.Item, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockQuerier) GetItem(ctx context.Context, id int32) (db.Item, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Item), args.Error(1)
}

func (m *MockQuerier) GetItemsByProcessingStatus(ctx context.Context, processingStatus *string) ([]db.Item, error) {
	args := m.Called(ctx, processingStatus)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockQuerier) GetItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockQuerier) GetPendingItems(ctx context.Context, limit int32) ([]db.Item, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockQuerier) GetUnreadItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockQuerier) GetUser(ctx context.Context, id int32) (db.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockQuerier) GetUserByEmail(ctx context.Context, email *string) (db.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockQuerier) ListUsers(ctx context.Context) ([]db.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.User), args.Error(1)
}

func (m *MockQuerier) MarkItemAsRead(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) UpdateItem(ctx context.Context, arg db.UpdateItemParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdateItemAsProcessing(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) UpdateItemProcessingStatus(ctx context.Context, arg db.UpdateItemProcessingStatusParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdateUser(ctx context.Context, arg db.UpdateUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) GetUnreadItemsFromPreviousDay(ctx context.Context) ([]db.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockQuerier) GetUnreadItemsFromPreviousDayByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Item), args.Error(1)
}
