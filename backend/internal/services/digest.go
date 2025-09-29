package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yamirghofran/briefbot/internal/db"
)

// DigestService handles both regular and integrated (podcast + email) digest functionality
type DigestService interface {
	// Regular digest methods (backward compatibility)
	SendDailyDigest(ctx context.Context) error
	SendDailyDigestForUser(ctx context.Context, userID int32) error
	GetDailyDigestItemsForUser(ctx context.Context, userID int32) ([]db.Item, error)

	// Integrated digest methods (with podcast)
	SendIntegratedDigest(ctx context.Context) error
	SendIntegratedDigestForUser(ctx context.Context, userID int32) (*DigestResult, error)

	// Configuration methods
	SetPodcastGenerationEnabled(enabled bool)
	IsPodcastGenerationEnabled() bool
}

// DigestResult contains the result of sending a digest (regular or integrated)
type DigestResult struct {
	EmailSent  bool
	PodcastURL *string
	ItemsCount int
	Error      error
	DigestType string // "regular" or "integrated"
}

type digestService struct {
	queries        db.Querier
	emailService   EmailService
	podcastService PodcastService
	config         DigestConfig
	podcastEnabled bool
}

type DigestConfig struct {
	Subject        string
	PodcastEnabled bool // Global setting for podcast generation
}

// NewDigestService creates a new unified digest service
func NewDigestService(queries db.Querier, emailService EmailService, podcastService PodcastService) DigestService {
	// Load configuration from environment variables
	subject := os.Getenv("DAILY_DIGEST_SUBJECT")
	if subject == "" {
		subject = "Your Daily BriefBot Digest - %s"
	}

	// Check if podcast generation is globally enabled
	podcastEnabled := os.Getenv("DIGEST_PODCAST_ENABLED") == "true"

	return &digestService{
		queries:        queries,
		emailService:   emailService,
		podcastService: podcastService,
		config: DigestConfig{
			Subject:        subject,
			PodcastEnabled: podcastEnabled,
		},
		podcastEnabled: podcastEnabled,
	}
}

// SetPodcastGenerationEnabled allows runtime configuration of podcast generation
func (s *digestService) SetPodcastGenerationEnabled(enabled bool) {
	s.podcastEnabled = enabled
}

// IsPodcastGenerationEnabled returns current podcast generation setting
func (s *digestService) IsPodcastGenerationEnabled() bool {
	return s.podcastEnabled
}

// SendDailyDigest sends regular daily digest to all users (backward compatibility)
func (s *digestService) SendDailyDigest(ctx context.Context) error {
	users, err := s.queries.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	for _, user := range users {
		if user.Email != nil && *user.Email != "" {
			if err := s.SendDailyDigestForUser(ctx, user.ID); err != nil {
				// Log error but continue with other users
				log.Printf("Failed to send daily digest to user %d: %v", user.ID, err)
			}
		}
	}

	return nil
}

// SendDailyDigestForUser sends regular daily digest to a specific user (backward compatibility)
func (s *digestService) SendDailyDigestForUser(ctx context.Context, userID int32) error {
	// Get items for specific user from previous day
	items, err := s.GetDailyDigestItemsForUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get daily digest items for user %d: %w", userID, err)
	}

	// If no items, don't send email
	if len(items) == 0 {
		return nil
	}

	// Get user info for email
	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user %d: %w", userID, err)
	}

	if user.Email == nil || *user.Email == "" {
		return fmt.Errorf("user %d has no email address", userID)
	}

	// Generate regular email content (no podcast)
	yesterday := time.Now().Add(-24 * time.Hour)
	htmlBody, textBody := GenerateDailyDigestEmail(items, yesterday)

	// Prepare subject with date
	subject := fmt.Sprintf(s.config.Subject, yesterday.Format("January 2, 2006"))

	// Send email to user
	emailRequest := EmailRequest{
		ToAddresses: []string{*user.Email},
		Subject:     subject,
		HTMLBody:    htmlBody,
		TextBody:    textBody,
	}

	if err := s.emailService.SendEmail(ctx, emailRequest); err != nil {
		return fmt.Errorf("failed to send email to user %d: %w", userID, err)
	}

	return nil
}

// SendIntegratedDigest sends integrated digest (with podcast) to all users
func (s *digestService) SendIntegratedDigest(ctx context.Context) error {
	users, err := s.queries.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	var results []*DigestResult
	for _, user := range users {
		if user.Email != nil && *user.Email != "" {
			result, err := s.SendIntegratedDigestForUser(ctx, user.ID)
			if err != nil {
				log.Printf("Failed to send integrated digest for user %d: %v", user.ID, err)
			}
			results = append(results, result)
		}
	}

	// Log summary
	totalUsers := len(users)
	successfulEmails := 0
	successfulPodcasts := 0

	for _, result := range results {
		if result.EmailSent {
			successfulEmails++
		}
		if result.PodcastURL != nil {
			successfulPodcasts++
		}
	}

	log.Printf("Integrated digest batch completed: totalUsers=%d, successfulEmails=%d, successfulPodcasts=%d",
		totalUsers, successfulEmails, successfulPodcasts)

	return nil
}

