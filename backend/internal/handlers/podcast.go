package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yamirghofran/briefbot/internal/services"
)

// PodcastHandler handles podcast-related HTTP requests
type PodcastHandler struct {
	podcastService services.PodcastService
}

// NewPodcastHandler creates a new podcast handler
func NewPodcastHandler(podcastService services.PodcastService) *PodcastHandler {
	return &PodcastHandler{
		podcastService: podcastService,
	}
}

// CreatePodcast creates a new podcast from items
func (h *PodcastHandler) CreatePodcast(c *gin.Context) {
	var req struct {
		UserID      int32   `json:"user_id" binding:"required"`
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		ItemIDs     []int32 `json:"item_ids" binding:"required,min=1"`
	}

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

// CreatePodcastFromSingleItem creates a podcast from a single item
func (h *PodcastHandler) CreatePodcastFromSingleItem(c *gin.Context) {
	var req struct {
		UserID int32 `json:"user_id" binding:"required"`
		ItemID int32 `json:"item_id" binding:"required"`
	}

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

// GetPodcast retrieves a podcast by ID
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

// GetPodcastsByUser retrieves all podcasts for a user
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

// GetPodcastsByStatus retrieves podcasts by status
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

// GetPendingPodcasts retrieves pending podcasts for processing
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

// GetPodcastItems retrieves all items in a podcast
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

// GetPodcastAudio retrieves the audio data for a podcast
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

// GeneratePodcastUploadURL generates a presigned URL for uploading podcast audio
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

// AddItemToPodcast adds an item to a podcast
func (h *PodcastHandler) AddItemToPodcast(c *gin.Context) {
	podcastIDStr := c.Param("id")
	podcastID, err := strconv.ParseInt(podcastIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	var req struct {
		ItemID int32 `json:"item_id" binding:"required"`
		Order  int   `json:"order"`
	}

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

// RemoveItemFromPodcast removes an item from a podcast
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

// UpdatePodcast updates podcast metadata
func (h *PodcastHandler) UpdatePodcast(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid podcast ID"})
		return
	}

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

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

// DeletePodcast deletes a podcast
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
