package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
	"github.com/yamirghofran/briefbot/internal/test"
)

func timeNow() *time.Time {
	t := time.Now()
	return &t
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

func TestSendDailyDigestForUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)

	ctx := context.Background()
	userID := int32(1)
	email := "test@example.com"

	items := []db.Item{
		{ID: 1, Title: "Test Item", Url: stringPtr("https://example.com"), CreatedAt: timeNow()},
	}

	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &userID).Return(items, nil)
	mockQuerier.On("GetUser", ctx, userID).Return(db.User{ID: userID, Email: &email}, nil)
	mockEmail.On("SendEmail", ctx, mock.MatchedBy(func(req EmailRequest) bool {
		return len(req.ToAddresses) == 1 && req.ToAddresses[0] == email
	})).Return(nil)

	err := service.SendDailyDigestForUser(ctx, userID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestSendDailyDigestForUser_NoItems(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)

	ctx := context.Background()
	userID := int32(1)

	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &userID).Return([]db.Item{}, nil)

	err := service.SendDailyDigestForUser(ctx, userID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
	// Email should NOT be called when no items
	mockEmail.AssertNotCalled(t, "SendEmail")
}

func TestSendDailyDigestForUser_NoEmail(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)

	ctx := context.Background()
	userID := int32(1)

	items := []db.Item{
		{ID: 1, Title: "Test Item"},
	}

	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &userID).Return(items, nil)
	mockQuerier.On("GetUser", ctx, userID).Return(db.User{ID: userID, Email: nil}, nil)

	err := service.SendDailyDigestForUser(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has no email address")
	mockQuerier.AssertExpectations(t)
}

func TestSendDailyDigest(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)

	ctx := context.Background()
	email1 := "user1@example.com"
	email2 := "user2@example.com"

	users := []db.User{
		{ID: 1, Email: &email1},
		{ID: 2, Email: &email2},
	}

	mockQuerier.On("ListUsers", ctx).Return(users, nil)

	// For each user
	for _, user := range users {
		mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &user.ID).Return([]db.Item{
			{ID: 1, Title: "Test", Url: stringPtr("https://example.com"), CreatedAt: timeNow()},
		}, nil)
		mockQuerier.On("GetUser", ctx, user.ID).Return(user, nil)
		mockEmail.On("SendEmail", ctx, mock.Anything).Return(nil)
	}

	err := service.SendDailyDigest(ctx)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestSendIntegratedDigestForUser_WithoutPodcast(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)
	service.SetPodcastGenerationEnabled(false)

	ctx := context.Background()
	userID := int32(1)
	email := "test@example.com"

	items := []db.Item{
		{ID: 1, Title: "Test Item", Url: stringPtr("https://example.com"), CreatedAt: timeNow()},
	}

	mockQuerier.On("GetUser", ctx, userID).Return(db.User{ID: userID, Email: &email}, nil)
	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &userID).Return(items, nil)
	mockEmail.On("SendEmail", ctx, mock.Anything).Return(nil)

	result, err := service.SendIntegratedDigestForUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.EmailSent)
	assert.Nil(t, result.PodcastURL)
	assert.Equal(t, 1, result.ItemsCount)
	mockQuerier.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestSendIntegratedDigestForUser_WithPodcast(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)
	service.SetPodcastGenerationEnabled(true)

	ctx := context.Background()
	userID := int32(1)
	email := "test@example.com"
	audioURL := "https://cdn.example.com/podcast.mp3"

	items := []db.Item{
		{ID: 1, Title: "Test Item", Url: stringPtr("https://example.com"), CreatedAt: timeNow()},
	}

	podcast := db.Podcast{
		ID:       1,
		Status:   "completed",
		AudioUrl: &audioURL,
	}

	mockQuerier.On("GetUser", ctx, userID).Return(db.User{ID: userID, Email: &email}, nil)
	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &userID).Return(items, nil)
	mockPodcast.On("CreatePodcastFromItems", ctx, userID, mock.Anything, mock.Anything, mock.Anything).Return(&podcast, nil)
	mockQuerier.On("GetPodcast", ctx, podcast.ID).Return(podcast, nil)
	mockEmail.On("SendEmail", ctx, mock.Anything).Return(nil)

	result, err := service.SendIntegratedDigestForUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.EmailSent)
	assert.NotNil(t, result.PodcastURL)
	assert.Equal(t, audioURL, *result.PodcastURL)
	assert.Equal(t, 1, result.ItemsCount)
	mockQuerier.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
	mockPodcast.AssertExpectations(t)
}

