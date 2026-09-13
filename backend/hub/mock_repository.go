package hub

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MockRepository provides a thread-safe in-memory test double for Hub storage.
type MockRepository struct {
	mu         sync.RWMutex
	dashboards map[string]*PersonaDashboard        // key: tenantID + ":" + userID + ":" + persona
	shortcuts  map[string][]QuickActionShortcut   // key: tenantID + ":" + persona
}

// NewMockRepository constructs a initialized mock repository with default shortcuts.
func NewMockRepository() *MockRepository {
	repo := &MockRepository{
		dashboards: make(map[string]*PersonaDashboard),
		shortcuts:  make(map[string][]QuickActionShortcut),
	}
	repo.seedDefaultShortcuts()
	return repo
}

func (m *MockRepository) seedDefaultShortcuts() {
	m.shortcuts["default:STUDENT"] = []QuickActionShortcut{
		{
			ID:          "sc-s1",
			TenantID:    "default",
			Persona:     PersonaStudent,
			ShortcutKey: "TRIGGER_SOS",
			Title:       "SOS Emergency",
			Description: "1-Tap emergency distress trigger",
			IconName:    "shield-alert",
			TargetRoute: "/sos",
			OrderIndex:  1,
			IsActive:    true,
		},
		{
			ID:          "sc-s2",
			TenantID:    "default",
			Persona:     PersonaStudent,
			ShortcutKey: "MARK_ATTENDANCE",
			Title:       "Geo Attendance",
			Description: "Mark geofenced biometric check-in",
			IconName:    "scan-face",
			TargetRoute: "/attendance",
			OrderIndex:  2,
			IsActive:    true,
		},
		{
			ID:          "sc-s3",
			TenantID:    "default",
			Persona:     PersonaStudent,
			ShortcutKey: "PAY_FEE",
			Title:       "Pay Invoice",
			Description: "Online tuition & hostel fee settlement",
			IconName:    "credit-card",
			TargetRoute: "/billing",
			OrderIndex:  3,
			IsActive:    true,
		},
		{
			ID:          "sc-s4",
			TenantID:    "default",
			Persona:     PersonaStudent,
			ShortcutKey: "GATE_PASS_APPLY",
			Title:       "Hostel Gate Pass",
			Description: "Apply for night-out or day leave",
			IconName:    "door-open",
			TargetRoute: "/hostel",
			OrderIndex:  4,
			IsActive:    true,
		},
	}

	m.shortcuts["default:FACULTY"] = []QuickActionShortcut{
		{
			ID:          "sc-f1",
			TenantID:    "default",
			Persona:     PersonaFaculty,
			ShortcutKey: "TAKE_ATTENDANCE",
			Title:       "Roll Call Session",
			Description: "Start digital roll call for lecture",
			IconName:    "users",
			TargetRoute: "/attendance",
			OrderIndex:  1,
			IsActive:    true,
		},
		{
			ID:          "sc-f2",
			TenantID:    "default",
			Persona:     PersonaFaculty,
			ShortcutKey: "PUBLISH_NOTICE",
			Title:       "Post Notice",
			Description: "Publish departmental circular",
			IconName:    "bell",
			TargetRoute: "/notices",
			OrderIndex:  2,
			IsActive:    true,
		},
		{
			ID:          "sc-f3",
			TenantID:    "default",
			Persona:     PersonaFaculty,
			ShortcutKey: "GRADE_ASSIGNMENT",
			Title:       "Grade Submissions",
			Description: "Review & score pending homework",
			IconName:    "file-check",
			TargetRoute: "/studyhub",
			OrderIndex:  3,
			IsActive:    true,
		},
	}
}

