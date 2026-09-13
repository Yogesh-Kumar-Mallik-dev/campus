/**
 * BLOCK_AUTH_DOMAIN_MODELS_001
 * Purpose: Authoritative domain entities, value objects, and capability sets for IAM.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"strings"
	"time"
)

// AccountTier represents the macro administrative classification level.
type AccountTier string

const (
	TierSuperAdmin AccountTier = "SUPER_ADMIN"
	TierAdmin      AccountTier = "ADMIN"
	TierUser       AccountTier = "USER"
)

// UserStatus represents the lifecycle state of a user identity.
type UserStatus string

const (
	StatusActive              UserStatus = "ACTIVE"
	StatusSuspended           UserStatus = "SUSPENDED"
	StatusPendingVerification UserStatus = "PENDING_VERIFICATION"
	StatusLocked              UserStatus = "LOCKED"
)

// MFAType defines the multi-factor authentication mechanism.
type MFAType string

const (
	MFATypeTOTP        MFAType = "TOTP"
	MFATypeBackupCodes MFAType = "BACKUP_CODES"
)

// User represents a central authenticated identity.
type User struct {
	ID                  string      `json:"id"`
	TenantID            string      `json:"tenant_id"`
	Tier                AccountTier `json:"tier"`
	Username            string      `json:"username"`
	Email               string      `json:"email"`
	EmailNormalized     string      `json:"email_normalized"`
	PasswordHash        string      `json:"-"`
	Status              UserStatus  `json:"status"`
	FailedLoginAttempts int         `json:"failed_login_attempts"`
	LockoutUntil        *time.Time  `json:"lockout_until,omitempty"`
	LastLoginAt         *time.Time  `json:"last_login_at,omitempty"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`

	// Attached Personas and direct capabilities
	Roles       []*UserRole `json:"roles,omitempty"`
	Permissions []string    `json:"permissions,omitempty"`
}

// IsActive returns true if the user account is in normal active status.
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsLockedOut checks if the user is temporarily locked due to consecutive failed attempts.
func (u *User) IsLockedOut(now time.Time) bool {
	if u.Status == StatusLocked {
		return true
	}
	if u.LockoutUntil != nil && now.Before(*u.LockoutUntil) {
		return true
	}
	return false
}

// Role represents a functional persona (e.g. Student, Faculty, Warden).
type Role struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission represents a fine-grained capability key (e.g. hostel:gatepass:approve).
type Permission struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Domain      string    `json:"domain"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRole maps a user to a persona role with an optional domain scope.
type UserRole struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	RoleKey   string    `json:"role_key"`
	RoleName  string    `json:"role_name"`
	ScopeType string    `json:"scope_type,omitempty"` // e.g. "DEPT", "HOSTEL", "CAMPUS"
	ScopeID   string    `json:"scope_id,omitempty"`   // e.g. "mech_dept_01", "block_b"
	CreatedAt time.Time `json:"created_at"`
}

// Session represents an active authenticated device session.
type Session struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	TenantID          string     `json:"tenant_id"`
	TokenFamilyID     string     `json:"token_family_id"`
	CurrentTokenHash  string     `json:"-"`
	IPAddress         string     `json:"ip_address,omitempty"`
	UserAgent         string     `json:"user_agent,omitempty"`
	DeviceFingerprint string     `json:"device_fingerprint,omitempty"`
	ExpiresAt         time.Time  `json:"expires_at"`
	RevokedAt         *time.Time `json:"revoked_at,omitempty"`
	RevokedReason     string     `json:"revoked_reason,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// IsActive returns true if the session has not expired and has not been revoked.
func (s *Session) IsActive(now time.Time) bool {
	return s.RevokedAt == nil && now.Before(s.ExpiresAt)
}

// RefreshTokenFamily tracks a sequence of single-use refresh tokens for a session.
type RefreshTokenFamily struct {
	FamilyID        string     `json:"family_id"`
	UserID          string     `json:"user_id"`
	ActiveTokenHash string     `json:"-"`
	RotationCount   int        `json:"rotation_count"`
	IsCompromised   bool       `json:"is_compromised"`
	RevokedAt       *time.Time `json:"revoked_at,omitempty"`
	RevokedReason   string     `json:"revoked_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// MFAEnrollment stores two-factor configuration for a user.
type MFAEnrollment struct {
	UserID            string    `json:"user_id"`
	MFAType           MFAType   `json:"mfa_type"`
	SecretEncrypted   string    `json:"-"`
	BackupCodesHashed string    `json:"-"`
	IsVerified        bool      `json:"is_verified"`
	EnrolledAt        time.Time `json:"enrolled_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Claims represents the verified payload inside a short-lived access token.
type Claims struct {
	UserID      string      `json:"sub"`
	TenantID    string      `json:"tenant_id"`
	Tier        AccountTier `json:"tier"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	Roles       []string    `json:"roles"`
	Permissions []string    `json:"permissions"`
	SessionID   string      `json:"session_id"`
	FamilyID    string      `json:"family_id"`
	IssuedAt    int64       `json:"iat"`
	ExpiresAt   int64       `json:"exp"`
}

// HasPermission performs an O(1) or linear capability check.
func (c *Claims) HasPermission(requiredPerm string) bool {
	if c.Tier == TierSuperAdmin {
		return true
	}
	for _, p := range c.Permissions {
		if p == requiredPerm || p == "*" || (strings.HasSuffix(p, ":*") && strings.HasPrefix(requiredPerm, strings.TrimSuffix(p, "*"))) {
			return true
		}
	}
	return false
}

// HasRole checks if the claims include a specific persona role key.
func (c *Claims) HasRole(roleKey string) bool {
	if c.Tier == TierSuperAdmin {
		return true
	}
	for _, r := range c.Roles {
		if r == roleKey {
			return true
		}
	}
	return false
}

// TokenPair represents the access and refresh token returned upon successful authentication.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // Seconds
}

// AuthResult represents the complete response returned by Login / Session bootstrap.
type AuthResult struct {
	TokenPair        *TokenPair `json:"tokens,omitempty"`
	User             *User      `json:"user,omitempty"`
	RequiresMFA      bool       `json:"requires_mfa"`
	MFATicket        string     `json:"mfa_ticket,omitempty"`
	AllowedWorkspaces []string  `json:"allowed_workspaces,omitempty"`
}
