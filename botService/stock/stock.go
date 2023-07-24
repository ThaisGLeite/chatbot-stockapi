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

	"botService/natsclient"

	"github.com/nats-io/nats.go"
)

type StockData struct {
	StockCode    string  `json:"stockCode"`
	Price        float64 `json:"price"`
	ChatroomName string  `json:"chatroom_name"`
}

type StockDataHandler struct {
	natsClient *natsclient.NATSClient
}

func NewStockDataHandler(nc *natsclient.NATSClient) *StockDataHandler {
	return &StockDataHandler{natsClient: nc}
}

// HandleRequest handles the stock request from the chatroom.
func (sdh *StockDataHandler) HandleRequest(m *nats.Msg) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in HandleRequest", r)
			// handle error or send error message back
		}
	}()

	// Unmarshal the stock request.
	var stockDataMessage StockData
	err := json.Unmarshal(m.Data, &stockDataMessage)
	if err != nil {
		errorMessage := fmt.Sprintf("Error unmarshalling stock request: %v", err)
		fmt.Println(errorMessage)
		sdh.natsClient.Publish("stock_errors", []byte(errorMessage))
		return
	}

	// Get the stock data from the API
	stockData, err := callStockAPI(stockDataMessage.StockCode)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting stock data: %v", err)
		fmt.Println(errorMessage)
		sdh.natsClient.Publish("stock_errors", []byte(errorMessage))
		return
	}

	// Set the chatroom name on the stock data and marshal it to JSON.
	stockData.ChatroomName = stockDataMessage.ChatroomName
	stockDataJSON, err := json.Marshal(stockData)
	if err != nil {
		errorMessage := fmt.Sprintf("Error encoding stock data to JSON: %v", err)
		fmt.Println(errorMessage)
		sdh.natsClient.Publish("stock_errors", []byte(errorMessage))
		return
	}
	// Publish the stock data to the stock data subject.
	err = sdh.natsClient.Publish("stock_data", stockDataJSON)
	if err != nil {
		errorMessage := fmt.Sprintf("Error publishing stock data: %v", err)
		fmt.Println(errorMessage)
		sdh.natsClient.Publish("stock_errors", []byte(errorMessage))
		return
	}
}

func (sdh *StockDataHandler) GetStockDataHTTPHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in GetStockDataHTTPHandler", r)
			// handle error or send error message back
		}
	}()

	stockCode := r.URL.Query().Get("stock_code")
	chatroomID := r.URL.Query().Get("chatroom_id")
	if stockCode == "" {
		http.Error(w, "Missing stock_code", http.StatusBadRequest)
		errorMessage := "Missing stock_code"
		fmt.Println(errorMessage)
		sdh.natsClient.Publish("stock_errors", []byte(errorMessage))
		return
	}

	stockData, err := callStockAPI(stockCode)
	if err != nil {
		http.Error(w, "Error getting stock data", http.StatusInternalServerError)
		errorMessage := fmt.Sprintf("Error getting stock data: %v", err)
		fmt.Println(errorMessage)
		sdh.natsClient.Publish("stock_errors", []byte(errorMessage))
		return
	}

	stockData.ChatroomName = chatroomID

	stockDataJSON, err := json.Marshal(stockData)
	if err != nil {
		http.Error(w, "Error encoding stock data to JSON", http.StatusInternalServerError)
		errorMessage := fmt.Sprintf("Error encoding stock data to JSON: %v", err)
		fmt.Println(errorMessage)
		sdh.natsClient.Publish("stock_errors", []byte(errorMessage))
		return
	}

	sdh.natsClient.Publish("stock_data", stockDataJSON)

	json.NewEncoder(w).Encode(stockData)
}

func callStockAPI(stockCode string) (*StockData, error) {
	// Make a GET request to the stock API.
	resp, err := http.Get(fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stockCode))
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
	price, err := parseCSV(records)
	if err != nil {
		return nil, err
	}
	// Return the stock data.
	return &StockData{StockCode: stockCode, Price: price}, nil
}

func parseCSV(records [][]string) (float64, error) {
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
