package services

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAIService struct {
	mock.Mock
}

func (m *MockAIService) ExtractContent(ctx context.Context, content string) (ItemExtraction, error) {
	args := m.Called(ctx, content)
	return args.Get(0).(ItemExtraction), args.Error(1)
}

func (m *MockAIService) SummarizeContent(ctx context.Context, content string) (ItemSummary, error) {
	args := m.Called(ctx, content)
	return args.Get(0).(ItemSummary), args.Error(1)
}

func (m *MockAIService) WritePodcast(content string) (Podcast, error) {
	args := m.Called(content)
	return args.Get(0).(Podcast), args.Error(1)
}
