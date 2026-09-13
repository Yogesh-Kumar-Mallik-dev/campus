/**
 * BLOCK_LIBRARY_ERRORS_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapping for library operations.
 */

package library

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBookNotFound          = errors.New("library book not found")
	ErrISBNExists            = errors.New("book with this ISBN already exists")
	ErrBookCopyNotFound      = errors.New("library book copy not found")
	ErrAccessionExists       = errors.New("book copy accession number already exists")
	ErrCopyUnavailable       = errors.New("book copy is not available for borrowing")
	ErrBorrowLimitExceeded   = errors.New("student has reached maximum allowed active book borrows")
	ErrBorrowRecordNotFound  = errors.New("library borrow record not found")
	ErrBorrowAlreadyReturned = errors.New("borrow record has already been marked returned")
	ErrMaxRenewalsExceeded   = errors.New("maximum allowed book renewals (2) exceeded")
	ErrBorrowOverdueRenewal  = errors.New("overdue books cannot be renewed; please return and settle fine")
	ErrReservationNotFound   = errors.New("book reservation record not found")
	ErrEbookNotFound         = errors.New("digital e-book attachment not found for this title")
)

type DomainError struct {
	Err        error
	HTTPStatus int
	Title      string
	Detail     string
	Type       string
}

func (e *DomainError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Title, e.Detail)
	}
	return e.Title
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewDomainError(err error, status int, title, detail, problemType string) *DomainError {
	return &DomainError{
		Err:        err,
		HTTPStatus: status,
		Title:      title,
		Detail:     detail,
		Type:       problemType,
	}
}

func MapErrorToProblem(err error) *DomainError {
	if err == nil {
		return nil
	}

	var domErr *DomainError
	if errors.As(err, &domErr) {
		return domErr
	}

	switch {
	case errors.Is(err, ErrBookNotFound), errors.Is(err, ErrBookCopyNotFound),
		errors.Is(err, ErrBorrowRecordNotFound), errors.Is(err, ErrReservationNotFound),
		errors.Is(err, ErrEbookNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrISBNExists), errors.Is(err, ErrAccessionExists),
		errors.Is(err, ErrCopyUnavailable), errors.Is(err, ErrBorrowLimitExceeded),
		errors.Is(err, ErrBorrowAlreadyReturned):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict / Policy Guard", err.Error(), "https://campus.internal/errors/conflict")

	case errors.Is(err, ErrMaxRenewalsExceeded), errors.Is(err, ErrBorrowOverdueRenewal):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected library error occurred", "https://campus.internal/errors/internal")
	}
}
