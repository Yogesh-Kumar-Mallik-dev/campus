/**
 * BLOCK_API_PROBLEM_DETAILS_002
 * Purpose: Authoritative RFC 7807 Problem Details response builder with trace ID and status helpers.
 * Invariants: Content-Type application/problem+json, RFC 9110 status codes, trace ID propagation.
 * Domain:  API Gateway & Transport Layer
 */

package problem

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// InvalidParam records field-level semantic validation failures.
type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// ProblemDetails represents an RFC 7807 error payload.
type ProblemDetails struct {
	Type          string          `json:"type"`
	Title         string          `json:"title"`
	Status        int             `json:"status"`
	Detail        string          `json:"detail"`
	Instance      string          `json:"instance,omitempty"`
	Code          string          `json:"code,omitempty"`
	InvalidParams []*InvalidParam `json:"invalid_params,omitempty"`
	TraceID       string          `json:"trace_id,omitempty"`
}

// Write emits an RFC 7807 formatted JSON response with trace ID from request context.
func Write(w http.ResponseWriter, r *http.Request, status int, title, detail, code string, invalidParams []*InvalidParam) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	instance := ""
	traceID := ""
	if r != nil {
		instance = r.URL.Path
		if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
			traceID = reqID
		} else if tid, ok := r.Context().Value("trace_id").(string); ok {
			traceID = tid
		}
	}

	p := ProblemDetails{
		Type:          "https://campus.institution.edu/errors/" + code,
		Title:         title,
		Status:        status,
		Detail:        detail,
		Instance:      instance,
		Code:          code,
		InvalidParams: invalidParams,
		TraceID:       traceID,
	}

	_ = json.NewEncoder(w).Encode(p)
}

// BadRequest emits a 400 Bad Request problem.
func BadRequest(w http.ResponseWriter, r *http.Request, detail string, code string) {
	if code == "" {
		code = "BAD_REQUEST"
	}
	Write(w, r, http.StatusBadRequest, "Bad Request", detail, code, nil)
}

// Unauthorized emits a 401 Unauthorized problem.
func Unauthorized(w http.ResponseWriter, r *http.Request, detail string, code string) {
	if code == "" {
		code = "UNAUTHORIZED"
	}
	Write(w, r, http.StatusUnauthorized, "Unauthorized", detail, code, nil)
}

// Forbidden emits a 403 Forbidden problem.
func Forbidden(w http.ResponseWriter, r *http.Request, detail string, code string) {
	if code == "" {
		code = "FORBIDDEN"
	}
	Write(w, r, http.StatusForbidden, "Forbidden", detail, code, nil)
}

// NotFound emits a 404 Not Found problem.
func NotFound(w http.ResponseWriter, r *http.Request, detail string, code string) {
	if code == "" {
		code = "NOT_FOUND"
	}
	Write(w, r, http.StatusNotFound, "Not Found", detail, code, nil)
}

// Conflict emits a 409 Conflict problem.
func Conflict(w http.ResponseWriter, r *http.Request, detail string, code string) {
	if code == "" {
		code = "CONFLICT"
	}
	Write(w, r, http.StatusConflict, "Conflict", detail, code, nil)
}

// PreconditionFailed emits a 412 Precondition Failed problem (Optimistic Concurrency Control failure).
func PreconditionFailed(w http.ResponseWriter, r *http.Request, detail string, code string) {
	if code == "" {
		code = "PRECONDITION_FAILED"
	}
	Write(w, r, http.StatusPreconditionFailed, "Precondition Failed", detail, code, nil)
}

// PayloadTooLarge emits a 413 Payload Too Large problem.
func PayloadTooLarge(w http.ResponseWriter, r *http.Request, detail string, code string) {
	if code == "" {
		code = "PAYLOAD_TOO_LARGE"
	}
	Write(w, r, http.StatusRequestEntityTooLarge, "Payload Too Large", detail, code, nil)
}

// UnsupportedMediaType emits a 415 Unsupported Media Type problem.
func UnsupportedMediaType(w http.ResponseWriter, r *http.Request, detail string, code string) {
	if code == "" {
		code = "UNSUPPORTED_MEDIA_TYPE"
	}
	Write(w, r, http.StatusUnsupportedMediaType, "Unsupported Media Type", detail, code, nil)
}

// UnprocessableEntity emits a 422 Unprocessable Content problem with validation params.
func UnprocessableEntity(w http.ResponseWriter, r *http.Request, detail string, code string, invalidParams []*InvalidParam) {
	if code == "" {
		code = "VALIDATION_FAILED"
	}
	Write(w, r, http.StatusUnprocessableEntity, "Unprocessable Content", detail, code, invalidParams)
}

// TooManyRequests emits a 429 Too Many Requests problem with Retry-After header.
func TooManyRequests(w http.ResponseWriter, r *http.Request, detail string, retryAfterSeconds int) {
	if retryAfterSeconds > 0 {
		w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfterSeconds))
	}
	Write(w, r, http.StatusTooManyRequests, "Too Many Requests", detail, "RATE_LIMIT_EXCEEDED", nil)
}

// InternalServerError emits a 500 Internal Server Error problem.
func InternalServerError(w http.ResponseWriter, r *http.Request, detail string) {
	Write(w, r, http.StatusInternalServerError, "Internal Server Error", detail, "INTERNAL_SERVER_ERROR", nil)
}

// ServiceUnavailable emits a 503 Service Unavailable problem.
func ServiceUnavailable(w http.ResponseWriter, r *http.Request, detail string) {
	Write(w, r, http.StatusServiceUnavailable, "Service Unavailable", detail, "SERVICE_UNAVAILABLE", nil)
}

// GatewayTimeout emits a 504 Gateway Timeout problem.
func GatewayTimeout(w http.ResponseWriter, r *http.Request, detail string) {
	Write(w, r, http.StatusGatewayTimeout, "Gateway Timeout", detail, "GATEWAY_TIMEOUT", nil)
}
