package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// maxChattersPages caps "Get Chatters" pagination so a broken/looping cursor
// from Twitch can never hang the Loyalty accrual tick forever - 5 pages of
// up to 1000 chatters each covers every realistically-sized channel on this
// platform.
const maxChattersPages = 5

// ChattersClient calls Helix endpoints that require a per-user (not app)
// access token - unlike TwitchAppTokenClient, there's no token of its own to
// manage here, the caller supplies the already-decrypted broadcaster token.
type ChattersClient struct {
	clientID   string
	httpClient *http.Client
}

func NewChattersClient(clientID string) *ChattersClient {
	return &ChattersClient{
		clientID:   clientID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Chatter is one viewer currently present in a channel's chat, as returned by
// Twitch's "Get Chatters" endpoint - this includes lurkers who never type,
// which is exactly why it's used for watchtime instead of message-derived
// presence.
type Chatter struct {
	UserID    string
	UserLogin string
}

// GetChatters returns everyone currently in broadcasterID's chat. Uses
// self-moderation (moderator_id = broadcaster_id), the same pattern already
// used for message deletion/timeouts - the broadcaster's own token always
// qualifies as its own moderator.
func (c *ChattersClient) GetChatters(ctx context.Context, broadcasterID, accessToken string) ([]Chatter, error) {
	var chatters []Chatter
	cursor := ""

	for page := 0; page < maxChattersPages; page++ {
		params := url.Values{}
		params.Set("broadcaster_id", broadcasterID)
		params.Set("moderator_id", broadcasterID)
		params.Set("first", "1000")
		if cursor != "" {
			params.Set("after", cursor)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.twitch.tv/helix/chat/chatters?"+params.Encode(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Client-Id", c.clientID)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Get-Chatters-Aufruf: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("twitch get chatters error: status=%d", resp.StatusCode)
		}

		var chattersResp struct {
			Data []struct {
				UserID    string `json:"user_id"`
				UserLogin string `json:"user_login"`
			} `json:"data"`
			Pagination struct {
				Cursor string `json:"cursor"`
			} `json:"pagination"`
		}
		err = json.NewDecoder(resp.Body).Decode(&chattersResp)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("fehler beim Decodieren der Chatters-Antwort: %w", err)
		}

		for _, ch := range chattersResp.Data {
			chatters = append(chatters, Chatter{UserID: ch.UserID, UserLogin: ch.UserLogin})
		}

		if chattersResp.Pagination.Cursor == "" {
			break
		}
		cursor = chattersResp.Pagination.Cursor
	}

	return chatters, nil
}
