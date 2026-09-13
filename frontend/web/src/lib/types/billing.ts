/**
 * BLOCK_WEB_BILLING_TYPES_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   TypeScript domain interfaces for fee structures, invoices, payment transactions, receipts, and double-entry ledger.
 */

export type FeeCategory =
  | 'TUITION'
  | 'ADMISSION'
  | 'EXAMINATION'
  | 'HOSTEL'
  | 'MESS'
  | 'LIBRARY'
  | 'TRANSPORTATION'
  | 'LABORATORY'
  | 'MISCELLANEOUS';

export type InvoiceStatus =
  | 'DRAFT'
  | 'ISSUED'
  | 'PARTIALLY_PAID'
  | 'PAID'
  | 'OVERDUE'
  | 'CANCELLED'
  | 'WRITTEN_OFF';

export type PaymentStatus =
  | 'INITIATED'
  | 'PROCESSING'
  | 'SUCCESS'
  | 'FAILED'
  | 'REFUNDED';

export type PaymentMethod =
  | 'ONLINE_GATEWAY'
  | 'BANK_TRANSFER'
  | 'UPI'
  | 'CHEQUE'
  | 'CASH'
  | 'SCHOLARSHIP_WAIVER';

export type LedgerAccountType =
  | 'ASSET'
  | 'LIABILITY'
  | 'EQUITY'
  | 'REVENUE'
  | 'EXPENSE';

export type LedgerEntryType = 'DEBIT' | 'CREDIT';

export interface FeeStructureItem {
  id: string;
  tenant_id: string;
  fee_structure_id: string;
  category: FeeCategory;
  name: string;
  amount: number; // in cents/paise
  is_optional: boolean;
  created_at: string;
  updated_at: string;
}

export interface FeeStructure {
  id: string;
  tenant_id: string;
  program_id: string;
  academic_year: string;
  semester: number;
  name: string;
  total_amount: number;
  currency: string;
  due_date: string;
  is_active: boolean;
  items?: FeeStructureItem[];
  created_at: string;
  updated_at: string;
}

export interface InvoiceItem {
  id: string;
  tenant_id: string;
  invoice_id: string;
  category: FeeCategory;
  name: string;
  amount: number;
  created_at: string;
  updated_at: string;
}

export interface StudentInvoice {
  id: string;
  tenant_id: string;
  student_id: string;
  student_name?: string;
  fee_structure_id?: string;
  invoice_number: string;
  academic_year: string;
  semester: number;
  subtotal_amount: number;
  discount_amount: number;
  tax_amount: number;
  total_amount: number;
  paid_amount: number;
  balance_amount: number;
  due_date: string;
  status: InvoiceStatus;
  notes?: string;
  items?: InvoiceItem[];
  created_at: string;
  updated_at: string;
}

export interface PaymentTransaction {
  id: string;
  tenant_id: string;
  invoice_id: string;
  student_id: string;
  transaction_ref: string;
  gateway_name: string;
  gateway_order_id?: string;
  gateway_payment_id?: string;
  amount: number;
  currency: string;
  method: PaymentMethod;
  status: PaymentStatus;
  idempotency_key?: string;
  receipt_number?: string;
  failure_reason?: string;
  paid_at?: string;
  created_at: string;
  updated_at: string;
}

export interface LedgerEntry {
  id: string;
  tenant_id: string;
  transaction_id?: string;
  account_id: string;
  account_name?: string;
  entry_type: LedgerEntryType;
  amount: number;
  description: string;
  posted_at: string;
  created_at: string;
}
