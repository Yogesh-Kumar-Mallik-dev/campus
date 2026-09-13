<script lang="ts">
  /**
   * BLOCK_WEB_AUDIT_LOG_TABLE_001
   * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
   * Purpose:   Interactive, searchable, tamper-evident audit ledger explorer with hash chain validation.
   * Viewports: Responsive across mobile cards (280px+) to wide desktop data table (4K).
   */

  export interface AuditLogItem {
    id: string;
    tenant_id: string;
    actor_id?: string;
    actor_type: string;
    actor_role?: string;
    action: string;
    resource_type: string;
    resource_id?: string;
    status: 'SUCCESS' | 'FAILURE' | 'ATTEMPTED';
    status_code?: number;
    ip_address?: string;
    trace_id?: string;
    metadata?: any;
    changes?: any;
    prev_hash?: string;
    hash: string;
    created_at: string;
  }

  let {
    logs = [],
    total = 0,
    isLoading = false,
    onVerifyChain = () => {},
    onFilterChange = (_filter: any) => {}
  }: {
    logs: AuditLogItem[];
    total: number;
    isLoading?: boolean;
    onVerifyChain?: () => void;
    onFilterChange?: (filter: any) => void;
  } = $props();

  let searchQuery = $state('');
  let statusFilter = $state('');
  let selectedLog = $state<AuditLogItem | null>(null);

  function handleSearchInput(e: Event) {
    const val = (e.target as HTMLInputElement).value;
    searchQuery = val;
    onFilterChange({ action: val, status: statusFilter });
  }

  function handleStatusChange(e: Event) {
    const val = (e.target as HTMLSelectElement).value;
    statusFilter = val;
    onFilterChange({ action: searchQuery, status: val });
  }
</script>

