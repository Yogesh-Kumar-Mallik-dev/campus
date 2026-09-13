/**
 * BLOCK_EVENTS_SERVICE_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   Core business logic service and audit dispatching for campus events, venue bookings, and single-use QR gate check-ins.
 */

package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"campus/backend/audit"
)

type Service interface {
	CreateEvent(ctx context.Context, req CreateEventRequest) (*CampusEvent, error)
	GetEvent(ctx context.Context, tenantID, id string) (*CampusEvent, error)
	ListEvents(ctx context.Context, tenantID string, category *EventCategory, status *EventStatus) ([]*CampusEvent, error)

	BookTicket(ctx context.Context, req BookTicketRequest) (*EventTicket, error)
	GetTicket(ctx context.Context, tenantID, id string) (*EventTicket, error)
	CheckInTicket(ctx context.Context, req CheckInTicketRequest) (*EventTicket, error)
	ListEventTickets(ctx context.Context, tenantID, eventID string) ([]*EventTicket, error)
	ListStudentTickets(ctx context.Context, tenantID, studentID string) ([]*EventTicket, error)
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

func (s *service) emitAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
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

type CreateEventRequest struct {
	TenantID             string        `json:"tenantId"`
	Title                string        `json:"title"`
	Description          string        `json:"description"`
	Category             EventCategory `json:"category"`
	VenueName            string        `json:"venueName"`
	VenueCapacity        int           `json:"venueCapacity"`
	StartTime            time.Time     `json:"startTime"`
	EndTime              time.Time     `json:"endTime"`
	RegistrationDeadline time.Time     `json:"registrationDeadline"`
	OrganizerID          string        `json:"organizerId"`
	IsTicketed           bool          `json:"isTicketed"`
	TicketPrice          float64       `json:"ticketPrice"`
	MaxTickets           int           `json:"maxTickets"`
}

func (s *service) CreateEvent(ctx context.Context, req CreateEventRequest) (*CampusEvent, error) {
	if req.EndTime.Before(req.StartTime) || req.EndTime.Equal(req.StartTime) {
		return nil, ErrInvalidEventTimeRange
	}
	if req.RegistrationDeadline.After(req.EndTime) {
		return nil, ErrInvalidEventTimeRange
	}

	// 1. Check Venue Overlap Conflict
	hasConflict, err := s.repo.CheckVenueConflict(ctx, req.TenantID, req.VenueName, req.StartTime, req.EndTime, nil)
	if err != nil {
		return nil, err
	}
	if hasConflict {
		return nil, ErrVenueConflict
	}

	cap := req.VenueCapacity
	if cap <= 0 {
		cap = 100
	}
	maxTix := req.MaxTickets
	if maxTix <= 0 || maxTix > cap {
		maxTix = cap
	}
	cat := req.Category
	if cat == "" {
		cat = EventCategoryTechnical
	}

	now := time.Now().UTC()
	event := &CampusEvent{
		ID:                   fmt.Sprintf("evnt_%d", now.UnixNano()),
		TenantID:             req.TenantID,
		Title:                req.Title,
		Description:          req.Description,
		Category:             cat,
		VenueName:            req.VenueName,
		VenueCapacity:        cap,
		StartTime:            req.StartTime,
		EndTime:              req.EndTime,
		RegistrationDeadline: req.RegistrationDeadline,
		OrganizerID:          req.OrganizerID,
		Status:               EventStatusPublished,
		IsTicketed:           req.IsTicketed,
		TicketPrice:          req.TicketPrice,
		MaxTickets:           maxTix,
		TicketsSold:          0,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, err
	}

	// 2. Reserve Venue Booking
	booking := &EventVenueBooking{
		ID:              fmt.Sprintf("vnb_%d", now.UnixNano()),
		TenantID:        req.TenantID,
		VenueName:       req.VenueName,
		EventID:         event.ID,
		BookedStartTime: req.StartTime,
		BookedEndTime:   req.EndTime,
		BookedByID:      req.OrganizerID,
		Status:          VenueBookingStatusConfirmed,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	_ = s.repo.CreateVenueBooking(ctx, booking)

	s.emitAudit(event.TenantID, req.OrganizerID, "USER", "events:event:created", "campus_event", event.ID, audit.StatusSuccess, map[string]interface{}{
		"title": event.Title, "venue": event.VenueName, "startTime": event.StartTime,
	})

	return event, nil
}

func (s *service) GetEvent(ctx context.Context, tenantID, id string) (*CampusEvent, error) {
	return s.repo.GetEventByID(ctx, tenantID, id)
}

func (s *service) ListEvents(ctx context.Context, tenantID string, category *EventCategory, status *EventStatus) ([]*CampusEvent, error) {
	return s.repo.ListEvents(ctx, tenantID, category, status)
}

type BookTicketRequest struct {
	TenantID  string `json:"tenantId"`
	EventID   string `json:"eventId"`
	StudentID string `json:"studentId"`
}

func (s *service) BookTicket(ctx context.Context, req BookTicketRequest) (*EventTicket, error) {
	event, err := s.repo.GetEventByID(ctx, req.TenantID, req.EventID)
	if err != nil {
		return nil, err
	}

	if event.Status != EventStatusPublished {
		return nil, ErrEventNotPublished
	}

	now := time.Now().UTC()
	if now.After(event.RegistrationDeadline) {
		return nil, ErrRegistrationClosed
	}

	if event.TicketsSold >= event.MaxTickets || event.TicketsSold >= event.VenueCapacity {
		return nil, ErrEventCapacityExceeded
	}

	if existing, _ := s.repo.GetTicketByStudent(ctx, req.TenantID, req.EventID, req.StudentID); existing != nil {
		return nil, ErrDuplicateTicketBooking
	}

	ticketCode := GenerateTicketCode(req.TenantID, req.EventID, req.StudentID, now)
	ticket := &EventTicket{
		ID:         fmt.Sprintf("tck_%d", now.UnixNano()),
		TenantID:   req.TenantID,
		EventID:    req.EventID,
		StudentID:  req.StudentID,
		TicketCode: ticketCode,
		Status:     TicketStatusConfirmed,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateTicket(ctx, ticket); err != nil {
		return nil, err
	}

	_ = s.repo.IncrementTicketsSold(ctx, req.TenantID, req.EventID)

	s.emitAudit(ticket.TenantID, req.StudentID, "STUDENT", "events:ticket:issued", "event_ticket", ticket.ID, audit.StatusSuccess, map[string]interface{}{
		"eventId": ticket.EventID, "ticketCode": ticket.TicketCode,
	})

	return ticket, nil
}

func (s *service) GetTicket(ctx context.Context, tenantID, id string) (*EventTicket, error) {
	return s.repo.GetTicketByID(ctx, tenantID, id)
}

type CheckInTicketRequest struct {
	TenantID   string `json:"tenantId"`
	TicketCode string `json:"ticketCode"`
	StaffID    string `json:"staffId"`
}

func (s *service) CheckInTicket(ctx context.Context, req CheckInTicketRequest) (*EventTicket, error) {
	ticket, err := s.repo.GetTicketByCode(ctx, req.TenantID, req.TicketCode)
	if err != nil {
		return nil, err
	}

	if err := ValidateCheckIn(ticket.Status); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := s.repo.CheckInTicket(ctx, req.TenantID, ticket.ID, req.StaffID, now); err != nil {
		return nil, err
	}

	ticket.Status = TicketStatusCheckedIn
	ticket.CheckedInAt = &now
	ticket.CheckedInByID = &req.StaffID
	ticket.UpdatedAt = now

	s.emitAudit(ticket.TenantID, req.StaffID, "STAFF", "events:ticket:checked_in", "event_ticket", ticket.ID, audit.StatusSuccess, map[string]interface{}{
		"eventId": ticket.EventID, "studentId": ticket.StudentID, "ticketCode": ticket.TicketCode,
	})

	return ticket, nil
}

func (s *service) ListEventTickets(ctx context.Context, tenantID, eventID string) ([]*EventTicket, error) {
	return s.repo.ListTicketsByEvent(ctx, tenantID, eventID)
}

func (s *service) ListStudentTickets(ctx context.Context, tenantID, studentID string) ([]*EventTicket, error) {
	return s.repo.ListTicketsByStudent(ctx, tenantID, studentID)
}
