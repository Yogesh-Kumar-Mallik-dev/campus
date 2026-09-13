<script lang="ts">
  /**
   * BLOCK_WEB_ATTENDANCE_ROSTER_001
   * Subsystem: Rank 4 - Attendance Management System (attendance)
   * Purpose:   Interactive multi-mode roll call roster with batch marking, RFID/biometric indicators, and real-time tally.
   */
  import type { AttendanceSession, AttendanceRecord, AttendanceRecordStatus } from '$lib/types/attendance';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';

  let {
    session,
    onSaveRecords = async (_r: AttendanceRecord[]) => {},
    onFinalize = async () => {}
  }: {
    session: AttendanceSession;
    onSaveRecords?: (records: AttendanceRecord[]) => Promise<void>;
    onFinalize?: () => Promise<void>;
  } = $props();

  let records = $state<AttendanceRecord[]>([]);

  $effect(() => {
    if (session.records && session.records.length > 0) {
      records = session.records;
    } else if (records.length === 0) {
      records = [
        {
          id: 'rec_1',
          tenant_id: session.tenant_id,
          session_id: session.id,
          student_id: 'stu_101',
          student_name: 'Aarav Sharma',
          roll_number: '2026-CS-0001',
          status: 'PRESENT',
          marked_at: new Date().toISOString(),
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        },
        {
          id: 'rec_2',
          tenant_id: session.tenant_id,
          session_id: session.id,
          student_id: 'stu_102',
          student_name: 'Diya Patel',
          roll_number: '2026-CS-0002',
          status: 'ABSENT',
          marked_at: new Date().toISOString(),
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        },
        {
          id: 'rec_3',
          tenant_id: session.tenant_id,
          session_id: session.id,
          student_id: 'stu_103',
          student_name: 'Rohan Iyer',
          roll_number: '2026-CS-0003',
          status: 'LATE',
          marked_at: new Date().toISOString(),
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      ];
    }
  });

  let isSubmitting = $state(false);

  let presentCount = $derived(records.filter(r => r.status === 'PRESENT' || r.status === 'LATE').length);
  let absentCount = $derived(records.filter(r => r.status === 'ABSENT').length);
  let excusedCount = $derived(records.filter(r => r.status === 'EXCUSED_MEDICAL').length);
  let attendancePct = $derived(records.length > 0 ? Math.round(((presentCount + excusedCount) / records.length) * 100) : 0);

  function updateStatus(index: number, status: AttendanceRecordStatus) {
    if (session.status === 'FINALIZED') return;
    records[index].status = status;
  }

  function markAll(status: AttendanceRecordStatus) {
    if (session.status === 'FINALIZED') return;
    records = records.map(r => ({ ...r, status }));
  }

  async function handleSave() {
    isSubmitting = true;
    try {
      await onSaveRecords(records);
    } finally {
      isSubmitting = false;
    }
  }

  async function handleFinalize() {
    isSubmitting = true;
    try {
      await onFinalize();
    } finally {
      isSubmitting = false;
    }
  }
</script>

<Card class="border shadow-sm">
  <CardHeader class="pb-4 border-b bg-muted/30">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div>
        <div class="flex items-center gap-2">
          <CardTitle class="text-xl font-bold">
            {session.subject_name || session.subject_id} — Roll Call
          </CardTitle>
          <Badge variant={session.status === 'FINALIZED' ? 'secondary' : session.status === 'OPEN' ? 'default' : 'outline'}>
            {session.status}
          </Badge>
        </div>
        <CardDescription class="mt-1">
          Date: <span class="font-medium text-foreground">{session.session_date}</span> |
          Time: <span class="font-medium text-foreground">{session.start_time} - {session.end_time}</span> |
          Mode: <span class="font-medium text-foreground">{session.mode}</span>
        </CardDescription>
      </div>

      <!-- Quick stats pill -->
      <div class="flex items-center gap-3 bg-background border px-3 py-1.5 rounded-lg text-sm">
        <div>
          <span class="text-muted-foreground">Present:</span>
          <span class="font-semibold text-emerald-600 dark:text-emerald-400 ml-1">{presentCount}</span>
        </div>
        <div class="h-4 w-px bg-border"></div>
        <div>
          <span class="text-muted-foreground">Absent:</span>
          <span class="font-semibold text-rose-600 dark:text-rose-400 ml-1">{absentCount}</span>
        </div>
        <div class="h-4 w-px bg-border"></div>
        <div>
          <span class="text-muted-foreground">Rate:</span>
          <span class="font-bold ml-1 {attendancePct < 75 ? 'text-amber-600' : 'text-emerald-600'}">{attendancePct}%</span>
        </div>
      </div>
    </div>
  </CardHeader>

  <CardContent class="p-0">
    {#if session.status !== 'FINALIZED'}
      <div class="flex items-center justify-between px-6 py-3 bg-muted/10 border-b text-xs text-muted-foreground">
        <span>Quick Actions:</span>
        <div class="flex gap-2">
          <Button variant="outline" size="sm" class="h-7 text-xs" onclick={() => markAll('PRESENT')}>
            Mark All Present
          </Button>
          <Button variant="outline" size="sm" class="h-7 text-xs" onclick={() => markAll('ABSENT')}>
            Mark All Absent
          </Button>
        </div>
      </div>
    {/if}

    <div class="overflow-x-auto">
      <table class="w-full text-sm text-left">
        <thead class="bg-muted/40 text-muted-foreground uppercase text-xs">
          <tr>
            <th class="px-6 py-3">Roll No</th>
            <th class="px-6 py-3">Student Name</th>
            <th class="px-6 py-3 text-center">Status</th>
            <th class="px-6 py-3">Remarks</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          {#each records as record, idx (record.student_id)}
            <tr class="hover:bg-muted/20 transition-colors">
              <td class="px-6 py-4 font-mono font-medium text-xs">
                {record.roll_number || record.student_id}
              </td>
              <td class="px-6 py-4 font-medium">
                {record.student_name || 'Student ' + (idx + 1)}
              </td>
              <td class="px-6 py-4">
                <div class="flex items-center justify-center gap-1.5">
                  <button
                    type="button"
                    disabled={session.status === 'FINALIZED'}
                    class="px-2.5 py-1 text-xs rounded-md font-medium transition-all {record.status === 'PRESENT' ? 'bg-emerald-600 text-white shadow-sm' : 'bg-muted hover:bg-emerald-100 dark:hover:bg-emerald-950/50 text-muted-foreground'}"
                    onclick={() => updateStatus(idx, 'PRESENT')}
                  >
                    P
                  </button>
                  <button
                    type="button"
                    disabled={session.status === 'FINALIZED'}
                    class="px-2.5 py-1 text-xs rounded-md font-medium transition-all {record.status === 'LATE' ? 'bg-amber-600 text-white shadow-sm' : 'bg-muted hover:bg-amber-100 dark:hover:bg-amber-950/50 text-muted-foreground'}"
                    onclick={() => updateStatus(idx, 'LATE')}
                  >
                    L
                  </button>
                  <button
                    type="button"
                    disabled={session.status === 'FINALIZED'}
                    class="px-2.5 py-1 text-xs rounded-md font-medium transition-all {record.status === 'ABSENT' ? 'bg-rose-600 text-white shadow-sm' : 'bg-muted hover:bg-rose-100 dark:hover:bg-rose-950/50 text-muted-foreground'}"
                    onclick={() => updateStatus(idx, 'ABSENT')}
                  >
                    A
                  </button>
                  {#if record.status === 'EXCUSED_MEDICAL'}
                    <span class="px-2.5 py-1 text-xs rounded-md bg-blue-600 text-white font-medium">
                      MED
                    </span>
                  {/if}
                </div>
              </td>
              <td class="px-6 py-4 text-xs text-muted-foreground">
                {record.remarks || (record.device_id ? `Synced from ${record.device_id}` : '—')}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </CardContent>

  <CardFooter class="flex items-center justify-between p-4 border-t bg-muted/10">
    <span class="text-xs text-muted-foreground">
      {records.length} students enrolled in cohort
    </span>
    {#if session.status !== 'FINALIZED'}
      <div class="flex gap-2">
        <Button variant="outline" disabled={isSubmitting} onclick={handleSave}>
          Save Draft
        </Button>
        <Button disabled={isSubmitting} onclick={handleFinalize}>
          Finalize Session
        </Button>
      </div>
    {:else}
      <Badge variant="secondary">Session Finalized & Locked</Badge>
    {/if}
  </CardFooter>
</Card>
