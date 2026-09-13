/**
 * BLOCK_ONBOARDING_ERRORS_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   Authoritative domain error types and sentinels for onboarding operations.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"errors"
	"fmt"
)

var (
	ErrApplicantNotFound         = errors.New("applicant not found")
	ErrDocumentNotFound          = errors.New("document not found")
	ErrTenantRequired            = errors.New("tenant_id is required")
	ErrInvalidStateTransition    = errors.New("invalid onboarding state transition")
	ErrMissingMandatoryDocs      = errors.New("missing mandatory KYC documents for submission or verification")
	ErrUnverifiedDocuments       = errors.New("cannot verify application with pending or rejected documents")
	ErrDuplicateEmail            = errors.New("an applicant with this email is already registered")
	ErrDuplicateRollNumber       = errors.New("roll number collision detected")
	ErrDuplicateEmployeeID       = errors.New("employee ID collision detected")
	ErrDepartmentNotFound        = errors.New("academic department not found")
	ErrProgramNotFound           = errors.New("academic program not found")
	ErrCohortNotFound            = errors.New("cohort batch not found")
	ErrCohortFull                = errors.New("cohort batch has reached maximum capacity")
	ErrInvalidInput              = errors.New("invalid onboarding input parameters")
	ErrUnauthorizedReviewer      = errors.New("reviewer_id is required")
	ErrAlreadyEnrolled           = errors.New("applicant is already enrolled")
	ErrRejectionReasonRequired   = errors.New("rejection reason is required")
)

// DomainError captures contextual error data for onboarding operations.
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError constructs a structured onboarding DomainError.
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
