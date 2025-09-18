package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yamirghofran/briefbot/internal/db"
)

func TestGenerateDailyDigestEmail(t *testing.T) {
	items := []db.Item{
		{
			ID:        1,
			Title:     "Test Article",
			Url:       stringPtr("https://example.com/article"),
			Summary:   stringPtr("This is a test summary of the article"),
			Platform:  stringPtr("Medium"),
			Type:      stringPtr("article"),
			CreatedAt: timePtr(time.Now().Add(-12 * time.Hour)),
		},
	}

	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	html, text := GenerateDailyDigestEmail(items, testDate)

	// Basic checks
	assert.Contains(t, html, "Test Article")
	assert.Contains(t, html, "https://example.com/article")
	assert.Contains(t, html, "This is a test summary")
	assert.Contains(t, html, "January 15, 2024")

	assert.Contains(t, text, "Test Article")
	assert.Contains(t, text, "https://example.com/article")
	assert.Contains(t, text, "This is a test summary")
	assert.Contains(t, text, "January 15, 2024")
}

func TestGenerateDailyDigestEmail_MultipleItems(t *testing.T) {
	items := []db.Item{
		{
			ID:        1,
			Title:     "First Article",
			Url:       stringPtr("https://example.com/article1"),
			Summary:   stringPtr("Summary of first article"),
			Platform:  stringPtr("Medium"),
			Type:      stringPtr("article"),
			CreatedAt: timePtr(time.Now().Add(-12 * time.Hour)),
		},
		{
			ID:        2,
			Title:     "Second Article",
			Url:       stringPtr("https://example.com/article2"),
			Summary:   stringPtr("Summary of second article"),
			Platform:  stringPtr("Blog"),
			Type:      stringPtr("article"),
			CreatedAt: timePtr(time.Now().Add(-6 * time.Hour)),
		},
	}

	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	html, text := GenerateDailyDigestEmail(items, testDate)

	// Check both items are included
	assert.Contains(t, html, "First Article")
	assert.Contains(t, html, "Second Article")
	assert.Contains(t, html, "https://example.com/article1")
	assert.Contains(t, html, "https://example.com/article2")

	assert.Contains(t, text, "First Article")
	assert.Contains(t, text, "Second Article")
}

func TestGenerateDailyDigestEmail_NoSummary(t *testing.T) {
	items := []db.Item{
		{
			ID:        1,
			Title:     "Article Without Summary",
			Url:       stringPtr("https://example.com/article"),
			Summary:   nil, // No summary
			Platform:  stringPtr("Medium"),
			Type:      stringPtr("article"),
			CreatedAt: timePtr(time.Now().Add(-12 * time.Hour)),
		},
	}

	testDate := time.Now()

	html, text := GenerateDailyDigestEmail(items, testDate)

	// Should still work without summary
	assert.Contains(t, html, "Article Without Summary")
	assert.Contains(t, html, "https://example.com/article")

	assert.Contains(t, text, "Article Without Summary")
	assert.Contains(t, text, "https://example.com/article")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
