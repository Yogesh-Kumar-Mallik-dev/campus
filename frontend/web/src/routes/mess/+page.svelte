<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_MESS_001
   * Subsystem: Rank 8 - Mess Management System (mess)
   * Purpose:   Campus dining portal: QR tokens, meal punch validator, weekly menu schedule, and leave rebates.
   */
  import DiningTokenScanner from '$lib/components/mess/DiningTokenScanner.svelte';
  import MessMenuDesk from '$lib/components/mess/MessMenuDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Mess & Dining Management
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Digital QR dining punches, weekly nutrition menus, diet preferences, and automated leave fee rebates.
      </p>
    </div>
  </div>

  <Tabs value="scanner" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-sm">
      <TabsTrigger value="scanner">QR Token & Punch Scanner</TabsTrigger>
      <TabsTrigger value="menu">Weekly Menus & Rebates</TabsTrigger>
    </TabsList>

    <TabsContent value="scanner" class="mt-6">
      <DiningTokenScanner />
    </TabsContent>

    <TabsContent value="menu" class="mt-6">
      <MessMenuDesk />
    </TabsContent>
  </Tabs>
</div>
