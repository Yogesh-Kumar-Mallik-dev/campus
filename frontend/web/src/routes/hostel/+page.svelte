<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_HOSTEL_001
   * Subsystem: Rank 7 - Hostel Management System (hostel)
   * Purpose:   Hostel residential management portal: room & bed allocations, warden gate passes, and curfew tracking.
   */
  import BedAllocationGrid from '$lib/components/hostel/BedAllocationGrid.svelte';
  import GatePassDesk from '$lib/components/hostel/GatePassDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Hostel & Residence Management
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Residential block allocations, room inventory, warden gate pass approvals, and curfew breach logging.
      </p>
    </div>
  </div>

  <Tabs value="allocations" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-sm">
      <TabsTrigger value="allocations">Bed Allocations & Matrix</TabsTrigger>
      <TabsTrigger value="gatepasses">Gate Passes & Curfew Desk</TabsTrigger>
    </TabsList>

    <TabsContent value="allocations" class="mt-6">
      <BedAllocationGrid />
    </TabsContent>

    <TabsContent value="gatepasses" class="mt-6">
      <GatePassDesk />
    </TabsContent>
  </Tabs>
</div>
