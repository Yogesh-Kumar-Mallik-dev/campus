/**
 * BLOCK_MENTORSHIP_REPOSITORY_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   Repository interface contracts for mentor allocations, counseling sessions, academic progress, and at-risk alerts.
 */

package mentorship

import (
	"context"
	"time"
)

type Repository interface {
	// Mentor Allocations
	CreateAllocation(ctx context.Context, alloc *MentorAllocation) error
	GetAllocationByID(ctx context.Context, tenantID, id string) (*MentorAllocation, error)
	GetActiveAllocationByStudent(ctx context.Context, tenantID, studentID string) (*MentorAllocation, error)
	ListAllocations(ctx context.Context, tenantID string, mentorStaffID *string, status *AllocationStatus) ([]*MentorAllocation, error)
	UpdateAllocationStatus(ctx context.Context, tenantID, id string, status AllocationStatus, completedAt *time.Time) error

	// Counseling Sessions
	CreateSession(ctx context.Context, sess *MentorshipSession) error
	GetSessionByID(ctx context.Context, tenantID, id string) (*MentorshipSession, error)
	ListSessionsByAllocation(ctx context.Context, tenantID, allocationID string) ([]*MentorshipSession, error)
	CompleteSession(ctx context.Context, tenantID, id string, summary, actionItems string, followUpDate *time.Time, completedAt time.Time) error

	// Academic Progress
	CreateProgress(ctx context.Context, prog *StudentAcademicProgress) error
	GetProgressByStudentSemester(ctx context.Context, tenantID, studentID string, semester int) (*StudentAcademicProgress, error)
	ListProgressByStudent(ctx context.Context, tenantID, studentID string) ([]*StudentAcademicProgress, error)

	// At-Risk Alerts
	CreateAlert(ctx context.Context, alert *MentorshipAtRiskAlert) error
	GetAlertByID(ctx context.Context, tenantID, id string) (*MentorshipAtRiskAlert, error)
	ListAlerts(ctx context.Context, tenantID string, studentID *string, resolved *bool) ([]*MentorshipAtRiskAlert, error)
	ResolveAlert(ctx context.Context, tenantID, id string, resolvedByID string, resolvedAt time.Time) error
}
