/**
 * BLOCK_API_HUB_HANDLER_001
 * Subsystem: Rank 17 - The Hub Root Super-App (hub)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for persona cockpits, live metrics aggregation, widget customizer, and 1-tap shortcuts.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package hub

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/hub"
)

type Handler struct {
	service hub.Service
}

func NewHandler(service hub.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/hub", func(r chi.Router) {
		r.Get("/cockpit", h.GetPersonaCockpit)
		r.Post("/widgets/reconfigure", h.ReconfigureWidgets)
		r.Post("/theme", h.UpdateLayoutTheme)
		r.Get("/shortcuts", h.ListShortcuts)
		r.Post("/shortcuts/trigger", h.TriggerShortcut)
	})
}

// GetPersonaCockpit handles GET /api/v1/hub/cockpit?persona=STUDENT
func (h *Handler) GetPersonaCockpit(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id is required", "MISSING_TENANT_ID")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = r.Header.Get("X-User-ID")
	}
	if userID == "" {
		userID = "anonymous-user"
	}

	personaParam := r.URL.Query().Get("persona")
	if personaParam == "" {
		personaParam = "STUDENT"
	}

	cockpit, err := h.service.GetPersonaCockpit(r.Context(), tenantID, userID, hub.PersonaType(personaParam))
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cockpit)
}

// ReconfigureWidgets handles POST /api/v1/hub/widgets/reconfigure
func (h *Handler) ReconfigureWidgets(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID    string             `json:"tenant_id"`
		UserID      string             `json:"user_id"`
		DashboardID string             `json:"dashboard_id"`
		Widgets     []hub.WidgetConfig `json:"widgets"`
		ActorID     string             `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}
	if req.UserID == "" {
		req.UserID = r.Header.Get("X-User-ID")
	}
	if req.ActorID == "" {
		req.ActorID = req.UserID
	}

	if req.TenantID == "" || req.DashboardID == "" {
		problem.BadRequest(w, r, "tenant_id and dashboard_id are required", "VALIDATION_ERROR")
		return
	}

	dash, err := h.service.ReconfigureWidgets(r.Context(), hub.ReconfigureWidgetsRequest{
		TenantID:    req.TenantID,
		UserID:      req.UserID,
		DashboardID: req.DashboardID,
		Widgets:     req.Widgets,
		ActorID:     req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dash)
}

// UpdateLayoutTheme handles POST /api/v1/hub/theme
func (h *Handler) UpdateLayoutTheme(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID    string `json:"tenant_id"`
		UserID      string `json:"user_id"`
		DashboardID string `json:"dashboard_id"`
		LayoutTheme string `json:"layout_theme"`
		ActorID     string `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}
	if req.UserID == "" {
		req.UserID = r.Header.Get("X-User-ID")
	}
	if req.ActorID == "" {
		req.ActorID = req.UserID
	}

	if req.TenantID == "" || req.DashboardID == "" || req.LayoutTheme == "" {
		problem.BadRequest(w, r, "tenant_id, dashboard_id, and layout_theme are required", "VALIDATION_ERROR")
		return
	}

	dash, err := h.service.UpdateLayoutTheme(r.Context(), hub.UpdateLayoutThemeRequest{
		TenantID:    req.TenantID,
		UserID:      req.UserID,
		DashboardID: req.DashboardID,
		LayoutTheme: req.LayoutTheme,
		ActorID:     req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dash)
}

// ListShortcuts handles GET /api/v1/hub/shortcuts?persona=STUDENT
func (h *Handler) ListShortcuts(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id is required", "MISSING_TENANT_ID")
		return
	}

	personaParam := r.URL.Query().Get("persona")
	if personaParam == "" {
		personaParam = "STUDENT"
	}

	shortcuts, err := h.service.ListShortcuts(r.Context(), tenantID, hub.PersonaType(personaParam))
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(shortcuts)
}

// TriggerShortcut handles POST /api/v1/hub/shortcuts/trigger
func (h *Handler) TriggerShortcut(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID    string          `json:"tenant_id"`
		UserID      string          `json:"user_id"`
		Persona     hub.PersonaType `json:"persona"`
		ShortcutKey string          `json:"shortcut_key"`
		ActorID     string          `json:"actor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "Unable to parse request payload", "INVALID_PAYLOAD")
		return
	}

	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}
	if req.UserID == "" {
		req.UserID = r.Header.Get("X-User-ID")
	}
	if req.ActorID == "" {
		req.ActorID = req.UserID
	}

	if req.TenantID == "" || req.Persona == "" || req.ShortcutKey == "" {
		problem.BadRequest(w, r, "tenant_id, persona, and shortcut_key are required", "VALIDATION_ERROR")
		return
	}

	sc, err := h.service.TriggerShortcut(r.Context(), hub.TriggerShortcutRequest{
		TenantID:    req.TenantID,
		UserID:      req.UserID,
		Persona:     req.Persona,
		ShortcutKey: req.ShortcutKey,
		ActorID:     req.ActorID,
	})
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sc)
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, hub.ErrDashboardNotFound),
		errors.Is(err, hub.ErrWidgetNotFound),
		errors.Is(err, hub.ErrShortcutNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, hub.ErrUnauthorizedAccess):
		problem.Forbidden(w, r, err.Error(), "UNAUTHORIZED_ACCESS")

	case errors.Is(err, hub.ErrInvalidPersona),
		errors.Is(err, hub.ErrInvalidWidgetLayout):
		problem.BadRequest(w, r, err.Error(), "BAD_REQUEST")

	default:
		var domErr *hub.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected hub error occurred")
	}
}
