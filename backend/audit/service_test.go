/**
 * BLOCK_AUDIT_TESTS_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Exhaustive unit test suite verifying cryptographic hash chains, tamper detection,
 *            asynchronous batching, compliance reports, and multi-tenant ledger isolation.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func ptr[T any](v T) *T {
	return &v
}

func setupAuditTest() (Service, *MockRepository, Hasher) {
	repo := NewMockRepository()
	hasher := NewSHA256Hasher()
	svc := NewService(repo, hasher, nil)
	return svc, repo, hasher
}

func TestAudit_Record_Nominal(t *testing.T) {
	svc, repo, _ := setupAuditTest()
	ctx := context.Background()

	actorID := "usr_100"
	actorRole := "faculty"
	req := RecordAuditRequest{
		TenantID:     "ten_alpha",
		ActorID:      &actorID,
		ActorType:    ActorTypeUser,
		ActorRole:    &actorRole,
		Action:       "course:syllabus:update",
		ResourceType: "syllabus",
		ResourceID:   ptr("syl_200"),
		Status:       StatusSuccess,
		StatusCode:   ptr(200),
		IPAddress:    ptr("192.168.1.50"),
		UserAgent:    ptr("Mozilla/5.0"),
		TraceID:      ptr("req_xyz123"),
		Metadata:     json.RawMessage(`{"module_count": 5}`),
	}

	record, err := svc.Record(ctx, req)
	if err != nil {
		t.Fatalf("Failed to record audit entry: %v", err)
	}

	if record.ID == "" {
		t.Error("Expected generated ID")
	}
	if record.Hash == "" {
		t.Error("Expected computed SHA-256 hash")
	}
	if record.PrevHash == nil || *record.PrevHash != GenesisHash {
		t.Errorf("Expected genesis prev_hash, got %v", record.PrevHash)
	}

	// Verify persistence in repository
	fetched, err := repo.GetByID(ctx, "ten_alpha", record.ID)
	if err != nil {
		t.Fatalf("Failed to fetch recorded log: %v", err)
	}
	if fetched.Hash != record.Hash {
		t.Errorf("Persisted hash mismatch: expected %s, got %s", record.Hash, fetched.Hash)
	}
}

func TestAudit_SequentialHashChaining(t *testing.T) {
	svc, _, _ := setupAuditTest()
	ctx := context.Background()
	tenantID := "ten_beta"

	// Record 3 sequential actions
	r1, err := svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "auth:login:success",
		ResourceType: "user",
		Status:       StatusSuccess,
	})
	if err != nil {
		t.Fatalf("Failed r1: %v", err)
	}

	r2, err := svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "hostel:gatepass:request",
		ResourceType: "gatepass",
		Status:       StatusSuccess,
	})
	if err != nil {
		t.Fatalf("Failed r2: %v", err)
	}

	r3, err := svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "hostel:gatepass:approve",
		ResourceType: "gatepass",
		Status:       StatusSuccess,
	})
	if err != nil {
		t.Fatalf("Failed r3: %v", err)
	}

	// Verify chaining links
	if *r1.PrevHash != GenesisHash {
		t.Errorf("r1 prev_hash should be GenesisHash, got %s", *r1.PrevHash)
	}
	if *r2.PrevHash != r1.Hash {
		t.Errorf("r2 prev_hash should match r1 hash, expected %s, got %s", r1.Hash, *r2.PrevHash)
	}
	if *r3.PrevHash != r2.Hash {
		t.Errorf("r3 prev_hash should match r2 hash, expected %s, got %s", r2.Hash, *r3.PrevHash)
	}

	// Verify chain integrity
	result, err := svc.VerifyChainIntegrity(ctx, tenantID, nil, nil)
	if err != nil {
		t.Fatalf("VerifyChainIntegrity error: %v", err)
	}
	if !result.IsChainIntact {
		t.Errorf("Expected chain to be intact, got errors: %+v", result.Errors)
	}
	if result.TotalVerified != 3 {
		t.Errorf("Expected 3 verified records, got %d", result.TotalVerified)
	}
	if result.HeadHash != r3.Hash {
		t.Errorf("Expected head hash %s, got %s", r3.Hash, result.HeadHash)
	}
}

func TestAudit_TamperDetection_PayloadModified(t *testing.T) {
	svc, repo, _ := setupAuditTest()
	ctx := context.Background()
	tenantID := "ten_security"

	r1, _ := svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "billing:fee:paid",
		ResourceType: "invoice",
		Status:       StatusSuccess,
	})
	r2, _ := svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "billing:fee:refund",
		ResourceType: "invoice",
		Status:       StatusSuccess,
	})
	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "billing:report:exported",
		ResourceType: "report",
		Status:       StatusSuccess,
	})

	// Malicious attacker alters r2 action directly in database
	repo.TamperWithRecord(r2.ID, func(log *AuditLog) {
		log.Action = "billing:fee:FRAUDULENT_MUTATION"
	})

	result, err := svc.VerifyChainIntegrity(ctx, tenantID, nil, nil)
	if err != nil {
		t.Fatalf("VerifyChainIntegrity error: %v", err)
	}

	if result.IsChainIntact {
		t.Fatal("Expected chain verification to FAIL due to payload tampering")
	}
	if len(result.Errors) == 0 {
		t.Fatal("Expected verification errors to be populated")
	}

	foundR2Error := false
	for _, e := range result.Errors {
		if e.RecordID == r2.ID {
			foundR2Error = true
			break
		}
	}
	if !foundR2Error {
		t.Errorf("Expected error to flag tampered record ID %s, but got: %+v", r2.ID, result.Errors)
	}
	_ = r1
}

func TestAudit_TamperDetection_BrokenLink(t *testing.T) {
	svc, repo, _ := setupAuditTest()
	ctx := context.Background()
	tenantID := "ten_integrity"

	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "event:create",
		ResourceType: "event",
		Status:       StatusSuccess,
	})
	r2, _ := svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "event:publish",
		ResourceType: "event",
		Status:       StatusSuccess,
	})

	// Attacker modifies prev_hash pointer of r2
	repo.TamperWithRecord(r2.ID, func(log *AuditLog) {
		fakeHash := "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
		log.PrevHash = &fakeHash
	})

	result, err := svc.VerifyChainIntegrity(ctx, tenantID, nil, nil)
	if err != nil {
		t.Fatalf("VerifyChainIntegrity error: %v", err)
	}

	if result.IsChainIntact {
		t.Fatal("Expected chain verification to FAIL due to broken hash pointer")
	}
}

func TestAudit_AsyncSubscriber_BatchingAndDrain(t *testing.T) {
	repo := NewMockRepository()
	hasher := NewSHA256Hasher()

	subscriber := NewAsyncSubscriber(repo, hasher, BatchConfig{
		QueueCapacity: 500,
		BatchSize:     5,
		FlushInterval: 20 * time.Millisecond,
	})

	ctx := context.Background()
	subscriber.Start(ctx)

	tenantID := "ten_async"
	for i := 0; i < 15; i++ {
		err := subscriber.Enqueue(RecordAuditRequest{
			TenantID:     tenantID,
			Action:       fmt.Sprintf("mess:meal:scanned_%d", i),
			ResourceType: "meal_token",
			Status:       StatusSuccess,
		})
		if err != nil {
			t.Fatalf("Enqueue failed at %d: %v", i, err)
		}
	}

	// Stop cleanly drains remaining items
	subscriber.Stop()

	// Verify all 15 records were processed and persisted
	records, err := repo.GetRecordsInRange(ctx, tenantID, nil, nil)
	if err != nil {
		t.Fatalf("GetRecordsInRange failed: %v", err)
	}
	if len(records) != 15 {
		t.Fatalf("Expected 15 records after flush, got %d", len(records))
	}

	// Verify chain integrity across all 15 batched records
	verification := hasher.VerifyChain(records, GenesisHash)
	if !verification.IsChainIntact {
		t.Errorf("Batched records failed chain integrity: %+v", verification.Errors)
	}
}

func TestAudit_QueryLogs_FiltersAndPagination(t *testing.T) {
	svc, _, _ := setupAuditTest()
	ctx := context.Background()
	tenantID := "ten_query"

	now := time.Now().UTC()
	t1 := now.Add(-3 * time.Hour)
	t2 := now.Add(-2 * time.Hour)
	t3 := now.Add(-1 * time.Hour)

	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      ptr("usr_prof_1"),
		ActorType:    ActorTypeUser,
		Action:       "attendance:mark:lecture",
		ResourceType: "attendance",
		Status:       StatusSuccess,
		Timestamp:    &t1,
	})
	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      ptr("usr_prof_1"),
		ActorType:    ActorTypeUser,
		Action:       "attendance:mark:lab",
		ResourceType: "attendance",
		Status:       StatusSuccess,
		Timestamp:    &t2,
	})
	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      ptr("usr_warden_2"),
		ActorType:    ActorTypeUser,
		Action:       "hostel:gatepass:reject",
		ResourceType: "gatepass",
		Status:       StatusFailure,
		Timestamp:    &t3,
	})

	// Filter by ActorID
	logs, total, err := svc.QueryLogs(ctx, AuditFilter{
		TenantID: tenantID,
		ActorID:  ptr("usr_prof_1"),
	})
	if err != nil {
		t.Fatalf("QueryLogs error: %v", err)
	}
	if total != 2 || len(logs) != 2 {
		t.Errorf("Expected 2 logs for usr_prof_1, got total=%d, count=%d", total, len(logs))
	}

	// Filter by ResourceType
	logsRes, totalRes, _ := svc.QueryLogs(ctx, AuditFilter{
		TenantID:     tenantID,
		ResourceType: ptr("gatepass"),
	})
	if totalRes != 1 || len(logsRes) != 1 {
		t.Errorf("Expected 1 gatepass log, got total=%d", totalRes)
	}

	// Filter by Action substring
	actionQuery := "attendance"
	logsAct, totalAct, _ := svc.QueryLogs(ctx, AuditFilter{
		TenantID: tenantID,
		Action:   &actionQuery,
	})
	if totalAct != 2 || len(logsAct) != 2 {
		t.Errorf("Expected 2 attendance logs, got total=%d", totalAct)
	}
}

func TestAudit_ComplianceReport(t *testing.T) {
	svc, _, _ := setupAuditTest()
	ctx := context.Background()
	tenantID := "ten_accreditation"

	now := time.Now().UTC()
	from := now.Add(-10 * time.Minute)
	to := now.Add(10 * time.Minute)

	// Create a mix of events
	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      ptr("usr_a"),
		ActorRole:    ptr("student"),
		Action:       "library:book:issue",
		ResourceType: "book",
		Status:       StatusSuccess,
	})
	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      ptr("usr_b"),
		ActorRole:    ptr("faculty"),
		Action:       "studyhub:assignment:grade",
		ResourceType: "assignment",
		Status:       StatusSuccess,
	})
	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      ptr("usr_c"),
		ActorRole:    ptr("admin"),
		Action:       "auth:token:breach_detected",
		ResourceType: "session",
		Status:       StatusFailure,
	})

	report, err := svc.GenerateComplianceReport(ctx, ComplianceReportRequest{
		TenantID:  tenantID,
		Framework: "NAAC_CRITERIA_6",
		FromTime:  from,
		ToTime:    to,
	})
	if err != nil {
		t.Fatalf("GenerateComplianceReport failed: %v", err)
	}

	if report.TotalEvents != 3 {
		t.Errorf("Expected 3 total events, got %d", report.TotalEvents)
	}
	if report.SuccessfulEvents != 2 {
		t.Errorf("Expected 2 successful events, got %d", report.SuccessfulEvents)
	}
	if report.SecurityIncidents != 1 {
		t.Errorf("Expected 1 security incident, got %d", report.SecurityIncidents)
	}
	if !report.ChainIntegrity {
		t.Error("Expected intact chain integrity in compliance report")
	}
}

func TestAudit_VerificationCheckpoint(t *testing.T) {
	svc, repo, _ := setupAuditTest()
	ctx := context.Background()
	tenantID := "ten_checkpoint"

	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "onboarding:student:registered",
		ResourceType: "student",
		Status:       StatusSuccess,
	})
	_, _ = svc.Record(ctx, RecordAuditRequest{
		TenantID:     tenantID,
		Action:       "onboarding:kyc:verified",
		ResourceType: "kyc_document",
		Status:       StatusSuccess,
	})

	cp, err := svc.CreateVerificationCheckpoint(ctx, tenantID)
	if err != nil {
		t.Fatalf("CreateVerificationCheckpoint failed: %v", err)
	}

	if cp.RecordCount != 2 {
		t.Errorf("Expected record count 2, got %d", cp.RecordCount)
	}
	if !cp.IsVerified {
		t.Error("Expected checkpoint to be verified")
	}

	checkpoints, err := repo.ListCheckpoints(ctx, tenantID, 10, 0)
	if err != nil {
		t.Fatalf("ListCheckpoints failed: %v", err)
	}
	if len(checkpoints) != 1 {
		t.Errorf("Expected 1 stored checkpoint, got %d", len(checkpoints))
	}
}

func TestAudit_ValidationErrors(t *testing.T) {
	svc, _, _ := setupAuditTest()
	ctx := context.Background()

	// Missing tenant
	_, err := svc.Record(ctx, RecordAuditRequest{
		Action:       "test",
		ResourceType: "test",
	})
	if err != ErrTenantRequired {
		t.Errorf("Expected ErrTenantRequired, got %v", err)
	}

	// Missing action
	_, err = svc.Record(ctx, RecordAuditRequest{
		TenantID:     "ten_1",
		ResourceType: "test",
	})
	if err != ErrActionRequired {
		t.Errorf("Expected ErrActionRequired, got %v", err)
	}

	// Missing resource type
	_, err = svc.Record(ctx, RecordAuditRequest{
		TenantID: "ten_1",
		Action:   "test",
	})
	if err != ErrResourceTypeRequired {
		t.Errorf("Expected ErrResourceTypeRequired, got %v", err)
	}
}

func TestAudit_MultiTenantIsolation(t *testing.T) {
	svc, _, _ := setupAuditTest()
	ctx := context.Background()

	tA := "tenant_engineering_college"
	tB := "tenant_medical_college"

	rA1, _ := svc.Record(ctx, RecordAuditRequest{TenantID: tA, Action: "eng:1", ResourceType: "r"})
	rB1, _ := svc.Record(ctx, RecordAuditRequest{TenantID: tB, Action: "med:1", ResourceType: "r"})
	rA2, _ := svc.Record(ctx, RecordAuditRequest{TenantID: tA, Action: "eng:2", ResourceType: "r"})
	rB2, _ := svc.Record(ctx, RecordAuditRequest{TenantID: tB, Action: "med:2", ResourceType: "r"})

	// Verify Tenant A chain links only within Tenant A
	if *rA1.PrevHash != GenesisHash {
		t.Errorf("rA1 should have genesis hash")
	}
	if *rA2.PrevHash != rA1.Hash {
		t.Errorf("rA2 should link to rA1 hash")
	}

	// Verify Tenant B chain links only within Tenant B
	if *rB1.PrevHash != GenesisHash {
		t.Errorf("rB1 should have genesis hash")
	}
	if *rB2.PrevHash != rB1.Hash {
		t.Errorf("rB2 should link to rB1 hash")
	}

	// Verify both chains pass verification independently
	resA, _ := svc.VerifyChainIntegrity(ctx, tA, nil, nil)
	resB, _ := svc.VerifyChainIntegrity(ctx, tB, nil, nil)

	if !resA.IsChainIntact || resA.TotalVerified != 2 {
		t.Errorf("Tenant A chain integrity failure: %+v", resA)
	}
	if !resB.IsChainIntact || resB.TotalVerified != 2 {
		t.Errorf("Tenant B chain integrity failure: %+v", resB)
	}
}
