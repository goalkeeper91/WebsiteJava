package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"demo/backend-go/internal/repository"
)

type TokenRefreshService struct {
	userRepo    repository.UserRepository
	tokenRepo   repository.AuthTokenRepository
	authService *AuthService
	interval    time.Duration
	stopChan    chan struct{}
}

func NewTokenRefreshService(
	userRepo repository.UserRepository,
	tokenRepo repository.AuthTokenRepository,
	authService *AuthService,
	interval time.Duration,
) *TokenRefreshService {
	return &TokenRefreshService{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		authService: authService,
		interval:    interval,
		stopChan:    make(chan struct{}),
	}
}

func (s *TokenRefreshService) Start() {
	log.Printf("🔄 Token Refresh Service gestartet (Interval: %v)", s.interval)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.refreshAllTokens()

	for {
		select {
		case <-ticker.C:
			s.refreshAllTokens()
		case <-s.stopChan:
			log.Println("🛑 Token Refresh Service gestoppt")
			return
		}
	}
}

func (s *TokenRefreshService) Stop() {
	close(s.stopChan)
}

func (s *TokenRefreshService) refreshAllTokens() {
	ctx := context.Background()

	botUsers, err := s.userRepo.GetBotUsers(ctx)
	if err != nil {
		log.Printf("❌ Fehler beim Laden der Bot Users: %v", err)
		return
	}

	if len(botUsers) == 0 {
		return
	}

	for _, botUser := range botUsers {
		token, err := s.tokenRepo.GetByTwitchUserID(ctx, botUser.TwitchID)
		if err != nil {
			log.Printf("⚠️ Token für Bot User %s nicht gefunden: %v", botUser.Username, err)
			continue
		}

		expiryTime := token.UpdatedAt.Add(time.Duration(token.ExpiresIn) * time.Second)
		timeUntilExpiry := time.Until(expiryTime)

		if timeUntilExpiry < 24*time.Hour {
			log.Printf("🔄 Token für Bot User %s läuft bald ab (in %v) - erneuere...",
				botUser.Username,
				timeUntilExpiry.Round(time.Minute),
			)

			if err := s.authService.RefreshToken(ctx, botUser.TwitchID); err != nil {
				log.Printf("❌ Fehler beim Erneuern des Tokens für %s: %v", botUser.Username, err)
			} else {
				log.Printf("✅ Token für Bot User %s erfolgreich erneuert", botUser.Username)
			}
		} else {
			log.Printf("✓ Token für Bot User %s ist noch gültig (läuft ab in %v)",
				botUser.Username,
				timeUntilExpiry.Round(time.Hour),
			)
		}
	}
}

func (s *TokenRefreshService) RefreshTokenNow(ctx context.Context, twitchUserID string) error {
	log.Printf("🔄 Manuelles Token Refresh angefordert für User: %s", twitchUserID)

	if err := s.authService.RefreshToken(ctx, twitchUserID); err != nil {
		return fmt.Errorf("fehler beim Erneuern des Tokens: %w", err)
	}

	log.Printf("✅ Token erfolgreich erneuert für User: %s", twitchUserID)
	return nil
}

func (s *TokenRefreshService) SetAuthService(authService *AuthService) {
    s.authService = authService
}