/**
 * BLOCK_AUTH_SERVICE_TESTS_001
 * Purpose: Complete AAA unit test suite for Central Auth & IAM.
 * Invariants: 100% domain coverage, zero DB dependency, millisecond execution.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func setupTestAuthService(t *testing.T) (*AuthService, *MockUserRepo, *MockSessionRepo, *MockRolePermRepo, *MockMFARepo) {
	userRepo := NewMockUserRepo()
	sessionRepo := NewMockSessionRepo()
	rolePermRepo := NewMockRolePermRepo()
	mfaRepo := NewMockMFARepo()
	hasher := NewArgon2idHasher(nil)
	signer := NewJWTSigner("test-secret-key-32-bytes-long!", "campus.test", 15*time.Minute)
	totp := NewStandardTOTPProvider()
	audit := &NoopAuditPublisher{}

	cfg := &ServiceConfig{
		MaxFailedLogins: 5,
		LockoutDuration: 15 * time.Minute,
		SessionLifespan: 24 * time.Hour,
		IssuerName:      "Test Campus",
	}

	service := NewAuthService(
		userRepo,
		sessionRepo,
		rolePermRepo,
		mfaRepo,
		hasher,
		signer,
		totp,
		audit,
		cfg,
	)

	return service, userRepo, sessionRepo, rolePermRepo, mfaRepo
}

func TestLogin_Success_Email(t *testing.T) {
	service, userRepo, _, rolePermRepo, _ := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, err := hasher.Hash("SecureP@ssw0rd123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	userID := uuid.NewString()
	tenantID := "tenant_01"

	user := &User{
		ID:              userID,
		TenantID:        tenantID,
		Tier:            TierUser,
		Username:        "rajesh.sharma",
		Email:           "rajesh.sharma@campus.edu",
		EmailNormalized: "rajesh.sharma@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusActive,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	_ = userRepo.Create(ctx, user)

	_ = rolePermRepo.AssignRoleToUser(ctx, &UserRole{
		ID:        uuid.NewString(),
		UserID:    userID,
		RoleID:    "role_prof",
		RoleKey:   "professor",
		RoleName:  "Professor",
		ScopeType: "DEPT",
		ScopeID:   "mech_dept",
	})

	rolePermRepo.Permissions[userID] = []string{"academics:attendance:mark", "academics:grades:submit"}

	// Act
	res, err := service.Login(ctx, LoginCommand{
		TenantID:   tenantID,
		Identifier: "rajesh.sharma@campus.edu",
		Password:   "SecureP@ssw0rd123",
	})

	// Assert
	if err != nil {
		t.Fatalf("Expected login to succeed, got error: %v", err)
	}
	if res.TokenPair == nil {
		t.Fatal("Expected TokenPair, got nil")
	}
	if res.TokenPair.AccessToken == "" || res.TokenPair.RefreshToken == "" {
		t.Fatal("Expected valid tokens in TokenPair")
	}
	if len(res.User.Roles) != 1 || res.User.Roles[0].RoleKey != "professor" {
		t.Errorf("Expected role 'professor', got %v", res.User.Roles)
	}
	if len(res.User.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(res.User.Permissions))
	}
}

func TestLogin_Success_Username(t *testing.T) {
	service, userRepo, _, _, _ := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("StudentP@ss1")

	user := &User{
		ID:              uuid.NewString(),
		TenantID:        "tenant_01",
		Tier:            TierUser,
		Username:        "anand.mishra",
		Email:           "anand@campus.edu",
		EmailNormalized: "anand@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	// Act: Login with username instead of email
	res, err := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "anand.mishra",
		Password:   "StudentP@ss1",
	})

	// Assert
	if err != nil {
		t.Fatalf("Expected login to succeed with username, got: %v", err)
	}
	if res.User.Username != "anand.mishra" {
		t.Errorf("Expected username 'anand.mishra', got '%s'", res.User.Username)
	}
}

func TestLogin_InvalidPassword_IncrementsFailedAttempts(t *testing.T) {
	service, userRepo, _, _, _ := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("CorrectPassword123")

	userID := uuid.NewString()
	user := &User{
		ID:              userID,
		TenantID:        "tenant_01",
		Tier:            TierUser,
		Username:        "target.user",
		Email:           "target@campus.edu",
		EmailNormalized: "target@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	// Act: Attempt login with wrong password
	_, err := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "target@campus.edu",
		Password:   "WrongPassword!",
	})

	// Assert
	if err != ErrInvalidCredentials {
		t.Fatalf("Expected ErrInvalidCredentials, got: %v", err)
	}
	savedUser, _ := userRepo.GetByID(ctx, "tenant_01", userID)
	if savedUser.FailedLoginAttempts != 1 {
		t.Errorf("Expected 1 failed login attempt, got %d", savedUser.FailedLoginAttempts)
	}
}

func TestLogin_AccountLockout_After5FailedAttempts(t *testing.T) {
	service, userRepo, _, _, _ := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("CorrectPassword123")

	userID := uuid.NewString()
	user := &User{
		ID:              userID,
		TenantID:        "tenant_01",
		Tier:            TierUser,
		Username:        "brute.target",
		Email:           "brute@campus.edu",
		EmailNormalized: "brute@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	// Act: Fail 5 times
	for i := 0; i < 5; i++ {
		_, _ = service.Login(ctx, LoginCommand{
			TenantID:   "tenant_01",
			Identifier: "brute@campus.edu",
			Password:   "WrongPassword",
		})
	}

	// 6th Attempt should trigger ErrAccountLocked
	_, err := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "brute@campus.edu",
		Password:   "CorrectPassword123",
	})

	// Assert
	if err != ErrAccountLocked {
		t.Fatalf("Expected ErrAccountLocked on 6th attempt, got: %v", err)
	}
}

func TestLogin_SuspendedAccount_Rejection(t *testing.T) {
	service, userRepo, _, _, _ := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("Password123")

	user := &User{
		ID:              uuid.NewString(),
		TenantID:        "tenant_01",
		Tier:            TierUser,
		Username:        "suspended.user",
		Email:           "suspended@campus.edu",
		EmailNormalized: "suspended@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusSuspended,
	}
	_ = userRepo.Create(ctx, user)

	// Act
	_, err := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "suspended@campus.edu",
		Password:   "Password123",
	})

	// Assert
	if err != ErrAccountSuspended {
		t.Fatalf("Expected ErrAccountSuspended, got: %v", err)
	}
}

func TestLogin_MFA_Flow(t *testing.T) {
	service, userRepo, _, _, mfaRepo := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("MfaUserPass123")

	totpProvider := NewStandardTOTPProvider()
	secret, _, _ := totpProvider.GenerateSecret("mfa@campus.edu", "TestCampus")

	userID := uuid.NewString()
	user := &User{
		ID:              userID,
		TenantID:        "tenant_01",
		Tier:            TierUser,
		Username:        "mfa.user",
		Email:           "mfa@campus.edu",
		EmailNormalized: "mfa@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	_ = mfaRepo.SaveEnrollment(ctx, &MFAEnrollment{
		UserID:          userID,
		MFAType:         MFATypeTOTP,
		SecretEncrypted: secret,
		IsVerified:      true,
		EnrolledAt:      time.Now().UTC(),
	})

	// Step 1: Login without MFA code -> should return RequiresMFA = true
	res1, err := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "mfa@campus.edu",
		Password:   "MfaUserPass123",
	})

	if err != nil {
		t.Fatalf("Expected stage 1 login to succeed, got: %v", err)
	}
	if !res1.RequiresMFA {
		t.Fatal("Expected RequiresMFA to be true")
	}

	// Step 2: Login with Invalid MFA code -> should fail
	_, err = service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "mfa@campus.edu",
		Password:   "MfaUserPass123",
		MFACode:    "000000",
	})
	if err != ErrInvalidMFACode {
		t.Fatalf("Expected ErrInvalidMFACode, got: %v", err)
	}
}

func TestRefreshToken_NominalRotation(t *testing.T) {
	service, userRepo, _, _, _ := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("RotationPass123")

	userID := uuid.NewString()
	user := &User{
		ID:              userID,
		TenantID:        "tenant_01",
		Tier:            TierUser,
		Username:        "rotate.user",
		Email:           "rotate@campus.edu",
		EmailNormalized: "rotate@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	// Step 1: Initial Login
	loginRes, err := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "rotate@campus.edu",
		Password:   "RotationPass123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	firstRefreshToken := loginRes.TokenPair.RefreshToken

	// Step 2: Rotate token
	rotatedRes, err := service.RefreshToken(ctx, RefreshCommand{
		RefreshToken: firstRefreshToken,
	})
	if err != nil {
		t.Fatalf("Expected token refresh to succeed, got: %v", err)
	}

	if rotatedRes.AccessToken == "" || rotatedRes.RefreshToken == "" {
		t.Fatal("Expected rotated tokens, got empty")
	}
	if rotatedRes.RefreshToken == firstRefreshToken {
		t.Fatal("Expected new refresh token string after rotation")
	}
}

func TestRefreshToken_ReplayAttack_RevokesFamily(t *testing.T) {
	service, userRepo, sessionRepo, _, _ := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("BreachPass123")

	userID := uuid.NewString()
	user := &User{
		ID:              userID,
		TenantID:        "tenant_01",
		Tier:            TierUser,
		Username:        "breach.target",
		Email:           "breach@campus.edu",
		EmailNormalized: "breach@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	// Step 1: Initial Login (Token A generated)
	loginRes, err := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "breach@campus.edu",
		Password:   "BreachPass123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	tokenA := loginRes.TokenPair.RefreshToken

	// Step 2: Legitimate rotation (Token A consumed, Token B issued)
	_, err = service.RefreshToken(ctx, RefreshCommand{
		RefreshToken: tokenA,
	})
	if err != nil {
		t.Fatalf("First rotation failed: %v", err)
	}

	// Step 3: Malicious Replay Attack (Attacker tries to reuse Token A)
	_, err = service.RefreshToken(ctx, RefreshCommand{
		RefreshToken: tokenA,
	})

	// Assert: Replay must return ErrTokenReused
	if err != ErrTokenReused {
		t.Fatalf("Expected ErrTokenReused on replay attack, got: %v", err)
	}

	// Step 4: Verify Token Family is compromised and revoked in database
	familyID := strings.Split(tokenA, ".")[0]
	family, err := sessionRepo.GetTokenFamily(ctx, familyID)
	if err != nil {
		t.Fatalf("Failed to lookup token family: %v", err)
	}
	if !family.IsCompromised || family.RevokedAt == nil {
		t.Error("Expected token family to be marked compromised and revoked")
	}
}

func TestLogout_And_RevokeAllSessions(t *testing.T) {
	service, userRepo, sessionRepo, _, _ := setupTestAuthService(t)
	ctx := context.Background()

	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("SessionPass123")

	userID := uuid.NewString()
	user := &User{
		ID:              userID,
		TenantID:        "tenant_01",
		Tier:            TierUser,
		Username:        "session.user",
		Email:           "session@campus.edu",
		EmailNormalized: "session@campus.edu",
		PasswordHash:    pwdHash,
		Status:          StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	// Login 2 devices
	login1, _ := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "session@campus.edu",
		Password:   "SessionPass123",
	})
	login2, _ := service.Login(ctx, LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "session@campus.edu",
		Password:   "SessionPass123",
	})

	sessions, _ := service.ListUserSessions(ctx, "tenant_01", userID)
	if len(sessions) != 2 {
		t.Fatalf("Expected 2 active sessions, got %d", len(sessions))
	}

	// Logout first session
	signer := NewJWTSigner("test-secret-key-32-bytes-long!", "campus.test", 15*time.Minute)
	claims1, _ := signer.ValidateAccessToken(login1.TokenPair.AccessToken)

	err := service.Logout(ctx, LogoutCommand{
		SessionID: claims1.SessionID,
	})
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	sessions, _ = service.ListUserSessions(ctx, "tenant_01", userID)
	if len(sessions) != 1 {
		t.Fatalf("Expected 1 active session after logout, got %d", len(sessions))
	}

	// Revoke all remaining sessions
	_ = service.RevokeAllSessions(ctx, "tenant_01", userID, "SECURITY_RESET")
	sessions, _ = service.ListUserSessions(ctx, "tenant_01", userID)
	if len(sessions) != 0 {
		t.Fatalf("Expected 0 active sessions after RevokeAllSessions, got %d", len(sessions))
	}
	_ = login2
	_ = sessionRepo
}

func TestClaims_HasPermission_And_SuperAdmin(t *testing.T) {
	// Standard User Claims
	userClaims := Claims{
		UserID:      "u1",
		Tier:        TierUser,
		Roles:       []string{"professor"},
		Permissions: []string{"academics:grades:submit", "hostel:gatepass:*"},
	}

	if !userClaims.HasPermission("academics:grades:submit") {
		t.Error("Expected true for exact permission match")
	}
	if !userClaims.HasPermission("hostel:gatepass:approve") {
		t.Error("Expected true for wildcard permission prefix match")
	}
	if userClaims.HasPermission("billing:invoices:create") {
		t.Error("Expected false for ungranted permission")
	}

	// Super Admin Claims (bypasses all checks)
	superClaims := Claims{
		UserID: "admin_root",
		Tier:   TierSuperAdmin,
	}
	if !superClaims.HasPermission("anything:anywhere:anytime") {
		t.Error("SuperAdmin must have all permissions implicitly")
	}
	if !superClaims.HasRole("any_role") {
		t.Error("SuperAdmin must satisfy any role check")
	}
}

func TestJWTSigner_Validation_Failures(t *testing.T) {
	signer := NewJWTSigner("test-secret-key-32-bytes-long!", "campus.test", 15*time.Minute)

	// Invalid format
	if _, err := signer.ValidateAccessToken("invalid.token"); err != ErrUnauthorized {
		t.Errorf("Expected ErrUnauthorized on invalid token format, got: %v", err)
	}

	// Tampered signature
	claims := Claims{UserID: "u1", TenantID: "t1"}
	token, _ := signer.GenerateAccessToken(claims)
	tampered := token + "tampered"
	if _, err := signer.ValidateAccessToken(tampered); err != ErrUnauthorized {
		t.Errorf("Expected ErrUnauthorized on tampered token, got: %v", err)
	}

	// Expired token
	expiredClaims := Claims{
		UserID:    "u1",
		TenantID:  "t1",
		IssuedAt:  time.Now().Unix() - 3600,
		ExpiresAt: time.Now().Unix() - 60,
	}
	expiredToken, _ := signer.GenerateAccessToken(expiredClaims)
	if _, err := signer.ValidateAccessToken(expiredToken); err != ErrTokenExpired {
		t.Errorf("Expected ErrTokenExpired on expired token, got: %v", err)
	}
}

func TestArgon2idHasher_EdgeCases(t *testing.T) {
	hasher := NewArgon2idHasher(nil)
	hash, err := hasher.Hash("ValidPass123")
	if err != nil {
		t.Fatalf("Failed to hash: %v", err)
	}

	if !hasher.Compare(hash, "ValidPass123") {
		t.Error("Expected password to match")
	}
	if hasher.Compare(hash, "WrongPass123") {
		t.Error("Expected wrong password not to match")
	}
	if hasher.Compare("invalid$hash$format", "ValidPass123") {
		t.Error("Expected invalid hash format to return false")
	}
}

func TestTOTPProvider_RFC6238_Validation(t *testing.T) {
	totp := NewStandardTOTPProvider()
	secret, urlStr, err := totp.GenerateSecret("user@campus.edu", "TestCampus")
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}
	if secret == "" || !strings.Contains(urlStr, "otpauth://") {
		t.Fatal("Expected valid secret and otpauth URL")
	}

	if totp.ValidatePasscode("123", secret) {
		t.Error("Expected short passcode to fail")
	}
	if totp.ValidatePasscode("123456", "invalid-base32-secret-!!!") {
		t.Error("Expected invalid base32 secret to fail")
	}
}

func TestDomainError_Format(t *testing.T) {
	err := NewDomainError("BLOCK_AUTH_001", "TEST_CODE", 400, "Test detail error", nil)
	if !strings.Contains(err.Error(), "[BLOCK_AUTH_001]") || !strings.Contains(err.Error(), "400") {
		t.Errorf("Expected formatted domain error string, got: %s", err.Error())
	}
}

func TestLogin_EdgeCases_PendingAndInvalidInputs(t *testing.T) {
	service, userRepo, _, _, _ := setupTestAuthService(t)
	ctx := context.Background()

	// Empty inputs
	if _, err := service.Login(ctx, LoginCommand{}); err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials on empty login command, got: %v", err)
	}

	// Pending verification user
	pendingUser := &User{
		ID:              uuid.NewString(),
		TenantID:        "t1",
		Username:        "pending",
		EmailNormalized: "pending@campus.edu",
		Status:          StatusPendingVerification,
		PasswordHash:    "$argon2id$v=19$m=65536,t=3,p=4$dummy$dummy",
	}
	_ = userRepo.Create(ctx, pendingUser)
	hasher := NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("Pass123")
	pendingUser.PasswordHash = pwdHash

	if _, err := service.Login(ctx, LoginCommand{
		TenantID:   "t1",
		Identifier: "pending@campus.edu",
		Password:   "Pass123",
	}); err != ErrAccountPending {
		t.Errorf("Expected ErrAccountPending, got: %v", err)
	}
}

