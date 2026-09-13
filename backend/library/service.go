/**
 * BLOCK_LIBRARY_SERVICE_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   Business domain service implementing book cataloging, checkout invariants, renewal rules, and fine calculation.
 */

package library

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"campus/backend/audit"
)

type Service struct {
	repo     Repository
	auditSub audit.Subscriber
}

func NewService(repo Repository, auditSub audit.Subscriber) *Service {
	return &Service{
		repo:     repo,
		auditSub: auditSub,
	}
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

// CatalogBook adds a new bibliographic title to the catalog
func (s *Service) CatalogBook(ctx context.Context, tenantID, isbn, title, author, publisher, edition string, category BookCategory, shelfLocation string, ebookKey, ebookFormat *string, ebookSize *int) (*Book, error) {
	if tenantID == "" || isbn == "" || title == "" || author == "" {
		return nil, NewDomainError(nil, 400, "Bad Request", "tenant_id, isbn, title and author are required", "https://campus.internal/errors/invalid-argument")
	}

	book := &Book{
		ID:              fmt.Sprintf("bk_%d", time.Now().UnixNano()),
		TenantID:        tenantID,
		ISBN:            isbn,
		Title:           title,
		Author:          author,
		Publisher:       publisher,
		Edition:         edition,
		Category:        category,
		TotalCopies:     0,
		AvailableCopies: 0,
		ShelfLocation:   shelfLocation,
		EbookKey:        ebookKey,
		EbookFormat:     ebookFormat,
		EbookSize:       ebookSize,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := s.repo.CreateBook(ctx, book); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, "SYSTEM", "SYSTEM", "library:book:cataloged", "library_book", book.ID, audit.StatusSuccess, map[string]interface{}{
		"isbn":  isbn,
		"title": title,
	})

	return book, nil
}

// AddBookCopy registers an accession copy barcode unit
func (s *Service) AddBookCopy(ctx context.Context, tenantID, bookID, accessionNumber, barcode string) (*BookCopy, error) {
	book, err := s.repo.GetBookByID(ctx, tenantID, bookID)
	if err != nil {
		return nil, err
	}

	copy := &BookCopy{
		ID:              fmt.Sprintf("cpy_%d", time.Now().UnixNano()),
		TenantID:        tenantID,
		BookID:          bookID,
		AccessionNumber: accessionNumber,
		Barcode:         barcode,
		Status:          CopyStatusAvailable,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := s.repo.CreateCopy(ctx, copy); err != nil {
		return nil, err
	}

	_ = s.repo.UpdateBookInventory(ctx, tenantID, bookID, book.TotalCopies+1, book.AvailableCopies+1)

	return copy, nil
}

// ListBooks searches and lists bibliographic titles
func (s *Service) ListBooks(ctx context.Context, tenantID string, category *BookCategory, query string) ([]*Book, error) {
	return s.repo.ListBooks(ctx, tenantID, category, query)
}

// ListCopies lists physical copies for a title
func (s *Service) ListCopies(ctx context.Context, tenantID, bookID string) ([]*BookCopy, error) {
	return s.repo.ListCopiesByBook(ctx, tenantID, bookID)
}

// BorrowBook loans a copy to a student with quota verification
func (s *Service) BorrowBook(ctx context.Context, tenantID, copyID, studentID, issuedByID string, loanDays int) (*BorrowRecord, error) {
	copy, err := s.repo.GetCopyByID(ctx, tenantID, copyID)
	if err != nil {
		return nil, err
	}

	if copy.Status != CopyStatusAvailable {
		return nil, ErrCopyUnavailable
	}

	activeBorrows, err := s.repo.GetActiveBorrowsByStudent(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}

	if err := ValidateBorrowEligibility(len(activeBorrows), 4); err != nil {
		return nil, err
	}

	book, err := s.repo.GetBookByID(ctx, tenantID, copy.BookID)
	if err != nil {
		return nil, err
	}

	if loanDays <= 0 {
		loanDays = 14 // default 14-day student checkout window
	}

	now := time.Now().UTC()
	dueDate := now.AddDate(0, 0, loanDays)

	var issuedByPtr *string
	if issuedByID != "" {
		issuedByPtr = &issuedByID
	}

	record := &BorrowRecord{
		ID:         fmt.Sprintf("brw_%d", time.Now().UnixNano()),
		TenantID:   tenantID,
		CopyID:     copyID,
		StudentID:  studentID,
		BorrowedAt: now,
		DueDate:    dueDate,
		RenewCount: 0,
		Status:     BorrowStatusIssued,
		FineAmount: 0,
		FinePaid:   false,
		IssuedByID: issuedByPtr,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateBorrowRecord(ctx, record); err != nil {
		return nil, err
	}

	_ = s.repo.UpdateCopyStatus(ctx, tenantID, copyID, CopyStatusIssued)
	if book.AvailableCopies > 0 {
		_ = s.repo.UpdateBookInventory(ctx, tenantID, book.ID, book.TotalCopies, book.AvailableCopies-1)
	}

	s.emitAudit(tenantID, studentID, "USER", "library:book:borrowed", "library_borrow_record", record.ID, audit.StatusSuccess, map[string]interface{}{
		"copy_id":  copyID,
		"due_date": dueDate.Format(time.RFC3339),
	})

	return record, nil
}

// ReturnBook checks in a borrowed copy and assesses overdue penalty fines
func (s *Service) ReturnBook(ctx context.Context, tenantID, borrowID string, dailyFineRate float64) (*BorrowRecord, error) {
	record, err := s.repo.GetBorrowRecordByID(ctx, tenantID, borrowID)
	if err != nil {
		return nil, err
	}

	if record.Status == BorrowStatusReturned {
		return nil, ErrBorrowAlreadyReturned
	}

	if dailyFineRate <= 0 {
		dailyFineRate = 5.0 // default ₹5 per day overdue
	}

	now := time.Now().UTC()
	daysOverdue, fineAmount := CalculateOverdueFine(record.DueDate, now, dailyFineRate)

	record.ReturnedAt = &now
	record.Status = BorrowStatusReturned
	record.FineAmount = fineAmount
	record.UpdatedAt = now

	if err := s.repo.UpdateBorrowRecord(ctx, record); err != nil {
		return nil, err
	}

	copy, err := s.repo.GetCopyByID(ctx, tenantID, record.CopyID)
	if err == nil {
		_ = s.repo.UpdateCopyStatus(ctx, tenantID, copy.ID, CopyStatusAvailable)
		if book, bErr := s.repo.GetBookByID(ctx, tenantID, copy.BookID); bErr == nil {
			_ = s.repo.UpdateBookInventory(ctx, tenantID, book.ID, book.TotalCopies, book.AvailableCopies+1)
		}
	}

	s.emitAudit(tenantID, record.StudentID, "USER", "library:book:returned", "library_borrow_record", record.ID, audit.StatusSuccess, map[string]interface{}{
		"days_overdue": daysOverdue,
		"fine_amount":  fineAmount,
	})

	return record, nil
}

// RenewBook extends loan due date if within renewal limits
func (s *Service) RenewBook(ctx context.Context, tenantID, borrowID string, additionalDays int) (*BorrowRecord, error) {
	record, err := s.repo.GetBorrowRecordByID(ctx, tenantID, borrowID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := CanRenewBorrow(record, now, 2); err != nil {
		return nil, err
	}

	if additionalDays <= 0 {
		additionalDays = 14
	}

	record.RenewCount++
	record.DueDate = record.DueDate.AddDate(0, 0, additionalDays)
	record.UpdatedAt = now

	if err := s.repo.UpdateBorrowRecord(ctx, record); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, record.StudentID, "USER", "library:book:renewed", "library_borrow_record", record.ID, audit.StatusSuccess, map[string]interface{}{
		"renew_count":  record.RenewCount,
		"new_due_date": record.DueDate.Format(time.RFC3339),
	})

	return record, nil
}

// ReserveBook places a hold request for a title
func (s *Service) ReserveBook(ctx context.Context, tenantID, bookID, studentID string) (*Reservation, error) {
	if _, err := s.repo.GetBookByID(ctx, tenantID, bookID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	res := &Reservation{
		ID:         fmt.Sprintf("res_%d", time.Now().UnixNano()),
		TenantID:   tenantID,
		BookID:     bookID,
		StudentID:  studentID,
		ReservedAt: now,
		ExpiresAt:  now.AddDate(0, 0, 7), // 7 days pickup window
		Status:     ReservationStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateReservation(ctx, res); err != nil {
		return nil, err
	}

	return res, nil
}

// RecordEbookProgress tracks digital reader progress
func (s *Service) RecordEbookProgress(ctx context.Context, tenantID, bookID, studentID string, pageRead, minutesSpent int) (*EbookAccess, error) {
	access, err := s.repo.GetEbookAccess(ctx, tenantID, bookID, studentID)
	if err != nil {
		return nil, err
	}

	if access == nil {
		access = &EbookAccess{
			ID:                  fmt.Sprintf("ebk_%d", time.Now().UnixNano()),
			TenantID:            tenantID,
			BookID:              bookID,
			StudentID:           studentID,
			LastPageRead:        pageRead,
			TotalReadingMinutes: minutesSpent,
			LastAccessedAt:      time.Now().UTC(),
		}
	} else {
		access.LastPageRead = pageRead
		access.TotalReadingMinutes += minutesSpent
		access.LastAccessedAt = time.Now().UTC()
	}

	if err := s.repo.RecordEbookAccess(ctx, access); err != nil {
		return nil, err
	}

	return access, nil
}

// ListStudentBorrows retrieves loan history for a student
func (s *Service) ListStudentBorrows(ctx context.Context, tenantID, studentID string, status *BorrowStatus) ([]*BorrowRecord, error) {
	return s.repo.ListBorrowRecords(ctx, tenantID, studentID, status)
}
