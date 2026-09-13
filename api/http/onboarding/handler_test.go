/**
 * BLOCK_API_ONBOARDING_HANDLER_TEST_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   HTTP transport integration tests verifying REST endpoints, RFC 7807 problem details, and state flows.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"campus/backend/onboarding"
)

func setupTestServer() (*Handler, *onboarding.Service, *onboarding.MockAcademicRepository, *chi.Mux) {
	appRepo := onboarding.NewMockApplicantRepository()
	docRepo := onboarding.NewMockDocumentRepository()
	acadRepo := onboarding.NewMockAcademicRepository()
	profRepo := onboarding.NewMockProfileRepository()
	seqRepo := onboarding.NewMockSequenceRepository()
	seqEngine := onboarding.NewSequenceEngine(seqRepo)

	svc := onboarding.NewService(appRepo, docRepo, acadRepo, profRepo, seqEngine, nil)
	handler := NewHandler(svc)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	return handler, svc, acadRepo, r
}

func TestHTTP_CreateDraft_Success(t *testing.T) {
	_, _, _, r := setupTestServer()

	payload := map[string]interface{}{
		"tenant_id":     "tenant-api",
		"type":          "STUDENT",
		"first_name":    "Nikhil",
		"last_name":     "Kumar",
		"email":         "nikhil.kumar@example.com",
		"academic_year": "2026-2027",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/onboarding/applicants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); loc == "" {
		t.Fatalf("expected Location header in 201 response")
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	app, ok := resp["applicant"].(map[string]interface{})
	if !ok || app["first_name"] != "Nikhil" {
		t.Fatalf("unexpected response structure: %+v", resp)
	}
}

func TestHTTP_CreateDraft_InvalidPayload(t *testing.T) {
	_, _, _, r := setupTestServer()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/onboarding/applicants", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on invalid JSON, got %d", rec.Code)
	}
}

func TestHTTP_GetApplicant_FoundAndNotFound(t *testing.T) {
	_, svc, _, r := setupTestServer()
	ctx := context.Background()

	app, _ := svc.CreateDraft(ctx, onboarding.CreateApplicantCommand{
		TenantID:  "tenant-api",
		Type:      onboarding.TypeStudent,
		FirstName: "Sunil",
		LastName:  "Gavali",
		Email:     "sunil@example.com",
	})

	// 1. Found
	req := httptest.NewRequest(http.MethodGet, "/api/v1/onboarding/applicants/"+app.ID+"?tenant_id=tenant-api", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. Not Found
	req = httptest.NewRequest(http.MethodGet, "/api/v1/onboarding/applicants/non-existent-id?tenant_id=tenant-api", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHTTP_ListApplicants_WithPagination(t *testing.T) {
	_, svc, _, r := setupTestServer()
	ctx := context.Background()
	tenantID := "tenant-list"

	for i := 0; i < 5; i++ {
		_, _ = svc.CreateDraft(ctx, onboarding.CreateApplicantCommand{
			TenantID:  tenantID,
			Type:      onboarding.TypeStudent,
			FirstName: "Test",
			LastName:  "Student",
			Email:     "student" + string(rune('a'+i)) + "@test.com",
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/onboarding/applicants?tenant_id="+tenantID+"&limit=2&offset=1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp["total"].(float64) != 5 {
		t.Fatalf("expected total=5, got %v", resp["total"])
	}
	apps := resp["applicants"].([]interface{})
	if len(apps) != 2 {
		t.Fatalf("expected 2 applicants returned for limit=2, got %d", len(apps))
	}
}

func TestHTTP_FullStudentEnrollment_Flow(t *testing.T) {
	_, _, acadRepo, r := setupTestServer()
	ctx := context.Background()
	tenantID := "tenant-flow"

	dept := &onboarding.AcademicDepartment{ID: "dept-cse", TenantID: tenantID, Code: "CSE", Name: "Computer Science", IsActive: true}
	_ = acadRepo.CreateDepartment(ctx, dept)
	prog := &onboarding.AcademicProgram{ID: "prog-cse", TenantID: tenantID, DepartmentID: dept.ID, Code: "BTECH_CSE", Name: "B.Tech CSE", DegreeType: onboarding.DegreeUG, IsActive: true}
	_ = acadRepo.CreateProgram(ctx, prog)
	cohort := &onboarding.CohortBatch{ID: "coh-cse", TenantID: tenantID, ProgramID: prog.ID, AcademicYear: "2026-2027", StartYear: 2026, EndYear: 2030, Section: "A", MaxCapacity: 60, IsActive: true}
	_ = acadRepo.CreateCohort(ctx, cohort)

	// Step 1: Create Draft
	draftPayload := map[string]interface{}{
		"tenant_id":     tenantID,
		"type":          "STUDENT",
		"first_name":    "Meera",
		"last_name":     "Rao",
		"email":         "meera.rao@example.com",
		"phone":         "+919876543111",
		"date_of_birth": "2006-05-12",
		"gender":        "FEMALE",
		"nationality":   "Indian",
		"emergency_contact": map[string]string{
			"name":     "Venkatesh Rao",
			"phone":    "+919876543112",
			"relation": "Father",
		},
		"address": map[string]string{
			"line1":       "123 Palm Grove",
			"city":        "Bengaluru",
			"state":       "Karnataka",
			"postal_code": "560034",
			"country":     "India",
		},
		"program_id":    prog.ID,
		"academic_year": "2026-2027",
		"guardian": map[string]string{
			"name":     "Venkatesh Rao",
			"phone":    "+919876543112",
			"relation": "Father",
		},
	}
	body, _ := json.Marshal(draftPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/onboarding/applicants", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on draft, got %d: %s", rec.Code, rec.Body.String())
	}
	var draftResp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&draftResp)
	appID := draftResp["applicant"].(map[string]interface{})["id"].(string)

	// Step 2: Attach 3 documents
	docs := []struct {
		dType string
		fKey  string
		fName string
		mime  string
	}{
		{"NATIONAL_ID", "uploads/aadhaar.pdf", "aadhaar.pdf", "application/pdf"},
		{"ACADEMIC_TRANSCRIPT", "uploads/marksheet.pdf", "marksheet.pdf", "application/pdf"},
		{"PASSPORT_PHOTO", "uploads/photo.jpg", "photo.jpg", "image/jpeg"},
	}

	docIDs := make([]string, 0)
	for _, d := range docs {
		docBody, _ := json.Marshal(map[string]interface{}{
			"tenant_id":     tenantID,
			"document_type": d.dType,
			"file_key":      d.fKey,
			"file_name":     d.fName,
			"file_size":     100000,
			"mime_type":     d.mime,
		})
		docReq := httptest.NewRequest(http.MethodPost, "/api/v1/onboarding/applicants/"+appID+"/documents", bytes.NewReader(docBody))
		docRec := httptest.NewRecorder()
		r.ServeHTTP(docRec, docReq)
		if docRec.Code != http.StatusCreated {
			t.Fatalf("failed to attach doc %s: %s", d.dType, docRec.Body.String())
		}
		var dResp map[string]interface{}
		_ = json.NewDecoder(docRec.Body).Decode(&dResp)
		docIDs = append(docIDs, dResp["document"].(map[string]interface{})["id"].(string))
	}

	// Step 3: Submit Application
	subReq := httptest.NewRequest(http.MethodPost, "/api/v1/onboarding/applicants/"+appID+"/submit?tenant_id="+tenantID, nil)
	subRec := httptest.NewRecorder()
	r.ServeHTTP(subRec, subReq)
	if subRec.Code != http.StatusOK {
		t.Fatalf("failed to submit application: %d: %s", subRec.Code, subRec.Body.String())
	}

	// Step 4: Assign Reviewer
	revBody, _ := json.Marshal(map[string]string{
		"tenant_id":   tenantID,
		"reviewer_id": "usr_admissions_admin",
	})
	revReq := httptest.NewRequest(http.MethodPost, "/api/v1/onboarding/applicants/"+appID+"/review", bytes.NewReader(revBody))
	revRec := httptest.NewRecorder()
	r.ServeHTTP(revRec, revReq)
	if revRec.Code != http.StatusOK {
		t.Fatalf("failed to assign reviewer: %s", revRec.Body.String())
	}

	// Step 5: Verify all documents
	for _, did := range docIDs {
		vBody, _ := json.Marshal(map[string]interface{}{
			"tenant_id":   tenantID,
			"approved":    true,
			"verifier_id": "usr_admissions_admin",
		})
		vReq := httptest.NewRequest(http.MethodPatch, "/api/v1/onboarding/applicants/"+appID+"/documents/"+did, bytes.NewReader(vBody))
		vRec := httptest.NewRecorder()
		r.ServeHTTP(vRec, vReq)
		if vRec.Code != http.StatusOK {
			t.Fatalf("failed to verify doc %s: %s", did, vRec.Body.String())
		}
	}

	// Step 6: Verify Application
	verReq := httptest.NewRequest(http.MethodPost, "/api/v1/onboarding/applicants/"+appID+"/verify?tenant_id="+tenantID, nil)
	verRec := httptest.NewRecorder()
	r.ServeHTTP(verRec, verReq)
	if verRec.Code != http.StatusOK {
		t.Fatalf("failed to verify application: %s", verRec.Body.String())
	}

	// Step 7: Enroll Student
	enrollBody, _ := json.Marshal(map[string]string{
		"tenant_id": tenantID,
		"user_id":   "usr_student_meera",
		"cohort_id": cohort.ID,
		"admin_id":  "usr_admissions_admin",
	})
	enrReq := httptest.NewRequest(http.MethodPost, "/api/v1/onboarding/applicants/"+appID+"/enroll", bytes.NewReader(enrollBody))
	enrRec := httptest.NewRecorder()
	r.ServeHTTP(enrRec, enrReq)
	if enrRec.Code != http.StatusCreated {
		t.Fatalf("failed to enroll student: %d: %s", enrRec.Code, enrRec.Body.String())
	}

	var enrResp map[string]interface{}
	_ = json.NewDecoder(enrRec.Body).Decode(&enrResp)
	if enrResp["roll_number"] != "2026-BTECH_CSE-0001" {
		t.Fatalf("expected roll number '2026-BTECH_CSE-0001', got %v", enrResp["roll_number"])
	}
}

func TestHTTP_AcademicCatalog_Endpoints(t *testing.T) {
	_, _, acadRepo, r := setupTestServer()
	ctx := context.Background()
	tenantID := "tenant-catalog"

	dept := &onboarding.AcademicDepartment{ID: "d1", TenantID: tenantID, Code: "MECH", Name: "Mechanical Engg", IsActive: true}
	_ = acadRepo.CreateDepartment(ctx, dept)
	prog := &onboarding.AcademicProgram{ID: "p1", TenantID: tenantID, DepartmentID: dept.ID, Code: "BTECH_MECH", Name: "B.Tech Mech", DegreeType: onboarding.DegreeUG, IsActive: true}
	_ = acadRepo.CreateProgram(ctx, prog)
	cohort := &onboarding.CohortBatch{ID: "c1", TenantID: tenantID, ProgramID: prog.ID, AcademicYear: "2026-2027", StartYear: 2026, EndYear: 2030, Section: "A", MaxCapacity: 60, IsActive: true}
	_ = acadRepo.CreateCohort(ctx, cohort)

	// GET /departments
	req := httptest.NewRequest(http.MethodGet, "/api/v1/onboarding/departments?tenant_id="+tenantID, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on /departments, got %d", rec.Code)
	}

	// GET /programs
	req = httptest.NewRequest(http.MethodGet, "/api/v1/onboarding/programs?tenant_id="+tenantID+"&department_id="+dept.ID, nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on /programs, got %d", rec.Code)
	}

	// GET /cohorts
	req = httptest.NewRequest(http.MethodGet, "/api/v1/onboarding/cohorts?tenant_id="+tenantID+"&program_id="+prog.ID+"&academic_year=2026-2027", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on /cohorts, got %d", rec.Code)
	}
}
