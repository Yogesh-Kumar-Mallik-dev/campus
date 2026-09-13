/**
 * BLOCK_LIBRARY_MOCK_REPOSITORY_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   In-memory mock repository implementing Repository interface for unit testing.
 */

package library

import (
	"context"
	"strings"
	"sync"
	"time"
)

type MockRepository struct {
	mu           sync.RWMutex
	books        map[string]*Book
	copies       map[string]*BookCopy
	borrows      map[string]*BorrowRecord
	reservations map[string]*Reservation
	ebooks       map[string]*EbookAccess
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		books:        make(map[string]*Book),
		copies:       make(map[string]*BookCopy),
		borrows:      make(map[string]*BorrowRecord),
		reservations: make(map[string]*Reservation),
		ebooks:       make(map[string]*EbookAccess),
	}
}

func (m *MockRepository) CreateBook(ctx context.Context, book *Book) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, b := range m.books {
		if b.TenantID == book.TenantID && b.ISBN == book.ISBN {
			return ErrISBNExists
		}
	}
	m.books[book.ID] = book
	return nil
}

func (m *MockRepository) GetBookByID(ctx context.Context, tenantID, id string) (*Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, exists := m.books[id]
	if !exists || b.TenantID != tenantID {
		return nil, ErrBookNotFound
	}
	return b, nil
}

func (m *MockRepository) GetBookByISBN(ctx context.Context, tenantID, isbn string) (*Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, b := range m.books {
		if b.TenantID == tenantID && b.ISBN == isbn {
			return b, nil
		}
	}
	return nil, ErrBookNotFound
}

func (m *MockRepository) ListBooks(ctx context.Context, tenantID string, category *BookCategory, query string) ([]*Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Book
	for _, b := range m.books {
		if b.TenantID != tenantID {
			continue
		}
		if category != nil && b.Category != *category {
			continue
		}
		if query != "" {
			q := strings.ToLower(query)
			if !strings.Contains(strings.ToLower(b.Title), q) &&
				!strings.Contains(strings.ToLower(b.Author), q) &&
				!strings.Contains(strings.ToLower(b.ISBN), q) {
				continue
			}
		}
		result = append(result, b)
	}
	return result, nil
}

func (m *MockRepository) UpdateBookInventory(ctx context.Context, tenantID, bookID string, totalCopies, availableCopies int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, exists := m.books[bookID]
	if !exists || b.TenantID != tenantID {
		return ErrBookNotFound
	}
	b.TotalCopies = totalCopies
	b.AvailableCopies = availableCopies
	b.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) CreateCopy(ctx context.Context, copy *BookCopy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, c := range m.copies {
		if c.TenantID == copy.TenantID && c.AccessionNumber == copy.AccessionNumber {
			return ErrAccessionExists
		}
	}
	m.copies[copy.ID] = copy
	return nil
}

func (m *MockRepository) GetCopyByID(ctx context.Context, tenantID, id string) (*BookCopy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, exists := m.copies[id]
	if !exists || c.TenantID != tenantID {
		return nil, ErrBookCopyNotFound
	}
	return c, nil
}

func (m *MockRepository) GetCopyByAccession(ctx context.Context, tenantID, accessionNumber string) (*BookCopy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, c := range m.copies {
		if c.TenantID == tenantID && c.AccessionNumber == accessionNumber {
			return c, nil
		}
	}
	return nil, ErrBookCopyNotFound
}

func (m *MockRepository) ListCopiesByBook(ctx context.Context, tenantID, bookID string) ([]*BookCopy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*BookCopy
	for _, c := range m.copies {
		if c.TenantID == tenantID && c.BookID == bookID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *MockRepository) UpdateCopyStatus(ctx context.Context, tenantID, id string, status BookCopyStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, exists := m.copies[id]
	if !exists || c.TenantID != tenantID {
		return ErrBookCopyNotFound
	}
	c.Status = status
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) CreateBorrowRecord(ctx context.Context, record *BorrowRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.borrows[record.ID] = record
	return nil
}

func (m *MockRepository) GetBorrowRecordByID(ctx context.Context, tenantID, id string) (*BorrowRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, exists := m.borrows[id]
	if !exists || r.TenantID != tenantID {
		return nil, ErrBorrowRecordNotFound
	}
	return r, nil
}

func (m *MockRepository) GetActiveBorrowsByStudent(ctx context.Context, tenantID, studentID string) ([]*BorrowRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*BorrowRecord
	for _, r := range m.borrows {
		if r.TenantID == tenantID && r.StudentID == studentID && (r.Status == BorrowStatusIssued || r.Status == BorrowStatusOverdue) {
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *MockRepository) UpdateBorrowRecord(ctx context.Context, record *BorrowRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.borrows[record.ID] = record
	return nil
}

func (m *MockRepository) ListBorrowRecords(ctx context.Context, tenantID, studentID string, status *BorrowStatus) ([]*BorrowRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*BorrowRecord
	for _, r := range m.borrows {
		if r.TenantID != tenantID {
			continue
		}
		if studentID != "" && r.StudentID != studentID {
			continue
		}
		if status != nil && r.Status != *status {
			continue
		}
		result = append(result, r)
	}
	return result, nil
}

func (m *MockRepository) CreateReservation(ctx context.Context, reservation *Reservation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reservations[reservation.ID] = reservation
	return nil
}

func (m *MockRepository) GetReservationByID(ctx context.Context, tenantID, id string) (*Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, exists := m.reservations[id]
	if !exists || r.TenantID != tenantID {
		return nil, ErrReservationNotFound
	}
	return r, nil
}

func (m *MockRepository) ListReservationsByBook(ctx context.Context, tenantID, bookID string) ([]*Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Reservation
	for _, r := range m.reservations {
		if r.TenantID == tenantID && r.BookID == bookID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *MockRepository) UpdateReservationStatus(ctx context.Context, tenantID, id string, status ReservationStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, exists := m.reservations[id]
	if !exists || r.TenantID != tenantID {
		return ErrReservationNotFound
	}
	r.Status = status
	r.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) RecordEbookAccess(ctx context.Context, access *EbookAccess) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := access.TenantID + ":" + access.BookID + ":" + access.StudentID
	m.ebooks[key] = access
	return nil
}

func (m *MockRepository) GetEbookAccess(ctx context.Context, tenantID, bookID, studentID string) (*EbookAccess, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := tenantID + ":" + bookID + ":" + studentID
	access, exists := m.ebooks[key]
	if !exists {
		return nil, nil
	}
	return access, nil
}
