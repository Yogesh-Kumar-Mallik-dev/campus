/**
 * BLOCK_EVENTS_ERRORS_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapping for event lifecycle, venue bookings, and ticketing.
 */

package events

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrEventNotFound            = errors.New("campus event not found")
	ErrEventNotPublished        = errors.New("event is not published and not accepting registrations")
	ErrRegistrationClosed       = errors.New("registration deadline has passed for this event")
	ErrEventCapacityExceeded    = errors.New("event capacity or maximum ticket quota has been reached")
	ErrTicketNotFound           = errors.New("event ticket not found")
	ErrTicketAlreadyCheckedIn   = errors.New("ticket has already been checked in; duplicate gate punch rejected")
	ErrTicketCancelled          = errors.New("ticket has been cancelled and cannot be used for entry")
	ErrDuplicateTicketBooking   = errors.New("student has already booked a confirmed ticket for this event")
	ErrVenueConflict            = errors.New("venue is already booked by another event during the requested time window")
	ErrInvalidEventTimeRange    = errors.New("event end time must be after start time, and registration deadline must precede event end")
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
	case errors.Is(err, ErrEventNotFound), errors.Is(err, ErrTicketNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrVenueConflict), errors.Is(err, ErrTicketAlreadyCheckedIn),
		errors.Is(err, ErrDuplicateTicketBooking), errors.Is(err, ErrEventCapacityExceeded):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict / Policy Guard", err.Error(), "https://campus.internal/errors/conflict")

	case errors.Is(err, ErrEventNotPublished), errors.Is(err, ErrRegistrationClosed),
		errors.Is(err, ErrTicketCancelled), errors.Is(err, ErrInvalidEventTimeRange):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected event organization error occurred", "https://campus.internal/errors/internal")
	}
}
