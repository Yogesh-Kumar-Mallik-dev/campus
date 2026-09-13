/**
 * BLOCK_NOTICES_TEST_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   Unit tests verifying notice authoring, priority pinning, target audience filtering, read receipts, and lifecycle transitions.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package notices_test

import (
	"context"
	"sync"
	"testing"

	"campus/backend/audit"
	"campus/backend/notices"
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

func setupNoticesTestService() (*notices.Service, *notices.MockNoticeRepository, *notices.MockAttachmentRepository, *notices.MockAcknowledgementRepository, *mockAuditSubscriber) {
	noticeRepo := notices.NewMockNoticeRepository()
	attRepo := notices.NewMockAttachmentRepository()
	ackRepo := notices.NewMockAcknowledgementRepository()
	auditSub := &mockAuditSubscriber{}

	svc := notices.NewService(noticeRepo, attRepo, ackRepo, auditSub)
	return svc, noticeRepo, attRepo, ackRepo, auditSub
}

func TestCreateNotice_DraftAndPublished(t *testing.T) {
	svc, _, _, _, auditSub := setupNoticesTestService()
	ctx := context.Background()

	// 1. Success Draft Notice
	draft, err := svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID:       "ten_default",
		AuthorID:       "usr_dean_01",
		Title:          "Midterm Examination Schedule Released",
		Content:        "The Midterm examinations will commence from 15th October 2026 across all departments.",
		Category:       notices.CategoryExamination,
		Priority:       notices.PriorityHigh,
		TargetAudience: notices.AudienceStudents,
		IsPinned:       true,
		PublishNow:     false,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if draft.Status != notices.StatusDraft {
		t.Errorf("expected status DRAFT, got %s", draft.Status)
	}
	if !draft.IsPinned {
		t.Errorf("expected IsPinned=true")
	}

	// 2. Missing Tenant
	_, err = svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID: "",
		Title:    "No Tenant Notice",
		Content:  "Content",
	})
	if err != notices.ErrTenantRequired {
		t.Errorf("expected ErrTenantRequired, got %v", err)
	}

	// 3. Empty Title/Content
	_, err = svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID: "ten_default",
		Title:    "",
		Content:  "",
	})
	if err == nil {
		t.Errorf("expected error for empty title/content, got nil")
	}

	// 4. Publish Now Notice
	pubNotice, err := svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID:       "ten_default",
		AuthorID:       "usr_admin_01",
		Title:          "Campus Power Maintenance Announcement",
		Content:        "Scheduled maintenance on Sunday 08:00 to 12:00.",
		Category:       notices.CategoryAdministrative,
		Priority:       notices.PriorityUrgent,
		TargetAudience: notices.AudienceAll,
		PublishNow:     true,
	})
	if err != nil {
		t.Fatalf("expected nil error for publishNow, got %v", err)
	}
	if pubNotice.Status != notices.StatusPublished {
		t.Errorf("expected status PUBLISHED, got %s", pubNotice.Status)
	}

	if len(auditSub.events) == 0 {
		t.Errorf("expected audit events to be enqueued")
	}
}

func TestNoticeLifecycle_PublishAndArchive(t *testing.T) {
	svc, _, _, _, _ := setupNoticesTestService()
	ctx := context.Background()

	// 1. Create Draft
	n, _ := svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID: "ten_default",
		AuthorID: "usr_dean_01",
		Title:    "Annual Tech Fest 2026",
		Content:  "Registrations open for hackathons and design challenges.",
		Category: notices.CategoryEvents,
	})

	// 2. Publish
	published, err := svc.PublishNotice(ctx, "ten_default", n.ID, "usr_dean_01")
	if err != nil {
		t.Fatalf("PublishNotice failed: %v", err)
	}
	if published.Status != notices.StatusPublished {
		t.Errorf("expected status PUBLISHED, got %s", published.Status)
	}

	// 3. Publishing already published notice fails
	_, err = svc.PublishNotice(ctx, "ten_default", n.ID, "usr_dean_01")
	if err != notices.ErrNoticeAlreadyPublished {
		t.Errorf("expected ErrNoticeAlreadyPublished, got %v", err)
	}

	// 4. Archive
	archived, err := svc.ArchiveNotice(ctx, "ten_default", n.ID, "usr_dean_01")
	if err != nil {
		t.Fatalf("ArchiveNotice failed: %v", err)
	}
	if archived.Status != notices.StatusArchived {
		t.Errorf("expected status ARCHIVED, got %s", archived.Status)
	}
}

func TestAcknowledgeAndReadReceipts(t *testing.T) {
	svc, _, _, _, _ := setupNoticesTestService()
	ctx := context.Background()

	n, _ := svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID:       "ten_default",
		AuthorID:       "usr_warden_01",
		Title:          "Mandatory Hostel Curfew Revision",
		Content:        "Curfew hours shifted to 22:00 starting October 1st.",
		Category:       notices.CategoryHostel,
		Priority:       notices.PriorityUrgent,
		TargetAudience: notices.AudienceHostelResidents,
		PublishNow:     true,
	})

	devID := "mobile_device_pixel_8"
	ack, err := svc.AcknowledgeNotice(ctx, notices.AcknowledgeNoticeCommand{
		TenantID: "ten_default",
		NoticeID: n.ID,
		UserID:   "usr_student_aarav",
		DeviceID: &devID,
	})
	if err != nil {
		t.Fatalf("AcknowledgeNotice failed: %v", err)
	}
	if ack.UserID != "usr_student_aarav" {
		t.Errorf("expected user_id usr_student_aarav, got %s", ack.UserID)
	}
	if ack.AcknowledgedAt == nil {
		t.Errorf("expected acknowledged_at to be populated")
	}
}

func TestAttachmentsAndViews(t *testing.T) {
	svc, _, _, _, _ := setupNoticesTestService()
	ctx := context.Background()

	n, _ := svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID:   "ten_default",
		AuthorID:   "usr_admin",
		Title:      "Academic Calendar 2026-2027",
		Content:    "Attached is the official institutional calendar approved by the Academic Council.",
		PublishNow: true,
	})

	// Attach PDF Circular
	att, err := svc.AttachDocument(ctx, notices.AttachDocumentCommand{
		TenantID: "ten_default",
		NoticeID: n.ID,
		FileName: "Academic_Calendar_2026.pdf",
		FileKey:  "s3://circulars/cal_2026.pdf",
		FileSize: 1048576,
		MimeType: "application/pdf",
	})
	if err != nil {
		t.Fatalf("AttachDocument failed: %v", err)
	}
	if att.FileName != "Academic_Calendar_2026.pdf" {
		t.Errorf("expected file name Academic_Calendar_2026.pdf, got %s", att.FileName)
	}

	// Get Notice with view count increment
	fetched, err := svc.GetNotice(ctx, "ten_default", n.ID, true)
	if err != nil {
		t.Fatalf("GetNotice failed: %v", err)
	}
	if fetched.ViewCount != 1 {
		t.Errorf("expected view count 1, got %d", fetched.ViewCount)
	}
	if len(fetched.Attachments) != 1 {
		t.Errorf("expected 1 attachment, got %d", len(fetched.Attachments))
	}
}

func TestTargetAudienceAndFilters(t *testing.T) {
	svc, _, _, _, _ := setupNoticesTestService()
	ctx := context.Background()

	// Notice 1: Exam - Students
	_, _ = svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID:       "ten_default",
		AuthorID:       "usr_1",
		Title:          "Endsem Hall Tickets",
		Content:        "Hall tickets available for download",
		Category:       notices.CategoryExamination,
		Priority:       notices.PriorityHigh,
		TargetAudience: notices.AudienceStudents,
		PublishNow:     true,
	})

	// Notice 2: Admin - Staff
	_, _ = svc.CreateNotice(ctx, notices.CreateNoticeCommand{
		TenantID:       "ten_default",
		AuthorID:       "usr_1",
		Title:          "Staff Quarterly Meeting",
		Content:        "Meeting at Auditorium",
		Category:       notices.CategoryAdministrative,
		Priority:       notices.PriorityNormal,
		TargetAudience: notices.AudienceStaff,
		PublishNow:     true,
	})

	// Filter by Category = EXAMINATION
	catExam := notices.CategoryExamination
	examNotices, total, err := svc.ListNotices(ctx, notices.NoticeFilter{
		TenantID: "ten_default",
		Category: &catExam,
	})
	if err != nil {
		t.Fatalf("ListNotices failed: %v", err)
	}
	if total != 1 || len(examNotices) != 1 {
		t.Errorf("expected 1 exam notice, got total=%d len=%d", total, len(examNotices))
	}
}
