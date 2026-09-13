<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_EVENTS_001
   * Subsystem: Rank 12 - Event Organisation System (events)
   * Purpose:   Campus events hub: discovery calendar, venue conflict resolution, digital ticketing, and QR turnstile check-in.
   */
  import EventCalendarDesk from '$lib/components/events/EventCalendarDesk.svelte';
  import TicketCheckInDesk from '$lib/components/events/TicketCheckInDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Events & Venue Booking
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Campus event calendar, conflict-free venue reservations, student ticketing quotas, and gate admissions scanner.
      </p>
    </div>
  </div>

  <Tabs value="calendar" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-sm">
      <TabsTrigger value="calendar">Events Calendar & Tickets</TabsTrigger>
      <TabsTrigger value="checkin">Gate Scanner & Admissions</TabsTrigger>
    </TabsList>

    <TabsContent value="calendar" class="mt-6">
      <EventCalendarDesk />
    </TabsContent>

    <TabsContent value="checkin" class="mt-6">
      <TicketCheckInDesk />
    </TabsContent>
  </Tabs>
</div>
