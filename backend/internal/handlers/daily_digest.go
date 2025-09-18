package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TriggerDailyDigest manually triggers the daily digest email for all users
// @Summary Trigger daily digest for all users
// @Description Manually trigger the daily digest email sending process for all users
// @Tags daily-digest
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "message": "Daily digest sent successfully"
// @Failure 500 {object} map[string]string "error": "Failed to send daily digest"
// @Router /daily-digest/trigger [post]
func (h *Handler) TriggerDailyDigest(c *gin.Context) {
	if h.dailyDigestService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Daily digest service not available"})
		return
	}

	ctx := c.Request.Context()

	// Send daily digest to all users
	if err := h.dailyDigestService.SendDailyDigest(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send daily digest: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Daily digest sent successfully to all users"})
}

// TriggerDailyDigestForUser manually triggers the daily digest email for a specific user
// @Summary Trigger daily digest for specific user
// @Description Manually trigger the daily digest email sending process for a specific user
// @Tags daily-digest
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {object} map[string]string "message": "Daily digest sent successfully"
// @Success 200 {object} map[string]string "message": "No unread items from previous day"
// @Failure 500 {object} map[string]string "error": "Failed to send daily digest"
// @Router /daily-digest/trigger/user/{userID} [post]
func (h *Handler) TriggerDailyDigestForUser(c *gin.Context) {
	if h.dailyDigestService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Daily digest service not available"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := c.Request.Context()

	// Get items first to check if there are any
	items, err := h.dailyDigestService.GetDailyDigestItemsForUser(ctx, int32(userID))
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
	if err := h.dailyDigestService.SendDailyDigestForUser(ctx, int32(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send daily digest: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Daily digest sent successfully to user"})
}
