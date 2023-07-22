package main

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

	"github.com/nats-io/nats.go"
)

type StockData struct {
	StockCode string  `json:"stockCode"`
	Price     float64 `json:"price"`
}

var natsclient *nats.Conn

func main() {
	var err error
	// Connect to a local NATS server
	natsclient, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println("Error connecting to NATS server: ", err)
		return
	}

	defer natsclient.Close()

	// Subscribe to the stock_codes topic
	_, err = natsclient.Subscribe("stock_codes", func(m *nats.Msg) {
		// This function will be executed when a new message arrives on the "stock_codes" topic
		stockCode := string(m.Data)

		stockData, err := callStockAPI(stockCode)
		if err != nil {
			fmt.Println("Error getting stock data: ", err)
			return
		}

		stockDataJSON, err := json.Marshal(stockData)
		if err != nil {
			fmt.Println("Error encoding stock data to JSON: ", err)
			return
		}

		// Publish stock data to a NATS subject
		natsclient.Publish("stock_data", stockDataJSON)
	})

	if err != nil {
		fmt.Println("Error subscribing to stock_codes: ", err)
		return
	}

	// Bind the handler function to the "/get-stock-data" endpoint.
	http.HandleFunc("/get-stock-data", getStockDataHandler)
	// Start the HTTP server on port 3000.
	http.ListenAndServe(":3000", nil)
}

func getStockDataHandler(w http.ResponseWriter, r *http.Request) {
	stockCode := r.URL.Query().Get("stock_code")
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

	stockDataJSON, err := json.Marshal(stockData)
	if err != nil {
		http.Error(w, "Error encoding stock data to JSON", http.StatusInternalServerError)
		fmt.Println("Error encoding stock data to JSON: ", err)
		return
	}

	// Publish stock data to a NATS subject
	natsclient.Publish("stock_data", stockDataJSON)

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
