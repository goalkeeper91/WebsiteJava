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
	"strconv"
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

// portalSessionResponse mirrors Paddle's real response shape for
// "Create a customer portal session" - there is no flat data.url field;
// links are nested under data.urls.general.overview (portal homepage) and
// data.urls.subscriptions[].cancel_subscription (deep link straight to a
// specific subscription's cancellation page).
type portalSessionResponse struct {
	Data struct {
		ID   string `json:"id"`
		URLs struct {
			General struct {
				Overview string `json:"overview"`
			} `json:"general"`
			Subscriptions []struct {
				ID                 string `json:"id"`
				CancelSubscription string `json:"cancel_subscription"`
			} `json:"subscriptions"`
		} `json:"urls"`
	} `json:"data"`
}

// GetOrCreatePortalSession mints a fresh, single-use, short-lived Paddle
// Customer Portal link for the given Paddle customer - per Paddle's docs
// these URLs must never be cached, so this is always a live API call, never
// backed by a stored value. If subscriptionID matches one of the returned
// per-subscription deep links, that direct cancellation-page link is
// returned instead of the generic portal overview - saves the customer a
// click and avoids ambiguity if they ever have more than one subscription.
func (c *Client) GetOrCreatePortalSession(ctx context.Context, customerID, subscriptionID string) (string, error) {
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

	for _, s := range parsed.Data.URLs.Subscriptions {
		if s.ID == subscriptionID {
			return s.CancelSubscription, nil
		}
	}

	return parsed.Data.URLs.General.Overview, nil
}

// priceResponse mirrors Paddle's "Get a price" response - unit_price.amount
// is a string holding the price in the currency's lowest denomination (e.g.
// "499" for 4.99 EUR), not a float, to avoid floating-point rounding on
// Paddle's side.
type priceResponse struct {
	Data struct {
		ID        string `json:"id"`
		UnitPrice struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"unit_price"`
	} `json:"data"`
}

// GetPrice fetches the current price configured in Paddle for priceID and
// returns it in major currency units (e.g. 4.99, not 499). Assumes a
// 2-decimal currency (EUR, the only currency this app's prices are quoted
// in) - does not special-case zero-decimal currencies like JPY.
func (c *Client) GetPrice(ctx context.Context, priceID string) (float64, error) {
	url := fmt.Sprintf("%s/prices/%s", c.baseURL, priceID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Erstellen der Paddle-Preis-Anfrage: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Aufruf der Paddle-API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("paddle-api antwortete mit status %d beim Laden von price %s", resp.StatusCode, priceID)
	}

	var parsed priceResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, fmt.Errorf("fehler beim Parsen der Paddle-Preis-Antwort: %w", err)
	}

	amountMinorUnits, err := strconv.ParseInt(parsed.Data.UnitPrice.Amount, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Parsen des Paddle-Preisbetrags %q: %w", parsed.Data.UnitPrice.Amount, err)
	}

	return float64(amountMinorUnits) / 100, nil
}
