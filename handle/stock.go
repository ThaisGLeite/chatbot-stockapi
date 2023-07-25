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
		fmt.Println("Error unmarshalling stock data: ", err)
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
	// Send the error message to WebSocket
	var stockData model.StockData
	err := json.Unmarshal(m.Data, &stockData)
	if err != nil {
		fmt.Println("Error unmarshalling stock data: ", err)
		ws.BroadcastMessage(m.Data, stockData.ChatroomName)
		return
	}

	ws.BroadcastMessage(m.Data, stockData.ChatroomName)
}
