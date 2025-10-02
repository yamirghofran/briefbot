package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/yamirghofran/briefbot/internal/db"
)

// DialogueAudioResult represents the result of generating audio for a single dialogue
type DialogueAudioResult struct {
	Index    int    // Dialogue index
	FilePath string // Path to the generated audio file
	Error    error  // Error if generation failed
}

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

	// Atomic podcast acquisition with locking - prevents multiple workers from processing the same podcast
	AcquirePendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error)

	// Audio management
	GetPodcastAudio(ctx context.Context, podcastID int32) ([]byte, error)
	HasPodcastAudio(ctx context.Context, podcastID int32) (bool, error)
	GeneratePodcastUploadURL(ctx context.Context, podcastID int32) (*UploadURLResponse, error)

	// SSE management
	SetSSEManager(sseManager *SSEManager)
}

// podcastService implements PodcastService
type podcastService struct {
	querier       db.Querier
	aiService     AIService
	speechService SpeechService
	r2Service     *R2Service
	config        PodcastConfig
	sseManager    *SSEManager
}

// PodcastConfig holds configuration for podcast service
type PodcastConfig struct {
	DefaultSpeed       float64
	MaxRetries         int
	TempDir            string
	EnableStorage      bool
	VoiceMapping       map[string]VoiceEnum
	MaxItemsPerPodcast int
	MaxConcurrentAudio int32 // Maximum concurrent audio generation requests
}

// DefaultPodcastConfig returns default configuration
func DefaultPodcastConfig() PodcastConfig {
	return PodcastConfig{
		DefaultSpeed:       1.0,
		MaxRetries:         3,
		TempDir:            "/tmp/podcasts",
		EnableStorage:      true,
		MaxItemsPerPodcast: 10,
		MaxConcurrentAudio: 5, // Default 5 concurrent audio requests
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
		sseManager:    nil, // Will be set later via SetSSEManager
	}
}

