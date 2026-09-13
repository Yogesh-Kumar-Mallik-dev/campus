/**
 * BLOCK_WHISTLEBLOWER_ERRORS_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapping for whistleblower operations.
 */

package whistleblower

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrReportNotFound               = errors.New("whistleblower report not found or invalid tracking token")
	ErrInvalidTrackingToken          = errors.New("tracking token format is invalid")
	ErrReportAlreadyClosed           = errors.New("whistleblower report is already resolved or closed")
	ErrInvalidReportTransition       = errors.New("invalid whistleblower report status transition")
	ErrUnauthorizedCommitteeAccess   = errors.New("unauthorized committee access to confidential whistleblower report")
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
	case errors.Is(err, ErrReportNotFound):
		return NewDomainError(err, http.StatusNotFound, "Report Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrUnauthorizedCommitteeAccess):
		return NewDomainError(err, http.StatusForbidden, "Forbidden", err.Error(), "https://campus.internal/errors/forbidden")

	case errors.Is(err, ErrInvalidTrackingToken), errors.Is(err, ErrInvalidReportTransition),
		errors.Is(err, ErrReportAlreadyClosed):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected whistleblower system error occurred", "https://campus.internal/errors/internal")
	}
}
