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

type MockItemService struct {
	mock.Mock
}

func (m *MockItemService) CreateItemAsync(ctx context.Context, userID int32, url string) (*db.Item, error) {
	args := m.Called(ctx, userID, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) ProcessURL(ctx context.Context, userID int32, url string) (*db.Item, error) {
	args := m.Called(ctx, userID, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) CreateItem(ctx context.Context, userID *int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string) (*db.Item, error) {
	args := m.Called(ctx, userID, title, url, textContent, summary, itemType, platform, tags, authors)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) GetItem(ctx context.Context, id int32) (*db.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) GetItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockItemService) GetUnreadItemsByUser(ctx context.Context, userID *int32) ([]db.Item, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockItemService) GetUnreadItemsFromPreviousDay(ctx context.Context) ([]db.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.Item), args.Error(1)
}

func (m *MockItemService) UpdateItem(ctx context.Context, id int32, title string, url *string, textContent *string, summary *string, itemType *string, platform *string, tags []string, authors []string, isRead *bool) error {
	args := m.Called(ctx, id, title, url, textContent, summary, itemType, platform, tags, authors, isRead)
	return args.Error(0)
}

func (m *MockItemService) MarkItemAsRead(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemService) ToggleItemReadStatus(ctx context.Context, id int32) (*db.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Item), args.Error(1)
}

func (m *MockItemService) DeleteItem(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemService) GetItemProcessingStatus(ctx context.Context, itemID int32) (*services.ItemStatus, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.ItemStatus), args.Error(1)
}

func (m *MockItemService) GetItemsByProcessingStatus(ctx context.Context, status *string) ([]db.Item, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]db.Item), args.Error(1)
}

