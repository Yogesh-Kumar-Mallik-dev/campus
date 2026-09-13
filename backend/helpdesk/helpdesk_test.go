/**
 * BLOCK_HELPDESK_TEST_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   Unit test suite verifying ticketing invariants, SLA radar calculations, state machine, and rating policies.
 */

package helpdesk_test

import (
	"context"
	"errors"
	"testing"

	"campus/backend/helpdesk"
)

func setupTestService() (helpdesk.Service, *helpdesk.MockRepository) {
	mockRepo := helpdesk.NewMockRepository()
	srv := helpdesk.NewService(mockRepo, nil)
	return srv, mockRepo
}

func TestCategoryCreationAndDuplicates(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupTestService()

	tenantID := "tenant-test-1"

	cat, err := srv.CreateCategory(ctx, helpdesk.CreateCategoryRequest{
		TenantID:           tenantID,
		Code:               "ACADEMIC",
		Name:               "Academic Inquiries",
		Description:        "Issues related to courses, exams, and grades",
		DefaultPriority:    helpdesk.PriorityHigh,
		SLAResponseHours:   12,
		SLAResolutionHours: 48,
		ActorID:            "admin-1",
	})
	if err != nil {
		t.Fatalf("expected successful category creation, got error: %v", err)
	}

	if cat.ID == "" || cat.Code != "ACADEMIC" || cat.SLAResolutionHours != 48 {
		t.Fatalf("unexpected category fields: %+v", cat)
	}

	// Duplicate code rejection
	_, err = srv.CreateCategory(ctx, helpdesk.CreateCategoryRequest{
		TenantID: tenantID,
		Code:     "ACADEMIC",
		Name:     "Duplicate Academic",
		ActorID:  "admin-1",
	})
	if !errors.Is(err, helpdesk.ErrCategoryCodeDuplicate) {
		t.Fatalf("expected ErrCategoryCodeDuplicate, got: %v", err)
	}

	// Listing categories
	list, err := srv.ListCategories(ctx, tenantID, true)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 category in list, got: %d, err: %v", len(list), err)
	}
}

func TestTicketCreationAndSLACalculation(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupTestService()

	tenantID := "tenant-test-1"

	cat, err := srv.CreateCategory(ctx, helpdesk.CreateCategoryRequest{
		TenantID:           tenantID,
		Code:               "IT_SUPPORT",
		Name:               "IT & Campus WiFi",
		SLAResolutionHours: 24,
		ActorID:            "admin-1",
	})
	if err != nil {
		t.Fatalf("category create failed: %v", err)
	}

	ticket, err := srv.CreateTicket(ctx, helpdesk.CreateTicketRequest{
		TenantID:    tenantID,
		CategoryID:  cat.ID,
		RequesterID: "student-100",
		Title:       "WiFi connection failure in Hostel Block B",
		Description: "Cannot connect to campus enterprise network since yesterday evening",
	})
	if err != nil {
		t.Fatalf("expected ticket creation success, got: %v", err)
	}

	if ticket.TicketNumber != "HD-2026-00001" && ticket.TicketNumber == "" {
		t.Fatalf("ticket number format invalid: %s", ticket.TicketNumber)
	}
	if ticket.Status != helpdesk.StatusOpen {
		t.Fatalf("expected initial status OPEN, got: %s", ticket.Status)
	}
	if ticket.SLADueAt.Before(ticket.CreatedAt) {
		t.Fatalf("SLA due date must be in the future")
	}

	// Ticket 2 should increment sequence
	ticket2, err := srv.CreateTicket(ctx, helpdesk.CreateTicketRequest{
		TenantID:    tenantID,
		CategoryID:  cat.ID,
		RequesterID: "student-101",
		Title:       "Hostel LAN port broken",
		Description: "Port RJ45 pin is bent",
	})
	if err != nil {
		t.Fatalf("expected ticket2 creation success, got: %v", err)
	}
	if ticket2.TicketNumber == ticket.TicketNumber {
		t.Fatalf("expected unique ticket numbers, got duplicate: %s", ticket2.TicketNumber)
	}
}

