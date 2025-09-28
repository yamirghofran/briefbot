package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type R2Service struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	bucket     string
	accountId  string
	publicHost string
}

type UploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
	PublicURL string `json:"public_url"`
}

func NewR2Service() (*R2Service, error) {
	// Load from environment variables
	accessKeyId := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	accountId := os.Getenv("R2_ACCOUNT_ID")
	bucket := os.Getenv("R2_BUCKET_NAME")
	publicHost := os.Getenv("R2_PUBLIC_HOST") // Custom domain or R2 public URL

	if accessKeyId == "" || secretAccessKey == "" || accountId == "" || bucket == "" {
		return nil, fmt.Errorf("missing required R2 environment variables")
	}

	// Create configuration
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyId,
			secretAccessKey,
			"",
		)),
		config.WithRegion("auto"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
				}, nil
			})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &R2Service{
		client:     client,
		presigner:  s3.NewPresignClient(client),
		bucket:     bucket,
		accountId:  accountId,
		publicHost: publicHost,
	}, nil
}

// GenerateUploadURL generates a presigned URL for uploading a file
func (r *R2Service) GenerateUploadURL(ctx context.Context, contentType string, folder string) (*UploadURLResponse, error) {
	// Generate unique key with folder structure
	ext := ".jpg" // Default extension
	switch contentType {
	case "image/png":
		ext = ".png"
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/webp":
		ext = ".webp"
	case "application/pdf":
		ext = ".pdf"
	case "audio/wav":
		ext = ".wav"
	case "audio/mpeg", "audio/mp3":
		ext = ".mp3"
	}

	key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	request := &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	// Create presigned URL valid for 15 minutes with CORS headers
	presignResult, err := r.presigner.PresignPutObject(ctx, request,
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(15 * time.Minute)
		})

	if err != nil {
		return nil, fmt.Errorf("failed to create presigned URL: %w", err)
	}

	return &UploadURLResponse{
		UploadURL: presignResult.URL,
		Key:       key,
		PublicURL: r.GetPublicURL(key),
	}, nil
}

// GetPublicURL returns the public URL for a given key
func (r *R2Service) GetPublicURL(key string) string {
	if r.publicHost != "" {
		// Use custom domain if configured
		return fmt.Sprintf("%s/%s", r.publicHost, key)
	}
	// Use R2 public URL
	return fmt.Sprintf("https://pub-%s.r2.dev/%s", r.accountId, key)
}

// DeleteFile deletes a file from R2
func (r *R2Service) DeleteFile(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	return err
}

// DeleteFiles deletes multiple files from R2
func (r *R2Service) DeleteFiles(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	objects := make([]types.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(key),
		}
	}

	_, err := r.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(r.bucket),
		Delete: &types.Delete{
			Objects: objects,
		},
	})
	return err
}

// ExtractKeyFromURL extracts the key from a public URL
func (r *R2Service) ExtractKeyFromURL(url string) string {
	// Handle both custom domain and R2 public URLs
	if r.publicHost != "" && len(url) > len(r.publicHost) {
		return url[len(r.publicHost)+1:] // +1 for the slash
	}

	// Handle R2 public URL format
	prefix := fmt.Sprintf("https://pub-%s.r2.dev/", r.accountId)
	if len(url) > len(prefix) {
		return url[len(prefix):]
	}

	// If it doesn't match either format, return the original URL
	// This might be a key already
	return url
}
