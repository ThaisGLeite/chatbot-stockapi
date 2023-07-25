package model

// ChatMessage represents a chat message with associated metadata
type ChatMessage struct {
	// Username of the sender
	Username string `json:"username"`
	// Content of the message
	Message string `json:"message"`
	// Timestamp when the message was sent
	Timestamp int64 `json:"timestamp"`
	// Name of the chatroom the message was sent in
	ChatroomName string `json:"chatroom_name"`
}

// StockData represents information about a particular stock
type StockData struct {
	// Code of the stock
	StockCode string `json:"stock_code"`
	// Price of the stock
	Price float64 `json:"price"`
	// Name of the chatroom where the stock information is used
	ChatroomName string `json:"chatroom_name"`
}
