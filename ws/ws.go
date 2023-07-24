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
	clients = make(map[*websocket.Conn]bool) // active clients
	mutex   sync.Mutex
)

func Connect(w http.ResponseWriter, r *http.Request) {
	var err error
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Conn.Close()

	// Register new client
	mutex.Lock()
	clients[Conn] = true
	mutex.Unlock()

	for {
		_, msg, err := Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			DeleteClient(Conn)
			break
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", Conn.RemoteAddr(), string(msg))

		// Broadcast message to all clients
		BroadcastMessage(msg)
	}
}

func DeleteClient(conn *websocket.Conn) {
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
}

func BroadcastMessage(msg []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Websocket error:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
