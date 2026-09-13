/**
 * BLOCK_API_NOTICES_HANDLER_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for campus notices, circular attachments, and read receipts.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package notices

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/notices"
)

type Handler struct {
	service *notices.Service
}

func NewHandler(service *notices.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/notices", func(r chi.Router) {
		r.Post("/", h.CreateNotice)
		r.Get("/", h.ListNotices)
		r.Get("/{id}", h.GetNotice)
		r.Post("/{id}/publish", h.PublishNotice)
		r.Post("/{id}/archive", h.ArchiveNotice)
		r.Post("/{id}/acknowledge", h.AcknowledgeNotice)
		r.Post("/{id}/attachments", h.AttachDocument)
	})
}

// CreateNotice handles POST /api/v1/notices
func (h *Handler) CreateNotice(w http.ResponseWriter, r *http.Request) {
	var cmd notices.CreateNoticeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	notice, err := h.service.CreateNotice(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v1/notices/"+notice.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"notice": notice,
	})
}

// ListNotices handles GET /api/v1/notices
func (h *Handler) ListNotices(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	filter := notices.NoticeFilter{
		TenantID:        tenantID,
		TargetDeptID:    r.URL.Query().Get("target_dept_id"),
		TargetProgramID: r.URL.Query().Get("target_program_id"),
		TargetCohortID:  r.URL.Query().Get("target_cohort_id"),
		SearchQuery:     r.URL.Query().Get("search"),
		ActiveOnly:      r.URL.Query().Get("active_only") == "true",
		Limit:           limit,
		Offset:          offset,
	}

	if cat := r.URL.Query().Get("category"); cat != "" {
		c := notices.NoticeCategory(strings.ToUpper(cat))
		filter.Category = &c
	}
	if pri := r.URL.Query().Get("priority"); pri != "" {
		p := notices.NoticePriority(strings.ToUpper(pri))
		filter.Priority = &p
	}
	if st := r.URL.Query().Get("status"); st != "" {
		s := notices.NoticeStatus(strings.ToUpper(st))
		filter.Status = &s
	}
	if aud := r.URL.Query().Get("target_audience"); aud != "" {
		a := notices.TargetAudience(strings.ToUpper(aud))
		filter.TargetAudience = &a
	}
	if pinned := r.URL.Query().Get("is_pinned"); pinned != "" {
		isP := pinned == "true"
		filter.IsPinned = &isP
	}

	noticeList, total, err := h.service.ListNotices(r.Context(), filter)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"notices": noticeList,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetNotice handles GET /api/v1/notices/{id}
func (h *Handler) GetNotice(w http.ResponseWriter, r *http.Request) {
	noticeID := chi.URLParam(r, "id")
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id parameter is required", "MISSING_TENANT_ID")
		return
	}

	recordView := r.URL.Query().Get("record_view") != "false"
	notice, err := h.service.GetNotice(r.Context(), tenantID, noticeID, recordView)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"notice": notice,
	})
}

// PublishNotice handles POST /api/v1/notices/{id}/publish
func (h *Handler) PublishNotice(w http.ResponseWriter, r *http.Request) {
	noticeID := chi.URLParam(r, "id")
	var body struct {
		TenantID string `json:"tenant_id"`
		AuthorID string `json:"author_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.TenantID == "" {
		body.TenantID = r.Header.Get("X-Tenant-ID")
	}

	notice, err := h.service.PublishNotice(r.Context(), body.TenantID, noticeID, body.AuthorID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"notice": notice,
	})
}

// ArchiveNotice handles POST /api/v1/notices/{id}/archive
func (h *Handler) ArchiveNotice(w http.ResponseWriter, r *http.Request) {
	noticeID := chi.URLParam(r, "id")
	var body struct {
		TenantID string `json:"tenant_id"`
		AuthorID string `json:"author_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.TenantID == "" {
		body.TenantID = r.Header.Get("X-Tenant-ID")
	}

	notice, err := h.service.ArchiveNotice(r.Context(), body.TenantID, noticeID, body.AuthorID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"notice": notice,
	})
}

// AcknowledgeNotice handles POST /api/v1/notices/{id}/acknowledge
func (h *Handler) AcknowledgeNotice(w http.ResponseWriter, r *http.Request) {
	noticeID := chi.URLParam(r, "id")
	var cmd notices.AcknowledgeNoticeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	cmd.NoticeID = noticeID
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	ack, err := h.service.AcknowledgeNotice(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"acknowledgement": ack,
	})
}

// AttachDocument handles POST /api/v1/notices/{id}/attachments
func (h *Handler) AttachDocument(w http.ResponseWriter, r *http.Request) {
	noticeID := chi.URLParam(r, "id")
	var cmd notices.AttachDocumentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	cmd.NoticeID = noticeID
	if cmd.TenantID == "" {
		cmd.TenantID = r.Header.Get("X-Tenant-ID")
	}

	att, err := h.service.AttachDocument(r.Context(), cmd)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"attachment": att,
	})
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, notices.ErrNoticeNotFound),
		errors.Is(err, notices.ErrAttachmentNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")
	case errors.Is(err, notices.ErrTenantRequired),
		errors.Is(err, notices.ErrInvalidInput),
		errors.Is(err, notices.ErrUnauthorizedAuthor):
		problem.BadRequest(w, r, err.Error(), "INVALID_INPUT")
	case errors.Is(err, notices.ErrNoticeAlreadyPublished),
		errors.Is(err, notices.ErrNoticeExpired),
		errors.Is(err, notices.ErrInvalidStateTransition):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)
	case errors.Is(err, notices.ErrDuplicateNoticeSlug):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")
	default:
		var domErr *notices.DomainError
		if errors.As(err, &domErr) {
			problem.UnprocessableEntity(w, r, domErr.Message, domErr.Code, nil)
			return
		}
		problem.InternalServerError(w, r, "an unexpected notices error occurred")
	}
}
