/**
 * BLOCK_MENTORSHIP_TEST_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   Unit test suite covering mentor allocation, sessions, academic progress, and automated at-risk alerts.
 */

package mentorship_test

import (
	"context"
	"testing"
	"time"

	"campus/backend/mentorship"
)

func setupMentorshipService() mentorship.Service {
	repo := mentorship.NewMockRepository()
	return mentorship.NewService(repo, nil)
}

func TestMentorship_AllocationLifecycleAndSingleActiveGuard(t *testing.T) {
	ctx := context.Background()
	svc := setupMentorshipService()

	tenantID := "tenant-alpha"
	studentID := "stu-001"
	mentorID := "staff-prof-rao"

	// 1. Allocate Mentor
	alloc, err := svc.AllocateMentor(ctx, mentorship.AllocateMentorRequest{
		TenantID:      tenantID,
		StudentID:     studentID,
		MentorStaffID: mentorID,
	})
	if err != nil {
		t.Fatalf("unexpected error allocating mentor: %v", err)
	}
	if alloc.Status != mentorship.AllocationStatusActive {
		t.Fatalf("expected status ACTIVE, got %s", alloc.Status)
	}

	// 2. Duplicate Active Allocation Guard
	_, err = svc.AllocateMentor(ctx, mentorship.AllocateMentorRequest{
		TenantID:      tenantID,
		StudentID:     studentID,
		MentorStaffID: "staff-prof-verma",
	})
	if err != mentorship.ErrActiveAllocationExists {
		t.Fatalf("expected ErrActiveAllocationExists, got %v", err)
	}

	// 3. Get Active Allocation
	fetched, err := svc.GetActiveAllocation(ctx, tenantID, studentID)
	if err != nil || fetched.ID != alloc.ID {
		t.Fatalf("failed to get active allocation: %v", err)
	}

	// 4. List Allocations by Mentor
	list, err := svc.ListAllocations(ctx, tenantID, &mentorID, nil)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 allocation in list, got %d, err: %v", len(list), err)
	}
}

func TestMentorship_SessionSchedulingAndCompletion(t *testing.T) {
	ctx := context.Background()
	svc := setupMentorshipService()

	tenantID := "tenant-alpha"

	alloc, _ := svc.AllocateMentor(ctx, mentorship.AllocateMentorRequest{
		TenantID:      tenantID,
		StudentID:     "stu-002",
		MentorStaffID: "staff-prof-rao",
	})

	schedTime := time.Now().UTC().Add(24 * time.Hour)

	// 1. Schedule Session
	sess, err := svc.ScheduleSession(ctx, mentorship.ScheduleSessionRequest{
		TenantID:     tenantID,
		AllocationID: alloc.ID,
		ScheduledAt:  schedTime,
		Location:     "Room 304, Academic Block A",
		MeetingType:  mentorship.MeetingTypeAcademicReview,
	})
	if err != nil {
		t.Fatalf("unexpected error scheduling session: %v", err)
	}
	if sess.Status != mentorship.SessionStatusScheduled {
		t.Fatalf("expected status SCHEDULED, got %s", sess.Status)
	}

	// 2. Complete Session
	followUp := time.Now().UTC().Add(14 * 24 * time.Hour)
	completedSess, err := svc.CompleteSession(ctx, mentorship.CompleteSessionRequest{
		TenantID:          tenantID,
		ID:                sess.ID,
		DiscussionSummary: "Reviewed mid-term results and discussed study strategy for Algorithms.",
		ActionItems:       "Submit weekly practice sets; attend peer tutoring lab on Thursdays.",
		FollowUpDate:      &followUp,
	})
	if err != nil {
		t.Fatalf("unexpected error completing session: %v", err)
	}
	if completedSess.Status != mentorship.SessionStatusCompleted {
		t.Fatalf("expected status COMPLETED, got %s", completedSess.Status)
	}

	// 3. Duplicate Completion Guard
	_, err = svc.CompleteSession(ctx, mentorship.CompleteSessionRequest{
		TenantID: tenantID,
		ID:       sess.ID,
	})
	if err != mentorship.ErrSessionAlreadyCompleted {
		t.Fatalf("expected ErrSessionAlreadyCompleted, got %v", err)
	}
}

