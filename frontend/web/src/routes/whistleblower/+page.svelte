<!--
  BLOCK_WEB_PAGE_WHISTLEBLOWER_001
  Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
  Purpose:   Institutional Whistleblower Portal: Zero-knowledge reporting, secret token tracking, and vigilance committee desk.
-->
<script lang="ts">
  import AnonymousReportForm from '$lib/components/whistleblower/AnonymousReportForm.svelte';
  import TrackingPortalDesk from '$lib/components/whistleblower/TrackingPortalDesk.svelte';
  import CommitteeInvestigationDesk from '$lib/components/whistleblower/CommitteeInvestigationDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Anonymity & Whistleblower Portal
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Zero-knowledge confidential reporting, encrypted tracking keys, UGC anti-ragging protection, and vigilance investigation desk.
      </p>
    </div>
  </div>

  <Tabs value="report" class="w-full">
    <TabsList class="grid w-full grid-cols-3 max-w-lg">
      <TabsTrigger value="report">Submit Secret Report</TabsTrigger>
      <TabsTrigger value="track">Track by Secret Key</TabsTrigger>
      <TabsTrigger value="committee">Committee Console</TabsTrigger>
    </TabsList>

    <TabsContent value="report" class="mt-6">
      <AnonymousReportForm />
    </TabsContent>

    <TabsContent value="track" class="mt-6">
      <TrackingPortalDesk />
    </TabsContent>

    <TabsContent value="committee" class="mt-6">
      <CommitteeInvestigationDesk />
    </TabsContent>
  </Tabs>
</div>
