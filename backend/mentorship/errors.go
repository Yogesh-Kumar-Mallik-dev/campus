/**
 * BLOCK_MENTORSHIP_ERRORS_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapping for mentorship operations.
 */

package mentorship

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrAllocationNotFound        = errors.New("mentor allocation not found")
	ErrActiveAllocationExists    = errors.New("student already has an active mentor allocation in this cohort")
	ErrSessionNotFound           = errors.New("mentorship session not found")
	ErrSessionAlreadyCompleted   = errors.New("mentorship session is already completed and cannot be modified")
	ErrInvalidSessionTransition  = errors.New("invalid session state transition; only scheduled sessions can be completed or cancelled")
	ErrProgressRecordNotFound    = errors.New("student academic progress record not found")
	ErrProgressRecordExists      = errors.New("academic progress record already evaluated for this semester")
	ErrAlertNotFound             = errors.New("at-risk alert not found")
	ErrAlertAlreadyResolved      = errors.New("at-risk alert is already marked resolved")
	ErrInvalidScoreBounds        = errors.New("GPA values must be between 0.0 and 10.0, and attendance between 0% and 100%")
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
	case errors.Is(err, ErrAllocationNotFound), errors.Is(err, ErrSessionNotFound),
		errors.Is(err, ErrProgressRecordNotFound), errors.Is(err, ErrAlertNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrActiveAllocationExists), errors.Is(err, ErrProgressRecordExists),
		errors.Is(err, ErrAlertAlreadyResolved), errors.Is(err, ErrSessionAlreadyCompleted):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict / Policy Guard", err.Error(), "https://campus.internal/errors/conflict")

	case errors.Is(err, ErrInvalidSessionTransition), errors.Is(err, ErrInvalidScoreBounds):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected mentorship error occurred", "https://campus.internal/errors/internal")
	}
}
