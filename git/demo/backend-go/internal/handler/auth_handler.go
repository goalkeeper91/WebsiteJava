package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
)

type AuthServiceInterface interface {
	GetAuthURL() (string, string, error)
	GetBotAuthURL() (string, string, error)
	HandleCallback(ctx context.Context, code string) (*domain.User, error)
	HandleBotCallback(ctx context.Context, code string) (*domain.User, error)
	GetUserByTwitchID(ctx context.Context, twitchID string) (*domain.User, error)
}

const (
	sessionUserKey  = "user_id"
	sessionStateKey = "oauth_state"
	sessionBotStateKey = "oauth_bot_state"
)

type AuthHandler struct {
	authService  AuthServiceInterface
	sessionStore *sessions.CookieStore
	sessionName  string
	frontendURL  string
}

func NewAuthHandler(
	authService AuthServiceInterface,
	sessionStore *sessions.CookieStore,
	sessionName string,
	frontendURL string,
) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		sessionStore: sessionStore,
		sessionName:  sessionName,
		frontendURL:  frontendURL,
	}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", h.Login).Methods("GET")
	router.HandleFunc("/auth/callback", h.Callback).Methods("GET")

	router.HandleFunc("/auth/login/bot", h.BotLogin).Methods("GET")
	router.HandleFunc("/auth/callback/bot", h.BotCallback).Methods("GET")

	router.HandleFunc("/auth/logout", h.Logout).Methods("POST")
	router.HandleFunc("/auth/me", h.Me).Methods("GET")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	authURL, state, err := h.authService.GetAuthURL()
	if err != nil {
		log.Printf("Fehler beim Generieren der Auth URL: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	session, _ := h.sessionStore.Get(r, h.sessionName)
	session.Values[sessionStateKey] = state
	if err := session.Save(r, w); err != nil {
		log.Printf("Fehler beim Speichern der Session: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		log.Printf("Fehler beim Laden der Session: %v", err)
		http.Redirect(w, r, h.frontendURL+"/error?message=session_error", http.StatusTemporaryRedirect)
		return
	}

	storedState, ok := session.Values[sessionStateKey].(string)
	if !ok || storedState == "" {
		log.Printf("Kein State in Session gefunden")
		http.Redirect(w, r, h.frontendURL+"/error?message=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	receivedState := r.URL.Query().Get("state")
	if receivedState != storedState {
		log.Printf("State Mismatch: erwartet %s, erhalten %s", storedState, receivedState)
		http.Redirect(w, r, h.frontendURL+"/error?message=state_mismatch", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("Kein Authorization Code erhalten")
		http.Redirect(w, r, h.frontendURL+"/error?message=no_code", http.StatusTemporaryRedirect)
		return
	}

	user, err := h.authService.HandleCallback(r.Context(), code)
	if err != nil {
		log.Printf("Fehler beim Verarbeiten des Callbacks: %v", err)
		http.Redirect(w, r, h.frontendURL+"/error?message=auth_failed", http.StatusTemporaryRedirect)
		return
	}

	session.Values[sessionUserKey] = user.TwitchID
	delete(session.Values, sessionStateKey)
	session.Options.MaxAge = 86400

	if err := session.Save(r, w); err != nil {
		log.Printf("Fehler beim Speichern der Session: %v", err)
		http.Redirect(w, r, h.frontendURL+"/error?message=session_save_error", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, h.frontendURL+"/dashboard", http.StatusTemporaryRedirect)
}

func (h *AuthHandler) BotLogin(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		log.Printf("Fehler beim Laden der Session: %v", err)
		http.Redirect(w, r, h.frontendURL+"/dashboard?error=not_authenticated", http.StatusTemporaryRedirect)
		return
	}

	userID, ok := session.Values[sessionUserKey].(string)
	if !ok || userID == "" {
		log.Printf("Kein User in Session - Bot Auth abgelehnt")
		http.Redirect(w, r, h.frontendURL+"/dashboard?error=not_authenticated", http.StatusTemporaryRedirect)
		return
	}

	user, err := h.authService.GetUserByTwitchID(r.Context(), userID)
	if err != nil || !user.IsAdmin {
		log.Printf("Unbefugter Zugriff auf Bot Auth von User: %s", userID)
		http.Redirect(w, r, h.frontendURL+"/dashboard?error=unauthorized", http.StatusTemporaryRedirect)
		return
	}

	authURL, state, err := h.authService.GetBotAuthURL()
	if err != nil {
		log.Printf("Fehler beim Generieren der Bot Auth URL: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	session.Values[sessionBotStateKey] = state
	if err := session.Save(r, w); err != nil {
		log.Printf("Fehler beim Speichern der Session: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Bot Auth initiiert für Admin User: %s", user.Username)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) BotCallback(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		log.Printf("Fehler beim Laden der Session: %v", err)
		http.Redirect(w, r, h.frontendURL+"/dashboard?error=session_error", http.StatusTemporaryRedirect)
		return
	}

	storedState, ok := session.Values[sessionBotStateKey].(string)
	if !ok || storedState == "" {
		log.Printf("Kein Bot State in Session gefunden")
		http.Redirect(w, r, h.frontendURL+"/dashboard?error=invalid_bot_state", http.StatusTemporaryRedirect)
		return
	}

	receivedState := r.URL.Query().Get("state")
	if receivedState != storedState {
		log.Printf("Bot State Mismatch: erwartet %s, erhalten %s", storedState, receivedState)
		http.Redirect(w, r, h.frontendURL+"/dashboard?error=bot_state_mismatch", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("Kein Bot Authorization Code erhalten")
		http.Redirect(w, r, h.frontendURL+"/dashboard?error=no_bot_code", http.StatusTemporaryRedirect)
		return
	}

	botUser, err := h.authService.HandleBotCallback(r.Context(), code)
	if err != nil {
		log.Printf("Fehler beim Verarbeiten des Bot Callbacks: %v", err)
		http.Redirect(w, r, h.frontendURL+"/dashboard?error=bot_auth_failed", http.StatusTemporaryRedirect)
		return
	}

	delete(session.Values, sessionBotStateKey)
	if err := session.Save(r, w); err != nil {
		log.Printf("Fehler beim Speichern der Session: %v", err)
	}

	log.Printf("✅ Bot Account erfolgreich authentifiziert: %s (ID: %s)", botUser.Username, botUser.TwitchID)

	http.Redirect(w, r, h.frontendURL+"/dashboard?success=bot_authenticated&bot_username="+botUser.Username, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		log.Printf("Fehler beim Laden der Session: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		log.Printf("Fehler beim Löschen der Session: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Erfolgreich ausgeloggt",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		log.Printf("Fehler beim Laden der Session: %v", err)
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values[sessionUserKey].(string)
	if !ok || userID == "" {
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUserByTwitchID(r.Context(), userID)
	if err != nil {
		log.Printf("Fehler beim Laden des Users: %v", err)
		http.Error(w, "User nicht gefunden", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}