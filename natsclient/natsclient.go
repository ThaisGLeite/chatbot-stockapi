package natsclient

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

var natsclient *nats.Conn

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Connect() {
	natsclient, err := nats.Connect(os.Getenv("NATS_URL"))
	failOnError(err, "Failed to connect to NATS")

	// Simple Async Subscriber
	_, _ = natsclient.Subscribe("stock_quotes", func(m *nats.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
	})
}

func Close() {
	natsclient.Close()
}
