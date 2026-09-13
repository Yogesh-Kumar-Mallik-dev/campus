/**
 * BLOCK_BILLING_DOMAIN_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   Domain entities, invoice state machines, payment transactions, double-entry ledger invariants, and fee templates.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package billing

import (
	"strings"
	"time"
)

type FeeCategory string

const (
	CategoryTuition        FeeCategory = "TUITION"
	CategoryAdmission      FeeCategory = "ADMISSION"
	CategoryExamination    FeeCategory = "EXAMINATION"
	CategoryHostel         FeeCategory = "HOSTEL"
	CategoryMess           FeeCategory = "MESS"
	CategoryLibrary        FeeCategory = "LIBRARY"
	CategoryTransportation FeeCategory = "TRANSPORTATION"
	CategoryLaboratory     FeeCategory = "LABORATORY"
	CategoryMiscellaneous  FeeCategory = "MISCELLANEOUS"
)

type InvoiceStatus string

const (
	InvoiceDraft         InvoiceStatus = "DRAFT"
	InvoiceIssued        InvoiceStatus = "ISSUED"
	InvoicePartiallyPaid InvoiceStatus = "PARTIALLY_PAID"
	InvoicePaid          InvoiceStatus = "PAID"
	InvoiceOverdue       InvoiceStatus = "OVERDUE"
	InvoiceCancelled     InvoiceStatus = "CANCELLED"
	InvoiceWrittenOff    InvoiceStatus = "WRITTEN_OFF"
)

type PaymentStatus string

const (
	PaymentInitiated  PaymentStatus = "INITIATED"
	PaymentProcessing PaymentStatus = "PROCESSING"
	PaymentSuccess    PaymentStatus = "SUCCESS"
	PaymentFailed     PaymentStatus = "FAILED"
	PaymentRefunded   PaymentStatus = "REFUNDED"
)

type PaymentMethod string

const (
	MethodOnlineGateway     PaymentMethod = "ONLINE_GATEWAY"
	MethodBankTransfer      PaymentMethod = "BANK_TRANSFER"
	MethodUPI               PaymentMethod = "UPI"
	MethodCheque            PaymentMethod = "CHEQUE"
	MethodCash              PaymentMethod = "CASH"
	MethodScholarshipWaiver PaymentMethod = "SCHOLARSHIP_WAIVER"
)

type LedgerAccountType string

const (
	AccountAsset     LedgerAccountType = "ASSET"
	AccountLiability LedgerAccountType = "LIABILITY"
	AccountEquity    LedgerAccountType = "EQUITY"
	AccountRevenue   LedgerAccountType = "REVENUE"
	AccountExpense   LedgerAccountType = "EXPENSE"
)

type LedgerEntryType string

const (
	EntryDebit  LedgerEntryType = "DEBIT"
	EntryCredit LedgerEntryType = "CREDIT"
)

type FeeStructureItem struct {
	ID             string      `json:"id"`
	TenantID       string      `json:"tenant_id"`
	FeeStructureID string      `json:"fee_structure_id"`
	Category       FeeCategory `json:"category"`
	Name           string      `json:"name"`
	Amount         int         `json:"amount"` // Smallest currency unit (cents/paise)
	IsOptional     bool        `json:"is_optional"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type FeeStructure struct {
	ID           string             `json:"id"`
	TenantID     string             `json:"tenant_id"`
	ProgramID    string             `json:"program_id"`
	AcademicYear string             `json:"academic_year"` // "2026-2027"
	Semester     int                `json:"semester"`
	Name         string             `json:"name"`
	TotalAmount  int                `json:"total_amount"`
	Currency     string             `json:"currency"`
	DueDate      string             `json:"due_date"` // YYYY-MM-DD
	IsActive     bool               `json:"is_active"`
	Items        []FeeStructureItem `json:"items,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type InvoiceItem struct {
	ID        string      `json:"id"`
	TenantID  string      `json:"tenant_id"`
	InvoiceID string      `json:"invoice_id"`
	Category  FeeCategory `json:"category"`
	Name      string      `json:"name"`
	Amount    int         `json:"amount"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type StudentInvoice struct {
	ID             string         `json:"id"`
	TenantID       string         `json:"tenant_id"`
	StudentID      string         `json:"student_id"`
	FeeStructureID *string        `json:"fee_structure_id,omitempty"`
	InvoiceNumber  string         `json:"invoice_number"` // "INV-2026-000001"
	AcademicYear   string         `json:"academic_year"`
	Semester       int            `json:"semester"`
	SubtotalAmount int            `json:"subtotal_amount"`
	DiscountAmount int            `json:"discount_amount"`
	TaxAmount      int            `json:"tax_amount"`
	TotalAmount    int            `json:"total_amount"`
	PaidAmount     int            `json:"paid_amount"`
	BalanceAmount  int            `json:"balance_amount"`
	DueDate        string         `json:"due_date"` // YYYY-MM-DD
	Status         InvoiceStatus  `json:"status"`
	Notes          string         `json:"notes,omitempty"`
	Items          []InvoiceItem  `json:"items,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type PaymentTransaction struct {
	ID               string        `json:"id"`
	TenantID         string        `json:"tenant_id"`
	InvoiceID        string        `json:"invoice_id"`
	StudentID        string        `json:"student_id"`
	TransactionRef   string        `json:"transaction_ref"` // "TXN-2026-000001"
	GatewayName      string        `json:"gateway_name"`
	GatewayOrderID   *string       `json:"gateway_order_id,omitempty"`
	GatewayPaymentID *string       `json:"gateway_payment_id,omitempty"`
	Amount           int           `json:"amount"`
	Currency         string        `json:"currency"`
	Method           PaymentMethod `json:"method"`
	Status           PaymentStatus `json:"status"`
	IdempotencyKey   *string       `json:"idempotency_key,omitempty"`
	ReceiptNumber    *string       `json:"receipt_number,omitempty"` // "REC-2026-000001"
	FailureReason    string        `json:"failure_reason,omitempty"`
	PaidAt           *time.Time    `json:"paid_at,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type LedgerAccount struct {
	ID        string            `json:"id"`
	TenantID  string            `json:"tenant_id"`
	Code      string            `json:"code"` // "1000-CASH", "2000-AR", "4000-TUITION"
	Name      string            `json:"name"`
	Type      LedgerAccountType `json:"type"`
	Balance   int               `json:"balance"` // Running balance
	IsActive  bool              `json:"is_active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type LedgerEntry struct {
	ID            string          `json:"id"`
	TenantID      string          `json:"tenant_id"`
	TransactionID *string         `json:"transaction_id,omitempty"`
	AccountID     string          `json:"account_id"`
	EntryType     LedgerEntryType `json:"entry_type"`
	Amount        int             `json:"amount"`
	Description   string          `json:"description"`
	PostedAt      time.Time       `json:"posted_at"`
	CreatedAt     time.Time       `json:"created_at"`
}

// Issue transitions DRAFT -> ISSUED.
func (inv *StudentInvoice) Issue() error {
	if inv.Status != InvoiceDraft {
		return NewDomainError("INVALID_STATUS", "only draft invoices can be issued", ErrInvalidStateTransition)
	}
	inv.Status = InvoiceIssued
	inv.UpdatedAt = time.Now().UTC()
	return nil
}

// ApplyPayment updates invoice paid amount, calculates balance, and updates status.
func (inv *StudentInvoice) ApplyPayment(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if inv.Status == InvoiceCancelled || inv.Status == InvoiceWrittenOff {
		return ErrInvoiceCancelled
	}
	if inv.Status == InvoicePaid {
		return ErrInvoiceAlreadyPaid
	}
	if amount > inv.BalanceAmount {
		return ErrOverpaymentNotAllowed
	}

	inv.PaidAmount += amount
	inv.BalanceAmount = inv.TotalAmount - inv.PaidAmount

	if inv.BalanceAmount == 0 {
		inv.Status = InvoicePaid
	} else {
		inv.Status = InvoicePartiallyPaid
	}
	inv.UpdatedAt = time.Now().UTC()
	return nil
}

// Cancel transitions ISSUED -> CANCELLED if no payments received.
func (inv *StudentInvoice) Cancel(reason string) error {
	if inv.PaidAmount > 0 {
		return NewDomainError("CANNOT_CANCEL", "cannot cancel invoice with existing payments", ErrInvalidStateTransition)
	}
	if inv.Status == InvoiceCancelled || inv.Status == InvoicePaid {
		return NewDomainError("INVALID_STATUS", "invoice is already finalized", ErrInvalidStateTransition)
	}
	inv.Status = InvoiceCancelled
	if reason != "" {
		if inv.Notes != "" {
			inv.Notes += " | Cancelled: " + reason
		} else {
			inv.Notes = "Cancelled: " + reason
		}
	}
	inv.UpdatedAt = time.Now().UTC()
	return nil
}

// MarkSuccess transitions INITIATED/PROCESSING -> SUCCESS.
func (tx *PaymentTransaction) MarkSuccess(receiptNumber string, paidAt time.Time) error {
	if tx.Status == PaymentSuccess {
		return nil // Idempotent
	}
	if tx.Status == PaymentRefunded {
		return NewDomainError("INVALID_STATUS", "refunded transaction cannot be marked successful", ErrInvalidStateTransition)
	}

	tx.Status = PaymentSuccess
	tx.ReceiptNumber = &receiptNumber
	tx.PaidAt = &paidAt
	tx.UpdatedAt = time.Now().UTC()
	return nil
}

// MarkFailed transitions INITIATED/PROCESSING -> FAILED.
func (tx *PaymentTransaction) MarkFailed(reason string) error {
	if tx.Status == PaymentSuccess {
		return NewDomainError("INVALID_STATUS", "successful transaction cannot be failed", ErrInvalidStateTransition)
	}
	tx.Status = PaymentFailed
	tx.FailureReason = reason
	tx.UpdatedAt = time.Now().UTC()
	return nil
}

// Refund transitions SUCCESS -> REFUNDED.
func (tx *PaymentTransaction) Refund() error {
	if tx.Status != PaymentSuccess {
		return NewDomainError("INVALID_STATUS", "only successful payments can be refunded", ErrInvalidStateTransition)
	}
	tx.Status = PaymentRefunded
	tx.UpdatedAt = time.Now().UTC()
	return nil
}

// ValidateBalancedLedger enforces the foundational double-entry invariant: Sum(Debit) == Sum(Credit).
func ValidateBalancedLedger(entries []LedgerEntry) error {
	if len(entries) < 2 {
		return NewDomainError("INCOMPLETE_LEDGER", "double-entry ledger requires at least two balanced entries", ErrLedgerUnbalanced)
	}

	totalDebit := 0
	totalCredit := 0

	for _, e := range entries {
		if e.Amount <= 0 {
			return ErrInvalidAmount
		}
		if strings.TrimSpace(e.AccountID) == "" {
			return NewDomainError("MISSING_ACCOUNT", "ledger entry must reference a valid account", ErrInvalidInput)
		}
		if e.EntryType == EntryDebit {
			totalDebit += e.Amount
		} else if e.EntryType == EntryCredit {
			totalCredit += e.Amount
		} else {
			return NewDomainError("INVALID_ENTRY_TYPE", "entry must be DEBIT or CREDIT", ErrInvalidInput)
		}
	}

	if totalDebit != totalCredit {
		return NewDomainError("UNBALANCED_LEDGER", "total debits do not equal total credits", ErrLedgerUnbalanced)
	}
	return nil
}
