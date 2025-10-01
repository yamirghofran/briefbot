package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPodcastConfig(t *testing.T) {
	config := DefaultPodcastConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 1.0, config.DefaultSpeed)
	assert.Equal(t, 3, config.MaxRetries)
	assert.NotEmpty(t, config.TempDir)
	assert.True(t, config.EnableStorage)
	assert.NotNil(t, config.VoiceMapping)
}

func TestPodcastConfig_Structure(t *testing.T) {
	config := PodcastConfig{
		DefaultSpeed:       1.5,
		MaxRetries:         5,
		TempDir:            "/tmp",
		EnableStorage:      true,
		VoiceMapping:       make(map[string]VoiceEnum),
		MaxItemsPerPodcast: 10,
		MaxConcurrentAudio: 5,
	}

	assert.Equal(t, 1.5, config.DefaultSpeed)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, "/tmp", config.TempDir)
	assert.True(t, config.EnableStorage)
	assert.Equal(t, 10, config.MaxItemsPerPodcast)
	assert.Equal(t, int32(5), config.MaxConcurrentAudio)
}

// Note: More comprehensive tests for PodcastService would require mocking:
// - AIService for script generation
// - SpeechService for audio generation
// - R2Service for storage
// - db.Querier for database operations
// These would be integration tests or require extensive mocking infrastructure.
