// Package twitch holds shared, dependency-free Twitch Helix API clients used
// by more than one service (the clip-detector and the scheduled-message
// scheduler both need to check a channel's live status via an App Access
// Token, so this lives outside any one feature's package).
package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TwitchAppTokenClient manages a Twitch App Access Token (Client Credentials
// Grant). Checking stream live-status doesn't need a per-user token.
type TwitchAppTokenClient struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

func NewTwitchAppTokenClient(clientID, clientSecret string) *TwitchAppTokenClient {
	return &TwitchAppTokenClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *TwitchAppTokenClient) getToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.expiresAt) {
		return c.token, nil
	}

	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fehler beim App-Token-Abruf: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("twitch app token error: status=%d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("fehler beim Decodieren der App-Token-Antwort: %w", err)
	}

	c.token = tokenResp.AccessToken
	// Etwas Puffer vor dem echten Ablauf abziehen, damit wir nie mit einem
	// gerade abgelaufenen Token anfragen.
	c.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)

	return c.token, nil
}

// LiveStream ist ein aktiver Stream, wie ihn Twitch "Get Streams" zurückgibt.
type LiveStream struct {
	UserID    string
	UserLogin string
	Title     string
	GameName  string
}

// GetLiveStreams prüft, welche der übergebenen Twitch User-IDs gerade live
// sind (bis zu 100 pro Aufruf, Twitch-Limit). Nicht enthaltene IDs sind offline.
func (c *TwitchAppTokenClient) GetLiveStreams(ctx context.Context, userIDs []string) (map[string]LiveStream, error) {
	if len(userIDs) == 0 {
		return map[string]LiveStream{}, nil
	}

	token, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	for _, id := range userIDs {
		params.Add("user_id", id)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.twitch.tv/helix/streams?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Get-Streams-Aufruf: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch get streams error: status=%d", resp.StatusCode)
	}

	var streamsResp struct {
		Data []struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			Title     string `json:"title"`
			GameName  string `json:"game_name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&streamsResp); err != nil {
		return nil, fmt.Errorf("fehler beim Decodieren der Streams-Antwort: %w", err)
	}

	result := make(map[string]LiveStream, len(streamsResp.Data))
	for _, s := range streamsResp.Data {
		result[s.UserID] = LiveStream{
			UserID:    s.UserID,
			UserLogin: s.UserLogin,
			Title:     s.Title,
			GameName:  s.GameName,
		}
	}

	return result, nil
}

// TwitchUser ist das Ergebnis einer Login->ID-Auflösung, wie sie Twitch
// "Get Users" zurückgibt.
type TwitchUser struct {
	ID    string
	Login string
}

// GetUserByLogin löst einen Twitch-Login zur numerischen User-ID auf - nutzt
// dasselbe App-Access-Token wie GetLiveStreams, sodass die aufgelöste Person
// sich vor der Auflösung noch nie auf der Plattform eingeloggt haben muss
// (Twitch kennt sie ohnehin). Gibt (nil, nil) zurück, wenn der Login nicht
// existiert - kein Fehler, einfach "nicht gefunden".
func (c *TwitchAppTokenClient) GetUserByLogin(ctx context.Context, login string) (*TwitchUser, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.twitch.tv/helix/users?login="+url.QueryEscape(strings.ToLower(login)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Get-Users-Aufruf: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch get users error: status=%d", resp.StatusCode)
	}

	var usersResp struct {
		Data []struct {
			ID    string `json:"id"`
			Login string `json:"login"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&usersResp); err != nil {
		return nil, fmt.Errorf("fehler beim Decodieren der Users-Antwort: %w", err)
	}
	if len(usersResp.Data) == 0 {
		return nil, nil
	}

	return &TwitchUser{ID: usersResp.Data[0].ID, Login: usersResp.Data[0].Login}, nil
}
