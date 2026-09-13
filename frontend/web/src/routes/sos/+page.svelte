<!--
  BLOCK_WEB_PAGE_SOS_001
  Subsystem: Rank 14 - SOS & Emergency Response (sos)
  Purpose:   Campus SOS & Emergency Operations Portal: Live telemetry radar, rapid response dispatch, and incident history.
-->
<script lang="ts">
  import SOSEmergencyRadar from '$lib/components/sos/SOSEmergencyRadar.svelte';
  import EmergencyResponderDesk from '$lib/components/sos/EmergencyResponderDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
  import type { SOSIncident } from '$lib/types/sos';

  let selectedIncident = $state<SOSIncident | null>(null);
  let activeTab = $state<string>('radar');

  function handleSelectIncident(inc: SOSIncident) {
    selectedIncident = inc;
    activeTab = 'dispatch';
  }

  function handleBackToRadar() {
    selectedIncident = null;
    activeTab = 'radar';
  }
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        SOS & Rapid Emergency Operations
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Central campus safety radar: 1-tap distress alerts, instantaneous GPS telemetry broadcast, paramedic dispatch, and security logs.
      </p>
    </div>
  </div>

  <Tabs bind:value={activeTab} class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-sm">
      <TabsTrigger value="radar">Active Emergency Radar</TabsTrigger>
      <TabsTrigger value="dispatch">Responder Command Desk</TabsTrigger>
    </TabsList>

    <TabsContent value="radar" class="mt-6">
      <SOSEmergencyRadar onSelectIncident={handleSelectIncident} />
    </TabsContent>

    <TabsContent value="dispatch" class="mt-6">
      <EmergencyResponderDesk incident={selectedIncident || undefined} onBack={handleBackToRadar} />
    </TabsContent>
  </Tabs>
</div>
