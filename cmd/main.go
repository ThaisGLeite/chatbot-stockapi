package main

import (
	"chatbot/handle"
	"chatbot/natsclient"
	"chatbot/redis"
	"fmt"
	"log"
	"net/http"
)

func main() {
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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		fmt.Println("Server stopped")
		defer natsclient.Close()
		defer redis.Close()
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
