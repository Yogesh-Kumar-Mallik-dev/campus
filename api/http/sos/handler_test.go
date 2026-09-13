/**
 * BLOCK_API_SOS_HANDLER_TEST_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   HTTP integration tests verifying REST endpoints for emergency trigger, responder dispatch, status updates, and resolution.
 */

package sos_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apiSOS "campus/api/http/sos"
	backendSOS "campus/backend/sos"
)

func setupTestRouter() (chi.Router, backendSOS.Service) {
	r := chi.NewRouter()
	repo := backendSOS.NewMockRepository()
	svc := backendSOS.NewService(repo, nil)
	handler := apiSOS.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_SOSTriggerDispatchAndResolution(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. 1-Tap Trigger SOS
	triggerPayload := map[string]interface{}{
		"tenant_id":            "tenant-sos-http",
		"user_id":              "student-danger-1",
		"emergency_type":       "MEDICAL",
		"latitude":             12.9716,
		"longitude":            77.5946,
		"location_description": "Library 2nd Floor Reading Room",
	}
	body, _ := json.Marshal(triggerPayload)
	req := httptest.NewRequest("POST", "/api/v1/sos/trigger", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for SOS trigger, got %d: %s", rec.Code, rec.Body.String())
	}

	var incident backendSOS.SOSIncident
	_ = json.Unmarshal(rec.Body.Bytes(), &incident)

	if incident.AlertNumber != "SOS-2026-00001" || incident.Status != backendSOS.IncidentTriggered {
		t.Fatalf("Unexpected incident fields: %+v", incident)
	}

	// 2. Acknowledge Incident
	ackPayload := map[string]interface{}{
		"tenant_id": "tenant-sos-http",
		"actor_id":  "control-room-officer",
	}
	ackBody, _ := json.Marshal(ackPayload)
	reqAck := httptest.NewRequest("POST", "/api/v1/sos/incidents/"+incident.ID+"/acknowledge", bytes.NewReader(ackBody))
	reqAck.Header.Set("Content-Type", "application/json")
	recAck := httptest.NewRecorder()
	router.ServeHTTP(recAck, reqAck)

	if recAck.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for acknowledge, got %d: %s", recAck.Code, recAck.Body.String())
	}

	// 3. Dispatch Paramedic Responder
	dispPayload := map[string]interface{}{
		"tenant_id":    "tenant-sos-http",
		"responder_id": "paramedic-unit-1",
		"role":         "PARAMEDIC",
		"actor_id":     "control-room-officer",
	}
	dispBody, _ := json.Marshal(dispPayload)
	reqDisp := httptest.NewRequest("POST", "/api/v1/sos/incidents/"+incident.ID+"/dispatch", bytes.NewReader(dispBody))
	reqDisp.Header.Set("Content-Type", "application/json")
	recDisp := httptest.NewRecorder()
	router.ServeHTTP(recDisp, reqDisp)

	if recDisp.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for responder dispatch, got %d: %s", recDisp.Code, recDisp.Body.String())
	}

	// 4. Update Responder Status to ON_SCENE
	statusPayload := map[string]interface{}{
		"tenant_id": "tenant-sos-http",
		"status":    "ON_SCENE",
	}
	statusBody, _ := json.Marshal(statusPayload)
	reqStatus := httptest.NewRequest("POST", "/api/v1/sos/incidents/"+incident.ID+"/responders/paramedic-unit-1/status", bytes.NewReader(statusBody))
	reqStatus.Header.Set("Content-Type", "application/json")
	recStatus := httptest.NewRecorder()
	router.ServeHTTP(recStatus, reqStatus)

	if recStatus.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for responder status, got %d: %s", recStatus.Code, recStatus.Body.String())
	}

	// 5. Resolve Incident
	resolvePayload := map[string]interface{}{
		"tenant_id":        "tenant-sos-http",
		"resolved_by_id":   "paramedic-unit-1",
		"status":           "RESOLVED",
		"resolution_notes": "Student safely transferred to medical clinic.",
	}
	resolveBody, _ := json.Marshal(resolvePayload)
	reqResolve := httptest.NewRequest("POST", "/api/v1/sos/incidents/"+incident.ID+"/resolve", bytes.NewReader(resolveBody))
	reqResolve.Header.Set("Content-Type", "application/json")
	recResolve := httptest.NewRecorder()
	router.ServeHTTP(recResolve, reqResolve)

	if recResolve.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for resolve incident, got %d: %s", recResolve.Code, recResolve.Body.String())
	}
}
