/**
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   Unit tests for domain, business service, shortage detection, and medical condonation cascade.
 */

package attendance_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"campus/backend/attendance"
	"campus/backend/audit"
)

type mockAuditSubscriber struct {
	mu     sync.Mutex
	events []audit.RecordAuditRequest
}

func (m *mockAuditSubscriber) Enqueue(req audit.RecordAuditRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, req)
	return nil
}

func (m *mockAuditSubscriber) Start(ctx context.Context) {}
func (m *mockAuditSubscriber) Stop()                      {}

func setupTestService() (*attendance.Service, *attendance.MockCatalogRepository, *attendance.MockSessionRepository, *attendance.MockRecordRepository, *attendance.MockMedicalLeaveRepository, *mockAuditSubscriber) {
	catRepo := attendance.NewMockCatalogRepository()
	sessRepo := attendance.NewMockSessionRepository()
	recRepo := attendance.NewMockRecordRepository()
	leaveRepo := attendance.NewMockMedicalLeaveRepository()
	auditMock := &mockAuditSubscriber{}

	svc := attendance.NewService(catRepo, sessRepo, recRepo, leaveRepo, auditMock)
	return svc, catRepo, sessRepo, recRepo, leaveRepo, auditMock
}

func TestCreateSubject(t *testing.T) {
	svc, _, _, _, _, _ := setupTestService()
	ctx := context.Background()

	// 1. Success
	sub, err := svc.CreateSubject(ctx, attendance.CreateSubjectCommand{
		TenantID:  "ten_default",
		ProgramID: "prog_cs",
		Code:      "cs101",
		Name:      "Introduction to Computer Science",
		Credits:   4,
		Semester:  1,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if sub.Code != "CS101" {
		t.Errorf("expected code to be uppercased CS101, got %s", sub.Code)
	}

	// 2. Missing Tenant
	_, err = svc.CreateSubject(ctx, attendance.CreateSubjectCommand{
		TenantID: "",
		Code:     "CS102",
		Name:     "Data Structures",
	})
	if err != attendance.ErrTenantRequired {
		t.Errorf("expected ErrTenantRequired, got %v", err)
	}

	// 3. Missing Code/Name
	_, err = svc.CreateSubject(ctx, attendance.CreateSubjectCommand{
		TenantID: "ten_default",
		Code:     "",
		Name:     "",
	})
	if err == nil {
		t.Errorf("expected error for empty code/name, got nil")
	}
}

func TestCreateTimetableSlot(t *testing.T) {
	svc, _, _, _, _, _ := setupTestService()
	ctx := context.Background()

	slot, err := svc.CreateTimetableSlot(ctx, attendance.CreateSlotCommand{
		TenantID:  "ten_default",
		SubjectID: "sub_cs101",
		CohortID:  "coh_2026",
		FacultyID: "fac_prof_smith",
		DayOfWeek: attendance.DayMonday,
		StartTime: "09:00",
		EndTime:   "10:00",
		RoomCode:  "LH-101",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if slot.RoomCode != "LH-101" {
		t.Errorf("expected room code LH-101, got %s", slot.RoomCode)
	}

	// Validation check
	_, err = svc.CreateTimetableSlot(ctx, attendance.CreateSlotCommand{
		TenantID:  "ten_default",
		SubjectID: "",
	})
	if err == nil {
		t.Errorf("expected error for missing subject ID")
	}
}

func TestScheduleAndSessionLifecycle(t *testing.T) {
	svc, _, _, _, _, _ := setupTestService()
	ctx := context.Background()

	// 1. Schedule session
	sess, err := svc.ScheduleSession(ctx, attendance.ScheduleSessionCommand{
		TenantID:    "ten_default",
		SubjectID:   "sub_cs101",
		CohortID:    "coh_2026",
		FacultyID:   "fac_prof_smith",
		SessionDate: "2026-09-15",
		StartTime:   "09:00",
		EndTime:     "10:00",
		Mode:        attendance.ModeManualFaculty,
	})
	if err != nil {
		t.Fatalf("ScheduleSession failed: %v", err)
	}
	if sess.Status != attendance.SessionScheduled {
		t.Errorf("expected status SCHEDULED, got %s", sess.Status)
	}

	// 2. Open session
	opened, err := svc.OpenSession(ctx, "ten_default", sess.ID, "fac_prof_smith")
	if err != nil {
		t.Fatalf("OpenSession failed: %v", err)
	}
	if opened.Status != attendance.SessionOpen {
		t.Errorf("expected status OPEN, got %s", opened.Status)
	}

	// 3. Opening already open session should fail
	_, err = svc.OpenSession(ctx, "ten_default", sess.ID, "fac_prof_smith")
	if err == nil {
		t.Errorf("expected error opening already open session")
	}

	// 4. Lock session
	locked, err := svc.LockSession(ctx, "ten_default", sess.ID, "fac_prof_smith")
	if err != nil {
		t.Fatalf("LockSession failed: %v", err)
	}
	if locked.Status != attendance.SessionLocked {
		t.Errorf("expected status LOCKED, got %s", locked.Status)
	}

	// 5. Finalize session
	finalized, err := svc.FinalizeSession(ctx, "ten_default", sess.ID, "fac_prof_smith")
	if err != nil {
		t.Fatalf("FinalizeSession failed: %v", err)
	}
	if finalized.Status != attendance.SessionFinalized {
		t.Errorf("expected status FINALIZED, got %s", finalized.Status)
	}
}

func TestMarkAttendance(t *testing.T) {
	svc, _, _, _, _, auditMock := setupTestService()
	ctx := context.Background()

	sess, err := svc.ScheduleSession(ctx, attendance.ScheduleSessionCommand{
		TenantID:    "ten_default",
		SubjectID:   "sub_cs101",
		CohortID:    "coh_2026",
		FacultyID:   "fac_prof_smith",
		SessionDate: "2026-09-15",
		StartTime:   "09:00",
		EndTime:     "10:00",
	})
	if err != nil {
		t.Fatalf("ScheduleSession failed: %v", err)
	}

	// Mark attendance for 3 students
	marked, err := svc.MarkAttendance(ctx, attendance.MarkAttendanceCommand{
		TenantID:  "ten_default",
		SessionID: sess.ID,
		FacultyID: "fac_prof_smith",
		Records: []attendance.RecordInput{
			{StudentID: "stu_1", Status: attendance.RecordPresent},
			{StudentID: "stu_2", Status: attendance.RecordAbsent},
			{StudentID: "stu_3", Status: attendance.RecordLate},
		},
	})
	if err != nil {
		t.Fatalf("MarkAttendance failed: %v", err)
	}

	if marked.TotalStudents != 3 {
		t.Errorf("expected TotalStudents=3, got %d", marked.TotalStudents)
	}
	if marked.PresentCount != 2 { // Present + Late count as present
		t.Errorf("expected PresentCount=2, got %d", marked.PresentCount)
	}
	if marked.AbsentCount != 1 {
		t.Errorf("expected AbsentCount=1, got %d", marked.AbsentCount)
	}
	if marked.Status != attendance.SessionFinalized {
		t.Errorf("expected status FINALIZED, got %s", marked.Status)
	}

	if len(auditMock.events) == 0 {
		t.Errorf("expected audit event to be enqueued")
	}
}

func TestBiometricBatchIngest(t *testing.T) {
	svc, _, _, _, _, _ := setupTestService()
	ctx := context.Background()

	sess, _ := svc.ScheduleSession(ctx, attendance.ScheduleSessionCommand{
		TenantID:    "ten_default",
		SubjectID:   "sub_cs101",
		CohortID:    "coh_2026",
		FacultyID:   "fac_prof_smith",
		SessionDate: "2026-09-16",
		Mode:        attendance.ModeBiometricTerminal,
	})

	now := time.Now().UTC()
	ingested, err := svc.IngestBiometricBatch(ctx, attendance.IngestBiometricCommand{
		TenantID:  "ten_default",
		SessionID: sess.ID,
		Checkins: []attendance.BiometricCheckinInput{
			{StudentID: "stu_1", DeviceID: "GATE_BIO_01", Timestamp: now},
			{StudentID: "stu_2", DeviceID: "GATE_BIO_01", Timestamp: now},
		},
	})
	if err != nil {
		t.Fatalf("IngestBiometricBatch failed: %v", err)
	}

	if ingested.PresentCount != 2 {
		t.Errorf("expected PresentCount=2, got %d", ingested.PresentCount)
	}
}

func TestShortageCalculationAndSummary(t *testing.T) {
	svc, _, sessRepo, recRepo, _, _ := setupTestService()
	ctx := context.Background()

	studentID := "stu_shortage_test"
	subjectID := "sub_cs101"
	tenantID := "ten_default"

	// Create dummy session for the subject
	sess := &attendance.AttendanceSession{
		ID:          "sess_shortage",
		TenantID:    tenantID,
		SubjectID:   subjectID,
		SessionDate: "2026-09-15",
		Status:      attendance.SessionFinalized,
	}
	_ = sessRepo.CreateSession(ctx, sess)

	// Scenario 1: 10 sessions, 8 present (80% -> no shortage)
	var records []attendance.AttendanceRecord
	for i := 1; i <= 8; i++ {
		records = append(records, attendance.AttendanceRecord{
			ID:        fmt.Sprintf("rec_%d", i),
			TenantID:  tenantID,
			SessionID: fmt.Sprintf("sess_%d", i),
			StudentID: studentID,
			Status:    attendance.RecordPresent,
		})
	}
	for i := 9; i <= 10; i++ {
		records = append(records, attendance.AttendanceRecord{
			ID:        fmt.Sprintf("rec_%d", i),
			TenantID:  tenantID,
			SessionID: fmt.Sprintf("sess_%d", i),
			StudentID: studentID,
			Status:    attendance.RecordAbsent,
		})
	}
	_ = recRepo.BatchUpsertRecords(ctx, records)

	summary, err := svc.GetStudentAttendanceSummary(ctx, tenantID, studentID, subjectID)
	if err != nil {
		t.Fatalf("GetStudentAttendanceSummary failed: %v", err)
	}
	if summary.AttendancePercentage != 80.0 {
		t.Errorf("expected 80%%, got %.2f%%", summary.AttendancePercentage)
	}
	if summary.IsShortage {
		t.Errorf("expected IsShortage=false for 80%%")
	}

	// Scenario 2: Add 3 more absences -> 8 present out of 13 total (61.53% -> shortage)
	var extraAbsences []attendance.AttendanceRecord
	for i := 11; i <= 13; i++ {
		extraAbsences = append(extraAbsences, attendance.AttendanceRecord{
			ID:        fmt.Sprintf("rec_%d", i),
			TenantID:  tenantID,
			SessionID: fmt.Sprintf("sess_%d", i),
			StudentID: studentID,
			Status:    attendance.RecordAbsent,
		})
	}
	_ = recRepo.BatchUpsertRecords(ctx, extraAbsences)

	summary2, err := svc.GetStudentAttendanceSummary(ctx, tenantID, studentID, subjectID)
	if err != nil {
		t.Fatalf("GetStudentAttendanceSummary failed: %v", err)
	}
	if !summary2.IsShortage {
		t.Errorf("expected IsShortage=true for %.2f%%", summary2.AttendancePercentage)
	}
}

func TestMedicalLeaveCondonationCascade(t *testing.T) {
	svc, _, sessRepo, recRepo, _, _ := setupTestService()
	ctx := context.Background()

	tenantID := "ten_default"
	studentID := "stu_sick_1"
	subjectID := "sub_cs101"

	// Create 2 sessions on dates 2026-09-10 and 2026-09-11
	s1 := &attendance.AttendanceSession{
		ID:          "sess_day1",
		TenantID:    tenantID,
		SubjectID:   subjectID,
		SessionDate: "2026-09-10",
		Status:      attendance.SessionFinalized,
	}
	s2 := &attendance.AttendanceSession{
		ID:          "sess_day2",
		TenantID:    tenantID,
		SubjectID:   subjectID,
		SessionDate: "2026-09-11",
		Status:      attendance.SessionFinalized,
	}
	_ = sessRepo.CreateSession(ctx, s1)
	_ = sessRepo.CreateSession(ctx, s2)

	// Mark student absent for both sessions
	_ = recRepo.BatchUpsertRecords(ctx, []attendance.AttendanceRecord{
		{
			ID:        "rec_1",
			TenantID:  tenantID,
			SessionID: s1.ID,
			StudentID: studentID,
			Status:    attendance.RecordAbsent,
		},
		{
			ID:        "rec_2",
			TenantID:  tenantID,
			SessionID: s2.ID,
			StudentID: studentID,
			Status:    attendance.RecordAbsent,
		},
	})

	// Check initial summary: 0% attendance, isShortage = true
	sumInitial, err := svc.GetStudentAttendanceSummary(ctx, tenantID, studentID, subjectID)
	if err != nil {
		t.Fatalf("GetStudentAttendanceSummary failed: %v", err)
	}
	if sumInitial.AttendancePercentage != 0.0 || !sumInitial.IsShortage {
		t.Errorf("expected 0%% and shortage=true, got %.2f and %v", sumInitial.AttendancePercentage, sumInitial.IsShortage)
	}

	// Apply medical leave for 2026-09-10 to 2026-09-11
	leave, err := svc.ApplyMedicalLeave(ctx, attendance.ApplyMedicalLeaveCommand{
		TenantID:       tenantID,
		StudentID:      studentID,
		FromDate:       "2026-09-10",
		ToDate:         "2026-09-11",
		Reason:         "Dengue fever hospitalization",
		CertificateKey: "s3://medical/cert_123.pdf",
	})
	if err != nil {
		t.Fatalf("ApplyMedicalLeave failed: %v", err)
	}
	if leave.Status != attendance.MedicalLeavePending {
		t.Errorf("expected PENDING status, got %s", leave.Status)
	}

	// Approve medical leave
	approvedLeave, err := svc.ApproveMedicalLeave(ctx, tenantID, leave.ID, "fac_dean_john", true, "")
	if err != nil {
		t.Fatalf("ApproveMedicalLeave failed: %v", err)
	}
	if approvedLeave.Status != attendance.MedicalLeaveApproved {
		t.Errorf("expected APPROVED, got %s", approvedLeave.Status)
	}

	// Verify attendance summary after approval: records are now excused medical -> 100% attendance, no shortage
	sumAfter, err := svc.GetStudentAttendanceSummary(ctx, tenantID, studentID, subjectID)
	if err != nil {
		t.Fatalf("GetStudentAttendanceSummary after approval failed: %v", err)
	}
	if sumAfter.ExcusedSessions != 2 {
		t.Errorf("expected 2 excused sessions, got %d", sumAfter.ExcusedSessions)
	}
	if sumAfter.AttendancePercentage != 100.0 {
		t.Errorf("expected 100.0%% after medical condonation, got %.2f%%", sumAfter.AttendancePercentage)
	}
	if sumAfter.IsShortage {
		t.Errorf("expected IsShortage=false after medical condonation")
	}
}
