/**
 * BLOCK_API_AUTH_TESTS_001
 * Purpose: Complete HTTP integration and contract tests for IAM API endpoints.
 * Invariants: RFC 7807 problem details verification, 200 OK empty array semantics.
 * Domain:  API Gateway & Transport Layer
 */

package authhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"campus/backend/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func setupTestHTTPRouter() (chi.Router, *auth.AuthService, *auth.JWTSigner, *auth.MockUserRepo) {
	userRepo := auth.NewMockUserRepo()
	sessionRepo := auth.NewMockSessionRepo()
	rolePermRepo := auth.NewMockRolePermRepo()
	mfaRepo := auth.NewMockMFARepo()
	hasher := auth.NewArgon2idHasher(nil)
	signer := auth.NewJWTSigner("test-secret-key-32-bytes-long!", "campus.test", 15*time.Minute)
	totp := auth.NewStandardTOTPProvider()
	audit := &auth.NoopAuditPublisher{}

	cfg := &auth.ServiceConfig{
		MaxFailedLogins: 5,
		LockoutDuration: 15 * time.Minute,
		SessionLifespan: 24 * time.Hour,
		IssuerName:      "Test Campus",
	}

	service := auth.NewAuthService(
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

	handler := NewAuthHandler(service)
	authMiddleware := NewAuthMiddleware(signer)

	r := chi.NewRouter()
	RegisterRoutes(r, handler, authMiddleware)

	return r, service, signer, userRepo
}

func TestHTTP_Login_Success(t *testing.T) {
	router, _, _, userRepo := setupTestHTTPRouter()
	ctx := context.Background()

	hasher := auth.NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("ValidPassword123")

	user := &auth.User{
		ID:              uuid.NewString(),
		TenantID:        "tenant_01",
		Tier:            auth.TierUser,
		Username:        "student.alpha",
		Email:           "student@campus.edu",
		EmailNormalized: "student@campus.edu",
		PasswordHash:    pwdHash,
		Status:          auth.StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	body, _ := json.Marshal(map[string]string{
		"tenant_id":  "tenant_01",
		"identifier": "student@campus.edu",
		"password":   "ValidPassword123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got: %d (%s)", rec.Code, rec.Body.String())
	}

	var res auth.AuthResult
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.TokenPair == nil || res.TokenPair.AccessToken == "" {
		t.Fatal("Expected valid TokenPair in response")
	}
}

func TestHTTP_Login_InvalidCredentials_RFC7807(t *testing.T) {
	router, _, _, _ := setupTestHTTPRouter()

	body, _ := json.Marshal(map[string]string{
		"tenant_id":  "tenant_01",
		"identifier": "unknown@campus.edu",
		"password":   "WrongPass",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got: %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/problem+json" {
		t.Errorf("Expected Content-Type application/problem+json, got: %s", contentType)
	}
}

func TestHTTP_Me_And_Sessions_Authenticated(t *testing.T) {
	router, _, signer, _ := setupTestHTTPRouter()

	claims := auth.Claims{
		UserID:      "usr_test_123",
		TenantID:    "tenant_01",
		Tier:        auth.TierUser,
		Username:    "prof.test",
		Email:       "prof@campus.edu",
		Roles:       []string{"professor"},
		Permissions: []string{"academics:attendance:mark"},
		SessionID:   "session_abc_123",
		FamilyID:    "family_xyz_123",
	}

	token, _ := signer.GenerateAccessToken(claims)

	// Test GET /api/v1/auth/me
	reqMe := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+token)
	recMe := httptest.NewRecorder()

	router.ServeHTTP(recMe, reqMe)

	if recMe.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK on /me, got: %d (%s)", recMe.Code, recMe.Body.String())
	}

	var meRes map[string]any
	_ = json.NewDecoder(recMe.Body).Decode(&meRes)
	if meRes["user_id"] != "usr_test_123" {
		t.Errorf("Expected user_id usr_test_123, got: %v", meRes["user_id"])
	}

	// Test GET /api/v1/auth/sessions (Empty collection query semantics)
	reqSessions := httptest.NewRequest(http.MethodGet, "/api/v1/auth/sessions", nil)
	reqSessions.Header.Set("Authorization", "Bearer "+token)
	recSessions := httptest.NewRecorder()

	router.ServeHTTP(recSessions, reqSessions)

	if recSessions.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK on /sessions, got: %d", recSessions.Code)
	}

	var sessRes map[string]any
	_ = json.NewDecoder(recSessions.Body).Decode(&sessRes)
	data, ok := sessRes["data"].([]any)
	if !ok || len(data) != 0 {
		t.Errorf("Expected empty array for data collection, got: %v", sessRes["data"])
	}
}

func TestHTTP_Logout(t *testing.T) {
	router, service, signer, userRepo := setupTestHTTPRouter()
	ctx := context.Background()

	hasher := auth.NewArgon2idHasher(nil)
	pwdHash, _ := hasher.Hash("LogoutPass123")

	userID := uuid.NewString()
	user := &auth.User{
		ID:              userID,
		TenantID:        "tenant_01",
		Tier:            auth.TierUser,
		Username:        "logout.user",
		EmailNormalized: "logout@campus.edu",
		PasswordHash:    pwdHash,
		Status:          auth.StatusActive,
	}
	_ = userRepo.Create(ctx, user)

	loginRes, _ := service.Login(ctx, auth.LoginCommand{
		TenantID:   "tenant_01",
		Identifier: "logout@campus.edu",
		Password:   "LogoutPass123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+loginRes.TokenPair.AccessToken)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("Expected 204 No Content, got: %d (%s)", rec.Code, rec.Body.String())
	}
	_ = signer
}
