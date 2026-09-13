/**
 * BLOCK_API_PORTAL_HANDLER_TEST_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   HTTP integration tests verifying REST endpoints for public landing, program catalog, and prospect lead capture.
 */

package portal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apiPortal "campus/api/http/portal"
	backendPortal "campus/backend/portal"
)

func setupTestRouter() (chi.Router, backendPortal.Service) {
	r := chi.NewRouter()
	repo := backendPortal.NewMockRepository()
	svc := backendPortal.NewService(repo, nil)
	handler := apiPortal.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_PortalLifecycle(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Update Landing Page
	landingPayload := map[string]interface{}{
		"tenant_id":        "tenant-portal-http",
		"hero_headline":    "Shape the Future of Computing",
		"hero_subheadline": "Admissions open for Academic Year 2026-2027.",
		"admissions_open":  true,
		"admissions_cycle": "2026-2027",
		"contact_email":    "admissions@campus.edu",
		"contact_phone":    "+91 80 1234 5678",
		"campus_address":   "Bangalore Campus",
		"is_published":     true,
		"actor_id":         "admin-super",
	}
	body, _ := json.Marshal(landingPayload)
	req := httptest.NewRequest("PUT", "/api/v1/portal/landing", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for update landing, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. Create Program
	progPayload := map[string]interface{}{
		"tenant_id":            "tenant-portal-http",
		"program_code":         "BTECH_CSE",
		"programName":          "B.Tech in Computer Science & Engineering",
		"program_name":         "B.Tech in Computer Science & Engineering",
		"degree_type":          "UG",
		"department_name":      "Computer Science & Engineering",
		"duration_years":       4,
		"total_semesters":      8,
		"eligibility_criteria": "10+2 with Physics, Mathematics, Chemistry >= 75%",
		"annual_fee":           200000.0,
		"is_featured":          true,
		"actor_id":             "admin-super",
	}
	progBody, _ := json.Marshal(progPayload)
	reqProg := httptest.NewRequest("POST", "/api/v1/portal/programs", bytes.NewReader(progBody))
	reqProg.Header.Set("Content-Type", "application/json")
	recProg := httptest.NewRecorder()
	router.ServeHTTP(recProg, reqProg)

	if recProg.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for program creation, got %d: %s", recProg.Code, recProg.Body.String())
	}

	// 3. Submit Public Prospect Inquiry
	inqPayload := map[string]interface{}{
		"tenant_id":           "tenant-portal-http",
		"prospect_name":       "Rohan Gupta",
		"prospect_email":      "rohan.gupta@example.com",
		"prospect_phone":      "+91 91234 56789",
		"program_of_interest": "B.Tech in Computer Science & Engineering",
		"message":             "Interested in hostel accommodation fee details.",
	}
	inqBody, _ := json.Marshal(inqPayload)
	reqInq := httptest.NewRequest("POST", "/api/v1/portal/inquiries", bytes.NewReader(inqBody))
	reqInq.Header.Set("Content-Type", "application/json")
	recInq := httptest.NewRecorder()
	router.ServeHTTP(recInq, reqInq)

	if recInq.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for inquiry submission, got %d: %s", recInq.Code, recInq.Body.String())
	}
}
