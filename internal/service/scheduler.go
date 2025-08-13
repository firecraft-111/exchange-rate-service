package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/firecraft-111/exchange-rate-service/internal/domain"
)

type Scheduler struct {
	service    *ExchangeService
	interval   time.Duration
	stopCtx    context.Context
	stopCancel context.CancelFunc
	wg         sync.WaitGroup
}

func NewScheduler(service *ExchangeService, interval time.Duration) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		service:    service,
		interval:   interval,
		stopCtx:    ctx,
		stopCancel: cancel,
	}
}

func (s *Scheduler) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		s.refreshAll()

		for {
			select {
			case <-s.stopCtx.Done():
				log.Println("Scheduler stopped")
				return

			case <-ticker.C:
				s.refreshAll()
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.stopCancel()
	s.wg.Wait()
}

func (s *Scheduler) refreshAll() {
	for base := range domain.SupportedCurrencies {
		_, err := s.service.GetRates(base)
		if err != nil {
			log.Printf("Scheduler: failed to refresh rates for %s: %v", base, err)
		} else {
			log.Printf("Scheduler: refreshed rates for %s", base)
		}
	}
}
