/**
 * BLOCK_HOSTEL_MOCK_REPOSITORY_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   In-memory mock repository implementing Repository interface for unit testing.
 */

package hostel

import (
	"context"
	"sync"
	"time"
)

type MockRepository struct {
	mu          sync.RWMutex
	blocks      map[string]*Block
	rooms       map[string]*Room
	beds        map[string]*Bed
	allocations map[string]*Allocation
	gatePasses  map[string]*GatePass
	incidents   map[string]*IncidentLog
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		blocks:      make(map[string]*Block),
		rooms:       make(map[string]*Room),
		beds:        make(map[string]*Bed),
		allocations: make(map[string]*Allocation),
		gatePasses:  make(map[string]*GatePass),
		incidents:   make(map[string]*IncidentLog),
	}
}

func (m *MockRepository) CreateBlock(ctx context.Context, block *Block) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, b := range m.blocks {
		if b.TenantID == block.TenantID && b.Code == block.Code {
			return ErrBlockCodeExists
		}
	}
	m.blocks[block.ID] = block
	return nil
}

func (m *MockRepository) GetBlockByID(ctx context.Context, tenantID, id string) (*Block, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, exists := m.blocks[id]
	if !exists || b.TenantID != tenantID {
		return nil, ErrBlockNotFound
	}
	return b, nil
}

func (m *MockRepository) GetBlockByCode(ctx context.Context, tenantID, code string) (*Block, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, b := range m.blocks {
		if b.TenantID == tenantID && b.Code == code {
			return b, nil
		}
	}
	return nil, ErrBlockNotFound
}

func (m *MockRepository) ListBlocks(ctx context.Context, tenantID string) ([]*Block, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Block
	for _, b := range m.blocks {
		if b.TenantID == tenantID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *MockRepository) CreateRoom(ctx context.Context, room *Room) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, r := range m.rooms {
		if r.TenantID == room.TenantID && r.BlockID == room.BlockID && r.RoomNumber == room.RoomNumber {
			return ErrRoomNumberExists
		}
	}
	m.rooms[room.ID] = room
	return nil
}

func (m *MockRepository) GetRoomByID(ctx context.Context, tenantID, id string) (*Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, exists := m.rooms[id]
	if !exists || r.TenantID != tenantID {
		return nil, ErrRoomNotFound
	}
	return r, nil
}

func (m *MockRepository) ListRoomsByBlock(ctx context.Context, tenantID, blockID string) ([]*Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Room
	for _, r := range m.rooms {
		if r.TenantID == tenantID && r.BlockID == blockID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *MockRepository) UpdateRoomStatus(ctx context.Context, tenantID, id string, status RoomStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, exists := m.rooms[id]
	if !exists || r.TenantID != tenantID {
		return ErrRoomNotFound
	}
	r.Status = status
	r.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) CreateBed(ctx context.Context, bed *Bed) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, b := range m.beds {
		if b.TenantID == bed.TenantID && b.RoomID == bed.RoomID && b.BedNumber == bed.BedNumber {
			return ErrBedNumberExists
		}
	}
	m.beds[bed.ID] = bed
	return nil
}

func (m *MockRepository) GetBedByID(ctx context.Context, tenantID, id string) (*Bed, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, exists := m.beds[id]
	if !exists || b.TenantID != tenantID {
		return nil, ErrBedNotFound
	}
	return b, nil
}

func (m *MockRepository) UpdateBedStatus(ctx context.Context, tenantID, id string, status BedStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, exists := m.beds[id]
	if !exists || b.TenantID != tenantID {
		return ErrBedNotFound
	}
	b.Status = status
	b.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) ListBedsByRoom(ctx context.Context, tenantID, roomID string) ([]*Bed, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Bed
	for _, b := range m.beds {
		if b.TenantID == tenantID && b.RoomID == roomID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *MockRepository) CreateAllocation(ctx context.Context, alloc *Allocation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.allocations[alloc.ID] = alloc
	return nil
}

func (m *MockRepository) GetAllocationByID(ctx context.Context, tenantID, id string) (*Allocation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, exists := m.allocations[id]
	if !exists || a.TenantID != tenantID {
		return nil, ErrAllocationNotFound
	}
	return a, nil
}

func (m *MockRepository) GetActiveAllocationByStudent(ctx context.Context, tenantID, studentID string) (*Allocation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, a := range m.allocations {
		if a.TenantID == tenantID && a.StudentID == studentID && a.Status == AllocationStatusAllocated {
			return a, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) UpdateAllocationStatus(ctx context.Context, tenantID, id string, status AllocationStatus, vacatedAt *time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	a, exists := m.allocations[id]
	if !exists || a.TenantID != tenantID {
		return ErrAllocationNotFound
	}
	a.Status = status
	a.VacatedAt = vacatedAt
	a.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) ListAllocations(ctx context.Context, tenantID, studentID, bedID string) ([]*Allocation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Allocation
	for _, a := range m.allocations {
		if a.TenantID != tenantID {
			continue
		}
		if studentID != "" && a.StudentID != studentID {
			continue
		}
		if bedID != "" && a.BedID != bedID {
			continue
		}
		result = append(result, a)
	}
	return result, nil
}

func (m *MockRepository) CreateGatePass(ctx context.Context, gp *GatePass) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gatePasses[gp.ID] = gp
	return nil
}

func (m *MockRepository) GetGatePassByID(ctx context.Context, tenantID, id string) (*GatePass, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	gp, exists := m.gatePasses[id]
	if !exists || gp.TenantID != tenantID {
		return nil, ErrGatePassNotFound
	}
	return gp, nil
}

func (m *MockRepository) UpdateGatePass(ctx context.Context, gp *GatePass) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gatePasses[gp.ID] = gp
	return nil
}

func (m *MockRepository) ListGatePasses(ctx context.Context, tenantID, studentID, blockID string, status *GatePassStatus) ([]*GatePass, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*GatePass
	for _, gp := range m.gatePasses {
		if gp.TenantID != tenantID {
			continue
		}
		if studentID != "" && gp.StudentID != studentID {
			continue
		}
		if blockID != "" && gp.BlockID != blockID {
			continue
		}
		if status != nil && gp.Status != *status {
			continue
		}
		result = append(result, gp)
	}
	return result, nil
}

func (m *MockRepository) CreateIncident(ctx context.Context, inc *IncidentLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.incidents[inc.ID] = inc
	return nil
}

func (m *MockRepository) GetIncidentByID(ctx context.Context, tenantID, id string) (*IncidentLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inc, exists := m.incidents[id]
	if !exists || inc.TenantID != tenantID {
		return nil, ErrIncidentNotFound
	}
	return inc, nil
}

func (m *MockRepository) ListIncidents(ctx context.Context, tenantID, studentID, blockID string) ([]*IncidentLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*IncidentLog
	for _, inc := range m.incidents {
		if inc.TenantID != tenantID {
			continue
		}
		if studentID != "" && inc.StudentID != studentID {
			continue
		}
		if blockID != "" && inc.BlockID != blockID {
			continue
		}
		result = append(result, inc)
	}
	return result, nil
}
