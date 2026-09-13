<script lang="ts">
  /**
   * BLOCK_WEB_COMPLIANCE_REPORT_001
   * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
   * Purpose:   Accreditation & compliance dashboard for NAAC, NIRF, ISO 27001, and ABET audit reviews.
   */
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';

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
  <!-- Framework Selection Header Card -->
  <Card class="shadow-sm">
    <CardHeader class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 pb-6">
      <div>
        <CardTitle class="text-xl font-bold tracking-tight">Institutional Compliance & Accreditation</CardTitle>
        <CardDescription class="mt-1">
          Generate tamper-evident, cryptographically verified audit summaries for regulatory standards.
        </CardDescription>
      </div>

      <div class="flex items-center gap-3">
        <select
          bind:value={selectedFramework}
          class="h-9 px-3 rounded-md border border-input bg-background text-foreground text-sm font-medium focus:outline-none focus:ring-2 focus:ring-ring"
        >
          <option value="NAAC">NAAC (Criterion 6 & Governance)</option>
          <option value="NIRF">NIRF (Institutional Data Verification)</option>
          <option value="ISO_27001">ISO/IEC 27001 (Information Security)</option>
          <option value="ABET">ABET (Computing & Engineering Accreditation)</option>
        </select>

        <Button
          onclick={() => onGenerate(selectedFramework)}
          disabled={isLoading}
          size="sm"
        >
          {isLoading ? 'Compiling...' : 'Generate Report'}
        </Button>
      </div>
    </CardHeader>
  </Card>

  {#if report}
    <!-- KPI Summary Metrics -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <Card class="shadow-sm">
        <CardContent class="p-5">
          <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Total Recorded Events</span>
          <div class="mt-2 text-3xl font-bold tracking-tight text-foreground">{report.total_events.toLocaleString()}</div>
          <div class="mt-1 text-xs text-muted-foreground">Within reporting window</div>
        </CardContent>
      </Card>

      <Card class="shadow-sm">
        <CardContent class="p-5">
          <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Successful Mutations</span>
          <div class="mt-2 text-3xl font-bold tracking-tight text-emerald-600 dark:text-emerald-400">{report.successful_events.toLocaleString()}</div>
          <div class="mt-1 text-xs text-muted-foreground">
            {((report.successful_events / (report.total_events || 1)) * 100).toFixed(1)}% nominal rate
          </div>
        </CardContent>
      </Card>

      <Card class="shadow-sm">
        <CardContent class="p-5">
          <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Security Incidents</span>
          <div class="mt-2 text-3xl font-bold tracking-tight {report.security_incidents > 0 ? 'text-amber-600 dark:text-amber-400' : 'text-foreground'}">
            {report.security_incidents.toLocaleString()}
          </div>
          <div class="mt-1 text-xs text-muted-foreground">Lockouts, breaches, failures</div>
        </CardContent>
      </Card>

      <Card class="shadow-sm">
        <CardContent class="p-5">
          <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Ledger Cryptographic Integrity</span>
          <div class="mt-2 flex items-center gap-2">
            {#if report.chain_integrity}
              <Badge variant="outline" class="bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/30">
                ✓ INTACT & VERIFIED
              </Badge>
            {:else}
              <Badge variant="destructive">
                ⚠ TAMPERING DETECTED
              </Badge>
            {/if}
          </div>
          <div class="mt-1 text-xs text-muted-foreground">SHA-256 Merkle chain check</div>
        </CardContent>
      </Card>
    </div>

    <!-- Top Activities Breakdown -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <Card class="shadow-sm">
        <CardHeader class="pb-3">
          <CardTitle class="text-sm font-bold uppercase tracking-wider text-muted-foreground">Top Audit Actions by Volume</CardTitle>
        </CardHeader>
        <CardContent class="divide-y divide-border">
          {#each report.top_actions as item}
            <div class="py-2.5 flex items-center justify-between text-xs">
              <span class="font-mono font-medium text-foreground">{item.action}</span>
              <Badge variant="secondary" class="font-mono">{item.count.toLocaleString()}</Badge>
            </div>
          {/each}
        </CardContent>
      </Card>

      <Card class="shadow-sm">
        <CardHeader class="pb-3">
          <CardTitle class="text-sm font-bold uppercase tracking-wider text-muted-foreground">Most Active Actors</CardTitle>
        </CardHeader>
        <CardContent class="divide-y divide-border">
          {#each report.top_actors as actor}
            <div class="py-2.5 flex items-center justify-between text-xs">
              <div class="flex items-center gap-2">
                <span class="font-medium text-foreground">{actor.actor_id}</span>
                <Badge variant="outline" class="text-[10px] uppercase">{actor.actor_role}</Badge>
              </div>
              <Badge variant="secondary" class="font-mono">{actor.count.toLocaleString()} ops</Badge>
            </div>
          {/each}
        </CardContent>
      </Card>
    </div>
  {/if}
</div>
