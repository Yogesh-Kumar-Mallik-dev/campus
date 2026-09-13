/**
 * BLOCK_HELPDESK_DOMAIN_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   Domain entities, value objects, and business invariant rules for ticketing, SLA monitoring, and resolution desk.
 */

package helpdesk

import (
	"fmt"
	"time"
)

type HelpdeskPriority string

const (
	PriorityLow    HelpdeskPriority = "LOW"
	PriorityMedium HelpdeskPriority = "MEDIUM"
	PriorityHigh   HelpdeskPriority = "HIGH"
	PriorityUrgent HelpdeskPriority = "URGENT"
)

type HelpdeskTicketStatus string

const (
	StatusOpen                HelpdeskTicketStatus = "OPEN"
	StatusInProgress          HelpdeskTicketStatus = "IN_PROGRESS"
	StatusWaitingForApplicant HelpdeskTicketStatus = "WAITING_FOR_APPLICANT"
	StatusResolved            HelpdeskTicketStatus = "RESOLVED"
	StatusClosed              HelpdeskTicketStatus = "CLOSED"
)

type HelpdeskEscalationStatus string

const (
	EscalationPending      HelpdeskEscalationStatus = "PENDING"
	EscalationAcknowledged HelpdeskEscalationStatus = "ACKNOWLEDGED"
	EscalationResolved     HelpdeskEscalationStatus = "RESOLVED"
)

type HelpdeskCategory struct {
	ID                 string           `json:"id"`
	TenantID           string           `json:"tenantId"`
	Code               string           `json:"code"`
	Name               string           `json:"name"`
	Description        string           `json:"description,omitempty"`
	DefaultPriority    HelpdeskPriority `json:"defaultPriority"`
	SLAResponseHours   int              `json:"slaResponseHours"`
	SLAResolutionHours int              `json:"slaResolutionHours"`
	IsActive           bool             `json:"isActive"`
	CreatedAt          time.Time        `json:"createdAt"`
	UpdatedAt          time.Time        `json:"updatedAt"`
}

type HelpdeskTicket struct {
	ID              string               `json:"id"`
	TenantID        string               `json:"tenantId"`
	TicketNumber    string               `json:"ticketNumber"`
	CategoryID      string               `json:"categoryId"`
	RequesterID     string               `json:"requesterId"`
	AssignedStaffID *string              `json:"assignedStaffId,omitempty"`
	Title           string               `json:"title"`
	Description     string               `json:"description"`
	Priority        HelpdeskPriority     `json:"priority"`
	Status          HelpdeskTicketStatus `json:"status"`
	SLADueAt        time.Time            `json:"slaDueAt"`
	FirstResponseAt *time.Time           `json:"firstResponseAt,omitempty"`
	ResolvedAt      *time.Time           `json:"resolvedAt,omitempty"`
	ClosedAt        *time.Time           `json:"closedAt,omitempty"`
	Rating          *int                 `json:"rating,omitempty"`
	Feedback        *string              `json:"feedback,omitempty"`
	CreatedAt       time.Time            `json:"createdAt"`
	UpdatedAt       time.Time            `json:"updatedAt"`
}

type HelpdeskMessage struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenantId"`
	TicketID       string    `json:"ticketId"`
	SenderID       string    `json:"senderId"`
	IsStaffReply   bool      `json:"isStaffReply"`
	IsInternalNote bool      `json:"isInternalNote"`
	Message        string    `json:"message"`
	Attachments    []string  `json:"attachments,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type HelpdeskEscalation struct {
	ID            string                   `json:"id"`
	TenantID      string                   `json:"tenantId"`
	TicketID      string                   `json:"ticketId"`
	EscalatedToID *string                  `json:"escalatedToId,omitempty"`
	Reason        string                   `json:"reason"`
	Status        HelpdeskEscalationStatus `json:"status"`
	CreatedAt     time.Time                `json:"createdAt"`
	UpdatedAt     time.Time                `json:"updatedAt"`
}

type TicketFilter struct {
	TenantID        string
	RequesterID     *string
	AssignedStaffID *string
	CategoryID      *string
	Status          *HelpdeskTicketStatus
	Priority        *HelpdeskPriority
}

// FormatTicketNumber generates deterministic, human-readable ticket identifiers.
func FormatTicketNumber(year int, seq int) string {
	return fmt.Sprintf("HD-%d-%05d", year, seq)
}

// CalculateSLADueAt computes SLA deadline based on category resolution hours.
func CalculateSLADueAt(created time.Time, resolutionHours int) time.Time {
	if resolutionHours <= 0 {
		resolutionHours = 72 // default 3 days
	}
	return created.Add(time.Duration(resolutionHours) * time.Hour)
}

// ValidateStateTransition checks whether a status transition is legally permitted.
func ValidateStateTransition(current, next HelpdeskTicketStatus) error {
	if current == next {
		return nil
	}

	switch current {
	case StatusOpen:
		if next == StatusInProgress || next == StatusResolved || next == StatusClosed {
			return nil
		}
	case StatusInProgress:
		if next == StatusWaitingForApplicant || next == StatusResolved || next == StatusClosed {
			return nil
		}
	case StatusWaitingForApplicant:
		if next == StatusInProgress || next == StatusResolved || next == StatusClosed {
			return nil
		}
	case StatusResolved:
		if next == StatusClosed || next == StatusInProgress { // Can reopen if applicant disputes
			return nil
		}
	case StatusClosed:
		if next == StatusOpen || next == StatusInProgress { // Super-admin reopen only
			return nil
		}
	}

	return fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidTicketTransition, current, next)
}

// ValidateMessagePost verifies that the message post satisfies security & lifecycle policies.
func ValidateMessagePost(ticketStatus HelpdeskTicketStatus, isStaff, isInternalNote bool) error {
	if ticketStatus == StatusClosed {
		return ErrTicketClosed
	}
	if isInternalNote && !isStaff {
		return ErrInternalNoteUnauthorized
	}
	return nil
}

// ValidateRating verifies rating bounds and ticket resolution requirement.
func ValidateRating(ticketStatus HelpdeskTicketStatus, rating int) error {
	if ticketStatus != StatusResolved && ticketStatus != StatusClosed {
		return ErrRatingNotAllowed
	}
	if rating < 1 || rating > 5 {
		return ErrInvalidRatingScore
	}
	return nil
}
