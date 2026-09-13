/**
 * BLOCK_TYPES_HELPDESK_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   TypeScript domain definitions for ticketing, SLA monitoring, live chat resolution, and escalations.
 */

export type HelpdeskPriority = 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';

export type HelpdeskTicketStatus =
  | 'OPEN'
  | 'IN_PROGRESS'
  | 'WAITING_FOR_APPLICANT'
  | 'RESOLVED'
  | 'CLOSED';

export type HelpdeskEscalationStatus = 'PENDING' | 'ACKNOWLEDGED' | 'RESOLVED';

export interface HelpdeskCategory {
  id: string;
  tenantId: string;
  code: string;
  name: string;
  description?: string;
  defaultPriority: HelpdeskPriority;
  slaResponseHours: number;
  slaResolutionHours: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface HelpdeskTicket {
  id: string;
  tenantId: string;
  ticketNumber: string;
  categoryId: string;
  categoryName?: string;
  requesterId: string;
  requesterName?: string;
  assignedStaffId?: string;
  assignedStaffName?: string;
  title: string;
  description: string;
  priority: HelpdeskPriority;
  status: HelpdeskTicketStatus;
  slaDueAt: string;
  firstResponseAt?: string;
  resolvedAt?: string;
  closedAt?: string;
  rating?: number;
  feedback?: string;
  createdAt: string;
  updatedAt: string;
}

export interface HelpdeskMessage {
  id: string;
  tenantId: string;
  ticketId: string;
  senderId: string;
  senderName?: string;
  isStaffReply: boolean;
  isInternalNote: boolean;
  message: string;
  attachments?: string[];
  createdAt: string;
}

export interface HelpdeskEscalation {
  id: string;
  tenantId: string;
  ticketId: string;
  escalatedToId?: string;
  escalatedToName?: string;
  reason: string;
  status: HelpdeskEscalationStatus;
  createdAt: string;
  updatedAt: string;
}
