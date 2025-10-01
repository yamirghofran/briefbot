package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, name, email, authProvider, oauthID, passwordHash *string) (*db.User, error) {
	args := m.Called(ctx, name, email, authProvider, oauthID, passwordHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id int32) (*db.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email *string) (*db.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context) ([]db.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id int32, name, email, authProvider, oauthID, passwordHash *string) error {
	args := m.Called(ctx, id, name, email, authProvider, oauthID, passwordHash)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCreateUser(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewHandler(mockUserService, nil, nil, nil)

	router := setupTestRouter()
	router.POST("/users", handler.CreateUser)

	name := "John Doe"
	email := "john@example.com"

	expectedUser := &db.User{
		ID:    1,
		Name:  &name,
		Email: &email,
	}

	mockUserService.On("CreateUser", mock.Anything, &name, &email, mock.Anything, mock.Anything, mock.Anything).Return(expectedUser, nil)

	reqBody := map[string]interface{}{
		"name":  name,
		"email": email,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestCreateUser_Error(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewHandler(mockUserService, nil, nil, nil)

	router := setupTestRouter()
	router.POST("/users", handler.CreateUser)

	mockUserService.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	reqBody := map[string]interface{}{
		"name":  "John",
		"email": "john@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewHandler(mockUserService, nil, nil, nil)

	router := setupTestRouter()
	router.GET("/users/:id", handler.GetUser)

	name := "John Doe"
	expectedUser := &db.User{
		ID:   1,
		Name: &name,
	}

	mockUserService.On("GetUser", mock.Anything, int32(1)).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestGetUser_InvalidID(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewHandler(mockUserService, nil, nil, nil)

	router := setupTestRouter()
	router.GET("/users/:id", handler.GetUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserByEmail(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewHandler(mockUserService, nil, nil, nil)

	router := setupTestRouter()
	router.GET("/users/email/:email", handler.GetUserByEmail)

	email := "john@example.com"
	expectedUser := &db.User{
		ID:    1,
		Email: &email,
	}

	mockUserService.On("GetUserByEmail", mock.Anything, &email).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/email/john@example.com", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewHandler(mockUserService, nil, nil, nil)

	router := setupTestRouter()
	router.GET("/users", handler.ListUsers)

	name1 := "User 1"
	name2 := "User 2"
	expectedUsers := []db.User{
		{ID: 1, Name: &name1},
		{ID: 2, Name: &name2},
	}

	mockUserService.On("ListUsers", mock.Anything).Return(expectedUsers, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewHandler(mockUserService, nil, nil, nil)

	router := setupTestRouter()
	router.PUT("/users/:id", handler.UpdateUser)

	name := "Updated Name"

	mockUserService.On("UpdateUser", mock.Anything, int32(1), &name, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqBody := map[string]interface{}{
		"name": name,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewHandler(mockUserService, nil, nil, nil)

	router := setupTestRouter()
	router.DELETE("/users/:id", handler.DeleteUser)

	mockUserService.On("DeleteUser", mock.Anything, int32(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserService.AssertExpectations(t)
}
