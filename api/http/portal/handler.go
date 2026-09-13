/**
 * BLOCK_API_PORTAL_HANDLER_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for public landing, program offerings, and admissions prospect leads.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package portal

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/portal"
)

type Handler struct {
	service portal.Service
}

func NewHandler(service portal.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/portal", func(r chi.Router) {
		// Landing Config
		r.Get("/landing", h.GetLandingPage)
		r.Put("/landing", h.UpdateLandingPage)

		// Program Catalog
		r.Get("/programs", h.ListPrograms)
		r.Post("/programs", h.CreateProgram)
		r.Get("/programs/{id}", h.GetProgram)

		// Prospect Inquiries
		r.Post("/inquiries", h.SubmitInquiry)
		r.Get("/inquiries", h.ListInquiries)
		r.Post("/inquiries/{id}/status", h.UpdateInquiryStatus)
	})
}

// GetLandingPage handles GET /api/v1/portal/landing
func (h *Handler) GetLandingPage(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	page, err := h.service.GetLandingPage(r.Context(), tenantID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(page)
}

// UpdateLandingPage handles PUT /api/v1/portal/landing
func (h *Handler) UpdateLandingPage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID        string  `json:"tenant_id"`
		HeroHeadline    string  `json:"hero_headline"`
		HeroSubheadline string  `json:"hero_subheadline"`
		AdmissionsOpen  bool    `json:"admissions_open"`
		AdmissionsCycle string  `json:"admissions_cycle"`
		HeroImageURL    *string `json:"hero_image_url"`
		ContactEmail    string  `json:"contact_email"`
		ContactPhone    string  `json:"contact_phone"`
		CampusAddress   string  `json:"campus_address"`
		IsPublished     bool    `json:"is_published"`
		ActorID         string  `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.HeroHeadline == "" || req.ContactEmail == "" {
		problem.BadRequest(w, r, "tenant_id, hero_headline, and contact_email are required", "VALIDATION_ERROR")
		return
	}

	page, err := h.service.UpdateLandingPage(r.Context(), portal.UpdateLandingPageRequest{
		TenantID:        req.TenantID,
		HeroHeadline:    req.HeroHeadline,
		HeroSubheadline: req.HeroSubheadline,
		AdmissionsOpen:  req.AdmissionsOpen,
		AdmissionsCycle: req.AdmissionsCycle,
		HeroImageURL:    req.HeroImageURL,
		ContactEmail:    req.ContactEmail,
		ContactPhone:    req.ContactPhone,
		CampusAddress:   req.CampusAddress,
		IsPublished:     req.IsPublished,
		ActorID:         req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(page)
}

// ListPrograms handles GET /api/v1/portal/programs
func (h *Handler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	filter := portal.ProgramCatalogFilter{
		TenantID: tenantID,
	}
	if dt := r.URL.Query().Get("degree_type"); dt != "" {
		deg := portal.PortalDegreeType(dt)
		filter.DegreeType = &deg
	}
	if f := r.URL.Query().Get("featured"); f != "" {
		feat := f == "true"
		filter.Featured = &feat
	}

	programs, err := h.service.ListPrograms(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(programs)
}

// CreateProgram handles POST /api/v1/portal/programs
func (h *Handler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID            string                  `json:"tenant_id"`
		ProgramCode         string                  `json:"program_code"`
		ProgramName         string                  `json:"program_name"`
		DegreeType          portal.PortalDegreeType `json:"degree_type"`
		DepartmentName      string                  `json:"department_name"`
		DurationYears       int                     `json:"duration_years"`
		TotalSemesters      int                     `json:"total_semesters"`
		EligibilityCriteria string                  `json:"eligibility_criteria"`
		AnnualFee           float64                 `json:"annual_fee"`
		IsFeatured          bool                    `json:"is_featured"`
		ActorID             string                  `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.ProgramCode == "" || req.ProgramName == "" {
		problem.BadRequest(w, r, "tenant_id, program_code, and program_name are required", "VALIDATION_ERROR")
		return
	}

	prog, err := h.service.CreateProgram(r.Context(), portal.CreateProgramRequest{
		TenantID:            req.TenantID,
		ProgramCode:         req.ProgramCode,
		ProgramName:         req.ProgramName,
		DegreeType:          req.DegreeType,
		DepartmentName:      req.DepartmentName,
		DurationYears:       req.DurationYears,
		TotalSemesters:      req.TotalSemesters,
		EligibilityCriteria: req.EligibilityCriteria,
		AnnualFee:           req.AnnualFee,
		IsFeatured:          req.IsFeatured,
		ActorID:             req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(prog)
}

// GetProgram handles GET /api/v1/portal/programs/{id}
func (h *Handler) GetProgram(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	id := chi.URLParam(r, "id")

	if tenantID == "" || id == "" {
		problem.BadRequest(w, r, "tenant_id and id are required", "VALIDATION_ERROR")
		return
	}

	prog, err := h.service.GetProgram(r.Context(), tenantID, id)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(prog)
}

// SubmitInquiry handles POST /api/v1/portal/inquiries
func (h *Handler) SubmitInquiry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID          string `json:"tenant_id"`
		ProspectName      string `json:"prospect_name"`
		ProspectEmail     string `json:"prospect_email"`
		ProspectPhone     string `json:"prospect_phone"`
		ProgramOfInterest string `json:"program_of_interest"`
		Message           string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || req.ProspectName == "" || req.ProspectEmail == "" {
		problem.BadRequest(w, r, "tenant_id, prospect_name, and prospect_email are required", "VALIDATION_ERROR")
		return
	}

	inq, err := h.service.SubmitInquiry(r.Context(), portal.SubmitInquiryRequest{
		TenantID:          req.TenantID,
		ProspectName:      req.ProspectName,
		ProspectEmail:     req.ProspectEmail,
		ProspectPhone:     req.ProspectPhone,
		ProgramOfInterest: req.ProgramOfInterest,
		Message:           req.Message,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(inq)
}

// ListInquiries handles GET /api/v1/portal/inquiries
func (h *Handler) ListInquiries(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	filter := portal.PublicInquiryFilter{
		TenantID: tenantID,
	}
	if st := r.URL.Query().Get("status"); st != "" {
		status := portal.PortalInquiryStatus(st)
		filter.Status = &status
	}
	if cid := r.URL.Query().Get("counselor_id"); cid != "" {
		filter.AssignedCounselorID = &cid
	}

	inquiries, err := h.service.ListInquiries(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(inquiries)
}

// UpdateInquiryStatus handles POST /api/v1/portal/inquiries/{id}/status
func (h *Handler) UpdateInquiryStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		TenantID            string                     `json:"tenant_id"`
		Status              portal.PortalInquiryStatus `json:"status"`
		AssignedCounselorID *string                    `json:"assigned_counselor_id"`
		Notes               *string                    `json:"notes"`
		ActorID             string                     `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	if req.TenantID == "" || id == "" || req.Status == "" {
		problem.BadRequest(w, r, "tenant_id, id, and status are required", "VALIDATION_ERROR")
		return
	}

	inq, err := h.service.UpdateInquiryStatus(r.Context(), portal.UpdateInquiryStatusRequest{
		TenantID:            req.TenantID,
		InquiryID:           id,
		Status:              req.Status,
		AssignedCounselorID: req.AssignedCounselorID,
		Notes:               req.Notes,
		ActorID:             req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(inq)
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, portal.ErrLandingPageNotFound),
		errors.Is(err, portal.ErrProgramNotFound),
		errors.Is(err, portal.ErrInquiryNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, portal.ErrDuplicateProgramCode):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	default:
		var domErr *portal.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected portal error occurred")
	}
}
