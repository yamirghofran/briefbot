package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yamirghofran/briefbot/internal/db"
)

// PodcastStatus represents the status of a podcast
type PodcastStatus string

const (
	PodcastStatusPending    PodcastStatus = "pending"
	PodcastStatusWriting    PodcastStatus = "writing"
	PodcastStatusGenerating PodcastStatus = "generating"
	PodcastStatusCompleted  PodcastStatus = "completed"
	PodcastStatusFailed     PodcastStatus = "failed"
)

// PodcastService handles podcast creation and management
type PodcastService interface {
	// Main podcast creation methods
	CreatePodcastFromItems(ctx context.Context, userID int32, title string, description string, itemIDs []int32) (*db.Podcast, error)
	CreatePodcastFromSingleItem(ctx context.Context, userID int32, itemID int32) (*db.Podcast, error)

	// Podcast generation workflow
	GeneratePodcastScript(ctx context.Context, podcastID int32) error
	GeneratePodcastAudio(ctx context.Context, podcastID int32) error
	ProcessPodcast(ctx context.Context, podcastID int32) error

	// CRUD operations
	GetPodcast(ctx context.Context, podcastID int32) (*db.Podcast, error)
	GetPodcastsByUser(ctx context.Context, userID int32) ([]db.Podcast, error)
	GetPodcastsByStatus(ctx context.Context, status PodcastStatus) ([]db.Podcast, error)
	UpdatePodcast(ctx context.Context, podcastID int32, title string, description string) error
	DeletePodcast(ctx context.Context, podcastID int32) error

	// Item management
	AddItemToPodcast(ctx context.Context, podcastID int32, itemID int32, order int) error
	RemoveItemFromPodcast(ctx context.Context, podcastID int32, itemID int32) error
	GetPodcastItems(ctx context.Context, podcastID int32) ([]db.GetPodcastItemsRow, error)

	// Status management
	UpdatePodcastStatus(ctx context.Context, podcastID int32, status PodcastStatus) error
	GetPendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error)
	GetProcessingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error)

	// Audio management
	GetPodcastAudio(ctx context.Context, podcastID int32) ([]byte, error)
	HasPodcastAudio(ctx context.Context, podcastID int32) (bool, error)
	GeneratePodcastUploadURL(ctx context.Context, podcastID int32) (*UploadURLResponse, error)
}

// podcastService implements PodcastService
type podcastService struct {
	querier       db.Querier
	aiService     AIService
	speechService SpeechService
	r2Service     *R2Service
	config        PodcastConfig
}

// PodcastConfig holds configuration for podcast service
type PodcastConfig struct {
	DefaultSpeed       float64
	MaxRetries         int
	TempDir            string
	EnableStorage      bool
	VoiceMapping       map[string]VoiceEnum
	MaxItemsPerPodcast int
}

// DefaultPodcastConfig returns default configuration
func DefaultPodcastConfig() PodcastConfig {
	return PodcastConfig{
		DefaultSpeed:       1.0,
		MaxRetries:         3,
		TempDir:            "/tmp/podcasts",
		EnableStorage:      true,
		MaxItemsPerPodcast: 10,
		VoiceMapping: map[string]VoiceEnum{
			"heart": VoiceAfHeart,
			"adam":  VoiceAmAdam,
		},
	}
}

// NewPodcastService creates a new podcast service
func NewPodcastService(querier db.Querier, aiService AIService, speechService SpeechService, r2Service *R2Service, config PodcastConfig) PodcastService {
	return &podcastService{
		querier:       querier,
		aiService:     aiService,
		speechService: speechService,
		r2Service:     r2Service,
		config:        config,
	}
}

