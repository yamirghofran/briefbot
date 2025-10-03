package handlers

import "github.com/yamirghofran/briefbot/internal/db"

// User request/response models

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name         *string `json:"name" example:"John Doe"`
	Email        *string `json:"email" example:"john@example.com"`
	AuthProvider *string `json:"auth_provider" example:"google"`
	OauthID      *string `json:"oauth_id" example:"google_123456"`
	PasswordHash *string `json:"password_hash"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name         *string `json:"name" example:"Jane Doe"`
	Email        *string `json:"email" example:"jane@example.com"`
	AuthProvider *string `json:"auth_provider"`
	OauthID      *string `json:"oauth_id"`
	PasswordHash *string `json:"password_hash"`
}

// Item request/response models

// CreateItemRequest represents the request body for creating an item
type CreateItemRequest struct {
	UserID *int32  `json:"user_id" binding:"required" example:"1"`
	URL    *string `json:"url" binding:"required" example:"https://example.com/article"`
}

// CreateItemResponse represents the response after creating an item
type CreateItemResponse struct {
	Item             db.Item `json:"item"`
	Message          string  `json:"message"`
	ProcessingStatus string  `json:"processing_status"`
}

// UpdateItemRequest represents the request body for updating an item
type UpdateItemRequest struct {
	Title       string   `json:"title" example:"Article Title"`
	URL         *string  `json:"url"`
	TextContent *string  `json:"text_content"`
	Summary     *string  `json:"summary"`
	Type        *string  `json:"type" example:"article"`
	Platform    *string  `json:"platform" example:"web"`
	Tags        []string `json:"tags"`
	Authors     []string `json:"authors"`
	IsRead      *bool    `json:"is_read"`
}

// PatchItemRequest represents the request body for patching an item
type PatchItemRequest struct {
	Title   *string  `json:"title" example:"Updated Title"`
	Summary *string  `json:"summary" example:"Updated summary"`
	Tags    []string `json:"tags" example:"tech,news"`
	Authors []string `json:"authors" example:"John Doe"`
}

// ItemProcessingStatusResponse represents the processing status of an item
type ItemProcessingStatusResponse struct {
	ItemID           int32   `json:"item_id"`
	ProcessingStatus string  `json:"processing_status"`
	IsProcessing     bool    `json:"is_processing"`
	IsCompleted      bool    `json:"is_completed"`
	IsFailed         bool    `json:"is_failed"`
	ProcessingError  *string `json:"processing_error"`
}

// ItemsByStatusResponse represents items filtered by processing status
type ItemsByStatusResponse struct {
	Status string    `json:"status"`
	Items  []db.Item `json:"items"`
	Count  int       `json:"count"`
}

// Podcast request/response models

// CreatePodcastRequest represents the request body for creating a podcast
type CreatePodcastRequest struct {
	UserID      int32   `json:"user_id" binding:"required" example:"1"`
	Title       string  `json:"title" binding:"required" example:"My Daily Digest"`
	Description string  `json:"description" example:"Today's top articles"`
	ItemIDs     []int32 `json:"item_ids" binding:"required,min=1" example:"1,2,3"`
}

// CreatePodcastFromItemRequest represents the request body for creating a podcast from a single item
type CreatePodcastFromItemRequest struct {
	UserID int32 `json:"user_id" binding:"required" example:"1"`
	ItemID int32 `json:"item_id" binding:"required" example:"1"`
}

// CreatePodcastResponse represents the response after creating a podcast
type CreatePodcastResponse struct {
	Podcast          db.Podcast `json:"podcast"`
	Message          string     `json:"message"`
	ProcessingStatus string     `json:"processing_status"`
}

// AddItemToPodcastRequest represents the request body for adding an item to a podcast
type AddItemToPodcastRequest struct {
	ItemID int32 `json:"item_id" binding:"required" example:"1"`
	Order  int   `json:"order" example:"1"`
}

// UpdatePodcastRequest represents the request body for updating a podcast
type UpdatePodcastRequest struct {
	Title       string `json:"title" binding:"required" example:"Updated Title"`
	Description string `json:"description" example:"Updated description"`
}

// PodcastsResponse represents a list of podcasts
type PodcastsResponse struct {
	Podcasts []db.Podcast `json:"podcasts"`
	Count    int          `json:"count"`
}

// PodcastItemsResponse represents items in a podcast
type PodcastItemsResponse struct {
	Items []db.Item `json:"items"`
	Count int       `json:"count"`
}

// PodcastAudioResponse represents the audio information for a podcast
type PodcastAudioResponse struct {
	AudioURL string `json:"audio_url"`
	Duration *int32 `json:"duration"`
	Message  string `json:"message"`
}

// PodcastProcessingStatusResponse represents the processing status of a podcast
type PodcastProcessingStatusResponse struct {
	PodcastID    int32   `json:"podcast_id"`
	Status       string  `json:"status"`
	IsPending    bool    `json:"is_pending"`
	IsWriting    bool    `json:"is_writing"`
	IsGenerating bool    `json:"is_generating"`
	IsProcessing bool    `json:"is_processing"`
	IsCompleted  bool    `json:"is_completed"`
	IsFailed     bool    `json:"is_failed"`
	AudioURL     *string `json:"audio_url"`
}

// PodcastUploadInfo represents upload URL information for a podcast
type PodcastUploadInfo struct {
	UploadURL string `json:"upload_url"`
	PublicURL string `json:"public_url"`
	ExpiresAt string `json:"expires_at"`
}

// Generic response models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// MessageResponse represents a success message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