func TestTicketAssignmentAndStatusTransitions(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupTestService()

	tenantID := "tenant-test-1"
	cat, _ := srv.CreateCategory(ctx, helpdesk.CreateCategoryRequest{
		TenantID: tenantID,
		Code:     "BILLING",
		Name:     "Fee Payments & Invoicing",
		ActorID:  "admin-1",
	})

	ticket, _ := srv.CreateTicket(ctx, helpdesk.CreateTicketRequest{
		TenantID:    tenantID,
		CategoryID:  cat.ID,
		RequesterID: "student-100",
		Title:       "Fee transaction pending",
		Description: "Amount debited but receipt not generated",
	})

	// Assign ticket to staff -> moves status to IN_PROGRESS
	assigned, err := srv.AssignTicket(ctx, tenantID, ticket.ID, "staff-tech-1", "admin-1")
	if err != nil {
		t.Fatalf("expected assign success, got: %v", err)
	}
	if assigned.Status != helpdesk.StatusInProgress || *assigned.AssignedStaffID != "staff-tech-1" {
		t.Fatalf("expected assigned ticket in progress, got: %+v", assigned)
	}

	// Move to WAITING_FOR_APPLICANT
	updated, err := srv.UpdateTicketStatus(ctx, tenantID, ticket.ID, helpdesk.StatusWaitingForApplicant, "staff-tech-1")
	if err != nil || updated.Status != helpdesk.StatusWaitingForApplicant {
		t.Fatalf("expected waiting for applicant status, got: %+v, err: %v", updated, err)
	}

	// Resolve ticket
	resolved, err := srv.UpdateTicketStatus(ctx, tenantID, ticket.ID, helpdesk.StatusResolved, "staff-tech-1")
	if err != nil || resolved.Status != helpdesk.StatusResolved || resolved.ResolvedAt == nil {
		t.Fatalf("expected resolved status with timestamp, got: %+v, err: %v", resolved, err)
	}

	// Close ticket
	closed, err := srv.UpdateTicketStatus(ctx, tenantID, ticket.ID, helpdesk.StatusClosed, "student-100")
	if err != nil || closed.Status != helpdesk.StatusClosed || closed.ClosedAt == nil {
		t.Fatalf("expected closed status with timestamp, got: %+v, err: %v", closed, err)
	}
}

func TestMessagesAndInternalNotesSecurity(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupTestService()

	tenantID := "tenant-test-1"
	cat, _ := srv.CreateCategory(ctx, helpdesk.CreateCategoryRequest{
		TenantID: tenantID,
		Code:     "EXAM",
		Name:     "Examination Wing",
		ActorID:  "admin-1",
	})

	ticket, _ := srv.CreateTicket(ctx, helpdesk.CreateTicketRequest{
		TenantID:    tenantID,
		CategoryID:  cat.ID,
		RequesterID: "student-100",
		Title:       "Exam hall ticket misprint",
		Description: "My name has spelling mistake on hall ticket",
	})

	// Non-staff trying to add internal note -> REJECTED
	_, err := srv.AddMessage(ctx, helpdesk.AddMessageRequest{
		TenantID:       tenantID,
		TicketID:       ticket.ID,
		SenderID:       "student-100",
		IsStaffReply:   false,
		IsInternalNote: true,
		Message:        "Trying to write secret note",
	})
	if !errors.Is(err, helpdesk.ErrInternalNoteUnauthorized) {
		t.Fatalf("expected ErrInternalNoteUnauthorized, got: %v", err)
	}

	// Staff adds internal note -> SUCCESS
	internalMsg, err := srv.AddMessage(ctx, helpdesk.AddMessageRequest{
		TenantID:       tenantID,
		TicketID:       ticket.ID,
		SenderID:       "staff-controller-1",
		IsStaffReply:   true,
		IsInternalNote: true,
		Message:        "Checking registration DB records before answering applicant",
	})
	if err != nil || !internalMsg.IsInternalNote {
		t.Fatalf("expected staff internal note success, got: %+v, err: %v", internalMsg, err)
	}

	// Staff lists messages -> sees internal note
	staffMsgs, err := srv.ListMessages(ctx, tenantID, ticket.ID, "staff-controller-1", true)
	if err != nil || len(staffMsgs) != 1 {
		t.Fatalf("expected 1 message for staff, got: %d, err: %v", len(staffMsgs), err)
	}

	// Student lists messages -> internal note filtered out
	studentMsgs, err := srv.ListMessages(ctx, tenantID, ticket.ID, "student-100", false)
	if err != nil || len(studentMsgs) != 0 {
		t.Fatalf("expected 0 public messages for student, got: %d, err: %v", len(studentMsgs), err)
	}

	// Other student tries to view ticket messages -> REJECTED
	_, err = srv.ListMessages(ctx, tenantID, ticket.ID, "student-999", false)
	if !errors.Is(err, helpdesk.ErrUnauthorizedTicketAccess) {
		t.Fatalf("expected ErrUnauthorizedTicketAccess for unauthorized student, got: %v", err)
	}

	// Staff sends public response -> ticket records firstResponseAt and sets status WAITING_FOR_APPLICANT
	_, err = srv.AssignTicket(ctx, tenantID, ticket.ID, "staff-controller-1", "admin-1")
	if err != nil {
		t.Fatalf("assign failed: %v", err)
	}

	_, err = srv.AddMessage(ctx, helpdesk.AddMessageRequest{
		TenantID:       tenantID,
		TicketID:       ticket.ID,
		SenderID:       "staff-controller-1",
		IsStaffReply:   true,
		IsInternalNote: false,
		Message:        "Please upload your national identity proof",
	})
	if err != nil {
		t.Fatalf("public reply failed: %v", err)
	}

	updatedTicket, _ := srv.GetTicket(ctx, tenantID, ticket.ID)
	if updatedTicket.FirstResponseAt == nil {
		t.Fatalf("expected firstResponseAt to be recorded")
	}
	if updatedTicket.Status != helpdesk.StatusWaitingForApplicant {
		t.Fatalf("expected status WAITING_FOR_APPLICANT, got: %s", updatedTicket.Status)
	}
}

