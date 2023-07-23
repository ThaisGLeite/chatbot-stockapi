package handle

import (
	"chatbot/model"
	"chatbot/natsclient"
	"chatbot/redis"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Chatroom represents a chatroom with a unique ID and a list of users in the chatroom
var chatroomCounter int64

// CheckChatroomExistHandler checks if a chatroom exists
func CheckChatroomExistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		chatroomName := r.FormValue("chatroomName")

		// Check if chatroom exists in Redis
		exists, err := redis.CheckChatroomExist(chatroomName)
		if err != nil {
			http.Error(w, "Error checking chatroom existence", http.StatusInternalServerError)
			fmt.Println("Error checking chatroom existence", err)
			return
		}

		// Return the result
		if exists {
			fmt.Fprint(w, "true")
		} else {
			fmt.Fprint(w, "false")
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

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

		chatroomName := r.FormValue("chatroomName")
		username := r.FormValue("username")
		message := r.FormValue("message")

		// Check if message starts with /stock=
		if strings.HasPrefix(message, "/stock=") {
			stockCode := strings.TrimPrefix(message, "/stock=")
			stockRequest := model.StockData{
				StockCode:    stockCode,
				ChatroomName: chatroomName,
			}
			requestBytes, err := json.Marshal(stockRequest)
			if err != nil {
				fmt.Println("Error marshalling stock request", err)
				return
			}
			// Post the stock code to a NATS queue
			natsclient.Client.Publish("stock_codes", requestBytes)
			natsclient.Client.Flush()
			if err := natsclient.Client.LastError(); err != nil {
				http.Error(w, "Error posting to NATS queue", http.StatusInternalServerError)
				fmt.Println("Error posting to NATS queue", err)
				return
			}

			return
		}

		// Store message in chatroom
		err := redis.StoreMessageInChatroom(chatroomName, username, message)
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

		chatroomName := r.FormValue("chatroomName")

		// Retrieve all messages from the chatroom
		messages, err := redis.RetrieveChatroomMessages(chatroomName)
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
