package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"demo/backend-go/internal/composio"
)

// Public "connect your TikTok account" demo page (/tiktok on the frontend):
// Composio holds the OAuth app + tokens, this handler is the thin server-side
// half so the Composio API key never reaches the browser. Nothing about the
// visitor is stored in our own database - only a random, HttpOnly cookie id
// that namespaces their Composio "user".
//
//	POST /api/tiktok/connect     -> {"redirectUrl": "<Composio connect link>"}
//	GET  /api/tiktok/me          -> {"connected": bool, "displayName", "avatarUrl"}
//	POST /api/tiktok/disconnect  -> {"disconnected": bool}

const (
	tiktokToolkitSlug  = "tiktok"
	tiktokProfileTool  = "TIKTOK_GET_USER_STATS"
	tiktokCookieName   = "tiktok_demo_uid"
	tiktokCookieMaxAge = 24 * 60 * 60 // seconds - just long enough to finish one connect flow
	tiktokUserPrefix   = "tiktok-demo-"
)

var tiktokUIDPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)

// tiktokComposio is the slice of *composio.Client this handler uses (an
// interface so tests can fake it).
type tiktokComposio interface {
	Configured() bool
	CreateConnectLink(ctx context.Context, authConfigID, userID, callbackURL string) (*composio.ConnectLink, error)
	FindActiveAccount(ctx context.Context, userID, toolkitSlug string) (*composio.Account, error)
	ExecuteTool(ctx context.Context, toolSlug, userID, connectedAccountID string, arguments map[string]any) (json.RawMessage, error)
	DeleteAccount(ctx context.Context, connectedAccountID string, revoke bool) error
}

type TikTokHandler struct {
	client       tiktokComposio
	authConfigID string
	callbackURL  string
	secureCookie bool

	connectLimiter    *rateLimiter
	meLimiter         *rateLimiter
	disconnectLimiter *rateLimiter
}

// NewTikTokHandler: frontendURL is FRONTEND_URL (e.g. https://goalkeeper91.de),
// the page Composio sends the visitor back to after TikTok consent.
func NewTikTokHandler(client tiktokComposio, authConfigID, frontendURL string, secureCookie bool) *TikTokHandler {
	return &TikTokHandler{
		client:            client,
		authConfigID:      authConfigID,
		callbackURL:       strings.TrimRight(frontendURL, "/") + "/tiktok",
		secureCookie:      secureCookie,
		connectLimiter:    newRateLimiter(10, time.Hour),
		meLimiter:         newRateLimiter(60, time.Minute),
		disconnectLimiter: newRateLimiter(10, time.Hour),
	}
}

func (h *TikTokHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/tiktok/connect", h.Connect).Methods("POST")
	router.HandleFunc("/api/tiktok/me", h.Me).Methods("GET")
	router.HandleFunc("/api/tiktok/disconnect", h.Disconnect).Methods("POST")
}

func (h *TikTokHandler) ready(w http.ResponseWriter) bool {
	if h.client == nil || !h.client.Configured() || h.authConfigID == "" {
		writeTikTokJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "TikTok-Anbindung ist nicht konfiguriert."})
		return false
	}
	return true
}

func (h *TikTokHandler) Connect(w http.ResponseWriter, r *http.Request) {
	if !h.ready(w) {
		return
	}
	if !h.connectLimiter.allow(clientIP(r)) {
		writeTikTokJSON(w, http.StatusTooManyRequests, map[string]any{"error": "Zu viele Anfragen. Bitte später erneut versuchen."})
		return
	}

	uid, err := h.ensureVisitorID(w, r)
	if err != nil {
		log.Printf("tiktok: visitor id generation failed: %v", err)
		writeTikTokJSON(w, http.StatusInternalServerError, map[string]any{"error": "Interner Serverfehler"})
		return
	}

	link, err := h.client.CreateConnectLink(r.Context(), h.authConfigID, tiktokUserPrefix+uid, h.callbackURL)
	if err != nil {
		log.Printf("tiktok: create connect link failed: %v", err)
		writeTikTokJSON(w, http.StatusBadGateway, map[string]any{"error": "Verbindung konnte nicht gestartet werden."})
		return
	}
	// Never hand the browser anything but an https URL to navigate to.
	if !strings.HasPrefix(link.RedirectURL, "https://") {
		log.Printf("tiktok: refusing non-https redirect url from composio")
		writeTikTokJSON(w, http.StatusBadGateway, map[string]any{"error": "Verbindung konnte nicht gestartet werden."})
		return
	}

	writeTikTokJSON(w, http.StatusOK, map[string]any{"redirectUrl": link.RedirectURL})
}

func (h *TikTokHandler) Me(w http.ResponseWriter, r *http.Request) {
	if !h.ready(w) {
		return
	}
	if !h.meLimiter.allow(clientIP(r)) {
		writeTikTokJSON(w, http.StatusTooManyRequests, map[string]any{"error": "Zu viele Anfragen. Bitte später erneut versuchen."})
		return
	}

	uid, ok := visitorIDFromCookie(r)
	if !ok {
		writeTikTokJSON(w, http.StatusOK, map[string]any{"connected": false})
		return
	}
	userID := tiktokUserPrefix + uid

	account, err := h.client.FindActiveAccount(r.Context(), userID, tiktokToolkitSlug)
	if err != nil {
		log.Printf("tiktok: account lookup failed: %v", err)
		writeTikTokJSON(w, http.StatusBadGateway, map[string]any{"error": "Kontostatus konnte nicht geladen werden."})
		return
	}
	if account == nil {
		writeTikTokJSON(w, http.StatusOK, map[string]any{"connected": false})
		return
	}

	data, err := h.client.ExecuteTool(r.Context(), tiktokProfileTool, userID, account.ID, map[string]any{
		// user.info.basic only - the stats/profile fields need extra scopes.
		"fields": []string{"display_name", "avatar_url"},
	})
	if err != nil {
		log.Printf("tiktok: profile fetch failed: %v", err)
		writeTikTokJSON(w, http.StatusBadGateway, map[string]any{"connected": true, "error": "Profil konnte nicht geladen werden."})
		return
	}

	displayName, avatarURL := extractTikTokProfile(data)
	// Only ever pass an https image URL on to the page.
	if !strings.HasPrefix(avatarURL, "https://") {
		avatarURL = ""
	}
	writeTikTokJSON(w, http.StatusOK, map[string]any{
		"connected":   true,
		"displayName": displayName,
		"avatarUrl":   avatarURL,
	})
}

