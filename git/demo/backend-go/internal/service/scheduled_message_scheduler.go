package service

import (
	"context"
	"log"
	"time"
)

// ScheduledMessageScheduler ticks on a fixed interval and asks
// ScheduledMessageService to process whichever per-channel messages are due.
// A single shared tick (rather than one goroutine/ticker per message) is
// enough since due-checking is a cheap indexed DB query - see the plan.
type ScheduledMessageScheduler struct {
	messageService *ScheduledMessageService
	interval       time.Duration
	stopChan       chan struct{}
}

func NewScheduledMessageScheduler(messageService *ScheduledMessageService, interval time.Duration) *ScheduledMessageScheduler {
	return &ScheduledMessageScheduler{
		messageService: messageService,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

func (s *ScheduledMessageScheduler) Start() {
	log.Printf("⏰ Scheduled-Message-Scheduler gestartet (Tick: %v)", s.interval)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.messageService.RunDueMessages(context.Background())
		case <-s.stopChan:
			log.Println("🛑 Scheduled-Message-Scheduler gestoppt")
			return
		}
	}
}

func (s *ScheduledMessageScheduler) Stop() {
	close(s.stopChan)
}
