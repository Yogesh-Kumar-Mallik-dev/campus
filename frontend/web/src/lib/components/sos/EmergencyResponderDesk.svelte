<!--
  BLOCK_WEB_SOS_RESPONDER_DESK_001
  Subsystem: Rank 14 - SOS & Emergency Response (sos)
  Purpose:   Emergency Command Desk, responder dispatch management, live arrival checkpoints, and incident resolution logging.
-->
<script lang="ts">
  import type { SOSIncident, SOSDispatchResponder, SOSResponderRole, SOSResponderStatus } from '$lib/types/sos';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let { incident, onBack }: { incident?: SOSIncident; onBack?: () => void } = $props();

  let currentIncident = $state<SOSIncident>(
    incident || {
      id: 'sos-001',
      tenantId: 'tenant-demo',
      alertNumber: 'SOS-2026-00012',
      userId: 'stu-danger-1',
      userName: 'Kavya Menon (2024ECE034)',
      userPhone: '+91 98765 43210',
      emergencyType: 'MEDICAL',
      latitude: 12.9718,
      longitude: 77.5948,
      locationDescription: 'Girls Hostel Block A - 2nd Floor Common Hall',
      status: 'DISPATCHED',
      triggeredAt: new Date(Date.now() - 7 * 60000).toISOString(),
      acknowledgedAt: new Date(Date.now() - 6 * 60000).toISOString(),
      createdAt: new Date(Date.now() - 7 * 60000).toISOString(),
      updatedAt: new Date().toISOString(),
    }
  );

  let responders = $state<SOSDispatchResponder[]>([
    {
      id: 'resp-1',
      tenantId: 'tenant-demo',
      incidentId: currentIncident.id,
      responderId: 'staff-paramedic-1',
      responderName: 'Dr. Anita Joshi (Campus Paramedic Unit 1)',
      role: 'PARAMEDIC',
      status: 'EN_ROUTE',
      dispatchedAt: new Date(Date.now() - 5 * 60000).toISOString(),
      createdAt: new Date(Date.now() - 5 * 60000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'resp-2',
      tenantId: 'tenant-demo',
      incidentId: currentIncident.id,
      responderId: 'staff-warden-1',
      responderName: 'Mrs. Sunita Rao (Hostel Senior Warden)',
      role: 'HOSTEL_WARDEN',
      status: 'ON_SCENE',
      dispatchedAt: new Date(Date.now() - 5 * 60000).toISOString(),
      arrivedAt: new Date(Date.now() - 2 * 60000).toISOString(),
      createdAt: new Date(Date.now() - 5 * 60000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ]);

  // Dispatch modal
  let isDispatchModalOpen = $state<boolean>(false);
  let newResponderName = $state<string>('');
  let newResponderRole = $state<SOSResponderRole>('CAMPUS_SECURITY');

  // Resolution modal
  let isResolveModalOpen = $state<boolean>(false);
  let resolutionNotes = $state<string>('');
  let isFalseAlarm = $state<boolean>(false);

  function handleAcknowledge() {
    currentIncident.status = 'ACKNOWLEDGED';
    currentIncident.acknowledgedAt = new Date().toISOString();
  }

  function handleDispatch() {
    if (!newResponderName.trim()) return;

    const newResp: SOSDispatchResponder = {
      id: `resp-${Date.now()}`,
      tenantId: currentIncident.tenantId,
      incidentId: currentIncident.id,
      responderId: `resp-user-${Date.now()}`,
      responderName: newResponderName.trim(),
      role: newResponderRole,
      status: 'ASSIGNED',
      dispatchedAt: new Date().toISOString(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    responders = [...responders, newResp];
    if (currentIncident.status === 'TRIGGERED' || currentIncident.status === 'ACKNOWLEDGED') {
      currentIncident.status = 'DISPATCHED';
    }
    newResponderName = '';
    isDispatchModalOpen = false;
  }

  function handleUpdateResponderStatus(respId: string, nextStatus: SOSResponderStatus) {
    responders = responders.map((r) => {
      if (r.id === respId) {
        const updated = { ...r, status: nextStatus, updatedAt: new Date().toISOString() };
        if (nextStatus === 'ON_SCENE' && !updated.arrivedAt) {
          updated.arrivedAt = new Date().toISOString();
          currentIncident.status = 'ON_SCENE';
        }
        return updated;
      }
      return r;
    });
  }

  function handleResolve() {
    if (!resolutionNotes.trim()) return;

    currentIncident.status = isFalseAlarm ? 'FALSE_ALARM' : 'RESOLVED';
    currentIncident.resolvedAt = new Date().toISOString();
    currentIncident.resolvedByName = 'Control Room Supervisor';
    currentIncident.resolutionNotes = resolutionNotes.trim();
    isResolveModalOpen = false;
  }
</script>

<div class="space-y-6">
  <!-- Incident Overview Banner -->
  <div class="p-5 bg-card rounded-2xl border shadow-sm space-y-4">
    <div class="flex flex-wrap items-center justify-between gap-4 border-b pb-4">
      <div class="flex items-center gap-3">
        {#if onBack}
          <Button variant="ghost" size="sm" onclick={onBack}>
            ← Back to Radar
          </Button>
        {/if}
        <div>
          <div class="flex items-center gap-2">
            <span class="font-mono text-xs font-black text-rose-600 bg-rose-50 px-2.5 py-0.5 rounded border border-rose-200">
              {currentIncident.alertNumber}
            </span>
            <Badge variant="destructive">{currentIncident.emergencyType}</Badge>
            <Badge variant="default">{currentIncident.status}</Badge>
          </div>
          <h2 class="text-xl font-bold text-foreground mt-1">
            📍 {currentIncident.locationDescription}
          </h2>
        </div>
      </div>

      <!-- Action Commands -->
      <div class="flex flex-wrap items-center gap-2">
        {#if currentIncident.status === 'TRIGGERED'}
          <Button size="sm" class="bg-amber-600 hover:bg-amber-700 text-white font-bold" onclick={handleAcknowledge}>
            ✓ Acknowledge Alert
          </Button>
        {/if}
        {#if currentIncident.status !== 'RESOLVED' && currentIncident.status !== 'FALSE_ALARM'}
          <Button size="sm" class="bg-blue-600 hover:bg-blue-700 text-white font-bold" onclick={() => (isDispatchModalOpen = true)}>
            + Dispatch Responder
          </Button>
          <Button size="sm" class="bg-emerald-600 hover:bg-emerald-700 text-white font-bold" onclick={() => (isResolveModalOpen = true)}>
            ✓ Resolve Incident
          </Button>
        {/if}
      </div>
    </div>

    <!-- Metadata Grid -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs">
      <div>
        <span class="text-muted-foreground block">Caller / Student:</span>
        <span class="font-bold text-foreground">{currentIncident.userName || 'Student'}</span>
      </div>
      <div>
        <span class="text-muted-foreground block">Phone Contact:</span>
        <span class="font-bold text-foreground">{currentIncident.userPhone || 'N/A'}</span>
      </div>
      <div>
        <span class="text-muted-foreground block">GPS Coordinates:</span>
        <span class="font-mono font-bold text-foreground">{currentIncident.latitude.toFixed(4)}°, {currentIncident.longitude.toFixed(4)}°</span>
      </div>
      <div>
        <span class="text-muted-foreground block">Triggered Timestamp:</span>
        <span class="font-bold text-foreground">{new Date(currentIncident.triggeredAt).toLocaleString()}</span>
      </div>
    </div>

    {#if currentIncident.resolutionNotes}
      <div class="p-3 bg-emerald-50 text-emerald-900 border border-emerald-200 rounded-xl text-xs space-y-1">
        <span class="font-bold block">✓ Incident Final Resolution Record:</span>
        <p>{currentIncident.resolutionNotes}</p>
      </div>
    {/if}
  </div>

  <!-- Dispatched Responders Roster -->
  <div class="p-5 bg-card rounded-2xl border shadow-sm space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-bold uppercase tracking-wider text-muted-foreground">
        Emergency Responder Teams ({responders.length})
      </h3>
      {#if currentIncident.status !== 'RESOLVED' && currentIncident.status !== 'FALSE_ALARM'}
        <Button size="sm" variant="outline" onclick={() => (isDispatchModalOpen = true)} class="text-xs font-semibold">
          + Add Team
        </Button>
      {/if}
    </div>

    <div class="space-y-3">
      {#if responders.length === 0}
        <div class="p-6 text-center text-xs text-muted-foreground border border-dashed rounded-xl">
          No response teams dispatched yet. Click "+ Dispatch Responder" above.
        </div>
      {:else}
        {#each responders as resp}
          <div class="p-4 bg-muted/40 rounded-xl border flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div class="space-y-1">
              <div class="flex items-center gap-2">
                <span class="font-bold text-sm text-foreground">{resp.responderName}</span>
                <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-primary/10 text-primary">
                  {resp.role}
                </span>
                <Badge variant={resp.status === 'ON_SCENE' ? 'default' : 'outline'}>
                  {resp.status}
                </Badge>
              </div>
              <div class="text-xs text-muted-foreground flex gap-4">
                <span>Dispatched: {new Date(resp.dispatchedAt).toLocaleTimeString()}</span>
                {#if resp.arrivedAt}
                  <span class="text-emerald-600 font-semibold">Arrived on scene: {new Date(resp.arrivedAt).toLocaleTimeString()}</span>
                {/if}
              </div>
            </div>

            {#if currentIncident.status !== 'RESOLVED' && currentIncident.status !== 'FALSE_ALARM'}
              <div class="flex items-center gap-2">
                {#if resp.status === 'ASSIGNED'}
                  <Button size="sm" variant="outline" onclick={() => handleUpdateResponderStatus(resp.id, 'EN_ROUTE')}>
                    Mark En Route
                  </Button>
                {/if}
                {#if resp.status === 'EN_ROUTE'}
                  <Button size="sm" class="bg-emerald-600 text-white font-bold" onclick={() => handleUpdateResponderStatus(resp.id, 'ON_SCENE')}>
                    Mark On Scene
                  </Button>
                {/if}
                {#if resp.status === 'ON_SCENE'}
                  <Button size="sm" variant="secondary" onclick={() => handleUpdateResponderStatus(resp.id, 'COMPLETED')}>
                    Complete Mission
                  </Button>
                {/if}
              </div>
            {/if}
          </div>
        {/each}
      {/if}
    </div>
  </div>

  <!-- Dispatch Modal -->
  {#if isDispatchModalOpen}
    <div class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-card w-full max-w-md rounded-2xl border shadow-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b pb-3">
          <h2 class="text-lg font-bold text-foreground">Dispatch Emergency Responder</h2>
          <button onclick={() => (isDispatchModalOpen = false)} class="text-muted-foreground hover:text-foreground">✕</button>
        </div>

        <div class="space-y-3">
          <div>
            <label for="resp-role" class="text-xs font-semibold text-muted-foreground block mb-1">Responder Unit Role</label>
            <select
              id="resp-role"
              bind:value={newResponderRole}
              class="w-full bg-background text-foreground text-sm rounded-lg border px-3 py-2"
            >
              <option value="PARAMEDIC">Paramedic / Ambulance Team</option>
              <option value="CAMPUS_SECURITY">Campus Security Patrol</option>
              <option value="HOSTEL_WARDEN">Hostel Senior Warden</option>
              <option value="FIRE_SAFETY_OFFICER">Fire Safety Unit</option>
              <option value="POLICE_LIAISON">Police Liaison Officer</option>
            </select>
          </div>

          <div>
            <label for="resp-name" class="text-xs font-semibold text-muted-foreground block mb-1">Responder Name / Vehicle ID</label>
            <Input id="resp-name" bind:value={newResponderName} placeholder="e.g. Paramedic Unit 2 (Nurse Sharma)" />
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-3 border-t">
          <Button variant="outline" onclick={() => (isDispatchModalOpen = false)}>Cancel</Button>
          <Button onclick={handleDispatch} class="bg-blue-600 hover:bg-blue-700 text-white font-bold">
            Confirm Dispatch
          </Button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Resolution Modal -->
  {#if isResolveModalOpen}
    <div class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-card w-full max-w-md rounded-2xl border shadow-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b pb-3">
          <h2 class="text-lg font-bold text-foreground">Resolve Emergency Incident</h2>
          <button onclick={() => (isResolveModalOpen = false)} class="text-muted-foreground hover:text-foreground">✕</button>
        </div>

        <div class="space-y-3">
          <label class="flex items-center gap-2 cursor-pointer text-xs font-semibold text-amber-800 bg-amber-100/70 px-3 py-2 rounded-lg">
            <input type="checkbox" bind:checked={isFalseAlarm} class="rounded" />
            <span>Mark as False Alarm (Test / Accidental Trigger)</span>
          </label>

          <div>
            <label for="res-notes" class="text-xs font-semibold text-muted-foreground block mb-1">Resolution Summary & Actions Taken</label>
            <textarea
              id="res-notes"
              bind:value={resolutionNotes}
              rows={4}
              placeholder="Describe actions taken by responders, student condition, and handoff..."
              class="w-full bg-background text-foreground text-xs rounded-xl border p-3"
            ></textarea>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-3 border-t">
          <Button variant="outline" onclick={() => (isResolveModalOpen = false)}>Cancel</Button>
          <Button onclick={handleResolve} class="bg-emerald-600 hover:bg-emerald-700 text-white font-bold">
            Complete Incident Resolution
          </Button>
        </div>
      </div>
    </div>
  {/if}
</div>
