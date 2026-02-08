package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/service"
)

const (
	sessionUserKey  = "user_id"
	sessionStateKey = "oauth_state"
)

type AuthHandler struct {
	authService  *service.AuthService
	sessionStore *sessions.CookieStore
	sessionName  string
	frontendURL  string
}

func NewAuthHandler(
	authService *service.AuthService,
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

// RegisterRoutes registriert alle Auth-Routen
func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", h.Login).Methods("GET")
	router.HandleFunc("/auth/callback", h.Callback).Methods("GET")
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