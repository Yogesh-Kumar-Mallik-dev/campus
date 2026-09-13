/**
 * BLOCK_HELPDESK_ERRORS_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapping for helpdesk operations.
 */

package helpdesk

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrCategoryNotFound          = errors.New("helpdesk category not found")
	ErrCategoryCodeDuplicate     = errors.New("helpdesk category code already exists for tenant")
	ErrTicketNotFound            = errors.New("helpdesk ticket not found")
	ErrInvalidTicketTransition   = errors.New("invalid helpdesk ticket state transition")
	ErrTicketClosed              = errors.New("ticket is closed and cannot accept new replies")
	ErrInternalNoteUnauthorized  = errors.New("only staff members can create internal notes")
	ErrInvalidRatingScore        = errors.New("rating score must be an integer between 1 and 5 stars")
	ErrRatingNotAllowed          = errors.New("rating and feedback can only be submitted for resolved or closed tickets")
	ErrEscalationAlreadyPending  = errors.New("ticket already has an active pending escalation")
	ErrUnauthorizedTicketAccess  = errors.New("unauthorized access to helpdesk ticket")
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
	case errors.Is(err, ErrCategoryNotFound), errors.Is(err, ErrTicketNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrCategoryCodeDuplicate), errors.Is(err, ErrEscalationAlreadyPending):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict", err.Error(), "https://campus.internal/errors/conflict")

	case errors.Is(err, ErrInternalNoteUnauthorized), errors.Is(err, ErrUnauthorizedTicketAccess):
		return NewDomainError(err, http.StatusForbidden, "Forbidden Action", err.Error(), "https://campus.internal/errors/forbidden")

	case errors.Is(err, ErrInvalidTicketTransition), errors.Is(err, ErrTicketClosed),
		errors.Is(err, ErrInvalidRatingScore), errors.Is(err, ErrRatingNotAllowed):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected helpdesk error occurred", "https://campus.internal/errors/internal")
	}
}
