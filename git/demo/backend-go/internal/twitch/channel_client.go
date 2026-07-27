package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ChannelClient calls Helix endpoints that need the broadcaster's own
// per-user access token - same style as ChattersClient, a small ad-hoc
// client per concern rather than one shared generic Helix abstraction.
type ChannelClient struct {
	clientID   string
	httpClient *http.Client
}

func NewChannelClient(clientID string) *ChannelClient {
	return &ChannelClient{
		clientID:   clientID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ChannelInfo is a channel's current title/category, as returned by
// "Get Channel Information".
type ChannelInfo struct {
	BroadcasterID    string
	BroadcasterLogin string
	Title            string
	GameID           string
	GameName         string
}

func (c *ChannelClient) GetChannelInfo(ctx context.Context, broadcasterID, accessToken string) (*ChannelInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.twitch.tv/helix/channels?broadcaster_id="+url.QueryEscape(broadcasterID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Get-Channel-Aufruf: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch get channel error: status=%d", resp.StatusCode)
	}

	var channelResp struct {
		Data []struct {
			BroadcasterID    string `json:"broadcaster_id"`
			BroadcasterLogin string `json:"broadcaster_login"`
			Title            string `json:"title"`
			GameID           string `json:"game_id"`
			GameName         string `json:"game_name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&channelResp); err != nil {
		return nil, fmt.Errorf("fehler beim Decodieren der Channel-Antwort: %w", err)
	}
	if len(channelResp.Data) == 0 {
		return nil, fmt.Errorf("kein Channel gefunden für broadcaster_id=%s", broadcasterID)
	}

	d := channelResp.Data[0]
	return &ChannelInfo{
		BroadcasterID:    d.BroadcasterID,
		BroadcasterLogin: d.BroadcasterLogin,
		Title:            d.Title,
		GameID:           d.GameID,
		GameName:         d.GameName,
	}, nil
}

// ModifyChannelInfo updates title and/or category - both optional, only
// non-nil fields are sent to Twitch.
func (c *ChannelClient) ModifyChannelInfo(ctx context.Context, broadcasterID, accessToken string, title, gameID *string) error {
	body := map[string]string{}
	if title != nil {
		body["title"] = *title
	}
	if gameID != nil {
		body["game_id"] = *gameID
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		"https://api.twitch.tv/helix/channels?broadcaster_id="+url.QueryEscape(broadcasterID), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fehler beim Modify-Channel-Aufruf: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("twitch modify channel error: status=%d", resp.StatusCode)
	}

	return nil
}

// GetFollowerCount returns the channel's total follower count via
// self-moderation (moderator_id = broadcaster_id, same pattern as
// ChattersClient.GetChatters) - moderator:read:followers already covers
// this, no new scope needed. first=1 keeps the payload minimal since only
// the total is needed.
func (c *ChannelClient) GetFollowerCount(ctx context.Context, broadcasterID, accessToken string) (int, error) {
	params := url.Values{}
	params.Set("broadcaster_id", broadcasterID)
	params.Set("moderator_id", broadcasterID)
	params.Set("first", "1")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.twitch.tv/helix/channels/followers?"+params.Encode(), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Get-Followers-Aufruf: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("twitch get followers error: status=%d", resp.StatusCode)
	}

	var followersResp struct {
		Total int `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&followersResp); err != nil {
		return 0, fmt.Errorf("fehler beim Decodieren der Followers-Antwort: %w", err)
	}

	return followersResp.Total, nil
}

// GetSubscriberCount returns the channel's total subscriber count -
// channel:read:subscriptions already covers this, no new scope needed.
func (c *ChannelClient) GetSubscriberCount(ctx context.Context, broadcasterID, accessToken string) (int, error) {
	params := url.Values{}
	params.Set("broadcaster_id", broadcasterID)
	params.Set("first", "1")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.twitch.tv/helix/subscriptions?"+params.Encode(), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Get-Subscriptions-Aufruf: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("twitch get subscriptions error: status=%d", resp.StatusCode)
	}

	var subsResp struct {
		Total int `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&subsResp); err != nil {
		return 0, fmt.Errorf("fehler beim Decodieren der Subscriptions-Antwort: %w", err)
	}

	return subsResp.Total, nil
}

// CommercialResult is Twitch's response to "Start Commercial" - Length is
// the actual granted length in seconds (may differ from what was
// requested), RetryAfter is the cooldown in seconds before another
// commercial can run.
type CommercialResult struct {
	Length     int    `json:"length"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retryAfter"`
}

// StartCommercial triggers an ad break. Deliberately does not pre-validate
// eligibility (Partner/Affiliate status, cooldown) - Twitch's own error
// response is surfaced to the caller as-is.
func (c *ChannelClient) StartCommercial(ctx context.Context, broadcasterID, accessToken string, lengthSeconds int) (*CommercialResult, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"broadcaster_id": broadcasterID,
		"length":         lengthSeconds,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.twitch.tv/helix/channels/commercial", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Start-Commercial-Aufruf: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch start commercial error: status=%d", resp.StatusCode)
	}

	// Twitch's raw response uses retry_after (snake_case) - decoded
	// separately from CommercialResult, whose retryAfter tag matches this
	// package's other frontend-facing types instead.
	var commercialResp struct {
		Data []struct {
			Length     int    `json:"length"`
			Message    string `json:"message"`
			RetryAfter int    `json:"retry_after"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&commercialResp); err != nil {
		return nil, fmt.Errorf("fehler beim Decodieren der Commercial-Antwort: %w", err)
	}
	if len(commercialResp.Data) == 0 {
		return nil, fmt.Errorf("leere Antwort von Twitch")
	}

	d := commercialResp.Data[0]
	return &CommercialResult{Length: d.Length, Message: d.Message, RetryAfter: d.RetryAfter}, nil
}
