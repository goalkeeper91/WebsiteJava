package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/infrastructure/redis"
	"demo/backend-go/internal/repository"
	"demo/backend-go/internal/service"
)

// AdminCustomerHandler stellt die Admin-"Kunden/Abos"-Übersicht bereit
// (Phase 5 der SaaS-Roadmap): Kundenliste inkl. Tarif, manueller Tarif-
// Override für Support-Fälle, sowie eine grobe MRR-/Clip-Detector-
// Auslastungs-Übersicht.
type AdminCustomerHandler struct {
	subscriptionService *service.SubscriptionService
	userRepo            repository.UserRepository
	redisService        *redis.RedisService
	sessionStore        *sessions.CookieStore
	sessionName         string
}

func NewAdminCustomerHandler(
	subscriptionService *service.SubscriptionService,
	userRepo repository.UserRepository,
	redisService *redis.RedisService,
	sessionStore *sessions.CookieStore,
	sessionName string,
) *AdminCustomerHandler {
	return &AdminCustomerHandler{
		subscriptionService: subscriptionService,
		userRepo:            userRepo,
		redisService:        redisService,
		sessionStore:        sessionStore,
		sessionName:         sessionName,
	}
}

func (h *AdminCustomerHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/admin/customers", h.handleListCustomers).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/customers/stats", h.handleStats).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/customers/{twitchId}/tier", h.handleSetTier).Methods(http.MethodPut)
}

// isAdmin: gleiches Muster wie DiscordGuildHandler/BotStatusHandler (kein
// Session-basiertes admin_helpers.go-Middleware wird im Rest des Codes
// tatsächlich genutzt).
func (h *AdminCustomerHandler) isAdmin(r *http.Request) bool {
	session, err := h.sessionStore.Get(r, h.sessionName)
	if err != nil {
		return false
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		return false
	}

	user, err := h.userRepo.GetByTwitchID(r.Context(), userID)
	if err != nil || user == nil {
		return false
	}

	return user.IsAdmin
}

type adminCustomerResponse struct {
	TwitchID             string     `json:"twitch_id"`
	Username             string     `json:"username"`
	Email                string     `json:"email,omitempty"`
	IsAdmin              bool       `json:"is_admin"`
	CreatedAt            time.Time  `json:"created_at"`
	TierID               string     `json:"tier_id"`
	TierName             string     `json:"tier_name"`
	Status               string     `json:"status"`
	IsActive             bool       `json:"is_active"`
	BillingCycle         *string    `json:"billing_cycle,omitempty"`
	ExpiresAt            *time.Time `json:"expires_at,omitempty"`
	PriceMonthly         float64    `json:"price_monthly"`
	PriceYearly          float64    `json:"price_yearly"`
	PaddleCustomerID     *string    `json:"paddle_customer_id,omitempty"`
	PaddleSubscriptionID *string    `json:"paddle_subscription_id,omitempty"`
}

func toAdminCustomerResponse(row *domain.AdminCustomerRow) adminCustomerResponse {
	resp := adminCustomerResponse{
		TwitchID:  row.User.TwitchID,
		Username:  row.User.Username,
		Email:     row.User.Email,
		IsAdmin:   row.User.IsAdmin,
		CreatedAt: row.User.CreatedAt,
	}

	sub := row.Subscription
	if sub == nil {
		// Kein bisheriger Zugriff auf GetSubscription -> implizit aktiver
		// Free-Kunde (jeder Nutzer bekommt Free ohne jede Zahlung).
		resp.TierID = string(domain.TierFree)
		resp.TierName = "Free"
		resp.Status = "none"
		resp.IsActive = true
		return resp
	}

	resp.TierID = string(sub.TierID)
	resp.Status = string(sub.Status)
	resp.IsActive = sub.IsActive()
	resp.ExpiresAt = sub.ExpiresAt
	resp.PaddleCustomerID = sub.PaddleCustomerID
	resp.PaddleSubscriptionID = sub.PaddleSubscriptionID
	if sub.BillingCycle != nil {
		v := string(*sub.BillingCycle)
		resp.BillingCycle = &v
	}
	if sub.Tier != nil {
		resp.TierName = sub.Tier.Name
		resp.PriceMonthly = sub.Tier.PriceMonthly
		resp.PriceYearly = sub.Tier.PriceYearly
	} else {
		resp.TierName = string(sub.TierID)
	}

	return resp
}

func (h *AdminCustomerHandler) handleListCustomers(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Error(w, "Zugriff verweigert - Admin-Rechte erforderlich", http.StatusForbidden)
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

	rows, total, err := h.subscriptionService.AdminListCustomers(r.Context(), page, pageSize)
	if err != nil {
		log.Printf("Fehler beim Laden der Kundenliste: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	data := make([]adminCustomerResponse, 0, len(rows))
	for _, row := range rows {
		data = append(data, toAdminCustomerResponse(row))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	h.respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

type setTierRequest struct {
	TierID string `json:"tier_id"`
}

func (h *AdminCustomerHandler) handleSetTier(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Error(w, "Zugriff verweigert - Admin-Rechte erforderlich", http.StatusForbidden)
		return
	}

	twitchID := mux.Vars(r)["twitchId"]

	user, err := h.userRepo.GetByTwitchID(r.Context(), twitchID)
	if err != nil || user == nil {
		http.Error(w, "Nutzer nicht gefunden", http.StatusNotFound)
		return
	}

	var req setTierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if err := h.subscriptionService.AdminSetTier(r.Context(), twitchID, domain.TierID(req.TierID)); err != nil {
		log.Printf("Fehler beim manuellen Tarif-Override für %s: %v", twitchID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type adminCustomerStatsResponse struct {
	MRR             float64                `json:"mrr"`
	ActiveCustomers int                    `json:"active_customers"`
	ClipDetector    *clipDetectorStatusDTO `json:"clip_detector"`
}

type clipDetectorStatusDTO struct {
	ActiveChannels []clipDetectorChannelDTO `json:"active_channels"`
	MaxConcurrent  int                      `json:"max_concurrent"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

type clipDetectorChannelDTO struct {
	TwitchUserID string `json:"twitch_user_id"`
	Login        string `json:"login"`
}

func (h *AdminCustomerHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Error(w, "Zugriff verweigert - Admin-Rechte erforderlich", http.StatusForbidden)
		return
	}

	mrr, activeCustomers, err := h.subscriptionService.AdminGetStats(r.Context())
	if err != nil {
		log.Printf("Fehler beim Berechnen der Admin-Kundenstatistik: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	resp := adminCustomerStatsResponse{
		MRR:             mrr,
		ActiveCustomers: activeCustomers,
	}

	raw, err := h.redisService.GetClient().Get(r.Context(), "clip_detector:status").Bytes()
	if err == nil {
		var status clipDetectorStatusDTO
		if err := json.Unmarshal(raw, &status); err == nil {
			resp.ClipDetector = &status
		}
	}

	h.respondJSON(w, http.StatusOK, resp)
}

func (h *AdminCustomerHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
