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

func ConnectToNATS() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	failOnError(err, "Failed to connect to NATS")

	return nc
}
