package handle_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"chatbot/handle"
	"chatbot/model"
	"chatbot/natsclient"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

type mockNatsClient struct {
	subCallbacks map[string]nats.MsgHandler
}

func (m *mockNatsClient) Connect() error {
	return nil
}

func (m *mockNatsClient) Close() {}

func (m *mockNatsClient) Publish(subject string, data []byte) error {
	return nil
}

func (m *mockNatsClient) Subscribe(subject string, cb nats.MsgHandler) (*nats.Subscription, error) {
	m.subCallbacks[subject] = cb
	return &nats.Subscription{}, nil
}

func (m *mockNatsClient) QueueSubscribe(subject, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	m.subCallbacks[subject] = cb
	return &nats.Subscription{}, nil
}

func (m *mockNatsClient) IsConnected() bool {
	return true
}

type mockWs struct {
	lastMessage map[string]string
}

func (m *mockWs) BroadcastMessage(message []byte, room string) {
	m.lastMessage[room] = string(message)
}

func TestListenStockData(t *testing.T) {
	ctx := context.Background()
	mockNats := &mockNatsClient{
		subCallbacks: make(map[string]nats.MsgHandler),
	}
	mockWebsocket := &mockWs{
		lastMessage: make(map[string]string),
	}
	natsclient.Client = mockNats

	// Pass the mockWebsocket to handle.ListenStockData
	handle.ListenStockData(ctx, mockWebsocket)

	// Simulate receiving a stock_data message
	stockData := model.StockData{
		StockCode:    "apple",
		Price:        100.0,
		ChatroomName: "room1",
	}
	data, _ := json.Marshal(stockData)
	mockNats.subCallbacks["stock_data"](&nats.Msg{Data: data})

	// Give it some time to process
	time.Sleep(time.Second)

	// Check that the correct message was sent to the websocket client
	expectedMessage := "APPLE quote is $100.00 per share"
	assert.Equal(t, expectedMessage, mockWebsocket.lastMessage[stockData.ChatroomName])

	// Simulate receiving a stock_errors message
	mockNats.subCallbacks["stock_errors"](&nats.Msg{Data: []byte("room1 | Error message")})

	// Give it some time to process
	time.Sleep(time.Second)

	// Check that the correct message was sent to the websocket client
	expectedMessage = "Error message"
	assert.Equal(t, expectedMessage, mockWebsocket.lastMessage[stockData.ChatroomName])
}
