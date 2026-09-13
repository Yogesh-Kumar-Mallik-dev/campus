/**
 * BLOCK_AUTH_REPOSITORY_INTERFACES_001
 * Purpose: Outbound Dependency Injection repository and persistence contracts.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"context"
	"time"
)

// UserRepository handles user entity retrieval, persistence, and brute-force lockout counters.
type UserRepository interface {
	GetByID(ctx context.Context, tenantID, userID string) (*User, error)
	GetByEmailNormalized(ctx context.Context, tenantID, emailNormalized string) (*User, error)
	GetByUsername(ctx context.Context, tenantID, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	RecordFailedLogin(ctx context.Context, tenantID, userID string, lockoutDuration time.Duration, maxAttempts int) error
	ResetFailedLogins(ctx context.Context, tenantID, userID string) error
	UpdateLastLogin(ctx context.Context, tenantID, userID string, loginTime time.Time) error
}

// SessionRepository manages active session records and refresh token family lifecycles.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session, family *RefreshTokenFamily) error
	GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
	GetTokenFamily(ctx context.Context, familyID string) (*RefreshTokenFamily, error)
	RotateTokenFamily(ctx context.Context, familyID string, oldHash, newHash string, newRotationCount int) error
	RevokeTokenFamily(ctx context.Context, familyID string, reason string) error
	RevokeSession(ctx context.Context, sessionID string, reason string) error
	RevokeAllUserSessions(ctx context.Context, tenantID, userID string, reason string) error
	ListUserSessions(ctx context.Context, tenantID, userID string) ([]*Session, error)
}

// RolePermissionRepository resolves assigned roles, persona scopes, and effective permissions.
type RolePermissionRepository interface {
	GetUserRoles(ctx context.Context, userID string) ([]*UserRole, error)
	GetEffectivePermissions(ctx context.Context, userID string) ([]string, error)
	AssignRoleToUser(ctx context.Context, userRole *UserRole) error
	RemoveRoleFromUser(ctx context.Context, userRoleID string) error
	GetRoleByKey(ctx context.Context, tenantID, roleKey string) (*Role, error)
}

// MFARepository manages TOTP secrets and recovery codes.
type MFARepository interface {
	GetEnrollment(ctx context.Context, userID string) (*MFAEnrollment, error)
	SaveEnrollment(ctx context.Context, enrollment *MFAEnrollment) error
	DeleteEnrollment(ctx context.Context, userID string) error
}

// PasswordHasher abstracts cryptographic password hashing and constant-time verification.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, plainPassword string) bool
}

// TokenSigner abstracts access token and refresh token generation, claims signing, and validation.
type TokenSigner interface {
	GenerateAccessToken(claims Claims) (string, error)
	GenerateRefreshToken() (plainToken string, tokenHash string, err error)
	HashRefreshToken(token string) string
	ValidateAccessToken(tokenString string) (*Claims, error)
}

// TOTPProvider abstracts RFC 6238 time-based one-time password generation and validation.
type TOTPProvider interface {
	GenerateSecret(accountEmail string, issuer string) (secret string, otpAuthURL string, err error)
	ValidatePasscode(passcode string, secret string) bool
}
