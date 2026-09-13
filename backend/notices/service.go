/**
 * BLOCK_NOTICES_SERVICE_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   Core business orchestration for targeted broadcasts, priority pinning, read receipts, and circular attachments.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package notices

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

type AttachmentInput struct {
	FileName string `json:"file_name"`
	FileKey  string `json:"file_key"`
	FileSize int    `json:"file_size"`
	MimeType string `json:"mime_type"`
}

type CreateNoticeCommand struct {
	TenantID        string            `json:"tenant_id"`
	AuthorID        string            `json:"author_id"`
	Title           string            `json:"title"`
	Content         string            `json:"content"`
	Category        NoticeCategory    `json:"category"`
	Priority        NoticePriority    `json:"priority"`
	TargetAudience  TargetAudience    `json:"target_audience"`
	TargetDeptID    *string           `json:"target_dept_id,omitempty"`
	TargetProgramID *string           `json:"target_program_id,omitempty"`
	TargetCohortID  *string           `json:"target_cohort_id,omitempty"`
	IsPinned        bool              `json:"is_pinned"`
	PublishNow      bool              `json:"publish_now"`
	ExpiresAt       *time.Time        `json:"expires_at,omitempty"`
	Attachments     []AttachmentInput `json:"attachments,omitempty"`
}

type AcknowledgeNoticeCommand struct {
	TenantID string  `json:"tenant_id"`
	NoticeID string  `json:"notice_id"`
	UserID   string  `json:"user_id"`
	DeviceID *string `json:"device_id,omitempty"`
}

type AttachDocumentCommand struct {
	TenantID string `json:"tenant_id"`
	NoticeID string `json:"notice_id"`
	FileName string `json:"file_name"`
	FileKey  string `json:"file_key"`
	FileSize int    `json:"file_size"`
	MimeType string `json:"mime_type"`
}

type Service struct {
	noticeRepo NoticeRepository
	attRepo    AttachmentRepository
	ackRepo    AcknowledgementRepository
	auditSub   audit.Subscriber
}

func NewService(
	noticeRepo NoticeRepository,
	attRepo AttachmentRepository,
	ackRepo AcknowledgementRepository,
	auditSub audit.Subscriber,
) *Service {
	return &Service{
		noticeRepo: noticeRepo,
		attRepo:    attRepo,
		ackRepo:    ackRepo,
		auditSub:   auditSub,
	}
}

// CreateNotice creates a new campus announcement draft or publishes it immediately.
func (s *Service) CreateNotice(ctx context.Context, cmd CreateNoticeCommand) (*CampusNotice, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(cmd.Title) == "" || strings.TrimSpace(cmd.Content) == "" {
		return nil, NewDomainError("INVALID_NOTICE", "title and content are required", ErrInvalidInput)
	}

	if cmd.Category == "" {
		cmd.Category = CategoryGeneral
	}
	if cmd.Priority == "" {
		cmd.Priority = PriorityNormal
	}
	if cmd.TargetAudience == "" {
		cmd.TargetAudience = AudienceAll
	}

	now := time.Now().UTC()
	status := StatusDraft
	if cmd.PublishNow {
		status = StatusPublished
	}

	baseSlug := Slugify(cmd.Title)
	if baseSlug == "" {
		baseSlug = "notice"
	}
	slug := fmt.Sprintf("%s-%s", baseSlug, generateID("")[0:6])
	noticeID := generateID("not_")

	attachments := make([]NoticeAttachment, 0, len(cmd.Attachments))
	for _, a := range cmd.Attachments {
		attachments = append(attachments, NoticeAttachment{
			ID:        generateID("att_"),
			TenantID:  cmd.TenantID,
			NoticeID:  noticeID,
			FileName:  a.FileName,
			FileKey:   a.FileKey,
			FileSize:  a.FileSize,
			MimeType:  a.MimeType,
			CreatedAt: now,
		})
	}

	notice := &CampusNotice{
		ID:              noticeID,
		TenantID:        cmd.TenantID,
		AuthorID:        cmd.AuthorID,
		Title:           strings.TrimSpace(cmd.Title),
		Slug:            slug,
		Content:         cmd.Content,
		Category:        cmd.Category,
		Priority:        cmd.Priority,
		Status:          status,
		TargetAudience:  cmd.TargetAudience,
		TargetDeptID:    cmd.TargetDeptID,
		TargetProgramID: cmd.TargetProgramID,
		TargetCohortID:  cmd.TargetCohortID,
		IsPinned:        cmd.IsPinned,
		PublishAt:       now,
		ExpiresAt:       cmd.ExpiresAt,
		ViewCount:       0,
		Attachments:     attachments,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.noticeRepo.CreateNotice(ctx, notice); err != nil {
		return nil, err
	}

	for _, a := range attachments {
		_ = s.attRepo.CreateAttachment(ctx, &a)
	}

	s.emitAudit(cmd.TenantID, cmd.AuthorID, "USER", "notices:notice:created", "campus_notice", notice.ID, audit.StatusSuccess, map[string]interface{}{
		"title":     notice.Title,
		"priority":  notice.Priority,
		"published": cmd.PublishNow,
	})

	return notice, nil
}

// PublishNotice transitions a draft notice to published status.
func (s *Service) PublishNotice(ctx context.Context, tenantID, noticeID, authorID string) (*CampusNotice, error) {
	notice, err := s.noticeRepo.GetNoticeByID(ctx, tenantID, noticeID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := notice.Publish(now); err != nil {
		return nil, err
	}

	if err := s.noticeRepo.UpdateNotice(ctx, notice); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, authorID, "USER", "notices:notice:published", "campus_notice", notice.ID, audit.StatusSuccess, map[string]interface{}{
		"title":    notice.Title,
		"priority": notice.Priority,
	})

	return notice, nil
}

// ArchiveNotice archives a published notice.
func (s *Service) ArchiveNotice(ctx context.Context, tenantID, noticeID, authorID string) (*CampusNotice, error) {
	notice, err := s.noticeRepo.GetNoticeByID(ctx, tenantID, noticeID)
	if err != nil {
		return nil, err
	}

	if err := notice.Archive(); err != nil {
		return nil, err
	}

	if err := s.noticeRepo.UpdateNotice(ctx, notice); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, authorID, "USER", "notices:notice:archived", "campus_notice", notice.ID, audit.StatusSuccess, nil)
	return notice, nil
}

// AcknowledgeNotice records a user's read receipt and confirmation of a notice.
func (s *Service) AcknowledgeNotice(ctx context.Context, cmd AcknowledgeNoticeCommand) (*NoticeAcknowledgement, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(cmd.NoticeID) == "" || strings.TrimSpace(cmd.UserID) == "" {
		return nil, NewDomainError("INVALID_ACK", "notice_id and user_id are required", ErrInvalidInput)
	}

	now := time.Now().UTC()
	ack := &NoticeAcknowledgement{
		ID:             generateID("ack_"),
		TenantID:       cmd.TenantID,
		NoticeID:       cmd.NoticeID,
		UserID:         cmd.UserID,
		ReadAt:         now,
		AcknowledgedAt: &now,
		DeviceID:       cmd.DeviceID,
		CreatedAt:      now,
	}

	if err := s.ackRepo.RecordAcknowledgement(ctx, ack); err != nil {
		return nil, err
	}

	return ack, nil
}

// AttachDocument adds a circular document to an existing notice.
func (s *Service) AttachDocument(ctx context.Context, cmd AttachDocumentCommand) (*NoticeAttachment, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(cmd.NoticeID) == "" || strings.TrimSpace(cmd.FileName) == "" {
		return nil, NewDomainError("INVALID_ATTACHMENT", "notice_id and file_name are required", ErrInvalidInput)
	}

	now := time.Now().UTC()
	att := &NoticeAttachment{
		ID:        generateID("att_"),
		TenantID:  cmd.TenantID,
		NoticeID:  cmd.NoticeID,
		FileName:  cmd.FileName,
		FileKey:   cmd.FileKey,
		FileSize:  cmd.FileSize,
		MimeType:  cmd.MimeType,
		CreatedAt: now,
	}

	if err := s.attRepo.CreateAttachment(ctx, att); err != nil {
		return nil, err
	}

	return att, nil
}

// GetNotice fetches a notice with attachments and optional view count increment.
func (s *Service) GetNotice(ctx context.Context, tenantID, noticeID string, recordView bool) (*CampusNotice, error) {
	notice, err := s.noticeRepo.GetNoticeByID(ctx, tenantID, noticeID)
	if err != nil {
		return nil, err
	}

	if recordView {
		_ = s.noticeRepo.IncrementViewCount(ctx, tenantID, noticeID)
		notice.ViewCount++
	}

	atts, _ := s.attRepo.ListAttachmentsByNotice(ctx, tenantID, noticeID)
	notice.Attachments = atts
	return notice, nil
}

// ListNotices returns filtered notices matching target parameters.
func (s *Service) ListNotices(ctx context.Context, filter NoticeFilter) ([]*CampusNotice, int, error) {
	return s.noticeRepo.ListNotices(ctx, filter)
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