// SendIntegratedDigestForUser sends integrated digest (with podcast) to a specific user
func (s *digestService) SendIntegratedDigestForUser(ctx context.Context, userID int32) (*DigestResult, error) {
	result := &DigestResult{
		EmailSent:  false,
		PodcastURL: nil,
		ItemsCount: 0,
		Error:      nil,
		DigestType: "integrated",
	}

	// Get user info
	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		result.Error = fmt.Errorf("failed to get user: %w", err)
		return result, result.Error
	}

	if user.Email == nil || *user.Email == "" {
		result.Error = fmt.Errorf("user %d has no email address", userID)
		return result, result.Error
	}

	// Get unread items from previous day
	items, err := s.GetDailyDigestItemsForUser(ctx, userID)
	if err != nil {
		result.Error = fmt.Errorf("failed to get digest items: %w", err)
		return result, result.Error
	}

	result.ItemsCount = len(items)

	if len(items) == 0 {
		log.Printf("No items for integrated digest for user %d, skipping", userID)
		return result, nil
	}

	// Generate podcast from items (if enabled)
	var podcastURL *string
	var durationSeconds *int32

	if s.podcastEnabled && s.podcastService != nil && len(items) > 0 {
		log.Printf("Generating podcast for integrated digest for user %d with %d items", userID, len(items))

		// Extract item IDs for podcast generation
		itemIDs := make([]int32, len(items))
		for i, item := range items {
			itemIDs[i] = item.ID
		}

		// Create podcast with meaningful title and description
		dateStr := time.Now().Format("January 2, 2006")
		title := fmt.Sprintf("Daily Digest for %s", dateStr)
		description := fmt.Sprintf("Your personalized daily digest with %d curated items", len(items))

		podcast, err := s.podcastService.CreatePodcastFromItems(ctx, userID, title, description, itemIDs)
		if err != nil {
			log.Printf("Failed to create podcast for user %d: %v, continuing without audio", userID, err)
			// Don't fail the entire digest if podcast generation fails
		} else {
			// Wait for podcast to complete (with timeout)
			originalPodcastID := podcast.ID // Store the original ID before the call
			podcast, err = s.waitForPodcastCompletion(ctx, originalPodcastID)
			if err != nil {
				log.Printf("Podcast generation failed for user %d, podcast %d: %v", userID, originalPodcastID, err)
			} else if podcast.AudioUrl != nil && *podcast.AudioUrl != "" {
				podcastURL = podcast.AudioUrl
				durationSeconds = podcast.DurationSeconds
				result.PodcastURL = podcastURL
				log.Printf("Podcast generated successfully for user %d: audioURL=%s", userID, *podcast.AudioUrl)
			}
		}
	}

	// Generate integrated email with podcast link (if available)
	htmlBody, textBody := GenerateIntegratedDigestEmail(items, podcastURL, durationSeconds, time.Now())

	// Send email
	subject := fmt.Sprintf("Daily Digest - %s", time.Now().Format("January 2, 2006"))
	emailReq := EmailRequest{
		ToAddresses: []string{*user.Email},
		Subject:     subject,
		HTMLBody:    htmlBody,
		TextBody:    textBody,
	}

	if err := s.emailService.SendEmail(ctx, emailReq); err != nil {
		result.Error = fmt.Errorf("failed to send email: %w", err)
		return result, result.Error
	}

	result.EmailSent = true
	log.Printf("Integrated digest sent successfully for user %d: podcastGenerated=%v", userID, podcastURL != nil)

	return result, nil
}

// GetDailyDigestItemsForUser gets unread items from previous day for a user
func (s *digestService) GetDailyDigestItemsForUser(ctx context.Context, userID int32) ([]db.Item, error) {
	items, err := s.queries.GetUnreadItemsFromPreviousDayByUser(ctx, &userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread items from previous day for user %d: %w", userID, err)
	}
	return items, nil
}

// waitForPodcastCompletion waits for podcast to complete with timeout
func (s *digestService) waitForPodcastCompletion(ctx context.Context, podcastID int32) (*db.Podcast, error) {
	timeout := time.After(5 * time.Minute) // 5 minute timeout
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("podcast generation timed out")
		case <-ticker.C:
			podcast, err := s.queries.GetPodcast(ctx, podcastID)
			if err != nil {
				return nil, err
			}

			switch podcast.Status {
			case "completed":
				return &podcast, nil
			case "failed":
				return nil, fmt.Errorf("podcast generation failed")
			case "pending", "writing", "generating":
				// Continue waiting
				continue
			default:
				return nil, fmt.Errorf("unknown podcast status: %s", podcast.Status)
			}
		}
	}
}
