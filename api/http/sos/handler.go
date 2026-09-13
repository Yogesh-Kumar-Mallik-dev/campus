/**
 * BLOCK_API_SOS_HANDLER_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for emergency trigger, responder dispatch, and incident resolution.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package sos

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/sos"
)

type Handler struct {
	service sos.Service
}

func NewHandler(service sos.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/sos", func(r chi.Router) {
		r.Post("/trigger", h.TriggerSOS)
		r.Get("/incidents", h.ListIncidents)
		r.Get("/incidents/{id}", h.GetIncident)
		r.Post("/incidents/{id}/acknowledge", h.AcknowledgeIncident)
		r.Post("/incidents/{id}/dispatch", h.DispatchResponder)
		r.Post("/incidents/{id}/responders/{responderId}/status", h.UpdateResponderStatus)
		r.Post("/incidents/{id}/resolve", h.ResolveIncident)
		r.Get("/incidents/{id}/responders", h.ListResponders)
	})
}

// TriggerSOS handles POST /api/v1/sos/trigger
func (h *Handler) TriggerSOS(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID            string               `json:"tenant_id"`
		UserID              string               `json:"user_id"`
		EmergencyType       sos.SOSEmergencyType `json:"emergency_type"`
		Latitude            float64              `json:"latitude"`
		Longitude           float64              `json:"longitude"`
		LocationDescription string               `json:"location_description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.UserID == "" || req.LocationDescription == "" {
		problem.BadRequest(w, r, "tenant_id, user_id, and location_description are required", "VALIDATION_ERROR")
		return
	}

	incident, err := h.service.TriggerSOS(r.Context(), sos.TriggerSOSRequest{
		TenantID:            req.TenantID,
		UserID:              req.UserID,
		EmergencyType:       req.EmergencyType,
		Latitude:            req.Latitude,
		Longitude:           req.Longitude,
		LocationDescription: req.LocationDescription,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(incident)
}

// ListIncidents handles GET /api/v1/sos/incidents
func (h *Handler) ListIncidents(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	filter := sos.SOSIncidentFilter{
		TenantID: tenantID,
	}
	if st := r.URL.Query().Get("status"); st != "" {
		status := sos.SOSIncidentStatus(st)
		filter.Status = &status
	}
	if em := r.URL.Query().Get("emergency_type"); em != "" {
		emType := sos.SOSEmergencyType(em)
		filter.EmergencyType = &emType
	}
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		filter.UserID = &uid
	}

	incidents, err := h.service.ListIncidents(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incidents)
}

// GetIncident handles GET /api/v1/sos/incidents/{id}
func (h *Handler) GetIncident(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	incidentID := chi.URLParam(r, "id")

	if tenantID == "" || incidentID == "" {
		problem.BadRequest(w, r, "tenant_id and incident id are required", "VALIDATION_ERROR")
		return
	}

	incident, err := h.service.GetIncident(r.Context(), tenantID, incidentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incident)
}

// AcknowledgeIncident handles POST /api/v1/sos/incidents/{id}/acknowledge
func (h *Handler) AcknowledgeIncident(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "id")
	var req struct {
		TenantID string `json:"tenant_id"`
		ActorID  string `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || incidentID == "" || req.ActorID == "" {
		problem.BadRequest(w, r, "tenant_id, incident id, and actor_id are required", "VALIDATION_ERROR")
		return
	}

	incident, err := h.service.AcknowledgeIncident(r.Context(), req.TenantID, incidentID, req.ActorID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incident)
}

// DispatchResponder handles POST /api/v1/sos/incidents/{id}/dispatch
func (h *Handler) DispatchResponder(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "id")
	var req struct {
		TenantID    string               `json:"tenant_id"`
		ResponderID string               `json:"responder_id"`
		Role        sos.SOSResponderRole `json:"role"`
		ActorID     string               `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || incidentID == "" || req.ResponderID == "" {
		problem.BadRequest(w, r, "tenant_id, incident id, and responder_id are required", "VALIDATION_ERROR")
		return
	}

	responder, err := h.service.DispatchResponder(r.Context(), sos.DispatchResponderRequest{
		TenantID:    req.TenantID,
		IncidentID:  incidentID,
		ResponderID: req.ResponderID,
		Role:        req.Role,
		ActorID:     req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(responder)
}

// UpdateResponderStatus handles POST /api/v1/sos/incidents/{id}/responders/{responderId}/status
func (h *Handler) UpdateResponderStatus(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "id")
	responderID := chi.URLParam(r, "responderId")
	var req struct {
		TenantID string                 `json:"tenant_id"`
		Status   sos.SOSResponderStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || incidentID == "" || responderID == "" || req.Status == "" {
		problem.BadRequest(w, r, "tenant_id, incident id, responder_id, and status are required", "VALIDATION_ERROR")
		return
	}

	responder, err := h.service.UpdateResponderStatus(r.Context(), sos.UpdateResponderStatusRequest{
		TenantID:    req.TenantID,
		IncidentID:  incidentID,
		ResponderID: responderID,
		Status:      req.Status,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responder)
}

// ResolveIncident handles POST /api/v1/sos/incidents/{id}/resolve
func (h *Handler) ResolveIncident(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "id")
	var req struct {
		TenantID        string                `json:"tenant_id"`
		ResolvedByID    string                `json:"resolved_by_id"`
		Status          sos.SOSIncidentStatus `json:"status"`
		ResolutionNotes string                `json:"resolution_notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || incidentID == "" || req.ResolvedByID == "" {
		problem.BadRequest(w, r, "tenant_id, incident id, and resolved_by_id are required", "VALIDATION_ERROR")
		return
	}

	incident, err := h.service.ResolveIncident(r.Context(), sos.ResolveIncidentRequest{
		TenantID:        req.TenantID,
		IncidentID:      incidentID,
		ResolvedByID:    req.ResolvedByID,
		Status:          req.Status,
		ResolutionNotes: req.ResolutionNotes,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incident)
}

// ListResponders handles GET /api/v1/sos/incidents/{id}/responders
func (h *Handler) ListResponders(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	incidentID := chi.URLParam(r, "id")

	if tenantID == "" || incidentID == "" {
		problem.BadRequest(w, r, "tenant_id and incident id are required", "VALIDATION_ERROR")
		return
	}

	responders, err := h.service.ListResponders(r.Context(), tenantID, incidentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responders)
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, sos.ErrIncidentNotFound),
		errors.Is(err, sos.ErrResponderNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, sos.ErrResponderAlreadyDispatched):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	case errors.Is(err, sos.ErrInvalidCoordinates),
		errors.Is(err, sos.ErrInvalidIncidentTransition),
		errors.Is(err, sos.ErrIncidentAlreadyClosed):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	default:
		var domErr *sos.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected emergency response error occurred")
	}
}
