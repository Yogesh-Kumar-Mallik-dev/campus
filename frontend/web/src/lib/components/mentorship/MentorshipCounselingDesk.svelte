<script lang="ts">
  /**
   * BLOCK_WEB_MENTORSHIP_COUNSELING_001
   * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
   * Purpose:   1-on-1 counseling log desk, semester SGPA/attendance progress tracker, and at-risk alert intervention panel.
   */
  import type {
    MentorshipSession,
    StudentAcademicProgress,
    MentorshipAtRiskAlert,
    MentorshipMeetingType
  } from '$lib/types/mentorship';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Badge } from '$lib/components/ui/badge';

  let sessions = $state<MentorshipSession[]>([
    {
      id: 'sess-01',
      tenantId: 'tenant-main',
      allocationId: 'alloc-01',
      studentName: 'Aarav Sharma (CS22B001)',
      scheduledAt: '2026-09-12T14:30:00Z',
      completedAt: '2026-09-12T15:15:00Z',
      location: 'Faculty Cabin 204',
      meetingType: 'ACADEMIC_REVIEW',
      status: 'COMPLETED',
      discussionSummary: 'Analyzed Semester 4 performance. Candidate scored 8.6 SGPA and is planning honors thesis in Distributed Systems.',
      actionItems: 'Apply for summer research internship; submit paper outline by end of month.',
      followUpDate: '2026-10-15T00:00:00Z',
      createdAt: '2026-09-10T00:00:00Z',
      updatedAt: '2026-09-12T15:15:00Z'
    },
    {
      id: 'sess-02',
      tenantId: 'tenant-main',
      allocationId: 'alloc-02',
      studentName: 'Diya Patel (CS22B014)',
      scheduledAt: '2026-09-20T11:00:00Z',
      location: 'Online Google Meet',
      meetingType: 'ONE_ON_ONE',
      status: 'SCHEDULED',
      createdAt: '2026-09-14T00:00:00Z',
      updatedAt: '2026-09-14T00:00:00Z'
    }
  ]);

  let progressRecords = $state<StudentAcademicProgress[]>([
    {
      id: 'prog-01',
      tenantId: 'tenant-main',
      studentId: 'stu-001',
      studentName: 'Aarav Sharma',
      rollNumber: 'CS22B001',
      semester: 4,
      sgpa: 8.6,
      cgpa: 8.4,
      attendancePercentage: 92.5,
      creditsEarned: 22,
      totalCredits: 22,
      atRiskStatus: 'NORMAL',
      remarks: 'Dean List qualifier.',
      evaluatedAt: '2026-08-15T00:00:00Z',
      createdAt: '2026-08-15T00:00:00Z',
      updatedAt: '2026-08-15T00:00:00Z'
    },
    {
      id: 'prog-02',
      tenantId: 'tenant-main',
      studentId: 'stu-004',
      studentName: 'Kunal Joshi',
      rollNumber: 'CS22B077',
      semester: 4,
      sgpa: 3.8,
      cgpa: 4.5,
      attendancePercentage: 61.0,
      creditsEarned: 12,
      totalCredits: 22,
      atRiskStatus: 'CRITICAL_INTERVENTION',
      remarks: 'Attendance deficit in Data Structures & OS. Academic probation notice issued.',
      evaluatedAt: '2026-08-15T00:00:00Z',
      createdAt: '2026-08-15T00:00:00Z',
      updatedAt: '2026-08-15T00:00:00Z'
    }
  ]);

  let alerts = $state<MentorshipAtRiskAlert[]>([
    {
      id: 'alrt-01',
      tenantId: 'tenant-main',
      studentId: 'stu-004',
      studentName: 'Kunal Joshi',
      rollNumber: 'CS22B077',
      riskType: 'ACADEMIC_PROBATION',
      severity: 'CRITICAL',
      description: 'Automated Flag: SGPA 3.8 and Attendance 61.0% in Semester 4 requires intervention.',
      createdAt: '2026-08-15T00:00:00Z',
      updatedAt: '2026-08-15T00:00:00Z'
    }
  ]);

  // Schedule Session Modal
  let isScheduleModalOpen = $state(false);
  let newStudent = $state('Diya Patel (CS22B014)');
  let newSchedDate = $state('2026-09-22T14:00');
  let newLocation = $state('Faculty Cabin 204');
  let newMeetingType = $state<MentorshipMeetingType>('ONE_ON_ONE');

  // Complete Session Modal
  let isCompleteModalOpen = $state(false);
  let activeCompleteSess = $state<MentorshipSession | null>(null);
  let completeSummary = $state('');
  let completeActionItems = $state('');

  function handleScheduleSession() {
    const newSess: MentorshipSession = {
      id: `sess-${Date.now()}`,
      tenantId: 'tenant-main',
      allocationId: 'alloc-01',
      studentName: newStudent,
      scheduledAt: new Date(newSchedDate).toISOString(),
      location: newLocation,
      meetingType: newMeetingType,
      status: 'SCHEDULED',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    sessions = [newSess, ...sessions];
    isScheduleModalOpen = false;
  }

  function openCompleteSession(sess: MentorshipSession) {
    activeCompleteSess = sess;
    completeSummary = sess.discussionSummary ?? '';
    completeActionItems = sess.actionItems ?? '';
    isCompleteModalOpen = true;
  }

  function handleSaveComplete() {
    if (!activeCompleteSess) return;
    sessions = sessions.map((s) => {
      if (s.id === activeCompleteSess!.id) {
        return {
          ...s,
          status: 'COMPLETED',
          discussionSummary: completeSummary,
          actionItems: completeActionItems,
          completedAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        };
      }
      return s;
    });
    isCompleteModalOpen = false;
    activeCompleteSess = null;
  }

  function handleResolveAlert(alertId: string) {
    alerts = alerts.map((a) => {
      if (a.id === alertId) {
        return {
          ...a,
          resolvedAt: new Date().toISOString(),
          resolvedById: 'staff-dean-01',
          updatedAt: new Date().toISOString()
        };
      }
      return a;
    });
  }

  function getAtRiskBadge(status: string) {
    switch (status) {
      case 'NORMAL':
        return 'bg-emerald-100 text-emerald-800 border-emerald-200';
      case 'WATCHLIST':
        return 'bg-amber-100 text-amber-800 border-amber-200';
      case 'CRITICAL_INTERVENTION':
        return 'bg-rose-100 text-rose-800 border-rose-200';
      default:
        return 'bg-slate-100 text-slate-800 border-slate-200';
    }
  }
</script>

<div class="space-y-6">
  <!-- Active Risk Intervention Banner if any alerts -->
  {#each alerts.filter((a) => !a.resolvedAt) as alert}
    <div class="bg-rose-50 border border-rose-200 rounded-2xl p-5 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="space-y-1">
        <div class="flex items-center gap-2">
          <Badge variant="destructive">CRITICAL AT-RISK ALERT</Badge>
          <span class="text-xs font-bold text-rose-900 font-mono">{alert.studentName} ({alert.rollNumber})</span>
        </div>
        <p class="text-xs text-rose-800 leading-relaxed">
          {alert.description}
        </p>
      </div>
      <Button size="sm" class="bg-rose-700 hover:bg-rose-800 text-white shrink-0" onclick={() => handleResolveAlert(alert.id)}>
        ✓ Acknowledge & Resolve Case
      </Button>
    </div>
  {/each}

  <!-- Sessions Header & Actions -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-card p-4 rounded-xl border">
    <div>
      <h3 class="text-sm font-bold text-foreground">Mentorship & Counseling Sessions</h3>
      <p class="text-xs text-muted-foreground">Log 1-on-1 counseling notes, remedial academic action items, and review meetings.</p>
    </div>
    <Button size="sm" onclick={() => (isScheduleModalOpen = true)}>
      📅 Schedule Session
    </Button>
  </div>

  <!-- Sessions List -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    {#each sessions as sess}
      <div class="bg-card border rounded-2xl p-5 space-y-3 shadow-sm flex flex-col justify-between">
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <Badge variant={sess.status === 'COMPLETED' ? 'default' : 'secondary'}>
              {sess.status}
            </Badge>
            <span class="text-[11px] font-semibold text-muted-foreground">
              {sess.meetingType.replace('_', ' ')}
            </span>
          </div>

          <h4 class="text-sm font-bold text-foreground">{sess.studentName}</h4>
          <p class="text-xs text-muted-foreground">
            📍 {sess.location ?? 'Online'} · ⏰ {new Date(sess.scheduledAt).toLocaleString()}
          </p>

          {#if sess.discussionSummary}
            <div class="bg-muted/40 p-3 rounded-xl border text-xs text-foreground space-y-1 mt-2">
              <span class="font-bold text-muted-foreground text-[10px] uppercase block">Discussion Notes:</span>
              <p class="text-xs leading-relaxed">{sess.discussionSummary}</p>
              {#if sess.actionItems}
                <div class="mt-2 pt-2 border-t text-[11px] text-muted-foreground">
                  <span class="font-bold text-foreground">Action Items:</span> {sess.actionItems}
                </div>
              {/if}
            </div>
          {/if}
        </div>

        <div class="pt-3 border-t flex items-center justify-end">
          {#if sess.status === 'SCHEDULED'}
            <Button size="sm" variant="outline" onclick={() => openCompleteSession(sess)}>
              ✓ Complete & Log Notes
            </Button>
          {:else}
            <span class="text-[11px] text-muted-foreground italic">Session Completed</span>
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <!-- Academic Performance Progress Matrix -->
  <div class="bg-card border rounded-2xl overflow-hidden shadow-sm space-y-0 mt-6">
    <div class="p-4 border-b">
      <h3 class="text-sm font-bold text-foreground">Semester Academic Performance & Risk Flags</h3>
      <p class="text-xs text-muted-foreground">Automated surveillance of GPA metrics and biometric attendance thresholds.</p>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead class="bg-muted/50 text-muted-foreground font-semibold uppercase tracking-wider border-b">
          <tr>
            <th class="p-3">Student</th>
            <th class="p-3">Semester</th>
            <th class="p-3">SGPA / CGPA</th>
            <th class="p-3">Attendance</th>
            <th class="p-3">Credits</th>
            <th class="p-3">At-Risk Status</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          {#each progressRecords as prog}
            <tr class="hover:bg-muted/30 transition-colors">
              <td class="p-3">
                <div class="font-bold text-foreground">{prog.studentName}</div>
                <div class="text-[10px] text-muted-foreground font-mono">{prog.rollNumber}</div>
              </td>
              <td class="p-3 font-semibold">Semester {prog.semester}</td>
              <td class="p-3">
                <span class="font-bold {prog.sgpa < 5.0 ? 'text-rose-600' : 'text-foreground'}">
                  {prog.sgpa.toFixed(2)}
                </span>
                <span class="text-[10px] text-muted-foreground block">CGPA: {prog.cgpa.toFixed(2)}</span>
              </td>
              <td class="p-3">
                <span class="font-bold {prog.attendancePercentage < 75.0 ? 'text-rose-600' : 'text-emerald-700'}">
                  {prog.attendancePercentage.toFixed(1)}%
                </span>
              </td>
              <td class="p-3 text-muted-foreground">
                {prog.creditsEarned} / {prog.totalCredits}
              </td>
              <td class="p-3">
                <span class="px-2 py-0.5 rounded text-[10px] font-bold border {getAtRiskBadge(prog.atRiskStatus)}">
                  {prog.atRiskStatus.replace('_', ' ')}
                </span>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>

<!-- Modal: Schedule Session -->
{#if isScheduleModalOpen}
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
    <div class="bg-card border max-w-lg w-full rounded-2xl p-6 shadow-2xl space-y-4">
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-lg font-bold text-foreground">Schedule Mentorship Session</h3>
        <button
          onclick={() => (isScheduleModalOpen = false)}
          class="text-muted-foreground hover:text-foreground text-sm font-bold"
        >
          ✕
        </button>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <Label for="sessStudent" class="text-xs font-semibold">Mentee Student</Label>
          <Input id="sessStudent" bind:value={newStudent} />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <Label for="sessDate" class="text-xs font-semibold">Date & Time</Label>
            <Input id="sessDate" type="datetime-local" bind:value={newSchedDate} />
          </div>

          <div class="space-y-1">
            <Label for="sessType" class="text-xs font-semibold">Meeting Type</Label>
            <select
              id="sessType"
              bind:value={newMeetingType}
              class="w-full bg-background border rounded-lg px-3 py-2 text-xs font-medium"
            >
              <option value="ONE_ON_ONE">1-on-1 Counseling</option>
              <option value="ACADEMIC_REVIEW">Academic Review</option>
              <option value="GROUP">Group Mentoring</option>
              <option value="EMERGENCY_COUNSELING">Emergency Intervention</option>
            </select>
          </div>
        </div>

        <div class="space-y-1">
          <Label for="sessLoc" class="text-xs font-semibold">Meeting Location / Room</Label>
          <Input id="sessLoc" bind:value={newLocation} />
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-3 border-t">
        <Button variant="outline" size="sm" onclick={() => (isScheduleModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleScheduleSession}>
          Confirm Schedule
        </Button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal: Complete Session -->
{#if isCompleteModalOpen && activeCompleteSess}
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
    <div class="bg-card border max-w-lg w-full rounded-2xl p-6 shadow-2xl space-y-4">
      <div class="flex items-center justify-between border-b pb-3">
        <div>
          <h3 class="text-lg font-bold text-foreground">Log Session Discussion</h3>
          <p class="text-xs text-muted-foreground">Mentee: {activeCompleteSess.studentName}</p>
        </div>
        <button
          onclick={() => (isCompleteModalOpen = false)}
          class="text-muted-foreground hover:text-foreground text-sm font-bold"
        >
          ✕
        </button>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <Label for="discSummary" class="text-xs font-semibold">Discussion Summary & Observations</Label>
          <textarea
            id="discSummary"
            bind:value={completeSummary}
            rows="3"
            class="w-full bg-background border rounded-lg p-2.5 text-xs focus:ring-2 focus:ring-primary"
            placeholder="Key topics discussed, student concerns, mental well-being..."
          ></textarea>
        </div>

        <div class="space-y-1">
          <Label for="actItems" class="text-xs font-semibold">Agreed Action Items & Goals</Label>
          <textarea
            id="actItems"
            bind:value={completeActionItems}
            rows="2"
            class="w-full bg-background border rounded-lg p-2.5 text-xs focus:ring-2 focus:ring-primary"
            placeholder="Target study hours, tutoring schedules, project deadlines..."
          ></textarea>
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-3 border-t">
        <Button variant="outline" size="sm" onclick={() => (isCompleteModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleSaveComplete}>
          Save & Complete Session
        </Button>
      </div>
    </div>
  </div>
{/if}
