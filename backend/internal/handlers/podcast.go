package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yamirghofran/briefbot/internal/services"
)

// PodcastHandler handles podcast-related HTTP requests
type PodcastHandler struct {
	podcastService services.PodcastService
	sseManager     *services.SSEManager
}

// NewPodcastHandler creates a new podcast handler
func NewPodcastHandler(podcastService services.PodcastService) *PodcastHandler {
	return &PodcastHandler{
		podcastService: podcastService,
		sseManager:     nil, // Will be set via SetSSEManager
	}
}

// SetSSEManager sets the SSE manager for the podcast handler
func (h *PodcastHandler) SetSSEManager(sseManager *services.SSEManager) {
	h.sseManager = sseManager
}

// CreatePodcast godoc
// @Summary      Create a new podcast
// @Description  Generate a podcast from multiple content items
// @Tags         podcasts
// @Accept       json
// @Produce      json
// @Param        podcast  body      CreatePodcastRequest  true  "Podcast creation request"
// @Success      201      {object}  CreatePodcastResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /podcasts [post]
func (h *PodcastHandler) CreatePodcast(c *gin.Context) {
	var req CreatePodcastRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create podcast from items
	podcast, err := h.podcastService.CreatePodcastFromItems(c.Request.Context(), req.UserID, req.Title, req.Description, req.ItemIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"podcast":           podcast,
		"message":           "Podcast created successfully and will be processed in the background",
		"processing_status": podcast.Status,
	})
}

// CreatePodcastFromSingleItem godoc
// @Summary      Create podcast from single item
// @Description  Generate a podcast from a single content item
// @Tags         podcasts
// @Accept       json
// @Produce      json
// @Param        podcast  body      CreatePodcastFromItemRequest  true  "Podcast from item request"
// @Success      201      {object}  CreatePodcastResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /podcasts/from-item [post]
func (h *PodcastHandler) CreatePodcastFromSingleItem(c *gin.Context) {
	var req CreatePodcastFromItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create podcast from single item
	podcast, err := h.podcastService.CreatePodcastFromSingleItem(c.Request.Context(), req.UserID, req.ItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"podcast":           podcast,
		"message":           "Podcast created successfully from single item and will be processed in the background",
		"processing_status": podcast.Status,
	})
}

// GetPodcast godoc
// @Summary      Get a podcast by ID
// @Description  Retrieve a podcast's information by its ID
// @Tags         podcasts
// @Produce      json
// @Param        id   path      int  true  "Podcast ID"
// @Success      200  {object}  github_com_yamirghofran_briefbot_internal_db.Podcast
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /podcasts/{id} [get]
func (h *PodcastHandler) GetPodcast(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	podcast, err := h.podcastService.GetPodcast(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, podcast)
}

// GetPodcastsByUser godoc
// @Summary      Get podcasts by user
// @Description  Retrieve all podcasts for a specific user
// @Tags         podcasts
// @Produce      json
// @Param        userID  path      int  true  "User ID"
// @Success      200     {object}  PodcastsResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /podcasts/user/{userID} [get]
func (h *PodcastHandler) GetPodcastsByUser(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	podcasts, err := h.podcastService.GetPodcastsByUser(c.Request.Context(), int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"podcasts": podcasts,
		"count":    len(podcasts),
	})
}

// GetPodcastsByStatus godoc
// @Summary      Get podcasts by status
// @Description  Retrieve podcasts filtered by their processing status
// @Tags         podcasts
// @Produce      json
// @Param        status  path      string  true  "Podcast status (pending, writing, generating, completed, failed)"
// @Success      200     {object}  PodcastsResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /podcasts/status/{status} [get]
func (h *PodcastHandler) GetPodcastsByStatus(c *gin.Context) {
	status := c.Param("status")

	// Validate status
	validStatuses := []string{"pending", "writing", "generating", "completed", "failed"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: pending, writing, generating, completed, failed"})
		return
	}

	podcasts, err := h.podcastService.GetPodcastsByStatus(c.Request.Context(), services.PodcastStatus(status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"podcasts": podcasts,
		"count":    len(podcasts),
	})
}

// GetPendingPodcasts godoc
// @Summary      Get pending podcasts
// @Description  Retrieve pending podcasts awaiting processing
// @Tags         podcasts
// @Produce      json
// @Param        limit  query     int  false  "Maximum number of podcasts to return"  default(10)
// @Success      200    {object}  PodcastsResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /podcasts/pending [get]
func (h *PodcastHandler) GetPendingPodcasts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	podcasts, err := h.podcastService.GetPendingPodcasts(c.Request.Context(), int32(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"podcasts": podcasts,
		"count":    len(podcasts),
	})
}

