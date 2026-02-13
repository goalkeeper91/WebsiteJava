package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository/postgres"
)

type DiscordAuthService struct {
	connectionRepo *postgres.DiscordConnectionRepository
	clientID       string
	clientSecret   string
	redirectURL    string
}

func NewDiscordAuthService(
	connectionRepo *postgres.DiscordConnectionRepository,
	clientID string,
	clientSecret string,
	redirectURL string,
) *DiscordAuthService {
	return &DiscordAuthService{
		connectionRepo: connectionRepo,
		clientID:       clientID,
		clientSecret:   clientSecret,
		redirectURL:    redirectURL,
	}
}

func (s *DiscordAuthService) GetAuthURL() (string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}

	params := url.Values{}
	params.Set("client_id", s.clientID)
	params.Set("redirect_uri", s.redirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "identify guilds")
	params.Set("state", state)

	authURL := "https://discord.com/api/oauth2/authorize?" + params.Encode()

	return authURL, state, nil
}

func (s *DiscordAuthService) HandleCallback(ctx context.Context, code string, userID int64) error {
	token, err := s.exchangeCodeForToken(code)
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}

	discordUser, err := s.getDiscordUser(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get discord user: %w", err)
	}

	existing, err := s.connectionRepo.GetByDiscordUserID(ctx, discordUser.ID)
	if err != nil {
		return fmt.Errorf("failed to check existing connection: %w", err)
	}

	if existing != nil && existing.UserID != userID {
		return fmt.Errorf("this Discord account is already connected to another user")
	}

	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	existingForUser, err := s.connectionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing connection: %w", err)
	}

	if existingForUser != nil {
		err = s.connectionRepo.Update(ctx, userID, domain.DiscordConnectionUpdateInput{
			DiscordUsername:      &discordUser.Username,
			DiscordDiscriminator: &discordUser.Discriminator,
			AccessToken:          &token.AccessToken,
			RefreshToken:         &token.RefreshToken,
			TokenExpiresAt:       &expiresAt,
		})
		if err != nil {
			return fmt.Errorf("failed to update connection: %w", err)
		}
	} else {
		// Create new connection
		_, err = s.connectionRepo.Create(ctx, domain.DiscordConnectionCreateInput{
			UserID:               userID,
			DiscordUserID:        discordUser.ID,
			DiscordUsername:      discordUser.Username,
			DiscordDiscriminator: discordUser.Discriminator,
			AccessToken:          token.AccessToken,
			RefreshToken:         token.RefreshToken,
			TokenExpiresAt:       expiresAt,
		})
		if err != nil {
			return fmt.Errorf("failed to create connection: %w", err)
		}
	}

	return nil
}

func (s *DiscordAuthService) GetConnection(ctx context.Context, userID int64) (*domain.DiscordConnection, error) {
	return s.connectionRepo.GetByUserID(ctx, userID)
}

func (s *DiscordAuthService) RefreshToken(ctx context.Context, userID int64) error {
	connection, err := s.connectionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	if connection == nil {
		return fmt.Errorf("no Discord connection found")
	}

	token, err := s.refreshAccessToken(connection.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	err = s.connectionRepo.Update(ctx, userID, domain.DiscordConnectionUpdateInput{
		AccessToken:    &token.AccessToken,
		RefreshToken:   &token.RefreshToken,
		TokenExpiresAt: &expiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to update connection: %w", err)
	}

	return nil
}

func (s *DiscordAuthService) Disconnect(ctx context.Context, userID int64) error {
	return s.connectionRepo.Delete(ctx, userID)
}

func (s *DiscordAuthService) IsConnected(ctx context.Context, userID int64) (bool, error) {
	connection, err := s.connectionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	return connection != nil, nil
}

func (s *DiscordAuthService) exchangeCodeForToken(code string) (*discordTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.redirectURL)

	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord token error: %s", string(body))
	}

	var token discordTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *DiscordAuthService) refreshAccessToken(refreshToken string) (*discordTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord refresh error: %s", string(body))
	}

	var token discordTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *DiscordAuthService) getDiscordUser(accessToken string) (*discordUser, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord user error: %s", string(body))
	}

	var user discordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

type discordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type discordUser struct {
	ID            int64  `json:"id,string"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
}