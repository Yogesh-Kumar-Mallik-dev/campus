/**
 * BLOCK_TYPES_PORTAL_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   TypeScript domain definitions for public landing, academic program offerings, and admissions prospect leads.
 */

export type PortalDegreeType = 'UG' | 'PG' | 'DIPLOMA' | 'PHD';

export type PortalInquiryStatus = 'NEW' | 'CONTACTED' | 'CONVERTED' | 'CLOSED';

export interface PortalLandingPage {
  id: string;
  tenantId: string;
  heroHeadline: string;
  heroSubheadline: string;
  admissionsOpen: boolean;
  admissionsCycle: string;
  heroImageUrl?: string;
  contactEmail: string;
  contactPhone: string;
  campusAddress: string;
  isPublished: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface PortalProgramCatalog {
  id: string;
  tenantId: string;
  programCode: string;
  programName: string;
  degreeType: PortalDegreeType;
  departmentName: string;
  durationYears: number;
  totalSemesters: number;
  eligibilityCriteria: string;
  annualFee: number;
  isFeatured: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface PortalPublicInquiry {
  id: string;
  tenantId: string;
  prospectName: string;
  prospectEmail: string;
  prospectPhone: string;
  programOfInterest: string;
  message: string;
  status: PortalInquiryStatus;
  assignedCounselorId?: string;
  assignedCounselorName?: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}
