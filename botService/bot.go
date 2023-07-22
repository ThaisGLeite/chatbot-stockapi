package main

import (
	"encoding/json"
	"net/http"
)

type StockData struct {
	StockCode string  `json:"stockCode"`
	Price     float64 `json:"price"`
	// Include other relevant stock data fields
}

func main() {
	http.HandleFunc("/get-stock-data", getStockDataHandler)
	http.ListenAndServe(":8080", nil)
}

func getStockDataHandler(w http.ResponseWriter, r *http.Request) {
	stockCode := r.URL.Query().Get("stock_code")
	if stockCode == "" {
		http.Error(w, "Missing stock_code", http.StatusBadRequest)
		return
	}

	stockData, err := callStockAPI(stockCode)
	if err != nil {
		http.Error(w, "Error getting stock data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stockData)
}

func callStockAPI(stockCode string) (*StockData, error) {

	resp := &StockData{StockCode: stockCode, Price: 100.00}
	return resp, nil
}
