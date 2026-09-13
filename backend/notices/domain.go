/**
 * BLOCK_NOTICES_DOMAIN_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   Domain entities, lifecycle state machines, priority pinning, target audience rules, and acknowledgement records.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package notices

import (
	"regexp"
	"strings"
	"time"
)

type NoticeCategory string

const (
	CategoryAcademic       NoticeCategory = "ACADEMIC"
	CategoryExamination    NoticeCategory = "EXAMINATION"
	CategoryAdministrative NoticeCategory = "ADMINISTRATIVE"
	CategoryHostel         NoticeCategory = "HOSTEL"
	CategoryEvents         NoticeCategory = "EVENTS"
	CategoryEmergency      NoticeCategory = "EMERGENCY"
	CategoryPlacement      NoticeCategory = "PLACEMENT"
	CategoryGeneral        NoticeCategory = "GENERAL"
)

type NoticePriority string

const (
	PriorityLow    NoticePriority = "LOW"
	PriorityNormal NoticePriority = "NORMAL"
	PriorityHigh   NoticePriority = "HIGH"
	PriorityUrgent NoticePriority = "URGENT"
)

type NoticeStatus string

const (
	StatusDraft     NoticeStatus = "DRAFT"
	StatusPublished NoticeStatus = "PUBLISHED"
	StatusArchived  NoticeStatus = "ARCHIVED"
	StatusExpired   NoticeStatus = "EXPIRED"
)

type TargetAudience string

const (
	AudienceAll             TargetAudience = "ALL"
	AudienceStudents        TargetAudience = "STUDENTS"
	AudienceFaculty         TargetAudience = "FACULTY"
	AudienceStaff           TargetAudience = "STAFF"
	AudienceHostelResidents TargetAudience = "HOSTEL_RESIDENTS"
)

type NoticeAttachment struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	NoticeID  string    `json:"notice_id"`
	FileName  string    `json:"file_name"`
	FileKey   string    `json:"file_key"`
	FileSize  int       `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

type NoticeAcknowledgement struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	NoticeID       string     `json:"notice_id"`
	UserID         string     `json:"user_id"`
	ReadAt         time.Time  `json:"read_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	DeviceID       *string    `json:"device_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type CampusNotice struct {
	ID               string                  `json:"id"`
	TenantID         string                  `json:"tenant_id"`
	AuthorID         string                  `json:"author_id"`
	Title            string                  `json:"title"`
	Slug             string                  `json:"slug"`
	Content          string                  `json:"content"`
	Category         NoticeCategory          `json:"category"`
	Priority         NoticePriority          `json:"priority"`
	Status           NoticeStatus            `json:"status"`
	TargetAudience   TargetAudience          `json:"target_audience"`
	TargetDeptID     *string                 `json:"target_dept_id,omitempty"`
	TargetProgramID  *string                 `json:"target_program_id,omitempty"`
	TargetCohortID   *string                 `json:"target_cohort_id,omitempty"`
	IsPinned         bool                    `json:"is_pinned"`
	PublishAt        time.Time               `json:"publish_at"`
	ExpiresAt        *time.Time              `json:"expires_at,omitempty"`
	ViewCount        int                     `json:"view_count"`
	Attachments      []NoticeAttachment      `json:"attachments,omitempty"`
	Acknowledgements []NoticeAcknowledgement `json:"acknowledgements,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

// Publish transitions DRAFT -> PUBLISHED.
func (n *CampusNotice) Publish(now time.Time) error {
	if n.Status == StatusPublished {
		return ErrNoticeAlreadyPublished
	}
	if n.Status == StatusArchived || n.Status == StatusExpired {
		return NewDomainError("INVALID_STATUS", "archived or expired notice cannot be published", ErrInvalidStateTransition)
	}
	n.Status = StatusPublished
	n.PublishAt = now
	n.UpdatedAt = now
	return nil
}

// Archive transitions PUBLISHED/EXPIRED -> ARCHIVED.
func (n *CampusNotice) Archive() error {
	if n.Status == StatusArchived {
		return nil // Idempotent
	}
	n.Status = StatusArchived
	n.UpdatedAt = time.Now().UTC()
	return nil
}

// Expire transitions PUBLISHED -> EXPIRED.
func (n *CampusNotice) Expire() error {
	if n.Status == StatusExpired {
		return nil
	}
	n.Status = StatusExpired
	n.UpdatedAt = time.Now().UTC()
	return nil
}

// IsActive returns whether the notice is currently published and not expired.
func (n *CampusNotice) IsActive(now time.Time) bool {
	if n.Status != StatusPublished {
		return false
	}
	if n.ExpiresAt != nil && now.After(*n.ExpiresAt) {
		return false
	}
	return true
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converts a notice title into an authoritative URL slug.
func Slugify(title string) string {
	clean := strings.ToLower(strings.TrimSpace(title))
	slug := nonAlphanumericRegex.ReplaceAllString(clean, "-")
	return strings.Trim(slug, "-")
}
