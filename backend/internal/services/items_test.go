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

func TestCreateItem(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	userID := int32(1)
	title := "Test Article"
	url := "https://example.com/article"
	content := "Article content"
	summary := "Article summary"
	itemType := "article"
	platform := "web"

	expectedItem := db.Item{
		ID:      1,
		UserID:  &userID,
		Title:   title,
		Url:     &url,
		Summary: &summary,
	}

	mockQuerier.On("CreateItem", ctx, mock.MatchedBy(func(params db.CreateItemParams) bool {
		return params.Title == title && *params.UserID == userID
	})).Return(expectedItem, nil)

	item, err := service.CreateItem(ctx, &userID, title, &url, &content, &summary, &itemType, &platform, []string{}, []string{})

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, title, item.Title)
	mockQuerier.AssertExpectations(t)
}

func TestGetItem(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	itemID := int32(1)
	title := "Test Item"

	expectedItem := db.Item{
		ID:    itemID,
		Title: title,
	}

	mockQuerier.On("GetItem", ctx, itemID).Return(expectedItem, nil)

	item, err := service.GetItem(ctx, itemID)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, itemID, item.ID)
	mockQuerier.AssertExpectations(t)
}

func TestGetItemsByUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	userID := int32(1)

	expectedItems := []db.Item{
		{ID: 1, Title: "Item 1"},
		{ID: 2, Title: "Item 2"},
	}

	mockQuerier.On("GetItemsByUser", ctx, &userID).Return(expectedItems, nil)

	items, err := service.GetItemsByUser(ctx, &userID)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	mockQuerier.AssertExpectations(t)
}

func TestGetUnreadItemsByUser(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	userID := int32(1)

	expectedItems := []db.Item{
		{ID: 1, Title: "Unread Item 1"},
		{ID: 2, Title: "Unread Item 2"},
	}

	mockQuerier.On("GetUnreadItemsByUser", ctx, &userID).Return(expectedItems, nil)

	items, err := service.GetUnreadItemsByUser(ctx, &userID)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateItem(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	itemID := int32(1)
	title := "Updated Title"
	isRead := true

	mockQuerier.On("UpdateItem", ctx, mock.MatchedBy(func(params db.UpdateItemParams) bool {
		return params.ID == itemID && params.Title == title
	})).Return(nil)

	err := service.UpdateItem(ctx, itemID, title, nil, nil, nil, nil, nil, []string{}, []string{}, &isRead)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestMarkItemAsRead(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	itemID := int32(1)

	mockQuerier.On("MarkItemAsRead", ctx, itemID).Return(nil)

	err := service.MarkItemAsRead(ctx, itemID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestToggleItemReadStatus(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	itemID := int32(1)
	isRead := true

	expectedItem := db.Item{
		ID:     itemID,
		IsRead: &isRead,
	}

	mockQuerier.On("ToggleItemReadStatus", ctx, itemID).Return(expectedItem, nil)

	item, err := service.ToggleItemReadStatus(ctx, itemID)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.True(t, *item.IsRead)
	mockQuerier.AssertExpectations(t)
}

func TestDeleteItem(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	itemID := int32(1)

	mockQuerier.On("DeleteItem", ctx, itemID).Return(nil)

	err := service.DeleteItem(ctx, itemID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestCreateItemAsync(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	userID := int32(1)
	url := "https://example.com"

	expectedItem := &db.Item{
		ID:     1,
		UserID: &userID,
		Url:    &url,
	}

	mockJobQueue.On("EnqueueItem", ctx, userID, url, url).Return(expectedItem, nil)

	item, err := service.CreateItemAsync(ctx, userID, url)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	mockJobQueue.AssertExpectations(t)
}

func TestProcessURL(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	userID := int32(1)
	url := "https://example.com"
	content := "scraped content"
	title := "Test Title"
	itemType := "article"
	platform := "web"

	extraction := ItemExtraction{
		Title:    title,
		Type:     itemType,
		Platform: platform,
		Tags:     []string{"tag1"},
		Authors:  []string{"author1"},
	}

	summary := ItemSummary{
		Overview:  "Overview text",
		KeyPoints: []string{"Point 1", "Point 2"},
	}

	mockScraper.On("Scrape", url).Return(content, nil)
	mockAI.On("ExtractContent", ctx, content).Return(extraction, nil)
	mockAI.On("SummarizeContent", ctx, content).Return(summary, nil)

	concatenatedSummary := "Overview text Point 1 Point 2"
	expectedItem := db.Item{
		ID:      1,
		UserID:  &userID,
		Title:   title,
		Summary: &concatenatedSummary,
	}

	mockQuerier.On("CreateItem", ctx, mock.Anything).Return(expectedItem, nil)

	item, err := service.ProcessURL(ctx, userID, url)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, title, item.Title)
	mockScraper.AssertExpectations(t)
	mockAI.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestProcessURL_ScrapingError(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	userID := int32(1)
	url := "https://example.com"

	mockScraper.On("Scrape", url).Return("", errors.New("scraping failed"))

	item, err := service.ProcessURL(ctx, userID, url)

	assert.Error(t, err)
	assert.Nil(t, item)
	mockScraper.AssertExpectations(t)
}

func TestGetItemProcessingStatus(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	itemID := int32(1)

	item := &db.Item{ID: itemID, Title: "Test"}
	expectedStatus := &ItemStatus{
		Item:        item,
		IsCompleted: true,
		IsFailed:    false,
	}

	mockJobQueue.On("GetItemStatus", ctx, itemID).Return(expectedStatus, nil)

	status, err := service.GetItemProcessingStatus(ctx, itemID)

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.True(t, status.IsCompleted)
	mockJobQueue.AssertExpectations(t)
}

func TestGetUnreadItemsFromPreviousDay(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()

	expectedItems := []db.Item{
		{ID: 1, Title: "Item 1"},
		{ID: 2, Title: "Item 2"},
	}

	mockQuerier.On("GetUnreadItemsFromPreviousDay", ctx).Return(expectedItems, nil)

	items, err := service.GetUnreadItemsFromPreviousDay(ctx)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	mockQuerier.AssertExpectations(t)
}

func TestConcatenateSummary(t *testing.T) {
	summary := ItemSummary{
		Overview:  "This is an overview.",
		KeyPoints: []string{"Point 1", "Point 2", "Point 3"},
	}

	result := ConcatenateSummary(summary)

	expected := "This is an overview. Point 1 Point 2 Point 3"
	assert.Equal(t, expected, result)
}

func TestGetItemsByProcessingStatus(t *testing.T) {
	mockQuerier := new(test.MockQuerier)
	mockAI := new(MockAIService)
	mockScraper := new(MockScrapingService)
	mockJobQueue := new(MockJobQueueService)

	service := NewItemService(mockQuerier, mockAI, mockScraper, mockJobQueue)

	ctx := context.Background()
	status := "completed"

	expectedItems := []db.Item{
		{ID: 1, Title: "Item 1", ProcessingStatus: &status},
		{ID: 2, Title: "Item 2", ProcessingStatus: &status},
	}

	// The service calls jobQueueService.GetItemsByStatus with dereferenced status
	mockJobQueue.On("GetItemsByStatus", ctx, status).Return(expectedItems, nil)

	items, err := service.GetItemsByProcessingStatus(ctx, &status)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	mockJobQueue.AssertExpectations(t)
}
