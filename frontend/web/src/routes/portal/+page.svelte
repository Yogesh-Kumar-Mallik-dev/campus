<!--
  BLOCK_WEB_PAGE_PORTAL_001
  Subsystem: Rank 16 - Public Web Portal (portal)
  Purpose:   Institutional Public Landing, Academic Program Catalog, and Admissions Prospect Lead Triage Portal.
-->
<script lang="ts">
  import PublicLandingShowcase from '$lib/components/portal/PublicLandingShowcase.svelte';
  import ProspectInquiryDesk from '$lib/components/portal/ProspectInquiryDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Public Institutional Portal
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Public campus showcase, accredited degree catalog, admissions cycle announcements, and prospect triage console.
      </p>
    </div>
  </div>

  <Tabs value="showcase" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-md">
      <TabsTrigger value="showcase">Public Landing & Catalog</TabsTrigger>
      <TabsTrigger value="inquiries">Admissions Lead Desk</TabsTrigger>
    </TabsList>

    <TabsContent value="showcase" class="mt-6">
      <PublicLandingShowcase />
    </TabsContent>

    <TabsContent value="inquiries" class="mt-6">
      <ProspectInquiryDesk />
    </TabsContent>
  </Tabs>
</div>
