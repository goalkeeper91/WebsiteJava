package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type ContactHandler struct {
	contactRepo repository.ContactRequestRepository
}

func NewContactHandler(contactRepo repository.ContactRequestRepository) *ContactHandler {
	return &ContactHandler{
		contactRepo: contactRepo,
	}
}

func (h *ContactHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/contact", h.SubmitContact).Methods("POST")
}

type ContactRequestDTO struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Subject      string `json:"subject"`
	Message      string `json:"message"`
	ConsentGiven bool   `json:"consentGiven"`
}

func (h *ContactHandler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	var dto ContactRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	contactRequest := domain.NewContactRequest(
		dto.Name,
		dto.Email,
		dto.Phone,
		dto.Subject,
		dto.Message,
		dto.ConsentGiven,
	)

	if err := contactRequest.Validate(); err != nil {
		http.Error(w, "Fehlende oder ungültige Felder", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.contactRepo.Create(ctx, contactRequest); err != nil {
		log.Printf("Fehler beim Speichern der Contact Request: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	// TODO: Optional - Discord Notification senden (später)
	log.Printf("✅ Contact Request erhalten von: %s (%s)", dto.Name, dto.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contactRequest)
}