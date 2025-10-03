package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TriggerDailyDigest godoc
// @Summary      Trigger daily digest for all users
// @Description  Manually trigger the daily digest email sending process for all users
// @Tags         digest
// @Accept       json
// @Produce      json
// @Success      200  {object}  MessageResponse
// @Failure      500  {object}  ErrorResponse
// @Failure      503  {object}  ErrorResponse
// @Router       /digest/trigger [post]
func (h *Handler) TriggerDailyDigest(c *gin.Context) {
	if h.digestService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Digest service not available"})
		return
	}

	ctx := c.Request.Context()

	// Send daily digest to all users
	if err := h.digestService.SendDailyDigest(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send daily digest: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Daily digest sent successfully to all users"})
}

// TriggerDailyDigestForUser godoc
// @Summary      Trigger daily digest for specific user
// @Description  Manually trigger the daily digest email sending process for a specific user
// @Tags         digest
// @Accept       json
// @Produce      json
// @Param        userID  path      int  true  "User ID"
// @Success      200     {object}  MessageResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Failure      503     {object}  ErrorResponse
// @Router       /digest/trigger/user/{userID} [post]
func (h *Handler) TriggerDailyDigestForUser(c *gin.Context) {
	if h.digestService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Digest service not available"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := c.Request.Context()

	// Get items first to check if there are any
	items, err := h.digestService.GetDailyDigestItemsForUser(ctx, int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily digest items"})
		return
	}

	// If no items, return a friendly message
	if len(items) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No unread items from previous day for this user"})
		return
	}

	// Send the daily digest for this user
	if err := h.digestService.SendDailyDigestForUser(ctx, int32(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send daily digest: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Daily digest sent successfully to user"})
}

// TriggerIntegratedDigest godoc
// @Summary      Trigger integrated digest for all users
// @Description  Manually trigger the integrated digest (podcast generation + email) process for all users
// @Tags         digest
// @Accept       json
// @Produce      json
// @Success      202  {object}  MessageResponse
// @Failure      503  {object}  ErrorResponse
// @Router       /digest/trigger/integrated [post]
func (h *Handler) TriggerIntegratedDigest(c *gin.Context) {
	if h.digestService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Digest service not available"})
		return
	}

	// Start processing in background and return immediately
	go func() {
		ctx := context.Background() // Use background context for async processing
		if err := h.digestService.SendIntegratedDigest(ctx); err != nil {
			log.Printf("Background integrated digest failed: %v", err)
		} else {
			log.Printf("Background integrated digest completed successfully")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Integrated digest processing started in background"})
}

// TriggerIntegratedDigestForUser godoc
// @Summary      Trigger integrated digest for specific user
// @Description  Manually trigger the integrated digest (podcast generation + email) for a specific user
// @Tags         digest
// @Accept       json
// @Produce      json
// @Param        userID  path      int  true  "User ID"
// @Success      202     {object}  MessageResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      503     {object}  ErrorResponse
// @Router       /digest/trigger/integrated/user/{userID} [post]
func (h *Handler) TriggerIntegratedDigestForUser(c *gin.Context) {
	if h.digestService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Digest service not available"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Start processing in background and return immediately
	go func() {
		ctx := context.Background() // Use background context for async processing
		result, err := h.digestService.SendIntegratedDigestForUser(ctx, int32(userID))
		if err != nil {
			log.Printf("Background integrated digest failed for user %d: %v", userID, err)
		} else if result.ItemsCount == 0 {
			log.Printf("No items to process for user %d", userID)
		} else {
			log.Printf("Background integrated digest completed for user %d: emailSent=%v, podcastGenerated=%v",
				userID, result.EmailSent, result.PodcastURL != nil)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Integrated digest processing started for user", "userID": userID})
}
