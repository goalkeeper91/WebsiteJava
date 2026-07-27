package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/handler"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository/postgres"
	"demo/backend-go/internal/security"
	"demo/backend-go/internal/service"
	"demo/backend-go/internal/twitch"
	"demo/backend-go/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Fehler beim Laden der Konfiguration: %v", err)
	}

	log.Printf("Server startet in %s Modus auf Port %s", cfg.Server.Environment, cfg.Server.Port)

	db, err := initDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Fehler beim Initialisieren der Datenbank: %v", err)
	}
	defer db.Close()

	crypto, err := security.NewCrypto(cfg.Security.AppSecretKey)
	if err != nil {
		log.Fatalf("Fehler beim Initialisieren der Verschlüsselung: %v", err)
	}

	// ============================================================
	// REPOSITORIES
	// ============================================================

	// Twitch
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewAuthTokenRepository(db, crypto)
	channelRepo := postgres.NewTwitchChannelRepository(db)
	commandRepo := postgres.NewChatCommandRepository(db)
	activityRepo := postgres.NewStreamActivityRepository(db)
	contactRepo := postgres.NewContactRequestRepository(db)

	// Discord
	discordConnRepo := postgres.NewDiscordConnectionRepository(db, crypto)
	discordGuildRepo := postgres.NewDiscordGuildRepository(db)
	discordSettingsRepo := postgres.NewDiscordGuildSettingsRepository(db)
	jtcRepo := postgres.NewJoinToCreateRepository(db)

	// Subathon Timer
	subathonRepo := postgres.NewSubathonRepository(db)

	// SaaS / n8n
	subscriptionTierRepo := postgres.NewSubscriptionTierRepository(db)
	userSubscriptionRepo := postgres.NewUserSubscriptionRepository(db)
	n8nIntegrationRepo := postgres.NewN8NIntegrationRepository(db)
	workflowTemplateRepo := postgres.NewWorkflowTemplateRepository(db)
	voteRepo := postgres.NewVoteRepository(db)
	analyticsRepo := postgres.NewUsageAnalyticsRepository(db)

	// Clip-Automatisierung
	automationSettingsRepo := postgres.NewAutomationSettingsRepository(db)
	clipLogRepo := postgres.NewClipLogRepository(db)

	// Automatisierte Chat-Nachrichten (Timer/Scheduler)
	scheduledMessageRepo := postgres.NewScheduledMessageRepository(db)

	// Automod
	automodRepo := postgres.NewAutomodRepository(db)

	// Loyalty/Watchtime
	loyaltyRepo := postgres.NewLoyaltyRepository(db)

	// Giveaways
	giveawayRepo := postgres.NewGiveawayRepository(db)

	// Team-Zugriff (Twitch-Mods bekommen delegierten Dashboard-Zugriff)
	teamMemberRepo := postgres.NewTeamMemberRepository(db)

	// ============================================================
	// REDIS
	// ============================================================

	activityService := service.NewActivityService(activityRepo)
	activityEventHandler := service.NewActivityEventHandler(activityService)

	var redisService *redis.RedisService
	redisAddr := fmt.Sprintf("%s:%s",
		getEnv("REDIS_HOST", "localhost"),
		getEnv("REDIS_PORT", "6379"),
	)
	redisService, err = redis.NewRedisService(
		redisAddr,
		getEnv("REDIS_PASSWORD", ""),
		getEnvAsInt("REDIS_DB", 0),
		activityEventHandler,
	)
	if err != nil {
		log.Printf("⚠️ Warnung: Redis-Verbindung fehlgeschlagen: %v", err)
		log.Printf("⚠️ Bot-Kommunikation wird nicht verfügbar sein")
		redisService = nil
	} else {
		defer redisService.Close()
		log.Println("✅ Redis Service erfolgreich verbunden")
	}

	// ============================================================
	// SERVICES
	// ============================================================

	// Twitch core
	channelService := service.NewTwitchChannelService(channelRepo, redisService)

	// SaaS layer (subscription first – andere Services hängen davon ab)
	subscriptionService := service.NewSubscriptionService(
		userSubscriptionRepo,
		subscriptionTierRepo,
		analyticsRepo,
		userRepo,
	)

	n8nService := service.NewN8NService(
		n8nIntegrationRepo,
		subscriptionService,
	)

	automationSettingsService := service.NewAutomationSettingsService(
		automationSettingsRepo,
		subscriptionService,
	)
	clipLogService := service.NewClipLogService(clipLogRepo)

	scheduledMessageAppToken := twitch.NewTwitchAppTokenClient(cfg.Twitch.ClientID, cfg.Twitch.ClientSecret)
	scheduledMessageService := service.NewScheduledMessageService(
		scheduledMessageRepo,
		channelRepo,
		commandRepo,
		redisService,
		scheduledMessageAppToken,
	)

	automodService := service.NewAutomodService(automodRepo, redisService, loyaltyRepo)

	loyaltyChattersClient := twitch.NewChattersClient(cfg.Twitch.ClientID)
	loyaltyService := service.NewLoyaltyService(loyaltyRepo, tokenRepo, scheduledMessageAppToken, loyaltyChattersClient)

	giveawayService := service.NewGiveawayService(giveawayRepo, redisService)

	teamService := service.NewTeamService(teamMemberRepo, userRepo, scheduledMessageAppToken)

	commandService := service.NewChatCommandService(
		commandRepo,
		channelRepo,
		redisService,
		subscriptionService,
		n8nService,
	)

	voteService := service.NewVoteService(
		voteRepo,
		n8nService,
		subscriptionService,
		analyticsRepo,
	)

	workflowTemplateService := service.NewWorkflowTemplateService(
		workflowTemplateRepo,
		subscriptionService,
	)

	analyticsService := service.NewAnalyticsService(
		analyticsRepo,
		subscriptionService,
	)
	_ = analyticsService // Available for future use / injection

	// Token refresh background service
	tokenRefreshService := service.NewTokenRefreshService(
		userRepo,
		tokenRepo,
		nil, // authService wird unten gesetzt
		1*time.Hour,
	)

	// Discord
	discordClientID := getEnv("DISCORD_CLIENT_ID", "")
	discordClientSecret := getEnv("DISCORD_CLIENT_SECRET", "")
	discordRedirectURL := getEnv("DISCORD_REDIRECT_URL", "http://localhost:3000/discord/callback")
	discordPermissions := getEnv("DISCORD_PERMISSIONS", "8")

	if discordClientID == "" {
		log.Printf("⚠️ Warnung: DISCORD_CLIENT_ID nicht gesetzt - Discord Features deaktiviert")
	}

	// Subathon Timer - EventSub webhook transport
	subathonWebhookSecret := getEnv("SUBATHON_EVENTSUB_SECRET", "")
	subathonCallbackURL := getEnv("SUBATHON_CALLBACK_URL", "https://goalkeeper91.de/api/subathon/eventsub/callback")
	if subathonWebhookSecret == "" {
		log.Printf("⚠️ Warnung: SUBATHON_EVENTSUB_SECRET nicht gesetzt - Subathon-Webhook-Signaturen koennen nicht verifiziert werden")
	}

	// Live-Dashboard - Activity-Feed EventSub webhook transport (follow/sub/gift-sub/cheer/raid)
	activityWebhookSecret := getEnv("ACTIVITY_EVENTSUB_SECRET", "")
	activityCallbackURL := getEnv("ACTIVITY_CALLBACK_URL", "https://goalkeeper91.de/api/activity/eventsub/callback")
	if activityWebhookSecret == "" {
		log.Printf("⚠️ Warnung: ACTIVITY_EVENTSUB_SECRET nicht gesetzt - Activity-Webhook-Signaturen koennen nicht verifiziert werden")
	}

	discordAuthService := service.NewDiscordAuthService(
		discordConnRepo,
		discordClientID,
		discordClientSecret,
		discordRedirectURL,
	)

	discordGuildService := service.NewDiscordGuildService(
		discordGuildRepo,
		discordSettingsRepo,
		redisService,
	)

	discordSettingsService := service.NewDiscordSettingsService(
		discordSettingsRepo,
		discordGuildRepo,
	)

	discordNotificationService := service.NewDiscordNotificationService(
		discordSettingsRepo,
		discordGuildRepo,
		redisService,
	)

	jtcService := service.NewJoinToCreateService(
		jtcRepo,
		discordGuildRepo,
		discordSettingsRepo,
		redisService,
	)

	subathonService := service.NewSubathonService(subathonRepo)

	// ============================================================
	// AUTH SERVICE (braucht channelService → kommt nach services)
	// ============================================================

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Twitch.ClientID,
		ClientSecret: cfg.Twitch.ClientSecret,
		RedirectURL:  cfg.Twitch.RedirectURL,
		Scopes: []string{
			"user:read:email", "clips:edit",
			// Subathon-Timer: EventSub-Subscriptions für Subs/Bits
			"channel:read:subscriptions", "bits:read",
			"channel:read:redemptions", "moderator:read:followers",
			// Eingebaute Chat-Commands: !title/!game setzen
			"channel:manage:broadcast",
			// Automod: Nachrichten löschen + Timeout verhängen
			"moderator:manage:chat_messages", "moderator:manage:banned_users",
			// Loyalty/Watchtime: Chatter-Liste für Punkte-Vergabe
			"moderator:read:chatters",
			// Live-Dashboard: Werbepause auslösen
			"channel:edit:commercial",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
			TokenURL: "https://id.twitch.tv/oauth2/token",
		},
	}

	baseAuthService := service.NewAuthService(
		userRepo,
		tokenRepo,
		oauthConfig,
		cfg.Frontend.URL,
		cfg.Twitch.AdminTwitchID,
	)

	// Inject authService into tokenRefreshService
	tokenRefreshService.SetAuthService(baseAuthService)

	// Live-Dashboard: Chat, Live-Stats, Titel/Kategorie, Werbepause
	channelClient := twitch.NewChannelClient(cfg.Twitch.ClientID)
	streamDashboardService := service.NewStreamDashboardService(baseAuthService, channelClient, scheduledMessageAppToken, activityRepo)

	// Clip-Erzeugung (Phase C): braucht baseAuthService für Token-Refresh
	// Phase D: Caption via selbstgehostetem Ollama, Discord-Post über den bestehenden Bot
	captionService := service.NewCaptionService(
		getEnv("OLLAMA_BASE_URL", "http://ollama:11434"),
		getEnv("OLLAMA_MODEL", "llama3.2:3b"),
	)

	clipCreationService := service.NewClipCreationService(
		clipLogRepo,
		tokenRepo,
		automationSettingsRepo,
		baseAuthService,
		n8nService,
		discordNotificationService,
		captionService,
		cfg.Twitch.ClientID,
	)

	// Wrapper: synct Channel nach Login
	authService := &AuthServiceWithChannelSync{
		authService:    baseAuthService,
		channelService: channelService,
	}

	// ============================================================
	// SESSION STORE
	// ============================================================

	sessionStore := sessions.NewCookieStore([]byte(cfg.Session.Secret))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.Session.MaxAge,
		HttpOnly: cfg.Session.HTTPOnly,
		Secure:   cfg.Session.Secure,
		SameSite: parseSameSite(cfg.Session.SameSite),
	}

	// ============================================================
	// HANDLERS
	// ============================================================

	// Twitch
	authHandler := handler.NewAuthHandler(authService, sessionStore, cfg.Session.Name, cfg.Frontend.URL)
	commandHandler := handler.NewChatCommandHandler(commandService, sessionStore, cfg.Session.Name, teamService)
	activityHandler := handler.NewActivityHandler(activityService, sessionStore, cfg.Session.Name, teamService)
	streamDashboardHandler := handler.NewStreamDashboardHandler(streamDashboardService, sessionStore, cfg.Session.Name, teamService)
	contactHandler := handler.NewContactHandler(contactRepo)
	botStatusHandler := handler.NewBotStatusHandler(userRepo, tokenRepo, sessionStore, redisService, cfg.Session.Name)
	botStatsHandler := handler.NewBotStatsHandler(redisService, sessionStore, cfg.Session.Name)

	// SaaS / n8n
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, sessionStore, cfg.Session.Name)
	n8nHandler := handler.NewN8NHandler(n8nService, sessionStore, cfg.Session.Name)
	voteHandler := handler.NewVoteHandler(voteService, sessionStore, cfg.Session.Name, cfg.Security.BotInternalSecret)
	workflowTemplateHandler := handler.NewWorkflowTemplateHandler(workflowTemplateService, sessionStore, cfg.Session.Name)
	botAnnounceHandler := handler.NewBotAnnounceHandler(redisService, cfg.Security.BotInternalSecret)
	automationSettingsHandler := handler.NewAutomationSettingsHandler(automationSettingsService, sessionStore, cfg.Session.Name)
	clipLogHandler := handler.NewClipLogHandler(clipLogService, sessionStore, cfg.Session.Name)
	scheduledMessageHandler := handler.NewScheduledMessageHandler(scheduledMessageService, sessionStore, cfg.Session.Name, teamService)
	automodHandler := handler.NewAutomodHandler(automodService, sessionStore, cfg.Session.Name, cfg.Security.BotInternalSecret, teamService)
	loyaltyHandler := handler.NewLoyaltyHandler(loyaltyService, sessionStore, cfg.Session.Name, cfg.Security.BotInternalSecret, teamService)
	giveawayHandler := handler.NewGiveawayHandler(giveawayService, sessionStore, cfg.Session.Name, cfg.Security.BotInternalSecret, redisService, teamService)
	teamHandler := handler.NewTeamHandler(teamService, sessionStore, cfg.Session.Name)

	// Discord
	discordAuthHandler := handler.NewDiscordAuthHandler(discordAuthService, discordGuildService, sessionStore, cfg.Session.Name, redisService)
	discordGuildHandler := handler.NewDiscordGuildHandler(discordGuildService, userRepo, sessionStore, cfg.Session.Name, discordClientID, discordPermissions)
	discordSettingsHandler := handler.NewDiscordSettingsHandler(discordSettingsService, discordNotificationService, sessionStore, cfg.Session.Name)
	jtcHandler := handler.NewJoinToCreateHandler(jtcService, sessionStore, cfg.Session.Name)
	discordBotStatusHandler := handler.NewDiscordBotStatusHandler(discordGuildService, redisService)

	// Subathon Timer
	subathonHandler := handler.NewSubathonHandler(subathonService, sessionStore, cfg.Session.Name)
	subathonWebhookHandler := handler.NewSubathonWebhookHandler(subathonService, subathonWebhookSecret)

	// Live-Dashboard: Activity-Feed EventSub webhook
	activityWebhookHandler := handler.NewActivityWebhookHandler(activityService, activityWebhookSecret)

	// ============================================================
	// ROUTER
	// ============================================================

	router := mux.NewRouter()
	router.Use(corsMiddleware(cfg.Frontend.AllowedOrigins))
	router.Use(loggingMiddleware)

	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Twitch routes
	authHandler.RegisterRoutes(router)
	commandHandler.RegisterRoutes(router)
	activityHandler.RegisterRoutes(router)
	contactHandler.RegisterRoutes(router)
	botStatusHandler.RegisterRoutes(router)
	botStatsHandler.RegisterRoutes(router)

	// SaaS / n8n routes
	subscriptionHandler.RegisterRoutes(router)
	n8nHandler.RegisterRoutes(router)
	voteHandler.RegisterRoutes(router)
	workflowTemplateHandler.RegisterRoutes(router)
	botAnnounceHandler.RegisterRoutes(router)
	automationSettingsHandler.RegisterRoutes(router)
	clipLogHandler.RegisterRoutes(router)
	scheduledMessageHandler.RegisterRoutes(router)
	automodHandler.RegisterRoutes(router)
	loyaltyHandler.RegisterRoutes(router)
	giveawayHandler.RegisterRoutes(router)
	teamHandler.RegisterRoutes(router)
	streamDashboardHandler.RegisterRoutes(router)

	// Discord routes
	discordAuthHandler.RegisterRoutes(router)
	discordGuildHandler.RegisterRoutes(router)
	discordSettingsHandler.RegisterRoutes(router)
	jtcHandler.RegisterRoutes(router)
	discordBotStatusHandler.RegisterRoutes(router)

	// Subathon Timer routes
	subathonHandler.RegisterRoutes(router)
	subathonWebhookHandler.RegisterRoutes(router)

	// Live-Dashboard Activity-Feed webhook route
	activityWebhookHandler.RegisterRoutes(router)

	// ============================================================
	// BACKGROUND SERVICES
	// ============================================================

	// Discord bot event listener
	discordBotEventListener := service.NewDiscordBotEventListener(redisService, discordGuildRepo, discordConnRepo)
	go func() {
		ctx := context.Background()
		log.Println("🚀 Discord Bot Event Listener gestartet...")
		if err := discordBotEventListener.Start(ctx); err != nil {
			log.Printf("Discord Bot Event Listener Fehler: %v", err)
		}
	}()

	// Subathon EventSub manager - keeps webhook subscriptions active per user
	// and retries any events that failed to persist (see subathon_eventsub_manager.go)
	subathonEventSubManager := service.NewSubathonEventSubManager(subathonService, tokenRepo, cfg.Twitch.ClientID, subathonWebhookSecret, subathonCallbackURL)
	go subathonEventSubManager.Start(context.Background())

	// Live-Dashboard Activity EventSub manager - keeps follow/sub/gift-sub/
	// cheer/raid webhook subscriptions active per user (see activity_eventsub_manager.go)
	activityEventSubManager := service.NewActivityEventSubManager(userRepo, tokenRepo, cfg.Twitch.ClientID, activityWebhookSecret, activityCallbackURL)
	go activityEventSubManager.Start(context.Background())

	// Clip trigger listener (Phase C)
	clipTriggerListener := service.NewClipTriggerListener(redisService, clipCreationService)
	go func() {
		ctx := context.Background()
		if err := clipTriggerListener.Start(ctx); err != nil {
			log.Printf("Clip Trigger Listener Fehler: %v", err)
		}
	}()

	// Token refresh service
	go tokenRefreshService.Start()
	defer tokenRefreshService.Stop()

	// Scheduled-message scheduler - checks every 30s for due automated chat messages
	scheduledMessageScheduler := service.NewScheduledMessageScheduler(scheduledMessageService, 30*time.Second)
	go scheduledMessageScheduler.Start()
	defer scheduledMessageScheduler.Stop()

	// Loyalty scheduler - checks every 60s for channels due for a points accrual tick
	loyaltyScheduler := service.NewLoyaltyScheduler(loyaltyService, 60*time.Second)
	go loyaltyScheduler.Start()
	defer loyaltyScheduler.Stop()

	log.Println("✅ Background Services gestartet")

	// ============================================================
	// HTTP SERVER
	// ============================================================

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server läuft auf http://localhost:%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Fehler: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server wird heruntergefahren...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Fehler: %v", err)
	}

	log.Println("✅ Server erfolgreich heruntergefahren")
}

