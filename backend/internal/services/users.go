package services

import (
	"context"

	"github.com/yamirghofran/briefbot/internal/db"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email, authProvider, oauthID, passwordHash *string) (*db.User, error)
	GetUser(ctx context.Context, id int32) (*db.User, error)
	GetUserByEmail(ctx context.Context, email *string) (*db.User, error)
	UpdateUser(ctx context.Context, id int32, name, email, authProvider, oauthID, passwordHash *string) error
	DeleteUser(ctx context.Context, id int32) error
}

type userService struct {
	querier db.Querier
}

func NewUserService(querier db.Querier) UserService {
	return &userService{querier: querier}
}

func (s *userService) CreateUser(ctx context.Context, name, email, authProvider, oauthID, passwordHash *string) (*db.User, error) {
	params := db.CreateUserParams{
		Name:         name,
		Email:        email,
		AuthProvider: authProvider,
		OauthID:      oauthID,
		PasswordHash: passwordHash,
	}
	user, err := s.querier.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUser(ctx context.Context, id int32) (*db.User, error) {
	user, err := s.querier.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email *string) (*db.User, error) {
	user, err := s.querier.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, name, email, authProvider, oauthID, passwordHash *string) error {
	params := db.UpdateUserParams{
		ID:           id,
		Name:         name,
		Email:        email,
		AuthProvider: authProvider,
		OauthID:      oauthID,
		PasswordHash: passwordHash,
	}
	return s.querier.UpdateUser(ctx, params)
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	return s.querier.DeleteUser(ctx, id)
}
