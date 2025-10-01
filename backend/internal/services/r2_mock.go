package services

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockR2Service struct {
	mock.Mock
}

func (m *MockR2Service) GenerateUploadURL(ctx context.Context, filename string, contentType string) (*UploadURLResponse, error) {
	args := m.Called(ctx, filename, contentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UploadURLResponse), args.Error(1)
}

func (m *MockR2Service) GetPublicURL(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockR2Service) DeleteFile(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockR2Service) DeleteFiles(ctx context.Context, keys []string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockR2Service) ExtractKeyFromURL(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func (m *MockR2Service) GenerateUploadURLForKey(ctx context.Context, key string, contentType string) (*UploadURLResponse, error) {
	args := m.Called(ctx, key, contentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UploadURLResponse), args.Error(1)
}

func (m *MockR2Service) UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	args := m.Called(ctx, key, data, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockR2Service) UploadFileWithMetadata(ctx context.Context, key string, data []byte, contentType string, metadata map[string]string) (string, error) {
	args := m.Called(ctx, key, data, contentType, metadata)
	return args.String(0), args.Error(1)
}
