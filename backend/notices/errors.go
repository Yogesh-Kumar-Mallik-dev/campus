/**
 * BLOCK_NOTICES_ERRORS_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   Authoritative domain error types and sentinels for notices, broadcasts, and acknowledgements.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package notices

import (
	"errors"
	"fmt"
)

var (
	ErrNoticeNotFound         = errors.New("campus notice not found")
	ErrAttachmentNotFound     = errors.New("notice attachment not found")
	ErrTenantRequired         = errors.New("tenant_id is required")
	ErrInvalidInput           = errors.New("invalid notice input parameters")
	ErrInvalidStateTransition = errors.New("invalid notice state transition")
	ErrNoticeAlreadyPublished = errors.New("notice is already published")
	ErrNoticeExpired          = errors.New("notice has expired")
	ErrDuplicateNoticeSlug    = errors.New("notice with this slug already exists")
	ErrUnauthorizedAuthor     = errors.New("unauthorized notice author")
)

// DomainError captures contextual error data for notice operations.
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

// NewDomainError constructs a structured notice DomainError.
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
