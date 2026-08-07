package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
)

type AuthServiceInterface interface {
	GetAuthURL(redirectURI string) (string, string, error)
	GetBotAuthURL() (string, string, error)
	HandleCallback(ctx context.Context, code, redirectURI string) (*domain.User, error)
	HandleBotCallback(ctx context.Context, code string) (*domain.User, error)
	GetUserByTwitchID(ctx context.Context, twitchID string) (*domain.User, error)
}

const (
	sessionUserKey       = "user_id"
	sessionStateKey      = "oauth_state"
	sessionBotStateKey   = "oauth_bot_state"
	sessionReturnToKey   = "oauth_return_to"
	sessionOriginKey     = "oauth_origin"
	defaultLoginRedirect = "/dashboard"
)

// botStorefrontHost mirrors frontend/src/lib/botDomain.ts's
// BOT_STOREFRONT_HOSTNAME - the Paddle-compliant Twitch Bot storefront
// subdomain. Login/Callback keep the entire OAuth round-trip (authorize ->
// Twitch -> callback -> final redirect) on whichever of these hosts the
// user actually started from, otherwise a login started on the bot
// storefront would silently bounce back to the main domain instead of
// staying there.
const botStorefrontHost = "bot.goalkeeper91.de"
const botStorefrontOrigin = "https://" + botStorefrontHost

// oauthOriginForRequest returns "" for the main domain (meaning: use the
// handler's/service's configured defaults, unchanged), or the bot
// storefront's own origin when the request came in on that hostname.
func oauthOriginForRequest(r *http.Request) string {
	if r.Host == botStorefrontHost {
		return botStorefrontOrigin
	}
	return ""
}

// isValidReturnTo guards against open-redirect: only a same-origin relative
// path is accepted (must start with exactly one '/', never '//' - which
// browsers treat as protocol-relative - and never contain '://').
func isValidReturnTo(path string) bool {
	if path == "" || path[0] != '/' {
		return false
	}
	if strings.HasPrefix(path, "//") {
		return false
	}
	if strings.Contains(path, "://") {
		return false
	}
	return true
}

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
	origin := oauthOriginForRequest(r)
	redirectURI := ""
	if origin != "" {
		redirectURI = origin + "/auth/callback"
	}

	authURL, state, err := h.authService.GetAuthURL(redirectURI)
	if err != nil {
		log.Printf("Fehler beim Generieren der Auth URL: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	session, _ := h.sessionStore.Get(r, h.sessionName)
	session.Values[sessionStateKey] = state
	session.Values[sessionOriginKey] = origin
	if returnTo := r.URL.Query().Get("returnTo"); isValidReturnTo(returnTo) {
		session.Values[sessionReturnToKey] = returnTo
	}
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

	// origin/redirectURI must mirror exactly what Login stored for this
	// flow - Twitch's token exchange rejects a redirect_uri mismatch, and a
	// user who started on the bot storefront must land back there, not on
	// the main domain.
	origin := h.frontendURL
	redirectURI := ""
	if stored, ok := session.Values[sessionOriginKey].(string); ok && stored != "" {
		origin = stored
		redirectURI = stored + "/auth/callback"
	}

	storedState, ok := session.Values[sessionStateKey].(string)
	if !ok || storedState == "" {
		log.Printf("Kein State in Session gefunden")
		http.Redirect(w, r, origin+"/error?message=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	receivedState := r.URL.Query().Get("state")
	if receivedState != storedState {
		log.Printf("State Mismatch: erwartet %s, erhalten %s", storedState, receivedState)
		http.Redirect(w, r, origin+"/error?message=state_mismatch", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("Kein Authorization Code erhalten")
		http.Redirect(w, r, origin+"/error?message=no_code", http.StatusTemporaryRedirect)
		return
	}

	user, err := h.authService.HandleCallback(r.Context(), code, redirectURI)
	if err != nil {
		log.Printf("Fehler beim Verarbeiten des Callbacks: %v", err)
		http.Redirect(w, r, origin+"/error?message=auth_failed", http.StatusTemporaryRedirect)
		return
	}

	session.Values[sessionUserKey] = user.TwitchID
	session.Values["is_admin"] = user.IsAdmin
	delete(session.Values, sessionStateKey)
	delete(session.Values, sessionOriginKey)

	returnTo := defaultLoginRedirect
	if stored, ok := session.Values[sessionReturnToKey].(string); ok && isValidReturnTo(stored) {
		returnTo = stored
	}
	delete(session.Values, sessionReturnToKey)

	session.Options.MaxAge = 86400

	if err := session.Save(r, w); err != nil {
		log.Printf("Fehler beim Speichern der Session: %v", err)
		http.Redirect(w, r, origin+"/error?message=session_save_error", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("✅ User %s erfolgreich eingeloggt (Admin: %v)", user.Username, user.IsAdmin)
	http.Redirect(w, r, origin+returnTo, http.StatusTemporaryRedirect)
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

	// Get is_admin from session (or fallback to user.IsAdmin)
	isAdmin, _ := session.Values["is_admin"].(bool)

	response := map[string]interface{}{
		"id":         user.TwitchID, // ← FIXED: use TwitchID instead of ID
		"twitch_id":  user.TwitchID,
		"username":   user.Username,
		"email":      user.Email,
		"is_admin":   isAdmin,
		"created_at": user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}