// GetPodcastItems godoc
// @Summary      Get podcast items
// @Description  Retrieve all content items in a podcast
// @Tags         podcasts
// @Produce      json
// @Param        id   path      int  true  "Podcast ID"
// @Success      200  {object}  PodcastItemsResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /podcasts/{id}/items [get]
func (h *PodcastHandler) GetPodcastItems(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	items, err := h.podcastService.GetPodcastItems(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

// GetPodcastAudio godoc
// @Summary      Get podcast audio
// @Description  Retrieve the audio URL for a completed podcast
// @Tags         podcasts
// @Produce      json
// @Param        id   path      int  true  "Podcast ID"
// @Success      200  {object}  PodcastAudioResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /podcasts/{id}/audio [get]
func (h *PodcastHandler) GetPodcastAudio(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	// Check if podcast has audio
	hasAudio, err := h.podcastService.HasPodcastAudio(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !hasAudio {
		c.JSON(http.StatusNotFound, gin.H{"error": "Podcast audio not available"})
		return
	}

	// Get the podcast to retrieve the audio URL
	podcast, err := h.podcastService.GetPodcast(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if podcast.AudioUrl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Podcast audio URL not found"})
		return
	}

	// Return the audio URL for client to download/stream
	c.JSON(http.StatusOK, gin.H{
		"audio_url": *podcast.AudioUrl,
		"duration":  podcast.DurationSeconds,
		"message":   "Audio available at the provided URL",
	})
}

// GeneratePodcastUploadURL godoc
// @Summary      Generate podcast upload URL
// @Description  Generate a presigned URL for uploading podcast audio to R2 storage
// @Tags         podcasts
// @Produce      json
// @Param        id   path      int  true  "Podcast ID"
// @Success      200  {object}  PodcastUploadInfo
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /podcasts/{id}/upload-url [get]
func (h *PodcastHandler) GeneratePodcastUploadURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	uploadInfo, err := h.podcastService.GeneratePodcastUploadURL(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, uploadInfo)
}

// AddItemToPodcast godoc
// @Summary      Add item to podcast
// @Description  Add a content item to an existing podcast
// @Tags         podcasts
// @Accept       json
// @Produce      json
// @Param        id    path      int                        true  "Podcast ID"
// @Param        item  body      AddItemToPodcastRequest    true  "Add item request"
// @Success      200   {object}  MessageResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /podcasts/{id}/items [post]
func (h *PodcastHandler) AddItemToPodcast(c *gin.Context) {
	podcastIDStr := c.Param("id")
	podcastID, err := strconv.ParseInt(podcastIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	var req AddItemToPodcastRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.podcastService.AddItemToPodcast(c.Request.Context(), int32(podcastID), req.ItemID, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to podcast successfully",
	})
}

// RemoveItemFromPodcast godoc
// @Summary      Remove item from podcast
// @Description  Remove a content item from a podcast
// @Tags         podcasts
// @Produce      json
// @Param        id      path      int  true  "Podcast ID"
// @Param        itemID  path      int  true  "Item ID"
// @Success      200     {object}  MessageResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /podcasts/{id}/items/{itemID} [delete]
func (h *PodcastHandler) RemoveItemFromPodcast(c *gin.Context) {
	podcastIDStr := c.Param("id")
	podcastID, err := strconv.ParseInt(podcastIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	itemIDStr := c.Param("itemID")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = h.podcastService.RemoveItemFromPodcast(c.Request.Context(), int32(podcastID), int32(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from podcast successfully",
	})
}

// UpdatePodcast godoc
// @Summary      Update podcast
// @Description  Update podcast metadata (title and description)
// @Tags         podcasts
// @Accept       json
// @Produce      json
// @Param        id       path      int                     true  "Podcast ID"
// @Param        podcast  body      UpdatePodcastRequest    true  "Podcast update request"
// @Success      200      {object}  MessageResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /podcasts/{id} [put]
func (h *PodcastHandler) UpdatePodcast(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	var req UpdatePodcastRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.podcastService.UpdatePodcast(c.Request.Context(), int32(id), req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Podcast updated successfully",
	})
}

// DeletePodcast godoc
// @Summary      Delete podcast
// @Description  Delete a podcast from the system
// @Tags         podcasts
// @Produce      json
// @Param        id   path      int  true  "Podcast ID"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /podcasts/{id} [delete]
func (h *PodcastHandler) DeletePodcast(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	err = h.podcastService.DeletePodcast(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Podcast deleted successfully",
	})
}

// GetPodcastProcessingStatus godoc
// @Summary      Get podcast processing status
// @Description  Retrieve the processing status of a podcast
// @Tags         podcasts
// @Produce      json
// @Param        id   path      int  true  "Podcast ID"
// @Success      200  {object}  PodcastProcessingStatusResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /podcasts/{id}/status [get]
func (h *PodcastHandler) GetPodcastProcessingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	podcast, err := h.podcastService.GetPodcast(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Determine processing flags based on status
	isPending := podcast.Status == string(services.PodcastStatusPending)
	isWriting := podcast.Status == string(services.PodcastStatusWriting)
	isGenerating := podcast.Status == string(services.PodcastStatusGenerating)
	isCompleted := podcast.Status == string(services.PodcastStatusCompleted)
	isFailed := podcast.Status == string(services.PodcastStatusFailed)
	isProcessing := isWriting || isGenerating

	c.JSON(http.StatusOK, gin.H{
		"podcast_id":    podcast.ID,
		"status":        podcast.Status,
		"is_pending":    isPending,
		"is_writing":    isWriting,
		"is_generating": isGenerating,
		"is_processing": isProcessing,
		"is_completed":  isCompleted,
		"is_failed":     isFailed,
		"audio_url":     podcast.AudioUrl,
	})
}

// StreamPodcastUpdates godoc
// @Summary      Stream podcast updates
// @Description  Server-Sent Events endpoint for real-time podcast processing updates
// @Tags         podcasts
// @Produce      text/event-stream
// @Param        userID  path      int  true  "User ID"
// @Success      200     {string}  string "SSE stream"
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /podcasts/user/{userID}/stream [get]
func (h *PodcastHandler) StreamPodcastUpdates(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if h.sseManager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SSE manager not available"})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Create a channel for writing SSE messages
	messageChan := make(chan string, 10)

	// Register client with SSE manager
	client := h.sseManager.AddClient(int32(userID))
	defer h.sseManager.RemoveClient(client)

	// Start goroutine to write messages from the client channel
	go func() {
		defer close(messageChan)
		services.WriteSSEMessage(c.Request.Context(), client, messageChan)
	}()

	// Stream messages to client
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// Send messages
	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			fmt.Fprint(c.Writer, msg)
			flusher.Flush()
		case <-c.Request.Context().Done():
			log.Printf("SSE: Client disconnected for user %d", userID)
			return
		}
	}
}
