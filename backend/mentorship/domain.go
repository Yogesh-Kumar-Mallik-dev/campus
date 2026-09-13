/**
 * BLOCK_MENTORSHIP_DOMAIN_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   Domain entities, value objects, and business invariant rules for mentor allocation, counseling sessions, and at-risk alerts.
 */

package mentorship

import (
	"time"
)

type AllocationStatus string

const (
	AllocationStatusActive     AllocationStatus = "ACTIVE"
	AllocationStatusCompleted  AllocationStatus = "COMPLETED"
	AllocationStatusReassigned AllocationStatus = "REASSIGNED"
)

type MeetingType string

const (
	MeetingTypeOneOnOne            MeetingType = "ONE_ON_ONE"
	MeetingTypeGroup               MeetingType = "GROUP"
	MeetingTypeAcademicReview      MeetingType = "ACADEMIC_REVIEW"
	MeetingTypeEmergencyCounseling MeetingType = "EMERGENCY_COUNSELING"
)

type SessionStatus string

const (
	SessionStatusScheduled SessionStatus = "SCHEDULED"
	SessionStatusCompleted SessionStatus = "COMPLETED"
	SessionStatusCancelled SessionStatus = "CANCELLED"
	SessionStatusNoShow    SessionStatus = "NO_SHOW"
)

type AtRiskStatus string

const (
	AtRiskStatusNormal               AtRiskStatus = "NORMAL"
	AtRiskStatusWatchlist            AtRiskStatus = "WATCHLIST"
	AtRiskStatusCriticalIntervention AtRiskStatus = "CRITICAL_INTERVENTION"
)

type AtRiskType string

const (
	AtRiskTypeAcademicProbation  AtRiskType = "ACADEMIC_PROBATION"
	AtRiskTypeAttendanceShortage AtRiskType = "ATTENDANCE_SHORTAGE"
	AtRiskTypeDisciplinary       AtRiskType = "DISCIPLINARY"
	AtRiskTypeMentalHealth       AtRiskType = "MENTAL_HEALTH"
)

type RiskSeverity string

const (
	RiskSeverityLow      RiskSeverity = "LOW"
	RiskSeverityMedium   RiskSeverity = "MEDIUM"
	RiskSeverityHigh     RiskSeverity = "HIGH"
	RiskSeverityCritical RiskSeverity = "CRITICAL"
)

type MentorAllocation struct {
	ID            string           `json:"id"`
	TenantID      string           `json:"tenantId"`
	StudentID     string           `json:"studentId"`
	MentorStaffID string           `json:"mentorStaffId"`
	CohortID      *string          `json:"cohortId,omitempty"`
	Status        AllocationStatus `json:"status"`
	AllocatedAt   time.Time        `json:"allocatedAt"`
	CompletedAt   *time.Time       `json:"completedAt,omitempty"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

type MentorshipSession struct {
	ID                string        `json:"id"`
	TenantID          string        `json:"tenantId"`
	AllocationID      string        `json:"allocationId"`
	ScheduledAt       time.Time     `json:"scheduledAt"`
	CompletedAt       *time.Time    `json:"completedAt,omitempty"`
	Location          string        `json:"location,omitempty"`
	MeetingType       MeetingType   `json:"meetingType"`
	Status            SessionStatus `json:"status"`
	DiscussionSummary string        `json:"discussionSummary,omitempty"`
	ActionItems       string        `json:"actionItems,omitempty"`
	FollowUpDate      *time.Time    `json:"followUpDate,omitempty"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

type StudentAcademicProgress struct {
	ID                   string       `json:"id"`
	TenantID             string       `json:"tenantId"`
	StudentID            string       `json:"studentId"`
	Semester             int          `json:"semester"`
	SGPA                 float64      `json:"sgpa"`
	CGPA                 float64      `json:"cgpa"`
	AttendancePercentage float64      `json:"attendancePercentage"`
	CreditsEarned        int          `json:"creditsEarned"`
	TotalCredits         int          `json:"totalCredits"`
	AtRiskStatus         AtRiskStatus `json:"atRiskStatus"`
	Remarks              string       `json:"remarks,omitempty"`
	EvaluatedAt          time.Time    `json:"evaluatedAt"`
	CreatedAt            time.Time    `json:"createdAt"`
	UpdatedAt            time.Time    `json:"updatedAt"`
}

type MentorshipAtRiskAlert struct {
	ID           string       `json:"id"`
	TenantID     string       `json:"tenantId"`
	StudentID    string       `json:"studentId"`
	RiskType     AtRiskType   `json:"riskType"`
	Severity     RiskSeverity `json:"severity"`
	Description  string       `json:"description"`
	ResolvedAt   *time.Time   `json:"resolvedAt,omitempty"`
	ResolvedByID *string      `json:"resolvedById,omitempty"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

// EvaluateAtRiskStatus computes the academic risk status and identifies if an alert should trigger.
func EvaluateAtRiskStatus(sgpa, cgpa, attendance float64) (AtRiskStatus, *AtRiskType, RiskSeverity, error) {
	if sgpa < 0 || sgpa > 10.0 || cgpa < 0 || cgpa > 10.0 || attendance < 0 || attendance > 100.0 {
		return "", nil, "", ErrInvalidScoreBounds
	}

	// Critical Thresholds
	if attendance < 65.0 {
		riskType := AtRiskTypeAttendanceShortage
		return AtRiskStatusCriticalIntervention, &riskType, RiskSeverityCritical, nil
	}
	if sgpa < 4.0 {
		riskType := AtRiskTypeAcademicProbation
		return AtRiskStatusCriticalIntervention, &riskType, RiskSeverityCritical, nil
	}

	// Watchlist Thresholds
	if attendance < 75.0 {
		riskType := AtRiskTypeAttendanceShortage
		return AtRiskStatusWatchlist, &riskType, RiskSeverityMedium, nil
	}
	if sgpa < 5.0 {
		riskType := AtRiskTypeAcademicProbation
		return AtRiskStatusWatchlist, &riskType, RiskSeverityLow, nil
	}

	return AtRiskStatusNormal, nil, RiskSeverityLow, nil
}

// ValidateSessionCompletion verifies if a session can be marked completed.
func ValidateSessionCompletion(currentStatus SessionStatus) error {
	if currentStatus == SessionStatusCompleted {
		return ErrSessionAlreadyCompleted
	}
	if currentStatus != SessionStatusScheduled {
		return ErrInvalidSessionTransition
	}
	return nil
}
