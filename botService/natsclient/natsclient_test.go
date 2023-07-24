package natsclient_test

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define a mock for NATSClient
type NATSClientMock struct {
	mock.Mock
}

func (m *NATSClientMock) Close() {}

func (m *NATSClientMock) Subscribe(topic string, handler nats.MsgHandler) (*nats.Subscription, error) {
	args := m.Called(topic, handler)
	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (m *NATSClientMock) Publish(topic string, data []byte) error {
	args := m.Called(topic, data)
	return args.Error(0)
}

func TestNATSClient(t *testing.T) {
	// Initialize mock client and set expectations
	client := new(NATSClientMock)
	client.On("Subscribe", mock.Anything, mock.Anything).Return(&nats.Subscription{}, nil)
	client.On("Publish", mock.Anything, mock.Anything).Return(nil)

	topic := "test_topic"
	message := []byte("Hello, NATS!")

	// Set up a subscriber to the topic.
	_, err := client.Subscribe(topic, func(msg *nats.Msg) {})
	assert.NoError(t, err, "Error subscribing to topic")

	// Publish a message to the topic.
	err = client.Publish(topic, message)
	assert.NoError(t, err, "Error publishing to topic")

	// Ensure that all expectations were met
	client.AssertExpectations(t)
}
