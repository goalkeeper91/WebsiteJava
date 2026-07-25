package service

import (
	"context"
	"log"
	"time"
)

// LoyaltyScheduler ticks on a fixed interval and asks LoyaltyService to
// process whichever channels are due for a points accrual - 1:1 the
// ScheduledMessageScheduler pattern. Each channel's actual accrual cadence
// is governed by its own next_accrual_at, not by this tick's frequency.
type LoyaltyScheduler struct {
	loyaltyService *LoyaltyService
	interval       time.Duration
	stopChan       chan struct{}
}

func NewLoyaltyScheduler(loyaltyService *LoyaltyService, interval time.Duration) *LoyaltyScheduler {
	return &LoyaltyScheduler{
		loyaltyService: loyaltyService,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

func (s *LoyaltyScheduler) Start() {
	log.Printf("⏰ Loyalty-Scheduler gestartet (Tick: %v)", s.interval)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.loyaltyService.RunAccrualTick(context.Background())
		case <-s.stopChan:
			log.Println("🛑 Loyalty-Scheduler gestoppt")
			return
		}
	}
}

func (s *LoyaltyScheduler) Stop() {
	close(s.stopChan)
}
