<script lang="ts">
 import PublicShell from "$public/components/PublicShell.svelte"; import PageHeader from "$public/components/PageHeader.svelte"; import { Button } from "$lib/components/ui/button"; import * as Empty from "$lib/components/ui/empty"; import IconNews from "@tabler/icons-svelte/icons/news";
</script>
<svelte:head><title>News & Announcements · BBDIT</title></svelte:head>
<PublicShell><main><PageHeader eyebrow="Newsroom" title="News and announcements from BBDIT" description="Official academic notices, placement updates, campus events and institutional announcements published during the last 30 days." />
  <section class="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8"><Empty.Root class="rounded-xl border border-dashed"><Empty.Header><Empty.Media variant="icon"><IconNews /></Empty.Media><Empty.Title>No announcements in the last 30 days</Empty.Title><Empty.Description>The official BBDIT announcements archive has no posts dated between 22 July and 21 August 2026.</Empty.Description></Empty.Header><Empty.Content><Button href="https://bbdit.edu.in/posts/announcements/" target="_blank" rel="noreferrer" variant="outline">Open the official archive</Button></Empty.Content></Empty.Root></section>
</main></PublicShell>
