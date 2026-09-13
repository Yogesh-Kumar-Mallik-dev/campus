/**
 * BLOCK_SOS_MOCK_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   Concurrent in-memory mock repository test double for SOS subsystem.
 */

package sos

import (
	"context"
	"sort"
	"sync"
)

type MockRepository struct {
	mu         sync.RWMutex
	incidents  map[string]SOSIncident
	responders []SOSDispatchResponder
	broadcasts []SOSTelemetryBroadcast
	sequences  map[string]int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		incidents:  make(map[string]SOSIncident),
		responders: make([]SOSDispatchResponder, 0),
		broadcasts: make([]SOSTelemetryBroadcast, 0),
		sequences:  make(map[string]int),
	}
}

func (m *MockRepository) CreateIncident(ctx context.Context, incident *SOSIncident) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.incidents[incident.ID] = *incident
	return nil
}

func (m *MockRepository) GetIncidentByID(ctx context.Context, tenantID, id string) (*SOSIncident, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inc, ok := m.incidents[id]
	if !ok || inc.TenantID != tenantID {
		return nil, ErrIncidentNotFound
	}
	return &inc, nil
}

func (m *MockRepository) GetIncidentByAlertNumber(ctx context.Context, tenantID, alertNumber string) (*SOSIncident, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, inc := range m.incidents {
		if inc.TenantID == tenantID && inc.AlertNumber == alertNumber {
			return &inc, nil
		}
	}
	return nil, ErrIncidentNotFound
}

func (m *MockRepository) ListIncidents(ctx context.Context, filter SOSIncidentFilter) ([]SOSIncident, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]SOSIncident, 0)
	for _, inc := range m.incidents {
		if inc.TenantID != filter.TenantID {
			continue
		}
		if filter.Status != nil && inc.Status != *filter.Status {
			continue
		}
		if filter.EmergencyType != nil && inc.EmergencyType != *filter.EmergencyType {
			continue
		}
		if filter.UserID != nil && inc.UserID != *filter.UserID {
			continue
		}
		res = append(res, inc)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].TriggeredAt.After(res[j].TriggeredAt)
	})

	return res, nil
}

func (m *MockRepository) UpdateIncident(ctx context.Context, incident *SOSIncident) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.incidents[incident.ID]; !ok {
		return ErrIncidentNotFound
	}
	m.incidents[incident.ID] = *incident
	return nil
}

func (m *MockRepository) GetNextAlertSequence(ctx context.Context, tenantID string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sequences[tenantID]++
	return m.sequences[tenantID], nil
}

func (m *MockRepository) CreateResponder(ctx context.Context, responder *SOSDispatchResponder) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, r := range m.responders {
		if r.TenantID == responder.TenantID && r.IncidentID == responder.IncidentID && r.ResponderID == responder.ResponderID {
			return ErrResponderAlreadyDispatched
		}
	}

	m.responders = append(m.responders, *responder)
	return nil
}

func (m *MockRepository) GetResponder(ctx context.Context, tenantID, incidentID, responderID string) (*SOSDispatchResponder, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, r := range m.responders {
		if r.TenantID == tenantID && r.IncidentID == incidentID && r.ResponderID == responderID {
			return &r, nil
		}
	}
	return nil, ErrResponderNotFound
}

func (m *MockRepository) ListResponders(ctx context.Context, tenantID, incidentID string) ([]SOSDispatchResponder, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]SOSDispatchResponder, 0)
	for _, r := range m.responders {
		if r.TenantID == tenantID && r.IncidentID == incidentID {
			res = append(res, r)
		}
	}
	return res, nil
}

func (m *MockRepository) UpdateResponder(ctx context.Context, responder *SOSDispatchResponder) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	found := false
	for i, r := range m.responders {
		if r.ID == responder.ID {
			m.responders[i] = *responder
			found = true
			break
		}
	}
	if !found {
		return ErrResponderNotFound
	}
	return nil
}

func (m *MockRepository) CreateBroadcast(ctx context.Context, broadcast *SOSTelemetryBroadcast) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.broadcasts = append(m.broadcasts, *broadcast)
	return nil
}

func (m *MockRepository) ListBroadcasts(ctx context.Context, tenantID, incidentID string) ([]SOSTelemetryBroadcast, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]SOSTelemetryBroadcast, 0)
	for _, b := range m.broadcasts {
		if b.TenantID == tenantID && b.IncidentID == incidentID {
			res = append(res, b)
		}
	}
	return res, nil
}
