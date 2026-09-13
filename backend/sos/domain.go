/**
 * BLOCK_SOS_DOMAIN_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   Domain entities, value objects, coordinate validation, and incident state machine.
 */

package sos

import (
	"fmt"
	"time"
)

type SOSEmergencyType string

const (
	EmergencyMedical           SOSEmergencyType = "MEDICAL"
	EmergencyFire              SOSEmergencyType = "FIRE"
	EmergencySecurityThreat    SOSEmergencyType = "SECURITY_THREAT"
	EmergencyNaturalHazard     SOSEmergencyType = "NATURAL_HAZARD"
	EmergencyHarassmentRagging SOSEmergencyType = "HARASSMENT_RAGGING"
	EmergencyOther             SOSEmergencyType = "OTHER"
)

type SOSIncidentStatus string

const (
	IncidentTriggered    SOSIncidentStatus = "TRIGGERED"
	IncidentAcknowledged SOSIncidentStatus = "ACKNOWLEDGED"
	IncidentDispatched   SOSIncidentStatus = "DISPATCHED"
	IncidentOnScene      SOSIncidentStatus = "ON_SCENE"
	IncidentResolved     SOSIncidentStatus = "RESOLVED"
	IncidentFalseAlarm   SOSIncidentStatus = "FALSE_ALARM"
)

type SOSResponderRole string

const (
	RoleCampusSecurity    SOSResponderRole = "CAMPUS_SECURITY"
	RoleParamedic         SOSResponderRole = "PARAMEDIC"
	RoleHostelWarden      SOSResponderRole = "HOSTEL_WARDEN"
	RoleFireSafetyOfficer SOSResponderRole = "FIRE_SAFETY_OFFICER"
	RolePoliceLiaison     SOSResponderRole = "POLICE_LIAISON"
)

type SOSResponderStatus string

const (
	ResponderAssigned  SOSResponderStatus = "ASSIGNED"
	ResponderEnRoute   SOSResponderStatus = "EN_ROUTE"
	ResponderOnScene   SOSResponderStatus = "ON_SCENE"
	ResponderCompleted SOSResponderStatus = "COMPLETED"
)

type SOSTelemetryChannel string

const (
	ChannelSMSGateway        SOSTelemetryChannel = "SMS_GATEWAY"
	ChannelPushNotification SOSTelemetryChannel = "PUSH_NOTIFICATION"
	ChannelCampusPA          SOSTelemetryChannel = "CAMPUS_PA"
	ChannelPublicSiren       SOSTelemetryChannel = "PUBLIC_SIREN"
)

type SOSIncident struct {
	ID                  string            `json:"id"`
	TenantID            string            `json:"tenantId"`
	AlertNumber         string            `json:"alertNumber"`
	UserID              string            `json:"userId"`
	EmergencyType       SOSEmergencyType  `json:"emergencyType"`
	Latitude            float64           `json:"latitude"`
	Longitude           float64           `json:"longitude"`
	LocationDescription string            `json:"locationDescription"`
	Status              SOSIncidentStatus `json:"status"`
	TriggeredAt         time.Time         `json:"triggeredAt"`
	AcknowledgedAt      *time.Time        `json:"acknowledgedAt,omitempty"`
	ResolvedAt          *time.Time        `json:"resolvedAt,omitempty"`
	ResolvedByID        *string           `json:"resolvedById,omitempty"`
	ResolutionNotes     *string           `json:"resolutionNotes,omitempty"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt"`
}

type SOSDispatchResponder struct {
	ID           string             `json:"id"`
	TenantID     string             `json:"tenantId"`
	IncidentID   string             `json:"incidentId"`
	ResponderID  string             `json:"responderId"`
	Role         SOSResponderRole   `json:"role"`
	Status       SOSResponderStatus `json:"status"`
	DispatchedAt time.Time          `json:"dispatchedAt"`
	ArrivedAt    *time.Time         `json:"arrivedAt,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}

type SOSTelemetryBroadcast struct {
	ID             string              `json:"id"`
	TenantID       string              `json:"tenantId"`
	IncidentID     string              `json:"incidentId"`
	Channel        SOSTelemetryChannel `json:"channel"`
	Payload        string              `json:"payload"`
	DeliveredCount int                 `json:"deliveredCount"`
	BroadcastAt    time.Time           `json:"broadcastAt"`
	CreatedAt      time.Time           `json:"createdAt"`
}

type SOSIncidentFilter struct {
	TenantID      string
	Status        *SOSIncidentStatus
	EmergencyType *SOSEmergencyType
	UserID        *string
}

// FormatAlertNumber generates deterministic SOS alert code.
func FormatAlertNumber(year, seq int) string {
	return fmt.Sprintf("SOS-%d-%05d", year, seq)
}

// ValidateCoordinates verifies geographical coordinate limits.
func ValidateCoordinates(lat, lng float64) error {
	if lat < -90.0 || lat > 90.0 || lng < -180.0 || lng > 180.0 {
		return ErrInvalidCoordinates
	}
	return nil
}

// ValidateIncidentTransition validates legal state transitions for SOS alerts.
func ValidateIncidentTransition(current, next SOSIncidentStatus) error {
	if current == next {
		return nil
	}

	if current == IncidentResolved || current == IncidentFalseAlarm {
		return fmt.Errorf("%w: incident is in terminal state %s", ErrIncidentAlreadyClosed, current)
	}

	switch current {
	case IncidentTriggered:
		if next == IncidentAcknowledged || next == IncidentDispatched || next == IncidentResolved || next == IncidentFalseAlarm {
			return nil
		}
	case IncidentAcknowledged:
		if next == IncidentDispatched || next == IncidentOnScene || next == IncidentResolved || next == IncidentFalseAlarm {
			return nil
		}
	case IncidentDispatched:
		if next == IncidentOnScene || next == IncidentResolved || next == IncidentFalseAlarm {
			return nil
		}
	case IncidentOnScene:
		if next == IncidentResolved || next == IncidentFalseAlarm {
			return nil
		}
	}

	return fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidIncidentTransition, current, next)
}
