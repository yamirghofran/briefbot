package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateItem(c *gin.Context) {
	var req struct {
		UserID *int32  `json:"user_id"`
		URL    *string `json:"url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use async creation - just save the URL and return immediately
	item, err := h.itemService.CreateItemAsync(c.Request.Context(), *req.UserID, *req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the item with pending status
	c.JSON(http.StatusCreated, gin.H{
		"item":              item,
		"message":           "Item created successfully and will be processed in the background",
		"processing_status": item.ProcessingStatus,
	})
}

func (h *Handler) GetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.itemService.GetItem(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) GetItemsByUser(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID32 := int32(userID)
	items, err := h.itemService.GetItemsByUser(c.Request.Context(), &userID32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetUnreadItemsByUser(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID32 := int32(userID)
	items, err := h.itemService.GetUnreadItemsByUser(c.Request.Context(), &userID32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		Title       string   `json:"title"`
		URL         *string  `json:"url"`
		TextContent *string  `json:"text_content"`
		Summary     *string  `json:"summary"`
		Type        *string  `json:"type"`
		Platform    *string  `json:"platform"`
		Tags        []string `json:"tags"`
		Authors     []string `json:"authors"`
		IsRead      *bool    `json:"is_read"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.itemService.UpdateItem(c.Request.Context(), int32(id), req.Title, req.URL, req.TextContent, req.Summary, req.Type, req.Platform, req.Tags, req.Authors, req.IsRead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

func (h *Handler) MarkItemAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = h.itemService.MarkItemAsRead(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated item
	item, err := h.itemService.GetItem(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) ToggleItemReadStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.itemService.ToggleItemReadStatus(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = h.itemService.DeleteItem(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

func (h *Handler) GetItemProcessingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	status, err := h.itemService.GetItemProcessingStatus(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item_id":           status.Item.ID,
		"processing_status": status.Item.ProcessingStatus,
		"is_processing":     status.IsProcessing,
		"is_completed":      status.IsCompleted,
		"is_failed":         status.IsFailed,
		"processing_error":  status.ProcessingError,
	})
}

func (h *Handler) GetItemsByProcessingStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		status = "pending"
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"completed":  true,
		"failed":     true,
	}

	if !validStatuses[status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: pending, processing, completed, failed"})
		return
	}

	items, err := h.itemService.GetItemsByProcessingStatus(c.Request.Context(), &status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"items":  items,
		"count":  len(items),
	})
}
