<script lang="ts">
  /**
   * BLOCK_WEB_HOSTEL_ALLOCATION_GRID_001
   * Subsystem: Rank 7 - Hostel Management System (hostel)
   * Purpose:   Interactive residential block room & bed matrix, real-time occupancy status, and bed assignment desk.
   */
  import type { HostelBlock, HostelRoom, HostelBed, HostelAllocation } from '$lib/types/hostel';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  interface RoomWithBeds extends HostelRoom {
    beds: (HostelBed & { allocation?: HostelAllocation })[];
  }

  let {
    blocks = [
      {
        id: 'blk_1',
        tenant_id: 'ten_default',
        name: 'Aryabhatta Hall of Residence (Block A)',
        code: 'BH-A',
        gender: 'MALE',
        total_floors: 3,
        total_rooms: 60,
        capacity: 120,
        is_active: true,
        created_at: '2026-08-01T00:00:00Z'
      },
      {
        id: 'blk_2',
        tenant_id: 'ten_default',
        name: 'Gargi Hall of Residence (Block B)',
        code: 'GH-B',
        gender: 'FEMALE',
        total_floors: 3,
        total_rooms: 60,
        capacity: 120,
        is_active: true,
        created_at: '2026-08-01T00:00:00Z'
      }
    ] as HostelBlock[],
    rooms = [
      {
        id: 'rm_101',
        tenant_id: 'ten_default',
        block_id: 'blk_1',
        room_number: '101',
        floor_number: 1,
        room_type: 'DOUBLE',
        is_ac: true,
        base_fee_per_semester: 32000,
        status: 'OCCUPIED',
        max_beds: 2,
        beds: [
          {
            id: 'bed_101_1',
            tenant_id: 'ten_default',
            room_id: 'rm_101',
            bed_number: 'A-101-1',
            status: 'ALLOCATED',
            allocation: {
              id: 'alloc_1',
              tenant_id: 'ten_default',
              bed_id: 'bed_101_1',
              student_id: 'stu_1',
              student_name: 'Rahul Sharma',
              roll_number: '2026CSE001',
              academic_year: '2026-2027',
              semester: 1,
              allocated_at: '2026-08-15T10:00:00Z',
              status: 'ALLOCATED'
            }
          },
          {
            id: 'bed_101_2',
            tenant_id: 'ten_default',
            room_id: 'rm_101',
            bed_number: 'A-101-2',
            status: 'ALLOCATED',
            allocation: {
              id: 'alloc_2',
              tenant_id: 'ten_default',
              bed_id: 'bed_101_2',
              student_id: 'stu_2',
              student_name: 'Aditya Verma',
              roll_number: '2026CSE002',
              academic_year: '2026-2027',
              semester: 1,
              allocated_at: '2026-08-15T11:00:00Z',
              status: 'ALLOCATED'
            }
          }
        ]
      },
      {
        id: 'rm_102',
        tenant_id: 'ten_default',
        block_id: 'blk_1',
        room_number: '102',
        floor_number: 1,
        room_type: 'DOUBLE',
        is_ac: false,
        base_fee_per_semester: 25000,
        status: 'AVAILABLE',
        max_beds: 2,
        beds: [
          {
            id: 'bed_102_1',
            tenant_id: 'ten_default',
            room_id: 'rm_102',
            bed_number: 'A-102-1',
            status: 'ALLOCATED',
            allocation: {
              id: 'alloc_3',
              tenant_id: 'ten_default',
              bed_id: 'bed_102_1',
              student_id: 'stu_3',
              student_name: 'Vikram Mehta',
              roll_number: '2026ECE014',
              academic_year: '2026-2027',
              semester: 1,
              allocated_at: '2026-08-16T09:30:00Z',
              status: 'ALLOCATED'
            }
          },
          {
            id: 'bed_102_2',
            tenant_id: 'ten_default',
            room_id: 'rm_102',
            bed_number: 'A-102-2',
            status: 'AVAILABLE'
          }
        ]
      },
      {
        id: 'rm_103',
        tenant_id: 'ten_default',
        block_id: 'blk_1',
        room_number: '103',
        floor_number: 1,
        room_type: 'SINGLE',
        is_ac: true,
        base_fee_per_semester: 45000,
        status: 'AVAILABLE',
        max_beds: 1,
        beds: [
          {
            id: 'bed_103_1',
            tenant_id: 'ten_default',
            room_id: 'rm_103',
            bed_number: 'A-103-1',
            status: 'AVAILABLE'
          }
        ]
      },
      {
        id: 'rm_201',
        tenant_id: 'ten_default',
        block_id: 'blk_1',
        room_number: '201',
        floor_number: 2,
        room_type: 'TRIPLE',
        is_ac: false,
        base_fee_per_semester: 22000,
        status: 'AVAILABLE',
        max_beds: 3,
        beds: [
          {
            id: 'bed_201_1',
            tenant_id: 'ten_default',
            room_id: 'rm_201',
            bed_number: 'A-201-1',
            status: 'AVAILABLE'
          },
          {
            id: 'bed_201_2',
            tenant_id: 'ten_default',
            room_id: 'rm_201',
            bed_number: 'A-201-2',
            status: 'AVAILABLE'
          },
          {
            id: 'bed_201_3',
            tenant_id: 'ten_default',
            room_id: 'rm_201',
            bed_number: 'A-201-3',
            status: 'AVAILABLE'
          }
        ]
      }
    ] as RoomWithBeds[]
  } = $props();

  let selectedBlockId = $state<string>('blk_1');
  let selectedFloor = $state<number | 'ALL'>('ALL');
  let filterAvailability = $state<'ALL' | 'AVAILABLE_ONLY'>('ALL');
  let searchQuery = $state<string>('');

  let allocationModalBed = $state<HostelBed | null>(null);
  let newStudentName = $state<string>('');
  let newStudentRoll = $state<string>('');

  let activeRooms = $derived(
    rooms.filter((rm) => {
      const matchBlock = rm.block_id === selectedBlockId;
      const matchFloor = selectedFloor === 'ALL' || rm.floor_number === selectedFloor;
      const hasAvailableBed = rm.beds.some((b) => b.status === 'AVAILABLE');
      const matchAvailability = filterAvailability === 'ALL' || hasAvailableBed;
      const matchSearch =
        searchQuery.trim() === '' ||
        rm.room_number.includes(searchQuery.trim()) ||
        rm.beds.some((b) => b.allocation?.student_name?.toLowerCase().includes(searchQuery.toLowerCase()) || b.allocation?.roll_number?.toLowerCase().includes(searchQuery.toLowerCase()));

      return matchBlock && matchFloor && matchAvailability && matchSearch;
    })
  );

  let totalBedsCount = $derived(rooms.filter(r => r.block_id === selectedBlockId).reduce((sum, r) => sum + r.beds.length, 0));
  let allocatedBedsCount = $derived(rooms.filter(r => r.block_id === selectedBlockId).reduce((sum, r) => sum + r.beds.filter(b => b.status === 'ALLOCATED').length, 0));
  let vacantBedsCount = $derived(totalBedsCount - allocatedBedsCount);
  let occupancyRate = $derived(totalBedsCount > 0 ? Math.round((allocatedBedsCount / totalBedsCount) * 100) : 0);

  function handleOpenAllocate(bed: HostelBed) {
    allocationModalBed = bed;
    newStudentName = '';
    newStudentRoll = '';
  }

  function handleSaveAllocation() {
    if (!allocationModalBed || !newStudentName || !newStudentRoll) return;

    allocationModalBed.status = 'ALLOCATED';
    (allocationModalBed as any).allocation = {
      id: `alloc_${Date.now()}`,
      tenant_id: 'ten_default',
      bed_id: allocationModalBed.id,
      student_id: `stu_${Date.now()}`,
      student_name: newStudentName,
      roll_number: newStudentRoll,
      academic_year: '2026-2027',
      semester: 1,
      allocated_at: new Date().toISOString(),
      status: 'ALLOCATED'
    };

    allocationModalBed = null;
  }

  function handleVacateBed(bed: HostelBed & { allocation?: HostelAllocation }) {
    if (confirm(`Are you sure you want to vacate bed ${bed.bed_number} for ${bed.allocation?.student_name}?`)) {
      bed.status = 'AVAILABLE';
      bed.allocation = undefined;
    }
  }
