/**
 * BLOCK_EVENTS_TEST_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   Unit test suite covering event creation, venue conflict detection, capacity limits, and QR gate check-in guards.
 */

package events_test

import (
	"context"
	"testing"
	"time"

	"campus/backend/events"
)

func setupEventsService() events.Service {
	repo := events.NewMockRepository()
	return events.NewService(repo, nil)
}

func TestEvents_EventCreationAndVenueConflictGuard(t *testing.T) {
	ctx := context.Background()
	svc := setupEventsService()

	tenantID := "tenant-alpha"
	venue := "Main Auditorium"
	startTime := time.Now().UTC().Add(48 * time.Hour)
	endTime := startTime.Add(3 * time.Hour)
	regDeadline := startTime.Add(-2 * time.Hour)

	// 1. Create First Event in Main Auditorium
	event1, err := svc.CreateEvent(ctx, events.CreateEventRequest{
		TenantID:             tenantID,
		Title:                "Annual Tech Symposium 2026",
		Description:          "Keynotes, hackathon showcases, and robotics competitions.",
		Category:             events.EventCategoryTechnical,
		VenueName:            venue,
		VenueCapacity:        500,
		StartTime:            startTime,
		EndTime:              endTime,
		RegistrationDeadline: regDeadline,
		OrganizerID:          "staff-prof-dean",
		MaxTickets:           500,
	})
	if err != nil {
		t.Fatalf("unexpected error creating event 1: %v", err)
	}
	if event1.Status != events.EventStatusPublished {
		t.Fatalf("expected status PUBLISHED, got %s", event1.Status)
	}

	// 2. Conflicting Booking in Same Venue
	overlappingStart := startTime.Add(1 * time.Hour)
	overlappingEnd := endTime.Add(1 * time.Hour)
	_, err = svc.CreateEvent(ctx, events.CreateEventRequest{
		TenantID:             tenantID,
		Title:                "Conflicting Cultural Play",
		Category:             events.EventCategoryCultural,
		VenueName:            venue,
		VenueCapacity:        300,
		StartTime:            overlappingStart,
		EndTime:              overlappingEnd,
		RegistrationDeadline: overlappingStart.Add(-1 * time.Hour),
		OrganizerID:          "staff-cultural-head",
	})
	if err != events.ErrVenueConflict {
		t.Fatalf("expected ErrVenueConflict, got %v", err)
	}

	// 3. Non-conflicting Event in Different Venue
	event3, err := svc.CreateEvent(ctx, events.CreateEventRequest{
		TenantID:             tenantID,
		Title:                "Guest Lecture on Quantum Computing",
		Category:             events.EventCategoryGuestLecture,
		VenueName:            "Seminar Hall B",
		VenueCapacity:        150,
		StartTime:            startTime,
		EndTime:              endTime,
		RegistrationDeadline: regDeadline,
		OrganizerID:          "staff-prof-physics",
	})
	if err != nil {
		t.Fatalf("unexpected error creating event 3: %v", err)
	}
	if event3.VenueName != "Seminar Hall B" {
		t.Fatalf("expected Seminar Hall B, got %s", event3.VenueName)
	}
}

func TestEvents_TicketBookingAndCapacityGuard(t *testing.T) {
	ctx := context.Background()
	svc := setupEventsService()

	tenantID := "tenant-alpha"
	startTime := time.Now().UTC().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	regDeadline := startTime.Add(-1 * time.Hour)

	event, _ := svc.CreateEvent(ctx, events.CreateEventRequest{
		TenantID:             tenantID,
		Title:                "AI Hands-on Workshop",
		VenueName:            "Computer Lab 4",
		VenueCapacity:        2, // Max 2 seats
		StartTime:            startTime,
		EndTime:              endTime,
		RegistrationDeadline: regDeadline,
		OrganizerID:          "staff-ai-club",
		MaxTickets:           2,
	})

	// 1. Student 1 Books
	tck1, err := svc.BookTicket(ctx, events.BookTicketRequest{
		TenantID:  tenantID,
		EventID:   event.ID,
		StudentID: "stu-001",
	})
	if err != nil {
		t.Fatalf("failed to book ticket 1: %v", err)
	}
	if tck1.Status != events.TicketStatusConfirmed {
		t.Fatalf("expected status CONFIRMED, got %s", tck1.Status)
	}

	// 2. Duplicate Booking Guard
	_, err = svc.BookTicket(ctx, events.BookTicketRequest{
		TenantID:  tenantID,
		EventID:   event.ID,
		StudentID: "stu-001",
	})
	if err != events.ErrDuplicateTicketBooking {
		t.Fatalf("expected ErrDuplicateTicketBooking, got %v", err)
	}

	// 3. Student 2 Books (Fills remaining seat)
	tck2, err := svc.BookTicket(ctx, events.BookTicketRequest{
		TenantID:  tenantID,
		EventID:   event.ID,
		StudentID: "stu-002",
	})
	if err != nil {
		t.Fatalf("failed to book ticket 2: %v", err)
	}
	_ = tck2

	// 4. Student 3 Books -> Capacity Exceeded
	_, err = svc.BookTicket(ctx, events.BookTicketRequest{
		TenantID:  tenantID,
		EventID:   event.ID,
		StudentID: "stu-003",
	})
	if err != events.ErrEventCapacityExceeded {
		t.Fatalf("expected ErrEventCapacityExceeded, got %v", err)
	}
}

