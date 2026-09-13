<script lang="ts">
  /**
   * BLOCK_WEB_MESS_MENU_DESK_001
   * Subsystem: Rank 8 - Mess Management System (mess)
   * Purpose:   Weekly meal menu schedule, nutrition breakdown, and leave-linked rebate applications.
   */
  import type { MessMenuItem, DayOfWeek, MealType, DietCategory, MessRebateApplication } from '$lib/types/mess';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let {
    menuItems = [
      {
        id: 'menu_1',
        tenant_id: 'ten_default',
        mess_hall_id: 'mess_1',
        day_of_week: 'MONDAY',
        meal_type: 'BREAKFAST',
        diet_category: 'VEG',
        title: 'Masala Idli, Medu Vada, Sambar & Coconut Chutney',
        description: 'Served with fresh cut fruits and hot filter coffee / masala tea',
        calories: 450
      },
      {
        id: 'menu_2',
        tenant_id: 'ten_default',
        mess_hall_id: 'mess_1',
        day_of_week: 'MONDAY',
        meal_type: 'LUNCH',
        diet_category: 'VEG',
        title: 'Paneer Butter Masala, Dal Tadka, Jeera Rice & Tawa Roti',
        description: 'Includes Boondi Raita, Roasted Papad, and Gulab Jamun',
        calories: 780
      },
      {
        id: 'menu_3',
        tenant_id: 'ten_default',
        mess_hall_id: 'mess_1',
        day_of_week: 'MONDAY',
        meal_type: 'SNACKS',
        diet_category: 'VEG',
        title: 'Veg Cutlet with Green Mint Chutney & Tea',
        description: 'Fresh evening snacks',
        calories: 320
      },
      {
        id: 'menu_4',
        tenant_id: 'ten_default',
        mess_hall_id: 'mess_1',
        day_of_week: 'MONDAY',
        meal_type: 'DINNER',
        diet_category: 'VEG',
        title: 'Aloo Gobi Matar, Mix Veg Curry, Steamed Basmati & Phulka',
        description: 'Served with hot Tomato Soup and Kheer dessert',
        calories: 710
      },
      {
        id: 'menu_5',
        tenant_id: 'ten_default',
        mess_hall_id: 'mess_1',
        day_of_week: 'TUESDAY',
        meal_type: 'BREAKFAST',
        diet_category: 'VEG',
        title: 'Aloo Paratha with White Butter, Curd & Pickle',
        description: 'Served with hot Chai / Coffee',
        calories: 520
      },
      {
        id: 'menu_6',
        tenant_id: 'ten_default',
        mess_hall_id: 'mess_1',
        day_of_week: 'TUESDAY',
        meal_type: 'LUNCH',
        diet_category: 'VEG',
        title: 'Rajma Masala, Kashmiri Pulao, Kadhi Pakora & Butter Roti',
        description: 'Served with mixed salad and Rasgulla',
        calories: 750
      }
    ] as MessMenuItem[],
    rebates = [
      {
        id: 'reb_1',
        tenant_id: 'ten_default',
        student_id: 'stu_1',
        student_name: 'Rahul Sharma',
        roll_number: '2026CSE001',
        subscription_id: 'sub_1',
        from_date: '2026-09-15T00:00:00Z',
        to_date: '2026-09-20T00:00:00Z',
        total_days: 5,
        daily_rate: 150,
        total_rebate_amount: 750,
        reason: 'Authorized festival break leave (Gate Pass #gp_1)',
        status: 'PENDING',
        created_at: '2026-09-13T09:00:00Z'
      }
    ] as MessRebateApplication[]
  } = $props();

  let selectedDay = $state<DayOfWeek>('MONDAY');
  let selectedDietFilter = $state<string>('ALL');
  let isRebateModalOpen = $state<boolean>(false);

  // Rebate Form
  let rebateStudentName = $state<string>('');
  let rebateRollNumber = $state<string>('');
  let rebateFromDate = $state<string>('');
  let rebateToDate = $state<string>('');
  let rebateReason = $state<string>('');

  let daysCount = $derived(() => {
    if (!rebateFromDate || !rebateToDate) return 0;
    const diff = new Date(rebateToDate).getTime() - new Date(rebateFromDate).getTime();
    return Math.max(0, Math.round(diff / (1000 * 60 * 60 * 24)));
  });

  let estimatedRebateCredit = $derived(daysCount() * 150);

  let filteredMenuItems = $derived(
    menuItems.filter((item) => {
      const matchDay = item.day_of_week === selectedDay;
      const matchDiet = selectedDietFilter === 'ALL' || item.diet_category === selectedDietFilter;
      return matchDay && matchDiet;
    })
  );

  const daysList: DayOfWeek[] = ['MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY'];

  function handleApplyRebate() {
    if (daysCount() < 3) {
      alert('Error: Mess fee rebate requires a minimum of 3 continuous days of leave.');
      return;
    }

    const newRebate: MessRebateApplication = {
      id: `reb_${Date.now()}`,
      tenant_id: 'ten_default',
      student_id: `stu_${Date.now()}`,
      student_name: rebateStudentName || 'Student Applicant',
      roll_number: rebateRollNumber || '2026CSE001',
      subscription_id: 'sub_active',
      from_date: new Date(rebateFromDate).toISOString(),
      to_date: new Date(rebateToDate).toISOString(),
      total_days: daysCount(),
      daily_rate: 150,
      total_rebate_amount: estimatedRebateCredit,
      reason: rebateReason,
      status: 'PENDING',
      created_at: new Date().toISOString()
    };

    rebates.unshift(newRebate);
    isRebateModalOpen = false;
    alert(`✓ Rebate Application Submitted: ₹${estimatedRebateCredit} fee credit requested for ${daysCount()} days.`);
  }

  function handleApproveRebate(r: MessRebateApplication) {
    r.status = 'APPROVED';
    alert(`✓ Rebate Approved: ₹${r.total_rebate_amount} credited to student billing ledger.`);
  }