// CreatePodcastFromItems creates a podcast from multiple items
func (s *podcastService) CreatePodcastFromItems(ctx context.Context, userID int32, title string, description string, itemIDs []int32) (*db.Podcast, error) {
	if len(itemIDs) == 0 {
		return nil, fmt.Errorf("no items provided for podcast creation")
	}

	if len(itemIDs) > s.config.MaxItemsPerPodcast {
		return nil, fmt.Errorf("too many items: maximum %d items per podcast", s.config.MaxItemsPerPodcast)
	}

	// Create the podcast
	params := db.CreatePodcastParams{
		UserID:      &userID,
		Title:       title,
		Description: &description,
		Status:      string(PodcastStatusPending),
	}

	podcast, err := s.querier.CreatePodcast(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create podcast: %w", err)
	}

	// Add items to podcast in order
	for i, itemID := range itemIDs {
		addParams := db.AddItemToPodcastParams{
			PodcastID: &podcast.ID,
			ItemID:    &itemID,
			ItemOrder: int32(i),
		}

		if _, err := s.querier.AddItemToPodcast(ctx, addParams); err != nil {
			// Clean up on error
			s.querier.DeletePodcast(ctx, podcast.ID)
			return nil, fmt.Errorf("failed to add item %d to podcast: %w", itemID, err)
		}
	}

	return &podcast, nil
}

// CreatePodcastFromSingleItem creates a podcast from a single item with auto-generated title
func (s *podcastService) CreatePodcastFromSingleItem(ctx context.Context, userID int32, itemID int32) (*db.Podcast, error) {
	// Get the item to create a meaningful title
	item, err := s.querier.GetItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	title := fmt.Sprintf("Podcast: %s", item.Title)
	description := "A podcast discussion about: " + item.Title
	if item.Summary != nil {
		description = *item.Summary
	}

	return s.CreatePodcastFromItems(ctx, userID, title, description, []int32{itemID})
}

// GeneratePodcastScript generates the podcast script using AI service
func (s *podcastService) GeneratePodcastScript(ctx context.Context, podcastID int32) error {
	// Get podcast items
	items, err := s.querier.GetPodcastItems(ctx, &podcastID)
	if err != nil {
		return fmt.Errorf("failed to get podcast items: %w", err)
	}

	if len(items) == 0 {
		return fmt.Errorf("no items found for podcast")
	}

	// Update status to writing
	if err := s.UpdatePodcastStatus(ctx, podcastID, PodcastStatusWriting); err != nil {
		return fmt.Errorf("failed to update podcast status: %w", err)
	}

	// Convert GetPodcastItemsRow to Items for content building
	content := s.buildPodcastContentFromRows(items)

	// Generate podcast script using AI service
	podcastData, err := s.aiService.WritePodcast(content, nil)
	if err != nil {
		return fmt.Errorf("failed to generate podcast script: %w", err)
	}

	// Convert dialogues to JSON and update podcast
	dialoguesJSON, err := json.Marshal(podcastData.Dialogues)
	if err != nil {
		return fmt.Errorf("failed to marshal dialogues: %w", err)
	}

	updateParams := db.UpdatePodcastDialoguesParams{
		ID:        podcastID,
		Dialogues: dialoguesJSON,
	}

	if err := s.querier.UpdatePodcastDialogues(ctx, updateParams); err != nil {
		return fmt.Errorf("failed to update podcast dialogues: %w", err)
	}

	return nil
}

// buildPodcastContentFromRows combines content from multiple items (from GetPodcastItemsRow) into a single content string
func (s *podcastService) buildPodcastContentFromRows(items []db.GetPodcastItemsRow) string {
	var content strings.Builder

	for i, item := range items {
		if i > 0 {
			content.WriteString("\n\n---\n\n")
		}

		content.WriteString(fmt.Sprintf("Title: %s\n", item.Title))
		if item.Summary != nil && *item.Summary != "" {
			content.WriteString(fmt.Sprintf("Summary: %s\n", *item.Summary))
		}
		if item.TextContent != nil && *item.TextContent != "" {
			// Truncate text content if too long
			text := *item.TextContent
			if len(text) > 2000 {
				text = text[:2000] + "..."
			}
			content.WriteString(fmt.Sprintf("Content: %s\n", text))
		}
	}

	return content.String()
}

