package ws

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHandleMessages(t *testing.T) {
	t.Log("Starting the test")
	server := httptest.NewServer(http.HandlerFunc(Connect))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	u.Scheme = "ws"
	u.RawQuery = "chatroom=test"

	// Open sender connection
	senderConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to open a sender WS connection: %v", err)
	}
	defer senderConn.Close()

	// Open receiver connection
	receiverConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to open a receiver WS connection: %v", err)
	}
	defer receiverConn.Close()

	// Wait for client registration
	time.Sleep(time.Millisecond * 100)

	// Send a message
	err = senderConn.WriteMessage(websocket.TextMessage, []byte("test message"))
	if err != nil {
		t.Fatalf("Failed to write message to connection: %v", err)
	}

	done := make(chan bool)
	go func() {
		// Attempt to read the message from the receiver connection
		_, p, err := receiverConn.ReadMessage()
		if err != nil {
			t.Error("Failed to read message from connection: ", err)
			done <- true
			return
		}

		receivedMsg := string(p)
		if receivedMsg != "test message" {
			t.Errorf("Unexpected message received: got %v, want 'test message'", receivedMsg)
		}

		done <- true
	}()

	select {
	case <-done:
		t.Log("Message read successfully.")
	case <-time.After(time.Second * 5):
		t.Error("Test timeout, possible deadlock.")
	}
}

func TestConnectWithoutChatroom(t *testing.T) {
	t.Log("Starting the test")
	server := httptest.NewServer(http.HandlerFunc(Connect))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	u.Scheme = "ws"
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer conn.Close()

	_, p, err := conn.NextReader()
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	bytes, err := io.ReadAll(p)
	if err != nil {
		t.Fatalf("Failed to read all bytes: %v", err)
	}

	body := string(bytes)
	if !strings.Contains(body, "No chatroom specified in URL") {
		t.Error("Expected an error due to missing chatroom, but got none")
	}
}

func TestBroadcastMessage(t *testing.T) {
	// In a real-world scenario, you would check if the message is correctly sent to all clients in the chatroom.
	// This is a placeholder for the test.
	t.Skip("Test not yet implemented")
}
