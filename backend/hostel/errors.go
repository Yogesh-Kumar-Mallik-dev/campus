/**
 * BLOCK_HOSTEL_ERRORS_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapping for hostel operations.
 */

package hostel

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBlockNotFound           = errors.New("hostel block not found")
	ErrBlockCodeExists         = errors.New("hostel block code already exists")
	ErrRoomNotFound            = errors.New("hostel room not found")
	ErrRoomNumberExists        = errors.New("room number already exists in block")
	ErrBedNotFound             = errors.New("hostel bed not found")
	ErrBedNumberExists         = errors.New("bed number already exists in room")
	ErrBedUnavailable          = errors.New("hostel bed is not available for allocation")
	ErrStudentAlreadyAllocated = errors.New("student already has an active hostel bed allocation")
	ErrAllocationNotFound      = errors.New("hostel bed allocation not found")
	ErrGatePassNotFound        = errors.New("hostel gate pass not found")
	ErrInvalidGatePassState    = errors.New("invalid gate pass state transition")
	ErrInvalidTimeRange        = errors.New("expected in time must be after expected out time")
	ErrIncidentNotFound        = errors.New("hostel incident record not found")
	ErrUnauthorizedWarden      = errors.New("actor is not authorized as a warden for this block")
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
	case errors.Is(err, ErrBlockNotFound), errors.Is(err, ErrRoomNotFound),
		errors.Is(err, ErrBedNotFound), errors.Is(err, ErrAllocationNotFound),
		errors.Is(err, ErrGatePassNotFound), errors.Is(err, ErrIncidentNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrBlockCodeExists), errors.Is(err, ErrRoomNumberExists),
		errors.Is(err, ErrBedNumberExists), errors.Is(err, ErrBedUnavailable),
		errors.Is(err, ErrStudentAlreadyAllocated):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict", err.Error(), "https://campus.internal/errors/conflict")

	case errors.Is(err, ErrInvalidGatePassState), errors.Is(err, ErrInvalidTimeRange):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	case errors.Is(err, ErrUnauthorizedWarden):
		return NewDomainError(err, http.StatusForbidden, "Forbidden Access", err.Error(), "https://campus.internal/errors/forbidden")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected error occurred in hostel management", "https://campus.internal/errors/internal")
	}
}
