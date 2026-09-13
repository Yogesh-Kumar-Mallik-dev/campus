<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_STUDYHUB_001
   * Subsystem: Rank 10 - Study Hub System (studyhub)
   * Purpose:   Integrated academic syllabus, lecture notes & lab manuals, assignment submission desk, and peer evaluation.
   */
  import CourseMaterialDesk from '$lib/components/studyhub/CourseMaterialDesk.svelte';
  import AssignmentSubmissionDesk from '$lib/components/studyhub/AssignmentSubmissionDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Study Hub & Assignments
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Course syllabi, unit-wise notes & lab manuals, assignment submissions, late penalty grading, and peer reviews.
      </p>
    </div>
  </div>

  <Tabs value="materials" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-sm">
      <TabsTrigger value="materials">Syllabus & Materials</TabsTrigger>
      <TabsTrigger value="assignments">Assignments & Grading</TabsTrigger>
    </TabsList>

    <TabsContent value="materials" class="mt-6">
      <CourseMaterialDesk />
    </TabsContent>

    <TabsContent value="assignments" class="mt-6">
      <AssignmentSubmissionDesk />
    </TabsContent>
  </Tabs>
</div>
