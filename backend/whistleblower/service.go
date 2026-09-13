/**
 * BLOCK_WHISTLEBLOWER_SERVICE_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   Zero-knowledge grievance intake, anonymous 2-way communication, and audit logging.
 */

package whistleblower

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"campus/backend/audit"
)

type Service interface {
	SubmitReport(ctx context.Context, req SubmitReportRequest) (*WhistleblowerReport, string, error)
	LookupReportByToken(ctx context.Context, tenantID, rawToken string) (*WhistleblowerReport, []WhistleblowerMessage, error)
	AddMessage(ctx context.Context, req AddMessageRequest) (*WhistleblowerMessage, error)
	AssignInvestigator(ctx context.Context, tenantID, reportID, officerID, actorID string) (*WhistleblowerReport, error)
	UpdateReportStatus(ctx context.Context, req UpdateReportStatusRequest) (*WhistleblowerReport, error)
	ListReports(ctx context.Context, filter WhistleblowerReportFilter) ([]WhistleblowerReport, error)
	GetReport(ctx context.Context, tenantID, reportID string) (*WhistleblowerReport, error)
}

type service struct {
	repo     Repository
	auditSub audit.Subscriber
}

func NewService(repo Repository, auditSub audit.Subscriber) Service {
	return &service{
		repo:     repo,
		auditSub: auditSub,
	}
}

func (s *service) logAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
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

type SubmitReportRequest struct {
	TenantID             string                `json:"tenantId"`
	Category             WhistleblowerCategory `json:"category"`
	Severity             WhistleblowerSeverity `json:"severity"`
	Title                string                `json:"title"`
	DescriptionEncrypted string                `json:"description"`
}

func (s *service) SubmitReport(ctx context.Context, req SubmitReportRequest) (*WhistleblowerReport, string, error) {
	if req.Category == "" {
		req.Category = CategoryAntiRagging
	}
	if req.Severity == "" {
		req.Severity = SeverityHigh
	}

	rawToken, tokenHash, err := GenerateSecretTrackingToken()
	if err != nil {
		return nil, "", err
	}

	seq, err := s.repo.GetNextReportSequence(ctx, req.TenantID)
	if err != nil {
		return nil, "", err
	}

	now := time.Now().UTC()
	reportNumber := FormatReportNumber(now.Year(), seq)

	report := &WhistleblowerReport{
		ID:                   uuid.New().String(),
		TenantID:             req.TenantID,
		ReportNumber:         reportNumber,
		TrackingHash:         tokenHash,
		Category:             req.Category,
		Severity:             req.Severity,
		Title:                req.Title,
		DescriptionEncrypted: req.DescriptionEncrypted,
		Status:               StatusSubmitted,
		SubmittedAt:          now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := s.repo.CreateReport(ctx, report); err != nil {
		return nil, "", err
	}

	// Zero-knowledge anonymous audit log
	s.logAudit(req.TenantID, "ANONYMOUS", "SYSTEM", "whistleblower:report:submitted", "whistleblower_report", report.ID, audit.StatusSuccess, map[string]interface{}{
		"reportNumber": report.ReportNumber,
		"category":     report.Category,
		"severity":     report.Severity,
	})

	return report, rawToken, nil
}

func (s *service) LookupReportByToken(ctx context.Context, tenantID, rawToken string) (*WhistleblowerReport, []WhistleblowerMessage, error) {
	if rawToken == "" {
		return nil, nil, ErrInvalidTrackingToken
	}

	tokenHash := HashTrackingToken(rawToken)
	report, err := s.repo.GetReportByTrackingHash(ctx, tenantID, tokenHash)
	if err != nil {
		return nil, nil, ErrReportNotFound
	}

	messages, err := s.repo.ListMessages(ctx, tenantID, report.ID)
	if err != nil {
		return nil, nil, err
	}

	return report, messages, nil
}

type AddMessageRequest struct {
	TenantID         string                  `json:"tenantId"`
	ReportID         string                  `json:"reportId"`
	SenderType       WhistleblowerSenderType `json:"senderType"`
	MessageEncrypted string                  `json:"message"`
	RawTrackingToken string                  `json:"rawTrackingToken,omitempty"` // For anonymous reporter verification
	ActorID          string                  `json:"actorId,omitempty"`          // For committee officer
}

func (s *service) AddMessage(ctx context.Context, req AddMessageRequest) (*WhistleblowerMessage, error) {
	report, err := s.repo.GetReportByID(ctx, req.TenantID, req.ReportID)
	if err != nil {
		return nil, err
	}

	if report.Status == StatusResolved || report.Status == StatusRejected {
		return nil, ErrReportAlreadyClosed
	}

	// Verify anonymous reporter token if sent by reporter
	if req.SenderType == SenderAnonymousReporter {
		if req.RawTrackingToken == "" {
			return nil, ErrInvalidTrackingToken
		}
		computedHash := HashTrackingToken(req.RawTrackingToken)
		if computedHash != report.TrackingHash {
			return nil, ErrReportNotFound
		}
	}

	now := time.Now().UTC()
	msg := &WhistleblowerMessage{
		ID:               uuid.New().String(),
		TenantID:         req.TenantID,
		ReportID:         req.ReportID,
		SenderType:       req.SenderType,
		MessageEncrypted: req.MessageEncrypted,
		CreatedAt:        now,
	}

	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		return nil, err
	}

	report.UpdatedAt = now
	_ = s.repo.UpdateReport(ctx, report)

	return msg, nil
}

