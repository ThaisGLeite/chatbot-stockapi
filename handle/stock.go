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

func ListenStockData(ctx context.Context, broadcaster ws.MessageBroadcaster) {
	log.Println("Listening to stock data")

	// Subscribe to the "stock_data" subject
	stockSub, err := natsclient.Client.Subscribe("stock_data", getStock(broadcaster))
	if err != nil {
		log.Fatal("Error subscribing to stock_data: ", err)
	}

	// Subscribe to the "stock_errors" subject
	errSub, err := natsclient.Client.Subscribe("stock_errors", getError(broadcaster))
	if err != nil {
		log.Fatal("Error subscribing to stock_errors: ", err)
	}

	go func() {
		<-ctx.Done()
		stockSub.Unsubscribe()
		errSub.Unsubscribe()
	}()
}

func getStock(broadcaster ws.MessageBroadcaster) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		var stockData model.StockData
		err := json.Unmarshal(m.Data, &stockData)
		if err != nil {
			log.Printf("getStock: Error unmarshalling stock data: %v", err)
			broadcaster.BroadcastMessage(m.Data, stockData.ChatroomName)
			return
		}
		// Create message in the required format
		stockData.StockCode = strings.ToUpper(stockData.StockCode)
		botMessage := fmt.Sprintf("%s quote is $%.2f per share", stockData.StockCode, stockData.Price)

		// Send message to WebSocket
		broadcaster.BroadcastMessage([]byte(botMessage), stockData.ChatroomName)
	}
}

func getError(broadcaster ws.MessageBroadcaster) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		msgStr := string(m.Data)
		msgParts := strings.SplitN(msgStr, " | ", 2)
		if len(msgParts) < 2 {
			log.Println("getError: received error message in unexpected format")
			return
		}

		chatroomName := msgParts[0]
		broadcaster.BroadcastMessage([]byte(msgParts[1]), chatroomName)
	}
}
