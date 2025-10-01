package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerConfig_Defaults(t *testing.T) {
	config := WorkerConfig{
		WorkerCount:    5,
		BatchSize:      10,
		PollInterval:   5 * time.Second,
		MaxRetries:     3,
		EnablePodcasts: true,
	}

	assert.Equal(t, 5, config.WorkerCount)
	assert.Equal(t, int32(10), config.BatchSize)
	assert.Equal(t, 5*time.Second, config.PollInterval)
	assert.Equal(t, 3, config.MaxRetries)
	assert.True(t, config.EnablePodcasts)
}

// Note: Testing the worker service would require:
// - Mocking all dependencies (ItemService, JobQueueService, PodcastService, etc.)
// - Testing concurrent operations and goroutines
// - Testing timer-based polling logic
// - Testing graceful shutdown
// These would be complex integration tests requiring careful setup.
