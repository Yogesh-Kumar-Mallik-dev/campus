/**
 * BLOCK_API_HELPDESK_HANDLER_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for categories, ticketing, live chat desk, and SLA escalations.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package helpdesk

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/helpdesk"
)

type Handler struct {
	service helpdesk.Service
}

func NewHandler(service helpdesk.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/helpdesk", func(r chi.Router) {
		r.Post("/categories", h.CreateCategory)
		r.Get("/categories", h.ListCategories)

		r.Post("/tickets", h.CreateTicket)
		r.Get("/tickets", h.ListTickets)
		r.Get("/tickets/{id}", h.GetTicket)

		r.Post("/tickets/{id}/assign", h.AssignTicket)
		r.Post("/tickets/{id}/status", h.UpdateTicketStatus)
		r.Post("/tickets/{id}/messages", h.AddMessage)
		r.Get("/tickets/{id}/messages", h.ListMessages)
		r.Post("/tickets/{id}/escalate", h.EscalateTicket)
		r.Post("/tickets/{id}/rate", h.RateTicket)
	})
}

// CreateCategory handles POST /api/v1/helpdesk/categories
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID           string                    `json:"tenant_id"`
		Code               string                    `json:"code"`
		Name               string                    `json:"name"`
		Description        string                    `json:"description"`
		DefaultPriority    helpdesk.HelpdeskPriority `json:"default_priority"`
		SLAResponseHours   int                       `json:"sla_response_hours"`
		SLAResolutionHours int                       `json:"sla_resolution_hours"`
		ActorID            string                    `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.Code == "" || req.Name == "" {
		problem.BadRequest(w, r, "tenant_id, code, and name are required", "VALIDATION_ERROR")
		return
	}

	cat, err := h.service.CreateCategory(r.Context(), helpdesk.CreateCategoryRequest{
		TenantID:           req.TenantID,
		Code:               req.Code,
		Name:               req.Name,
		Description:        req.Description,
		DefaultPriority:    req.DefaultPriority,
		SLAResponseHours:   req.SLAResponseHours,
		SLAResolutionHours: req.SLAResolutionHours,
		ActorID:            req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cat)
}

// ListCategories handles GET /api/v1/helpdesk/categories
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	activeOnly := r.URL.Query().Get("active_only") != "false"
	categories, err := h.service.ListCategories(r.Context(), tenantID, activeOnly)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(categories)
}

// CreateTicket handles POST /api/v1/helpdesk/tickets
func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID    string                     `json:"tenant_id"`
		CategoryID  string                     `json:"category_id"`
		RequesterID string                     `json:"requester_id"`
		Title       string                     `json:"title"`
		Description string                     `json:"description"`
		Priority    *helpdesk.HelpdeskPriority `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.CategoryID == "" || req.RequesterID == "" || req.Title == "" || req.Description == "" {
		problem.BadRequest(w, r, "tenant_id, category_id, requester_id, title, and description are required", "VALIDATION_ERROR")
		return
	}

	ticket, err := h.service.CreateTicket(r.Context(), helpdesk.CreateTicketRequest{
		TenantID:    req.TenantID,
		CategoryID:  req.CategoryID,
		RequesterID: req.RequesterID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ticket)
}

// ListTickets handles GET /api/v1/helpdesk/tickets
func (h *Handler) ListTickets(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	filter := helpdesk.TicketFilter{
		TenantID: tenantID,
	}
	if reqID := r.URL.Query().Get("requester_id"); reqID != "" {
		filter.RequesterID = &reqID
	}
	if staffID := r.URL.Query().Get("assigned_staff_id"); staffID != "" {
		filter.AssignedStaffID = &staffID
	}
	if catID := r.URL.Query().Get("category_id"); catID != "" {
		filter.CategoryID = &catID
	}
	if st := r.URL.Query().Get("status"); st != "" {
		status := helpdesk.HelpdeskTicketStatus(st)
		filter.Status = &status
	}
	if pr := r.URL.Query().Get("priority"); pr != "" {
		priority := helpdesk.HelpdeskPriority(pr)
		filter.Priority = &priority
	}

	tickets, err := h.service.ListTickets(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tickets)
}

// GetTicket handles GET /api/v1/helpdesk/tickets/{id}
func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	ticketID := chi.URLParam(r, "id")

	if tenantID == "" || ticketID == "" {
		problem.BadRequest(w, r, "tenant_id and ticket id are required", "VALIDATION_ERROR")
		return
	}

	ticket, err := h.service.GetTicket(r.Context(), tenantID, ticketID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ticket)
}

// AssignTicket handles POST /api/v1/helpdesk/tickets/{id}/assign
func (h *Handler) AssignTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	var req struct {
		TenantID string `json:"tenant_id"`
		StaffID  string `json:"staff_id"`
		ActorID  string `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || ticketID == "" || req.StaffID == "" {
		problem.BadRequest(w, r, "tenant_id, ticket id, and staff_id are required", "VALIDATION_ERROR")
		return
	}

	ticket, err := h.service.AssignTicket(r.Context(), req.TenantID, ticketID, req.StaffID, req.ActorID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ticket)
}

// UpdateTicketStatus handles POST /api/v1/helpdesk/tickets/{id}/status
func (h *Handler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	var req struct {
		TenantID string                        `json:"tenant_id"`
		Status   helpdesk.HelpdeskTicketStatus `json:"status"`
		ActorID  string                        `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || ticketID == "" || req.Status == "" {
		problem.BadRequest(w, r, "tenant_id, ticket id, and status are required", "VALIDATION_ERROR")
		return
	}

	ticket, err := h.service.UpdateTicketStatus(r.Context(), req.TenantID, ticketID, req.Status, req.ActorID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ticket)
}

// AddMessage handles POST /api/v1/helpdesk/tickets/{id}/messages
func (h *Handler) AddMessage(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	var req struct {
		TenantID       string   `json:"tenant_id"`
		SenderID       string   `json:"sender_id"`
		IsStaffReply   bool     `json:"is_staff_reply"`
		IsInternalNote bool     `json:"is_internal_note"`
		Message        string   `json:"message"`
		Attachments    []string `json:"attachments"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || ticketID == "" || req.SenderID == "" || req.Message == "" {
		problem.BadRequest(w, r, "tenant_id, ticket id, sender_id, and message are required", "VALIDATION_ERROR")
		return
	}

	msg, err := h.service.AddMessage(r.Context(), helpdesk.AddMessageRequest{
		TenantID:       req.TenantID,
		TicketID:       ticketID,
		SenderID:       req.SenderID,
		IsStaffReply:   req.IsStaffReply,
		IsInternalNote: req.IsInternalNote,
		Message:        req.Message,
		Attachments:    req.Attachments,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(msg)
}

// ListMessages handles GET /api/v1/helpdesk/tickets/{id}/messages
func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	ticketID := chi.URLParam(r, "id")
	userID := r.URL.Query().Get("user_id")
	isStaff := r.URL.Query().Get("is_staff") == "true"

	if tenantID == "" || ticketID == "" || userID == "" {
		problem.BadRequest(w, r, "tenant_id, ticket id, and user_id are required", "VALIDATION_ERROR")
		return
	}

	messages, err := h.service.ListMessages(r.Context(), tenantID, ticketID, userID, isStaff)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(messages)
}

// EscalateTicket handles POST /api/v1/helpdesk/tickets/{id}/escalate
func (h *Handler) EscalateTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	var req struct {
		TenantID      string  `json:"tenant_id"`
		EscalatedToID *string `json:"escalated_to_id"`
		Reason        string  `json:"reason"`
		ActorID       string  `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || ticketID == "" || req.Reason == "" {
		problem.BadRequest(w, r, "tenant_id, ticket id, and reason are required", "VALIDATION_ERROR")
		return
	}

	esc, err := h.service.EscalateTicket(r.Context(), helpdesk.EscalateTicketRequest{
		TenantID:      req.TenantID,
		TicketID:      ticketID,
		EscalatedToID: req.EscalatedToID,
		Reason:        req.Reason,
		ActorID:       req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(esc)
}

// RateTicket handles POST /api/v1/helpdesk/tickets/{id}/rate
func (h *Handler) RateTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	var req struct {
		TenantID string `json:"tenant_id"`
		UserID   string `json:"user_id"`
		Rating   int    `json:"rating"`
		Feedback string `json:"feedback"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || ticketID == "" || req.UserID == "" || req.Rating == 0 {
		problem.BadRequest(w, r, "tenant_id, ticket id, user_id, and rating are required", "VALIDATION_ERROR")
		return
	}

	ticket, err := h.service.RateTicket(r.Context(), helpdesk.RateTicketRequest{
		TenantID: req.TenantID,
		TicketID: ticketID,
		UserID:   req.UserID,
		Rating:   req.Rating,
		Feedback: req.Feedback,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ticket)
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, helpdesk.ErrCategoryNotFound),
		errors.Is(err, helpdesk.ErrTicketNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, helpdesk.ErrCategoryCodeDuplicate),
		errors.Is(err, helpdesk.ErrEscalationAlreadyPending):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	case errors.Is(err, helpdesk.ErrInternalNoteUnauthorized),
		errors.Is(err, helpdesk.ErrUnauthorizedTicketAccess):
		problem.Forbidden(w, r, err.Error(), "FORBIDDEN")

	case errors.Is(err, helpdesk.ErrInvalidTicketTransition),
		errors.Is(err, helpdesk.ErrTicketClosed),
		errors.Is(err, helpdesk.ErrInvalidRatingScore),
		errors.Is(err, helpdesk.ErrRatingNotAllowed):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	default:
		var domErr *helpdesk.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected helpdesk error occurred")
	}
}
