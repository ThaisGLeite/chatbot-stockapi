package main

import (
	"botService/natsclient"
	"botService/stock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// Mocked interfaces for NATS client and Stock Data Handler
type mockNATSClient struct {
	natsclient.NATSClientInterface
	mockCtrl *gomock.Controller
}

func newMockNATSClient(mockCtrl *gomock.Controller) *mockNATSClient {
	return &mockNATSClient{
		mockCtrl: mockCtrl,
	}
}

type mockStockDataHandler struct {
	*stock.StockDataHandler
	mockCtrl          *gomock.Controller
	HandleRequestFunc func(m *nats.Msg)
}

func newMockStockDataHandler(mockCtrl *gomock.Controller) *mockStockDataHandler {
	return &mockStockDataHandler{
		mockCtrl: mockCtrl,
	}
}

// Implementing Subscribe method for the mock NATS client
func (m *mockNATSClient) Subscribe(topic string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return nil, nil
}

// Implementing HandleRequest method for the mock Stock Data Handler
func (m *mockStockDataHandler) HandleRequest(msg *nats.Msg) {
	if m.HandleRequestFunc != nil {
		m.HandleRequestFunc(msg)
	}
}

// Test for initializeServices function
func TestInitializeServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNatsClient := newMockNATSClient(ctrl)
	mockStockDataHandler := newMockStockDataHandler(ctrl)
	logger = logrus.New()
	err := InitializeServices(mockNatsClient, mockStockDataHandler)
	require.NoError(t, err)
}
