/**
 * BLOCK_BILLING_REPOSITORY_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   Data access interface definitions for fee structures, student invoices, transactions, ledger, and sequencing.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package billing

import (
	"context"
)

type InvoiceFilter struct {
	TenantID     string
	StudentID    string
	AcademicYear string
	Semester     int
	Status       *InvoiceStatus
	Limit        int
	Offset       int
}

type FeeRepository interface {
	CreateFeeStructure(ctx context.Context, fs *FeeStructure) error
	GetFeeStructureByID(ctx context.Context, tenantID, id string) (*FeeStructure, error)
	GetFeeStructureByProgram(ctx context.Context, tenantID, programID, academicYear string, semester int) (*FeeStructure, error)
	ListFeeStructures(ctx context.Context, tenantID string) ([]*FeeStructure, error)
}

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, inv *StudentInvoice) error
	GetInvoiceByID(ctx context.Context, tenantID, id string) (*StudentInvoice, error)
	GetInvoiceByNumber(ctx context.Context, tenantID, invoiceNumber string) (*StudentInvoice, error)
	UpdateInvoice(ctx context.Context, inv *StudentInvoice) error
	ListInvoices(ctx context.Context, filter InvoiceFilter) ([]*StudentInvoice, int, error)
}

type PaymentRepository interface {
	CreateTransaction(ctx context.Context, tx *PaymentTransaction) error
	GetTransactionByID(ctx context.Context, tenantID, id string) (*PaymentTransaction, error)
	GetTransactionByRef(ctx context.Context, tenantID, ref string) (*PaymentTransaction, error)
	GetTransactionByIdempotencyKey(ctx context.Context, tenantID, key string) (*PaymentTransaction, error)
	UpdateTransaction(ctx context.Context, tx *PaymentTransaction) error
	ListTransactionsByInvoice(ctx context.Context, tenantID, invoiceID string) ([]*PaymentTransaction, error)
	ListTransactionsByStudent(ctx context.Context, tenantID, studentID string) ([]*PaymentTransaction, error)
}

type LedgerRepository interface {
	CreateAccount(ctx context.Context, acc *LedgerAccount) error
	GetAccountByCode(ctx context.Context, tenantID, code string) (*LedgerAccount, error)
	UpdateAccountBalance(ctx context.Context, tenantID, accountID string, balanceDelta int) error
	AppendEntries(ctx context.Context, entries []LedgerEntry) error
	ListEntriesByTransaction(ctx context.Context, tenantID, txID string) ([]LedgerEntry, error)
}

type SequenceRepository interface {
	NextInvoiceSequence(ctx context.Context, tenantID, academicYear string) (int, error)
	NextReceiptSequence(ctx context.Context, tenantID, academicYear string) (int, error)
	NextTransactionSequence(ctx context.Context, tenantID, academicYear string) (int, error)
}
