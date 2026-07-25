package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type GiveawayService struct {
	giveawayRepo repository.GiveawayRepository
}

func NewGiveawayService(giveawayRepo repository.GiveawayRepository) *GiveawayService {
	return &GiveawayService{giveawayRepo: giveawayRepo}
}

// StartGiveaway is bot-facing - fails if a giveaway is already open for this
// channel (only one at a time, enforced here rather than by a DB
// constraint, since giveaways is an episodic record table).
func (s *GiveawayService) StartGiveaway(ctx context.Context, userTwitchID string, subBonus bool) (*domain.Giveaway, error) {
	existing, err := s.giveawayRepo.GetOpenGiveaway(ctx, userTwitchID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrGiveawayAlreadyOpen
	}

	giveaway := domain.NewGiveaway(userTwitchID, subBonus)
	if err := s.giveawayRepo.CreateGiveaway(ctx, giveaway); err != nil {
		return nil, fmt.Errorf("fehler beim Starten des Giveaways: %w", err)
	}

	return giveaway, nil
}

// EnterGiveaway is bot-facing - inserted=false means the viewer had already
// entered, so the bot can stay silent instead of confirm-spamming on every
// repeat `!giveaway`.
func (s *GiveawayService) EnterGiveaway(ctx context.Context, userTwitchID, viewerTwitchID, viewerLogin string, isSubscriber bool) (bool, error) {
	giveaway, err := s.giveawayRepo.GetOpenGiveaway(ctx, userTwitchID)
	if err != nil {
		return false, err
	}
	if giveaway == nil {
		return false, domain.ErrNoOpenGiveaway
	}

	entries := 1
	if giveaway.SubBonus && isSubscriber {
		entries = 2
	}

	return s.giveawayRepo.AddEntry(ctx, giveaway.ID, viewerTwitchID, viewerLogin, entries)
}

// DrawWinner is bot-facing - picks a winner weighted by each entrant's
// ticket count using crypto/rand (same package as auth_service.go's OAuth
// state generation), then closes the giveaway.
func (s *GiveawayService) DrawWinner(ctx context.Context, userTwitchID string) (*domain.Giveaway, error) {
	giveaway, err := s.giveawayRepo.GetOpenGiveaway(ctx, userTwitchID)
	if err != nil {
		return nil, err
	}
	if giveaway == nil {
		return nil, domain.ErrNoOpenGiveaway
	}

	entries, err := s.giveawayRepo.GetEntries(ctx, giveaway.ID)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, domain.ErrNoGiveawayEntries
	}

	totalTickets := 0
	for _, e := range entries {
		totalTickets += e.Entries
	}

	drawn, err := rand.Int(rand.Reader, big.NewInt(int64(totalTickets)))
	if err != nil {
		return nil, fmt.Errorf("fehler bei der zufälligen Ziehung: %w", err)
	}
	target := drawn.Int64()

	winner := entries[len(entries)-1]
	cumulative := int64(0)
	for _, e := range entries {
		cumulative += int64(e.Entries)
		if target < cumulative {
			winner = e
			break
		}
	}

	if err := s.giveawayRepo.CloseGiveaway(ctx, giveaway.ID, winner.ViewerTwitchID, winner.ViewerLogin); err != nil {
		return nil, err
	}

	giveaway.Status = domain.GiveawayStatusClosed
	giveaway.WinnerTwitchID = &winner.ViewerTwitchID
	giveaway.WinnerLogin = &winner.ViewerLogin

	return giveaway, nil
}

// GetStatus is used by both the dashboard's live-status card and the bot's
// `!giveaway status` command. Returns nil giveaway (not an error) if none
// is open.
func (s *GiveawayService) GetStatus(ctx context.Context, userTwitchID string) (*domain.Giveaway, int, error) {
	giveaway, err := s.giveawayRepo.GetOpenGiveaway(ctx, userTwitchID)
	if err != nil {
		return nil, 0, err
	}
	if giveaway == nil {
		return nil, 0, nil
	}

	count, err := s.giveawayRepo.GetEntryCount(ctx, giveaway.ID)
	if err != nil {
		return giveaway, 0, err
	}

	return giveaway, count, nil
}

// GetHistory is dashboard-facing, paginated.
func (s *GiveawayService) GetHistory(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.Giveaway, int64, error) {
	return s.giveawayRepo.GetHistory(ctx, userTwitchID, limit, offset)
}
