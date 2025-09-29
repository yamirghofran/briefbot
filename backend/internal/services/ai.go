package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

type AIService interface {
	ExtractContent(ctx context.Context, content string) (ItemExtraction, error)
	SummarizeContent(ctx context.Context, content string) (ItemSummary, error)
	WritePodcast(content string, schemaParam *openai.ResponseFormatJSONSchemaJSONSchemaParam) (Podcast, error)
}

type aiService struct {
	textClient openai.Client
}

func NewAIService(oaiClient *openai.Client) (AIService, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	return &aiService{
		textClient: *oaiClient,
	}, nil
}

type ItemExtraction struct {
	Title    string   `json:"title" jsonschema_description:"The title for this item."`
	Authors  []string `json:"authors" jsonschema_description:"The authors of this item"`
	Tags     []string `json:"tags" jsonschema_description: "Broad tags that match this item"`
	Platform string   `json:"platform" jsonschema_description:"The platform the item is published on." jsonschema:"enum=Youtube,enum=Github,enum=Arxiv,enum=WSJ,enum=Blog,enum=Medium,enum=Substack"`
	Type     string   `json:"type" jsonschema:"enum=article,enum=github-repo,enum=research-paper,enum=podcast,enum=video"`
}

type ItemSummary struct {
	Overview  string   `json:"overview" jsonschema_description:"Brief overview about the item."`
	KeyPoints []string `json:"key_points" jsonschema_description:"A list of key points that succinctly deliver the most important facts from the item."`
}

type Podcast struct {
	Dialogues []Dialogue `json:"dialogues" jsonschema:"required" jsonschema_description:"The dialogues that make up the podcast"`
}

type Dialogue struct {
	Speaker string `json:"speaker" jsonschema:"required,enum=heart,enum=adam" jsonschema_description:"The speaker identifier"`
	Content string `json:"content" jsonschema:"required" jsonschema_description:"The content of what is spoken"`
}

// SectionedPodcast represents a podcast broken into introduction, body, and conclusion
type SectionedPodcast struct {
	Introduction []Dialogue `json:"introduction" jsonschema:"required"`
	Body         []Dialogue `json:"body" jsonschema:"required"`
	Conclusion   []Dialogue `json:"conclusion" jsonschema:"required"`
}

// PodcastSectionResult represents the result of generating a podcast section
type PodcastSectionResult struct {
	Section   string     // "introduction", "body", or "conclusion"
	Dialogues []Dialogue // The generated dialogues for this section
	Error     error
}

func GenerateSchema[T any]() any {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var ItemExtractionSchema = GenerateSchema[ItemExtraction]()
var ItemSummarySchema = GenerateSchema[ItemSummary]()
var PodcastSchema = GenerateSchema[Podcast]()

type Choice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type UsageAccounting struct {
	Include bool `json:"include"`
}

func (s *aiService) ExtractContent(ctx context.Context, content string) (ItemExtraction, error) {
	chatCompletion, err := s.textClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are an accountant and your job is to analyze invoices and summarize them into the key metrics and facts."),
			openai.UserMessage("Summarize this invoice"),
			openai.UserMessage(content),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "item_extraction",
					Description: openai.String("Extraction of the item."),
					Schema:      ItemExtractionSchema,
					Strict:      openai.Bool(true),
				},
			},
		},
		Model: "moonshotai/kimi-k2-instruct-0905",
	})
	if err != nil {
		return ItemExtraction{}, err
	}
	// extract into a well-typed struct
	var itemExtraction ItemExtraction
	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &itemExtraction); err != nil {
		panic(err)
	}
	return itemExtraction, nil
}

func (s *aiService) SummarizeContent(ctx context.Context, content string) (ItemSummary, error) {
	chatCompletion, err := s.textClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are an expert summarizer that has to summarize the provided material."),
			openai.UserMessage("Summarize this content."),
			openai.UserMessage(content),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "item_summary",
					Description: openai.String("Summary of the item content"),
					Schema:      ItemSummarySchema,
					Strict:      openai.Bool(true),
				},
			},
		},
		Model: "moonshotai/kimi-k2-instruct-0905",
	})
	if err != nil {
		return ItemSummary{}, err
	}
	// extract into a well-typed struct
	var itemSummary ItemSummary
	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &itemSummary); err != nil {
		panic(err)
	}
	return itemSummary, nil
}

