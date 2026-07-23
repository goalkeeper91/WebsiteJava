package service

import (
	"context"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository/postgres"
)

// SubathonService ports the original Next.js app's lib/timerStore.ts logic
// to a DB-backed, per-user model (the original kept a single in-memory
// global object, fine for one streamer, wrong for a multi-tenant dashboard).
type SubathonService struct {
	repo *postgres.SubathonRepository
}

func NewSubathonService(repo *postgres.SubathonRepository) *SubathonService {
	return &SubathonService{repo: repo}
}

// GetState returns the user's timer state, recomputing time_remaining live
// if it's currently running (matches timerStore.ts's getTimerState()) and
// persisting an auto-stop once it hits zero.
func (s *SubathonService) GetState(ctx context.Context, userID string) (*domain.SubathonState, error) {
	state, err := s.repo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	if state.IsRunning && state.TargetTimestamp != nil {
		remaining := int((*state.TargetTimestamp - time.Now().UnixMilli()) / 1000)
		if remaining < 0 {
			remaining = 0
		}
		state.TimeRemaining = remaining

		if remaining == 0 {
			state.IsRunning = false
			state.TargetTimestamp = nil
			zero := 0
			updated, err := s.repo.Update(ctx, userID, domain.SubathonStateUpdateInput{
				IsRunning:       boolPtr(false),
				TargetTimestamp: nil,
				TimeRemaining:   &zero,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to persist auto-stop: %w", err)
			}
			return updated, nil
		}
	}

	return state, nil
}

// StartTimer is a no-op if already running (matches the original's guard).
func (s *SubathonService) StartTimer(ctx context.Context, userID string) (*domain.SubathonState, error) {
	state, err := s.GetState(ctx, userID)
	if err != nil {
		return nil, err
	}
	if state.IsRunning {
		return state, nil
	}

	target := time.Now().UnixMilli() + int64(state.TimeRemaining)*1000
	return s.repo.Update(ctx, userID, domain.SubathonStateUpdateInput{
		IsRunning:       boolPtr(true),
		TargetTimestamp: &target,
	})
}

func (s *SubathonService) PauseTimer(ctx context.Context, userID string) (*domain.SubathonState, error) {
	state, err := s.GetState(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !state.IsRunning {
		return state, nil
	}

	return s.repo.Update(ctx, userID, domain.SubathonStateUpdateInput{
		IsRunning:       boolPtr(false),
		TargetTimestamp: nil,
		TimeRemaining:   &state.TimeRemaining,
	})
}

// ResetTimer stops the timer, restores time_remaining from the configured
// initial time, and clears stats/log - matches the dashboard's Reset button.
func (s *SubathonService) ResetTimer(ctx context.Context, userID string) (*domain.SubathonState, error) {
	state, err := s.repo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	remaining := state.InitialTimeMinutes * 60
	zero := 0
	return s.repo.Update(ctx, userID, domain.SubathonStateUpdateInput{
		IsRunning:       boolPtr(false),
		TargetTimestamp: nil,
		TimeRemaining:   &remaining,
		TotalSubs:       &zero,
		TotalBits:       &zero,
		TotalEvents:     &zero,
		EventLog:        domain.SubathonEventLog{},
	})
}

// UpdateSettings only touches the fields that were actually provided (nil =
// leave unchanged) and always passes the current target_timestamp through
// unmodified, since the repository never COALESCEs that column.
func (s *SubathonService) UpdateSettings(ctx context.Context, userID string, initialTimeMinutes, subTimeSeconds, bitsTimeSeconds *int) (*domain.SubathonState, error) {
	state, err := s.repo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, userID, domain.SubathonStateUpdateInput{
		TargetTimestamp:    state.TargetTimestamp,
		InitialTimeMinutes: initialTimeMinutes,
		SubTimeSeconds:     subTimeSeconds,
		BitsTimeSeconds:    bitsTimeSeconds,
	})
}

// ProcessEvent applies a sub or bits cheer to the timer (matches
// timerStore.ts's processEvent): adds time, bumps stats, and prepends a
// line to the capped 50-entry event log. Called by the EventSub client
// whenever Twitch reports a real sub/cheer.
func (s *SubathonService) ProcessEvent(ctx context.Context, userID, eventType, userName string, amount int, subTier string) error {
	state, err := s.repo.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}

	var secondsToAdd int
	var logText string

	switch eventType {
	case "sub":
		multiplier := 1
		switch subTier {
		case "2000":
			multiplier = 2
		case "3000":
			multiplier = 5
		}
		secondsToAdd = state.SubTimeSeconds * multiplier
		logText = fmt.Sprintf("Sub: %s (+%ds)", userName, secondsToAdd)
	case "bits":
		secondsToAdd = (amount / 100) * state.BitsTimeSeconds
		logText = fmt.Sprintf("%d Bits: %s (+%ds)", amount, userName, secondsToAdd)
	default:
		return fmt.Errorf("unknown subathon event type: %s", eventType)
	}

	newRemaining := state.TimeRemaining + secondsToAdd

	var newTarget *int64
	if state.IsRunning && state.TargetTimestamp != nil {
		t := *state.TargetTimestamp + int64(secondsToAdd)*1000
		newTarget = &t
	}

	entry := domain.SubathonEventLogEntry{Time: time.Now().Format("15:04:05"), Text: logText}
	newLog := append(domain.SubathonEventLog{entry}, state.EventLog...)
	if len(newLog) > 50 {
		newLog = newLog[:50]
	}

	totalSubs := state.TotalSubs
	totalBits := state.TotalBits
	if eventType == "sub" {
		totalSubs++
	} else {
		totalBits += amount
	}
	totalEvents := state.TotalEvents + 1

	_, err = s.repo.Update(ctx, userID, domain.SubathonStateUpdateInput{
		TimeRemaining:   &newRemaining,
		TargetTimestamp: newTarget,
		TotalSubs:       &totalSubs,
		TotalBits:       &totalBits,
		TotalEvents:     &totalEvents,
		EventLog:        newLog,
	})
	if err != nil {
		return fmt.Errorf("failed to persist subathon event: %w", err)
	}
	return nil
}

// GetAllUserIDs is used by the EventSub background client to know which
// users need a live Twitch connection.
func (s *SubathonService) GetAllUserIDs(ctx context.Context) ([]string, error) {
	return s.repo.GetAllUserIDs(ctx)
}
