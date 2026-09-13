/**
 * BLOCK_API_WHISTLEBLOWER_HANDLER_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for zero-knowledge grievances and anonymous 2-way investigation.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package whistleblower

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/whistleblower"
)

type Handler struct {
	service whistleblower.Service
}

func NewHandler(service whistleblower.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/whistleblower", func(r chi.Router) {
		// Public / Anonymous routes
		r.Post("/reports", h.SubmitReport)
		r.Post("/reports/lookup", h.LookupReportByToken)
		r.Post("/reports/{id}/messages", h.AddMessage)

		// Committee / Authority routes
		r.Get("/committee/reports", h.ListReports)
		r.Post("/committee/reports/{id}/assign", h.AssignInvestigator)
		r.Post("/committee/reports/{id}/status", h.UpdateReportStatus)
	})
}

// SubmitReport handles POST /api/v1/whistleblower/reports
func (h *Handler) SubmitReport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID             string                                `json:"tenant_id"`
		Category             whistleblower.WhistleblowerCategory  `json:"category"`
		Severity             whistleblower.WhistleblowerSeverity  `json:"severity"`
		Title                string                                `json:"title"`
		DescriptionEncrypted string                                `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.Title == "" || req.DescriptionEncrypted == "" {
		problem.BadRequest(w, r, "tenant_id, title, and description are required", "VALIDATION_ERROR")
		return
	}

	report, rawToken, err := h.service.SubmitReport(r.Context(), whistleblower.SubmitReportRequest{
		TenantID:             req.TenantID,
		Category:             req.Category,
		Severity:             req.Severity,
		Title:                req.Title,
		DescriptionEncrypted: req.DescriptionEncrypted,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"report":               report,
		"secret_tracking_token": rawToken,
		"warning":              "Store this secret tracking token safely. It is your only access key to view status updates and communicate anonymously.",
	})
}

// LookupReportByToken handles POST /api/v1/whistleblower/reports/lookup
func (h *Handler) LookupReportByToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID         string `json:"tenant_id"`
		RawTrackingToken string `json:"raw_tracking_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.RawTrackingToken == "" {
		problem.BadRequest(w, r, "tenant_id and raw_tracking_token are required", "VALIDATION_ERROR")
		return
	}

	report, messages, err := h.service.LookupReportByToken(r.Context(), req.TenantID, req.RawTrackingToken)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"report":   report,
		"messages": messages,
	})
}

// AddMessage handles POST /api/v1/whistleblower/reports/{id}/messages
func (h *Handler) AddMessage(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "id")
	var req struct {
		TenantID         string                                  `json:"tenant_id"`
		SenderType       whistleblower.WhistleblowerSenderType   `json:"sender_type"`
		MessageEncrypted string                                  `json:"message"`
		RawTrackingToken string                                  `json:"raw_tracking_token,omitempty"`
		ActorID          string                                  `json:"actor_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || reportID == "" || req.MessageEncrypted == "" {
		problem.BadRequest(w, r, "tenant_id, report id, and message are required", "VALIDATION_ERROR")
		return
	}

	msg, err := h.service.AddMessage(r.Context(), whistleblower.AddMessageRequest{
		TenantID:         req.TenantID,
		ReportID:         reportID,
		SenderType:       req.SenderType,
		MessageEncrypted: req.MessageEncrypted,
		RawTrackingToken: req.RawTrackingToken,
		ActorID:          req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(msg)
}

// ListReports handles GET /api/v1/whistleblower/committee/reports
func (h *Handler) ListReports(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	filter := whistleblower.WhistleblowerReportFilter{
		TenantID: tenantID,
	}
	if cat := r.URL.Query().Get("category"); cat != "" {
		category := whistleblower.WhistleblowerCategory(cat)
		filter.Category = &category
	}
	if sev := r.URL.Query().Get("severity"); sev != "" {
		severity := whistleblower.WhistleblowerSeverity(sev)
		filter.Severity = &severity
	}
	if st := r.URL.Query().Get("status"); st != "" {
		status := whistleblower.WhistleblowerStatus(st)
		filter.Status = &status
	}

	reports, err := h.service.ListReports(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reports)
}

// AssignInvestigator handles POST /api/v1/whistleblower/committee/reports/{id}/assign
func (h *Handler) AssignInvestigator(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "id")
	var req struct {
		TenantID  string `json:"tenant_id"`
		OfficerID string `json:"officer_id"`
		ActorID   string `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || reportID == "" || req.OfficerID == "" {
		problem.BadRequest(w, r, "tenant_id, report id, and officer_id are required", "VALIDATION_ERROR")
		return
	}

	report, err := h.service.AssignInvestigator(r.Context(), req.TenantID, reportID, req.OfficerID, req.ActorID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(report)
}

// UpdateReportStatus handles POST /api/v1/whistleblower/committee/reports/{id}/status
func (h *Handler) UpdateReportStatus(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "id")
	var req struct {
		TenantID        string                              `json:"tenant_id"`
		Status          whistleblower.WhistleblowerStatus   `json:"status"`
		FindingsSummary string                              `json:"findings_summary"`
		ActorID         string                              `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || reportID == "" || req.Status == "" {
		problem.BadRequest(w, r, "tenant_id, report id, and status are required", "VALIDATION_ERROR")
		return
	}

	report, err := h.service.UpdateReportStatus(r.Context(), whistleblower.UpdateReportStatusRequest{
		TenantID:        req.TenantID,
		ReportID:        reportID,
		Status:          req.Status,
		FindingsSummary: req.FindingsSummary,
		ActorID:         req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(report)
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, whistleblower.ErrReportNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, whistleblower.ErrUnauthorizedCommitteeAccess):
		problem.Forbidden(w, r, err.Error(), "FORBIDDEN")

	case errors.Is(err, whistleblower.ErrInvalidTrackingToken),
		errors.Is(err, whistleblower.ErrInvalidReportTransition),
		errors.Is(err, whistleblower.ErrReportAlreadyClosed):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	default:
		var domErr *whistleblower.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected whistleblower error occurred")
	}
}
