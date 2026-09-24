// Package composio is a deliberately tiny client for the handful of Composio
// REST endpoints the TikTok connect page needs (https://docs.composio.dev,
// API v3.1): create a Connect Link, look up a visitor's connected account,
// execute one tool, delete a connected account.
//
// It exists (instead of the Composio SDK) because the API key must only ever
// live server-side - the browser talks to our own /api/tiktok/* endpoints,
// never to Composio directly.
package composio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://backend.composio.dev/api/v3.1"

const (
	requestTimeout = 20 * time.Second
	maxBodyBytes   = 1 << 20 // 1 MiB - far above any legitimate response here
)

// Client talks to the Composio REST API with a project API key.
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		http:    &http.Client{Timeout: requestTimeout},
	}
}

// Configured reports whether an API key is present. The rest of the
// backend deliberately boots without one (like Paddle/Discord) - callers
// check this and answer 503 instead.
func (c *Client) Configured() bool {
	return c != nil && c.apiKey != ""
}

// APIError is a non-2xx answer from Composio. Body is truncated and never
// contains the API key (it's only ever sent as a request header).
type APIError struct {
	Status int
	Body   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("composio API %d: %s", e.Status, e.Body)
}

// ConnectLink is the result of CreateConnectLink.
type ConnectLink struct {
	RedirectURL        string
	ConnectedAccountID string
}

// Account is one connected account (a user's authorized link to a toolkit).
type Account struct {
	ID        string
	Status    string
	UserID    string
	CreatedAt string
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	if !c.Configured() {
		return errors.New("composio: API key not configured")
	}

	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("composio request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return fmt.Errorf("composio response read failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := string(raw)
		if len(msg) > 300 {
			msg = msg[:300]
		}
		return &APIError{Status: resp.StatusCode, Body: msg}
	}

	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("composio response decode failed: %w", err)
	}
	return nil
}

// CreateConnectLink starts an OAuth flow for userID against an auth config.
// After the provider consent Composio redirects the browser to callbackURL
// (appending ?status=...&connected_account_id=...).
func (c *Client) CreateConnectLink(ctx context.Context, authConfigID, userID, callbackURL string) (*ConnectLink, error) {
	body := map[string]any{
		"auth_config_id": authConfigID,
		"user_id":        userID,
	}
	if callbackURL != "" {
		body["callback_url"] = callbackURL
	}

	var resp struct {
		RedirectURL        string `json:"redirect_url"`
		ConnectedAccountID string `json:"connected_account_id"`
	}
	if err := c.do(ctx, http.MethodPost, "/connected_accounts/link", nil, body, &resp); err != nil {
		return nil, err
	}
	if resp.RedirectURL == "" {
		return nil, errors.New("composio: response contained no redirect_url")
	}
	return &ConnectLink{RedirectURL: resp.RedirectURL, ConnectedAccountID: resp.ConnectedAccountID}, nil
}

// FindActiveAccount returns userID's newest ACTIVE connected account for the
// toolkit, or (nil, nil) if there is none. The result is re-checked locally
// (user, toolkit, status) so a misbehaving server-side filter can never hand
// one visitor another visitor's account.
func (c *Client) FindActiveAccount(ctx context.Context, userID, toolkitSlug string) (*Account, error) {
	query := url.Values{}
	query.Add("user_ids", userID)
	query.Add("toolkit_slugs", toolkitSlug)
	query.Add("statuses", "ACTIVE")

	var resp struct {
		Items []struct {
			ID        string `json:"id"`
			Status    string `json:"status"`
			UserID    string `json:"user_id"`
			CreatedAt string `json:"created_at"`
			Toolkit   struct {
				Slug string `json:"slug"`
			} `json:"toolkit"`
		} `json:"items"`
	}
	if err := c.do(ctx, http.MethodGet, "/connected_accounts", query, nil, &resp); err != nil {
		return nil, err
	}

	var best *Account
	for _, item := range resp.Items {
		if item.UserID != userID ||
			!strings.EqualFold(item.Toolkit.Slug, toolkitSlug) ||
			!strings.EqualFold(item.Status, "ACTIVE") {
			continue
		}
		candidate := &Account{ID: item.ID, Status: item.Status, UserID: item.UserID, CreatedAt: item.CreatedAt}
		// ISO-8601 timestamps sort lexicographically; without one, keep the first.
		if best == nil || (candidate.CreatedAt != "" && candidate.CreatedAt > best.CreatedAt) {
			best = candidate
		}
	}
	return best, nil
}

// ExecuteTool runs one Composio tool as the given connected account and
// returns the raw `data` payload of the response.
func (c *Client) ExecuteTool(ctx context.Context, toolSlug, userID, connectedAccountID string, arguments map[string]any) (json.RawMessage, error) {
	body := map[string]any{
		"connected_account_id": connectedAccountID,
		"user_id":              userID,
		"arguments":            arguments,
	}

	var resp struct {
		Data       json.RawMessage `json:"data"`
		Error      *string         `json:"error"`
		Successful bool            `json:"successful"`
	}
	if err := c.do(ctx, http.MethodPost, "/tools/execute/"+url.PathEscape(toolSlug), nil, body, &resp); err != nil {
		return nil, err
	}
	if !resp.Successful {
		msg := "unknown error"
		if resp.Error != nil && *resp.Error != "" {
			msg = *resp.Error
		}
		return nil, fmt.Errorf("composio tool %s failed: %s", toolSlug, msg)
	}
	return resp.Data, nil
}

// DeleteAccount removes a connected account and, with revoke=true, also
// revokes the upstream (TikTok) credentials.
func (c *Client) DeleteAccount(ctx context.Context, connectedAccountID string, revoke bool) error {
	query := url.Values{}
	if revoke {
		query.Set("revoke_on_delete", "true")
	}
	return c.do(ctx, http.MethodDelete, "/connected_accounts/"+url.PathEscape(connectedAccountID), query, nil, nil)
}
