/**
 * BLOCK_ATTENDANCE_SERVICE_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   Core business orchestration for multi-mode roll call, shortage tracking (<75%), and medical condonations.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package attendance

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"campus/backend/audit"
)

type CreateSubjectCommand struct {
	TenantID  string `json:"tenant_id"`
	ProgramID string `json:"program_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Credits   int    `json:"credits"`
	Semester  int    `json:"semester"`
}

type CreateSlotCommand struct {
	TenantID  string    `json:"tenant_id"`
	SubjectID string    `json:"subject_id"`
	CohortID  string    `json:"cohort_id"`
	FacultyID string    `json:"faculty_id"`
	DayOfWeek DayOfWeek `json:"day_of_week"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	RoomCode  string    `json:"room_code"`
}

type ScheduleSessionCommand struct {
	TenantID    string         `json:"tenant_id"`
	SlotID      string         `json:"slot_id,omitempty"`
	SubjectID   string         `json:"subject_id"`
	CohortID    string         `json:"cohort_id"`
	FacultyID   string         `json:"faculty_id"`
	SessionDate string         `json:"session_date"` // YYYY-MM-DD
	StartTime   string         `json:"start_time"`
	EndTime     string         `json:"end_time"`
	Mode        AttendanceMode `json:"mode"`
}

type RecordInput struct {
	StudentID string                 `json:"student_id"`
	Status    AttendanceRecordStatus `json:"status"`
	Remarks   string                 `json:"remarks,omitempty"`
	DeviceID  string                 `json:"device_id,omitempty"`
	Latitude  *float64               `json:"latitude,omitempty"`
	Longitude *float64               `json:"longitude,omitempty"`
}

type MarkAttendanceCommand struct {
	TenantID  string        `json:"tenant_id"`
	SessionID string        `json:"session_id"`
	FacultyID string        `json:"faculty_id"`
	Records   []RecordInput `json:"records"`
}

type BiometricCheckinInput struct {
	StudentID string    `json:"student_id"`
	DeviceID  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
}

type IngestBiometricCommand struct {
	TenantID  string                  `json:"tenant_id"`
	SessionID string                  `json:"session_id"`
	Checkins  []BiometricCheckinInput `json:"checkins"`
}

type ApplyMedicalLeaveCommand struct {
	TenantID       string `json:"tenant_id"`
	StudentID      string `json:"student_id"`
	FromDate       string `json:"from_date"`
	ToDate         string `json:"to_date"`
	Reason         string `json:"reason"`
	CertificateKey string `json:"certificate_key"`
}

type Service struct {
	catalogRepo CatalogRepository
	sessionRepo SessionRepository
	recordRepo  RecordRepository
	leaveRepo   MedicalLeaveRepository
	auditSub    audit.Subscriber
}

func NewService(
	catalogRepo CatalogRepository,
	sessionRepo SessionRepository,
	recordRepo RecordRepository,
	leaveRepo MedicalLeaveRepository,
	auditSub audit.Subscriber,
) *Service {
	return &Service{
		catalogRepo: catalogRepo,
		sessionRepo: sessionRepo,
		recordRepo:  recordRepo,
		leaveRepo:   leaveRepo,
		auditSub:    auditSub,
	}
}

func (s *Service) CreateSubject(ctx context.Context, cmd CreateSubjectCommand) (*CourseSubject, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	cleanCode := strings.ToUpper(strings.TrimSpace(cmd.Code))
	if cleanCode == "" || strings.TrimSpace(cmd.Name) == "" {
		return nil, NewDomainError("INVALID_SUBJECT", "subject code and name are required", ErrInvalidInput)
	}

	subject := &CourseSubject{
		ID:        generateID("sub_"),
		TenantID:  cmd.TenantID,
		ProgramID: cmd.ProgramID,
		Code:      cleanCode,
		Name:      strings.TrimSpace(cmd.Name),
		Credits:   cmd.Credits,
		Semester:  cmd.Semester,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.catalogRepo.CreateSubject(ctx, subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *Service) CreateTimetableSlot(ctx context.Context, cmd CreateSlotCommand) (*TimetableSlot, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(cmd.SubjectID) == "" || strings.TrimSpace(cmd.CohortID) == "" || strings.TrimSpace(cmd.FacultyID) == "" {
		return nil, NewDomainError("INVALID_SLOT", "subject, cohort, and faculty are required", ErrInvalidInput)
	}

	slot := &TimetableSlot{
		ID:        generateID("slot_"),
		TenantID:  cmd.TenantID,
		SubjectID: cmd.SubjectID,
		CohortID:  cmd.CohortID,
		FacultyID: cmd.FacultyID,
		DayOfWeek: cmd.DayOfWeek,
		StartTime: cmd.StartTime,
		EndTime:   cmd.EndTime,
		RoomCode:  cmd.RoomCode,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.catalogRepo.CreateSlot(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

func (s *Service) ScheduleSession(ctx context.Context, cmd ScheduleSessionCommand) (*AttendanceSession, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(cmd.SubjectID) == "" || strings.TrimSpace(cmd.CohortID) == "" || strings.TrimSpace(cmd.SessionDate) == "" {
		return nil, NewDomainError("INVALID_SESSION", "subject, cohort, and session date are required", ErrInvalidInput)
	}
	if cmd.Mode == "" {
		cmd.Mode = ModeManualFaculty
	}

	session := &AttendanceSession{
		ID:          generateID("sess_"),
		TenantID:    cmd.TenantID,
		SlotID:      cmd.SlotID,
		SubjectID:   cmd.SubjectID,
		CohortID:    cmd.CohortID,
		FacultyID:   cmd.FacultyID,
		SessionDate: cmd.SessionDate,
		StartTime:   cmd.StartTime,
		EndTime:     cmd.EndTime,
		Mode:        cmd.Mode,
		Status:      SessionScheduled,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	s.emitAudit(cmd.TenantID, cmd.FacultyID, "USER", "attendance:session:scheduled", "attendance_session", session.ID, audit.StatusSuccess, map[string]interface{}{
		"subject_id":   cmd.SubjectID,
		"cohort_id":    cmd.CohortID,
		"session_date": cmd.SessionDate,
	})

	return session, nil
}

func (s *Service) OpenSession(ctx context.Context, tenantID, sessionID, facultyID string) (*AttendanceSession, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}

	if err := session.Open(); err != nil {
		return nil, err
	}

	if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, facultyID, "USER", "attendance:session:opened", "attendance_session", session.ID, audit.StatusSuccess, nil)
	return session, nil
}

func (s *Service) MarkAttendance(ctx context.Context, cmd MarkAttendanceCommand) (*AttendanceSession, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	session, err := s.sessionRepo.GetSessionByID(ctx, cmd.TenantID, cmd.SessionID)
	if err != nil {
		return nil, err
	}
	if session.Status != SessionOpen && session.Status != SessionScheduled {
		return nil, ErrSessionLocked
	}

	now := time.Now().UTC()
	records := make([]AttendanceRecord, 0, len(cmd.Records))
	for _, rec := range cmd.Records {
		if rec.Status == "" {
			rec.Status = RecordPresent
		}
		records = append(records, AttendanceRecord{
			ID:        generateID("att_rec_"),
			TenantID:  cmd.TenantID,
			SessionID: session.ID,
			StudentID: rec.StudentID,
			Status:    rec.Status,
			Remarks:   rec.Remarks,
			DeviceID:  rec.DeviceID,
			Latitude:  rec.Latitude,
			Longitude: rec.Longitude,
			MarkedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	if err := s.recordRepo.BatchUpsertRecords(ctx, records); err != nil {
		return nil, err
	}

	_ = session.Finalize(records)
	if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}

	session.Records = records

	s.emitAudit(cmd.TenantID, cmd.FacultyID, "USER", "attendance:session:marked", "attendance_session", session.ID, audit.StatusSuccess, map[string]interface{}{
		"total_students": session.TotalStudents,
		"present_count":  session.PresentCount,
		"absent_count":   session.AbsentCount,
	})

	return session, nil
}

func (s *Service) IngestBiometricBatch(ctx context.Context, cmd IngestBiometricCommand) (*AttendanceSession, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	session, err := s.sessionRepo.GetSessionByID(ctx, cmd.TenantID, cmd.SessionID)
	if err != nil {
		return nil, err
	}
	if session.Status == SessionFinalized {
		return nil, ErrSessionLocked
	}

	now := time.Now().UTC()
	records := make([]AttendanceRecord, 0, len(cmd.Checkins))
	for _, chk := range cmd.Checkins {
		records = append(records, AttendanceRecord{
			ID:        generateID("bio_rec_"),
			TenantID:  cmd.TenantID,
			SessionID: session.ID,
			StudentID: chk.StudentID,
			Status:    RecordPresent,
			DeviceID:  chk.DeviceID,
			Remarks:   fmt.Sprintf("Biometric terminal sync from %s", chk.DeviceID),
			MarkedAt:  chk.Timestamp,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	if err := s.recordRepo.BatchUpsertRecords(ctx, records); err != nil {
		return nil, err
	}

	_ = session.Finalize(records)
	_ = s.sessionRepo.UpdateSession(ctx, session)
	session.Records = records
	return session, nil
}

func (s *Service) LockSession(ctx context.Context, tenantID, sessionID, facultyID string) (*AttendanceSession, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}

	if err := session.Lock(); err != nil {
		return nil, err
	}

	if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, facultyID, "USER", "attendance:session:locked", "attendance_session", session.ID, audit.StatusSuccess, nil)
	return session, nil
}

func (s *Service) FinalizeSession(ctx context.Context, tenantID, sessionID, facultyID string) (*AttendanceSession, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}

	records, err := s.recordRepo.GetRecordsBySession(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}

	if err := session.Finalize(records); err != nil {
		return nil, err
	}

	if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, facultyID, "USER", "attendance:session:finalized", "attendance_session", session.ID, audit.StatusSuccess, map[string]interface{}{
		"present_count": session.PresentCount,
		"absent_count":  session.AbsentCount,
	})

	return session, nil
}

func (s *Service) GetStudentAttendanceSummary(ctx context.Context, tenantID, studentID, subjectID string) (*StudentAttendanceSummary, error) {
	records, err := s.recordRepo.GetSubjectRecordsForStudent(ctx, tenantID, studentID, subjectID)
	if err != nil {
		return nil, err
	}

	total := len(records)
	present := 0
	excused := 0
	absent := 0

	for _, r := range records {
		switch r.Status {
		case RecordPresent, RecordLate:
			present++
		case RecordExcusedMedical:
			excused++
		default:
			absent++
		}
	}

	pct, isShortage := CalculatePercentage(total, present, excused)

	return &StudentAttendanceSummary{
		StudentID:            studentID,
		SubjectID:            subjectID,
		TotalSessions:        total,
		PresentSessions:      present,
		ExcusedSessions:      excused,
		AbsentSessions:       absent,
		AttendancePercentage: pct,
		IsShortage:           isShortage,
	}, nil
}

func (s *Service) ApplyMedicalLeave(ctx context.Context, cmd ApplyMedicalLeaveCommand) (*MedicalLeaveApplication, error) {
	leave := &MedicalLeaveApplication{
		ID:             generateID("med_leave_"),
		TenantID:       cmd.TenantID,
		StudentID:      cmd.StudentID,
		FromDate:       cmd.FromDate,
		ToDate:         cmd.ToDate,
		Reason:         cmd.Reason,
		CertificateKey: cmd.CertificateKey,
		Status:         MedicalLeavePending,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	if err := leave.Validate(); err != nil {
		return nil, err
	}

	if err := s.leaveRepo.CreateLeave(ctx, leave); err != nil {
		return nil, err
	}

	s.emitAudit(cmd.TenantID, cmd.StudentID, "USER", "attendance:medical_leave:applied", "medical_leave_application", leave.ID, audit.StatusSuccess, nil)
	return leave, nil
}

func (s *Service) ApproveMedicalLeave(ctx context.Context, tenantID, leaveID, approverID string, approved bool, rejectionReason string) (*MedicalLeaveApplication, error) {
	leave, err := s.leaveRepo.GetLeaveByID(ctx, tenantID, leaveID)
	if err != nil {
		return nil, err
	}

	status := MedicalLeaveApproved
	if !approved {
		status = MedicalLeaveRejected
	}

	now := time.Now().UTC()
	if err := s.leaveRepo.UpdateLeaveStatus(ctx, tenantID, leaveID, status, approverID, rejectionReason, &now); err != nil {
		return nil, err
	}

	leave.Status = status
	leave.ApprovedByID = approverID
	leave.ApprovedAt = &now
	leave.RejectionReason = rejectionReason
	leave.UpdatedAt = now

	// If approved, automatically convert affected attendance records to EXCUSED_MEDICAL
	if approved {
		_, _ = s.recordRepo.MarkExcusedMedical(ctx, tenantID, leave.StudentID, leave.FromDate, leave.ToDate)
	}

	s.emitAudit(tenantID, approverID, "USER", "attendance:medical_leave:decision", "medical_leave_application", leave.ID, audit.StatusSuccess, map[string]interface{}{
		"status":   status,
		"approved": approved,
	})

	return leave, nil
}

func (s *Service) ListSessions(ctx context.Context, filter SessionFilter) ([]*AttendanceSession, int, error) {
	return s.sessionRepo.ListSessions(ctx, filter)
}

func (s *Service) GetSession(ctx context.Context, tenantID, sessionID string) (*AttendanceSession, error) {
	sess, err := s.sessionRepo.GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}
	recs, _ := s.recordRepo.GetRecordsBySession(ctx, tenantID, sessionID)
	sess.Records = recs
	return sess, nil
}

func (s *Service) emitAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
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

func generateID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + hex.EncodeToString(b)
}
