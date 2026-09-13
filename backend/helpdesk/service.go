/**
 * BLOCK_HELPDESK_SERVICE_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   Business logic orchestration, SLA calculation, ticket workflows, and audit publishing.
 */

package helpdesk

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"campus/backend/audit"
)

type Service interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*HelpdeskCategory, error)
	ListCategories(ctx context.Context, tenantID string, activeOnly bool) ([]HelpdeskCategory, error)

	CreateTicket(ctx context.Context, req CreateTicketRequest) (*HelpdeskTicket, error)
	GetTicket(ctx context.Context, tenantID, id string) (*HelpdeskTicket, error)
	ListTickets(ctx context.Context, filter TicketFilter) ([]HelpdeskTicket, error)

	AssignTicket(ctx context.Context, tenantID, ticketID, staffID, actorID string) (*HelpdeskTicket, error)
	UpdateTicketStatus(ctx context.Context, tenantID, ticketID string, nextStatus HelpdeskTicketStatus, actorID string) (*HelpdeskTicket, error)

	AddMessage(ctx context.Context, req AddMessageRequest) (*HelpdeskMessage, error)
	ListMessages(ctx context.Context, tenantID, ticketID, userID string, isStaff bool) ([]HelpdeskMessage, error)

	EscalateTicket(ctx context.Context, req EscalateTicketRequest) (*HelpdeskEscalation, error)
	RateTicket(ctx context.Context, req RateTicketRequest) (*HelpdeskTicket, error)
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

func (s *service) logAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
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

type CreateCategoryRequest struct {
	TenantID           string           `json:"tenantId"`
	Code               string           `json:"code"`
	Name               string           `json:"name"`
	Description        string           `json:"description,omitempty"`
	DefaultPriority    HelpdeskPriority `json:"defaultPriority"`
	SLAResponseHours   int              `json:"slaResponseHours"`
	SLAResolutionHours int              `json:"slaResolutionHours"`
	ActorID            string           `json:"actorId"`
}

func (s *service) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*HelpdeskCategory, error) {
	if req.DefaultPriority == "" {
		req.DefaultPriority = PriorityMedium
	}
	if req.SLAResponseHours <= 0 {
		req.SLAResponseHours = 24
	}
	if req.SLAResolutionHours <= 0 {
		req.SLAResolutionHours = 72
	}

	now := time.Now().UTC()
	cat := &HelpdeskCategory{
		ID:                 uuid.New().String(),
		TenantID:           req.TenantID,
		Code:               req.Code,
		Name:               req.Name,
		Description:        req.Description,
		DefaultPriority:    req.DefaultPriority,
		SLAResponseHours:   req.SLAResponseHours,
		SLAResolutionHours: req.SLAResolutionHours,
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, req.ActorID, "STAFF", "helpdesk:category:created", "helpdesk_category", cat.ID, audit.StatusSuccess, map[string]interface{}{
		"code": cat.Code,
		"name": cat.Name,
	})

	return cat, nil
}

func (s *service) ListCategories(ctx context.Context, tenantID string, activeOnly bool) ([]HelpdeskCategory, error) {
	return s.repo.ListCategories(ctx, tenantID, activeOnly)
}