func (s *service) AssignInvestigator(ctx context.Context, tenantID, reportID, officerID, actorID string) (*WhistleblowerReport, error) {
	report, err := s.repo.GetReportByID(ctx, tenantID, reportID)
	if err != nil {
		return nil, err
	}

	report.AssignedOfficerID = &officerID
	if report.Status == StatusSubmitted {
		report.Status = StatusUnderInvestigation
	}
	report.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateReport(ctx, report); err != nil {
		return nil, err
	}

	s.logAudit(tenantID, actorID, "STAFF", "whistleblower:report:investigating", "whistleblower_report", report.ID, audit.StatusSuccess, map[string]interface{}{
		"assignedOfficerId": officerID,
	})

	return report, nil
}

type UpdateReportStatusRequest struct {
	TenantID        string              `json:"tenantId"`
	ReportID        string              `json:"reportId"`
	Status          WhistleblowerStatus `json:"status"`
	FindingsSummary string              `json:"findingsSummary,omitempty"`
	ActorID         string              `json:"actorId"`
}

func (s *service) UpdateReportStatus(ctx context.Context, req UpdateReportStatusRequest) (*WhistleblowerReport, error) {
	report, err := s.repo.GetReportByID(ctx, req.TenantID, req.ReportID)
	if err != nil {
		return nil, err
	}

	if err := ValidateReportTransition(report.Status, req.Status); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	report.Status = req.Status
	report.UpdatedAt = now

	if req.FindingsSummary != "" {
		report.FindingsSummary = &req.FindingsSummary
	}
	if req.Status == StatusResolved || req.Status == StatusRejected {
		report.ResolvedAt = &now
	}

	if err := s.repo.UpdateReport(ctx, report); err != nil {
		return nil, err
	}

	s.logAudit(req.TenantID, req.ActorID, "STAFF", "whistleblower:report:status_changed", "whistleblower_report", report.ID, audit.StatusSuccess, map[string]interface{}{
		"status": req.Status,
	})

	return report, nil
}

func (s *service) ListReports(ctx context.Context, filter WhistleblowerReportFilter) ([]WhistleblowerReport, error) {
	return s.repo.ListReports(ctx, filter)
}

func (s *service) GetReport(ctx context.Context, tenantID, reportID string) (*WhistleblowerReport, error) {
	return s.repo.GetReportByID(ctx, tenantID, reportID)
}
