/**
 * BLOCK_WHISTLEBLOWER_REPO_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   Repository interface contracts for confidential reports, anonymous messages, and evidence.
 */

package whistleblower

import "context"

type Repository interface {
	CreateReport(ctx context.Context, report *WhistleblowerReport) error
	GetReportByID(ctx context.Context, tenantID, id string) (*WhistleblowerReport, error)
	GetReportByTrackingHash(ctx context.Context, tenantID, trackingHash string) (*WhistleblowerReport, error)
	ListReports(ctx context.Context, filter WhistleblowerReportFilter) ([]WhistleblowerReport, error)
	UpdateReport(ctx context.Context, report *WhistleblowerReport) error
	GetNextReportSequence(ctx context.Context, tenantID string) (int, error)

	CreateMessage(ctx context.Context, message *WhistleblowerMessage) error
	ListMessages(ctx context.Context, tenantID, reportID string) ([]WhistleblowerMessage, error)

	CreateEvidence(ctx context.Context, evidence *WhistleblowerEvidence) error
	ListEvidences(ctx context.Context, tenantID, reportID string) ([]WhistleblowerEvidence, error)
}
