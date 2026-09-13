/**
 * BLOCK_TYPES_STUDYHUB_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   TypeScript domain models for courses, materials, assignments, submissions, and peer reviews.
 */

export type StudyMaterialType =
  | 'NOTE'
  | 'SLIDE'
  | 'LAB_MANUAL'
  | 'RECORDING'
  | 'SAMPLE_PAPER';

export type StudyAssignmentStatus = 'DRAFT' | 'PUBLISHED' | 'CLOSED';

export type StudySubmissionStatus = 'SUBMITTED' | 'GRADED' | 'LATE' | 'REJECTED';

export interface StudyCourse {
  id: string;
  tenantId: string;
  code: string;
  title: string;
  description?: string;
  department: string;
  credits: number;
  semester: number;
  instructorId?: string;
  syllabusText?: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface StudyMaterial {
  id: string;
  tenantId: string;
  courseId: string;
  title: string;
  description?: string;
  unitNumber: number;
  materialType: StudyMaterialType;
  fileUrl: string;
  fileSizeBytes: number;
  publishedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface StudyAssignment {
  id: string;
  tenantId: string;
  courseId: string;
  title: string;
  description: string;
  maxMarks: number;
  dueDate: string;
  allowLateSubmission: boolean;
  latePenaltyPercentPerDay: number;
  status: StudyAssignmentStatus;
  createdAt: string;
  updatedAt: string;
}

export interface StudySubmission {
  id: string;
  tenantId: string;
  assignmentId: string;
  studentId: string;
  studentName?: string;
  rollNumber?: string;
  fileUrl?: string;
  contentText?: string;
  submittedAt: string;
  status: StudySubmissionStatus;
  marksObtained?: number;
  feedback?: string;
  gradedById?: string;
  gradedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface StudyPeerReview {
  id: string;
  submissionId: string;
  reviewerStudentId: string;
  score: number;
  comments: string;
  reviewedAt: string;
}
