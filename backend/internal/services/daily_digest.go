package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/yamirghofran/briefbot/internal/db"
)

type DailyDigestService interface {
	SendDailyDigest(ctx context.Context) error
	SendDailyDigestForUser(ctx context.Context, userID int32) error
	GetDailyDigestItemsForUser(ctx context.Context, userID int32) ([]db.Item, error)
}

type dailyDigestService struct {
	querier      db.Querier
	emailService EmailService
	config       DailyDigestConfig
}

type DailyDigestConfig struct {
	Recipients []string
	Subject    string
}

func NewDailyDigestService(querier db.Querier, emailService EmailService) DailyDigestService {
	// Load configuration from environment variables
	subject := os.Getenv("DAILY_DIGEST_SUBJECT")

	if subject == "" {
		subject = "Your Daily BriefBot Digest - %s"
	}

	return &dailyDigestService{
		querier:      querier,
		emailService: emailService,
		config: DailyDigestConfig{
			Subject: subject,
		},
	}
}

func (s *dailyDigestService) SendDailyDigest(ctx context.Context) error {
	// Get all users
	users, err := s.querier.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	// Send digest to each user
	for _, user := range users {
		if user.Email != nil && *user.Email != "" {
			if err := s.SendDailyDigestForUser(ctx, user.ID); err != nil {
				// Log error but continue with other users
				fmt.Printf("Failed to send daily digest to user %d: %v\n", user.ID, err)
			}
		}
	}

	return nil
}

func (s *dailyDigestService) SendDailyDigestForUser(ctx context.Context, userID int32) error {
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
	user, err := s.querier.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user %d: %w", userID, err)
	}

	if user.Email == nil || *user.Email == "" {
		return fmt.Errorf("user %d has no email address", userID)
	}

	// Generate email content
	yesterday := time.Now().Add(-24 * time.Hour)
	htmlBody, textBody := GenerateDailyDigestEmail(items, yesterday)

	// Prepare subject with date
	subject := fmt.Sprintf("Your Daily BriefBot Digest - %s", yesterday.Format("January 2, 2006"))

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

func (s *dailyDigestService) GetDailyDigestItemsForUser(ctx context.Context, userID int32) ([]db.Item, error) {
	items, err := s.querier.GetUnreadItemsFromPreviousDayByUser(ctx, &userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread items from previous day for user %d: %w", userID, err)
	}
	return items, nil
}

// Deprecated: Use SendDailyDigestForUser instead
func (s *dailyDigestService) GetDailyDigestItems(ctx context.Context) ([]db.Item, error) {
	// For backward compatibility, return items for user ID 1
	return s.GetDailyDigestItemsForUser(ctx, 1)
}
