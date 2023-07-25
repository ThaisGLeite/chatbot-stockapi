package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	Broadcaster = &WsBroadcaster{
		Clients: NewClients(),
	}
)

type MessageBroadcaster interface {
	BroadcastMessage(msg []byte, chatroom string)
}

type WsBroadcaster struct {
	Clients *Clients
}

type Clients struct {
	sync.Mutex
	clients map[string]map[*websocket.Conn]bool
}

func NewClients() *Clients {
	return &Clients{
		clients: make(map[string]map[*websocket.Conn]bool),
	}
}

func (c *Clients) GetClient(chatroom string) (map[*websocket.Conn]bool, bool) {
	// Print all connected clients
	c.PrintClients()

	c.Lock()
	defer c.Unlock()
	client, ok := c.clients[chatroom]
	return client, ok
}

func (c *Clients) SetClient(chatroom string, conn *websocket.Conn) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.clients[chatroom]; !ok {
		c.clients[chatroom] = make(map[*websocket.Conn]bool)
	}
	if _, ok := c.clients[chatroom][conn]; ok {
		return fmt.Errorf("client already exists")
	}
	c.clients[chatroom][conn] = true
	return nil
}

func (c *Clients) DeleteClient(conn *websocket.Conn, chatroom string) {
	c.Lock()
	defer c.Unlock()
	conns, ok := c.clients[chatroom]
	if !ok {
		log.Printf("Attempted to delete a client from non-existent chatroom: %s", chatroom)
		return
	}

	delete(conns, conn)
}

func (w *WsBroadcaster) BroadcastMessage(msg []byte, chatroom string) {

	conns, ok := w.Clients.GetClient(chatroom)
	if !ok {
		log.Printf("Attempted to broadcast a message to non-existent chatroom: %s.", chatroom)
	}

	for conn := range conns {
		if conn == nil {
			continue
		}
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Failed to broadcast message due to WebSocket error: %v", err)
			w.Clients.DeleteClient(conn, chatroom)
		}
	}
}

// Connect upgrades the HTTP server connection to the WebSocket protocol.
func Connect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade HTTP connection to WebSocket protocol: %v", err)
		return
	}

	// Get chatroom from request URL
	chatroom := r.URL.Query().Get("chatroom")
	if chatroom == "" {
		errMessage := "No chatroom specified in URL. Please specify a chatroom as a query parameter."
		log.Print(errMessage)
		// Send the error message to client over the WebSocket connection
		if err := conn.WriteMessage(websocket.TextMessage, []byte(errMessage)); err != nil {
			log.Printf("Failed to send error message to client due to WebSocket error: %v", err)
		}
		return
	}

	// Register new client
	if err := Broadcaster.Clients.SetClient(chatroom, conn); err != nil {
		log.Printf("Failed to set client: %v", err)
		return
	}
	Broadcaster.Clients.PrintClients()
	go handleMessages(conn, chatroom)
}

// handleMessages listens for new messages broadcast to our chatroom.
func handleMessages(conn *websocket.Conn, chatroom string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			// If the error was an unexpected closure or the going away error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
				strings.Contains(err.Error(), "websocket: close 1001 (going away)") {
				// If the error was an unexpected closure, log it and stop reading messages
				log.Printf("Unexpected WebSocket closure detected: %v", err)
				Broadcaster.Clients.DeleteClient(conn, chatroom)
				cancel()
				break
			} else {
				// If the error was something else, log it and continue reading messages
				log.Printf("Error reading WebSocket message: %v", err)

				continue
			}
		}

		// broadcast received messages back to all clients in the chatroom
		if msgType == websocket.TextMessage {
			Broadcaster.BroadcastMessage(msg, chatroom)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}

}

func (c *Clients) PrintClients() {
	c.Lock()
	defer c.Unlock()
	log.Println("Clientes conectados: ")
	for chatroom, conns := range c.clients {
		log.Printf("Chatroom: %s\n", chatroom)
		for conn := range conns {
			log.Printf("Connected client: %v\n", conn.LocalAddr())
		}
	}
}