// WritePodcastSection generates a specific section of the podcast concurrently
func (s *aiService) WritePodcastSection(content string, section string, schemaParam *openai.ResponseFormatJSONSchemaJSONSchemaParam, resultChan chan<- PodcastSectionResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create section-specific prompts
	sectionPrompts := map[string]string{
		"introduction": "Write an engaging introduction for a podcast discussing the given content. The introduction should introduce the topic, set the context, and get listeners interested. Use 2 co-hosts named 'heart' and 'adam'.",
		"body":         "Write the main body discussion for a podcast about the given content. This should be the core content where the hosts discuss the key points, provide insights, and have a natural conversation. Use 2 co-hosts named 'heart' and 'adam'.",
		"conclusion":   "Write a conclusion for a podcast discussing the given content. This should summarize key points, provide final thoughts, and give listeners a sense of closure. Use 2 co-hosts named 'heart' and 'adam'.",
	}

	// Generate section-specific dialogue
	chatCompletion, err := s.textClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sectionPrompts[section] + " CRITICAL: You MUST return valid JSON with exactly this structure: {\"dialogues\": [{\"speaker\":\"heart\",\"content\":\"dialogue text\"},{\"speaker\":\"adam\",\"content\":\"dialogue text\"}]}. Each dialogue MUST have both 'speaker' and 'content' fields. The speaker MUST be either 'heart' or 'adam'. Never use null values. Always use lowercase field names."),
			openai.UserMessage("Create " + section + " dialogue between 2 cohosts (heart and adam) discussing this content:"),
			openai.UserMessage(content),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: *schemaParam,
			},
		},
		Model: "moonshotai/kimi-k2-instruct-0905",
	})

	if err != nil {
		resultChan <- PodcastSectionResult{Section: section, Error: err}
		return
	}

	var sectionPodcast Podcast
	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &sectionPodcast); err != nil {
		resultChan <- PodcastSectionResult{Section: section, Error: fmt.Errorf("failed to unmarshal section JSON: %w", err)}
		return
	}

	resultChan <- PodcastSectionResult{
		Section:   section,
		Dialogues: sectionPodcast.Dialogues,
		Error:     nil,
	}
}

// WritePodcast generates the complete podcast by writing introduction, body, and conclusion concurrently
func (s *aiService) WritePodcast(content string, schemaParam *openai.ResponseFormatJSONSchemaJSONSchemaParam) (Podcast, error) {
	var wg sync.WaitGroup
	resultChan := make(chan PodcastSectionResult, 3) // Buffer for 3 sections

	sections := []string{"introduction", "body", "conclusion"}

	// Launch concurrent section generation
	for _, section := range sections {
		wg.Add(1)
		go s.WritePodcastSection(content, section, schemaParam, resultChan, &wg)
	}

	// Wait for all goroutines in separate goroutine
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results and maintain order
	sectionResults := make(map[string][]Dialogue)
	var errors []error

	for result := range resultChan {
		if result.Error != nil {
			errors = append(errors, fmt.Errorf("%s section error: %w", result.Section, result.Error))
		} else {
			sectionResults[result.Section] = result.Dialogues
		}
	}

	if len(errors) > 0 {
		return Podcast{}, fmt.Errorf("failed to generate podcast sections: %v", errors)
	}

	// Combine sections in proper order: introduction -> body -> conclusion
	var combinedDialogues []Dialogue

	// Add introduction
	if intro, exists := sectionResults["introduction"]; exists {
		combinedDialogues = append(combinedDialogues, intro...)
	}

	// Add body
	if body, exists := sectionResults["body"]; exists {
		combinedDialogues = append(combinedDialogues, body...)
	}

	// Add conclusion
	if conclusion, exists := sectionResults["conclusion"]; exists {
		combinedDialogues = append(combinedDialogues, conclusion...)
	}

	return Podcast{Dialogues: combinedDialogues}, nil
}
