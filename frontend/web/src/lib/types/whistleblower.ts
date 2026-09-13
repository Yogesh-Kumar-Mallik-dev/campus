/**
 * BLOCK_TYPES_WHISTLEBLOWER_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   TypeScript domain definitions for zero-knowledge grievances, secret token lookups, and committee investigation.
 */

export type WhistleblowerCategory =
  | 'ANTI_RAGGING'
  | 'FINANCIAL_FRAUD'
  | 'ACADEMIC_CORRUPTION'
  | 'HARASSMENT'
  | 'SAFETY_VIOLATION'
  | 'OTHER';

export type WhistleblowerSeverity = 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';

export type WhistleblowerStatus =
  | 'SUBMITTED'
  | 'UNDER_INVESTIGATION'
  | 'EVIDENCE_REQUESTED'
  | 'RESOLVED'
  | 'REJECTED';

export type WhistleblowerSenderType = 'ANONYMOUS_REPORTER' | 'INVESTIGATION_COMMITTEE';

export interface WhistleblowerReport {
  id: string;
  tenantId: string;
  reportNumber: string;
  category: WhistleblowerCategory;
  severity: WhistleblowerSeverity;
  title: string;
  description: string;
  status: WhistleblowerStatus;
  assignedOfficerId?: string;
  assignedOfficerName?: string;
  findingsSummary?: string;
  submittedAt: string;
  resolvedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface WhistleblowerMessage {
  id: string;
  tenantId: string;
  reportId: string;
  senderType: WhistleblowerSenderType;
  message: string;
  createdAt: string;
}

export interface WhistleblowerEvidence {
  id: string;
  tenantId: string;
  reportId: string;
  fileUrl: string;
  fileHash: string;
  mimeType: string;
  createdAt: string;
}
