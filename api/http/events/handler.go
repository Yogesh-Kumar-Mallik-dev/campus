/**
 * BLOCK_API_EVENTS_HANDLER_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for campus events, venue conflict prevention, and QR gate admissions.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package events

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/events"
)

type Handler struct {
	service events.Service
}

func NewHandler(service events.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/events", func(r chi.Router) {
		r.Post("/", h.CreateEvent)
		r.Get("/", h.ListEvents)
		r.Get("/{id}", h.GetEvent)

		r.Post("/{id}/tickets", h.BookTicket)
		r.Get("/{id}/tickets", h.ListEventTickets)
		r.Get("/tickets/my", h.ListStudentTickets)
		r.Post("/tickets/checkin", h.CheckInTicket)
	})
}

// CreateEvent handles POST /api/v1/events
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID             string               `json:"tenant_id"`
		Title                string               `json:"title"`
		Description          string               `json:"description"`
		Category             events.EventCategory `json:"category"`
		VenueName            string               `json:"venue_name"`
		VenueCapacity        int                  `json:"venue_capacity"`
		StartTime            string               `json:"start_time"`
		EndTime              string               `json:"end_time"`
		RegistrationDeadline string               `json:"registration_deadline"`
		OrganizerID          string               `json:"organizer_id"`
		IsTicketed           bool                 `json:"is_ticketed"`
		TicketPrice          float64              `json:"ticket_price"`
		MaxTickets           int                  `json:"max_tickets"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.Title == "" || req.VenueName == "" || req.StartTime == "" || req.EndTime == "" {
		problem.BadRequest(w, r, "tenant_id, title, venue_name, start_time, and end_time are required", "VALIDATION_ERROR")
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		problem.BadRequest(w, r, "start_time must be formatted in RFC3339", "INVALID_START_TIME")
		return
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		problem.BadRequest(w, r, "end_time must be formatted in RFC3339", "INVALID_END_TIME")
		return
	}

	regDeadline := start
	if req.RegistrationDeadline != "" {
		if d, err := time.Parse(time.RFC3339, req.RegistrationDeadline); err == nil {
			regDeadline = d
		}
	}

	event, err := h.service.CreateEvent(r.Context(), events.CreateEventRequest{
		TenantID:             req.TenantID,
		Title:                req.Title,
		Description:          req.Description,
		Category:             req.Category,
		VenueName:            req.VenueName,
		VenueCapacity:        req.VenueCapacity,
		StartTime:            start,
		EndTime:              end,
		RegistrationDeadline: regDeadline,
		OrganizerID:          req.OrganizerID,
		IsTicketed:           req.IsTicketed,
		TicketPrice:          req.TicketPrice,
		MaxTickets:           req.MaxTickets,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(event)
}

// ListEvents handles GET /api/v1/events
func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	var catPtr *events.EventCategory
	if c := r.URL.Query().Get("category"); c != "" {
		cat := events.EventCategory(c)
		catPtr = &cat
	}

	var statusPtr *events.EventStatus
	if s := r.URL.Query().Get("status"); s != "" {
		st := events.EventStatus(s)
		statusPtr = &st
	}

	eventList, err := h.service.ListEvents(r.Context(), tenantID, catPtr, statusPtr)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(eventList)
}

// GetEvent handles GET /api/v1/events/{id}
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || eventID == "" {
		problem.BadRequest(w, r, "tenant_id and event id are required", "MISSING_PARAMETERS")
		return
	}

	event, err := h.service.GetEvent(r.Context(), tenantID, eventID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(event)
}

// BookTicket handles POST /api/v1/events/{id}/tickets
func (h *Handler) BookTicket(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	var req struct {
		TenantID  string `json:"tenant_id"`
		StudentID string `json:"student_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.StudentID == "" {
		problem.BadRequest(w, r, "tenant_id and student_id are required", "VALIDATION_ERROR")
		return
	}

	ticket, err := h.service.BookTicket(r.Context(), events.BookTicketRequest{
		TenantID:  req.TenantID,
		EventID:   eventID,
		StudentID: req.StudentID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ticket)
}

// ListEventTickets handles GET /api/v1/events/{id}/tickets
func (h *Handler) ListEventTickets(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || eventID == "" {
		problem.BadRequest(w, r, "tenant_id and event id are required", "MISSING_PARAMETERS")
		return
	}

	tickets, err := h.service.ListEventTickets(r.Context(), tenantID, eventID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tickets)
}

// ListStudentTickets handles GET /api/v1/events/tickets/my
func (h *Handler) ListStudentTickets(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	studentID := r.URL.Query().Get("student_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || studentID == "" {
		problem.BadRequest(w, r, "tenant_id and student_id query parameters are required", "MISSING_PARAMETERS")
		return
	}

	tickets, err := h.service.ListStudentTickets(r.Context(), tenantID, studentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tickets)
}

// CheckInTicket handles POST /api/v1/events/tickets/checkin
func (h *Handler) CheckInTicket(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID   string `json:"tenant_id"`
		TicketCode string `json:"ticket_code"`
		StaffID    string `json:"staff_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.TicketCode == "" || req.StaffID == "" {
		problem.BadRequest(w, r, "tenant_id, ticket_code, and staff_id are required", "VALIDATION_ERROR")
		return
	}

	ticket, err := h.service.CheckInTicket(r.Context(), events.CheckInTicketRequest{
		TenantID:   req.TenantID,
		TicketCode: req.TicketCode,
		StaffID:    req.StaffID,
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
	case errors.Is(err, events.ErrEventNotFound),
		errors.Is(err, events.ErrTicketNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, events.ErrVenueConflict),
		errors.Is(err, events.ErrTicketAlreadyCheckedIn),
		errors.Is(err, events.ErrDuplicateTicketBooking),
		errors.Is(err, events.ErrEventCapacityExceeded):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	case errors.Is(err, events.ErrEventNotPublished),
		errors.Is(err, events.ErrRegistrationClosed),
		errors.Is(err, events.ErrTicketCancelled),
		errors.Is(err, events.ErrInvalidEventTimeRange):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	default:
		var domErr *events.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected events error occurred")
	}
}
