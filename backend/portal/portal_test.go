/**
 * BLOCK_PORTAL_TEST_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   Unit test suite covering landing page config, program catalog filtering, duplicate code guards, and prospect inquiry triage.
 */

package portal_test

import (
	"context"
	"errors"
	"testing"

	"campus/backend/portal"
)

func setupPortalService() (portal.Service, *portal.MockRepository) {
	mockRepo := portal.NewMockRepository()
	srv := portal.NewService(mockRepo, nil)
	return srv, mockRepo
}

func TestPortal_LandingPageConfiguration(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupPortalService()
	tenantID := "tenant-portal-1"

	// Initial fetch returns fallback defaults
	initial, err := srv.GetLandingPage(ctx, tenantID)
	if err != nil || initial.TenantID != tenantID {
		t.Fatalf("expected initial default landing page, got err: %v, page: %+v", err, initial)
	}

	// Update landing page
	updated, err := srv.UpdateLandingPage(ctx, portal.UpdateLandingPageRequest{
		TenantID:        tenantID,
		HeroHeadline:    "Excellence in Global Technology & AI",
		HeroSubheadline: "Admissions for Fall 2026 are now open across 24 multidisciplinary programs.",
		AdmissionsOpen:  true,
		AdmissionsCycle: "2026-2027",
		ContactEmail:    "apply@campus.edu",
		ContactPhone:    "+91 80 9876 5432",
		CampusAddress:   "Silicon Plateau, Outer Ring Road, Bangalore",
		IsPublished:     true,
		ActorID:         "admin-super",
	})
	if err != nil || updated.HeroHeadline != "Excellence in Global Technology & AI" {
		t.Fatalf("expected updated landing page, got err: %v, page: %+v", err, updated)
	}
}

func TestPortal_ProgramCatalogManagementAndDuplicates(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupPortalService()
	tenantID := "tenant-portal-1"

	// 1. Create Program
	p1, err := srv.CreateProgram(ctx, portal.CreateProgramRequest{
		TenantID:            tenantID,
		ProgramCode:         "BTECH_AI_ML",
		ProgramName:         "B.Tech in Artificial Intelligence & Machine Learning",
		DegreeType:          portal.DegreeUG,
		DepartmentName:      "Computer Science & Engineering",
		DurationYears:       4,
		TotalSemesters:      8,
		EligibilityCriteria: "10+2 with Physics, Mathematics, and Chemistry with minimum 75% aggregate.",
		AnnualFee:           250000.0,
		IsFeatured:          true,
		ActorID:             "admin-super",
	})
	if err != nil || p1.ProgramCode != "BTECH_AI_ML" {
		t.Fatalf("expected program creation success, got err: %v, prog: %+v", err, p1)
	}

	// 2. Duplicate Program Code -> REJECTED
	_, err = srv.CreateProgram(ctx, portal.CreateProgramRequest{
		TenantID:       tenantID,
		ProgramCode:    "BTECH_AI_ML",
		ProgramName:    "Duplicate Program",
		DepartmentName: "CSE",
		ActorID:        "admin-super",
	})
	if !errors.Is(err, portal.ErrDuplicateProgramCode) {
		t.Fatalf("expected ErrDuplicateProgramCode, got: %v", err)
	}

	// 3. List Programs with Filters
	featured := true
	list, err := srv.ListPrograms(ctx, portal.ProgramCatalogFilter{
		TenantID: tenantID,
		Featured: &featured,
	})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 featured program, got: %d", len(list))
	}
}

func TestPortal_ProspectInquiryLifecycle(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupPortalService()
	tenantID := "tenant-portal-1"

	// 1. Submit Public Prospect Inquiry
	inq, err := srv.SubmitInquiry(ctx, portal.SubmitInquiryRequest{
		TenantID:          tenantID,
		ProspectName:      "Aman Sharma",
		ProspectEmail:     "aman.sharma@example.com",
		ProspectPhone:     "+91 99887 76655",
		ProgramOfInterest: "B.Tech in Artificial Intelligence & Machine Learning",
		Message:           "Looking for scholarship eligibility details for state board rank holders.",
	})
	if err != nil || inq.Status != portal.InquiryNew {
		t.Fatalf("expected inquiry submission success, got err: %v, inq: %+v", err, inq)
	}

	// 2. Admissions Counselor Triage (Update Status to CONTACTED)
	counselorID := "counselor-admissions-1"
	notes := "Called prospect. Explained 50% merit waiver criteria. Application form link shared via email."
	updatedInq, err := srv.UpdateInquiryStatus(ctx, portal.UpdateInquiryStatusRequest{
		TenantID:            tenantID,
		InquiryID:           inq.ID,
		Status:              portal.InquiryContacted,
		AssignedCounselorID: &counselorID,
		Notes:               &notes,
		ActorID:             counselorID,
	})
	if err != nil || updatedInq.Status != portal.InquiryContacted || *updatedInq.AssignedCounselorID != counselorID {
		t.Fatalf("expected updated inquiry with counselor, got err: %v, inq: %+v", err, updatedInq)
	}
}
