<script lang="ts">
  /**
   * BLOCK_WEB_NOTICES_FEED_001
   * Subsystem: Rank 6 - Notice & Announcement System (notices)
   * Purpose:   Interactive bulletin feed with pinned alerts, category filters, circular PDF downloads, and read receipts.
   */
  import type { CampusNotice, NoticeCategory } from '$lib/types/notices';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let {
    notices = [
      {
        id: 'not_1',
        tenant_id: 'ten_default',
        author_id: 'usr_dean',
        author_name: 'Dr. Ramesh Rao (Dean Academics)',
        title: 'Midterm Examination Schedule Autumn 2026',
        slug: 'midterm-exam-schedule-autumn-2026',
        content: 'The Midterm examinations for all UG and PG programs will commence from Monday, 19th October 2026. Hall tickets will be issued digitally on 12th October. Students with attendance shortage (<75%) must report to the academic office immediately.',
        category: 'EXAMINATION',
        priority: 'URGENT',
        status: 'PUBLISHED',
        target_audience: 'STUDENTS',
        is_pinned: true,
        publish_at: '2026-09-12T09:00:00Z',
        view_count: 1420,
        is_acknowledged: false,
        attachments: [
          {
            id: 'att_1',
            tenant_id: 'ten_default',
            notice_id: 'not_1',
            file_name: 'Midterm_Exam_Datesheet_2026.pdf',
            file_key: 's3://bulletin/midterm_2026.pdf',
            file_size: 2450000,
            mime_type: 'application/pdf',
            created_at: '2026-09-12T09:00:00Z'
          }
        ],
        created_at: '2026-09-12T09:00:00Z',
        updated_at: '2026-09-12T09:00:00Z'
      },
      {
        id: 'not_2',
        tenant_id: 'ten_default',
        author_id: 'usr_warden',
        author_name: 'Hostel Administration Office',
        title: 'Hostel Curfew & Dining Hall Timings Revision',
        slug: 'hostel-curfew-dining-timings-revision',
        content: 'Effective October 1st, hostel entry gates will close strictly at 22:30. Night mess dining will remain open until 23:15. Gate passes must be approved by wardens 24 hours in advance.',
        category: 'HOSTEL',
        priority: 'HIGH',
        status: 'PUBLISHED',
        target_audience: 'HOSTEL_RESIDENTS',
        is_pinned: false,
        publish_at: '2026-09-11T14:30:00Z',
        view_count: 850,
        is_acknowledged: true,
        created_at: '2026-09-11T14:30:00Z',
        updated_at: '2026-09-11T14:30:00Z'
      },
      {
        id: 'not_3',
        tenant_id: 'ten_default',
        author_id: 'usr_cdc',
        author_name: 'Career & Placement Cell',
        title: 'Google & Microsoft Campus Placement Drive Registration',
        slug: 'google-microsoft-campus-drive-2026',
        content: 'Registrations are open for the 2026 Campus Placement recruitment season. Eligible B.Tech CSE/ECE/IT students with CGPA >= 7.5 must submit verified resumes before September 25th.',
        category: 'PLACEMENT',
        priority: 'NORMAL',
        status: 'PUBLISHED',
        target_audience: 'STUDENTS',
        is_pinned: false,
        publish_at: '2026-09-10T11:00:00Z',
        view_count: 2100,
        is_acknowledged: false,
        created_at: '2026-09-10T11:00:00Z',
        updated_at: '2026-09-10T11:00:00Z'
      }
    ],
    onAcknowledge = async (_noticeID: string) => {}
  }: {
    notices?: CampusNotice[];
    onAcknowledge?: (noticeID: string) => Promise<void>;
  } = $props();

  let searchQuery = $state('');
  let selectedCategory = $state<NoticeCategory | 'ALL'>('ALL');

  let filteredNotices = $derived(
    notices.filter((n) => {
      if (selectedCategory !== 'ALL' && n.category !== selectedCategory) return false;
      if (searchQuery.trim() !== '') {
        const q = searchQuery.toLowerCase();
        return n.title.toLowerCase().includes(q) || n.content.toLowerCase().includes(q);
      }
      return true;
    }).sort((a, b) => {
      // Pinned notices first, then newest
      if (a.is_pinned !== b.is_pinned) return a.is_pinned ? -1 : 1;
      return new Date(b.publish_at).getTime() - new Date(a.publish_at).getTime();
    })
  );

  const categories: (NoticeCategory | 'ALL')[] = [
    'ALL',
    'EXAMINATION',
    'ACADEMIC',
    'HOSTEL',
    'PLACEMENT',
    'EVENTS',
    'EMERGENCY',
    'ADMINISTRATIVE'
  ];

  async function handleAcknowledgeClick(notice: CampusNotice) {
    if (notice.is_acknowledged) return;
    await onAcknowledge(notice.id);
    notice.is_acknowledged = true;
  }
