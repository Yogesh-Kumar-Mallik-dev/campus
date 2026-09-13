/**
 * BLOCK_HOSTEL_REPOSITORY_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   Repository interface contracts for hostel data storage and retrieval.
 */

package hostel

import (
	"context"
	"time"
)

type Repository interface {
	// Block operations
	CreateBlock(ctx context.Context, block *Block) error
	GetBlockByID(ctx context.Context, tenantID, id string) (*Block, error)
	GetBlockByCode(ctx context.Context, tenantID, code string) (*Block, error)
	ListBlocks(ctx context.Context, tenantID string) ([]*Block, error)

	// Room operations
	CreateRoom(ctx context.Context, room *Room) error
	GetRoomByID(ctx context.Context, tenantID, id string) (*Room, error)
	ListRoomsByBlock(ctx context.Context, tenantID, blockID string) ([]*Room, error)
	UpdateRoomStatus(ctx context.Context, tenantID, id string, status RoomStatus) error

	// Bed operations
	CreateBed(ctx context.Context, bed *Bed) error
	GetBedByID(ctx context.Context, tenantID, id string) (*Bed, error)
	UpdateBedStatus(ctx context.Context, tenantID, id string, status BedStatus) error
	ListBedsByRoom(ctx context.Context, tenantID, roomID string) ([]*Bed, error)

	// Allocation operations
	CreateAllocation(ctx context.Context, alloc *Allocation) error
	GetAllocationByID(ctx context.Context, tenantID, id string) (*Allocation, error)
	GetActiveAllocationByStudent(ctx context.Context, tenantID, studentID string) (*Allocation, error)
	UpdateAllocationStatus(ctx context.Context, tenantID, id string, status AllocationStatus, vacatedAt *time.Time) error
	ListAllocations(ctx context.Context, tenantID, studentID, bedID string) ([]*Allocation, error)

	// Gate Pass operations
	CreateGatePass(ctx context.Context, gp *GatePass) error
	GetGatePassByID(ctx context.Context, tenantID, id string) (*GatePass, error)
	UpdateGatePass(ctx context.Context, gp *GatePass) error
	ListGatePasses(ctx context.Context, tenantID, studentID, blockID string, status *GatePassStatus) ([]*GatePass, error)

	// Disciplinary / Incident Log operations
	CreateIncident(ctx context.Context, inc *IncidentLog) error
	GetIncidentByID(ctx context.Context, tenantID, id string) (*IncidentLog, error)
	ListIncidents(ctx context.Context, tenantID, studentID, blockID string) ([]*IncidentLog, error)
}
