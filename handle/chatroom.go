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

const (
	stockCommandPrefix = "/stock="
)

var chatroomCounter int64

// Error messages
var (
	ErrInvalidRequest         = "Invalid request method"
	ErrChatroomExistenceCheck = "Error checking chatroom existence"
	ErrCreatingChatroom       = "Error creating chatroom"
	ErrSendingMessage         = "Error sending message"
	ErrRetrieveMessages       = "Error retrieving messages"
	ErrGettingChatrooms       = "Error getting chatrooms"
)

// checkChatroomExistHandler checks if a chatroom exists
func CheckChatroomExistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrInvalidRequest, http.StatusMethodNotAllowed)
		return
	}

	if chatroomName, err := getFormValue(r, "chatroomName"); err == nil {
		if exists, err := redis.CheckChatroomExist(chatroomName); err == nil {
			fmt.Fprint(w, strconv.FormatBool(exists))
			return
		}
	}
	respondWithError(w, ErrChatroomExistenceCheck)
}

// CreateChatroomHandler handles chatroom creation requests
func CreateChatroomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrInvalidRequest, http.StatusMethodNotAllowed)
		return
	}

	if chatroomName, err := getFormValue(r, "chatroomName"); err == nil {
		chatroomCounter++
		if err := redis.StoreChatroomData(chatroomName); err == nil {
			w.Write([]byte(strconv.FormatInt(chatroomCounter, 10)))
			return
		}
	}
	respondWithError(w, ErrCreatingChatroom)
}

// SendMessageHandler handles sending messages
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrInvalidRequest, http.StatusMethodNotAllowed)
		return
	}

	chatroomName, _ := getFormValue(r, "chatroomName")
	username, _ := getFormValue(r, "username")
	message, _ := getFormValue(r, "message")

	if strings.HasPrefix(message, stockCommandPrefix) {
		stockRequest := model.StockData{
			StockCode:    strings.TrimPrefix(message, stockCommandPrefix),
			ChatroomName: chatroomName,
		}
		publishToNATS(w, stockRequest)
	} else if err := redis.StoreMessageInChatroom(chatroomName, username, message); err != nil {
		respondWithError(w, ErrSendingMessage)
	}
}

// RetrieveMessagesHandler handles retrieving all messages from a chatroom
func RetrieveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrInvalidRequest, http.StatusMethodNotAllowed)
		return
	}

	if chatroomName, err := getFormValue(r, "chatroomName"); err == nil {
		if messages, err := redis.RetrieveChatroomMessages(chatroomName); err == nil {
			respondWithJSON(w, messages)
			return
		}
	}
	respondWithError(w, ErrRetrieveMessages)
}

func GetAllChatroomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, ErrInvalidRequest, http.StatusMethodNotAllowed)
		return
	}

	if chatrooms, err := redis.GetAllChatrooms(); err == nil {
		respondWithJSON(w, chatrooms)
		return
	}
	respondWithError(w, ErrGettingChatrooms)
}

func getFormValue(r *http.Request, key string) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}
	return r.FormValue(key), nil
}

func publishToNATS(w http.ResponseWriter, request model.StockData) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling stock request", err)
		return
	}
	if err := natsclient.Client.Publish("stock_codes", requestBytes); err != nil {
		respondWithError(w, "Error posting to NATS queue")
		fmt.Println("Error posting to NATS queue", err)
		return
	}
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error converting data to JSON", err)
		http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func respondWithError(w http.ResponseWriter, message string) {
	fmt.Println(message)
	http.Error(w, message, http.StatusInternalServerError)
}
