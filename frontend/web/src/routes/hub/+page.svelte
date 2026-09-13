<!--
  BLOCK_WEB_PAGE_HUB_001
  Subsystem: Rank 17 - The Hub Root Super-App (hub)
  Purpose:   Central Institutional Super-App Cockpit with Multi-Persona Dashboards, Cross-Domain Aggregation, and Widget Customization.
-->
<script lang="ts">
  import PersonaCockpitDesk from '$lib/components/hub/PersonaCockpitDesk.svelte';
  import WidgetConfiguratorDesk from '$lib/components/hub/WidgetConfiguratorDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-7xl mx-auto space-y-6">
  <Tabs value="cockpit" class="w-full">
    <div class="flex items-center justify-between border-b pb-4">
      <div>
        <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
          Campus Super-App Hub
        </h1>
        <p class="text-xs text-muted-foreground mt-0.5">
          Unified persona cockpits, topological telemetry, 1-tap shortcuts, and live cross-subsystem metrics.
        </p>
      </div>

      <TabsList class="grid grid-cols-2 w-72">
        <TabsTrigger value="cockpit">Live Cockpit</TabsTrigger>
        <TabsTrigger value="customize">Customize Widgets</TabsTrigger>
      </TabsList>
    </div>

    <TabsContent value="cockpit" class="mt-6">
      <PersonaCockpitDesk />
    </TabsContent>

    <TabsContent value="customize" class="mt-6">
      <WidgetConfiguratorDesk />
    </TabsContent>
  </Tabs>
</div>
