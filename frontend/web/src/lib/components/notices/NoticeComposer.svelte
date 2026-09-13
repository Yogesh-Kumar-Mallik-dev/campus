<script lang="ts">
  /**
   * BLOCK_WEB_NOTICES_COMPOSER_001
   * Subsystem: Rank 6 - Notice & Announcement System (notices)
   * Purpose:   Administrative circular composer with audience targeting, priority pinning, and attachment upload.
   */
  import type { NoticeCategory, NoticePriority, TargetAudience, CampusNotice } from '$lib/types/notices';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';

  let {
    onCreated = (_notice: CampusNotice) => {}
  }: {
    onCreated?: (notice: CampusNotice) => void;
  } = $props();

  let title = $state('');
  let content = $state('');
  let category = $state<NoticeCategory>('GENERAL');
  let priority = $state<NoticePriority>('NORMAL');
  let targetAudience = $state<TargetAudience>('ALL');
  let isPinned = $state(false);
  let isSubmitting = $state(false);
  let successMessage = $state<string | null>(null);

  async function handleSubmit(publishNow: boolean) {
    if (!title.trim() || !content.trim()) return;
    isSubmitting = true;
    try {
      const createdNotice: CampusNotice = {
        id: `not_${Date.now()}`,
        tenant_id: 'ten_default',
        author_id: 'usr_admin',
        author_name: 'Registrar Office',
        title: title.trim(),
        slug: title.toLowerCase().replace(/[^a-z0-9]+/g, '-'),
        content: content.trim(),
        category,
        priority,
        status: publishNow ? 'PUBLISHED' : 'DRAFT',
        target_audience: targetAudience,
        is_pinned: isPinned,
        publish_at: new Date().toISOString(),
        view_count: 0,
        is_acknowledged: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      };

      onCreated(createdNotice);
      successMessage = publishNow
        ? 'Notice published immediately to the campus bulletin!'
        : 'Notice saved to drafts queue.';
      title = '';
      content = '';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<Card class="border shadow-sm">
  <CardHeader class="pb-4 border-b bg-muted/20">
    <CardTitle class="text-xl font-bold">Compose Institutional Notice</CardTitle>
    <CardDescription>
      Broadcast announcements to targeted student cohorts, faculty departments, or entire institution.
    </CardDescription>
  </CardHeader>

  <CardContent class="p-6 space-y-4">
    {#if successMessage}
      <div class="p-4 rounded-xl bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 text-emerald-800 dark:text-emerald-200 text-xs font-bold">
        ✓ {successMessage}
      </div>
    {/if}

    <div class="space-y-1.5">
      <Label for="noticeTitle">Notice Title / Headline</Label>
      <Input
        id="noticeTitle"
        bind:value={title}
        placeholder="e.g. Autumn Semester 2026 Academic Calendar Notification"
      />
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="space-y-1.5">
        <Label for="noticeCategory">Category</Label>
        <select
          id="noticeCategory"
          bind:value={category}
          class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm"
        >
          <option value="GENERAL">General Notice</option>
          <option value="ACADEMIC">Academic</option>
          <option value="EXAMINATION">Examination</option>
          <option value="HOSTEL">Hostel & Housing</option>
          <option value="PLACEMENT">Career & Placements</option>
          <option value="EVENTS">Events & Cultural</option>
          <option value="EMERGENCY">Emergency / SOS</option>
          <option value="ADMINISTRATIVE">Administrative</option>
        </select>
      </div>

      <div class="space-y-1.5">
        <Label for="noticePriority">Priority Level</Label>
        <select
          id="noticePriority"
          bind:value={priority}
          class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm"
        >
          <option value="NORMAL">Normal</option>
          <option value="LOW">Low</option>
          <option value="HIGH">High Priority</option>
          <option value="URGENT">Urgent Broadcast</option>
        </select>
      </div>

      <div class="space-y-1.5">
        <Label for="targetAudience">Target Audience</Label>
        <select
          id="targetAudience"
          bind:value={targetAudience}
          class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm"
        >
          <option value="ALL">Entire Campus (All)</option>
          <option value="STUDENTS">Students Only</option>
          <option value="FACULTY">Faculty Only</option>
          <option value="STAFF">Administrative Staff</option>
          <option value="HOSTEL_RESIDENTS">Hostel Boarders</option>
        </select>
      </div>
    </div>

    <div class="space-y-1.5">
      <Label for="noticeContent">Announcement Body</Label>
      <textarea
        id="noticeContent"
        bind:value={content}
        rows="6"
        placeholder="Draft the detailed official notification content here..."
        class="w-full rounded-md border border-input bg-background p-3 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
      ></textarea>
    </div>

    <div class="flex items-center gap-2 pt-2">
      <input
        type="checkbox"
        id="isPinned"
        bind:checked={isPinned}
        class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
      />
      <Label for="isPinned" class="cursor-pointer text-xs font-medium text-foreground">
        Pin this notice to the top of the campus bulletin feed
      </Label>
    </div>
  </CardContent>

  <CardFooter class="flex items-center justify-between p-4 border-t bg-muted/10">
    <Button
      variant="outline"
      disabled={isSubmitting || !title.trim() || !content.trim()}
      onclick={() => handleSubmit(false)}
    >
      Save as Draft
    </Button>
    <Button
      disabled={isSubmitting || !title.trim() || !content.trim()}
      onclick={() => handleSubmit(true)}
    >
      {isSubmitting ? 'Publishing...' : 'Publish Announcement'}
    </Button>
  </CardFooter>
</Card>
