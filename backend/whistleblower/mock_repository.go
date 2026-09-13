/**
 * BLOCK_WHISTLEBLOWER_MOCK_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   Concurrent in-memory mock repository test double for whistleblower subsystem.
 */

package whistleblower

import (
	"context"
	"sort"
	"sync"
)

type MockRepository struct {
	mu        sync.RWMutex
	reports   map[string]WhistleblowerReport
	messages  []WhistleblowerMessage
	evidences []WhistleblowerEvidence
	sequences map[string]int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		reports:   make(map[string]WhistleblowerReport),
		messages:  make([]WhistleblowerMessage, 0),
		evidences: make([]WhistleblowerEvidence, 0),
		sequences: make(map[string]int),
	}
}

func (m *MockRepository) CreateReport(ctx context.Context, report *WhistleblowerReport) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reports[report.ID] = *report
	return nil
}

func (m *MockRepository) GetReportByID(ctx context.Context, tenantID, id string) (*WhistleblowerReport, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.reports[id]
	if !ok || r.TenantID != tenantID {
		return nil, ErrReportNotFound
	}
	return &r, nil
}

func (m *MockRepository) GetReportByTrackingHash(ctx context.Context, tenantID, trackingHash string) (*WhistleblowerReport, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, r := range m.reports {
		if r.TenantID == tenantID && r.TrackingHash == trackingHash {
			return &r, nil
		}
	}
	return nil, ErrReportNotFound
}

func (m *MockRepository) ListReports(ctx context.Context, filter WhistleblowerReportFilter) ([]WhistleblowerReport, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]WhistleblowerReport, 0)
	for _, r := range m.reports {
		if r.TenantID != filter.TenantID {
			continue
		}
		if filter.Category != nil && r.Category != *filter.Category {
			continue
		}
		if filter.Severity != nil && r.Severity != *filter.Severity {
			continue
		}
		if filter.Status != nil && r.Status != *filter.Status {
			continue
		}
		res = append(res, r)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].SubmittedAt.After(res[j].SubmittedAt)
	})

	return res, nil
}

func (m *MockRepository) UpdateReport(ctx context.Context, report *WhistleblowerReport) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.reports[report.ID]; !ok {
		return ErrReportNotFound
	}
	m.reports[report.ID] = *report
	return nil
}

func (m *MockRepository) GetNextReportSequence(ctx context.Context, tenantID string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sequences[tenantID]++
	return m.sequences[tenantID], nil
}

func (m *MockRepository) CreateMessage(ctx context.Context, message *WhistleblowerMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, *message)
	return nil
}

func (m *MockRepository) ListMessages(ctx context.Context, tenantID, reportID string) ([]WhistleblowerMessage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]WhistleblowerMessage, 0)
	for _, msg := range m.messages {
		if msg.TenantID == tenantID && msg.ReportID == reportID {
			res = append(res, msg)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.Before(res[j].CreatedAt)
	})

	return res, nil
}

func (m *MockRepository) CreateEvidence(ctx context.Context, evidence *WhistleblowerEvidence) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.evidences = append(m.evidences, *evidence)
	return nil
}

func (m *MockRepository) ListEvidences(ctx context.Context, tenantID, reportID string) ([]WhistleblowerEvidence, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]WhistleblowerEvidence, 0)
	for _, ev := range m.evidences {
		if ev.TenantID == tenantID && ev.ReportID == reportID {
			res = append(res, ev)
		}
	}
	return res, nil
}
