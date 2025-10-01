package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
	"github.com/yamirghofran/briefbot/internal/test"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(ctx context.Context, request EmailRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

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

func (m *MockPodcastService) GetPodcast(ctx context.Context, id int32) (*db.Podcast, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GetPodcastsByUser(ctx context.Context, userID int32) ([]db.Podcast, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GetPodcastsByStatus(ctx context.Context, status PodcastStatus) ([]db.Podcast, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) UpdatePodcast(ctx context.Context, podcastID int32, title string, description string) error {
	args := m.Called(ctx, podcastID, title, description)
	return args.Error(0)
}

func (m *MockPodcastService) DeletePodcast(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
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
	return args.Get(0).([]db.GetPodcastItemsRow), args.Error(1)
}

func (m *MockPodcastService) GetPendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) AcquirePendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GetProcessingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GeneratePodcastUploadURL(ctx context.Context, podcastID int32) (*UploadURLResponse, error) {
	args := m.Called(ctx, podcastID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UploadURLResponse), args.Error(1)
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

func (m *MockPodcastService) UpdatePodcastStatus(ctx context.Context, podcastID int32, status PodcastStatus) error {
	args := m.Called(ctx, podcastID, status)
	return args.Error(0)
}


func TestSetPodcastGenerationEnabled(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)

	// Initially should be false (unless env var is set)
	service.SetPodcastGenerationEnabled(true)
	assert.True(t, service.IsPodcastGenerationEnabled())

	service.SetPodcastGenerationEnabled(false)
	assert.False(t, service.IsPodcastGenerationEnabled())
}

func TestGetDailyDigestItemsForUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)

	ctx := context.Background()
	userID := int32(1)

	expectedItems := []db.Item{
		{ID: 1, Title: "Item 1"},
		{ID: 2, Title: "Item 2"},
	}

	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &userID).Return(expectedItems, nil)

	items, err := service.GetDailyDigestItemsForUser(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	mockQuerier.AssertExpectations(t)
}
