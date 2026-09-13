/**
 * BLOCK_PORTAL_DOMAIN_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   Domain entities, value objects, and business invariant rules for institutional landing, program catalog, and admissions inquiries.
 */

package portal

import "time"

type PortalDegreeType string

const (
	DegreeUG      PortalDegreeType = "UG"
	DegreePG      PortalDegreeType = "PG"
	DegreeDiploma PortalDegreeType = "DIPLOMA"
	DegreePhD     PortalDegreeType = "PHD"
)

type PortalInquiryStatus string

const (
	InquiryNew       PortalInquiryStatus = "NEW"
	InquiryContacted PortalInquiryStatus = "CONTACTED"
	InquiryConverted PortalInquiryStatus = "CONVERTED"
	InquiryClosed    PortalInquiryStatus = "CLOSED"
)

type PortalLandingPage struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenantId"`
	HeroHeadline    string    `json:"heroHeadline"`
	HeroSubheadline string    `json:"heroSubheadline"`
	AdmissionsOpen  bool      `json:"admissionsOpen"`
	AdmissionsCycle string    `json:"admissionsCycle"`
	HeroImageURL    *string   `json:"heroImageUrl,omitempty"`
	ContactEmail    string    `json:"contactEmail"`
	ContactPhone    string    `json:"contactPhone"`
	CampusAddress   string    `json:"campusAddress"`
	IsPublished     bool      `json:"isPublished"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type PortalProgramCatalog struct {
	ID                  string           `json:"id"`
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
	CreatedAt           time.Time        `json:"createdAt"`
	UpdatedAt           time.Time        `json:"updatedAt"`
}

type PortalPublicInquiry struct {
	ID                  string              `json:"id"`
	TenantID            string              `json:"tenantId"`
	ProspectName        string              `json:"prospectName"`
	ProspectEmail       string              `json:"prospectEmail"`
	ProspectPhone       string              `json:"prospectPhone"`
	ProgramOfInterest   string              `json:"programOfInterest"`
	Message             string              `json:"message"`
	Status              PortalInquiryStatus `json:"status"`
	AssignedCounselorID *string             `json:"assignedCounselorId,omitempty"`
	Notes               *string             `json:"notes,omitempty"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
}

type ProgramCatalogFilter struct {
	TenantID   string
	DegreeType *PortalDegreeType
	Featured   *bool
}

type PublicInquiryFilter struct {
	TenantID            string
	Status              *PortalInquiryStatus
	AssignedCounselorID *string
}
