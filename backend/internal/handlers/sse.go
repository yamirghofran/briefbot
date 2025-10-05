package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yamirghofran/briefbot/internal/services"
)

// StreamItemUpdates godoc
// @Summary      Stream item processing updates
// @Description  Server-Sent Events endpoint for real-time item processing updates
// @Tags         items
// @Produce      text/event-stream
// @Param        userID  path      int  true  "User ID"
// @Success      200     {string}  string "SSE stream"
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /items/user/{userID}/stream [get]
func (h *Handler) StreamItemUpdates(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
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
