package main

import (
	"chatbot/handle"
	"chatbot/redis"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the NATS server
	nats.Connect(os.Getenv("NATS_URL"))

	// Create redis client
	redis.InitializeRedisClient()

	handle.StaticFilesHandler()
	http.HandleFunc("/login", handle.LoginHandler)
	http.HandleFunc("/register", handle.RegisterHandler)
	http.HandleFunc("/createChatroom", handle.CreateChatroomHandler)
	http.HandleFunc("/sendMessage", handle.SendMessageHandler)
	http.HandleFunc("/retrieveMessages", handle.RetrieveMessagesHandler)
	http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
