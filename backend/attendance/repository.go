/**
 * BLOCK_ATTENDANCE_REPOSITORY_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   Data access interface definitions for timetable slots, sessions, student records, and medical leaves.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package attendance

import (
	"context"
	"time"
)

type SessionFilter struct {
	TenantID     string
	SubjectID    string
	CohortID     string
	FacultyID    string
	SessionDate  string
	Status       *AttendanceSessionStatus
	FromDate     string
	ToDate       string
	Limit        int
	Offset       int
}

type CatalogRepository interface {
	CreateSubject(ctx context.Context, subject *CourseSubject) error
	GetSubjectByID(ctx context.Context, tenantID, id string) (*CourseSubject, error)
	GetSubjectByCode(ctx context.Context, tenantID, programID, code string) (*CourseSubject, error)
	ListSubjects(ctx context.Context, tenantID, programID string, semester int) ([]*CourseSubject, error)

	CreateSlot(ctx context.Context, slot *TimetableSlot) error
	GetSlotByID(ctx context.Context, tenantID, id string) (*TimetableSlot, error)
	ListSlotsByCohort(ctx context.Context, tenantID, cohortID string) ([]*TimetableSlot, error)
	ListSlotsByFaculty(ctx context.Context, tenantID, facultyID string) ([]*TimetableSlot, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *AttendanceSession) error
	GetSessionByID(ctx context.Context, tenantID, id string) (*AttendanceSession, error)
	UpdateSession(ctx context.Context, session *AttendanceSession) error
	ListSessions(ctx context.Context, filter SessionFilter) ([]*AttendanceSession, int, error)
}

type RecordRepository interface {
	BatchUpsertRecords(ctx context.Context, records []AttendanceRecord) error
	GetRecordsBySession(ctx context.Context, tenantID, sessionID string) ([]AttendanceRecord, error)
	GetRecordsByStudent(ctx context.Context, tenantID, studentID string) ([]AttendanceRecord, error)
	GetSubjectRecordsForStudent(ctx context.Context, tenantID, studentID, subjectID string) ([]AttendanceRecord, error)
	MarkExcusedMedical(ctx context.Context, tenantID, studentID, fromDate, toDate string) (int, error)
}

type MedicalLeaveRepository interface {
	CreateLeave(ctx context.Context, leave *MedicalLeaveApplication) error
	GetLeaveByID(ctx context.Context, tenantID, id string) (*MedicalLeaveApplication, error)
	UpdateLeaveStatus(ctx context.Context, tenantID, id string, status MedicalLeaveStatus, approvedByID, rejectionReason string, approvedAt *time.Time) error
	ListLeavesByStudent(ctx context.Context, tenantID, studentID string) ([]*MedicalLeaveApplication, error)
	ListPendingLeaves(ctx context.Context, tenantID string) ([]*MedicalLeaveApplication, error)
}
