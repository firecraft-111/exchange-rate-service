package infrastructure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/firecraft-111/exchange-rate-service/internal/config"
	"github.com/firecraft-111/exchange-rate-service/internal/domain"
)

type ApiResponse struct {
	Result            string             `json:"result"`
	Documentation     string             `json:"documentation"`
	TermsOfUse        string             `json:"terms_of_use"`
	TimeLastUpdate    int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC string             `json:"time_last_update_utc"`
	TimeNextUpdate    int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC string             `json:"time_next_update_utc"`
	BaseCode          string             `json:"base_code"`
	ConversionRates   map[string]float64 `json:"conversion_rates"`
	ErrorType         string             `json:"error-type,omitempty"` 
}

func FetchLatestRates(base string) (*ApiResponse, error) {
	if !domain.SupportedCurrencies[base] {
		return nil, fmt.Errorf("unsupported base currency: %s", base)
	}

	apiKey := ""
	if config.App != nil {
		apiKey = config.App.ExchangeRate.APIKey
	}

	if apiKey == "" {
		return nil, fmt.Errorf("missing API key: set exchange_rate.api_key in config or EXCHANGE_RATE_API_KEY environment variable")
	}

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, base)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status: %d", resp.StatusCode)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if apiResp.Result == "error" {
		return nil, fmt.Errorf("API error: %s", apiResp.ErrorType)
	}

	return &apiResp, nil
}
