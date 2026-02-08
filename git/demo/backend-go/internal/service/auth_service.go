package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type AuthService struct {
	userRepo      repository.UserRepository
	tokenRepo     repository.AuthTokenRepository
	oauthConfig   *oauth2.Config
	frontendURL   string
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.AuthTokenRepository,
	oauthConfig *oauth2.Config,
	frontendURL string,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		oauthConfig: oauthConfig,
		frontendURL: frontendURL,
	}
}

func (s *AuthService) GetAuthURL() (string, string, error) {
	state, err := generateStateToken()
	if err != nil {
		return "", "", fmt.Errorf("fehler beim Generieren des State Tokens: %w", err)
	}

	url := s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url, state, nil
}

func (s *AuthService) HandleCallback(ctx context.Context, code string) (*domain.User, error) {
	token, err := s.oauthConfig.Exchange(ctx, code)
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

	user, err := s.upsertUser(ctx, twitchUserResp)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Speichern des Users: %w", err)
	}

	authToken := domain.NewAuthToken(
		userData.ID,
		userData.Login,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		extractScope(token),
		int64(token.Expiry.Sub(time.Now()).Seconds()),
	)

	if err := s.tokenRepo.Upsert(ctx, authToken); err != nil {
		return nil, fmt.Errorf("fehler beim Speichern des Auth Tokens: %w", err)
	}

	return user, nil
}

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

func (s *AuthService) upsertUser(ctx context.Context, twitchUserResp *TwitchUserResponse) (*domain.User, error) {
	if len(twitchUserResp.Data) == 0 {
		return nil, fmt.Errorf("keine User-Daten in Twitch Response")
	}

	userData := twitchUserResp.Data[0]

	existingUser, err := s.userRepo.GetByTwitchID(ctx, userData.ID)
	if err == nil {
		existingUser.Update(userData.Login, userData.Email)
		if err := s.userRepo.Update(ctx, existingUser); err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	newUser := domain.NewUser(userData.ID, userData.Login, userData.Email)
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