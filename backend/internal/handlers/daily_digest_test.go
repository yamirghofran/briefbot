package handlers

import (
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

type MockDigestService struct {
	mock.Mock
}

func (m *MockDigestService) SendDailyDigest(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDigestService) SendDailyDigestForUser(ctx context.Context, userID int32) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockDigestService) GetDailyDigestItemsForUser(ctx context.Context, userID int32) ([]db.Item, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockDigestService) SendIntegratedDigest(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDigestService) SendIntegratedDigestForUser(ctx context.Context, userID int32) (*services.DigestResult, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.DigestResult), args.Error(1)
}

func (m *MockDigestService) SetPodcastGenerationEnabled(enabled bool) {
	m.Called(enabled)
}

func (m *MockDigestService) IsPodcastGenerationEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestTriggerDailyDigest(t *testing.T) {
	mockDigestService := new(MockDigestService)
	handler := NewHandler(nil, nil, mockDigestService, nil)

	router := setupTestRouter()
	router.POST("/digest/trigger", handler.TriggerDailyDigest)

	mockDigestService.On("SendDailyDigest", mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/digest/trigger", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDigestService.AssertExpectations(t)
}

func TestTriggerDailyDigest_ServiceUnavailable(t *testing.T) {
	handler := NewHandler(nil, nil, nil, nil) // No digest service

	router := setupTestRouter()
	router.POST("/digest/trigger", handler.TriggerDailyDigest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/digest/trigger", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestTriggerDailyDigest_Error(t *testing.T) {
	mockDigestService := new(MockDigestService)
	handler := NewHandler(nil, nil, mockDigestService, nil)

	router := setupTestRouter()
	router.POST("/digest/trigger", handler.TriggerDailyDigest)

	mockDigestService.On("SendDailyDigest", mock.Anything).Return(errors.New("digest error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/digest/trigger", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockDigestService.AssertExpectations(t)
}

func TestTriggerDailyDigestForUser(t *testing.T) {
	mockDigestService := new(MockDigestService)
	handler := NewHandler(nil, nil, mockDigestService, nil)

	router := setupTestRouter()
	router.POST("/digest/trigger/user/:userID", handler.TriggerDailyDigestForUser)

	expectedItems := []db.Item{
		{ID: 1, Title: "Item 1"},
	}

	mockDigestService.On("GetDailyDigestItemsForUser", mock.Anything, int32(1)).Return(expectedItems, nil)
	mockDigestService.On("SendDailyDigestForUser", mock.Anything, int32(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/digest/trigger/user/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDigestService.AssertExpectations(t)
}

func TestTriggerDailyDigestForUser_NoItems(t *testing.T) {
	mockDigestService := new(MockDigestService)
	handler := NewHandler(nil, nil, mockDigestService, nil)

	router := setupTestRouter()
	router.POST("/digest/trigger/user/:userID", handler.TriggerDailyDigestForUser)

	mockDigestService.On("GetDailyDigestItemsForUser", mock.Anything, int32(1)).Return([]db.Item{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/digest/trigger/user/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["message"], "No unread items")
}

func TestTriggerDailyDigestForUser_InvalidID(t *testing.T) {
	mockDigestService := new(MockDigestService)
	handler := NewHandler(nil, nil, mockDigestService, nil)

	router := setupTestRouter()
	router.POST("/digest/trigger/user/:userID", handler.TriggerDailyDigestForUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/digest/trigger/user/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Note: TriggerIntegratedDigest and TriggerIntegratedDigestForUser tests are skipped
// because they spawn goroutines that make them difficult to test with mocks.
// These would be better tested with integration tests.