</script>

<div class="space-y-6">
  <!-- Day Picker & Dietary Category Filters -->
  <Card class="p-4">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <!-- Days of Week Tabs -->
      <div class="flex flex-wrap items-center gap-1 bg-muted p-1 rounded-lg">
        {#each daysList as day}
          <Button
            size="sm"
            variant={selectedDay === day ? 'default' : 'ghost'}
            class="h-8 text-xs font-semibold px-3"
            onclick={() => (selectedDay = day)}
          >
            {day.substring(0, 3)}
          </Button>
        {/each}
      </div>

      <div class="flex items-center gap-3">
        <!-- Diet Filter -->
        <select
          bind:value={selectedDietFilter}
          class="h-9 rounded-md border border-input bg-background px-3 py-1 text-xs shadow-sm focus:outline-none focus:ring-1 focus:ring-ring"
        >
          <option value="ALL">All Dietary Preferences</option>
          <option value="VEG">🌱 Vegetarian</option>
          <option value="NON_VEG">🍗 Non-Vegetarian</option>
          <option value="JAIN">✨ Jain Food</option>
        </select>

        <Button size="sm" variant="outline" class="h-9 text-xs font-semibold" onclick={() => (isRebateModalOpen = true)}>
          + Apply Mess Rebate (≥3 Days)
        </Button>
      </div>
    </div>
  </Card>

  <!-- 4-Meal Slot Menu Grid for Selected Day -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
    {#each ['BREAKFAST', 'LUNCH', 'SNACKS', 'DINNER'] as MealType slot}
      {@const slotItems = filteredMenuItems.filter((i) => i.meal_type === slot)}
      <Card class="shadow-sm border hover:border-primary/40 transition-colors flex flex-col justify-between">
        <CardHeader class="pb-3 border-b bg-muted/10">
          <div class="flex items-center justify-between">
            <CardTitle class="text-base font-bold">{slot}</CardTitle>
            <Badge variant="outline" class="text-[10px] font-semibold">
              {#if slot === 'BREAKFAST'}
                07:30 - 09:30
              {:else if slot === 'LUNCH'}
                12:30 - 14:30
              {:else if slot === 'SNACKS'}
                17:00 - 18:00
              {:else}
                19:30 - 21:30
              {/if}
            </Badge>
          </div>
        </CardHeader>

        <CardContent class="p-4 space-y-3 flex-1">
          {#if slotItems.length === 0}
            <div class="p-6 text-center text-muted-foreground text-xs">
              <p>No dishes scheduled for this slot.</p>
            </div>
          {:else}
            {#each slotItems as item}
              <div class="space-y-1.5">
                <div class="flex items-center gap-2">
                  <Badge variant="secondary" class="bg-emerald-100 text-emerald-800 text-[10px]">
                    {item.diet_category}
                  </Badge>
                  <span class="text-xs font-bold text-foreground leading-tight">{item.title}</span>
                </div>
                {#if item.description}
                  <p class="text-xs text-muted-foreground leading-relaxed pl-1">{item.description}</p>
                {/if}
                <div class="flex items-center justify-between pt-2 text-[11px] text-muted-foreground border-t">
                  <span>Nutritional Energy:</span>
                  <span class="font-bold text-foreground">🔥 {item.calories} kcal</span>
                </div>
              </div>
            {/each}
          {/if}
        </CardContent>

        <CardFooter class="p-3 bg-muted/20 border-t flex items-center justify-between text-[11px] text-muted-foreground">
          <span>Serving Hall: Central North Mess</span>
          <span>Unlimited Buffet</span>
        </CardFooter>
      </Card>
    {/each}
  </div>

  <!-- Rebates & Fee Deductions Queue -->
  <Card class="p-5 border shadow-sm">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h3 class="text-base font-bold">Leave-Linked Mess Fee Rebates</h3>
        <p class="text-xs text-muted-foreground">Automated fee deductions for continuous absence ≥ 3 days</p>
      </div>
      <Badge variant="outline">{rebates.length} Applications</Badge>
    </div>

    <div class="space-y-3">
      {#each rebates as reb (reb.id)}
        <div class="p-4 rounded-xl border bg-muted/10 text-xs flex items-center justify-between">
          <div class="space-y-1">
            <div class="flex items-center gap-2">
              <span class="font-bold text-foreground text-sm">{reb.student_name}</span>
              <Badge variant="outline" class="font-mono">{reb.roll_number}</Badge>
              {#if reb.status === 'PENDING'}
                <Badge class="bg-amber-100 text-amber-800 border-amber-300">PENDING WARDEN APPROVAL</Badge>
              {:else}
                <Badge class="bg-emerald-100 text-emerald-800 border-emerald-300">CREDITED TO LEDGER ✓</Badge>
              {/if}
            </div>
            <p class="text-muted-foreground">{reb.reason}</p>
            <p class="text-muted-foreground">
              Period: {new Date(reb.from_date).toLocaleDateString()} to {new Date(reb.to_date).toLocaleDateString()} ({reb.total_days} Days @ ₹{reb.daily_rate}/day)
            </p>
          </div>

          <div class="text-right space-y-2">
            <div>
              <span class="text-muted-foreground text-[11px]">Rebate Credit:</span>
              <p class="text-base font-extrabold text-emerald-600">₹{reb.total_rebate_amount.toLocaleString('en-IN')}</p>
            </div>
            {#if reb.status === 'PENDING'}
              <Button size="sm" class="h-7 text-xs bg-slate-900 text-white" onclick={() => handleApproveRebate(reb)}>
                Approve & Credit
              </Button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </Card>
</div>

<!-- Apply Rebate Modal -->
{#if isRebateModalOpen}
  <div class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
    <Card class="w-full max-w-lg shadow-2xl bg-card animate-in fade-in zoom-in-95">
      <CardHeader>
        <CardTitle class="text-lg">Apply for Mess Fee Rebate</CardTitle>
        <CardDescription>
          Absence of 3 or more continuous days entitles student to a ₹150/day dining fee refund.
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-3">
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="reb_name">Student Name</label>
            <Input id="reb_name" placeholder="Full name" bind:value={rebateStudentName} />
          </div>
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="reb_roll">Roll Number</label>
            <Input id="reb_roll" placeholder="Roll No" bind:value={rebateRollNumber} />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="reb_from">Absence From Date</label>
            <Input id="reb_from" type="date" bind:value={rebateFromDate} />
          </div>
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="reb_to">Absence To Date</label>
            <Input id="reb_to" type="date" bind:value={rebateToDate} />
          </div>
        </div>

        <div class="space-y-1">
          <label class="text-xs font-semibold" for="reb_reason">Reason / Gate Pass Reference</label>
          <Input id="reb_reason" placeholder="e.g. Approved festival leave" bind:value={rebateReason} />
        </div>

        <!-- Live Rebate Calculation Box -->
        <div class="p-3 rounded-lg bg-emerald-50 border border-emerald-200 text-xs flex justify-between items-center">
          <div>
            <p class="font-bold text-emerald-900">Calculated Absence: {daysCount()} Days</p>
            <p class="text-[11px] text-emerald-700">Daily rate: ₹150 per day</p>
          </div>
          <div class="text-right">
            <p class="text-[11px] text-emerald-700">Estimated Credit:</p>
            <p class="text-base font-extrabold text-emerald-900">₹{estimatedRebateCredit}</p>
          </div>
        </div>
      </CardContent>
      <CardFooter class="flex justify-end gap-2 border-t pt-4">
        <Button variant="ghost" size="sm" onclick={() => (isRebateModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleApplyRebate} disabled={daysCount() < 3}>
          Submit Rebate Application
        </Button>
      </CardFooter>
    </Card>
  </div>
{/if}
