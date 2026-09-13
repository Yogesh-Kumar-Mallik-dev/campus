/**
 * BLOCK_API_AUTH_MIDDLEWARE_001
 * Purpose: Authentication & capability authorization middlewares.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package authhttp

import (
	"context"
	"net/http"
	"strings"

	"campus/api/problem"
	"campus/backend/auth"
)

type contextKey string

const claimsContextKey contextKey = "auth_claims"

// WithClaims injects claims into a context.
func WithClaims(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// ClaimsFromContext extracts verified claims from context.
func ClaimsFromContext(ctx context.Context) *auth.Claims {
	if val, ok := ctx.Value(claimsContextKey).(*auth.Claims); ok {
		return val
	}
	return nil
}

// NewAuthMiddleware validates Bearer JWT access tokens and injects verified claims.
func NewAuthMiddleware(signer auth.TokenSigner) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				problem.Unauthorized(w, r, "missing or malformed Authorization header", "AUTH_MISSING_TOKEN")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := signer.ValidateAccessToken(tokenString)
			if err != nil {
				if err == auth.ErrTokenExpired {
					problem.Unauthorized(w, r, "access token is expired; please refresh", "AUTH_TOKEN_EXPIRED")
					return
				}
				problem.Unauthorized(w, r, "invalid or unverified access token", "AUTH_INVALID_TOKEN")
				return
			}

			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireCapability enforces that the authenticated identity possesses a specific permission.
func RequireCapability(requiredPerm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				problem.Unauthorized(w, r, "authentication required", "AUTH_UNAUTHENTICATED")
				return
			}

			if !claims.HasPermission(requiredPerm) {
				problem.Forbidden(w, r, "insufficient capabilities to access this resource", "AUTH_FORBIDDEN")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole enforces that the authenticated identity possesses a specific persona role.
func RequireRole(roleKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				problem.Unauthorized(w, r, "authentication required", "AUTH_UNAUTHENTICATED")
				return
			}

			if !claims.HasRole(roleKey) {
				problem.Forbidden(w, r, "required persona role missing", "AUTH_FORBIDDEN")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
