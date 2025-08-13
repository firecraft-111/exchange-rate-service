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

func (s *ExchangeService) GetRates(base string) (map[string]float64, error) {
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
