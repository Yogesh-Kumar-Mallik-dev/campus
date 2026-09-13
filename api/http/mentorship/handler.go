/**
 * BLOCK_API_MENTORSHIP_HANDLER_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for mentor allocations, counseling sessions, academic progress, and risk interventions.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package mentorship

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/mentorship"
)

type Handler struct {
	service mentorship.Service
}

func NewHandler(service mentorship.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/mentorship", func(r chi.Router) {
		// Allocations
		r.Post("/allocations", h.AllocateMentor)
		r.Get("/allocations", h.ListAllocations)
		r.Get("/allocations/active", h.GetActiveAllocation)

		// Sessions
		r.Post("/sessions", h.ScheduleSession)
		r.Get("/allocations/{id}/sessions", h.ListSessions)
		r.Post("/sessions/{id}/complete", h.CompleteSession)

		// Academic Progress
		r.Post("/progress", h.RecordProgress)
		r.Get("/progress", h.ListProgress)

		// At-Risk Alerts
		r.Get("/alerts", h.ListAlerts)
		r.Post("/alerts/{id}/resolve", h.ResolveAlert)
	})
}

// AllocateMentor handles POST /api/v1/mentorship/allocations
func (h *Handler) AllocateMentor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID      string  `json:"tenant_id"`
		StudentID     string  `json:"student_id"`
		MentorStaffID string  `json:"mentor_staff_id"`
		CohortID      *string `json:"cohort_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.StudentID == "" || req.MentorStaffID == "" {
		problem.BadRequest(w, r, "tenant_id, student_id, and mentor_staff_id are required", "VALIDATION_ERROR")
		return
	}

	alloc, err := h.service.AllocateMentor(r.Context(), mentorship.AllocateMentorRequest{
		TenantID:      req.TenantID,
		StudentID:     req.StudentID,
		MentorStaffID: req.MentorStaffID,
		CohortID:      req.CohortID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(alloc)
}

// ListAllocations handles GET /api/v1/mentorship/allocations
func (h *Handler) ListAllocations(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	var mentorPtr *string
	if m := r.URL.Query().Get("mentor_staff_id"); m != "" {
		mentorPtr = &m
	}

	var statusPtr *mentorship.AllocationStatus
	if s := r.URL.Query().Get("status"); s != "" {
		st := mentorship.AllocationStatus(s)
		statusPtr = &st
	}

	allocs, err := h.service.ListAllocations(r.Context(), tenantID, mentorPtr, statusPtr)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(allocs)
}

// GetActiveAllocation handles GET /api/v1/mentorship/allocations/active
func (h *Handler) GetActiveAllocation(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	studentID := r.URL.Query().Get("student_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || studentID == "" {
		problem.BadRequest(w, r, "tenant_id and student_id query parameters are required", "MISSING_PARAMETERS")
		return
	}

	alloc, err := h.service.GetActiveAllocation(r.Context(), tenantID, studentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(alloc)
}

// ScheduleSession handles POST /api/v1/mentorship/sessions
func (h *Handler) ScheduleSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID     string                 `json:"tenant_id"`
		AllocationID string                 `json:"allocation_id"`
		ScheduledAt  string                 `json:"scheduled_at"`
		Location     string                 `json:"location"`
		MeetingType  mentorship.MeetingType `json:"meeting_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.AllocationID == "" || req.ScheduledAt == "" {
		problem.BadRequest(w, r, "tenant_id, allocation_id, and scheduled_at are required", "VALIDATION_ERROR")
		return
	}

	sched, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		problem.BadRequest(w, r, "scheduled_at must be in RFC3339 format (e.g. 2026-09-20T15:00:00Z)", "INVALID_SCHEDULE_TIME")
		return
	}

	sess, err := h.service.ScheduleSession(r.Context(), mentorship.ScheduleSessionRequest{
		TenantID:     req.TenantID,
		AllocationID: req.AllocationID,
		ScheduledAt:  sched,
		Location:     req.Location,
		MeetingType:  req.MeetingType,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sess)
}

// ListSessions handles GET /api/v1/mentorship/allocations/{id}/sessions
func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	allocationID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || allocationID == "" {
		problem.BadRequest(w, r, "tenant_id and allocation id are required", "MISSING_PARAMETERS")
		return
	}

	sessions, err := h.service.ListSessions(r.Context(), tenantID, allocationID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sessions)
}

// CompleteSession handles POST /api/v1/mentorship/sessions/{id}/complete
func (h *Handler) CompleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var req struct {
		TenantID          string  `json:"tenant_id"`
		DiscussionSummary string  `json:"discussion_summary"`
		ActionItems       string  `json:"action_items"`
		FollowUpDate      *string `json:"follow_up_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	var followUpPtr *time.Time
	if req.FollowUpDate != nil && *req.FollowUpDate != "" {
		fUp, err := time.Parse(time.RFC3339, *req.FollowUpDate)
		if err == nil {
			followUpPtr = &fUp
		}
	}

	sess, err := h.service.CompleteSession(r.Context(), mentorship.CompleteSessionRequest{
		TenantID:          req.TenantID,
		ID:                sessionID,
		DiscussionSummary: req.DiscussionSummary,
		ActionItems:       req.ActionItems,
		FollowUpDate:      followUpPtr,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sess)
}

// RecordProgress handles POST /api/v1/mentorship/progress
func (h *Handler) RecordProgress(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID             string  `json:"tenant_id"`
		StudentID            string  `json:"student_id"`
		Semester             int     `json:"semester"`
		SGPA                 float64 `json:"sgpa"`
		CGPA                 float64 `json:"cgpa"`
		AttendancePercentage float64 `json:"attendance_percentage"`
		CreditsEarned        int     `json:"credits_earned"`
		TotalCredits         int     `json:"total_credits"`
		Remarks              string  `json:"remarks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.StudentID == "" || req.Semester <= 0 {
		problem.BadRequest(w, r, "tenant_id, student_id, and semester are required", "VALIDATION_ERROR")
		return
	}

	prog, err := h.service.RecordProgress(r.Context(), mentorship.RecordProgressRequest{
		TenantID:             req.TenantID,
		StudentID:            req.StudentID,
		Semester:             req.Semester,
		SGPA:                 req.SGPA,
		CGPA:                 req.CGPA,
		AttendancePercentage: req.AttendancePercentage,
		CreditsEarned:        req.CreditsEarned,
		TotalCredits:         req.TotalCredits,
		Remarks:              req.Remarks,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(prog)
}

// ListProgress handles GET /api/v1/mentorship/progress
func (h *Handler) ListProgress(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	studentID := r.URL.Query().Get("student_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || studentID == "" {
		problem.BadRequest(w, r, "tenant_id and student_id query parameters are required", "MISSING_PARAMETERS")
		return
	}

	progressList, err := h.service.ListProgress(r.Context(), tenantID, studentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(progressList)
}

// ListAlerts handles GET /api/v1/mentorship/alerts
func (h *Handler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	var studentPtr *string
	if s := r.URL.Query().Get("student_id"); s != "" {
		studentPtr = &s
	}

	var resolvedPtr *bool
	if resStr := r.URL.Query().Get("resolved"); resStr != "" {
		if b, err := strconv.ParseBool(resStr); err == nil {
			resolvedPtr = &b
		}
	}

	alerts, err := h.service.ListAlerts(r.Context(), tenantID, studentPtr, resolvedPtr)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(alerts)
}

// ResolveAlert handles POST /api/v1/mentorship/alerts/{id}/resolve
func (h *Handler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	alertID := chi.URLParam(r, "id")
	var req struct {
		TenantID     string `json:"tenant_id"`
		ResolvedByID string `json:"resolved_by_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.ResolvedByID == "" {
		problem.BadRequest(w, r, "tenant_id and resolved_by_id are required", "VALIDATION_ERROR")
		return
	}

	alert, err := h.service.ResolveAlert(r.Context(), mentorship.ResolveAlertRequest{
		TenantID:     req.TenantID,
		ID:           alertID,
		ResolvedByID: req.ResolvedByID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(alert)
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, mentorship.ErrAllocationNotFound),
		errors.Is(err, mentorship.ErrSessionNotFound),
		errors.Is(err, mentorship.ErrProgressRecordNotFound),
		errors.Is(err, mentorship.ErrAlertNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, mentorship.ErrActiveAllocationExists),
		errors.Is(err, mentorship.ErrProgressRecordExists),
		errors.Is(err, mentorship.ErrAlertAlreadyResolved),
		errors.Is(err, mentorship.ErrSessionAlreadyCompleted):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	case errors.Is(err, mentorship.ErrInvalidSessionTransition),
		errors.Is(err, mentorship.ErrInvalidScoreBounds):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	default:
		var domErr *mentorship.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected mentorship error occurred")
	}
}
