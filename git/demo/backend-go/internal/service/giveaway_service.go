package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
)

type GiveawayService struct {
	giveawayRepo repository.GiveawayRepository
	redisService *redis.RedisService
}

func NewGiveawayService(giveawayRepo repository.GiveawayRepository, redisService *redis.RedisService) *GiveawayService {
	return &GiveawayService{giveawayRepo: giveawayRepo, redisService: redisService}
}

// notifyGiveawayChanged keeps the bot's in-memory codeword cache in sync
// regardless of whether the change came from a chat command or the
// dashboard (Giveaways Teil 3 made the dashboard a second writer) - see
// GiveawayCommands.refresh_single on the bot side.
func (s *GiveawayService) notifyGiveawayChanged(userTwitchID string) {
	if s.redisService != nil {
		if err := s.redisService.SendRefreshGiveawaySignal(userTwitchID); err != nil {
			fmt.Printf("⚠️ Fehler beim Senden des Giveaway-Reload-Signals: %v\n", err)
		}
	}
}

// StartGiveaway is bot-facing - fails if a giveaway is already open for this
// channel (only one at a time, enforced here rather than by a DB
// constraint, since giveaways is an episodic record table). keyword is
// normalized (trimmed, lowercased) before validation and storage so chat
// message matching is later a plain string comparison.
func (s *GiveawayService) StartGiveaway(ctx context.Context, userTwitchID, keyword string, subBonus bool) (*domain.Giveaway, error) {
	existing, err := s.giveawayRepo.GetOpenGiveaway(ctx, userTwitchID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrGiveawayAlreadyOpen
	}

	normalizedKeyword := domain.NormalizeGiveawayKeyword(keyword)
	if err := domain.ValidateGiveawayKeyword(normalizedKeyword); err != nil {
		return nil, err
	}

	giveaway := domain.NewGiveaway(userTwitchID, normalizedKeyword, subBonus)
	if err := s.giveawayRepo.CreateGiveaway(ctx, giveaway); err != nil {
		return nil, fmt.Errorf("fehler beim Starten des Giveaways: %w", err)
	}

	s.notifyGiveawayChanged(userTwitchID)

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

	s.notifyGiveawayChanged(userTwitchID)

	return giveaway, nil
}

// CancelGiveaway closes an open giveaway without drawing a winner - lets a
// streamer end one with zero (or unwanted) entries instead of being stuck
// until someone enters.
func (s *GiveawayService) CancelGiveaway(ctx context.Context, userTwitchID string) (*domain.Giveaway, error) {
	giveaway, err := s.giveawayRepo.GetOpenGiveaway(ctx, userTwitchID)
	if err != nil {
		return nil, err
	}
	if giveaway == nil {
		return nil, domain.ErrNoOpenGiveaway
	}

	if err := s.giveawayRepo.CancelGiveaway(ctx, giveaway.ID); err != nil {
		return nil, err
	}

	giveaway.Status = domain.GiveawayStatusClosed

	s.notifyGiveawayChanged(userTwitchID)

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

// GetOpenGiveawaysForBot backs the bot's startup keyword-cache warmup - no
// session auth at this layer, the handler gates it via the shared internal
// secret instead (same as AutomodService.GetAllEnabledSettingsForBot).
func (s *GiveawayService) GetOpenGiveawaysForBot(ctx context.Context) ([]*domain.Giveaway, error) {
	return s.giveawayRepo.GetAllOpenGiveaways(ctx)
}

// GetHistory is dashboard-facing, paginated.
func (s *GiveawayService) GetHistory(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.Giveaway, int64, error) {
	return s.giveawayRepo.GetHistory(ctx, userTwitchID, limit, offset)
}
