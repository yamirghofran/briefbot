package test

import (
	"time"

	"github.com/yamirghofran/briefbot/internal/db"
)

// TestDataBuilder provides convenient methods to create test data
type TestDataBuilder struct{}

// NewTestDataBuilder creates a new test data builder
func NewTestDataBuilder() *TestDataBuilder {
	return &TestDataBuilder{}
}

// BuildUser creates a test user
func (b *TestDataBuilder) BuildUser() *db.User {
	now := time.Now()
	return &db.User{
		ID:           1,
		Name:         strPtr("Test User"),
		Email:        strPtr("test@example.com"),
		AuthProvider: strPtr("google"),
		OauthID:      strPtr("oauth123"),
		PasswordHash: nil,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
}

// BuildItem creates a test item
func (b *TestDataBuilder) BuildItem() *db.Item {
	now := time.Now()
	return &db.Item{
		ID:               1,
		UserID:           int32Ptr(1),
		Url:              strPtr("https://example.com/article"),
		IsRead:           boolPtr(false),
		TextContent:      strPtr("This is test content"),
		Summary:          strPtr("This is a test summary"),
		Type:             strPtr("article"),
		Tags:             []string{"test", "example"},
		Platform:         strPtr("web"),
		Authors:          []string{"Test Author"},
		CreatedAt:        &now,
		ModifiedAt:       &now,
		Title:            "Test Article",
		ProcessingStatus: strPtr("pending"),
		ProcessingError:  nil,
	}
}

// BuildPendingItem creates a test pending item
func (b *TestDataBuilder) BuildPendingItem() *db.Item {
	item := b.BuildItem()
	item.ProcessingStatus = strPtr("pending")
	item.TextContent = nil
	item.Summary = nil
	item.Type = nil
	item.Tags = nil
	item.Platform = nil
	item.Authors = nil
	return item
}

// BuildProcessingItem creates a test processing item
func (b *TestDataBuilder) BuildProcessingItem() *db.Item {
	item := b.BuildItem()
	item.ProcessingStatus = strPtr("processing")
	return item
}

// BuildCompletedItem creates a test completed item
func (b *TestDataBuilder) BuildCompletedItem() *db.Item {
	item := b.BuildItem()
	item.ProcessingStatus = strPtr("completed")
	return item
}

// BuildFailedItem creates a test failed item
func (b *TestDataBuilder) BuildFailedItem() *db.Item {
	item := b.BuildItem()
	item.ProcessingStatus = strPtr("failed")
	item.ProcessingError = strPtr("Test error message")
	return item
}

// Helper functions
func strPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32 { return &i }
func boolPtr(b bool) *bool    { return &b }