func TestCreateItem(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.POST("/items", handler.CreateItem)

	userID := int32(1)
	url := "https://example.com"
	status := "pending"

	expectedItem := &db.Item{
		ID:               1,
		UserID:           &userID,
		Url:              &url,
		ProcessingStatus: &status,
	}

	mockItemService.On("CreateItemAsync", mock.Anything, userID, url).Return(expectedItem, nil)

	reqBody := map[string]interface{}{
		"user_id": userID,
		"url":     url,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetItem(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/:id", handler.GetItem)

	expectedItem := &db.Item{
		ID:    1,
		Title: "Test Item",
	}

	mockItemService.On("GetItem", mock.Anything, int32(1)).Return(expectedItem, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetItemsByUser(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/user/:userID", handler.GetItemsByUser)

	userID := int32(1)
	expectedItems := []db.Item{
		{ID: 1, Title: "Item 1"},
		{ID: 2, Title: "Item 2"},
	}

	mockItemService.On("GetItemsByUser", mock.Anything, &userID).Return(expectedItems, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/user/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetUnreadItemsByUser(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/user/:userID/unread", handler.GetUnreadItemsByUser)

	userID := int32(1)
	expectedItems := []db.Item{
		{ID: 1, Title: "Unread Item 1"},
	}

	mockItemService.On("GetUnreadItemsByUser", mock.Anything, &userID).Return(expectedItems, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/user/1/unread", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestUpdateItem(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PUT("/items/:id", handler.UpdateItem)

	mockItemService.On("UpdateItem", mock.Anything, int32(1), "Updated Title", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqBody := map[string]interface{}{
		"title": "Updated Title",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/items/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestMarkItemAsRead(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PATCH("/items/:id/read", handler.MarkItemAsRead)

	isRead := true
	expectedItem := &db.Item{
		ID:     1,
		IsRead: &isRead,
	}

	mockItemService.On("MarkItemAsRead", mock.Anything, int32(1)).Return(nil)
	mockItemService.On("GetItem", mock.Anything, int32(1)).Return(expectedItem, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/items/1/read", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestToggleItemReadStatus(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PATCH("/items/:id/toggle-read", handler.ToggleItemReadStatus)

	isRead := true
	expectedItem := &db.Item{
		ID:     1,
		IsRead: &isRead,
	}

	mockItemService.On("ToggleItemReadStatus", mock.Anything, int32(1)).Return(expectedItem, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/items/1/toggle-read", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestDeleteItem(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.DELETE("/items/:id", handler.DeleteItem)

	mockItemService.On("DeleteItem", mock.Anything, int32(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/items/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetItemProcessingStatus(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/:id/status", handler.GetItemProcessingStatus)

	status := "completed"
	expectedStatus := &services.ItemStatus{
		Item: &db.Item{
			ID:               1,
			ProcessingStatus: &status,
		},
		IsCompleted: true,
	}

	mockItemService.On("GetItemProcessingStatus", mock.Anything, int32(1)).Return(expectedStatus, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/1/status", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetItemsByProcessingStatus(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/status", handler.GetItemsByProcessingStatus)

	status := "pending"
	expectedItems := []db.Item{
		{ID: 1, Title: "Pending Item"},
	}

	mockItemService.On("GetItemsByProcessingStatus", mock.Anything, &status).Return(expectedItems, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/status?status=pending", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetItemsByProcessingStatus_InvalidStatus(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/status", handler.GetItemsByProcessingStatus)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/status?status=invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateItem_Error(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.POST("/items", handler.CreateItem)

	mockItemService.On("CreateItemAsync", mock.Anything, int32(1), "https://example.com").Return(nil, errors.New("service error"))

	reqBody := map[string]interface{}{
		"user_id": 1,
		"url":     "https://example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestCreateItem_InvalidJSON(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.POST("/items", handler.CreateItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetItem_InvalidID(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/:id", handler.GetItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetItem_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/:id", handler.GetItem)

	mockItemService.On("GetItem", mock.Anything, int32(1)).Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetItemsByUser_InvalidID(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/user/:userID", handler.GetItemsByUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/user/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetItemsByUser_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/user/:userID", handler.GetItemsByUser)

	userID := int32(1)
	mockItemService.On("GetItemsByUser", mock.Anything, &userID).Return([]db.Item{}, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/user/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetUnreadItemsByUser_InvalidID(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/user/:userID/unread", handler.GetUnreadItemsByUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/user/invalid/unread", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUnreadItemsByUser_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/user/:userID/unread", handler.GetUnreadItemsByUser)

	userID := int32(1)
	mockItemService.On("GetUnreadItemsByUser", mock.Anything, &userID).Return([]db.Item{}, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/user/1/unread", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestUpdateItem_InvalidID(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PUT("/items/:id", handler.UpdateItem)

	reqBody := map[string]interface{}{
		"title": "Updated Title",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/items/invalid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateItem_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PUT("/items/:id", handler.UpdateItem)

	mockItemService.On("UpdateItem", mock.Anything, int32(1), "Updated Title", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("service error"))

	reqBody := map[string]interface{}{
		"title": "Updated Title",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/items/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestMarkItemAsRead_InvalidID(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PATCH("/items/:id/read", handler.MarkItemAsRead)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/items/invalid/read", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMarkItemAsRead_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PATCH("/items/:id/read", handler.MarkItemAsRead)

	mockItemService.On("MarkItemAsRead", mock.Anything, int32(1)).Return(errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/items/1/read", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestMarkItemAsRead_GetItemError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PATCH("/items/:id/read", handler.MarkItemAsRead)

	mockItemService.On("MarkItemAsRead", mock.Anything, int32(1)).Return(nil)
	mockItemService.On("GetItem", mock.Anything, int32(1)).Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/items/1/read", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestToggleItemReadStatus_InvalidID(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PATCH("/items/:id/toggle-read", handler.ToggleItemReadStatus)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/items/invalid/toggle-read", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestToggleItemReadStatus_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.PATCH("/items/:id/toggle-read", handler.ToggleItemReadStatus)

	mockItemService.On("ToggleItemReadStatus", mock.Anything, int32(1)).Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/items/1/toggle-read", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestDeleteItem_InvalidID(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.DELETE("/items/:id", handler.DeleteItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/items/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteItem_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.DELETE("/items/:id", handler.DeleteItem)

	mockItemService.On("DeleteItem", mock.Anything, int32(1)).Return(errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/items/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetItemProcessingStatus_InvalidID(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/:id/status", handler.GetItemProcessingStatus)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/invalid/status", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetItemProcessingStatus_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/:id/status", handler.GetItemProcessingStatus)

	mockItemService.On("GetItemProcessingStatus", mock.Anything, int32(1)).Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/1/status", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}

func TestGetItemsByProcessingStatus_ServiceError(t *testing.T) {
	mockItemService := new(MockItemService)
	handler := NewHandler(nil, mockItemService, nil, nil)

	router := setupTestRouter()
	router.GET("/items/status", handler.GetItemsByProcessingStatus)

	status := "pending"
	mockItemService.On("GetItemsByProcessingStatus", mock.Anything, &status).Return([]db.Item{}, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/status?status=pending", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockItemService.AssertExpectations(t)
}
