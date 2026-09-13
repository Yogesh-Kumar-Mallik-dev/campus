package hub

import (
	"context"
)

// Repository defines the contract for persisting Hub dashboard configurations, widgets, and shortcuts.
type Repository interface {
	GetOrCreateDashboard(ctx context.Context, tenantID, userID string, persona PersonaType) (*PersonaDashboard, error)
	GetDashboard(ctx context.Context, tenantID, dashboardID string) (*PersonaDashboard, error)
	UpdateDashboardLayout(ctx context.Context, tenantID, dashboardID string, layoutTheme string) error
	SaveWidgetConfigs(ctx context.Context, tenantID, dashboardID string, widgets []WidgetConfig) error
	ListShortcuts(ctx context.Context, tenantID string, persona PersonaType) ([]QuickActionShortcut, error)
	CreateShortcut(ctx context.Context, shortcut QuickActionShortcut) (*QuickActionShortcut, error)
	GetAggregatedMetrics(ctx context.Context, tenantID, userID string, persona PersonaType) (*AggregatedPersonaMetrics, error)
}
