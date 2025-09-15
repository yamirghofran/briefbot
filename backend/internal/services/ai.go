package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

type AIService interface {
	ExtractContent(ctx context.Context, content string) (ItemExtraction, error)
	SummarizeContent(ctx context.Context, content string) (ItemSummary, error)
}

type aiService struct {
	client openai.Client
}

func NewAIService(oaiClient *openai.Client) (AIService, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	return &aiService{
		client: *oaiClient,
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
	chatCompletion, err := s.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
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
	chatCompletion, err := s.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
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
