/**
 * BLOCK_SOS_TEST_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   Unit test suite covering 1-tap SOS trigger, coordinate bounds, responder dispatch, and incident state transitions.
 */

package sos_test

import (
	"context"
	"errors"
	"testing"

	"campus/backend/sos"
)

func setupSOSService() (sos.Service, *sos.MockRepository) {
	mockRepo := sos.NewMockRepository()
	srv := sos.NewService(mockRepo, nil)
	return srv, mockRepo
}

func TestSOSTriggerAndCoordinateValidation(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupSOSService()
	tenantID := "tenant-alpha"

	// Invalid coordinates (latitude > 90)
	_, err := srv.TriggerSOS(ctx, sos.TriggerSOSRequest{
		TenantID:            tenantID,
		UserID:              "student-1",
		EmergencyType:       sos.EmergencyMedical,
		Latitude:            95.123,
		Longitude:           77.123,
		LocationDescription: "Invalid coords",
	})
	if !errors.Is(err, sos.ErrInvalidCoordinates) {
		t.Fatalf("expected ErrInvalidCoordinates, got: %v", err)
	}

	// Valid SOS trigger
	inc, err := srv.TriggerSOS(ctx, sos.TriggerSOSRequest{
		TenantID:            tenantID,
		UserID:              "student-1",
		EmergencyType:       sos.EmergencyMedical,
		Latitude:            12.9716,
		Longitude:           77.5946,
		LocationDescription: "Hostel Block B 3rd Floor Corridor",
	})
	if err != nil {
		t.Fatalf("expected trigger success, got: %v", err)
	}

	if inc.ID == "" || inc.AlertNumber != "SOS-2026-00001" || inc.Status != sos.IncidentTriggered {
		t.Fatalf("unexpected incident payload: %+v", inc)
	}

	// Second trigger increments alert sequence
	inc2, err := srv.TriggerSOS(ctx, sos.TriggerSOSRequest{
		TenantID:            tenantID,
		UserID:              "student-2",
		EmergencyType:       sos.EmergencyFire,
		Latitude:            12.9720,
		Longitude:           77.5950,
		LocationDescription: "Chemistry Lab 204",
	})
	if err != nil || inc2.AlertNumber != "SOS-2026-00002" {
		t.Fatalf("expected sequential alert number SOS-2026-00002, got: %v", inc2)
	}
}

func TestSOSIncidentLifecycle_AcknowledgeAndDispatch(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupSOSService()
	tenantID := "tenant-alpha"

	inc, _ := srv.TriggerSOS(ctx, sos.TriggerSOSRequest{
		TenantID:            tenantID,
		UserID:              "student-1",
		EmergencyType:       sos.EmergencyMedical,
		Latitude:            12.9716,
		Longitude:           77.5946,
		LocationDescription: "Main Ground",
	})

	// 1. Acknowledge
	ack, err := srv.AcknowledgeIncident(ctx, tenantID, inc.ID, "control-room-officer-1")
	if err != nil || ack.Status != sos.IncidentAcknowledged || ack.AcknowledgedAt == nil {
		t.Fatalf("expected acknowledged incident, got: %+v, err: %v", ack, err)
	}

	// 2. Dispatch Responder
	responder, err := srv.DispatchResponder(ctx, sos.DispatchResponderRequest{
		TenantID:    tenantID,
		IncidentID:  inc.ID,
		ResponderID: "security-guard-1",
		Role:        sos.RoleCampusSecurity,
		ActorID:     "control-room-officer-1",
	})
	if err != nil || responder.Status != sos.ResponderAssigned {
		t.Fatalf("expected responder dispatch success, got: %+v, err: %v", responder, err)
	}

	// Incident status should now be DISPATCHED
	incDispatched, _ := srv.GetIncident(ctx, tenantID, inc.ID)
	if incDispatched.Status != sos.IncidentDispatched {
		t.Fatalf("expected incident status DISPATCHED, got: %s", incDispatched.Status)
	}

	// Duplicate dispatch of same responder to same incident -> REJECTED
	_, err = srv.DispatchResponder(ctx, sos.DispatchResponderRequest{
		TenantID:    tenantID,
		IncidentID:  inc.ID,
		ResponderID: "security-guard-1",
		Role:        sos.RoleCampusSecurity,
		ActorID:     "control-room-officer-1",
	})
	if !errors.Is(err, sos.ErrResponderAlreadyDispatched) {
		t.Fatalf("expected ErrResponderAlreadyDispatched, got: %v", err)
	}
}

