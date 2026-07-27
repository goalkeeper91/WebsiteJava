package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

// TeamHandler manages who has delegated access to a channel's Twitch-chatbot
// dashboard. Unlike the 5 feature handlers, these routes always act "as
// yourself" (X-Acting-As is irrelevant here) - the one exception is
// RemoveMember, which also lets a member remove themselves from a team they
// don't own ("leave").
type TeamHandler struct {
	teamService  *service.TeamService
	sessionStore *sessions.CookieStore
	sessionName  string
}

func NewTeamHandler(teamService *service.TeamService, sessionStore *sessions.CookieStore, sessionName string) *TeamHandler {
	return &TeamHandler{teamService: teamService, sessionStore: sessionStore, sessionName: sessionName}
}

func (h *TeamHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/dashboard/team/members", h.ListMembers).Methods("GET")
	router.HandleFunc("/api/dashboard/team/invite", h.Invite).Methods("POST")
	router.HandleFunc("/api/dashboard/team/members/{memberTwitchID}", h.RemoveMember).Methods("DELETE")
	router.HandleFunc("/api/dashboard/team/leave/{ownerTwitchID}", h.Leave).Methods("DELETE")
	router.HandleFunc("/api/dashboard/team/managed-channels", h.ManagedChannels).Methods("GET")
}

func (h *TeamHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r, h.sessionStore, h.sessionName)
	if user == nil {
		return
	}

	members, err := h.teamService.ListMembers(r.Context(), user.ID)
	if err != nil {
		log.Printf("Team Fehler (Liste): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, members)
}

type InviteMemberRequest struct {
	Login string `json:"login"`
}

func (h *TeamHandler) Invite(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r, h.sessionStore, h.sessionName)
	if user == nil {
		return
	}

	var req InviteMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Login == "" {
		http.Error(w, "login ist erforderlich", http.StatusBadRequest)
		return
	}

	member, err := h.teamService.InviteMember(r.Context(), user.ID, req.Login)
	if err == domain.ErrTwitchUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err == domain.ErrCannotInviteSelf || err == domain.ErrAlreadyTeamMember {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("Team Fehler (Einladen): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, member)
}

// RemoveMember is the owner-side action: remove memberTwitchID from MY
// (the caller's) team. See Leave for the member-side "leave a team I don't
// own" action - a bare member ID alone can't disambiguate which team to
// leave, since a member can have access to more than one channel.
func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r, h.sessionStore, h.sessionName)
	if user == nil {
		return
	}

	memberTwitchID := mux.Vars(r)["memberTwitchID"]

	if err := h.teamService.RemoveMember(r.Context(), user.ID, memberTwitchID); err != nil {
		log.Printf("Team Fehler (Entfernen): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Leave is the member-side action: remove the caller (as member) from
// ownerTwitchID's team.
func (h *TeamHandler) Leave(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r, h.sessionStore, h.sessionName)
	if user == nil {
		return
	}

	ownerTwitchID := mux.Vars(r)["ownerTwitchID"]

	if err := h.teamService.RemoveMember(r.Context(), ownerTwitchID, user.ID); err != nil {
		log.Printf("Team Fehler (Verlassen): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) ManagedChannels(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r, h.sessionStore, h.sessionName)
	if user == nil {
		return
	}

	channels, err := h.teamService.ListManagedChannels(r.Context(), user.ID)
	if err != nil {
		log.Printf("Team Fehler (Verwaltbare Kanäle): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, channels)
}

func (h *TeamHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