// GeneratePodcastAudio converts the podcast script to audio
func (s *podcastService) GeneratePodcastAudio(ctx context.Context, podcastID int32) error {
	// Get podcast with dialogues
	podcast, err := s.querier.GetPodcast(ctx, podcastID)
	if err != nil {
		return fmt.Errorf("failed to get podcast: %w", err)
	}

	if podcast.Dialogues == nil {
		return fmt.Errorf("no dialogues found for podcast")
	}

	// Update status to generating
	if err := s.UpdatePodcastStatus(ctx, podcastID, PodcastStatusGenerating); err != nil {
		return fmt.Errorf("failed to update podcast status: %w", err)
	}

	// Parse dialogues
	var dialogues []Dialogue
	if err := json.Unmarshal(podcast.Dialogues, &dialogues); err != nil {
		return fmt.Errorf("failed to unmarshal dialogues: %w", err)
	}

	if len(dialogues) == 0 {
		return fmt.Errorf("no dialogues to convert to audio")
	}

	// Create temp directory for audio files
	tempDir := filepath.Join(s.config.TempDir, fmt.Sprintf("podcast_%d", podcastID))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Convert each dialogue to audio
	audioFiles := make([]string, 0, len(dialogues))
	for i, dialogue := range dialogues {
		audioFile, err := s.convertDialogueToAudio(ctx, dialogue, i, tempDir)
		if err != nil {
			return fmt.Errorf("failed to convert dialogue %d to audio: %w", i, err)
		}
		audioFiles = append(audioFiles, audioFile)
	}

	// Stitch audio files together
	finalAudioFile := filepath.Join(tempDir, "final_podcast.mp3")
	if err := StitchAudioFiles(audioFiles, finalAudioFile); err != nil {
		return fmt.Errorf("failed to stitch audio files: %w", err)
	}

	// Read final audio file
	audioData, err := os.ReadFile(finalAudioFile)
	if err != nil {
		return fmt.Errorf("failed to read final audio file: %w", err)
	}

	// Store audio (either locally or upload to R2)
	audioURL, duration, err := s.storePodcastAudio(ctx, podcastID, audioData)
	if err != nil {
		return fmt.Errorf("failed to store podcast audio: %w", err)
	}

	// Update podcast with audio information and mark as completed
	updateParams := db.UpdatePodcastStatusWithAudioParams{
		ID:              podcastID,
		Status:          string(PodcastStatusCompleted),
		AudioUrl:        &audioURL,
		DurationSeconds: &duration,
	}

	if err := s.querier.UpdatePodcastStatusWithAudio(ctx, updateParams); err != nil {
		return fmt.Errorf("failed to update podcast with audio: %w", err)
	}

	return nil
}

// convertDialogueToAudio converts a single dialogue to audio
func (s *podcastService) convertDialogueToAudio(ctx context.Context, dialogue Dialogue, index int, tempDir string) (string, error) {
	// Map speaker to voice
	voice, ok := s.config.VoiceMapping[dialogue.Speaker]
	if !ok {
		return "", fmt.Errorf("unknown speaker: %s", dialogue.Speaker)
	}

	// Generate audio using speech service
	audioFile, err := s.speechService.TextToSpeech(dialogue.Content, voice, s.config.DefaultSpeed)
	if err != nil {
		return "", fmt.Errorf("failed to generate speech: %w", err)
	}

	// Download audio data
	audioData, err := s.speechService.DownloadAudio(audioFile.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download audio: %w", err)
	}

	// Save to temp file
	tempFile := filepath.Join(tempDir, fmt.Sprintf("dialogue_%d.mp3", index))
	if err := s.speechService.SaveAudio(audioData, tempFile); err != nil {
		return "", fmt.Errorf("failed to save audio file: %w", err)
	}

	return tempFile, nil
}

