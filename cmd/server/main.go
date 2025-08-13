package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/firecraft-111/exchange-rate-service/internal/config"
	"github.com/firecraft-111/exchange-rate-service/internal/handler"
	"github.com/firecraft-111/exchange-rate-service/internal/infrastructure"
	"github.com/firecraft-111/exchange-rate-service/internal/service"
)

func main() {
	if err := config.Load("config.yaml"); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cache := infrastructure.NewRatesCache(time.Hour)
	exchangeService := service.NewExchangeService(cache, time.Hour)
	scheduler := service.NewScheduler(exchangeService, time.Hour)

	scheduler.Start()
	defer scheduler.Stop()

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, exchangeService)

	port := config.App.Server.Port
	log.Printf("Server starting on port %d..", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
