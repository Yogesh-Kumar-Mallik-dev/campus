/**
 * BLOCK_WHISTLEBLOWER_TEST_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   Unit test suite covering zero-knowledge report submission, token lookup, anonymous messaging, and state machine.
 */

package whistleblower_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"campus/backend/whistleblower"
)

func setupWhistleblowerService() (whistleblower.Service, *whistleblower.MockRepository) {
	mockRepo := whistleblower.NewMockRepository()
	srv := whistleblower.NewService(mockRepo, nil)
	return srv, mockRepo
}

func TestWhistleblower_ReportSubmissionAndTokenLookup(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupWhistleblowerService()
	tenantID := "tenant-alpha"

	// 1. Submit Anonymous Report
	report, rawToken, err := srv.SubmitReport(ctx, whistleblower.SubmitReportRequest{
		TenantID:             tenantID,
		Category:             whistleblower.CategoryAntiRagging,
		Severity:             whistleblower.SeverityCritical,
		Title:                "Late night ragging in Hostel Block C 4th Floor",
		DescriptionEncrypted: "Seniors forcing 1st year students to do sit-ups past 01:00 AM.",
	})
	if err != nil {
		t.Fatalf("expected report submission success, got: %v", err)
	}

	if report.ID == "" || report.ReportNumber != "WB-2026-00001" || report.Status != whistleblower.StatusSubmitted {
		t.Fatalf("unexpected report fields: %+v", report)
	}

	if !strings.HasPrefix(rawToken, "WB-tok-") {
		t.Fatalf("expected secret token format WB-tok-*, got: %s", rawToken)
	}

	// 2. Lookup Report using correct Secret Tracking Token
	lookedUpReport, messages, err := srv.LookupReportByToken(ctx, tenantID, rawToken)
	if err != nil || lookedUpReport.ID != report.ID {
		t.Fatalf("expected successful lookup with token, got err: %v, report: %+v", err, lookedUpReport)
	}
	if len(messages) != 0 {
		t.Fatalf("expected 0 initial messages, got: %d", len(messages))
	}

	// 3. Lookup with Wrong Token -> REJECTED
	_, _, err = srv.LookupReportByToken(ctx, tenantID, "WB-tok-invalidfakehash123")
	if !errors.Is(err, whistleblower.ErrReportNotFound) {
		t.Fatalf("expected ErrReportNotFound for wrong token, got: %v", err)
	}
}

func TestWhistleblower_TwoWayAnonymousMessaging(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupWhistleblowerService()
	tenantID := "tenant-alpha"

	report, rawToken, _ := srv.SubmitReport(ctx, whistleblower.SubmitReportRequest{
		TenantID:             tenantID,
		Category:             whistleblower.CategoryFinancialFraud,
		Severity:             whistleblower.SeverityHigh,
		Title:                "Fake billing in canteen supplies",
		DescriptionEncrypted: "Overbilling detected in grocery receipts.",
	})

	// 1. Committee requests more evidence
	_, err := srv.AddMessage(ctx, whistleblower.AddMessageRequest{
		TenantID:         tenantID,
		ReportID:         report.ID,
		SenderType:       whistleblower.SenderInvestigationCommittee,
		MessageEncrypted: "Can you provide scanned copies or dates of the affected invoices?",
		ActorID:          "committee-head-1",
	})
	if err != nil {
		t.Fatalf("committee message failed: %v", err)
	}

	// 2. Anonymous reporter replies using Secret Token
	_, err = srv.AddMessage(ctx, whistleblower.AddMessageRequest{
		TenantID:         tenantID,
		ReportID:         report.ID,
		SenderType:       whistleblower.SenderAnonymousReporter,
		MessageEncrypted: "Invoices are dated 10th and 14th of August from Vendor XYZ.",
		RawTrackingToken: rawToken,
	})
	if err != nil {
		t.Fatalf("reporter message failed: %v", err)
	}

	// 3. Unauthorized entity tries to message without valid token -> REJECTED
	_, err = srv.AddMessage(ctx, whistleblower.AddMessageRequest{
		TenantID:         tenantID,
		ReportID:         report.ID,
		SenderType:       whistleblower.SenderAnonymousReporter,
		MessageEncrypted: "Fake spoof message",
		RawTrackingToken: "WB-tok-fake-token",
	})
	if !errors.Is(err, whistleblower.ErrReportNotFound) {
		t.Fatalf("expected ErrReportNotFound on wrong token, got: %v", err)
	}

	// 4. Verify message history
	_, messages, _ := srv.LookupReportByToken(ctx, tenantID, rawToken)
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages in thread, got: %d", len(messages))
	}
}

func TestWhistleblower_InvestigationLifecycleAndResolution(t *testing.T) {
	ctx := context.Background()
	srv, _ := setupWhistleblowerService()
	tenantID := "tenant-alpha"

	report, _, _ := srv.SubmitReport(ctx, whistleblower.SubmitReportRequest{
		TenantID:             tenantID,
		Category:             whistleblower.CategoryAcademicCorruption,
		Severity:             whistleblower.SeverityHigh,
		Title:                "Unauthorized question paper leak rumor",
		DescriptionEncrypted: "Circulating on Telegram groups prior to exam.",
	})

	// 1. Assign Investigator -> Moves to UNDER_INVESTIGATION
	assigned, err := srv.AssignInvestigator(ctx, tenantID, report.ID, "officer-vigilance-1", "admin-1")
	if err != nil || assigned.Status != whistleblower.StatusUnderInvestigation {
		t.Fatalf("expected UNDER_INVESTIGATION status, got: %+v, err: %v", assigned, err)
	}

	// 2. Evidence Requested
	evidReq, err := srv.UpdateReportStatus(ctx, whistleblower.UpdateReportStatusRequest{
		TenantID: tenantID,
		ReportID: report.ID,
		Status:   whistleblower.StatusEvidenceRequested,
		ActorID:  "officer-vigilance-1",
	})
	if err != nil || evidReq.Status != whistleblower.StatusEvidenceRequested {
		t.Fatalf("expected EVIDENCE_REQUESTED status, got: %+v, err: %v", evidReq, err)
	}

	// 3. Resolve Report with Findings Summary
	resolved, err := srv.UpdateReportStatus(ctx, whistleblower.UpdateReportStatusRequest{
		TenantID:        tenantID,
		ReportID:        report.ID,
		Status:          whistleblower.StatusResolved,
		FindingsSummary: "Investigation completed. Telegram channel admin identified and disciplinary action initiated.",
		ActorID:         "officer-vigilance-1",
	})
	if err != nil || resolved.Status != whistleblower.StatusResolved || resolved.ResolvedAt == nil {
		t.Fatalf("expected RESOLVED status with findings, got: %+v, err: %v", resolved, err)
	}

	// 4. Modifying closed report -> REJECTED
	_, err = srv.UpdateReportStatus(ctx, whistleblower.UpdateReportStatusRequest{
		TenantID: tenantID,
		ReportID: report.ID,
		Status:   whistleblower.StatusUnderInvestigation,
		ActorID:  "officer-vigilance-1",
	})
	if !errors.Is(err, whistleblower.ErrReportAlreadyClosed) {
		t.Fatalf("expected ErrReportAlreadyClosed, got: %v", err)
	}
}
