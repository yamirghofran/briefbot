package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/invopop/jsonschema"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ItemExtraction struct {
	Title   string   `json:"title" jsonschema_description:"The title for this item."`
	Authors []string `json:"authors" jsonschema_description:"The authors of this item"`
	Tags    []string `json:"tags" jsonschema_description: "Broad tags that match this item"`
	Type    string   `json:"type" jsonschema:"enum=article,enum=github-repo,enum=research-paper"`
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

func extractItemFromContent(itemContent string, schemaParam *openai.ResponseFormatJSONSchemaJSONSchemaParam, oaiClient *openai.Client) (ItemExtraction, error) {
	chatCompletion, err := oaiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are an accountant and your job is to analyze invoices and summarize them into the key metrics and facts."),
			openai.UserMessage("Summarize this invoice"),
			openai.UserMessage(itemContent),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: *schemaParam,
			},
		},
		Model: "openai/gpt-oss-20b",
	})
	if err != nil {
		return ItemExtraction{}, err
	}
	// extract into a well-typed struct
	var invoiceAnalysis ItemExtraction
	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &invoiceAnalysis); err != nil {
		panic(err)
	}
	return invoiceAnalysis, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load .env file.")
	}

	//fmt.Printf("%v", markdownResult)
	//fmt.Printf("%v", filepaths)
	oaiClient := openai.NewClient(
		option.WithBaseURL("https://api.groq.com/openai/v1"),
		option.WithAPIKey(os.Getenv("GROQ_API_KEY")),
	)

	// Support environment variable for custom content file path, fallback to default location
	contentFile := os.Getenv("CONTENT_FILE")
	if contentFile == "" {
		// Get the directory where the source file is located
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("Failed to get current file path")
		}
		dir := filepath.Dir(filename)
		contentFile = filepath.Join(dir, "content.txt")
	}

	// Read content from file
	content, err := os.ReadFile(contentFile)
	if err != nil {
		log.Fatalf("Failed to read content file \"%s\": %v", contentFile, err)
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "item_extraction",
		Description: openai.String("Extraction of the item."),
		Schema:      ItemExtractionSchema,
		Strict:      openai.Bool(true),
	}

	item, err := extractItemFromContent(string(content), &schemaParam, &oaiClient)
	if err != nil {
		log.Fatal("Failed to extract item:", err)
	}

	// Pretty print the JSON output
	jsonBytes, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal JSON:", err)
	}

	fmt.Println(string(jsonBytes))
}
