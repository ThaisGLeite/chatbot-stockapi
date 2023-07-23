package main

import (
	"chatbot/handle"
	"chatbot/natsclient"
	"chatbot/redis"
	"fmt"
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
	if !natsclient.Client.IsConnected() {
		fmt.Println("Not connected to NATS server")
	}
	// Create redis client
	redis.InitializeRedisClient()

	// Listen to stock data and update redis cache with new data
	// Run in a goroutine
	go handle.ListenStockData()

	fmt.Println("Server started on port 8080")

	handle.StaticFilesHandler()
	http.HandleFunc("/login", handle.LoginHandler)
	http.HandleFunc("/register", handle.RegisterHandler)
	http.HandleFunc("/createChatroom", handle.CreateChatroomHandler)
	http.HandleFunc("/sendMessage", handle.SendMessageHandler)
	http.HandleFunc("/retrieveMessages", handle.RetrieveMessagesHandler)
	http.HandleFunc("/getAllChatrooms", handle.GetAllChatroomsHandler)
	http.HandleFunc("/checkChatroomExist", handle.CheckChatroomExistHandler)
	http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server stopped")
	defer natsclient.Close()
	defer redis.Close()
}
