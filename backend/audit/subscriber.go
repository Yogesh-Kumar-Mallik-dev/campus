/**
 * BLOCK_AUDIT_SUBSCRIBER_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   High-throughput asynchronous event subscriber and batched ledger ingestion queue.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Subscriber defines the async ingestion interface for publishing audit events.
type Subscriber interface {
	Enqueue(req RecordAuditRequest) error
	Start(ctx context.Context)
	Stop()
}

// BatchConfig configures the asynchronous worker behavior.
type BatchConfig struct {
	QueueCapacity int           // e.g. 10000
	BatchSize     int           // e.g. 100
	FlushInterval time.Duration // e.g. 50ms
}

// AsyncSubscriber processes audit logs in the background with batching and chain integrity.
type AsyncSubscriber struct {
	repo    Repository
	hasher  Hasher
	config  BatchConfig
	queue   chan RecordAuditRequest
	stopCh  chan struct{}
	doneCh  chan struct{}
	mu      sync.Mutex
	running bool
}

// NewAsyncSubscriber constructs an asynchronous event subscriber.
func NewAsyncSubscriber(repo Repository, hasher Hasher, config BatchConfig) *AsyncSubscriber {
	if config.QueueCapacity <= 0 {
		config.QueueCapacity = 10000
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}
	if config.FlushInterval <= 0 {
		config.FlushInterval = 50 * time.Millisecond
	}

	return &AsyncSubscriber{
		repo:    repo,
		hasher:  hasher,
		config:  config,
		queue:   make(chan RecordAuditRequest, config.QueueCapacity),
		stopCh:  make(chan struct{}),
		doneCh:  make(chan struct{}),
	}
}

// Enqueue puts a record request into the non-blocking ingestion buffer.
func (s *AsyncSubscriber) Enqueue(req RecordAuditRequest) error {
	select {
	case s.queue <- req:
		return nil
	default:
		return ErrQueueFull
	}
}

// Start launches the background processing goroutine.
func (s *AsyncSubscriber) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go s.workerLoop(ctx)
}

// Stop gracefully signals the worker to flush remaining records and terminate.
func (s *AsyncSubscriber) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	close(s.stopCh)
	s.mu.Unlock()

	<-s.doneCh
}

func (s *AsyncSubscriber) workerLoop(ctx context.Context) {
	defer close(s.doneCh)

	ticker := time.NewTicker(s.config.FlushInterval)
	defer ticker.Stop()

	var batch []RecordAuditRequest

	flush := func() {
		if len(batch) == 0 {
			return
		}
		s.processBatch(ctx, batch)
		batch = make([]RecordAuditRequest, 0, s.config.BatchSize)
	}

	for {
		select {
		case <-s.stopCh:
			// Drain remaining items from queue
			for {
				select {
				case req := <-s.queue:
					batch = append(batch, req)
					if len(batch) >= s.config.BatchSize {
						flush()
					}
				default:
					flush()
					return
				}
			}

		case <-ticker.C:
			flush()

		case req := <-s.queue:
			batch = append(batch, req)
			if len(batch) >= s.config.BatchSize {
				flush()
			}
		}
	}
}

func (s *AsyncSubscriber) processBatch(ctx context.Context, requests []RecordAuditRequest) {
	if len(requests) == 0 {
		return
	}

	// Group requests by tenant to maintain deterministic hash chains per tenant
	tenantGroups := make(map[string][]RecordAuditRequest)
	for _, req := range requests {
		tenantGroups[req.TenantID] = append(tenantGroups[req.TenantID], req)
	}

	for tenantID, reqs := range tenantGroups {
		logs := make([]*AuditLog, 0, len(reqs))

		// Fetch latest record to get the current head hash
		latest, _ := s.repo.GetLatestRecord(ctx, tenantID)
		currentPrevHash := GenesisHash
		if latest != nil && latest.Hash != "" {
			currentPrevHash = latest.Hash
		}

		for _, req := range reqs {
			createdAt := time.Now().UTC()
			if req.Timestamp != nil {
				createdAt = req.Timestamp.UTC()
			}

			id := generateAuditID()
			prev := currentPrevHash

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
			currentPrevHash = logEntry.Hash
			logs = append(logs, logEntry)
		}

		_ = s.repo.AppendBatch(ctx, logs)
	}
}

func generateAuditID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return "aud_" + hex.EncodeToString(b)
}
