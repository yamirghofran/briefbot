package services

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
)

type MockPodcastService struct {
	mock.Mock
}

func (m *MockPodcastService) CreatePodcastFromItems(ctx context.Context, userID int32, title string, description string, itemIDs []int32) (*db.Podcast, error) {
	args := m.Called(ctx, userID, title, description, itemIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Podcast), args.Error(1)
}

func (m *MockPodcastService) CreatePodcastFromSingleItem(ctx context.Context, userID int32, itemID int32) (*db.Podcast, error) {
	args := m.Called(ctx, userID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GeneratePodcastScript(ctx context.Context, podcastID int32) error {
	args := m.Called(ctx, podcastID)
	return args.Error(0)
}

func (m *MockPodcastService) GeneratePodcastAudio(ctx context.Context, podcastID int32) error {
	args := m.Called(ctx, podcastID)
	return args.Error(0)
}

func (m *MockPodcastService) ProcessPodcast(ctx context.Context, podcastID int32) error {
	args := m.Called(ctx, podcastID)
	return args.Error(0)
}

func (m *MockPodcastService) GetPodcast(ctx context.Context, podcastID int32) (*db.Podcast, error) {
	args := m.Called(ctx, podcastID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GetPodcastsByUser(ctx context.Context, userID int32) ([]db.Podcast, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GetPodcastsByStatus(ctx context.Context, status PodcastStatus) ([]db.Podcast, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) UpdatePodcast(ctx context.Context, podcastID int32, title string, description string) error {
	args := m.Called(ctx, podcastID, title, description)
	return args.Error(0)
}

func (m *MockPodcastService) DeletePodcast(ctx context.Context, podcastID int32) error {
	args := m.Called(ctx, podcastID)
	return args.Error(0)
}

func (m *MockPodcastService) AddItemToPodcast(ctx context.Context, podcastID int32, itemID int32, order int) error {
	args := m.Called(ctx, podcastID, itemID, order)
	return args.Error(0)
}

func (m *MockPodcastService) RemoveItemFromPodcast(ctx context.Context, podcastID int32, itemID int32) error {
	args := m.Called(ctx, podcastID, itemID)
	return args.Error(0)
}

func (m *MockPodcastService) GetPodcastItems(ctx context.Context, podcastID int32) ([]db.GetPodcastItemsRow, error) {
	args := m.Called(ctx, podcastID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.GetPodcastItemsRow), args.Error(1)
}

func (m *MockPodcastService) UpdatePodcastStatus(ctx context.Context, podcastID int32, status PodcastStatus) error {
	args := m.Called(ctx, podcastID, status)
	return args.Error(0)
}

func (m *MockPodcastService) GetPendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GetProcessingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) AcquirePendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GetPodcastAudio(ctx context.Context, podcastID int32) ([]byte, error) {
	args := m.Called(ctx, podcastID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPodcastService) HasPodcastAudio(ctx context.Context, podcastID int32) (bool, error) {
	args := m.Called(ctx, podcastID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPodcastService) GeneratePodcastUploadURL(ctx context.Context, podcastID int32) (*UploadURLResponse, error) {
	args := m.Called(ctx, podcastID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UploadURLResponse), args.Error(1)
}
