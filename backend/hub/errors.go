/**
 * BLOCK_HUB_ERRORS_001
 * Subsystem: Rank 17 - The Hub Root Super-App (hub)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapper for Hub cockpit operations.
 */

package hub

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidPersona      = errors.New("invalid or unsupported persona type")
	ErrDashboardNotFound   = errors.New("persona dashboard not found")
	ErrWidgetNotFound      = errors.New("widget configuration not found")
	ErrShortcutNotFound    = errors.New("quick action shortcut not found")
	ErrUnauthorizedAccess  = errors.New("unauthorized access to persona dashboard")
	ErrInvalidWidgetLayout = errors.New("invalid widget layout or size parameters")
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

// MapErrorToProblem converts domain errors into RFC 7807 problem details.
func MapErrorToProblem(err error) *DomainError {
	if err == nil {
		return nil
	}

	var domErr *DomainError
	if errors.As(err, &domErr) {
		return domErr
	}

	switch {
	case errors.Is(err, ErrInvalidPersona), errors.Is(err, ErrInvalidWidgetLayout):
		return NewDomainError(err, http.StatusBadRequest, "Bad Request", err.Error(), "https://campus.internal/errors/bad-request")

	case errors.Is(err, ErrDashboardNotFound), errors.Is(err, ErrWidgetNotFound), errors.Is(err, ErrShortcutNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrUnauthorizedAccess):
		return NewDomainError(err, http.StatusForbidden, "Forbidden", err.Error(), "https://campus.internal/errors/forbidden")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected hub error occurred", "https://campus.internal/errors/internal")
	}
}
