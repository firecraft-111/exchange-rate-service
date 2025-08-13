package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/firecraft-111/exchange-rate-service/internal/service"
)

type LatestRateResponse struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

type ConvertResponse struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Converted float64 `json:"converted"`
}

// Register all routes
func RegisterRoutes(mux *http.ServeMux, svc *service.ExchangeService) {
	mux.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		handleLatestRate(w, r, svc)
	})

	mux.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		handleConvert(w, r, svc)
	})

	mux.HandleFunc("/historical", func(w http.ResponseWriter, r *http.Request) {
		handleHistorical(w, r, svc)
	})
}

// Helper to write JSON responses
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper to write error responses
func writeError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}

// Helper to get and validate common query params
func getFromToParams(r *http.Request) (string, string, bool) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	return from, to, from != "" && to != ""
}

// Handler for /latest
func handleLatestRate(w http.ResponseWriter, r *http.Request, svc *service.ExchangeService) {
	from, to, ok := getFromToParams(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "`from` and `to` query parameters are required")
		return
	}

	rates, err := svc.GetLatestRates(from)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get latest rates: "+err.Error())
		return
	}

	rate, exists := rates[to]
	if !exists {
		writeError(w, http.StatusBadRequest, "Unsupported target currency: "+to)
		return
	}

	writeJSON(w, http.StatusOK, LatestRateResponse{From: from, To: to, Rate: rate})
}

// Handler for /convert
func handleConvert(w http.ResponseWriter, r *http.Request, svc *service.ExchangeService) {
	from, to, ok := getFromToParams(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "`from`, `to`, and `amount` are required")
		return
	}

	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount < 0 {
		writeError(w, http.StatusBadRequest, "Invalid amount")
		return
	}

	dateStr := r.URL.Query().Get("date")
	var rate float64

	if dateStr == "" {
		rates, err := svc.GetLatestRates(from)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to get rates: "+err.Error())
			return
		}
		var exists bool
		rate, exists = rates[to]
		if !exists {
			writeError(w, http.StatusBadRequest, "Unsupported target currency: "+to)
			return
		}
	} else {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
			return
		}
		if time.Since(date) > 90*24*time.Hour {
			writeError(w, http.StatusBadRequest, "Date too old. Only last 90 days allowed")
			return
		}

		rate, err = svc.GetHistoricalRates(from, to, date)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to get historical rate: "+err.Error())
			return
		}
	}

	converted := amount * rate
	writeJSON(w, http.StatusOK, ConvertResponse{
		From:      from,
		To:        to,
		Amount:    amount,
		Converted: converted,
	})
}

// Handler for /historical
func handleHistorical(w http.ResponseWriter, r *http.Request, svc *service.ExchangeService) {
	from, to, ok := getFromToParams(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "`from`, `to`, and `date` are required")
		return
	}

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	if time.Since(date) > 90*24*time.Hour {
		writeError(w, http.StatusBadRequest, "Date too old. Only last 90 days allowed")
		return
	}

	rate, err := svc.GetHistoricalRates(from, to, date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get historical rate: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, LatestRateResponse{From: from, To: to, Rate: rate})
}