// ============================================================
// AUTH WRAPPER (Channel Sync nach Login)
// ============================================================

type AuthServiceWithChannelSync struct {
	authService    *service.AuthService
	channelService *service.TwitchChannelService
}

func (s *AuthServiceWithChannelSync) GetAuthURL() (string, string, error) {
	return s.authService.GetAuthURL()
}

func (s *AuthServiceWithChannelSync) HandleCallback(ctx context.Context, code string) (*domain.User, error) {
	user, err := s.authService.HandleCallback(ctx, code)
	if err != nil {
		return nil, err
	}

	if err := s.channelService.SyncChannel(ctx, user.TwitchID, user.Username); err != nil {
		log.Printf("⚠️ Warnung: Channel Sync fehlgeschlagen: %v", err)
	}

	return user, nil
}

func (s *AuthServiceWithChannelSync) HandleBotCallback(ctx context.Context, code string) (*domain.User, error) {
	return s.authService.HandleBotCallback(ctx, code)
}

func (s *AuthServiceWithChannelSync) GetUserByTwitchID(ctx context.Context, twitchID string) (*domain.User, error) {
	return s.authService.GetUserByTwitchID(ctx, twitchID)
}

func (s *AuthServiceWithChannelSync) GetBotAuthURL() (string, string, error) {
	return s.authService.GetBotAuthURL()
}

// ============================================================
// DATABASE
// ============================================================

func initDatabase(cfg *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("fehler beim Öffnen der Datenbankverbindung: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("fehler beim Verbinden zur Datenbank: %w", err)
	}

	log.Println("✅ Datenbankverbindung erfolgreich hergestellt")
	return db, nil
}

// ============================================================
// MIDDLEWARE
// ============================================================

func corsMiddleware(allowedOrigins []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

// ============================================================
// HELPERS
// ============================================================

func parseSameSite(sameSite string) http.SameSite {
	switch sameSite {
	case "strict":
		return http.SameSiteStrictMode
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}