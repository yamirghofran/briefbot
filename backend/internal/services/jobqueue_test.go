package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yamirghofran/briefbot/internal/db"
	"github.com/yamirghofran/briefbot/internal/test"
)

// TestEnqueueItem tests the basic enqueue functionality
func TestEnqueueItem(t *testing.T) {
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	userID := int32(1)
	title := "Test Article"
	url := "https://example.com/test"

	expectedParams := db.CreatePendingItemParams{
		UserID: &userID,
		Title:  title,
		Url:    &url,
	}

	expectedItem := test.NewTestDataBuilder().BuildPendingItem()
	expectedItem.UserID = &userID
	expectedItem.Title = title
	expectedItem.Url = &url

	mockQuerier.On("CreatePendingItem", ctx, expectedParams).Return(*expectedItem, nil)

	item, err := jobQueueService.EnqueueItem(ctx, userID, title, url)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, title, item.Title)
	assert.Equal(t, url, *item.Url)
	assert.Equal(t, userID, *item.UserID)
	assert.Equal(t, "pending", *item.ProcessingStatus)
	mockQuerier.AssertExpectations(t)
}

// TestDequeuePendingItems tests dequeuing pending items
func TestDequeuePendingItems(t *testing.T) {
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	limit := int32(10)
	expectedItems := []db.Item{
		*test.NewTestDataBuilder().BuildPendingItem(),
		*test.NewTestDataBuilder().BuildPendingItem(),
	}
	expectedItems[1].ID = 2 // Make second item different

	mockQuerier.On("GetPendingItems", ctx, limit).Return(expectedItems, nil)

	items, err := jobQueueService.DequeuePendingItems(ctx, limit)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, expectedItems, items)
	mockQuerier.AssertExpectations(t)
}