func (m *MockRepository) GetOrCreateDashboard(ctx context.Context, tenantID, userID string, persona PersonaType) (*PersonaDashboard, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s:%s:%s", tenantID, userID, persona)
	if dash, exists := m.dashboards[key]; exists {
		return dash, nil
	}

	// Create default dashboard with starter widgets
	dashID := fmt.Sprintf("dash-%d", time.Now().UnixNano())
	defaultWidgets := []WidgetConfig{
		{
			ID:          fmt.Sprintf("w-1-%d", time.Now().UnixNano()),
			TenantID:    tenantID,
			DashboardID: dashID,
			WidgetKey:   "kpi_overview",
			Title:       "Vital KPIs Overview",
			Category:    WidgetCategoryMetrics,
			Size:        WidgetSizeFullWidth,
			OrderIndex:  1,
			IsEnabled:   true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          fmt.Sprintf("w-2-%d", time.Now().UnixNano()),
			TenantID:    tenantID,
			DashboardID: dashID,
			WidgetKey:   "quick_actions",
			Title:       "Quick Action Shortcuts",
			Category:    WidgetCategoryActions,
			Size:        WidgetSizeMedium,
			OrderIndex:  2,
			IsEnabled:   true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          fmt.Sprintf("w-3-%d", time.Now().UnixNano()),
			TenantID:    tenantID,
			DashboardID: dashID,
			WidgetKey:   "safety_emergency",
			Title:       "Campus Safety & SOS Status",
			Category:    WidgetCategorySafety,
			Size:        WidgetSizeMedium,
			OrderIndex:  3,
			IsEnabled:   true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	newDash := &PersonaDashboard{
		ID:          dashID,
		TenantID:    tenantID,
		UserID:      userID,
		Persona:     persona,
		LayoutTheme: "default",
		IsDefault:   true,
		Widgets:     defaultWidgets,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.dashboards[key] = newDash
	return newDash, nil
}

func (m *MockRepository) GetDashboard(ctx context.Context, tenantID, dashboardID string) (*PersonaDashboard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, d := range m.dashboards {
		if d.TenantID == tenantID && d.ID == dashboardID {
			return d, nil
		}
	}
	return nil, ErrDashboardNotFound
}

func (m *MockRepository) UpdateDashboardLayout(ctx context.Context, tenantID, dashboardID string, layoutTheme string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, d := range m.dashboards {
		if d.TenantID == tenantID && d.ID == dashboardID {
			d.LayoutTheme = layoutTheme
			d.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrDashboardNotFound
}

func (m *MockRepository) SaveWidgetConfigs(ctx context.Context, tenantID, dashboardID string, widgets []WidgetConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, d := range m.dashboards {
		if d.TenantID == tenantID && d.ID == dashboardID {
			d.Widgets = widgets
			d.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrDashboardNotFound
}

func (m *MockRepository) ListShortcuts(ctx context.Context, tenantID string, persona PersonaType) ([]QuickActionShortcut, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", tenantID, persona)
	if items, exists := m.shortcuts[key]; exists {
		return items, nil
	}

	// Fallback to default
	defaultKey := fmt.Sprintf("default:%s", persona)
	if items, exists := m.shortcuts[defaultKey]; exists {
		return items, nil
	}

	return []QuickActionShortcut{}, nil
}

func (m *MockRepository) CreateShortcut(ctx context.Context, shortcut QuickActionShortcut) (*QuickActionShortcut, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s:%s", shortcut.TenantID, shortcut.Persona)
	shortcut.ID = fmt.Sprintf("sc-%d", time.Now().UnixNano())
	shortcut.CreatedAt = time.Now()
	shortcut.UpdatedAt = time.Now()

	m.shortcuts[key] = append(m.shortcuts[key], shortcut)
	return &shortcut, nil
}

func (m *MockRepository) GetAggregatedMetrics(ctx context.Context, tenantID, userID string, persona PersonaType) (*AggregatedPersonaMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return realistic synthesized KPIs based on persona
	switch persona {
	case PersonaStudent:
		return &AggregatedPersonaMetrics{
			Persona:              PersonaStudent,
			AttendanceRate:       89.2,
			PendingInvoicesCount: 1,
			UnpaidBalanceTotal:   25000.0,
			ActiveLibraryBorrows: 2,
			OverdueBooksCount:    0,
			ActiveGatePassStatus: "APPROVED",
			UpcomingEventsCount:  3,
			OpenHelpdeskTickets:  1,
			ActiveSOSCount:       0,
			UnreadNoticesCount:   4,
			PendingAssignments:   2,
			MentorshipStatus:     "ON_TRACK",
		}, nil
	case PersonaFaculty:
		return &AggregatedPersonaMetrics{
			Persona:              PersonaFaculty,
			AttendanceRate:       94.5,
			PendingInvoicesCount: 0,
			UnpaidBalanceTotal:   0.0,
			ActiveLibraryBorrows: 5,
			OverdueBooksCount:    0,
			ActiveGatePassStatus: "N/A",
			UpcomingEventsCount:  5,
			OpenHelpdeskTickets:  2,
			ActiveSOSCount:       0,
			UnreadNoticesCount:   2,
			PendingAssignments:   18,
			MentorshipStatus:     "8_MENTEES_ACTIVE",
		}, nil
	case PersonaWarden:
		return &AggregatedPersonaMetrics{
			Persona:              PersonaWarden,
			AttendanceRate:       0.0,
			PendingInvoicesCount: 0,
			UnpaidBalanceTotal:   0.0,
			ActiveLibraryBorrows: 0,
			OverdueBooksCount:    0,
			ActiveGatePassStatus: "14_PENDING_APPROVAL",
			UpcomingEventsCount:  1,
			OpenHelpdeskTickets:  4,
			ActiveSOSCount:       0,
			UnreadNoticesCount:   1,
			PendingAssignments:   0,
			MentorshipStatus:     "N/A",
		}, nil
	default:
		return &AggregatedPersonaMetrics{
			Persona:              persona,
			AttendanceRate:       85.0,
			PendingInvoicesCount: 0,
			UnpaidBalanceTotal:   0.0,
			ActiveLibraryBorrows: 0,
			OverdueBooksCount:    0,
			ActiveGatePassStatus: "N/A",
			UpcomingEventsCount:  2,
			OpenHelpdeskTickets:  0,
			ActiveSOSCount:       0,
			UnreadNoticesCount:   3,
			PendingAssignments:   0,
			MentorshipStatus:     "ACTIVE",
		}, nil
	}
}
