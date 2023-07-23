package model

// Define a struct to hold your message data
type ChatMessage struct {
	Username     string `json:"username"`
	Message      string `json:"message"`
	Timestamp    int64  `json:"timestamp"`
	ChatroomName string `json:"chatroom_name"`
}

// Define a struct to hold stock data
type StockData struct {
	StockCode    string  `json:"stockCode"`
	Price        float64 `json:"price"`
	ChatroomName string  `json:"chatroom_name"`
}
