package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	Conn *websocket.Conn
)

func Connect(w http.ResponseWriter, r *http.Request) {
	var err error
	Conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		_, msg, err := Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", Conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println(err)
			return
		}
	}
}
