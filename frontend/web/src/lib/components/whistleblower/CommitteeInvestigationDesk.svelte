<!--
  BLOCK_WEB_WHISTLEBLOWER_COMMITTEE_001
  Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
  Purpose:   Vigilance & Anti-Ragging Committee Investigation Console, status updates, and findings logging.
-->
<script lang="ts">
  import type { WhistleblowerReport, WhistleblowerCategory, WhistleblowerStatus } from '$lib/types/whistleblower';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let reports = $state<WhistleblowerReport[]>([
    {
      id: 'wb-001',
      tenantId: 'tenant-demo',
      reportNumber: 'WB-2026-00042',
      category: 'ANTI_RAGGING',
      severity: 'CRITICAL',
      title: 'Late night intimidation in Hostel Block C 4th floor',
      description: 'First year students being subjected to late night roll-call and pushups past 01:00 AM.',
      status: 'UNDER_INVESTIGATION',
      assignedOfficerName: 'Anti-Ragging Squad (Unit 2)',
      submittedAt: new Date(Date.now() - 36 * 3600000).toISOString(),
      createdAt: new Date(Date.now() - 36 * 3600000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'wb-002',
      tenantId: 'tenant-demo',
      reportNumber: 'WB-2026-00043',
      category: 'FINANCIAL_FRAUD',
      severity: 'HIGH',
      title: 'Discrepancies in Central Library purchase receipts',
      description: 'Books billed at list prices while bulk publisher discount was omitted in institutional invoice.',
      status: 'SUBMITTED',
      submittedAt: new Date(Date.now() - 12 * 3600000).toISOString(),
      createdAt: new Date(Date.now() - 12 * 3600000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'wb-003',
      tenantId: 'tenant-demo',
      reportNumber: 'WB-2026-00040',
      category: 'ACADEMIC_CORRUPTION',
      severity: 'HIGH',
      title: 'Fake certificate submission for credit waiver',
      description: 'Online MOOC course completion certificate verified as counterfeit.',
      status: 'RESOLVED',
      assignedOfficerName: 'Academic Ethics Panel',
      findingsSummary: 'Certificate verified against Coursera API. Counterfeit confirmed; student credits cancelled and disciplinary notice issued.',
      submittedAt: new Date(Date.now() - 72 * 3600000).toISOString(),
      resolvedAt: new Date(Date.now() - 24 * 3600000).toISOString(),
      createdAt: new Date(Date.now() - 72 * 3600000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ]);

  let selectedReport = $state<WhistleblowerReport | null>(null);
  let categoryFilter = $state<string>('ALL');
  let statusFilter = $state<string>('ALL');

  // Resolve modal
  let isResolveModalOpen = $state<boolean>(false);
  let findingsText = $state<string>('');
  let resolveStatus = $state<WhistleblowerStatus>('RESOLVED');

  let filteredReports = $derived(
    reports.filter((r) => {
      if (categoryFilter !== 'ALL' && r.category !== categoryFilter) return false;
      if (statusFilter !== 'ALL' && r.status !== statusFilter) return false;
      return true;
    })
  );

  function handleAssign(reportId: string) {
    reports = reports.map((r) =>
      r.id === reportId ? { ...r, status: 'UNDER_INVESTIGATION', assignedOfficerName: 'Vigilance Officer (Assigned)', updatedAt: new Date().toISOString() } : r
    );
  }

  function handleCompleteResolution() {
    if (!findingsText.trim() || !selectedReport) return;

    reports = reports.map((r) =>
      r.id === selectedReport!.id
        ? {
            ...r,
            status: resolveStatus,
            findingsSummary: findingsText.trim(),
            resolvedAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          }
        : r
    );

    isResolveModalOpen = false;
    findingsText = '';
    selectedReport = null;
  }
</script>

<div class="space-y-6">
  <!-- Top Stat Cards -->
  <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Active Inquiries</div>
      <div class="text-2xl font-black mt-1 text-foreground">
        {reports.filter((r) => r.status !== 'RESOLVED' && r.status !== 'REJECTED').length}
      </div>
      <div class="text-xs text-indigo-600 font-medium mt-1">Confidential vigilance active</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-bold text-rose-600 uppercase tracking-wider">Anti-Ragging Cases</div>
      <div class="text-2xl font-black mt-1 text-rose-600">
        {reports.filter((r) => r.category === 'ANTI_RAGGING' && r.status !== 'RESOLVED').length}
      </div>
      <div class="text-xs text-rose-500 font-medium mt-1">Zero Tolerance Mandate</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-bold text-amber-600 uppercase tracking-wider">Pending Assignment</div>
      <div class="text-2xl font-black mt-1 text-amber-600">
        {reports.filter((r) => r.status === 'SUBMITTED').length}
      </div>
      <div class="text-xs text-amber-500 font-medium mt-1">Requires investigator allocation</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-bold text-emerald-600 uppercase tracking-wider">Resolved & Closed</div>
      <div class="text-2xl font-black mt-1 text-emerald-600">
        {reports.filter((r) => r.status === 'RESOLVED').length}
      </div>
      <div class="text-xs text-emerald-600 font-medium mt-1">Action taken verified</div>
    </div>
  </div>

  <!-- Filters -->
  <div class="p-4 bg-card rounded-xl border shadow-sm flex flex-wrap items-center gap-4">
    <div class="flex items-center gap-2 text-xs font-semibold">
      <span>Category:</span>
      <select bind:value={categoryFilter} class="bg-background text-foreground text-xs rounded-md border px-2 py-1">
        <option value="ALL">All Categories</option>
        <option value="ANTI_RAGGING">Anti-Ragging</option>
        <option value="FINANCIAL_FRAUD">Financial Fraud</option>
        <option value="ACADEMIC_CORRUPTION">Academic Corruption</option>
        <option value="HARASSMENT">Harassment</option>
        <option value="SAFETY_VIOLATION">Safety Violation</option>
      </select>
    </div>

    <div class="flex items-center gap-2 text-xs font-semibold">
      <span>Status:</span>
      <select bind:value={statusFilter} class="bg-background text-foreground text-xs rounded-md border px-2 py-1">
        <option value="ALL">All Statuses</option>
        <option value="SUBMITTED">Submitted</option>
        <option value="UNDER_INVESTIGATION">Under Investigation</option>
        <option value="EVIDENCE_REQUESTED">Evidence Requested</option>
        <option value="RESOLVED">Resolved</option>
        <option value="REJECTED">Rejected</option>
      </select>
    </div>
  </div>

  <!-- Reports Roster -->
  <div class="space-y-3">
    {#each filteredReports as r}
      <div class="p-4 bg-card rounded-xl border shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div class="space-y-1.5 flex-1">
          <div class="flex flex-wrap items-center gap-2">
            <span class="font-mono text-xs font-bold text-primary bg-primary/10 px-2 py-0.5 rounded">
              {r.reportNumber}
            </span>
            <span class="text-xs font-bold px-2 py-0.5 rounded-full {r.severity === 'CRITICAL' ? 'bg-rose-500 text-white' : 'bg-amber-500 text-white'}">
              {r.severity}
            </span>
            <Badge variant="outline">{r.category}</Badge>
            <Badge variant={r.status === 'RESOLVED' ? 'secondary' : 'default'}>{r.status}</Badge>
          </div>

          <h3 class="text-base font-bold text-foreground">{r.title}</h3>
          <p class="text-xs text-muted-foreground line-clamp-2">{r.description}</p>

          <div class="text-xs text-muted-foreground flex gap-4 pt-1">
            <span>Investigator: <strong>{r.assignedOfficerName || 'Unassigned'}</strong></span>
            <span>Lodged: {new Date(r.submittedAt).toLocaleDateString()}</span>
          </div>

          {#if r.findingsSummary}
            <div class="p-2.5 bg-emerald-50 text-emerald-900 border border-emerald-200 rounded-lg text-xs">
              <span class="font-bold">✓ Findings Log:</span> {r.findingsSummary}
            </div>
          {/if}
        </div>

        <div class="flex flex-col md:items-end justify-center gap-2">
          {#if r.status === 'SUBMITTED'}
            <Button size="sm" onclick={() => handleAssign(r.id)} class="font-bold">
              Assign to Squad
            </Button>
          {/if}
          {#if r.status !== 'RESOLVED' && r.status !== 'REJECTED'}
            <Button
              size="sm"
              variant="outline"
              onclick={() => {
                selectedReport = r;
                isResolveModalOpen = true;
              }}
              class="text-emerald-600 border-emerald-300 font-bold"
            >
              ✓ Complete Findings
            </Button>
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <!-- Resolution Modal -->
  {#if isResolveModalOpen && selectedReport}
    <div class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-card w-full max-w-md rounded-2xl border shadow-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b pb-3">
          <h2 class="text-lg font-bold text-foreground">Record Investigation Findings</h2>
          <button onclick={() => (isResolveModalOpen = false)} class="text-muted-foreground hover:text-foreground">✕</button>
        </div>

        <div class="space-y-3">
          <div>
            <label for="res-verdict" class="text-xs font-semibold text-muted-foreground block mb-1">Investigation Outcome</label>
            <select
              id="res-verdict"
              bind:value={resolveStatus}
              class="w-full bg-background text-foreground text-sm rounded-lg border px-3 py-2"
            >
              <option value="RESOLVED">Resolved - Corrective Action / Disciplinary Measures Taken</option>
              <option value="REJECTED">Dismissed - Insufficient Evidence / Unsubstantiated</option>
            </select>
          </div>

          <div>
            <label for="find-summary" class="text-xs font-semibold text-muted-foreground block mb-1">Official Findings Summary</label>
            <textarea
              id="find-summary"
              bind:value={findingsText}
              rows={4}
              placeholder="Detail actions taken, disciplinary measures, and committee conclusion..."
              class="w-full bg-background text-foreground text-xs rounded-xl border p-3"
            ></textarea>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-3 border-t">
          <Button variant="outline" onclick={() => (isResolveModalOpen = false)}>Cancel</Button>
          <Button onclick={handleCompleteResolution} class="bg-emerald-600 hover:bg-emerald-700 text-white font-bold">
            Finalize & Close Report
          </Button>
        </div>
      </div>
    </div>
  {/if}
</div>
