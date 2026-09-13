<script lang="ts">
  /**
   * BLOCK_WEB_COMPLIANCE_REPORT_001
   * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
   * Purpose:   Accreditation & compliance dashboard for NAAC, NIRF, ISO 27001, and ABET audit reviews.
   */

  export interface ComplianceReport {
    tenant_id: string;
    framework: string;
    from_time: string;
    to_time: string;
    total_events: number;
    successful_events: number;
    failed_events: number;
    security_incidents: number;
    top_actions: { action: string; count: number }[];
    top_actors: { actor_id: string; actor_role: string; count: number }[];
    chain_integrity: boolean;
    generated_at: string;
  }

  let {
    report = null,
    isLoading = false,
    onGenerate = (_framework: string) => {}
  }: {
    report: ComplianceReport | null;
    isLoading?: boolean;
    onGenerate?: (framework: string) => void;
  } = $props();

  let selectedFramework = $state('NAAC');
</script>

<div class="space-y-6">
  <!-- Framework Selection Header -->
  <div class="p-6 bg-white rounded-xl border border-slate-200 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
    <div>
      <h2 class="text-xl font-bold text-slate-900">Institutional Compliance & Accreditation</h2>
      <p class="text-sm text-slate-500 mt-0.5">
        Generate tamper-evident, cryptographically verified audit summaries for regulatory standards.
      </p>
    </div>

    <div class="flex items-center gap-3">
      <select
        bind:value={selectedFramework}
        class="h-10 px-3 rounded-lg border border-slate-300 text-sm bg-white font-medium focus:outline-none focus:ring-2 focus:ring-slate-900"
      >
        <option value="NAAC">NAAC (Criterion 6 & Governance)</option>
        <option value="NIRF">NIRF (Institutional Data Verification)</option>
        <option value="ISO_27001">ISO/IEC 27001 (Information Security)</option>
        <option value="ABET">ABET (Computing & Engineering Accreditation)</option>
      </select>

      <button
        onclick={() => onGenerate(selectedFramework)}
        disabled={isLoading}
        class="h-10 px-4 bg-slate-900 hover:bg-slate-800 text-white text-xs font-semibold uppercase tracking-wider rounded-lg transition disabled:opacity-50"
      >
        {isLoading ? 'Compiling...' : 'Generate Report'}
      </button>
    </div>
  </div>

  {#if report}
    <!-- KPI Summary Metrics -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="p-5 bg-white rounded-xl border border-slate-200 shadow-sm">
        <span class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Recorded Events</span>
        <div class="mt-2 text-3xl font-bold text-slate-900">{report.total_events.toLocaleString()}</div>
        <div class="mt-1 text-xs text-slate-400">Within reporting window</div>
      </div>

      <div class="p-5 bg-white rounded-xl border border-slate-200 shadow-sm">
        <span class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Successful Mutations</span>
        <div class="mt-2 text-3xl font-bold text-emerald-600">{report.successful_events.toLocaleString()}</div>
        <div class="mt-1 text-xs text-slate-400">
          {((report.successful_events / (report.total_events || 1)) * 100).toFixed(1)}% nominal rate
        </div>
      </div>

      <div class="p-5 bg-white rounded-xl border border-slate-200 shadow-sm">
        <span class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Security Incidents</span>
        <div class="mt-2 text-3xl font-bold {report.security_incidents > 0 ? 'text-amber-600' : 'text-slate-900'}">
          {report.security_incidents.toLocaleString()}
        </div>
        <div class="mt-1 text-xs text-slate-400">Lockouts, breaches, failures</div>
      </div>

      <div class="p-5 bg-white rounded-xl border border-slate-200 shadow-sm">
        <span class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Ledger Cryptographic Integrity</span>
        <div class="mt-2 flex items-center gap-2">
          {#if report.chain_integrity}
            <span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-100 text-emerald-800">
              ✓ INTACT & VERIFIED
            </span>
          {:else}
            <span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-bold bg-red-100 text-red-800">
              ⚠ TAMPERING DETECTED
            </span>
          {/if}
        </div>
        <div class="mt-1 text-xs text-slate-400">SHA-256 Merkle chain check</div>
      </div>
    </div>

    <!-- Top Activities Breakdown -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="p-6 bg-white rounded-xl border border-slate-200 shadow-sm space-y-4">
        <h3 class="text-sm font-bold uppercase tracking-wider text-slate-700">Top Audit Actions by Volume</h3>
        <div class="divide-y divide-slate-100">
          {#each report.top_actions as item}
            <div class="py-2.5 flex items-center justify-between text-xs">
              <span class="font-mono font-medium text-slate-800">{item.action}</span>
              <span class="px-2 py-0.5 rounded bg-slate-100 text-slate-700 font-semibold">{item.count.toLocaleString()}</span>
            </div>
          {/each}
        </div>
      </div>

      <div class="p-6 bg-white rounded-xl border border-slate-200 shadow-sm space-y-4">
        <h3 class="text-sm font-bold uppercase tracking-wider text-slate-700">Most Active Actors</h3>
        <div class="divide-y divide-slate-100">
          {#each report.top_actors as actor}
            <div class="py-2.5 flex items-center justify-between text-xs">
              <div>
                <span class="font-medium text-slate-900">{actor.actor_id}</span>
                <span class="ml-1.5 text-[10px] px-1.5 py-0.5 rounded bg-slate-100 text-slate-600 font-semibold uppercase">{actor.actor_role}</span>
              </div>
              <span class="px-2 py-0.5 rounded bg-slate-100 text-slate-700 font-semibold">{actor.count.toLocaleString()} ops</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}
</div>
