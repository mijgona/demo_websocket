package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

// Client simulation for the test
type TestClient struct {
	conn *websocket.Conn
	wg   *sync.WaitGroup
}

// Connect to the WebSocket server
func (tc *TestClient) connect(serverURL, token string) error {
	header := http.Header{}
	header.Set("token", token)

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, header)
	if err != nil {
		return err
	}

	tc.conn = conn
	return nil
}

// Listen for messages from the server
func (tc *TestClient) listen(messageCount int) {
	defer tc.wg.Done()
	received := 0

	for {
		_, _, err := tc.conn.ReadMessage()
		if err != nil {
			log.Printf("Client read error: %v", err)
			return
		}
		received++
		if received >= messageCount {
			break
		}
	}
}

// Benchmark function
func BenchmarkWebSocketServer(b *testing.B) {
	serverURL := "ws://localhost:8080/ws"
	clientCount := 1000
	messageCount := 10

	// WaitGroup to ensure all clients finish
	var wg sync.WaitGroup

	// Start server in a separate goroutine
	go func() {
		main()
	}()

	// Give the server time to start
	time.Sleep(2 * time.Second)

	// Create and connect clients
	clients := make([]*TestClient, clientCount)
	for i := 0; i < clientCount; i++ {
		client := &TestClient{wg: &wg}
		err := client.connect(serverURL, fmt.Sprintf("token-%d", i))
		if err != nil {
			b.Fatalf("Failed to connect client %d: %v", i, err)
		}
		clients[i] = client
		wg.Add(1)
		go client.listen(messageCount)
	}

	// Give clients time to connect
	time.Sleep(1 * time.Second)

	// Start sending messages
	startTime := time.Now()
	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Test message %d", i+1)
		log.Printf("Broadcasting: %s", message)

		for _, client := range clients {
			err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				b.Fatalf("Failed to send message to client: %v", err)
			}
		}
	}

	// Wait for all clients to finish
	wg.Wait()
	elapsedTime := time.Since(startTime)

	log.Printf("Benchmark complete: Sent %d messages to %d clients in %v", messageCount, clientCount, elapsedTime)
}
