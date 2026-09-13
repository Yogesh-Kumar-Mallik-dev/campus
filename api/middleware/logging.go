/**
 * BLOCK_API_MIDDLEWARE_LOGGING_001
 * Purpose: Structured JSON request logging with latency tracking and trace ID correlation.
 * Invariants: Zero credential logging (passwords, tokens, cookies redacted).
 * Domain:  API Gateway & Transport Layer
 */

package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

type StructuredLogEntry struct {
	Timestamp string  `json:"timestamp"`
	Level     string  `json:"level"`
	TraceID   string  `json:"trace_id"`
	Method    string  `json:"method"`
	Route     string  `json:"route"`
	Status    int     `json:"status"`
	LatencyMs float64 `json:"latency_ms"`
	ClientIP  string  `json:"client_ip"`
}

// StructuredLoggingMiddleware logs every HTTP request as structured JSON.
func StructuredLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		latency := float64(time.Since(start).Nanoseconds()) / 1e6

		traceID := ""
		if tid, ok := r.Context().Value(TraceIDKey).(string); ok {
			traceID = tid
		}

		level := "INFO"
		if rec.statusCode >= 500 {
			level = "ERROR"
		} else if rec.statusCode >= 400 {
			level = "WARN"
		}

		entry := StructuredLogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Level:     level,
			TraceID:   traceID,
			Method:    r.Method,
			Route:     r.URL.Path,
			Status:    rec.statusCode,
			LatencyMs: latency,
			ClientIP:  getClientIP(r),
		}

		// Stream JSON line to stdout
		_ = json.NewEncoder(os.Stdout).Encode(entry)
	})
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
