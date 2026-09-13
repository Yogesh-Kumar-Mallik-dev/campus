/**
 * BLOCK_HOSTEL_UNIT_TESTS_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   Unit tests covering block hierarchies, bed allocation invariants, gate pass state machines, and curfew breach incident cascades.
 */

package hostel_test

import (
	"context"
	"testing"
	"time"

	"campus/backend/hostel"
)

func TestHostel_BlockRoomBedCreationAndListing(t *testing.T) {
	repo := hostel.NewMockRepository()
	svc := hostel.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	// 1. Create Block
	block, err := svc.CreateBlock(ctx, tenantID, "Aryabhatta Block A", "BH-A", hostel.GenderMale, 4, 80, 160, nil)
	if err != nil {
		t.Fatalf("CreateBlock failed: %v", err)
	}
	if block.Code != "BH-A" {
		t.Errorf("Expected block code BH-A, got %s", block.Code)
	}

	// 2. Room Capacity Validation
	_, err = svc.CreateRoom(ctx, tenantID, block.ID, "101", 1, hostel.RoomTypeDouble, true, 30000, 3)
	if err == nil {
		t.Fatalf("Expected error creating Double room with 3 beds, got nil")
	}

	// 3. Create Valid Room
	room, err := svc.CreateRoom(ctx, tenantID, block.ID, "101", 1, hostel.RoomTypeDouble, true, 30000, 2)
	if err != nil {
		t.Fatalf("CreateRoom failed: %v", err)
	}

	// 4. Create Beds
	bed1, err := svc.CreateBed(ctx, tenantID, room.ID, "A-101-1")
	if err != nil {
		t.Fatalf("CreateBed 1 failed: %v", err)
	}
	bed2, err := svc.CreateBed(ctx, tenantID, room.ID, "A-101-2")
	if err != nil {
		t.Fatalf("CreateBed 2 failed: %v", err)
	}

	// 5. List Beds
	beds, err := svc.ListBeds(ctx, tenantID, room.ID)
	if err != nil {
		t.Fatalf("ListBeds failed: %v", err)
	}
	if len(beds) != 2 {
		t.Errorf("Expected 2 beds, got %d", len(beds))
	}
	_ = bed1
	_ = bed2
}

func TestHostel_BedAllocationAndVacateLifecycle(t *testing.T) {
	repo := hostel.NewMockRepository()
	svc := hostel.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	block, _ := svc.CreateBlock(ctx, tenantID, "Gargi Girls Hostel", "GH-1", hostel.GenderFemale, 3, 60, 120, nil)
	room, _ := svc.CreateRoom(ctx, tenantID, block.ID, "201", 2, hostel.RoomTypeSingle, false, 35000, 1)
	bed, _ := svc.CreateBed(ctx, tenantID, room.ID, "G-201-1")

	studentID := "stu_anjali_01"

	// 1. Allocate bed
	alloc, err := svc.AllocateBed(ctx, tenantID, bed.ID, studentID, "2026-2027", 1, "Semester 1 allotment")
	if err != nil {
		t.Fatalf("AllocateBed failed: %v", err)
	}
	if alloc.Status != hostel.AllocationStatusAllocated {
		t.Errorf("Expected status ALLOCATED, got %s", alloc.Status)
	}

	// 2. Invariant: Same student cannot be allocated another bed
	bed2, _ := svc.CreateBed(ctx, tenantID, room.ID, "G-201-2-dummy")
	_, err = svc.AllocateBed(ctx, tenantID, bed2.ID, studentID, "2026-2027", 1, "Second allotment attempt")
	if err == nil {
		t.Fatalf("Expected ErrStudentAlreadyAllocated, got nil")
	}

	// 3. Invariant: Same bed cannot be allocated to another student
	_, err = svc.AllocateBed(ctx, tenantID, bed.ID, "stu_pooja_02", "2026-2027", 1, "Clash allotment")
	if err == nil {
		t.Fatalf("Expected ErrBedUnavailable, got nil")
	}

	// 4. Vacate Bed
	vacatedAlloc, err := svc.VacateBed(ctx, tenantID, alloc.ID, "Semester ended")
	if err != nil {
		t.Fatalf("VacateBed failed: %v", err)
	}
	if vacatedAlloc.Status != hostel.AllocationStatusVacated {
		t.Errorf("Expected status VACATED, got %s", vacatedAlloc.Status)
	}
	if vacatedAlloc.VacatedAt == nil {
		t.Errorf("Expected VacatedAt to be set")
	}

	// 5. Now bed can be re-allocated to another student
	alloc2, err := svc.AllocateBed(ctx, tenantID, bed.ID, "stu_pooja_02", "2026-2027", 2, "Spring semester")
	if err != nil {
		t.Fatalf("Re-allocating vacated bed failed: %v", err)
	}
	if alloc2.StudentID != "stu_pooja_02" {
		t.Errorf("Expected student stu_pooja_02, got %s", alloc2.StudentID)
	}
}

