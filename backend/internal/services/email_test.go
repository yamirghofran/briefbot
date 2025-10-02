package services

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yamirghofran/briefbot/internal/db"
)

func TestGenerateIntegratedDigestEmail(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	url1 := "https://example.com/article1"
	url2 := "https://example.com/article2"
	platform1 := "HackerNews"
	platform2 := "Reddit"
	itemType1 := "article"
	itemType2 := "video"
	summary1 := "First article summary"
	summary2 := "Second article summary"
	podcastURL := "https://example.com/podcast.mp3"
	duration := int32(125)
	createdAt1 := time.Date(2024, 1, 14, 10, 30, 0, 0, time.UTC)
	createdAt2 := time.Date(2024, 1, 14, 15, 45, 0, 0, time.UTC)

	items := []db.Item{
		{
			Title:     "First Article",
			Url:       &url1,
			Platform:  &platform1,
			Type:      &itemType1,
			Summary:   &summary1,
			CreatedAt: &createdAt1,
		},
		{
			Title:     "Second Article",
			Url:       &url2,
			Platform:  &platform2,
			Type:      &itemType2,
			Summary:   &summary2,
			CreatedAt: &createdAt2,
		},
	}

	t.Run("with podcast URL and duration", func(t *testing.T) {
		htmlContent, textContent := GenerateIntegratedDigestEmail(items, &podcastURL, &duration, date)

		// Verify HTML content
		if !strings.Contains(htmlContent, "January 15, 2024") {
			t.Errorf("HTML content missing date")
		}
		if !strings.Contains(htmlContent, "First Article") {
			t.Errorf("HTML content missing first article title")
		}
		if !strings.Contains(htmlContent, "Second Article") {
			t.Errorf("HTML content missing second article title")
		}
		if !strings.Contains(htmlContent, podcastURL) {
			t.Errorf("HTML content missing podcast URL")
		}
		if !strings.Contains(htmlContent, "2:05") {
			t.Errorf("HTML content missing formatted duration")
		}
		if !strings.Contains(htmlContent, "🎧 Listen to Today's Digest") {
			t.Errorf("HTML content missing podcast section")
		}

		// Verify text content
		if !strings.Contains(textContent, "January 15, 2024") {
			t.Errorf("Text content missing date")
		}
		if !strings.Contains(textContent, "First Article") {
			t.Errorf("Text content missing first article title")
		}
		if !strings.Contains(textContent, "Second Article") {
			t.Errorf("Text content missing second article title")
		}
		if !strings.Contains(textContent, podcastURL) {
			t.Errorf("Text content missing podcast URL")
		}
		if !strings.Contains(textContent, "2:05") {
			t.Errorf("Text content missing formatted duration")
		}
		if !strings.Contains(textContent, "🎧 Listen to Today's Digest") {
			t.Errorf("Text content missing podcast section")
		}
	})

	t.Run("without podcast URL", func(t *testing.T) {
		htmlContent, textContent := GenerateIntegratedDigestEmail(items, nil, nil, date)

		// Verify HTML doesn't contain podcast section
		if strings.Contains(htmlContent, "🎧 Listen to Today's Digest") {
			t.Errorf("HTML content should not contain podcast section when URL is nil")
		}
		if strings.Contains(htmlContent, `<div class="podcast-section">`) {
			t.Errorf("HTML content should not contain podcast section div when URL is nil")
		}

		// Verify text doesn't contain podcast section
		if strings.Contains(textContent, "🎧 Listen to Today's Digest") {
			t.Errorf("Text content should not contain podcast section when URL is nil")
		}

		// Should still contain items
		if !strings.Contains(htmlContent, "First Article") {
			t.Errorf("HTML content missing first article")
		}
		if !strings.Contains(textContent, "First Article") {
			t.Errorf("Text content missing first article")
		}
	})

	t.Run("with empty podcast URL", func(t *testing.T) {
		emptyURL := ""
		htmlContent, textContent := GenerateIntegratedDigestEmail(items, &emptyURL, nil, date)

		// Verify HTML doesn't contain podcast section
		if strings.Contains(htmlContent, "🎧 Listen to Today's Digest") {
			t.Errorf("HTML content should not contain podcast section when URL is empty")
		}

		// Verify text doesn't contain podcast section
		if strings.Contains(textContent, "🎧 Listen to Today's Digest") {
			t.Errorf("Text content should not contain podcast section when URL is empty")
		}
	})

	t.Run("with podcast URL but no duration", func(t *testing.T) {
		htmlContent, textContent := GenerateIntegratedDigestEmail(items, &podcastURL, nil, date)

		// Should contain podcast section without duration
		if !strings.Contains(htmlContent, "🎧 Listen to Today's Digest") {
			t.Errorf("HTML content missing podcast section")
		}
		if strings.Contains(htmlContent, "2:05") {
			t.Errorf("HTML content should not contain duration when not provided")
		}

		if !strings.Contains(textContent, "🎧 Listen to Today's Digest") {
			t.Errorf("Text content missing podcast section")
		}
		if strings.Contains(textContent, "2:05") {
			t.Errorf("Text content should not contain duration when not provided")
		}
	})

	t.Run("with empty items", func(t *testing.T) {
		htmlContent, textContent := GenerateIntegratedDigestEmail([]db.Item{}, &podcastURL, &duration, date)

		// Should still contain date and podcast section
		if !strings.Contains(htmlContent, "January 15, 2024") {
			t.Errorf("HTML content missing date")
		}
		if !strings.Contains(htmlContent, podcastURL) {
			t.Errorf("HTML content missing podcast URL")
		}

		if !strings.Contains(textContent, "January 15, 2024") {
			t.Errorf("Text content missing date")
		}
		if !strings.Contains(textContent, podcastURL) {
			t.Errorf("Text content missing podcast URL")
		}
	})
}

