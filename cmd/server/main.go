package main

import (
	"fmt"
	"log"

	"github.com/firecraft-111/exchange-rate-service/internal/config"
	"github.com/firecraft-111/exchange-rate-service/internal/infrastructure"
)

func main() {
	// http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "pong")
	// })

	// fmt.Println("Server running on :8080")
	// http.ListenAndServe(":8080", nil)

	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
		panic(err)
	}

	base := "USD"
	resp, err := infrastructure.FetchLatestRates(base)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Base: %s\nDate: %s\nRates:\n", resp.BaseCode, resp.TimeLastUpdateUTC)

	for k, v := range resp.ConversionRates {
		fmt.Printf("%s: %.4f\n", k, v)
	}
}
