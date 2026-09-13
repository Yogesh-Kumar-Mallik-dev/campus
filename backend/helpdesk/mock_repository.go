/**
 * BLOCK_HELPDESK_MOCK_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   In-memory concurrent test double for helpdesk repository contracts.
 */

package helpdesk

import (
	"context"
	"sort"
	"sync"
)

type MockRepository struct {
	mu          sync.RWMutex
	categories  map[string]HelpdeskCategory
	tickets     map[string]HelpdeskTicket
	messages    []HelpdeskMessage
	escalations []HelpdeskEscalation
	sequences   map[string]int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		categories:  make(map[string]HelpdeskCategory),
		tickets:     make(map[string]HelpdeskTicket),
		messages:    make([]HelpdeskMessage, 0),
		escalations: make([]HelpdeskEscalation, 0),
		sequences:   make(map[string]int),
	}
}

func (m *MockRepository) CreateCategory(ctx context.Context, cat *HelpdeskCategory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, existing := range m.categories {
		if existing.TenantID == cat.TenantID && existing.Code == cat.Code {
			return ErrCategoryCodeDuplicate
		}
	}

	m.categories[cat.ID] = *cat
	return nil
}

func (m *MockRepository) GetCategoryByID(ctx context.Context, tenantID, id string) (*HelpdeskCategory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cat, ok := m.categories[id]
	if !ok || cat.TenantID != tenantID {
		return nil, ErrCategoryNotFound
	}
	return &cat, nil
}

func (m *MockRepository) GetCategoryByCode(ctx context.Context, tenantID, code string) (*HelpdeskCategory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, cat := range m.categories {
		if cat.TenantID == tenantID && cat.Code == code {
			return &cat, nil
		}
	}
	return nil, ErrCategoryNotFound
}

func (m *MockRepository) ListCategories(ctx context.Context, tenantID string, activeOnly bool) ([]HelpdeskCategory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]HelpdeskCategory, 0)
	for _, cat := range m.categories {
		if cat.TenantID == tenantID {
			if !activeOnly || cat.IsActive {
				res = append(res, cat)
			}
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res, nil
}

func (m *MockRepository) CreateTicket(ctx context.Context, ticket *HelpdeskTicket) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tickets[ticket.ID] = *ticket
	return nil
}

func (m *MockRepository) GetTicketByID(ctx context.Context, tenantID, id string) (*HelpdeskTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ticket, ok := m.tickets[id]
	if !ok || ticket.TenantID != tenantID {
		return nil, ErrTicketNotFound
	}
	return &ticket, nil
}

func (m *MockRepository) GetTicketByNumber(ctx context.Context, tenantID, number string) (*HelpdeskTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, ticket := range m.tickets {
		if ticket.TenantID == tenantID && ticket.TicketNumber == number {
			return &ticket, nil
		}
	}
	return nil, ErrTicketNotFound
}

func (m *MockRepository) ListTickets(ctx context.Context, filter TicketFilter) ([]HelpdeskTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]HelpdeskTicket, 0)
	for _, t := range m.tickets {
		if t.TenantID != filter.TenantID {
			continue
		}
		if filter.RequesterID != nil && t.RequesterID != *filter.RequesterID {
			continue
		}
		if filter.AssignedStaffID != nil && (t.AssignedStaffID == nil || *t.AssignedStaffID != *filter.AssignedStaffID) {
			continue
		}
		if filter.CategoryID != nil && t.CategoryID != *filter.CategoryID {
			continue
		}
		if filter.Status != nil && t.Status != *filter.Status {
			continue
		}
		if filter.Priority != nil && t.Priority != *filter.Priority {
			continue
		}
		res = append(res, t)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(res[j].CreatedAt)
	})

	return res, nil
}

func (m *MockRepository) UpdateTicket(ctx context.Context, ticket *HelpdeskTicket) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tickets[ticket.ID]; !ok {
		return ErrTicketNotFound
	}
	m.tickets[ticket.ID] = *ticket
	return nil
}

func (m *MockRepository) GetNextTicketSequence(ctx context.Context, tenantID string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sequences[tenantID]++
	return m.sequences[tenantID], nil
}

func (m *MockRepository) CreateMessage(ctx context.Context, msg *HelpdeskMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, *msg)
	return nil
}

func (m *MockRepository) ListMessages(ctx context.Context, tenantID, ticketID string, includeInternal bool) ([]HelpdeskMessage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]HelpdeskMessage, 0)
	for _, msg := range m.messages {
		if msg.TenantID == tenantID && msg.TicketID == ticketID {
			if !includeInternal && msg.IsInternalNote {
				continue
			}
			res = append(res, msg)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.Before(res[j].CreatedAt)
	})

	return res, nil
}

func (m *MockRepository) CreateEscalation(ctx context.Context, esc *HelpdeskEscalation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.escalations = append(m.escalations, *esc)
	return nil
}

func (m *MockRepository) ListEscalations(ctx context.Context, tenantID, ticketID string) ([]HelpdeskEscalation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]HelpdeskEscalation, 0)
	for _, esc := range m.escalations {
		if esc.TenantID == tenantID && esc.TicketID == ticketID {
			res = append(res, esc)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(res[j].CreatedAt)
	})

	return res, nil
}