func TestMentorship_AcademicProgressEvaluationNormal(t *testing.T) {
	ctx := context.Background()
	svc := setupMentorshipService()

	tenantID := "tenant-alpha"
	studentID := "stu-003"

	prog, err := svc.RecordProgress(ctx, mentorship.RecordProgressRequest{
		TenantID:             tenantID,
		StudentID:            studentID,
		Semester:             3,
		SGPA:                 8.6,
		CGPA:                 8.4,
		AttendancePercentage: 92.5,
		CreditsEarned:        22,
		TotalCredits:         22,
		Remarks:              "Consistent high performance.",
	})
	if err != nil {
		t.Fatalf("failed to record progress: %v", err)
	}
	if prog.AtRiskStatus != mentorship.AtRiskStatusNormal {
		t.Fatalf("expected status NORMAL, got %s", prog.AtRiskStatus)
	}

	// Verify no alerts created
	alerts, err := svc.ListAlerts(ctx, tenantID, &studentID, nil)
	if err != nil || len(alerts) != 0 {
		t.Fatalf("expected 0 alerts for normal student, got %d", len(alerts))
	}
}

func TestMentorship_AcademicProgressAtRiskAutomatedAlert(t *testing.T) {
	ctx := context.Background()
	svc := setupMentorshipService()

	tenantID := "tenant-alpha"
	studentID := "stu-004"

	// 1. Record Progress with low SGPA and low attendance -> Critical Intervention
	prog, err := svc.RecordProgress(ctx, mentorship.RecordProgressRequest{
		TenantID:             tenantID,
		StudentID:            studentID,
		Semester:             4,
		SGPA:                 3.8,
		CGPA:                 4.5,
		AttendancePercentage: 61.0,
		CreditsEarned:        12,
		TotalCredits:         22,
		Remarks:              "Severe attendance drop and backlog in 3 core subjects.",
	})
	if err != nil {
		t.Fatalf("failed to record progress: %v", err)
	}
	if prog.AtRiskStatus != mentorship.AtRiskStatusCriticalIntervention {
		t.Fatalf("expected CRITICAL_INTERVENTION, got %s", prog.AtRiskStatus)
	}

	// 2. Check Automated Alert
	alerts, err := svc.ListAlerts(ctx, tenantID, &studentID, nil)
	if err != nil || len(alerts) != 1 {
		t.Fatalf("expected 1 automated alert, got %d, err: %v", len(alerts), err)
	}
	alert := alerts[0]
	if alert.Severity != mentorship.RiskSeverityCritical {
		t.Fatalf("expected CRITICAL severity, got %s", alert.Severity)
	}

	// 3. Resolve Alert
	resolvedAlert, err := svc.ResolveAlert(ctx, mentorship.ResolveAlertRequest{
		TenantID:     tenantID,
		ID:           alert.ID,
		ResolvedByID: "staff-dean-student-affairs",
	})
	if err != nil {
		t.Fatalf("failed to resolve alert: %v", err)
	}
	if resolvedAlert.ResolvedAt == nil {
		t.Fatalf("expected resolvedAt timestamp to be set")
	}

	// 4. Duplicate Resolve Guard
	_, err = svc.ResolveAlert(ctx, mentorship.ResolveAlertRequest{
		TenantID:     tenantID,
		ID:           alert.ID,
		ResolvedByID: "staff-dean-student-affairs",
	})
	if err != mentorship.ErrAlertAlreadyResolved {
		t.Fatalf("expected ErrAlertAlreadyResolved, got %v", err)
	}
}

func TestMentorship_ScoreBoundsGuard(t *testing.T) {
	ctx := context.Background()
	svc := setupMentorshipService()

	tenantID := "tenant-alpha"

	// SGPA > 10
	_, err := svc.RecordProgress(ctx, mentorship.RecordProgressRequest{
		TenantID:             tenantID,
		StudentID:            "stu-005",
		Semester:             1,
		SGPA:                 11.5,
		CGPA:                 8.0,
		AttendancePercentage: 80.0,
	})
	if err != mentorship.ErrInvalidScoreBounds {
		t.Fatalf("expected ErrInvalidScoreBounds, got %v", err)
	}

	// Attendance > 100
	_, err = svc.RecordProgress(ctx, mentorship.RecordProgressRequest{
		TenantID:             tenantID,
		StudentID:            "stu-005",
		Semester:             1,
		SGPA:                 8.0,
		CGPA:                 8.0,
		AttendancePercentage: 105.0,
	})
	if err != mentorship.ErrInvalidScoreBounds {
		t.Fatalf("expected ErrInvalidScoreBounds, got %v", err)
	}
}
