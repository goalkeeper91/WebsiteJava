package service

import (
	"context"
	"log"
	"time"
)

// PriceSyncScheduler ticks once a day and asks PriceSyncService to pull the
// current Pro/Premium prices from Paddle - 1:1 the TokenRefreshService
// pattern (an initial run at startup, then on a fixed ticker), so a fresh
// deploy fixes any drift immediately instead of waiting up to a full day.
type PriceSyncScheduler struct {
	priceSyncService *PriceSyncService
	interval         time.Duration
	stopChan         chan struct{}
}

func NewPriceSyncScheduler(priceSyncService *PriceSyncService, interval time.Duration) *PriceSyncScheduler {
	return &PriceSyncScheduler{
		priceSyncService: priceSyncService,
		interval:         interval,
		stopChan:         make(chan struct{}),
	}
}

func (s *PriceSyncScheduler) Start() {
	log.Printf("💶 Preis-Sync-Scheduler gestartet (Interval: %v)", s.interval)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.priceSyncService.SyncPrices(context.Background())

	for {
		select {
		case <-ticker.C:
			s.priceSyncService.SyncPrices(context.Background())
		case <-s.stopChan:
			log.Println("🛑 Preis-Sync-Scheduler gestoppt")
			return
		}
	}
}

func (s *PriceSyncScheduler) Stop() {
	close(s.stopChan)
}
