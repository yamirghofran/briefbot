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

func (m *MockQuerier) ToggleItemReadStatus(ctx context.Context, id int32) (db.Item, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Item), args.Error(1)
}

func (m *MockQuerier) UpdateItem(ctx context.Context, arg db.UpdateItemParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) PatchItem(ctx context.Context, arg db.PatchItemParams) (db.Item, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Item), args.Error(1)
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

// Podcast-related methods
func (m *MockQuerier) AddItemToPodcast(ctx context.Context, arg db.AddItemToPodcastParams) (db.PodcastItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.PodcastItem), args.Error(1)
}

func (m *MockQuerier) ClearPodcastItems(ctx context.Context, podcastID *int32) error {
	args := m.Called(ctx, podcastID)
	return args.Error(0)
}

func (m *MockQuerier) CountPodcastItems(ctx context.Context, podcastID *int32) (int64, error) {
	args := m.Called(ctx, podcastID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) CreatePodcast(ctx context.Context, arg db.CreatePodcastParams) (db.Podcast, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Podcast), args.Error(1)
}

func (m *MockQuerier) CreatePodcastWithDialogues(ctx context.Context, arg db.CreatePodcastWithDialoguesParams) (db.Podcast, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Podcast), args.Error(1)
}

func (m *MockQuerier) DeletePodcast(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) GetCompletedPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetPendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetPodcast(ctx context.Context, id int32) (db.Podcast, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetPodcastByUser(ctx context.Context, userID *int32) ([]db.Podcast, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetPodcastItemIDs(ctx context.Context, podcastID *int32) ([]*int32, error) {
	args := m.Called(ctx, podcastID)
	return args.Get(0).([]*int32), args.Error(1)
}

func (m *MockQuerier) GetPodcastItems(ctx context.Context, podcastID *int32) ([]db.GetPodcastItemsRow, error) {
	args := m.Called(ctx, podcastID)
	return args.Get(0).([]db.GetPodcastItemsRow), args.Error(1)
}

func (m *MockQuerier) GetPodcastWithItems(ctx context.Context, id int32) (db.GetPodcastWithItemsRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.GetPodcastWithItemsRow), args.Error(1)
}

func (m *MockQuerier) GetPodcastsByStatus(ctx context.Context, status string) ([]db.Podcast, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetPodcastsByUserAndStatus(ctx context.Context, arg db.GetPodcastsByUserAndStatusParams) ([]db.Podcast, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetPodcastsForItem(ctx context.Context, itemID *int32) ([]db.Podcast, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetProcessingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetRecentPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockQuerier) GetUserPodcastStats(ctx context.Context, userID *int32) (db.GetUserPodcastStatsRow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(db.GetUserPodcastStatsRow), args.Error(1)
}

func (m *MockQuerier) RemoveItemFromPodcast(ctx context.Context, arg db.RemoveItemFromPodcastParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdatePodcast(ctx context.Context, arg db.UpdatePodcastParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdatePodcastAudio(ctx context.Context, arg db.UpdatePodcastAudioParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdatePodcastDialogues(ctx context.Context, arg db.UpdatePodcastDialoguesParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdatePodcastItemOrder(ctx context.Context, arg db.UpdatePodcastItemOrderParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdatePodcastStatus(ctx context.Context, arg db.UpdatePodcastStatusParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdatePodcastStatusWithAudio(ctx context.Context, arg db.UpdatePodcastStatusWithAudioParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdatePodcastsStatus(ctx context.Context, arg db.UpdatePodcastsStatusParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}
