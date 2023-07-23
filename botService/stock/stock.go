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

func (sdh *StockDataHandler) HandleRequest(m *nats.Msg) {
	var stockDataMessage StockData
	err := json.Unmarshal(m.Data, &stockDataMessage)
	if err != nil {
		fmt.Println("Error unmarshalling stock request: ", err)
		return
	}

	stockData, err := callStockAPI(stockDataMessage.StockCode)
	if err != nil {
		fmt.Println("Error getting stock data: ", err)
		return
	}

	stockData.ChatroomName = stockDataMessage.ChatroomName

	stockDataJSON, err := json.Marshal(stockData)
	if err != nil {
		fmt.Println("Error encoding stock data to JSON: ", err)
		return
	}
	fmt.Println("Sending stock data: ", string(stockDataJSON))
	sdh.natsClient.Publish("stock_data."+stockDataMessage.ChatroomName, stockDataJSON)
}

func (sdh *StockDataHandler) GetStockDataHTTPHandler(w http.ResponseWriter, r *http.Request) {
	stockCode := r.URL.Query().Get("stock_code")
	chatroomID := r.URL.Query().Get("chatroom_id")
	if stockCode == "" {
		http.Error(w, "Missing stock_code", http.StatusBadRequest)
		fmt.Println("Missing stock_code")
		return
	}

	stockData, err := callStockAPI(stockCode)
	if err != nil {
		http.Error(w, "Error getting stock data", http.StatusInternalServerError)
		fmt.Println("Error getting stock data: ", err)
		return
	}

	stockData.ChatroomName = chatroomID

	stockDataJSON, err := json.Marshal(stockData)
	if err != nil {
		http.Error(w, "Error encoding stock data to JSON", http.StatusInternalServerError)
		fmt.Println("Error encoding stock data to JSON: ", err)
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
