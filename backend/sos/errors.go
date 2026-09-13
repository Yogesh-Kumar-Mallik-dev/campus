/**
 * BLOCK_SOS_ERRORS_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapper for SOS operations.
 */

package sos

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrIncidentNotFound           = errors.New("emergency SOS incident not found")
	ErrInvalidCoordinates         = errors.New("invalid geographical coordinates; latitude must be [-90, 90] and longitude [-180, 180]")
	ErrInvalidIncidentTransition  = errors.New("invalid SOS incident state transition")
	ErrIncidentAlreadyClosed      = errors.New("SOS incident is already closed or resolved")
	ErrResponderAlreadyDispatched = errors.New("responder has already been dispatched to this emergency incident")
	ErrResponderNotFound          = errors.New("dispatch responder mission not found")
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
	case errors.Is(err, ErrIncidentNotFound), errors.Is(err, ErrResponderNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrResponderAlreadyDispatched):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict", err.Error(), "https://campus.internal/errors/conflict")

	case errors.Is(err, ErrInvalidCoordinates), errors.Is(err, ErrInvalidIncidentTransition),
		errors.Is(err, ErrIncidentAlreadyClosed):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected emergency response error occurred", "https://campus.internal/errors/internal")
	}
}
