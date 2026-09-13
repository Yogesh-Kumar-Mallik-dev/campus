/**
 * BLOCK_API_MESS_HANDLER_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for campus dining, QR punches, and leave rebates.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package mess

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/mess"
)

type Handler struct {
	service *mess.Service
}

func NewHandler(service *mess.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/mess", func(r chi.Router) {
		// Halls
		r.Post("/halls", h.CreateMessHall)
		r.Get("/halls", h.ListMessHalls)

		// Menu Items
		r.Post("/menu-items", h.AddMenuItem)
		r.Get("/menu-items", h.ListMenu)

		// Subscriptions
		r.Post("/subscriptions", h.SubscribeStudent)

		// Dining Tokens & QR Punch
		r.Post("/tokens/daily", h.GenerateDailyToken)
		r.Post("/tokens/redeem", h.RedeemDiningToken)
		r.Get("/tokens", h.ListTokens)

		// Rebates
		r.Post("/rebates", h.ApplyRebate)
		r.Post("/rebates/{id}/approve", h.ApproveRebate)
		r.Get("/rebates", h.ListRebates)

		// Feedback
		r.Post("/feedbacks", h.SubmitFeedback)
		r.Get("/feedbacks", h.ListFeedbacks)
	})
}

// CreateMessHall handles POST /api/v1/mess/halls
func (h *Handler) CreateMessHall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID    string  `json:"tenant_id"`
		Name        string  `json:"name"`
		Code        string  `json:"code"`
		Capacity    int     `json:"capacity"`
		Location    string  `json:"location"`
		CatererName string  `json:"caterer_name"`
		ManagerID   *string `json:"manager_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	hall, err := h.service.CreateMessHall(r.Context(), req.TenantID, req.Name, req.Code, req.Location, req.CatererName, req.Capacity, req.ManagerID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"mess_hall": hall})
}

// ListMessHalls handles GET /api/v1/mess/halls
func (h *Handler) ListMessHalls(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	halls, err := h.service.ListMessHalls(r.Context(), tenantID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"mess_halls": halls})
}

// AddMenuItem handles POST /api/v1/mess/menu-items
func (h *Handler) AddMenuItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID     string            `json:"tenant_id"`
		MessHallID   string            `json:"mess_hall_id"`
		DayOfWeek    mess.DayOfWeek    `json:"day_of_week"`
		MealType     mess.MealType     `json:"meal_type"`
		DietCategory mess.DietCategory `json:"diet_category"`
		Title        string            `json:"title"`
		Description  string            `json:"description"`
		Calories     int               `json:"calories"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	item, err := h.service.AddMenuItem(r.Context(), req.TenantID, req.MessHallID, req.DayOfWeek, req.MealType, req.DietCategory, req.Title, req.Description, req.Calories)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"menu_item": item})
}

// ListMenu handles GET /api/v1/mess/menu-items
func (h *Handler) ListMenu(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	messHallID := r.URL.Query().Get("mess_hall_id")
	day := mess.DayOfWeek(r.URL.Query().Get("day_of_week"))
	mealType := mess.MealType(r.URL.Query().Get("meal_type"))

	if tenantID == "" || messHallID == "" {
		problem.BadRequest(w, r, "tenant_id and mess_hall_id are required", "MISSING_PARAMS")
		return
	}

	menu, err := h.service.ListMenu(r.Context(), tenantID, messHallID, day, mealType)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"menu_items": menu})
}

// SubscribeStudent handles POST /api/v1/mess/subscriptions
func (h *Handler) SubscribeStudent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID       string            `json:"tenant_id"`
		StudentID      string            `json:"student_id"`
		MessHallID     string            `json:"mess_hall_id"`
		PlanType       mess.PlanType     `json:"plan_type"`
		DietPreference mess.DietCategory `json:"diet_preference"`
		StartDate      time.Time         `json:"start_date"`
		EndDate        time.Time         `json:"end_date"`
		MonthlyFee     float64           `json:"monthly_fee"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	sub, err := h.service.SubscribeStudent(r.Context(), req.TenantID, req.StudentID, req.MessHallID, req.PlanType, req.DietPreference, req.StartDate, req.EndDate, req.MonthlyFee)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"subscription": sub})
}

// GenerateDailyToken handles POST /api/v1/mess/tokens/daily
func (h *Handler) GenerateDailyToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID  string        `json:"tenant_id"`
		StudentID string        `json:"student_id"`
		MealDate  string        `json:"meal_date"`
		MealType  mess.MealType `json:"meal_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	tok, err := h.service.GenerateDailyToken(r.Context(), req.TenantID, req.StudentID, req.MealDate, req.MealType)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"token": tok})
}

// RedeemDiningToken handles POST /api/v1/mess/tokens/redeem
func (h *Handler) RedeemDiningToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID       string `json:"tenant_id"`
		TokenCode      string `json:"token_code"`
		DeviceReaderID string `json:"device_reader_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	tok, err := h.service.RedeemDiningToken(r.Context(), req.TenantID, req.TokenCode, req.DeviceReaderID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"token": tok, "redeemed": true})
}

// ListTokens handles GET /api/v1/mess/tokens
func (h *Handler) ListTokens(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	studentID := r.URL.Query().Get("student_id")
	if tenantID == "" || studentID == "" {
		problem.BadRequest(w, r, "tenant_id and student_id are required", "MISSING_PARAMS")
		return
	}

	tokens, err := h.service.ListStudentTokens(r.Context(), tenantID, studentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"tokens": tokens})
}

// ApplyRebate handles POST /api/v1/mess/rebates
func (h *Handler) ApplyRebate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID   string    `json:"tenant_id"`
		StudentID  string    `json:"student_id"`
		FromDate   time.Time `json:"from_date"`
		ToDate     time.Time `json:"to_date"`
		Reason     string    `json:"reason"`
		GatePassID string    `json:"gate_pass_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	rebate, err := h.service.ApplyRebate(r.Context(), req.TenantID, req.StudentID, req.FromDate, req.ToDate, req.Reason, req.GatePassID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"rebate": rebate})
}

// ApproveRebate handles POST /api/v1/mess/rebates/{id}/approve
func (h *Handler) ApproveRebate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		TenantID   string `json:"tenant_id"`
		ApproverID string `json:"approver_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	rebate, err := h.service.ApproveRebate(r.Context(), req.TenantID, id, req.ApproverID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"rebate": rebate})
}

// ListRebates handles GET /api/v1/mess/rebates
func (h *Handler) ListRebates(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	studentID := r.URL.Query().Get("student_id")

	rebates, err := h.service.ListStudentRebates(r.Context(), tenantID, studentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"rebates": rebates})
}

// SubmitFeedback handles POST /api/v1/mess/feedbacks
func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID      string        `json:"tenant_id"`
		StudentID     string        `json:"student_id"`
		MessHallID    string        `json:"mess_hall_id"`
		MealDate      string        `json:"meal_date"`
		MealType      mess.MealType `json:"meal_type"`
		FoodRating    int           `json:"food_rating"`
		HygieneRating int           `json:"hygiene_rating"`
		Comments      string        `json:"comments"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	fb, err := h.service.SubmitFeedback(r.Context(), req.TenantID, req.StudentID, req.MessHallID, req.MealDate, req.MealType, req.FoodRating, req.HygieneRating, req.Comments)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"feedback": fb})
}

// ListFeedbacks handles GET /api/v1/mess/feedbacks
func (h *Handler) ListFeedbacks(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	messHallID := r.URL.Query().Get("mess_hall_id")
	if tenantID == "" || messHallID == "" {
		problem.BadRequest(w, r, "tenant_id and mess_hall_id are required", "MISSING_PARAMS")
		return
	}

	feedbacks, err := h.service.ListFeedbacks(r.Context(), tenantID, messHallID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"feedbacks": feedbacks})
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, mess.ErrMessHallNotFound),
		errors.Is(err, mess.ErrMenuItemNotFound),
		errors.Is(err, mess.ErrSubscriptionNotFound),
		errors.Is(err, mess.ErrTokenNotFound),
		errors.Is(err, mess.ErrRebateNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, mess.ErrMessHallCodeExists),
		errors.Is(err, mess.ErrActiveSubscriptionExists),
		errors.Is(err, mess.ErrTokenAlreadyRedeemed),
		errors.Is(err, mess.ErrDoubleRedemptionAttempt):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	case errors.Is(err, mess.ErrInsufficientRebateDays),
		errors.Is(err, mess.ErrInvalidDietPreference),
		errors.Is(err, mess.ErrInvalidDateRange),
		errors.Is(err, mess.ErrTokenExpired):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	default:
		var domErr *mess.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected dining error occurred")
	}
}
