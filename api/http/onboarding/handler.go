/**
 * BLOCK_API_ONBOARDING_HANDLER_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   HTTP transport handlers implementing RFC 7807 problem details, KYC flows, and cohort provisioning.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/onboarding"
)

type Handler struct {
	service *onboarding.Service
}

func NewHandler(service *onboarding.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/onboarding", func(r chi.Router) {
		r.Get("/departments", h.ListDepartments)
		r.Get("/programs", h.ListPrograms)
		r.Get("/cohorts", h.ListCohorts)

		r.Post("/applicants", h.CreateDraft)
		r.Get("/applicants", h.ListApplicants)
		r.Get("/applicants/{id}", h.GetApplicant)
		r.Patch("/applicants/{id}", h.UpdateDraft)

		r.Post("/applicants/{id}/documents", h.AttachDocument)
		r.Patch("/applicants/{id}/documents/{doc_id}", h.VerifyDocument)

		r.Post("/applicants/{id}/submit", h.SubmitApplication)
		r.Post("/applicants/{id}/review", h.AssignReviewer)
		r.Post("/applicants/{id}/verify", h.VerifyApplication)
		r.Post("/applicants/{id}/reject", h.RejectApplication)
		r.Post("/applicants/{id}/enroll", h.EnrollStudent)
		r.Post("/applicants/{id}/provision", h.ProvisionStaff)
	})
}

// CreateDraft handles POST /api/v1/onboarding/applicants
func (h *Handler) CreateDraft(w http.ResponseWriter, r *http.Request) {
	var cmd onboarding.CreateApplicantCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}

	app, err := h.service.CreateDraft(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v1/onboarding/applicants/"+app.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"applicant": app,
	})
}

// UpdateDraft handles PATCH /api/v1/onboarding/applicants/{id}
func (h *Handler) UpdateDraft(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	var cmd onboarding.UpdateApplicantCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}

	app, err := h.service.UpdateDraft(r.Context(), tenantID, applicantID, cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"applicant": app,
	})
}

// GetApplicant handles GET /api/v1/onboarding/applicants/{id}
func (h *Handler) GetApplicant(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	app, err := h.service.GetApplicant(r.Context(), tenantID, applicantID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"applicant":           app,
		"mandatory_documents": app.MandatoryDocuments(),
	})
}

// ListApplicants handles GET /api/v1/onboarding/applicants
func (h *Handler) ListApplicants(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	tenantID := q.Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	filter := onboarding.ApplicantFilter{
		TenantID:     tenantID,
		ProgramID:    q.Get("program_id"),
		DepartmentID: q.Get("department_id"),
		AcademicYear: q.Get("academic_year"),
		Search:       q.Get("search"),
		Limit:        50,
		Offset:       0,
	}

	if t := q.Get("type"); t != "" {
		ot := onboarding.OnboardingType(strings.ToUpper(t))
		filter.Type = &ot
	}
	if s := q.Get("status"); s != "" {
		os := onboarding.OnboardingStatus(strings.ToUpper(s))
		filter.Status = &os
	}

	if limitStr := q.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			filter.Limit = l
		}
	}
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	items, total, err := h.service.ListApplicants(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"applicants": items,
		"total":      total,
		"limit":      filter.Limit,
		"offset":     filter.Offset,
	})
}

// AttachDocument handles POST /api/v1/onboarding/applicants/{id}/documents
func (h *Handler) AttachDocument(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")

	var cmd onboarding.AttachDocumentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	cmd.ApplicantID = applicantID
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	doc, err := h.service.AttachDocument(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"document": doc,
	})
}

type verifyDocRequest struct {
	TenantID        string `json:"tenant_id"`
	Approved        bool   `json:"approved"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	VerifierID      string `json:"verifier_id"`
}

// VerifyDocument handles PATCH /api/v1/onboarding/applicants/{id}/documents/{doc_id}
func (h *Handler) VerifyDocument(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")
	docID := chi.URLParam(r, "doc_id")

	var req verifyDocRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	doc, err := h.service.VerifyDocument(r.Context(), req.TenantID, applicantID, docID, req.VerifierID, req.Approved, req.RejectionReason)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"document": doc,
	})
}

// SubmitApplication handles POST /api/v1/onboarding/applicants/{id}/submit
func (h *Handler) SubmitApplication(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}

	app, err := h.service.SubmitApplication(r.Context(), tenantID, applicantID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"applicant": app,
		"status":    app.Status,
	})
}

type assignReviewerRequest struct {
	TenantID   string `json:"tenant_id"`
	ReviewerID string `json:"reviewer_id"`
}

// AssignReviewer handles POST /api/v1/onboarding/applicants/{id}/review
func (h *Handler) AssignReviewer(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")

	var req assignReviewerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	app, err := h.service.AssignReviewer(r.Context(), req.TenantID, applicantID, req.ReviewerID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"applicant": app,
		"status":    app.Status,
	})
}

// VerifyApplication handles POST /api/v1/onboarding/applicants/{id}/verify
func (h *Handler) VerifyApplication(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}

	app, err := h.service.VerifyApplication(r.Context(), tenantID, applicantID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"applicant": app,
		"status":    app.Status,
	})
}

type rejectAppRequest struct {
	TenantID   string `json:"tenant_id"`
	ReviewerID string `json:"reviewer_id"`
	Reason     string `json:"reason"`
}

// RejectApplication handles POST /api/v1/onboarding/applicants/{id}/reject
func (h *Handler) RejectApplication(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")

	var req rejectAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	app, err := h.service.RejectApplication(r.Context(), req.TenantID, applicantID, req.ReviewerID, req.Reason)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"applicant": app,
		"status":    app.Status,
	})
}

// EnrollStudent handles POST /api/v1/onboarding/applicants/{id}/enroll
func (h *Handler) EnrollStudent(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")

	var cmd onboarding.EnrollStudentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	cmd.ApplicantID = applicantID
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	profile, err := h.service.EnrollStudent(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"student_profile": profile,
		"roll_number":     profile.RollNumber,
	})
}

// ProvisionStaff handles POST /api/v1/onboarding/applicants/{id}/provision
func (h *Handler) ProvisionStaff(w http.ResponseWriter, r *http.Request) {
	applicantID := chi.URLParam(r, "id")

	var cmd onboarding.ProvisionStaffCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	cmd.ApplicantID = applicantID
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	profile, err := h.service.ProvisionStaff(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"staff_profile": profile,
		"employee_id":   profile.EmployeeID,
	})
}

// ListDepartments handles GET /api/v1/onboarding/departments
func (h *Handler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	depts, err := h.service.ListDepartments(r.Context(), tenantID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"departments": depts,
	})
}

// ListPrograms handles GET /api/v1/onboarding/programs
func (h *Handler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}
	deptID := r.URL.Query().Get("department_id")

	progs, err := h.service.ListPrograms(r.Context(), tenantID, deptID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"programs": progs,
	})
}

// ListCohorts handles GET /api/v1/onboarding/cohorts
func (h *Handler) ListCohorts(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}
	progID := r.URL.Query().Get("program_id")
	academicYear := r.URL.Query().Get("academic_year")

	cohorts, err := h.service.ListCohorts(r.Context(), tenantID, progID, academicYear)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"cohorts": cohorts,
	})
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, onboarding.ErrApplicantNotFound),
		errors.Is(err, onboarding.ErrDocumentNotFound),
		errors.Is(err, onboarding.ErrDepartmentNotFound),
		errors.Is(err, onboarding.ErrProgramNotFound),
		errors.Is(err, onboarding.ErrCohortNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")
	case errors.Is(err, onboarding.ErrTenantRequired),
		errors.Is(err, onboarding.ErrInvalidInput),
		errors.Is(err, onboarding.ErrRejectionReasonRequired),
		errors.Is(err, onboarding.ErrUnauthorizedReviewer):
		problem.BadRequest(w, r, err.Error(), "INVALID_INPUT")
	case errors.Is(err, onboarding.ErrMissingMandatoryDocs),
		errors.Is(err, onboarding.ErrUnverifiedDocuments),
		errors.Is(err, onboarding.ErrInvalidStateTransition),
		errors.Is(err, onboarding.ErrCohortFull),
		errors.Is(err, onboarding.ErrAlreadyEnrolled):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)
	case errors.Is(err, onboarding.ErrDuplicateEmail),
		errors.Is(err, onboarding.ErrDuplicateRollNumber),
		errors.Is(err, onboarding.ErrDuplicateEmployeeID):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")
	default:
		var domErr *onboarding.DomainError
		if errors.As(err, &domErr) {
			problem.UnprocessableEntity(w, r, domErr.Message, domErr.Code, nil)
			return
		}
		problem.InternalServerError(w, r, "an unexpected onboarding error occurred")
	}
}
