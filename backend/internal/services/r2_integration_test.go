//go:build integration
// +build integration

package services

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestR2Integration_UploadAndDelete tests the complete upload and delete flow
func TestR2Integration_UploadAndDelete(t *testing.T) {
	// Load R2 configuration from environment
	config := R2Config{
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		PublicHost:      os.Getenv("R2_PUBLIC_HOST"),
	}

	// Skip test if credentials are not configured
	if config.AccessKeyID == "" || config.SecretAccessKey == "" {
		t.Skip("Skipping R2 integration test: R2_ACCESS_KEY_ID or R2_SECRET_ACCESS_KEY not set")
	}

	service, err := NewR2Service(config)
	require.NoError(t, err, "Failed to create R2 service")

	ctx := context.Background()

	// Test data
	key := "integration-tests/" + uuid.New().String() + ".txt"
	data := []byte("This is test data for integration testing")
	contentType := "text/plain"

	// Test 1: Upload file
	t.Run("UploadFile", func(t *testing.T) {
		publicURL, err := service.UploadFile(ctx, key, data, contentType)
		assert.NoError(t, err)
		assert.NotEmpty(t, publicURL)
		assert.Contains(t, publicURL, key)
		t.Logf("Uploaded file to: %s", publicURL)
	})

	// Test 2: Delete file
	t.Run("DeleteFile", func(t *testing.T) {
		err := service.DeleteFile(ctx, key)
		assert.NoError(t, err)
		t.Logf("Deleted file: %s", key)
	})
}

// TestR2Integration_UploadWithMetadata tests uploading with custom metadata
func TestR2Integration_UploadWithMetadata(t *testing.T) {
	config := R2Config{
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		PublicHost:      os.Getenv("R2_PUBLIC_HOST"),
	}

	if config.AccessKeyID == "" || config.SecretAccessKey == "" {
		t.Skip("Skipping R2 integration test: R2_ACCESS_KEY_ID or R2_SECRET_ACCESS_KEY not set")
	}

	service, err := NewR2Service(config)
	require.NoError(t, err)

	ctx := context.Background()
	key := "integration-tests/" + uuid.New().String() + "-metadata.txt"
	data := []byte("Test data with metadata")
	contentType := "text/plain"
	metadata := map[string]string{
		"test-id":   uuid.New().String(),
		"source":    "integration-test",
		"timestamp": "2024-01-01",
	}

	// Upload with metadata
	publicURL, err := service.UploadFileWithMetadata(ctx, key, data, contentType, metadata)
	assert.NoError(t, err)
	assert.NotEmpty(t, publicURL)
	t.Logf("Uploaded file with metadata to: %s", publicURL)

	// Cleanup
	defer func() {
		err := service.DeleteFile(ctx, key)
		if err != nil {
			t.Logf("Warning: Failed to cleanup file %s: %v", key, err)
		}
	}()
}

// TestR2Integration_DeleteMultipleFiles tests batch deletion
func TestR2Integration_DeleteMultipleFiles(t *testing.T) {
	config := R2Config{
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		PublicHost:      os.Getenv("R2_PUBLIC_HOST"),
	}

	if config.AccessKeyID == "" || config.SecretAccessKey == "" {
		t.Skip("Skipping R2 integration test: R2_ACCESS_KEY_ID or R2_SECRET_ACCESS_KEY not set")
	}

	service, err := NewR2Service(config)
	require.NoError(t, err)

	ctx := context.Background()

	// Create multiple test files
	keys := make([]string, 3)
	for i := range keys {
		keys[i] = "integration-tests/batch-" + uuid.New().String() + ".txt"
		_, err := service.UploadFile(ctx, keys[i], []byte("test data"), "text/plain")
		require.NoError(t, err)
	}

	t.Logf("Uploaded %d files for batch deletion test", len(keys))

	// Delete all files at once
	err = service.DeleteFiles(ctx, keys)
	assert.NoError(t, err)
	t.Logf("Batch deleted %d files", len(keys))
}

// TestR2Integration_GenerateUploadURL tests presigned URL generation
func TestR2Integration_GenerateUploadURL(t *testing.T) {
	config := R2Config{
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		PublicHost:      os.Getenv("R2_PUBLIC_HOST"),
	}

	if config.AccessKeyID == "" || config.SecretAccessKey == "" {
		t.Skip("Skipping R2 integration test: R2_ACCESS_KEY_ID or R2_SECRET_ACCESS_KEY not set")
	}

	service, err := NewR2Service(config)
	require.NoError(t, err)

	ctx := context.Background()

	// Generate upload URL
	response, err := service.GenerateUploadURL(ctx, "image/png", "integration-tests")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.UploadURL)
	assert.NotEmpty(t, response.Key)
	assert.NotEmpty(t, response.PublicURL)
	assert.Contains(t, response.Key, "integration-tests/")
	assert.Contains(t, response.Key, ".png")

	t.Logf("Generated upload URL for key: %s", response.Key)
	t.Logf("Upload URL: %s", response.UploadURL)
	t.Logf("Public URL: %s", response.PublicURL)
}

// TestR2Integration_GenerateUploadURLForKey tests presigned URL for specific key
func TestR2Integration_GenerateUploadURLForKey(t *testing.T) {
	config := R2Config{
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		PublicHost:      os.Getenv("R2_PUBLIC_HOST"),
	}

	if config.AccessKeyID == "" || config.SecretAccessKey == "" {
		t.Skip("Skipping R2 integration test: R2_ACCESS_KEY_ID or R2_SECRET_ACCESS_KEY not set")
	}

	service, err := NewR2Service(config)
	require.NoError(t, err)

	ctx := context.Background()
	customKey := "integration-tests/custom-" + uuid.New().String() + ".mp3"

	// Generate upload URL for specific key
	response, err := service.GenerateUploadURLForKey(ctx, customKey, "audio/mpeg")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, customKey, response.Key)
	assert.NotEmpty(t, response.UploadURL)
	assert.NotEmpty(t, response.PublicURL)

	t.Logf("Generated upload URL for custom key: %s", customKey)
}
