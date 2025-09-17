package services

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestScrapingService tests the basic scraping functionality
func TestScrapingService(t *testing.T) {
	scraper := NewScraper()

	// Test with a simple HTML content
	t.Run("BasicScraping", func(t *testing.T) {
		// This test would normally require a mock server or test HTML
		// For now, we'll test the service creation and basic structure
		assert.NotNil(t, scraper)
	})
}

// TestScrapingServiceScrape tests the actual scraping functionality
func TestScrapingServiceScrape(t *testing.T) {
	t.Skip("Skipping actual scraping test - requires network access")

	scraper := NewScraper()

	// Test with a known good URL
	testCases := []struct {
		name     string
		url      string
		contains []string
	}{
		{
			name: "ScrapeExampleCom",
			url:  "https://example.com",
			contains: []string{
				"Example Domain",
				"example",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := scraper.Scrape(tc.url)

			require.NoError(t, err)
			assert.NotEmpty(t, content)

			// Check that expected content is present
			for _, expected := range tc.contains {
				assert.True(t, strings.Contains(content, expected),
					"Content should contain '%s'", expected)
			}
		})
	}
}

// TestScrapingServiceErrorHandling tests error scenarios
func TestScrapingServiceErrorHandling(t *testing.T) {
	t.Skip("Skipping error handling test - requires network access")

	scraper := NewScraper()

	testCases := []struct {
		name        string
		url         string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "InvalidURL",
			url:         "not-a-valid-url",
			expectError: true,
		},
		{
			name:        "NonExistentDomain",
			url:         "https://this-domain-does-not-exist-12345.com",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := scraper.Scrape(tc.url)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, content)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, content)
			}
		})
	}
}
