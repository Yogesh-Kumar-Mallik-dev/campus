/**
 * BLOCK_TYPES_LIBRARY_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   TypeScript domain models for bibliographic catalog, physical copies, checkout loans, and overdue fines.
 */

export type BookCategory =
  | 'COMPUTER_SCIENCE'
  | 'ELECTRONICS'
  | 'MECHANICAL'
  | 'MATHEMATICS'
  | 'PHYSICS'
  | 'LITERATURE'
  | 'MANAGEMENT'
  | 'GENERAL';

export type BookCopyStatus = 'AVAILABLE' | 'ISSUED' | 'RESERVED' | 'MAINTENANCE' | 'LOST';
export type BorrowStatus = 'ISSUED' | 'RETURNED' | 'OVERDUE' | 'LOST';
export type ReservationStatus = 'PENDING' | 'READY_FOR_PICKUP' | 'FULFILLED' | 'EXPIRED' | 'CANCELLED';

export interface LibraryBook {
  id: string;
  tenant_id: string;
  isbn: string;
  title: string;
  author: string;
  publisher: string;
  edition?: string;
  category: BookCategory;
  total_copies: number;
  available_copies: number;
  shelf_location: string;
  ebook_key?: string;
  ebook_format?: string;
  ebook_size?: number;
  created_at: string;
}

export interface LibraryBookCopy {
  id: string;
  tenant_id: string;
  book_id: string;
  accession_number: string;
  barcode: string;
  status: BookCopyStatus;
}

export interface LibraryBorrowRecord {
  id: string;
  tenant_id: string;
  copy_id: string;
  book_title?: string;
  isbn?: string;
  accession_number?: string;
  student_id: string;
  student_name?: string;
  roll_number?: string;
  borrowed_at: string;
  due_date: string;
  returned_at?: string;
  renew_count: number;
  status: BorrowStatus;
  fine_amount: number;
  fine_paid: boolean;
  issued_by_id?: string;
  created_at: string;
}

export interface LibraryReservation {
  id: string;
  tenant_id: string;
  book_id: string;
  book_title?: string;
  student_id: string;
  student_name?: string;
  roll_number?: string;
  reserved_at: string;
  expires_at: string;
  status: ReservationStatus;
  created_at: string;
}

export interface LibraryEbookAccess {
  id: string;
  tenant_id: string;
  book_id: string;
  student_id: string;
  last_page_read: number;
  total_reading_minutes: number;
  last_accessed_at: string;
}
