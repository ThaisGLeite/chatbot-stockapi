package ws

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[string]map[*websocket.Conn]bool) // active clients categorized by chatroom
	mutex   sync.Mutex
)

// Connect upgrades the HTTP server connection to the WebSocket protocol.
func Connect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", http.StatusBadRequest)
		log.Printf("Failed to open a WS connection: %v", err)
		return
	}

	// Get chatroom from request URL
	chatroom := r.URL.Query().Get("chatroom")
	if chatroom == "" {
		http.Error(w, "No chatroom specified in URL", http.StatusBadRequest)
		log.Print("No chatroom specified in URL")
		return
	}

	// Register new client
	mutex.Lock()
	if _, ok := clients[chatroom]; !ok {
		clients[chatroom] = make(map[*websocket.Conn]bool)
	}
	clients[chatroom][conn] = true
	mutex.Unlock()

	go handleMessages(conn, chatroom)
}

// handleMessages listens for new messages broadcast to our chatroom.
func handleMessages(conn *websocket.Conn, chatroom string) {
	defer conn.Close()
	defer DeleteClient(conn, chatroom)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			cancel() // Cancel the context when an error occurs
			break
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// DeleteClient removes a client from the clients map.
func DeleteClient(conn *websocket.Conn, chatroom string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := clients[chatroom]; ok {
		delete(clients[chatroom], conn)
	}
}

// BroadcastMessage sends a message to all clients in the same chatroom.
func BroadcastMessage(msg []byte, chatroom string) {
	mutex.Lock()
	defer mutex.Unlock()
	log.Printf("Active chatrooms and clients:\n")
	for chatroom, conns := range clients {
		log.Printf("Chatroom %s:\n", chatroom)
		for conn := range conns {
			log.Printf("\tClient %v\n", conn.RemoteAddr())
		}
	}
	if conns, ok := clients[chatroom]; ok {
		for conn := range conns {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Websocket error: %v", err)
				delete(clients[chatroom], conn)
			}
		}
	}
}
