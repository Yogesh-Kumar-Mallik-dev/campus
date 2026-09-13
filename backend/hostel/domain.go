/**
 * BLOCK_HOSTEL_DOMAIN_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   Domain entities, business invariants, bed allocation state machines, and gate pass lifecycles.
 */

package hostel

import (
	"fmt"
	"time"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderCoed   Gender = "COED"
)

type RoomType string

const (
	RoomTypeSingle      RoomType = "SINGLE"
	RoomTypeDouble      RoomType = "DOUBLE"
	RoomTypeTriple      RoomType = "TRIPLE"
	RoomTypeFourSharing RoomType = "FOUR_SHARING"
	RoomTypeDormitory   RoomType = "DORMITORY"
)

type RoomStatus string

const (
	RoomStatusAvailable        RoomStatus = "AVAILABLE"
	RoomStatusOccupied         RoomStatus = "OCCUPIED"
	RoomStatusUnderMaintenance RoomStatus = "UNDER_MAINTENANCE"
)

type BedStatus string

const (
	BedStatusAvailable   BedStatus = "AVAILABLE"
	BedStatusAllocated   BedStatus = "ALLOCATED"
	BedStatusReserved    BedStatus = "RESERVED"
	BedStatusMaintenance BedStatus = "MAINTENANCE"
)

type AllocationStatus string

const (
	AllocationStatusAllocated   AllocationStatus = "ALLOCATED"
	AllocationStatusVacated     AllocationStatus = "VACATED"
	AllocationStatusTransferred AllocationStatus = "TRANSFERRED"
	AllocationStatusCancelled   AllocationStatus = "CANCELLED"
)

type GatePassStatus string

const (
	GatePassStatusPending   GatePassStatus = "PENDING"
	GatePassStatusApproved  GatePassStatus = "APPROVED"
	GatePassStatusRejected  GatePassStatus = "REJECTED"
	GatePassStatusOutCampus GatePassStatus = "OUT_CAMPUS"
	GatePassStatusReturned  GatePassStatus = "RETURNED"
	GatePassStatusExpired   GatePassStatus = "EXPIRED"
	GatePassStatusCancelled GatePassStatus = "CANCELLED"
)

type IncidentType string

const (
	IncidentCurfewViolation    IncidentType = "CURFEW_VIOLATION"
	IncidentUnauthorizedGuest  IncidentType = "UNAUTHORIZED_GUEST"
	IncidentNoiseDisturbance   IncidentType = "NOISE_DISTURBANCE"
	IncidentPropertyDamage     IncidentType = "PROPERTY_DAMAGE"
	IncidentSubstanceViolation IncidentType = "SUBSTANCE_VIOLATION"
	IncidentOther              IncidentType = "OTHER"
)

type IncidentSeverity string

const (
	SeverityLow      IncidentSeverity = "LOW"
	SeverityMedium   IncidentSeverity = "MEDIUM"
	SeverityHigh     IncidentSeverity = "HIGH"
	SeverityCritical IncidentSeverity = "CRITICAL"
)

// Block represents a physical hostel residence building
type Block struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Gender      Gender    `json:"gender"`
	TotalFloors int       `json:"total_floors"`
	TotalRooms  int       `json:"total_rooms"`
	Capacity    int       `json:"capacity"`
	WardenID    *string   `json:"warden_id,omitempty"`
	IsActive    Boolean   `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Boolean = bool

// Room represents an accommodation unit within a hostel block
type Room struct {
	ID                 string     `json:"id"`
	TenantID           string     `json:"tenant_id"`
	BlockID            string     `json:"block_id"`
	RoomNumber         string     `json:"room_number"`
	FloorNumber        int        `json:"floor_number"`
	RoomType           RoomType   `json:"room_type"`
	IsAC               bool       `json:"is_ac"`
	BaseFeePerSemester float64    `json:"base_fee_per_semester"`
	Status             RoomStatus `json:"status"`
	MaxBeds            int        `json:"max_beds"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// Bed represents a single bed within a hostel room
type Bed struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	RoomID    string    `json:"room_id"`
	BedNumber string    `json:"bed_number"`
	Status    BedStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Allocation represents a student's active or historic housing assignment
