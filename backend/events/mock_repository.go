/**
 * BLOCK_EVENTS_MOCK_REPO_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   Thread-safe in-memory mock repository test doubles for events, venue bookings, and ticketing.
 */

package events

import (
	"context"
	"sync"
	"time"
)

type MockRepository struct {
	mu            sync.RWMutex
	events        map[string]*CampusEvent
	venueBookings map[string]*EventVenueBooking
	tickets       map[string]*EventTicket
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		events:        make(map[string]*CampusEvent),
		venueBookings: make(map[string]*EventVenueBooking),
		tickets:       make(map[string]*EventTicket),
	}
}

// --- Events ---

func (m *MockRepository) CreateEvent(_ context.Context, event *CampusEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events[event.ID] = event
	return nil
}

func (m *MockRepository) GetEventByID(_ context.Context, tenantID, id string) (*CampusEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ev, ok := m.events[id]
	if !ok || ev.TenantID != tenantID {
		return nil, ErrEventNotFound
	}
	return ev, nil
}

func (m *MockRepository) ListEvents(_ context.Context, tenantID string, category *EventCategory, status *EventStatus) ([]*CampusEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*CampusEvent
	for _, ev := range m.events {
		if ev.TenantID != tenantID {
			continue
		}
		if category != nil && ev.Category != *category {
			continue
		}
		if status != nil && ev.Status != *status {
			continue
		}
		res = append(res, ev)
	}
	return res, nil
}

func (m *MockRepository) UpdateEventStatus(_ context.Context, tenantID, id string, status EventStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ev, ok := m.events[id]
	if !ok || ev.TenantID != tenantID {
		return ErrEventNotFound
	}
	ev.Status = status
	ev.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) IncrementTicketsSold(_ context.Context, tenantID, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ev, ok := m.events[id]
	if !ok || ev.TenantID != tenantID {
		return ErrEventNotFound
	}
	ev.TicketsSold++
	ev.UpdatedAt = time.Now().UTC()
	return nil
}

// --- Venue Bookings ---

func (m *MockRepository) CreateVenueBooking(_ context.Context, booking *EventVenueBooking) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.venueBookings[booking.ID] = booking
	return nil
}

func (m *MockRepository) ListVenueBookings(_ context.Context, tenantID, venueName string) ([]*EventVenueBooking, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*EventVenueBooking
	for _, vb := range m.venueBookings {
		if vb.TenantID == tenantID && vb.VenueName == venueName {
			res = append(res, vb)
		}
	}
	return res, nil
}

func (m *MockRepository) CheckVenueConflict(_ context.Context, tenantID, venueName string, startTime, endTime time.Time, excludeEventID *string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, vb := range m.venueBookings {
		if vb.TenantID != tenantID || vb.VenueName != venueName || vb.Status != VenueBookingStatusConfirmed {
			continue
		}
		if excludeEventID != nil && vb.EventID == *excludeEventID {
			continue
		}
		if CheckVenueOverlap(vb.BookedStartTime, vb.BookedEndTime, startTime, endTime) {
			return true, nil
		}
	}
	return false, nil
}

// --- Tickets ---

func (m *MockRepository) CreateTicket(_ context.Context, ticket *EventTicket) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tickets[ticket.ID] = ticket
	return nil
}

func (m *MockRepository) GetTicketByID(_ context.Context, tenantID, id string) (*EventTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tck, ok := m.tickets[id]
	if !ok || tck.TenantID != tenantID {
		return nil, ErrTicketNotFound
	}
	return tck, nil
}

func (m *MockRepository) GetTicketByCode(_ context.Context, tenantID, code string) (*EventTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, t := range m.tickets {
		if t.TenantID == tenantID && t.TicketCode == code {
			return t, nil
		}
	}
	return nil, ErrTicketNotFound
}

func (m *MockRepository) GetTicketByStudent(_ context.Context, tenantID, eventID, studentID string) (*EventTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, t := range m.tickets {
		if t.TenantID == tenantID && t.EventID == eventID && t.StudentID == studentID && t.Status != TicketStatusCancelled {
			return t, nil
		}
	}
	return nil, ErrTicketNotFound
}

func (m *MockRepository) ListTicketsByEvent(_ context.Context, tenantID, eventID string) ([]*EventTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*EventTicket
	for _, t := range m.tickets {
		if t.TenantID == tenantID && t.EventID == eventID {
			res = append(res, t)
		}
	}
	return res, nil
}

func (m *MockRepository) ListTicketsByStudent(_ context.Context, tenantID, studentID string) ([]*EventTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*EventTicket
	for _, t := range m.tickets {
		if t.TenantID == tenantID && t.StudentID == studentID {
			res = append(res, t)
		}
	}
	return res, nil
}

func (m *MockRepository) CheckInTicket(_ context.Context, tenantID, id string, staffID string, checkedInAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	tck, ok := m.tickets[id]
	if !ok || tck.TenantID != tenantID {
		return ErrTicketNotFound
	}
	tck.Status = TicketStatusCheckedIn
	tck.CheckedInAt = &checkedInAt
	tck.CheckedInByID = &staffID
	tck.UpdatedAt = checkedInAt
	return nil
}
