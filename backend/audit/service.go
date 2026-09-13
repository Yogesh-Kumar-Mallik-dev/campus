/**
 * BLOCK_AUDIT_SERVICE_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Domain service orchestrating audit logging, chain verification, and compliance summaries.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sort"
	"strings"
	"sync"
	"time"
)

// Service defines the public domain contract for the Central Audit & Compliance System.
type Service interface {
	Record(ctx context.Context, req RecordAuditRequest) (*AuditLog, error)
	RecordAsync(req RecordAuditRequest) error
	QueryLogs(ctx context.Context, filter AuditFilter) ([]*AuditLog, int64, error)
	GetLogByID(ctx context.Context, tenantID, id string) (*AuditLog, error)
	VerifyChainIntegrity(ctx context.Context, tenantID string, fromTime, toTime *time.Time) (*VerificationResult, error)
	CreateVerificationCheckpoint(ctx context.Context, tenantID string) (*AuditVerificationCheckpoint, error)
	GenerateComplianceReport(ctx context.Context, req ComplianceReportRequest) (*ComplianceReportSummary, error)
}

type service struct {
	repo       Repository
	hasher     Hasher
	subscriber Subscriber
	chainMu    sync.Mutex // Serializes record hash computation per tenant
}

// NewService constructs a new Central Audit & Compliance domain service.
func NewService(repo Repository, hasher Hasher, subscriber Subscriber) Service {
	return &service{
		repo:       repo,
		hasher:     hasher,
		subscriber: subscriber,
	}
}

// Record creates and immediately commits a single cryptographically chained audit entry.
func (s *service) Record(ctx context.Context, req RecordAuditRequest) (*AuditLog, error) {
	if strings.TrimSpace(req.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(req.Action) == "" {
		return nil, ErrActionRequired
	}
	if strings.TrimSpace(req.ResourceType) == "" {
		return nil, ErrResourceTypeRequired
	}

	s.chainMu.Lock()
	defer s.chainMu.Unlock()

	latest, err := s.repo.GetLatestRecord(ctx, req.TenantID)
	if err != nil {
		return nil, NewDomainError("AUDIT_FETCH_FAILED", "failed to read previous audit record", err)
	}

	prevHash := GenesisHash
	if latest != nil && latest.Hash != "" {
		prevHash = latest.Hash
	}

	createdAt := time.Now().UTC()
	if req.Timestamp != nil {
		createdAt = req.Timestamp.UTC()
	}

	id := generateID()
	prev := prevHash

	logEntry := &AuditLog{
		ID:           id,
		TenantID:     req.TenantID,
		ActorID:      req.ActorID,
		ActorType:    req.ActorType,
		ActorRole:    req.ActorRole,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Status:       req.Status,
		StatusCode:   req.StatusCode,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		TraceID:      req.TraceID,
		Metadata:     req.Metadata,
		Changes:      req.Changes,
		PrevHash:     &prev,
		CreatedAt:    createdAt,
	}

	logEntry.Hash = s.hasher.ComputeRecordHash(logEntry, prev)

	if err := s.repo.Append(ctx, logEntry); err != nil {
		return nil, NewDomainError("AUDIT_APPEND_FAILED", "failed to persist audit log", err)
	}

	return logEntry, nil
}

// RecordAsync enqueues an audit event to the asynchronous ingestion worker.
func (s *service) RecordAsync(req RecordAuditRequest) error {
	if strings.TrimSpace(req.TenantID) == "" {
		return ErrTenantRequired
	}
	if strings.TrimSpace(req.Action) == "" {
		return ErrActionRequired
	}
	if strings.TrimSpace(req.ResourceType) == "" {
		return ErrResourceTypeRequired
	}

	if s.subscriber == nil {
		// Fallback to synchronous if subscriber not configured
		_, err := s.Record(context.Background(), req)
		return err
	}

	return s.subscriber.Enqueue(req)
}

// QueryLogs executes a filtered search over audit records with pagination.
func (s *service) QueryLogs(ctx context.Context, filter AuditFilter) ([]*AuditLog, int64, error) {
	if strings.TrimSpace(filter.TenantID) == "" {
		return nil, 0, ErrTenantRequired
	}

	if filter.Limit <= 0 {
		filter.Limit = 50
	} else if filter.Limit > 200 {
		filter.Limit = 200
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return s.repo.Query(ctx, filter)
}

// GetLogByID retrieves a single audit entry.
func (s *service) GetLogByID(ctx context.Context, tenantID, id string) (*AuditLog, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(id) == "" {
		return nil, ErrAuditNotFound
	}

	return s.repo.GetByID(ctx, tenantID, id)
}

// VerifyChainIntegrity verifies the cryptographic SHA-256 chain across records.
func (s *service) VerifyChainIntegrity(ctx context.Context, tenantID string, fromTime, toTime *time.Time) (*VerificationResult, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, ErrTenantRequired
	}

	records, err := s.repo.GetRecordsInRange(ctx, tenantID, fromTime, toTime)
	if err != nil {
		return nil, NewDomainError("VERIFY_RANGE_FAILED", "failed to fetch audit records for verification", err)
	}

	if len(records) == 0 {
		return &VerificationResult{
			TenantID:      tenantID,
			TotalVerified: 0,
			IsChainIntact: true,
			HeadHash:      GenesisHash,
			VerifiedAt:    time.Now().UTC(),
		}, nil
	}

	// Determine starting previous hash
	initialPrevHash := GenesisHash
	if records[0].PrevHash != nil {
		initialPrevHash = *records[0].PrevHash
	}

	return s.hasher.VerifyChain(records, initialPrevHash), nil
}

// CreateVerificationCheckpoint creates an immutable periodic snapshot anchor.
func (s *service) CreateVerificationCheckpoint(ctx context.Context, tenantID string) (*AuditVerificationCheckpoint, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, ErrTenantRequired
	}

	records, err := s.repo.GetRecordsInRange(ctx, tenantID, nil, nil)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, ErrAuditNotFound
	}

	initialPrevHash := GenesisHash
	if records[0].PrevHash != nil {
		initialPrevHash = *records[0].PrevHash
	}

	verification := s.hasher.VerifyChain(records, initialPrevHash)

	checkpoint := &AuditVerificationCheckpoint{
		ID:                 generateCheckpointID(),
		TenantID:           tenantID,
		StartAuditID:       records[0].ID,
		EndAuditID:         records[len(records)-1].ID,
		RecordCount:        len(records),
		HeadHash:           verification.HeadHash,
		IsVerified:         verification.IsChainIntact,
		VerifiedAt:         time.Now().UTC(),
		VerificationErrors: verification.Errors,
	}

	if err := s.repo.CreateCheckpoint(ctx, checkpoint); err != nil {
		return nil, err
	}

	return checkpoint, nil
}

// GenerateComplianceReport compiles accreditation and regulatory compliance analytics.
func (s *service) GenerateComplianceReport(ctx context.Context, req ComplianceReportRequest) (*ComplianceReportSummary, error) {
	if strings.TrimSpace(req.TenantID) == "" {
		return nil, ErrTenantRequired
	}

	from := req.FromTime
	to := req.ToTime
	records, err := s.repo.GetRecordsInRange(ctx, req.TenantID, &from, &to)
	if err != nil {
		return nil, err
	}

	var (
		totalEvents       int64
		successfulEvents  int64
		failedEvents      int64
		securityIncidents int64
		actionCounts      = make(map[string]int64)
		actorCounts       = make(map[string]*ActorActivitySummary)
	)

	for _, rec := range records {
		totalEvents++
		if rec.Status == StatusSuccess {
			successfulEvents++
		} else {
			failedEvents++
		}

		// Security incidents include auth breach, lockouts, failures, and permission escalations
		if strings.HasPrefix(rec.Action, "auth:token:breach") ||
			strings.HasPrefix(rec.Action, "auth:lockout") ||
			strings.HasPrefix(rec.Action, "security:") ||
			rec.Status == StatusFailure {
			securityIncidents++
		}

		actionCounts[rec.Action]++

		actorKey := "anonymous"
		if rec.ActorID != nil {
			actorKey = *rec.ActorID
		}
		role := "unknown"
		if rec.ActorRole != nil {
			role = *rec.ActorRole
		}

		if item, exists := actorCounts[actorKey]; exists {
			item.Count++
		} else {
			actorCounts[actorKey] = &ActorActivitySummary{
				ActorID:   actorKey,
				ActorRole: role,
				Count:     1,
			}
		}
	}

	// Sort top actions
	topActions := make([]ActionStats, 0, len(actionCounts))
	for action, count := range actionCounts {
		topActions = append(topActions, ActionStats{Action: action, Count: count})
	}
	sort.Slice(topActions, func(i, j int) bool {
		return topActions[i].Count > topActions[j].Count
	})
	if len(topActions) > 10 {
		topActions = topActions[:10]
	}

	// Sort top actors
	topActors := make([]ActorActivitySummary, 0, len(actorCounts))
	for _, summary := range actorCounts {
		topActors = append(topActors, *summary)
	}
	sort.Slice(topActors, func(i, j int) bool {
		return topActors[i].Count > topActors[j].Count
	})
	if len(topActors) > 10 {
		topActors = topActors[:10]
	}

	// Verify chain integrity over reporting window
	initialPrevHash := GenesisHash
	if len(records) > 0 && records[0].PrevHash != nil {
		initialPrevHash = *records[0].PrevHash
	}
	verification := s.hasher.VerifyChain(records, initialPrevHash)

	return &ComplianceReportSummary{
		TenantID:          req.TenantID,
		Framework:         req.Framework,
		FromTime:          req.FromTime,
		ToTime:            req.ToTime,
		TotalEvents:       totalEvents,
		SuccessfulEvents:  successfulEvents,
		FailedEvents:      failedEvents,
		SecurityIncidents: securityIncidents,
		TopActions:        topActions,
		TopActors:         topActors,
		ChainIntegrity:    verification.IsChainIntact,
		GeneratedAt:       time.Now().UTC(),
	}, nil
}

func generateID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return "aud_" + hex.EncodeToString(b)
}

func generateCheckpointID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return "chk_" + hex.EncodeToString(b)
}
