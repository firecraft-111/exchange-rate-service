package infrastructure

import (
	"errors"
	"sync"
	"time"
)

type CacheItem struct {
	Rates      map[string]float64
	Expiration time.Time
}

type RatesCache struct {
	data map[string]*CacheItem
	mu   sync.RWMutex
	ttl  time.Duration
}

func NewRatesCache(ttl time.Duration) *RatesCache {
	return &RatesCache{
		data: make(map[string]*CacheItem),
		ttl:  ttl,
	}
}

func (c *RatesCache) Set(base string, rates map[string]float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[base] = &CacheItem{
		Rates:      rates,
		Expiration: time.Now().Add(c.ttl),
	}
}

func (c *RatesCache) Get(base string) (map[string]float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.data[base]
	if !exists {
		return nil, errors.New("cache miss")
	}
	if time.Now().After(item.Expiration) {
		return nil, errors.New("cache expired")
	}

	return item.Rates, nil
}

func (c *RatesCache) IsExpires(base string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.data[base]
	if !exists || time.Now().After(item.Expiration) {
		return true
	}
	return false
}
