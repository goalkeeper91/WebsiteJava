package service

import (
	"context"
	"log"
	"time"
)

type StreamSessionScheduler struct {
	sessionService *StreamSessionService
	interval       time.Duration
	stopChan       chan struct{}
}

func NewStreamSessionScheduler(sessionService *StreamSessionService, interval time.Duration) *StreamSessionScheduler {
	return &StreamSessionScheduler{
		sessionService: sessionService,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

func (s *StreamSessionScheduler) Start() {
	log.Printf("⏰ Stream-Session-Sampler gestartet (Tick: %v)", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.sessionService.RunSampleTick(context.Background())
		case <-s.stopChan:
			log.Println("🛑 Stream-Session-Sampler gestoppt")
			return
		}
	}
}

func (s *StreamSessionScheduler) Stop() {
	close(s.stopChan)
}
