package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPublicURL(t *testing.T) {
	testCases := []struct {
		name       string
		publicHost string
		accountId  string
		key        string
		expected   string
	}{
		{
			name:       "WithCustomDomain",
			publicHost: "https://cdn.example.com",
			accountId:  "test-account",
			key:        "podcasts/test.mp3",
			expected:   "https://cdn.example.com/podcasts/test.mp3",
		},
		{
			name:       "WithoutCustomDomain",
			publicHost: "",
			accountId:  "test-account",
			key:        "podcasts/test.mp3",
			expected:   "https://pub-test-account.r2.dev/podcasts/test.mp3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &R2Service{
				publicHost: tc.publicHost,
				accountId:  tc.accountId,
			}

			result := service.GetPublicURL(tc.key)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExtractKeyFromURL(t *testing.T) {
	testCases := []struct {
		name       string
		publicHost string
		accountId  string
		url        string
		expected   string
	}{
		{
			name:       "CustomDomainURL",
			publicHost: "https://cdn.example.com",
			accountId:  "test-account",
			url:        "https://cdn.example.com/podcasts/test.mp3",
			expected:   "podcasts/test.mp3",
		},
		{
			name:       "R2PublicURL",
			publicHost: "",
			accountId:  "test-account",
			url:        "https://pub-test-account.r2.dev/podcasts/test.mp3",
			expected:   "podcasts/test.mp3",
		},
		{
			name:       "PlainKey",
			publicHost: "",
			accountId:  "test-account",
			url:        "podcasts/test.mp3",
			expected:   "podcasts/test.mp3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &R2Service{
				publicHost: tc.publicHost,
				accountId:  tc.accountId,
			}

			result := service.ExtractKeyFromURL(tc.url)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNewR2Service_MissingConfig(t *testing.T) {
	testCases := []struct {
		name   string
		config R2Config
	}{
		{
			name: "MissingAccessKeyID",
			config: R2Config{
				SecretAccessKey: "secret",
				AccountID:       "account",
				BucketName:      "bucket",
			},
		},
		{
			name: "MissingSecretAccessKey",
			config: R2Config{
				AccessKeyID: "key",
				AccountID:   "account",
				BucketName:  "bucket",
			},
		},
		{
			name: "MissingAccountID",
			config: R2Config{
				AccessKeyID:     "key",
				SecretAccessKey: "secret",
				BucketName:      "bucket",
			},
		},
		{
			name: "MissingBucketName",
			config: R2Config{
				AccessKeyID:     "key",
				SecretAccessKey: "secret",
				AccountID:       "account",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, err := NewR2Service(tc.config)
			assert.Error(t, err)
			assert.Nil(t, service)
			assert.Contains(t, err.Error(), "missing required R2 configuration fields")
		})
	}
}

func TestNewR2Service_ValidConfig(t *testing.T) {
	config := R2Config{
		AccessKeyID:     "test-key",
		SecretAccessKey: "test-secret",
		AccountID:       "test-account",
		BucketName:      "test-bucket",
		PublicHost:      "https://cdn.example.com",
	}

	service, err := NewR2Service(config)
	require.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "test-bucket", service.bucket)
	assert.Equal(t, "test-account", service.accountId)
	assert.Equal(t, "https://cdn.example.com", service.publicHost)
}
