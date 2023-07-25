package ws

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients = make(map[string]map[*websocket.Conn]bool) // active clients categorized by chatroom
	mutex   sync.Mutex
)

func Connect(w http.ResponseWriter, r *http.Request) {
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get chatroom from request URL
	chatroom := r.URL.Query().Get("chatroom")
	if chatroom == "" {
		fmt.Println("No chatroom specified in URL")
		Conn.Close()
		return
	}

	// Register new client
	mutex.Lock()
	if _, ok := clients[chatroom]; !ok {
		clients[chatroom] = make(map[*websocket.Conn]bool)
	}
	clients[chatroom][Conn] = true
	mutex.Unlock()

	// Make a channel to signal connection close
	closeConnChan := make(chan struct{})

	go func() {
		// Read messages from the client
		for {
			_, _, err := Conn.ReadMessage()
			if err != nil {
				closeConnChan <- struct{}{} // Signal to close the connection
				break
			}
		}
	}()

	go func() {
		<-closeConnChan // Wait for signal to close the connection
		// Unregister the client and close the connection.
		mutex.Lock()
		delete(clients[chatroom], Conn)
		mutex.Unlock()
		Conn.Close()
	}()
}

func DeleteClient(conn *websocket.Conn, chatroom string) {
	mutex.Lock()
	if _, ok := clients[chatroom]; ok {
		delete(clients[chatroom], conn)
	}
	mutex.Unlock()
}

func BroadcastMessage(msg []byte, chatroom string) {
	mutex.Lock()
	defer mutex.Unlock()

	if conns, ok := clients[chatroom]; ok {
		for conn := range conns {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("Websocket error:", err)
				conn.Close()
				delete(clients[chatroom], conn)
			}
		}
	}
}