type Allocation struct {
	ID           string           `json:"id"`
	TenantID     string           `json:"tenant_id"`
	BedID        string           `json:"bed_id"`
	StudentID    string           `json:"student_id"`
	AcademicYear string           `json:"academic_year"`
	Semester     int              `json:"semester"`
	AllocatedAt  time.Time        `json:"allocated_at"`
	VacatedAt    *time.Time       `json:"vacated_at,omitempty"`
	Status       AllocationStatus `json:"status"`
	Remarks      string           `json:"remarks,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// GatePass represents an authorized digital exit and re-entry pass
type GatePass struct {
	ID               string         `json:"id"`
	TenantID         string         `json:"tenant_id"`
	StudentID        string         `json:"student_id"`
	BlockID          string         `json:"block_id"`
	Reason           string         `json:"reason"`
	Destination      string         `json:"destination"`
	EmergencyContact string         `json:"emergency_contact"`
	ExpectedOutAt    time.Time      `json:"expected_out_at"`
	ExpectedInAt     time.Time      `json:"expected_in_at"`
	ActualOutAt      *time.Time     `json:"actual_out_at,omitempty"`
	ActualInAt       *time.Time     `json:"actual_in_at,omitempty"`
	Status           GatePassStatus `json:"status"`
	ApprovedByID     *string        `json:"approved_by_id,omitempty"`
	RejectionReason  *string        `json:"rejection_reason,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// IncidentLog represents disciplinary or curfew infraction logs
type IncidentLog struct {
	ID           string           `json:"id"`
	TenantID     string           `json:"tenant_id"`
	StudentID    string           `json:"student_id"`
	BlockID      string           `json:"block_id"`
	WardenID     string           `json:"warden_id"`
	IncidentType IncidentType     `json:"incident_type"`
	Severity     IncidentSeverity `json:"severity"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	ActionTaken  string           `json:"action_taken,omitempty"`
	FineAmount   float64          `json:"fine_amount"`
	ResolvedAt   *time.Time       `json:"resolved_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// ValidateRoomCapacity verifies bed capacity constraints for room type
func ValidateRoomCapacity(roomType RoomType, maxBeds int) error {
	switch roomType {
	case RoomTypeSingle:
		if maxBeds != 1 {
			return fmt.Errorf("single room must have exactly 1 bed, got %d", maxBeds)
		}
	case RoomTypeDouble:
		if maxBeds != 2 {
			return fmt.Errorf("double room must have exactly 2 beds, got %d", maxBeds)
		}
	case RoomTypeTriple:
		if maxBeds != 3 {
			return fmt.Errorf("triple room must have exactly 3 beds, got %d", maxBeds)
		}
	case RoomTypeFourSharing:
		if maxBeds != 4 {
			return fmt.Errorf("four-sharing room must have exactly 4 beds, got %d", maxBeds)
		}
	case RoomTypeDormitory:
		if maxBeds < 5 {
			return fmt.Errorf("dormitory room must have at least 5 beds, got %d", maxBeds)
		}
	}
	return nil
}

// ValidateGatePassTimes verifies out and in timestamps
func ValidateGatePassTimes(outTime, inTime time.Time) error {
	if !inTime.After(outTime) {
		return ErrInvalidTimeRange
	}
	return nil
}

// CanTransitionGatePass evaluates validity of state transitions
func CanTransitionGatePass(current, target GatePassStatus) bool {
	switch current {
	case GatePassStatusPending:
		return target == GatePassStatusApproved || target == GatePassStatusRejected || target == GatePassStatusCancelled
	case GatePassStatusApproved:
		return target == GatePassStatusOutCampus || target == GatePassStatusExpired || target == GatePassStatusCancelled
	case GatePassStatusOutCampus:
		return target == GatePassStatusReturned
	default:
		return false
	}
}

// CheckCurfewBreach checks if actual return exceeded expected return time
func CheckCurfewBreach(expectedIn, actualIn time.Time) (isBreach bool, overdueMinutes int) {
	if actualIn.After(expectedIn) {
		diff := actualIn.Sub(expectedIn)
		return true, int(diff.Minutes())
	}
	return false, 0
}