// storePodcastAudio stores the podcast audio in R2 and returns URL and duration
func (s *podcastService) storePodcastAudio(ctx context.Context, podcastID int32, audioData []byte) (string, int32, error) {
	if s.r2Service == nil {
		return "", 0, fmt.Errorf("R2 service not available")
	}

	// Generate R2 key for podcast audio
	key := fmt.Sprintf("generated/podcasts/podcast_%d_%d.mp3", podcastID, time.Now().Unix())

	// Upload to R2
	publicURL, err := s.r2Service.UploadFile(ctx, key, audioData, "audio/mpeg")
	if err != nil {
		return "", 0, fmt.Errorf("failed to upload podcast audio to R2: %w", err)
	}

	// Calculate duration (simplified - in real implementation, use audio metadata)
	// Assuming ~150 words per minute, average 5 characters per word
	duration := int32(len(audioData) / 10000) // Rough estimate
	if duration < 1 {
		duration = 1
	}

	return publicURL, duration, nil
}

// ProcessPodcast processes a complete podcast (script + audio generation)
func (s *podcastService) ProcessPodcast(ctx context.Context, podcastID int32) error {
	// Generate script first
	if err := s.GeneratePodcastScript(ctx, podcastID); err != nil {
		return fmt.Errorf("failed to generate podcast script: %w", err)
	}

	// Generate audio
	if err := s.GeneratePodcastAudio(ctx, podcastID); err != nil {
		return fmt.Errorf("failed to generate podcast audio: %w", err)
	}

	return nil
}

// GetPodcast retrieves a podcast by ID
func (s *podcastService) GetPodcast(ctx context.Context, podcastID int32) (*db.Podcast, error) {
	podcast, err := s.querier.GetPodcast(ctx, podcastID)
	if err != nil {
		return nil, fmt.Errorf("failed to get podcast: %w", err)
	}
	return &podcast, nil
}

// GetPodcastsByUser retrieves all podcasts for a user
func (s *podcastService) GetPodcastsByUser(ctx context.Context, userID int32) ([]db.Podcast, error) {
	userIDPtr := int32(userID)
	podcasts, err := s.querier.GetPodcastByUser(ctx, &userIDPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get user podcasts: %w", err)
	}
	return podcasts, nil
}

// GetPodcastsByStatus retrieves podcasts by status
func (s *podcastService) GetPodcastsByStatus(ctx context.Context, status PodcastStatus) ([]db.Podcast, error) {
	podcasts, err := s.querier.GetPodcastsByStatus(ctx, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to get podcasts by status: %w", err)
	}
	return podcasts, nil
}

// UpdatePodcast updates podcast metadata
func (s *podcastService) UpdatePodcast(ctx context.Context, podcastID int32, title string, description string) error {
	params := db.UpdatePodcastParams{
		ID:          podcastID,
		Title:       title,
		Description: &description,
	}

	if err := s.querier.UpdatePodcast(ctx, params); err != nil {
		return fmt.Errorf("failed to update podcast: %w", err)
	}

	return nil
}

// DeletePodcast deletes a podcast and its associated data
func (s *podcastService) DeletePodcast(ctx context.Context, podcastID int32) error {
	// Clear podcast items first
	podcastIDPtr := int32(podcastID)
	if err := s.querier.ClearPodcastItems(ctx, &podcastIDPtr); err != nil {
		return fmt.Errorf("failed to clear podcast items: %w", err)
	}

	// Delete the podcast
	if err := s.querier.DeletePodcast(ctx, podcastID); err != nil {
		return fmt.Errorf("failed to delete podcast: %w", err)
	}

	return nil
}

