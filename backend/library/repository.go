/**
 * BLOCK_LIBRARY_REPOSITORY_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   Repository interface contracts for bibliographic catalog, copies, loans, and reservations.
 */

package library

import (
	"context"
)

type Repository interface {
	// Book operations
	CreateBook(ctx context.Context, book *Book) error
	GetBookByID(ctx context.Context, tenantID, id string) (*Book, error)
	GetBookByISBN(ctx context.Context, tenantID, isbn string) (*Book, error)
	ListBooks(ctx context.Context, tenantID string, category *BookCategory, query string) ([]*Book, error)
	UpdateBookInventory(ctx context.Context, tenantID, bookID string, totalCopies, availableCopies int) error

	// Copy operations
	CreateCopy(ctx context.Context, copy *BookCopy) error
	GetCopyByID(ctx context.Context, tenantID, id string) (*BookCopy, error)
	GetCopyByAccession(ctx context.Context, tenantID, accessionNumber string) (*BookCopy, error)
	ListCopiesByBook(ctx context.Context, tenantID, bookID string) ([]*BookCopy, error)
	UpdateCopyStatus(ctx context.Context, tenantID, id string, status BookCopyStatus) error

	// Borrow operations
	CreateBorrowRecord(ctx context.Context, record *BorrowRecord) error
	GetBorrowRecordByID(ctx context.Context, tenantID, id string) (*BorrowRecord, error)
	GetActiveBorrowsByStudent(ctx context.Context, tenantID, studentID string) ([]*BorrowRecord, error)
	UpdateBorrowRecord(ctx context.Context, record *BorrowRecord) error
	ListBorrowRecords(ctx context.Context, tenantID, studentID string, status *BorrowStatus) ([]*BorrowRecord, error)

	// Reservation operations
	CreateReservation(ctx context.Context, reservation *Reservation) error
	GetReservationByID(ctx context.Context, tenantID, id string) (*Reservation, error)
	ListReservationsByBook(ctx context.Context, tenantID, bookID string) ([]*Reservation, error)
	UpdateReservationStatus(ctx context.Context, tenantID, id string, status ReservationStatus) error

	// Ebook Access operations
	RecordEbookAccess(ctx context.Context, access *EbookAccess) error
	GetEbookAccess(ctx context.Context, tenantID, bookID, studentID string) (*EbookAccess, error)
}
