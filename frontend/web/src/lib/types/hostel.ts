/**
 * BLOCK_TYPES_HOSTEL_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   TypeScript domain models for hostel residential blocks, bed allocations, and gate passes.
 */

export type HostelGender = 'MALE' | 'FEMALE' | 'COED';
export type HostelRoomType = 'SINGLE' | 'DOUBLE' | 'TRIPLE' | 'FOUR_SHARING' | 'DORMITORY';
export type HostelRoomStatus = 'AVAILABLE' | 'OCCUPIED' | 'UNDER_MAINTENANCE';
export type HostelBedStatus = 'AVAILABLE' | 'ALLOCATED' | 'RESERVED' | 'MAINTENANCE';
export type HostelAllocationStatus = 'ALLOCATED' | 'VACATED' | 'TRANSFERRED' | 'CANCELLED';
export type GatePassStatus = 'PENDING' | 'APPROVED' | 'REJECTED' | 'OUT_CAMPUS' | 'RETURNED' | 'EXPIRED' | 'CANCELLED';
export type HostelIncidentType = 'CURFEW_VIOLATION' | 'UNAUTHORIZED_GUEST' | 'NOISE_DISTURBANCE' | 'PROPERTY_DAMAGE' | 'SUBSTANCE_VIOLATION' | 'OTHER';
export type IncidentSeverity = 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';

export interface HostelBlock {
  id: string;
  tenant_id: string;
  name: string;
  code: string;
  gender: HostelGender;
  total_floors: number;
  total_rooms: number;
  capacity: number;
  warden_id?: string;
  is_active: boolean;
  created_at: string;
}

export interface HostelRoom {
  id: string;
  tenant_id: string;
  block_id: string;
  room_number: string;
  floor_number: number;
  room_type: HostelRoomType;
  is_ac: boolean;
  base_fee_per_semester: number;
  status: HostelRoomStatus;
  max_beds: number;
}

export interface HostelBed {
  id: string;
  tenant_id: string;
  room_id: string;
  bed_number: string;
  status: HostelBedStatus;
}

export interface HostelAllocation {
  id: string;
  tenant_id: string;
  bed_id: string;
  student_id: string;
  student_name?: string;
  roll_number?: string;
  academic_year: string;
  semester: number;
  allocated_at: string;
  vacated_at?: string;
  status: HostelAllocationStatus;
  remarks?: string;
}

export interface HostelGatePass {
  id: string;
  tenant_id: string;
  student_id: string;
  student_name?: string;
  roll_number?: string;
  block_id: string;
  block_name?: string;
  reason: string;
  destination: string;
  emergency_contact: string;
  expected_out_at: string;
  expected_in_at: string;
  actual_out_at?: string;
  actual_in_at?: string;
  status: GatePassStatus;
  approved_by_id?: string;
  rejection_reason?: string;
  created_at: string;
}

export interface HostelIncidentLog {
  id: string;
  tenant_id: string;
  student_id: string;
  student_name?: string;
  block_id: string;
  warden_id: string;
  warden_name?: string;
  incident_type: HostelIncidentType;
  severity: IncidentSeverity;
  title: string;
  description: string;
  action_taken?: string;
  fine_amount: number;
  created_at: string;
}
