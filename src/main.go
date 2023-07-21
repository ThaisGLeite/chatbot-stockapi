package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/nats-io/nats.go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	stock_code := "aapl.us"
	url := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stock_code)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "YOUR_USER_AGENT")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making the request:", err)
	}

	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}

	if len(records) < 2 || len(records[1]) < 7 {
		log.Fatal("Unexpected data format")
	}

	quote := fmt.Sprintf("%s quote is $%s per share", strings.ToUpper(stock_code), records[1][6])

	nc, err := nats.Connect(nats.DefaultURL)
	failOnError(err, "Failed to connect to NATS")
	defer nc.Close()

	// Simple Publisher
	err = nc.Publish("stock_quotes", []byte(quote))
	failOnError(err, "Failed to publish a message")

	// Flush the connection to the server, which ensures all messages
	// have been processed by the server.
	nc.Flush()

	fmt.Println("Sent: ", quote)
}
