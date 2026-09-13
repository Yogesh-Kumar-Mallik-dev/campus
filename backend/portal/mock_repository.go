/**
 * BLOCK_PORTAL_MOCK_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   Concurrent in-memory test double for portal repository contracts.
 */

package portal

import (
	"context"
	"sort"
	"sync"
)

type MockRepository struct {
	mu        sync.RWMutex
	landings  map[string]PortalLandingPage
	programs  map[string]PortalProgramCatalog
	inquiries map[string]PortalPublicInquiry
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		landings:  make(map[string]PortalLandingPage),
		programs:  make(map[string]PortalProgramCatalog),
		inquiries: make(map[string]PortalPublicInquiry),
	}
}

func (m *MockRepository) GetLandingPage(ctx context.Context, tenantID string) (*PortalLandingPage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.landings[tenantID]
	if !ok {
		return nil, ErrLandingPageNotFound
	}
	return &p, nil
}

func (m *MockRepository) UpsertLandingPage(ctx context.Context, page *PortalLandingPage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.landings[page.TenantID] = *page
	return nil
}

func (m *MockRepository) CreateProgram(ctx context.Context, program *PortalProgramCatalog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.programs {
		if p.TenantID == program.TenantID && p.ProgramCode == program.ProgramCode {
			return ErrDuplicateProgramCode
		}
	}

	m.programs[program.ID] = *program
	return nil
}

func (m *MockRepository) GetProgramByID(ctx context.Context, tenantID, id string) (*PortalProgramCatalog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.programs[id]
	if !ok || p.TenantID != tenantID {
		return nil, ErrProgramNotFound
	}
	return &p, nil
}

func (m *MockRepository) GetProgramByCode(ctx context.Context, tenantID, code string) (*PortalProgramCatalog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.programs {
		if p.TenantID == tenantID && p.ProgramCode == code {
			return &p, nil
		}
	}
	return nil, ErrProgramNotFound
}

func (m *MockRepository) ListPrograms(ctx context.Context, filter ProgramCatalogFilter) ([]PortalProgramCatalog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]PortalProgramCatalog, 0)
	for _, p := range m.programs {
		if p.TenantID != filter.TenantID {
			continue
		}
		if filter.DegreeType != nil && p.DegreeType != *filter.DegreeType {
			continue
		}
		if filter.Featured != nil && p.IsFeatured != *filter.Featured {
			continue
		}
		res = append(res, p)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ProgramName < res[j].ProgramName
	})

	return res, nil
}

func (m *MockRepository) UpdateProgram(ctx context.Context, program *PortalProgramCatalog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.programs[program.ID]; !ok {
		return ErrProgramNotFound
	}
	m.programs[program.ID] = *program
	return nil
}

func (m *MockRepository) CreateInquiry(ctx context.Context, inquiry *PortalPublicInquiry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.inquiries[inquiry.ID] = *inquiry
	return nil
}

func (m *MockRepository) GetInquiryByID(ctx context.Context, tenantID, id string) (*PortalPublicInquiry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inq, ok := m.inquiries[id]
	if !ok || inq.TenantID != tenantID {
		return nil, ErrInquiryNotFound
	}
	return &inq, nil
}

func (m *MockRepository) ListInquiries(ctx context.Context, filter PublicInquiryFilter) ([]PortalPublicInquiry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]PortalPublicInquiry, 0)
	for _, inq := range m.inquiries {
		if inq.TenantID != filter.TenantID {
			continue
		}
		if filter.Status != nil && inq.Status != *filter.Status {
			continue
		}
		if filter.AssignedCounselorID != nil && (inq.AssignedCounselorID == nil || *inq.AssignedCounselorID != *filter.AssignedCounselorID) {
			continue
		}
		res = append(res, inq)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(res[j].CreatedAt)
	})

	return res, nil
}

func (m *MockRepository) UpdateInquiry(ctx context.Context, inquiry *PortalPublicInquiry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.inquiries[inquiry.ID]; !ok {
		return ErrInquiryNotFound
	}
	m.inquiries[inquiry.ID] = *inquiry
	return nil
}
