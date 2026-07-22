package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

const (
	clipPollInterval = 3 * time.Second
	clipPollAttempts = 10
)

// ClipCreationService turns a clip-detector trigger into an actual Twitch
// clip, tracks it via ClipLog until it's ready, and hands it off to n8n for
// distribution (Phase D).
type ClipCreationService struct {
	clipRepo                    repository.ClipLogRepository
	tokenRepo                   repository.AuthTokenRepository
	settingsRepo                repository.AutomationSettingsRepository
	authService                 *AuthService
	n8nService                  *N8NService
	discordNotificationService  *DiscordNotificationService
	captionService              *CaptionService
	clientID                    string
	httpClient                  *http.Client
}

func NewClipCreationService(
	clipRepo repository.ClipLogRepository,
	tokenRepo repository.AuthTokenRepository,
	settingsRepo repository.AutomationSettingsRepository,
	authService *AuthService,
	n8nService *N8NService,
	discordNotificationService *DiscordNotificationService,
	captionService *CaptionService,
	clientID string,
) *ClipCreationService {
	return &ClipCreationService{
		clipRepo:                   clipRepo,
		tokenRepo:                  tokenRepo,
		settingsRepo:               settingsRepo,
		authService:                authService,
		n8nService:                 n8nService,
		discordNotificationService: discordNotificationService,
		captionService:             captionService,
		clientID:                   clientID,
		httpClient:                 &http.Client{Timeout: 15 * time.Second},
	}
}

// ProcessTrigger handles a single clips:trigger message end to end. It never
// returns an error — everything is logged, since there's no HTTP caller
// waiting on this (it's dispatched from a Redis listener).
func (s *ClipCreationService) ProcessTrigger(ctx context.Context, userTwitchID, reason string) {
	log.Printf("🎬 Verarbeite Clip-Trigger für %s (Grund: %s)", userTwitchID, reason)

	busy, err := s.hasActiveClip(ctx, userTwitchID)
	if err != nil {
		log.Printf("⚠️ Fehler beim Debounce-Check für %s: %v", userTwitchID, err)
		return
	}
	if busy {
		log.Printf("⏭️ Überspringe Trigger für %s — es läuft bereits ein Clip-Vorgang", userTwitchID)
		return
	}

	settings, err := s.settingsRepo.GetByUser(ctx, userTwitchID)
	if err != nil || settings == nil || !settings.IsEnabled {
		log.Printf("⏭️ Automatisierung für %s nicht aktiv, überspringe Trigger", userTwitchID)
		return
	}

	accessToken, err := s.getValidAccessToken(ctx, userTwitchID)
	if err != nil {
		log.Printf("⚠️ Kein gültiger Token für %s: %v", userTwitchID, err)
		return
	}

	clipID, err := s.createClip(ctx, userTwitchID, accessToken)
	if err != nil {
		log.Printf("❌ Clip-Erstellung fehlgeschlagen für %s: %v", userTwitchID, err)
		return
	}

	clipLog := domain.NewClipLog(userTwitchID, clipID, fmt.Sprintf("https://clips.twitch.tv/%s", clipID))
	clipLog.MarkAsProcessing("", "")

	if err := s.clipRepo.Create(ctx, clipLog); err != nil {
		log.Printf("❌ Fehler beim Speichern des Clip-Logs für %s: %v", userTwitchID, err)
		return
	}

	clipData, err := s.waitForClipReady(ctx, clipID, accessToken)
	if err != nil {
		log.Printf("⚠️ Clip %s wurde nicht rechtzeitig fertig: %v", clipID, err)
		clipLog.MarkAsFailed("Clip-Verarbeitung bei Twitch hat zu lange gedauert")
		if uerr := s.clipRepo.Update(ctx, clipLog); uerr != nil {
			log.Printf("❌ Fehler beim Aktualisieren des Clip-Logs %s: %v", clipID, uerr)
		}
		return
	}

	s.fillClipDetails(clipLog, clipData)
	clipLog.GameName = s.getGameName(ctx, clipData.GameID, accessToken)
	clipLog.MarkAsCompleted()

	if err := s.clipRepo.Update(ctx, clipLog); err != nil {
		log.Printf("❌ Fehler beim Aktualisieren des Clip-Logs %s: %v", clipID, err)
		return
	}
	log.Printf("✅ Clip fertig: %s (%s)", clipID, clipLog.ClipURL)

	generated := s.captionService.GenerateCaption(ctx, clipLog.ClipTitle, clipLog.GameName, settings)
	clipLog.SetAIContent(generated.Caption, generated.Hashtags)
	if err := s.clipRepo.Update(ctx, clipLog); err != nil {
		log.Printf("⚠️ Fehler beim Speichern der Caption für %s: %v", clipID, err)
	}

	if containsPlatform(settings.TargetPlatforms, domain.PlatformDiscord) {
		s.postToDiscord(ctx, userTwitchID, clipLog, generated)
	}
}