func TestTicketEscalationAndDuplicateGuard(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupTestService()

	tenantID := "tenant-test-1"
	cat, _ := srv.CreateCategory(ctx, helpdesk.CreateCategoryRequest{
		TenantID: tenantID,
		Code:     "HOSTEL",
		Name:     "Hostel Facilities",
		ActorID:  "admin-1",
	})

	ticket, _ := srv.CreateTicket(ctx, helpdesk.CreateTicketRequest{
		TenantID:    tenantID,
		CategoryID:  cat.ID,
		RequesterID: "student-100",
		Title:       "Water leakage in room",
		Description: "Pipe burst in 3rd floor bathroom",
	})

	esc, err := srv.EscalateTicket(ctx, helpdesk.EscalateTicketRequest{
		TenantID: tenantID,
		TicketID: ticket.ID,
		Reason:   "Urgent health and safety hazard",
		ActorID:  "student-100",
	})
	if err != nil || esc.Status != helpdesk.EscalationPending {
		t.Fatalf("expected successful escalation, got: %+v, err: %v", esc, err)
	}

	// Ticket priority should have escalated to URGENT
	updatedTicket, _ := srv.GetTicket(ctx, tenantID, ticket.ID)
	if updatedTicket.Priority != helpdesk.PriorityUrgent {
		t.Fatalf("expected priority URGENT, got: %s", updatedTicket.Priority)
	}

	// Duplicate pending escalation attempt -> REJECTED
	_, err = srv.EscalateTicket(ctx, helpdesk.EscalateTicketRequest{
		TenantID: tenantID,
		TicketID: ticket.ID,
		Reason:   "Still leaking",
		ActorID:  "student-100",
	})
	if !errors.Is(err, helpdesk.ErrEscalationAlreadyPending) {
		t.Fatalf("expected ErrEscalationAlreadyPending, got: %v", err)
	}
}

func TestTicketRatingAndFeedback(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupTestService()

	tenantID := "tenant-test-1"
	cat, _ := srv.CreateCategory(ctx, helpdesk.CreateCategoryRequest{
		TenantID: tenantID,
		Code:     "MESS",
		Name:     "Mess & Dining",
		ActorID:  "admin-1",
	})

	ticket, _ := srv.CreateTicket(ctx, helpdesk.CreateTicketRequest{
		TenantID:    tenantID,
		CategoryID:  cat.ID,
		RequesterID: "student-100",
		Title:       "Special dietary request",
		Description: "Gluten free options required",
	})

	// Rating while still OPEN -> REJECTED
	_, err := srv.RateTicket(ctx, helpdesk.RateTicketRequest{
		TenantID: tenantID,
		TicketID: ticket.ID,
		UserID:   "student-100",
		Rating:   5,
	})
	if !errors.Is(err, helpdesk.ErrRatingNotAllowed) {
		t.Fatalf("expected ErrRatingNotAllowed for open ticket, got: %v", err)
	}

	// Resolve ticket
	_, _ = srv.UpdateTicketStatus(ctx, tenantID, ticket.ID, helpdesk.StatusResolved, "staff-1")

	// Invalid score out of bounds (< 1 or > 5)
	_, err = srv.RateTicket(ctx, helpdesk.RateTicketRequest{
		TenantID: tenantID,
		TicketID: ticket.ID,
		UserID:   "student-100",
		Rating:   6,
	})
	if !errors.Is(err, helpdesk.ErrInvalidRatingScore) {
		t.Fatalf("expected ErrInvalidRatingScore for 6 stars, got: %v", err)
	}

	// Valid rating & feedback
	rated, err := srv.RateTicket(ctx, helpdesk.RateTicketRequest{
		TenantID: tenantID,
		TicketID: ticket.ID,
		UserID:   "student-100",
		Rating:   5,
		Feedback: "Resolved immediately with dining manager. Excellent!",
	})
	if err != nil || *rated.Rating != 5 || *rated.Feedback == "" {
		t.Fatalf("expected rating success, got: %+v, err: %v", rated, err)
	}
}
