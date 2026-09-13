/**
 * BLOCK_AUTH_ERRORS_001
 * Purpose: Domain error definitions and RFC 7807 problem details status mappings.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials       = errors.New("AUTH_INVALID_CREDENTIALS: username/email or password is incorrect")
	ErrAccountLocked            = errors.New("AUTH_ACCOUNT_LOCKED: account is temporarily locked due to consecutive failed attempts")
	ErrAccountSuspended         = errors.New("AUTH_ACCOUNT_SUSPENDED: account has been suspended by administration")
	ErrAccountPending           = errors.New("AUTH_ACCOUNT_PENDING: account is pending verification")
	ErrUserNotFound             = errors.New("AUTH_USER_NOT_FOUND: requested user identity does not exist")
	ErrTenantNotFound           = errors.New("AUTH_TENANT_NOT_FOUND: requested institution tenant does not exist")
	ErrMFARequired              = errors.New("AUTH_MFA_REQUIRED: multi-factor authentication passcode required")
	ErrInvalidMFACode           = errors.New("AUTH_INVALID_MFA_CODE: provided TOTP passcode is invalid or expired")
	ErrMFAAlreadyEnrolled       = errors.New("AUTH_MFA_ALREADY_ENROLLED: MFA is already active on this account")
	ErrMFANotEnrolled           = errors.New("AUTH_MFA_NOT_ENROLLED: MFA has not been configured on this account")
	ErrTokenExpired             = errors.New("AUTH_TOKEN_EXPIRED: refresh token is expired")
	ErrTokenReused              = errors.New("AUTH_TOKEN_REUSED: refresh token replay breach detected; session terminated")
	ErrSessionRevoked           = errors.New("AUTH_SESSION_REVOKED: session has been terminated")
	ErrSessionNotFound          = errors.New("AUTH_SESSION_NOT_FOUND: session record does not exist")
	ErrDuplicateEmail           = errors.New("AUTH_DUPLICATE_EMAIL: a user with this email address already exists in this tenant")
	ErrDuplicateUsername        = errors.New("AUTH_DUPLICATE_USERNAME: a user with this username already exists in this tenant")
	ErrUnauthorized             = errors.New("AUTH_UNAUTHORIZED: authentication token is missing or invalid")
	ErrForbidden                = errors.New("AUTH_FORBIDDEN: caller lacks required capability permissions")
	ErrInvalidPasswordStrength  = errors.New("AUTH_INVALID_PASSWORD_STRENGTH: password does not meet security complexity criteria")
)

// DomainError wraps standard errors with contextual values and modular block IDs.
type DomainError struct {
	BlockID string
	Code    string
	Status  int
	Detail  string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s (%d): %s - %v", e.BlockID, e.Code, e.Status, e.Detail, e.Err)
	}
	return fmt.Sprintf("[%s] %s (%d): %s", e.BlockID, e.Code, e.Status, e.Detail)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates an RFC 7807 compatible domain error.
func NewDomainError(blockID, code string, status int, detail string, originalErr error) *DomainError {
	return &DomainError{
		BlockID: blockID,
		Code:    code,
		Status:  status,
		Detail:  detail,
		Err:     originalErr,
	}
}
