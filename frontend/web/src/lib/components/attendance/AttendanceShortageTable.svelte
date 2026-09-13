<script lang="ts">
  /**
   * BLOCK_WEB_ATTENDANCE_SHORTAGE_TABLE_001
   * Subsystem: Rank 4 - Attendance Management System (attendance)
   * Purpose:   Student shortage tracker flagging attendance rate <75% and managing medical leave condonations.
   */
  import type { StudentAttendanceSummary } from '$lib/types/attendance';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';

  let {
    summaries = [
      {
        student_id: 'stu_101',
        subject_id: 'sub_cs101',
        total_sessions: 40,
        present_sessions: 36,
        excused_sessions: 0,
        absent_sessions: 4,
        attendance_percentage: 90.0,
        is_shortage: false
      },
      {
        student_id: 'stu_102',
        subject_id: 'sub_cs101',
        total_sessions: 40,
        present_sessions: 24,
        excused_sessions: 2,
        absent_sessions: 14,
        attendance_percentage: 65.0,
        is_shortage: true
      },
      {
        student_id: 'stu_103',
        subject_id: 'sub_cs101',
        total_sessions: 40,
        present_sessions: 28,
        excused_sessions: 4,
        absent_sessions: 8,
        attendance_percentage: 80.0,
        is_shortage: false
      }
    ],
    onApplyMedicalLeave = (_studentID: string) => {}
  }: {
    summaries?: StudentAttendanceSummary[];
    onApplyMedicalLeave?: (studentID: string) => void;
  } = $props();

  let filterShortageOnly = $state(false);

  let filteredSummaries = $derived(filterShortageOnly
    ? summaries.filter(s => s.is_shortage)
    : summaries);

  let totalShortages = $derived(summaries.filter(s => s.is_shortage).length);
</script>

<Card class="border shadow-sm">
  <CardHeader class="pb-4 border-b bg-muted/30">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div>
        <CardTitle class="text-xl font-bold flex items-center gap-2">
          Cohort Attendance & Shortage Monitor
          {#if totalShortages > 0}
            <Badge variant="destructive" class="ml-1">
              {totalShortages} At Risk (&lt;75%)
            </Badge>
          {/if}
        </CardTitle>
        <CardDescription class="mt-1">
          Academic threshold: Minimum 75.0% attendance required to sit for semester examinations.
        </CardDescription>
      </div>

      <div class="flex items-center gap-2">
        <Button
          variant={filterShortageOnly ? 'default' : 'outline'}
          size="sm"
          onclick={() => (filterShortageOnly = !filterShortageOnly)}
        >
          {filterShortageOnly ? 'Show All Students' : 'Filter Shortages Only'}
        </Button>
      </div>
    </div>
  </CardHeader>

  <CardContent class="p-0">
    <div class="overflow-x-auto">
      <table class="w-full text-sm text-left">
        <thead class="bg-muted/40 text-muted-foreground uppercase text-xs">
          <tr>
            <th class="px-6 py-3">Student ID</th>
            <th class="px-6 py-3 text-center">Total Sessions</th>
            <th class="px-6 py-3 text-center">Present</th>
            <th class="px-6 py-3 text-center">Excused Med</th>
            <th class="px-6 py-3 text-center">Absent</th>
            <th class="px-6 py-3 text-center">Effective Rate</th>
            <th class="px-6 py-3 text-center">Status</th>
            <th class="px-6 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          {#each filteredSummaries as item (item.student_id)}
            <tr class="hover:bg-muted/20 transition-colors {item.is_shortage ? 'bg-rose-50/50 dark:bg-rose-950/10' : ''}">
              <td class="px-6 py-4 font-mono font-medium text-xs">
                {item.student_id}
              </td>
              <td class="px-6 py-4 text-center font-medium">
                {item.total_sessions}
              </td>
              <td class="px-6 py-4 text-center text-emerald-600 dark:text-emerald-400 font-semibold">
                {item.present_sessions}
              </td>
              <td class="px-6 py-4 text-center text-blue-600 dark:text-blue-400 font-medium">
                {item.excused_sessions}
              </td>
              <td class="px-6 py-4 text-center text-rose-600 dark:text-rose-400 font-medium">
                {item.absent_sessions}
              </td>
              <td class="px-6 py-4 text-center">
                <span class="font-bold {item.is_shortage ? 'text-rose-600' : 'text-emerald-600'}">
                  {item.attendance_percentage.toFixed(1)}%
                </span>
              </td>
              <td class="px-6 py-4 text-center">
                {#if item.is_shortage}
                  <Badge variant="destructive">SHORTAGE</Badge>
                {:else}
                  <Badge variant="secondary" class="bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-300">
                    ELIGIBLE
                  </Badge>
                {/if}
              </td>
              <td class="px-6 py-4 text-right">
                <Button
                  variant="outline"
                  size="sm"
                  class="h-8 text-xs"
                  onclick={() => onApplyMedicalLeave(item.student_id)}
                >
                  Apply Leave
                </Button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </CardContent>
</Card>
