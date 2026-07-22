package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"demo/backend-go/internal/detector"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/repository/postgres"
	"demo/backend-go/pkg/config"
)

const pollInterval = 60 * time.Second

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Fehler beim Laden der Konfiguration: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der Datenbankverbindung: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("Fehler beim Verbinden zur Datenbank: %v", err)
	}
	log.Println("✅ Datenbankverbindung erfolgreich hergestellt")

	redisAddr := getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379")
	redisService, err := redis.NewRedisService(redisAddr, getEnv("REDIS_PASSWORD", ""), 0, noopEventHandler{})
	if err != nil {
		log.Fatalf("Fehler beim Verbinden zu Redis: %v", err)
	}
	defer redisService.Close()

	automationRepo := postgres.NewAutomationSettingsRepository(db)
	appToken := detector.NewTwitchAppTokenClient(cfg.Twitch.ClientID, cfg.Twitch.ClientSecret)

	allowlist := parseAllowlist(getEnv("DETECTOR_ALLOWED_USERS", ""))
	if len(allowlist) == 0 {
		log.Println("⚠️ DETECTOR_ALLOWED_USERS ist leer — der Detector überwacht aktuell niemanden (v1-Sicherheitsgate, siehe Phase-B-Plan).")
	} else {
		log.Printf("🔒 Detector-Allowlist aktiv für %d Twitch-User-ID(s)", len(allowlist))
	}

	active := map[string]*detector.StreamIngestor{}

	log.Println("🚀 Clip-Detector gestartet")

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		pollOnce(ctx, automationRepo, appToken, allowlist, active, redisService)
		<-ticker.C
	}
}

// pollOnce lädt die aktivierten Automation-Settings, schränkt sie auf die
// Allowlist ein, prüft den Live-Status via Twitch und startet/stoppt
// Ingestoren entsprechend.
func pollOnce(
	ctx context.Context,
	automationRepo repository.AutomationSettingsRepository,
	appToken *detector.TwitchAppTokenClient,
	allowlist map[string]bool,
	active map[string]*detector.StreamIngestor,
	redisService *redis.RedisService,
) {
	candidates, err := automationRepo.GetAllEnabled(ctx)
	if err != nil {
		log.Printf("⚠️ Fehler beim Laden der aktivierten Automation-Settings: %v", err)
		return
	}

	ids := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if allowlist[c.UserTwitchID] {
			ids = append(ids, c.UserTwitchID)
		}
	}

	live, err := appToken.GetLiveStreams(ctx, ids)
	if err != nil {
		log.Printf("⚠️ Fehler beim Live-Status-Check: %v", err)
		return
	}
	log.Printf("🔎 Live-Status geprüft: %d Kandidat(en), %d live", len(ids), len(live))

	for userTwitchID, stream := range live {
		if _, alreadyActive := active[userTwitchID]; alreadyActive {
			continue
		}
		active[userTwitchID] = detector.StartIngestor(ctx, userTwitchID, stream.UserLogin, redisService)
		log.Printf("🟢 %s ist live (%s) — Ingestion gestartet", stream.UserLogin, stream.GameName)
	}

	for userTwitchID, ing := range active {
		if _, stillLive := live[userTwitchID]; stillLive {
			continue
		}
		ing.Stop()
		delete(active, userTwitchID)
		log.Printf("🔴 Stream nicht mehr live, Ingestion gestoppt: %s", userTwitchID)
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func parseAllowlist(raw string) map[string]bool {
	result := map[string]bool{}
	for _, id := range strings.Split(raw, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			result[id] = true
		}
	}
	return result
}

// noopEventHandler ignoriert die "backend:events"-Nachrichten, die
// RedisService intern konsumiert — der Detector braucht nur Publish/GetClient,
// nicht die Aktivitäts-/Chat-/Status-Weiterverarbeitung des Web-Backends.
type noopEventHandler struct{}

func (noopEventHandler) HandleActivityEvent(event *redis.BotEvent) error { return nil }
func (noopEventHandler) HandleChatMessage(event *redis.BotEvent) error   { return nil }
func (noopEventHandler) HandleBotStatus(event *redis.BotEvent) error     { return nil }
