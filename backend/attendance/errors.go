/**
 * BLOCK_ATTENDANCE_ERRORS_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   Authoritative domain error types and sentinels for attendance operations.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package attendance

import (
	"errors"
	"fmt"
)

var (
	ErrSessionNotFound        = errors.New("attendance session not found")
	ErrSessionLocked          = errors.New("attendance session is locked or finalized")
	ErrSessionNotOpen         = errors.New("attendance session is not currently open for marking")
	ErrInvalidStateTransition = errors.New("invalid attendance session state transition")
	ErrSubjectNotFound        = errors.New("course subject not found")
	ErrSlotNotFound           = errors.New("timetable slot not found")
	ErrStudentNotFound        = errors.New("student not found in cohort")
	ErrMedicalLeaveNotFound   = errors.New("medical leave application not found")
	ErrTenantRequired         = errors.New("tenant_id is required")
	ErrInvalidInput           = errors.New("invalid attendance input parameters")
	ErrUnauthorizedMarker     = errors.New("unauthorized attendance marker")
	ErrDuplicateSession       = errors.New("attendance session for slot/date already exists")
)

// DomainError captures contextual error data for attendance operations.
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

// NewDomainError constructs a structured attendance DomainError.
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
