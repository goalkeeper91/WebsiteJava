package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type ContactHandler struct {
	contactRepo      repository.ContactRequestRepository
	discordBotToken  string
	discordChannelID string
}

// discordBotToken/discordChannelID are DISCORD_BOT_TOKEN/DISCORD_ADMIN_CONTACT_CHANNEL
// (pkg/config/config.go's DiscordConfig - both fields already existed, unused,
// clearly scaffolded for exactly this feature already). Deliberately the bot's
// own token + Discord's REST API directly, not the separate discord-bot
// container (bot-plattform-discord-bot, a Python service for unrelated
// Twitch-command notifications) - no reason to touch that working service
// just to post one admin notification.
func NewContactHandler(contactRepo repository.ContactRequestRepository, discordBotToken, discordChannelID string) *ContactHandler {
	return &ContactHandler{
		contactRepo:      contactRepo,
		discordBotToken:  discordBotToken,
		discordChannelID: discordChannelID,
	}
}

func (h *ContactHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/contact", h.SubmitContact).Methods("POST")
}

type ContactRequestDTO struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Subject      string `json:"subject"`
	Message      string `json:"message"`
	ConsentGiven bool   `json:"consentGiven"`
}

func (h *ContactHandler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	var dto ContactRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	contactRequest := domain.NewContactRequest(
		dto.Name,
		dto.Email,
		dto.Phone,
		dto.Subject,
		dto.Message,
		dto.ConsentGiven,
	)

	if err := contactRequest.Validate(); err != nil {
		http.Error(w, "Fehlende oder ungültige Felder", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.contactRepo.Create(ctx, contactRequest); err != nil {
		log.Printf("Fehler beim Speichern der Contact Request: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Contact Request erhalten von: %s (%s)", dto.Name, dto.Email)

	// Fire-and-forget: eine langsame/fehlerhafte Discord-API darf die
	// Antwort an den Besucher niemals verzögern oder zum Scheitern bringen -
	// die Anfrage ist zu diesem Zeitpunkt bereits sicher gespeichert.
	go h.notifyDiscord(contactRequest)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contactRequest)
}

// notifyDiscord postet eine neue Kontaktanfrage als Embed in einen privaten
// Discord-Channel, direkt über Discords REST-API mit dem Bot-Token. No-op,
// solange Token/Channel nicht konfiguriert sind (DISCORD_BOT_TOKEN /
// DISCORD_ADMIN_CONTACT_CHANNEL) - das Kontaktformular funktioniert auch
// ohne diese Konfiguration weiterhin normal, nur eben ohne Benachrichtigung.
func (h *ContactHandler) notifyDiscord(c *domain.ContactRequest) {
	if h.discordBotToken == "" || h.discordChannelID == "" {
		return
	}

	phone := c.Phone
	if phone == "" {
		phone = "–"
	}

	payload := map[string]any{
		"embeds": []map[string]any{
			{
				"title": "📬 Neue Kontaktanfrage – dev.goalkeeper91.de",
				"color": 0x0058e6,
				"fields": []map[string]any{
					{"name": "Name", "value": c.Name, "inline": true},
					{"name": "E-Mail", "value": c.Email, "inline": true},
					{"name": "Telefon", "value": phone, "inline": true},
					{"name": "Betreff", "value": c.Subject},
					{"name": "Nachricht", "value": truncateForDiscord(c.Message, 1000)},
				},
				"timestamp": c.CreatedAt.Format(time.RFC3339),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("⚠️ Discord-Benachrichtigung: Payload-Fehler: %v", err)
		return
	}

	url := "https://discord.com/api/v10/channels/" + h.discordChannelID + "/messages"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Printf("⚠️ Discord-Benachrichtigung: Request-Fehler: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bot "+h.discordBotToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️ Discord-Benachrichtigung fehlgeschlagen: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		log.Printf("⚠️ Discord-Benachrichtigung: Discord antwortete mit HTTP %d", resp.StatusCode)
	}
}

func truncateForDiscord(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "…"
}
