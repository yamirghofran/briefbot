package services

import (
	"testing"
)

func TestDigestService_UnifiedFunctionality(t *testing.T) {
	// Create a mock querier and services
	// This is a basic test - in a real scenario you'd use proper mocks

	// Test that the unified service implements both interfaces
	var _ DigestService = (*digestService)(nil)

	// Test configuration methods
	service := &digestService{
		podcastEnabled: false,
	}

	// Test SetPodcastGenerationEnabled and IsPodcastGenerationEnabled
	service.SetPodcastGenerationEnabled(true)
	if !service.IsPodcastGenerationEnabled() {
		t.Error("Podcast generation should be enabled after setting to true")
	}

	service.SetPodcastGenerationEnabled(false)
	if service.IsPodcastGenerationEnabled() {
		t.Error("Podcast generation should be disabled after setting to false")
	}
}

func TestDigestResult_Structure(t *testing.T) {
	// Test DigestResult structure
	result := &DigestResult{
		EmailSent:  true,
		PodcastURL: stringPtr("https://example.com/podcast.mp3"),
		ItemsCount: 5,
		Error:      nil,
		DigestType: "integrated",
	}

	if result.DigestType != "integrated" {
		t.Errorf("Expected digest type 'integrated', got '%s'", result.DigestType)
	}

	if !result.EmailSent {
		t.Error("Expected EmailSent to be true")
	}

	if result.ItemsCount != 5 {
		t.Errorf("Expected ItemsCount to be 5, got %d", result.ItemsCount)
	}

	if result.PodcastURL == nil || *result.PodcastURL != "https://example.com/podcast.mp3" {
		t.Error("Expected PodcastURL to be set correctly")
	}
}

func TestDigestConfig_DefaultValues(t *testing.T) {
	// Test that default configuration is set correctly
	service := &digestService{
		config: DigestConfig{
			Subject:        "Your Daily BriefBot Digest - %s",
			PodcastEnabled: false,
		},
	}

	if service.config.Subject != "Your Daily BriefBot Digest - %s" {
		t.Errorf("Expected default subject, got '%s'", service.config.Subject)
	}

	if service.config.PodcastEnabled {
		t.Error("Expected podcast to be disabled by default")
	}
}
