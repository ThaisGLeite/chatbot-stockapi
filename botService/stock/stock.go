package stock

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	natsclient "botService/natsclient"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// StockData represents the data for a particular stock
type StockData struct {
	StockCode    string  `json:"stockcode"`
	Price        float64 `json:"price"`
	ChatroomName string  `json:"chatroom_name"`
}

// StockDataHandler handles the stock data
type StockDataHandler struct {
	natsClient natsclient.NATSClientInterface
	logger     *logrus.Logger
	baseURL    string
}

type StockDataHandlerInterface interface {
	HandleRequest(m *nats.Msg)
}

// NewStockDataHandler creates a new StockDataHandler
func NewStockDataHandler(nc natsclient.NATSClientInterface, logger *logrus.Logger, baseURL string) *StockDataHandler {
	return &StockDataHandler{natsClient: nc, logger: logger, baseURL: baseURL}
}

// HandleRequest handles the stock request from the chatroom
func (sdh *StockDataHandler) HandleRequest(m *nats.Msg) {
	defer sdh.recoverPanic("HandleRequest")

	var stockDataMessage StockData
	err := json.Unmarshal(m.Data, &stockDataMessage)
	if err != nil {
		sdh.handleError("Error unmarshalling stock request", err, "stock_errors")
		return
	}

	stockData, err := sdh.FetchStockData(stockDataMessage.StockCode)
	if err != nil {
		sdh.handleError("Error getting stock data", err, "stock_errors")
		return
	}

	stockData.ChatroomName = stockDataMessage.ChatroomName
	stockDataJSON, err := json.Marshal(stockData)
	if err != nil {
		sdh.handleError("Error encoding stock data to JSON", err, "stock_errors")
		return
	}

	err = sdh.natsClient.Publish("stock_data", stockDataJSON)
	if err != nil {
		sdh.handleError("Error publishing stock data", err, "stock_errors")
		return
	}

	logrus.Info(stockData)
}

// GetStockDataHTTPHandler is the HTTP handler for fetching stock data
func (sdh *StockDataHandler) GetStockDataHTTPHandler(w http.ResponseWriter, r *http.Request) {
	defer sdh.recoverPanic("GetStockDataHTTPHandler")

	stockCode := r.URL.Query().Get("stock_code")
	chatroomName := r.URL.Query().Get("chatroom_name")
	if stockCode == "" {
		sdh.handleHTTPError(w, "Missing stock_code", http.StatusBadRequest, "stock_errors")
		return
	}

	stockData, err := sdh.FetchStockData(stockCode)
	if err != nil {
		sdh.handleHTTPError(w, "Error getting stock data", http.StatusInternalServerError, "stock_errors")
		return
	}

	stockData.ChatroomName = chatroomName

	stockDataJSON, err := json.Marshal(stockData)
	if err != nil {
		sdh.handleHTTPError(w, "Error encoding stock data to JSON", http.StatusInternalServerError, "stock_errors")
		return
	}

	err = sdh.natsClient.Publish("stock_data", stockDataJSON)
	if err != nil {
		sdh.handleError("Error publishing stock data", err, "stock_errors")
		return
	}

	json.NewEncoder(w).Encode(stockData)
}

// recoverPanic recovers from a panic and logs the error
func (sdh *StockDataHandler) recoverPanic(source string) {
	if r := recover(); r != nil {
		sdh.logger.Errorf("Recovered from panic in %s: %v", source, r)
	}
}

// handleError handles an error by logging it and publishing it on NATS
func (sdh *StockDataHandler) handleError(message string, err error, topic string) {
	sdh.logger.Errorf("%s: %v", message, err)
	errMsg := fmt.Sprintf("%s: %v", message, err)
	sdh.natsClient.Publish(topic, []byte(errMsg))
}

// handleHTTPError is like handleError but also sends an HTTP error response
func (sdh *StockDataHandler) handleHTTPError(w http.ResponseWriter, message string, statusCode int, topic string) {
	http.Error(w, message, statusCode)
	sdh.handleError(message, nil, topic)
}

// fetchStockData fetches the stock data for a particular stock code
func (sdh *StockDataHandler) FetchStockData(stockCode string) (*StockData, error) {

	// Make a GET request to the stock API.
	url := fmt.Sprintf("%s/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", sdh.baseURL, stockCode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the body of the response.
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the CSV data.
	records, err := csv.NewReader(strings.NewReader(buf.String())).ReadAll()
	if err != nil {
		return nil, err
	}

	// Check that the CSV data is in the expected format.
	if len(records) < 2 || len(records[1]) < 7 {
		return nil, errors.New("malformed CSV")
	}

	// Extract the closing price and convert it to a float.
	price, err := ParseCSV(records)
	if err != nil {
		return nil, err
	}
	// Return the stock data.
	return &StockData{StockCode: stockCode, Price: price}, nil
}

// ParseCSV parses the CSV data into a float64.
func ParseCSV(records [][]string) (float64, error) {
	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		for _, cell := range row {
			if strings.Contains(cell, "N/D") {
				return 0, errors.New("no data available")
			}
		}
	}

	priceStr := records[1][4] // Closing price is at index 4
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
