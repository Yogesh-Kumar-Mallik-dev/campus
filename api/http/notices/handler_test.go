/**
 * BLOCK_API_NOTICES_HANDLER_TEST_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   HTTP transport tests for notices, broadcasts, and read receipts.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package notices

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"campus/backend/notices"
)

func setupNoticesHTTPServer() (*Handler, *notices.Service, *chi.Mux) {
	noticeRepo := notices.NewMockNoticeRepository()
	attRepo := notices.NewMockAttachmentRepository()
	ackRepo := notices.NewMockAcknowledgementRepository()

	svc := notices.NewService(noticeRepo, attRepo, ackRepo, nil)
	handler := NewHandler(svc)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	return handler, svc, r
}

func TestHTTP_CreateAndPublishNotice(t *testing.T) {
	_, _, r := setupNoticesHTTPServer()

	payload := map[string]interface{}{
		"tenant_id":       "ten_http_notices",
		"author_id":       "usr_registrar",
		"title":           "Graduation Convocation 2026",
		"content":         "The 15th Annual Convocation will be held at Main Auditorium.",
		"category":        "EVENTS",
		"priority":        "HIGH",
		"target_audience": "ALL",
		"is_pinned":       true,
		"publish_now":     false,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for notice, got %d: %s", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); loc == "" {
		t.Fatalf("expected Location header in response")
	}

	var resp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	n := resp["notice"].(map[string]interface{})
	noticeID := n["id"].(string)

	// Publish Notice
	req = httptest.NewRequest(http.MethodPost, "/api/v1/notices/"+noticeID+"/publish", bytes.NewReader([]byte(`{"tenant_id":"ten_http_notices","author_id":"usr_registrar"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for publish, got %d: %s", rec.Code, rec.Body.String())
	}

	// Acknowledge Notice
	ackPayload := map[string]interface{}{
		"tenant_id": "ten_http_notices",
		"user_id":   "usr_student_diya",
	}
	body, _ = json.Marshal(ackPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/notices/"+noticeID+"/acknowledge", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for acknowledge, got %d: %s", rec.Code, rec.Body.String())
	}
}
