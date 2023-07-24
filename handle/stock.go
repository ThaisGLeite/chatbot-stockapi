package handle

import (
	"chatbot/model"
	"chatbot/natsclient"
	"chatbot/ws"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
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
		ws.Conn.WriteMessage(websocket.TextMessage, []byte(m.Data))
		return
	}

	// Create message in the required format
	stockData.StockCode = strings.ToUpper(stockData.StockCode)
	botMessage := fmt.Sprintf("%s quote is $%.2f per share", stockData.StockCode, stockData.Price)
	fmt.Println("Stock data: ", botMessage)

	// Send message to WebSocket
	if err := ws.Conn.WriteMessage(websocket.TextMessage, []byte(botMessage)); err != nil {
		fmt.Println("Error sending message to WebSocket:", err)
	}
}

func getErro(m *nats.Msg) {
	// Send message to WebSocket
	if err := ws.Conn.WriteMessage(websocket.TextMessage, []byte(m.Data)); err != nil {
		fmt.Println("Error from botservices:", err)
	}
}
