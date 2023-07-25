package handle

import (
	"chatbot/model"
	"chatbot/natsclient"
	"chatbot/ws"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

func ListenStockData(ctx context.Context) {
	log.Println("Listening to stock data")

	// Subscribe to the "stock_data" subject
	stockSub, err := natsclient.Client.Subscribe("stock_data", getStock)
	if err != nil {
		log.Fatal("Error subscribing to stock_data: ", err)
	}

	// Subscribe to the "stock_errors" subject
	errSub, err := natsclient.Client.Subscribe("stock_errors", getError)
	if err != nil {
		log.Fatal("Error subscribing to stock_errors: ", err)
	}

	go func() {
		<-ctx.Done()
		stockSub.Unsubscribe()
		errSub.Unsubscribe()
	}()
}

func getStock(m *nats.Msg) {
	var stockData model.StockData
	err := json.Unmarshal(m.Data, &stockData)
	if err != nil {
		log.Printf("getStock: Error unmarshalling stock data: %v", err)
		ws.BroadcastMessage(m.Data, stockData.ChatroomName)
		return
	}

	// Create message in the required format
	stockData.StockCode = strings.ToUpper(stockData.StockCode)
	botMessage := fmt.Sprintf("%s quote is $%.2f per share", stockData.StockCode, stockData.Price)
	log.Printf("Stock data: %s", botMessage)

	// Send message to WebSocket
	ws.BroadcastMessage([]byte(botMessage), stockData.ChatroomName)
}

func getError(m *nats.Msg) {
	msgStr := string(m.Data)
	msgParts := strings.SplitN(msgStr, " | ", 2)
	if len(msgParts) < 2 {
		log.Println("getError: received error message in unexpected format")
		return
	}

	chatroomName := msgParts[0]
	ws.BroadcastMessage([]byte(msgParts[1]), chatroomName)
}
