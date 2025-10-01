package services

import (
	"testing"

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

// Note: Testing ExtractContent, SummarizeContent, and WritePodcast requires mocking
// the OpenAI client, which is complex. These would be better tested with integration tests
// or by creating a mock OpenAI client interface.
