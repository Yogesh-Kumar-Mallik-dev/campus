/**
 * BLOCK_PORTAL_ERRORS_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapper for portal operations.
 */

package portal

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrLandingPageNotFound = errors.New("portal landing page configuration not found")
	ErrProgramNotFound     = errors.New("program catalog entry not found")
	ErrDuplicateProgramCode = errors.New("program code already exists in catalog")
	ErrInquiryNotFound     = errors.New("prospect inquiry not found")
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
	case errors.Is(err, ErrLandingPageNotFound), errors.Is(err, ErrProgramNotFound), errors.Is(err, ErrInquiryNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrDuplicateProgramCode):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict", err.Error(), "https://campus.internal/errors/conflict")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected public portal error occurred", "https://campus.internal/errors/internal")
	}
}
