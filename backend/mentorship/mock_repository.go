/**
 * BLOCK_MENTORSHIP_MOCK_REPO_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   Thread-safe in-memory mock repository test doubles for mentorship entities.
 */

package mentorship

import (
	"context"
	"sync"
	"time"
)

type MockRepository struct {
	mu          sync.RWMutex
	allocations map[string]*MentorAllocation
	sessions    map[string]*MentorshipSession
	progress    map[string]*StudentAcademicProgress
	alerts      map[string]*MentorshipAtRiskAlert
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		allocations: make(map[string]*MentorAllocation),
		sessions:    make(map[string]*MentorshipSession),
		progress:    make(map[string]*StudentAcademicProgress),
		alerts:      make(map[string]*MentorshipAtRiskAlert),
	}
}

// --- Mentor Allocations ---

func (m *MockRepository) CreateAllocation(_ context.Context, alloc *MentorAllocation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.allocations[alloc.ID] = alloc
	return nil
}

func (m *MockRepository) GetAllocationByID(_ context.Context, tenantID, id string) (*MentorAllocation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	alloc, ok := m.allocations[id]
	if !ok || alloc.TenantID != tenantID {
		return nil, ErrAllocationNotFound
	}
	return alloc, nil
}

func (m *MockRepository) GetActiveAllocationByStudent(_ context.Context, tenantID, studentID string) (*MentorAllocation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, a := range m.allocations {
		if a.TenantID == tenantID && a.StudentID == studentID && a.Status == AllocationStatusActive {
			return a, nil
		}
	}
	return nil, ErrAllocationNotFound
}

func (m *MockRepository) ListAllocations(_ context.Context, tenantID string, mentorStaffID *string, status *AllocationStatus) ([]*MentorAllocation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*MentorAllocation
	for _, a := range m.allocations {
		if a.TenantID != tenantID {
			continue
		}
		if mentorStaffID != nil && a.MentorStaffID != *mentorStaffID {
			continue
		}
		if status != nil && a.Status != *status {
			continue
		}
		res = append(res, a)
	}
	return res, nil
}

func (m *MockRepository) UpdateAllocationStatus(_ context.Context, tenantID, id string, status AllocationStatus, completedAt *time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	alloc, ok := m.allocations[id]
	if !ok || alloc.TenantID != tenantID {
		return ErrAllocationNotFound
	}
	alloc.Status = status
	alloc.CompletedAt = completedAt
	alloc.UpdatedAt = time.Now().UTC()
	return nil
}

// --- Sessions ---

func (m *MockRepository) CreateSession(_ context.Context, sess *MentorshipSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[sess.ID] = sess
	return nil
}

func (m *MockRepository) GetSessionByID(_ context.Context, tenantID, id string) (*MentorshipSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sess, ok := m.sessions[id]
	if !ok || sess.TenantID != tenantID {
		return nil, ErrSessionNotFound
	}
	return sess, nil
}

func (m *MockRepository) ListSessionsByAllocation(_ context.Context, tenantID, allocationID string) ([]*MentorshipSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*MentorshipSession
	for _, s := range m.sessions {
		if s.TenantID == tenantID && s.AllocationID == allocationID {
			res = append(res, s)
		}
	}
	return res, nil
}

func (m *MockRepository) CompleteSession(_ context.Context, tenantID, id string, summary, actionItems string, followUpDate *time.Time, completedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[id]
	if !ok || sess.TenantID != tenantID {
		return ErrSessionNotFound
	}
	sess.Status = SessionStatusCompleted
	sess.DiscussionSummary = summary
	sess.ActionItems = actionItems
	sess.FollowUpDate = followUpDate
	sess.CompletedAt = &completedAt
	sess.UpdatedAt = completedAt
	return nil
}

// --- Progress ---

func (m *MockRepository) CreateProgress(_ context.Context, prog *StudentAcademicProgress) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.progress[prog.ID] = prog
	return nil
}

func (m *MockRepository) GetProgressByStudentSemester(_ context.Context, tenantID, studentID string, semester int) (*StudentAcademicProgress, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.progress {
		if p.TenantID == tenantID && p.StudentID == studentID && p.Semester == semester {
			return p, nil
		}
	}
	return nil, ErrProgressRecordNotFound
}

func (m *MockRepository) ListProgressByStudent(_ context.Context, tenantID, studentID string) ([]*StudentAcademicProgress, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*StudentAcademicProgress
	for _, p := range m.progress {
		if p.TenantID == tenantID && p.StudentID == studentID {
			res = append(res, p)
		}
	}
	return res, nil
}

// --- Alerts ---

func (m *MockRepository) CreateAlert(_ context.Context, alert *MentorshipAtRiskAlert) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alerts[alert.ID] = alert
	return nil
}

func (m *MockRepository) GetAlertByID(_ context.Context, tenantID, id string) (*MentorshipAtRiskAlert, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	alert, ok := m.alerts[id]
	if !ok || alert.TenantID != tenantID {
		return nil, ErrAlertNotFound
	}
	return alert, nil
}

func (m *MockRepository) ListAlerts(_ context.Context, tenantID string, studentID *string, resolved *bool) ([]*MentorshipAtRiskAlert, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*MentorshipAtRiskAlert
	for _, a := range m.alerts {
		if a.TenantID != tenantID {
			continue
		}
		if studentID != nil && a.StudentID != *studentID {
			continue
		}
		if resolved != nil {
			isResolved := a.ResolvedAt != nil
			if isResolved != *resolved {
				continue
			}
		}
		res = append(res, a)
	}
	return res, nil
}

func (m *MockRepository) ResolveAlert(_ context.Context, tenantID, id string, resolvedByID string, resolvedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	alert, ok := m.alerts[id]
	if !ok || alert.TenantID != tenantID {
		return ErrAlertNotFound
	}
	alert.ResolvedAt = &resolvedAt
	alert.ResolvedByID = &resolvedByID
	alert.UpdatedAt = resolvedAt
	return nil
}
