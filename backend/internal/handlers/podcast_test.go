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

func TestGetPodcastsByUser(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/user/:userID", handler.GetPodcastsByUser)

	expectedPodcasts := []db.Podcast{
		{ID: 1, Title: "Podcast 1"},
		{ID: 2, Title: "Podcast 2"},
	}

	mockPodcastService.On("GetPodcastsByUser", mock.Anything, int32(1)).Return(expectedPodcasts, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/user/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastsByUser_InvalidID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/user/:userID", handler.GetPodcastsByUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/user/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPodcastsByStatus(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/status/:status", handler.GetPodcastsByStatus)

	expectedPodcasts := []db.Podcast{
		{ID: 1, Status: "pending"},
	}

	mockPodcastService.On("GetPodcastsByStatus", mock.Anything, services.PodcastStatus("pending")).Return(expectedPodcasts, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/status/pending", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastsByStatus_InvalidStatus(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/status/:status", handler.GetPodcastsByStatus)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/status/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPendingPodcasts(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/pending", handler.GetPendingPodcasts)

	expectedPodcasts := []db.Podcast{
		{ID: 1, Status: "pending"},
		{ID: 2, Status: "pending"},
	}

	mockPodcastService.On("GetPendingPodcasts", mock.Anything, int32(10)).Return(expectedPodcasts, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/pending", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPendingPodcasts_WithLimit(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/pending", handler.GetPendingPodcasts)

	expectedPodcasts := []db.Podcast{
		{ID: 1, Status: "pending"},
	}

	mockPodcastService.On("GetPendingPodcasts", mock.Anything, int32(5)).Return(expectedPodcasts, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/pending?limit=5", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastItems(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/items", handler.GetPodcastItems)

	title1 := "Item 1"
	title2 := "Item 2"
	expectedItems := []db.GetPodcastItemsRow{
		{ID: 1, Title: title1},
		{ID: 2, Title: title2},
	}

	mockPodcastService.On("GetPodcastItems", mock.Anything, int32(1)).Return(expectedItems, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/items", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastAudio(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/audio", handler.GetPodcastAudio)

	audioURL := "https://example.com/audio.mp3"
	duration := int32(300)
	expectedPodcast := &db.Podcast{
		ID:              1,
		AudioUrl:        &audioURL,
		DurationSeconds: &duration,
	}

	mockPodcastService.On("HasPodcastAudio", mock.Anything, int32(1)).Return(true, nil)
	mockPodcastService.On("GetPodcast", mock.Anything, int32(1)).Return(expectedPodcast, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/audio", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastAudio_NotAvailable(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/audio", handler.GetPodcastAudio)

	mockPodcastService.On("HasPodcastAudio", mock.Anything, int32(1)).Return(false, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/audio", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGeneratePodcastUploadURL(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/upload-url", handler.GeneratePodcastUploadURL)

	expectedResponse := &services.UploadURLResponse{
		UploadURL: "https://example.com/upload",
		Key:       "podcasts/1/audio.mp3",
	}

	mockPodcastService.On("GeneratePodcastUploadURL", mock.Anything, int32(1)).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/upload-url", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestAddItemToPodcast(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts/:id/items", handler.AddItemToPodcast)

	mockPodcastService.On("AddItemToPodcast", mock.Anything, int32(1), int32(5), 0).Return(nil)

	reqBody := map[string]interface{}{
		"item_id": 5,
		"order":   0,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts/1/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestRemoveItemFromPodcast(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.DELETE("/podcasts/:id/items/:itemID", handler.RemoveItemFromPodcast)

	mockPodcastService.On("RemoveItemFromPodcast", mock.Anything, int32(1), int32(5)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/podcasts/1/items/5", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestUpdatePodcast(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.PUT("/podcasts/:id", handler.UpdatePodcast)

	mockPodcastService.On("UpdatePodcast", mock.Anything, int32(1), "Updated Title", "Updated Description").Return(nil)

	reqBody := map[string]interface{}{
		"title":       "Updated Title",
		"description": "Updated Description",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/podcasts/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestDeletePodcast(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.DELETE("/podcasts/:id", handler.DeletePodcast)

	mockPodcastService.On("DeletePodcast", mock.Anything, int32(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/podcasts/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcast_InvalidID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id", handler.GetPodcast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPodcast_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id", handler.GetPodcast)

	mockPodcastService.On("GetPodcast", mock.Anything, int32(1)).Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastsByUser_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/user/:userID", handler.GetPodcastsByUser)

	mockPodcastService.On("GetPodcastsByUser", mock.Anything, int32(1)).Return([]db.Podcast{}, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/user/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastsByStatus_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/status/:status", handler.GetPodcastsByStatus)

	mockPodcastService.On("GetPodcastsByStatus", mock.Anything, services.PodcastStatus("pending")).Return([]db.Podcast{}, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/status/pending", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPendingPodcasts_InvalidLimit(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/pending", handler.GetPendingPodcasts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/pending?limit=invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPendingPodcasts_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/pending", handler.GetPendingPodcasts)

	mockPodcastService.On("GetPendingPodcasts", mock.Anything, int32(10)).Return([]db.Podcast{}, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/pending", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastItems_InvalidID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/items", handler.GetPodcastItems)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/invalid/items", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPodcastItems_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/items", handler.GetPodcastItems)

	mockPodcastService.On("GetPodcastItems", mock.Anything, int32(1)).Return([]db.GetPodcastItemsRow{}, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/items", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastAudio_InvalidID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/audio", handler.GetPodcastAudio)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/invalid/audio", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPodcastAudio_HasAudioError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/audio", handler.GetPodcastAudio)

	mockPodcastService.On("HasPodcastAudio", mock.Anything, int32(1)).Return(false, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/audio", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastAudio_GetPodcastError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/audio", handler.GetPodcastAudio)

	mockPodcastService.On("HasPodcastAudio", mock.Anything, int32(1)).Return(true, nil)
	mockPodcastService.On("GetPodcast", mock.Anything, int32(1)).Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/audio", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGetPodcastAudio_NoAudioURL(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/audio", handler.GetPodcastAudio)

	mockPodcastService.On("HasPodcastAudio", mock.Anything, int32(1)).Return(true, nil)
	mockPodcastService.On("GetPodcast", mock.Anything, int32(1)).Return(&db.Podcast{ID: 1, AudioUrl: nil}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/audio", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestGeneratePodcastUploadURL_InvalidID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/upload-url", handler.GeneratePodcastUploadURL)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/invalid/upload-url", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGeneratePodcastUploadURL_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.GET("/podcasts/:id/upload-url", handler.GeneratePodcastUploadURL)

	mockPodcastService.On("GeneratePodcastUploadURL", mock.Anything, int32(1)).Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/podcasts/1/upload-url", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestAddItemToPodcast_InvalidPodcastID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts/:id/items", handler.AddItemToPodcast)

	reqBody := map[string]interface{}{
		"item_id": 5,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts/invalid/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddItemToPodcast_ValidationError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts/:id/items", handler.AddItemToPodcast)

	reqBody := map[string]interface{}{
		"order": 0,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts/1/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddItemToPodcast_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts/:id/items", handler.AddItemToPodcast)

	mockPodcastService.On("AddItemToPodcast", mock.Anything, int32(1), int32(5), 0).Return(errors.New("service error"))

	reqBody := map[string]interface{}{
		"item_id": 5,
		"order":   0,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts/1/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestRemoveItemFromPodcast_InvalidPodcastID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.DELETE("/podcasts/:id/items/:itemID", handler.RemoveItemFromPodcast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/podcasts/invalid/items/5", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemoveItemFromPodcast_InvalidItemID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.DELETE("/podcasts/:id/items/:itemID", handler.RemoveItemFromPodcast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/podcasts/1/items/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemoveItemFromPodcast_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.DELETE("/podcasts/:id/items/:itemID", handler.RemoveItemFromPodcast)

	mockPodcastService.On("RemoveItemFromPodcast", mock.Anything, int32(1), int32(5)).Return(errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/podcasts/1/items/5", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestUpdatePodcast_InvalidID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.PUT("/podcasts/:id", handler.UpdatePodcast)

	reqBody := map[string]interface{}{
		"title":       "Updated Title",
		"description": "Updated Description",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/podcasts/invalid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePodcast_ValidationError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.PUT("/podcasts/:id", handler.UpdatePodcast)

	reqBody := map[string]interface{}{
		"description": "Updated Description",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/podcasts/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePodcast_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.PUT("/podcasts/:id", handler.UpdatePodcast)

	mockPodcastService.On("UpdatePodcast", mock.Anything, int32(1), "Updated Title", "Updated Description").Return(errors.New("service error"))

	reqBody := map[string]interface{}{
		"title":       "Updated Title",
		"description": "Updated Description",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/podcasts/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestDeletePodcast_InvalidID(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.DELETE("/podcasts/:id", handler.DeletePodcast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/podcasts/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeletePodcast_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.DELETE("/podcasts/:id", handler.DeletePodcast)

	mockPodcastService.On("DeletePodcast", mock.Anything, int32(1)).Return(errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/podcasts/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}

func TestCreatePodcastFromSingleItem_ValidationError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts/from-item", handler.CreatePodcastFromSingleItem)

	reqBody := map[string]interface{}{
		"user_id": 1,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts/from-item", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePodcastFromSingleItem_ServiceError(t *testing.T) {
	mockPodcastService := new(MockPodcastService)
	handler := NewPodcastHandler(mockPodcastService)

	router := setupTestRouter()
	router.POST("/podcasts/from-item", handler.CreatePodcastFromSingleItem)

	mockPodcastService.On("CreatePodcastFromSingleItem", mock.Anything, int32(1), int32(1)).Return(nil, errors.New("service error"))

	reqBody := map[string]interface{}{
		"user_id": 1,
		"item_id": 1,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/podcasts/from-item", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPodcastService.AssertExpectations(t)
}
