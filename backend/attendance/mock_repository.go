/**
 * BLOCK_ATTENDANCE_MOCK_REPO_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   In-memory, concurrency-safe mock repositories for attendance unit testing.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package attendance

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

type MockCatalogRepository struct {
	mu       sync.RWMutex
	subjects map[string]*CourseSubject
	slots    map[string]*TimetableSlot
}

func NewMockCatalogRepository() *MockCatalogRepository {
	return &MockCatalogRepository{
		subjects: make(map[string]*CourseSubject),
		slots:    make(map[string]*TimetableSlot),
	}
}

func (m *MockCatalogRepository) CreateSubject(ctx context.Context, subject *CourseSubject) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subjects[subject.TenantID+":"+subject.ID] = subject
	return nil
}

func (m *MockCatalogRepository) GetSubjectByID(ctx context.Context, tenantID, id string) (*CourseSubject, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sub, ok := m.subjects[tenantID+":"+id]
	if !ok {
		return nil, ErrSubjectNotFound
	}
	c := *sub
	return &c, nil
}

func (m *MockCatalogRepository) GetSubjectByCode(ctx context.Context, tenantID, programID, code string) (*CourseSubject, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.subjects {
		if s.TenantID == tenantID && s.ProgramID == programID && strings.EqualFold(s.Code, code) {
			c := *s
			return &c, nil
		}
	}
	return nil, ErrSubjectNotFound
}

func (m *MockCatalogRepository) ListSubjects(ctx context.Context, tenantID, programID string, semester int) ([]*CourseSubject, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*CourseSubject
	for _, s := range m.subjects {
		if s.TenantID == tenantID {
			if programID != "" && s.ProgramID != programID {
				continue
			}
			if semester > 0 && s.Semester != semester {
				continue
			}
			c := *s
			res = append(res, &c)
		}
	}
	return res, nil
}

func (m *MockCatalogRepository) CreateSlot(ctx context.Context, slot *TimetableSlot) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.slots[slot.TenantID+":"+slot.ID] = slot
	return nil
}

func (m *MockCatalogRepository) GetSlotByID(ctx context.Context, tenantID, id string) (*TimetableSlot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	slot, ok := m.slots[tenantID+":"+id]
	if !ok {
		return nil, ErrSlotNotFound
	}
	c := *slot
	return &c, nil
}

func (m *MockCatalogRepository) ListSlotsByCohort(ctx context.Context, tenantID, cohortID string) ([]*TimetableSlot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*TimetableSlot
	for _, s := range m.slots {
		if s.TenantID == tenantID && s.CohortID == cohortID {
			c := *s
			res = append(res, &c)
		}
	}
	return res, nil
}

func (m *MockCatalogRepository) ListSlotsByFaculty(ctx context.Context, tenantID, facultyID string) ([]*TimetableSlot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*TimetableSlot
	for _, s := range m.slots {
		if s.TenantID == tenantID && s.FacultyID == facultyID {
			c := *s
			res = append(res, &c)
		}
	}
	return res, nil
}

type MockSessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]*AttendanceSession
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions: make(map[string]*AttendanceSession),
	}
}

func (m *MockSessionRepository) CreateSession(ctx context.Context, session *AttendanceSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[session.TenantID+":"+session.ID] = session
	return nil
}

func (m *MockSessionRepository) GetSessionByID(ctx context.Context, tenantID, id string) (*AttendanceSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[tenantID+":"+id]
	if !ok {
		return nil, ErrSessionNotFound
	}
	c := *s
	return &c, nil
}

func (m *MockSessionRepository) UpdateSession(ctx context.Context, session *AttendanceSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := session.TenantID + ":" + session.ID
	if _, ok := m.sessions[key]; !ok {
		return ErrSessionNotFound
	}
	m.sessions[key] = session
	return nil
}

func (m *MockSessionRepository) ListSessions(ctx context.Context, filter SessionFilter) ([]*AttendanceSession, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*AttendanceSession
	for _, s := range m.sessions {
		if filter.TenantID != "" && s.TenantID != filter.TenantID {
			continue
		}
		if filter.SubjectID != "" && s.SubjectID != filter.SubjectID {
			continue
		}
		if filter.CohortID != "" && s.CohortID != filter.CohortID {
			continue
		}
		if filter.FacultyID != "" && s.FacultyID != filter.FacultyID {
			continue
		}
		if filter.SessionDate != "" && s.SessionDate != filter.SessionDate {
			continue
		}
		if filter.Status != nil && s.Status != *filter.Status {
			continue
		}
		c := *s
		results = append(results, &c)
	}

	total := len(results)
	if filter.Offset >= total {
		return []*AttendanceSession{}, total, nil
	}
	end := total
	if filter.Limit > 0 && filter.Offset+filter.Limit < end {
		end = filter.Offset + filter.Limit
	}
	return results[filter.Offset:end], total, nil
}

type MockRecordRepository struct {
	mu      sync.RWMutex
	records map[string]*AttendanceRecord // Key: sessionID + ":" + studentID
}

func NewMockRecordRepository() *MockRecordRepository {
	return &MockRecordRepository{
		records: make(map[string]*AttendanceRecord),
	}
}

func (m *MockRecordRepository) BatchUpsertRecords(ctx context.Context, records []AttendanceRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range records {
		recCopy := r
		key := r.SessionID + ":" + r.StudentID
		m.records[key] = &recCopy
	}
	return nil
}

func (m *MockRecordRepository) GetRecordsBySession(ctx context.Context, tenantID, sessionID string) ([]AttendanceRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []AttendanceRecord
	for _, r := range m.records {
		if r.TenantID == tenantID && r.SessionID == sessionID {
			res = append(res, *r)
		}
	}
	return res, nil
}

func (m *MockRecordRepository) GetRecordsByStudent(ctx context.Context, tenantID, studentID string) ([]AttendanceRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []AttendanceRecord
	for _, r := range m.records {
		if r.TenantID == tenantID && r.StudentID == studentID {
			res = append(res, *r)
		}
	}
	return res, nil
}

func (m *MockRecordRepository) GetSubjectRecordsForStudent(ctx context.Context, tenantID, studentID, subjectID string) ([]AttendanceRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []AttendanceRecord
	for _, r := range m.records {
		if r.TenantID == tenantID && r.StudentID == studentID {
			res = append(res, *r)
		}
	}
	return res, nil
}

func (m *MockRecordRepository) MarkExcusedMedical(ctx context.Context, tenantID, studentID, fromDate, toDate string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for _, r := range m.records {
		if r.TenantID == tenantID && r.StudentID == studentID {
			r.Status = RecordExcusedMedical
			r.Remarks = fmt.Sprintf("Medical Leave Condonation (%s to %s)", fromDate, toDate)
			count++
		}
	}
	return count, nil
}

type MockMedicalLeaveRepository struct {
	mu     sync.RWMutex
	leaves map[string]*MedicalLeaveApplication
}

func NewMockMedicalLeaveRepository() *MockMedicalLeaveRepository {
	return &MockMedicalLeaveRepository{
		leaves: make(map[string]*MedicalLeaveApplication),
	}
}

func (m *MockMedicalLeaveRepository) CreateLeave(ctx context.Context, leave *MedicalLeaveApplication) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.leaves[leave.TenantID+":"+leave.ID] = leave
	return nil
}

func (m *MockMedicalLeaveRepository) GetLeaveByID(ctx context.Context, tenantID, id string) (*MedicalLeaveApplication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.leaves[tenantID+":"+id]
	if !ok {
		return nil, ErrMedicalLeaveNotFound
	}
	c := *l
	return &c, nil
}

func (m *MockMedicalLeaveRepository) UpdateLeaveStatus(ctx context.Context, tenantID, id string, status MedicalLeaveStatus, approvedByID, rejectionReason string, approvedAt *time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := tenantID + ":" + id
	l, ok := m.leaves[key]
	if !ok {
		return ErrMedicalLeaveNotFound
	}
	l.Status = status
	l.ApprovedByID = approvedByID
	l.ApprovedAt = approvedAt
	l.RejectionReason = rejectionReason
	l.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockMedicalLeaveRepository) ListLeavesByStudent(ctx context.Context, tenantID, studentID string) ([]*MedicalLeaveApplication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*MedicalLeaveApplication
	for _, l := range m.leaves {
		if l.TenantID == tenantID && l.StudentID == studentID {
			c := *l
			res = append(res, &c)
		}
	}
	return res, nil
}

func (m *MockMedicalLeaveRepository) ListPendingLeaves(ctx context.Context, tenantID string) ([]*MedicalLeaveApplication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*MedicalLeaveApplication
	for _, l := range m.leaves {
		if l.TenantID == tenantID && l.Status == MedicalLeavePending {
			c := *l
			res = append(res, &c)
		}
	}
	return res, nil
}
