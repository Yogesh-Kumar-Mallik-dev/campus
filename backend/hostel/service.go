/**
 * BLOCK_HOSTEL_SERVICE_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   Business domain service implementing bed allocations, gate passes, and curfew enforcement.
 */

package hostel

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"campus/backend/audit"
)

type Service struct {
	repo     Repository
	auditSub audit.Subscriber
}

func NewService(repo Repository, auditSub audit.Subscriber) *Service {
	return &Service{
		repo:     repo,
		auditSub: auditSub,
	}
}

func (s *Service) emitAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
	if s.auditSub == nil {
		return
	}
	var metaRaw json.RawMessage
	if meta != nil {
		if b, err := json.Marshal(meta); err == nil {
			metaRaw = b
		}
	}
	_ = s.auditSub.Enqueue(audit.RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      &actorID,
		ActorType:    audit.ActorType(actorType),
		Action:       action,
		ResourceType: resType,
		ResourceID:   &resID,
		Status:       status,
		Metadata:     metaRaw,
	})
}

// CreateBlock registers a physical residential block
func (s *Service) CreateBlock(ctx context.Context, tenantID, name, code string, gender Gender, totalFloors, totalRooms, capacity int, wardenID *string) (*Block, error) {
	if tenantID == "" || name == "" || code == "" {
		return nil, NewDomainError(nil, 400, "Bad Request", "tenant_id, name and code are required", "https://campus.internal/errors/invalid-argument")
	}

	block := &Block{
		ID:          fmt.Sprintf("blk_%d", time.Now().UnixNano()),
		TenantID:    tenantID,
		Name:        name,
		Code:        code,
		Gender:      gender,
		TotalFloors: totalFloors,
		TotalRooms:  totalRooms,
		Capacity:    capacity,
		WardenID:    wardenID,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.CreateBlock(ctx, block); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, "SYSTEM", "SYSTEM", "hostel:block:created", "hostel_block", block.ID, audit.StatusSuccess, map[string]interface{}{
		"code": code,
		"name": name,
	})

	return block, nil
}

// ListBlocks retrieves all blocks for a tenant
func (s *Service) ListBlocks(ctx context.Context, tenantID string) ([]*Block, error) {
	return s.repo.ListBlocks(ctx, tenantID)
}

