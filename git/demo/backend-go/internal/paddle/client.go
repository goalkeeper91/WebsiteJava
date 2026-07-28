// Package paddle is a small, dependency-free outbound client for Paddle's
// Billing API - mirrors the shape of internal/twitch's Helix clients
// (base URL + auth header + JSON decode into a clean domain type), but
// without token-refresh caching since Paddle uses a static Bearer API key.
package paddle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type portalSessionResponse struct {
	Data struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"data"`
}

// GetOrCreatePortalSession mints a fresh, single-use, short-lived Paddle
// Customer Portal link for the given Paddle customer - per Paddle's docs
// these URLs must never be cached, so this is always a live API call, never
// backed by a stored value.
func (c *Client) GetOrCreatePortalSession(ctx context.Context, customerID string) (string, error) {
	url := fmt.Sprintf("%s/customers/%s/portal-sessions", c.baseURL, customerID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("fehler beim Erstellen der Paddle-Portal-Anfrage: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fehler beim Aufruf der Paddle-API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("paddle-api antwortete mit status %d", resp.StatusCode)
	}

	var parsed portalSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("fehler beim Parsen der Paddle-Antwort: %w", err)
	}

	return parsed.Data.URL, nil
}