func TestGenerateIntegratedDigestHTML(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	url := "https://example.com/article"
	platform := "HackerNews"
	itemType := "article"
	summary := "Article summary"
	podcastURL := "https://example.com/podcast.mp3"
	duration := int32(125)
	createdAt := time.Date(2024, 1, 14, 10, 30, 0, 0, time.UTC)

	items := []db.Item{
		{
			Title:     "Test Article",
			Url:       &url,
			Platform:  &platform,
			Type:      &itemType,
			Summary:   &summary,
			CreatedAt: &createdAt,
		},
	}

	t.Run("HTML structure", func(t *testing.T) {
		html := generateIntegratedDigestHTML(items, &podcastURL, &duration, date)

		// Verify HTML document structure
		if !strings.Contains(html, "<!DOCTYPE html>") {
			t.Errorf("HTML missing DOCTYPE declaration")
		}
		if !strings.Contains(html, "<html>") || !strings.Contains(html, "</html>") {
			t.Errorf("HTML missing html tags")
		}
		if !strings.Contains(html, "<head>") || !strings.Contains(html, "</head>") {
			t.Errorf("HTML missing head tags")
		}
		if !strings.Contains(html, "<body>") || !strings.Contains(html, "</body>") {
			t.Errorf("HTML missing body tags")
		}

		// Verify CSS styles are present
		if !strings.Contains(html, "<style>") {
			t.Errorf("HTML missing CSS styles")
		}
		if !strings.Contains(html, ".podcast-section") {
			t.Errorf("HTML missing podcast-section style")
		}
	})

	t.Run("includes all item fields", func(t *testing.T) {
		html := generateIntegratedDigestHTML(items, nil, nil, date)

		if !strings.Contains(html, "Test Article") {
			t.Errorf("HTML missing article title")
		}
		if !strings.Contains(html, url) {
			t.Errorf("HTML missing article URL")
		}
		if !strings.Contains(html, platform) {
			t.Errorf("HTML missing platform")
		}
		if !strings.Contains(html, itemType) {
			t.Errorf("HTML missing item type")
		}
		if !strings.Contains(html, summary) {
			t.Errorf("HTML missing summary")
		}
	})

	t.Run("duration formatting", func(t *testing.T) {
		// Test various durations
		tests := []struct {
			duration int32
			expected string
		}{
			{65, "1:05"},
			{125, "2:05"},
			{3661, "61:01"},
			{59, "0:59"},
		}

		for _, tt := range tests {
			html := generateIntegratedDigestHTML(items, &podcastURL, &tt.duration, date)
			if !strings.Contains(html, tt.expected) {
				t.Errorf("For duration %d, expected %s in HTML, but not found", tt.duration, tt.expected)
			}
		}
	})

	t.Run("items without optional fields", func(t *testing.T) {
		minimalCreatedAt := time.Date(2024, 1, 14, 10, 30, 0, 0, time.UTC)
		minimalItems := []db.Item{
			{
				Title:     "Minimal Item",
				Url:       &url,
				CreatedAt: &minimalCreatedAt,
			},
		}

		html := generateIntegratedDigestHTML(minimalItems, nil, nil, date)

		if !strings.Contains(html, "Minimal Item") {
			t.Errorf("HTML missing minimal item title")
		}
		if !strings.Contains(html, url) {
			t.Errorf("HTML missing minimal item URL")
		}
	})
}

