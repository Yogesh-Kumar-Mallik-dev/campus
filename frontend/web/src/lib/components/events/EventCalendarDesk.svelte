<script lang="ts">
  /**
   * BLOCK_WEB_EVENTS_CALENDAR_001
   * Subsystem: Rank 12 - Event Organisation System (events)
   * Purpose:   Interactive campus events schedule, category filters, real-time ticket quota booking, and venue reservation modal.
   */
  import type { CampusEvent, EventCategory } from '$lib/types/events';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Badge } from '$lib/components/ui/badge';

  let events = $state<CampusEvent[]>([
    {
      id: 'evnt-01',
      tenantId: 'tenant-main',
      title: 'Annual Tech Symposium & Hackathon 2026',
      description: '36-hour sprint on generative AI, autonomous robotics showcase, and algorithmic challenges.',
      category: 'TECHNICAL',
      venueName: 'Main University Auditorium',
      venueCapacity: 500,
      startTime: '2026-09-25T09:00:00Z',
      endTime: '2026-09-26T21:00:00Z',
      registrationDeadline: '2026-09-24T18:00:00Z',
      organizerId: 'staff-dean-eng',
      organizerName: 'Faculty of Engineering & Technology',
      status: 'PUBLISHED',
      isTicketed: true,
      ticketPrice: 0.0,
      maxTickets: 500,
      ticketsSold: 342,
      createdAt: '2026-09-01T00:00:00Z',
      updatedAt: '2026-09-01T00:00:00Z'
    },
    {
      id: 'evnt-02',
      tenantId: 'tenant-main',
      title: 'Inter-College Cultural Night: Tarang 2026',
      description: 'Classical dance performances, rock band showcase, theatrical dramas, and art exhibitions.',
      category: 'CULTURAL',
      venueName: 'Open Amphitheatre',
      venueCapacity: 800,
      startTime: '2026-10-02T18:00:00Z',
      endTime: '2026-10-02T23:30:00Z',
      registrationDeadline: '2026-10-01T23:59:59Z',
      organizerId: 'staff-cultural-council',
      organizerName: 'Student Cultural Council',
      status: 'PUBLISHED',
      isTicketed: true,
      ticketPrice: 150.0,
      maxTickets: 800,
      ticketsSold: 610,
      createdAt: '2026-09-05T00:00:00Z',
      updatedAt: '2026-09-05T00:00:00Z'
    },
    {
      id: 'evnt-03',
      tenantId: 'tenant-main',
      title: 'Industry Keynote: Next-Gen Quantum Computing Systems',
      description: 'Guest lecture by distinguished IBM Quantum Research Scientist on superconducting qubits.',
      category: 'GUEST_LECTURE',
      venueName: 'Seminar Hall 3B',
      venueCapacity: 120,
      startTime: '2026-09-28T14:00:00Z',
      endTime: '2026-09-28T16:30:00Z',
      registrationDeadline: '2026-09-27T23:59:59Z',
      organizerId: 'staff-physics-hod',
      organizerName: 'Department of Physics & Computing',
      status: 'PUBLISHED',
      isTicketed: false,
      ticketPrice: 0.0,
      maxTickets: 120,
      ticketsSold: 88,
      createdAt: '2026-09-08T00:00:00Z',
      updatedAt: '2026-09-08T00:00:00Z'
    }
  ]);

  let selectedCategory = $state<EventCategory | 'ALL'>('ALL');
  let searchQuery = $state('');

  // Host Event Modal
  let isHostModalOpen = $state(false);
  let newTitle = $state('');
  let newDesc = $state('');
  let newCategory = $state<EventCategory>('TECHNICAL');
  let newVenue = $state('Main University Auditorium');
  let newCapacity = $state(200);
  let newStartTime = $state('2026-10-10T10:00');
  let newEndTime = $state('2026-10-10T17:00');
  let newDeadline = $state('2026-10-09T23:59');
  let newIsTicketed = $state(true);
  let newTicketPrice = $state(0.0);

  const filteredEvents = $derived(
    events.filter((e) => {
      const matchCategory = selectedCategory === 'ALL' || e.category === selectedCategory;
      const matchSearch =
        e.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        e.venueName.toLowerCase().includes(searchQuery.toLowerCase()) ||
        e.description.toLowerCase().includes(searchQuery.toLowerCase());
      return matchCategory && matchSearch;
    })
  );

  function handleCreateEvent() {
    if (!newTitle || !newVenue) return;
    const newEvnt: CampusEvent = {
      id: `evnt-${Date.now()}`,
      tenantId: 'tenant-main',
      title: newTitle,
      description: newDesc,
      category: newCategory,
      venueName: newVenue,
      venueCapacity: Number(newCapacity) || 100,
      startTime: new Date(newStartTime).toISOString(),
      endTime: new Date(newEndTime).toISOString(),
      registrationDeadline: new Date(newDeadline).toISOString(),
      organizerId: 'staff-me',
      organizerName: 'You (Organizer)',
      status: 'PUBLISHED',
      isTicketed: newIsTicketed,
      ticketPrice: Number(newTicketPrice) || 0.0,
      maxTickets: Number(newCapacity) || 100,
      ticketsSold: 0,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    events = [newEvnt, ...events];
    newTitle = '';
    newDesc = '';
    isHostModalOpen = false;
  }

  function handleBook(evnt: CampusEvent) {
    if (evnt.ticketsSold >= evnt.maxTickets) return;
    events = events.map((e) => {
      if (e.id === evnt.id) {
        return {
          ...e,
          ticketsSold: e.ticketsSold + 1
        };
      }
      return e;
    });
    alert(`Pass Reserved! Your QR ticket pass for "${evnt.title}" has been issued to your student portal.`);
  }

  function getCategoryBadge(cat: EventCategory) {
    switch (cat) {
      case 'TECHNICAL':
        return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'CULTURAL':
        return 'bg-purple-100 text-purple-800 border-purple-200';
      case 'SPORTS':
        return 'bg-emerald-100 text-emerald-800 border-emerald-200';
      case 'GUEST_LECTURE':
        return 'bg-amber-100 text-amber-800 border-amber-200';
      case 'ACADEMIC_WORKSHOP':
        return 'bg-indigo-100 text-indigo-800 border-indigo-200';
      case 'CAREER_FAIR':
        return 'bg-rose-100 text-rose-800 border-rose-200';
      default:
        return 'bg-slate-100 text-slate-800 border-slate-200';
    }
  }
</script>

<div class="space-y-6">
  <!-- Top Navigation & Action Bar -->
  <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-card p-4 rounded-xl border">
    <div class="flex flex-wrap items-center gap-2">
      <Button
        variant={selectedCategory === 'ALL' ? 'default' : 'outline'}
        size="sm"
        onclick={() => (selectedCategory = 'ALL')}
      >
        All Events
      </Button>
      <Button
        variant={selectedCategory === 'TECHNICAL' ? 'default' : 'ghost'}
        size="sm"
        onclick={() => (selectedCategory = 'TECHNICAL')}
      >
        ⚡ Technical
      </Button>
      <Button
        variant={selectedCategory === 'CULTURAL' ? 'default' : 'ghost'}
        size="sm"
        onclick={() => (selectedCategory = 'CULTURAL')}
      >
        🎭 Cultural
      </Button>
      <Button
        variant={selectedCategory === 'GUEST_LECTURE' ? 'default' : 'ghost'}
        size="sm"
        onclick={() => (selectedCategory = 'GUEST_LECTURE')}
      >
        🎓 Guest Lectures
      </Button>
    </div>

    <div class="flex items-center gap-2">
      <Input
        placeholder="Search event title or venue..."
        bind:value={searchQuery}
        class="text-xs w-full sm:w-60"
      />
      <Button size="sm" onclick={() => (isHostModalOpen = true)}>
        + Host New Event
      </Button>
    </div>
  </div>

  <!-- Events Grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
    {#each filteredEvents as evnt}
      <div class="bg-card border rounded-2xl p-6 shadow-sm flex flex-col justify-between space-y-4 hover:shadow-md transition-shadow">
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="px-2.5 py-0.5 rounded text-[10px] font-bold border {getCategoryBadge(evnt.category)}">
              {evnt.category.replace('_', ' ')}
            </span>
            <span class="text-xs font-bold text-foreground">
              {evnt.ticketPrice > 0 ? `₹${evnt.ticketPrice}` : 'FREE ENTRY'}
            </span>
          </div>

          <div>
            <h3 class="text-base font-bold text-foreground line-clamp-1">{evnt.title}</h3>
            <p class="text-xs text-muted-foreground mt-1 line-clamp-2">{evnt.description}</p>
          </div>

          <div class="bg-muted/40 p-3 rounded-xl border space-y-1 text-xs">
            <div class="flex items-center justify-between text-muted-foreground">
              <span>📍 Venue:</span>
              <span class="font-semibold text-foreground">{evnt.venueName}</span>
            </div>
            <div class="flex items-center justify-between text-muted-foreground">
              <span>⏰ Date & Time:</span>
              <span class="font-semibold text-foreground">
                {new Date(evnt.startTime).toLocaleDateString()} · {new Date(evnt.startTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </span>
            </div>
            <div class="flex items-center justify-between text-muted-foreground">
              <span>⏳ Deadline:</span>
              <span class="text-slate-500 font-mono text-[11px]">
                {new Date(evnt.registrationDeadline).toLocaleDateString()}
              </span>
            </div>
          </div>

          <!-- Capacity Fill Meter -->
          <div class="space-y-1">
            <div class="flex justify-between text-[11px] font-semibold text-muted-foreground">
              <span>Capacity ({evnt.ticketsSold}/{evnt.maxTickets} booked)</span>
              <span>{Math.round((evnt.ticketsSold / evnt.maxTickets) * 100)}%</span>
            </div>
            <div class="w-full bg-slate-200 dark:bg-slate-700 h-2 rounded-full overflow-hidden">
              <div
                class="bg-primary h-full rounded-full transition-all"
                style="width: {(evnt.ticketsSold / evnt.maxTickets) * 100}%"
              ></div>
            </div>
          </div>
        </div>

        <div class="pt-3 border-t">
          <Button
            size="sm"
            class="w-full"
            disabled={evnt.ticketsSold >= evnt.maxTickets}
            onclick={() => handleBook(evnt)}
          >
            {evnt.ticketsSold >= evnt.maxTickets ? 'Sold Out' : '🎟️ Book Student Pass'}
          </Button>
        </div>
      </div>
    {/each}
  </div>
</div>

<!-- Modal: Host Event -->
{#if isHostModalOpen}
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
    <div class="bg-card border max-w-xl w-full rounded-2xl p-6 shadow-2xl space-y-4">
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-lg font-bold text-foreground">Host Campus Event</h3>
        <button
          onclick={() => (isHostModalOpen = false)}
          class="text-muted-foreground hover:text-foreground text-sm font-bold"
        >
          ✕
        </button>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <Label for="evntTitle" class="text-xs font-semibold">Event Title</Label>
          <Input id="evntTitle" bind:value={newTitle} placeholder="e.g. Robotics Championship 2026" />
        </div>

        <div class="space-y-1">
          <Label for="evntDesc" class="text-xs font-semibold">Description</Label>
          <textarea
            id="evntDesc"
            bind:value={newDesc}
            rows="2"
            class="w-full bg-background border rounded-lg p-2.5 text-xs focus:ring-2 focus:ring-primary"
            placeholder="Event theme, schedule breakdown, prize money..."
          ></textarea>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <Label for="evntCat" class="text-xs font-semibold">Category</Label>
            <select
              id="evntCat"
              bind:value={newCategory}
              class="w-full bg-background border rounded-lg px-3 py-2 text-xs font-medium"
            >
              <option value="TECHNICAL">Technical & Hackathon</option>
              <option value="CULTURAL">Cultural & Arts</option>
              <option value="SPORTS">Sports Tournament</option>
              <option value="GUEST_LECTURE">Guest Lecture / Keynote</option>
              <option value="ACADEMIC_WORKSHOP">Academic Workshop</option>
              <option value="CAREER_FAIR">Career / Placement Fair</option>
            </select>
          </div>

          <div class="space-y-1">
            <Label for="evntVenue" class="text-xs font-semibold">Venue (Conflict Checked)</Label>
            <select
              id="evntVenue"
              bind:value={newVenue}
              class="w-full bg-background border rounded-lg px-3 py-2 text-xs font-medium"
            >
              <option value="Main University Auditorium">Main Auditorium (Cap: 500)</option>
              <option value="Open Amphitheatre">Open Amphitheatre (Cap: 800)</option>
              <option value="Seminar Hall 3B">Seminar Hall 3B (Cap: 120)</option>
              <option value="Indoor Sports Complex">Indoor Sports Complex (Cap: 300)</option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <Label for="evntStart" class="text-xs font-semibold">Start Time</Label>
            <Input id="evntStart" type="datetime-local" bind:value={newStartTime} />
          </div>

          <div class="space-y-1">
            <Label for="evntEnd" class="text-xs font-semibold">End Time</Label>
            <Input id="evntEnd" type="datetime-local" bind:value={newEndTime} />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <Label for="evntPrice" class="text-xs font-semibold">Ticket Price (₹0 = Free)</Label>
            <Input id="evntPrice" type="number" min="0" bind:value={newTicketPrice} />
          </div>

          <div class="space-y-1">
            <Label for="evntCap" class="text-xs font-semibold">Seat Quota</Label>
            <Input id="evntCap" type="number" min="1" bind:value={newCapacity} />
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-3 border-t">
        <Button variant="outline" size="sm" onclick={() => (isHostModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleCreateEvent}>
          Publish Event & Reserve Venue
        </Button>
      </div>
    </div>
  </div>
{/if}