type CreateTicketRequest struct {
	TenantID    string            `json:"tenantId"`
	CategoryID  string            `json:"categoryId"`
	RequesterID string            `json:"requesterId"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Priority    *HelpdeskPriority `json:"priority,omitempty"`
}

func (s *service) CreateTicket(ctx context.Context, req CreateTicketRequest) (*HelpdeskTicket, error) {
	cat, err := s.repo.GetCategoryByID(ctx, req.TenantID, req.CategoryID)
	if err != nil {
		return nil, err
	}

	priority := cat.DefaultPriority
	if req.Priority != nil && *req.Priority != "" {
		priority = *req.Priority
	}

	seq, err := s.repo.GetNextTicketSequence(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	ticketNumber := FormatTicketNumber(now.Year(), seq)
	slaDueAt := CalculateSLADueAt(now, cat.SLAResolutionHours)

	ticket := &HelpdeskTicket{
		ID:           uuid.New().String(),
		TenantID:     req.TenantID,
		TicketNumber: ticketNumber,
		CategoryID:   req.CategoryID,
		RequesterID:  req.RequesterID,
		Title:        req.Title,
		Description:  req.Description,
		Priority:     priority,
		Status:       StatusOpen,
		SLADueAt:     slaDueAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateTicket(ctx, ticket); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, req.RequesterID, "STUDENT", "helpdesk:ticket:created", "helpdesk_ticket", ticket.ID, audit.StatusSuccess, map[string]interface{}{
		"ticketNumber": ticket.TicketNumber,
		"priority":     ticket.Priority,
		"categoryId":   ticket.CategoryID,
		"slaDueAt":     ticket.SLADueAt,
	})

	return ticket, nil
}

func (s *service) GetTicket(ctx context.Context, tenantID, id string) (*HelpdeskTicket, error) {
	return s.repo.GetTicketByID(ctx, tenantID, id)
}

func (s *service) ListTickets(ctx context.Context, filter TicketFilter) ([]HelpdeskTicket, error) {
	return s.repo.ListTickets(ctx, filter)
}

func (s *service) AssignTicket(ctx context.Context, tenantID, ticketID, staffID, actorID string) (*HelpdeskTicket, error) {
	ticket, err := s.repo.GetTicketByID(ctx, tenantID, ticketID)
	if err != nil {
		return nil, err
	}

	ticket.AssignedStaffID = &staffID
	if ticket.Status == StatusOpen {
		ticket.Status = StatusInProgress
	}
	ticket.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateTicket(ctx, ticket); err != nil {
		return nil, err
	}

	s.logAudit(tenantID, actorID, "STAFF", "helpdesk:ticket:assigned", "helpdesk_ticket", ticket.ID, audit.StatusSuccess, map[string]interface{}{
		"assignedStaffId": staffID,
		"status":          ticket.Status,
	})

	return ticket, nil
}

func (s *service) UpdateTicketStatus(ctx context.Context, tenantID, ticketID string, nextStatus HelpdeskTicketStatus, actorID string) (*HelpdeskTicket, error) {
	ticket, err := s.repo.GetTicketByID(ctx, tenantID, ticketID)
	if err != nil {
		return nil, err
	}

	if err := ValidateStateTransition(ticket.Status, nextStatus); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	prevStatus := ticket.Status
	ticket.Status = nextStatus
	ticket.UpdatedAt = now

	if nextStatus == StatusResolved && ticket.ResolvedAt == nil {
		ticket.ResolvedAt = &now
	}
	if nextStatus == StatusClosed && ticket.ClosedAt == nil {
		ticket.ClosedAt = &now
	}

	if err := s.repo.UpdateTicket(ctx, ticket); err != nil {
		return nil, err
	}

	s.logAudit(tenantID, actorID, "STAFF", "helpdesk:ticket:status_changed", "helpdesk_ticket", ticket.ID, audit.StatusSuccess, map[string]interface{}{
		"from": prevStatus,
		"to":   nextStatus,
	})

	return ticket, nil
}

type AddMessageRequest struct {
	TenantID       string   `json:"tenantId"`
	TicketID       string   `json:"ticketId"`
	SenderID       string   `json:"senderId"`
	IsStaffReply   bool     `json:"isStaffReply"`
	IsInternalNote bool     `json:"isInternalNote"`
	Message        string   `json:"message"`
	Attachments    []string `json:"attachments,omitempty"`
}

func (s *service) AddMessage(ctx context.Context, req AddMessageRequest) (*HelpdeskMessage, error) {
	ticket, err := s.repo.GetTicketByID(ctx, req.TenantID, req.TicketID)
	if err != nil {
		return nil, err
	}

	if err := ValidateMessagePost(ticket.Status, req.IsStaffReply, req.IsInternalNote); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	msg := &HelpdeskMessage{
		ID:             uuid.New().String(),
		TenantID:       req.TenantID,
		TicketID:       req.TicketID,
		SenderID:       req.SenderID,
		IsStaffReply:   req.IsStaffReply,
		IsInternalNote: req.IsInternalNote,
		Message:        req.Message,
		Attachments:    req.Attachments,
		CreatedAt:      now,
	}

	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		return nil, err
	}

	// Update ticket activity & SLA response tracker
	ticket.UpdatedAt = now
	if req.IsStaffReply && ticket.FirstResponseAt == nil && !req.IsInternalNote {
		ticket.FirstResponseAt = &now
	}

	// Dynamic conversation state updates:
	// If staff replies with public message, status -> WAITING_FOR_APPLICANT
	// If applicant replies, status -> IN_PROGRESS
	if !req.IsInternalNote {
		if req.IsStaffReply && ticket.Status == StatusInProgress {
			ticket.Status = StatusWaitingForApplicant
		} else if !req.IsStaffReply && ticket.Status == StatusWaitingForApplicant {
			ticket.Status = StatusInProgress
		}
	}
	_ = s.repo.UpdateTicket(ctx, ticket)

	s.logAudit(req.TenantID, req.SenderID, "USER", "helpdesk:ticket:message_added", "helpdesk_message", msg.ID, audit.StatusSuccess, map[string]interface{}{
		"ticketId":       ticket.ID,
		"isStaffReply":   req.IsStaffReply,
		"isInternalNote": req.IsInternalNote,
	})

	return msg, nil
}

func (s *service) ListMessages(ctx context.Context, tenantID, ticketID, userID string, isStaff bool) ([]HelpdeskMessage, error) {
	ticket, err := s.repo.GetTicketByID(ctx, tenantID, ticketID)
	if err != nil {
		return nil, err
	}

	if !isStaff && ticket.RequesterID != userID {
		return nil, ErrUnauthorizedTicketAccess
	}

	return s.repo.ListMessages(ctx, tenantID, ticketID, isStaff)
}

type EscalateTicketRequest struct {
	TenantID      string  `json:"tenantId"`
	TicketID      string  `json:"ticketId"`
	EscalatedToID *string `json:"escalatedToId,omitempty"`
	Reason        string  `json:"reason"`
	ActorID       string  `json:"actorId"`
}

func (s *service) EscalateTicket(ctx context.Context, req EscalateTicketRequest) (*HelpdeskEscalation, error) {
	ticket, err := s.repo.GetTicketByID(ctx, req.TenantID, req.TicketID)
	if err != nil {
		return nil, err
	}

	existingEscs, err := s.repo.ListEscalations(ctx, req.TenantID, req.TicketID)
	if err == nil {
		for _, e := range existingEscs {
			if e.Status == EscalationPending {
				return nil, ErrEscalationAlreadyPending
			}
		}
	}

	now := time.Now().UTC()
	esc := &HelpdeskEscalation{
		ID:            uuid.New().String(),
		TenantID:      req.TenantID,
		TicketID:      req.TicketID,
		EscalatedToID: req.EscalatedToID,
		Reason:        req.Reason,
		Status:        EscalationPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.CreateEscalation(ctx, esc); err != nil {
		return nil, err
	}

	// Bump priority to Urgent if not already
	ticket.Priority = PriorityUrgent
	ticket.UpdatedAt = now
	_ = s.repo.UpdateTicket(ctx, ticket)

	s.logAudit(req.TenantID, req.ActorID, "USER", "helpdesk:ticket:escalated", "helpdesk_escalation", esc.ID, audit.StatusSuccess, map[string]interface{}{
		"ticketId": req.TicketID,
		"reason":   req.Reason,
	})

	return esc, nil
}

type RateTicketRequest struct {
	TenantID string `json:"tenantId"`
	TicketID string `json:"ticketId"`
	UserID   string `json:"userId"`
	Rating   int    `json:"rating"`
	Feedback string `json:"feedback,omitempty"`
}

func (s *service) RateTicket(ctx context.Context, req RateTicketRequest) (*HelpdeskTicket, error) {
	ticket, err := s.repo.GetTicketByID(ctx, req.TenantID, req.TicketID)
	if err != nil {
		return nil, err
	}

	if ticket.RequesterID != req.UserID {
		return nil, ErrUnauthorizedTicketAccess
	}

	if err := ValidateRating(ticket.Status, req.Rating); err != nil {
		return nil, err
	}

	ticket.Rating = &req.Rating
	if req.Feedback != "" {
		ticket.Feedback = &req.Feedback
	}
	ticket.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateTicket(ctx, ticket); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, req.UserID, "STUDENT", "helpdesk:ticket:rated", "helpdesk_ticket", ticket.ID, audit.StatusSuccess, map[string]interface{}{
		"rating": req.Rating,
	})

	return ticket, nil
}