<div class="w-full space-y-4">
  <!-- Controls Bar -->
  <div class="flex flex-col sm:flex-row gap-3 items-stretch sm:items-center justify-between">
    <div class="flex flex-1 gap-2 items-center">
      <input
        type="text"
        placeholder="Filter by action or keyword (e.g. auth:login)..."
        value={searchQuery}
        oninput={handleSearchInput}
        class="w-full sm:max-w-xs h-10 px-3 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-slate-900"
      />
      <select
        value={statusFilter}
        onchange={handleStatusChange}
        class="h-10 px-3 rounded-lg border border-slate-300 text-sm bg-white focus:outline-none focus:ring-2 focus:ring-slate-900"
      >
        <option value="">All Statuses</option>
        <option value="SUCCESS">SUCCESS</option>
        <option value="FAILURE">FAILURE</option>
        <option value="ATTEMPTED">ATTEMPTED</option>
      </select>
    </div>

    <div class="flex items-center gap-2">
      <button
        onclick={onVerifyChain}
        class="h-10 px-4 rounded-lg bg-slate-900 hover:bg-slate-800 text-white text-xs font-semibold tracking-wide uppercase transition flex items-center gap-1.5 shadow-sm"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
        </svg>
        Verify Hash Chain
      </button>
    </div>
  </div>

  <!-- Table (Desktop) / Cards (Mobile) -->
  <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
    {#if isLoading}
      <div class="p-8 text-center text-slate-500 text-sm">
        <span class="animate-spin inline-block mr-2">⟳</span> Loading audit ledger entries...
      </div>
    {:else if logs.length === 0}
      <div class="p-8 text-center text-slate-500 text-sm">
        No audit entries match the current filter criteria.
      </div>
    {:else}
      <div class="hidden md:block overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600 uppercase tracking-wider">
            <tr>
              <th class="px-4 py-3">Timestamp (UTC)</th>
              <th class="px-4 py-3">Action</th>
              <th class="px-4 py-3">Actor</th>
              <th class="px-4 py-3">Resource</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3 font-mono">Hash (SHA-256)</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            {#each logs as log (log.id)}
              <tr class="hover:bg-slate-50/70 transition">
                <td class="px-4 py-3 text-slate-600 whitespace-nowrap text-xs">
                  {new Date(log.created_at).toLocaleString()}
                </td>
                <td class="px-4 py-3 font-mono font-medium text-slate-900 text-xs">
                  {log.action}
                </td>
                <td class="px-4 py-3 text-slate-700 text-xs">
                  <span class="font-medium">{log.actor_id || 'System'}</span>
                  {#if log.actor_role}
                    <span class="ml-1 text-[10px] px-1.5 py-0.5 rounded bg-slate-100 text-slate-600 uppercase font-semibold">
                      {log.actor_role}
                    </span>
                  {/if}
                </td>
                <td class="px-4 py-3 text-slate-600 text-xs">
                  {log.resource_type}{log.resource_id ? `:${log.resource_id}` : ''}
                </td>
                <td class="px-4 py-3 text-xs">
                  {#if log.status === 'SUCCESS'}
                    <span class="px-2 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-100 text-emerald-800">
                      SUCCESS
                    </span>
                  {:else if log.status === 'FAILURE'}
                    <span class="px-2 py-0.5 rounded-full text-[11px] font-semibold bg-red-100 text-red-800">
                      FAILURE
                    </span>
                  {:else}
                    <span class="px-2 py-0.5 rounded-full text-[11px] font-semibold bg-amber-100 text-amber-800">
                      {log.status}
                    </span>
                  {/if}
                </td>
                <td class="px-4 py-3 font-mono text-[11px] text-slate-500 truncate max-w-[120px]" title={log.hash}>
                  {log.hash.substring(0, 10)}...
                </td>
                <td class="px-4 py-3 text-right">
                  <button
                    onclick={() => selectedLog = log}
                    class="text-xs font-semibold text-blue-600 hover:text-blue-800 transition"
                  >
                    Inspect
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <!-- Mobile Card View -->
      <div class="md:hidden divide-y divide-slate-100">
        {#each logs as log (log.id)}
          <div class="p-4 space-y-2">
            <div class="flex items-center justify-between">
              <span class="font-mono text-xs font-bold text-slate-900">{log.action}</span>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold {log.status === 'SUCCESS' ? 'bg-emerald-100 text-emerald-800' : 'bg-red-100 text-red-800'}">
                {log.status}
              </span>
            </div>
            <div class="text-xs text-slate-500 flex justify-between">
              <span>Actor: {log.actor_id || 'System'}</span>
              <span>{new Date(log.created_at).toLocaleTimeString()}</span>
            </div>
            <div class="flex justify-between items-center pt-1">
              <span class="text-[10px] font-mono text-slate-400">Hash: {log.hash.substring(0, 8)}...</span>
              <button
                onclick={() => selectedLog = log}
                class="text-xs font-medium text-blue-600"
              >
                View Details →
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}

    <div class="px-4 py-3 bg-slate-50 border-t border-slate-200 text-xs text-slate-500 flex items-center justify-between">
      <span>Showing {logs.length} of {total} records</span>
      <span class="text-[11px] font-mono">Immutable SHA-256 Chained</span>
    </div>
  </div>

  <!-- Detail Modal -->
  {#if selectedLog}
    <div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-xl max-w-2xl w-full max-h-[85vh] overflow-y-auto p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b pb-3">
          <h3 class="text-base font-bold text-slate-900">Audit Record Inspector</h3>
          <button onclick={() => selectedLog = null} class="text-slate-400 hover:text-slate-600 text-lg font-bold">✕</button>
        </div>

        <div class="grid grid-cols-2 gap-4 text-xs">
          <div><strong class="text-slate-500">Record ID:</strong> <span class="font-mono">{selectedLog.id}</span></div>
          <div><strong class="text-slate-500">Tenant ID:</strong> <span class="font-mono">{selectedLog.tenant_id}</span></div>
          <div><strong class="text-slate-500">Action:</strong> <span class="font-mono font-semibold">{selectedLog.action}</span></div>
          <div><strong class="text-slate-500">Status:</strong> {selectedLog.status}</div>
          <div><strong class="text-slate-500">Actor:</strong> {selectedLog.actor_id || 'System'} ({selectedLog.actor_type})</div>
          <div><strong class="text-slate-500">Trace ID:</strong> <span class="font-mono">{selectedLog.trace_id || 'N/A'}</span></div>
          <div><strong class="text-slate-500">IP Address:</strong> {selectedLog.ip_address || 'N/A'}</div>
          <div><strong class="text-slate-500">Created:</strong> {new Date(selectedLog.created_at).toISOString()}</div>
        </div>

        <div class="space-y-1">
          <span class="text-xs font-semibold text-slate-700">Cryptographic Hash Signature:</span>
          <p class="p-2 bg-slate-100 rounded text-[11px] font-mono break-all text-slate-800">{selectedLog.hash}</p>
        </div>

        <div class="space-y-1">
          <span class="text-xs font-semibold text-slate-700">Previous Record Hash (prev_hash):</span>
          <p class="p-2 bg-slate-100 rounded text-[11px] font-mono break-all text-slate-600">{selectedLog.prev_hash || 'Genesis'}</p>
        </div>

        {#if selectedLog.metadata}
          <div class="space-y-1">
            <span class="text-xs font-semibold text-slate-700">Contextual Metadata:</span>
            <pre class="p-2.5 bg-slate-950 text-slate-100 rounded text-[11px] font-mono overflow-x-auto">{JSON.stringify(selectedLog.metadata, null, 2)}</pre>
          </div>
        {/if}

        <div class="pt-2 text-right">
          <button
            onclick={() => selectedLog = null}
            class="px-4 py-2 bg-slate-200 hover:bg-slate-300 rounded-lg text-xs font-semibold text-slate-800"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>
