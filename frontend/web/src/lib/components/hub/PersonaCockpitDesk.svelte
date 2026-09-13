<!--
  BLOCK_WEB_HUB_COCKPIT_001
  Subsystem: Rank 17 - The Hub Root Super-App (hub)
  Purpose:   Unified Multi-Persona Root Cockpit aggregating KPIs, 1-tap shortcuts, and live alerts across all 16 campus subsystems.
-->
<script lang="ts">
  import type {
    HubPersonaType,
    HubAggregatedMetrics,
    HubQuickActionShortcut,
    HubWidgetConfig,
  } from '$lib/types/hub';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';

  let activePersona = $state<HubPersonaType>('STUDENT');

  let metrics = $derived<HubAggregatedMetrics>(getMetricsForPersona(activePersona));
  let shortcuts = $derived<HubQuickActionShortcut[]>(getShortcutsForPersona(activePersona));

  function getMetricsForPersona(persona: HubPersonaType): HubAggregatedMetrics {
    switch (persona) {
      case 'STUDENT':
        return {
          persona: 'STUDENT',
          attendanceRate: 89.2,
          pendingInvoicesCount: 1,
          unpaidBalanceTotal: 25000,
          activeLibraryBorrows: 2,
          overdueBooksCount: 0,
          activeGatePassStatus: 'APPROVED',
          upcomingEventsCount: 3,
          openHelpdeskTickets: 1,
          activeSOSCount: 0,
          unreadNoticesCount: 4,
          pendingAssignments: 2,
          mentorshipStatus: 'ON_TRACK',
        };
      case 'FACULTY':
        return {
          persona: 'FACULTY',
          attendanceRate: 94.5,
          pendingInvoicesCount: 0,
          unpaidBalanceTotal: 0,
          activeLibraryBorrows: 5,
          overdueBooksCount: 0,
          activeGatePassStatus: 'N/A',
          upcomingEventsCount: 5,
          openHelpdeskTickets: 2,
          activeSOSCount: 0,
          unreadNoticesCount: 2,
          pendingAssignments: 18,
          mentorshipStatus: '8_MENTEES_ACTIVE',
        };
      case 'WARDEN':
        return {
          persona: 'WARDEN',
          attendanceRate: 0,
          pendingInvoicesCount: 0,
          unpaidBalanceTotal: 0,
          activeLibraryBorrows: 0,
          overdueBooksCount: 0,
          activeGatePassStatus: '14_PENDING_APPROVAL',
          upcomingEventsCount: 1,
          openHelpdeskTickets: 4,
          activeSOSCount: 0,
          unreadNoticesCount: 1,
          pendingAssignments: 0,
          mentorshipStatus: 'N/A',
        };
      case 'LIBRARIAN':
        return {
          persona: 'LIBRARIAN',
          attendanceRate: 0,
          pendingInvoicesCount: 0,
          unpaidBalanceTotal: 0,
          activeLibraryBorrows: 342,
          overdueBooksCount: 12,
          activeGatePassStatus: 'N/A',
          upcomingEventsCount: 0,
          openHelpdeskTickets: 1,
          activeSOSCount: 0,
          unreadNoticesCount: 2,
          pendingAssignments: 0,
          mentorshipStatus: 'N/A',
        };
      case 'ADMIN':
      case 'SUPER_ADMIN':
        return {
          persona: persona,
          attendanceRate: 91.8,
          pendingInvoicesCount: 42,
          unpaidBalanceTotal: 1050000,
          activeLibraryBorrows: 1240,
          overdueBooksCount: 38,
          activeGatePassStatus: 'ALL_NORMAL',
          upcomingEventsCount: 8,
          openHelpdeskTickets: 9,
          activeSOSCount: 0,
          unreadNoticesCount: 6,
          pendingAssignments: 154,
          mentorshipStatus: 'INSTITUTE_HEALTHY',
        };
    }
  }

  function getShortcutsForPersona(persona: HubPersonaType): HubQuickActionShortcut[] {
    switch (persona) {
      case 'STUDENT':
        return [
          {
            id: 'sc-1',
            tenantId: 'tenant-demo',
            persona: 'STUDENT',
            shortcutKey: 'TRIGGER_SOS',
            title: 'SOS Emergency',
            description: '1-Tap emergency distress trigger',
            iconName: '🚨',
            targetRoute: '/sos',
            orderIndex: 1,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-2',
            tenantId: 'tenant-demo',
            persona: 'STUDENT',
            shortcutKey: 'MARK_ATTENDANCE',
            title: 'Attendance Scan',
            description: 'Geofenced biometric check-in',
            iconName: '📍',
            targetRoute: '/attendance',
            orderIndex: 2,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-3',
            tenantId: 'tenant-demo',
            persona: 'STUDENT',
            shortcutKey: 'PAY_FEE',
            title: 'Pay Semester Fee',
            description: 'Direct UPI & netbanking gateway',
            iconName: '💳',
            targetRoute: '/billing',
            orderIndex: 3,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-4',
            tenantId: 'tenant-demo',
            persona: 'STUDENT',
            shortcutKey: 'GATE_PASS_APPLY',
            title: 'Hostel Gate Pass',
            description: 'Request day/night out pass',
            iconName: '🚪',
            targetRoute: '/hostel',
            orderIndex: 4,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-5',
            tenantId: 'tenant-demo',
            persona: 'STUDENT',
            shortcutKey: 'MESS_TOKEN',
            title: 'Mess Meal QR',
            description: 'Scan meal coupon token',
            iconName: '🍽️',
            targetRoute: '/mess',
            orderIndex: 5,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-6',
            tenantId: 'tenant-demo',
            persona: 'STUDENT',
            shortcutKey: 'CONFIDENTIAL_WB',
            title: 'Anonymous Report',
            description: 'Zero-knowledge anti-ragging report',
            iconName: '🛡️',
            targetRoute: '/whistleblower',
            orderIndex: 6,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
        ];
      case 'FACULTY':
        return [
          {
            id: 'sc-f1',
            tenantId: 'tenant-demo',
            persona: 'FACULTY',
            shortcutKey: 'TAKE_ATTENDANCE',
            title: 'Class Roll Call',
            description: 'Open QR roll call session',
            iconName: '📋',
            targetRoute: '/attendance',
            orderIndex: 1,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-f2',
            tenantId: 'tenant-demo',
            persona: 'FACULTY',
            shortcutKey: 'GRADE_HOMEWORK',
            title: 'Grade Submissions',
            description: 'Score 18 pending assignments',
            iconName: '📝',
            targetRoute: '/studyhub',
            orderIndex: 2,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-f3',
            tenantId: 'tenant-demo',
            persona: 'FACULTY',
            shortcutKey: 'PUBLISH_NOTICE',
            title: 'Post Notice',
            description: 'Broadcast academic circular',
            iconName: '📢',
            targetRoute: '/notices',
            orderIndex: 3,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-f4',
            tenantId: 'tenant-demo',
            persona: 'FACULTY',
            shortcutKey: 'SCHEDULE_MENTORSHIP',
            title: 'Mentee Reviews',
            description: 'Check 8 assigned student scores',
            iconName: '🌱',
            targetRoute: '/mentorship',
            orderIndex: 4,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
        ];
      case 'WARDEN':
        return [
          {
            id: 'sc-w1',
            tenantId: 'tenant-demo',
            persona: 'WARDEN',
            shortcutKey: 'APPROVE_GATEPASS',
            title: 'Approve Gate Passes',
            description: '14 student applications pending',
            iconName: '🔑',
            targetRoute: '/hostel',
            orderIndex: 1,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-w2',
            tenantId: 'tenant-demo',
            persona: 'WARDEN',
            shortcutKey: 'LOG_INCIDENT',
            title: 'Log Curfew Breach',
            description: 'Disciplinary infraction registry',
            iconName: '⚠️',
            targetRoute: '/hostel',
            orderIndex: 2,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-w3',
            tenantId: 'tenant-demo',
            persona: 'WARDEN',
            shortcutKey: 'ROOM_AUDIT',
            title: 'Hostel Bed Occupancy',
            description: 'Verify vacant beds & room status',
            iconName: '🛏️',
            targetRoute: '/hostel',
            orderIndex: 3,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
        ];
      default:
        return [
          {
            id: 'sc-a1',
            tenantId: 'tenant-demo',
            persona: persona,
            shortcutKey: 'APPLICANT_PIPELINE',
            title: 'Onboarding Pipeline',
            description: 'Verify fresh enrollments',
            iconName: '🎓',
            targetRoute: '/onboarding',
            orderIndex: 1,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-a2',
            tenantId: 'tenant-demo',
            persona: persona,
            shortcutKey: 'FINANCIAL_LEDGER',
            title: 'Fee Collection Ledger',
            description: 'Review fee reconciliations',
            iconName: '📊',
            targetRoute: '/billing',
            orderIndex: 2,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-a3',
            tenantId: 'tenant-demo',
            persona: persona,
            shortcutKey: 'AUDIT_LEDGER',
            title: 'Compliance Audit Trail',
            description: 'SHA-256 Merkle checkpoint proofs',
            iconName: '🔒',
            targetRoute: '/audit',
            orderIndex: 3,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          {
            id: 'sc-a4',
            tenantId: 'tenant-demo',
            persona: persona,
            shortcutKey: 'PUBLIC_PORTAL_CRM',
            title: 'Prospect Leads CRM',
            description: 'Admissions inquiry triage',
            iconName: '🌐',
            targetRoute: '/portal',
            orderIndex: 4,
            isActive: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
        ];
    }
  }

  let toastMessage = $state<string | null>(null);

  function triggerAction(sc: HubQuickActionShortcut) {
    toastMessage = `Triggered 1-Tap Action: ${sc.title} → Routing to ${sc.targetRoute}`;
    setTimeout(() => {
      toastMessage = null;
    }, 3000);
  }
</script>

<div class="space-y-6">
  <!-- Top Persona Switcher Header -->
  <div class="p-6 bg-gradient-to-br from-slate-900 via-primary/90 to-slate-950 rounded-3xl text-white shadow-xl flex flex-col md:flex-row md:items-center justify-between gap-4">
    <div class="space-y-1">
      <div class="flex items-center gap-2">
        <span class="px-2.5 py-0.5 bg-white/20 text-white border border-white/30 rounded-full text-[10px] font-extrabold uppercase tracking-wider">
          ✦ The Hub Super-App
        </span>
        <span class="text-xs text-emerald-400 font-bold flex items-center gap-1">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          16 Subsystems Online
        </span>
      </div>
      <h1 class="text-2xl md:text-3xl font-black tracking-tight">
        Enterprise Cockpit & Aggregator
      </h1>
      <p class="text-xs text-slate-300">
        Active Persona: <strong class="text-white uppercase tracking-wider">{activePersona}</strong> · Real-time cross-domain synchronization.
      </p>
    </div>

    <!-- Persona Switcher Pills -->
    <div class="flex flex-wrap gap-1.5 bg-black/40 p-1.5 rounded-2xl border border-white/10 self-start md:self-center">
      {#each ['STUDENT', 'FACULTY', 'WARDEN', 'LIBRARIAN', 'ADMIN', 'SUPER_ADMIN'] as p}
        <button
          onclick={() => (activePersona = p as HubPersonaType)}
          class="px-3 py-1 rounded-xl text-xs font-bold transition-all {activePersona === p ? 'bg-primary text-primary-foreground shadow-md' : 'text-slate-300 hover:text-white hover:bg-white/10'}"
        >
          {p}
        </button>
      {/each}
    </div>
  </div>

  {#if toastMessage}
    <div class="p-3.5 bg-emerald-500/15 border border-emerald-500/30 rounded-2xl text-emerald-700 dark:text-emerald-300 text-xs font-bold flex items-center justify-between shadow-sm animate-fade-in">
      <span>🚀 {toastMessage}</span>
      <button onclick={() => (toastMessage = null)} class="text-xs opacity-70 hover:opacity-100">✕</button>
    </div>
  {/if}

  <!-- Vital KPI Metrics Grid -->
  <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
    <!-- Metric 1: Attendance Rate -->
    <div class="p-4 bg-card rounded-2xl border shadow-sm space-y-1">
      <span class="text-[10px] font-bold text-muted-foreground uppercase">Attendance</span>
      <div class="flex items-baseline gap-1">
        <span class="text-xl font-black text-foreground">{metrics.attendanceRate}%</span>
      </div>
      <span class="text-[10px] text-emerald-600 dark:text-emerald-400 font-semibold block">
        {metrics.attendanceRate >= 75 ? '✓ Eligible for exams' : '⚠️ Low Attendance'}
      </span>
    </div>

    <!-- Metric 2: Financial Invoices -->
    <div class="p-4 bg-card rounded-2xl border shadow-sm space-y-1">
      <span class="text-[10px] font-bold text-muted-foreground uppercase">Pending Dues</span>
      <div class="text-xl font-black text-foreground">
        ₹{metrics.unpaidBalanceTotal.toLocaleString()}
      </div>
      <span class="text-[10px] text-muted-foreground block">
        {metrics.pendingInvoicesCount} invoice(s) pending
      </span>
    </div>

    <!-- Metric 3: Library Borrows -->
    <div class="p-4 bg-card rounded-2xl border shadow-sm space-y-1">
      <span class="text-[10px] font-bold text-muted-foreground uppercase">Library Books</span>
      <div class="text-xl font-black text-foreground">{metrics.activeLibraryBorrows}</div>
      <span class="text-[10px] {metrics.overdueBooksCount > 0 ? 'text-rose-500 font-bold' : 'text-muted-foreground'} block">
        {metrics.overdueBooksCount} overdue fine
      </span>
    </div>

    <!-- Metric 4: Gate Pass / Hostel -->
    <div class="p-4 bg-card rounded-2xl border shadow-sm space-y-1">
      <span class="text-[10px] font-bold text-muted-foreground uppercase">Gate Pass</span>
      <div class="text-sm font-black text-foreground truncate mt-1">
        {metrics.activeGatePassStatus}
      </div>
      <span class="text-[10px] text-muted-foreground block">Curfew: 09:30 PM</span>
    </div>

    <!-- Metric 5: Helpdesk & Grievance -->
    <div class="p-4 bg-card rounded-2xl border shadow-sm space-y-1">
      <span class="text-[10px] font-bold text-muted-foreground uppercase">Open Tickets</span>
      <div class="text-xl font-black text-foreground">{metrics.openHelpdeskTickets}</div>
      <span class="text-[10px] text-muted-foreground block">SLA active</span>
    </div>

    <!-- Metric 6: Unread Notices -->
    <div class="p-4 bg-card rounded-2xl border shadow-sm space-y-1">
      <span class="text-[10px] font-bold text-muted-foreground uppercase">Campus Notices</span>
      <div class="text-xl font-black text-primary">{metrics.unreadNoticesCount} unread</div>
      <span class="text-[10px] text-muted-foreground block">Circulars & Events</span>
    </div>
  </div>

  <!-- 1-Tap Quick Actions Section -->
  <div class="space-y-3">
    <div class="flex items-center justify-between">
      <h2 class="text-base font-black text-foreground flex items-center gap-2">
        <span>⚡</span> 1-Tap Quick Action Shortcuts ({activePersona})
      </h2>
      <span class="text-xs text-muted-foreground">Instant cross-subsystem routing</span>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
      {#each shortcuts as sc}
        <button
          onclick={() => triggerAction(sc)}
          class="p-4 bg-card rounded-2xl border shadow-sm hover:border-primary/60 hover:shadow-md transition-all text-left flex items-start gap-3 group"
        >
          <span class="text-2xl p-2 bg-muted/60 rounded-xl group-hover:scale-110 transition-transform">
            {sc.iconName}
          </span>
          <div class="space-y-0.5 flex-1 min-w-0">
            <h3 class="font-bold text-sm text-foreground truncate group-hover:text-primary transition-colors">
              {sc.title}
            </h3>
            <p class="text-[11px] text-muted-foreground line-clamp-1 leading-snug">
              {sc.description}
            </p>
            <span class="text-[10px] font-mono text-primary font-semibold block pt-1">
              {sc.targetRoute} →
            </span>
          </div>
        </button>
      {/each}
    </div>
  </div>

  <!-- Realtime Systems Radar Map -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <!-- Subsystems Health Strip -->
    <div class="p-5 bg-card rounded-2xl border shadow-sm md:col-span-2 space-y-3">
      <h3 class="text-sm font-black text-foreground">Topological Subsystem Telemetry</h3>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs">
        {#each [
          { name: '1. Auth & MFA', status: 'HEALTHY' },
          { name: '2. Audit & Merkle', status: 'SYNCED' },
          { name: '3. Onboarding', status: 'ACTIVE' },
          { name: '4. Attendance', status: 'RUNNING' },
          { name: '5. Billing & LEDGER', status: 'BALANCED' },
          { name: '6. Campus Notices', status: 'PUBLISHED' },
          { name: '7. Hostel & Curfew', status: 'NORMAL' },
          { name: '8. Mess & QR Dining', status: 'SERVED' },
          { name: '9. E-Library & RFID', status: 'CIRCULATING' },
          { name: '10. Study Hub & LMS', status: 'ONLINE' },
          { name: '11. Mentorship', status: 'ENGAGED' },
          { name: '12. Campus Events', status: 'SCHEDULED' },
          { name: '13. Helpdesk SLA', status: 'ON_TRACK' },
          { name: '14. SOS Distress', status: 'STANDBY' },
          { name: '15. Whistleblower', status: 'SHIELDED' },
          { name: '16. Public Portal', status: 'PUBLIC' },
        ] as sub}
          <div class="p-2.5 bg-muted/40 rounded-xl border flex items-center justify-between">
            <span class="font-medium text-[11px] truncate text-foreground">{sub.name}</span>
            <span class="text-[9px] font-extrabold px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-400">
              {sub.status}
            </span>
          </div>
        {/each}
      </div>
    </div>

    <!-- Live Emergency & Compliance Widget -->
    <div class="p-5 bg-card rounded-2xl border shadow-sm space-y-3">
      <h3 class="text-sm font-black text-foreground">Emergency & Vigilance</h3>
      <div class="p-4 bg-emerald-500/10 border border-emerald-500/30 rounded-2xl space-y-1">
        <div class="flex items-center gap-2">
          <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
          <span class="text-xs font-black text-emerald-700 dark:text-emerald-400">All Safe · SOS Standby</span>
        </div>
        <p class="text-[11px] text-muted-foreground">
          0 active distress dispatches. Institutional security perimeter nominal.
        </p>
      </div>

      <div class="pt-2 text-xs space-y-2 border-t">
        <div class="flex justify-between items-center text-[11px]">
          <span class="text-muted-foreground">UGC Anti-Ragging Vault</span>
          <span class="font-bold text-emerald-600">Zero Pending</span>
        </div>
        <div class="flex justify-between items-center text-[11px]">
          <span class="text-muted-foreground">Audit Merkle Root</span>
          <span class="font-mono text-[10px] text-foreground">0x9a8f...4e1b</span>
        </div>
      </div>
    </div>
  </div>
</div>
