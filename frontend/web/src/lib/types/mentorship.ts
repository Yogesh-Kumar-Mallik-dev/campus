/**
 * BLOCK_TYPES_MENTORSHIP_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   TypeScript domain models for mentor allocations, counseling sessions, academic progress, and risk alerts.
 */

export type MentorAllocationStatus = 'ACTIVE' | 'COMPLETED' | 'REASSIGNED';

export type MentorshipMeetingType =
  | 'ONE_ON_ONE'
  | 'GROUP'
  | 'ACADEMIC_REVIEW'
  | 'EMERGENCY_COUNSELING';

export type MentorshipSessionStatus =
  | 'SCHEDULED'
  | 'COMPLETED'
  | 'CANCELLED'
  | 'NO_SHOW';

export type AtRiskStatus = 'NORMAL' | 'WATCHLIST' | 'CRITICAL_INTERVENTION';

export type AtRiskType =
  | 'ACADEMIC_PROBATION'
  | 'ATTENDANCE_SHORTAGE'
  | 'DISCIPLINARY'
  | 'MENTAL_HEALTH';

export type RiskSeverity = 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';

export interface MentorAllocation {
  id: string;
  tenantId: string;
  studentId: string;
  studentName?: string;
  rollNumber?: string;
  mentorStaffId: string;
  mentorName?: string;
  mentorDesignation?: string;
  cohortId?: string;
  status: MentorAllocationStatus;
  allocatedAt: string;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface MentorshipSession {
  id: string;
  tenantId: string;
  allocationId: string;
  studentName?: string;
  scheduledAt: string;
  completedAt?: string;
  location?: string;
  meetingType: MentorshipMeetingType;
  status: MentorshipSessionStatus;
  discussionSummary?: string;
  actionItems?: string;
  followUpDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface StudentAcademicProgress {
  id: string;
  tenantId: string;
  studentId: string;
  studentName?: string;
  rollNumber?: string;
  semester: number;
  sgpa: number;
  cgpa: number;
  attendancePercentage: number;
  creditsEarned: number;
  totalCredits: number;
  atRiskStatus: AtRiskStatus;
  remarks?: string;
  evaluatedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface MentorshipAtRiskAlert {
  id: string;
  tenantId: string;
  studentId: string;
  studentName?: string;
  rollNumber?: string;
  riskType: AtRiskType;
  severity: RiskSeverity;
  description: string;
  resolvedAt?: string;
  resolvedById?: string;
  createdAt: string;
  updatedAt: string;
}