// postToDiscord announces the clip in the user's configured Discord
// notification channel(s), reusing the existing bot infrastructure — no n8n
// needed for this v1 target.
func (s *ClipCreationService) postToDiscord(ctx context.Context, userTwitchID string, clipLog *domain.ClipLog, generated GeneratedCaption) {
	posted, err := s.discordNotificationService.SendClipNotification(ctx, userTwitchID, clipLog.ClipTitle, clipLog.ClipURL, generated.Caption, generated.Hashtags)
	if err != nil {
		log.Printf("⚠️ Discord-Post fehlgeschlagen für %s: %v", userTwitchID, err)
		return
	}
	if !posted {
		log.Printf("ℹ️ Kein Discord-Notification-Channel konfiguriert für %s, überspringe Post", userTwitchID)
		return
	}

	clipLog.AddPlatformPost(domain.PlatformDiscord, "", clipLog.ClipURL)
	if err := s.clipRepo.Update(ctx, clipLog); err != nil {
		log.Printf("❌ Fehler beim Aktualisieren des Clip-Logs nach Discord-Post für %s: %v", userTwitchID, err)
	}
}

func containsPlatform(platforms domain.PlatformArray, target domain.Platform) bool {
	for _, p := range platforms {
		if p == target {
			return true
		}
	}
	return false
}

func (s *ClipCreationService) hasActiveClip(ctx context.Context, userTwitchID string) (bool, error) {
	for _, status := range []domain.ClipStatus{domain.StatusPending, domain.StatusProcessing} {
		clips, err := s.clipRepo.GetByUserAndStatus(ctx, userTwitchID, status, 1, 0)
		if err != nil {
			return false, err
		}
		if len(clips) > 0 {
			return true, nil
		}
	}
	return false, nil
}

// getValidAccessToken loads the user's stored Twitch token, refreshes it if
// expired, and verifies it actually carries the clips:edit scope.
func (s *ClipCreationService) getValidAccessToken(ctx context.Context, userTwitchID string) (string, error) {
	token, err := s.tokenRepo.GetByTwitchUserID(ctx, userTwitchID)
	if err != nil {
		return "", fmt.Errorf("kein gespeicherter Token: %w", err)
	}

	if token.IsExpired() {
		if err := s.authService.RefreshToken(ctx, userTwitchID); err != nil {
			return "", fmt.Errorf("token-refresh fehlgeschlagen: %w", err)
		}
		token, err = s.tokenRepo.GetByTwitchUserID(ctx, userTwitchID)
		if err != nil {
			return "", fmt.Errorf("token nach refresh nicht ladbar: %w", err)
		}
	}

	if !strings.Contains(token.Scope, "clips:edit") {
		return "", fmt.Errorf("token hat keine clips:edit-Berechtigung — bitte einmal aus- und wieder einloggen")
	}

	return token.AccessToken, nil
}

