package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// requireUser extracts user from session
// This is the base method that all handlers should use
func requireUser(w http.ResponseWriter, r *http.Request, sessionStore *sessions.CookieStore, sessionName string) *UserSession {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Printf("Fehler beim Laden der Session: %v", err)
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return nil
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Nicht authentifiziert", http.StatusUnauthorized)
		return nil
	}

	// Check if admin flag is in session
	isAdmin, _ := session.Values["is_admin"].(bool)

	return &UserSession{
		ID:      userID,
		IsAdmin: isAdmin,
	}
}

// requireAdmin checks if user is admin
func requireAdmin(w http.ResponseWriter, r *http.Request, sessionStore *sessions.CookieStore, sessionName string) *UserSession {
	user := requireUser(w, r, sessionStore, sessionName)
	if user == nil {
		return nil
	}

	if !user.IsAdmin {
		http.Error(w, "Zugriff verweigert - Admin-Rechte erforderlich", http.StatusForbidden)
		return nil
	}

	return user
}

// AdminMiddleware can be used as middleware for admin-only routes
func AdminMiddleware(sessionStore *sessions.CookieStore, sessionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := requireAdmin(w, r, sessionStore, sessionName)
			if user == nil {
				return
			}

			// Store user in context for downstream handlers
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves user from context (set by middleware)
func GetUserFromContext(ctx context.Context) *UserSession {
	user, ok := ctx.Value("user").(*UserSession)
	if !ok {
		return nil
	}
	return user
}