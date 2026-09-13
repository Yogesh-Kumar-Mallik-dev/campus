<!--
  BLOCK_WEB_PAGE_HELPDESK_001
  Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
  Purpose:   Central Helpdesk portal: SLA radar surveillance, live resolution chat desk, and category SLA configuration.
-->
<script lang="ts">
  import TicketKanbanDesk from '$lib/components/helpdesk/TicketKanbanDesk.svelte';
  import TicketDetailChat from '$lib/components/helpdesk/TicketDetailChat.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
  import type { HelpdeskTicket } from '$lib/types/helpdesk';

  let selectedTicket = $state<HelpdeskTicket | null>(null);
  let activeTab = $state<string>('tickets');

  function handleSelectTicket(t: HelpdeskTicket) {
    selectedTicket = t;
    activeTab = 'chat';
  }

  function handleBackToTickets() {
    selectedTicket = null;
    activeTab = 'tickets';
  }
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Application & Query Helpdesk
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Central institutional ticketing hub, SLA breach radar surveillance, staff internal notes, and applicant resolution desk.
      </p>
    </div>
  </div>

  <Tabs bind:value={activeTab} class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-sm">
      <TabsTrigger value="tickets">All Tickets & Radar</TabsTrigger>
      <TabsTrigger value="chat">Active Resolution Desk</TabsTrigger>
    </TabsList>

    <TabsContent value="tickets" class="mt-6">
      <TicketKanbanDesk onSelectTicket={handleSelectTicket} />
    </TabsContent>

    <TabsContent value="chat" class="mt-6">
      <TicketDetailChat ticket={selectedTicket || undefined} onBack={handleBackToTickets} />
    </TabsContent>
  </Tabs>
</div>
