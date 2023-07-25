package stock_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	natsclient "botService/natsclient"
	stock "botService/stock"
)

var _ natsclient.NATSClientInterface = (*NATSClientMock)(nil)

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

func TestHandleRequest(t *testing.T) {
	mockClient := new(NATSClientMock)
	mockClient.On("Publish", mock.Anything, mock.Anything).Return(nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Symbol,Date,Time,Open,High,Low,Close,Volume
AAPL.US,2023-07-23,22:00:10,147.27,149.04,147.00,148.56,14077800`))
	}))
	defer ts.Close()

	logger := logrus.New()

	handler := stock.NewStockDataHandler(mockClient, logger, ts.URL)

	testMsg := &nats.Msg{
		Subject: "subject",
		Reply:   "reply",
		Data:    []byte(`{"stockCode": "AAPL", "chatroomName": "room1"}`),
		Sub:     nil,
	}

	handler.HandleRequest(testMsg)

	mockClient.AssertExpectations(t)
}

func TestGetStockDataHTTPHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Symbol,Date,Time,Open,High,Low,Close,Volume
AAPL.US,2023-07-23,22:00:10,147.27,149.04,147.00,148.56,14077800`))
	}))
	defer ts.Close()

	mockClient := new(NATSClientMock)
	mockClient.On("Publish", mock.Anything, mock.Anything).Return(nil)

	logger := logrus.New()

	handler := stock.NewStockDataHandler(mockClient, logger, ts.URL)

	req, err := http.NewRequest("GET", "/?stock_code=AAPL&chatroom_name=room1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.GetStockDataHTTPHandler)

	handlerFunc.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var sd stock.StockData
	err = json.Unmarshal(rr.Body.Bytes(), &sd)
	if err != nil {
		t.Fatal(err)
	}

	expected := stock.StockData{
		StockCode:    "AAPL",
		Price:        148.56,
		ChatroomName: "room1",
	}

	assert.Equal(t, expected, sd)

	mockClient.AssertExpectations(t)
}
