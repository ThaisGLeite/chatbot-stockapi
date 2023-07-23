package main

import (
	"fmt"
	"net/http"
	"os"

	"botService/natsclient"
	"botService/stock"
)

func main() {
	// Connect to NATS server
	natsClient, err := natsclient.NewNATSClient(os.Getenv("NATS_URL"))
	if err != nil {
		fmt.Println("Error connecting to NATS server: ", err)
		os.Exit(1)
	}
	defer natsClient.Close()

	// Create stock data handler
	stockDataHandler := stock.NewStockDataHandler(natsClient)

	// Subscribe to stock_codes subject
	_, err = natsClient.Subscribe("stock_codes", stockDataHandler.HandleRequest)
	if err != nil {
		fmt.Println("Error subscribing to stock_codes: ", err)
		os.Exit(1)
	}

	fmt.Println("Bot service started")

	// Start HTTP server
	http.HandleFunc("/get-stock-data", stockDataHandler.GetStockDataHTTPHandler)
	http.ListenAndServe(":3000", nil)
}
