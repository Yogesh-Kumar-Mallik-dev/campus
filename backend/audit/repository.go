/**
 * BLOCK_AUDIT_REPOSITORY_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Data access interface contracts for immutable audit log ledger and checkpoints.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"context"
	"time"
)

// Repository defines the persistent storage operations for the immutable audit ledger.
type Repository interface {
	// Append writes a single immutable audit entry.
	Append(ctx context.Context, log *AuditLog) error

	// AppendBatch writes a sequence of immutable audit entries atomically or in batch.
	AppendBatch(ctx context.Context, logs []*AuditLog) error

	// GetByID retrieves a specific audit log by tenant and ID.
	GetByID(ctx context.Context, tenantID, id string) (*AuditLog, error)

	// GetLatestRecord returns the most recent audit entry for a given tenant (to get prev_hash).
	GetLatestRecord(ctx context.Context, tenantID string) (*AuditLog, error)

	// Query returns filtered audit logs and the total matching count.
	Query(ctx context.Context, filter AuditFilter) ([]*AuditLog, int64, error)

	// GetRecordsInRange returns all records in chronological order for integrity verification.
	GetRecordsInRange(ctx context.Context, tenantID string, fromTime, toTime *time.Time) ([]*AuditLog, error)

	// CreateCheckpoint stores an immutable verification checkpoint.
	CreateCheckpoint(ctx context.Context, cp *AuditVerificationCheckpoint) error

	// ListCheckpoints returns historic verification snapshots for a tenant.
	ListCheckpoints(ctx context.Context, tenantID string, limit, offset int) ([]*AuditVerificationCheckpoint, error)
}