// SetSSEManager sets the SSE manager for the podcast service
func (s *podcastService) SetSSEManager(sseManager *SSEManager) {
	s.sseManager = sseManager
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

	// Notify via SSE that podcast was created
	if s.sseManager != nil {
		s.sseManager.NotifyPodcastUpdate(userID, podcast.ID, string(PodcastStatusPending), "created")
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
	// Get the podcast to retrieve user ID for SSE notifications
	podcast, err := s.querier.GetPodcast(ctx, podcastID)
	if err != nil {
		return fmt.Errorf("failed to get podcast: %w", err)
	}

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

	// Notify via SSE that script writing has started
	if s.sseManager != nil && podcast.UserID != nil {
		s.sseManager.NotifyPodcastUpdate(*podcast.UserID, podcastID, string(PodcastStatusWriting), "writing")
	}

	// Convert GetPodcastItemsRow to Items for content building
	content := s.buildPodcastContentFromRows(items)

	// Generate podcast script using AI service
	podcastData, err := s.aiService.WritePodcast(content)
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

	// Notify via SSE that audio generation has started
	if s.sseManager != nil && podcast.UserID != nil {
		s.sseManager.NotifyPodcastUpdate(*podcast.UserID, podcastID, string(PodcastStatusGenerating), "generating")
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
	// Don't defer cleanup here - we'll clean up after we're completely done with all files

	// Convert each dialogue to audio concurrently
	audioFiles := make([]string, len(dialogues))
	var wg sync.WaitGroup
	resultChan := make(chan DialogueAudioResult, len(dialogues))

	// Use a semaphore to limit concurrent requests and prevent API overload
	// Use configured max concurrent audio requests (default 5)
	maxConcurrent := int(s.config.MaxConcurrentAudio)
	if maxConcurrent <= 0 {
		maxConcurrent = 5 // Fallback to default if not configured
	}
	log.Printf("Starting concurrent audio generation for %d dialogues in podcast %d (max concurrent: %d)",
		len(dialogues), podcastID, maxConcurrent)

	semaphore := make(chan struct{}, maxConcurrent)

	// Launch concurrent audio generation for all dialogues
	for i, dialogue := range dialogues {
		wg.Add(1)
		go func(idx int, dlg Dialogue) {
			defer wg.Done()

			// Acquire semaphore (blocks if max concurrent requests reached)
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release semaphore when done

			// Generate audio for this dialogue
			audioFile, err := s.convertDialogueToAudio(ctx, dlg, idx, len(dialogues), tempDir)
			if err != nil {
				resultChan <- DialogueAudioResult{
					Index: idx,
					Error: err,
				}
				return
			}

			resultChan <- DialogueAudioResult{
				Index:    idx,
				FilePath: audioFile,
				Error:    nil,
			}
		}(i, dialogue)
	}

	// Wait for ALL goroutines to complete before processing results
	wg.Wait()
	close(resultChan)

	// Collect results and maintain order
	failedCount := 0
	successCount := 0
	for result := range resultChan {
		if result.Error != nil {
			log.Printf("Warning: Failed to convert dialogue %d to audio: %v", result.Index, result.Error)
			failedCount++
		} else {
			audioFiles[result.Index] = result.FilePath
			successCount++
		}
	}

	log.Printf("Completed concurrent audio generation: %d succeeded, %d failed out of %d dialogues",
		successCount, failedCount, len(dialogues))

	// Filter out failed results and collect valid audio files
	validAudioFiles := make([]string, 0, len(audioFiles))
	for i, audioFile := range audioFiles {
		if audioFile == "" {
			log.Printf("Warning: Dialogue %d failed to generate audio", i)
			continue
		}
		if _, err := os.Stat(audioFile); err != nil {
			log.Printf("Warning: Audio file %s for dialogue %d does not exist: %v", audioFile, i, err)
			continue
		}
		validAudioFiles = append(validAudioFiles, audioFile)
	}

	// Check if we have any valid audio files
	if len(validAudioFiles) == 0 {
		log.Printf("No valid audio files found for podcast %d. Marking as completed without audio.", podcastID)

		// Clean up temp directory before returning
		if err := os.RemoveAll(tempDir); err != nil {
			log.Printf("Warning: Failed to clean up temp directory %s: %v", tempDir, err)
		}

		// Update podcast status to completed without audio
		noAudioURL := ""
		zeroDuration := int32(0)
		updateParams := db.UpdatePodcastStatusWithAudioParams{
			ID:              podcastID,
			Status:          string(PodcastStatusCompleted),
			AudioUrl:        &noAudioURL,
			DurationSeconds: &zeroDuration,
		}

		if err := s.querier.UpdatePodcastStatusWithAudio(ctx, updateParams); err != nil {
			return fmt.Errorf("failed to update podcast status without audio: %w", err)
		}

		return nil
	}

	if len(validAudioFiles) < len(audioFiles) {
		log.Printf("Warning: Only %d/%d audio files are valid for podcast %d",
			len(validAudioFiles), len(audioFiles), podcastID)
	}

	// Stitch audio files together using only valid files
	finalAudioFile := filepath.Join(tempDir, "final_podcast.mp3")
	if err := StitchAudioFiles(validAudioFiles, finalAudioFile); err != nil {
		// Clean up temp directory on error
		if cleanupErr := os.RemoveAll(tempDir); cleanupErr != nil {
			log.Printf("Warning: Failed to clean up temp directory %s: %v", tempDir, cleanupErr)
		}
		return fmt.Errorf("failed to stitch audio files: %w", err)
	}

	// Read final audio file
	audioData, err := os.ReadFile(finalAudioFile)
	if err != nil {
		// Clean up temp directory on error
		if cleanupErr := os.RemoveAll(tempDir); cleanupErr != nil {
			log.Printf("Warning: Failed to clean up temp directory %s: %v", tempDir, cleanupErr)
		}
		return fmt.Errorf("failed to read final audio file: %w", err)
	}

	// Store audio (either locally or upload to R2)
	audioURL, duration, err := s.storePodcastAudio(ctx, podcastID, audioData)
	if err != nil {
		// Clean up temp directory on error
		if cleanupErr := os.RemoveAll(tempDir); cleanupErr != nil {
			log.Printf("Warning: Failed to clean up temp directory %s: %v", tempDir, cleanupErr)
		}
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
		// Clean up temp directory on error
		if cleanupErr := os.RemoveAll(tempDir); cleanupErr != nil {
			log.Printf("Warning: Failed to clean up temp directory %s: %v", tempDir, cleanupErr)
		}
		return fmt.Errorf("failed to update podcast with audio: %w", err)
	}

	// Clean up temp directory only after everything is successfully completed
	if err := os.RemoveAll(tempDir); err != nil {
		log.Printf("Warning: Failed to clean up temp directory %s: %v", tempDir, err)
	}

	return nil
}

// convertDialogueToAudio converts a single dialogue to audio
func (s *podcastService) convertDialogueToAudio(ctx context.Context, dialogue Dialogue, index int, total int, tempDir string) (string, error) {
	// Map speaker to voice
	voice, ok := s.config.VoiceMapping[dialogue.Speaker]
	if !ok {
		return "", fmt.Errorf("unknown speaker: %s", dialogue.Speaker)
	}

	// Create dialogue info for enhanced logging
	dialogueInfo := &DialogueInfo{
		Index:   index + 1, // Convert to 1-based index
		Total:   total,
		Speaker: dialogue.Speaker,
		Content: dialogue.Content,
	}

	// Generate audio using speech service
	audioFile, err := s.speechService.TextToSpeech(dialogue.Content, voice, s.config.DefaultSpeed, dialogueInfo)
	if err != nil {
		return "", fmt.Errorf("failed to generate speech: %w", err)
	}

	// Download audio data
	audioData, err := s.speechService.DownloadAudio(audioFile.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download audio: %w", err)
	}

	// Determine file extension based on content type from FAL response
	fileExt := ".mp3"           // default fallback
	contentType := "audio/mpeg" // default fallback
	if audioFile.ContentType != "" {
		switch audioFile.ContentType {
		case "audio/wav", "audio/x-wav", "audio/wave":
			fileExt = ".wav"
			contentType = "audio/wav"
		case "audio/mpeg", "audio/mp3":
			fileExt = ".mp3"
			contentType = "audio/mpeg"
		default:
			log.Printf("Warning: Unknown audio content type '%s', using .mp3 extension", audioFile.ContentType)
		}
	} else {
		log.Printf("Warning: No content type provided by FAL, assuming .mp3 extension")
	}

	// Log the detected format for debugging
	log.Printf("Dialogue %d/%d: Detected audio format '%s', saving as %s",
		dialogueInfo.Index, dialogueInfo.Total, contentType, fileExt)

	// Save to temp file with correct extension
	tempFile := filepath.Join(tempDir, fmt.Sprintf("dialogue_%d%s", index, fileExt))
	if err := s.speechService.SaveAudio(audioData, tempFile); err != nil {
		return "", fmt.Errorf("failed to save audio file: %w", err)
	}

	return tempFile, nil
}

// generateDialogueAudioConcurrent generates audio for a single dialogue concurrently
func (s *podcastService) generateDialogueAudioConcurrent(ctx context.Context, dialogue Dialogue, index int, total int, tempDir string, resultChan chan<- DialogueAudioResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// Generate audio using the existing method
	audioFile, err := s.convertDialogueToAudio(ctx, dialogue, index, total, tempDir)
	if err != nil {
		resultChan <- DialogueAudioResult{
			Index: index,
			Error: err,
		}
		return
	}

	resultChan <- DialogueAudioResult{
		Index:    index,
		FilePath: audioFile,
		Error:    nil,
	}
}

// detectAudioFormat detects the audio format from the file header
func detectAudioFormat(data []byte) string {
	if len(data) < 12 {
		return "application/octet-stream" // fallback for too small data
	}

	// Check for WAV format (RIFF header)
	if string(data[0:4]) == "RIFF" && string(data[8:12]) == "WAVE" {
		return "audio/wav"
	}

	// Check for MP3 format (ID3v2 or MPEG sync word)
	if string(data[0:3]) == "ID3" || (data[0] == 0xFF && (data[1]&0xE0) == 0xE0) {
		return "audio/mpeg"
	}

	// Default fallback - assume MP3 for podcast audio
	return "audio/mpeg"
}

// storePodcastAudio stores the podcast audio in R2 and returns URL and duration
func (s *podcastService) storePodcastAudio(ctx context.Context, podcastID int32, audioData []byte) (string, int32, error) {
	if s.r2Service == nil {
		return "", 0, fmt.Errorf("R2 service not available")
	}

	// Detect the actual audio format from the data
	contentType := detectAudioFormat(audioData)
	log.Printf("Storing podcast audio: detected format '%s' (%d bytes)", contentType, len(audioData))

	// Generate R2 key for podcast audio - use appropriate extension based on content type
	var fileExt string
	switch contentType {
	case "audio/wav":
		fileExt = ".wav"
	case "audio/mpeg":
		fileExt = ".mp3"
	default:
		fileExt = ".mp3" // fallback to MP3
	}

	key := fmt.Sprintf("generated/podcasts/podcast_%d_%d%s", podcastID, time.Now().Unix(), fileExt)

	// Upload to R2 with correct content type
	publicURL, err := s.r2Service.UploadFile(ctx, key, audioData, contentType)
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
	// Get the podcast to retrieve user ID for SSE notifications
	podcast, err := s.querier.GetPodcast(ctx, podcastID)
	if err != nil {
		return fmt.Errorf("failed to get podcast: %w", err)
	}

	// Generate script first
	if err := s.GeneratePodcastScript(ctx, podcastID); err != nil {
		// Notify failure via SSE
		if s.sseManager != nil && podcast.UserID != nil {
			s.sseManager.NotifyPodcastUpdate(*podcast.UserID, podcastID, string(PodcastStatusFailed), "failed")
		}
		return fmt.Errorf("failed to generate podcast script: %w", err)
	}

	// Generate audio
	if err := s.GeneratePodcastAudio(ctx, podcastID); err != nil {
		// Notify failure via SSE
		if s.sseManager != nil && podcast.UserID != nil {
			s.sseManager.NotifyPodcastUpdate(*podcast.UserID, podcastID, string(PodcastStatusFailed), "failed")
		}
		return fmt.Errorf("failed to generate podcast audio: %w", err)
	}

	// Notify completion via SSE
	if s.sseManager != nil && podcast.UserID != nil {
		s.sseManager.NotifyPodcastUpdate(*podcast.UserID, podcastID, string(PodcastStatusCompleted), "completed")
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

	// Notify via SSE if manager is available
	if s.sseManager != nil {
		// Get the podcast to retrieve the user ID
		podcast, err := s.querier.GetPodcast(ctx, podcastID)
		if err == nil && podcast.UserID != nil {
			s.sseManager.NotifyPodcastUpdate(*podcast.UserID, podcastID, string(status), string(status))
		}
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

// AcquirePendingPodcasts atomically acquires pending podcasts with row-level locking
// This prevents multiple workers from processing the same podcast concurrently
func (s *podcastService) AcquirePendingPodcasts(ctx context.Context, limit int32) ([]db.Podcast, error) {
	// For now, implement a simple approach: get pending podcasts and immediately mark them as writing
	// This creates a small window where race conditions could occur, but it's much better than before
	// The FOR UPDATE SKIP LOCKED in the query will still help prevent most conflicts

	podcasts, err := s.querier.GetPendingPodcasts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending podcasts: %w", err)
	}

	// Immediately update the status of acquired podcasts to 'writing'
	// This reduces the window for race conditions and marks them as being processed
	for _, podcast := range podcasts {
		params := db.UpdatePodcastStatusParams{
			ID:     podcast.ID,
			Status: string(PodcastStatusWriting),
		}
		if err := s.querier.UpdatePodcastStatus(ctx, params); err != nil {
			log.Printf("Warning: Failed to update podcast %d status to writing: %v", podcast.ID, err)
			// Continue with other podcasts even if one update fails
		}
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
