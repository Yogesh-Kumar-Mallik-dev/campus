/**
 * BLOCK_BILLING_MOCK_REPOSITORY_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   In-memory, concurrency-safe mock repositories for fee structures, invoices, transactions, ledger, and sequences.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package billing

import (
	"context"
	"strings"
	"sync"
)

type MockFeeRepository struct {
	mu         sync.RWMutex
	structures map[string]*FeeStructure
}

func NewMockFeeRepository() *MockFeeRepository {
	return &MockFeeRepository{
		structures: make(map[string]*FeeStructure),
	}
}

func (m *MockFeeRepository) CreateFeeStructure(ctx context.Context, fs *FeeStructure) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.structures[fs.TenantID+":"+fs.ID] = fs
	return nil
}

func (m *MockFeeRepository) GetFeeStructureByID(ctx context.Context, tenantID, id string) (*FeeStructure, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fs, ok := m.structures[tenantID+":"+id]
	if !ok {
		return nil, ErrFeeStructureNotFound
	}
	c := *fs
	return &c, nil
}

func (m *MockFeeRepository) GetFeeStructureByProgram(ctx context.Context, tenantID, programID, academicYear string, semester int) (*FeeStructure, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, fs := range m.structures {
		if fs.TenantID == tenantID && fs.ProgramID == programID && fs.AcademicYear == academicYear && fs.Semester == semester {
			c := *fs
			return &c, nil
		}
	}
	return nil, ErrFeeStructureNotFound
}

func (m *MockFeeRepository) ListFeeStructures(ctx context.Context, tenantID string) ([]*FeeStructure, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*FeeStructure
	for _, fs := range m.structures {
		if fs.TenantID == tenantID {
			c := *fs
			res = append(res, &c)
		}
	}
	return res, nil
}

type MockInvoiceRepository struct {
	mu       sync.RWMutex
	invoices map[string]*StudentInvoice
}

func NewMockInvoiceRepository() *MockInvoiceRepository {
	return &MockInvoiceRepository{
		invoices: make(map[string]*StudentInvoice),
	}
}

func (m *MockInvoiceRepository) CreateInvoice(ctx context.Context, inv *StudentInvoice) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.invoices[inv.TenantID+":"+inv.ID] = inv
	return nil
}

func (m *MockInvoiceRepository) GetInvoiceByID(ctx context.Context, tenantID, id string) (*StudentInvoice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	inv, ok := m.invoices[tenantID+":"+id]
	if !ok {
		return nil, ErrInvoiceNotFound
	}
	c := *inv
	return &c, nil
}

func (m *MockInvoiceRepository) GetInvoiceByNumber(ctx context.Context, tenantID, invoiceNumber string) (*StudentInvoice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, inv := range m.invoices {
		if inv.TenantID == tenantID && strings.EqualFold(inv.InvoiceNumber, invoiceNumber) {
			c := *inv
			return &c, nil
		}
	}
	return nil, ErrInvoiceNotFound
}

func (m *MockInvoiceRepository) UpdateInvoice(ctx context.Context, inv *StudentInvoice) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := inv.TenantID + ":" + inv.ID
	if _, ok := m.invoices[key]; !ok {
		return ErrInvoiceNotFound
	}
	m.invoices[key] = inv
	return nil
}

func (m *MockInvoiceRepository) ListInvoices(ctx context.Context, filter InvoiceFilter) ([]*StudentInvoice, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*StudentInvoice
	for _, inv := range m.invoices {
		if filter.TenantID != "" && inv.TenantID != filter.TenantID {
			continue
		}
		if filter.StudentID != "" && inv.StudentID != filter.StudentID {
			continue
		}
		if filter.AcademicYear != "" && inv.AcademicYear != filter.AcademicYear {
			continue
		}
		if filter.Semester > 0 && inv.Semester != filter.Semester {
			continue
		}
		if filter.Status != nil && inv.Status != *filter.Status {
			continue
		}
		c := *inv
		results = append(results, &c)
	}

	total := len(results)
	if filter.Offset >= total {
		return []*StudentInvoice{}, total, nil
	}
	end := total
	if filter.Limit > 0 && filter.Offset+filter.Limit < end {
		end = filter.Offset + filter.Limit
	}
	return results[filter.Offset:end], total, nil
}

type MockPaymentRepository struct {
	mu           sync.RWMutex
	transactions map[string]*PaymentTransaction
}

func NewMockPaymentRepository() *MockPaymentRepository {
	return &MockPaymentRepository{
		transactions: make(map[string]*PaymentTransaction),
	}
}

func (m *MockPaymentRepository) CreateTransaction(ctx context.Context, tx *PaymentTransaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transactions[tx.TenantID+":"+tx.ID] = tx
	return nil
}

func (m *MockPaymentRepository) GetTransactionByID(ctx context.Context, tenantID, id string) (*PaymentTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tx, ok := m.transactions[tenantID+":"+id]
	if !ok {
		return nil, ErrTransactionNotFound
	}
	c := *tx
	return &c, nil
}

func (m *MockPaymentRepository) GetTransactionByRef(ctx context.Context, tenantID, ref string) (*PaymentTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, tx := range m.transactions {
		if tx.TenantID == tenantID && strings.EqualFold(tx.TransactionRef, ref) {
			c := *tx
			return &c, nil
		}
	}
	return nil, ErrTransactionNotFound
}

func (m *MockPaymentRepository) GetTransactionByIdempotencyKey(ctx context.Context, tenantID, key string) (*PaymentTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, tx := range m.transactions {
		if tx.TenantID == tenantID && tx.IdempotencyKey != nil && *tx.IdempotencyKey == key {
			c := *tx
			return &c, nil
		}
	}
	return nil, ErrTransactionNotFound
}

func (m *MockPaymentRepository) UpdateTransaction(ctx context.Context, tx *PaymentTransaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := tx.TenantID + ":" + tx.ID
	if _, ok := m.transactions[key]; !ok {
		return ErrTransactionNotFound
	}
	m.transactions[key] = tx
	return nil
}

func (m *MockPaymentRepository) ListTransactionsByInvoice(ctx context.Context, tenantID, invoiceID string) ([]*PaymentTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*PaymentTransaction
	for _, tx := range m.transactions {
		if tx.TenantID == tenantID && tx.InvoiceID == invoiceID {
			c := *tx
			res = append(res, &c)
		}
	}
	return res, nil
}

func (m *MockPaymentRepository) ListTransactionsByStudent(ctx context.Context, tenantID, studentID string) ([]*PaymentTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*PaymentTransaction
	for _, tx := range m.transactions {
		if tx.TenantID == tenantID && tx.StudentID == studentID {
			c := *tx
			res = append(res, &c)
		}
	}
	return res, nil
}

type MockLedgerRepository struct {
	mu       sync.RWMutex
	accounts map[string]*LedgerAccount
	entries  []LedgerEntry
}

func NewMockLedgerRepository() *MockLedgerRepository {
	return &MockLedgerRepository{
		accounts: make(map[string]*LedgerAccount),
		entries:  make([]LedgerEntry, 0),
	}
}

func (m *MockLedgerRepository) CreateAccount(ctx context.Context, acc *LedgerAccount) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.accounts[acc.TenantID+":"+acc.Code] = acc
	return nil
}

func (m *MockLedgerRepository) GetAccountByCode(ctx context.Context, tenantID, code string) (*LedgerAccount, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	acc, ok := m.accounts[tenantID+":"+code]
	if !ok {
		return nil, ErrAccountNotFound
	}
	c := *acc
	return &c, nil
}

func (m *MockLedgerRepository) UpdateAccountBalance(ctx context.Context, tenantID, accountID string, balanceDelta int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.accounts {
		if a.TenantID == tenantID && a.ID == accountID {
			a.Balance += balanceDelta
			return nil
		}
	}
	return nil
}

func (m *MockLedgerRepository) AppendEntries(ctx context.Context, entries []LedgerEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, entries...)
	return nil
}

func (m *MockLedgerRepository) ListEntriesByTransaction(ctx context.Context, tenantID, txID string) ([]LedgerEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []LedgerEntry
	for _, e := range m.entries {
		if e.TenantID == tenantID && e.TransactionID != nil && *e.TransactionID == txID {
			res = append(res, e)
		}
	}
	return res, nil
}

type MockSequenceRepository struct {
	mu           sync.Mutex
	invoiceSeq   map[string]int
	receiptSeq   map[string]int
	txnSeq       map[string]int
}

func NewMockSequenceRepository() *MockSequenceRepository {
	return &MockSequenceRepository{
		invoiceSeq: make(map[string]int),
		receiptSeq: make(map[string]int),
		txnSeq:     make(map[string]int),
	}
}

func (m *MockSequenceRepository) NextInvoiceSequence(ctx context.Context, tenantID, academicYear string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := tenantID + ":" + academicYear
	m.invoiceSeq[key]++
	return m.invoiceSeq[key], nil
}

func (m *MockSequenceRepository) NextReceiptSequence(ctx context.Context, tenantID, academicYear string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := tenantID + ":" + academicYear
	m.receiptSeq[key]++
	return m.receiptSeq[key], nil
}

func (m *MockSequenceRepository) NextTransactionSequence(ctx context.Context, tenantID, academicYear string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := tenantID + ":" + academicYear
	m.txnSeq[key]++
	return m.txnSeq[key], nil
}
