/**
 * BLOCK_API_BILLING_HANDLER_TEST_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   HTTP transport tests for fee structures, invoices, idempotency payments, and balance clearance.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package billing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/backend/billing"
)

func setupBillingHTTPServer() (*Handler, *billing.Service, *chi.Mux) {
	feeRepo := billing.NewMockFeeRepository()
	invRepo := billing.NewMockInvoiceRepository()
	payRepo := billing.NewMockPaymentRepository()
	ledgerRepo := billing.NewMockLedgerRepository()
	seqRepo := billing.NewMockSequenceRepository()

	svc := billing.NewService(feeRepo, invRepo, payRepo, ledgerRepo, seqRepo, nil)
	handler := NewHandler(svc)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	return handler, svc, r
}

func TestHTTP_CreateFeeStructure(t *testing.T) {
	_, _, r := setupBillingHTTPServer()

	payload := map[string]interface{}{
		"tenant_id":     "ten_http_bill",
		"program_id":    "prog_cs",
		"academic_year": "2026-2027",
		"semester":      1,
		"name":          "Semester 1 Standard Fee",
		"due_date":      "2026-10-31",
		"items": []map[string]interface{}{
			{"category": "TUITION", "name": "Tuition", "amount": 6000000},
			{"category": "EXAMINATION", "name": "Exam Fee", "amount": 400000},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/fee-structures", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); loc == "" {
		t.Fatalf("expected Location header in response")
	}

	var resp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	fs := resp["fee_structure"].(map[string]interface{})
	if int(fs["total_amount"].(float64)) != 6400000 {
		t.Fatalf("expected total 6400000, got %v", fs["total_amount"])
	}
}

func TestHTTP_InvoiceAndPaymentWorkflow(t *testing.T) {
	_, _, r := setupBillingHTTPServer()

	// 1. Create Invoice
	invPayload := map[string]interface{}{
		"tenant_id":       "ten_http_bill",
		"student_id":      "stu_rohan",
		"academic_year":   "2026-2027",
		"semester":        1,
		"due_date":        "2026-10-31",
		"discount_amount": 0,
		"tax_amount":      0,
		"custom_items": []map[string]interface{}{
			{"category": "TUITION", "name": "Semester Tuition", "amount": 5000000},
		},
	}
	body, _ := json.Marshal(invPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/invoices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for invoice, got %d: %s", rec.Code, rec.Body.String())
	}

	var invResp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&invResp)
	invoiceObj := invResp["invoice"].(map[string]interface{})
	invoiceID := invoiceObj["id"].(string)

	// 2. Initiate Payment with Idempotency Key
	idempKey := "idemp_http_txn_001"
	payPayload := map[string]interface{}{
		"tenant_id":       "ten_http_bill",
		"invoice_id":      invoiceID,
		"student_id":      "stu_rohan",
		"amount":          5000000,
		"method":          "ONLINE_GATEWAY",
		"idempotency_key": idempKey,
	}
	body, _ = json.Marshal(payPayload)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/billing/payments/initiate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for payment init, got %d: %s", rec.Code, rec.Body.String())
	}

	var payResp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&payResp)
	txObj := payResp["transaction"].(map[string]interface{})
	txID := txObj["id"].(string)

	// 3. Process Payment Success Callback
	callbackPayload := map[string]interface{}{
		"tenant_id":          "ten_http_bill",
		"transaction_id":     txID,
		"gateway_payment_id": "pay_razorpay_9999",
		"paid_at":            time.Now().UTC(),
	}
	body, _ = json.Marshal(callbackPayload)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/billing/payments/callback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for payment callback, got %d: %s", rec.Code, rec.Body.String())
	}

	// 4. Verify Student Balance is now 0
	req = httptest.NewRequest(http.MethodGet, "/api/v1/billing/balances/students/stu_rohan?tenant_id=ten_http_bill", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for balance query, got %d", rec.Code)
	}

	var balResp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&balResp)
	if int(balResp["pending_balance"].(float64)) != 0 {
		t.Fatalf("expected pending balance 0 after full payment, got %v", balResp["pending_balance"])
	}
}
