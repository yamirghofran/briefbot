package services

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// EmailRequest contains the properties needed to send an email
type EmailRequest struct {
	ToAddresses    []string // Recipient email addresses
	Subject        string   // Email subject
	HTMLBody       string   // HTML content (optional)
	TextBody       string   // Plain text content (optional)
	ReplyToAddress string   // Reply-to address (optional, defaults to service config)
}

type EmailService interface {
	SendEmail(ctx context.Context, request EmailRequest) error
}

type emailService struct {
	client       *ses.Client
	fromEmail    string
	fromName     string
	replyToEmail string
}

func NewEmailService() (EmailService, error) {
	// Load from environment variables
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	fromEmail := os.Getenv("SES_FROM_EMAIL")
	fromName := os.Getenv("SES_FROM_NAME")
	replyToEmail := os.Getenv("SES_REPLY_TO_EMAIL")

	if accessKeyId == "" || secretAccessKey == "" || region == "" || fromEmail == "" {
		return nil, fmt.Errorf("missing required AWS SES environment variables")
	}

	if fromName == "" {
		fromName = "BriefBot"
	}
	if replyToEmail == "" {
		replyToEmail = fromEmail
	}

	// Create configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyId,
			secretAccessKey,
			"",
		)),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create SES client
	client := ses.NewFromConfig(cfg)

	return &emailService{
		client:       client,
		fromEmail:    fromEmail,
		fromName:     fromName,
		replyToEmail: replyToEmail,
	}, nil
}

func (s *emailService) SendEmail(ctx context.Context, request EmailRequest) error {
	// Validate required fields
	if len(request.ToAddresses) == 0 {
		return fmt.Errorf("to addresses cannot be empty")
	}
	if request.Subject == "" {
		return fmt.Errorf("subject cannot be empty")
	}
	if request.HTMLBody == "" && request.TextBody == "" {
		return fmt.Errorf("either HTML body or text body must be provided")
	}

	// Create the email input
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: request.ToAddresses,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(request.Subject),
			},
		},
		Source: aws.String(fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)),
	}

	// Set up the message body
	body := &types.Body{}
	if request.HTMLBody != "" {
		body.Html = &types.Content{
			Charset: aws.String("UTF-8"),
			Data:    aws.String(request.HTMLBody),
		}
	}
	if request.TextBody != "" {
		body.Text = &types.Content{
			Charset: aws.String("UTF-8"),
			Data:    aws.String(request.TextBody),
		}
	}
	input.Message.Body = body

	// Set reply-to address if provided, otherwise use service default
	replyTo := s.replyToEmail
	if request.ReplyToAddress != "" {
		replyTo = request.ReplyToAddress
	}
	input.ReplyToAddresses = []string{replyTo}

	// Send the email
	_, err := s.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
