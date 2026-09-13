/**
 * BLOCK_API_MIDDLEWARE_CONTENT_TYPE_001
 * Purpose: Enforces Content-Type negotiation on state-mutating HTTP requests.
 * Domain:  API Gateway & Transport Layer
 */

package middleware

import (
	"net/http"
	"strings"

	"campus/api/problem"
)

// EnforceJSONContentType rejects POST, PUT, and PATCH requests lacking application/json.
func EnforceJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			// Skip if body is empty (e.g. standard logout)
			if r.ContentLength > 0 {
				ct := r.Header.Get("Content-Type")
				if !strings.HasPrefix(ct, "application/json") {
					problem.UnsupportedMediaType(w, r, "requests with payload must declare Content-Type: application/json", "UNSUPPORTED_MEDIA_TYPE")
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
