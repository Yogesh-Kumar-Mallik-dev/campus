<!--
  BLOCK_WEB_HELPDESK_KANBAN_001
  Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
  Purpose:   Interactive Ticket Operations Radar, SLA status tracker, filterable list, and ticket submission desk.
-->
<script lang="ts">
  import type { HelpdeskTicket, HelpdeskCategory, HelpdeskPriority, HelpdeskTicketStatus } from '$lib/types/helpdesk';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let { onSelectTicket }: { onSelectTicket?: (ticket: HelpdeskTicket) => void } = $props();

  let categories = $state<HelpdeskCategory[]>([
    {
      id: 'cat-1',
      tenantId: 'tenant-demo',
      code: 'ACADEMIC',
      name: 'Academic Affairs & Transcripts',
      description: 'Course registration, SGPA queries, transcripts and exam timetable conflicts',
      defaultPriority: 'HIGH',
      slaResponseHours: 12,
      slaResolutionHours: 48,
      isActive: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'cat-2',
      tenantId: 'tenant-demo',
      code: 'HOSTEL_MAINTENANCE',
      name: 'Hostel & Infrastructure',
      description: 'Plumbing, electrical, furniture, and room maintenance queries',
      defaultPriority: 'MEDIUM',
      slaResponseHours: 24,
      slaResolutionHours: 72,
      isActive: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'cat-3',
      tenantId: 'tenant-demo',
      code: 'FEE_BILLING',
      name: 'Central Fee & Scholarships',
      description: 'Fee payment reconciliation, scholarship disbursement and fine waivers',
      defaultPriority: 'HIGH',
      slaResponseHours: 12,
      slaResolutionHours: 36,
      isActive: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'cat-4',
      tenantId: 'tenant-demo',
      code: 'IT_SUPPORT',
      name: 'IT, Campus WiFi & LMS',
      description: 'SSO login, portal access, campus Wi-Fi network and LMS synchronization',
      defaultPriority: 'MEDIUM',
      slaResponseHours: 24,
      slaResolutionHours: 48,
      isActive: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ]);

  let tickets = $state<HelpdeskTicket[]>([
    {
      id: 't-101',
      tenantId: 'tenant-demo',
      ticketNumber: 'HD-2026-00042',
      categoryId: 'cat-1',
      categoryName: 'Academic Affairs & Transcripts',
      requesterId: 'stu-1',
      requesterName: 'Aarav Sharma (2023CSE042)',
      assignedStaffId: 'staff-exam-1',
      assignedStaffName: 'Dr. Ramesh Nair (Exam Officer)',
      title: 'End-sem Exam slot overlap with Lab Exam',
      description: 'My CS302 Advanced Algorithms exam conflicts with CS392 Compiler Lab at 10:00 AM on 18 Oct.',
      priority: 'URGENT',
      status: 'IN_PROGRESS',
      slaDueAt: new Date(Date.now() + 18 * 3600000).toISOString(),
      firstResponseAt: new Date(Date.now() - 2 * 3600000).toISOString(),
      createdAt: new Date(Date.now() - 6 * 3600000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 't-102',
      tenantId: 'tenant-demo',
      ticketNumber: 'HD-2026-00043',
      categoryId: 'cat-2',
      categoryName: 'Hostel & Infrastructure',
      requesterId: 'stu-2',
      requesterName: 'Priya Patel (2024ECE019)',
      assignedStaffId: 'staff-warden-1',
      assignedStaffName: 'Chief Warden Office',
      title: 'Water filter motor faulty in Block C Floor 2',
      description: 'Drinking water dispenser unit is not chilling and making loud vibration sounds.',
      priority: 'MEDIUM',
      status: 'WAITING_FOR_APPLICANT',
      slaDueAt: new Date(Date.now() + 44 * 3600000).toISOString(),
      firstResponseAt: new Date(Date.now() - 4 * 3600000).toISOString(),
      createdAt: new Date(Date.now() - 8 * 3600000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 't-103',
      tenantId: 'tenant-demo',
      ticketNumber: 'HD-2026-00044',
      categoryId: 'cat-3',
      categoryName: 'Central Fee & Scholarships',
      requesterId: 'stu-3',
      requesterName: 'Vikram Singhania (2022MECH088)',
      title: 'State Merit Scholarship deduction missing in fee challan',
      description: 'My state merit scholarship of ₹25,000 has not been adjusted in the Autumn 2026 invoice.',
      priority: 'HIGH',
      status: 'OPEN',
      slaDueAt: new Date(Date.now() + 28 * 3600000).toISOString(),
      createdAt: new Date(Date.now() - 2 * 3600000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 't-104',
      tenantId: 'tenant-demo',
      ticketNumber: 'HD-2026-00041',
      categoryId: 'cat-4',
      categoryName: 'IT, Campus WiFi & LMS',
      requesterId: 'stu-4',
      requesterName: 'Ananya Roy (2023BIO012)',
      assignedStaffId: 'staff-it-1',
      assignedStaffName: 'IT Helpdesk Team',
      title: 'Portal MFA Authenticator device change',
      description: 'I changed my phone and need to reset my TOTP authenticator key.',
      priority: 'MEDIUM',
      status: 'RESOLVED',
      slaDueAt: new Date(Date.now() - 10 * 3600000).toISOString(),
      firstResponseAt: new Date(Date.now() - 18 * 3600000).toISOString(),
      resolvedAt: new Date(Date.now() - 3 * 3600000).toISOString(),
      rating: 5,
      feedback: 'Resolved within 2 hours. Super quick response!',
      createdAt: new Date(Date.now() - 24 * 3600000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ]);

  let statusFilter = $state<string>('ALL');
  let priorityFilter = $state<string>('ALL');
  let categoryFilter = $state<string>('ALL');
  let searchQuery = $state<string>('');

  // New ticket form state
  let isCreateModalOpen = $state<boolean>(false);
  let newCategoryId = $state<string>('cat-1');
  let newTitle = $state<string>('');
  let newDescription = $state<string>('');
  let newPriority = $state<HelpdeskPriority>('MEDIUM');

  let filteredTickets = $derived(
    tickets.filter((t) => {
      if (statusFilter !== 'ALL' && t.status !== statusFilter) return false;
      if (priorityFilter !== 'ALL' && t.priority !== priorityFilter) return false;
      if (categoryFilter !== 'ALL' && t.categoryId !== categoryFilter) return false;
      if (searchQuery.trim() !== '') {
        const q = searchQuery.toLowerCase();
        return (
          t.title.toLowerCase().includes(q) ||
          t.ticketNumber.toLowerCase().includes(q) ||
          t.description.toLowerCase().includes(q) ||
          (t.requesterName && t.requesterName.toLowerCase().includes(q))
        );
      }
      return true;
    })
  );

  function handleCreateTicket() {
    if (!newTitle.trim() || !newDescription.trim()) return;

    const cat = categories.find((c) => c.id === newCategoryId) || categories[0];
    const newSeq = tickets.length + 45;
    const now = new Date();
    const slaDue = new Date(now.getTime() + cat.slaResolutionHours * 3600000);

    const newTicket: HelpdeskTicket = {
      id: `t-${Date.now()}`,
      tenantId: 'tenant-demo',
      ticketNumber: `HD-2026-${String(newSeq).padStart(5, '0')}`,
      categoryId: cat.id,
      categoryName: cat.name,
      requesterId: 'current-user',
      requesterName: 'Yogesh K (Logged In)',
      title: newTitle.trim(),
      description: newDescription.trim(),
      priority: newPriority,
      status: 'OPEN',
      slaDueAt: slaDue.toISOString(),
      createdAt: now.toISOString(),
      updatedAt: now.toISOString(),
    };

    tickets = [newTicket, ...tickets];
    newTitle = '';
    newDescription = '';
    isCreateModalOpen = false;
  }

  function getPriorityColor(p: HelpdeskPriority): string {
    switch (p) {
      case 'URGENT':
        return 'bg-rose-500 text-white border-rose-600';
      case 'HIGH':
        return 'bg-amber-500 text-white border-amber-600';
      case 'MEDIUM':
        return 'bg-blue-500 text-white border-blue-600';
      case 'LOW':
        return 'bg-slate-500 text-white border-slate-600';
    }
  }

  function getStatusBadge(s: HelpdeskTicketStatus): { variant: 'default' | 'secondary' | 'outline' | 'destructive'; label: string } {
    switch (s) {
      case 'OPEN':
        return { variant: 'secondary', label: 'Open / Unassigned' };
      case 'IN_PROGRESS':
        return { variant: 'default', label: 'In Progress' };
      case 'WAITING_FOR_APPLICANT':
        return { variant: 'outline', label: 'Waiting for Student' };
      case 'RESOLVED':
        return { variant: 'secondary', label: 'Resolved' };
      case 'CLOSED':
        return { variant: 'outline', label: 'Closed' };
    }
  }

  function formatSLARemaining(dueStr: string): { text: string; isBreached: boolean } {
    const due = new Date(dueStr).getTime();
    const diff = due - Date.now();
    if (diff <= 0) {
      return { text: 'SLA Breached', isBreached: true };
    }
    const hours = Math.floor(diff / 3600000);
    const mins = Math.floor((diff % 3600000) / 60000);
    return { text: `${hours}h ${mins}m left`, isBreached: false };
  }
</script>

<div class="space-y-6">
  <!-- Top Stat Cards -->
  <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Active Tickets</div>
      <div class="text-2xl font-bold mt-1 text-foreground">
        {tickets.filter((t) => t.status !== 'CLOSED' && t.status !== 'RESOLVED').length}
      </div>
      <div class="text-xs text-blue-600 font-medium mt-1">SLA radar surveillance active</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Urgent Escalations</div>
      <div class="text-2xl font-bold mt-1 text-rose-600">
        {tickets.filter((t) => t.priority === 'URGENT' && t.status !== 'CLOSED').length}
      </div>
      <div class="text-xs text-rose-500 font-medium mt-1">Requires senior officer attention</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Waiting for Student</div>
      <div class="text-2xl font-bold mt-1 text-amber-600">
        {tickets.filter((t) => t.status === 'WAITING_FOR_APPLICANT').length}
      </div>
      <div class="text-xs text-amber-500 font-medium mt-1">Applicant response pending</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Resolved (Satisfaction)</div>
      <div class="text-2xl font-bold mt-1 text-emerald-600">
        {tickets.filter((t) => t.status === 'RESOLVED' || t.status === 'CLOSED').length}
      </div>
      <div class="text-xs text-emerald-600 font-medium mt-1">★ 4.9/5 Average CSAT Score</div>
    </div>
  </div>

  <!-- Action Bar & Filter Controls -->
  <div class="p-4 bg-card rounded-xl border shadow-sm space-y-4">
    <div class="flex flex-col sm:flex-row items-center justify-between gap-3">
      <div class="flex-1 w-full max-w-md">
        <Input
          placeholder="Search by ticket #, title, student name..."
          bind:value={searchQuery}
          class="w-full text-sm"
        />
      </div>
      <div class="flex items-center gap-2 w-full sm:w-auto justify-end">
        <Button onclick={() => (isCreateModalOpen = true)} class="font-semibold shadow-sm">
          + Raise New Ticket
        </Button>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-3 pt-2 border-t">
      <div class="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
        <span>Status:</span>
        <select
          bind:value={statusFilter}
          class="bg-background text-foreground text-xs rounded-md border px-2 py-1 focus:ring-1 focus:ring-primary"
        >
          <option value="ALL">All Statuses</option>
          <option value="OPEN">Open</option>
          <option value="IN_PROGRESS">In Progress</option>
          <option value="WAITING_FOR_APPLICANT">Waiting for Applicant</option>
          <option value="RESOLVED">Resolved</option>
          <option value="CLOSED">Closed</option>
        </select>
      </div>

      <div class="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
        <span>Priority:</span>
        <select
          bind:value={priorityFilter}
          class="bg-background text-foreground text-xs rounded-md border px-2 py-1 focus:ring-1 focus:ring-primary"
        >
          <option value="ALL">All Priorities</option>
          <option value="URGENT">Urgent</option>
          <option value="HIGH">High</option>
          <option value="MEDIUM">Medium</option>
          <option value="LOW">Low</option>
        </select>
      </div>

      <div class="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
        <span>Category:</span>
        <select
          bind:value={categoryFilter}
          class="bg-background text-foreground text-xs rounded-md border px-2 py-1 focus:ring-1 focus:ring-primary"
        >
          <option value="ALL">All Categories</option>
          {#each categories as cat}
            <option value={cat.id}>{cat.name}</option>
          {/each}
        </select>
      </div>
    </div>
  </div>

  <!-- Ticket Grid / List -->
  <div class="space-y-3">
    {#if filteredTickets.length === 0}
      <div class="p-8 text-center bg-card rounded-xl border border-dashed text-muted-foreground">
        No tickets match the selected query criteria.
      </div>
    {:else}
      {#each filteredTickets as ticket}
        {@const badge = getStatusBadge(ticket.status)}
        {@const sla = formatSLARemaining(ticket.slaDueAt)}
        <div
          role="button"
          tabindex="0"
          onclick={() => onSelectTicket?.(ticket)}
          onkeydown={(e) => { if (e.key === 'Enter') onSelectTicket?.(ticket); }}
          class="p-4 bg-card hover:bg-muted/40 transition-colors rounded-xl border shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-4 cursor-pointer"
        >
          <div class="space-y-1.5 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <span class="font-mono text-xs font-bold text-primary bg-primary/10 px-2 py-0.5 rounded">
                {ticket.ticketNumber}
              </span>
              <span class="text-xs font-semibold px-2 py-0.5 rounded-full {getPriorityColor(ticket.priority)}">
                {ticket.priority}
              </span>
              <Badge variant={badge.variant}>{badge.label}</Badge>
              <span class="text-xs text-muted-foreground">· {ticket.categoryName}</span>
            </div>

            <h3 class="text-base font-bold text-foreground">
              {ticket.title}
            </h3>

            <p class="text-xs text-muted-foreground line-clamp-1">
              {ticket.description}
            </p>

            <div class="text-xs text-muted-foreground flex flex-wrap items-center gap-3 pt-1">
              <span>👤 {ticket.requesterName || 'Student'}</span>
              <span>🛠️ Assigned: {ticket.assignedStaffName || 'Unassigned'}</span>
              <span>📅 Raised: {new Date(ticket.createdAt).toLocaleDateString()}</span>
            </div>
          </div>

          <div class="flex flex-col md:items-end justify-center gap-2 border-t md:border-t-0 pt-2 md:pt-0">
            <div class="text-xs font-semibold px-2.5 py-1 rounded-md {sla.isBreached ? 'bg-rose-100 text-rose-700 border border-rose-200' : 'bg-blue-50 text-blue-700 border border-blue-200'}">
              ⏱️ {sla.text}
            </div>

            {#if ticket.rating}
              <div class="text-xs text-amber-500 font-bold">
                {'★'.repeat(ticket.rating)}{'☆'.repeat(5 - ticket.rating)} ({ticket.rating}/5)
              </div>
            {/if}

            <Button variant="outline" size="sm" class="text-xs font-semibold">
              Open Live Desk →
            </Button>
          </div>
        </div>
      {/each}
    {/if}
  </div>

  <!-- New Ticket Modal -->
  {#if isCreateModalOpen}
    <div class="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-card w-full max-w-lg rounded-2xl border shadow-xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b pb-3">
          <h2 class="text-lg font-bold text-foreground">Raise Application / Query Ticket</h2>
          <button onclick={() => (isCreateModalOpen = false)} class="text-muted-foreground hover:text-foreground">
            ✕
          </button>
        </div>

        <div class="space-y-3">
          <div>
            <label for="cat-select" class="text-xs font-semibold text-muted-foreground block mb-1">Category & Subsystem Domain</label>
            <select
              id="cat-select"
              bind:value={newCategoryId}
              class="w-full bg-background text-foreground text-sm rounded-lg border px-3 py-2"
            >
              {#each categories as cat}
                <option value={cat.id}>{cat.name} (SLA: {cat.slaResolutionHours}h)</option>
              {/each}
            </select>
          </div>

          <div>
            <label for="prio-select" class="text-xs font-semibold text-muted-foreground block mb-1">Priority Level</label>
            <select
              id="prio-select"
              bind:value={newPriority}
              class="w-full bg-background text-foreground text-sm rounded-lg border px-3 py-2"
            >
              <option value="LOW">Low (General Query)</option>
              <option value="MEDIUM">Medium (Standard Request)</option>
              <option value="HIGH">High (Impacts Academics/Billing)</option>
              <option value="URGENT">Urgent (Immediate Blocker)</option>
            </select>
          </div>

          <div>
            <label for="ticket-title" class="text-xs font-semibold text-muted-foreground block mb-1">Subject / Issue Summary</label>
            <Input id="ticket-title" bind:value={newTitle} placeholder="e.g. Hosteller WiFi connectivity issue in Room 304" />
          </div>

          <div>
            <label for="ticket-desc" class="text-xs font-semibold text-muted-foreground block mb-1">Detailed Description & Context</label>
            <textarea
              id="ticket-desc"
              bind:value={newDescription}
              rows={4}
              placeholder="Provide exact details, relevant IDs, and steps to reproduce..."
              class="w-full bg-background text-foreground text-sm rounded-lg border p-3"
            ></textarea>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-3 border-t">
          <Button variant="outline" onclick={() => (isCreateModalOpen = false)}>Cancel</Button>
          <Button onclick={handleCreateTicket} class="font-bold">Submit Query Ticket</Button>
        </div>
      </div>
    </div>
  {/if}
</div>
