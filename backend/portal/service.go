/**
 * BLOCK_PORTAL_SERVICE_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   Business logic orchestration for public landing showcase, program offerings, and prospect intake.
 */

package portal

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"campus/backend/audit"
)

type Service interface {
	GetLandingPage(ctx context.Context, tenantID string) (*PortalLandingPage, error)
	UpdateLandingPage(ctx context.Context, req UpdateLandingPageRequest) (*PortalLandingPage, error)

	CreateProgram(ctx context.Context, req CreateProgramRequest) (*PortalProgramCatalog, error)
	ListPrograms(ctx context.Context, filter ProgramCatalogFilter) ([]PortalProgramCatalog, error)
	GetProgram(ctx context.Context, tenantID, id string) (*PortalProgramCatalog, error)

	SubmitInquiry(ctx context.Context, req SubmitInquiryRequest) (*PortalPublicInquiry, error)
	ListInquiries(ctx context.Context, filter PublicInquiryFilter) ([]PortalPublicInquiry, error)
	UpdateInquiryStatus(ctx context.Context, req UpdateInquiryStatusRequest) (*PortalPublicInquiry, error)
}

type service struct {
	repo     Repository
	auditSub audit.Subscriber
}

func NewService(repo Repository, auditSub audit.Subscriber) Service {
	return &service{
		repo:     repo,
		auditSub: auditSub,
	}
}

func (s *service) logAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
	if s.auditSub == nil {
		return
	}
	var metaRaw json.RawMessage
	if meta != nil {
		if b, err := json.Marshal(meta); err == nil {
			metaRaw = b
		}
	}
	_ = s.auditSub.Enqueue(audit.RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      &actorID,
		ActorType:    audit.ActorType(actorType),
		Action:       action,
		ResourceType: resType,
		ResourceID:   &resID,
		Status:       status,
		Metadata:     metaRaw,
	})
}

type UpdateLandingPageRequest struct {
	TenantID        string  `json:"tenantId"`
	HeroHeadline    string  `json:"heroHeadline"`
	HeroSubheadline string  `json:"heroSubheadline"`
	AdmissionsOpen  bool    `json:"admissionsOpen"`
	AdmissionsCycle string  `json:"admissionsCycle"`
	HeroImageURL    *string `json:"heroImageUrl,omitempty"`
	ContactEmail    string  `json:"contactEmail"`
	ContactPhone    string  `json:"contactPhone"`
	CampusAddress   string  `json:"campusAddress"`
	IsPublished     bool    `json:"isPublished"`
	ActorID         string  `json:"actorId"`
}

func (s *service) GetLandingPage(ctx context.Context, tenantID string) (*PortalLandingPage, error) {
	page, err := s.repo.GetLandingPage(ctx, tenantID)
	if err != nil {
		// Return sane default for initial tenant render
		now := time.Now().UTC()
		return &PortalLandingPage{
			ID:              uuid.New().String(),
			TenantID:        tenantID,
			HeroHeadline:    "Welcome to Premier Global Campus",
			HeroSubheadline: "Empowering Next-Gen Leaders, Innovators, and Engineers",
			AdmissionsOpen:  true,
			AdmissionsCycle: "2026-2027",
			ContactEmail:    "admissions@campus.edu",
			ContactPhone:    "+91 80 1234 5678",
			CampusAddress:   "Academic City, Bangalore - 560100, India",
			IsPublished:     true,
			CreatedAt:       now,
			UpdatedAt:       now,
		}, nil
	}
	return page, nil
}

func (s *service) UpdateLandingPage(ctx context.Context, req UpdateLandingPageRequest) (*PortalLandingPage, error) {
	now := time.Now().UTC()
	page := &PortalLandingPage{
		ID:              uuid.New().String(),
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
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.UpsertLandingPage(ctx, page); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, req.ActorID, "ADMIN", "portal:landing:updated", "portal_landing_page", page.ID, audit.StatusSuccess, map[string]interface{}{
		"admissionsOpen": req.AdmissionsOpen,
		"cycle":          req.AdmissionsCycle,
	})

	return page, nil
}

