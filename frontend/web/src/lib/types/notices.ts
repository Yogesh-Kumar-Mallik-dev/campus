/**
 * BLOCK_WEB_NOTICES_TYPES_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   TypeScript domain interfaces for notices, categories, priority levels, attachments, and read receipts.
 */

export type NoticeCategory =
  | 'ACADEMIC'
  | 'EXAMINATION'
  | 'ADMINISTRATIVE'
  | 'HOSTEL'
  | 'EVENTS'
  | 'EMERGENCY'
  | 'PLACEMENT'
  | 'GENERAL';

export type NoticePriority = 'LOW' | 'NORMAL' | 'HIGH' | 'URGENT';
export type NoticeStatus = 'DRAFT' | 'PUBLISHED' | 'ARCHIVED' | 'EXPIRED';
export type TargetAudience = 'ALL' | 'STUDENTS' | 'FACULTY' | 'STAFF' | 'HOSTEL_RESIDENTS';

export interface NoticeAttachment {
  id: string;
  tenant_id: string;
  notice_id: string;
  file_name: string;
  file_key: string;
  file_size: number;
  mime_type: string;
  created_at: string;
}

export interface NoticeAcknowledgement {
  id: string;
  tenant_id: string;
  notice_id: string;
  user_id: string;
  read_at: string;
  acknowledged_at?: string;
  device_id?: string;
  created_at: string;
}

export interface CampusNotice {
  id: string;
  tenant_id: string;
  author_id: string;
  author_name?: string;
  title: string;
  slug: string;
  content: string;
  category: NoticeCategory;
  priority: NoticePriority;
  status: NoticeStatus;
  target_audience: TargetAudience;
  target_dept_id?: string;
  target_program_id?: string;
  target_cohort_id?: string;
  is_pinned: boolean;
  publish_at: string;
  expires_at?: string;
  view_count: number;
  attachments?: NoticeAttachment[];
  acknowledgements?: NoticeAcknowledgement[];
  is_acknowledged?: boolean;
  created_at: string;
  updated_at: string;
}
