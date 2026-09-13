/**
 * BLOCK_SOS_REPO_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   Repository interface contracts for SOS incident records, responder dispatches, and broadcasts.
 */

package sos

import "context"

type Repository interface {
	CreateIncident(ctx context.Context, incident *SOSIncident) error
	GetIncidentByID(ctx context.Context, tenantID, id string) (*SOSIncident, error)
	GetIncidentByAlertNumber(ctx context.Context, tenantID, alertNumber string) (*SOSIncident, error)
	ListIncidents(ctx context.Context, filter SOSIncidentFilter) ([]SOSIncident, error)
	UpdateIncident(ctx context.Context, incident *SOSIncident) error
	GetNextAlertSequence(ctx context.Context, tenantID string) (int, error)

	CreateResponder(ctx context.Context, responder *SOSDispatchResponder) error
	GetResponder(ctx context.Context, tenantID, incidentID, responderID string) (*SOSDispatchResponder, error)
	ListResponders(ctx context.Context, tenantID, incidentID string) ([]SOSDispatchResponder, error)
	UpdateResponder(ctx context.Context, responder *SOSDispatchResponder) error

	CreateBroadcast(ctx context.Context, broadcast *SOSTelemetryBroadcast) error
	ListBroadcasts(ctx context.Context, tenantID, incidentID string) ([]SOSTelemetryBroadcast, error)
}
