<!--
  BLOCK_WEB_WHISTLEBLOWER_FORM_001
  Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
  Purpose:   Zero-Knowledge Anonymous Grievance Intake & Anti-Ragging submission with one-time tracking token generation.
-->
<script lang="ts">
  import type { WhistleblowerCategory, WhistleblowerSeverity, WhistleblowerReport } from '$lib/types/whistleblower';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let { onReportSubmitted }: { onReportSubmitted?: (report: WhistleblowerReport, token: string) => void } = $props();

  let category = $state<WhistleblowerCategory>('ANTI_RAGGING');
  let severity = $state<WhistleblowerSeverity>('HIGH');
  let title = $state<string>('');
  let description = $state<string>('');
  let isSubmitting = $state<boolean>(false);

  // Success state modal
  let generatedToken = $state<string | null>(null);
  let createdReport = $state<WhistleblowerReport | null>(null);
  let isCopied = $state<boolean>(false);

  function handleSubmit() {
    if (!title.trim() || !description.trim()) return;

    isSubmitting = true;

    // Simulate cryptographic token generation
    const hex = Array.from({ length: 16 }, () => Math.floor(Math.random() * 16).toString(16)).join('');
    const rawTok = `WB-tok-${hex}`;
    const repNum = `WB-2026-000${Math.floor(10 + Math.random() * 90)}`;

    const newReport: WhistleblowerReport = {
      id: `wb-${Date.now()}`,
      tenantId: 'tenant-demo',
      reportNumber: repNum,
      category: category,
      severity: severity,
      title: title.trim(),
      description: description.trim(),
      status: 'SUBMITTED',
      submittedAt: new Date().toISOString(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    createdReport = newReport;
    generatedToken = rawTok;
    isSubmitting = false;

    onReportSubmitted?.(newReport, rawTok);
  }

  function handleCopyToken() {
    if (!generatedToken) return;
    navigator.clipboard.writeText(generatedToken);
    isCopied = true;
    setTimeout(() => (isCopied = false), 2000);
  }

  function handleReset() {
    title = '';
    description = '';
    generatedToken = null;
    createdReport = null;
  }
</script>

<div class="max-w-2xl mx-auto space-y-6">
  <!-- Anonymous Privacy Notice -->
  <div class="p-5 bg-gradient-to-r from-slate-900 to-indigo-950 rounded-2xl border border-indigo-800 text-white shadow-md space-y-2">
    <div class="flex items-center gap-2">
      <span class="text-indigo-400 font-bold text-sm">🔒 ZERO-KNOWLEDGE ANONYMOUS PORTAL</span>
    </div>
    <h2 class="text-xl font-bold">Confidential Grievance & Anti-Ragging Intake</h2>
    <p class="text-xs text-slate-300 leading-relaxed">
      Your submission is completely detached from your user identity. No IP addresses, user credentials, or timestamps are tied to your identity. You will be provided with a unique cryptographic tracking key to follow up securely.
    </p>
  </div>

  {#if !generatedToken}
    <!-- Form Body -->
    <div class="p-6 bg-card rounded-2xl border shadow-sm space-y-4">
      <div class="space-y-3">
        <div>
          <label for="wb-cat" class="text-xs font-semibold text-muted-foreground block mb-1">Grievance Category</label>
          <select
            id="wb-cat"
            bind:value={category}
            class="w-full bg-background text-foreground text-sm rounded-lg border px-3 py-2"
          >
            <option value="ANTI_RAGGING">⚠️ Anti-Ragging & Senior Intimidation</option>
            <option value="HARASSMENT">🚨 Physical / Verbal Harassment</option>
            <option value="ACADEMIC_CORRUPTION">📚 Academic Integrity & Examination Corruption</option>
            <option value="FINANCIAL_FRAUD">💰 Financial Irregularity / Embezzlement</option>
            <option value="SAFETY_VIOLATION">🦺 Campus Health & Safety Breach</option>
            <option value="OTHER">🔍 Other Sensitive Concern</option>
          </select>
        </div>

        <div>
          <label for="wb-sev" class="text-xs font-semibold text-muted-foreground block mb-1">Urgency & Severity</label>
          <select
            id="wb-sev"
            bind:value={severity}
            class="w-full bg-background text-foreground text-sm rounded-lg border px-3 py-2"
          >
            <option value="LOW">Low (Informational Observation)</option>
            <option value="MEDIUM">Medium (Policy Infraction)</option>
            <option value="HIGH">High (Severe Harassment / Corruption)</option>
            <option value="CRITICAL">Critical (Immediate Danger / Urgent Intervention)</option>
          </select>
        </div>

        <div>
          <label for="wb-title" class="text-xs font-semibold text-muted-foreground block mb-1">Report Subject Summary</label>
          <Input id="wb-title" bind:value={title} placeholder="e.g. Unofficial mandatory midnight gathering in Block C" />
        </div>

        <div>
          <label for="wb-desc" class="text-xs font-semibold text-muted-foreground block mb-1">Detailed Incident Narrative</label>
          <textarea
            id="wb-desc"
            bind:value={description}
            rows={6}
            placeholder="Please detail dates, approximate times, locations, and nature of the incident. Do NOT include your own name or identifiable contact info..."
            class="w-full bg-background text-foreground text-xs rounded-xl border p-3.5 focus:ring-1 focus:ring-primary"
          ></textarea>
        </div>
      </div>

      <div class="pt-3 border-t flex justify-end">
        <Button
          onclick={handleSubmit}
          disabled={isSubmitting || !title.trim() || !description.trim()}
          class="bg-indigo-600 hover:bg-indigo-700 text-white font-bold px-6"
        >
          Submit Anonymous Report →
        </Button>
      </div>
    </div>
  {:else}
    <!-- One-Time Secret Token Presentation Card -->
    <div class="p-6 bg-card rounded-2xl border-2 border-indigo-500 shadow-xl space-y-4 text-center">
      <div class="w-12 h-12 rounded-full bg-emerald-100 text-emerald-700 font-bold text-xl flex items-center justify-center mx-auto">
        ✓
      </div>

      <div class="space-y-1">
        <h3 class="text-xl font-black text-foreground">Grievance Successfully Lodged</h3>
        <p class="text-xs text-muted-foreground">
          Assigned Reference Number: <span class="font-mono font-bold text-primary">{createdReport?.reportNumber}</span>
        </p>
      </div>

      <div class="p-4 bg-amber-50 border border-amber-200 rounded-xl text-left space-y-2">
        <div class="text-xs font-bold text-amber-900 flex items-center gap-1.5">
          ⚠️ CRITICAL: SAVE YOUR SECRET TRACKING TOKEN
        </div>
        <p class="text-[11px] text-amber-800 leading-relaxed">
          Because this report is entirely anonymous, the token below is the ONLY way to check investigation progress, receive messages from the committee, or provide further evidence. We cannot recover this token for you.
        </p>

        <div class="flex items-center gap-2 pt-2">
          <Input readonly value={generatedToken} class="font-mono text-xs font-bold bg-white text-indigo-900" />
          <Button size="sm" onclick={handleCopyToken} class="font-bold bg-indigo-600 hover:bg-indigo-700 text-white shrink-0">
            {isCopied ? 'Copied! ✓' : 'Copy Key'}
          </Button>
        </div>
      </div>

      <div class="pt-2 flex justify-center gap-3">
        <Button variant="outline" onclick={handleReset}>
          Submit Another Report
        </Button>
      </div>
    </div>
  {/if}
</div>
