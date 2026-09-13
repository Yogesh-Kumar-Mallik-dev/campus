/**
 * BLOCK_NOTICES_MOCK_REPOSITORY_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   In-memory, concurrency-safe mock repositories for notices, attachments, and acknowledgements.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package notices

import (
	"context"
	"strings"
	"sync"
	"time"
)

type MockNoticeRepository struct {
	mu      sync.RWMutex
	notices map[string]*CampusNotice
}

func NewMockNoticeRepository() *MockNoticeRepository {
	return &MockNoticeRepository{
		notices: make(map[string]*CampusNotice),
	}
}

func (m *MockNoticeRepository) CreateNotice(ctx context.Context, n *CampusNotice) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notices[n.TenantID+":"+n.ID] = n
	return nil
}

func (m *MockNoticeRepository) GetNoticeByID(ctx context.Context, tenantID, id string) (*CampusNotice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n, ok := m.notices[tenantID+":"+id]
	if !ok {
		return nil, ErrNoticeNotFound
	}
	c := *n
	return &c, nil
}

func (m *MockNoticeRepository) GetNoticeBySlug(ctx context.Context, tenantID, slug string) (*CampusNotice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, n := range m.notices {
		if n.TenantID == tenantID && strings.EqualFold(n.Slug, slug) {
			c := *n
			return &c, nil
		}
	}
	return nil, ErrNoticeNotFound
}

func (m *MockNoticeRepository) UpdateNotice(ctx context.Context, n *CampusNotice) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := n.TenantID + ":" + n.ID
	if _, ok := m.notices[key]; !ok {
		return ErrNoticeNotFound
	}
	m.notices[key] = n
	return nil
}

func (m *MockNoticeRepository) IncrementViewCount(ctx context.Context, tenantID, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if n, ok := m.notices[tenantID+":"+id]; ok {
		n.ViewCount++
		return nil
	}
	return ErrNoticeNotFound
}

func (m *MockNoticeRepository) ListNotices(ctx context.Context, filter NoticeFilter) ([]*CampusNotice, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now().UTC()
	var results []*CampusNotice

	for _, n := range m.notices {
		if filter.TenantID != "" && n.TenantID != filter.TenantID {
			continue
		}
		if filter.Category != nil && n.Category != *filter.Category {
			continue
		}
		if filter.Priority != nil && n.Priority != *filter.Priority {
			continue
		}
		if filter.Status != nil && n.Status != *filter.Status {
			continue
		}
		if filter.TargetAudience != nil && n.TargetAudience != AudienceAll && n.TargetAudience != *filter.TargetAudience {
			continue
		}
		if filter.TargetDeptID != "" && n.TargetDeptID != nil && *n.TargetDeptID != filter.TargetDeptID {
			continue
		}
		if filter.IsPinned != nil && n.IsPinned != *filter.IsPinned {
			continue
		}
		if filter.ActiveOnly && !n.IsActive(now) {
			continue
		}
		if filter.SearchQuery != "" {
			q := strings.ToLower(filter.SearchQuery)
			if !strings.Contains(strings.ToLower(n.Title), q) && !strings.Contains(strings.ToLower(n.Content), q) {
				continue
			}
		}

		c := *n
		results = append(results, &c)
	}

	total := len(results)
	if filter.Offset >= total {
		return []*CampusNotice{}, total, nil
	}
	end := total
	if filter.Limit > 0 && filter.Offset+filter.Limit < end {
		end = filter.Offset + filter.Limit
	}
	return results[filter.Offset:end], total, nil
}

type MockAttachmentRepository struct {
	mu          sync.RWMutex
	attachments map[string]*NoticeAttachment
}

func NewMockAttachmentRepository() *MockAttachmentRepository {
	return &MockAttachmentRepository{
		attachments: make(map[string]*NoticeAttachment),
	}
}

func (m *MockAttachmentRepository) CreateAttachment(ctx context.Context, att *NoticeAttachment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.attachments[att.TenantID+":"+att.ID] = att
	return nil
}

func (m *MockAttachmentRepository) ListAttachmentsByNotice(ctx context.Context, tenantID, noticeID string) ([]NoticeAttachment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []NoticeAttachment
	for _, a := range m.attachments {
		if a.TenantID == tenantID && a.NoticeID == noticeID {
			res = append(res, *a)
		}
	}
	return res, nil
}

func (m *MockAttachmentRepository) DeleteAttachment(ctx context.Context, tenantID, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := tenantID + ":" + id
	if _, ok := m.attachments[key]; !ok {
		return ErrAttachmentNotFound
	}
	delete(m.attachments, key)
	return nil
}

type MockAcknowledgementRepository struct {
	mu   sync.RWMutex
	acks map[string]*NoticeAcknowledgement
}

func NewMockAcknowledgementRepository() *MockAcknowledgementRepository {
	return &MockAcknowledgementRepository{
		acks: make(map[string]*NoticeAcknowledgement),
	}
}

func (m *MockAcknowledgementRepository) RecordAcknowledgement(ctx context.Context, ack *NoticeAcknowledgement) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.acks[ack.NoticeID+":"+ack.UserID] = ack
	return nil
}

func (m *MockAcknowledgementRepository) GetAcknowledgement(ctx context.Context, tenantID, noticeID, userID string) (*NoticeAcknowledgement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ack, ok := m.acks[noticeID+":"+userID]
	if !ok || ack.TenantID != tenantID {
		return nil, nil
	}
	c := *ack
	return &c, nil
}

func (m *MockAcknowledgementRepository) ListAcknowledgementsByNotice(ctx context.Context, tenantID, noticeID string) ([]NoticeAcknowledgement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []NoticeAcknowledgement
	for _, a := range m.acks {
		if a.TenantID == tenantID && a.NoticeID == noticeID {
			res = append(res, *a)
		}
	}
	return res, nil
}
