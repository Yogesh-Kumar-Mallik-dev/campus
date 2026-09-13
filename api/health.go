/**
 * BLOCK_API_HEALTH_PROBES_001
 * Purpose: Kubernetes / Container health probes for liveness and readiness.
 * Domain:  API Gateway & Transport Layer
 */

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// RegisterHealthProbes registers /healthz/live and /healthz/ready routes.
func RegisterHealthProbes(r chi.Router, dbPingFn func() error) {
	r.Get("/healthz/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "UP",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	r.Get("/healthz/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if dbPingFn != nil {
			if err := dbPingFn(); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"status":    "DEGRADED",
					"error":     err.Error(),
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				})
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "READY",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
}
