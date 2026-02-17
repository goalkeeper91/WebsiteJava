package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Twitch   TwitchConfig
	Security SecurityConfig
	Session  SessionConfig
	Discord  DiscordConfig
	Frontend FrontendConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

type TwitchConfig struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	AdminTwitchID string
}

type SecurityConfig struct {
	AppSecretKey string
}

type SessionConfig struct {
	Name     string
	Secret   string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
	SameSite string
}

type FrontendConfig struct {
	URL            string
	AllowedOrigins []string
}

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	BotToken     string
	RedirectURL  string
	Permissions  string
	AdminGuildID string
	AdminContactChannel string
}

// Load lädt die Konfiguration aus Umgebungsvariablen
func Load() (*Config, error) {
	// Versuche .env Datei zu laden (optional in Production)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			Name:         getEnv("DB_NAME", "twitch_auth"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "postgres"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		},
		Twitch: TwitchConfig{
			ClientID:     getEnv("TWITCH_CLIENT_ID", ""),
			ClientSecret: getEnv("TWITCH_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("TWITCH_REDIRECT_URL", ""),
			AdminTwitchID: getEnv("ADMIN_TWITCH_ID", ""),
		},
		Security: SecurityConfig{
			AppSecretKey: getEnv("GO_APP_SECRET_KEY", ""),
		},
		Session: SessionConfig{
			Name:     getEnv("SESSION_NAME", "twitch_auth_session"),
			Secret:   getEnv("SESSION_SECRET", ""),
			MaxAge:   getEnvAsInt("SESSION_MAX_AGE", 86400),
			Secure:   getEnvAsBool("SESSION_SECURE", false),
			HTTPOnly: getEnvAsBool("SESSION_HTTP_ONLY", true),
			SameSite: getEnv("SESSION_SAME_SITE", "lax"),
		},
		Frontend: FrontendConfig{
			URL:            getEnv("FRONTEND_URL", "http://localhost:3000"),
			AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", ",", []string{"http://localhost:3000"}),
		},
        Discord: DiscordConfig{
    		ClientID:            getEnv("DISCORD_CLIENT_ID", ""),
    		ClientSecret:        getEnv("DISCORD_CLIENT_SECRET", ""),
    		BotToken:            getEnv("DISCORD_BOT_TOKEN", ""),
    		RedirectURL:         getEnv("DISCORD_REDIRECT_URL", "http://localhost:3000/discord/callback"),
    		Permissions:         getEnv("DISCORD_PERMISSIONS", "8"),
    		AdminGuildID:        getEnv("DISCORD_ADMIN_GUILD_ID", ""),
    		AdminContactChannel: getEnv("DISCORD_ADMIN_CONTACT_CHANNEL", ""),
    	},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate überprüft die Konfiguration auf Vollständigkeit
func (c *Config) Validate() error {
	if c.Twitch.ClientID == "" {
		return fmt.Errorf("TWITCH_CLIENT_ID ist erforderlich")
	}
	if c.Twitch.ClientSecret == "" {
		return fmt.Errorf("TWITCH_CLIENT_SECRET ist erforderlich")
	}
	if c.Twitch.RedirectURL == "" {
		return fmt.Errorf("TWITCH_REDIRECT_URL ist erforderlich")
	}
	if len(c.Security.AppSecretKey) != 32 {
		return fmt.Errorf("GO_APP_SECRET_KEY muss genau 32 Zeichen lang sein (aktuell: %d)", len(c.Security.AppSecretKey))
	}
	if len(c.Session.Secret) < 32 {
		return fmt.Errorf("SESSION_SECRET muss mindestens 32 Zeichen lang sein")
	}
	return nil
}

// GetDSN gibt den PostgreSQL Connection String zurück
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// Helper-Funktionen
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

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsSlice(key, separator string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, separator)
	}
	return defaultValue
}