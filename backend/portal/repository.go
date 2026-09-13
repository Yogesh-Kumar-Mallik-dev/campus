/**
 * BLOCK_PORTAL_REPO_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   Repository interface contracts for landing configuration, program catalog, and prospective applicant inquiries.
 */

package portal

import "context"

type Repository interface {
	// Landing Page Config
	GetLandingPage(ctx context.Context, tenantID string) (*PortalLandingPage, error)
	UpsertLandingPage(ctx context.Context, page *PortalLandingPage) error

	// Program Catalog
	CreateProgram(ctx context.Context, program *PortalProgramCatalog) error
	GetProgramByID(ctx context.Context, tenantID, id string) (*PortalProgramCatalog, error)
	GetProgramByCode(ctx context.Context, tenantID, code string) (*PortalProgramCatalog, error)
	ListPrograms(ctx context.Context, filter ProgramCatalogFilter) ([]PortalProgramCatalog, error)
	UpdateProgram(ctx context.Context, program *PortalProgramCatalog) error

	// Public Inquiries
	CreateInquiry(ctx context.Context, inquiry *PortalPublicInquiry) error
	GetInquiryByID(ctx context.Context, tenantID, id string) (*PortalPublicInquiry, error)
	ListInquiries(ctx context.Context, filter PublicInquiryFilter) ([]PortalPublicInquiry, error)
	UpdateInquiry(ctx context.Context, inquiry *PortalPublicInquiry) error
}