func (s *ClipCreationService) createClip(ctx context.Context, userTwitchID, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.twitch.tv/helix/clips?broadcaster_id="+url.QueryEscape(userTwitchID), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", s.clientID)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("create-clip request fehlgeschlagen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("twitch create-clip error: status=%d body=%s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			ID      string `json:"id"`
			EditURL string `json:"edit_url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("fehler beim Decodieren der Create-Clip-Antwort: %w", err)
	}
	if len(result.Data) == 0 {
		return "", fmt.Errorf("twitch lieferte keine Clip-ID zurück")
	}

	return result.Data[0].ID, nil
}

type twitchClipData struct {
	URL             string  `json:"url"`
	BroadcasterName string  `json:"broadcaster_name"`
	GameID          string  `json:"game_id"`
	Title           string  `json:"title"`
	ViewCount       int     `json:"view_count"`
	CreatedAt       string  `json:"created_at"`
	ThumbnailURL    string  `json:"thumbnail_url"`
	Duration        float64 `json:"duration"`
}

// waitForClipReady polls Get Clip until Twitch has finished rendering the
// clip (signaled by a non-empty thumbnail_url), per Twitch's own guidance.
func (s *ClipCreationService) waitForClipReady(ctx context.Context, clipID, accessToken string) (*twitchClipData, error) {
	for attempt := 0; attempt < clipPollAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(clipPollInterval):
		}

		clip, err := s.getClip(ctx, clipID, accessToken)
		if err != nil {
			return nil, err
		}
		if clip != nil && clip.ThumbnailURL != "" {
			return clip, nil
		}
	}
	return nil, fmt.Errorf("clip nach %d Versuchen nicht fertig", clipPollAttempts)
}

func (s *ClipCreationService) getClip(ctx context.Context, clipID, accessToken string) (*twitchClipData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.twitch.tv/helix/clips?id="+url.QueryEscape(clipID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", s.clientID)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get-clip request fehlgeschlagen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("twitch get-clip error: status=%d body=%s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []twitchClipData `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("fehler beim Decodieren der Get-Clip-Antwort: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, nil
	}

	return &result.Data[0], nil
}

// getGameName resolves a Twitch game_id to its display name (best effort —
// an empty result just means the dashboard shows no game for this clip).
func (s *ClipCreationService) getGameName(ctx context.Context, gameID, accessToken string) string {
	if gameID == "" {
		return ""
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.twitch.tv/helix/games?id="+url.QueryEscape(gameID), nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", s.clientID)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result struct {
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result.Data) == 0 {
		return ""
	}

	return result.Data[0].Name
}

func (s *ClipCreationService) fillClipDetails(clipLog *domain.ClipLog, clip *twitchClipData) {
	clipLog.ClipURL = clip.URL
	clipLog.ClipTitle = clip.Title
	clipLog.BroadcasterName = clip.BroadcasterName
	clipLog.DurationSeconds = int(clip.Duration)
	clipLog.ViewCount = clip.ViewCount

	if createdAt, err := time.Parse(time.RFC3339, clip.CreatedAt); err == nil {
		clipLog.CreatedAtTwitch = &createdAt
	}
}

// postToN8N hands the completed clip off for distribution. Missing/not-ready
// n8n integration is a soft skip, not a failure — the clip itself is already
// tracked as completed regardless of whether it gets posted anywhere.
func (s *ClipCreationService) postToN8N(ctx context.Context, userTwitchID string, clipLog *domain.ClipLog, settings *domain.AutomationSettings) {
	platforms := filterLinkSharePlatforms(settings.TargetPlatforms)
	if len(platforms) == 0 {
		log.Printf("ℹ️ Keine Link-Sharing-Plattformen für %s konfiguriert, überspringe n8n-Übergabe", userTwitchID)
		return
	}

	payload := N8NClipPostPayload{
		ClipID:          clipLog.ClipID,
		ClipURL:         clipLog.ClipURL,
		Title:           clipLog.ClipTitle,
		GameName:        clipLog.GameName,
		DurationSeconds: clipLog.DurationSeconds,
		AITone:          string(settings.AITone),
		AIStyle:         string(settings.AIStyle),
		UseHashtags:     settings.UseHashtags,
		MaxHashtags:     settings.MaxHashtags,
		TargetPlatforms: platforms,
	}

	if _, err := s.n8nService.TriggerClipPostWebhook(ctx, userTwitchID, payload); err != nil {
		if errors.Is(err, domain.ErrN8NIntegrationNotReady) {
			log.Printf("ℹ️ n8n nicht konfiguriert für %s, Clip bleibt ohne automatisches Posting", userTwitchID)
			return
		}
		log.Printf("⚠️ n8n-Übergabe fehlgeschlagen für %s: %v", userTwitchID, err)
	}
}

// filterLinkSharePlatforms enforces the v1 distribution scope server-side
// (discord/twitter only), regardless of what's stored in AutomationSettings.
func filterLinkSharePlatforms(platforms domain.PlatformArray) []string {
	result := make([]string, 0, len(platforms))
	for _, p := range platforms {
		if p == domain.PlatformDiscord || p == domain.PlatformTwitter {
			result = append(result, string(p))
		}
	}
	return result
}