func (h *TikTokHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	if !h.ready(w) {
		return
	}
	if !h.disconnectLimiter.allow(clientIP(r)) {
		writeTikTokJSON(w, http.StatusTooManyRequests, map[string]any{"error": "Zu viele Anfragen. Bitte später erneut versuchen."})
		return
	}

	uid, ok := visitorIDFromCookie(r)
	if !ok {
		writeTikTokJSON(w, http.StatusOK, map[string]any{"disconnected": false})
		return
	}
	userID := tiktokUserPrefix + uid

	account, err := h.client.FindActiveAccount(r.Context(), userID, tiktokToolkitSlug)
	if err != nil {
		log.Printf("tiktok: account lookup failed: %v", err)
		writeTikTokJSON(w, http.StatusBadGateway, map[string]any{"error": "Kontostatus konnte nicht geladen werden."})
		return
	}
	if account == nil {
		writeTikTokJSON(w, http.StatusOK, map[string]any{"disconnected": false})
		return
	}

	// revoke=true: also revoke the token on TikTok's side, not just forget it.
	if err := h.client.DeleteAccount(r.Context(), account.ID, true); err != nil {
		log.Printf("tiktok: delete account failed: %v", err)
		writeTikTokJSON(w, http.StatusBadGateway, map[string]any{"error": "Verbindung konnte nicht getrennt werden."})
		return
	}
	writeTikTokJSON(w, http.StatusOK, map[string]any{"disconnected": true})
}

// ---------------------------------------------------------------- helpers

func writeTikTokJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func visitorIDFromCookie(r *http.Request) (string, bool) {
	c, err := r.Cookie(tiktokCookieName)
	if err != nil || !tiktokUIDPattern.MatchString(c.Value) {
		return "", false
	}
	return c.Value, true
}

// ensureVisitorID returns the visitor's id, minting + setting the cookie the
// first time. HttpOnly + SameSite=Lax: sent on the same-origin fetches the
// page makes (also after the cross-site redirect back from TikTok), never
// readable by scripts, and not sent on cross-site POSTs.
func (h *TikTokHandler) ensureVisitorID(w http.ResponseWriter, r *http.Request) (string, error) {
	if uid, ok := visitorIDFromCookie(r); ok {
		return uid, nil
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	uid := hex.EncodeToString(buf)
	http.SetCookie(w, &http.Cookie{
		Name:     tiktokCookieName,
		Value:    uid,
		Path:     "/",
		MaxAge:   tiktokCookieMaxAge,
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})
	return uid, nil
}

// extractTikTokProfile pulls display_name / avatar URL out of a tool result
// whatever its exact nesting is (Composio's `data` may be an object or a
// JSON-encoded string, and the user fields may sit a level or two down).
func extractTikTokProfile(data json.RawMessage) (displayName, avatarURL string) {
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return "", ""
	}
	if s, ok := decoded.(string); ok { // double-encoded JSON
		var inner any
		if err := json.Unmarshal([]byte(s), &inner); err == nil {
			decoded = inner
		}
	}
	displayName = findStringField(decoded, "display_name")
	for _, key := range []string{"avatar_url", "avatar_large_url", "avatar_url_100"} {
		if avatarURL = findStringField(decoded, key); avatarURL != "" {
			break
		}
	}
	return displayName, avatarURL
}

func findStringField(v any, key string) string {
	switch t := v.(type) {
	case map[string]any:
		if s, ok := t[key].(string); ok && s != "" {
			return s
		}
		for _, child := range t {
			if s := findStringField(child, key); s != "" {
				return s
			}
		}
	case []any:
		for _, child := range t {
			if s := findStringField(child, key); s != "" {
				return s
			}
		}
	}
	return ""
}

// clientIP: nginx sets X-Real-IP to the real peer (see nginx/proxy.conf);
// fall back to the socket address when running without the proxy.
func clientIP(r *http.Request) string {
	if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// rateLimiter is a small in-memory fixed-window limiter, per key. Enough to
// stop one client hammering the Composio-backed endpoints; state is per
// process (single backend instance) and simply resets on restart.
type rateLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	hits   map[string][]time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{limit: limit, window: window, hits: make(map[string][]time.Time)}
}

func (l *rateLimiter) allow(key string) bool {
	now := time.Now()
	cutoff := now.Add(-l.window)

	l.mu.Lock()
	defer l.mu.Unlock()

	// Bound memory: drop keys whose hits have all expired.
	if len(l.hits) > 5000 {
		for k, times := range l.hits {
			if len(times) == 0 || times[len(times)-1].Before(cutoff) {
				delete(l.hits, k)
			}
		}
	}

	recent := l.hits[key][:0]
	for _, t := range l.hits[key] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	if len(recent) >= l.limit {
		l.hits[key] = recent
		return false
	}
	l.hits[key] = append(recent, now)
	return true
}
