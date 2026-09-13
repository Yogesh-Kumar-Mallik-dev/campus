/**
 * BLOCK_API_ATTENDANCE_HANDLER_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for timetable, sessions, roll call marking, biometric sync, shortage alerts, and medical leave.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package attendance

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/attendance"
)

type Handler struct {
	service *attendance.Service
}

func NewHandler(service *attendance.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/attendance", func(r chi.Router) {
		// Subjects & Timetable
		r.Post("/subjects", h.CreateSubject)
		r.Post("/slots", h.CreateSlot)

		// Attendance Sessions
		r.Post("/sessions", h.ScheduleSession)
		r.Get("/sessions", h.ListSessions)
		r.Get("/sessions/{id}", h.GetSession)
		r.Post("/sessions/{id}/open", h.OpenSession)
		r.Post("/sessions/{id}/mark", h.MarkAttendance)
		r.Post("/sessions/{id}/biometric", h.IngestBiometric)
		r.Post("/sessions/{id}/lock", h.LockSession)
		r.Post("/sessions/{id}/finalize", h.FinalizeSession)

		// Summaries & Shortage Metrics
		r.Get("/summary/students/{student_id}", h.GetStudentSummary)

		// Medical Leave Condonations
		r.Post("/leaves", h.ApplyMedicalLeave)
		r.Post("/leaves/{id}/decision", h.DecideMedicalLeave)
	})
}

// CreateSubject handles POST /api/v1/attendance/subjects
func (h *Handler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	var cmd attendance.CreateSubjectCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}

	subject, err := h.service.CreateSubject(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v1/attendance/subjects/"+subject.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"subject": subject,
	})
}

// CreateSlot handles POST /api/v1/attendance/slots
func (h *Handler) CreateSlot(w http.ResponseWriter, r *http.Request) {
	var cmd attendance.CreateSlotCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}

	slot, err := h.service.CreateTimetableSlot(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"slot": slot,
	})
}

// ScheduleSession handles POST /api/v1/attendance/sessions
func (h *Handler) ScheduleSession(w http.ResponseWriter, r *http.Request) {
	var cmd attendance.ScheduleSessionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}

	session, err := h.service.ScheduleSession(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v1/attendance/sessions/"+session.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"session": session,
	})
}

// ListSessions handles GET /api/v1/attendance/sessions
func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	filter := attendance.SessionFilter{
		TenantID:    tenantID,
		SubjectID:   r.URL.Query().Get("subject_id"),
		CohortID:    r.URL.Query().Get("cohort_id"),
		FacultyID:   r.URL.Query().Get("faculty_id"),
		SessionDate: r.URL.Query().Get("session_date"),
		FromDate:    r.URL.Query().Get("from_date"),
		ToDate:      r.URL.Query().Get("to_date"),
		Limit:       limit,
		Offset:      offset,
	}

	if st := r.URL.Query().Get("status"); st != "" {
		status := attendance.AttendanceSessionStatus(strings.ToUpper(st))
		filter.Status = &status
	}

	sessions, total, err := h.service.ListSessions(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"sessions": sessions,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetSession handles GET /api/v1/attendance/sessions/{id}
func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	session, err := h.service.GetSession(r.Context(), tenantID, sessionID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"session": session,
	})
}

// OpenSession handles POST /api/v1/attendance/sessions/{id}/open
func (h *Handler) OpenSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var body struct {
		TenantID  string `json:"tenant_id"`
		FacultyID string `json:"faculty_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if body.TenantID == "" {
		body.TenantID = r.Header.Get("X-Tenant-ID")
	}

	session, err := h.service.OpenSession(r.Context(), body.TenantID, sessionID, body.FacultyID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"session": session,
	})
}

// MarkAttendance handles POST /api/v1/attendance/sessions/{id}/mark
func (h *Handler) MarkAttendance(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var cmd attendance.MarkAttendanceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	cmd.SessionID = sessionID
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	session, err := h.service.MarkAttendance(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"session": session,
	})
}

// IngestBiometric handles POST /api/v1/attendance/sessions/{id}/biometric
func (h *Handler) IngestBiometric(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var cmd attendance.IngestBiometricCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	cmd.SessionID = sessionID
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	session, err := h.service.IngestBiometricBatch(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"session": session,
	})
}

// LockSession handles POST /api/v1/attendance/sessions/{id}/lock
func (h *Handler) LockSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var body struct {
		TenantID  string `json:"tenant_id"`
		FacultyID string `json:"faculty_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.TenantID == "" {
		body.TenantID = r.Header.Get("X-Tenant-ID")
	}

	session, err := h.service.LockSession(r.Context(), body.TenantID, sessionID, body.FacultyID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"session": session,
	})
}

// FinalizeSession handles POST /api/v1/attendance/sessions/{id}/finalize
func (h *Handler) FinalizeSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var body struct {
		TenantID  string `json:"tenant_id"`
		FacultyID string `json:"faculty_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.TenantID == "" {
		body.TenantID = r.Header.Get("X-Tenant-ID")
	}

	session, err := h.service.FinalizeSession(r.Context(), body.TenantID, sessionID, body.FacultyID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"session": session,
	})
}

// GetStudentSummary handles GET /api/v1/attendance/summary/students/{student_id}
func (h *Handler) GetStudentSummary(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "student_id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}
	subjectID := r.URL.Query().Get("subject_id")

	summary, err := h.service.GetStudentAttendanceSummary(r.Context(), tenantID, studentID, subjectID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"summary": summary,
	})
}

// ApplyMedicalLeave handles POST /api/v1/attendance/leaves
func (h *Handler) ApplyMedicalLeave(w http.ResponseWriter, r *http.Request) {
	var cmd attendance.ApplyMedicalLeaveCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}

	leave, err := h.service.ApplyMedicalLeave(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v1/attendance/leaves/"+leave.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"leave": leave,
	})
}

// DecideMedicalLeave handles POST /api/v1/attendance/leaves/{id}/decision
func (h *Handler) DecideMedicalLeave(w http.ResponseWriter, r *http.Request) {
	leaveID := chi.URLParam(r, "id")
	var body struct {
		TenantID        string `json:"tenant_id"`
		ApproverID      string `json:"approver_id"`
		Approved        bool   `json:"approved"`
		RejectionReason string `json:"rejection_reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if body.TenantID == "" {
		body.TenantID = r.Header.Get("X-Tenant-ID")
	}

	leave, err := h.service.ApproveMedicalLeave(r.Context(), body.TenantID, leaveID, body.ApproverID, body.Approved, body.RejectionReason)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"leave": leave,
	})
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, attendance.ErrSubjectNotFound),
		errors.Is(err, attendance.ErrSlotNotFound),
		errors.Is(err, attendance.ErrSessionNotFound),
		errors.Is(err, attendance.ErrStudentNotFound),
		errors.Is(err, attendance.ErrMedicalLeaveNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")
	case errors.Is(err, attendance.ErrTenantRequired),
		errors.Is(err, attendance.ErrInvalidInput),
		errors.Is(err, attendance.ErrUnauthorizedMarker):
		problem.BadRequest(w, r, err.Error(), "INVALID_INPUT")
	case errors.Is(err, attendance.ErrSessionLocked),
		errors.Is(err, attendance.ErrSessionNotOpen),
		errors.Is(err, attendance.ErrInvalidStateTransition):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)
	case errors.Is(err, attendance.ErrDuplicateSession):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")
	default:
		var domErr *attendance.DomainError
		if errors.As(err, &domErr) {
			problem.UnprocessableEntity(w, r, domErr.Message, domErr.Code, nil)
			return
		}
		problem.InternalServerError(w, r, "an unexpected attendance error occurred")
	}
}
