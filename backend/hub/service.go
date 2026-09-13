/**
 * BLOCK_HUB_SERVICE_001
 * Subsystem: Rank 17 - The Hub Root Super-App (hub)
 * Purpose:   Unified Persona Cockpit orchestration, cross-subsystem metrics aggregation, widget customizer, and 1-tap action shortcuts.
 */

package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"campus/backend/audit"
)

type Service interface {
	GetPersonaCockpit(ctx context.Context, tenantID, userID string, persona PersonaType) (*PersonaDashboardView, error)
	ReconfigureWidgets(ctx context.Context, req ReconfigureWidgetsRequest) (*PersonaDashboard, error)
	UpdateLayoutTheme(ctx context.Context, req UpdateLayoutThemeRequest) (*PersonaDashboard, error)
	TriggerShortcut(ctx context.Context, req TriggerShortcutRequest) (*QuickActionShortcut, error)
	ListShortcuts(ctx context.Context, tenantID string, persona PersonaType) ([]QuickActionShortcut, error)
}

type service struct {
	repo     Repository
	auditSub audit.Subscriber
}

func NewService(repo Repository, auditSub audit.Subscriber) Service {
	return &service{
		repo:     repo,
		auditSub: auditSub,
	}
}

func (s *service) logAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
	if s.auditSub == nil {
		return
	}
	var metaRaw json.RawMessage
	if meta != nil {
		if b, err := json.Marshal(meta); err == nil {
			metaRaw = b
		}
	}
	_ = s.auditSub.Enqueue(audit.RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      &actorID,
		ActorType:    audit.ActorType(actorType),
		Action:       action,
		ResourceType: resType,
		ResourceID:   &resID,
		Status:       status,
		Metadata:     metaRaw,
	})
}

func validatePersona(p PersonaType) bool {
	switch p {
	case PersonaStudent, PersonaFaculty, PersonaWarden, PersonaLibrarian, PersonaAdmin, PersonaSuperAdmin:
		return true
	default:
		return false
	}
}

func (s *service) GetPersonaCockpit(ctx context.Context, tenantID, userID string, persona PersonaType) (*PersonaDashboardView, error) {
	if !validatePersona(persona) {
		return nil, ErrInvalidPersona
	}

	dashboard, err := s.repo.GetOrCreateDashboard(ctx, tenantID, userID, persona)
	if err != nil {
		return nil, err
	}

	metrics, err := s.repo.GetAggregatedMetrics(ctx, tenantID, userID, persona)
	if err != nil {
		return nil, err
	}

	shortcuts, err := s.repo.ListShortcuts(ctx, tenantID, persona)
	if err != nil {
		return nil, err
	}

	view := &PersonaDashboardView{
		Dashboard:            dashboard,
		Metrics:              metrics,
		Shortcuts:            shortcuts,
		UnreadNoticesCount:   metrics.UnreadNoticesCount,
		ActiveEmergencyCount: metrics.ActiveSOSCount,
	}

	s.logAudit(tenantID, userID, string(persona), "hub:cockpit:viewed", "hub_persona_dashboard", dashboard.ID, audit.StatusSuccess, map[string]interface{}{
		"persona": string(persona),
	})

	return view, nil
}

type ReconfigureWidgetsRequest struct {
	TenantID    string         `json:"tenantId"`
	UserID      string         `json:"userId"`
	DashboardID string         `json:"dashboardId"`
	Widgets     []WidgetConfig `json:"widgets"`
	ActorID     string         `json:"actorId"`
}

func (s *service) ReconfigureWidgets(ctx context.Context, req ReconfigureWidgetsRequest) (*PersonaDashboard, error) {
	dashboard, err := s.repo.GetDashboard(ctx, req.TenantID, req.DashboardID)
	if err != nil {
		return nil, err
	}

	if dashboard.UserID != req.UserID && req.ActorID != req.UserID {
		return nil, ErrUnauthorizedAccess
	}

	for i := range req.Widgets {
		req.Widgets[i].DashboardID = req.DashboardID
		req.Widgets[i].TenantID = req.TenantID
		req.Widgets[i].UpdatedAt = time.Now().UTC()
	}

	if err := s.repo.SaveWidgetConfigs(ctx, req.TenantID, req.DashboardID, req.Widgets); err != nil {
		return nil, err
	}

	dashboard.Widgets = req.Widgets
	dashboard.UpdatedAt = time.Now().UTC()

	s.logAudit(req.TenantID, req.ActorID, "USER", "hub:widgets:reconfigured", "hub_persona_dashboard", req.DashboardID, audit.StatusSuccess, map[string]interface{}{
		"widgetCount": len(req.Widgets),
	})

	return dashboard, nil
}

type UpdateLayoutThemeRequest struct {
	TenantID    string `json:"tenantId"`
	UserID      string `json:"userId"`
	DashboardID string `json:"dashboardId"`
	LayoutTheme string `json:"layoutTheme"`
	ActorID     string `json:"actorId"`
}

func (s *service) UpdateLayoutTheme(ctx context.Context, req UpdateLayoutThemeRequest) (*PersonaDashboard, error) {
	dashboard, err := s.repo.GetDashboard(ctx, req.TenantID, req.DashboardID)
	if err != nil {
		return nil, err
	}

	if dashboard.UserID != req.UserID && req.ActorID != req.UserID {
		return nil, ErrUnauthorizedAccess
	}

	if err := s.repo.UpdateDashboardLayout(ctx, req.TenantID, req.DashboardID, req.LayoutTheme); err != nil {
		return nil, err
	}

	dashboard.LayoutTheme = req.LayoutTheme
	dashboard.UpdatedAt = time.Now().UTC()

	s.logAudit(req.TenantID, req.ActorID, "USER", "hub:theme:updated", "hub_persona_dashboard", req.DashboardID, audit.StatusSuccess, map[string]interface{}{
		"layoutTheme": req.LayoutTheme,
	})

	return dashboard, nil
}

type TriggerShortcutRequest struct {
	TenantID    string      `json:"tenantId"`
	UserID      string      `json:"userId"`
	Persona     PersonaType `json:"persona"`
	ShortcutKey string      `json:"shortcutKey"`
	ActorID     string      `json:"actorId"`
}

func (s *service) TriggerShortcut(ctx context.Context, req TriggerShortcutRequest) (*QuickActionShortcut, error) {
	if !validatePersona(req.Persona) {
		return nil, ErrInvalidPersona
	}

	shortcuts, err := s.repo.ListShortcuts(ctx, req.TenantID, req.Persona)
	if err != nil {
		return nil, err
	}

	var matched *QuickActionShortcut
	for _, sc := range shortcuts {
		if sc.ShortcutKey == req.ShortcutKey {
			c := sc
			matched = &c
			break
		}
	}

	if matched == nil {
		return nil, fmt.Errorf("%w: key=%s", ErrShortcutNotFound, req.ShortcutKey)
	}

	s.logAudit(req.TenantID, req.ActorID, string(req.Persona), "hub:shortcut:triggered", "hub_quick_action_shortcut", matched.ID, audit.StatusSuccess, map[string]interface{}{
		"shortcutKey": req.ShortcutKey,
		"targetRoute": matched.TargetRoute,
	})

	return matched, nil
}

func (s *service) ListShortcuts(ctx context.Context, tenantID string, persona PersonaType) ([]QuickActionShortcut, error) {
	if !validatePersona(persona) {
		return nil, ErrInvalidPersona
	}
	return s.repo.ListShortcuts(ctx, tenantID, persona)
}
