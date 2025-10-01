package services

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
	"github.com/yamirghofran/briefbot/internal/test"
)

func TestDefaultPodcastConfig(t *testing.T) {
	config := DefaultPodcastConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 1.0, config.DefaultSpeed)
	assert.Equal(t, 3, config.MaxRetries)
	assert.NotEmpty(t, config.TempDir)
	assert.True(t, config.EnableStorage)
	assert.NotNil(t, config.VoiceMapping)
}

func TestPodcastConfig_Structure(t *testing.T) {
	config := PodcastConfig{
		DefaultSpeed:       1.5,
		MaxRetries:         5,
		TempDir:            "/tmp",
		EnableStorage:      true,
		VoiceMapping:       make(map[string]VoiceEnum),
		MaxItemsPerPodcast: 10,
		MaxConcurrentAudio: 5,
	}

	assert.Equal(t, 1.5, config.DefaultSpeed)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, "/tmp", config.TempDir)
	assert.True(t, config.EnableStorage)
	assert.Equal(t, 10, config.MaxItemsPerPodcast)
	assert.Equal(t, int32(5), config.MaxConcurrentAudio)
}

func TestCreatePodcastFromItems(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	userID := int32(1)
	title := "Test Podcast"
	description := "Test Description"
	itemIDs := []int32{1, 2, 3}

	expectedPodcast := db.Podcast{
		ID:          1,
		UserID:      &userID,
		Title:       title,
		Description: &description,
		Status:      "pending",
	}

	mockQuerier.On("CreatePodcast", ctx, mock.MatchedBy(func(params db.CreatePodcastParams) bool {
		return params.Title == title && *params.UserID == userID
	})).Return(expectedPodcast, nil)

	for i, itemID := range itemIDs {
		mockQuerier.On("AddItemToPodcast", ctx, mock.MatchedBy(func(params db.AddItemToPodcastParams) bool {
			return *params.PodcastID == expectedPodcast.ID && *params.ItemID == itemID && params.ItemOrder == int32(i)
		})).Return(db.PodcastItem{}, nil)
	}

	podcast, err := service.CreatePodcastFromItems(ctx, userID, title, description, itemIDs)

	assert.NoError(t, err)
	assert.NotNil(t, podcast)
	assert.Equal(t, title, podcast.Title)
	mockQuerier.AssertExpectations(t)
}

func TestCreatePodcastFromItems_TooManyItems(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	userID := int32(1)

	// Create more items than allowed
	itemIDs := make([]int32, config.MaxItemsPerPodcast+1)
	for i := range itemIDs {
		itemIDs[i] = int32(i + 1)
	}

	podcast, err := service.CreatePodcastFromItems(ctx, userID, "Test", "Desc", itemIDs)

	assert.Error(t, err)
	assert.Nil(t, podcast)
	assert.Contains(t, err.Error(), "too many items")
}

func TestCreatePodcastFromSingleItem(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	userID := int32(1)
	itemID := int32(1)
	itemTitle := "Test Item"
	itemSummary := "Test Summary"

	item := db.Item{
		ID:      itemID,
		Title:   itemTitle,
		Summary: &itemSummary,
	}

	expectedPodcast := db.Podcast{
		ID:     1,
		UserID: &userID,
		Title:  fmt.Sprintf("Podcast: %s", itemTitle),
		Status: "pending",
	}

	mockQuerier.On("GetItem", ctx, itemID).Return(item, nil)
	mockQuerier.On("CreatePodcast", ctx, mock.Anything).Return(expectedPodcast, nil)
	mockQuerier.On("AddItemToPodcast", ctx, mock.Anything).Return(db.PodcastItem{}, nil)

	podcast, err := service.CreatePodcastFromSingleItem(ctx, userID, itemID)

	assert.NoError(t, err)
	assert.NotNil(t, podcast)
	mockQuerier.AssertExpectations(t)
}

func TestGetPodcast(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)

	expectedPodcast := db.Podcast{
		ID:     podcastID,
		Title:  "Test Podcast",
		Status: "completed",
	}

	mockQuerier.On("GetPodcast", ctx, podcastID).Return(expectedPodcast, nil)

	podcast, err := service.GetPodcast(ctx, podcastID)

	assert.NoError(t, err)
	assert.NotNil(t, podcast)
	assert.Equal(t, podcastID, podcast.ID)
	mockQuerier.AssertExpectations(t)
}

func TestGetPodcastsByUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	userID := int32(1)

	expectedPodcasts := []db.Podcast{
		{ID: 1, Title: "Podcast 1"},
		{ID: 2, Title: "Podcast 2"},
	}

	mockQuerier.On("GetPodcastByUser", ctx, &userID).Return(expectedPodcasts, nil)

	podcasts, err := service.GetPodcastsByUser(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, podcasts, 2)
	mockQuerier.AssertExpectations(t)
}

