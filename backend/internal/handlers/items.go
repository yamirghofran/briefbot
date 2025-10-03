package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateItem godoc
// @Summary      Create a new content item
// @Description  Create a new content item from URL with async processing
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        item  body      CreateItemRequest  true  "Item creation request"
// @Success      201   {object}  CreateItemResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /items [post]
func (h *Handler) CreateItem(c *gin.Context) {
	var req CreateItemRequest

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

// GetItem godoc
// @Summary      Get an item by ID
// @Description  Retrieve a content item's information by its ID
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {object}  github_com_yamirghofran_briefbot_internal_db.Item
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /items/{id} [get]
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

// GetItemsByUser godoc
// @Summary      Get items by user
// @Description  Retrieve all content items for a specific user
// @Tags         items
// @Produce      json
// @Param        userID  path      int  true  "User ID"
// @Success      200     {array}   github_com_yamirghofran_briefbot_internal_db.Item
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /items/user/{userID} [get]
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

// GetUnreadItemsByUser godoc
// @Summary      Get unread items by user
// @Description  Retrieve all unread content items for a specific user
// @Tags         items
// @Produce      json
// @Param        userID  path      int  true  "User ID"
// @Success      200     {array}   github_com_yamirghofran_briefbot_internal_db.Item
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /items/user/{userID}/unread [get]
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

// UpdateItem godoc
// @Summary      Update an item
// @Description  Update a content item's information
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id    path      int                 true  "Item ID"
// @Param        item  body      UpdateItemRequest   true  "Item update request"
// @Success      200   {object}  MessageResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /items/{id} [put]
func (h *Handler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req UpdateItemRequest

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

// MarkItemAsRead godoc
// @Summary      Mark item as read
// @Description  Mark a content item as read
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {object}  github_com_yamirghofran_briefbot_internal_db.Item
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /items/{id}/read [patch]
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

// ToggleItemReadStatus godoc
// @Summary      Toggle item read status
// @Description  Toggle a content item's read/unread status
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {object}  github_com_yamirghofran_briefbot_internal_db.Item
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /items/{id}/toggle-read [patch]
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

// PatchItem godoc
// @Summary      Patch an item
// @Description  Partially update a content item's information
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "Item ID"
// @Param        item  body      PatchItemRequest   true  "Item patch request"
// @Success      200   {object}  github_com_yamirghofran_briefbot_internal_db.Item
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /items/{id} [patch]
func (h *Handler) PatchItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req PatchItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.itemService.PatchItem(c.Request.Context(), int32(id), req.Title, req.Summary, req.Tags, req.Authors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem godoc
// @Summary      Delete an item
// @Description  Delete a content item from the system
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /items/{id} [delete]
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

// GetItemProcessingStatus godoc
// @Summary      Get item processing status
// @Description  Retrieve the processing status of a content item
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {object}  ItemProcessingStatusResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /items/{id}/status [get]
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

// GetItemsByProcessingStatus godoc
// @Summary      Get items by processing status
// @Description  Retrieve content items filtered by their processing status
// @Tags         items
// @Produce      json
// @Param        status  query     string  false  "Processing status (pending, processing, completed, failed)"  default(pending)
// @Success      200     {object}  ItemsByStatusResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /items/status [get]
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
