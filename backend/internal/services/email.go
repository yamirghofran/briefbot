package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/yamirghofran/briefbot/internal/db"
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

// GenerateDailyDigestEmail generates HTML and text content for a daily digest email
func GenerateDailyDigestEmail(items []db.Item, date time.Time) (string, string) {
	// Generate HTML content
	htmlContent := generateDailyDigestHTML(items, date)

	// Generate text content
	textContent := generateDailyDigestText(items, date)

	return htmlContent, textContent
}

func generateDailyDigestHTML(items []db.Item, date time.Time) string {
	var html strings.Builder

	html.WriteString(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Daily Digest - %s</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .header h1 { color: #2c3e50; margin: 0; }
        .item { background-color: #ffffff; border: 1px solid #e9ecef; border-radius: 8px; padding: 20px; margin-bottom: 15px; }
        .item-title { font-size: 18px; font-weight: bold; color: #2c3e50; margin-bottom: 8px; }
        .item-meta { color: #6c757d; font-size: 14px; margin-bottom: 10px; }
        .item-summary { color: #495057; line-height: 1.5; }
        .item-link { color: #007bff; text-decoration: none; }
        .item-link:hover { text-decoration: underline; }
        .footer { text-align: center; color: #6c757d; font-size: 12px; margin-top: 30px; padding-top: 20px; border-top: 1px solid #e9ecef; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Daily Digest - %s</h1>
        <p>Your unread items from yesterday</p>
    </div>
`, date.Format("January 2, 2006"), date.Format("January 2, 2006")))

	for _, item := range items {
		html.WriteString(fmt.Sprintf(`
    <div class="item">
        <div class="item-title"><a href="%s" class="item-link">%s</a></div>
        <div class="item-meta">`,
			*item.Url, item.Title))

		if item.Platform != nil && *item.Platform != "" {
			html.WriteString(fmt.Sprintf("%s | ", *item.Platform))
		}
		if item.Type != nil && *item.Type != "" {
			html.WriteString(fmt.Sprintf("%s", *item.Type))
		}
		html.WriteString("</div>")

		if item.Summary != nil && *item.Summary != "" {
			html.WriteString(fmt.Sprintf(`
        <div class="item-summary">%s</div>`, *item.Summary))
		}

		html.WriteString("\n    </div>")
	}

	html.WriteString(`
    <div class="footer">
        <p>Sent by BriefBot - Your personal content curator</p>
    </div>
</body>
</html>`)

	return html.String()
}

func generateDailyDigestText(items []db.Item, date time.Time) string {
	var text strings.Builder

	text.WriteString(fmt.Sprintf("Daily Digest - %s\n", date.Format("January 2, 2006")))
	text.WriteString("Your unread items from yesterday\n")
	text.WriteString(strings.Repeat("=", 50) + "\n\n")

	for i, item := range items {
		text.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Title))
		text.WriteString(fmt.Sprintf("   Link: %s\n", *item.Url))

		var meta []string
		if item.Platform != nil && *item.Platform != "" {
			meta = append(meta, fmt.Sprintf("Platform: %s", *item.Platform))
		}
		if item.Type != nil && *item.Type != "" {
			meta = append(meta, fmt.Sprintf("Type: %s", *item.Type))
		}
		meta = append(meta, fmt.Sprintf("Added: %s", item.CreatedAt.Format("Jan 2, 3:04 PM")))

		text.WriteString(fmt.Sprintf("   %s\n", strings.Join(meta, " | ")))

		if item.Summary != nil && *item.Summary != "" {
			text.WriteString(fmt.Sprintf("   Summary: %s\n", *item.Summary))
		}

		text.WriteString("\n")
	}

	text.WriteString(strings.Repeat("-", 50) + "\n")
	text.WriteString("Sent by BriefBot - Your personal content curator\n")

	return text.String()
}
