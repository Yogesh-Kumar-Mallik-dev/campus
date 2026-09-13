/**
 * BLOCK_BILLING_ERRORS_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   Authoritative domain error types and sentinels for billing, invoices, and ledger operations.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package billing

import (
	"errors"
	"fmt"
)

var (
	ErrFeeStructureNotFound   = errors.New("fee structure not found")
	ErrInvoiceNotFound        = errors.New("student invoice not found")
	ErrTransactionNotFound    = errors.New("payment transaction not found")
	ErrAccountNotFound        = errors.New("ledger account not found")
	ErrTenantRequired         = errors.New("tenant_id is required")
	ErrInvalidInput           = errors.New("invalid billing input parameters")
	ErrInvalidStateTransition = errors.New("invalid invoice or payment state transition")
	ErrInvoiceAlreadyPaid     = errors.New("invoice is already fully paid")
	ErrInvoiceCancelled       = errors.New("invoice is cancelled or written off")
	ErrDuplicateInvoiceNumber = errors.New("invoice number already exists")
	ErrDuplicateTxnRef        = errors.New("transaction reference already exists")
	ErrOverpaymentNotAllowed  = errors.New("payment amount exceeds invoice balance")
	ErrLedgerUnbalanced       = errors.New("double-entry ledger entries must be balanced (debits equal credits)")
	ErrIdempotencyConflict    = errors.New("idempotency key conflict: transaction already processed with different payload")
	ErrInvalidAmount          = errors.New("amount must be strictly positive")
)

// DomainError captures contextual error data for billing operations.
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError constructs a structured billing DomainError.
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