func TestUpdatePodcast(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)
	title := "Updated Title"
	description := "Updated Description"

	mockQuerier.On("UpdatePodcast", ctx, mock.MatchedBy(func(params db.UpdatePodcastParams) bool {
		return params.ID == podcastID && params.Title == title
	})).Return(nil)

	err := service.UpdatePodcast(ctx, podcastID, title, description)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestDeletePodcast(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)

	mockQuerier.On("ClearPodcastItems", ctx, &podcastID).Return(nil)
	mockQuerier.On("DeletePodcast", ctx, podcastID).Return(nil)

	err := service.DeletePodcast(ctx, podcastID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestAddItemToPodcast(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)
	itemID := int32(1)
	order := 0

	mockQuerier.On("AddItemToPodcast", ctx, mock.MatchedBy(func(params db.AddItemToPodcastParams) bool {
		return *params.PodcastID == podcastID && *params.ItemID == itemID
	})).Return(db.PodcastItem{}, nil)

	err := service.AddItemToPodcast(ctx, podcastID, itemID, order)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestRemoveItemFromPodcast(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)
	itemID := int32(1)

	mockQuerier.On("RemoveItemFromPodcast", ctx, mock.MatchedBy(func(params db.RemoveItemFromPodcastParams) bool {
		return *params.PodcastID == podcastID && *params.ItemID == itemID
	})).Return(nil)

	err := service.RemoveItemFromPodcast(ctx, podcastID, itemID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestUpdatePodcastStatus(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)
	status := PodcastStatusCompleted

	mockQuerier.On("UpdatePodcastStatus", ctx, mock.MatchedBy(func(params db.UpdatePodcastStatusParams) bool {
		return params.ID == podcastID && params.Status == string(status)
	})).Return(nil)

	err := service.UpdatePodcastStatus(ctx, podcastID, status)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestGetPodcastItems(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)

	expectedItems := []db.GetPodcastItemsRow{
		{ID: 1, Title: "Item 1"},
		{ID: 2, Title: "Item 2"},
	}

	mockQuerier.On("GetPodcastItems", ctx, &podcastID).Return(expectedItems, nil)

	items, err := service.GetPodcastItems(ctx, podcastID)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	mockQuerier.AssertExpectations(t)
}

func TestGetPendingPodcasts(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	limit := int32(10)

	expectedPodcasts := []db.Podcast{
		{ID: 1, Status: "pending"},
		{ID: 2, Status: "pending"},
	}

	mockQuerier.On("GetPendingPodcasts", ctx, limit).Return(expectedPodcasts, nil)

	podcasts, err := service.GetPendingPodcasts(ctx, limit)

	assert.NoError(t, err)
	assert.Len(t, podcasts, 2)
	mockQuerier.AssertExpectations(t)
}

func TestAcquirePendingPodcasts(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	limit := int32(5)

	expectedPodcasts := []db.Podcast{
		{ID: 1, Status: "pending"},
		{ID: 2, Status: "pending"},
	}

	mockQuerier.On("GetPendingPodcasts", ctx, limit).Return(expectedPodcasts, nil)

	// Mock UpdatePodcastStatus for each podcast
	for _, podcast := range expectedPodcasts {
		mockQuerier.On("UpdatePodcastStatus", ctx, mock.MatchedBy(func(params db.UpdatePodcastStatusParams) bool {
			return params.ID == podcast.ID && params.Status == "writing"
		})).Return(nil)
	}

	podcasts, err := service.AcquirePendingPodcasts(ctx, limit)

	assert.NoError(t, err)
	assert.Len(t, podcasts, 2)
	mockQuerier.AssertExpectations(t)
}

func TestGeneratePodcastScript(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)

	items := []db.GetPodcastItemsRow{
		{ID: 1, Title: "Item 1", Summary: stringPtr("Summary 1")},
		{ID: 2, Title: "Item 2", Summary: stringPtr("Summary 2")},
	}

	podcastData := Podcast{
		Dialogues: []Dialogue{
			{Speaker: "heart", Content: "Hello"},
			{Speaker: "adam", Content: "Hi there"},
		},
	}

	dialoguesJSON, _ := json.Marshal(podcastData.Dialogues)

	mockQuerier.On("GetPodcastItems", ctx, &podcastID).Return(items, nil)
	mockQuerier.On("UpdatePodcastStatus", ctx, mock.MatchedBy(func(params db.UpdatePodcastStatusParams) bool {
		return params.ID == podcastID && params.Status == "writing"
	})).Return(nil)
	mockAI.On("WritePodcast", mock.Anything).Return(podcastData, nil)
	mockQuerier.On("UpdatePodcastDialogues", ctx, mock.MatchedBy(func(params db.UpdatePodcastDialoguesParams) bool {
		return params.ID == podcastID
	})).Return(nil)

	err := service.GeneratePodcastScript(ctx, podcastID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
	mockAI.AssertExpectations(t)
	_ = dialoguesJSON // Use the variable to avoid unused error
}

func TestProcessPodcast(t *testing.T) {
	t.Skip("Skipping - ProcessPodcast requires complex mocking of audio generation with file I/O")
	// This test would need to mock:
	// - Script generation (done)
	// - Audio file generation and download
	// - Temp file creation and cleanup
	// - R2 upload
	// - Multiple DB calls for status updates
	// Better tested in integration tests
}

func TestHasPodcastAudio(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockSpeech := new(MockSpeechService)
	var mockR2 *R2Service = nil

	config := DefaultPodcastConfig()
	service := NewPodcastService(mockQuerier, mockAI, mockSpeech, mockR2, config)

	ctx := context.Background()
	podcastID := int32(1)
	audioURL := "https://example.com/podcast.mp3"

	podcast := db.Podcast{
		ID:       podcastID,
		AudioUrl: &audioURL,
	}

	mockQuerier.On("GetPodcast", ctx, podcastID).Return(podcast, nil)

	hasAudio, err := service.HasPodcastAudio(ctx, podcastID)

	assert.NoError(t, err)
	assert.True(t, hasAudio)
	mockQuerier.AssertExpectations(t)
}
