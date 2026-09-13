<!--
  BLOCK_WEB_PORTAL_DESK_001
  Subsystem: Rank 16 - Public Web Portal (portal)
  Purpose:   Admissions Counselor Triage Desk for Reviewing, Contacting, and Converting Prospective Student Leads.
-->
<script lang="ts">
  import type { PortalPublicInquiry, PortalInquiryStatus } from '$lib/types/portal';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let inquiries = $state<PortalPublicInquiry[]>([
    {
      id: 'inq-101',
      tenantId: 'tenant-demo',
      prospectName: 'Aarav Mehta',
      prospectEmail: 'aarav.mehta@gmail.com',
      prospectPhone: '+91 98234 11223',
      programOfInterest: 'B.Tech in Computer Science & Artificial Intelligence',
      message: 'Interested in sports scholarship quotas and lab infrastructure for robotics.',
      status: 'NEW',
      createdAt: new Date(Date.now() - 3600000 * 2).toISOString(),
      updatedAt: new Date(Date.now() - 3600000 * 2).toISOString(),
    },
    {
      id: 'inq-102',
      tenantId: 'tenant-demo',
      prospectName: 'Pooja Iyer',
      prospectEmail: 'pooja.iyer@outlook.com',
      prospectPhone: '+91 94455 66778',
      programOfInterest: 'M.Tech in Cybersecurity & Cloud Resilience',
      message: 'Checking for weekend batch options and fee installment schedules.',
      status: 'CONTACTED',
      assignedCounselorId: 'counselor-01',
      assignedCounselorName: 'Dr. Suresh Nair',
      notes: 'Initial counseling call done. Sent brochure and fee schedule PDF.',
      createdAt: new Date(Date.now() - 86400000 * 1).toISOString(),
      updatedAt: new Date(Date.now() - 3600000 * 4).toISOString(),
    },
    {
      id: 'inq-103',
      tenantId: 'tenant-demo',
      prospectName: 'Vikram Rajput',
      prospectEmail: 'vikram.r@gmail.com',
      prospectPhone: '+91 91234 98765',
      programOfInterest: 'B.Tech in Electronics & VLSI Semiconductor Design',
      message: 'Hostel accommodation inquiry and placement stats for VLSI branch.',
      status: 'CONVERTED',
      assignedCounselorId: 'counselor-02',
      assignedCounselorName: 'Meera Deshmukh',
      notes: 'Application submitted & entrance test cleared. Seat blocked in Round 1.',
      createdAt: new Date(Date.now() - 86400000 * 3).toISOString(),
      updatedAt: new Date(Date.now() - 86400000 * 1).toISOString(),
    }
  ]);

  let statusFilter = $state<string>('ALL');
  let searchQuery = $state<string>('');
  let selectedInquiry = $state<PortalPublicInquiry | null>(null);
  let counselorNote = $state<string>('');

  let filteredInquiries = $derived(
    inquiries.filter((inq) => {
      if (statusFilter !== 'ALL' && inq.status !== statusFilter) return false;
      if (searchQuery.trim() !== '') {
        const q = searchQuery.toLowerCase();
        return (
          inq.prospectName.toLowerCase().includes(q) ||
          inq.prospectEmail.toLowerCase().includes(q) ||
          inq.programOfInterest.toLowerCase().includes(q) ||
          inq.prospectPhone.includes(q)
        );
      }
      return true;
    })
  );

  function getStatusBadge(status: PortalInquiryStatus) {
    switch (status) {
      case 'NEW':
        return 'bg-blue-100 text-blue-800 border-blue-200 dark:bg-blue-900/40 dark:text-blue-300';
      case 'CONTACTED':
        return 'bg-amber-100 text-amber-800 border-amber-200 dark:bg-amber-900/40 dark:text-amber-300';
      case 'CONVERTED':
        return 'bg-emerald-100 text-emerald-800 border-emerald-200 dark:bg-emerald-900/40 dark:text-emerald-300';
      case 'CLOSED':
        return 'bg-slate-100 text-slate-800 border-slate-200 dark:bg-slate-800 dark:text-slate-300';
    }
  }

  function handleSelect(inq: PortalPublicInquiry) {
    selectedInquiry = inq;
    counselorNote = inq.notes || '';
  }

  function updateStatus(newStatus: PortalInquiryStatus) {
    if (!selectedInquiry) return;
    selectedInquiry.status = newStatus;
    selectedInquiry.updatedAt = new Date().toISOString();
    const idx = inquiries.findIndex(i => i.id === selectedInquiry?.id);
    if (idx !== -1) {
      inquiries[idx] = { ...selectedInquiry };
    }
  }

  function saveNotes() {
    if (!selectedInquiry) return;
    selectedInquiry.notes = counselorNote;
    selectedInquiry.updatedAt = new Date().toISOString();
    const idx = inquiries.findIndex(i => i.id === selectedInquiry?.id);
    if (idx !== -1) {
      inquiries[idx] = { ...selectedInquiry };
    }
  }
</script>

