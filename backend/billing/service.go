/**
 * BLOCK_BILLING_SERVICE_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   Core business orchestration for fee templates, invoice lifecycle, idempotency payments, double-entry ledger, and receipt generation.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package billing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"campus/backend/audit"
)

type FeeItemInput struct {
	Category   FeeCategory `json:"category"`
	Name       string      `json:"name"`
	Amount     int         `json:"amount"`
	IsOptional bool        `json:"is_optional"`
}

type CreateFeeStructureCommand struct {
	TenantID     string         `json:"tenant_id"`
	ProgramID    string         `json:"program_id"`
	AcademicYear string         `json:"academic_year"`
	Semester     int            `json:"semester"`
	Name         string         `json:"name"`
	DueDate      string         `json:"due_date"`
	Currency     string         `json:"currency"`
	Items        []FeeItemInput `json:"items"`
}

type GenerateInvoiceCommand struct {
	TenantID       string         `json:"tenant_id"`
	StudentID      string         `json:"student_id"`
	FeeStructureID *string        `json:"fee_structure_id,omitempty"`
	AcademicYear   string         `json:"academic_year"`
	Semester       int            `json:"semester"`
	DueDate        string         `json:"due_date"`
	DiscountAmount int            `json:"discount_amount"`
	TaxAmount      int            `json:"tax_amount"`
	CustomItems    []FeeItemInput `json:"custom_items,omitempty"`
	Notes          string         `json:"notes,omitempty"`
}

type InitiatePaymentCommand struct {
	TenantID       string        `json:"tenant_id"`
	InvoiceID      string        `json:"invoice_id"`
	StudentID      string        `json:"student_id"`
	Amount         int           `json:"amount"`
	Currency       string        `json:"currency"`
	Method         PaymentMethod `json:"method"`
	GatewayName    string        `json:"gateway_name"`
	IdempotencyKey *string       `json:"idempotency_key,omitempty"`
}

type ProcessPaymentSuccessCommand struct {
	TenantID         string    `json:"tenant_id"`
	TransactionID    string    `json:"transaction_id"`
	GatewayPaymentID *string   `json:"gateway_payment_id,omitempty"`
	PaidAt           time.Time `json:"paid_at"`
}

type Service struct {
	feeRepo     FeeRepository
	invoiceRepo InvoiceRepository
	paymentRepo PaymentRepository
	ledgerRepo  LedgerRepository
	seqRepo     SequenceRepository
	auditSub    audit.Subscriber
}

func NewService(
	feeRepo FeeRepository,
	invoiceRepo InvoiceRepository,
	paymentRepo PaymentRepository,
	ledgerRepo LedgerRepository,
	seqRepo SequenceRepository,
	auditSub audit.Subscriber,
) *Service {
	return &Service{
		feeRepo:     feeRepo,
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
		ledgerRepo:  ledgerRepo,
		seqRepo:     seqRepo,
		auditSub:    auditSub,
	}
}

// CreateFeeStructure registers an institutional fee schedule template.
func (s *Service) CreateFeeStructure(ctx context.Context, cmd CreateFeeStructureCommand) (*FeeStructure, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(cmd.ProgramID) == "" || strings.TrimSpace(cmd.Name) == "" || cmd.Semester <= 0 {
		return nil, NewDomainError("INVALID_FEE_STRUCTURE", "program, semester, and name are required", ErrInvalidInput)
	}
	if cmd.Currency == "" {
		cmd.Currency = "INR"
	}

	totalAmount := 0
	items := make([]FeeStructureItem, 0, len(cmd.Items))
	now := time.Now().UTC()
	fsID := generateID("fs_")

	for _, it := range cmd.Items {
		if it.Amount <= 0 {
			return nil, ErrInvalidAmount
		}
		totalAmount += it.Amount
		items = append(items, FeeStructureItem{
			ID:             generateID("fsi_"),
			TenantID:       cmd.TenantID,
			FeeStructureID: fsID,
			Category:       it.Category,
			Name:           it.Name,
			Amount:         it.Amount,
			IsOptional:     it.IsOptional,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}

	fs := &FeeStructure{
		ID:           fsID,
		TenantID:     cmd.TenantID,
		ProgramID:    cmd.ProgramID,
		AcademicYear: cmd.AcademicYear,
		Semester:     cmd.Semester,
		Name:         cmd.Name,
		TotalAmount:  totalAmount,
		Currency:     cmd.Currency,
		DueDate:      cmd.DueDate,
		IsActive:     true,
		Items:        items,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.feeRepo.CreateFeeStructure(ctx, fs); err != nil {
		return nil, err
	}

	s.emitAudit(cmd.TenantID, "", "SYSTEM", "billing:fee_structure:created", "fee_structure", fs.ID, audit.StatusSuccess, map[string]interface{}{
		"program_id":   cmd.ProgramID,
		"semester":     cmd.Semester,
		"total_amount": totalAmount,
	})

	return fs, nil
}

// GenerateInvoice computes fee items, applies adjustments, and generates a student invoice.
func (s *Service) GenerateInvoice(ctx context.Context, cmd GenerateInvoiceCommand) (*StudentInvoice, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(cmd.StudentID) == "" || cmd.Semester <= 0 {
		return nil, NewDomainError("INVALID_INVOICE", "student_id and semester are required", ErrInvalidInput)
	}

	subtotal := 0
	invoiceID := generateID("inv_")
	now := time.Now().UTC()
	var items []InvoiceItem

	if len(cmd.CustomItems) > 0 {
		for _, cit := range cmd.CustomItems {
			if cit.Amount <= 0 {
				return nil, ErrInvalidAmount
			}
			subtotal += cit.Amount
			items = append(items, InvoiceItem{
				ID:        generateID("initm_"),
				TenantID:  cmd.TenantID,
				InvoiceID: invoiceID,
				Category:  cit.Category,
				Name:      cit.Name,
				Amount:    cit.Amount,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}
	} else if cmd.FeeStructureID != nil && *cmd.FeeStructureID != "" {
		fs, err := s.feeRepo.GetFeeStructureByID(ctx, cmd.TenantID, *cmd.FeeStructureID)
		if err != nil {
			return nil, err
		}
		for _, fit := range fs.Items {
			subtotal += fit.Amount
			items = append(items, InvoiceItem{
				ID:        generateID("initm_"),
				TenantID:  cmd.TenantID,
				InvoiceID: invoiceID,
				Category:  fit.Category,
				Name:      fit.Name,
				Amount:    fit.Amount,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}
	} else {
		return nil, NewDomainError("MISSING_ITEMS", "either custom items or a fee structure must be provided", ErrInvalidInput)
	}

	total := subtotal - cmd.DiscountAmount + cmd.TaxAmount
	if total < 0 {
		return nil, NewDomainError("INVALID_TOTAL", "invoice total cannot be negative after discount", ErrInvalidAmount)
	}

	seq, err := s.seqRepo.NextInvoiceSequence(ctx, cmd.TenantID, cmd.AcademicYear)
	if err != nil {
		return nil, err
	}
	invoiceNumber := fmt.Sprintf("INV-%s-%06d", strings.ReplaceAll(cmd.AcademicYear, "/", "-"), seq)

	invoice := &StudentInvoice{
		ID:             invoiceID,
		TenantID:       cmd.TenantID,
		StudentID:      cmd.StudentID,
		FeeStructureID: cmd.FeeStructureID,
		InvoiceNumber:  invoiceNumber,
		AcademicYear:   cmd.AcademicYear,
		Semester:       cmd.Semester,
		SubtotalAmount: subtotal,
		DiscountAmount: cmd.DiscountAmount,
		TaxAmount:      cmd.TaxAmount,
		TotalAmount:    total,
		PaidAmount:     0,
		BalanceAmount:  total,
		DueDate:        cmd.DueDate,
		Status:         InvoiceIssued,
		Notes:          cmd.Notes,
		Items:          items,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.invoiceRepo.CreateInvoice(ctx, invoice); err != nil {
		return nil, err
	}

	s.emitAudit(cmd.TenantID, cmd.StudentID, "SYSTEM", "billing:invoice:issued", "student_invoice", invoice.ID, audit.StatusSuccess, map[string]interface{}{
		"invoice_number": invoiceNumber,
		"total_amount":   total,
		"semester":       cmd.Semester,
	})

	return invoice, nil
}

// InitiatePayment starts a payment transaction with idempotency verification.
func (s *Service) InitiatePayment(ctx context.Context, cmd InitiatePaymentCommand) (*PaymentTransaction, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if cmd.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Idempotency check
	if cmd.IdempotencyKey != nil && strings.TrimSpace(*cmd.IdempotencyKey) != "" {
		existing, err := s.paymentRepo.GetTransactionByIdempotencyKey(ctx, cmd.TenantID, *cmd.IdempotencyKey)
		if err == nil && existing != nil {
			if existing.Amount != cmd.Amount || existing.InvoiceID != cmd.InvoiceID {
				return nil, ErrIdempotencyConflict
			}
			return existing, nil
		}
	}

	invoice, err := s.invoiceRepo.GetInvoiceByID(ctx, cmd.TenantID, cmd.InvoiceID)
	if err != nil {
		return nil, err
	}
	if invoice.Status == InvoicePaid {
		return nil, ErrInvoiceAlreadyPaid
	}
	if invoice.Status == InvoiceCancelled || invoice.Status == InvoiceWrittenOff {
		return nil, ErrInvoiceCancelled
	}
	if cmd.Amount > invoice.BalanceAmount {
		return nil, ErrOverpaymentNotAllowed
	}

	now := time.Now().UTC()
	yearStr := time.Now().UTC().Format("2006")
	seq, err := s.seqRepo.NextTransactionSequence(ctx, cmd.TenantID, yearStr)
	if err != nil {
		return nil, err
	}
	txnRef := fmt.Sprintf("TXN-%s-%06d", yearStr, seq)

	if cmd.Currency == "" {
		cmd.Currency = "INR"
	}
	if cmd.GatewayName == "" {
		cmd.GatewayName = "INTERNAL"
	}

	tx := &PaymentTransaction{
		ID:             generateID("txn_"),
		TenantID:       cmd.TenantID,
		InvoiceID:      cmd.InvoiceID,
		StudentID:      cmd.StudentID,
		TransactionRef: txnRef,
		GatewayName:    cmd.GatewayName,
		Amount:         cmd.Amount,
		Currency:       cmd.Currency,
		Method:         cmd.Method,
		Status:         PaymentInitiated,
		IdempotencyKey: cmd.IdempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.paymentRepo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	s.emitAudit(cmd.TenantID, cmd.StudentID, "USER", "billing:payment:initiated", "payment_transaction", tx.ID, audit.StatusSuccess, map[string]interface{}{
		"transaction_ref": txnRef,
		"amount":          cmd.Amount,
		"invoice_id":      cmd.InvoiceID,
	})

	return tx, nil
}

// ProcessPaymentSuccess finalizes a payment, updates the invoice, writes balanced double-entry ledger entries, and issues a receipt.
func (s *Service) ProcessPaymentSuccess(ctx context.Context, cmd ProcessPaymentSuccessCommand) (*PaymentTransaction, error) {
	tx, err := s.paymentRepo.GetTransactionByID(ctx, cmd.TenantID, cmd.TransactionID)
	if err != nil {
		return nil, err
	}
	if tx.Status == PaymentSuccess {
		return tx, nil // Idempotent return
	}

	invoice, err := s.invoiceRepo.GetInvoiceByID(ctx, cmd.TenantID, tx.InvoiceID)
	if err != nil {
		return nil, err
	}

	if err := invoice.ApplyPayment(tx.Amount); err != nil {
		return nil, err
	}

	yearStr := time.Now().UTC().Format("2006")
	recSeq, err := s.seqRepo.NextReceiptSequence(ctx, cmd.TenantID, yearStr)
	if err != nil {
		return nil, err
	}
	receiptNumber := fmt.Sprintf("REC-%s-%06d", yearStr, recSeq)

	if err := tx.MarkSuccess(receiptNumber, cmd.PaidAt); err != nil {
		return nil, err
	}
	tx.GatewayPaymentID = cmd.GatewayPaymentID

	// Create Double-Entry Ledger Records:
	// Debit: Asset (1000-CASH or 1100-BANK)
	// Credit: Accounts Receivable / Revenue (2000-AR)
	cashAccountCode := "1100-BANK"
	if tx.Method == MethodCash {
		cashAccountCode = "1000-CASH"
	}
	arAccountCode := "2000-AR"

	cashAcc, err := s.ledgerRepo.GetAccountByCode(ctx, cmd.TenantID, cashAccountCode)
	if err != nil {
		// Auto-initialize standard account if missing
		cashAcc = &LedgerAccount{
			ID:        generateID("acc_"),
			TenantID:  cmd.TenantID,
			Code:      cashAccountCode,
			Name:      "Cash / Bank Asset Account",
			Type:      AccountAsset,
			Balance:   0,
			IsActive:  true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		_ = s.ledgerRepo.CreateAccount(ctx, cashAcc)
	}

	arAcc, err := s.ledgerRepo.GetAccountByCode(ctx, cmd.TenantID, arAccountCode)
	if err != nil {
		arAcc = &LedgerAccount{
			ID:        generateID("acc_"),
			TenantID:  cmd.TenantID,
			Code:      arAccountCode,
			Name:      "Accounts Receivable",
			Type:      AccountAsset,
			Balance:   0,
			IsActive:  true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		_ = s.ledgerRepo.CreateAccount(ctx, arAcc)
	}

	ledgerEntries := []LedgerEntry{
		{
			ID:            generateID("lent_"),
			TenantID:      cmd.TenantID,
			TransactionID: &tx.ID,
			AccountID:     cashAcc.ID,
			EntryType:     EntryDebit,
			Amount:        tx.Amount,
			Description:   fmt.Sprintf("Payment received for %s (%s)", invoice.InvoiceNumber, receiptNumber),
			PostedAt:      cmd.PaidAt,
			CreatedAt:     time.Now().UTC(),
		},
		{
			ID:            generateID("lent_"),
			TenantID:      cmd.TenantID,
			TransactionID: &tx.ID,
			AccountID:     arAcc.ID,
			EntryType:     EntryCredit,
			Amount:        tx.Amount,
			Description:   fmt.Sprintf("AR clearance for invoice %s", invoice.InvoiceNumber),
			PostedAt:      cmd.PaidAt,
			CreatedAt:     time.Now().UTC(),
		},
	}

	if err := ValidateBalancedLedger(ledgerEntries); err != nil {
		return nil, err
	}

	if err := s.ledgerRepo.AppendEntries(ctx, ledgerEntries); err != nil {
		return nil, err
	}

	_ = s.ledgerRepo.UpdateAccountBalance(ctx, cmd.TenantID, cashAcc.ID, tx.Amount)
	_ = s.ledgerRepo.UpdateAccountBalance(ctx, cmd.TenantID, arAcc.ID, -tx.Amount)

	if err := s.invoiceRepo.UpdateInvoice(ctx, invoice); err != nil {
		return nil, err
	}
	if err := s.paymentRepo.UpdateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	s.emitAudit(cmd.TenantID, tx.StudentID, "SYSTEM", "billing:payment:success", "payment_transaction", tx.ID, audit.StatusSuccess, map[string]interface{}{
		"receipt_number": receiptNumber,
		"amount":         tx.Amount,
		"balance_amount": invoice.BalanceAmount,
		"invoice_status": invoice.Status,
	})

	return tx, nil
}

// ProcessRefund reverses a payment, updates the invoice balance, and posts reversal ledger entries.
func (s *Service) ProcessRefund(ctx context.Context, tenantID, transactionID, reason string) (*PaymentTransaction, error) {
	tx, err := s.paymentRepo.GetTransactionByID(ctx, tenantID, transactionID)
	if err != nil {
		return nil, err
	}
	if tx.Status != PaymentSuccess {
		return nil, NewDomainError("INVALID_REFUND", "only successful transactions can be refunded", ErrInvalidStateTransition)
	}

	invoice, err := s.invoiceRepo.GetInvoiceByID(ctx, tenantID, tx.InvoiceID)
	if err != nil {
		return nil, err
	}

	if err := tx.Refund(); err != nil {
		return nil, err
	}
	if reason != "" {
		tx.FailureReason = "Refund: " + reason
	}

	// Adjust invoice paid amount and balance
	invoice.PaidAmount -= tx.Amount
	invoice.BalanceAmount = invoice.TotalAmount - invoice.PaidAmount
	if invoice.PaidAmount == 0 {
		invoice.Status = InvoiceIssued
	} else {
		invoice.Status = InvoicePartiallyPaid
	}

	cashAccountCode := "1100-BANK"
	if tx.Method == MethodCash {
		cashAccountCode = "1000-CASH"
	}
	arAccountCode := "2000-AR"

	cashAcc, _ := s.ledgerRepo.GetAccountByCode(ctx, tenantID, cashAccountCode)
	arAcc, _ := s.ledgerRepo.GetAccountByCode(ctx, tenantID, arAccountCode)

	if cashAcc != nil && arAcc != nil {
		now := time.Now().UTC()
		reversalEntries := []LedgerEntry{
			{
				ID:            generateID("lent_"),
				TenantID:      tenantID,
				TransactionID: &tx.ID,
				AccountID:     arAcc.ID,
				EntryType:     EntryDebit,
				Amount:        tx.Amount,
				Description:   fmt.Sprintf("Refund AR restoration for %s", invoice.InvoiceNumber),
				PostedAt:      now,
				CreatedAt:     now,
			},
			{
				ID:            generateID("lent_"),
				TenantID:      tenantID,
				TransactionID: &tx.ID,
				AccountID:     cashAcc.ID,
				EntryType:     EntryCredit,
				Amount:        tx.Amount,
				Description:   fmt.Sprintf("Refund payout for txn %s", tx.TransactionRef),
				PostedAt:      now,
				CreatedAt:     now,
			},
		}
		_ = ValidateBalancedLedger(reversalEntries)
		_ = s.ledgerRepo.AppendEntries(ctx, reversalEntries)
		_ = s.ledgerRepo.UpdateAccountBalance(ctx, tenantID, cashAcc.ID, -tx.Amount)
		_ = s.ledgerRepo.UpdateAccountBalance(ctx, tenantID, arAcc.ID, tx.Amount)
	}

	_ = s.invoiceRepo.UpdateInvoice(ctx, invoice)
	_ = s.paymentRepo.UpdateTransaction(ctx, tx)

	s.emitAudit(tenantID, tx.StudentID, "ADMIN", "billing:payment:refunded", "payment_transaction", tx.ID, audit.StatusSuccess, map[string]interface{}{
		"amount":         tx.Amount,
		"reason":         reason,
		"balance_amount": invoice.BalanceAmount,
	})

	return tx, nil
}

func (s *Service) GetInvoice(ctx context.Context, tenantID, invoiceID string) (*StudentInvoice, error) {
	return s.invoiceRepo.GetInvoiceByID(ctx, tenantID, invoiceID)
}

func (s *Service) ListInvoices(ctx context.Context, filter InvoiceFilter) ([]*StudentInvoice, int, error) {
	return s.invoiceRepo.ListInvoices(ctx, filter)
}

func (s *Service) GetStudentBalance(ctx context.Context, tenantID, studentID string) (int, error) {
	invoices, _, err := s.invoiceRepo.ListInvoices(ctx, InvoiceFilter{
		TenantID:  tenantID,
		StudentID: studentID,
		Limit:     1000,
	})
	if err != nil {
		return 0, err
	}

	totalBalance := 0
	for _, inv := range invoices {
		if inv.Status != InvoiceCancelled && inv.Status != InvoiceWrittenOff {
			totalBalance += inv.BalanceAmount
		}
	}
	return totalBalance, nil
}

func (s *Service) emitAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
	if s.auditSub == nil {
		return
	}
	var metaRaw json.RawMessage
	if meta != nil {
		if b, err := json.Marshal(meta); err == nil {
			metaRaw = b
		}
	}
	_ = s.auditSub.Enqueue(audit.RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      &actorID,
		ActorType:    audit.ActorType(actorType),
		Action:       action,
		ResourceType: resType,
		ResourceID:   &resID,
		Status:       status,
		Metadata:     metaRaw,
	})
}

func generateID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + hex.EncodeToString(b)
}
