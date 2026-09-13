<script lang="ts">
  /**
   * BLOCK_WEB_EVENTS_CHECKIN_001
   * Subsystem: Rank 12 - Event Organisation System (events)
   * Purpose:   Event entrance scanner terminal, single-use QR ticket punch validation, and real-time turnstile attendee tracker.
   */
  import type { EventTicket } from '$lib/types/events';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Badge } from '$lib/components/ui/badge';

  let tickets = $state<EventTicket[]>([
    {
      id: 'tck-01',
      tenantId: 'tenant-main',
      eventId: 'evnt-01',
      eventTitle: 'Annual Tech Symposium 2026',
      venueName: 'Main University Auditorium',
      studentId: 'stu-001',
      studentName: 'Aarav Sharma',
      rollNumber: 'CS22B001',
      ticketCode: 'TCK-8f921a',
      status: 'CHECKED_IN',
      checkedInAt: '2026-09-25T09:14:00Z',
      checkedInById: 'staff-gate-01',
      createdAt: '2026-09-10T00:00:00Z',
      updatedAt: '2026-09-25T09:14:00Z'
    },
    {
      id: 'tck-02',
      tenantId: 'tenant-main',
      eventId: 'evnt-01',
      eventTitle: 'Annual Tech Symposium 2026',
      venueName: 'Main University Auditorium',
      studentId: 'stu-002',
      studentName: 'Diya Patel',
      rollNumber: 'CS22B014',
      ticketCode: 'TCK-4c11b0',
      status: 'CONFIRMED',
      createdAt: '2026-09-12T00:00:00Z',
      updatedAt: '2026-09-12T00:00:00Z'
    }
  ]);

  let scanCodeInput = $state('TCK-4c11b0');
  let scanMessage = $state<{ type: 'SUCCESS' | 'ERROR' | 'INFO'; text: string } | null>(null);

  function handleScan() {
    if (!scanCodeInput.trim()) return;
    const found = tickets.find((t) => t.ticketCode.toUpperCase() === scanCodeInput.trim().toUpperCase());
    if (!found) {
      scanMessage = { type: 'ERROR', text: `Invalid Ticket: "${scanCodeInput}" not found in tenant registry.` };
      return;
    }
    if (found.status === 'CHECKED_IN') {
      scanMessage = {
        type: 'ERROR',
        text: `Double-Punch Conflict! Ticket ${found.ticketCode} was already redeemed at ${new Date(found.checkedInAt!).toLocaleTimeString()}.`
      };
      return;
    }

    tickets = tickets.map((t) => {
      if (t.id === found.id) {
        return {
          ...t,
          status: 'CHECKED_IN',
          checkedInAt: new Date().toISOString(),
          checkedInById: 'staff-scanner-web',
          updatedAt: new Date().toISOString()
        };
      }
      return t;
    });

    scanMessage = {
      type: 'SUCCESS',
      text: `Admitted! Valid Pass for ${found.studentName} (${found.rollNumber}) to ${found.eventTitle}.`
    };
    scanCodeInput = '';
  }

  const checkedInCount = $derived(tickets.filter((t) => t.status === 'CHECKED_IN').length);
</script>

<div class="space-y-6">
  <!-- Scanner Console -->
  <div class="bg-card border rounded-2xl p-6 shadow-sm space-y-4">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b pb-4">
      <div>
        <h3 class="text-base font-bold text-foreground">Gate Turnstile QR Scanner</h3>
        <p class="text-xs text-muted-foreground">Scan 2D barcode or enter 6-byte hexadecimal ticket tokens.</p>
      </div>
      <div class="flex items-center gap-2">
        <Badge variant="outline" class="font-mono text-xs">
          Turnstile Checked-in: {checkedInCount} / {tickets.length}
        </Badge>
      </div>
    </div>

    <!-- Scan Input Form -->
    <div class="flex flex-col sm:flex-row items-center gap-3">
      <div class="relative flex-1 w-full">
        <Input
          placeholder="Scan or type Ticket Code (e.g. TCK-4c11b0)..."
          bind:value={scanCodeInput}
          class="font-mono text-xs pl-8 uppercase"
        />
        <span class="absolute left-2.5 top-2 text-xs">📷</span>
      </div>
      <Button onclick={handleScan} class="w-full sm:w-auto">
        Verify & Admit
      </Button>
    </div>

    <!-- Scan Result Banner -->
    {#if scanMessage}
      <div
        class="p-4 rounded-xl border text-xs font-semibold flex items-center justify-between {scanMessage.type ===
        'SUCCESS'
          ? 'bg-emerald-50 text-emerald-900 border-emerald-200'
          : 'bg-rose-50 text-rose-900 border-rose-200'}"
      >
        <span>{scanMessage.text}</span>
        <button onclick={() => (scanMessage = null)} class="text-xs font-bold px-2">✕</button>
      </div>
    {/if}
  </div>

  <!-- Recent Scans Matrix -->
  <div class="bg-card border rounded-2xl overflow-hidden shadow-sm">
    <div class="p-4 border-b">
      <h3 class="text-sm font-bold text-foreground">Turnstile Access Logs ({tickets.length})</h3>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead class="bg-muted/50 text-muted-foreground font-semibold uppercase tracking-wider border-b">
          <tr>
            <th class="p-3">Ticket Code</th>
            <th class="p-3">Attendee</th>
            <th class="p-3">Event & Venue</th>
            <th class="p-3">Status</th>
            <th class="p-3">Punch Time</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          {#each tickets as t}
            <tr class="hover:bg-muted/30 transition-colors">
              <td class="p-3 font-mono font-bold text-foreground">{t.ticketCode}</td>
              <td class="p-3">
                <div class="font-semibold text-foreground">{t.studentName ?? t.studentId}</div>
                <div class="text-[10px] text-muted-foreground font-mono">{t.rollNumber ?? 'N/A'}</div>
              </td>
              <td class="p-3">
                <div class="font-medium text-foreground">{t.eventTitle ?? 'Campus Event'}</div>
                <div class="text-[10px] text-muted-foreground">{t.venueName ?? 'Main Venue'}</div>
              </td>
              <td class="p-3">
                <Badge variant={t.status === 'CHECKED_IN' ? 'default' : 'secondary'}>
                  {t.status}
                </Badge>
              </td>
              <td class="p-3 text-muted-foreground">
                {t.checkedInAt ? new Date(t.checkedInAt).toLocaleTimeString() : 'Not Scanned'}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>
