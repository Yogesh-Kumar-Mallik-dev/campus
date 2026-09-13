/**
 * BLOCK_LIBRARY_DOMAIN_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   Domain entities, checkout state machines, renewal limits, and overdue fine calculation.
 */

package library

import (
	"math"
	"time"
)

type BookCategory string

const (
	CategoryComputerScience BookCategory = "COMPUTER_SCIENCE"
	CategoryElectronics     BookCategory = "ELECTRONICS"
	CategoryMechanical      BookCategory = "MECHANICAL"
	CategoryMathematics     BookCategory = "MATHEMATICS"
	CategoryPhysics         BookCategory = "PHYSICS"
	CategoryLiterature      BookCategory = "LITERATURE"
	CategoryManagement      BookCategory = "MANAGEMENT"
	CategoryGeneral         BookCategory = "GENERAL"
)

type BookCopyStatus string

const (
	CopyStatusAvailable   BookCopyStatus = "AVAILABLE"
	CopyStatusIssued      BookCopyStatus = "ISSUED"
	CopyStatusReserved    BookCopyStatus = "RESERVED"
	CopyStatusMaintenance BookCopyStatus = "MAINTENANCE"
	CopyStatusLost        BookCopyStatus = "LOST"
)

type BorrowStatus string

const (
	BorrowStatusIssued   BorrowStatus = "ISSUED"
	BorrowStatusReturned BorrowStatus = "RETURNED"
	BorrowStatusOverdue  BorrowStatus = "OVERDUE"
	BorrowStatusLost     BorrowStatus = "LOST"
)

type ReservationStatus string

const (
	ReservationStatusPending        ReservationStatus = "PENDING"
	ReservationStatusReadyForPickup ReservationStatus = "READY_FOR_PICKUP"
	ReservationStatusFulfilled      ReservationStatus = "FULFILLED"
	ReservationStatusExpired        ReservationStatus = "EXPIRED"
	ReservationStatusCancelled      ReservationStatus = "CANCELLED"
)

// Book represents a bibliographic catalog title
type Book struct {
	ID              string       `json:"id"`
	TenantID        string       `json:"tenant_id"`
	ISBN            string       `json:"isbn"`
	Title           string       `json:"title"`
	Author          string       `json:"author"`
	Publisher       string       `json:"publisher"`
	Edition         string       `json:"edition,omitempty"`
	Category        BookCategory `json:"category"`
	TotalCopies     int          `json:"total_copies"`
	AvailableCopies int          `json:"available_copies"`
	ShelfLocation   string       `json:"shelf_location"`
	EbookKey        *string      `json:"ebook_key,omitempty"`
	EbookFormat     *string      `json:"ebook_format,omitempty"`
	EbookSize       *int         `json:"ebook_size,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

// BookCopy represents an individual physical physical accession unit
type BookCopy struct {
	ID              string         `json:"id"`
	TenantID        string         `json:"tenant_id"`
	BookID          string         `json:"book_id"`
	AccessionNumber string         `json:"accession_number"`
	Barcode         string         `json:"barcode"`
	Status          BookCopyStatus `json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// BorrowRecord represents an active or historic loan transaction
type BorrowRecord struct {
	ID         string       `json:"id"`
	TenantID   string       `json:"tenant_id"`
	CopyID     string       `json:"copy_id"`
	StudentID  string       `json:"student_id"`
	BorrowedAt time.Time    `json:"borrowed_at"`
	DueDate    time.Time    `json:"due_date"`
	ReturnedAt *time.Time   `json:"returned_at,omitempty"`
	RenewCount int          `json:"renew_count"`
	Status     BorrowStatus `json:"status"`
	FineAmount float64      `json:"fine_amount"`
	FinePaid   bool         `json:"fine_paid"`
	IssuedByID *string      `json:"issued_by_id,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// Reservation represents a title hold request in the reservation queue
type Reservation struct {
	ID         string            `json:"id"`
	TenantID   string            `json:"tenant_id"`
	BookID     string            `json:"book_id"`
	StudentID  string            `json:"student_id"`
	ReservedAt time.Time         `json:"reserved_at"`
	ExpiresAt  time.Time         `json:"expires_at"`
	Status     ReservationStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// EbookAccess represents digital e-reader telemetry and progress
type EbookAccess struct {
	ID                  string    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	BookID              string    `json:"book_id"`
	StudentID           string    `json:"student_id"`
	LastPageRead        int       `json:"last_page_read"`
	TotalReadingMinutes int       `json:"total_reading_minutes"`
	LastAccessedAt      time.Time `json:"last_accessed_at"`
}

// CalculateOverdueFine computes penalty fees for overdue returns
func CalculateOverdueFine(dueDate, returnedAt time.Time, dailyRate float64) (int, float64) {
	if !returnedAt.After(dueDate) {
		return 0, 0.0
	}

	diffHours := returnedAt.Sub(dueDate).Hours()
	daysOverdue := int(math.Ceil(diffHours / 24.0))
	if daysOverdue <= 0 {
		return 0, 0.0
	}

	fineAmount := float64(daysOverdue) * dailyRate
	return daysOverdue, fineAmount
}

// ValidateBorrowEligibility ensures student has not exceeded active quota
func ValidateBorrowEligibility(activeCount, maxAllowed int) error {
	if activeCount >= maxAllowed {
		return ErrBorrowLimitExceeded
	}
	return nil
}

// CanRenewBorrow validates renewal limits and overdue status
func CanRenewBorrow(record *BorrowRecord, now time.Time, maxRenewals int) error {
	if record.RenewCount >= maxRenewals {
		return ErrMaxRenewalsExceeded
	}
	if now.After(record.DueDate) {
		return ErrBorrowOverdueRenewal
	}
	return nil
}
