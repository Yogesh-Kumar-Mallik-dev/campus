<script lang="ts">
  /**
   * BLOCK_WEB_HOSTEL_GATE_PASS_DESK_001
   * Subsystem: Rank 7 - Hostel Management System (hostel)
   * Purpose:   Outing and gate pass desk, warden approvals, live security gate entry/exit logging, and curfew infraction tracking.
   */
  import type { HostelGatePass, GatePassStatus, HostelIncidentLog } from '$lib/types/hostel';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let {
    gatePasses = [
      {
        id: 'gp_1',
        tenant_id: 'ten_default',
        student_id: 'stu_1',
        student_name: 'Rahul Sharma',
        roll_number: '2026CSE001',
        block_id: 'blk_1',
        block_name: 'Aryabhatta Block A',
        reason: 'Family visit over the weekend',
        destination: 'New Delhi (Home)',
        emergency_contact: '+91 98765 43210',
        expected_out_at: '2026-09-15T16:00:00Z',
        expected_in_at: '2026-09-17T20:00:00Z',
        status: 'PENDING',
        created_at: '2026-09-13T08:00:00Z'
      },
      {
        id: 'gp_2',
        tenant_id: 'ten_default',
        student_id: 'stu_2',
        student_name: 'Aditya Verma',
        roll_number: '2026CSE002',
        block_id: 'blk_1',
        block_name: 'Aryabhatta Block A',
        reason: 'Medical consultation at Apollo Hospital',
        destination: 'Apollo Hospital, Sector 62',
        emergency_contact: '+91 98112 23344',
        expected_out_at: '2026-09-13T10:00:00Z',
        expected_in_at: '2026-09-13T14:00:00Z',
        actual_out_at: '2026-09-13T10:15:00Z',
        status: 'OUT_CAMPUS',
        approved_by_id: 'usr_warden_1',
        created_at: '2026-09-12T19:00:00Z'
      },
      {
        id: 'gp_3',
        tenant_id: 'ten_default',
        student_id: 'stu_3',
        student_name: 'Vikram Mehta',
        roll_number: '2026ECE014',
        block_id: 'blk_1',
        block_name: 'Aryabhatta Block A',
        reason: 'Project meeting with industry mentor',
        destination: 'Cyber City, Gurugram',
        emergency_contact: '+91 99887 76655',
        expected_out_at: '2026-09-12T14:00:00Z',
        expected_in_at: '2026-09-12T18:00:00Z',
        actual_out_at: '2026-09-12T14:05:00Z',
        actual_in_at: '2026-09-12T19:45:00Z',
        status: 'RETURNED',
        approved_by_id: 'usr_warden_1',
        created_at: '2026-09-11T12:00:00Z'
      }
    ] as HostelGatePass[],
    incidents = [
      {
        id: 'inc_1',
        tenant_id: 'ten_default',
        student_id: 'stu_3',
        student_name: 'Vikram Mehta',
        block_id: 'blk_1',
        warden_id: 'usr_warden_1',
        warden_name: 'Prof. Ananya Roy (Chief Warden)',
        incident_type: 'CURFEW_VIOLATION',
        severity: 'MEDIUM',
        title: 'Gate Pass Curfew Breach: 105 mins overdue',
        description: 'Student returned at 19:45 exceeding the permitted 18:00 curfew on Gate Pass #gp_3.',
        action_taken: 'Fine of ₹100 charged to student billing invoice; parent notified.',
        fine_amount: 100,
        created_at: '2026-09-12T19:50:00Z'
      }
    ] as HostelIncidentLog[]
  } = $props();

  let activeTab = $state<'ALL' | 'PENDING' | 'APPROVED' | 'OUT_CAMPUS' | 'RETURNED' | 'INCIDENTS'>('PENDING');
  let searchQuery = $state<string>('');
  let isApplyModalOpen = $state<boolean>(false);

  // New Gate Pass Application Form State
  let formStudentName = $state<string>('');
  let formRollNumber = $state<string>('');
  let formReason = $state<string>('');
  let formDestination = $state<string>('');
  let formEmergencyContact = $state<string>('');
  let formOutAt = $state<string>('');
  let formInAt = $state<string>('');

  let filteredGatePasses = $derived(
    gatePasses.filter((gp) => {
      const matchTab = activeTab === 'ALL' || gp.status === activeTab;
      const matchSearch =
        searchQuery.trim() === '' ||
        gp.student_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        gp.roll_number?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        gp.destination.toLowerCase().includes(searchQuery.toLowerCase());
      return matchTab && matchSearch;
    })
  );

  function handleApprove(gp: HostelGatePass) {
    gp.status = 'APPROVED';
    gp.approved_by_id = 'usr_warden_active';
  }

  function handleReject(gp: HostelGatePass) {
    const reason = prompt('Please enter rejection reason:');
    if (reason !== null) {
      gp.status = 'REJECTED';
      gp.rejection_reason = reason;
    }
  }

  function handleRecordExit(gp: HostelGatePass) {
    gp.status = 'OUT_CAMPUS';
    gp.actual_out_at = new Date().toISOString();
  }

  function handleRecordReturn(gp: HostelGatePass) {
    const now = new Date();
    gp.actual_in_at = now.toISOString();
    gp.status = 'RETURNED';

    const expectedIn = new Date(gp.expected_in_at);
    if (now > expectedIn) {
      const diffMins = Math.round((now.getTime() - expectedIn.getTime()) / (1000 * 60));
      incidents.unshift({
        id: `inc_${Date.now()}`,
        tenant_id: gp.tenant_id,
        student_id: gp.student_id,
        student_name: gp.student_name,
        block_id: gp.block_id,
        warden_id: 'usr_warden_active',
        warden_name: 'Duty Warden',
        incident_type: 'CURFEW_VIOLATION',
        severity: 'MEDIUM',
        title: `Curfew Violation: ${diffMins} mins overdue`,
        description: `Student returned late past designated return time of ${expectedIn.toLocaleTimeString()}.`,
        action_taken: 'Fine of ₹100 charged to student ledger.',
        fine_amount: 100,
        created_at: now.toISOString()
      });
      alert(`⚠️ Curfew Infraction Logged: Student returned ${diffMins} minutes late. Fine charged.`);
    } else {
      alert(`✓ Gate Pass Returned on time.`);
    }
  }

  function handleCreateGatePass() {
    if (!formStudentName || !formRollNumber || !formReason || !formDestination || !formOutAt || !formInAt) {
      alert('Please fill all required fields');
      return;
    }

    const newGp: HostelGatePass = {
      id: `gp_${Date.now()}`,
      tenant_id: 'ten_default',
      student_id: `stu_${Date.now()}`,
      student_name: formStudentName,
      roll_number: formRollNumber,
      block_id: 'blk_1',
      block_name: 'Aryabhatta Block A',
      reason: formReason,
      destination: formDestination,
      emergency_contact: formEmergencyContact || '+91 99999 88888',
      expected_out_at: new Date(formOutAt).toISOString(),
      expected_in_at: new Date(formInAt).toISOString(),
      status: 'PENDING',
      created_at: new Date().toISOString()
    };

    gatePasses.unshift(newGp);
    isApplyModalOpen = false;
    // reset
    formStudentName = '';
    formRollNumber = '';
    formReason = '';
    formDestination = '';
    formEmergencyContact = '';
  }
