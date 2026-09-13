<!--
  BLOCK_WEB_WHISTLEBLOWER_TRACKING_001
  Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
  Purpose:   Secret Token Report Lookup & 2-Way Anonymous Investigation Dialogue.
-->
<script lang="ts">
  import type { WhistleblowerReport, WhistleblowerMessage } from '$lib/types/whistleblower';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let lookupToken = $state<string>('');
  let activeReport = $state<WhistleblowerReport | null>(null);
  let messages = $state<WhistleblowerMessage[]>([]);
  let newMessageText = $state<string>('');
  let isSearching = $state<boolean>(false);

  function handleLookup() {
    if (!lookupToken.trim()) return;

    isSearching = true;

    // Simulate looking up by token
    setTimeout(() => {
      activeReport = {
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
      };

      messages = [
        {
          id: 'm-1',
          tenantId: 'tenant-demo',
          reportId: 'wb-001',
          senderType: 'INVESTIGATION_COMMITTEE',
          message: 'The Anti-Ragging Committee has taken cognizance of this report. A surprise inspection of Block C 4th floor has been scheduled with the Dean of Student Welfare. Please let us know if specific room numbers were involved.',
          createdAt: new Date(Date.now() - 12 * 3600000).toISOString(),
        },
      ];

      isSearching = false;
    }, 300);
  }

  function handleSendAnonymousReply() {
    if (!newMessageText.trim() || !activeReport) return;

    const newMsg: WhistleblowerMessage = {
      id: `m-${Date.now()}`,
      tenantId: activeReport.tenantId,
      reportId: activeReport.id,
      senderType: 'ANONYMOUS_REPORTER',
      message: newMessageText.trim(),
      createdAt: new Date().toISOString(),
    };

    messages = [...messages, newMsg];
    newMessageText = '';
  }

  function getStatusBadge(st: string): { variant: 'default' | 'outline' | 'secondary' | 'destructive'; label: string } {
    switch (st) {
      case 'SUBMITTED':
        return { variant: 'outline', label: 'Submitted (Pending Review)' };
      case 'UNDER_INVESTIGATION':
        return { variant: 'default', label: '🔍 Active Investigation' };
      case 'EVIDENCE_REQUESTED':
        return { variant: 'destructive', label: '⚠️ Additional Evidence Requested' };
      case 'RESOLVED':
        return { variant: 'secondary', label: '✓ Resolved & Action Taken' };
      case 'REJECTED':
        return { variant: 'outline', label: '✕ Dismissed' };
      default:
        return { variant: 'outline', label: st };
    }
  }
</script>

<div class="max-w-3xl mx-auto space-y-6">
  <!-- Search Card -->
  <div class="p-6 bg-card rounded-2xl border shadow-sm space-y-3">
    <h3 class="text-base font-bold text-foreground">Track Anonymous Grievance Status</h3>
    <p class="text-xs text-muted-foreground">
      Enter the secret tracking token provided upon submission to view real-time investigation findings and exchange confidential messages with the committee.
    </p>

    <div class="flex flex-col sm:flex-row gap-2 pt-1">
      <Input
        bind:value={lookupToken}
        placeholder="e.g. WB-tok-4f81c9a02e..."
        class="font-mono text-xs flex-1"
      />
      <Button onclick={handleLookup} disabled={isSearching || !lookupToken.trim()} class="font-bold">
        {isSearching ? 'Searching...' : 'Lookup Report →'}
      </Button>
    </div>
  </div>

  {#if activeReport}
    {@const badge = getStatusBadge(activeReport.status)}
    <!-- Report Details -->
    <div class="p-6 bg-card rounded-2xl border shadow-sm space-y-4">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b pb-4">
        <div>
          <div class="flex items-center gap-2">
            <span class="font-mono text-xs font-bold text-primary bg-primary/10 px-2 py-0.5 rounded">
              {activeReport.reportNumber}
            </span>
            <Badge variant="outline">{activeReport.category}</Badge>
            <Badge variant={badge.variant}>{badge.label}</Badge>
          </div>
          <h2 class="text-lg font-bold text-foreground mt-1.5">{activeReport.title}</h2>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
        <div>
          <span class="text-muted-foreground block">Assigned Investigation Unit:</span>
          <span class="font-semibold text-foreground">{activeReport.assignedOfficerName || 'Under Committee Allocation'}</span>
        </div>
        <div>
          <span class="text-muted-foreground block">Lodged Timestamp:</span>
          <span class="font-semibold text-foreground">{new Date(activeReport.submittedAt).toLocaleDateString()}</span>
        </div>
        <div>
          <span class="text-muted-foreground block">Identity Masking:</span>
          <span class="font-bold text-indigo-600">100% Zero-Knowledge Shield</span>
        </div>
      </div>

      <div class="p-3 bg-muted/40 rounded-xl text-xs space-y-1">
        <span class="font-bold text-muted-foreground block">Original Submission Summary:</span>
        <p class="text-foreground leading-relaxed">{activeReport.description}</p>
      </div>

      {#if activeReport.findingsSummary}
        <div class="p-3.5 bg-emerald-50 text-emerald-900 border border-emerald-200 rounded-xl text-xs space-y-1">
          <span class="font-bold block">✓ Committee Final Findings & Corrective Measures:</span>
          <p>{activeReport.findingsSummary}</p>
        </div>
      {/if}
    </div>

    <!-- 2-Way Message Thread -->
    <div class="p-6 bg-card rounded-2xl border shadow-sm space-y-4">
      <h3 class="text-sm font-bold uppercase tracking-wider text-muted-foreground">
        Anonymous 2-Way Channel ({messages.length} messages)
      </h3>

      <div class="space-y-3 max-h-[350px] overflow-y-auto pr-1">
        {#each messages as msg}
          <div class="p-3.5 rounded-xl border {msg.senderType === 'INVESTIGATION_COMMITTEE' ? 'bg-primary/5 border-primary/20' : 'bg-muted/40'}">
            <div class="flex items-center justify-between text-xs mb-1">
              <span class="font-bold {msg.senderType === 'INVESTIGATION_COMMITTEE' ? 'text-primary' : 'text-indigo-600'}">
                {msg.senderType === 'INVESTIGATION_COMMITTEE' ? '🏛️ Official Investigation Committee' : '👤 Anonymous Reporter (You)'}
              </span>
              <span class="text-muted-foreground text-[10px]">
                {new Date(msg.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </span>
            </div>
            <p class="text-xs text-foreground leading-relaxed">{msg.message}</p>
          </div>
        {/each}
      </div>

      <!-- Anonymous Reply Box -->
      {#if activeReport.status !== 'RESOLVED' && activeReport.status !== 'REJECTED'}
        <div class="pt-3 border-t space-y-2">
          <textarea
            bind:value={newMessageText}
            rows={2}
            placeholder="Send anonymous follow-up or provide requested details..."
            class="w-full bg-background text-foreground text-xs rounded-xl border p-3 focus:ring-1 focus:ring-primary"
          ></textarea>
          <div class="flex justify-end">
            <Button size="sm" onclick={handleSendAnonymousReply} class="bg-indigo-600 hover:bg-indigo-700 text-white font-bold">
              Send Anonymous Reply
            </Button>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
