/**
 * BLOCK_STUDYHUB_DOMAIN_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   Domain entities, value objects, and business invariant evaluation for courses, materials, assignments, and peer grading.
 */

package studyhub

import (
	"math"
	"strings"
	"time"
)

type MaterialType string

const (
	MaterialTypeNote        MaterialType = "NOTE"
	MaterialTypeSlide       MaterialType = "SLIDE"
	MaterialTypeLabManual   MaterialType = "LAB_MANUAL"
	MaterialTypeRecording   MaterialType = "RECORDING"
	MaterialTypeSamplePaper MaterialType = "SAMPLE_PAPER"
)

type AssignmentStatus string

const (
	AssignmentStatusDraft     AssignmentStatus = "DRAFT"
	AssignmentStatusPublished AssignmentStatus = "PUBLISHED"
	AssignmentStatusClosed    AssignmentStatus = "CLOSED"
)

type SubmissionStatus string

const (
	SubmissionStatusSubmitted SubmissionStatus = "SUBMITTED"
	SubmissionStatusGraded    SubmissionStatus = "GRADED"
	SubmissionStatusLate      SubmissionStatus = "LATE"
	SubmissionStatusRejected  SubmissionStatus = "REJECTED"
)

type StudyCourse struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenantId"`
	Code         string    `json:"code"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Department   string    `json:"department"`
	Credits      int       `json:"credits"`
	Semester     int       `json:"semester"`
	InstructorID *string   `json:"instructorId,omitempty"`
	SyllabusText string    `json:"syllabusText,omitempty"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type StudyMaterial struct {
	ID            string       `json:"id"`
	TenantID      string       `json:"tenantId"`
	CourseID      string       `json:"courseId"`
	Title         string       `json:"title"`
	Description   string       `json:"description,omitempty"`
	UnitNumber    int          `json:"unitNumber"`
	MaterialType  MaterialType `json:"materialType"`
	FileURL       string       `json:"fileUrl"`
	FileSizeBytes int64        `json:"fileSizeBytes"`
	PublishedAt   time.Time    `json:"publishedAt"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

type StudyAssignment struct {
	ID                       string           `json:"id"`
	TenantID                 string           `json:"tenantId"`
	CourseID                 string           `json:"courseId"`
	Title                    string           `json:"title"`
	Description              string           `json:"description"`
	MaxMarks                 int              `json:"maxMarks"`
	DueDate                  time.Time        `json:"dueDate"`
	AllowLateSubmission      bool             `json:"allowLateSubmission"`
	LatePenaltyPercentPerDay int              `json:"latePenaltyPercentPerDay"`
	Status                   AssignmentStatus `json:"status"`
	CreatedAt                time.Time        `json:"createdAt"`
	UpdatedAt                time.Time        `json:"updatedAt"`
}

type StudySubmission struct {
	ID            string           `json:"id"`
	TenantID      string           `json:"tenantId"`
	AssignmentID  string           `json:"assignmentId"`
	StudentID     string           `json:"studentId"`
	FileURL       string           `json:"fileUrl,omitempty"`
	ContentText   string           `json:"contentText,omitempty"`
	SubmittedAt   time.Time        `json:"submittedAt"`
	Status        SubmissionStatus `json:"status"`
	MarksObtained *float64         `json:"marksObtained,omitempty"`
	Feedback      string           `json:"feedback,omitempty"`
	GradedByID    *string          `json:"gradedById,omitempty"`
	GradedAt      *time.Time       `json:"gradedAt,omitempty"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

type StudyPeerReview struct {
	ID                string    `json:"id"`
	SubmissionID      string    `json:"submissionId"`
	ReviewerStudentID string    `json:"reviewerStudentId"`
	Score             float64   `json:"score"`
	Comments          string    `json:"comments"`
	ReviewedAt        time.Time `json:"reviewedAt"`
}

// EvaluateSubmissionTiming checks due date policy and determines if submission is on-time, late, or disallowed.
func EvaluateSubmissionTiming(dueDate time.Time, allowLate bool, now time.Time) (SubmissionStatus, error) {
	if now.Before(dueDate) || now.Equal(dueDate) {
		return SubmissionStatusSubmitted, nil
	}

	if !allowLate {
		return SubmissionStatusRejected, ErrLateSubmissionDisallowed
	}

	return SubmissionStatusLate, nil
}

// CalculateAdjustedMarks validates marks and applies late penalty deduction if applicable.
func CalculateAdjustedMarks(rawMarks float64, maxMarks int, dueDate, submittedAt time.Time, penaltyPercentPerDay int) (float64, error) {
	if rawMarks < 0 || rawMarks > float64(maxMarks) {
		return 0, ErrMarksExceedMaxMarks
	}

	if submittedAt.Before(dueDate) || submittedAt.Equal(dueDate) || penaltyPercentPerDay <= 0 {
		return rawMarks, nil
	}

	// Calculate full 24-hour periods late (at least 1 day if after due date)
	durationLate := submittedAt.Sub(dueDate)
	if durationLate <= 0 {
		return rawMarks, nil
	}
	hoursLate := durationLate.Hours()
	daysLate := int(math.Ceil(hoursLate/24.0 - 1e-4))
	if daysLate < 1 {
		daysLate = 1
	}

	penaltyRate := float64(daysLate*penaltyPercentPerDay) / 100.0
	if penaltyRate > 1.0 {
		penaltyRate = 1.0
	}

	deduction := float64(maxMarks) * penaltyRate
	adjusted := rawMarks - deduction
	if adjusted < 0 {
		adjusted = 0
	}

	return adjusted, nil
}

// NormalizeCourseCode standardizes course codes to uppercase trimmed string.
func NormalizeCourseCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
