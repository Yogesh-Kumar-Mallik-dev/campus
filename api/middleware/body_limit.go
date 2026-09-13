/**
 * BLOCK_API_MIDDLEWARE_BODY_LIMIT_001
 * Purpose: Enforces payload size caps to prevent memory exhaustion and DoS attacks.
 * Domain:  API Gateway & Transport Layer
 */

package middleware

import (
	"net/http"

	"campus/api/problem"
)

// BodyLimitMiddleware caps incoming request body bytes (e.g. 250 KB for JSON).
func BodyLimitMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	if maxBytes <= 0 {
		maxBytes = 250 * 1024 // 250 KB default
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				problem.PayloadTooLarge(w, r, "request body exceeds maximum allowed size of 250 KB", "PAYLOAD_TOO_LARGE")
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
