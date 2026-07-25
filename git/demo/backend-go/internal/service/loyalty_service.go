package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/twitch"
)

// dueLoyaltyBatchSize caps how many due channels one scheduler tick
// processes - keeps a single slow tick bounded, same role as
// dueMessagesBatchSize for scheduled messages.
const dueLoyaltyBatchSize = 100

type LoyaltyService struct {
	loyaltyRepo repository.LoyaltyRepository
	tokenRepo   repository.AuthTokenRepository
	appToken    *twitch.TwitchAppTokenClient
	chatters    *twitch.ChattersClient
}

func NewLoyaltyService(
	loyaltyRepo repository.LoyaltyRepository,
	tokenRepo repository.AuthTokenRepository,
	appToken *twitch.TwitchAppTokenClient,
	chatters *twitch.ChattersClient,
) *LoyaltyService {
	return &LoyaltyService{
		loyaltyRepo: loyaltyRepo,
		tokenRepo:   tokenRepo,
		appToken:    appToken,
		chatters:    chatters,
	}
}

// GetSettings is dashboard-facing - userTwitchID is directly the settings
// row's PK, same convention as AutomodService.
func (s *LoyaltyService) GetSettings(ctx context.Context, userTwitchID string) (*domain.LoyaltySettings, error) {
	return s.loyaltyRepo.GetSettings(ctx, userTwitchID)
}

func (s *LoyaltyService) UpdateSettings(ctx context.Context, userTwitchID string, input domain.LoyaltySettingsUpdateInput) (*domain.LoyaltySettings, error) {
	settings, err := s.loyaltyRepo.GetSettings(ctx, userTwitchID)
	if err != nil {
		return nil, err
	}

	settings.ApplyUpdate(input)

	if err := s.loyaltyRepo.UpsertSettings(ctx, settings); err != nil {
		return nil, fmt.Errorf("fehler beim Speichern der Loyalty-Settings: %w", err)
	}

	return settings, nil
}

// GetLeaderboard is dashboard-facing, paginated.
func (s *LoyaltyService) GetLeaderboard(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.LoyaltyPoints, int64, error) {
	return s.loyaltyRepo.GetLeaderboard(ctx, userTwitchID, limit, offset)
}

// GetViewerPointsForBot backs the `!points` chat command - no session auth
// at this layer, the handler gates it via the shared internal secret.
func (s *LoyaltyService) GetViewerPointsForBot(ctx context.Context, userTwitchID, viewerLogin string) (int64, string, error) {
	settings, err := s.loyaltyRepo.GetSettings(ctx, userTwitchID)
	if err != nil {
		return 0, "", err
	}

	points, err := s.loyaltyRepo.GetViewerPoints(ctx, userTwitchID, viewerLogin)
	if err != nil {
		return 0, "", err
	}
	if points == nil {
		return 0, settings.PointsName, nil
	}

	return points.Points, settings.PointsName, nil
}

// GetLeaderboardForBot backs the `!top` chat command.
func (s *LoyaltyService) GetLeaderboardForBot(ctx context.Context, userTwitchID string, limit int) ([]*domain.LoyaltyPoints, string, error) {
	settings, err := s.loyaltyRepo.GetSettings(ctx, userTwitchID)
	if err != nil {
		return nil, "", err
	}

	entries, _, err := s.loyaltyRepo.GetLeaderboard(ctx, userTwitchID, limit, 0)
	if err != nil {
		return nil, "", err
	}

	return entries, settings.PointsName, nil
}

// GetRegularsForBot backs the `!regulars` chat command - returns the
// currently-qualifying logins, or an empty slice if Regulars is disabled
// (RegularsThreshold <= 0).
func (s *LoyaltyService) GetRegularsForBot(ctx context.Context, userTwitchID string) ([]string, error) {
	settings, err := s.loyaltyRepo.GetSettings(ctx, userTwitchID)
	if err != nil {
		return nil, err
	}
	if settings.RegularsThreshold <= 0 {
		return []string{}, nil
	}

	return s.loyaltyRepo.GetViewerLoginsAboveThreshold(ctx, userTwitchID, settings.RegularsThreshold)
}

// RunAccrualTick is called on every scheduler tick. It loads due, enabled
// channels, batch-checks live status (same pattern as
// ScheduledMessageService.RunDueMessages), then for live channels fetches
// the current chatter list via the channel's own decrypted broadcaster
// token and credits points to everyone present except the broadcaster
// themself. Offline channels still have their schedule advanced - no
// catch-up once they go live again.
func (s *LoyaltyService) RunAccrualTick(ctx context.Context) {
	due, err := s.loyaltyRepo.GetDueForAccrual(ctx, dueLoyaltyBatchSize)
	if err != nil {
		log.Printf("⚠️ Fehler beim Laden fälliger Loyalty-Settings: %v", err)
		return
	}
	if len(due) == 0 {
		return
	}

	ids := make([]string, 0, len(due))
	for _, settings := range due {
		ids = append(ids, settings.UserTwitchID)
	}

	liveStreams, err := s.appToken.GetLiveStreams(ctx, ids)
	if err != nil {
		log.Printf("⚠️ Fehler beim Live-Status-Check für Loyalty: %v", err)
		liveStreams = map[string]twitch.LiveStream{}
	}

	for _, settings := range due {
		liveStream, isLive := liveStreams[settings.UserTwitchID]
		if isLive {
			s.creditChannel(ctx, settings, liveStream.UserLogin)
		}

		settings.MarkAccrued()
		if err := s.loyaltyRepo.MarkAccrued(ctx, settings.UserTwitchID, settings.NextAccrualAt); err != nil {
			log.Printf("⚠️ Fehler beim Fortschreiben des Loyalty-Zeitplans für %s: %v", settings.UserTwitchID, err)
		}
	}
}

func (s *LoyaltyService) creditChannel(ctx context.Context, settings *domain.LoyaltySettings, broadcasterLogin string) {
	authToken, err := s.tokenRepo.GetByTwitchUserID(ctx, settings.UserTwitchID)
	if err != nil {
		log.Printf("⚠️ Kein Broadcaster-Token für Loyalty-Kanal %s: %v", settings.UserTwitchID, err)
		return
	}

	chatters, err := s.chatters.GetChatters(ctx, settings.UserTwitchID, authToken.AccessToken)
	if err != nil {
		log.Printf("⚠️ Fehler beim Get-Chatters-Aufruf für %s: %v", settings.UserTwitchID, err)
		return
	}
	if len(chatters) == 0 {
		return
	}

	credits := make([]domain.ViewerCredit, 0, len(chatters))
	for _, chatter := range chatters {
		if strings.EqualFold(chatter.UserLogin, broadcasterLogin) {
			continue
		}
		credits = append(credits, domain.ViewerCredit{
			ViewerTwitchID: chatter.UserID,
			ViewerLogin:    chatter.UserLogin,
			Amount:         int64(settings.PointsPerInterval),
		})
	}

	if err := s.loyaltyRepo.CreditPoints(ctx, settings.UserTwitchID, credits); err != nil {
		log.Printf("⚠️ Fehler beim Gutschreiben von Loyalty-Punkten für %s: %v", settings.UserTwitchID, err)
	}
}
