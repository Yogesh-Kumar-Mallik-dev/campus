/**
 * BLOCK_WEB_ONBOARDING_TYPES_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   TypeScript interfaces matching Go domain models & REST API schemas.
 */

export type OnboardingType = 'STUDENT' | 'FACULTY' | 'STAFF' | 'WARDEN' | 'LIBRARIAN';
export type OnboardingStatus = 'DRAFT' | 'SUBMITTED' | 'UNDER_REVIEW' | 'VERIFIED' | 'REJECTED' | 'ENROLLED';
export type DocumentType =
  | 'NATIONAL_ID'
  | 'PASSPORT'
  | 'BIRTH_CERTIFICATE'
  | 'ACADEMIC_TRANSCRIPT'
  | 'DEGREE_CERTIFICATE'
  | 'TRANSFER_CERTIFICATE'
  | 'MEDICAL_FITNESS'
  | 'PASSPORT_PHOTO';

export type DocumentStatus = 'PENDING_REVIEW' | 'VERIFIED' | 'REJECTED';
export type Gender = 'MALE' | 'FEMALE' | 'OTHER';

export interface EmergencyContact {
  name: string;
  phone: string;
  relation: string;
}

export interface Address {
  line1: string;
  line2?: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
}

export interface Guardian {
  name: string;
  email?: string;
  phone: string;
  relation: string;
}

export interface OnboardingDocument {
  id: string;
  tenant_id: string;
  applicant_id: string;
  document_type: DocumentType;
  file_key: string;
  file_name: string;
  file_size: number;
  mime_type: string;
  status: DocumentStatus;
  rejection_reason?: string;
  verified_by_id?: string;
  verified_at?: string;
  created_at: string;
  updated_at: string;
}

export interface OnboardingApplicant {
  id: string;
  tenant_id: string;
  type: OnboardingType;
  status: OnboardingStatus;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  date_of_birth: string;
  gender: Gender;
  blood_group?: string;
  nationality: string;
  emergency_contact: EmergencyContact;
  address: Address;
  program_id?: string;
  department_id?: string;
  academic_year: string;
  target_designation?: string;
  guardian?: Guardian;
  reviewer_id?: string;
  reviewed_at?: string;
  rejection_reason?: string;
  documents: OnboardingDocument[];
  created_at: string;
  updated_at: string;
}

export interface AcademicDepartment {
  id: string;
  tenant_id: string;
  code: string;
  name: string;
  description?: string;
  is_active: boolean;
}

export interface AcademicProgram {
  id: string;
  tenant_id: string;
  department_id: string;
  code: string;
  name: string;
  degree_type: 'UG' | 'PG' | 'DIPLOMA' | 'DOCTORAL';
  duration_years: number;
  total_semesters: number;
  is_active: boolean;
}

export interface CohortBatch {
  id: string;
  tenant_id: string;
  program_id: string;
  academic_year: string;
  start_year: number;
  end_year: number;
  section: string;
  max_capacity: number;
  current_enrolled: number;
  is_active: boolean;
}

export interface StudentProfile {
  id: string;
  tenant_id: string;
  user_id: string;
  applicant_id: string;
  roll_number: string;
  registration_number: string;
  program_id: string;
  cohort_id: string;
  current_semester: number;
  enrollment_date: string;
}

export interface StaffProfile {
  id: string;
  tenant_id: string;
  user_id: string;
  applicant_id: string;
  employee_id: string;
  department_id: string;
  designation: string;
  role_key: string;
  joining_date: string;
}
