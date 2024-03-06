package server

import (
	"client-server-application/internal/orderedmap"
	"sync"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

// MockQueueClient for unit tests
type MockQueueClient struct {
	handleMsg func(amqp.Delivery)
}

func (m *MockQueueClient) ConsumeMessages(handleMsg func(amqp.Delivery)) {
	m.handleMsg = handleMsg
}

func (m *MockQueueClient) SimulateMessageReceive(body string) {
	if m.handleMsg != nil {
		m.handleMsg(amqp.Delivery{Body: []byte(body)})
	}
}

func TestServer_HandleCommand(t *testing.T) {
	orderedMap := orderedmap.NewOrderedMap()
	mockQueueClient := &MockQueueClient{}
	commandHandler := NewOrderedMapCommandHandler(orderedMap, &sync.Mutex{})

	srv := NewServer(mockQueueClient, commandHandler)

	// Start server in a goroutine to listen for messages
	go srv.Run()

	// Allow some time for the server to set up message consumption
	time.Sleep(100 * time.Millisecond)

	// Simulate receiving a command
	mockQueueClient.SimulateMessageReceive("addItem(firstKey,testValue)")
	// Allow some time for the command to be processed
	time.Sleep(100 * time.Millisecond)

	// Check the result
	val, exists := orderedMap.Get("firstKey")
	if !exists || val.Value != "testValue" {
		t.Errorf("Expected 'firstKey' to be added with value 'testValue', got %v, exists: %t", val, exists)
	}

	mockQueueClient.SimulateMessageReceive("addItem(secondKey,testValue)")
	time.Sleep(100 * time.Millisecond)
	mockQueueClient.SimulateMessageReceive("deleteItem(secondKey)")
	time.Sleep(100 * time.Millisecond)
	val, exists = orderedMap.Get("secondKey")
	if exists {
		t.Errorf("Expected 'secondKey' is removed after deleted")
	}

}
