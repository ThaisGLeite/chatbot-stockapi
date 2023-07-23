package natsclient

import (
	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	conn *nats.Conn
}

func NewNATSClient(url string) (*NATSClient, error) {
	if url == "" {
		url = nats.DefaultURL
	}
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NATSClient{conn: conn}, nil
}

func (nc *NATSClient) Close() {
	nc.conn.Close()
}

func (nc *NATSClient) Subscribe(topic string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return nc.conn.Subscribe(topic, handler)
}

func (nc *NATSClient) Publish(topic string, data []byte) error {
	return nc.conn.Publish(topic, data)
}
