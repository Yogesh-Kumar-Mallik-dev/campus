/**
 * BLOCK_EVENTS_REPOSITORY_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   Repository interface contracts for campus events, venue bookings, and gate ticketing.
 */

package events

import (
	"context"
	"time"
)

type Repository interface {
	// Events
	CreateEvent(ctx context.Context, event *CampusEvent) error
	GetEventByID(ctx context.Context, tenantID, id string) (*CampusEvent, error)
	ListEvents(ctx context.Context, tenantID string, category *EventCategory, status *EventStatus) ([]*CampusEvent, error)
	UpdateEventStatus(ctx context.Context, tenantID, id string, status EventStatus) error
	IncrementTicketsSold(ctx context.Context, tenantID, id string) error

	// Venue Bookings & Conflicts
	CreateVenueBooking(ctx context.Context, booking *EventVenueBooking) error
	ListVenueBookings(ctx context.Context, tenantID, venueName string) ([]*EventVenueBooking, error)
	CheckVenueConflict(ctx context.Context, tenantID, venueName string, startTime, endTime time.Time, excludeEventID *string) (bool, error)

	// Tickets
	CreateTicket(ctx context.Context, ticket *EventTicket) error
	GetTicketByID(ctx context.Context, tenantID, id string) (*EventTicket, error)
	GetTicketByCode(ctx context.Context, tenantID, code string) (*EventTicket, error)
	GetTicketByStudent(ctx context.Context, tenantID, eventID, studentID string) (*EventTicket, error)
	ListTicketsByEvent(ctx context.Context, tenantID, eventID string) ([]*EventTicket, error)
	ListTicketsByStudent(ctx context.Context, tenantID, studentID string) ([]*EventTicket, error)
	CheckInTicket(ctx context.Context, tenantID, id string, staffID string, checkedInAt time.Time) error
}
