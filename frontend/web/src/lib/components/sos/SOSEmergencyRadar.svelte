<!--
  BLOCK_WEB_SOS_RADAR_001
  Subsystem: Rank 14 - SOS & Emergency Response (sos)
  Purpose:   Real-time Emergency Radar, live telemetry surveillance, active geo-location alert feed, and instant trigger console.
-->
<script lang="ts">
  import type { SOSIncident, SOSEmergencyType } from '$lib/types/sos';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let { onSelectIncident }: { onSelectIncident?: (incident: SOSIncident) => void } = $props();

  let incidents = $state<SOSIncident[]>([
    {
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
    },
    {
      id: 'sos-002',
      tenantId: 'tenant-demo',
      alertNumber: 'SOS-2026-00013',
      userId: 'stu-danger-2',
      userName: 'Rohan Deshmukh (2023CSE112)',
      userPhone: '+91 91234 56789',
      emergencyType: 'FIRE',
      latitude: 12.9732,
      longitude: 77.5961,
      locationDescription: 'Physics Department Optic Lab 102',
      status: 'ON_SCENE',
      triggeredAt: new Date(Date.now() - 14 * 60000).toISOString(),
      acknowledgedAt: new Date(Date.now() - 13 * 60000).toISOString(),
      createdAt: new Date(Date.now() - 14 * 60000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'sos-003',
      tenantId: 'tenant-demo',
      alertNumber: 'SOS-2026-00011',
      userId: 'stu-danger-3',
      userName: 'Rahul Verma (2022MECH055)',
      userPhone: '+91 98450 11223',
      emergencyType: 'SECURITY_THREAT',
      latitude: 12.9705,
      longitude: 77.5932,
      locationDescription: 'South Campus Gate 3 Outer Perimeter',
      status: 'RESOLVED',
      triggeredAt: new Date(Date.now() - 45 * 60000).toISOString(),
      acknowledgedAt: new Date(Date.now() - 44 * 60000).toISOString(),
      resolvedAt: new Date(Date.now() - 20 * 60000).toISOString(),
      resolvedById: 'sec-officer-1',
      resolvedByName: 'Chief Security Officer (Unit 4)',
      resolutionNotes: 'Unauthorized trespasser intercepted and escorted off campus premises.',
      createdAt: new Date(Date.now() - 45 * 60000).toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ]);

  // Trigger modal state
  let isTriggerModalOpen = $state<boolean>(false);
  let triggerType = $state<SOSEmergencyType>('MEDICAL');
  let triggerLocation = $state<string>('');
  let triggerLat = $state<number>(12.9716);
  let triggerLng = $state<number>(77.5946);

  function handleTriggerEmergency() {
    if (!triggerLocation.trim()) return;

    const newSeq = incidents.length + 14;
    const now = new Date();
    const newInc: SOSIncident = {
      id: `sos-${Date.now()}`,
      tenantId: 'tenant-demo',
      alertNumber: `SOS-2026-${String(newSeq).padStart(5, '0')}`,
      userId: 'current-user',
      userName: 'Yogesh K (1-Tap SOS)',
      userPhone: '+91 99000 11111',
      emergencyType: triggerType,
      latitude: triggerLat,
      longitude: triggerLng,
      locationDescription: triggerLocation.trim(),
      status: 'TRIGGERED',
      triggeredAt: now.toISOString(),
      createdAt: now.toISOString(),
      updatedAt: now.toISOString(),
    };

    incidents = [newInc, ...incidents];
    triggerLocation = '';
    isTriggerModalOpen = false;
  }

  function getEmergencyTypeBadge(type: SOSEmergencyType): { color: string; icon: string } {
    switch (type) {
      case 'MEDICAL':
        return { color: 'bg-rose-500 text-white', icon: '🚑 MEDICAL EMERGENCY' };
      case 'FIRE':
        return { color: 'bg-orange-500 text-white', icon: '🔥 FIRE HAZARD' };
      case 'SECURITY_THREAT':
        return { color: 'bg-red-700 text-white', icon: '🛡️ SECURITY THREAT' };
      case 'HARASSMENT_RAGGING':
        return { color: 'bg-purple-600 text-white', icon: '⚠️ RAGGING / HARASSMENT' };
      case 'NATURAL_HAZARD':
        return { color: 'bg-amber-600 text-white', icon: '🌪️ NATURAL HAZARD' };
      default:
        return { color: 'bg-slate-600 text-white', icon: '🚨 OTHER INCIDENT' };
    }
  }

  function getStatusBadge(status: string): { variant: 'destructive' | 'default' | 'outline' | 'secondary'; label: string } {
    switch (status) {
      case 'TRIGGERED':
        return { variant: 'destructive', label: '🔴 TRIGGERED (Awaiting Ack)' };
      case 'ACKNOWLEDGED':
        return { variant: 'default', label: '🟡 ACKNOWLEDGED' };
      case 'DISPATCHED':
        return { variant: 'default', label: '🔵 RESPONDERS DISPATCHED' };
      case 'ON_SCENE':
        return { variant: 'default', label: '🟢 RESPONDERS ON SCENE' };
      case 'RESOLVED':
        return { variant: 'secondary', label: '✓ RESOLVED' };
      case 'FALSE_ALARM':
        return { variant: 'outline', label: '✕ FALSE ALARM' };
      default:
        return { variant: 'outline', label: status };
    }
  }

  function formatMinutesAgo(timestamp: string): string {
    const mins = Math.floor((Date.now() - new Date(timestamp).getTime()) / 60000);
    if (mins < 1) return 'Just now';
    return `${mins}m ago`;
  }
</script>

<div class="space-y-6">
  <!-- Emergency Control Header Banner -->
  <div class="p-5 bg-gradient-to-r from-red-950 via-slate-900 to-red-900 rounded-2xl border border-red-800 text-white shadow-lg flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div class="space-y-1">
      <div class="flex items-center gap-2">
        <span class="w-3 h-3 rounded-full bg-rose-500 animate-ping"></span>
        <span class="text-xs font-bold uppercase tracking-widest text-rose-400">Live Telemetry Control Hub</span>
      </div>
      <h2 class="text-2xl font-black tracking-tight">SOS & Rapid Emergency Dispatch</h2>
      <p class="text-xs text-slate-300 max-w-xl">
        24/7 real-time geolocation alerts, rapid siren and paramedic dispatch, on-scene coordination, and zero-latency campus surveillance.
      </p>
    </div>

    <div>
      <Button
        onclick={() => (isTriggerModalOpen = true)}
        class="bg-rose-600 hover:bg-rose-700 text-white font-extrabold shadow-md px-6 py-3 rounded-xl border border-rose-400 animate-pulse text-sm"
      >
        🚨 Trigger 1-Tap SOS Alert
      </Button>
    </div>
  </div>

  <!-- Incident Status Counters -->
  <div class="grid grid-cols-1 sm:grid-cols-4 gap-4">
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-bold text-rose-600 uppercase tracking-wider">Active Emergencies</div>
      <div class="text-3xl font-black mt-1 text-foreground">
        {incidents.filter((i) => i.status === 'TRIGGERED' || i.status === 'ACKNOWLEDGED' || i.status === 'DISPATCHED' || i.status === 'ON_SCENE').length}
      </div>
      <div class="text-[11px] text-muted-foreground mt-1">Live responders active</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-bold text-amber-600 uppercase tracking-wider">Awaiting Dispatch</div>
      <div class="text-3xl font-black mt-1 text-foreground">
        {incidents.filter((i) => i.status === 'TRIGGERED' || i.status === 'ACKNOWLEDGED').length}
      </div>
      <div class="text-[11px] text-amber-600 font-medium mt-1">Control room queue</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-bold text-blue-600 uppercase tracking-wider">Responders On Scene</div>
      <div class="text-3xl font-black mt-1 text-foreground">
        {incidents.filter((i) => i.status === 'ON_SCENE').length}
      </div>
      <div class="text-[11px] text-blue-600 font-medium mt-1">First-aid / security active</div>
    </div>
    <div class="p-4 bg-card rounded-xl border shadow-sm">
      <div class="text-xs font-bold text-emerald-600 uppercase tracking-wider">Resolved Today</div>
      <div class="text-3xl font-black mt-1 text-foreground">
        {incidents.filter((i) => i.status === 'RESOLVED').length}
      </div>
      <div class="text-[11px] text-emerald-600 font-medium mt-1">All clear verified</div>
    </div>
  </div>

  <!-- Incident Feed List -->
  <div class="space-y-4">
    <h3 class="text-sm font-bold uppercase tracking-wider text-muted-foreground">
      Live Incident Telemetry Feed ({incidents.length})
    </h3>

    <div class="space-y-3">
      {#each incidents as inc}
        {@const emBadge = getEmergencyTypeBadge(inc.emergencyType)}
        {@const statusBadge = getStatusBadge(inc.status)}
        <div
          role="button"
          tabindex="0"
          onclick={() => onSelectIncident?.(inc)}
          onkeydown={(e) => { if (e.key === 'Enter') onSelectIncident?.(inc); }}
          class="p-4 bg-card hover:bg-muted/40 transition-colors rounded-xl border shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-4 cursor-pointer {inc.status === 'TRIGGERED' ? 'border-rose-500 ring-2 ring-rose-500/20' : ''}"
        >
          <div class="space-y-2 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <span class="font-mono text-xs font-black text-rose-600 bg-rose-50 px-2 py-0.5 rounded border border-rose-200">
                {inc.alertNumber}
              </span>
              <span class="text-xs font-bold px-2.5 py-0.5 rounded-full {emBadge.color}">
                {emBadge.icon}
              </span>
              <Badge variant={statusBadge.variant}>{statusBadge.label}</Badge>
              <span class="text-xs text-muted-foreground">⏱️ Triggered {formatMinutesAgo(inc.triggeredAt)}</span>
            </div>

            <div class="text-sm font-bold text-foreground flex items-center gap-2">
              📍 {inc.locationDescription}
            </div>

            <div class="flex flex-wrap items-center gap-4 text-xs text-muted-foreground">
              <span>👤 {inc.userName}</span>
              <span>📞 {inc.userPhone}</span>
              <span class="font-mono text-[11px] bg-muted px-2 py-0.5 rounded">
                GPS: {inc.latitude.toFixed(4)}° N, {inc.longitude.toFixed(4)}° E
              </span>
            </div>

            {#if inc.resolutionNotes}
              <div class="p-2.5 bg-emerald-50 text-emerald-900 border border-emerald-200 rounded-lg text-xs">
                <span class="font-bold">✓ Resolution Log:</span> {inc.resolutionNotes}
              </div>
            {/if}
          </div>

          <div class="flex flex-col md:items-end justify-center gap-2">
            <Button size="sm" class="text-xs font-bold">
              Open Command Desk →
            </Button>
          </div>
        </div>
      {/each}
    </div>
  </div>

  <!-- Manual / Simulator Trigger Modal -->
  {#if isTriggerModalOpen}
    <div class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-card w-full max-w-lg rounded-2xl border shadow-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b pb-3">
          <div class="flex items-center gap-2">
            <span class="text-rose-600 text-lg">🚨</span>
            <h2 class="text-lg font-bold text-foreground">Trigger Campus Emergency SOS</h2>
          </div>
          <button onclick={() => (isTriggerModalOpen = false)} class="text-muted-foreground hover:text-foreground">✕</button>
        </div>

        <div class="space-y-3">
          <div>
            <label for="em-type-select" class="text-xs font-semibold text-muted-foreground block mb-1">Emergency Category</label>
            <select
              id="em-type-select"
              bind:value={triggerType}
              class="w-full bg-background text-foreground text-sm rounded-lg border px-3 py-2"
            >
              <option value="MEDICAL">Medical Emergency (Paramedic / Ambulance)</option>
              <option value="FIRE">Fire & Smoke Hazard</option>
              <option value="SECURITY_THREAT">Security Threat / Intrusion</option>
              <option value="HARASSMENT_RAGGING">Anti-Ragging / Harassment Alert</option>
              <option value="NATURAL_HAZARD">Natural / Structural Hazard</option>
              <option value="OTHER">Other Urgent Incident</option>
            </select>
          </div>

          <div>
            <label for="location-desc" class="text-xs font-semibold text-muted-foreground block mb-1">Exact Campus Location Description</label>
            <Input id="location-desc" bind:value={triggerLocation} placeholder="e.g. Science Block B - 3rd Floor Lab 302" />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label for="latitude-input" class="text-xs font-semibold text-muted-foreground block mb-1">Latitude</label>
              <Input id="latitude-input" type="number" step="0.0001" bind:value={triggerLat} />
            </div>
            <div>
              <label for="longitude-input" class="text-xs font-semibold text-muted-foreground block mb-1">Longitude</label>
              <Input id="longitude-input" type="number" step="0.0001" bind:value={triggerLng} />
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-3 border-t">
          <Button variant="outline" onclick={() => (isTriggerModalOpen = false)}>Cancel</Button>
          <Button onclick={handleTriggerEmergency} class="bg-rose-600 hover:bg-rose-700 text-white font-bold">
            🚨 Broadcast Emergency SOS
          </Button>
        </div>
      </div>
    </div>
  {/if}
</div>
