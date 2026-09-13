/**
 * BLOCK_API_ATTENDANCE_HANDLER_TEST_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   HTTP transport tests for timetable, sessions, roll call, student summary, and medical leave condonations.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package attendance

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"campus/backend/attendance"
)

func setupTestHTTPServer() (*Handler, *attendance.Service, *attendance.MockCatalogRepository, *chi.Mux) {
	catRepo := attendance.NewMockCatalogRepository()
	sessRepo := attendance.NewMockSessionRepository()
	recRepo := attendance.NewMockRecordRepository()
	leaveRepo := attendance.NewMockMedicalLeaveRepository()

	svc := attendance.NewService(catRepo, sessRepo, recRepo, leaveRepo, nil)
	handler := NewHandler(svc)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	return handler, svc, catRepo, r
}

func TestHTTP_CreateSubject(t *testing.T) {
	_, _, _, r := setupTestHTTPServer()

	payload := map[string]interface{}{
		"tenant_id":  "ten_http",
		"program_id": "prog_cs",
		"code":       "cs201",
		"name":       "Algorithms and Analysis",
		"credits":    4,
		"semester":   2,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/attendance/subjects", bytes.NewReader(body))
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
	sub, ok := resp["subject"].(map[string]interface{})
	if !ok || sub["code"] != "CS201" {
		t.Fatalf("unexpected subject payload: %+v", resp)
	}
}

func TestHTTP_ScheduleAndMarkAttendance(t *testing.T) {
	_, _, _, r := setupTestHTTPServer()

	// 1. Schedule Session
	schedPayload := map[string]interface{}{
		"tenant_id":    "ten_http",
		"subject_id":   "sub_algo",
		"cohort_id":    "coh_cs_2026",
		"facultyID":    "fac_algo_prof",
		"session_date": "2026-09-18",
		"start_time":   "10:00",
		"end_time":     "11:00",
		"mode":         "MANUAL_FACULTY",
	}
	body, _ := json.Marshal(schedPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/attendance/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for session, got %d: %s", rec.Code, rec.Body.String())
	}

	var sessResp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&sessResp)
	sessObj := sessResp["session"].(map[string]interface{})
	sessionID := sessObj["id"].(string)

	// 2. Mark Attendance for 2 students
	markPayload := map[string]interface{}{
		"tenant_id":  "ten_http",
		"session_id": sessionID,
		"faculty_id": "fac_algo_prof",
		"records": []map[string]interface{}{
			{"student_id": "stu_nikhil", "status": "PRESENT"},
			{"student_id": "stu_aarav", "status": "ABSENT"},
		},
	}
	body, _ = json.Marshal(markPayload)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/attendance/sessions/"+sessionID+"/mark", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for mark attendance, got %d: %s", rec.Code, rec.Body.String())
	}

	var markResp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&markResp)
	markedSess := markResp["session"].(map[string]interface{})
	if int(markedSess["present_count"].(float64)) != 1 || int(markedSess["absent_count"].(float64)) != 1 {
		t.Fatalf("unexpected counts in marked session: %+v", markedSess)
	}
}

func TestHTTP_StudentSummaryAndShortage(t *testing.T) {
	_, svc, _, r := setupTestHTTPServer()
	ctx := context.Background()

	tenantID := "ten_http"
	studentID := "stu_aarav"
	subjectID := "sub_algo"

	// Create 2 sessions where student is absent (0% -> shortage)
	s1, _ := svc.ScheduleSession(ctx, attendance.ScheduleSessionCommand{
		TenantID:    tenantID,
		SubjectID:   subjectID,
		CohortID:    "coh_cs_2026",
		FacultyID:   "fac_1",
		SessionDate: "2026-09-01",
	})
	s2, _ := svc.ScheduleSession(ctx, attendance.ScheduleSessionCommand{
		TenantID:    tenantID,
		SubjectID:   subjectID,
		CohortID:    "coh_cs_2026",
		FacultyID:   "fac_1",
		SessionDate: "2026-09-02",
	})
	_, _ = svc.MarkAttendance(ctx, attendance.MarkAttendanceCommand{
		TenantID:  tenantID,
		SessionID: s1.ID,
		FacultyID: "fac_1",
		Records:   []attendance.RecordInput{{StudentID: studentID, Status: attendance.RecordAbsent}},
	})
	_, _ = svc.MarkAttendance(ctx, attendance.MarkAttendanceCommand{
		TenantID:  tenantID,
		SessionID: s2.ID,
		FacultyID: "fac_1",
		Records:   []attendance.RecordInput{{StudentID: studentID, Status: attendance.RecordAbsent}},
	})

	// Call GET summary endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/attendance/summary/students/"+studentID+"?tenant_id="+tenantID+"&subject_id="+subjectID, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	sum := resp["summary"].(map[string]interface{})
	if sum["is_shortage"] != true || sum["attendance_percentage"].(float64) != 0.0 {
		t.Fatalf("expected shortage=true and 0.0%%, got %+v", sum)
	}

	// Apply medical leave via HTTP
	leavePayload := map[string]interface{}{
		"tenant_id":       tenantID,
		"student_id":      studentID,
		"from_date":       "2026-09-01",
		"to_date":         "2026-09-02",
		"reason":          "Viral Fever",
		"certificate_key": "s3://medical/fever.pdf",
	}
	body, _ := json.Marshal(leavePayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/attendance/leaves", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for leave, got %d: %s", rec.Code, rec.Body.String())
	}

	var leaveResp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&leaveResp)
	leaveObj := leaveResp["leave"].(map[string]interface{})
	leaveID := leaveObj["id"].(string)

	// Approve medical leave via HTTP
	decisionPayload := map[string]interface{}{
		"tenant_id":   tenantID,
		"approver_id": "fac_dean",
		"approved":    true,
	}
	body, _ = json.Marshal(decisionPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/attendance/leaves/"+leaveID+"/decision", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for decision, got %d: %s", rec.Code, rec.Body.String())
	}

	// Re-verify summary via HTTP: now 100% attendance and shortage = false
	req = httptest.NewRequest(http.MethodGet, "/api/v1/attendance/summary/students/"+studentID+"?tenant_id="+tenantID+"&subject_id="+subjectID, nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	var resp2 map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&resp2)
	sum2 := resp2["summary"].(map[string]interface{})
	if sum2["is_shortage"] != false || sum2["attendance_percentage"].(float64) != 100.0 {
		t.Fatalf("expected shortage=false and 100.0%% after approval, got %+v", sum2)
	}
}
