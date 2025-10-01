package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSSEManager_AddClient(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)

	client := manager.AddClient(userID)

	assert.NotNil(t, client)
	assert.Equal(t, userID, client.UserID)
	assert.NotNil(t, client.Channel)
	assert.Len(t, manager.clients[userID], 1)
}

func TestSSEManager_AddMultipleClientsForSameUser(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)

	client1 := manager.AddClient(userID)
	client2 := manager.AddClient(userID)

	assert.NotNil(t, client1)
	assert.NotNil(t, client2)
	assert.NotEqual(t, client1, client2)
	assert.Len(t, manager.clients[userID], 2)
}

func TestSSEManager_RemoveClient(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)

	client := manager.AddClient(userID)
	assert.Len(t, manager.clients[userID], 1)

	manager.RemoveClient(client)

	// Should remove user entry completely when last client is removed
	_, exists := manager.clients[userID]
	assert.False(t, exists)
}

func TestSSEManager_RemoveOneOfMultipleClients(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)

	client1 := manager.AddClient(userID)
	client2 := manager.AddClient(userID)
	assert.Len(t, manager.clients[userID], 2)

	manager.RemoveClient(client1)

	// Should still have one client
	assert.Len(t, manager.clients[userID], 1)
	assert.Equal(t, client2, manager.clients[userID][0])
}

func TestSSEManager_BroadcastToUser(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)

	client := manager.AddClient(userID)

	message := SSEMessage{
		Event: "test-event",
		Data:  map[string]string{"key": "value"},
	}

	// Broadcast in a goroutine since it blocks
	go manager.BroadcastToUser(userID, message)

	// Receive the message
	select {
	case received := <-client.Channel:
		assert.Equal(t, message.Event, received.Event)
		assert.Equal(t, message.Data, received.Data)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

func TestSSEManager_BroadcastToNonExistentUser(t *testing.T) {
	manager := NewSSEManager()

	message := SSEMessage{
		Event: "test-event",
		Data:  map[string]string{"key": "value"},
	}

	// Should not panic
	manager.BroadcastToUser(999, message)
}

func TestSSEManager_BroadcastToMultipleClients(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)

	client1 := manager.AddClient(userID)
	client2 := manager.AddClient(userID)

	message := SSEMessage{
		Event: "test-event",
		Data:  map[string]string{"key": "value"},
	}

	// Broadcast in a goroutine
	go manager.BroadcastToUser(userID, message)

	// Both clients should receive the message
	select {
	case received := <-client1.Channel:
		assert.Equal(t, message.Event, received.Event)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message on client1")
	}

	select {
	case received := <-client2.Channel:
		assert.Equal(t, message.Event, received.Event)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message on client2")
	}
}

func TestSSEManager_NotifyItemUpdate(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)
	itemID := int32(100)
	status := "completed"

	client := manager.AddClient(userID)

	// Notify in a goroutine
	go manager.NotifyItemUpdate(userID, itemID, &status, "completed")

	// Receive the notification
	select {
	case received := <-client.Channel:
		assert.Equal(t, "item-update", received.Event)

		data, ok := received.Data.(ItemUpdateEvent)
		assert.True(t, ok)
		assert.Equal(t, itemID, data.ItemID)
		assert.Equal(t, status, *data.ProcessingStatus)
		assert.Equal(t, "completed", data.UpdateType)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for item update notification")
	}
}

func TestWriteSSEMessage(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)

	client := manager.AddClient(userID)
	defer manager.RemoveClient(client)

	writer := make(chan string, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start writing in background
	go WriteSSEMessage(ctx, client, writer)

	// Wait for initial connection message
	select {
	case msg := <-writer:
		assert.Contains(t, msg, "connected")
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for connection message")
	}

	// Send a message to the client
	testMessage := SSEMessage{
		Event: "test-event",
		Data:  map[string]string{"test": "data"},
	}
	client.Channel <- testMessage

	// Should receive formatted SSE message
	select {
	case msg := <-writer:
		assert.Contains(t, msg, "event: test-event")
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for event line")
	}

	select {
	case msg := <-writer:
		assert.Contains(t, msg, "data:")
		assert.Contains(t, msg, "test")
		assert.Contains(t, msg, "data")
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for data line")
	}

	// Cancel context and verify it stops
	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestWriteSSEMessage_Keepalive(t *testing.T) {
	manager := NewSSEManager()
	userID := int32(1)

	client := manager.AddClient(userID)
	defer manager.RemoveClient(client)

	writer := make(chan string, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start writing in background
	go WriteSSEMessage(ctx, client, writer)

	// Wait for initial connection message
	select {
	case msg := <-writer:
		assert.Contains(t, msg, "connected")
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for connection message")
	}

	// Wait for context to be done
	<-ctx.Done()
}

func TestSSEManager_ConcurrentAccess(t *testing.T) {
	manager := NewSSEManager()

	// Add clients concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		userID := int32(i)
		go func(uid int32) {
			client := manager.AddClient(uid)
			time.Sleep(10 * time.Millisecond)
			manager.RemoveClient(client)
			done <- true
		}(userID)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have no clients left
	assert.Empty(t, manager.clients)
}
