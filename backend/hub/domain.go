package hub

import (
	"time"
)

// PersonaType identifies the role cockpit.
type PersonaType string

const (
	PersonaStudent    PersonaType = "STUDENT"
	PersonaFaculty    PersonaType = "FACULTY"
	PersonaWarden     PersonaType = "WARDEN"
	PersonaLibrarian  PersonaType = "LIBRARIAN"
	PersonaAdmin      PersonaType = "ADMIN"
	PersonaSuperAdmin PersonaType = "SUPER_ADMIN"
)

// WidgetSize defines the grid span of a dashboard card.
type WidgetSize string

const (
	WidgetSizeSmall     WidgetSize = "SMALL"
	WidgetSizeMedium    WidgetSize = "MEDIUM"
	WidgetSizeLarge     WidgetSize = "LARGE"
	WidgetSizeFullWidth WidgetSize = "FULL_WIDTH"
)

// WidgetCategory organizes widgets by functional domain.
type WidgetCategory string

const (
	WidgetCategoryMetrics       WidgetCategory = "METRICS"
	WidgetCategoryActions       WidgetCategory = "ACTIONS"
	WidgetCategorySchedule      WidgetCategory = "SCHEDULE"
	WidgetCategoryCommunication WidgetCategory = "COMMUNICATION"
	WidgetCategorySafety        WidgetCategory = "SAFETY"
	WidgetCategoryFinance       WidgetCategory = "FINANCE"
)

// PersonaDashboard represents a user's customized cockpit layout.
type PersonaDashboard struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenantId"`
	UserID      string         `json:"userId"`
	Persona     PersonaType    `json:"persona"`
	LayoutTheme string         `json:"layoutTheme"`
	IsDefault   BooleanWrapper `json:"isDefault"`
	Widgets     []WidgetConfig `json:"widgets"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// BooleanWrapper provides consistent JSON boolean representation.
type BooleanWrapper bool

// WidgetConfig defines an individual widget card on the dashboard.
type WidgetConfig struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenantId"`
	DashboardID string         `json:"dashboardId"`
	WidgetKey   string         `json:"widgetKey"`
	Title       string         `json:"title"`
	Category    WidgetCategory `json:"category"`
	Size        WidgetSize     `json:"size"`
	OrderIndex  int            `json:"orderIndex"`
	IsEnabled   bool           `json:"isEnabled"`
	ConfigJSON  string         `json:"configJson,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// QuickActionShortcut defines a 1-tap navigation or action trigger.
type QuickActionShortcut struct {
	ID          string      `json:"id"`
	TenantID    string      `json:"tenantId"`
	Persona     PersonaType `json:"persona"`
	ShortcutKey string      `json:"shortcutKey"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	IconName    string      `json:"iconName"`
	TargetRoute string      `json:"targetRoute"`
	OrderIndex  int         `json:"orderIndex"`
	IsActive    bool        `json:"isActive"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

// AggregatedPersonaMetrics aggregates live KPIs across all 16 subsystems.
type AggregatedPersonaMetrics struct {
	Persona              PersonaType `json:"persona"`
	AttendanceRate       float64     `json:"attendanceRate"`
	PendingInvoicesCount int         `json:"pendingInvoicesCount"`
	UnpaidBalanceTotal   float64     `json:"unpaidBalanceTotal"`
	ActiveLibraryBorrows int         `json:"activeLibraryBorrows"`
	OverdueBooksCount    int         `json:"overdueBooksCount"`
	ActiveGatePassStatus string      `json:"activeGatePassStatus"`
	UpcomingEventsCount  int         `json:"upcomingEventsCount"`
	OpenHelpdeskTickets  int         `json:"openHelpdeskTickets"`
	ActiveSOSCount       int         `json:"activeSOSCount"`
	UnreadNoticesCount   int         `json:"unreadNoticesCount"`
	PendingAssignments   int         `json:"pendingAssignments"`
	MentorshipStatus     string      `json:"mentorshipStatus"`
}

// PersonaDashboardView returns the consolidated cockpit payload.
type PersonaDashboardView struct {
	Dashboard            *PersonaDashboard        `json:"dashboard"`
	Metrics              *AggregatedPersonaMetrics `json:"metrics"`
	Shortcuts            []QuickActionShortcut    `json:"shortcuts"`
	UnreadNoticesCount   int                      `json:"unreadNoticesCount"`
	ActiveEmergencyCount int                      `json:"activeEmergencyCount"`
}
