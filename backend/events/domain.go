/**
 * BLOCK_EVENTS_DOMAIN_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   Domain entities, value objects, and business invariant rules for campus events, venue conflict checks, and QR gate ticketing.
 */

package events

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type EventCategory string

const (
	EventCategoryCultural         EventCategory = "CULTURAL"
	EventCategoryTechnical        EventCategory = "TECHNICAL"
	EventCategorySports           EventCategory = "SPORTS"
	EventCategoryAcademicWorkshop EventCategory = "ACADEMIC_WORKSHOP"
	EventCategoryGuestLecture     EventCategory = "GUEST_LECTURE"
	EventCategoryCareerFair       EventCategory = "CAREER_FAIR"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "DRAFT"
	EventStatusPublished EventStatus = "PUBLISHED"
	EventStatusCancelled EventStatus = "CANCELLED"
	EventStatusCompleted EventStatus = "COMPLETED"
)

type TicketStatus string

const (
	TicketStatusConfirmed TicketStatus = "CONFIRMED"
	TicketStatusCheckedIn TicketStatus = "CHECKED_IN"
	TicketStatusCancelled TicketStatus = "CANCELLED"
)

type VenueBookingStatus string

const (
	VenueBookingStatusConfirmed VenueBookingStatus = "CONFIRMED"
	VenueBookingStatusCancelled VenueBookingStatus = "CANCELLED"
)

type CampusEvent struct {
	ID                   string        `json:"id"`
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
	Status               EventStatus   `json:"status"`
	IsTicketed           bool          `json:"isTicketed"`
	TicketPrice          float64       `json:"ticketPrice"`
	MaxTickets           int           `json:"maxTickets"`
	TicketsSold          int           `json:"ticketsSold"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt"`
}

type EventTicket struct {
	ID            string       `json:"id"`
	TenantID      string       `json:"tenantId"`
	EventID       string       `json:"eventId"`
	StudentID     string       `json:"studentId"`
	TicketCode    string       `json:"ticketCode"`
	Status        TicketStatus `json:"status"`
	CheckedInAt   *time.Time   `json:"checkedInAt,omitempty"`
	CheckedInByID *string      `json:"checkedInById,omitempty"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

type EventVenueBooking struct {
	ID              string             `json:"id"`
	TenantID        string             `json:"tenantId"`
	VenueName       string             `json:"venueName"`
	EventID         string             `json:"eventId"`
	BookedStartTime time.Time          `json:"bookedStartTime"`
	BookedEndTime   time.Time          `json:"bookedEndTime"`
	BookedByID      string             `json:"bookedById"`
	Status          VenueBookingStatus `json:"status"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

// CheckVenueOverlap returns true if two time intervals intersect.
func CheckVenueOverlap(existingStart, existingEnd, newStart, newEnd time.Time) bool {
	return newStart.Before(existingEnd) && newEnd.After(existingStart)
}

// GenerateTicketCode generates a unique cryptographic token for QR ticketing.
func GenerateTicketCode(tenantID, eventID, studentID string, timestamp time.Time) string {
	raw := fmt.Sprintf("%s:%s:%s:%d", tenantID, eventID, studentID, timestamp.UnixNano())
	hash := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("TCK-%s", hex.EncodeToString(hash[:6]))
}

// ValidateCheckIn verifies ticket validity prior to gate admission.
func ValidateCheckIn(ticketStatus TicketStatus) error {
	if ticketStatus == TicketStatusCheckedIn {
		return ErrTicketAlreadyCheckedIn
	}
	if ticketStatus == TicketStatusCancelled {
		return ErrTicketCancelled
	}
	return nil
}
