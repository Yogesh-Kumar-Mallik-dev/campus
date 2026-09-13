/**
 * BLOCK_API_AUDIT_HANDLER_TEST_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   HTTP transport integration tests verifying REST endpoints, RFC 7807 responses, and headers.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/backend/audit"
)

func setupTestRouter() (chi.Router, audit.Service, *audit.MockRepository) {
	repo := audit.NewMockRepository()
	hasher := audit.NewSHA256Hasher()
	svc := audit.NewService(repo, hasher, nil)

	r := chi.NewRouter()
	RegisterRoutes(r, svc)

	return r, svc, repo
}

func TestHTTP_QueryLogs_Success(t *testing.T) {
	r, svc, _ := setupTestRouter()
	ctx := context.Background()

	tenantID := "ten_audit_http"
	_, _ = svc.Record(ctx, audit.RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "auth:login:success",
		ResourceType: "user",
		Status:       audit.StatusSuccess,
	})
	_, _ = svc.Record(ctx, audit.RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "hostel:room:allocated",
		ResourceType: "room",
		Status:       audit.StatusSuccess,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit/logs?tenant_id=ten_audit_http&limit=10", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got: %d", rec.Code)
	}

	totalHeader := rec.Header().Get("X-Total-Count")
	if totalHeader != "2" {
		t.Errorf("Expected X-Total-Count: 2, got: %s", totalHeader)
	}

	var resp struct {
		Data   []audit.AuditLog `json:"data"`
		Total  int64            `json:"total"`
		Limit  int              `json:"limit"`
		Offset int              `json:"offset"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Total != 2 || len(resp.Data) != 2 {
		t.Errorf("Expected 2 records in data array, got %d (total: %d)", len(resp.Data), resp.Total)
	}
}

func TestHTTP_GetLogByID_FoundAndNotFound(t *testing.T) {
	r, svc, _ := setupTestRouter()
	ctx := context.Background()
	tenantID := "ten_get_id"

	entry, _ := svc.Record(ctx, audit.RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "library:book:issued",
		ResourceType: "book",
		Status:       audit.StatusSuccess,
	})

	// Found
	reqFound := httptest.NewRequest(http.MethodGet, "/api/v1/audit/logs/"+entry.ID+"?tenant_id="+tenantID, nil)
	recFound := httptest.NewRecorder()
	r.ServeHTTP(recFound, reqFound)

	if recFound.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got: %d", recFound.Code)
	}

	// Not Found
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/v1/audit/logs/aud_nonexistent?tenant_id="+tenantID, nil)
	recNotFound := httptest.NewRecorder()
	r.ServeHTTP(recNotFound, reqNotFound)

	if recNotFound.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 Not Found, got: %d", recNotFound.Code)
	}
	if recNotFound.Header().Get("Content-Type") != "application/problem+json" {
		t.Error("Expected application/problem+json content-type")
	}
}

func TestHTTP_VerifyIntegrity(t *testing.T) {
	r, svc, _ := setupTestRouter()
	ctx := context.Background()
	tenantID := "ten_verify_api"

	_, _ = svc.Record(ctx, audit.RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "fee:invoice:created",
		ResourceType: "invoice",
		Status:       audit.StatusSuccess,
	})
	_, _ = svc.Record(ctx, audit.RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "fee:invoice:paid",
		ResourceType: "invoice",
		Status:       audit.StatusSuccess,
	})

	payload := map[string]any{
		"tenant_id": tenantID,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/audit/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got: %d", rec.Code)
	}

	var resp struct {
		Data audit.VerificationResult `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode verification response: %v", err)
	}

	if !resp.Data.IsChainIntact || resp.Data.TotalVerified != 2 {
		t.Errorf("Expected chain intact with 2 verified, got: %+v", resp.Data)
	}
}

func TestHTTP_CreateCheckpoint_Success(t *testing.T) {
	r, svc, _ := setupTestRouter()
	ctx := context.Background()
	tenantID := "ten_checkpoint_api"

	_, _ = svc.Record(ctx, audit.RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "attendance:mark",
		ResourceType: "attendance",
		Status:       audit.StatusSuccess,
	})

	body, _ := json.Marshal(map[string]string{"tenant_id": tenantID})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/audit/checkpoints", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got: %d", rec.Code)
	}

	var resp struct {
		Data audit.AuditVerificationCheckpoint `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode checkpoint response: %v", err)
	}

	if resp.Data.RecordCount != 1 || !resp.Data.IsVerified {
		t.Errorf("Checkpoint data mismatch: %+v", resp.Data)
	}
}

func TestHTTP_ComplianceReport_Success(t *testing.T) {
	r, svc, _ := setupTestRouter()
	ctx := context.Background()
	tenantID := "ten_report_api"

	_, _ = svc.Record(ctx, audit.RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "mentorship:session:logged",
		ResourceType: "mentorship_log",
		Status:       audit.StatusSuccess,
	})

	from := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit/compliance-reports?tenant_id="+tenantID+"&framework=NAAC&from="+from+"&to="+to, nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got: %d", rec.Code)
	}

	var resp struct {
		Data audit.ComplianceReportSummary `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode compliance report: %v", err)
	}

	if resp.Data.TotalEvents != 1 || resp.Data.Framework != "NAAC" {
		t.Errorf("Compliance report mismatch: %+v", resp.Data)
	}
}