func TestSendIntegratedDigestForUser_NoItems(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)

	ctx := context.Background()
	userID := int32(1)
	email := "test@example.com"

	mockQuerier.On("GetUser", ctx, userID).Return(db.User{ID: userID, Email: &email}, nil)
	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &userID).Return([]db.Item{}, nil)

	result, err := service.SendIntegratedDigestForUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.EmailSent)
	assert.Equal(t, 0, result.ItemsCount)
	mockQuerier.AssertExpectations(t)
	mockEmail.AssertNotCalled(t, "SendEmail")
}

func TestSendIntegratedDigest(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)
	service.SetPodcastGenerationEnabled(false)

	ctx := context.Background()
	email1 := "user1@example.com"

	users := []db.User{
		{ID: 1, Email: &email1},
	}

	mockQuerier.On("ListUsers", ctx).Return(users, nil)
	mockQuerier.On("GetUser", ctx, int32(1)).Return(users[0], nil)
	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, mock.Anything).Return([]db.Item{
		{ID: 1, Title: "Test", Url: stringPtr("https://example.com"), CreatedAt: timeNow()},
	}, nil)
	mockEmail.On("SendEmail", ctx, mock.Anything).Return(nil)

	err := service.SendIntegratedDigest(ctx)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestWaitForPodcastCompletion_Success(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast).(*digestService)

	ctx := context.Background()
	podcastID := int32(1)
	audioURL := "https://example.com/podcast.mp3"

	completedPodcast := db.Podcast{
		ID:       podcastID,
		Status:   "completed",
		AudioUrl: &audioURL,
	}

	mockQuerier.On("GetPodcast", ctx, podcastID).Return(completedPodcast, nil)

	result, err := service.waitForPodcastCompletion(ctx, podcastID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "completed", result.Status)
	mockQuerier.AssertExpectations(t)
}

func TestWaitForPodcastCompletion_Failed(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast).(*digestService)

	ctx := context.Background()
	podcastID := int32(1)

	failedPodcast := db.Podcast{
		ID:     podcastID,
		Status: "failed",
	}

	mockQuerier.On("GetPodcast", ctx, podcastID).Return(failedPodcast, nil)

	result, err := service.waitForPodcastCompletion(ctx, podcastID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "podcast generation failed")
	mockQuerier.AssertExpectations(t)
}

// TestWaitForPodcastCompletion_Timeout is skipped because it takes 5 minutes to run
// The actual timeout behavior is tested in integration tests
func TestWaitForPodcastCompletion_Timeout(t *testing.T) {
	t.Skip("Skipping timeout test - takes 5 minutes to run")
}

func TestSendIntegratedDigestForUser_PodcastCreationFails(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockEmail := new(MockEmailService)
	mockPodcast := new(MockPodcastService)

	service := NewDigestService(mockQuerier, mockEmail, mockPodcast)
	service.SetPodcastGenerationEnabled(true)

	ctx := context.Background()
	userID := int32(1)
	email := "test@example.com"

	items := []db.Item{
		{ID: 1, Title: "Test Item", Url: stringPtr("https://example.com"), CreatedAt: timeNow()},
	}

	mockQuerier.On("GetUser", ctx, userID).Return(db.User{ID: userID, Email: &email}, nil)
	mockQuerier.On("GetUnreadItemsFromPreviousDayByUser", ctx, &userID).Return(items, nil)
	mockPodcast.On("CreatePodcastFromItems", ctx, userID, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("podcast creation failed"))
	mockEmail.On("SendEmail", ctx, mock.Anything).Return(nil)

	result, err := service.SendIntegratedDigestForUser(ctx, userID)

	// Should continue without podcast even if creation fails
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.EmailSent)
	assert.Nil(t, result.PodcastURL)
	mockQuerier.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
	mockPodcast.AssertExpectations(t)
}
