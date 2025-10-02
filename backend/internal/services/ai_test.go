package services

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSchema(t *testing.T) {
	schema := GenerateSchema[ItemExtraction]()
	assert.NotNil(t, schema)

	schema2 := GenerateSchema[ItemSummary]()
	assert.NotNil(t, schema2)

	schema3 := GenerateSchema[Podcast]()
	assert.NotNil(t, schema3)
}

func TestSchemaGeneration(t *testing.T) {
	assert.NotNil(t, ItemExtractionSchema)
	assert.NotNil(t, ItemSummarySchema)
	assert.NotNil(t, PodcastSchema)
	assert.NotNil(t, PodcastSectionSchema)
}

func TestNewAIService_MissingAPIKey(t *testing.T) {
	// Temporarily unset GROQ_API_KEY to test error handling
	originalKey := os.Getenv("GROQ_API_KEY")
	os.Unsetenv("GROQ_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("GROQ_API_KEY", originalKey)
		}
	}()

	_, err := NewAIService(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GROQ_API_KEY")
}

func TestNewAIService_Success(t *testing.T) {
	// Set a dummy API key
	originalKey := os.Getenv("GROQ_API_KEY")
	os.Setenv("GROQ_API_KEY", "test-key")
	defer func() {
		if originalKey != "" {
			os.Setenv("GROQ_API_KEY", originalKey)
		} else {
			os.Unsetenv("GROQ_API_KEY")
		}
	}()

	mockClient := &openai.Client{}
	svc, err := NewAIService(mockClient)
	assert.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestExtractContent_ValidatesStructure(t *testing.T) {
	// Set a dummy API key
	originalKey := os.Getenv("GROQ_API_KEY")
	os.Setenv("GROQ_API_KEY", "test-key")
	defer func() {
		if originalKey != "" {
			os.Setenv("GROQ_API_KEY", originalKey)
		} else {
			os.Unsetenv("GROQ_API_KEY")
		}
	}()

	// Create a mock service with a nil client to test structure
	svc := &aiService{
		textClient: openai.Client{},
	}

	// Test that the method accepts proper parameters
	ctx := context.Background()
	content := "Test content for extraction"

	// This will fail to connect but validates the structure exists
	_, err := svc.ExtractContent(ctx, content)
	// We expect an error since we don't have a real client configured
	assert.Error(t, err)
}

func TestSummarizeContent_ValidatesStructure(t *testing.T) {
	// Set a dummy API key
	originalKey := os.Getenv("GROQ_API_KEY")
	os.Setenv("GROQ_API_KEY", "test-key")
	defer func() {
		if originalKey != "" {
			os.Setenv("GROQ_API_KEY", originalKey)
		} else {
			os.Unsetenv("GROQ_API_KEY")
		}
	}()

	// Create a mock service with a nil client to test structure
	svc := &aiService{
		textClient: openai.Client{},
	}

	// Test that the method accepts proper parameters
	ctx := context.Background()
	content := "Test content for summarization"

	// This will fail to connect but validates the structure exists
	_, err := svc.SummarizeContent(ctx, content)
	// We expect an error since we don't have a real client configured
	assert.Error(t, err)
}

func TestWritePodcast_ValidatesStructure(t *testing.T) {
	// Set a dummy API key
	originalKey := os.Getenv("GROQ_API_KEY")
	os.Setenv("GROQ_API_KEY", "test-key")
	defer func() {
		if originalKey != "" {
			os.Setenv("GROQ_API_KEY", originalKey)
		} else {
			os.Unsetenv("GROQ_API_KEY")
		}
	}()

	// Create a mock service with a nil client to test structure
	svc := &aiService{
		textClient: openai.Client{},
	}

	content := "Test content for podcast generation"

	// This will fail to connect but validates the structure exists
	_, err := svc.WritePodcast(content)
	// We expect an error since we don't have a real client configured
	assert.Error(t, err)
}

func TestWritePodcastSection_ValidatesStructure(t *testing.T) {
	// Set a dummy API key
	originalKey := os.Getenv("GROQ_API_KEY")
	os.Setenv("GROQ_API_KEY", "test-key")
	defer func() {
		if originalKey != "" {
			os.Setenv("GROQ_API_KEY", originalKey)
		} else {
			os.Unsetenv("GROQ_API_KEY")
		}
	}()

	// Create a mock service with a nil client to test structure
	svc := &aiService{
		textClient: openai.Client{},
	}

	content := "Test content for podcast section"
	section := "introduction"
	resultChan := make(chan PodcastSectionResult, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	// Call the method
	go svc.WritePodcastSection(content, section, resultChan, &wg)

	// Wait for completion
	wg.Wait()
	close(resultChan)

	// Should receive a result with an error
	result := <-resultChan
	assert.Equal(t, section, result.Section)
	assert.Error(t, result.Error)
}

// Note: Full integration testing of ExtractContent, SummarizeContent, and WritePodcast
// would require a real OpenAI API key and are better suited for integration tests.
// The tests above validate that the methods exist, accept correct parameters, and
// handle errors appropriately.
