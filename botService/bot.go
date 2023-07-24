package main

import (
	"botService/natsclient"
	"botService/stock"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Initialize services
func InitializeServices(nc natsclient.NATSClientInterface, sdh stock.StockDataHandlerInterface) error {
	// Subscribe to stock_codes subject
	if _, err := nc.Subscribe("stock_codes", sdh.HandleRequest); err != nil {
		return err
	}

	logger.Info("Bot service initialized")

	return nil
}

func main() {
	// Create logger
	logger = logrus.New()
	// Connect to NATS server
	natsURL := os.Getenv("NATS_URL")
	opt := natsclient.ConnectionOptions{
		URL:            natsURL,
		ReconnectDelay: 5 * time.Second,
		MaxReconnects:  3,
	}

	natsClient, err := natsclient.NewNATSClient(opt)
	if err != nil {
		logger.Fatalf("Server could not be started: %v", err)
	}
	defer natsClient.Close()

	// Create stock data handler
	stockDataHandler := stock.NewStockDataHandler(natsClient, logger, os.Getenv("API_URL"))

	// Initialize services
	if err := InitializeServices(natsClient, stockDataHandler); err != nil {
		logger.Fatalf("Server could not be started: %v", err)
	}

	// Start HTTP server
	porta := os.Getenv("SERVER_PORT")
	http.HandleFunc("/get-stock-data", stockDataHandler.GetStockDataHTTPHandler)
	err = http.ListenAndServe(porta, nil)
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
