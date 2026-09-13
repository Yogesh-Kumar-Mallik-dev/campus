/**
 * BLOCK_API_HEALTH_TESTS_001
 * Purpose: Unit tests verifying Kubernetes health probe invariants (/healthz/live, /healthz/ready).
 * Domain:  API Gateway & Transport Layer
 */

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHealthProbes_Liveness(t *testing.T) {
	r := chi.NewRouter()
	RegisterHealthProbes(r, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got: %d", rec.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode liveness response: %v", err)
	}

	if resp["status"] != "UP" {
		t.Errorf("Expected status UP, got: %v", resp["status"])
	}
}

func TestHealthProbes_Readiness_Healthy(t *testing.T) {
	r := chi.NewRouter()
	RegisterHealthProbes(r, func() error {
		return nil // DB is healthy
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got: %d", rec.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode readiness response: %v", err)
	}

	if resp["status"] != "READY" {
		t.Errorf("Expected status READY, got: %v", resp["status"])
	}
}

func TestHealthProbes_Readiness_Degraded(t *testing.T) {
	r := chi.NewRouter()
	RegisterHealthProbes(r, func() error {
		return errors.New("database connection pool exhausted")
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("Expected 503 Service Unavailable, got: %d", rec.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode degraded response: %v", err)
	}

	if resp["status"] != "DEGRADED" {
		t.Errorf("Expected status DEGRADED, got: %v", resp["status"])
	}
	if resp["error"] != "database connection pool exhausted" {
		t.Errorf("Expected error string in response, got: %v", resp["error"])
	}
}
