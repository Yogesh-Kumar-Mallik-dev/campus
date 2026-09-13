/**
 * BLOCK_AUDIT_MOCK_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Thread-safe in-memory test doubles for repository and subscriber testing.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"
)

// MockRepository provides an in-memory thread-safe implementation of Repository.
type MockRepository struct {
	mu          sync.RWMutex
	logs        []*AuditLog
	checkpoints []*AuditVerificationCheckpoint
}

// NewMockRepository initializes a clean in-memory audit repository.
func NewMockRepository() *MockRepository {
	return &MockRepository{
		logs:        make([]*AuditLog, 0),
		checkpoints: make([]*AuditVerificationCheckpoint, 0),
	}
}

func (m *MockRepository) Append(ctx context.Context, log *AuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clone to simulate isolated persistence
	cloned := *log
	m.logs = append(m.logs, &cloned)
	return nil
}

func (m *MockRepository) AppendBatch(ctx context.Context, logs []*AuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, l := range logs {
		cloned := *l
		m.logs = append(m.logs, &cloned)
	}
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, tenantID, id string) (*AuditLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, l := range m.logs {
		if l.TenantID == tenantID && l.ID == id {
			cloned := *l
			return &cloned, nil
		}
	}
	return nil, ErrAuditNotFound
}

func (m *MockRepository) GetLatestRecord(ctx context.Context, tenantID string) (*AuditLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for i := len(m.logs) - 1; i >= 0; i-- {
		if m.logs[i].TenantID == tenantID {
			cloned := *m.logs[i]
			return &cloned, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) Query(ctx context.Context, filter AuditFilter) ([]*AuditLog, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var matched []*AuditLog

	for _, l := range m.logs {
		if l.TenantID != filter.TenantID {
			continue
		}
		if filter.ActorID != nil && (l.ActorID == nil || *l.ActorID != *filter.ActorID) {
			continue
		}
		if filter.ActorType != nil && l.ActorType != *filter.ActorType {
			continue
		}
		if filter.Action != nil && !strings.Contains(strings.ToLower(l.Action), strings.ToLower(*filter.Action)) {
			continue
		}
		if filter.ResourceType != nil && l.ResourceType != *filter.ResourceType {
			continue
		}
		if filter.ResourceID != nil && (l.ResourceID == nil || *l.ResourceID != *filter.ResourceID) {
			continue
		}
		if filter.Status != nil && l.Status != *filter.Status {
			continue
		}
		if filter.TraceID != nil && (l.TraceID == nil || *l.TraceID != *filter.TraceID) {
			continue
		}
		if filter.FromTime != nil && l.CreatedAt.Before(*filter.FromTime) {
			continue
		}
		if filter.ToTime != nil && l.CreatedAt.After(*filter.ToTime) {
			continue
		}

		cloned := *l
		matched = append(matched, &cloned)
	}

	total := int64(len(matched))

	// Sort descending by created_at for browsing
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})

	start := filter.Offset
	if start > len(matched) {
		return []*AuditLog{}, total, nil
	}

	end := start + filter.Limit
	if end > len(matched) {
		end = len(matched)
	}

	return matched[start:end], total, nil
}

func (m *MockRepository) GetRecordsInRange(ctx context.Context, tenantID string, fromTime, toTime *time.Time) ([]*AuditLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*AuditLog

	for _, l := range m.logs {
		if l.TenantID != tenantID {
			continue
		}
		if fromTime != nil && l.CreatedAt.Before(*fromTime) {
			continue
		}
		if toTime != nil && l.CreatedAt.After(*toTime) {
			continue
		}

		cloned := *l
		result = append(result, &cloned)
	}

	// Order chronologically for chain verification
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}

func (m *MockRepository) CreateCheckpoint(ctx context.Context, cp *AuditVerificationCheckpoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cloned := *cp
	m.checkpoints = append(m.checkpoints, &cloned)
	return nil
}

func (m *MockRepository) ListCheckpoints(ctx context.Context, tenantID string, limit, offset int) ([]*AuditVerificationCheckpoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var tenantCheckpoints []*AuditVerificationCheckpoint
	for _, cp := range m.checkpoints {
		if cp.TenantID == tenantID {
			cloned := *cp
			tenantCheckpoints = append(tenantCheckpoints, &cloned)
		}
	}

	sort.Slice(tenantCheckpoints, func(i, j int) bool {
		return tenantCheckpoints[i].VerifiedAt.After(tenantCheckpoints[j].VerifiedAt)
	})

	if offset > len(tenantCheckpoints) {
		return []*AuditVerificationCheckpoint{}, nil
	}
	end := offset + limit
	if end > len(tenantCheckpoints) {
		end = len(tenantCheckpoints)
	}

	return tenantCheckpoints[offset:end], nil
}

// TamperWithRecord simulates a database tampering attack for testing tamper detection.
func (m *MockRepository) TamperWithRecord(recordID string, mutateFunc func(log *AuditLog)) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, l := range m.logs {
		if l.ID == recordID {
			mutateFunc(l)
			return true
		}
	}
	return false
}
