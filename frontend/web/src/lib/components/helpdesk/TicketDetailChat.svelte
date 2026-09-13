<!--
  BLOCK_WEB_HELPDESK_CHAT_001
  Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
  Purpose:   Interactive Ticket Resolution Desk, Staff Internal Notes, Live Chat Thread, SLA Escalation, and CSAT Feedback.
-->
<script lang="ts">
  import type { HelpdeskTicket, HelpdeskMessage, HelpdeskTicketStatus } from '$lib/types/helpdesk';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let { ticket, onBack }: { ticket?: HelpdeskTicket; onBack?: () => void } = $props();

  // Fallback demo ticket if none selected
  let currentTicket = $state<HelpdeskTicket>(
    ticket || {
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
    }
  );

  let messages = $state<HelpdeskMessage[]>([
    {
      id: 'm-1',
      tenantId: 'tenant-demo',
      ticketId: currentTicket.id,
      senderId: 'stu-1',
      senderName: 'Aarav Sharma (Student)',
      isStaffReply: false,
      isInternalNote: false,
      message: 'Hello, I noticed the clash in the newly published schedule timetable. Kindly advise which exam I should appear for first.',
      createdAt: new Date(Date.now() - 5 * 3600000).toISOString(),
    },
    {
      id: 'm-2',
      tenantId: 'tenant-demo',
      ticketId: currentTicket.id,
      senderId: 'staff-exam-1',
      senderName: 'Dr. Ramesh Nair (Staff)',
      isStaffReply: true,
      isInternalNote: true,
      message: 'INTERNAL NOTE: Checked CS Department Lab Slot B. We have a buffer slot on 19 Oct morning for CS392. Can reschedule the lab.',
      createdAt: new Date(Date.now() - 3 * 3600000).toISOString(),
    },
    {
      id: 'm-3',
      tenantId: 'tenant-demo',
      ticketId: currentTicket.id,
      senderId: 'staff-exam-1',
      senderName: 'Dr. Ramesh Nair (Staff)',
      isStaffReply: true,
      isInternalNote: false,
      message: 'Dear Aarav, we are arranging a special slot for your CS392 Compiler Lab on 19 Oct 09:00 AM. Please confirm if that slot works for you.',
      createdAt: new Date(Date.now() - 2 * 3600000).toISOString(),
    },
  ]);

  let newMessageText = $state<string>('');
  let isInternalNote = $state<boolean>(false);
  let isEscalating = $state<boolean>(false);
  let escalationReason = $state<string>('');
  let isRatingSubmitted = $state<boolean>(false);
  let selectedRating = $state<number>(5);
  let feedbackText = $state<string>('');

  function handleSendMessage() {
    if (!newMessageText.trim()) return;

    const newMsg: HelpdeskMessage = {
      id: `m-${Date.now()}`,
      tenantId: currentTicket.tenantId,
      ticketId: currentTicket.id,
      senderId: 'current-user',
      senderName: isInternalNote ? 'Staff Officer (Internal)' : 'Staff Officer',
      isStaffReply: true,
      isInternalNote: isInternalNote,
      message: newMessageText.trim(),
      createdAt: new Date().toISOString(),
    };

    messages = [...messages, newMsg];
    newMessageText = '';

    if (!isInternalNote && currentTicket.status === 'IN_PROGRESS') {
      currentTicket.status = 'WAITING_FOR_APPLICANT';
    }
  }

  function handleStatusChange(nextStatus: HelpdeskTicketStatus) {
    currentTicket.status = nextStatus;
    if (nextStatus === 'RESOLVED') {
      currentTicket.resolvedAt = new Date().toISOString();
    } else if (nextStatus === 'CLOSED') {
      currentTicket.closedAt = new Date().toISOString();
    }
  }

  function handleEscalate() {
    if (!escalationReason.trim()) return;
    currentTicket.priority = 'URGENT';
    isEscalating = false;
    escalationReason = '';

    const escMsg: HelpdeskMessage = {
      id: `m-${Date.now()}`,
      tenantId: currentTicket.tenantId,
      ticketId: currentTicket.id,
      senderId: 'system',
      senderName: 'SLA Radar Automation',
      isStaffReply: true,
      isInternalNote: true,
      message: `🚨 ESCALATION TRIGGERED: Priority raised to URGENT. Reason: ${escalationReason}`,
      createdAt: new Date().toISOString(),
    };
    messages = [...messages, escMsg];
  }

  function handleRateTicket() {
    currentTicket.rating = selectedRating;
    currentTicket.feedback = feedbackText;
    isRatingSubmitted = true;
  }
</script>

