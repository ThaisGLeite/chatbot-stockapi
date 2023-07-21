package natsclient

import (
	"log"

	"github.com/nats-io/nats.go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Connect() {
	nc, err := nats.Connect(nats.DefaultURL)
	failOnError(err, "Failed to connect to NATS")

	// Simple Async Subscriber
	_, _ = nc.Subscribe("stock_quotes", func(m *nats.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
	})
}