func TestGenerateIntegratedDigestText(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	url := "https://example.com/article"
	platform := "HackerNews"
	itemType := "article"
	summary := "Article summary"
	podcastURL := "https://example.com/podcast.mp3"
	duration := int32(125)
	createdAt := time.Date(2024, 1, 14, 10, 30, 0, 0, time.UTC)

	items := []db.Item{
		{
			Title:     "Test Article",
			Url:       &url,
			Platform:  &platform,
			Type:      &itemType,
			Summary:   &summary,
			CreatedAt: &createdAt,
		},
	}

	t.Run("text structure with podcast", func(t *testing.T) {
		text := generateIntegratedDigestText(items, &podcastURL, &duration, date)

		// Verify basic structure
		if !strings.HasPrefix(text, "Daily Digest - January 15, 2024\n") {
			t.Errorf("Text missing proper header")
		}
		if !strings.Contains(text, "Your unread items from yesterday") {
			t.Errorf("Text missing subheader")
		}
		if !strings.Contains(text, "🎧 Listen to Today's Digest") {
			t.Errorf("Text missing podcast section")
		}
		if !strings.Contains(text, "Sent by BriefBot") {
			t.Errorf("Text missing footer")
		}
	})

	t.Run("includes all item fields", func(t *testing.T) {
		text := generateIntegratedDigestText(items, nil, nil, date)

		if !strings.Contains(text, "Test Article") {
			t.Errorf("Text missing article title")
		}
		if !strings.Contains(text, url) {
			t.Errorf("Text missing article URL")
		}
		if !strings.Contains(text, "Platform: "+platform) {
			t.Errorf("Text missing platform")
		}
		if !strings.Contains(text, "Type: "+itemType) {
			t.Errorf("Text missing item type")
		}
		if !strings.Contains(text, "Summary: "+summary) {
			t.Errorf("Text missing summary")
		}
	})

	t.Run("numbered list format", func(t *testing.T) {
		multipleCreatedAt := time.Date(2024, 1, 14, 10, 30, 0, 0, time.UTC)
		multipleItems := []db.Item{
			{
				Title:     "First Article",
				Url:       &url,
				CreatedAt: &multipleCreatedAt,
			},
			{
				Title:     "Second Article",
				Url:       &url,
				CreatedAt: &multipleCreatedAt,
			},
		}

		text := generateIntegratedDigestText(multipleItems, nil, nil, date)

		if !strings.Contains(text, "1. First Article") {
			t.Errorf("Text missing first numbered item")
		}
		if !strings.Contains(text, "2. Second Article") {
			t.Errorf("Text missing second numbered item")
		}
	})

	t.Run("duration formatting", func(t *testing.T) {
		tests := []struct {
			duration int32
			expected string
		}{
			{65, "1:05"},
			{125, "2:05"},
			{3661, "61:01"},
			{59, "0:59"},
		}

		for _, tt := range tests {
			text := generateIntegratedDigestText(items, &podcastURL, &tt.duration, date)
			if !strings.Contains(text, tt.expected) {
				t.Errorf("For duration %d, expected %s in text, but not found", tt.duration, tt.expected)
			}
		}
	})

	t.Run("separator lines with and without podcast", func(t *testing.T) {
		textWithPodcast := generateIntegratedDigestText(items, &podcastURL, &duration, date)
		textWithoutPodcast := generateIntegratedDigestText(items, nil, nil, date)

		// With podcast should have different separator after podcast section
		if !strings.Contains(textWithPodcast, strings.Repeat("-", 50)) {
			t.Errorf("Text with podcast missing dash separator")
		}

		// Without podcast should have equals separator
		if !strings.Contains(textWithoutPodcast, strings.Repeat("=", 50)) {
			t.Errorf("Text without podcast missing equals separator")
		}
	})

	t.Run("items without optional fields", func(t *testing.T) {
		minimalCreatedAt := time.Date(2024, 1, 14, 10, 30, 0, 0, time.UTC)
		minimalItems := []db.Item{
			{
				Title:     "Minimal Item",
				Url:       &url,
				CreatedAt: &minimalCreatedAt,
			},
		}

		text := generateIntegratedDigestText(minimalItems, nil, nil, date)

		if !strings.Contains(text, "Minimal Item") {
			t.Errorf("Text missing minimal item title")
		}
		if !strings.Contains(text, url) {
			t.Errorf("Text missing minimal item URL")
		}
		// Should still have the "Added:" timestamp
		if !strings.Contains(text, "Added:") {
			t.Errorf("Text missing timestamp metadata")
		}
	})
}

