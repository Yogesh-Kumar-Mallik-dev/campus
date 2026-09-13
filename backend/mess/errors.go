/**
 * BLOCK_MESS_ERRORS_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapping for campus dining operations.
 */

package mess

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrMessHallNotFound         = errors.New("mess hall not found")
	ErrMessHallCodeExists       = errors.New("mess hall code already exists")
	ErrMenuItemNotFound         = errors.New("mess menu item not found")
	ErrSubscriptionNotFound     = errors.New("mess subscription not found")
	ErrActiveSubscriptionExists = errors.New("student already has an active mess subscription")
	ErrTokenNotFound            = errors.New("dining token not found")
	ErrTokenAlreadyRedeemed     = errors.New("dining token has already been redeemed for this meal")
	ErrTokenExpired             = errors.New("dining token has expired")
	ErrDoubleRedemptionAttempt  = errors.New("student already redeemed dining token for this meal slot today")
	ErrRebateNotFound           = errors.New("mess rebate application not found")
	ErrInsufficientRebateDays   = errors.New("mess rebate requires a minimum absence of 3 continuous days")
	ErrInvalidDietPreference    = errors.New("invalid dietary preference specified")
	ErrInvalidDateRange         = errors.New("end date must be after start date")
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
	case errors.Is(err, ErrMessHallNotFound), errors.Is(err, ErrMenuItemNotFound),
		errors.Is(err, ErrSubscriptionNotFound), errors.Is(err, ErrTokenNotFound),
		errors.Is(err, ErrRebateNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrMessHallCodeExists), errors.Is(err, ErrActiveSubscriptionExists),
		errors.Is(err, ErrTokenAlreadyRedeemed), errors.Is(err, ErrDoubleRedemptionAttempt):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict / Double Punch", err.Error(), "https://campus.internal/errors/conflict")

	case errors.Is(err, ErrInsufficientRebateDays), errors.Is(err, ErrInvalidDietPreference),
		errors.Is(err, ErrInvalidDateRange), errors.Is(err, ErrTokenExpired):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected dining error occurred", "https://campus.internal/errors/internal")
	}
}
