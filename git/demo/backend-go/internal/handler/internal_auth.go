package handler

import (
	"crypto/subtle"
	"net/http"
)

// InternalSecretHeader is the header server-to-server callers (the Twitch bot, n8n)
// must set to authenticate against endpoints that aren't backed by a browser session.
const InternalSecretHeader = "X-Bot-Secret"

// requireInternalSecret checks the request against a pre-shared secret using a
// constant-time comparison. Used for endpoints called by trusted backend services
// (the bot, n8n) rather than by a logged-in browser session.
func requireInternalSecret(w http.ResponseWriter, r *http.Request, expectedSecret string) bool {
	if expectedSecret == "" {
		http.Error(w, "Service nicht konfiguriert", http.StatusServiceUnavailable)
		return false
	}

	provided := r.Header.Get(InternalSecretHeader)
	if provided == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(expectedSecret)) != 1 {
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return false
	}

	return true
}
