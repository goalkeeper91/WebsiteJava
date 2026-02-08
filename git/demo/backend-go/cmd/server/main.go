package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"

	"demo/backend-go/internal/handler"
	"demo/backend-go/internal/repository/postgres"
	"demo/backend-go/internal/security"
	"demo/backend-go/internal/service"
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

	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewAuthTokenRepository(db, crypto)

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Twitch.ClientID,
		ClientSecret: cfg.Twitch.ClientSecret,
		RedirectURL:  cfg.Twitch.RedirectURL,
		Scopes:       []string{"user:read:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
			TokenURL: "https://id.twitch.tv/oauth2/token",
		},
	}

	authService := service.NewAuthService(userRepo, tokenRepo, oauthConfig, cfg.Frontend.URL)

	sessionStore := sessions.NewCookieStore([]byte(cfg.Session.Secret))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.Session.MaxAge,
		HttpOnly: cfg.Session.HTTPOnly,
		Secure:   cfg.Session.Secure,
		SameSite: parseSameSite(cfg.Session.SameSite),
	}

	authHandler := handler.NewAuthHandler(authService, sessionStore, cfg.Session.Name, cfg.Frontend.URL)

	router := mux.NewRouter()

	router.Use(corsMiddleware(cfg.Frontend.AllowedOrigins))

	router.Use(loggingMiddleware)

	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	authHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server läuft auf http://localhost:%s", cfg.Server.Port)
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

	log.Println("Server erfolgreich heruntergefahren")
}

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

	log.Println("Datenbankverbindung erfolgreich hergestellt")
	return db, nil
}

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

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

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