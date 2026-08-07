package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type TwitchUserResponse struct {
	Data []struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		Email           string `json:"email"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

var (
	UserScopes = []string{
		"user:read:email",
		"moderator:read:chat_messages",
		"user:read:email",
        "moderator:read:chat_messages",   // Erledigt: Chat lesen via EventSub
        "channel:read:subscriptions",    // Wichtig: Um auf neue Subs zu reagieren
        "channel:read:redemptions",      // Wichtig: Für Kanalpunkte-Belohnungen
        "channel:read:predictions",      // Optional: Falls der Bot Predictions anzeigt
        "moderator:read:followers",      // Wichtig: Um auf neue Follower zu reagieren
        "channel:moderate",
	}

	BotScopes = []string{
		"user:bot",                      // Bot-specific scope
		"user:read:email",               // Basic user info
		"chat:read",                     // Chat lesen
		"chat:edit",                     // Chat schreiben
		"user:read:chat",                // User Chat Einstellungen lesen
		"user:write:chat",               // User Chat Einstellungen schreiben
		"channel:manage:broadcast",      // Stream Info ändern
		"moderator:read:followers",      // Follower lesen
		"moderator:read:chatters",       // Chatter lesen
		"channel:read:subscriptions",    // Subscriptions lesen
		"bits:read",                     // Bits/Cheers lesen
	}
)

type AuthService struct {
	userRepo      repository.UserRepository
	tokenRepo     repository.AuthTokenRepository
	oauthConfig   *oauth2.Config
	botOauthConfig *oauth2.Config
	frontendURL   string
	adminTwitchID string
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.AuthTokenRepository,
	oauthConfig *oauth2.Config,
	frontendURL string,
	adminTwitchID string,
) *AuthService {
	botOauthConfig := &oauth2.Config{
		ClientID:     oauthConfig.ClientID,
		ClientSecret: oauthConfig.ClientSecret,
		RedirectURL:  strings.TrimSuffix(oauthConfig.RedirectURL, "/callback") + "/callback/bot",
		Scopes:       BotScopes,
		Endpoint:     oauthConfig.Endpoint,
	}

	return &AuthService{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		oauthConfig:    oauthConfig,
		botOauthConfig: botOauthConfig,
		frontendURL:    frontendURL,
		adminTwitchID:  adminTwitchID,
	}
}

// GetAuthURL builds the Twitch authorize URL. redirectURI overrides the
// service's configured default (used when the login started on
// bot.goalkeeper91.de, see AuthHandler.oauthOriginForRequest) - pass "" to
// use the default, main-domain redirect_uri unchanged.
func (s *AuthService) GetAuthURL(redirectURI string) (string, string, error) {
	state, err := generateStateToken()
	if err != nil {
		return "", "", fmt.Errorf("fehler beim Generieren des State Tokens: %w", err)
	}

	cfg := s.oauthConfig
	if redirectURI != "" {
		cfgCopy := *s.oauthConfig
		cfgCopy.RedirectURL = redirectURI
		cfg = &cfgCopy
	}

	url := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url, state, nil
}

// HandleCallback exchanges the Twitch authorization code. redirectURI must
// be the exact same value passed to GetAuthURL for this flow - Twitch's
// token endpoint rejects a mismatch (RFC 6749 redirect_uri validation).
func (s *AuthService) HandleCallback(ctx context.Context, code, redirectURI string) (*domain.User, error) {
	return s.handleCallback(ctx, code, false, redirectURI)
}


func (s *AuthService) GetBotAuthURL() (string, string, error) {
	state, err := generateStateToken()
	if err != nil {
		return "", "", fmt.Errorf("fehler beim Generieren des State Tokens: %w", err)
	}

	url := s.botOauthConfig.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("force_verify", "true"),
	)

	return url, state, nil
}

func (s *AuthService) HandleBotCallback(ctx context.Context, code string) (*domain.User, error) {
	return s.handleCallback(ctx, code, true, "")
}

// ===== SHARED CALLBACK LOGIC =====

func (s *AuthService) handleCallback(ctx context.Context, code string, isBot bool, redirectURI string) (*domain.User, error) {
	config := s.oauthConfig
	if isBot {
		config = s.botOauthConfig
	}
	if redirectURI != "" {
		cfgCopy := *config
		cfgCopy.RedirectURL = redirectURI
		config = &cfgCopy
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Token-Austausch: %w", err)
	}

	twitchUserResp, err := s.getTwitchUser(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abrufen der Twitch User-Daten: %w", err)
	}

	if len(twitchUserResp.Data) == 0 {
		return nil, fmt.Errorf("keine User-Daten von Twitch erhalten")
	}

	userData := twitchUserResp.Data[0]

	user, err := s.upsertUser(ctx, twitchUserResp, isBot)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Speichern des Users: %w", err)
	}

	// ✅ FIX: Scopes aus Token extrahieren, oder Config Scopes als Fallback
	scope := extractScope(token)
	if scope == "" {
		if isBot {
			scope = strings.Join(BotScopes, " ")
			log.Printf("⚠️  Token hat keine Scopes, verwende BotScopes: %s", scope)
		} else {
			scope = strings.Join(UserScopes, " ")
			log.Printf("⚠️  Token hat keine Scopes, verwende UserScopes: %s", scope)
		}
	} else {
		log.Printf("✅ Scopes aus Token: %s", scope)
	}

	authToken := domain.NewAuthToken(
		userData.ID,
		userData.Login,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		scope,  // ← Jetzt garantiert nicht leer!
		int64(token.Expiry.Sub(time.Now()).Seconds()),
	)

	if err := s.tokenRepo.Upsert(ctx, authToken); err != nil {
		return nil, fmt.Errorf("fehler beim Speichern des Auth Tokens: %w", err)
	}

	return user, nil
}

// ===== HELPER METHODS =====

func (s *AuthService) GetUserByTwitchID(ctx context.Context, twitchID string) (*domain.User, error) {
	return s.userRepo.GetByTwitchID(ctx, twitchID)
}

func (s *AuthService) RefreshToken(ctx context.Context, twitchUserID string) error {
	authToken, err := s.tokenRepo.GetByTwitchUserID(ctx, twitchUserID)
	if err != nil {
		return fmt.Errorf("fehler beim Laden des Auth Tokens: %w", err)
	}

	tokenSource := s.oauthConfig.TokenSource(ctx, &oauth2.Token{
		RefreshToken: authToken.RefreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("fehler beim Erneuern des Tokens: %w", err)
	}

	authToken.Update(
		newToken.AccessToken,
		newToken.RefreshToken,
		int64(newToken.Expiry.Sub(time.Now()).Seconds()),
	)

	if err := s.tokenRepo.Update(ctx, authToken); err != nil {
		return fmt.Errorf("fehler beim Aktualisieren des Auth Tokens: %w", err)
	}

	return nil
}

// GetFreshAccessToken returns a broadcaster's access token, proactively
// refreshing it first if it's within 5 minutes of expiry. Unlike
// TokenRefreshService.refreshAllTokens (which only proactively refreshes
// bot-flagged accounts on a slow background tick), this exists for the
// Live-Dashboard's much more frequent (5-30s polling) broadcaster-token
// Helix calls, where hitting an expired token mid-request would otherwise
// just fail with no retry anywhere in this codebase.
func (s *AuthService) GetFreshAccessToken(ctx context.Context, twitchUserID string) (string, error) {
	authToken, err := s.tokenRepo.GetByTwitchUserID(ctx, twitchUserID)
	if err != nil {
		return "", fmt.Errorf("fehler beim Laden des Auth Tokens: %w", err)
	}

	expiryTime := authToken.UpdatedAt.Add(time.Duration(authToken.ExpiresIn) * time.Second)
	if time.Until(expiryTime) < 5*time.Minute {
		if err := s.RefreshToken(ctx, twitchUserID); err != nil {
			return "", fmt.Errorf("fehler beim proaktiven Token-Refresh: %w", err)
		}
		authToken, err = s.tokenRepo.GetByTwitchUserID(ctx, twitchUserID)
		if err != nil {
			return "", fmt.Errorf("fehler beim Laden des erneuerten Auth Tokens: %w", err)
		}
	}

	return authToken.AccessToken, nil
}

func (s *AuthService) getTwitchUser(ctx context.Context, accessToken string) (*TwitchUserResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", s.oauthConfig.ClientID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("twitch API Fehler: %d - %s", resp.StatusCode, string(body))
	}

	var userResp TwitchUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	if len(userResp.Data) == 0 {
		return nil, fmt.Errorf("keine User-Daten von Twitch erhalten")
	}

	return &userResp, nil
}

func (s *AuthService) upsertUser(ctx context.Context, twitchUserResp *TwitchUserResponse, isBot bool) (*domain.User, error) {
	if len(twitchUserResp.Data) == 0 {
		return nil, fmt.Errorf("keine User-Daten in Twitch Response")
	}

	userData := twitchUserResp.Data[0]

	idFromTwitch := strings.TrimSpace(userData.ID)
	idFromConfig := strings.TrimSpace(s.adminTwitchID)
	isAdmin := idFromTwitch == idFromConfig

	existingUser, err := s.userRepo.GetByTwitchID(ctx, userData.ID)
	if err == nil {
		existingUser.Update(userData.Login, userData.Email)
		existingUser.IsAdmin = isAdmin

		if isBot {
			existingUser.MarkAsBot()
		}

		if err := s.userRepo.Update(ctx, existingUser); err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	var newUser *domain.User
	if isBot {
		newUser = domain.NewBotUser(userData.ID, userData.Login, userData.Email)
	} else {
		newUser = domain.NewUser(userData.ID, userData.Login, userData.Email)
	}

	newUser.IsAdmin = isAdmin

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func generateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func extractScope(token *oauth2.Token) string {
	if scope, ok := token.Extra("scope").(string); ok {
		return scope
	}
	return ""
}