/**
 * BLOCK_STUDYHUB_REPOSITORY_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   Repository interface contracts for courses, study materials, assignments, submissions, and peer reviews.
 */

package studyhub

import (
	"context"
	"time"
)

type Repository interface {
	// Courses
	CreateCourse(ctx context.Context, course *StudyCourse) error
	GetCourseByID(ctx context.Context, tenantID, id string) (*StudyCourse, error)
	GetCourseByCode(ctx context.Context, tenantID, code string) (*StudyCourse, error)
	ListCourses(ctx context.Context, tenantID string, department *string, semester *int) ([]*StudyCourse, error)

	// Materials
	CreateMaterial(ctx context.Context, mat *StudyMaterial) error
	GetMaterialByID(ctx context.Context, tenantID, id string) (*StudyMaterial, error)
	ListMaterialsByCourse(ctx context.Context, tenantID, courseID string, unitNumber *int) ([]*StudyMaterial, error)

	// Assignments
	CreateAssignment(ctx context.Context, asgn *StudyAssignment) error
	GetAssignmentByID(ctx context.Context, tenantID, id string) (*StudyAssignment, error)
	ListAssignmentsByCourse(ctx context.Context, tenantID, courseID string) ([]*StudyAssignment, error)
	UpdateAssignmentStatus(ctx context.Context, tenantID, id string, status AssignmentStatus) error

	// Submissions
	CreateSubmission(ctx context.Context, sub *StudySubmission) error
	GetSubmissionByID(ctx context.Context, tenantID, id string) (*StudySubmission, error)
	GetSubmissionByStudent(ctx context.Context, tenantID, assignmentID, studentID string) (*StudySubmission, error)
	ListSubmissionsByAssignment(ctx context.Context, tenantID, assignmentID string) ([]*StudySubmission, error)
	GradeSubmission(ctx context.Context, tenantID, id string, marks float64, feedback string, gradedByID string, gradedAt time.Time) error

	// Peer Reviews
	CreatePeerReview(ctx context.Context, review *StudyPeerReview) error
	ListPeerReviewsBySubmission(ctx context.Context, submissionID string) ([]*StudyPeerReview, error)
}
