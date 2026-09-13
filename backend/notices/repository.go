/**
 * BLOCK_NOTICES_REPOSITORY_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   Data access interface definitions for notices, circular attachments, and user acknowledgements.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package notices

import (
	"context"
)

type NoticeFilter struct {
	TenantID        string
	Category        *NoticeCategory
	Priority        *NoticePriority
	Status          *NoticeStatus
	TargetAudience  *TargetAudience
	TargetDeptID    string
	TargetProgramID string
	TargetCohortID  string
	IsPinned        *bool
	SearchQuery     string
	ActiveOnly      bool
	Limit           int
	Offset          int
}

type NoticeRepository interface {
	CreateNotice(ctx context.Context, n *CampusNotice) error
	GetNoticeByID(ctx context.Context, tenantID, id string) (*CampusNotice, error)
	GetNoticeBySlug(ctx context.Context, tenantID, slug string) (*CampusNotice, error)
	UpdateNotice(ctx context.Context, n *CampusNotice) error
	IncrementViewCount(ctx context.Context, tenantID, id string) error
	ListNotices(ctx context.Context, filter NoticeFilter) ([]*CampusNotice, int, error)
}

type AttachmentRepository interface {
	CreateAttachment(ctx context.Context, att *NoticeAttachment) error
	ListAttachmentsByNotice(ctx context.Context, tenantID, noticeID string) ([]NoticeAttachment, error)
	DeleteAttachment(ctx context.Context, tenantID, id string) error
}

type AcknowledgementRepository interface {
	RecordAcknowledgement(ctx context.Context, ack *NoticeAcknowledgement) error
	GetAcknowledgement(ctx context.Context, tenantID, noticeID, userID string) (*NoticeAcknowledgement, error)
	ListAcknowledgementsByNotice(ctx context.Context, tenantID, noticeID string) ([]NoticeAcknowledgement, error)
}
