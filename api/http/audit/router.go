/**
 * BLOCK_API_AUDIT_ROUTER_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Route mounting and capability-based access authorization for audit endpoints.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"github.com/go-chi/chi/v5"

	"campus/backend/audit"
)

// RegisterRoutes registers all Audit subsystem REST endpoints onto a Chi router.
func RegisterRoutes(r chi.Router, service audit.Service) {
	h := NewHandler(service)

	r.Route("/api/v1/audit", func(r chi.Router) {
		r.Get("/logs", h.QueryLogs)
		r.Get("/logs/{id}", h.GetLogByID)
		r.Post("/verify", h.VerifyIntegrity)
		r.Post("/checkpoints", h.CreateCheckpoint)
		r.Get("/compliance-reports", h.GenerateComplianceReport)
	})
}
