/**
 * BLOCK_STUDYHUB_SERVICE_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   Core business logic service and audit dispatching for courses, study materials, assignments, submissions, and peer grading.
 */

package studyhub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"campus/backend/audit"
)

type Service interface {
	CreateCourse(ctx context.Context, req CreateCourseRequest) (*StudyCourse, error)
	GetCourse(ctx context.Context, tenantID, id string) (*StudyCourse, error)
	ListCourses(ctx context.Context, tenantID string, department *string, semester *int) ([]*StudyCourse, error)

	PublishMaterial(ctx context.Context, req PublishMaterialRequest) (*StudyMaterial, error)
	ListMaterials(ctx context.Context, tenantID, courseID string, unitNumber *int) ([]*StudyMaterial, error)

	CreateAssignment(ctx context.Context, req CreateAssignmentRequest) (*StudyAssignment, error)
	GetAssignment(ctx context.Context, tenantID, id string) (*StudyAssignment, error)
	ListAssignments(ctx context.Context, tenantID, courseID string) ([]*StudyAssignment, error)

	SubmitAssignment(ctx context.Context, req SubmitAssignmentRequest) (*StudySubmission, error)
	GetSubmission(ctx context.Context, tenantID, id string) (*StudySubmission, error)
	ListSubmissions(ctx context.Context, tenantID, assignmentID string) ([]*StudySubmission, error)
	GradeSubmission(ctx context.Context, req GradeSubmissionRequest) (*StudySubmission, error)

	SubmitPeerReview(ctx context.Context, req SubmitPeerReviewRequest) (*StudyPeerReview, error)
	ListPeerReviews(ctx context.Context, submissionID string) ([]*StudyPeerReview, error)
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

func (s *service) emitAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
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

type CreateCourseRequest struct {
	TenantID     string  `json:"tenantId"`
	Code         string  `json:"code"`
	Title        string  `json:"title"`
	Description  string  `json:"description,omitempty"`
	Department   string  `json:"department"`
	Credits      int     `json:"credits"`
	Semester     int     `json:"semester"`
	InstructorID *string `json:"instructorId,omitempty"`
	SyllabusText string  `json:"syllabusText,omitempty"`
}

func (s *service) CreateCourse(ctx context.Context, req CreateCourseRequest) (*StudyCourse, error) {
	normCode := NormalizeCourseCode(req.Code)
	if existing, _ := s.repo.GetCourseByCode(ctx, req.TenantID, normCode); existing != nil {
		return nil, ErrCourseCodeExists
	}

	credits := req.Credits
	if credits <= 0 {
		credits = 4
	}
	sem := req.Semester
	if sem <= 0 {
		sem = 1
	}

	now := time.Now().UTC()
	course := &StudyCourse{
		ID:           fmt.Sprintf("crs_%d", now.UnixNano()),
		TenantID:     req.TenantID,
		Code:         normCode,
		Title:        req.Title,
		Description:  req.Description,
		Department:   req.Department,
		Credits:      credits,
		Semester:     sem,
		InstructorID: req.InstructorID,
		SyllabusText: req.SyllabusText,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateCourse(ctx, course); err != nil {
		return nil, err
	}

	s.emitAudit(course.TenantID, "system", "USER", "studyhub:course:created", "study_course", course.ID, audit.StatusSuccess, map[string]interface{}{
		"code": course.Code, "title": course.Title, "department": course.Department,
	})

	return course, nil
}

func (s *service) GetCourse(ctx context.Context, tenantID, id string) (*StudyCourse, error) {
	return s.repo.GetCourseByID(ctx, tenantID, id)
}

func (s *service) ListCourses(ctx context.Context, tenantID string, department *string, semester *int) ([]*StudyCourse, error) {
	return s.repo.ListCourses(ctx, tenantID, department, semester)
}

type PublishMaterialRequest struct {
	TenantID      string       `json:"tenantId"`
	CourseID      string       `json:"courseId"`
	Title         string       `json:"title"`
	Description   string       `json:"description,omitempty"`
	UnitNumber    int          `json:"unitNumber"`
	MaterialType  MaterialType `json:"materialType"`
	FileURL       string       `json:"fileUrl"`
	FileSizeBytes int64        `json:"fileSizeBytes"`
}

func (s *service) PublishMaterial(ctx context.Context, req PublishMaterialRequest) (*StudyMaterial, error) {
	if _, err := s.repo.GetCourseByID(ctx, req.TenantID, req.CourseID); err != nil {
		return nil, err
	}

	unit := req.UnitNumber
	if unit <= 0 {
		unit = 1
	}

	mType := req.MaterialType
	if mType == "" {
		mType = MaterialTypeNote
	}

	now := time.Now().UTC()
	mat := &StudyMaterial{
		ID:            fmt.Sprintf("mat_%d", now.UnixNano()),
		TenantID:      req.TenantID,
		CourseID:      req.CourseID,
		Title:         req.Title,
		Description:   req.Description,
		UnitNumber:    unit,
		MaterialType:  mType,
		FileURL:       req.FileURL,
		FileSizeBytes: req.FileSizeBytes,
		PublishedAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.CreateMaterial(ctx, mat); err != nil {
		return nil, err
	}

	s.emitAudit(mat.TenantID, "system", "USER", "studyhub:material:published", "study_material", mat.ID, audit.StatusSuccess, map[string]interface{}{
		"courseId": mat.CourseID, "title": mat.Title, "type": string(mat.MaterialType),
	})

	return mat, nil
}

func (s *service) ListMaterials(ctx context.Context, tenantID, courseID string, unitNumber *int) ([]*StudyMaterial, error) {
	return s.repo.ListMaterialsByCourse(ctx, tenantID, courseID, unitNumber)
}

type CreateAssignmentRequest struct {
	TenantID                 string    `json:"tenantId"`
	CourseID                 string    `json:"courseId"`
	Title                    string    `json:"title"`
	Description              string    `json:"description"`
	MaxMarks                 int       `json:"maxMarks"`
	DueDate                  time.Time `json:"dueDate"`
	AllowLateSubmission      bool      `json:"allowLateSubmission"`
	LatePenaltyPercentPerDay int       `json:"latePenaltyPercentPerDay"`
}

func (s *service) CreateAssignment(ctx context.Context, req CreateAssignmentRequest) (*StudyAssignment, error) {
	if _, err := s.repo.GetCourseByID(ctx, req.TenantID, req.CourseID); err != nil {
		return nil, err
	}

	maxMarks := req.MaxMarks
	if maxMarks <= 0 {
		maxMarks = 100
	}

	now := time.Now().UTC()
	asgn := &StudyAssignment{
		ID:                       fmt.Sprintf("asgn_%d", now.UnixNano()),
		TenantID:                 req.TenantID,
		CourseID:                 req.CourseID,
		Title:                    req.Title,
		Description:              req.Description,
		MaxMarks:                 maxMarks,
		DueDate:                  req.DueDate,
		AllowLateSubmission:      req.AllowLateSubmission,
		LatePenaltyPercentPerDay: req.LatePenaltyPercentPerDay,
		Status:                   AssignmentStatusPublished,
		CreatedAt:                now,
		UpdatedAt:                now,
	}

	if err := s.repo.CreateAssignment(ctx, asgn); err != nil {
		return nil, err
	}

	s.emitAudit(asgn.TenantID, "system", "USER", "studyhub:assignment:created", "study_assignment", asgn.ID, audit.StatusSuccess, map[string]interface{}{
		"courseId": asgn.CourseID, "title": asgn.Title, "dueDate": asgn.DueDate,
	})

	return asgn, nil
}

func (s *service) GetAssignment(ctx context.Context, tenantID, id string) (*StudyAssignment, error) {
	return s.repo.GetAssignmentByID(ctx, tenantID, id)
}

func (s *service) ListAssignments(ctx context.Context, tenantID, courseID string) ([]*StudyAssignment, error) {
	return s.repo.ListAssignmentsByCourse(ctx, tenantID, courseID)
}

type SubmitAssignmentRequest struct {
	TenantID       string     `json:"tenantId"`
	AssignmentID   string     `json:"assignmentId"`
	StudentID      string     `json:"studentId"`
	FileURL        string     `json:"fileUrl,omitempty"`
	ContentText    string     `json:"contentText,omitempty"`
	SubmissionTime *time.Time `json:"submissionTime,omitempty"`
}

func (s *service) SubmitAssignment(ctx context.Context, req SubmitAssignmentRequest) (*StudySubmission, error) {
	asgn, err := s.repo.GetAssignmentByID(ctx, req.TenantID, req.AssignmentID)
	if err != nil {
		return nil, err
	}

	if asgn.Status == AssignmentStatusClosed {
		return nil, ErrAssignmentClosed
	}

	if existing, _ := s.repo.GetSubmissionByStudent(ctx, req.TenantID, req.AssignmentID, req.StudentID); existing != nil {
		return nil, ErrDuplicateSubmission
	}

	now := time.Now().UTC()
	if req.SubmissionTime != nil {
		now = *req.SubmissionTime
	}

	status, err := EvaluateSubmissionTiming(asgn.DueDate, asgn.AllowLateSubmission, now)
	if err != nil {
		return nil, err
	}

	sub := &StudySubmission{
		ID:           fmt.Sprintf("sub_%d", now.UnixNano()),
		TenantID:     req.TenantID,
		AssignmentID: req.AssignmentID,
		StudentID:    req.StudentID,
		FileURL:      req.FileURL,
		ContentText:  req.ContentText,
		SubmittedAt:  now,
		Status:       status,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateSubmission(ctx, sub); err != nil {
		return nil, err
	}

	s.emitAudit(sub.TenantID, sub.StudentID, "STUDENT", "studyhub:submission:created", "study_submission", sub.ID, audit.StatusSuccess, map[string]interface{}{
		"assignmentId": sub.AssignmentID, "status": string(sub.Status),
	})

	return sub, nil
}

func (s *service) GetSubmission(ctx context.Context, tenantID, id string) (*StudySubmission, error) {
	return s.repo.GetSubmissionByID(ctx, tenantID, id)
}

func (s *service) ListSubmissions(ctx context.Context, tenantID, assignmentID string) ([]*StudySubmission, error) {
	return s.repo.ListSubmissionsByAssignment(ctx, tenantID, assignmentID)
}

type GradeSubmissionRequest struct {
	TenantID   string  `json:"tenantId"`
	ID         string  `json:"id"`
	RawMarks   float64 `json:"rawMarks"`
	Feedback   string  `json:"feedback"`
	GradedByID string  `json:"gradedById"`
}

func (s *service) GradeSubmission(ctx context.Context, req GradeSubmissionRequest) (*StudySubmission, error) {
	sub, err := s.repo.GetSubmissionByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}

	asgn, err := s.repo.GetAssignmentByID(ctx, req.TenantID, sub.AssignmentID)
	if err != nil {
		return nil, err
	}

	adjustedMarks, err := CalculateAdjustedMarks(req.RawMarks, asgn.MaxMarks, asgn.DueDate, sub.SubmittedAt, asgn.LatePenaltyPercentPerDay)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := s.repo.GradeSubmission(ctx, req.TenantID, req.ID, adjustedMarks, req.Feedback, req.GradedByID, now); err != nil {
		return nil, err
	}

	sub.Status = SubmissionStatusGraded
	sub.MarksObtained = &adjustedMarks
	sub.Feedback = req.Feedback
	sub.GradedByID = &req.GradedByID
	sub.GradedAt = &now
	sub.UpdatedAt = now

	s.emitAudit(sub.TenantID, req.GradedByID, "FACULTY", "studyhub:submission:graded", "study_submission", sub.ID, audit.StatusSuccess, map[string]interface{}{
		"rawMarks": req.RawMarks, "adjustedMarks": adjustedMarks, "maxMarks": asgn.MaxMarks,
	})

	return sub, nil
}

type SubmitPeerReviewRequest struct {
	TenantID          string  `json:"tenantId"`
	SubmissionID      string  `json:"submissionId"`
	ReviewerStudentID string  `json:"reviewerStudentId"`
	Score             float64 `json:"score"`
	Comments          string  `json:"comments"`
}

func (s *service) SubmitPeerReview(ctx context.Context, req SubmitPeerReviewRequest) (*StudyPeerReview, error) {
	sub, err := s.repo.GetSubmissionByID(ctx, req.TenantID, req.SubmissionID)
	if err != nil {
		return nil, err
	}

	if sub.StudentID == req.ReviewerStudentID {
		return nil, ErrSelfPeerReviewNotAllowed
	}

	if req.Score < 0 || req.Score > 100 {
		return nil, ErrInvalidScoreRange
	}

	now := time.Now().UTC()
	review := &StudyPeerReview{
		ID:                fmt.Sprintf("prev_%d", now.UnixNano()),
		SubmissionID:      req.SubmissionID,
		ReviewerStudentID: req.ReviewerStudentID,
		Score:             req.Score,
		Comments:          req.Comments,
		ReviewedAt:        now,
	}

	if err := s.repo.CreatePeerReview(ctx, review); err != nil {
		return nil, err
	}

	s.emitAudit(req.TenantID, req.ReviewerStudentID, "STUDENT", "studyhub:peer_review:submitted", "study_peer_review", review.ID, audit.StatusSuccess, map[string]interface{}{
		"submissionId": review.SubmissionID, "score": review.Score,
	})

	return review, nil
}

func (s *service) ListPeerReviews(ctx context.Context, submissionID string) ([]*StudyPeerReview, error) {
	return s.repo.ListPeerReviewsBySubmission(ctx, submissionID)
}
