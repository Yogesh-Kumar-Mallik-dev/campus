/**
 * BLOCK_TYPES_HUB_001
 * Subsystem: Rank 17 - The Hub Root Super-App (hub)
 * Purpose:   TypeScript domain definitions for Unified Persona Cockpits, Cross-subsystem KPIs, Widget Configs, and 1-Tap Shortcuts.
 */

export type HubPersonaType = 'STUDENT' | 'FACULTY' | 'WARDEN' | 'LIBRARIAN' | 'ADMIN' | 'SUPER_ADMIN';

export type HubWidgetSize = 'SMALL' | 'MEDIUM' | 'LARGE' | 'FULL_WIDTH';

export type HubWidgetCategory = 'METRICS' | 'ACTIONS' | 'SCHEDULE' | 'COMMUNICATION' | 'SAFETY' | 'FINANCE';

export interface HubWidgetConfig {
  id: string;
  tenantId: string;
  dashboardId: string;
  widgetKey: string;
  title: string;
  category: HubWidgetCategory;
  size: HubWidgetSize;
  orderIndex: number;
  isEnabled: boolean;
  configJson?: string;
  createdAt: string;
  updatedAt: string;
}

export interface HubPersonaDashboard {
  id: string;
  tenantId: string;
  userId: string;
  persona: HubPersonaType;
  layoutTheme: string;
  isDefault: boolean;
  widgets: HubWidgetConfig[];
  createdAt: string;
  updatedAt: string;
}

export interface HubQuickActionShortcut {
  id: string;
  tenantId: string;
  persona: HubPersonaType;
  shortcutKey: string;
  title: string;
  description?: string;
  iconName: string;
  targetRoute: string;
  orderIndex: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface HubAggregatedMetrics {
  persona: HubPersonaType;
  attendanceRate: number;
  pendingInvoicesCount: number;
  unpaidBalanceTotal: number;
  activeLibraryBorrows: number;
  overdueBooksCount: number;
  activeGatePassStatus: string;
  upcomingEventsCount: number;
  openHelpdeskTickets: number;
  activeSOSCount: number;
  unreadNoticesCount: number;
  pendingAssignments: number;
  mentorshipStatus: string;
}

export interface HubPersonaDashboardView {
  dashboard: HubPersonaDashboard;
  metrics: HubAggregatedMetrics;
  shortcuts: HubQuickActionShortcut[];
  unreadNoticesCount: number;
  activeEmergencyCount: number;
}
