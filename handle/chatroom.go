package handle

import (
	"chatbot/model"
	"chatbot/natsclient"
	"chatbot/redis"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// Chatroom represents a chatroom with a unique ID and a list of users in the chatroom
var chatroomCounter int64

// CreateChatroomHandler handles chatroom creation requests
func CreateChatroomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		chatroomName := r.FormValue("chatroomName")

		// Increment chatroomCounter
		chatroomCounter++

		// Store chatroomName with chatroomCounter as ID in Redis
		err := redis.StoreChatroomData(chatroomName)
		if err != nil {
			http.Error(w, "Error creating chatroom", http.StatusInternalServerError)
			fmt.Println("Error creating chatroom", err)
			return
		}

		// Write the chatroom ID to the response
		w.Write([]byte(strconv.FormatInt(chatroomCounter, 10)))
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

// SendMessageHandler handles sending messages
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		chatroomID := r.FormValue("chatroomID")
		username := r.FormValue("username")
		message := r.FormValue("message")

		// Check if message starts with /stock=
		if strings.HasPrefix(message, "/stock=") {
			stockCode := strings.TrimPrefix(message, "/stock=")

			// Post the stock code to a NATS queue
			natsclient.Client.Publish("stock_codes", []byte(stockCode))
			natsclient.Client.Flush()

			if err := natsclient.Client.LastError(); err != nil {
				http.Error(w, "Error posting to NATS queue", http.StatusInternalServerError)
				fmt.Println("Error posting to NATS queue", err)
				return
			}

			return
		}

		// Store message in chatroom
		err := redis.StoreMessageInChatroom(chatroomID, username, message)
		if err != nil {
			http.Error(w, "Error sending message", http.StatusInternalServerError)
			fmt.Println("Error sending message", err)
			return
		}

		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

// RetrieveMessagesHandler handles retrieving all messages from a chatroom
func RetrieveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		chatroomID := r.FormValue("chatroomID")

		// Retrieve all messages from the chatroom
		messages, err := redis.RetrieveChatroomMessages(chatroomID)
		if err != nil {
			http.Error(w, "Error retrieving messages", http.StatusInternalServerError)
			fmt.Println("Error retrieving messages", err)
			return
		}

		// Convert the slice of messages to JSON
		messagesJson, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, "Error converting messages to JSON", http.StatusInternalServerError)
			fmt.Println("Error converting messages to JSON", err)
			return
		}

		// Write the messages to the response
		w.Write(messagesJson)
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func GetAllChatroomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	chatrooms, err := redis.GetAllChatrooms()
	if err != nil {
		http.Error(w, "Error getting chatrooms", http.StatusInternalServerError)
		fmt.Println("Error getting chatrooms", err)
		return
	}

	// Extract only chatroom names
	var chatroomNames []string
	chatroomNames = append(chatroomNames, chatrooms...)

	chatroomsJson, err := json.Marshal(chatroomNames)
	if err != nil {
		http.Error(w, "Error converting chatrooms to JSON", http.StatusInternalServerError)
		fmt.Println("Error converting chatrooms to JSON", err)
		return
	}

	w.Write(chatroomsJson)
}

func ListenStockData() {
	// Create a new ticker that ticks every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Subscribe to the stock_data topic
	_, err := natsclient.Client.Subscribe("stock_data", func(m *nats.Msg) {
		// Unmarshal the stock data from the message
		var stockData model.StockData
		err := json.Unmarshal(m.Data, &stockData)
		if err != nil {
			fmt.Println("Error unmarshalling stock data: ", err)
			return
		}

		// Store the stock data in all chatrooms
		chatrooms, err := redis.GetAllChatrooms()
		if err != nil {
			fmt.Println("Error getting all chatrooms: ", err)
			return
		}

		for _, chatroom := range chatrooms {
			err := redis.StoreStockDataInChatroom(chatroom, stockData)
			if err != nil {
				fmt.Println("Error storing stock data in chatroom: ", err)
				return
			}
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	for range ticker.C {
		natsclient.Client.Flush()
	}
}
