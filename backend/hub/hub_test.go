/**
 * BLOCK_HUB_TEST_001
 * Subsystem: Rank 17 - The Hub Root Super-App (hub)
 * Purpose:   Unit test suite covering persona cockpits, KPI aggregation, widget customization, and 1-tap shortcuts.
 */

package hub_test

import (
	"context"
	"errors"
	"testing"

	"campus/backend/hub"
)

func setupHubService() (hub.Service, *hub.MockRepository) {
	mockRepo := hub.NewMockRepository()
	srv := hub.NewService(mockRepo, nil)
	return srv, mockRepo
}

func TestHub_GetPersonaCockpit_Student(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupHubService()
	tenantID := "tenant-hub-1"
	userID := "user-student-1"

	view, err := srv.GetPersonaCockpit(ctx, tenantID, userID, hub.PersonaStudent)
	if err != nil {
		t.Fatalf("expected cockpit retrieval success, got err: %v", err)
	}
	if view == nil || view.Dashboard == nil || view.Metrics == nil {
		t.Fatalf("expected non-nil view, dashboard, and metrics")
	}

	if view.Dashboard.Persona != hub.PersonaStudent {
		t.Errorf("expected PersonaStudent, got: %s", view.Dashboard.Persona)
	}
	if view.Dashboard.UserID != userID {
		t.Errorf("expected userID %s, got: %s", userID, view.Dashboard.UserID)
	}
	if len(view.Dashboard.Widgets) == 0 {
		t.Errorf("expected default widgets to be seeded")
	}
	if view.Metrics.AttendanceRate != 89.2 {
		t.Errorf("expected AttendanceRate 89.2, got: %f", view.Metrics.AttendanceRate)
	}
	if view.Metrics.PendingInvoicesCount != 1 {
		t.Errorf("expected 1 pending invoice, got: %d", view.Metrics.PendingInvoicesCount)
	}
	if view.Metrics.ActiveGatePassStatus != "APPROVED" {
		t.Errorf("expected APPROVED gatepass status, got: %s", view.Metrics.ActiveGatePassStatus)
	}
	if len(view.Shortcuts) == 0 {
		t.Errorf("expected quick action shortcuts")
	}
}

func TestHub_GetPersonaCockpit_FacultyAndWarden(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupHubService()
	tenantID := "tenant-hub-1"

	// 1. Faculty View
	facView, err := srv.GetPersonaCockpit(ctx, tenantID, "user-fac-1", hub.PersonaFaculty)
	if err != nil {
		t.Fatalf("expected faculty cockpit success, got err: %v", err)
	}
	if facView.Dashboard.Persona != hub.PersonaFaculty {
		t.Errorf("expected PersonaFaculty, got: %s", facView.Dashboard.Persona)
	}
	if facView.Metrics.AttendanceRate != 94.5 {
		t.Errorf("expected faculty attendance 94.5, got: %f", facView.Metrics.AttendanceRate)
	}
	if facView.Metrics.PendingAssignments != 18 {
		t.Errorf("expected 18 pending assignments, got: %d", facView.Metrics.PendingAssignments)
	}

	// 2. Warden View
	wardenView, err := srv.GetPersonaCockpit(ctx, tenantID, "user-warden-1", hub.PersonaWarden)
	if err != nil {
		t.Fatalf("expected warden cockpit success, got err: %v", err)
	}
	if wardenView.Dashboard.Persona != hub.PersonaWarden {
		t.Errorf("expected PersonaWarden, got: %s", wardenView.Dashboard.Persona)
	}
	if wardenView.Metrics.ActiveGatePassStatus != "14_PENDING_APPROVAL" {
		t.Errorf("expected 14_PENDING_APPROVAL, got: %s", wardenView.Metrics.ActiveGatePassStatus)
	}
}

func TestHub_ReconfigureWidgetsAndTheme(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupHubService()
	tenantID := "tenant-hub-1"
	userID := "user-student-1"

	cockpit, err := srv.GetPersonaCockpit(ctx, tenantID, userID, hub.PersonaStudent)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 1. Reconfigure widgets
	customWidgets := []hub.WidgetConfig{
		{
			ID:         "custom-w1",
			WidgetKey:  "kpi_overview",
			Title:      "My Custom KPIs",
			Category:   hub.WidgetCategoryMetrics,
			Size:       hub.WidgetSizeLarge,
			OrderIndex: 1,
			IsEnabled:  true,
		},
		{
			ID:         "custom-w2",
			WidgetKey:  "safety_emergency",
			Title:      "SOS Quick Radar",
			Category:   hub.WidgetCategorySafety,
			Size:       hub.WidgetSizeSmall,
			OrderIndex: 2,
			IsEnabled:  false,
		},
	}

	updatedDash, err := srv.ReconfigureWidgets(ctx, hub.ReconfigureWidgetsRequest{
		TenantID:    tenantID,
		UserID:      userID,
		DashboardID: cockpit.Dashboard.ID,
		Widgets:     customWidgets,
		ActorID:     userID,
	})
	if err != nil {
		t.Fatalf("expected reconfigure success, got err: %v", err)
	}
	if len(updatedDash.Widgets) != 2 {
		t.Errorf("expected 2 widgets, got: %d", len(updatedDash.Widgets))
	}
	if updatedDash.Widgets[0].Title != "My Custom KPIs" {
		t.Errorf("expected title 'My Custom KPIs', got: %s", updatedDash.Widgets[0].Title)
	}

	// 2. Update Layout Theme
	themedDash, err := srv.UpdateLayoutTheme(ctx, hub.UpdateLayoutThemeRequest{
		TenantID:    tenantID,
		UserID:      userID,
		DashboardID: cockpit.Dashboard.ID,
		LayoutTheme: "compact_dark",
		ActorID:     userID,
	})
	if err != nil {
		t.Fatalf("expected theme update success, got err: %v", err)
	}
	if themedDash.LayoutTheme != "compact_dark" {
		t.Errorf("expected layout theme 'compact_dark', got: %s", themedDash.LayoutTheme)
	}
}

func TestHub_TriggerShortcutAndGuards(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupHubService()

	// Valid trigger
	sc, err := srv.TriggerShortcut(ctx, hub.TriggerShortcutRequest{
		TenantID:    "default",
		UserID:      "user-1",
		Persona:     hub.PersonaStudent,
		ShortcutKey: "TRIGGER_SOS",
		ActorID:     "user-1",
	})
	if err != nil {
		t.Fatalf("expected shortcut trigger success, got err: %v", err)
	}
	if sc.ShortcutKey != "TRIGGER_SOS" || sc.TargetRoute != "/sos" {
		t.Errorf("unexpected shortcut: %+v", sc)
	}

	// Non-existent shortcut key
	_, err = srv.TriggerShortcut(ctx, hub.TriggerShortcutRequest{
		TenantID:    "default",
		UserID:      "user-1",
		Persona:     hub.PersonaStudent,
		ShortcutKey: "NON_EXISTENT_KEY",
		ActorID:     "user-1",
	})
	if !errors.Is(err, hub.ErrShortcutNotFound) {
		t.Fatalf("expected ErrShortcutNotFound, got: %v", err)
	}
}

func TestHub_InvalidPersonaError(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupHubService()

	_, err := srv.GetPersonaCockpit(ctx, "tenant-1", "user-1", hub.PersonaType("INVALID_ROLE"))
	if !errors.Is(err, hub.ErrInvalidPersona) {
		t.Fatalf("expected ErrInvalidPersona, got: %v", err)
	}
}
