package natsclient

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

type NatsClient interface {
	Connect() error
	Close()
	IsConnected() bool
	Publish(subject string, data []byte) error
	Subscribe(subject string, cb nats.MsgHandler) (*nats.Subscription, error)
	QueueSubscribe(subject, queue string, cb nats.MsgHandler) (*nats.Subscription, error)
}

type NatsConn struct {
	Conn *nats.Conn
}

// Client is the shared NATS client instance.
var Client NatsClient

// Connect establishes a connection to NATS server using the connection string specified by the NATS_URL environment variable.
func (nc *NatsConn) Connect() error {
	var err error
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return fmt.Errorf("the NATS_URL environment variable is not set")
	}
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS at %s: %v", natsURL, err)
	}
	nc.Conn = conn
	Client = nc
	return nil
}

// Close disconnects the shared NATS client instance.
func (nc *NatsConn) Close() {
	nc.Conn.Close()
}

func (nc *NatsConn) Subscribe(subject string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return nc.Conn.Subscribe(subject, cb)
}

func (c *NatsConn) Publish(subject string, data []byte) error {
	return c.Conn.Publish(subject, data)
}

func (nc *NatsConn) IsConnected() bool {
	return nc.Conn.IsConnected()
}

func (nc *NatsConn) QueueSubscribe(subject, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return nc.Conn.QueueSubscribe(subject, queue, cb)
}
