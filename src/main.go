package main

import (
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to the NATS server
	nc, _ := nats.Connect(nats.DefaultURL)

	// Simple Async Subscriber
	_, _ = nc.Subscribe("stock_quotes", func(m *nats.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
	})

	// Serving static files
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	log.Println("Serving on localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