type CreateProgramRequest struct {
	TenantID            string           `json:"tenantId"`
	ProgramCode         string           `json:"programCode"`
	ProgramName         string           `json:"programName"`
	DegreeType          PortalDegreeType `json:"degreeType"`
	DepartmentName      string           `json:"departmentName"`
	DurationYears       int              `json:"durationYears"`
	TotalSemesters      int              `json:"totalSemesters"`
	EligibilityCriteria string           `json:"eligibilityCriteria"`
	AnnualFee           float64          `json:"annualFee"`
	IsFeatured          bool             `json:"isFeatured"`
	ActorID             string           `json:"actorId"`
}

func (s *service) CreateProgram(ctx context.Context, req CreateProgramRequest) (*PortalProgramCatalog, error) {
	if req.DegreeType == "" {
		req.DegreeType = DegreeUG
	}
	if req.DurationYears <= 0 {
		req.DurationYears = 4
	}
	if req.TotalSemesters <= 0 {
		req.TotalSemesters = req.DurationYears * 2
	}

	now := time.Now().UTC()
	program := &PortalProgramCatalog{
		ID:                  uuid.New().String(),
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
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.repo.CreateProgram(ctx, program); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, req.ActorID, "ADMIN", "portal:program:created", "portal_program_catalog", program.ID, audit.StatusSuccess, map[string]interface{}{
		"programCode": program.ProgramCode,
		"programName": program.ProgramName,
	})

	return program, nil
}

func (s *service) ListPrograms(ctx context.Context, filter ProgramCatalogFilter) ([]PortalProgramCatalog, error) {
	return s.repo.ListPrograms(ctx, filter)
}

func (s *service) GetProgram(ctx context.Context, tenantID, id string) (*PortalProgramCatalog, error) {
	return s.repo.GetProgramByID(ctx, tenantID, id)
}

type SubmitInquiryRequest struct {
	TenantID          string `json:"tenantId"`
	ProspectName      string `json:"prospectName"`
	ProspectEmail     string `json:"prospectEmail"`
	ProspectPhone     string `json:"prospectPhone"`
	ProgramOfInterest string `json:"programOfInterest"`
	Message           string `json:"message"`
}

func (s *service) SubmitInquiry(ctx context.Context, req SubmitInquiryRequest) (*PortalPublicInquiry, error) {
	now := time.Now().UTC()
	inquiry := &PortalPublicInquiry{
		ID:                uuid.New().String(),
		TenantID:          req.TenantID,
		ProspectName:      req.ProspectName,
		ProspectEmail:     req.ProspectEmail,
		ProspectPhone:     req.ProspectPhone,
		ProgramOfInterest: req.ProgramOfInterest,
		Message:           req.Message,
		Status:            InquiryNew,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.repo.CreateInquiry(ctx, inquiry); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, "PROSPECT", "USER", "portal:inquiry:submitted", "portal_public_inquiry", inquiry.ID, audit.StatusSuccess, map[string]interface{}{
		"prospectEmail": req.ProspectEmail,
		"program":       req.ProgramOfInterest,
	})

	return inquiry, nil
}

func (s *service) ListInquiries(ctx context.Context, filter PublicInquiryFilter) ([]PortalPublicInquiry, error) {
	return s.repo.ListInquiries(ctx, filter)
}

type UpdateInquiryStatusRequest struct {
	TenantID            string              `json:"tenantId"`
	InquiryID           string              `json:"inquiryId"`
	Status              PortalInquiryStatus `json:"status"`
	AssignedCounselorID *string             `json:"assignedCounselorId,omitempty"`
	Notes               *string             `json:"notes,omitempty"`
	ActorID             string              `json:"actorId"`
}

func (s *service) UpdateInquiryStatus(ctx context.Context, req UpdateInquiryStatusRequest) (*PortalPublicInquiry, error) {
	inquiry, err := s.repo.GetInquiryByID(ctx, req.TenantID, req.InquiryID)
	if err != nil {
		return nil, err
	}

	inquiry.Status = req.Status
	if req.AssignedCounselorID != nil {
		inquiry.AssignedCounselorID = req.AssignedCounselorID
	}
	if req.Notes != nil {
		inquiry.Notes = req.Notes
	}
	inquiry.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateInquiry(ctx, inquiry); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, req.ActorID, "STAFF", "portal:inquiry:status_changed", "portal_public_inquiry", inquiry.ID, audit.StatusSuccess, map[string]interface{}{
		"status": req.Status,
	})

	return inquiry, nil
}
