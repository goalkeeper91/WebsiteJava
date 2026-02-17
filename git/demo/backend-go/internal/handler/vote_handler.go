package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

type VoteHandler struct {
	voteService  *service.VoteService
	sessionStore *sessions.CookieStore
	sessionName  string
}

func NewVoteHandler(
	voteService *service.VoteService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *VoteHandler {
	return &VoteHandler{
		voteService:  voteService,
		sessionStore: sessionStore,
		sessionName:  sessionName,
	}
}

func (h *VoteHandler) RegisterRoutes(router *mux.Router) {
	// Session management (owner)
	router.HandleFunc("/api/votes/sessions", h.GetSessions).Methods("GET")
	router.HandleFunc("/api/votes/sessions", h.StartSession).Methods("POST")
	router.HandleFunc("/api/votes/sessions/active", h.GetActiveSession).Methods("GET")
	router.HandleFunc("/api/votes/sessions/{id}", h.GetSession).Methods("GET")
	router.HandleFunc("/api/votes/sessions/{id}/close", h.CloseSession).Methods("POST")

	// Voting (viewers/participants – identified by query param or body)
	router.HandleFunc("/api/votes/sessions/{id}/vote", h.CastVote).Methods("POST")
	router.HandleFunc("/api/votes/sessions/{id}/results", h.GetResults).Methods("GET")
	router.HandleFunc("/api/votes/sessions/{id}/voted", h.HasVoted).Methods("GET")
}

// ============================================================
// SESSION ENDPOINTS (owner)
// ============================================================

func (h *VoteHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sessions, total, err := h.voteService.GetSessions(r.Context(), user.ID, pageSize, (page-1)*pageSize)
	if err != nil {
		log.Printf("Fehler beim Laden der Vote Sessions: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	h.respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       sessions,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *VoteHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	var req domain.VoteSessionCreateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	session, err := h.voteService.StartSession(r.Context(), user.ID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, session)
}

func (h *VoteHandler) GetActiveSession(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	session, err := h.voteService.GetActiveSession(r.Context(), user.ID)
	if err != nil {
		log.Printf("Fehler beim Laden der aktiven Session: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	if session == nil {
		h.respondJSON(w, http.StatusOK, nil)
		return
	}

	h.respondJSON(w, http.StatusOK, session)
}

func (h *VoteHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige Session ID", http.StatusBadRequest)
		return
	}

	session, err := h.voteService.GetSession(r.Context(), sessionID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, session)
}

func (h *VoteHandler) CloseSession(w http.ResponseWriter, r *http.Request) {
	user := h.requireUser(w, r)
	if user == nil {
		return
	}

	vars := mux.Vars(r)
	sessionID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige Session ID", http.StatusBadRequest)
		return
	}

	session, err := h.voteService.CloseSession(r.Context(), sessionID, user.ID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, session)
}

// ============================================================
// VOTING ENDPOINTS (participants)
// ============================================================

func (h *VoteHandler) CastVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige Session ID", http.StatusBadRequest)
		return
	}

	var req struct {
		VoterID  string `json:"voter_id"`  // Twitch username of voter
		Option   string `json:"option"`
		Weight   int    `json:"weight"`    // 1 = normal, 2 = subscriber bonus
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if req.VoterID == "" || req.Option == "" {
		http.Error(w, "voter_id und option sind erforderlich", http.StatusBadRequest)
		return
	}

	vote, err := h.voteService.CastVote(r.Context(), sessionID, req.VoterID, req.Option, req.Weight)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, vote)
}

func (h *VoteHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige Session ID", http.StatusBadRequest)
		return
	}

	results, total, err := h.voteService.GetResults(r.Context(), sessionID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"results":    results,
		"totalVotes": total,
	})
}

func (h *VoteHandler) HasVoted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige Session ID", http.StatusBadRequest)
		return
	}

	voterID := r.URL.Query().Get("voter_id")
	if voterID == "" {
		http.Error(w, "voter_id ist erforderlich", http.StatusBadRequest)
		return
	}

	hasVoted, err := h.voteService.HasVoted(r.Context(), sessionID, voterID)
	if err != nil {
		log.Printf("Fehler beim Prüfen des Votes: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]bool{"hasVoted": hasVoted})
}

// ============================================================
// HELPERS
// ============================================================

func (h *VoteHandler) requireUser(w http.ResponseWriter, r *http.Request) *UserSession {
	session, err := h.sessionStore.Get(r, h.sessionName)
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

	return &UserSession{ID: userID}
}

func (h *VoteHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrVoteSessionNotFound:
		http.Error(w, "Vote Session nicht gefunden", http.StatusNotFound)
	case domain.ErrVoteSessionClosed:
		http.Error(w, "Vote Session ist bereits geschlossen", http.StatusConflict)
	case domain.ErrVoteSessionNotActive:
		http.Error(w, "Vote Session ist nicht aktiv", http.StatusConflict)
	case domain.ErrAlreadyVoted:
		http.Error(w, "Du hast bereits abgestimmt", http.StatusConflict)
	case domain.ErrInvalidVoteOption:
		http.Error(w, "Ungültige Vote-Option", http.StatusBadRequest)
	case domain.ErrVoteLimitReached:
		http.Error(w, "Monatliches Vote-Limit erreicht", http.StatusForbidden)
	case domain.ErrN8NIntegrationNotReady:
		http.Error(w, "n8n Integration ist nicht konfiguriert", http.StatusBadRequest)
	case domain.ErrForbidden:
		http.Error(w, "Keine Berechtigung", http.StatusForbidden)
	default:
		log.Printf("Vote Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
	}
}

func (h *VoteHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}