// AddItemToPodcast adds an item to a podcast
func (s *podcastService) AddItemToPodcast(ctx context.Context, podcastID int32, itemID int32, order int) error {
	// Get current item count to determine order if not specified
	if order < 0 {
		count, err := s.querier.CountPodcastItems(ctx, &podcastID)
		if err != nil {
			return fmt.Errorf("failed to count podcast items: %w", err)
		}
		order = int(count)
	}

	params := db.AddItemToPodcastParams{
		PodcastID: &podcastID,
		ItemID:    &itemID,
		ItemOrder: int32(order),
	}

	if _, err := s.querier.AddItemToPodcast(ctx, params); err != nil {
		return fmt.Errorf("failed to add item to podcast: %w", err)
	}

	return nil
}

// RemoveItemFromPodcast removes an item from a podcast
func (s *podcastService) RemoveItemFromPodcast(ctx context.Context, podcastID int32, itemID int32) error {
	params := db.RemoveItemFromPodcastParams{
		PodcastID: &podcastID,
		ItemID:    &itemID,
	}

	if err := s.querier.RemoveItemFromPodcast(ctx, params); err != nil {
		return fmt.Errorf("failed to remove item from podcast: %w", err)
	}

	return nil
}

// GetPodcastItems retrieves all items in a podcast
func (s *podcastService) GetPodcastItems(ctx context.Context, podcastID int32) ([]db.GetPodcastItemsRow, error) {
	items, err := s.querier.GetPodcastItems(ctx, &podcastID)
	if err != nil {
		return nil, fmt.Errorf("failed to get podcast items: %w", err)
	}
	return items, nil
}

// UpdatePodcastStatus updates the status of a podcast
func (s *podcastService) UpdatePodcastStatus(ctx context.Context, podcastID int32, status PodcastStatus) error {
	params := db.UpdatePodcastStatusParams{
		ID:     podcastID,
		Status: string(status),
	}

	if err := s.querier.UpdatePodcastStatus(ctx, params); err != nil {
		return fmt.Errorf("failed to update podcast status: %w", err)
	}

	return nil
}

// GetPendingPodcasts retrieves pending podcasts for processing
func (s *podcastService) GetPendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	podcasts, err := s.querier.GetPendingPodcasts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending podcasts: %w", err)
	}
	return podcasts, nil
}

// GetProcessingPodcasts retrieves currently processing podcasts
func (s *podcastService) GetProcessingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	podcasts, err := s.querier.GetProcessingPodcasts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get processing podcasts: %w", err)
	}
	return podcasts, nil
}

// GetPodcastAudio retrieves the audio data for a podcast
func (s *podcastService) GetPodcastAudio(ctx context.Context, podcastID int32) ([]byte, error) {
	podcast, err := s.querier.GetPodcast(ctx, podcastID)
	if err != nil {
		return nil, fmt.Errorf("failed to get podcast: %w", err)
	}

	if podcast.AudioUrl == nil {
		return nil, fmt.Errorf("podcast has no audio")
	}

	// For now, return an error indicating that audio should be fetched from the URL
	// In a real implementation, you might want to download the audio from R2
	return nil, fmt.Errorf("audio available at URL: %s - use HTTP client to download", *podcast.AudioUrl)
}

// HasPodcastAudio checks if a podcast has audio available
func (s *podcastService) HasPodcastAudio(ctx context.Context, podcastID int32) (bool, error) {
	podcast, err := s.querier.GetPodcast(ctx, podcastID)
	if err != nil {
		return false, fmt.Errorf("failed to get podcast: %w", err)
	}

	return podcast.AudioUrl != nil && *podcast.AudioUrl != "", nil
}

// GeneratePodcastUploadURL generates a presigned URL for uploading podcast audio
func (s *podcastService) GeneratePodcastUploadURL(ctx context.Context, podcastID int32) (*UploadURLResponse, error) {
	if s.r2Service == nil {
		return nil, fmt.Errorf("R2 service not available")
	}

	// Generate R2 key for podcast audio
	key := fmt.Sprintf("generated/podcasts/podcast_%d_%d.mp3", podcastID, time.Now().Unix())

	// Generate presigned URL for upload
	return s.r2Service.GenerateUploadURLForKey(ctx, key, "audio/mpeg")
}