func TestNewEmailService_MissingEnvVars(t *testing.T) {
	// Save original env vars
	originalAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	originalSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	originalRegion := os.Getenv("AWS_REGION")
	originalFromEmail := os.Getenv("SES_FROM_EMAIL")
	originalFromName := os.Getenv("SES_FROM_NAME")
	originalReplyTo := os.Getenv("SES_REPLY_TO_EMAIL")

	// Unset all required env vars
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("SES_FROM_EMAIL")

	defer func() {
		// Restore original env vars
		if originalAccessKey != "" {
			os.Setenv("AWS_ACCESS_KEY_ID", originalAccessKey)
		}
		if originalSecretKey != "" {
			os.Setenv("AWS_SECRET_ACCESS_KEY", originalSecretKey)
		}
		if originalRegion != "" {
			os.Setenv("AWS_REGION", originalRegion)
		}
		if originalFromEmail != "" {
			os.Setenv("SES_FROM_EMAIL", originalFromEmail)
		}
		if originalFromName != "" {
			os.Setenv("SES_FROM_NAME", originalFromName)
		}
		if originalReplyTo != "" {
			os.Setenv("SES_REPLY_TO_EMAIL", originalReplyTo)
		}
	}()

	_, err := NewEmailService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required AWS SES environment variables")
}

func TestNewEmailService_WithDefaults(t *testing.T) {
	// Save original env vars
	originalAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	originalSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	originalRegion := os.Getenv("AWS_REGION")
	originalFromEmail := os.Getenv("SES_FROM_EMAIL")
	originalFromName := os.Getenv("SES_FROM_NAME")
	originalReplyTo := os.Getenv("SES_REPLY_TO_EMAIL")

	// Set required vars without optional ones
	os.Setenv("AWS_ACCESS_KEY_ID", "test-key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("SES_FROM_EMAIL", "test@example.com")
	os.Unsetenv("SES_FROM_NAME")
	os.Unsetenv("SES_REPLY_TO_EMAIL")

	defer func() {
		// Restore original env vars
		if originalAccessKey != "" {
			os.Setenv("AWS_ACCESS_KEY_ID", originalAccessKey)
		} else {
			os.Unsetenv("AWS_ACCESS_KEY_ID")
		}
		if originalSecretKey != "" {
			os.Setenv("AWS_SECRET_ACCESS_KEY", originalSecretKey)
		} else {
			os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		}
		if originalRegion != "" {
			os.Setenv("AWS_REGION", originalRegion)
		} else {
			os.Unsetenv("AWS_REGION")
		}
		if originalFromEmail != "" {
			os.Setenv("SES_FROM_EMAIL", originalFromEmail)
		} else {
			os.Unsetenv("SES_FROM_EMAIL")
		}
		if originalFromName != "" {
			os.Setenv("SES_FROM_NAME", originalFromName)
		}
		if originalReplyTo != "" {
			os.Setenv("SES_REPLY_TO_EMAIL", originalReplyTo)
		}
	}()

	svc, err := NewEmailService()
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	// Cast to concrete type to check defaults
	emailSvc, ok := svc.(*emailService)
	assert.True(t, ok)
	assert.Equal(t, "BriefBot", emailSvc.fromName)
	assert.Equal(t, "test@example.com", emailSvc.replyToEmail)
}
