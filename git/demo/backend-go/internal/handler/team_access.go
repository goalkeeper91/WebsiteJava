package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"

	"demo/backend-go/internal/service"
)

// resolveEffectiveTwitchID returns the Twitch ID whose Twitch-chatbot
// settings this request should act on: the caller's own ID normally, or -
// if the caller sent X-Acting-As and currently has team access to that
// channel - the delegated owner's ID instead. Access is re-checked on every
// call (never cached in the session), so a removed team member loses access
// on their very next request, not just at their next login.
func resolveEffectiveTwitchID(w http.ResponseWriter, r *http.Request, sessionStore *sessions.CookieStore, sessionName string, teamService *service.TeamService) (callerID string, effectiveID string, ok bool) {
	user := requireUser(w, r, sessionStore, sessionName)
	if user == nil {
		return "", "", false
	}

	actingAs := r.Header.Get("X-Acting-As")
	if actingAs == "" || actingAs == user.ID {
		return user.ID, user.ID, true
	}

	hasAccess, err := teamService.HasAccess(r.Context(), actingAs, user.ID)
	if err != nil {
		log.Printf("Team-Zugriffsprüfung fehlgeschlagen: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return "", "", false
	}
	if !hasAccess {
		http.Error(w, "Kein Zugriff auf diesen Kanal", http.StatusForbidden)
		return "", "", false
	}

	return user.ID, actingAs, true
}
