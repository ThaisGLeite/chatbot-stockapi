package main

import (
	"chatbot/redis"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	// Connect to the NATS server
	nc, _ := nats.Connect(nats.DefaultURL)
	_ = redis.GetRedisClient()

	// Simple Async Subscriber
	_, _ = nc.Subscribe("stock_quotes", func(m *nats.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
	})

	// Serving static files
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	log.Println("Serving on localhost:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
