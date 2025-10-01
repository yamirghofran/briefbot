package services

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(ctx context.Context, request EmailRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}
