/**
 * BLOCK_ATTENDANCE_DOMAIN_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   Domain entities, timetable slots, attendance sessions, shortage metrics, and medical condonation.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package attendance

import (
	"strings"
	"time"
)

type AttendanceSessionStatus string

const (
	SessionScheduled AttendanceSessionStatus = "SCHEDULED"
	SessionOpen      AttendanceSessionStatus = "OPEN"
	SessionLocked    AttendanceSessionStatus = "LOCKED"
	SessionFinalized AttendanceSessionStatus = "FINALIZED"
)

type AttendanceMode string

const (
	ModeManualFaculty     AttendanceMode = "MANUAL_FACULTY"
	ModeBiometricTerminal AttendanceMode = "BIOMETRIC_TERMINAL"
	ModeRFIDScan          AttendanceMode = "RFID_SCAN"
	ModeGeofenceMobile    AttendanceMode = "GEOFENCE_MOBILE"
)

type AttendanceRecordStatus string

const (
	RecordPresent        AttendanceRecordStatus = "PRESENT"
	RecordAbsent         AttendanceRecordStatus = "ABSENT"
	RecordLate           AttendanceRecordStatus = "LATE"
	RecordExcusedMedical AttendanceRecordStatus = "EXCUSED_MEDICAL"
)

type MedicalLeaveStatus string

const (
	MedicalLeavePending  MedicalLeaveStatus = "PENDING"
	MedicalLeaveApproved MedicalLeaveStatus = "APPROVED"
	MedicalLeaveRejected MedicalLeaveStatus = "REJECTED"
)

type DayOfWeek string

const (
	DayMonday    DayOfWeek = "MONDAY"
	DayTuesday   DayOfWeek = "TUESDAY"
	DayWednesday DayOfWeek = "WEDNESDAY"
	DayThursday  DayOfWeek = "THURSDAY"
	DayFriday    DayOfWeek = "FRIDAY"
	DaySaturday  DayOfWeek = "SATURDAY"
	DaySunday    DayOfWeek = "SUNDAY"
)

type CourseSubject struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	ProgramID string    `json:"program_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Credits   int       `json:"credits"`
	Semester  int       `json:"semester"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TimetableSlot struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	SubjectID string    `json:"subject_id"`
	CohortID  string    `json:"cohort_id"`
	FacultyID string    `json:"faculty_id"`
	DayOfWeek DayOfWeek `json:"day_of_week"`
	StartTime string    `json:"start_time"` // "09:00"
	EndTime   string    `json:"end_time"`   // "10:00"
	RoomCode  string    `json:"room_code"`  // "LH-101"
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AttendanceRecord struct {
	ID        string                 `json:"id"`
	TenantID  string                 `json:"tenant_id"`
	SessionID string                 `json:"session_id"`
	StudentID string                 `json:"student_id"`
	Status    AttendanceRecordStatus `json:"status"`
	Remarks   string                 `json:"remarks,omitempty"`
	MarkedAt  time.Time              `json:"marked_at"`
	DeviceID  string                 `json:"device_id,omitempty"`
	Latitude  *float64               `json:"latitude,omitempty"`
	Longitude *float64               `json:"longitude,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type AttendanceSession struct {
	ID            string                  `json:"id"`
	TenantID      string                  `json:"tenant_id"`
	SlotID        string                  `json:"slot_id,omitempty"`
	SubjectID     string                  `json:"subject_id"`
	FacultyID     string                  `json:"faculty_id"`
	CohortID      string                  `json:"cohort_id"`
	SessionDate   string                  `json:"session_date"` // YYYY-MM-DD
	StartTime     string                  `json:"start_time"`
	EndTime       string                  `json:"end_time"`
	Mode          AttendanceMode          `json:"mode"`
	Status        AttendanceSessionStatus `json:"status"`
	OpenedAt      *time.Time              `json:"opened_at,omitempty"`
	LockedAt      *time.Time              `json:"locked_at,omitempty"`
	TotalStudents int                     `json:"total_students"`
	PresentCount  int                     `json:"present_count"`
	AbsentCount   int                     `json:"absent_count"`
	Records       []AttendanceRecord      `json:"records,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

type MedicalLeaveApplication struct {
	ID              string             `json:"id"`
	TenantID        string             `json:"tenant_id"`
	StudentID       string             `json:"student_id"`
	FromDate        string             `json:"from_date"` // YYYY-MM-DD
	ToDate          string             `json:"to_date"`   // YYYY-MM-DD
	Reason          string             `json:"reason"`
	CertificateKey  string             `json:"certificate_key"`
	Status          MedicalLeaveStatus `json:"status"`
	ApprovedByID    string             `json:"approved_by_id,omitempty"`
	ApprovedAt      *time.Time         `json:"approved_at,omitempty"`
	RejectionReason string             `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type StudentAttendanceSummary struct {
	StudentID            string  `json:"student_id"`
	SubjectID            string  `json:"subject_id"`
	TotalSessions        int     `json:"total_sessions"`
	PresentSessions      int     `json:"present_sessions"`
	ExcusedSessions      int     `json:"excused_sessions"`
	AbsentSessions       int     `json:"absent_sessions"`
	AttendancePercentage float64 `json:"attendance_percentage"`
	IsShortage           bool    `json:"is_shortage"` // True when < 75.0%
}

// Open transitions SCHEDULED -> OPEN.
func (s *AttendanceSession) Open() error {
	if s.Status != SessionScheduled {
		return NewDomainError("INVALID_STATUS", "only scheduled sessions can be opened", ErrInvalidStateTransition)
	}
	s.Status = SessionOpen
	now := time.Now().UTC()
	s.OpenedAt = &now
	s.UpdatedAt = now
	return nil
}

// Lock transitions OPEN -> LOCKED.
func (s *AttendanceSession) Lock() error {
	if s.Status != SessionOpen {
		return NewDomainError("INVALID_STATUS", "only open sessions can be locked", ErrInvalidStateTransition)
	}
	s.Status = SessionLocked
	now := time.Now().UTC()
	s.LockedAt = &now
	s.UpdatedAt = now
	return nil
}

// Finalize transitions LOCKED -> FINALIZED and calculates final counts.
func (s *AttendanceSession) Finalize(records []AttendanceRecord) error {
	if s.Status != SessionLocked && s.Status != SessionOpen && s.Status != SessionScheduled {
		return NewDomainError("INVALID_STATUS", "session must be scheduled, open, or locked to finalize", ErrInvalidStateTransition)
	}

	present := 0
	absent := 0
	for _, r := range records {
		if r.Status == RecordPresent || r.Status == RecordLate || r.Status == RecordExcusedMedical {
			present++
		} else {
			absent++
		}
	}

	s.TotalStudents = len(records)
	s.PresentCount = present
	s.AbsentCount = absent
	s.Status = SessionFinalized
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// CalculatePercentage computes effective attendance rate considering excused medical leaves.
func CalculatePercentage(total, present, excused int) (float64, bool) {
	if total == 0 {
		return 100.0, false
	}
	effectivePresent := float64(present + excused)
	pct := (effectivePresent / float64(total)) * 100.0
	isShortage := pct < 75.0
	return pct, isShortage
}

func (m *MedicalLeaveApplication) Validate() error {
	if strings.TrimSpace(m.TenantID) == "" {
		return ErrTenantRequired
	}
	if strings.TrimSpace(m.StudentID) == "" {
		return NewDomainError("INVALID_STUDENT", "student ID is required", ErrInvalidInput)
	}
	if strings.TrimSpace(m.FromDate) == "" || strings.TrimSpace(m.ToDate) == "" {
		return NewDomainError("INVALID_DATES", "from and to dates are required", ErrInvalidInput)
	}
	if strings.TrimSpace(m.Reason) == "" {
		return NewDomainError("INVALID_REASON", "leave medical reason is required", ErrInvalidInput)
	}
	return nil
}
