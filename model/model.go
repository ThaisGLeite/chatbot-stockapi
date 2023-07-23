package model

// Define a struct to hold your message data
type ChatMessage struct {
	Username   string `json:"username"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
	ChatroomID string `json:"chatroom_id"`
}

// Define a struct to hold stock data
type StockData struct {
	StockCode  string  `json:"stockCode"`
	Price      float64 `json:"price"`
	ChatroomID string  `json:"chatroom_id"`
}
