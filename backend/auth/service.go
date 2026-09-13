/**
 * BLOCK_AUTH_SERVICE_001
 * Purpose: Central Auth & Permissions Domain Service implementing core IAM workflows.
 * Invariants: Single-use token family rotation, brute-force lockouts, constant-time verification.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LoginCommand contains parameters to authenticate a user identity.
type LoginCommand struct {
	TenantID          string `json:"tenant_id"`
	Identifier        string `json:"identifier"` // Email or Username
	Password          string `json:"password"`
	MFACode           string `json:"mfa_code,omitempty"`
	MFATicket         string `json:"mfa_ticket,omitempty"`
	IPAddress         string `json:"ip_address,omitempty"`
	UserAgent         string `json:"user_agent,omitempty"`
	DeviceFingerprint string `json:"device_fingerprint,omitempty"`
}

// RefreshCommand contains parameters to rotate a refresh token family.
type RefreshCommand struct {
	RefreshToken string `json:"refresh_token"`
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
}

// LogoutCommand terminates an active session.
type LogoutCommand struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason,omitempty"`
}

// ChangePasswordCommand updates user credentials.
type ChangePasswordCommand struct {
	TenantID    string `json:"tenant_id"`
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ServiceConfig holds domain service tunables.
type ServiceConfig struct {
	MaxFailedLogins  int
	LockoutDuration  time.Duration
	SessionLifespan  time.Duration
	IssuerName       string
}

// DefaultServiceConfig returns production-ready security defaults.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		MaxFailedLogins: 5,
		LockoutDuration: 15 * time.Minute,
		SessionLifespan: 30 * 24 * time.Hour,
		IssuerName:      "Campus Management System",
	}
}

// AuthService coordinates all identity, session, and credential operations.
type AuthService struct {
	userRepo     UserRepository
	sessionRepo  SessionRepository
	rolePermRepo RolePermissionRepository
	mfaRepo      MFARepository
	hasher       PasswordHasher
	signer       TokenSigner
	totp         TOTPProvider
	audit        AuditPublisher
	config       ServiceConfig
}

// NewAuthService constructs a new AuthService instance with all injected dependencies.
func NewAuthService(
	userRepo UserRepository,
	sessionRepo SessionRepository,
	rolePermRepo RolePermissionRepository,
	mfaRepo MFARepository,
	hasher PasswordHasher,
	signer TokenSigner,
	totp TOTPProvider,
	audit AuditPublisher,
	config *ServiceConfig,
) *AuthService {
	cfg := DefaultServiceConfig()
	if config != nil {
		cfg = *config
	}
	if audit == nil {
		audit = &NoopAuditPublisher{}
	}
	return &AuthService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		rolePermRepo: rolePermRepo,
		mfaRepo:      mfaRepo,
		hasher:       hasher,
		signer:       signer,
		totp:         totp,
		audit:        audit,
		config:       cfg,
	}
}

// Login authenticates credentials, verifies MFA if enabled, and issues a session token pair.
func (s *AuthService) Login(ctx context.Context, cmd LoginCommand) (*AuthResult, error) {
	// Guard: Validate inputs
	if strings.TrimSpace(cmd.TenantID) == "" || strings.TrimSpace(cmd.Identifier) == "" || cmd.Password == "" {
		return nil, ErrInvalidCredentials
	}

	identifier := strings.TrimSpace(cmd.Identifier)
	normalized := strings.ToLower(identifier)

	// Step 1: Lookup user by normalized email or username
	user, err := s.userRepo.GetByEmailNormalized(ctx, cmd.TenantID, normalized)
	if err != nil || user == nil {
		user, err = s.userRepo.GetByUsername(ctx, cmd.TenantID, identifier)
		if err != nil || user == nil {
			_ = s.audit.PublishAuthEvent(ctx, &AuthAuditEvent{
				EventID:   uuid.NewString(),
				Type:      EventLoginFailed,
				TenantID:  cmd.TenantID,
				IPAddress: cmd.IPAddress,
				UserAgent: cmd.UserAgent,
				Severity:  "WARN",
				Detail:    fmt.Sprintf("Login failed: unknown identifier '%s'", identifier),
				Timestamp: time.Now().UTC(),
			})
			return nil, ErrInvalidCredentials
		}
	}

	// Step 2: Check Lockout status
	now := time.Now().UTC()
	if user.IsLockedOut(now) {
		_ = s.audit.PublishAuthEvent(ctx, &AuthAuditEvent{
			EventID:   uuid.NewString(),
			Type:      EventAccountLocked,
			TenantID:  cmd.TenantID,
			UserID:    user.ID,
			Username:  user.Username,
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			Severity:  "WARN",
			Detail:    "Login rejected: account is locked out",
			Timestamp: now,
		})
		return nil, ErrAccountLocked
	}

	// Step 3: Check Account lifecycle status
	if user.Status == StatusSuspended {
		return nil, ErrAccountSuspended
	}
	if user.Status == StatusPendingVerification {
		return nil, ErrAccountPending
	}

	// Step 4: Verify Password with constant-time comparison
	if !s.hasher.Compare(user.PasswordHash, cmd.Password) {
		_ = s.userRepo.RecordFailedLogin(ctx, cmd.TenantID, user.ID, s.config.LockoutDuration, s.config.MaxFailedLogins)
		_ = s.audit.PublishAuthEvent(ctx, &AuthAuditEvent{
			EventID:   uuid.NewString(),
			Type:      EventLoginFailed,
			TenantID:  cmd.TenantID,
			UserID:    user.ID,
			Username:  user.Username,
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			Severity:  "WARN",
			Detail:    "Login failed: invalid password",
			Timestamp: now,
		})
		return nil, ErrInvalidCredentials
	}

	// Step 5: Check MFA Enrollment
	if s.mfaRepo != nil {
		mfa, err := s.mfaRepo.GetEnrollment(ctx, user.ID)
		if err == nil && mfa != nil && mfa.IsVerified {
			if strings.TrimSpace(cmd.MFACode) == "" {
				// MFA required - emit ticket
				mfaTicket := uuid.NewString()
				return &AuthResult{
					RequiresMFA: true,
					MFATicket:   mfaTicket,
					User:        user,
				}, nil
			}

			// Validate provided TOTP passcode
			if !s.totp.ValidatePasscode(cmd.MFACode, mfa.SecretEncrypted) {
				return nil, ErrInvalidMFACode
			}
		}
	}

	// Step 6: Reset failed login attempts on successful authentication
	_ = s.userRepo.ResetFailedLogins(ctx, cmd.TenantID, user.ID)
	_ = s.userRepo.UpdateLastLogin(ctx, cmd.TenantID, user.ID, now)

	// Step 7: Resolve User Roles, Scopes, and Effective Capabilities
	roles, err := s.rolePermRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		roles = []*UserRole{}
	}
	user.Roles = roles

	perms, err := s.rolePermRepo.GetEffectivePermissions(ctx, user.ID)
	if err != nil {
		perms = []string{}
	}
	user.Permissions = perms

	// Step 8: Create Session & Refresh Token Family
	sessionID := uuid.NewString()
	familyID := uuid.NewString()

	plainRefreshToken, tokenHash, err := s.signer.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("BLOCK_AUTH_SERVICE_002: failed to generate refresh token: %w", err)
	}

	session := &Session{
		ID:                sessionID,
		UserID:            user.ID,
		TenantID:          user.TenantID,
		TokenFamilyID:     familyID,
		CurrentTokenHash:  tokenHash,
		IPAddress:         cmd.IPAddress,
		UserAgent:         cmd.UserAgent,
		DeviceFingerprint: cmd.DeviceFingerprint,
		ExpiresAt:         now.Add(s.config.SessionLifespan),
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	family := &RefreshTokenFamily{
		FamilyID:        familyID,
		UserID:          user.ID,
		ActiveTokenHash: tokenHash,
		RotationCount:   0,
		IsCompromised:   false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.sessionRepo.CreateSession(ctx, session, family); err != nil {
		return nil, fmt.Errorf("BLOCK_AUTH_SERVICE_003: failed to persist session: %w", err)
	}

	// Step 9: Issue Short-Lived Access Token
	roleKeys := make([]string, len(roles))
	for i, r := range roles {
		roleKeys[i] = r.RoleKey
	}

	claims := Claims{
		UserID:      user.ID,
		TenantID:    user.TenantID,
		Tier:        user.Tier,
		Username:    user.Username,
		Email:       user.Email,
		Roles:       roleKeys,
		Permissions: perms,
		SessionID:   sessionID,
		FamilyID:    familyID,
	}

	accessToken, err := s.signer.GenerateAccessToken(claims)
	if err != nil {
		return nil, fmt.Errorf("BLOCK_AUTH_SERVICE_004: failed to sign access token: %w", err)
	}

	// Step 10: Publish Audit Log
	_ = s.audit.PublishAuthEvent(ctx, &AuthAuditEvent{
		EventID:   uuid.NewString(),
		Type:      EventLoginSuccess,
		TenantID:  user.TenantID,
		UserID:    user.ID,
		Username:  user.Username,
		IPAddress: cmd.IPAddress,
		UserAgent: cmd.UserAgent,
		Severity:  "INFO",
		Detail:    "User authenticated successfully",
		Timestamp: now,
	})

	return &AuthResult{
		TokenPair: &TokenPair{
			AccessToken:  accessToken,
			RefreshToken: fmt.Sprintf("%s.%s", familyID, plainRefreshToken),
			TokenType:    "Bearer",
			ExpiresIn:    900, // 15 minutes
		},
		User: user,
	}, nil
}

// RefreshToken executes atomic token family rotation with breach detection.
func (s *AuthService) RefreshToken(ctx context.Context, cmd RefreshCommand) (*TokenPair, error) {
	parts := strings.Split(cmd.RefreshToken, ".")
	if len(parts) != 2 {
		return nil, ErrUnauthorized
	}

	familyID := parts[0]
	plainToken := parts[1]
	providedHash := s.signer.HashRefreshToken(plainToken)

	// Step 1: Lookup Token Family
	family, err := s.sessionRepo.GetTokenFamily(ctx, familyID)
	if err != nil || family == nil {
		return nil, ErrUnauthorized
	}

	now := time.Now().UTC()

	// Step 2: Check if family is compromised or already revoked
	if family.IsCompromised || family.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}

	// Step 3: Check Token Hash (Replay Attack Breach Detection)
	if family.ActiveTokenHash != providedHash {
		// Stolen/Old token reused! Immediately revoke entire token family and all sessions!
		_ = s.sessionRepo.RevokeTokenFamily(ctx, familyID, "TOKEN_REUSE_BREACH_DETECTED")
		_ = s.audit.PublishAuthEvent(ctx, &AuthAuditEvent{
			EventID:   uuid.NewString(),
			Type:      EventTokenBreachDetected,
			UserID:    family.UserID,
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			Severity:  "CRITICAL",
			Detail:    fmt.Sprintf("Refresh token replay attack detected on family %s; family terminated", familyID),
			Timestamp: now,
		})
		return nil, ErrTokenReused
	}

	// Step 4: Generate Fresh Token Pair
	newPlainToken, newTokenHash, err := s.signer.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("BLOCK_AUTH_SERVICE_005: failed to generate new refresh token: %w", err)
	}

	newRotationCount := family.RotationCount + 1
	if err := s.sessionRepo.RotateTokenFamily(ctx, familyID, providedHash, newTokenHash, newRotationCount); err != nil {
		return nil, fmt.Errorf("BLOCK_AUTH_SERVICE_006: failed to rotate token family: %w", err)
	}

	// Step 5: Lookup User to compile updated claims
	user, err := s.userRepo.GetByID(ctx, "", family.UserID)
	if err != nil || user == nil || !user.IsActive() {
		return nil, ErrUnauthorized
	}

	roles, _ := s.rolePermRepo.GetUserRoles(ctx, user.ID)
	perms, _ := s.rolePermRepo.GetEffectivePermissions(ctx, user.ID)

	roleKeys := make([]string, len(roles))
	for i, r := range roles {
		roleKeys[i] = r.RoleKey
	}

	claims := Claims{
		UserID:      user.ID,
		TenantID:    user.TenantID,
		Tier:        user.Tier,
		Username:    user.Username,
		Email:       user.Email,
		Roles:       roleKeys,
		Permissions: perms,
		FamilyID:    familyID,
	}

	accessToken, err := s.signer.GenerateAccessToken(claims)
	if err != nil {
		return nil, fmt.Errorf("BLOCK_AUTH_SERVICE_007: failed to sign rotated access token: %w", err)
	}

	_ = s.audit.PublishAuthEvent(ctx, &AuthAuditEvent{
		EventID:   uuid.NewString(),
		Type:      EventTokenRefreshed,
		TenantID:  user.TenantID,
		UserID:    user.ID,
		IPAddress: cmd.IPAddress,
		UserAgent: cmd.UserAgent,
		Severity:  "INFO",
		Detail:    "Token family rotated successfully",
		Timestamp: now,
	})

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: fmt.Sprintf("%s.%s", familyID, newPlainToken),
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}, nil
}

// Logout terminates an active session.
func (s *AuthService) Logout(ctx context.Context, cmd LogoutCommand) error {
	if strings.TrimSpace(cmd.SessionID) == "" {
		return ErrSessionNotFound
	}
	reason := cmd.Reason
	if reason == "" {
		reason = "USER_INITIATED_LOGOUT"
	}
	return s.sessionRepo.RevokeSession(ctx, cmd.SessionID, reason)
}

// RevokeAllSessions revokes all sessions for a user across all devices.
func (s *AuthService) RevokeAllSessions(ctx context.Context, tenantID, userID string, reason string) error {
	if reason == "" {
		reason = "GLOBAL_LOGOUT"
	}
	return s.sessionRepo.RevokeAllUserSessions(ctx, tenantID, userID, reason)
}

// ListUserSessions returns active device sessions for a user.
func (s *AuthService) ListUserSessions(ctx context.Context, tenantID, userID string) ([]*Session, error) {
	return s.sessionRepo.ListUserSessions(ctx, tenantID, userID)
}
