package handler

import (
	"encoding/json"
	"net/http"

	"github.com/firecraft-111/exchange-rate-service/internal/service"
)

type LatestRateResponse struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

func RegisterRoutes(mux *http.ServeMux, svc *service.ExchangeService) {
	mux.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		handleLatestRate(w, r, svc)
	})
}

func handleLatestRate(w http.ResponseWriter, r *http.Request, svc *service.ExchangeService) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		http.Error(w, "`from` and `to` query parameters are required", http.StatusBadRequest)
		return
	}

	rates, err := svc.GetRates(from)
	if err != nil {
		http.Error(w, "Failed to get latest rates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rate, ok := rates[to]
	if !ok {
		http.Error(w, "Unsupported target currency: "+to, http.StatusBadRequest)
		return
	}

	resp := LatestRateResponse{
		From: from,
		To:   to,
		Rate: rate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
