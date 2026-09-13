/**
 * BLOCK_API_BILLING_HANDLER_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for fee structures, student invoices, idempotency payments, ledger, and balance clearance.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package billing

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/billing"
)

type Handler struct {
	service *billing.Service
}

func NewHandler(service *billing.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/billing", func(r chi.Router) {
		// Fee Templates & Schedules
		r.Post("/fee-structures", h.CreateFeeStructure)

		// Student Invoices
		r.Post("/invoices", h.GenerateInvoice)
		r.Get("/invoices", h.ListInvoices)
		r.Get("/invoices/{id}", h.GetInvoice)

		// Payment Transactions & Gateways
		r.Post("/payments/initiate", h.InitiatePayment)
		r.Post("/payments/callback", h.ProcessPaymentCallback)
		r.Post("/payments/{id}/refund", h.ProcessRefund)

		// Student Balance
		r.Get("/balances/students/{student_id}", h.GetStudentBalance)
	})
}

// CreateFeeStructure handles POST /api/v1/billing/fee-structures
func (h *Handler) CreateFeeStructure(w http.ResponseWriter, r *http.Request) {
	var cmd billing.CreateFeeStructureCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	fs, err := h.service.CreateFeeStructure(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v1/billing/fee-structures/"+fs.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"fee_structure": fs,
	})
}

// GenerateInvoice handles POST /api/v1/billing/invoices
func (h *Handler) GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	var cmd billing.GenerateInvoiceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	inv, err := h.service.GenerateInvoice(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v1/billing/invoices/"+inv.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"invoice": inv,
	})
}

// ListInvoices handles GET /api/v1/billing/invoices
func (h *Handler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	sem, _ := strconv.Atoi(r.URL.Query().Get("semester"))

	filter := billing.InvoiceFilter{
		TenantID:     tenantID,
		StudentID:    r.URL.Query().Get("student_id"),
		AcademicYear: r.URL.Query().Get("academic_year"),
		Semester:     sem,
		Limit:        limit,
		Offset:       offset,
	}

	if st := r.URL.Query().Get("status"); st != "" {
		status := billing.InvoiceStatus(strings.ToUpper(st))
		filter.Status = &status
	}

	invoices, total, err := h.service.ListInvoices(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"invoices": invoices,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetInvoice handles GET /api/v1/billing/invoices/{id}
func (h *Handler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	inv, err := h.service.GetInvoice(r.Context(), tenantID, invoiceID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"invoice": inv,
	})
}

// InitiatePayment handles POST /api/v1/billing/payments/initiate
func (h *Handler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	var cmd billing.InitiatePaymentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	// Capture idempotency header if not in payload
	if cmd.IdempotencyKey == nil || *cmd.IdempotencyKey == "" {
		if idempHdr := r.Header.Get("Idempotency-Key"); idempHdr != "" {
			cmd.IdempotencyKey = &idempHdr
		}
	}

	tx, err := h.service.InitiatePayment(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v1/billing/payments/"+tx.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction": tx,
	})
}

// ProcessPaymentCallback handles POST /api/v1/billing/payments/callback
func (h *Handler) ProcessPaymentCallback(w http.ResponseWriter, r *http.Request) {
	var cmd billing.ProcessPaymentSuccessCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}
	if cmd.PaidAt.IsZero() {
		cmd.PaidAt = time.Now().UTC()
	}

	tx, err := h.service.ProcessPaymentSuccess(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction": tx,
	})
}

// ProcessRefund handles POST /api/v1/billing/payments/{id}/refund
func (h *Handler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "id")
	var body struct {
		TenantID string `json:"tenant_id"`
		Reason   string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.TenantID == "" {
		body.TenantID = r.Header.Get("X-Tenant-ID")
	}

	tx, err := h.service.ProcessRefund(r.Context(), body.TenantID, transactionID, body.Reason)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction": tx,
	})
}

// GetStudentBalance handles GET /api/v1/billing/balances/students/{student_id}
func (h *Handler) GetStudentBalance(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "student_id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	balance, err := h.service.GetStudentBalance(r.Context(), tenantID, studentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"student_id":      studentID,
		"pending_balance": balance,
	})
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, billing.ErrFeeStructureNotFound),
		errors.Is(err, billing.ErrInvoiceNotFound),
		errors.Is(err, billing.ErrTransactionNotFound),
		errors.Is(err, billing.ErrAccountNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")
	case errors.Is(err, billing.ErrTenantRequired),
		errors.Is(err, billing.ErrInvalidInput),
		errors.Is(err, billing.ErrInvalidAmount):
		problem.BadRequest(w, r, err.Error(), "INVALID_INPUT")
	case errors.Is(err, billing.ErrInvoiceAlreadyPaid),
		errors.Is(err, billing.ErrInvoiceCancelled),
		errors.Is(err, billing.ErrOverpaymentNotAllowed),
		errors.Is(err, billing.ErrInvalidStateTransition),
		errors.Is(err, billing.ErrLedgerUnbalanced):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)
	case errors.Is(err, billing.ErrDuplicateInvoiceNumber),
		errors.Is(err, billing.ErrDuplicateTxnRef),
		errors.Is(err, billing.ErrIdempotencyConflict):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")
	default:
		var domErr *billing.DomainError
		if errors.As(err, &domErr) {
			problem.UnprocessableEntity(w, r, domErr.Message, domErr.Code, nil)
			return
		}
		problem.InternalServerError(w, r, "an unexpected billing error occurred")
	}
}
