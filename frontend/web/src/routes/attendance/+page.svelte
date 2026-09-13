<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_ATTENDANCE_001
   * Purpose: Institutional Attendance & Shortage Management Portal for Faculty, Administrators, and Students.
   */
  import AttendanceRoster from '$lib/components/attendance/AttendanceRoster.svelte';
  import AttendanceShortageTable from '$lib/components/attendance/AttendanceShortageTable.svelte';
  import type { AttendanceSession, AttendanceRecord } from '$lib/types/attendance';
  import { Button } from '$lib/components/ui/button';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';

  let currentSession = $state<AttendanceSession>({
    id: 'sess_live_101',
    tenant_id: 'ten_default',
    subject_id: 'CS101',
    subject_name: 'Introduction to Computer Science',
    faculty_id: 'fac_smith',
    cohort_id: 'coh_2026_cs',
    session_date: new Date().toISOString().split('T')[0],
    start_time: '09:00',
    end_time: '10:00',
    mode: 'MANUAL_FACULTY',
    status: 'OPEN',
    total_students: 3,
    present_count: 2,
    absent_count: 1,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  });

  async function handleSaveRecords(records: AttendanceRecord[]) {
    // In actual app, sends POST to /api/v1/attendance/sessions/{id}/mark
    currentSession.present_count = records.filter(r => r.status === 'PRESENT' || r.status === 'LATE').length;
    currentSession.absent_count = records.filter(r => r.status === 'ABSENT').length;
    currentSession.records = records;
    alert('Attendance draft saved successfully!');
  }

  async function handleFinalize() {
    currentSession.status = 'FINALIZED';
    alert('Session attendance finalized and locked into ledger!');
  }

  function handleApplyLeave(studentID: string) {
    alert(`Initiated medical condonation workflow for student: ${studentID}`);
  }
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Attendance Command Center
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Multi-modal roll call marking, real-time biometrics sync, shortage compliance (&lt;75%), and medical condonation.
      </p>
    </div>

    <div class="flex gap-2">
      <Button variant="outline" size="sm">Schedule Session</Button>
      <Button size="sm">Biometric Terminal Sync</Button>
    </div>
  </div>

  <Tabs value="roster" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-md">
      <TabsTrigger value="roster">Live Roll Call</TabsTrigger>
      <TabsTrigger value="shortage">Shortage Tracker (&lt;75%)</TabsTrigger>
    </TabsList>

    <TabsContent value="roster" class="mt-6">
      <AttendanceRoster
        session={currentSession}
        onSaveRecords={handleSaveRecords}
        onFinalize={handleFinalize}
      />
    </TabsContent>

    <TabsContent value="shortage" class="mt-6">
      <AttendanceShortageTable
        onApplyMedicalLeave={handleApplyLeave}
      />
    </TabsContent>
  </Tabs>
</div>
