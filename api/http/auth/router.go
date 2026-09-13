/**
 * BLOCK_API_AUTH_ROUTER_001
 * Purpose: Chi HTTP route registration for Central Auth & IAM.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package authhttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes sets up public and authenticated IAM endpoints.
func RegisterRoutes(r chi.Router, handler *AuthHandler, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		// Public Authentication Endpoints
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.Refresh)

		// Authenticated Endpoints (Guarded by JWT Auth Middleware)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Post("/logout", handler.Logout)
			r.Get("/me", handler.Me)
			r.Get("/sessions", handler.ListSessions)
			r.Delete("/sessions/{id}", handler.RevokeSession)
			r.Delete("/sessions", handler.RevokeAllSessions)
		})
	})
}