</script>

<div class="space-y-6">
  <!-- Controls and Tab Pills -->
  <Card class="p-4">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div class="flex flex-wrap items-center gap-1.5 bg-muted p-1 rounded-lg">
        <Button
          size="sm"
          variant={activeTab === 'PENDING' ? 'default' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'PENDING')}
        >
          Pending Review ({gatePasses.filter((g) => g.status === 'PENDING').length})
        </Button>
        <Button
          size="sm"
          variant={activeTab === 'APPROVED' ? 'default' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'APPROVED')}
        >
          Approved ({gatePasses.filter((g) => g.status === 'APPROVED').length})
        </Button>
        <Button
          size="sm"
          variant={activeTab === 'OUT_CAMPUS' ? 'default' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'OUT_CAMPUS')}
        >
          Out of Campus ({gatePasses.filter((g) => g.status === 'OUT_CAMPUS').length})
        </Button>
        <Button
          size="sm"
          variant={activeTab === 'RETURNED' ? 'default' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'RETURNED')}
        >
          Returned ({gatePasses.filter((g) => g.status === 'RETURNED').length})
        </Button>
        <Button
          size="sm"
          variant={activeTab === 'ALL' ? 'default' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'ALL')}
        >
          All Passes ({gatePasses.length})
        </Button>
        <Button
          size="sm"
          variant={activeTab === 'INCIDENTS' ? 'destructive' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'INCIDENTS')}
        >
          Curfew / Disciplinary Logs ({incidents.length})
        </Button>
      </div>

      <div class="flex items-center gap-3">
        <Input
          type="search"
          placeholder="Filter by student, roll no, destination..."
          bind:value={searchQuery}
          class="w-64 h-9"
        />
        <Button size="sm" class="h-9 font-semibold" onclick={() => (isApplyModalOpen = true)}>
          + Apply Gate Pass
        </Button>
      </div>
    </div>
  </Card>

  <!-- Content based on Active Tab -->
  {#if activeTab === 'INCIDENTS'}
    <!-- Curfew & Disciplinary Logs -->
    <div class="space-y-4">
      {#each incidents as inc (inc.id)}
        <Card class="border-red-200 bg-red-50/20 dark:bg-red-950/10 shadow-sm">
          <CardHeader class="pb-2 flex flex-row items-center justify-between">
            <div class="flex items-center gap-2">
              <Badge variant="destructive" class="text-[10px] font-bold uppercase">
                {inc.incident_type}
              </Badge>
              <Badge variant="outline" class="text-[10px] border-red-300 text-red-800 font-semibold">
                Severity: {inc.severity}
              </Badge>
            </div>
            <span class="text-xs text-muted-foreground">
              {new Date(inc.created_at).toLocaleString()}
            </span>
          </CardHeader>

          <CardContent class="p-4 space-y-2">
            <h4 class="font-bold text-foreground text-sm">{inc.title}</h4>
            <p class="text-xs text-muted-foreground leading-relaxed">{inc.description}</p>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-2 pt-2 text-xs border-t">
              <div>
                <span class="text-muted-foreground">Student:</span>
                <span class="font-semibold text-foreground ml-1">{inc.student_name}</span>
              </div>
              <div>
                <span class="text-muted-foreground">Logged By:</span>
                <span class="font-semibold text-foreground ml-1">{inc.warden_name}</span>
              </div>
              <div>
                <span class="text-muted-foreground">Fine Charged:</span>
                <span class="font-bold text-red-600 ml-1">₹{inc.fine_amount}</span>
              </div>
            </div>

            {#if inc.action_taken}
              <div class="p-2 rounded bg-muted/60 text-[11px] text-foreground font-medium mt-1">
                <strong>Action Taken:</strong> {inc.action_taken}
              </div>
            {/if}
          </CardContent>
        </Card>
      {/each}
    </div>
  {:else}
    <!-- Gate Passes Feed -->
    <div class="space-y-4">
      {#if filteredGatePasses.length === 0}
        <Card class="p-8 text-center text-muted-foreground">
          <p class="text-sm">No gate pass applications found in this view.</p>
        </Card>
      {/if}

      {#each filteredGatePasses as pass (pass.id)}
        <Card class="shadow-sm border hover:border-primary/40 transition-colors">
          <CardHeader class="pb-3 border-b bg-muted/10">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <span class="font-bold text-base text-foreground">{pass.student_name}</span>
                <Badge variant="outline" class="font-mono text-xs">{pass.roll_number}</Badge>
                <Badge variant="secondary" class="text-[11px]">{pass.block_name}</Badge>
              </div>

              <div>
                {#if pass.status === 'PENDING'}
                  <Badge class="bg-amber-100 text-amber-800 border-amber-300">PENDING WARDEN REVIEW</Badge>
                {:else if pass.status === 'APPROVED'}
                  <Badge class="bg-blue-100 text-blue-800 border-blue-300">APPROVED · READY FOR EXIT</Badge>
                {:else if pass.status === 'OUT_CAMPUS'}
                  <Badge class="bg-purple-100 text-purple-800 border-purple-300">OUT OF CAMPUS</Badge>
                {:else if pass.status === 'RETURNED'}
                  <Badge class="bg-emerald-100 text-emerald-800 border-emerald-300">RETURNED</Badge>
                {:else if pass.status === 'REJECTED'}
                  <Badge variant="destructive">REJECTED</Badge>
                {/if}
              </div>
            </div>
          </CardHeader>

          <CardContent class="p-4 space-y-3">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
              <div class="space-y-1">
                <p><strong class="text-muted-foreground">Reason:</strong> {pass.reason}</p>
                <p><strong class="text-muted-foreground">Destination:</strong> {pass.destination}</p>
                <p><strong class="text-muted-foreground">Emergency Contact:</strong> {pass.emergency_contact}</p>
              </div>

              <div class="space-y-1">
                <p>
                  <strong class="text-muted-foreground">Expected Out:</strong>
                  {new Date(pass.expected_out_at).toLocaleString()}
                </p>
                <p>
                  <strong class="text-muted-foreground">Expected In:</strong>
                  {new Date(pass.expected_in_at).toLocaleString()}
                </p>
                {#if pass.actual_out_at}
                  <p class="text-blue-600">
                    <strong>Actual Departure:</strong> {new Date(pass.actual_out_at).toLocaleTimeString()}
                  </p>
                {/if}
                {#if pass.actual_in_at}
                  <p class="text-emerald-600">
                    <strong>Actual Re-entry:</strong> {new Date(pass.actual_in_at).toLocaleTimeString()}
                  </p>
                {/if}
              </div>
            </div>

            {#if pass.rejection_reason}
              <div class="p-2 rounded bg-destructive/10 text-destructive text-xs">
                <strong>Rejection Reason:</strong> {pass.rejection_reason}
              </div>
            {/if}
          </CardContent>

          <CardFooter class="p-3 bg-muted/20 border-t flex items-center justify-between">
            <span class="text-[11px] text-muted-foreground">
              Applied on {new Date(pass.created_at).toLocaleDateString()}
            </span>

            <div class="flex items-center gap-2">
              {#if pass.status === 'PENDING'}
                <Button
                  size="sm"
                  variant="outline"
                  class="h-8 text-xs text-destructive hover:bg-destructive/10"
                  onclick={() => handleReject(pass)}
                >
                  Reject
                </Button>
                <Button
                  size="sm"
                  class="h-8 text-xs bg-slate-900 text-white"
                  onclick={() => handleApprove(pass)}
                >
                  Approve Pass
                </Button>
              {:else if pass.status === 'APPROVED'}
                <Button
                  size="sm"
                  class="h-8 text-xs bg-purple-700 hover:bg-purple-800 text-white"
                  onclick={() => handleRecordExit(pass)}
                >
                  Security Gate: Record Exit
                </Button>
              {:else if pass.status === 'OUT_CAMPUS'}
                <Button
                  size="sm"
                  class="h-8 text-xs bg-emerald-700 hover:bg-emerald-800 text-white"
                  onclick={() => handleRecordReturn(pass)}
                >
                  Security Gate: Record Return
                </Button>
              {/if}
            </div>
          </CardFooter>
        </Card>
      {/each}
    </div>
  {/if}
</div>

<!-- Apply Gate Pass Modal -->
{#if isApplyModalOpen}
  <div class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
    <Card class="w-full max-w-lg shadow-2xl bg-card animate-in fade-in zoom-in-95">
      <CardHeader>
        <CardTitle class="text-lg">Apply for Hostel Gate Pass</CardTitle>
        <CardDescription>Submit an official outing authorization request to the warden desk.</CardDescription>
      </CardHeader>
      <CardContent class="space-y-3">
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="stu_name">Student Name</label>
            <Input id="stu_name" placeholder="Full name" bind:value={formStudentName} />
          </div>
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="stu_roll">Roll Number</label>
            <Input id="stu_roll" placeholder="Roll No" bind:value={formRollNumber} />
          </div>
        </div>

        <div class="space-y-1">
          <label class="text-xs font-semibold" for="out_reason">Purpose / Reason for Leave</label>
          <Input id="out_reason" placeholder="e.g. Family wedding / Weekend home visit" bind:value={formReason} />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="out_dest">Destination Address</label>
            <Input id="out_dest" placeholder="City / Address" bind:value={formDestination} />
          </div>
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="out_contact">Emergency Contact</label>
            <Input id="out_contact" placeholder="+91 XXXXX XXXXX" bind:value={formEmergencyContact} />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="out_dep">Expected Departure</label>
            <Input id="out_dep" type="datetime-local" bind:value={formOutAt} />
          </div>
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="out_ret">Expected Return</label>
            <Input id="out_ret" type="datetime-local" bind:value={formInAt} />
          </div>
        </div>
      </CardContent>
      <CardFooter class="flex justify-end gap-2 border-t pt-4">
        <Button variant="ghost" size="sm" onclick={() => (isApplyModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleCreateGatePass}>
          Submit Gate Pass
        </Button>
      </CardFooter>
    </Card>
  </div>
{/if}
