package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"demo/backend-go/internal/domain"
)

// CaptionService generates short social-media captions for finished clips via
// a self-hosted Ollama instance — no external API costs, no per-user n8n
// dependency. AI captioning is a nice-to-have, never a hard dependency: any
// failure falls back to the clip's own Twitch title.
type CaptionService struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewCaptionService(baseURL, model string) *CaptionService {
	return &CaptionService{
		baseURL:    baseURL,
		model:      model,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

type GeneratedCaption struct {
	Caption  string
	Hashtags []string
}

// GenerateCaption asks Ollama for a caption. On any problem (unreachable,
// timeout, unexpected output) it returns the clip title as a plain fallback.
func (s *CaptionService) GenerateCaption(ctx context.Context, clipTitle, gameName string, settings *domain.AutomationSettings) GeneratedCaption {
	fallback := GeneratedCaption{Caption: clipTitle}

	if s.baseURL == "" {
		return fallback
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"model":  s.model,
		"prompt": buildCaptionPrompt(clipTitle, gameName, settings),
		"stream": false,
	})
	if err != nil {
		return fallback
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/api/generate", bytes.NewReader(reqBody))
	if err != nil {
		return fallback
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fallback
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fallback
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fallback
	}

	caption, hashtags, ok := parseCaptionResponse(result.Response)
	if !ok {
		return fallback
	}

	if !settings.UseHashtags {
		hashtags = nil
	} else if len(hashtags) > settings.MaxHashtags {
		hashtags = hashtags[:settings.MaxHashtags]
	}

	return GeneratedCaption{Caption: caption, Hashtags: hashtags}
}

func buildCaptionPrompt(clipTitle, gameName string, settings *domain.AutomationSettings) string {
	game := gameName
	if game == "" {
		game = "einem Spiel"
	}

	return fmt.Sprintf(
		`Du schreibst eine kurze Social-Media-Caption fuer einen Twitch-Clip. Tonalitaet: %s. Stil: %s.
Clip-Titel: %q
Spiel: %s
Antworte NUR in exakt diesem Format, ohne weitere Erklaerungen:
CAPTION: <ein bis zwei kurze Saetze>
HASHTAGS: #tag1 #tag2 #tag3`,
		settings.AITone, settings.AIStyle, clipTitle, game,
	)
}

// parseCaptionResponse prefers the requested "CAPTION: ...\nHASHTAGS: ..."
// format, but small models don't reliably stick to it — live testing against
// llama3.2:3b showed it happily ignoring the format and writing free-form
// text with hashtags scattered inline instead. So this falls back to using
// the whole response as the caption and pulls out any "#word" tokens found
// anywhere in it as hashtags (removed from the caption body to avoid them
// appearing twice once we append them separately).
func parseCaptionResponse(raw string) (caption string, hashtags []string, ok bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil, false
	}

	body := raw
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "CAPTION:") {
			body = strings.TrimSpace(strings.TrimPrefix(trimmed, "CAPTION:"))
			break
		}
	}

	// Hashtags können in einer eigenen HASHTAGS:-Zeile stehen (striktes
	// Format) oder irgendwo im Fließtext auftauchen (kleine Modelle halten
	// sich oft nicht ans angeforderte Format) — wir suchen deshalb in der
	// gesamten Rohantwort danach, nicht nur im Caption-Teil.
	for _, word := range strings.Fields(raw) {
		clean := strings.Trim(word, ".,!?\"'()")
		if strings.HasPrefix(clean, "#") && len(clean) > 1 {
			hashtags = append(hashtags, clean)
		}
	}

	bodyWords := strings.Fields(body)
	captionWords := make([]string, 0, len(bodyWords))
	for _, word := range bodyWords {
		clean := strings.Trim(word, ".,!?\"'()")
		if strings.HasPrefix(clean, "#") && len(clean) > 1 {
			continue // schon als Hashtag erfasst, nicht doppelt in der Caption
		}
		captionWords = append(captionWords, word)
	}

	caption = strings.TrimSpace(strings.Join(captionWords, " "))
	if caption == "" {
		caption = body
	}

	return caption, hashtags, true
}
