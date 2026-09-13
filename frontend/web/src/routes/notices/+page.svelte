<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_NOTICES_001
   * Purpose: Institutional Notice Board & Campus Announcement Portal.
   */
  import NoticeFeed from '$lib/components/notices/NoticeFeed.svelte';
  import NoticeComposer from '$lib/components/notices/NoticeComposer.svelte';
  import type { CampusNotice } from '$lib/types/notices';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';

  async function handleAcknowledge(noticeID: string) {
    // In production, sends POST /api/v1/notices/{id}/acknowledge
    console.log('Acknowledged notice:', noticeID);
  }

  function handleNoticeCreated(notice: CampusNotice) {
    console.log('Notice created:', notice);
  }
</script>

<div class="max-w-5xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Campus Bulletins & Notices
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Official circulars, urgent announcements, exam notifications, and campus-wide broadcasts.
      </p>
    </div>
  </div>

  <Tabs value="feed" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-xs">
      <TabsTrigger value="feed">Bulletin Feed</TabsTrigger>
      <TabsTrigger value="compose">Compose Notice</TabsTrigger>
    </TabsList>

    <TabsContent value="feed" class="mt-6">
      <NoticeFeed onAcknowledge={handleAcknowledge} />
    </TabsContent>

    <TabsContent value="compose" class="mt-6">
      <NoticeComposer onCreated={handleNoticeCreated} />
    </TabsContent>
  </Tabs>
</div>