func TestEvents_GateCheckInLifecycleAndDoublePunchGuard(t *testing.T) {
	ctx := context.Background()
	svc := setupEventsService()

	tenantID := "tenant-alpha"
	startTime := time.Now().UTC().Add(12 * time.Hour)
	endTime := startTime.Add(4 * time.Hour)
	regDeadline := startTime.Add(-2 * time.Hour)

	event, _ := svc.CreateEvent(ctx, events.CreateEventRequest{
		TenantID:             tenantID,
		Title:                "Sports Gala Finals",
		VenueName:            "University Stadium",
		VenueCapacity:        1000,
		StartTime:            startTime,
		EndTime:              endTime,
		RegistrationDeadline: regDeadline,
		OrganizerID:          "staff-sports-officer",
		MaxTickets:           1000,
	})

	ticket, _ := svc.BookTicket(ctx, events.BookTicketRequest{
		TenantID:  tenantID,
		EventID:   event.ID,
		StudentID: "stu-004",
	})

	// 1. Initial Gate Admission Check-In
	staffID := "staff-security-gate1"
	checkedInTck, err := svc.CheckInTicket(ctx, events.CheckInTicketRequest{
		TenantID:   tenantID,
		TicketCode: ticket.TicketCode,
		StaffID:    staffID,
	})
	if err != nil {
		t.Fatalf("failed to check in ticket: %v", err)
	}
	if checkedInTck.Status != events.TicketStatusCheckedIn {
		t.Fatalf("expected status CHECKED_IN, got %s", checkedInTck.Status)
	}
	if checkedInTck.CheckedInAt == nil || *checkedInTck.CheckedInByID != staffID {
		t.Fatalf("check-in metadata mismatch: %+v", checkedInTck)
	}

	// 2. Double-Punch Replay Attack Rejection
	_, err = svc.CheckInTicket(ctx, events.CheckInTicketRequest{
		TenantID:   tenantID,
		TicketCode: ticket.TicketCode,
		StaffID:    "staff-security-gate2",
	})
	if err != events.ErrTicketAlreadyCheckedIn {
		t.Fatalf("expected ErrTicketAlreadyCheckedIn, got %v", err)
	}
}

func TestEvents_RegistrationDeadlineGuard(t *testing.T) {
	ctx := context.Background()
	svc := setupEventsService()

	tenantID := "tenant-alpha"
	startTime := time.Now().UTC().Add(2 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	pastRegDeadline := time.Now().UTC().Add(-1 * time.Hour) // Closed 1 hour ago

	event, _ := svc.CreateEvent(ctx, events.CreateEventRequest{
		TenantID:             tenantID,
		Title:                "Flash Mob Rehearsal",
		VenueName:            "Open Amphitheatre",
		VenueCapacity:        100,
		StartTime:            startTime,
		EndTime:              endTime,
		RegistrationDeadline: pastRegDeadline,
		OrganizerID:          "staff-cultural",
	})

	_, err := svc.BookTicket(ctx, events.BookTicketRequest{
		TenantID:  tenantID,
		EventID:   event.ID,
		StudentID: "stu-005",
	})
	if err != events.ErrRegistrationClosed {
		t.Fatalf("expected ErrRegistrationClosed, got %v", err)
	}
}

func TestEvents_TimeRangeGuard(t *testing.T) {
	ctx := context.Background()
	svc := setupEventsService()

	tenantID := "tenant-alpha"
	startTime := time.Now().UTC().Add(24 * time.Hour)
	invalidEndTime := startTime.Add(-2 * time.Hour) // End before start

	_, err := svc.CreateEvent(ctx, events.CreateEventRequest{
		TenantID:             tenantID,
		Title:                "Invalid Time Event",
		VenueName:            "Main Hall",
		StartTime:            startTime,
		EndTime:              invalidEndTime,
		RegistrationDeadline: startTime.Add(-1 * time.Hour),
		OrganizerID:          "staff-01",
	})
	if err != events.ErrInvalidEventTimeRange {
		t.Fatalf("expected ErrInvalidEventTimeRange, got %v", err)
	}
}
