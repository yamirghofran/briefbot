package services

import (
	"context"

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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
