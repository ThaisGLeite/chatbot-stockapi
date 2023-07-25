package main

import (
	"chatbot/handle"
	"chatbot/natsclient"
	"chatbot/redis"
	"chatbot/ws"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Connect to the NATS server
	natsclient.Connect()
	if !natsclient.Client.IsConnected() {
		fmt.Println("Not connected to NATS server")
	}
	// Create redis client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := redis.InitializeRedisClient(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Listen to stock data and update redis cache with new data
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	// Run in a goroutine
	go handle.ListenStockData(ctx)

	http.Handle("/", handle.StaticFilesHandler())
	http.HandleFunc("/login", handle.LoginHandler)
	http.HandleFunc("/register", handle.RegisterHandler)
	http.HandleFunc("/createChatroom", handle.CreateChatroomHandler)
	http.HandleFunc("/sendMessage", handle.SendMessageHandler)
	http.HandleFunc("/retrieveMessages", handle.RetrieveMessagesHandler)
	http.HandleFunc("/getAllChatrooms", handle.GetAllChatroomsHandler)
	http.HandleFunc("/checkChatroomExist", handle.CheckChatroomExistHandler)
	http.HandleFunc("/stockUpdates", ws.Connect)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		fmt.Println("Server stopped")
		defer natsclient.Close()
		defer redis.Close()
		log.Fatalf("Failed to start server: %s", err.Error())
		<-ctx.Done()
	}
}
