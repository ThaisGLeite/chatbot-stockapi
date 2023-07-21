package handle

import (
	"chatbot/redis"
	"net/http"
	"strconv"
	"strings"
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
		err := redis.StoreChatroomData(strconv.FormatInt(chatroomCounter, 10), chatroomName)
		if err != nil {
			http.Error(w, "Error creating chatroom", http.StatusInternalServerError)
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

		// Store message in chatroom
		err := redis.StoreChatroomMessage(chatroomID, username, message)
		if err != nil {
			http.Error(w, "Error sending message", http.StatusInternalServerError)
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
			return
		}

		// Convert the slice of messages to a single string with each message on a new line
		messagesString := strings.Join(messages, "\n")

		// Write the messages to the response
		w.Write([]byte(messagesString))
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
