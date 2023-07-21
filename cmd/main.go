package main

import (
	"chatbot/handle"
	"chatbot/register"
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

	// Handle "/register" route
	http.HandleFunc("/register", register.Handler)

	// Handle static files
	handle.Handle()

	log.Println("Serving on localhost:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
