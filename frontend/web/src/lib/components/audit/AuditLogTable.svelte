<script lang="ts">
  /**
   * BLOCK_WEB_AUDIT_LOG_TABLE_001
   * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
   * Purpose:   Interactive, searchable, tamper-evident audit ledger explorer with hash chain validation.
   * Viewports: Responsive across mobile cards (280px+) to wide desktop data table (4K).
   */
  import * as Table from '$lib/components/ui/table';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import * as Card from '$lib/components/ui/card';

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
      <Input
        type="text"
        placeholder="Filter by action keyword (e.g. auth:login)..."
        value={searchQuery}
        oninput={handleSearchInput}
        class="w-full sm:max-w-xs h-10 text-sm bg-card"
      />
      <select
        value={statusFilter}
        onchange={handleStatusChange}
        class="h-10 px-3 rounded-lg border border-border text-sm bg-card text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
      >
        <option value="">All Statuses</option>
        <option value="SUCCESS">SUCCESS</option>
        <option value="FAILURE">FAILURE</option>
        <option value="ATTEMPTED">ATTEMPTED</option>
      </select>
    </div>

    <div class="flex items-center gap-2">
      <Button
        onclick={onVerifyChain}
        class="h-10 gap-2 font-semibold text-xs uppercase tracking-wider"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
        </svg>
        Verify Hash Chain
      </Button>
    </div>
  </div>

  <!-- Table (Desktop) / Cards (Mobile) -->
  <Card.Root class="border-border bg-card overflow-hidden">
    {#if isLoading}
      <div class="p-8 text-center text-muted-foreground text-sm">
        <span class="animate-spin inline-block mr-2">⟳</span> Loading audit ledger entries...
      </div>
    {:else if logs.length === 0}
      <div class="p-8 text-center text-muted-foreground text-sm">
        No audit entries match the current filter criteria.
      </div>
    {:else}
      <div class="hidden md:block overflow-x-auto">
        <Table.Root>
          <Table.Header>
            <Table.Row class="hover:bg-transparent">
              <Table.Head class="text-xs uppercase font-bold text-muted-foreground">Timestamp (UTC)</Table.Head>
              <Table.Head class="text-xs uppercase font-bold text-muted-foreground">Action</Table.Head>
              <Table.Head class="text-xs uppercase font-bold text-muted-foreground">Actor</Table.Head>
              <Table.Head class="text-xs uppercase font-bold text-muted-foreground">Resource</Table.Head>
              <Table.Head class="text-xs uppercase font-bold text-muted-foreground">Status</Table.Head>
              <Table.Head class="text-xs uppercase font-bold text-muted-foreground font-mono">Hash (SHA-256)</Table.Head>
              <Table.Head class="text-right text-xs uppercase font-bold text-muted-foreground">Inspect</Table.Head>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {#each logs as log (log.id)}
              <Table.Row>
                <Table.Cell class="text-xs text-muted-foreground whitespace-nowrap">
                  {new Date(log.created_at).toLocaleString()}
                </Table.Cell>
                <Table.Cell class="font-mono font-semibold text-xs text-foreground">
                  {log.action}
                </Table.Cell>
                <Table.Cell class="text-xs text-foreground">
                  <span class="font-medium">{log.actor_id || 'System'}</span>
                  {#if log.actor_role}
                    <Badge variant="secondary" class="ml-1 text-[10px] py-0 font-bold uppercase">
                      {log.actor_role}
                    </Badge>
                  {/if}
                </Table.Cell>
                <Table.Cell class="text-xs text-muted-foreground">
                  {log.resource_type}{log.resource_id ? `:${log.resource_id}` : ''}
                </Table.Cell>
                <Table.Cell>
                  {#if log.status === 'SUCCESS'}
                    <Badge variant="default" class="text-[10px] font-bold">
                      SUCCESS
                    </Badge>
                  {:else if log.status === 'FAILURE'}
                    <Badge variant="destructive" class="text-[10px] font-bold">
                      FAILURE
                    </Badge>
                  {:else}
                    <Badge variant="outline" class="text-[10px] font-bold">
                      {log.status}
                    </Badge>
                  {/if}
                </Table.Cell>
                <Table.Cell class="font-mono text-[11px] text-muted-foreground truncate max-w-[120px]" title={log.hash}>
                  {log.hash.substring(0, 10)}...
                </Table.Cell>
                <Table.Cell class="text-right">
                  <Button
                    variant="ghost"
                    size="sm"
                    onclick={() => selectedLog = log}
                    class="text-xs font-semibold text-primary"
                  >
                    View
                  </Button>
                </Table.Cell>
              </Table.Row>
            {/each}
          </Table.Body>
        </Table.Root>
      </div>

      <!-- Mobile Card View -->
      <div class="md:hidden divide-y divide-border">
        {#each logs as log (log.id)}
          <div class="p-4 space-y-2">
            <div class="flex items-center justify-between">
              <span class="font-mono text-xs font-bold text-foreground">{log.action}</span>
              <Badge variant={log.status === 'SUCCESS' ? 'default' : 'destructive'} class="text-[10px]">
                {log.status}
              </Badge>
            </div>
            <div class="text-xs text-muted-foreground flex justify-between">
              <span>Actor: {log.actor_id || 'System'}</span>
              <span>{new Date(log.created_at).toLocaleTimeString()}</span>
            </div>
            <div class="flex justify-between items-center pt-1">
              <span class="text-[10px] font-mono text-muted-foreground">Hash: {log.hash.substring(0, 8)}...</span>
              <Button
                variant="link"
                size="sm"
                onclick={() => selectedLog = log}
                class="text-xs p-0 h-auto"
              >
                Inspect →
              </Button>
            </div>
          </div>
        {/each}
      </div>
    {/if}

    <div class="px-4 py-3 bg-muted/40 border-t border-border text-xs text-muted-foreground flex items-center justify-between">
      <span>Showing {logs.length} of {total} records</span>
      <span class="text-[11px] font-mono">Immutable SHA-256 Merkle Chain</span>
    </div>
  </Card.Root>

  <!-- Detail Modal -->
  {#if selectedLog}
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-card border border-border rounded-xl max-w-2xl w-full max-h-[85vh] overflow-y-auto p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-border pb-3">
          <h3 class="text-base font-bold text-foreground">Audit Record Inspector</h3>
          <button onclick={() => selectedLog = null} class="text-muted-foreground hover:text-foreground text-lg font-bold">✕</button>
        </div>

        <div class="grid grid-cols-2 gap-4 text-xs">
          <div><strong class="text-muted-foreground">Record ID:</strong> <span class="font-mono">{selectedLog.id}</span></div>
          <div><strong class="text-muted-foreground">Tenant ID:</strong> <span class="font-mono">{selectedLog.tenant_id}</span></div>
          <div><strong class="text-muted-foreground">Action:</strong> <span class="font-mono font-semibold">{selectedLog.action}</span></div>
          <div><strong class="text-muted-foreground">Status:</strong> {selectedLog.status}</div>
          <div><strong class="text-muted-foreground">Actor:</strong> {selectedLog.actor_id || 'System'} ({selectedLog.actor_type})</div>
          <div><strong class="text-muted-foreground">Trace ID:</strong> <span class="font-mono">{selectedLog.trace_id || 'N/A'}</span></div>
          <div><strong class="text-muted-foreground">IP Address:</strong> {selectedLog.ip_address || 'N/A'}</div>
          <div><strong class="text-muted-foreground">Created:</strong> {new Date(selectedLog.created_at).toISOString()}</div>
        </div>

        <div class="space-y-1">
          <span class="text-xs font-semibold text-foreground">Cryptographic Hash Signature:</span>
          <p class="p-2.5 bg-muted rounded-lg text-[11px] font-mono break-all text-foreground">{selectedLog.hash}</p>
        </div>

        <div class="space-y-1">
          <span class="text-xs font-semibold text-foreground">Previous Record Hash (prev_hash):</span>
          <p class="p-2.5 bg-muted rounded-lg text-[11px] font-mono break-all text-muted-foreground">{selectedLog.prev_hash || 'Genesis'}</p>
        </div>

        {#if selectedLog.metadata}
          <div class="space-y-1">
            <span class="text-xs font-semibold text-foreground">Contextual Metadata:</span>
            <pre class="p-3 bg-muted rounded-lg text-[11px] font-mono overflow-x-auto text-foreground">{JSON.stringify(selectedLog.metadata, null, 2)}</pre>
          </div>
        {/if}

        <div class="pt-2 text-right">
          <Button
            variant="outline"
            size="sm"
            onclick={() => selectedLog = null}
          >
            Close
          </Button>
        </div>
      </div>
    </div>
  {/if}
</div>
