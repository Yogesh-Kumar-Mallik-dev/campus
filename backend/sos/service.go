/**
 * BLOCK_SOS_SERVICE_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   Core business logic orchestration, telemetry broadcasts, responder dispatches, and audit publishing.
 */

package sos

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"campus/backend/audit"
)

type Service interface {
	TriggerSOS(ctx context.Context, req TriggerSOSRequest) (*SOSIncident, error)
	AcknowledgeIncident(ctx context.Context, tenantID, incidentID, actorID string) (*SOSIncident, error)
	DispatchResponder(ctx context.Context, req DispatchResponderRequest) (*SOSDispatchResponder, error)
	UpdateResponderStatus(ctx context.Context, req UpdateResponderStatusRequest) (*SOSDispatchResponder, error)
	ResolveIncident(ctx context.Context, req ResolveIncidentRequest) (*SOSIncident, error)
	GetIncident(ctx context.Context, tenantID, id string) (*SOSIncident, error)
	ListIncidents(ctx context.Context, filter SOSIncidentFilter) ([]SOSIncident, error)
	ListResponders(ctx context.Context, tenantID, incidentID string) ([]SOSDispatchResponder, error)
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

type TriggerSOSRequest struct {
	TenantID            string           `json:"tenantId"`
	UserID              string           `json:"userId"`
	EmergencyType       SOSEmergencyType `json:"emergencyType"`
	Latitude            float64          `json:"latitude"`
	Longitude           float64          `json:"longitude"`
	LocationDescription string           `json:"locationDescription"`
}

func (s *service) TriggerSOS(ctx context.Context, req TriggerSOSRequest) (*SOSIncident, error) {
	if err := ValidateCoordinates(req.Latitude, req.Longitude); err != nil {
		return nil, err
	}

	if req.EmergencyType == "" {
		req.EmergencyType = EmergencyMedical
	}

	seq, err := s.repo.GetNextAlertSequence(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	alertNumber := FormatAlertNumber(now.Year(), seq)

	incident := &SOSIncident{
		ID:                  uuid.New().String(),
		TenantID:            req.TenantID,
		AlertNumber:         alertNumber,
		UserID:              req.UserID,
		EmergencyType:       req.EmergencyType,
		Latitude:            req.Latitude,
		Longitude:           req.Longitude,
		LocationDescription: req.LocationDescription,
		Status:              IncidentTriggered,
		TriggeredAt:         now,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.repo.CreateIncident(ctx, incident); err != nil {
		return nil, err
	}

	// Auto-create initial emergency push broadcast
	_ = s.repo.CreateBroadcast(ctx, &SOSTelemetryBroadcast{
		ID:             uuid.New().String(),
		TenantID:       req.TenantID,
		IncidentID:     incident.ID,
		Channel:        ChannelPushNotification,
		Payload:        "EMERGENCY ALERT: " + req.LocationDescription,
		DeliveredCount: 1,
		BroadcastAt:    now,
		CreatedAt:      now,
	})

	s.logAudit(req.TenantID, req.UserID, "STUDENT", "sos:alert:triggered", "sos_incident", incident.ID, audit.StatusSuccess, map[string]interface{}{
		"alertNumber":   incident.AlertNumber,
		"emergencyType": incident.EmergencyType,
		"latitude":      incident.Latitude,
		"longitude":     incident.Longitude,
	})

	return incident, nil
}

func (s *service) AcknowledgeIncident(ctx context.Context, tenantID, incidentID, actorID string) (*SOSIncident, error) {
	incident, err := s.repo.GetIncidentByID(ctx, tenantID, incidentID)
	if err != nil {
		return nil, err
	}

	if err := ValidateIncidentTransition(incident.Status, IncidentAcknowledged); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	incident.Status = IncidentAcknowledged
	incident.AcknowledgedAt = &now
	incident.UpdatedAt = now

	if err := s.repo.UpdateIncident(ctx, incident); err != nil {
		return nil, err
	}

	s.logAudit(tenantID, actorID, "STAFF", "sos:alert:acknowledged", "sos_incident", incident.ID, audit.StatusSuccess, map[string]interface{}{
		"alertNumber": incident.AlertNumber,
	})

	return incident, nil
}

type DispatchResponderRequest struct {
	TenantID    string           `json:"tenantId"`
	IncidentID  string           `json:"incidentId"`
	ResponderID string           `json:"responderId"`
	Role        SOSResponderRole `json:"role"`
	ActorID     string           `json:"actorId"`
}

func (s *service) DispatchResponder(ctx context.Context, req DispatchResponderRequest) (*SOSDispatchResponder, error) {
	incident, err := s.repo.GetIncidentByID(ctx, req.TenantID, req.IncidentID)
	if err != nil {
		return nil, err
	}

	if incident.Status == IncidentResolved || incident.Status == IncidentFalseAlarm {
		return nil, ErrIncidentAlreadyClosed
	}

	now := time.Now().UTC()
	responder := &SOSDispatchResponder{
		ID:           uuid.New().String(),
		TenantID:     req.TenantID,
		IncidentID:   req.IncidentID,
		ResponderID:  req.ResponderID,
		Role:         req.Role,
		Status:       ResponderAssigned,
		DispatchedAt: now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateResponder(ctx, responder); err != nil {
		return nil, err
	}

	// Move incident status to DISPATCHED if still TRIGGERED or ACKNOWLEDGED
	if incident.Status == IncidentTriggered || incident.Status == IncidentAcknowledged {
		incident.Status = IncidentDispatched
		incident.UpdatedAt = now
		_ = s.repo.UpdateIncident(ctx, incident)
	}

	s.logAudit(req.TenantID, req.ActorID, "STAFF", "sos:alert:dispatched", "sos_dispatch_responder", responder.ID, audit.StatusSuccess, map[string]interface{}{
		"incidentId":  req.IncidentID,
		"responderId": req.ResponderID,
		"role":        req.Role,
	})

	return responder, nil
}

type UpdateResponderStatusRequest struct {
	TenantID    string             `json:"tenantId"`
	IncidentID  string             `json:"incidentId"`
	ResponderID string             `json:"responderId"`
	Status      SOSResponderStatus `json:"status"`
}

func (s *service) UpdateResponderStatus(ctx context.Context, req UpdateResponderStatusRequest) (*SOSDispatchResponder, error) {
	responder, err := s.repo.GetResponder(ctx, req.TenantID, req.IncidentID, req.ResponderID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	responder.Status = req.Status
	responder.UpdatedAt = now
	if req.Status == ResponderOnScene && responder.ArrivedAt == nil {
		responder.ArrivedAt = &now
	}

	if err := s.repo.UpdateResponder(ctx, responder); err != nil {
		return nil, err
	}

	// If responder arrived on scene, update incident to ON_SCENE
	if req.Status == ResponderOnScene {
		if incident, err := s.repo.GetIncidentByID(ctx, req.TenantID, req.IncidentID); err == nil {
			if incident.Status == IncidentDispatched || incident.Status == IncidentAcknowledged {
				incident.Status = IncidentOnScene
				incident.UpdatedAt = now
				_ = s.repo.UpdateIncident(ctx, incident)
			}
		}
	}

	return responder, nil
}

type ResolveIncidentRequest struct {
	TenantID        string            `json:"tenantId"`
	IncidentID      string            `json:"incidentId"`
	ResolvedByID    string            `json:"resolvedById"`
	Status          SOSIncidentStatus `json:"status"` // RESOLVED or FALSE_ALARM
	ResolutionNotes string            `json:"resolutionNotes"`
}

func (s *service) ResolveIncident(ctx context.Context, req ResolveIncidentRequest) (*SOSIncident, error) {
	incident, err := s.repo.GetIncidentByID(ctx, req.TenantID, req.IncidentID)
	if err != nil {
		return nil, err
	}

	targetStatus := req.Status
	if targetStatus != IncidentResolved && targetStatus != IncidentFalseAlarm {
		targetStatus = IncidentResolved
	}

	if err := ValidateIncidentTransition(incident.Status, targetStatus); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	incident.Status = targetStatus
	incident.ResolvedAt = &now
	incident.ResolvedByID = &req.ResolvedByID
	incident.ResolutionNotes = &req.ResolutionNotes
	incident.UpdatedAt = now

	if err := s.repo.UpdateIncident(ctx, incident); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, req.ResolvedByID, "STAFF", "sos:alert:resolved", "sos_incident", incident.ID, audit.StatusSuccess, map[string]interface{}{
		"status":          targetStatus,
		"resolutionNotes": req.ResolutionNotes,
	})

	return incident, nil
}

func (s *service) GetIncident(ctx context.Context, tenantID, id string) (*SOSIncident, error) {
	return s.repo.GetIncidentByID(ctx, tenantID, id)
}

func (s *service) ListIncidents(ctx context.Context, filter SOSIncidentFilter) ([]SOSIncident, error) {
	return s.repo.ListIncidents(ctx, filter)
}

func (s *service) ListResponders(ctx context.Context, tenantID, incidentID string) ([]SOSDispatchResponder, error) {
	return s.repo.ListResponders(ctx, tenantID, incidentID)
}
