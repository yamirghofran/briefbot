package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// SSEMessage represents a message sent via Server-Sent Events
type SSEMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// ItemUpdateEvent represents an item update notification
type ItemUpdateEvent struct {
	ItemID           int32   `json:"item_id"`
	ProcessingStatus *string `json:"processing_status"`
	UpdateType       string  `json:"update_type"` // "created", "updated", "completed", "failed"
}

// SSEClient represents a single SSE connection
type SSEClient struct {
	UserID  int32
	Channel chan SSEMessage
}

// SSEManager manages SSE connections for multiple users
type SSEManager struct {
	clients map[int32][]*SSEClient
	mu      sync.RWMutex
}

// NewSSEManager creates a new SSE manager
func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[int32][]*SSEClient),
	}
}

// AddClient adds a new SSE client for a user
func (m *SSEManager) AddClient(userID int32) *SSEClient {
	m.mu.Lock()
	defer m.mu.Unlock()

	client := &SSEClient{
		UserID:  userID,
		Channel: make(chan SSEMessage, 10), // Buffer to prevent blocking
	}

	m.clients[userID] = append(m.clients[userID], client)
	log.Printf("SSE: Added client for user %d (total: %d)", userID, len(m.clients[userID]))

	return client
}

// RemoveClient removes an SSE client
func (m *SSEManager) RemoveClient(client *SSEClient) {
	m.mu.Lock()
	defer m.mu.Unlock()

	userClients := m.clients[client.UserID]
	for i, c := range userClients {
		if c == client {
			// Remove client from slice
			m.clients[client.UserID] = append(userClients[:i], userClients[i+1:]...)
			close(client.Channel)
			log.Printf("SSE: Removed client for user %d (remaining: %d)", client.UserID, len(m.clients[client.UserID]))
			break
		}
	}

	// Clean up empty user entries
	if len(m.clients[client.UserID]) == 0 {
		delete(m.clients, client.UserID)
	}
}

// BroadcastToUser sends a message to all clients for a specific user
func (m *SSEManager) BroadcastToUser(userID int32, message SSEMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clients := m.clients[userID]
	if len(clients) == 0 {
		return
	}

	log.Printf("SSE: Broadcasting to %d clients for user %d: %s", len(clients), userID, message.Event)

	for _, client := range clients {
		select {
		case client.Channel <- message:
			// Message sent successfully
		case <-time.After(1 * time.Second):
			// Client is not reading, skip (prevents blocking)
			log.Printf("SSE: Client for user %d is slow, skipping message", userID)
		}
	}
}

// NotifyItemUpdate notifies all clients for a user about an item update
func (m *SSEManager) NotifyItemUpdate(userID int32, itemID int32, processingStatus *string, updateType string) {
	event := ItemUpdateEvent{
		ItemID:           itemID,
		ProcessingStatus: processingStatus,
		UpdateType:       updateType,
	}

	message := SSEMessage{
		Event: "item-update",
		Data:  event,
	}

	m.BroadcastToUser(userID, message)
}

// WriteSSEMessage writes a SSE message to a channel with proper formatting
func WriteSSEMessage(ctx context.Context, client *SSEClient, writer chan<- string) {
	// Send initial comment to establish connection
	writer <- ": connected\n\n"

	// Send keepalive ping every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Send keepalive comment
			writer <- ": keepalive\n\n"
		case msg, ok := <-client.Channel:
			if !ok {
				return
			}

			// Marshal data to JSON
			data, err := json.Marshal(msg.Data)
			if err != nil {
				log.Printf("SSE: Error marshaling message: %v", err)
				continue
			}

			// Format SSE message
			if msg.Event != "" {
				writer <- fmt.Sprintf("event: %s\n", msg.Event)
			}
			writer <- fmt.Sprintf("data: %s\n\n", string(data))
		}
	}
}
