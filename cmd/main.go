package main

import (
	"chatbot/handle"
	"chatbot/natsclient"
	"chatbot/redis"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the NATS server
	natsclient.Connect()

	// Create redis client
	redis.InitializeRedisClient()

	handle.StaticFilesHandler()
	http.HandleFunc("/login", handle.LoginHandler)
	http.HandleFunc("/register", handle.RegisterHandler)
	http.HandleFunc("/createChatroom", handle.CreateChatroomHandler)
	http.HandleFunc("/sendMessage", handle.SendMessageHandler)
	http.HandleFunc("/retrieveMessages", handle.RetrieveMessagesHandler)
	http.HandleFunc("/getAllChatrooms", handle.GetAllChatroomsHandler)
	http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer natsclient.Close()
	defer redis.Close()
}
