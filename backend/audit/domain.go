/**
 * BLOCK_AUDIT_DOMAIN_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Core domain entities, filtering models, and compliance reporting schemas.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"encoding/json"
	"time"
)

// ActorType enumerates the entity category executing an audited action.
type ActorType string

const (
	ActorTypeUser      ActorType = "USER"
	ActorTypeSystem    ActorType = "SYSTEM"
	ActorTypeService   ActorType = "SERVICE"
	ActorTypeAnonymous ActorType = "ANONYMOUS"
)

// Status represents the terminal outcome of an audited event.
type Status string

const (
	StatusSuccess   Status = "SUCCESS"
	StatusFailure   Status = "FAILURE"
	StatusAttempted Status = "ATTEMPTED"
)

// AuditLog is the immutable, cryptographically chained audit ledger record.
type AuditLog struct {
	ID           string          `json:"id"`
	TenantID     string          `json:"tenant_id"`
	ActorID      *string         `json:"actor_id,omitempty"`
	ActorType    ActorType       `json:"actor_type"`
	ActorRole    *string         `json:"actor_role,omitempty"`
	Action       string          `json:"action"`        // e.g. "auth:login:success", "auth:token:breach_detected"
	ResourceType string          `json:"resource_type"` // e.g. "user", "session", "invoice"
	ResourceID   *string         `json:"resource_id,omitempty"`
	Status       Status          `json:"status"`
	StatusCode   *int            `json:"status_code,omitempty"`
	IPAddress    *string         `json:"ip_address,omitempty"`
	UserAgent    *string         `json:"user_agent,omitempty"`
	TraceID      *string         `json:"trace_id,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	Changes      json.RawMessage `json:"changes,omitempty"`   // Before / after state diff
	PrevHash     *string         `json:"prev_hash,omitempty"` // SHA-256 hash of previous record
	Hash         string          `json:"hash"`                // SHA-256 cryptographic signature of this record
	CreatedAt    time.Time       `json:"created_at"`
}

// RecordAuditRequest encapsulates input parameters for recording an audit event.
type RecordAuditRequest struct {
	TenantID     string          `json:"tenant_id"`
	ActorID      *string         `json:"actor_id,omitempty"`
	ActorType    ActorType       `json:"actor_type"`
	ActorRole    *string         `json:"actor_role,omitempty"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resource_type"`
	ResourceID   *string         `json:"resource_id,omitempty"`
	Status       Status          `json:"status"`
	StatusCode   *int            `json:"status_code,omitempty"`
	IPAddress    *string         `json:"ip_address,omitempty"`
	UserAgent    *string         `json:"user_agent,omitempty"`
	TraceID      *string         `json:"trace_id,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	Changes      json.RawMessage `json:"changes,omitempty"`
	Timestamp    *time.Time      `json:"timestamp,omitempty"`
}

// AuditFilter provides multi-criteria query parameters for audit log exploration.
type AuditFilter struct {
	TenantID     string     `json:"tenant_id"`
	ActorID      *string    `json:"actor_id,omitempty"`
	ActorType    *ActorType `json:"actor_type,omitempty"`
	Action       *string    `json:"action,omitempty"`
	ResourceType *string    `json:"resource_type,omitempty"`
	ResourceID   *string    `json:"resource_id,omitempty"`
	Status       *Status    `json:"status,omitempty"`
	TraceID      *string    `json:"trace_id,omitempty"`
	FromTime     *time.Time `json:"from_time,omitempty"`
	ToTime       *time.Time `json:"to_time,omitempty"`
	Limit        int        `json:"limit"`
	Offset       int        `json:"offset"`
}

// VerificationError describes a specific tampering or integrity break in the hash chain.
type VerificationError struct {
	RecordID     string `json:"record_id"`
	Index        int    `json:"index"`
	ExpectedHash string `json:"expected_hash"`
	ActualHash   string `json:"actual_hash"`
	PrevHash     string `json:"prev_hash"`
	Reason       string `json:"reason"`
}

// VerificationResult contains the cryptographic audit chain verification status.
type VerificationResult struct {
	TenantID       string              `json:"tenant_id"`
	TotalVerified  int                 `json:"total_verified"`
	IsChainIntact  bool                `json:"is_chain_intact"`
	HeadHash       string              `json:"head_hash"`
	VerifiedAt     time.Time           `json:"verified_at"`
	Errors         []VerificationError `json:"errors,omitempty"`
	TimeRangeStart *time.Time          `json:"time_range_start,omitempty"`
	TimeRangeEnd   *time.Time          `json:"time_range_end,omitempty"`
}

// AuditVerificationCheckpoint is a periodic cryptographic anchor snapshot.
type AuditVerificationCheckpoint struct {
	ID                 string              `json:"id"`
	TenantID           string              `json:"tenant_id"`
	StartAuditID       string              `json:"start_audit_id"`
	EndAuditID         string              `json:"end_audit_id"`
	RecordCount        int                 `json:"record_count"`
	HeadHash           string              `json:"head_hash"`
	IsVerified         bool                `json:"is_verified"`
	VerifiedAt         time.Time           `json:"verified_at"`
	VerificationErrors []VerificationError `json:"verification_errors,omitempty"`
}

// ComplianceReportRequest specifies parameters for accreditation summary generation.
type ComplianceReportRequest struct {
	TenantID  string    `json:"tenant_id"`
	Framework string    `json:"framework"` // e.g. "NAAC", "NIRF", "ISO_27001", "ABET"
	FromTime  time.Time `json:"from_time"`
	ToTime    time.Time `json:"to_time"`
}

// ActionStats summarizes event volumes by category/action.
type ActionStats struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

// ActorActivitySummary captures activity aggregated per actor.
type ActorActivitySummary struct {
	ActorID   string `json:"actor_id"`
	ActorRole string `json:"actor_role"`
	Count     int64  `json:"count"`
}

// ComplianceReportSummary provides formal compliance statistics for accreditation boards.
type ComplianceReportSummary struct {
	TenantID          string                 `json:"tenant_id"`
	Framework         string                 `json:"framework"`
	FromTime          time.Time              `json:"from_time"`
	ToTime            time.Time              `json:"to_time"`
	TotalEvents       int64                  `json:"total_events"`
	SuccessfulEvents  int64                  `json:"successful_events"`
	FailedEvents      int64                  `json:"failed_events"`
	SecurityIncidents int64                  `json:"security_incidents"`
	TopActions        []ActionStats          `json:"top_actions"`
	TopActors         []ActorActivitySummary `json:"top_actors"`
	ChainIntegrity    bool                   `json:"chain_integrity"`
	GeneratedAt       time.Time              `json:"generated_at"`
}
