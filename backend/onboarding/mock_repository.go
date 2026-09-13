/**
 * BLOCK_ONBOARDING_MOCK_REPO_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   In-memory, concurrency-safe mock repositories for onboarding unit and integration testing.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

type MockApplicantRepository struct {
	mu         sync.RWMutex
	applicants map[string]*Applicant // Key: tenantID + ":" + id
}

func NewMockApplicantRepository() *MockApplicantRepository {
	return &MockApplicantRepository{
		applicants: make(map[string]*Applicant),
	}
}

func (m *MockApplicantRepository) Create(ctx context.Context, applicant *Applicant) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := applicant.TenantID + ":" + applicant.ID
	m.applicants[key] = applicant
	return nil
}

func (m *MockApplicantRepository) GetByID(ctx context.Context, tenantID, id string) (*Applicant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := tenantID + ":" + id
	app, ok := m.applicants[key]
	if !ok {
		return nil, ErrApplicantNotFound
	}
	// Return a copy
	c := *app
	return &c, nil
}

func (m *MockApplicantRepository) GetByEmail(ctx context.Context, tenantID, email string) (*Applicant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clean := strings.ToLower(strings.TrimSpace(email))
	for _, app := range m.applicants {
		if app.TenantID == tenantID && strings.ToLower(app.Email) == clean {
			c := *app
			return &c, nil
		}
	}
	return nil, ErrApplicantNotFound
}

func (m *MockApplicantRepository) Update(ctx context.Context, applicant *Applicant) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := applicant.TenantID + ":" + applicant.ID
	if _, ok := m.applicants[key]; !ok {
		return ErrApplicantNotFound
	}
	m.applicants[key] = applicant
	return nil
}

func (m *MockApplicantRepository) List(ctx context.Context, filter ApplicantFilter) ([]*Applicant, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*Applicant
	for _, app := range m.applicants {
		if filter.TenantID != "" && app.TenantID != filter.TenantID {
			continue
		}
		if filter.Type != nil && app.Type != *filter.Type {
			continue
		}
		if filter.Status != nil && app.Status != *filter.Status {
			continue
		}
		if filter.ProgramID != "" && app.ProgramID != filter.ProgramID {
			continue
		}
		if filter.DepartmentID != "" && app.DepartmentID != filter.DepartmentID {
			continue
		}
		if filter.AcademicYear != "" && app.AcademicYear != filter.AcademicYear {
			continue
		}
		if filter.Search != "" {
			s := strings.ToLower(filter.Search)
			if !strings.Contains(strings.ToLower(app.FirstName), s) &&
				!strings.Contains(strings.ToLower(app.LastName), s) &&
				!strings.Contains(strings.ToLower(app.Email), s) {
				continue
			}
		}
		c := *app
		results = append(results, &c)
	}

	total := len(results)
	if filter.Offset >= total {
		return []*Applicant{}, total, nil
	}
	end := total
	if filter.Limit > 0 && filter.Offset+filter.Limit < end {
		end = filter.Offset + filter.Limit
	}

	return results[filter.Offset:end], total, nil
}

type MockDocumentRepository struct {
	mu        sync.RWMutex
	documents map[string]*Document // Key: tenantID + ":" + id
}

func NewMockDocumentRepository() *MockDocumentRepository {
	return &MockDocumentRepository{
		documents: make(map[string]*Document),
	}
}

func (m *MockDocumentRepository) Attach(ctx context.Context, doc *Document) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := doc.TenantID + ":" + doc.ID
	m.documents[key] = doc
	return nil
}

func (m *MockDocumentRepository) GetByID(ctx context.Context, tenantID, id string) (*Document, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := tenantID + ":" + id
	doc, ok := m.documents[key]
	if !ok {
		return nil, ErrDocumentNotFound
	}
	c := *doc
	return &c, nil
}

func (m *MockDocumentRepository) GetByApplicant(ctx context.Context, tenantID, applicantID string) ([]Document, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []Document
	for _, doc := range m.documents {
		if doc.TenantID == tenantID && doc.ApplicantID == applicantID {
			results = append(results, *doc)
		}
	}
	return results, nil
}

func (m *MockDocumentRepository) UpdateStatus(ctx context.Context, tenantID, docID string, status DocumentStatus, reason, verifierID string, verifiedAt *time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := tenantID + ":" + docID
	doc, ok := m.documents[key]
	if !ok {
		return ErrDocumentNotFound
	}
	doc.Status = status
	doc.RejectionReason = reason
	doc.VerifiedByID = verifierID
	doc.VerifiedAt = verifiedAt
	doc.UpdatedAt = time.Now().UTC()
	return nil
}

type MockAcademicRepository struct {
	mu          sync.RWMutex
	departments map[string]*AcademicDepartment
	programs    map[string]*AcademicProgram
	cohorts     map[string]*CohortBatch
}

func NewMockAcademicRepository() *MockAcademicRepository {
	return &MockAcademicRepository{
		departments: make(map[string]*AcademicDepartment),
		programs:    make(map[string]*AcademicProgram),
		cohorts:     make(map[string]*CohortBatch),
	}
}

func (m *MockAcademicRepository) CreateDepartment(ctx context.Context, dept *AcademicDepartment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.departments[dept.TenantID+":"+dept.ID] = dept
	return nil
}

func (m *MockAcademicRepository) GetDepartmentByID(ctx context.Context, tenantID, id string) (*AcademicDepartment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.departments[tenantID+":"+id]
	if !ok {
		return nil, ErrDepartmentNotFound
	}
	c := *d
	return &c, nil
}

func (m *MockAcademicRepository) GetDepartmentByCode(ctx context.Context, tenantID, code string) (*AcademicDepartment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, d := range m.departments {
		if d.TenantID == tenantID && strings.EqualFold(d.Code, code) {
			c := *d
			return &c, nil
		}
	}
	return nil, ErrDepartmentNotFound
}

func (m *MockAcademicRepository) ListDepartments(ctx context.Context, tenantID string) ([]*AcademicDepartment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*AcademicDepartment
	for _, d := range m.departments {
		if d.TenantID == tenantID {
			c := *d
			res = append(res, &c)
		}
	}
	return res, nil
}

func (m *MockAcademicRepository) CreateProgram(ctx context.Context, program *AcademicProgram) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.programs[program.TenantID+":"+program.ID] = program
	return nil
}

func (m *MockAcademicRepository) GetProgramByID(ctx context.Context, tenantID, id string) (*AcademicProgram, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.programs[tenantID+":"+id]
	if !ok {
		return nil, ErrProgramNotFound
	}
	c := *p
	return &c, nil
}

func (m *MockAcademicRepository) GetProgramByCode(ctx context.Context, tenantID, code string) (*AcademicProgram, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.programs {
		if p.TenantID == tenantID && strings.EqualFold(p.Code, code) {
			c := *p
			return &c, nil
		}
	}
	return nil, ErrProgramNotFound
}

func (m *MockAcademicRepository) ListPrograms(ctx context.Context, tenantID, departmentID string) ([]*AcademicProgram, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*AcademicProgram
	for _, p := range m.programs {
		if p.TenantID == tenantID {
			if departmentID != "" && p.DepartmentID != departmentID {
				continue
			}
			c := *p
			res = append(res, &c)
		}
	}
	return res, nil
}

func (m *MockAcademicRepository) CreateCohort(ctx context.Context, cohort *CohortBatch) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cohorts[cohort.TenantID+":"+cohort.ID] = cohort
	return nil
}

func (m *MockAcademicRepository) GetCohortByID(ctx context.Context, tenantID, id string) (*CohortBatch, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.cohorts[tenantID+":"+id]
	if !ok {
		return nil, ErrCohortNotFound
	}
	res := *c
	return &res, nil
}

func (m *MockAcademicRepository) ListCohorts(ctx context.Context, tenantID, programID, academicYear string) ([]*CohortBatch, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*CohortBatch
	for _, c := range m.cohorts {
		if c.TenantID == tenantID {
			if programID != "" && c.ProgramID != programID {
				continue
			}
			if academicYear != "" && c.AcademicYear != academicYear {
				continue
			}
			cp := *c
			res = append(res, &cp)
		}
	}
	return res, nil
}

func (m *MockAcademicRepository) IncrementCohortEnrolled(ctx context.Context, tenantID, cohortID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.cohorts[tenantID+":"+cohortID]
	if !ok {
		return ErrCohortNotFound
	}
	c.CurrentEnrolled++
	return nil
}

type MockProfileRepository struct {
	mu       sync.RWMutex
	students map[string]*StudentProfile
	staff    map[string]*StaffProfile
}

func NewMockProfileRepository() *MockProfileRepository {
	return &MockProfileRepository{
		students: make(map[string]*StudentProfile),
		staff:    make(map[string]*StaffProfile),
	}
}

func (m *MockProfileRepository) CreateStudentProfile(ctx context.Context, profile *StudentProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.students {
		if s.TenantID == profile.TenantID && s.RollNumber == profile.RollNumber {
			return ErrDuplicateRollNumber
		}
	}
	m.students[profile.TenantID+":"+profile.ID] = profile
	return nil
}

func (m *MockProfileRepository) GetStudentProfileByRollNumber(ctx context.Context, tenantID, rollNumber string) (*StudentProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.students {
		if s.TenantID == tenantID && s.RollNumber == rollNumber {
			c := *s
			return &c, nil
		}
	}
	return nil, fmt.Errorf("student profile not found")
}

func (m *MockProfileRepository) GetStudentProfileByApplicantID(ctx context.Context, tenantID, applicantID string) (*StudentProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.students {
		if s.TenantID == tenantID && s.ApplicantID == applicantID {
			c := *s
			return &c, nil
		}
	}
	return nil, fmt.Errorf("student profile not found")
}

func (m *MockProfileRepository) CreateStaffProfile(ctx context.Context, profile *StaffProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.staff {
		if s.TenantID == profile.TenantID && s.EmployeeID == profile.EmployeeID {
			return ErrDuplicateEmployeeID
		}
	}
	m.staff[profile.TenantID+":"+profile.ID] = profile
	return nil
}

func (m *MockProfileRepository) GetStaffProfileByEmployeeID(ctx context.Context, tenantID, employeeID string) (*StaffProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.staff {
		if s.TenantID == tenantID && s.EmployeeID == employeeID {
			c := *s
			return &c, nil
		}
	}
	return nil, fmt.Errorf("staff profile not found")
}

func (m *MockProfileRepository) GetStaffProfileByApplicantID(ctx context.Context, tenantID, applicantID string) (*StaffProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.staff {
		if s.TenantID == tenantID && s.ApplicantID == applicantID {
			c := *s
			return &c, nil
		}
	}
	return nil, fmt.Errorf("staff profile not found")
}

type MockSequenceRepository struct {
	mu       sync.Mutex
	counters map[string]int
}

func NewMockSequenceRepository() *MockSequenceRepository {
	return &MockSequenceRepository{
		counters: make(map[string]int),
	}
}

func (m *MockSequenceRepository) NextSequence(ctx context.Context, tenantID, sequenceType, prefix string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s:%s:%s", tenantID, sequenceType, prefix)
	m.counters[key]++
	return m.counters[key], nil
}
