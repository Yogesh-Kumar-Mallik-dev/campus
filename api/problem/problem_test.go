/**
 * BLOCK_API_PROBLEM_TESTS_001
 * Purpose: Unit tests verifying RFC 7807 problem details payloads, headers, and trace ID propagation.
 * Domain:  API Gateway & Transport Layer
 */

package problem

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProblem_BadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/test", nil)
	req.Header.Set("X-Request-ID", "req_test_123")
	rec := httptest.NewRecorder()

	BadRequest(rec, req, "Invalid input payload", "INVALID_PAYLOAD")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request, got: %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/problem+json" {
		t.Errorf("Expected application/problem+json content-type, got: %s", rec.Header().Get("Content-Type"))
	}

	var p ProblemDetails
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
		t.Fatalf("Failed to decode problem JSON: %v", err)
	}

	if p.Status != 400 || p.Code != "INVALID_PAYLOAD" || p.Detail != "Invalid input payload" {
		t.Errorf("ProblemDetails payload mismatch: %+v", p)
	}
	if p.TraceID != "req_test_123" {
		t.Errorf("Expected trace_id req_test_123, got: %s", p.TraceID)
	}
	if p.Instance != "/api/v1/test" {
		t.Errorf("Expected instance /api/v1/test, got: %s", p.Instance)
	}
}

func TestProblem_Unauthorized_And_Forbidden(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/secure", nil)
	rec401 := httptest.NewRecorder()
	Unauthorized(rec401, req, "Missing bearer token", "UNAUTHORIZED")
	if rec401.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got: %d", rec401.Code)
	}

	rec403 := httptest.NewRecorder()
	Forbidden(rec403, req, "Insufficient capability permissions", "FORBIDDEN")
	if rec403.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got: %d", rec403.Code)
	}
}

func TestProblem_NotFound_And_Conflict(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/usr_999", nil)
	rec404 := httptest.NewRecorder()
	NotFound(rec404, req, "User not found", "USER_NOT_FOUND")
	if rec404.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got: %d", rec404.Code)
	}

	rec409 := httptest.NewRecorder()
	Conflict(rec409, req, "Duplicate email in tenant", "EMAIL_EXISTS")
	if rec409.Code != http.StatusConflict {
		t.Errorf("Expected 409 Conflict, got: %d", rec409.Code)
	}
}

func TestProblem_PreconditionFailed_And_PayloadTooLarge(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/resources/1", nil)
	rec412 := httptest.NewRecorder()
	PreconditionFailed(rec412, req, "ETag precondition mismatch", "OCC_CONFLICT")
	if rec412.Code != http.StatusPreconditionFailed {
		t.Errorf("Expected 412 Precondition Failed, got: %d", rec412.Code)
	}

	rec413 := httptest.NewRecorder()
	PayloadTooLarge(rec413, req, "Payload exceeds 250KB limit", "PAYLOAD_TOO_LARGE")
	if rec413.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected 413 Payload Too Large, got: %d", rec413.Code)
	}
}

func TestProblem_UnprocessableEntity_WithInvalidParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", nil)
	rec := httptest.NewRecorder()

	params := []*InvalidParam{
		{Name: "email", Reason: "Must be valid institutional email"},
		{Name: "rollNumber", Reason: "Cannot be empty"},
	}

	UnprocessableEntity(rec, req, "Validation failed", "VALIDATION_FAILED", params)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Expected 422 Unprocessable Content, got: %d", rec.Code)
	}

	var p ProblemDetails
	_ = json.NewDecoder(rec.Body).Decode(&p)

	if len(p.InvalidParams) != 2 {
		t.Fatalf("Expected 2 invalid params, got: %d", len(p.InvalidParams))
	}
	if p.InvalidParams[0].Name != "email" || p.InvalidParams[1].Name != "rollNumber" {
		t.Errorf("InvalidParams mismatch: %+v", p.InvalidParams)
	}
}

func TestProblem_TooManyRequests_WithRetryAfter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/news", nil)
	rec := httptest.NewRecorder()

	TooManyRequests(rec, req, "Rate limit of 100 req/min exceeded", 60)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected 429 Too Many Requests, got: %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") != "60" {
		t.Errorf("Expected Retry-After: 60 header, got: %s", rec.Header().Get("Retry-After"))
	}
}

func TestProblem_ServerErrors_And_TraceContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", "ctx_trace_abc")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/data", nil).WithContext(ctx)

	rec500 := httptest.NewRecorder()
	InternalServerError(rec500, req, "Database connection failure")
	if rec500.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Server Error, got: %d", rec500.Code)
	}

	var p500 ProblemDetails
	_ = json.NewDecoder(rec500.Body).Decode(&p500)
	if p500.TraceID != "ctx_trace_abc" {
		t.Errorf("Expected trace ID from context ctx_trace_abc, got: %s", p500.TraceID)
	}

	rec503 := httptest.NewRecorder()
	ServiceUnavailable(rec503, req, "Service is in maintenance mode")
	if rec503.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 Service Unavailable, got: %d", rec503.Code)
	}

	rec504 := httptest.NewRecorder()
	GatewayTimeout(rec504, req, "Upstream timed out")
	if rec504.Code != http.StatusGatewayTimeout {
		t.Errorf("Expected 504 Gateway Timeout, got: %d", rec504.Code)
	}
}
