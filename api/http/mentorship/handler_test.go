/**
 * BLOCK_API_MENTORSHIP_HANDLER_TEST_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   HTTP integration tests verifying REST endpoints for mentor allocations, sessions, academic progress, and risk alerts.
 */

package mentorship_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	apiMentorship "campus/api/http/mentorship"
	backendMentorship "campus/backend/mentorship"
)

func setupTestRouter() (chi.Router, backendMentorship.Service) {
	r := chi.NewRouter()
	repo := backendMentorship.NewMockRepository()
	svc := backendMentorship.NewService(repo, nil)
	handler := apiMentorship.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_MentorshipAllocationAndSession(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Allocate Mentor
	allocPayload := map[string]interface{}{
		"tenant_id":       "tenant-test",
		"student_id":      "student-001",
		"mentor_staff_id": "staff-mentor-01",
		"cohort_id":       "cohort-2026-cse",
	}
	body, _ := json.Marshal(allocPayload)
	req := httptest.NewRequest("POST", "/api/v1/mentorship/allocations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for mentor allocation, got %d: %s", rec.Code, rec.Body.String())
	}

	var alloc backendMentorship.MentorAllocation
	_ = json.Unmarshal(rec.Body.Bytes(), &alloc)

	// 2. Schedule Session
	schedPayload := map[string]interface{}{
		"tenant_id":     "tenant-test",
		"allocation_id": alloc.ID,
		"scheduled_at":  time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339),
		"location":      "Dean Office, Room 102",
		"meeting_type":  "ONE_ON_ONE",
	}
	schedBody, _ := json.Marshal(schedPayload)
	reqSched := httptest.NewRequest("POST", "/api/v1/mentorship/sessions", bytes.NewReader(schedBody))
	reqSched.Header.Set("Content-Type", "application/json")
	recSched := httptest.NewRecorder()
	router.ServeHTTP(recSched, reqSched)

	if recSched.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for session, got %d: %s", recSched.Code, recSched.Body.String())
	}

	var sess backendMentorship.MentorshipSession
	_ = json.Unmarshal(recSched.Body.Bytes(), &sess)

	// 3. Complete Session
	compPayload := map[string]interface{}{
		"tenant_id":          "tenant-test",
		"discussion_summary": "Discussed career paths, semester electives, and time management.",
		"action_items":       "Register for Machine Learning elective before Friday.",
	}
	compBody, _ := json.Marshal(compPayload)
	reqComp := httptest.NewRequest("POST", "/api/v1/mentorship/sessions/"+sess.ID+"/complete", bytes.NewReader(compBody))
	reqComp.Header.Set("Content-Type", "application/json")
	recComp := httptest.NewRecorder()
	router.ServeHTTP(recComp, reqComp)

	if recComp.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for complete session, got %d: %s", recComp.Code, recComp.Body.String())
	}
}

func TestHTTP_MentorshipProgressAndAlerts(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Record Progress (Critical At-Risk)
	progPayload := map[string]interface{}{
		"tenant_id":             "tenant-test",
		"student_id":            "student-002",
		"semester":              2,
		"sgpa":                  3.6,
		"cgpa":                  4.2,
		"attendance_percentage": 58.0,
		"credits_earned":        10,
		"total_credits":         22,
		"remarks":               "Attendance shortage and failed 2 courses.",
	}
	body, _ := json.Marshal(progPayload)
	req := httptest.NewRequest("POST", "/api/v1/mentorship/progress", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for progress, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. Query Alerts
	reqAlerts := httptest.NewRequest("GET", "/api/v1/mentorship/alerts?tenant_id=tenant-test&student_id=student-002&resolved=false", nil)
	recAlerts := httptest.NewRecorder()
	router.ServeHTTP(recAlerts, reqAlerts)

	if recAlerts.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for list alerts, got %d: %s", recAlerts.Code, recAlerts.Body.String())
	}

	var alerts []*backendMentorship.MentorshipAtRiskAlert
	_ = json.Unmarshal(recAlerts.Body.Bytes(), &alerts)
	if len(alerts) != 1 {
		t.Fatalf("Expected 1 automated alert, got %d", len(alerts))
	}

	// 3. Resolve Alert
	resPayload := map[string]interface{}{
		"tenant_id":      "tenant-test",
		"resolved_by_id": "staff-dean-01",
	}
	resBody, _ := json.Marshal(resPayload)
	reqRes := httptest.NewRequest("POST", "/api/v1/mentorship/alerts/"+alerts[0].ID+"/resolve", bytes.NewReader(resBody))
	reqRes.Header.Set("Content-Type", "application/json")
	recRes := httptest.NewRecorder()
	router.ServeHTTP(recRes, reqRes)

	if recRes.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for resolve alert, got %d: %s", recRes.Code, recRes.Body.String())
	}
}
