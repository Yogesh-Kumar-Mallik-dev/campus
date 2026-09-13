<script lang="ts">
  /**
   * BLOCK_WEB_MENTORSHIP_ALLOCATION_001
   * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
   * Purpose:   Faculty mentor-mentee allocation matrix, load balancing, and active relationship lifecycle management.
   */
  import type { MentorAllocation, MentorAllocationStatus } from '$lib/types/mentorship';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Badge } from '$lib/components/ui/badge';

  let allocations = $state<MentorAllocation[]>([
    {
      id: 'alloc-01',
      tenantId: 'tenant-main',
      studentId: 'stu-001',
      studentName: 'Aarav Sharma',
      rollNumber: 'CS22B001',
      mentorStaffId: 'staff-01',
      mentorName: 'Dr. Suresh Rao',
      mentorDesignation: 'Professor (CSE)',
      cohortId: '2022-2026-CSE',
      status: 'ACTIVE',
      allocatedAt: '2026-08-01T09:00:00Z',
      createdAt: '2026-08-01T09:00:00Z',
      updatedAt: '2026-08-01T09:00:00Z'
    },
    {
      id: 'alloc-02',
      tenantId: 'tenant-main',
      studentId: 'stu-002',
      studentName: 'Diya Patel',
      rollNumber: 'CS22B014',
      mentorStaffId: 'staff-01',
      mentorName: 'Dr. Suresh Rao',
      mentorDesignation: 'Professor (CSE)',
      cohortId: '2022-2026-CSE',
      status: 'ACTIVE',
      allocatedAt: '2026-08-01T09:00:00Z',
      createdAt: '2026-08-01T09:00:00Z',
      updatedAt: '2026-08-01T09:00:00Z'
    },
    {
      id: 'alloc-03',
      tenantId: 'tenant-main',
      studentId: 'stu-003',
      studentName: 'Rohan Gupta',
      rollNumber: 'CS22B045',
      mentorStaffId: 'staff-02',
      mentorName: 'Prof. Ananya Verma',
      mentorDesignation: 'Associate Professor (CSE)',
      cohortId: '2022-2026-CSE',
      status: 'ACTIVE',
      allocatedAt: '2026-08-01T09:00:00Z',
      createdAt: '2026-08-01T09:00:00Z',
      updatedAt: '2026-08-01T09:00:00Z'
    }
  ]);

  let filterStatus = $state<MentorAllocationStatus | 'ALL'>('ALL');
  let searchQuery = $state('');

  // New Allocation Modal
  let isModalOpen = $state(false);
  let newStudentName = $state('');
  let newRollNumber = $state('');
  let newMentorName = $state('Dr. Suresh Rao');
  let newCohort = $state('2022-2026-CSE');

  const filteredAllocations = $derived(
    allocations.filter((a) => {
      const matchStatus = filterStatus === 'ALL' || a.status === filterStatus;
      const matchQuery =
        (a.studentName?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false) ||
        (a.rollNumber?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false) ||
        (a.mentorName?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false);
      return matchStatus && matchQuery;
    })
  );

  function handleCreateAllocation() {
    if (!newStudentName || !newRollNumber) return;
    const newAlloc: MentorAllocation = {
      id: `alloc-${Date.now()}`,
      tenantId: 'tenant-main',
      studentId: `stu-${Date.now()}`,
      studentName: newStudentName,
      rollNumber: newRollNumber,
      mentorStaffId: 'staff-01',
      mentorName: newMentorName,
      mentorDesignation: 'Faculty Mentor',
      cohortId: newCohort,
      status: 'ACTIVE',
      allocatedAt: new Date().toISOString(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    allocations = [newAlloc, ...allocations];
    newStudentName = '';
    newRollNumber = '';
    isModalOpen = false;
  }
</script>

<div class="space-y-6">
  <!-- Top Stat Cards -->
  <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
    <div class="bg-card border rounded-2xl p-5 shadow-sm space-y-1">
      <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Active Mentorships</span>
      <div class="text-2xl font-black text-foreground">
        {allocations.filter((a) => a.status === 'ACTIVE').length}
      </div>
      <p class="text-[11px] text-emerald-700 font-semibold">100% cohorts assigned faculty guides</p>
    </div>

    <div class="bg-card border rounded-2xl p-5 shadow-sm space-y-1">
      <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Faculty Mentors</span>
      <div class="text-2xl font-black text-foreground">18</div>
      <p class="text-[11px] text-muted-foreground">Average ~15 mentees / faculty</p>
    </div>

    <div class="bg-card border rounded-2xl p-5 shadow-sm space-y-1">
      <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Cohort Coverage</span>
      <div class="text-2xl font-black text-foreground">4 Batches</div>
      <p class="text-[11px] text-muted-foreground">B.Tech, M.Tech, Ph.D cohorts</p>
    </div>
  </div>

  <!-- Search & Action Bar -->
  <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 bg-card p-4 rounded-xl border">
    <div class="flex items-center gap-3 w-full sm:w-auto">
      <Input
        placeholder="Search student, roll number, or mentor..."
        bind:value={searchQuery}
        class="text-xs w-full sm:w-72"
      />
      <select
        bind:value={filterStatus}
        class="bg-background border rounded-lg px-3 py-2 text-xs font-medium"
      >
        <option value="ALL">All Statuses</option>
        <option value="ACTIVE">Active</option>
        <option value="COMPLETED">Completed</option>
        <option value="REASSIGNED">Reassigned</option>
      </select>
    </div>

    <Button size="sm" onclick={() => (isModalOpen = true)}>
      + Allocate Mentor
    </Button>
  </div>

  <!-- Allocations Table -->
  <div class="bg-card border rounded-2xl overflow-hidden shadow-sm">
    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead class="bg-muted/50 text-muted-foreground font-semibold uppercase tracking-wider border-b">
          <tr>
            <th class="p-3">Student Mentee</th>
            <th class="p-3">Assigned Mentor</th>
            <th class="p-3">Cohort</th>
            <th class="p-3">Status</th>
            <th class="p-3">Allocated On</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          {#each filteredAllocations as a}
            <tr class="hover:bg-muted/30 transition-colors">
              <td class="p-3">
                <div class="font-bold text-foreground">{a.studentName}</div>
                <div class="text-[10px] text-muted-foreground font-mono">{a.rollNumber}</div>
              </td>
              <td class="p-3">
                <div class="font-semibold text-foreground">{a.mentorName}</div>
                <div class="text-[10px] text-muted-foreground">{a.mentorDesignation}</div>
              </td>
              <td class="p-3 font-mono text-muted-foreground">
                {a.cohortId ?? 'General'}
              </td>
              <td class="p-3">
                <Badge variant={a.status === 'ACTIVE' ? 'default' : 'secondary'}>
                  {a.status}
                </Badge>
              </td>
              <td class="p-3 text-muted-foreground">
                {new Date(a.allocatedAt).toLocaleDateString()}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>

<!-- Modal: Allocate Mentor -->
{#if isModalOpen}
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
    <div class="bg-card border max-w-lg w-full rounded-2xl p-6 shadow-2xl space-y-4">
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-lg font-bold text-foreground">Allocate Faculty Mentor</h3>
        <button
          onclick={() => (isModalOpen = false)}
          class="text-muted-foreground hover:text-foreground text-sm font-bold"
        >
          ✕
        </button>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <Label for="stuName" class="text-xs font-semibold">Student Full Name</Label>
          <Input id="stuName" bind:value={newStudentName} placeholder="e.g. Maya Krishnan" />
        </div>

        <div class="space-y-1">
          <Label for="stuRoll" class="text-xs font-semibold">Roll Number / Registration ID</Label>
          <Input id="stuRoll" bind:value={newRollNumber} placeholder="e.g. CS22B088" />
        </div>

        <div class="space-y-1">
          <Label for="mentorSelect" class="text-xs font-semibold">Select Faculty Mentor</Label>
          <select
            id="mentorSelect"
            bind:value={newMentorName}
            class="w-full bg-background border rounded-lg px-3 py-2 text-xs font-medium"
          >
            <option value="Dr. Suresh Rao">Dr. Suresh Rao (Prof, CSE)</option>
            <option value="Prof. Ananya Verma">Prof. Ananya Verma (Assoc Prof, CSE)</option>
            <option value="Dr. Vikram Sen">Dr. Vikram Sen (HOD, Electronics)</option>
          </select>
        </div>

        <div class="space-y-1">
          <Label for="cohortInput" class="text-xs font-semibold">Cohort Batch</Label>
          <Input id="cohortInput" bind:value={newCohort} />
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-3 border-t">
        <Button variant="outline" size="sm" onclick={() => (isModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleCreateAllocation}>
          Confirm Allocation
        </Button>
      </div>
    </div>
  </div>
{/if}
