/**
 * BLOCK_API_MIDDLEWARE_SECURITY_HEADERS_001
 * Purpose: Injects defense-in-depth HTTP security headers and strips server information banners.
 * Domain:  API Gateway & Transport Layer
 */

package middleware

import "net/http"

// securityHeaderWriter wraps http.ResponseWriter to strip sensitive leak headers before response is sent.
type securityHeaderWriter struct {
	http.ResponseWriter
}

func (w *securityHeaderWriter) WriteHeader(statusCode int) {
	w.Header().Del("Server")
	w.Header().Del("X-Powered-By")
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *securityHeaderWriter) Write(b []byte) (int, error) {
	w.Header().Del("Server")
	w.Header().Del("X-Powered-By")
	return w.ResponseWriter.Write(b)
}

// SecurityHeadersMiddleware sets hardened security headers across all HTTP responses.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()

		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		h.Set("Cross-Origin-Resource-Policy", "same-origin")
		h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Strip information leak headers
		h.Del("Server")
		h.Del("X-Powered-By")

		wrapped := &securityHeaderWriter{ResponseWriter: w}
		next.ServeHTTP(wrapped, r)
	})
}

