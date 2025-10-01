package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
	"github.com/yamirghofran/briefbot/internal/services"
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
	return args.Get(0).([]db.Podcast), args.Error(1)
}

func (m *MockPodcastService) GetPodcastsByStatus(ctx context.Context, status services.PodcastStatus) ([]db.Podcast, error) {
	args := m.Called(ctx, status)
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
	return args.Get(0).([]db.GetPodcastItemsRow), args.Error(1)
}

func (m *MockPodcastService) UpdatePodcastStatus(ctx context.Context, podcastID int32, status services.PodcastStatus) error {
	args := m.Called(ctx, podcastID, status)
	return args.Error(0)
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

func (m *MockPodcastService) GeneratePodcastUploadURL(ctx context.Context, podcastID int32) (*services.UploadURLResponse, error) {
	args := m.Called(ctx, podcastID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.UploadURLResponse), args.Error(1)
}

func TestCreatePodcast(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts", handler.CreatePodcast)

	status := "pending"
	expectedPodcast := &db.Podcast{
		ID:     1,
		Status: status,
	}

	mockPodcastService.On("CreatePodcastFromItems", mock.Anything, int32(1), "Test Podcast", "Description", []int32{1, 2}).Return(expectedPodcast, nil)

	reqBody := map[string]interface{}{
		"user_id":     1,
		"title":       "Test Podcast",
		"description": "Description",
		"item_ids":    []int{1, 2},
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestCreatePodcastFromSingleItem(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts/from-item", handler.CreatePodcastFromSingleItem)

	status := "pending"
	expectedPodcast := &db.Podcast{
		ID:     1,
		Status: status,
	}

	mockPodcastService.On("CreatePodcastFromSingleItem", mock.Anything, int32(1), int32(1)).Return(expectedPodcast, nil)

	reqBody := map[string]interface{}{
		"user_id": 1,
		"item_id": 1,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts/from-item", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcast(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id", handler.GetPodcast)

	expectedPodcast := &db.Podcast{
		ID: 1,
	}

	mockPodcastService.On("GetPodcast", mock.Anything, int32(1)).Return(expectedPodcast, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestCreatePodcast_ValidationError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts", handler.CreatePodcast)

	// Missing required fields
	reqBody := map[string]interface{}{
		"user_id": 1,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePodcast_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts", handler.CreatePodcast)

	mockPodcastService.On("CreatePodcastFromItems", mock.Anything, int32(1), "Test", "Desc", []int32{1}).Return(nil, errors.New("service error"))

	reqBody := map[string]interface{}{
		"user_id":     1,
		"title":       "Test",
		"description": "Desc",
		"item_ids":    []int{1},
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}
