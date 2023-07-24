package natsclient

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	conn *nats.Conn
}

type ConnectionOptions struct {
	URL            string
	ReconnectDelay time.Duration
	MaxReconnects  int
}

type NATSClientInterface interface {
	Close()
	Subscribe(topic string, handler nats.MsgHandler) (*nats.Subscription, error)
	Publish(topic string, data []byte) error
}

func NewNATSClient(opts ConnectionOptions) (*NATSClient, error) {
	if opts.URL == "" {
		opts.URL = nats.DefaultURL
	}

	// Set up options with callbacks
	options := []nats.Option{
		nats.ReconnectWait(opts.ReconnectDelay),
		nats.MaxReconnects(opts.MaxReconnects),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Printf("Disconnected due to:%s, will attempt reconnects every %.2fs", err, opts.ReconnectDelay.Seconds())
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected [%s]", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Fatalf("Exiting: %v", nc.LastError())
		}),
	}

	// Connect to NATS server with the defined options
	conn, err := nats.Connect(opts.URL, options...)
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
