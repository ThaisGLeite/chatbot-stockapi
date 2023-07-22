package model

// Define a struct to hold your message data
type ChatMessage struct {
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
