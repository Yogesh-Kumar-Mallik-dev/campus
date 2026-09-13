/**
 * BLOCK_API_PROBLEM_DETAILS_001
 * Purpose: RFC 7807 Problem Details for HTTP APIs (application/problem+json).
 * Invariants: Content-Type application/problem+json, standardized error envelopes.
 * Domain:  API Gateway & Transport Layer
 */

package problem

import (
	"encoding/json"
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
}

// Write emits an RFC 7807 formatted JSON response.
func Write(w http.ResponseWriter, r *http.Request, status int, title, detail, code string, invalidParams []*InvalidParam) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	instance := ""
	if r != nil {
		instance = r.URL.Path
	}

	p := ProblemDetails{
		Type:          "https://campus.institution.edu/errors/" + code,
		Title:         title,
		Status:        status,
		Detail:        detail,
		Instance:      instance,
		Code:          code,
		InvalidParams: invalidParams,
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

// UnprocessableEntity emits a 422 Unprocessable Content problem with field errors.
func UnprocessableEntity(w http.ResponseWriter, r *http.Request, detail string, code string, invalidParams []*InvalidParam) {
	if code == "" {
		code = "VALIDATION_FAILED"
	}
	Write(w, r, http.StatusUnprocessableEntity, "Unprocessable Content", detail, code, invalidParams)
}

// InternalServerError emits a 500 Internal Server Error problem.
func InternalServerError(w http.ResponseWriter, r *http.Request, detail string) {
	Write(w, r, http.StatusInternalServerError, "Internal Server Error", detail, "INTERNAL_SERVER_ERROR", nil)
}
