/**
 * BLOCK_AUDIT_HASHER_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Cryptographic SHA-256 hash chaining engine for tamper-evident ledger integrity.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const GenesisHash = "0000000000000000000000000000000000000000000000000000000000000000"

// Hasher provides cryptographic hashing and integrity verification for audit chains.
type Hasher interface {
	ComputeRecordHash(record *AuditLog, prevHash string) string
	VerifyRecord(record *AuditLog, expectedPrevHash string) bool
	VerifyChain(records []*AuditLog, initialPrevHash string) *VerificationResult
}

type sha256Hasher struct{}

// NewSHA256Hasher constructs a production SHA-256 audit chain hasher.
func NewSHA256Hasher() Hasher {
	return &sha256Hasher{}
}

// ComputeRecordHash calculates the deterministic SHA-256 digest of an audit entry.
func (h *sha256Hasher) ComputeRecordHash(record *AuditLog, prevHash string) string {
	hasher := sha256.New()

	// Canonical serialization format
	hasher.Write([]byte(record.ID))
	hasher.Write([]byte("|"))
	hasher.Write([]byte(record.TenantID))
	hasher.Write([]byte("|"))

	if record.ActorID != nil {
		hasher.Write([]byte(*record.ActorID))
	}
	hasher.Write([]byte("|"))
	hasher.Write([]byte(string(record.ActorType)))
	hasher.Write([]byte("|"))

	if record.ActorRole != nil {
		hasher.Write([]byte(*record.ActorRole))
	}
	hasher.Write([]byte("|"))
	hasher.Write([]byte(record.Action))
	hasher.Write([]byte("|"))
	hasher.Write([]byte(record.ResourceType))
	hasher.Write([]byte("|"))

	if record.ResourceID != nil {
		hasher.Write([]byte(*record.ResourceID))
	}
	hasher.Write([]byte("|"))
	hasher.Write([]byte(string(record.Status)))
	hasher.Write([]byte("|"))

	if record.StatusCode != nil {
		hasher.Write([]byte(strconv.Itoa(*record.StatusCode)))
	}
	hasher.Write([]byte("|"))

	if record.IPAddress != nil {
		hasher.Write([]byte(*record.IPAddress))
	}
	hasher.Write([]byte("|"))

	if record.TraceID != nil {
		hasher.Write([]byte(*record.TraceID))
	}
	hasher.Write([]byte("|"))

	if len(record.Metadata) > 0 {
		hasher.Write(record.Metadata)
	}
	hasher.Write([]byte("|"))

	if len(record.Changes) > 0 {
		hasher.Write(record.Changes)
	}
	hasher.Write([]byte("|"))
	hasher.Write([]byte(record.CreatedAt.UTC().Format(time.RFC3339Nano)))
	hasher.Write([]byte("|"))
	hasher.Write([]byte(prevHash))

	return hex.EncodeToString(hasher.Sum(nil))
}

// VerifyRecord checks if a single record's hash matches its contents and expected previous hash.
func (h *sha256Hasher) VerifyRecord(record *AuditLog, expectedPrevHash string) bool {
	computed := h.ComputeRecordHash(record, expectedPrevHash)
	return computed == record.Hash
}

// VerifyChain executes sequential verification over a ordered slice of audit records.
func (h *sha256Hasher) VerifyChain(records []*AuditLog, initialPrevHash string) *VerificationResult {
	result := &VerificationResult{
		TotalVerified: len(records),
		IsChainIntact: true,
		VerifiedAt:    time.Now().UTC(),
		Errors:        make([]VerificationError, 0),
	}

	if len(records) == 0 {
		result.HeadHash = initialPrevHash
		return result
	}

	result.TenantID = records[0].TenantID
	start := records[0].CreatedAt
	end := records[len(records)-1].CreatedAt
	result.TimeRangeStart = &start
	result.TimeRangeEnd = &end

	currentPrevHash := initialPrevHash
	if currentPrevHash == "" {
		currentPrevHash = GenesisHash
	}

	for i, record := range records {
		recPrev := ""
		if record.PrevHash != nil {
			recPrev = *record.PrevHash
		} else {
			recPrev = GenesisHash
		}

		// Verify previous hash pointer continuity
		if recPrev != currentPrevHash {
			result.IsChainIntact = false
			result.Errors = append(result.Errors, VerificationError{
				RecordID:     record.ID,
				Index:        i,
				ExpectedHash: currentPrevHash,
				ActualHash:   recPrev,
				PrevHash:     currentPrevHash,
				Reason:       fmt.Sprintf("broken link: record %s expected prev_hash %s, got %s", record.ID, currentPrevHash, recPrev),
			})
		}

		// Verify record content integrity
		computedHash := h.ComputeRecordHash(record, recPrev)
		if computedHash != record.Hash {
			result.IsChainIntact = false
			result.Errors = append(result.Errors, VerificationError{
				RecordID:     record.ID,
				Index:        i,
				ExpectedHash: computedHash,
				ActualHash:   record.Hash,
				PrevHash:     recPrev,
				Reason:       fmt.Sprintf("payload tampered: record %s hash mismatch", record.ID),
			})
		}

		currentPrevHash = record.Hash
	}

	result.HeadHash = currentPrevHash
	return result
}
