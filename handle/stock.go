package handle

import (
	"chatbot/model"
	"chatbot/natsclient"
	"chatbot/redis"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

func ListenStockData() {
	// Subscribe to the "stock_data" subject
	_, err := natsclient.Client.Subscribe("stock_data", func(m *nats.Msg) {
		var stockData model.StockData
		err := json.Unmarshal(m.Data, &stockData)
		if err != nil {
			fmt.Println("Error unmarshalling stock data: ", err)
			return
		}
		fmt.Println("Stock data: ", stockData)
		// Create message in the required format
		stockData.StockCode = strings.ToUpper(stockData.StockCode)
		botMessage := fmt.Sprintf("%s quote is $%.2f per share", stockData.StockCode, stockData.Price)
		fmt.Println("Bot message: ", botMessage)
		// Store the bot's message in the specific chatroom
		err = redis.StoreMessageInChatroom(stockData.ChatroomName, "Bot", botMessage)
		if err != nil {
			fmt.Println("Error storing message in chatroom: ", err)
			return
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}
