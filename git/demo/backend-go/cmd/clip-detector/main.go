package main

import (
	"context"
	"database/sql"
	"encoding/json"
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

	// statusRedisKey/statusTTL: Snapshot des aktuellen Detector-Zustands für
	// die Admin-Kundenübersicht (Phase 5) - TTL knapp über pollInterval, damit
	// ein gestoppter/abgestürzter Detector nach kurzer Zeit als "keine Daten"
	// statt dauerhaft veraltet "aktiv" erscheint.
	statusRedisKey = "clip_detector:status"
	statusTTL      = 90 * time.Second
)

// detectorStatus ist der JSON-Snapshot, den cmd/server für die Admin-
// Kundenübersicht ausliest (siehe SubscriptionService.AdminGetStats).
type detectorStatus struct {
	ActiveChannels []detectorStatusChannel `json:"active_channels"`
	MaxConcurrent  int                     `json:"max_concurrent"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

type detectorStatusChannel struct {
	UserTwitchID string `json:"twitch_user_id"`
	Login        string `json:"login"`
}

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

	publishDetectorStatus(ctx, redisService, active, maxConcurrentIngestors)
}

// publishDetectorStatus schreibt den aktuellen Ingestor-Zustand als
// TTL-behafteten Redis-Key, damit cmd/server (separater Prozess) für die
// Admin-Kundenübersicht sehen kann, wie ausgelastet der Detector gerade ist.
func publishDetectorStatus(ctx context.Context, redisService *redis.RedisService, active map[string]*detector.StreamIngestor, maxConcurrentIngestors int) {
	channels := make([]detectorStatusChannel, 0, len(active))
	for _, ing := range active {
		channels = append(channels, detectorStatusChannel{
			UserTwitchID: ing.UserTwitchID(),
			Login:        ing.ChannelLogin(),
		})
	}

	payload, err := json.Marshal(detectorStatus{
		ActiveChannels: channels,
		MaxConcurrent:  maxConcurrentIngestors,
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		log.Printf("⚠️ Fehler beim Serialisieren des Detector-Status: %v", err)
		return
	}

	if err := redisService.GetClient().Set(ctx, statusRedisKey, payload, statusTTL).Err(); err != nil {
		log.Printf("⚠️ Fehler beim Veröffentlichen des Detector-Status in Redis: %v", err)
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