// TestMarkItemAsProcessing tests marking an item as processing
func TestMarkItemAsProcessing(t *testing.T) {
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	itemID := int32(1)
	userID := int32(1)

	// Mock GetItem to return an item with user ID (needed for SSE notification)
	testItem := test.NewTestDataBuilder().BuildPendingItem()
	testItem.ID = itemID
	testItem.UserID = &userID
	mockQuerier.On("GetItem", ctx, itemID).Return(*testItem, nil)

	mockQuerier.On("UpdateItemAsProcessing", ctx, itemID).Return(nil)

	err := jobQueueService.MarkItemAsProcessing(ctx, itemID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

// TestFailItem tests marking an item as failed
func TestFailItem(t *testing.T) {
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	itemID := int32(1)
	userID := int32(1)
	errorMsg := "Test error message"

	// Mock GetItem to return an item with user ID (needed for SSE notification)
	testItem := test.NewTestDataBuilder().BuildPendingItem()
	testItem.ID = itemID
	testItem.UserID = &userID
	mockQuerier.On("GetItem", ctx, itemID).Return(*testItem, nil)

	mockQuerier.On("UpdateItemProcessingStatus", ctx, mock.MatchedBy(func(params db.UpdateItemProcessingStatusParams) bool {
		return params.ID == itemID &&
			params.ProcessingStatus != nil &&
			*params.ProcessingStatus == "failed" &&
			params.ProcessingError != nil &&
			*params.ProcessingError == errorMsg
	})).Return(nil)

	err := jobQueueService.FailItem(ctx, itemID, errorMsg)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

// TestGetItemStatus tests getting item status
func TestGetItemStatus(t *testing.T) {
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	itemID := int32(1)

	// Test with completed item
	testItem := test.NewTestDataBuilder().BuildCompletedItem()
	mockQuerier.On("GetItem", ctx, itemID).Return(*testItem, nil)

	status, err := jobQueueService.GetItemStatus(ctx, itemID)

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, testItem, status.Item)
	assert.False(t, status.IsProcessing)
	assert.True(t, status.IsCompleted)
	assert.False(t, status.IsFailed)
	assert.Nil(t, status.ProcessingError)
	mockQuerier.AssertExpectations(t)
}

// TestGetItemsByStatus tests getting items by processing status
func TestGetItemsByStatus(t *testing.T) {
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	status := "pending"
	expectedItems := []db.Item{
		*test.NewTestDataBuilder().BuildPendingItem(),
		*test.NewTestDataBuilder().BuildPendingItem(),
	}
	expectedItems[1].ID = 2

	mockQuerier.On("GetItemsByProcessingStatus", ctx, &status).Return(expectedItems, nil)

	items, err := jobQueueService.GetItemsByStatus(ctx, status)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, expectedItems, items)
	mockQuerier.AssertExpectations(t)
}

// TestRetryItem tests retrying a failed item
func TestRetryItem(t *testing.T) {
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	itemID := int32(1)

	mockQuerier.On("UpdateItemProcessingStatus", ctx, mock.MatchedBy(func(params db.UpdateItemProcessingStatusParams) bool {
		return params.ID == itemID &&
			params.ProcessingStatus != nil &&
			*params.ProcessingStatus == "pending" &&
			params.ProcessingError == nil
	})).Return(nil)

	err := jobQueueService.RetryItem(ctx, itemID)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

// TestCompleteItemPreservesIsRead tests the critical bug fix where is_read was being set to NULL
func TestCompleteItemPreservesIsRead(t *testing.T) {
	// Setup
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	// Test data - item with is_read = false (the bug would set this to NULL)
	testItem := test.NewTestDataBuilder().BuildItem()
	testItem.IsRead = boolPtr(false) // Explicitly set to false

	// Expected parameters for UpdateItem call
	expectedUpdateParams := db.UpdateItemParams{
		ID:          testItem.ID,
		Title:       "AI Extracted Title",
		Url:         testItem.Url,
		IsRead:      testItem.IsRead, // This should preserve the original value, not be nil
		TextContent: strPtr("Extracted content"),
		Summary:     strPtr("Extracted summary"),
		Type:        strPtr("article"),
		Tags:        []string{"tag1", "tag2"},
		Platform:    strPtr("web"),
		Authors:     []string{"Author 1"},
	}

	// Mock expectations
	mockQuerier.On("GetItem", ctx, testItem.ID).Return(*testItem, nil)
	mockQuerier.On("UpdateItem", ctx, expectedUpdateParams).Return(nil)

	completedStatus := "completed"
	mockQuerier.On("UpdateItemProcessingStatus", ctx, mock.MatchedBy(func(params db.UpdateItemProcessingStatusParams) bool {
		return params.ID == testItem.ID &&
			params.ProcessingStatus != nil &&
			*params.ProcessingStatus == completedStatus &&
			params.ProcessingError == nil
	})).Return(nil)

	// Execute
	err := jobQueueService.CompleteItem(
		ctx,
		testItem.ID,
		"AI Extracted Title",
		"Extracted content",
		"Extracted summary",
		"article",
		"web",
		[]string{"tag1", "tag2"},
		[]string{"Author 1"},
	)

	// Assert
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

// TestCompleteItemWithNullIsRead tests the case where is_read is initially NULL
func TestCompleteItemWithNullIsRead(t *testing.T) {
	// Setup
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	// Test data - item with is_read = nil (unread)
	testItem := test.NewTestDataBuilder().BuildItem()
	testItem.IsRead = nil // Initially NULL

	// Expected parameters for UpdateItem call
	expectedUpdateParams := db.UpdateItemParams{
		ID:          testItem.ID,
		Title:       "AI Extracted Title",
		Url:         testItem.Url,
		IsRead:      testItem.IsRead, // Should preserve nil, not overwrite with false
		TextContent: strPtr("Extracted content"),
		Summary:     strPtr("Extracted summary"),
		Type:        strPtr("article"),
		Tags:        []string{"tag1", "tag2"},
		Platform:    strPtr("web"),
		Authors:     []string{"Author 1"},
	}

	// Mock expectations
	mockQuerier.On("GetItem", ctx, testItem.ID).Return(*testItem, nil)
	mockQuerier.On("UpdateItem", ctx, expectedUpdateParams).Return(nil)

	completedStatus := "completed"
	mockQuerier.On("UpdateItemProcessingStatus", ctx, mock.MatchedBy(func(params db.UpdateItemProcessingStatusParams) bool {
		return params.ID == testItem.ID &&
			params.ProcessingStatus != nil &&
			*params.ProcessingStatus == completedStatus &&
			params.ProcessingError == nil
	})).Return(nil)

	// Execute
	err := jobQueueService.CompleteItem(
		ctx,
		testItem.ID,
		"AI Extracted Title",
		"Extracted content",
		"Extracted summary",
		"article",
		"web",
		[]string{"tag1", "tag2"},
		[]string{"Author 1"},
	)

	// Assert
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

// TestCompleteItemWithTrueIsRead tests the case where is_read is true
func TestCompleteItemWithTrueIsRead(t *testing.T) {
	// Setup
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	// Test data - item with is_read = true (already read)
	testItem := test.NewTestDataBuilder().BuildItem()
	testItem.IsRead = boolPtr(true) // Already read

	// Expected parameters for UpdateItem call
	expectedUpdateParams := db.UpdateItemParams{
		ID:          testItem.ID,
		Title:       "AI Extracted Title",
		Url:         testItem.Url,
		IsRead:      testItem.IsRead, // Should preserve true, not change to false
		TextContent: strPtr("Extracted content"),
		Summary:     strPtr("Extracted summary"),
		Type:        strPtr("article"),
		Tags:        []string{"tag1", "tag2"},
		Platform:    strPtr("web"),
		Authors:     []string{"Author 1"},
	}

	// Mock expectations
	mockQuerier.On("GetItem", ctx, testItem.ID).Return(*testItem, nil)
	mockQuerier.On("UpdateItem", ctx, expectedUpdateParams).Return(nil)

	completedStatus := "completed"
	mockQuerier.On("UpdateItemProcessingStatus", ctx, mock.MatchedBy(func(params db.UpdateItemProcessingStatusParams) bool {
		return params.ID == testItem.ID &&
			params.ProcessingStatus != nil &&
			*params.ProcessingStatus == completedStatus &&
			params.ProcessingError == nil
	})).Return(nil)

	// Execute
	err := jobQueueService.CompleteItem(
		ctx,
		testItem.ID,
		"AI Extracted Title",
		"Extracted content",
		"Extracted summary",
		"article",
		"web",
		[]string{"tag1", "tag2"},
		[]string{"Author 1"},
	)

	// Assert
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

// TestCompleteItemErrorHandling tests error scenarios
func TestCompleteItemErrorHandling(t *testing.T) {
	t.Run("GetItemError", func(t *testing.T) {
		mockQuerier := &test.MockQuerier{}
		jobQueueService := NewJobQueueService(mockQuerier)
		ctx := context.Background()

		expectedError := assert.AnError
		mockQuerier.On("GetItem", ctx, int32(999)).Return(db.Item{}, expectedError)

		err := jobQueueService.CompleteItem(ctx, 999, "title", "content", "summary", "type", "platform", []string{}, []string{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get item for completion")
		mockQuerier.AssertExpectations(t)
	})

	t.Run("UpdateItemError", func(t *testing.T) {
		mockQuerier := &test.MockQuerier{}
		jobQueueService := NewJobQueueService(mockQuerier)
		ctx := context.Background()

		testItem := test.NewTestDataBuilder().BuildItem()
		expectedError := assert.AnError

		mockQuerier.On("GetItem", ctx, testItem.ID).Return(*testItem, nil)
		mockQuerier.On("UpdateItem", ctx, mock.Anything).Return(expectedError)

		err := jobQueueService.CompleteItem(ctx, testItem.ID, "title", "content", "summary", "type", "platform", []string{}, []string{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update item with processed data")
		mockQuerier.AssertExpectations(t)
	})
}

// TestGetFailedItemsForRetry tests getting failed items eligible for retry
func TestGetFailedItemsForRetry(t *testing.T) {
	mockQuerier := &test.MockQuerier{}
	jobQueueService := NewJobQueueService(mockQuerier)
	ctx := context.Background()

	limit := int32(5)

	// NOTE: Current implementation just returns empty slice
	// This is a placeholder for future implementation
	items, err := jobQueueService.GetFailedItemsForRetry(ctx, limit)

	assert.NoError(t, err)
	assert.NotNil(t, items)
	assert.Len(t, items, 0) // Current implementation returns empty
}

// Helper functions (these would normally be in a test utilities file)
func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }
