/**
 * BLOCK_API_WHISTLEBLOWER_HANDLER_TEST_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   HTTP integration tests verifying zero-knowledge report submission, secret token lookup, anonymous messaging, and committee investigations.
 */

package whistleblower_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apiWB "campus/api/http/whistleblower"
	backendWB "campus/backend/whistleblower"
)

func setupTestRouter() (chi.Router, backendWB.Service) {
	r := chi.NewRouter()
	repo := backendWB.NewMockRepository()
	svc := backendWB.NewService(repo, nil)
	handler := apiWB.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_WhistleblowerLifecycle(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Submit Anonymous Report
	submitPayload := map[string]interface{}{
		"tenant_id":   "tenant-wb-http",
		"category":    "ANTI_RAGGING",
		"severity":    "CRITICAL",
		"title":       "Hostel midnight ragging report",
		"description": "Forced physical drills in 4th floor corridor.",
	}
	body, _ := json.Marshal(submitPayload)
	req := httptest.NewRequest("POST", "/api/v1/whistleblower/reports", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for anonymous submission, got %d: %s", rec.Code, rec.Body.String())
	}

	var submitResp struct {
		Report              backendWB.WhistleblowerReport `json:"report"`
		SecretTrackingToken string                        `json:"secret_tracking_token"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &submitResp)

	if submitResp.Report.ReportNumber != "WB-2026-00001" || submitResp.SecretTrackingToken == "" {
		t.Fatalf("Unexpected submission response: %+v", submitResp)
	}

	// 2. Lookup Report with Token
	lookupPayload := map[string]interface{}{
		"tenant_id":          "tenant-wb-http",
		"raw_tracking_token": submitResp.SecretTrackingToken,
	}
	lookupBody, _ := json.Marshal(lookupPayload)
	reqLookup := httptest.NewRequest("POST", "/api/v1/whistleblower/reports/lookup", bytes.NewReader(lookupBody))
	reqLookup.Header.Set("Content-Type", "application/json")
	recLookup := httptest.NewRecorder()
	router.ServeHTTP(recLookup, reqLookup)

	if recLookup.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for report lookup, got %d: %s", recLookup.Code, recLookup.Body.String())
	}

	// 3. Post Anonymous Follow-up Message
	msgPayload := map[string]interface{}{
		"tenant_id":          "tenant-wb-http",
		"sender_type":        "ANONYMOUS_REPORTER",
		"message":            "Incident occurred around 01:30 AM specifically in room 408.",
		"raw_tracking_token": submitResp.SecretTrackingToken,
	}
	msgBody, _ := json.Marshal(msgPayload)
	reqMsg := httptest.NewRequest("POST", "/api/v1/whistleblower/reports/"+submitResp.Report.ID+"/messages", bytes.NewReader(msgBody))
	reqMsg.Header.Set("Content-Type", "application/json")
	recMsg := httptest.NewRecorder()
	router.ServeHTTP(recMsg, reqMsg)

	if recMsg.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for anonymous reply, got %d: %s", recMsg.Code, recMsg.Body.String())
	}

	// 4. Committee Assigns Investigator
	assignPayload := map[string]interface{}{
		"tenant_id":  "tenant-wb-http",
		"officer_id": "anti-ragging-head-1",
		"actor_id":   "dean-student-affairs",
	}
	assignBody, _ := json.Marshal(assignPayload)
	reqAssign := httptest.NewRequest("POST", "/api/v1/whistleblower/committee/reports/"+submitResp.Report.ID+"/assign", bytes.NewReader(assignBody))
	reqAssign.Header.Set("Content-Type", "application/json")
	recAssign := httptest.NewRecorder()
	router.ServeHTTP(recAssign, reqAssign)

	if recAssign.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for committee assignment, got %d: %s", recAssign.Code, recAssign.Body.String())
	}

	// 5. Committee Resolves Report
	resolvePayload := map[string]interface{}{
		"tenant_id":        "tenant-wb-http",
		"status":           "RESOLVED",
		"findings_summary": "Disciplinary committee summoned students. Warning issued and hostel room changed.",
		"actor_id":         "anti-ragging-head-1",
	}
	resolveBody, _ := json.Marshal(resolvePayload)
	reqResolve := httptest.NewRequest("POST", "/api/v1/whistleblower/committee/reports/"+submitResp.Report.ID+"/status", bytes.NewReader(resolveBody))
	reqResolve.Header.Set("Content-Type", "application/json")
	recResolve := httptest.NewRecorder()
	router.ServeHTTP(recResolve, reqResolve)

	if recResolve.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for report resolution, got %d: %s", recResolve.Code, recResolve.Body.String())
	}
}
