/**
 * BLOCK_STUDYHUB_MOCK_REPO_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   Thread-safe in-memory mock repository test doubles for study hub entities.
 */

package studyhub

import (
	"context"
	"sync"
	"time"
)

type MockRepository struct {
	mu          sync.RWMutex
	courses     map[string]*StudyCourse
	materials   map[string]*StudyMaterial
	assignments map[string]*StudyAssignment
	submissions map[string]*StudySubmission
	peerReviews map[string]*StudyPeerReview
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		courses:     make(map[string]*StudyCourse),
		materials:   make(map[string]*StudyMaterial),
		assignments: make(map[string]*StudyAssignment),
		submissions: make(map[string]*StudySubmission),
		peerReviews: make(map[string]*StudyPeerReview),
	}
}

// --- Courses ---

func (m *MockRepository) CreateCourse(_ context.Context, course *StudyCourse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.courses[course.ID] = course
	return nil
}

func (m *MockRepository) GetCourseByID(_ context.Context, tenantID, id string) (*StudyCourse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	course, ok := m.courses[id]
	if !ok || course.TenantID != tenantID {
		return nil, ErrCourseNotFound
	}
	return course, nil
}

func (m *MockRepository) GetCourseByCode(_ context.Context, tenantID, code string) (*StudyCourse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, c := range m.courses {
		if c.TenantID == tenantID && c.Code == code {
			return c, nil
		}
	}
	return nil, ErrCourseNotFound
}

func (m *MockRepository) ListCourses(_ context.Context, tenantID string, department *string, semester *int) ([]*StudyCourse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*StudyCourse
	for _, c := range m.courses {
		if c.TenantID != tenantID {
			continue
		}
		if department != nil && c.Department != *department {
			continue
		}
		if semester != nil && c.Semester != *semester {
			continue
		}
		res = append(res, c)
	}
	return res, nil
}

// --- Materials ---

func (m *MockRepository) CreateMaterial(_ context.Context, mat *StudyMaterial) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.materials[mat.ID] = mat
	return nil
}

func (m *MockRepository) GetMaterialByID(_ context.Context, tenantID, id string) (*StudyMaterial, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mat, ok := m.materials[id]
	if !ok || mat.TenantID != tenantID {
		return nil, ErrMaterialNotFound
	}
	return mat, nil
}

func (m *MockRepository) ListMaterialsByCourse(_ context.Context, tenantID, courseID string, unitNumber *int) ([]*StudyMaterial, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*StudyMaterial
	for _, mat := range m.materials {
		if mat.TenantID != tenantID || mat.CourseID != courseID {
			continue
		}
		if unitNumber != nil && mat.UnitNumber != *unitNumber {
			continue
		}
		res = append(res, mat)
	}
	return res, nil
}

// --- Assignments ---

func (m *MockRepository) CreateAssignment(_ context.Context, asgn *StudyAssignment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.assignments[asgn.ID] = asgn
	return nil
}

func (m *MockRepository) GetAssignmentByID(_ context.Context, tenantID, id string) (*StudyAssignment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	asgn, ok := m.assignments[id]
	if !ok || asgn.TenantID != tenantID {
		return nil, ErrAssignmentNotFound
	}
	return asgn, nil
}

func (m *MockRepository) ListAssignmentsByCourse(_ context.Context, tenantID, courseID string) ([]*StudyAssignment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*StudyAssignment
	for _, asgn := range m.assignments {
		if asgn.TenantID == tenantID && asgn.CourseID == courseID {
			res = append(res, asgn)
		}
	}
	return res, nil
}

func (m *MockRepository) UpdateAssignmentStatus(_ context.Context, tenantID, id string, status AssignmentStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	asgn, ok := m.assignments[id]
	if !ok || asgn.TenantID != tenantID {
		return ErrAssignmentNotFound
	}
	asgn.Status = status
	asgn.UpdatedAt = time.Now().UTC()
	return nil
}

// --- Submissions ---

func (m *MockRepository) CreateSubmission(_ context.Context, sub *StudySubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.submissions[sub.ID] = sub
	return nil
}

func (m *MockRepository) GetSubmissionByID(_ context.Context, tenantID, id string) (*StudySubmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sub, ok := m.submissions[id]
	if !ok || sub.TenantID != tenantID {
		return nil, ErrSubmissionNotFound
	}
	return sub, nil
}

func (m *MockRepository) GetSubmissionByStudent(_ context.Context, tenantID, assignmentID, studentID string) (*StudySubmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.submissions {
		if s.TenantID == tenantID && s.AssignmentID == assignmentID && s.StudentID == studentID {
			return s, nil
		}
	}
	return nil, ErrSubmissionNotFound
}

func (m *MockRepository) ListSubmissionsByAssignment(_ context.Context, tenantID, assignmentID string) ([]*StudySubmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*StudySubmission
	for _, s := range m.submissions {
		if s.TenantID == tenantID && s.AssignmentID == assignmentID {
			res = append(res, s)
		}
	}
	return res, nil
}

func (m *MockRepository) GradeSubmission(_ context.Context, tenantID, id string, marks float64, feedback string, gradedByID string, gradedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sub, ok := m.submissions[id]
	if !ok || sub.TenantID != tenantID {
		return ErrSubmissionNotFound
	}
	sub.MarksObtained = &marks
	sub.Feedback = feedback
	sub.GradedByID = &gradedByID
	sub.GradedAt = &gradedAt
	sub.Status = SubmissionStatusGraded
	sub.UpdatedAt = gradedAt
	return nil
}

// --- Peer Reviews ---

func (m *MockRepository) CreatePeerReview(_ context.Context, review *StudyPeerReview) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.peerReviews {
		if r.SubmissionID == review.SubmissionID && r.ReviewerStudentID == review.ReviewerStudentID {
			return ErrDuplicatePeerReview
		}
	}
	m.peerReviews[review.ID] = review
	return nil
}

func (m *MockRepository) ListPeerReviewsBySubmission(_ context.Context, submissionID string) ([]*StudyPeerReview, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*StudyPeerReview
	for _, r := range m.peerReviews {
		if r.SubmissionID == submissionID {
			res = append(res, r)
		}
	}
	return res, nil
}
