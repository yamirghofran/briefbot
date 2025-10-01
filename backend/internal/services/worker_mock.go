package services

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockWorkerService struct {
	mock.Mock
}

func (m *MockWorkerService) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockWorkerService) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWorkerService) IsRunning() bool {
	args := m.Called()
	return args.Bool(0)
}
