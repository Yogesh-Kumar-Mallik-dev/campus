/**
 * BLOCK_API_STUDYHUB_HANDLER_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for courses, materials, assignments, submissions, and peer grading.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package studyhub

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/studyhub"
)

type Handler struct {
	service studyhub.Service
}

func NewHandler(service studyhub.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/studyhub", func(r chi.Router) {
		// Courses
		r.Post("/courses", h.CreateCourse)
		r.Get("/courses", h.ListCourses)
		r.Get("/courses/{id}", h.GetCourse)

		// Materials
		r.Post("/courses/{id}/materials", h.PublishMaterial)
		r.Get("/courses/{id}/materials", h.ListMaterials)

		// Assignments
		r.Post("/courses/{id}/assignments", h.CreateAssignment)
		r.Get("/courses/{id}/assignments", h.ListAssignments)
		r.Get("/assignments/{id}", h.GetAssignment)

		// Submissions & Grading
		r.Post("/assignments/{id}/submissions", h.SubmitAssignment)
		r.Get("/assignments/{id}/submissions", h.ListSubmissions)
		r.Post("/submissions/{id}/grade", h.GradeSubmission)

		// Peer Reviews
		r.Post("/submissions/{id}/reviews", h.SubmitPeerReview)
		r.Get("/submissions/{id}/reviews", h.ListPeerReviews)
	})
}

// CreateCourse handles POST /api/v1/studyhub/courses
func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID     string  `json:"tenant_id"`
		Code         string  `json:"code"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		Department   string  `json:"department"`
		Credits      int     `json:"credits"`
		Semester     int     `json:"semester"`
		InstructorID *string `json:"instructor_id"`
		SyllabusText string  `json:"syllabus_text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.Code == "" || req.Title == "" || req.Department == "" {
		problem.BadRequest(w, r, "tenant_id, code, title, and department are required", "VALIDATION_ERROR")
		return
	}

	course, err := h.service.CreateCourse(r.Context(), studyhub.CreateCourseRequest{
		TenantID:     req.TenantID,
		Code:         req.Code,
		Title:        req.Title,
		Description:  req.Description,
		Department:   req.Department,
		Credits:      req.Credits,
		Semester:     req.Semester,
		InstructorID: req.InstructorID,
		SyllabusText: req.SyllabusText,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(course)
}

// ListCourses handles GET /api/v1/studyhub/courses
func (h *Handler) ListCourses(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	var deptPtr *string
	if dept := r.URL.Query().Get("department"); dept != "" {
		deptPtr = &dept
	}

	var semPtr *int
	if semStr := r.URL.Query().Get("semester"); semStr != "" {
		if sem, err := strconv.Atoi(semStr); err == nil {
			semPtr = &sem
		}
	}

	courses, err := h.service.ListCourses(r.Context(), tenantID, deptPtr, semPtr)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(courses)
}

// GetCourse handles GET /api/v1/studyhub/courses/{id}
func (h *Handler) GetCourse(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	courseID := chi.URLParam(r, "id")
	if tenantID == "" || courseID == "" {
		problem.BadRequest(w, r, "tenant_id and course id are required", "MISSING_PARAMETERS")
		return
	}

	course, err := h.service.GetCourse(r.Context(), tenantID, courseID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(course)
}

// PublishMaterial handles POST /api/v1/studyhub/courses/{id}/materials
func (h *Handler) PublishMaterial(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "id")
	var req struct {
		TenantID      string                `json:"tenant_id"`
		Title         string                `json:"title"`
		Description   string                `json:"description"`
		UnitNumber    int                   `json:"unit_number"`
		MaterialType  studyhub.MaterialType `json:"material_type"`
		FileURL       string                `json:"file_url"`
		FileSizeBytes int64                 `json:"file_size_bytes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.Title == "" || req.FileURL == "" {
		problem.BadRequest(w, r, "tenant_id, title, and file_url are required", "VALIDATION_ERROR")
		return
	}

	mat, err := h.service.PublishMaterial(r.Context(), studyhub.PublishMaterialRequest{
		TenantID:      req.TenantID,
		CourseID:      courseID,
		Title:         req.Title,
		Description:   req.Description,
		UnitNumber:    req.UnitNumber,
		MaterialType:  req.MaterialType,
		FileURL:       req.FileURL,
		FileSizeBytes: req.FileSizeBytes,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(mat)
}

// ListMaterials handles GET /api/v1/studyhub/courses/{id}/materials
func (h *Handler) ListMaterials(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || courseID == "" {
		problem.BadRequest(w, r, "tenant_id and course id are required", "MISSING_PARAMETERS")
		return
	}

	var unitPtr *int
	if unitStr := r.URL.Query().Get("unit_number"); unitStr != "" {
		if u, err := strconv.Atoi(unitStr); err == nil {
			unitPtr = &u
		}
	}

	materials, err := h.service.ListMaterials(r.Context(), tenantID, courseID, unitPtr)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(materials)
}

// CreateAssignment handles POST /api/v1/studyhub/courses/{id}/assignments
func (h *Handler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "id")
	var req struct {
		TenantID                 string `json:"tenant_id"`
		Title                    string `json:"title"`
		Description              string `json:"description"`
		MaxMarks                 int    `json:"max_marks"`
		DueDate                  string `json:"due_date"`
		AllowLateSubmission      bool   `json:"allow_late_submission"`
		LatePenaltyPercentPerDay int    `json:"late_penalty_percent_per_day"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.Title == "" || req.DueDate == "" {
		problem.BadRequest(w, r, "tenant_id, title, and due_date are required", "VALIDATION_ERROR")
		return
	}

	due, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		problem.BadRequest(w, r, "due_date must be formatted in RFC3339 (e.g. 2026-09-20T23:59:59Z)", "INVALID_DUE_DATE")
		return
	}

	asgn, err := h.service.CreateAssignment(r.Context(), studyhub.CreateAssignmentRequest{
		TenantID:                 req.TenantID,
		CourseID:                 courseID,
		Title:                    req.Title,
		Description:              req.Description,
		MaxMarks:                 req.MaxMarks,
		DueDate:                  due,
		AllowLateSubmission:      req.AllowLateSubmission,
		LatePenaltyPercentPerDay: req.LatePenaltyPercentPerDay,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(asgn)
}

// ListAssignments handles GET /api/v1/studyhub/courses/{id}/assignments
func (h *Handler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || courseID == "" {
		problem.BadRequest(w, r, "tenant_id and course id are required", "MISSING_PARAMETERS")
		return
	}

	assignments, err := h.service.ListAssignments(r.Context(), tenantID, courseID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(assignments)
}

// GetAssignment handles GET /api/v1/studyhub/assignments/{id}
func (h *Handler) GetAssignment(w http.ResponseWriter, r *http.Request) {
	assignmentID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || assignmentID == "" {
		problem.BadRequest(w, r, "tenant_id and assignment id are required", "MISSING_PARAMETERS")
		return
	}

	asgn, err := h.service.GetAssignment(r.Context(), tenantID, assignmentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(asgn)
}

// SubmitAssignment handles POST /api/v1/studyhub/assignments/{id}/submissions
func (h *Handler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	assignmentID := chi.URLParam(r, "id")
	var req struct {
		TenantID    string `json:"tenant_id"`
		StudentID   string `json:"student_id"`
		FileURL     string `json:"file_url"`
		ContentText string `json:"content_text"`
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

	sub, err := h.service.SubmitAssignment(r.Context(), studyhub.SubmitAssignmentRequest{
		TenantID:     req.TenantID,
		AssignmentID: assignmentID,
		StudentID:    req.StudentID,
		FileURL:      req.FileURL,
		ContentText:  req.ContentText,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sub)
}

// ListSubmissions handles GET /api/v1/studyhub/assignments/{id}/submissions
func (h *Handler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	assignmentID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" || assignmentID == "" {
		problem.BadRequest(w, r, "tenant_id and assignment id are required", "MISSING_PARAMETERS")
		return
	}

	submissions, err := h.service.ListSubmissions(r.Context(), tenantID, assignmentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(submissions)
}

// GradeSubmission handles POST /api/v1/studyhub/submissions/{id}/grade
func (h *Handler) GradeSubmission(w http.ResponseWriter, r *http.Request) {
	subID := chi.URLParam(r, "id")
	var req struct {
		TenantID   string  `json:"tenant_id"`
		RawMarks   float64 `json:"raw_marks"`
		Feedback   string  `json:"feedback"`
		GradedByID string  `json:"graded_by_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.GradedByID == "" {
		problem.BadRequest(w, r, "tenant_id and graded_by_id are required", "VALIDATION_ERROR")
		return
	}

	sub, err := h.service.GradeSubmission(r.Context(), studyhub.GradeSubmissionRequest{
		TenantID:   req.TenantID,
		ID:         subID,
		RawMarks:   req.RawMarks,
		Feedback:   req.Feedback,
		GradedByID: req.GradedByID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sub)
}

// SubmitPeerReview handles POST /api/v1/studyhub/submissions/{id}/reviews
func (h *Handler) SubmitPeerReview(w http.ResponseWriter, r *http.Request) {
	subID := chi.URLParam(r, "id")
	var req struct {
		TenantID          string  `json:"tenant_id"`
		ReviewerStudentID string  `json:"reviewer_student_id"`
		Score             float64 `json:"score"`
		Comments          string  `json:"comments"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.ReviewerStudentID == "" {
		problem.BadRequest(w, r, "tenant_id and reviewer_student_id are required", "VALIDATION_ERROR")
		return
	}

	review, err := h.service.SubmitPeerReview(r.Context(), studyhub.SubmitPeerReviewRequest{
		TenantID:          req.TenantID,
		SubmissionID:      subID,
		ReviewerStudentID: req.ReviewerStudentID,
		Score:             req.Score,
		Comments:          req.Comments,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(review)
}

// ListPeerReviews handles GET /api/v1/studyhub/submissions/{id}/reviews
func (h *Handler) ListPeerReviews(w http.ResponseWriter, r *http.Request) {
	subID := chi.URLParam(r, "id")
	if subID == "" {
		problem.BadRequest(w, r, "submission id is required", "MISSING_SUBMISSION_ID")
		return
	}

	reviews, err := h.service.ListPeerReviews(r.Context(), subID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, studyhub.ErrCourseNotFound),
		errors.Is(err, studyhub.ErrMaterialNotFound),
		errors.Is(err, studyhub.ErrAssignmentNotFound),
		errors.Is(err, studyhub.ErrSubmissionNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, studyhub.ErrCourseCodeExists),
		errors.Is(err, studyhub.ErrDuplicateSubmission),
		errors.Is(err, studyhub.ErrAssignmentClosed),
		errors.Is(err, studyhub.ErrDuplicatePeerReview):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	case errors.Is(err, studyhub.ErrLateSubmissionDisallowed),
		errors.Is(err, studyhub.ErrMarksExceedMaxMarks),
		errors.Is(err, studyhub.ErrSelfPeerReviewNotAllowed),
		errors.Is(err, studyhub.ErrInvalidScoreRange):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	default:
		var domErr *studyhub.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected study hub error occurred")
	}
}