func TestHostel_GatePassLifecycle_ApprovedAndReturned(t *testing.T) {
	repo := hostel.NewMockRepository()
	svc := hostel.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"
	studentID := "stu_rahul_01"
	blockID := "blk_01"
	wardenID := "warden_sharma"

	now := time.Now().UTC()
	outTime := now.Add(2 * time.Hour)
	inTime := now.Add(6 * time.Hour)

	// 1. Apply Gate Pass
	gp, err := svc.ApplyGatePass(ctx, tenantID, studentID, blockID, "Home visit for festival", "New Delhi", "+919876543210", outTime, inTime)
	if err != nil {
		t.Fatalf("ApplyGatePass failed: %v", err)
	}
	if gp.Status != hostel.GatePassStatusPending {
		t.Errorf("Expected status PENDING, got %s", gp.Status)
	}

	// 2. Warden Approve
	gpApproved, err := svc.ReviewGatePass(ctx, tenantID, gp.ID, wardenID, true, "")
	if err != nil {
		t.Fatalf("ReviewGatePass Approve failed: %v", err)
	}
	if gpApproved.Status != hostel.GatePassStatusApproved {
		t.Errorf("Expected status APPROVED, got %s", gpApproved.Status)
	}

	// 3. Security Guard records Exit
	actualOut := outTime.Add(10 * time.Minute)
	gpOut, err := svc.RecordGatePassExit(ctx, tenantID, gp.ID, actualOut)
	if err != nil {
		t.Fatalf("RecordGatePassExit failed: %v", err)
	}
	if gpOut.Status != hostel.GatePassStatusOutCampus {
		t.Errorf("Expected status OUT_CAMPUS, got %s", gpOut.Status)
	}

	// 4. Security Guard records Return on time (no curfew breach)
	actualIn := inTime.Add(-15 * time.Minute)
	gpReturned, incident, err := svc.RecordGatePassReturn(ctx, tenantID, gp.ID, wardenID, actualIn)
	if err != nil {
		t.Fatalf("RecordGatePassReturn failed: %v", err)
	}
	if gpReturned.Status != hostel.GatePassStatusReturned {
		t.Errorf("Expected status RETURNED, got %s", gpReturned.Status)
	}
	if incident != nil {
		t.Errorf("Expected no curfew incident for on-time return, got %+v", incident)
	}
}

func TestHostel_GatePassCurfewBreach_AutomaticIncidentLogging(t *testing.T) {
	repo := hostel.NewMockRepository()
	svc := hostel.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"
	studentID := "stu_karan_02"
	blockID := "blk_01"
	wardenID := "warden_sharma"

	now := time.Now().UTC()
	outTime := now.Add(1 * time.Hour)
	inTime := now.Add(4 * time.Hour)

	gp, err := svc.ApplyGatePass(ctx, tenantID, studentID, blockID, "City Library visit", "Central Library", "+919876543211", outTime, inTime)
	if err != nil {
		t.Fatalf("ApplyGatePass failed: %v", err)
	}

	_, _ = svc.ReviewGatePass(ctx, tenantID, gp.ID, wardenID, true, "")
	_, _ = svc.RecordGatePassExit(ctx, tenantID, gp.ID, outTime)

	// Return 90 minutes late!
	actualIn := inTime.Add(90 * time.Minute)
	gpReturned, incident, err := svc.RecordGatePassReturn(ctx, tenantID, gp.ID, wardenID, actualIn)
	if err != nil {
		t.Fatalf("RecordGatePassReturn failed: %v", err)
	}
	if gpReturned.Status != hostel.GatePassStatusReturned {
		t.Errorf("Expected status RETURNED, got %s", gpReturned.Status)
	}
	if incident == nil {
		t.Fatalf("Expected automatic curfew breach incident, got nil")
	}
	if incident.IncidentType != hostel.IncidentCurfewViolation {
		t.Errorf("Expected incident type CURFEW_VIOLATION, got %s", incident.IncidentType)
	}
	if incident.FineAmount <= 0 {
		t.Errorf("Expected fine amount to be imposed, got %f", incident.FineAmount)
	}

	// Verify incident in repository
	incidents, err := svc.ListIncidents(ctx, tenantID, studentID, "")
	if err != nil {
		t.Fatalf("ListIncidents failed: %v", err)
	}
	if len(incidents) != 1 {
		t.Errorf("Expected 1 incident record, got %d", len(incidents))
	}
}

func TestHostel_GatePassInvalidTransitionsAndRejections(t *testing.T) {
	repo := hostel.NewMockRepository()
	svc := hostel.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	now := time.Now().UTC()
	// Invalid time range: in before out
	_, err := svc.ApplyGatePass(ctx, tenantID, "stu_01", "blk_01", "Reason", "Dest", "123", now.Add(2*time.Hour), now.Add(1*time.Hour))
	if err == nil {
		t.Fatalf("Expected ErrInvalidTimeRange, got nil")
	}

	// Rejection flow
	gp, _ := svc.ApplyGatePass(ctx, tenantID, "stu_01", "blk_01", "Reason", "Dest", "123", now.Add(1*time.Hour), now.Add(3*time.Hour))
	rejectedGP, err := svc.ReviewGatePass(ctx, tenantID, gp.ID, "warden_1", false, "Disciplinary hold")
	if err != nil {
		t.Fatalf("ReviewGatePass reject failed: %v", err)
	}
	if rejectedGP.Status != hostel.GatePassStatusRejected {
		t.Errorf("Expected status REJECTED, got %s", rejectedGP.Status)
	}

	// Attempting exit on rejected pass should fail
	_, err = svc.RecordGatePassExit(ctx, tenantID, gp.ID, now.Add(1*time.Hour))
	if err == nil {
		t.Fatalf("Expected ErrInvalidGatePassState on rejected pass exit, got nil")
	}
}

func TestHostel_DisciplinaryIncidentLogging(t *testing.T) {
	repo := hostel.NewMockRepository()
	svc := hostel.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	inc, err := svc.LogIncident(
		ctx,
		tenantID,
		"stu_vikas_03",
		"blk_01",
		"warden_01",
		hostel.IncidentNoiseDisturbance,
		hostel.SeverityMedium,
		"Loud speakers past quiet hours",
		"Playing music at 1:00 AM after multiple warnings",
		"Verbal warning and room inspection",
		250.0,
	)
	if err != nil {
		t.Fatalf("LogIncident failed: %v", err)
	}
	if inc.FineAmount != 250.0 {
		t.Errorf("Expected fine 250.0, got %f", inc.FineAmount)
	}
}
