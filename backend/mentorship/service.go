/**
 * BLOCK_MENTORSHIP_SERVICE_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   Core business logic service and audit dispatching for mentor allocations, counseling sessions, progress tracking, and risk interventions.
 */

package mentorship

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"campus/backend/audit"
)

type Service interface {
	AllocateMentor(ctx context.Context, req AllocateMentorRequest) (*MentorAllocation, error)
	GetActiveAllocation(ctx context.Context, tenantID, studentID string) (*MentorAllocation, error)
	ListAllocations(ctx context.Context, tenantID string, mentorStaffID *string, status *AllocationStatus) ([]*MentorAllocation, error)

	ScheduleSession(ctx context.Context, req ScheduleSessionRequest) (*MentorshipSession, error)
	CompleteSession(ctx context.Context, req CompleteSessionRequest) (*MentorshipSession, error)
	ListSessions(ctx context.Context, tenantID, allocationID string) ([]*MentorshipSession, error)

	RecordProgress(ctx context.Context, req RecordProgressRequest) (*StudentAcademicProgress, error)
	ListProgress(ctx context.Context, tenantID, studentID string) ([]*StudentAcademicProgress, error)

	ResolveAlert(ctx context.Context, req ResolveAlertRequest) (*MentorshipAtRiskAlert, error)
	ListAlerts(ctx context.Context, tenantID string, studentID *string, resolved *bool) ([]*MentorshipAtRiskAlert, error)
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

type AllocateMentorRequest struct {
	TenantID      string  `json:"tenantId"`
	StudentID     string  `json:"studentId"`
	MentorStaffID string  `json:"mentorStaffId"`
	CohortID      *string `json:"cohortId,omitempty"`
}

func (s *service) AllocateMentor(ctx context.Context, req AllocateMentorRequest) (*MentorAllocation, error) {
	if existing, _ := s.repo.GetActiveAllocationByStudent(ctx, req.TenantID, req.StudentID); existing != nil {
		return nil, ErrActiveAllocationExists
	}

	now := time.Now().UTC()
	alloc := &MentorAllocation{
		ID:            fmt.Sprintf("alloc_%d", now.UnixNano()),
		TenantID:      req.TenantID,
		StudentID:     req.StudentID,
		MentorStaffID: req.MentorStaffID,
		CohortID:      req.CohortID,
		Status:        AllocationStatusActive,
		AllocatedAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.CreateAllocation(ctx, alloc); err != nil {
		return nil, err
	}

	s.emitAudit(alloc.TenantID, req.MentorStaffID, "STAFF", "mentorship:mentor:allocated", "mentor_allocation", alloc.ID, audit.StatusSuccess, map[string]interface{}{
		"studentId": alloc.StudentID, "mentorStaffId": alloc.MentorStaffID,
	})

	return alloc, nil
}

func (s *service) GetActiveAllocation(ctx context.Context, tenantID, studentID string) (*MentorAllocation, error) {
	return s.repo.GetActiveAllocationByStudent(ctx, tenantID, studentID)
}

func (s *service) ListAllocations(ctx context.Context, tenantID string, mentorStaffID *string, status *AllocationStatus) ([]*MentorAllocation, error) {
	return s.repo.ListAllocations(ctx, tenantID, mentorStaffID, status)
}

type ScheduleSessionRequest struct {
	TenantID     string      `json:"tenantId"`
	AllocationID string      `json:"allocationId"`
	ScheduledAt  time.Time   `json:"scheduledAt"`
	Location     string      `json:"location,omitempty"`
	MeetingType  MeetingType `json:"meetingType"`
}

func (s *service) ScheduleSession(ctx context.Context, req ScheduleSessionRequest) (*MentorshipSession, error) {
	if _, err := s.repo.GetAllocationByID(ctx, req.TenantID, req.AllocationID); err != nil {
		return nil, err
	}

	mType := req.MeetingType
	if mType == "" {
		mType = MeetingTypeOneOnOne
	}

	now := time.Now().UTC()
	sess := &MentorshipSession{
		ID:           fmt.Sprintf("sess_%d", now.UnixNano()),
		TenantID:     req.TenantID,
		AllocationID: req.AllocationID,
		ScheduledAt:  req.ScheduledAt,
		Location:     req.Location,
		MeetingType:  mType,
		Status:       SessionStatusScheduled,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	s.emitAudit(sess.TenantID, "system", "STAFF", "mentorship:session:scheduled", "mentorship_session", sess.ID, audit.StatusSuccess, map[string]interface{}{
		"allocationId": sess.AllocationID, "scheduledAt": sess.ScheduledAt, "meetingType": string(sess.MeetingType),
	})

	return sess, nil
}

type CompleteSessionRequest struct {
	TenantID          string     `json:"tenantId"`
	ID                string     `json:"id"`
	DiscussionSummary string     `json:"discussionSummary"`
	ActionItems       string     `json:"actionItems"`
	FollowUpDate      *time.Time `json:"followUpDate,omitempty"`
}

func (s *service) CompleteSession(ctx context.Context, req CompleteSessionRequest) (*MentorshipSession, error) {
	sess, err := s.repo.GetSessionByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}

	if err := ValidateSessionCompletion(sess.Status); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := s.repo.CompleteSession(ctx, req.TenantID, req.ID, req.DiscussionSummary, req.ActionItems, req.FollowUpDate, now); err != nil {
		return nil, err
	}

	sess.Status = SessionStatusCompleted
	sess.DiscussionSummary = req.DiscussionSummary
	sess.ActionItems = req.ActionItems
	sess.FollowUpDate = req.FollowUpDate
	sess.CompletedAt = &now
	sess.UpdatedAt = now

	s.emitAudit(sess.TenantID, "system", "STAFF", "mentorship:session:completed", "mentorship_session", sess.ID, audit.StatusSuccess, map[string]interface{}{
		"allocationId": sess.AllocationID, "completedAt": now,
	})

	return sess, nil
}

func (s *service) ListSessions(ctx context.Context, tenantID, allocationID string) ([]*MentorshipSession, error) {
	return s.repo.ListSessionsByAllocation(ctx, tenantID, allocationID)
}

type RecordProgressRequest struct {
	TenantID             string  `json:"tenantId"`
	StudentID            string  `json:"studentId"`
	Semester             int     `json:"semester"`
	SGPA                 float64 `json:"sgpa"`
	CGPA                 float64 `json:"cgpa"`
	AttendancePercentage float64 `json:"attendancePercentage"`
	CreditsEarned        int     `json:"creditsEarned"`
	TotalCredits         int     `json:"totalCredits"`
	Remarks              string  `json:"remarks,omitempty"`
}

func (s *service) RecordProgress(ctx context.Context, req RecordProgressRequest) (*StudentAcademicProgress, error) {
	if existing, _ := s.repo.GetProgressByStudentSemester(ctx, req.TenantID, req.StudentID, req.Semester); existing != nil {
		return nil, ErrProgressRecordExists
	}

	riskStatus, riskType, severity, err := EvaluateAtRiskStatus(req.SGPA, req.CGPA, req.AttendancePercentage)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	prog := &StudentAcademicProgress{
		ID:                   fmt.Sprintf("prog_%d", now.UnixNano()),
		TenantID:             req.TenantID,
		StudentID:            req.StudentID,
		Semester:             req.Semester,
		SGPA:                 req.SGPA,
		CGPA:                 req.CGPA,
		AttendancePercentage: req.AttendancePercentage,
		CreditsEarned:        req.CreditsEarned,
		TotalCredits:         req.TotalCredits,
		AtRiskStatus:         riskStatus,
		Remarks:              req.Remarks,
		EvaluatedAt:          now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := s.repo.CreateProgress(ctx, prog); err != nil {
		return nil, err
	}

	s.emitAudit(prog.TenantID, "system", "SYSTEM", "mentorship:progress:evaluated", "student_academic_progress", prog.ID, audit.StatusSuccess, map[string]interface{}{
		"studentId": prog.StudentID, "semester": prog.Semester, "sgpa": prog.SGPA, "atRiskStatus": string(prog.AtRiskStatus),
	})

	// Generate automated At-Risk Alert if flagged
	if riskStatus != AtRiskStatusNormal && riskType != nil {
		desc := fmt.Sprintf("Automated Flag: SGPA %.2f and Attendance %.1f%% in Semester %d requires intervention.", req.SGPA, req.AttendancePercentage, req.Semester)
		alert := &MentorshipAtRiskAlert{
			ID:          fmt.Sprintf("alrt_%d", now.UnixNano()),
			TenantID:    req.TenantID,
			StudentID:   req.StudentID,
			RiskType:    *riskType,
			Severity:    severity,
			Description: desc,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		_ = s.repo.CreateAlert(ctx, alert)

		s.emitAudit(alert.TenantID, "system", "SYSTEM", "mentorship:risk_alert:flagged", "mentorship_at_risk_alert", alert.ID, audit.StatusSuccess, map[string]interface{}{
			"studentId": alert.StudentID, "riskType": string(alert.RiskType), "severity": string(alert.Severity),
		})
	}

	return prog, nil
}

func (s *service) ListProgress(ctx context.Context, tenantID, studentID string) ([]*StudentAcademicProgress, error) {
	return s.repo.ListProgressByStudent(ctx, tenantID, studentID)
}

type ResolveAlertRequest struct {
	TenantID     string `json:"tenantId"`
	ID           string `json:"id"`
	ResolvedByID string `json:"resolvedById"`
}

func (s *service) ResolveAlert(ctx context.Context, req ResolveAlertRequest) (*MentorshipAtRiskAlert, error) {
	alert, err := s.repo.GetAlertByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}

	if alert.ResolvedAt != nil {
		return nil, ErrAlertAlreadyResolved
	}

	now := time.Now().UTC()
	if err := s.repo.ResolveAlert(ctx, req.TenantID, req.ID, req.ResolvedByID, now); err != nil {
		return nil, err
	}

	alert.ResolvedAt = &now
	alert.ResolvedByID = &req.ResolvedByID
	alert.UpdatedAt = now

	s.emitAudit(alert.TenantID, req.ResolvedByID, "STAFF", "mentorship:risk_alert:resolved", "mentorship_at_risk_alert", alert.ID, audit.StatusSuccess, map[string]interface{}{
		"alertId": alert.ID, "studentId": alert.StudentID,
	})

	return alert, nil
}

func (s *service) ListAlerts(ctx context.Context, tenantID string, studentID *string, resolved *bool) ([]*MentorshipAtRiskAlert, error) {
	return s.repo.ListAlerts(ctx, tenantID, studentID, resolved)
}
