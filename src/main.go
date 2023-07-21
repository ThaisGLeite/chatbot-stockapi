package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	stock_code := "aapl.us"
	url := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stock_code)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "YOUR_USER_AGENT")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return
	}

	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)

	records, err := r.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	if len(records) < 2 || len(records[1]) < 7 {
		fmt.Println("Unexpected data format")
		return
	}

	fmt.Printf("%s quote is $%s per share\n", strings.ToUpper(stock_code), records[1][6])
}
