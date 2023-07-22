package natsclient

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

var Client *nats.Conn

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Connect() {
	var err error
	Client, err = nats.Connect(os.Getenv("NATS_URL"))
	failOnError(err, "Failed to connect to NATS")
}

func Close() {
	Client.Close()
}
