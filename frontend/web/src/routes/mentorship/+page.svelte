<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_MENTORSHIP_001
   * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
   * Purpose:   Faculty mentor allocations, 1-on-1 counseling log matrix, and automated at-risk GPA/attendance early warning radar.
   */
  import MentorAllocationDesk from '$lib/components/mentorship/MentorAllocationDesk.svelte';
  import MentorshipCounselingDesk from '$lib/components/mentorship/MentorshipCounselingDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Progress Tracker & Mentorship
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Faculty mentor allocations, 1-on-1 counseling logs, automated at-risk surveillance, and academic remediation workflows.
      </p>
    </div>
  </div>

  <Tabs value="counseling" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-sm">
      <TabsTrigger value="counseling">Counseling & Academic Radar</TabsTrigger>
      <TabsTrigger value="allocations">Mentor Allocations</TabsTrigger>
    </TabsList>

    <TabsContent value="counseling" class="mt-6">
      <MentorshipCounselingDesk />
    </TabsContent>

    <TabsContent value="allocations" class="mt-6">
      <MentorAllocationDesk />
    </TabsContent>
  </Tabs>
</div>
