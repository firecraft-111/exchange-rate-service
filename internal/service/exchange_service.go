package service

import (
	"fmt"
	"time"

	"github.com/firecraft-111/exchange-rate-service/internal/domain"
	"github.com/firecraft-111/exchange-rate-service/internal/infrastructure"
)

type ExchangeService struct {
	cache *infrastructure.RatesCache
	ttl   time.Duration
}

func NewExchangeService(cache *infrastructure.RatesCache, ttl time.Duration) *ExchangeService {
	return &ExchangeService{
		cache: cache,
		ttl:   ttl,
	}
}

func (s *ExchangeService) GetLatestRates(base string) (map[string]float64, error) {
	if !domain.SupportedCurrencies[base] {
		return nil, fmt.Errorf("unsupported currency: %s", base)
	}

	rates, err := s.cache.Get(base)
	if err == nil {
		return rates, nil
	}

	apiResp, err := infrastructure.FetchLatestRates(base)
	if err != nil {
		return nil, err
	}

	s.cache.Set(base, apiResp.ConversionRates)

	return apiResp.ConversionRates, nil
}

func (s *ExchangeService) GetHistoricalRates(from, to string, date time.Time) (float64, error) {
	if !domain.SupportedCurrencies[from] || !domain.SupportedCurrencies[to] {
		return 0, fmt.Errorf("unsupported currency: %s or %s", from, to)
	}

	if time.Since(date) > 90*24*time.Hour {
		return 0, fmt.Errorf("date %s too old: only last 90 days allowed", date.Format("2006-01-02"))
	}

	today := time.Now().Truncate(24 * time.Hour)
	if date.Equal(today) {
		rates, err := s.GetLatestRates(from)
		if err != nil {
			return 0, err
		}
		rate, ok := rates[to]
		if !ok {
			return 0, fmt.Errorf("rate not found for currency: %s", to)
		}
		return rate, nil
	}

	apiResp, err := infrastructure.FetchHistoricalRates(from, date)
	if err != nil {
		return 0, err
	}

	rate, ok := apiResp.ConversionRates[to]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency: %s", to)
	}

	return rate, nil
}
