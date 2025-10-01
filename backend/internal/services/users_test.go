package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
	"github.com/yamirghofran/briefbot/internal/test"
)

func TestCreateUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	service := NewUserService(mockQuerier)

	ctx := context.Background()
	name := "John Doe"
	email := "john@example.com"
	authProvider := "google"
	oauthID := "oauth123"
	var passwordHash *string

	expectedUser := db.User{
		ID:           1,
		Name:         &name,
		Email:        &email,
		AuthProvider: &authProvider,
		OauthID:      &oauthID,
	}

	mockQuerier.On("CreateUser", ctx, mock.MatchedBy(func(params db.CreateUserParams) bool {
		return *params.Name == name && *params.Email == email
	})).Return(expectedUser, nil)

	user, err := service.CreateUser(ctx, &name, &email, &authProvider, &oauthID, passwordHash)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int32(1), user.ID)
	assert.Equal(t, name, *user.Name)
	mockQuerier.AssertExpectations(t)
}

func TestCreateUser_Error(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	service := NewUserService(mockQuerier)

	ctx := context.Background()
	name := "John Doe"
	email := "john@example.com"

	mockQuerier.On("CreateUser", ctx, mock.Anything).Return(db.User{}, errors.New("database error"))

	user, err := service.CreateUser(ctx, &name, &email, nil, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockQuerier.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	service := NewUserService(mockQuerier)

	ctx := context.Background()
	userID := int32(1)
	name := "John Doe"

	expectedUser := db.User{
		ID:   userID,
		Name: &name,
	}

	mockQuerier.On("GetUser", ctx, userID).Return(expectedUser, nil)

	user, err := service.GetUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	mockQuerier.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	service := NewUserService(mockQuerier)

	ctx := context.Background()
	email := "john@example.com"
	name := "John Doe"

	expectedUser := db.User{
		ID:    1,
		Name:  &name,
		Email: &email,
	}

	mockQuerier.On("GetUserByEmail", ctx, &email).Return(expectedUser, nil)

	user, err := service.GetUserByEmail(ctx, &email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, *user.Email)
	mockQuerier.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	service := NewUserService(mockQuerier)

	ctx := context.Background()
	name1 := "User 1"
	name2 := "User 2"

	expectedUsers := []db.User{
		{ID: 1, Name: &name1},
		{ID: 2, Name: &name2},
	}

	mockQuerier.On("ListUsers", ctx).Return(expectedUsers, nil)

	users, err := service.ListUsers(ctx)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	service := NewUserService(mockQuerier)

	ctx := context.Background()
	userID := int32(1)
	name := "Updated Name"
	email := "updated@example.com"

	mockQuerier.On("UpdateUser", ctx, mock.MatchedBy(func(params db.UpdateUserParams) bool {
		return params.ID == userID && *params.Name == name
	})).Return(nil)

	err := service.UpdateUser(ctx, userID, &name, &email, nil, nil, nil)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	service := NewUserService(mockQuerier)

	ctx := context.Background()
	userID := int32(1)

	mockQuerier.On("DeleteUser", ctx, userID).Return(nil)

	err := service.DeleteUser(ctx, userID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}
