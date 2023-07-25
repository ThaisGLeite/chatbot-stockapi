package natsclient

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

// Client is the shared NATS client instance.
var Client *nats.Conn

// Connect establishes a connection to NATS server using the connection string specified by the NATS_URL environment variable.
func Connect() error {
	var err error
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return fmt.Errorf("the NATS_URL environment variable is not set")
	}
	Client, err = nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS at %s: %v", natsURL, err)
	}
	return nil
}

// Close disconnects the shared NATS client instance.
func Close() {
	if Client != nil {
		Client.Close()
	}
}
