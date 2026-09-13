/**
 * BLOCK_AUDIT_ERRORS_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Authoritative domain error types and sentinels for audit operations.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"errors"
	"fmt"
)

var (
	ErrAuditNotFound              = errors.New("audit record not found")
	ErrTenantRequired             = errors.New("tenant_id is required")
	ErrActionRequired             = errors.New("action is required")
	ErrResourceTypeRequired       = errors.New("resource_type is required")
	ErrChainIntegrityCompromised  = errors.New("audit hash chain integrity verification failed")
	ErrInvalidFilter              = errors.New("invalid audit filter criteria")
	ErrQueueFull                  = errors.New("audit ingestion queue is saturated")
)

// DomainError captures contextual error data for auditing operations.
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

// NewDomainError constructs a structured audit DomainError.
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
