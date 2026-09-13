/**
 * BLOCK_API_HUB_HANDLER_TEST_001
 * Subsystem: Rank 17 - The Hub Root Super-App (hub)
 * Purpose:   HTTP integration tests verifying REST endpoints for persona cockpits, metrics aggregation, widget layout customizer, and shortcuts.
 */

package hub_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apiHub "campus/api/http/hub"
	backendHub "campus/backend/hub"
)

func setupTestRouter() (chi.Router, backendHub.Service) {
	r := chi.NewRouter()
	repo := backendHub.NewMockRepository()
	svc := backendHub.NewService(repo, nil)
	handler := apiHub.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_HubCockpitLifecycle(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Get Student Cockpit
	reqCockpit := httptest.NewRequest("GET", "/api/v1/hub/cockpit?tenant_id=tenant-hub-http&user_id=user-student-1&persona=STUDENT", nil)
	recCockpit := httptest.NewRecorder()
	router.ServeHTTP(recCockpit, reqCockpit)

	if recCockpit.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for student cockpit, got %d: %s", recCockpit.Code, recCockpit.Body.String())
	}

	var cockpitResp backendHub.PersonaDashboardView
	if err := json.NewDecoder(recCockpit.Body).Decode(&cockpitResp); err != nil {
		t.Fatalf("Failed to parse cockpit response: %v", err)
	}
	if cockpitResp.Dashboard.Persona != backendHub.PersonaStudent {
		t.Errorf("Expected PersonaStudent, got: %s", cockpitResp.Dashboard.Persona)
	}
	if cockpitResp.Metrics.AttendanceRate != 89.2 {
		t.Errorf("Expected AttendanceRate 89.2, got: %f", cockpitResp.Metrics.AttendanceRate)
	}

	// 2. Reconfigure Dashboard Widgets
	reconfigPayload := map[string]interface{}{
		"tenant_id":    "tenant-hub-http",
		"user_id":      "user-student-1",
		"dashboard_id": cockpitResp.Dashboard.ID,
		"actor_id":     "user-student-1",
		"widgets": []map[string]interface{}{
			{
				"widgetKey":  "kpi_overview",
				"title":      "My Vital KPIs",
				"category":   "METRICS",
				"size":       "FULL_WIDTH",
				"orderIndex": 1,
				"isEnabled":  true,
			},
		},
	}
	body, _ := json.Marshal(reconfigPayload)
	reqReconfig := httptest.NewRequest("POST", "/api/v1/hub/widgets/reconfigure", bytes.NewReader(body))
	reqReconfig.Header.Set("Content-Type", "application/json")
	recReconfig := httptest.NewRecorder()
	router.ServeHTTP(recReconfig, reqReconfig)

	if recReconfig.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for widget reconfigure, got %d: %s", recReconfig.Code, recReconfig.Body.String())
	}

	// 3. Trigger Quick Action Shortcut
	triggerPayload := map[string]interface{}{
		"tenant_id":    "default",
		"user_id":      "user-student-1",
		"persona":      "STUDENT",
		"shortcut_key": "TRIGGER_SOS",
		"actor_id":     "user-student-1",
	}
	bodyTrigger, _ := json.Marshal(triggerPayload)
	reqTrigger := httptest.NewRequest("POST", "/api/v1/hub/shortcuts/trigger", bytes.NewReader(bodyTrigger))
	reqTrigger.Header.Set("Content-Type", "application/json")
	recTrigger := httptest.NewRecorder()
	router.ServeHTTP(recTrigger, reqTrigger)

	if recTrigger.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for shortcut trigger, got %d: %s", recTrigger.Code, recTrigger.Body.String())
	}
}
