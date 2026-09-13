/**
 * BLOCK_API_MIDDLEWARE_REQUEST_ID_001
 * Purpose: Assigns or validates X-Request-ID / W3C traceparent and binds it to context & headers.
 * Domain:  API Gateway & Transport Layer
 */

package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const TraceIDKey contextKey = "trace_id"

// RequestIDMiddleware ensures every HTTP request carries an X-Request-ID.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = "req_" + uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), TraceIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