func TestSOSResponderOnSceneAndResolution(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupSOSService()
	tenantID := "tenant-alpha"

	inc, _ := srv.TriggerSOS(ctx, sos.TriggerSOSRequest{
		TenantID:            tenantID,
		UserID:              "student-1",
		EmergencyType:       sos.EmergencySecurityThreat,
		Latitude:            12.9716,
		Longitude:           77.5946,
		LocationDescription: "Library Back Alley",
	})

	_, _ = srv.DispatchResponder(ctx, sos.DispatchResponderRequest{
		TenantID:    tenantID,
		IncidentID:  inc.ID,
		ResponderID: "paramedic-1",
		Role:        sos.RoleParamedic,
		ActorID:     "control-room-1",
	})

	// Paramedic arrives ON_SCENE
	updatedResp, err := srv.UpdateResponderStatus(ctx, sos.UpdateResponderStatusRequest{
		TenantID:    tenantID,
		IncidentID:  inc.ID,
		ResponderID: "paramedic-1",
		Status:      sos.ResponderOnScene,
	})
	if err != nil || updatedResp.Status != sos.ResponderOnScene || updatedResp.ArrivedAt == nil {
		t.Fatalf("expected responder on scene with arrival time, got: %+v, err: %v", updatedResp, err)
	}

	// Incident status should auto-update to ON_SCENE
	incOnScene, _ := srv.GetIncident(ctx, tenantID, inc.ID)
	if incOnScene.Status != sos.IncidentOnScene {
		t.Fatalf("expected incident ON_SCENE, got: %s", incOnScene.Status)
	}

	// Resolve incident
	resolved, err := srv.ResolveIncident(ctx, sos.ResolveIncidentRequest{
		TenantID:        tenantID,
		IncidentID:      inc.ID,
		ResolvedByID:    "paramedic-1",
		Status:          sos.IncidentResolved,
		ResolutionNotes: "First aid administered. Student in stable condition accompanied to health center.",
	})
	if err != nil || resolved.Status != sos.IncidentResolved || resolved.ResolvedAt == nil {
		t.Fatalf("expected incident resolved with notes, got: %+v, err: %v", resolved, err)
	}

	// Modifying resolved incident -> REJECTED
	_, err = srv.AcknowledgeIncident(ctx, tenantID, inc.ID, "officer-2")
	if !errors.Is(err, sos.ErrIncidentAlreadyClosed) {
		t.Fatalf("expected ErrIncidentAlreadyClosed on resolved incident, got: %v", err)
	}
}

func TestSOSFalseAlarmResolution(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupSOSService()
	tenantID := "tenant-alpha"

	inc, _ := srv.TriggerSOS(ctx, sos.TriggerSOSRequest{
		TenantID:            tenantID,
		UserID:              "student-1",
		EmergencyType:       sos.EmergencyOther,
		Latitude:            12.9716,
		Longitude:           77.5946,
		LocationDescription: "Cafeteria",
	})

	resolved, err := srv.ResolveIncident(ctx, sos.ResolveIncidentRequest{
		TenantID:        tenantID,
		IncidentID:      inc.ID,
		ResolvedByID:    "security-admin",
		Status:          sos.IncidentFalseAlarm,
		ResolutionNotes: "Accidental button press tested by student.",
	})
	if err != nil || resolved.Status != sos.IncidentFalseAlarm {
		t.Fatalf("expected FALSE_ALARM status, got: %+v, err: %v", resolved, err)
	}
}
