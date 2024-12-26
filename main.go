package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"sync"
)

type Client struct {
	conn    *websocket.Conn
	send    chan []byte
	receive chan []byte
	token   string
	server  *Server
}

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type Server struct {
	clients    sync.Map
	register   chan *Client
	unregister chan *Client
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var logger *log.Logger

func init() {
	// Set up logging to a file
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error setting up logging: %v\n", err)
		os.Exit(1)
	}
	logger = log.New(file, "", log.LstdFlags)
	log.SetOutput(file)
}

func (server *Server) run() {
	for {
		select {
		case client := <-server.register:
			server.clients.Store(client.token, client)
			logger.Printf("Client connected: %s\n", client.token)

		case client := <-server.unregister:
			if _, ok := server.clients.Load(client.token); ok {
				server.clients.Delete(client.token)
				close(client.send)
				close(client.receive)
			}
			logger.Printf("Client disconnected: %s\n", client.token)
		}
	}
}

func (server *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Printf("Connection upgrade error: %v\n", err)
		return
	}

	client := &Client{
		conn:    conn,
		send:    make(chan []byte),
		receive: make(chan []byte),
		token:   token,
	}

	server.register <- client

	go client.writeMessages()
	go client.processMessages()
}

func (client *Client) processMessages() {
	for message := range client.receive {
		logger.Printf("Processing message for client %s: %s\n", client.token, message)
		client.send <- message // Echo back
	}
}

func (client *Client) writeMessages() {
	defer client.conn.Close()

	for message := range client.send {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			logger.Printf("Write message error: %v\n", err)
			break
		}
	}
}

func main() {
	server := &Server{
		clients:    sync.Map{},
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go server.run()

	http.HandleFunc("/ws", server.handleConnections)

	logger.Println("Server started on :8080")

	go func() {
		for {
			var (
				message = ""
			)
			fmt.Println("Message:  ")
			fmt.Scan(&message)

			m, err := json.Marshal(Message{
				Type:    "user",
				Payload: message,
			})
			if err != nil {
				logger.Printf("parsing error error: %v\n", err)
			}

			lenth := 0
			server.clients.Range(func(key, value any) bool {
				lenth++
				client, _ := value.(Client)
				client.send <- m
				return true
			},
			)

			fmt.Printf("finish: set to %v clients \n", lenth)
		}
	}()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Printf("Server error: %v\n", err)
	}
}
