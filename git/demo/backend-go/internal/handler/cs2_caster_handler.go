package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/service"
)

type CS2CasterHandler struct {
	casterService *service.CS2CasterService
	sessionStore  *sessions.CookieStore
	sessionName   string
	teamService   *service.TeamService
}

func NewCS2CasterHandler(
	casterService *service.CS2CasterService,
	sessionStore *sessions.CookieStore,
	sessionName string,
	teamService *service.TeamService,
) *CS2CasterHandler {
	return &CS2CasterHandler{
		casterService: casterService,
		sessionStore:  sessionStore,
		sessionName:   sessionName,
		teamService:   teamService,
	}
}

func (h *CS2CasterHandler) RegisterRoutes(router *mux.Router) {
	// Öffentlich - der Pfad-Token selbst ist die Authentifizierung, kein Session-Cookie nötig.
	router.HandleFunc("/api/cs2/gsi/{token}", h.IngestGSI).Methods("POST")

	router.HandleFunc("/api/dashboard/cs2/settings", h.GetSettings).Methods("GET")
	router.HandleFunc("/api/dashboard/cs2/settings", h.UpdateSettings).Methods("PUT")
	router.HandleFunc("/api/dashboard/cs2/live-status", h.GetLiveStatus).Methods("GET")
	router.HandleFunc("/api/dashboard/cs2/notes", h.ListNotes).Methods("GET")
	router.HandleFunc("/api/dashboard/cs2/notes", h.CreateNote).Methods("POST")
	router.HandleFunc("/api/dashboard/cs2/notes/{id}", h.UpdateNote).Methods("PUT")
	router.HandleFunc("/api/dashboard/cs2/notes/{id}", h.DeleteNote).Methods("DELETE")
}

// ============================================================
// GSI-INGESTION (öffentlich, Token-Auth über den Pfad)
// ============================================================

func (h *CS2CasterHandler) IngestGSI(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	if token == "" {
		http.Error(w, "Token fehlt", http.StatusBadRequest)
		return
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Fehler beim Lesen des Payloads", http.StatusBadRequest)
		return
	}

	if err := h.casterService.IngestGSIPayload(r.Context(), token, raw); err != nil {
		if err == domain.ErrCS2InvalidGSIToken {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		log.Printf("CS2 GSI Fehler: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================
// DASHBOARD (Session-Auth, Team-Zugriff-kompatibel)
// ============================================================

func (h *CS2CasterHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	settings, err := h.casterService.GetOrCreateSettings(r.Context(), channelID)
	if err != nil {
		log.Printf("CS2 Fehler (Settings laden): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *CS2CasterHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	var input domain.CS2CasterSettingsUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if err := h.casterService.UpdateSettings(r.Context(), channelID, input); err != nil {
		log.Printf("CS2 Fehler (Settings aktualisieren): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	settings, err := h.casterService.GetOrCreateSettings(r.Context(), channelID)
	if err != nil {
		log.Printf("CS2 Fehler (Settings nachladen): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *CS2CasterHandler) GetLiveStatus(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	status := h.casterService.GetLiveStatus(channelID)
	h.respondJSON(w, http.StatusOK, status)
}

func (h *CS2CasterHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	notes, err := h.casterService.ListNotes(r.Context(), channelID)
	if err != nil {
		log.Printf("CS2 Fehler (Notizen laden): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, notes)
}

type createCS2NoteRequest struct {
	SubjectType domain.CS2NoteSubjectType `json:"subject_type"`
	SubjectName string                    `json:"subject_name"`
	Content     string                    `json:"content"`
}

func (h *CS2CasterHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	var req createCS2NoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	note, err := h.casterService.CreateNote(r.Context(), channelID, domain.CS2NoteCreateInput{
		SubjectType: req.SubjectType,
		SubjectName: req.SubjectName,
		Content:     req.Content,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondJSON(w, http.StatusOK, note)
}

type updateCS2NoteRequest struct {
	SubjectName *string `json:"subject_name,omitempty"`
	Content     *string `json:"content,omitempty"`
}

func (h *CS2CasterHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	noteID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige Notiz-ID", http.StatusBadRequest)
		return
	}

	var req updateCS2NoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if err := h.casterService.UpdateNote(r.Context(), channelID, noteID, domain.CS2NoteUpdateInput{
		SubjectName: req.SubjectName,
		Content:     req.Content,
	}); err != nil {
		if err == domain.ErrCS2NoteNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("CS2 Fehler (Notiz aktualisieren): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CS2CasterHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	_, channelID, ok := resolveEffectiveTwitchID(w, r, h.sessionStore, h.sessionName, h.teamService)
	if !ok {
		return
	}

	noteID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Ungültige Notiz-ID", http.StatusBadRequest)
		return
	}

	if err := h.casterService.DeleteNote(r.Context(), channelID, noteID); err != nil {
		if err == domain.ErrCS2NoteNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("CS2 Fehler (Notiz löschen): %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================
// HELPERS
// ============================================================

func (h *CS2CasterHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
