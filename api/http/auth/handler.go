/**
 * BLOCK_API_AUTH_HANDLERS_001
 * Purpose: HTTP Handlers for Central Auth & IAM endpoints with RFC 7807 problem details.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package authhttp

import (
	"encoding/json"
	"net/http"
	"strings"

	"campus/api/problem"
	"campus/backend/auth"
	"github.com/go-chi/chi/v5"
)

// AuthHandler exposes HTTP handlers for IAM operations.
type AuthHandler struct {
	service *auth.AuthService
}

// NewAuthHandler constructs a new AuthHandler instance.
func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var cmd auth.LoginCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed JSON request payload", "MALFORMED_JSON")
		return
	}

	cmd.IPAddress = getClientIP(r)
	cmd.UserAgent = r.UserAgent()

	result, err := h.service.Login(r.Context(), cmd)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			problem.Unauthorized(w, r, "invalid identifier or password", "INVALID_CREDENTIALS")
		case auth.ErrAccountLocked:
			problem.Forbidden(w, r, "account is locked out due to multiple failed login attempts; try again in 15 minutes", "ACCOUNT_LOCKED")
		case auth.ErrAccountSuspended:
			problem.Forbidden(w, r, "account has been suspended", "ACCOUNT_SUSPENDED")
		case auth.ErrAccountPending:
			problem.Forbidden(w, r, "account is pending verification", "ACCOUNT_PENDING_VERIFICATION")
		case auth.ErrInvalidMFACode:
			problem.UnprocessableEntity(w, r, "invalid or expired TOTP passcode", "INVALID_MFA_CODE", []*problem.InvalidParam{
				{Name: "mfa_code", Reason: "Passcode does not match active authenticator token"},
			})
		default:
			problem.InternalServerError(w, r, "failed to authenticate session")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var cmd auth.RefreshCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed JSON request payload", "MALFORMED_JSON")
		return
	}

	cmd.IPAddress = getClientIP(r)
	cmd.UserAgent = r.UserAgent()

	tokenPair, err := h.service.RefreshToken(r.Context(), cmd)
	if err != nil {
		switch err {
		case auth.ErrTokenReused:
			problem.Forbidden(w, r, "refresh token reuse breach detected; session terminated immediately", "TOKEN_REUSE_BREACH_DETECTED")
		case auth.ErrSessionRevoked:
			problem.Unauthorized(w, r, "session has been revoked", "SESSION_REVOKED")
		default:
			problem.Unauthorized(w, r, "invalid, expired, or compromised refresh token", "INVALID_REFRESH_TOKEN")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tokenPair)
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil || claims.SessionID == "" {
		problem.Unauthorized(w, r, "unauthenticated session", "AUTH_UNAUTHENTICATED")
		return
	}

	err := h.service.Logout(r.Context(), auth.LogoutCommand{
		SessionID: claims.SessionID,
		Reason:    "USER_LOGOUT",
	})
	if err != nil {
		problem.InternalServerError(w, r, "failed to terminate session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Me handles GET /api/v1/auth/me.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		problem.Unauthorized(w, r, "unauthenticated session", "AUTH_UNAUTHENTICATED")
		return
	}

	// Compute accessible Hub workspaces based on roles
	workspaces := make([]string, 0)
	if claims.HasRole("student") {
		workspaces = append(workspaces, "student_hub")
	}
	if claims.HasRole("faculty") || claims.HasRole("professor") || claims.HasRole("teacher") {
		workspaces = append(workspaces, "teacher_hub")
	}
	if claims.Tier == auth.TierAdmin || claims.Tier == auth.TierSuperAdmin || claims.HasRole("admin") {
		workspaces = append(workspaces, "admin_hub")
	}
	if claims.HasRole("librarian") {
		workspaces = append(workspaces, "librarian_hub")
	}
	if claims.HasRole("warden") {
		workspaces = append(workspaces, "warden_hub")
	}
	if claims.HasRole("parent") {
		workspaces = append(workspaces, "parent_hub")
	}

	response := map[string]any{
		"user_id":            claims.UserID,
		"tenant_id":          claims.TenantID,
		"tier":               claims.Tier,
		"username":           claims.Username,
		"email":              claims.Email,
		"roles":              claims.Roles,
		"permissions":        claims.Permissions,
		"allowed_workspaces": workspaces,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// ListSessions handles GET /api/v1/auth/sessions.
func (h *AuthHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		problem.Unauthorized(w, r, "unauthenticated session", "AUTH_UNAUTHENTICATED")
		return
	}

	sessions, err := h.service.ListUserSessions(r.Context(), claims.TenantID, claims.UserID)
	if err != nil {
		problem.InternalServerError(w, r, "failed to query active sessions")
		return
	}

	// Standard 200 OK empty array semantics for collection queries
	if sessions == nil {
		sessions = make([]*auth.Session, 0)
	}

	response := map[string]any{
		"data": sessions,
		"pagination": map[string]any{
			"page":       1,
			"pageSize":   10,
			"totalItems": len(sessions),
			"totalPages": 1,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// RevokeSession handles DELETE /api/v1/auth/sessions/{id}.
func (h *AuthHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if strings.TrimSpace(sessionID) == "" {
		problem.BadRequest(w, r, "missing session ID parameter", "MISSING_PARAM")
		return
	}

	err := h.service.Logout(r.Context(), auth.LogoutCommand{
		SessionID: sessionID,
		Reason:    "USER_REVOKED_DEVICE",
	})
	if err != nil {
		problem.InternalServerError(w, r, "failed to revoke session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RevokeAllSessions handles DELETE /api/v1/auth/sessions.
func (h *AuthHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		problem.Unauthorized(w, r, "unauthenticated session", "AUTH_UNAUTHENTICATED")
		return
	}

	err := h.service.RevokeAllSessions(r.Context(), claims.TenantID, claims.UserID, "USER_GLOBAL_LOGOUT")
	if err != nil {
		problem.InternalServerError(w, r, "failed to revoke all sessions")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getClientIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		parts := strings.Split(xf, ",")
		return strings.TrimSpace(parts[0])
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return strings.TrimSpace(xr)
	}
	return r.RemoteAddr
}
