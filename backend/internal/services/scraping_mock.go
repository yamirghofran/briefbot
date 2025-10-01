package services

import (
	"github.com/stretchr/testify/mock"
)

type MockScrapingService struct {
	mock.Mock
}

func (m *MockScrapingService) Scrape(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}