</script>

<div class="space-y-6">
  <!-- Stats Summary Bar -->
  <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
    <Card class="bg-card">
      <CardContent class="p-4 flex items-center justify-between">
        <div>
          <p class="text-xs text-muted-foreground font-medium uppercase">Total Capacity</p>
          <p class="text-2xl font-bold tracking-tight">{totalBedsCount} Beds</p>
        </div>
        <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-bold">
          🛏️
        </div>
      </CardContent>
    </Card>

    <Card class="bg-card">
      <CardContent class="p-4 flex items-center justify-between">
        <div>
          <p class="text-xs text-muted-foreground font-medium uppercase">Occupied Beds</p>
          <p class="text-2xl font-bold tracking-tight text-blue-600">{allocatedBedsCount}</p>
        </div>
        <div class="w-10 h-10 rounded-full bg-blue-50 flex items-center justify-center text-blue-600 font-bold">
          👥
        </div>
      </CardContent>
    </Card>

    <Card class="bg-card">
      <CardContent class="p-4 flex items-center justify-between">
        <div>
          <p class="text-xs text-muted-foreground font-medium uppercase">Vacant Beds</p>
          <p class="text-2xl font-bold tracking-tight text-emerald-600">{vacantBedsCount}</p>
        </div>
        <div class="w-10 h-10 rounded-full bg-emerald-50 flex items-center justify-center text-emerald-600 font-bold">
          ✨
        </div>
      </CardContent>
    </Card>

    <Card class="bg-card">
      <CardContent class="p-4 flex items-center justify-between">
        <div>
          <p class="text-xs text-muted-foreground font-medium uppercase">Occupancy Rate</p>
          <p class="text-2xl font-bold tracking-tight">{occupancyRate}%</p>
        </div>
        <div class="w-10 h-10 rounded-full bg-amber-50 flex items-center justify-center text-amber-600 font-bold">
          📊
        </div>
      </CardContent>
    </Card>
  </div>

  <!-- Filter & Action Controls -->
  <Card class="p-4">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <!-- Block & Floor Selector -->
      <div class="flex items-center gap-3">
        <select
          bind:value={selectedBlockId}
          class="h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus:outline-none focus:ring-1 focus:ring-ring"
        >
          {#each blocks as blk}
            <option value={blk.id}>{blk.name} ({blk.code})</option>
          {/each}
        </select>

        <div class="flex items-center gap-1 bg-muted p-1 rounded-md">
          <Button
            size="sm"
            variant={selectedFloor === 'ALL' ? 'default' : 'ghost'}
            class="h-7 px-3 text-xs"
            onclick={() => (selectedFloor = 'ALL')}
          >
            All Floors
          </Button>
          <Button
            size="sm"
            variant={selectedFloor === 1 ? 'default' : 'ghost'}
            class="h-7 px-3 text-xs"
            onclick={() => (selectedFloor = 1)}
          >
            1st Floor
          </Button>
          <Button
            size="sm"
            variant={selectedFloor === 2 ? 'default' : 'ghost'}
            class="h-7 px-3 text-xs"
            onclick={() => (selectedFloor = 2)}
          >
            2nd Floor
          </Button>
          <Button
            size="sm"
            variant={selectedFloor === 3 ? 'default' : 'ghost'}
            class="h-7 px-3 text-xs"
            onclick={() => (selectedFloor = 3)}
          >
            3rd Floor
          </Button>
        </div>
      </div>

      <!-- Search & Availability Toggle -->
      <div class="flex items-center gap-3">
        <Input
          type="search"
          placeholder="Search room, student, roll no..."
          bind:value={searchQuery}
          class="w-64 h-9"
        />

        <Button
          size="sm"
          variant={filterAvailability === 'AVAILABLE_ONLY' ? 'secondary' : 'outline'}
          class="h-9 text-xs"
          onclick={() =>
            (filterAvailability = filterAvailability === 'ALL' ? 'AVAILABLE_ONLY' : 'ALL')}
        >
          {filterAvailability === 'AVAILABLE_ONLY' ? 'Showing Vacant' : 'All Rooms'}
        </Button>
      </div>
    </div>
  </Card>

  <!-- Room Grid Matrix -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
    {#each activeRooms as room (room.id)}
      <Card class="border shadow-sm hover:border-primary/40 transition-colors">
        <CardHeader class="pb-3 border-b bg-muted/20">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-lg font-bold">Room {room.room_number}</span>
              <Badge variant="outline" class="text-[11px] font-mono">Floor {room.floor_number}</Badge>
              {#if room.is_ac}
                <Badge variant="secondary" class="bg-blue-100 text-blue-800 text-[10px]">❄️ AC</Badge>
              {/if}
            </div>
            <span class="text-xs text-muted-foreground font-semibold">
              ₹{room.base_fee_per_semester.toLocaleString('en-IN')}/sem
            </span>
          </div>
          <CardDescription class="text-xs flex items-center justify-between mt-1">
            <span>Type: {room.room_type} ({room.max_beds} Beds)</span>
            <span class="font-medium text-foreground">
              {room.beds.filter((b) => b.status === 'ALLOCATED').length}/{room.max_beds} Occupied
            </span>
          </CardDescription>
        </CardHeader>

        <CardContent class="p-4 space-y-3">
          {#each room.beds as bed (bed.id)}
            <div
              class="p-3 rounded-lg border text-xs flex items-center justify-between {bed.status ===
              'ALLOCATED'
                ? 'bg-blue-50/40 border-blue-200 dark:bg-blue-950/20'
                : 'bg-emerald-50/40 border-emerald-200 dark:bg-emerald-950/20'}"
            >
              <div class="space-y-1">
                <div class="flex items-center gap-2">
                  <span class="font-bold">{bed.bed_number}</span>
                  {#if bed.status === 'ALLOCATED'}
                    <Badge variant="secondary" class="bg-blue-100 text-blue-800 text-[9px] px-1.5 py-0">
                      ALLOCATED
                    </Badge>
                  {:else}
                    <Badge variant="secondary" class="bg-emerald-100 text-emerald-800 text-[9px] px-1.5 py-0">
                      VACANT
                    </Badge>
                  {/if}
                </div>

                {#if bed.allocation}
                  <p class="font-semibold text-foreground">{bed.allocation.student_name}</p>
                  <p class="text-[11px] text-muted-foreground font-mono">Roll: {bed.allocation.roll_number}</p>
                {:else}
                  <p class="text-[11px] text-muted-foreground">Ready for housing allotment</p>
                {/if}
              </div>

              <div>
                {#if bed.status === 'AVAILABLE'}
                  <Button
                    size="sm"
                    class="h-7 text-xs bg-emerald-700 hover:bg-emerald-800 text-white"
                    onclick={() => handleOpenAllocate(bed)}
                  >
                    Allocate
                  </Button>
                {:else}
                  <Button
                    size="sm"
                    variant="outline"
                    class="h-7 text-xs text-destructive hover:bg-destructive/10"
                    onclick={() => handleVacateBed(bed)}
                  >
                    Vacate
                  </Button>
                {/if}
              </div>
            </div>
          {/each}
        </CardContent>
      </Card>
    {/each}
  </div>
</div>

<!-- Bed Allocation Modal -->
{#if allocationModalBed}
  <div class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
    <Card class="w-full max-w-md shadow-2xl bg-card animate-in fade-in zoom-in-95">
      <CardHeader>
        <CardTitle class="text-lg">Allot Bed: {allocationModalBed.bed_number}</CardTitle>
        <CardDescription>Assign residential accommodation to an enrolled student.</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-1.5">
          <label class="text-xs font-semibold" for="student_name_input">Student Full Name</label>
          <Input id="student_name_input" placeholder="e.g. Siddharth Sen" bind:value={newStudentName} />
        </div>
        <div class="space-y-1.5">
          <label class="text-xs font-semibold" for="student_roll_input">Student Roll Number</label>
          <Input id="student_roll_input" placeholder="e.g. 2026MECH042" bind:value={newStudentRoll} />
        </div>
      </CardContent>
      <CardFooter class="flex justify-end gap-2 border-t pt-4">
        <Button variant="ghost" size="sm" onclick={() => (allocationModalBed = null)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleSaveAllocation} disabled={!newStudentName || !newStudentRoll}>
          Confirm Allotment
        </Button>
      </CardFooter>
    </Card>
  </div>
{/if}