<div class="space-y-6">
  <!-- Top Bar -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h2 class="text-xl font-black text-foreground">Admissions Lead Triage Desk</h2>
      <p class="text-xs text-muted-foreground mt-0.5">
        Manage prospective candidate queries, track counselor follow-ups, and convert leads.
      </p>
    </div>

    <div class="flex items-center gap-2">
      <Input
        placeholder="Search prospect name, phone, program..."
        bind:value={searchQuery}
        class="w-full sm:w-64 text-xs"
      />
    </div>
  </div>

  <!-- Status Filter Pills -->
  <div class="flex flex-wrap gap-2">
    {#each ['ALL', 'NEW', 'CONTACTED', 'CONVERTED', 'CLOSED'] as st}
      <button
        onclick={() => (statusFilter = st)}
        class="px-3 py-1 rounded-xl text-xs font-bold transition-colors border {statusFilter === st ? 'bg-primary text-primary-foreground border-primary' : 'bg-card text-muted-foreground hover:bg-muted'}"
      >
        {st}
        <span class="ml-1 text-[10px] opacity-80">
          ({st === 'ALL' ? inquiries.length : inquiries.filter(i => i.status === st).length})
        </span>
      </button>
    {/each}
  </div>

  <!-- Main Split Layout -->
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Leads List -->
    <div class="lg:col-span-2 space-y-3">
      {#if filteredInquiries.length === 0}
        <div class="p-8 text-center bg-card rounded-2xl border text-muted-foreground text-xs">
          No inquiries found matching your filters.
        </div>
      {/if}

      {#each filteredInquiries as inq}
        <button
          onclick={() => handleSelect(inq)}
          class="w-full text-left p-4 bg-card rounded-2xl border shadow-sm hover:border-primary/50 transition-all flex flex-col sm:flex-row sm:items-center justify-between gap-3 {selectedInquiry?.id === inq.id ? 'ring-2 ring-primary border-transparent' : ''}"
        >
          <div class="space-y-1">
            <div class="flex items-center gap-2">
              <span class="font-bold text-sm text-foreground">{inq.prospectName}</span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded-full border {getStatusBadge(inq.status)}">
                {inq.status}
              </span>
            </div>
            <div class="text-xs text-muted-foreground">
              <span>🎯 {inq.programOfInterest}</span>
            </div>
            <div class="text-[11px] text-muted-foreground flex gap-3 pt-0.5">
              <span>📧 {inq.prospectEmail}</span>
              <span>📞 {inq.prospectPhone}</span>
            </div>
          </div>

          <div class="text-right flex flex-col sm:items-end justify-between self-start sm:self-center">
            <span class="text-[10px] text-muted-foreground font-mono">
              {new Date(inq.createdAt).toLocaleDateString()}
            </span>
            {#if inq.assignedCounselorName}
              <span class="text-[11px] text-primary font-medium mt-1">
                👤 {inq.assignedCounselorName}
              </span>
            {/if}
          </div>
        </button>
      {/each}
    </div>

    <!-- Lead Detail & Action Panel -->
    <div class="space-y-4">
      {#if selectedInquiry}
        <div class="p-6 bg-card rounded-2xl border shadow-md space-y-4 sticky top-6">
          <div class="flex items-center justify-between border-b pb-3">
            <div>
              <h3 class="text-base font-black text-foreground">{selectedInquiry.prospectName}</h3>
              <p class="text-xs text-muted-foreground font-mono">{selectedInquiry.prospectEmail}</p>
            </div>
            <span class="text-xs font-bold px-2.5 py-0.5 rounded-full border {getStatusBadge(selectedInquiry.status)}">
              {selectedInquiry.status}
            </span>
          </div>

          <div class="space-y-2 text-xs">
            <div>
              <span class="text-muted-foreground block text-[11px]">Program of Interest:</span>
              <span class="font-bold text-foreground">{selectedInquiry.programOfInterest}</span>
            </div>
            <div>
              <span class="text-muted-foreground block text-[11px]">Phone Contact:</span>
              <span class="font-mono text-foreground font-bold">{selectedInquiry.prospectPhone}</span>
            </div>
            <div>
              <span class="text-muted-foreground block text-[11px]">Inquiry / Note from Candidate:</span>
              <div class="p-3 bg-muted/40 rounded-xl mt-1 text-xs text-foreground italic border">
                "{selectedInquiry.message}"
              </div>
            </div>
          </div>

          <!-- Counselor Followup Section -->
          <div class="space-y-2 pt-2 border-t">
            <label for="counselor-notes-input" class="text-xs font-bold text-foreground block">Counselor Notes & Next Steps</label>
            <textarea
              id="counselor-notes-input"
              bind:value={counselorNote}
              rows={3}
              placeholder="Record call summary, scholarship discussion, or verification status..."
              class="w-full bg-background text-foreground text-xs rounded-xl border p-2.5"
            ></textarea>
            <Button size="sm" onclick={saveNotes} class="w-full font-bold">
              Save Counselor Notes
            </Button>
          </div>

          <!-- Triage Actions -->
          <div class="space-y-1.5 pt-2 border-t">
            <span class="text-[11px] font-bold text-muted-foreground block">Update Status:</span>
            <div class="grid grid-cols-2 gap-2">
              <Button
                variant={selectedInquiry.status === 'CONTACTED' ? 'default' : 'outline'}
                size="sm"
                onclick={() => updateStatus('CONTACTED')}
                class="text-xs font-bold"
              >
                Mark Contacted
              </Button>
              <Button
                variant={selectedInquiry.status === 'CONVERTED' ? 'default' : 'outline'}
                size="sm"
                onclick={() => updateStatus('CONVERTED')}
                class="text-xs font-bold text-emerald-600 dark:text-emerald-400"
              >
                Mark Converted
              </Button>
              <Button
                variant={selectedInquiry.status === 'CLOSED' ? 'default' : 'outline'}
                size="sm"
                onclick={() => updateStatus('CLOSED')}
                class="text-xs font-bold col-span-2"
              >
                Close Inquiry
              </Button>
            </div>
          </div>
        </div>
      {:else}
        <div class="p-8 text-center bg-card rounded-2xl border text-muted-foreground text-xs">
          Select an inquiry from the list to triage candidate details and record counselor notes.
        </div>
      {/if}
    </div>
  </div>
</div>
