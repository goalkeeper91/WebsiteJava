package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"demo/backend-go/internal/detector"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/repository/postgres"
	"demo/backend-go/internal/twitch"
	"demo/backend-go/pkg/config"
)

const (
	pollInterval                  = 60 * time.Second
	defaultMaxConcurrentIngestors = 20
)

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
	appToken := twitch.NewTwitchAppTokenClient(cfg.Twitch.ClientID, cfg.Twitch.ClientSecret)

	maxConcurrentIngestors := defaultMaxConcurrentIngestors
	if raw := getEnv("MAX_CONCURRENT_INGESTORS", ""); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			maxConcurrentIngestors = parsed
		} else {
			log.Printf("⚠️ Ungültiger MAX_CONCURRENT_INGESTORS-Wert %q, verwende Default %d", raw, defaultMaxConcurrentIngestors)
		}
	}
	log.Printf("📊 Kapazitätsgrenze: max. %d gleichzeitige Ingestoren", maxConcurrentIngestors)

	active := map[string]*detector.StreamIngestor{}

	log.Println("🚀 Clip-Detector gestartet")

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		pollOnce(ctx, automationRepo, appToken, maxConcurrentIngestors, active, redisService)
		<-ticker.C
	}
}

// pollOnce lädt die aktivierten Automation-Settings mitsamt Abo-Daten, prüft
// den Live-Status via Twitch und startet/stoppt Ingestoren entsprechend —
// begrenzt durch maxConcurrentIngestors.
func pollOnce(
	ctx context.Context,
	automationRepo repository.AutomationSettingsRepository,
	appToken *twitch.TwitchAppTokenClient,
	maxConcurrentIngestors int,
	active map[string]*detector.StreamIngestor,
	redisService *redis.RedisService,
) {
	candidates, err := automationRepo.GetAllEnabledWithSubscription(ctx)
	if err != nil {
		log.Printf("⚠️ Fehler beim Laden der aktivierten Automation-Settings: %v", err)
		return
	}

	ids := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if !c.IsAdmin && !(c.Subscription.IsActive() && c.Subscription.Tier.HasFeature("clip_automation")) {
			continue
		}
		ids = append(ids, c.Settings.UserTwitchID)
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
		if len(active) >= maxConcurrentIngestors {
			log.Printf("⚠️ Kapazitätsgrenze erreicht (%d/%d) — %s wird in diesem Tick übersprungen", len(active), maxConcurrentIngestors, stream.UserLogin)
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

// noopEventHandler ignoriert die "backend:events"-Nachrichten, die
// RedisService intern konsumiert — der Detector braucht nur Publish/GetClient,
// nicht die Aktivitäts-/Chat-/Status-Weiterverarbeitung des Web-Backends.
type noopEventHandler struct{}

func (noopEventHandler) HandleActivityEvent(event *redis.BotEvent) error { return nil }
func (noopEventHandler) HandleChatMessage(event *redis.BotEvent) error   { return nil }
func (noopEventHandler) HandleBotStatus(event *redis.BotEvent) error     { return nil }
