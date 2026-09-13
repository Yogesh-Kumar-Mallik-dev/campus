/**
 * BLOCK_WHISTLEBLOWER_DOMAIN_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   Domain entities, zero-knowledge cryptographic token generation, and grievance state machine.
 */

package whistleblower

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type WhistleblowerCategory string

const (
	CategoryAntiRagging         WhistleblowerCategory = "ANTI_RAGGING"
	CategoryFinancialFraud       WhistleblowerCategory = "FINANCIAL_FRAUD"
	CategoryAcademicCorruption  WhistleblowerCategory = "ACADEMIC_CORRUPTION"
	CategoryHarassment          WhistleblowerCategory = "HARASSMENT"
	CategorySafetyViolation     WhistleblowerCategory = "SAFETY_VIOLATION"
	CategoryOther               WhistleblowerCategory = "OTHER"
)

type WhistleblowerSeverity string

const (
	SeverityLow      WhistleblowerSeverity = "LOW"
	SeverityMedium   WhistleblowerSeverity = "MEDIUM"
	SeverityHigh     WhistleblowerSeverity = "HIGH"
	SeverityCritical WhistleblowerSeverity = "CRITICAL"
)

type WhistleblowerStatus string

const (
	StatusSubmitted          WhistleblowerStatus = "SUBMITTED"
	StatusUnderInvestigation WhistleblowerStatus = "UNDER_INVESTIGATION"
	StatusEvidenceRequested  WhistleblowerStatus = "EVIDENCE_REQUESTED"
	StatusResolved           WhistleblowerStatus = "RESOLVED"
	StatusRejected           WhistleblowerStatus = "REJECTED"
)

type WhistleblowerSenderType string

const (
	SenderAnonymousReporter     WhistleblowerSenderType = "ANONYMOUS_REPORTER"
	SenderInvestigationCommittee WhistleblowerSenderType = "INVESTIGATION_COMMITTEE"
)

type WhistleblowerReport struct {
	ID                   string                `json:"id"`
	TenantID             string                `json:"tenantId"`
	ReportNumber         string                `json:"reportNumber"`
	TrackingHash         string                `json:"-"` // Never expose hash in JSON
	Category             WhistleblowerCategory `json:"category"`
	Severity             WhistleblowerSeverity `json:"severity"`
	Title                string                `json:"title"`
	DescriptionEncrypted string                `json:"description"`
	Status               WhistleblowerStatus   `json:"status"`
	AssignedOfficerID    *string               `json:"assignedOfficerId,omitempty"`
	FindingsSummary      *string               `json:"findingsSummary,omitempty"`
	SubmittedAt          time.Time             `json:"submittedAt"`
	ResolvedAt           *time.Time            `json:"resolvedAt,omitempty"`
	CreatedAt            time.Time             `json:"createdAt"`
	UpdatedAt            time.Time             `json:"updatedAt"`
}

type WhistleblowerMessage struct {
	ID               string                  `json:"id"`
	TenantID         string                  `json:"tenantId"`
	ReportID         string                  `json:"reportId"`
	SenderType       WhistleblowerSenderType `json:"senderType"`
	MessageEncrypted string                  `json:"message"`
	CreatedAt        time.Time               `json:"createdAt"`
}

type WhistleblowerEvidence struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	ReportID  string    `json:"reportId"`
	FileURL   string    `json:"fileUrl"`
	FileHash  string    `json:"fileHash"`
	MimeType  string    `json:"mimeType"`
	CreatedAt time.Time `json:"createdAt"`
}

type WhistleblowerReportFilter struct {
	TenantID string
	Category *WhistleblowerCategory
	Severity *WhistleblowerSeverity
	Status   *WhistleblowerStatus
}

// GenerateSecretTrackingToken generates a high-entropy secret token and its SHA-256 hash.
func GenerateSecretTrackingToken() (rawToken, tokenHash string, err error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	rawToken = fmt.Sprintf("WB-tok-%s", hex.EncodeToString(bytes))
	tokenHash = HashTrackingToken(rawToken)
	return rawToken, tokenHash, nil
}

// HashTrackingToken computes the SHA-256 digest of a secret tracking token.
func HashTrackingToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

// FormatReportNumber generates deterministic whistleblower report identifier.
func FormatReportNumber(year, seq int) string {
	return fmt.Sprintf("WB-%d-%05d", year, seq)
}

// ValidateReportTransition validates legal investigation state machine progression.
func ValidateReportTransition(current, next WhistleblowerStatus) error {
	if current == next {
		return nil
	}

	if current == StatusResolved || current == StatusRejected {
		return fmt.Errorf("%w: report is in terminal state %s", ErrReportAlreadyClosed, current)
	}

	switch current {
	case StatusSubmitted:
		if next == StatusUnderInvestigation || next == StatusRejected {
			return nil
		}
	case StatusUnderInvestigation:
		if next == StatusEvidenceRequested || next == StatusResolved || next == StatusRejected {
			return nil
		}
	case StatusEvidenceRequested:
		if next == StatusUnderInvestigation || next == StatusResolved || next == StatusRejected {
			return nil
		}
	}

	return fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidReportTransition, current, next)
}