</script>

<div class="space-y-6">
  <!-- Search & Category Filters -->
  <div class="space-y-3">
    <div class="max-w-md">
      <Input
        placeholder="Search campus bulletins and circulars..."
        bind:value={searchQuery}
        class="bg-background"
      />
    </div>

    <!-- Category Pills -->
    <div class="flex flex-wrap gap-1.5">
      {#each categories as cat}
        <button
          type="button"
          class="px-3 py-1 rounded-full text-xs font-semibold transition-all {selectedCategory === cat ? 'bg-primary text-primary-foreground shadow-sm' : 'bg-muted hover:bg-muted/80 text-muted-foreground'}"
          onclick={() => (selectedCategory = cat)}
        >
          {cat}
        </button>
      {/each}
    </div>
  </div>

  <!-- Notices Cards -->
  <div class="space-y-4">
    {#each filteredNotices as notice (notice.id)}
      <Card class="border shadow-sm transition-all {notice.priority === 'URGENT' ? 'border-rose-300 dark:border-rose-900 bg-rose-50/20 dark:bg-rose-950/10' : notice.is_pinned ? 'border-amber-300 dark:border-amber-900 bg-amber-50/20 dark:bg-amber-950/10' : ''}">
        <CardHeader class="pb-3">
          <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-2">
            <div class="space-y-1">
              <div class="flex flex-wrap items-center gap-2">
                {#if notice.is_pinned}
                  <Badge variant="outline" class="bg-amber-100 text-amber-800 dark:bg-amber-950 dark:text-amber-300 border-amber-300 text-[10px] font-bold">
                    📌 PINNED
                  </Badge>
                {/if}
                <Badge variant={notice.priority === 'URGENT' ? 'destructive' : notice.priority === 'HIGH' ? 'default' : 'secondary'}>
                  {notice.priority}
                </Badge>
                <Badge variant="outline" class="text-xs">
                  {notice.category}
                </Badge>
                <span class="text-xs text-muted-foreground">· Target: {notice.target_audience}</span>
              </div>

              <CardTitle class="text-lg font-bold text-foreground mt-1">
                {notice.title}
              </CardTitle>
            </div>

            <span class="text-xs text-muted-foreground whitespace-nowrap">
              {new Date(notice.publish_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })}
            </span>
          </div>

          <CardDescription class="text-xs text-muted-foreground mt-1">
            Issued by: <span class="font-medium text-foreground">{notice.author_name || notice.author_id}</span>
          </CardDescription>
        </CardHeader>

        <CardContent class="text-sm text-foreground leading-relaxed">
          <p class="whitespace-pre-line">{notice.content}</p>

          {#if notice.attachments && notice.attachments.length > 0}
            <div class="mt-4 pt-3 border-t flex flex-wrap gap-2">
              {#each notice.attachments as att}
                <div class="flex items-center gap-2 px-3 py-1.5 rounded-lg border bg-background text-xs">
                  <span>📄</span>
                  <span class="font-medium text-foreground">{att.file_name}</span>
                  <span class="text-muted-foreground">({(att.file_size / (1024 * 1024)).toFixed(1)} MB)</span>
                  <Button variant="ghost" size="sm" class="h-6 px-2 text-xs font-semibold text-primary">
                    Download
                  </Button>
                </div>
              {/each}
            </div>
          {/if}
        </CardContent>

        <CardFooter class="flex items-center justify-between pt-3 pb-3 border-t bg-muted/20 text-xs">
          <span class="text-muted-foreground">
            👁 {notice.view_count.toLocaleString()} views
          </span>

          <div>
            {#if notice.is_acknowledged}
              <span class="inline-flex items-center gap-1 text-emerald-600 dark:text-emerald-400 font-medium">
                ✓ Read & Acknowledged
              </span>
            {:else}
              <Button
                variant="outline"
                size="sm"
                class="h-7 text-xs"
                onclick={() => handleAcknowledgeClick(notice)}
              >
                Acknowledge Notice
              </Button>
            {/if}
          </div>
        </CardFooter>
      </Card>
    {/each}
  </div>
</div>
