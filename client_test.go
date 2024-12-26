package main

//
//import (
//	"fmt"
//	"github.com/gorilla/websocket"
//	"log"
//	"net/http"
//	"sync"
//	"testing"
//	"time"
//)
//
//// Test client structure
//type TestClient struct {
//	conn *websocket.Conn
//	wg   *sync.WaitGroup
//	id   int
//}
//
//// Connect a client to the WebSocket server
//func (tc *TestClient) connect(serverURL string, token string) error {
//	header := http.Header{}
//	header.Set("token", token)
//
//	conn, _, err := websocket.DefaultDialer.Dial(serverURL, header)
//	if err != nil {
//		return err
//	}
//
//	tc.conn = conn
//	return nil
//}
//
//// Send and receive messages
//func (tc *TestClient) sendAndReceive(messageCount int) {
//	defer tc.wg.Done()
//	for i := 0; i < messageCount; i++ {
//		message := fmt.Sprintf("Client %d - Message %d", tc.id, i+1)
//
//		// Send a message
//		if err := tc.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
//			log.Printf("Error sending message from client %d: %v", tc.id, err)
//			return
//		}
//
//		// Receive a message (echoed back)
//		_, response, err := tc.conn.ReadMessage()
//		if err != nil {
//			log.Printf("Error receiving message for client %d: %v", tc.id, err)
//			return
//		}
//
//		log.Printf("Client %d received: %s", tc.id, response)
//	}
//}
//
//// Benchmark Test for 2000 Clients Sending 10 Messages Each
//func TestWebSocketServer(t *testing.T) {
//	serverURL := "ws://localhost:8080/ws"
//	clientCount := 50
//	messageCount := 10
//
//	// WaitGroup to wait for all clients to finish
//	var wg sync.WaitGroup
//	wg.Add(clientCount)
//
//	start := time.Now()
//
//	// Launch clients
//	for i := 1; i <= clientCount; i++ {
//		go func(id int) {
//			client := &TestClient{id: id, wg: &wg}
//			token := fmt.Sprintf("token-%d", id)
//
//			// Connect to the server
//			if err := client.connect(serverURL, token); err != nil {
//				t.Fatalf("Client %d failed to connect: %v", id, err)
//			}
//			defer client.conn.Close()
//
//			// Send and receive messages
//			client.sendAndReceive(messageCount)
//		}(i)
//	}
//
//	// Wait for all clients to complete
//	wg.Wait()
//
//	elapsed := time.Since(start)
//	t.Logf("Test completed: 2000 clients each sent %d messages in %s", messageCount, elapsed)
//}