<div class="space-y-6">
  <!-- Top Navigation & Ticket Summary -->
  <div class="p-5 bg-card rounded-xl border shadow-sm space-y-4">
    <div class="flex flex-wrap items-center justify-between gap-4 border-b pb-4">
      <div class="flex items-center gap-3">
        {#if onBack}
          <Button variant="ghost" size="sm" onclick={onBack}>
            ← Back to Tickets
          </Button>
        {/if}
        <div>
          <div class="flex items-center gap-2">
            <span class="font-mono text-xs font-bold text-primary bg-primary/10 px-2.5 py-0.5 rounded">
              {currentTicket.ticketNumber}
            </span>
            <Badge variant="outline">{currentTicket.categoryName}</Badge>
            <span class="text-xs font-bold px-2.5 py-0.5 rounded-full {currentTicket.priority === 'URGENT' ? 'bg-rose-500 text-white' : 'bg-blue-500 text-white'}">
              {currentTicket.priority}
            </span>
            <Badge variant="default">{currentTicket.status}</Badge>
          </div>
          <h2 class="text-xl font-bold text-foreground mt-1">
            {currentTicket.title}
          </h2>
        </div>
      </div>

      <!-- State Transition Controls -->
      <div class="flex flex-wrap items-center gap-2">
        {#if currentTicket.status === 'OPEN'}
          <Button size="sm" onclick={() => handleStatusChange('IN_PROGRESS')}>
            Start Working (In Progress)
          </Button>
        {/if}
        {#if currentTicket.status === 'IN_PROGRESS' || currentTicket.status === 'WAITING_FOR_APPLICANT'}
          <Button size="sm" variant="outline" class="text-emerald-600 border-emerald-300 hover:bg-emerald-50" onclick={() => handleStatusChange('RESOLVED')}>
            ✓ Mark as Resolved
          </Button>
        {/if}
        {#if currentTicket.status === 'RESOLVED'}
          <Button size="sm" variant="secondary" onclick={() => handleStatusChange('CLOSED')}>
            Close Ticket
          </Button>
          <Button size="sm" variant="outline" onclick={() => handleStatusChange('IN_PROGRESS')}>
            Reopen Ticket
          </Button>
        {/if}
        {#if currentTicket.priority !== 'URGENT'}
          <Button size="sm" variant="destructive" onclick={() => (isEscalating = true)}>
            ⚡ Escalate SLA
          </Button>
        {/if}
      </div>
    </div>

    <!-- Ticket Metadata Grid -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs">
      <div>
        <span class="text-muted-foreground block">Applicant / Student:</span>
        <span class="font-semibold text-foreground">{currentTicket.requesterName || 'Student'}</span>
      </div>
      <div>
        <span class="text-muted-foreground block">Assigned Staff Officer:</span>
        <span class="font-semibold text-foreground">{currentTicket.assignedStaffName || 'Unassigned'}</span>
      </div>
      <div>
        <span class="text-muted-foreground block">SLA Deadline:</span>
        <span class="font-semibold text-foreground">{new Date(currentTicket.slaDueAt).toLocaleString()}</span>
      </div>
      <div>
        <span class="text-muted-foreground block">Created At:</span>
        <span class="font-semibold text-foreground">{new Date(currentTicket.createdAt).toLocaleString()}</span>
      </div>
    </div>

    <div class="p-3 bg-muted/40 rounded-lg text-xs text-foreground leading-relaxed">
      <span class="font-bold text-muted-foreground block mb-0.5">Original Inquiry:</span>
      {currentTicket.description}
    </div>
  </div>

  <!-- SLA Escalation Modal / Drawer -->
  {#if isEscalating}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-xl space-y-3">
      <div class="flex items-center justify-between">
        <h4 class="text-sm font-bold text-rose-800">⚡ Escalate Ticket to Priority URGENT</h4>
        <button onclick={() => (isEscalating = false)} class="text-rose-600 font-bold text-sm">✕</button>
      </div>
      <p class="text-xs text-rose-700">
        Escalating will immediately alert the department supervisor, bump priority to URGENT, and log an audit trail.
      </p>
      <Input bind:value={escalationReason} placeholder="Reason for escalation (e.g. SLA timeline breached, academic blocker)..." />
      <div class="flex justify-end gap-2">
        <Button variant="outline" size="sm" onclick={() => (isEscalating = false)}>Cancel</Button>
        <Button variant="destructive" size="sm" onclick={handleEscalate}>Confirm Escalation</Button>
      </div>
    </div>
  {/if}

  <!-- Conversation Thread -->
  <div class="p-5 bg-card rounded-xl border shadow-sm space-y-4">
    <h3 class="text-sm font-bold uppercase tracking-wider text-muted-foreground">
      Live Resolution Desk ({messages.length} messages)
    </h3>

    <div class="space-y-3 max-h-[450px] overflow-y-auto pr-1">
      {#each messages as msg}
        <div class="flex flex-col {msg.isInternalNote ? 'border-l-4 border-amber-500 bg-amber-50/60 p-3 rounded-r-xl' : msg.isStaffReply ? 'bg-primary/5 p-3.5 rounded-xl border' : 'bg-muted/50 p-3.5 rounded-xl border'}">
          <div class="flex items-center justify-between text-xs mb-1.5">
            <div class="flex items-center gap-2">
              <span class="font-bold {msg.isInternalNote ? 'text-amber-800' : 'text-foreground'}">
                {msg.senderName}
              </span>
              {#if msg.isInternalNote}
                <span class="px-2 py-0.5 bg-amber-200 text-amber-900 rounded font-semibold text-[10px]">
                  🔒 STAFF INTERNAL NOTE (HIDDEN FROM STUDENT)
                </span>
              {:else if msg.isStaffReply}
                <span class="px-2 py-0.5 bg-primary/10 text-primary rounded font-semibold text-[10px]">
                  OFFICIAL STAFF RESPONSE
                </span>
              {/if}
            </div>
            <span class="text-muted-foreground text-[11px]">
              {new Date(msg.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
            </span>
          </div>

          <p class="text-xs text-foreground leading-relaxed">
            {msg.message}
          </p>
        </div>
      {/each}
    </div>

    <!-- Message Composer -->
    {#if currentTicket.status !== 'CLOSED'}
      <div class="pt-4 border-t space-y-3">
        <div class="flex items-center justify-between">
          <label class="flex items-center gap-2 cursor-pointer text-xs font-semibold text-amber-800 bg-amber-100/70 px-2.5 py-1 rounded-md">
            <input type="checkbox" bind:checked={isInternalNote} class="rounded" />
            <span>🔒 Write as Staff Internal Note (Restricted to Staff View)</span>
          </label>
        </div>

        <div class="flex gap-2">
          <textarea
            bind:value={newMessageText}
            rows={2}
            placeholder={isInternalNote ? "Write staff-only note for internal review..." : "Type response to applicant..."}
            class="flex-1 bg-background text-foreground text-xs rounded-xl border p-3 focus:outline-none focus:ring-1 focus:ring-primary"
          ></textarea>
          <Button onclick={handleSendMessage} class="self-end px-5 font-bold">
            Send Reply
          </Button>
        </div>
      </div>
    {:else}
      <div class="p-3 bg-muted/40 rounded-xl text-center text-xs text-muted-foreground">
        This ticket is closed. Reopen the ticket to send additional messages.
      </div>
    {/if}
  </div>

  <!-- CSAT Rating Widget for Resolved/Closed Tickets -->
  {#if currentTicket.status === 'RESOLVED' || currentTicket.status === 'CLOSED'}
    <div class="p-5 bg-gradient-to-r from-emerald-50 to-teal-50 border border-emerald-200 rounded-xl shadow-sm space-y-3">
      <div class="flex items-center justify-between">
        <h4 class="text-sm font-bold text-emerald-900">
          🌟 Student Satisfaction & Resolution Feedback
        </h4>
        {#if currentTicket.rating}
          <span class="text-xs font-bold text-emerald-700 bg-emerald-100 px-2.5 py-0.5 rounded-full">
            Feedback Submitted ({currentTicket.rating}/5 Stars)
          </span>
        {/if}
      </div>

      {#if !isRatingSubmitted && !currentTicket.rating}
        <div class="space-y-3">
          <p class="text-xs text-emerald-800">
            How would you rate the resolution speed and support quality for this ticket?
          </p>

          <div class="flex items-center gap-2">
            {#each [1, 2, 3, 4, 5] as star}
              <button
                type="button"
                onclick={() => (selectedRating = star)}
                class="text-2xl transition-transform hover:scale-125 {star <= selectedRating ? 'text-amber-500' : 'text-slate-300'}"
              >
                ★
              </button>
            {/each}
            <span class="text-xs font-bold text-emerald-900 ml-2">({selectedRating} / 5 Stars)</span>
          </div>

          <Input bind:value={feedbackText} placeholder="Optional feedback remarks..." class="bg-white text-xs" />

          <Button size="sm" onclick={handleRateTicket} class="bg-emerald-600 hover:bg-emerald-700 text-white font-bold">
            Submit Rating
          </Button>
        </div>
      {:else}
        <div class="text-xs text-emerald-800 space-y-1">
          <div class="font-bold text-amber-600 text-base">
            {'★'.repeat(currentTicket.rating || selectedRating)}{'☆'.repeat(5 - (currentTicket.rating || selectedRating))}
          </div>
          {#if currentTicket.feedback || feedbackText}
            <p class="italic">"{currentTicket.feedback || feedbackText}"</p>
          {/if}
        </div>
      {/if}
    </div>
  {/if}
</div>
