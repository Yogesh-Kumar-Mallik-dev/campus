/**
 * BLOCK_API_AUDIT_HANDLER_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   HTTP transport handlers implementing RFC 7807 error semantics, pagination, and chain verification.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/audit"
)

// Handler handles HTTP requests for audit trail exploration, verification, and compliance reporting.
type Handler struct {
	service audit.Service
}

// NewHandler constructs an HTTP handler for the audit subsystem.
func NewHandler(service audit.Service) *Handler {
	return &Handler{service: service}
}

// QueryLogs handles GET /api/v1/audit/logs
func (h *Handler) QueryLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()
	tenantID := q.Get("tenant_id")
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	filter := audit.AuditFilter{
		TenantID: tenantID,
		Limit:    50,
		Offset:   0,
	}

	if actorID := q.Get("actor_id"); actorID != "" {
		filter.ActorID = &actorID
	}
	if actorType := q.Get("actor_type"); actorType != "" {
		at := audit.ActorType(actorType)
		filter.ActorType = &at
	}
	if action := q.Get("action"); action != "" {
		filter.Action = &action
	}
	if resType := q.Get("resource_type"); resType != "" {
		filter.ResourceType = &resType
	}
	if resID := q.Get("resource_id"); resID != "" {
		filter.ResourceID = &resID
	}
	if status := q.Get("status"); status != "" {
		st := audit.Status(status)
		filter.Status = &st
	}
	if reqTrace := q.Get("trace_id"); reqTrace != "" {
		filter.TraceID = &reqTrace
	}

	if fromStr := q.Get("from"); fromStr != "" {
		if parsed, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filter.FromTime = &parsed
		}
	}
	if toStr := q.Get("to"); toStr != "" {
		if parsed, err := time.Parse(time.RFC3339, toStr); err == nil {
			filter.ToTime = &parsed
		}
	}

	if limitStr := q.Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			filter.Limit = val
		}
	}
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			filter.Offset = val
		}
	}

	logs, total, err := h.service.QueryLogs(ctx, filter)
	if err != nil {
		problem.InternalServerError(w, r, "Failed to query audit logs")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"data":   logs,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetLogByID handles GET /api/v1/audit/logs/{id}
func (h *Handler) GetLogByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	logEntry, err := h.service.GetLogByID(ctx, tenantID, id)
	if err != nil {
		if err == audit.ErrAuditNotFound {
			problem.NotFound(w, r, "Audit log entry not found", "AUDIT_NOT_FOUND")
			return
		}
		problem.InternalServerError(w, r, "Failed to fetch audit log")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": logEntry,
	})
}

// VerifyIntegrity handles POST /api/v1/audit/verify
func (h *Handler) VerifyIntegrity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload struct {
		TenantID string     `json:"tenant_id"`
		FromTime *time.Time `json:"from_time,omitempty"`
		ToTime   *time.Time `json:"to_time,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		problem.BadRequest(w, r, "Malformed verification payload", "INVALID_REQUEST_PAYLOAD")
		return
	}

	if payload.TenantID == "" {
		problem.BadRequest(w, r, "tenant_id is required", "MISSING_TENANT_ID")
		return
	}

	result, err := h.service.VerifyChainIntegrity(ctx, payload.TenantID, payload.FromTime, payload.ToTime)
	if err != nil {
		problem.InternalServerError(w, r, "Failed to verify audit hash chain")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": result,
	})
}

// CreateCheckpoint handles POST /api/v1/audit/checkpoints
func (h *Handler) CreateCheckpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload struct {
		TenantID string `json:"tenant_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		problem.BadRequest(w, r, "Malformed checkpoint payload", "INVALID_REQUEST_PAYLOAD")
		return
	}

	if payload.TenantID == "" {
		problem.BadRequest(w, r, "tenant_id is required", "MISSING_TENANT_ID")
		return
	}

	checkpoint, err := h.service.CreateVerificationCheckpoint(ctx, payload.TenantID)
	if err != nil {
		if err == audit.ErrAuditNotFound {
			problem.NotFound(w, r, "No audit records found to anchor", "NO_AUDIT_RECORDS")
			return
		}
		problem.InternalServerError(w, r, "Failed to create verification checkpoint")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": checkpoint,
	})
}

// GenerateComplianceReport handles GET /api/v1/audit/compliance-reports
func (h *Handler) GenerateComplianceReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()
	tenantID := q.Get("tenant_id")
	framework := q.Get("framework")
	if tenantID == "" || framework == "" {
		problem.BadRequest(w, r, "tenant_id and framework query parameters are required", "MISSING_REQUIRED_PARAMS")
		return
	}

	fromTime := time.Now().Add(-30 * 24 * time.Hour) // Default last 30 days
	toTime := time.Now()

	if f := q.Get("from"); f != "" {
		if parsed, err := time.Parse(time.RFC3339, f); err == nil {
			fromTime = parsed
		}
	}
	if t := q.Get("to"); t != "" {
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			toTime = parsed
		}
	}

	report, err := h.service.GenerateComplianceReport(ctx, audit.ComplianceReportRequest{
		TenantID:  tenantID,
		Framework: framework,
		FromTime:  fromTime,
		ToTime:    toTime,
	})
	if err != nil {
		problem.InternalServerError(w, r, "Failed to compile compliance report")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": report,
	})
}
