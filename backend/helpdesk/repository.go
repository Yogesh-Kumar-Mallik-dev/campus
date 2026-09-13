/**
 * BLOCK_HELPDESK_REPO_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   Repository interface contracts for ticket persistence, categories, messages, and escalations.
 */

package helpdesk

import "context"

type Repository interface {
	// Categories
	CreateCategory(ctx context.Context, cat *HelpdeskCategory) error
	GetCategoryByID(ctx context.Context, tenantID, id string) (*HelpdeskCategory, error)
	GetCategoryByCode(ctx context.Context, tenantID, code string) (*HelpdeskCategory, error)
	ListCategories(ctx context.Context, tenantID string, activeOnly bool) ([]HelpdeskCategory, error)

	// Tickets
	CreateTicket(ctx context.Context, ticket *HelpdeskTicket) error
	GetTicketByID(ctx context.Context, tenantID, id string) (*HelpdeskTicket, error)
	GetTicketByNumber(ctx context.Context, tenantID, number string) (*HelpdeskTicket, error)
	ListTickets(ctx context.Context, filter TicketFilter) ([]HelpdeskTicket, error)
	UpdateTicket(ctx context.Context, ticket *HelpdeskTicket) error
	GetNextTicketSequence(ctx context.Context, tenantID string) (int, error)

	// Messages
	CreateMessage(ctx context.Context, msg *HelpdeskMessage) error
	ListMessages(ctx context.Context, tenantID, ticketID string, includeInternal bool) ([]HelpdeskMessage, error)

	// Escalations
	CreateEscalation(ctx context.Context, esc *HelpdeskEscalation) error
	ListEscalations(ctx context.Context, tenantID, ticketID string) ([]HelpdeskEscalation, error)
}