// CreateRoom configures a room within a block
func (s *Service) CreateRoom(ctx context.Context, tenantID, blockID, roomNumber string, floorNumber int, roomType RoomType, isAc bool, baseFee float64, maxBeds int) (*Room, error) {
	if err := ValidateRoomCapacity(roomType, maxBeds); err != nil {
		return nil, NewDomainError(err, 400, "Invalid Room Capacity", err.Error(), "https://campus.internal/errors/invalid-argument")
	}

	// Verify block exists
	if _, err := s.repo.GetBlockByID(ctx, tenantID, blockID); err != nil {
		return nil, err
	}

	room := &Room{
		ID:                 fmt.Sprintf("rm_%d", time.Now().UnixNano()),
		TenantID:           tenantID,
		BlockID:            blockID,
		RoomNumber:         roomNumber,
		FloorNumber:        floorNumber,
		RoomType:           roomType,
		IsAC:               isAc,
		BaseFeePerSemester: baseFee,
		Status:             RoomStatusAvailable,
		MaxBeds:            maxBeds,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	if err := s.repo.CreateRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

// ListRooms lists rooms in a block
func (s *Service) ListRooms(ctx context.Context, tenantID, blockID string) ([]*Room, error) {
	return s.repo.ListRoomsByBlock(ctx, tenantID, blockID)
}

// CreateBed creates a bed in a room
func (s *Service) CreateBed(ctx context.Context, tenantID, roomID, bedNumber string) (*Bed, error) {
	if _, err := s.repo.GetRoomByID(ctx, tenantID, roomID); err != nil {
		return nil, err
	}

	bed := &Bed{
		ID:        fmt.Sprintf("bed_%d", time.Now().UnixNano()),
		TenantID:  tenantID,
		RoomID:    roomID,
		BedNumber: bedNumber,
		Status:    BedStatusAvailable,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateBed(ctx, bed); err != nil {
		return nil, err
	}

	return bed, nil
}

// ListBeds lists beds in a room
func (s *Service) ListBeds(ctx context.Context, tenantID, roomID string) ([]*Bed, error) {
	return s.repo.ListBedsByRoom(ctx, tenantID, roomID)
}

// AllocateBed assigns a bed to a student
func (s *Service) AllocateBed(ctx context.Context, tenantID, bedID, studentID, academicYear string, semester int, remarks string) (*Allocation, error) {
	bed, err := s.repo.GetBedByID(ctx, tenantID, bedID)
	if err != nil {
		return nil, err
	}

	if bed.Status != BedStatusAvailable {
		return nil, ErrBedUnavailable
	}

	// Verify student does not already have an active allocation
	activeAlloc, err := s.repo.GetActiveAllocationByStudent(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if activeAlloc != nil {
		return nil, ErrStudentAlreadyAllocated
	}

	alloc := &Allocation{
		ID:           fmt.Sprintf("alloc_%d", time.Now().UnixNano()),
		TenantID:     tenantID,
		BedID:        bedID,
		StudentID:    studentID,
		AcademicYear: academicYear,
		Semester:     semester,
		AllocatedAt:  time.Now().UTC(),
		Status:       AllocationStatusAllocated,
		Remarks:      remarks,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := s.repo.CreateAllocation(ctx, alloc); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateBedStatus(ctx, tenantID, bedID, BedStatusAllocated); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, "SYSTEM", "SYSTEM", "hostel:bed:allocated", "hostel_bed_allocation", alloc.ID, audit.StatusSuccess, map[string]interface{}{
		"student_id":    studentID,
		"bed_id":        bedID,
		"academic_year": academicYear,
	})

	return alloc, nil
}

// VacateBed marks an allocation as vacated and frees the bed
func (s *Service) VacateBed(ctx context.Context, tenantID, allocationID string, remarks string) (*Allocation, error) {
	alloc, err := s.repo.GetAllocationByID(ctx, tenantID, allocationID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := s.repo.UpdateAllocationStatus(ctx, tenantID, allocationID, AllocationStatusVacated, &now); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateBedStatus(ctx, tenantID, alloc.BedID, BedStatusAvailable); err != nil {
		return nil, err
	}

	alloc.Status = AllocationStatusVacated
	alloc.VacatedAt = &now
	alloc.Remarks = remarks

	s.emitAudit(tenantID, "SYSTEM", "SYSTEM", "hostel:bed:vacated", "hostel_bed_allocation", alloc.ID, audit.StatusSuccess, map[string]interface{}{
		"student_id": alloc.StudentID,
		"bed_id":     alloc.BedID,
	})

	return alloc, nil
}

// ApplyGatePass creates a new out-pass request
func (s *Service) ApplyGatePass(ctx context.Context, tenantID, studentID, blockID, reason, destination, emergencyContact string, outTime, inTime time.Time) (*GatePass, error) {
	if err := ValidateGatePassTimes(outTime, inTime); err != nil {
		return nil, err
	}

	gp := &GatePass{
		ID:               fmt.Sprintf("gp_%d", time.Now().UnixNano()),
		TenantID:         tenantID,
		StudentID:        studentID,
		BlockID:          blockID,
		Reason:           reason,
		Destination:      destination,
		EmergencyContact: emergencyContact,
		ExpectedOutAt:    outTime,
		ExpectedInAt:     inTime,
		Status:           GatePassStatusPending,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	if err := s.repo.CreateGatePass(ctx, gp); err != nil {
		return nil, err
	}

	return gp, nil
}

// ReviewGatePass approves or rejects a gate pass
func (s *Service) ReviewGatePass(ctx context.Context, tenantID, passID, wardenID string, approve bool, rejectionReason string) (*GatePass, error) {
	gp, err := s.repo.GetGatePassByID(ctx, tenantID, passID)
	if err != nil {
		return nil, err
	}

	targetStatus := GatePassStatusApproved
	if !approve {
		targetStatus = GatePassStatusRejected
	}

	if !CanTransitionGatePass(gp.Status, targetStatus) {
		return nil, ErrInvalidGatePassState
	}

	gp.Status = targetStatus
	gp.ApprovedByID = &wardenID
	if !approve && rejectionReason != "" {
		gp.RejectionReason = &rejectionReason
	}
	gp.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateGatePass(ctx, gp); err != nil {
		return nil, err
	}

	return gp, nil
}

// RecordGatePassExit marks student departed
func (s *Service) RecordGatePassExit(ctx context.Context, tenantID, passID string, actualOut time.Time) (*GatePass, error) {
	gp, err := s.repo.GetGatePassByID(ctx, tenantID, passID)
	if err != nil {
		return nil, err
	}

	if !CanTransitionGatePass(gp.Status, GatePassStatusOutCampus) {
		return nil, ErrInvalidGatePassState
	}

	gp.Status = GatePassStatusOutCampus
	gp.ActualOutAt = &actualOut
	gp.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateGatePass(ctx, gp); err != nil {
		return nil, err
	}

	return gp, nil
}

// RecordGatePassReturn marks student returned and checks for curfew violations
func (s *Service) RecordGatePassReturn(ctx context.Context, tenantID, passID, wardenID string, actualIn time.Time) (*GatePass, *IncidentLog, error) {
	gp, err := s.repo.GetGatePassByID(ctx, tenantID, passID)
	if err != nil {
		return nil, nil, err
	}

	if !CanTransitionGatePass(gp.Status, GatePassStatusReturned) {
		return nil, nil, ErrInvalidGatePassState
	}

	gp.Status = GatePassStatusReturned
	gp.ActualInAt = &actualIn
	gp.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateGatePass(ctx, gp); err != nil {
		return nil, nil, err
	}

	var incident *IncidentLog
	isBreach, overdueMins := CheckCurfewBreach(gp.ExpectedInAt, actualIn)
	if isBreach {
		inc, err := s.LogIncident(
			ctx,
			tenantID,
			gp.StudentID,
			gp.BlockID,
			wardenID,
			IncidentCurfewViolation,
			SeverityMedium,
			"Curfew Breach on Gate Pass Return",
			fmt.Sprintf("Student returned %d minutes late past expected time (%s)", overdueMins, gp.ExpectedInAt.Format(time.RFC3339)),
			"Automatic curfew infraction recorded; warning issued",
			100.0,
		)
		if err == nil {
			incident = inc
		}
	}

	return gp, incident, nil
}

// LogIncident creates a disciplinary or incident report
func (s *Service) LogIncident(ctx context.Context, tenantID, studentID, blockID, wardenID string, incType IncidentType, severity IncidentSeverity, title, description, actionTaken string, fineAmount float64) (*IncidentLog, error) {
	inc := &IncidentLog{
		ID:           fmt.Sprintf("inc_%d", time.Now().UnixNano()),
		TenantID:     tenantID,
		StudentID:    studentID,
		BlockID:      blockID,
		WardenID:     wardenID,
		IncidentType: incType,
		Severity:     severity,
		Title:        title,
		Description:  description,
		ActionTaken:  actionTaken,
		FineAmount:   fineAmount,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := s.repo.CreateIncident(ctx, inc); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, wardenID, "USER", "hostel:incident:logged", "hostel_incident_log", inc.ID, audit.StatusSuccess, map[string]interface{}{
		"student_id":    studentID,
		"incident_type": string(incType),
		"fine_amount":   fineAmount,
	})

	return inc, nil
}

// ListGatePasses lists gate passes matching filter
func (s *Service) ListGatePasses(ctx context.Context, tenantID, studentID, blockID string, status *GatePassStatus) ([]*GatePass, error) {
	return s.repo.ListGatePasses(ctx, tenantID, studentID, blockID, status)
}

// ListIncidents lists disciplinary incidents
func (s *Service) ListIncidents(ctx context.Context, tenantID, studentID, blockID string) ([]*IncidentLog, error) {
	return s.repo.ListIncidents(ctx, tenantID, studentID, blockID)
}
