package handle

import (
	"chatbot/model"
	"chatbot/natsclient"
	"chatbot/ws"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

func ListenStockData() {
	fmt.Println("Listening to stock data")
	// Subscribe to the "stock_data" subject
	_, err := natsclient.Client.Subscribe("stock_data", getStock)
	if err != nil {
		log.Fatal(err)
	}
	// Subscribe to the "stock_errors" subject
	_, err = natsclient.Client.Subscribe("stock_errors", getErro)
	if err != nil {
		log.Fatal(err)
	}
}

func getStock(m *nats.Msg) {

	// Unmarshal the message
	var stockData model.StockData
	err := json.Unmarshal(m.Data, &stockData)
	if err != nil {
		fmt.Println("getStock: Error unmarshalling stock data: ", err)
		ws.BroadcastMessage(m.Data, stockData.ChatroomName)
		return
	}

	// Create message in the required format
	stockData.StockCode = strings.ToUpper(stockData.StockCode)
	botMessage := fmt.Sprintf("%s quote is $%.2f per share", stockData.StockCode, stockData.Price)
	fmt.Println("Stock data: ", botMessage)

	// Send message to WebSocket
	ws.BroadcastMessage([]byte(botMessage), stockData.ChatroomName)
}

func getErro(m *nats.Msg) {
	// Convert message data to string
	msgStr := string(m.Data)

	// Split the message by " | " to get the chatroomName
	msgParts := strings.SplitN(msgStr, " | ", 2)

	if len(msgParts) < 2 {
		fmt.Println("Error: received error message in unexpected format")
		return
	}

	// The first part is chatroomName
	chatroomName := msgParts[0]

	// Send the error message to WebSocket
	ws.BroadcastMessage([]byte(msgParts[1]), chatroomName)
}
