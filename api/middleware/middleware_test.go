/**
 * BLOCK_API_MIDDLEWARE_TESTS_001
 * Purpose: Unit tests verifying gateway middleware invariants and RFC 7807 problem responses.
 * Domain:  API Gateway & Transport Layer
 */

package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := r.Context().Value(TraceIDKey).(string)
		if !ok || traceID == "" {
			t.Error("Expected TraceIDKey in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	respID := rec.Header().Get("X-Request-ID")
	if respID == "" || !strings.HasPrefix(respID, "req_") {
		t.Errorf("Expected X-Request-ID header starting with req_, got: %s", respID)
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	handler := SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "LeakyServer/1.0")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	h := rec.Header()
	if h.Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options: nosniff")
	}
	if h.Get("X-Frame-Options") != "DENY" {
		t.Error("Expected X-Frame-Options: DENY")
	}
	if h.Get("Server") != "" {
		t.Error("Expected Server header to be stripped")
	}
}

func TestBodyLimitMiddleware_RejectsLargePayload(t *testing.T) {
	limitMiddleware := BodyLimitMiddleware(100) // 100 bytes limit
	handler := limitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	largePayload := bytes.Repeat([]byte("A"), 200)
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(largePayload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("Expected 413 Payload Too Large, got: %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/problem+json" {
		t.Error("Expected RFC 7807 application/problem+json response")
	}
}

func TestEnforceJSONContentType(t *testing.T) {
	handler := EnforceJSONContentType(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Invalid content type on POST with body
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("text=data"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("Expected 415 Unsupported Media Type, got: %d", rec.Code)
	}

	// Valid content type
	reqValid := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
	reqValid.Header.Set("Content-Type", "application/json")
	recValid := httptest.NewRecorder()

	handler.ServeHTTP(recValid, reqValid)

	if recValid.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for valid JSON content type, got: %d", recValid.Code)
	}
}
