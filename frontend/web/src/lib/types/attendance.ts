/**
 * BLOCK_WEB_ATTENDANCE_TYPES_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   TypeScript domain interfaces for timetable, sessions, multi-mode roll call, shortage tracking, and medical leave.
 */

export type AttendanceSessionStatus = 'SCHEDULED' | 'OPEN' | 'LOCKED' | 'FINALIZED';
export type AttendanceMode = 'MANUAL_FACULTY' | 'BIOMETRIC_TERMINAL' | 'RFID_SCAN' | 'GEOFENCE_MOBILE';
export type AttendanceRecordStatus = 'PRESENT' | 'ABSENT' | 'LATE' | 'EXCUSED_MEDICAL';
export type MedicalLeaveStatus = 'PENDING' | 'APPROVED' | 'REJECTED';
export type DayOfWeek = 'MONDAY' | 'TUESDAY' | 'WEDNESDAY' | 'THURSDAY' | 'FRIDAY' | 'SATURDAY' | 'SUNDAY';

export interface CourseSubject {
  id: string;
  tenant_id: string;
  program_id: string;
  code: string;
  name: string;
  credits: number;
  semester: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface TimetableSlot {
  id: string;
  tenant_id: string;
  subject_id: string;
  cohort_id: string;
  faculty_id: string;
  day_of_week: DayOfWeek;
  start_time: string;
  end_time: string;
  room_code: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AttendanceRecord {
  id: string;
  tenant_id: string;
  session_id: string;
  student_id: string;
  student_name?: string;
  roll_number?: string;
  status: AttendanceRecordStatus;
  remarks?: string;
  marked_at: string;
  device_id?: string;
  latitude?: number;
  longitude?: number;
  created_at: string;
  updated_at: string;
}

export interface AttendanceSession {
  id: string;
  tenant_id: string;
  slot_id?: string;
  subject_id: string;
  subject_name?: string;
  faculty_id: string;
  cohort_id: string;
  session_date: string;
  start_time: string;
  end_time: string;
  mode: AttendanceMode;
  status: AttendanceSessionStatus;
  opened_at?: string;
  locked_at?: string;
  total_students: number;
  present_count: number;
  absent_count: number;
  records?: AttendanceRecord[];
  created_at: string;
  updated_at: string;
}

export interface MedicalLeaveApplication {
  id: string;
  tenant_id: string;
  student_id: string;
  student_name?: string;
  from_date: string;
  to_date: string;
  reason: string;
  certificate_key: string;
  status: MedicalLeaveStatus;
  approved_by_id?: string;
  approved_at?: string;
  rejection_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface StudentAttendanceSummary {
  student_id: string;
  subject_id: string;
  total_sessions: number;
  present_sessions: number;
  excused_sessions: number;
  absent_sessions: number;
  attendance_percentage: number;
  is_shortage: boolean; // true if percentage < 75.0%
}
