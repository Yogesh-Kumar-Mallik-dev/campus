<script lang="ts">
  /**
   * BLOCK_WEB_MESS_TOKEN_SCANNER_001
   * Subsystem: Rank 8 - Mess Management System (mess)
   * Purpose:   Interactive QR Dining Token generator, live countdown validity, and entrance punch scanner.
   */
  import type { MessDiningToken, MealType } from '$lib/types/mess';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let {
    activeTokens = [
      {
        id: 'tok_001',
        tenant_id: 'ten_default',
        student_id: 'stu_1',
        mess_hall_id: 'mess_1',
        meal_type: 'LUNCH',
        token_code: 'TOK-9d7a2f10b4',
        meal_date: '2026-09-13',
        status: 'GENERATED',
        created_at: new Date().toISOString()
      },
      {
        id: 'tok_002',
        tenant_id: 'ten_default',
        student_id: 'stu_2',
        mess_hall_id: 'mess_1',
        meal_type: 'BREAKFAST',
        token_code: 'TOK-4c8e11a2f9',
        meal_date: '2026-09-13',
        status: 'REDEEMED',
        redeemed_at: '2026-09-13T08:15:00Z',
        device_reader_id: 'READER_GATE_01',
        created_at: '2026-09-13T07:30:00Z'
      }
    ] as MessDiningToken[]
  } = $props();

  let selectedMealType = $state<MealType>('LUNCH');
  let simulatedScanCode = $state<string>('');
  let scanMessage = $state<{ type: 'success' | 'error'; text: string } | null>(null);

  // Student active token for selected meal
  let currentStudentToken = $derived(
    activeTokens.find(
      (t) => t.meal_type === selectedMealType && t.student_id === 'stu_1' && t.meal_date === '2026-09-13'
    )
  );

  function handleGenerateToken() {
    if (currentStudentToken) return;

    const hash = Math.random().toString(16).substring(2, 10);
    const newToken: MessDiningToken = {
      id: `tok_${Date.now()}`,
      tenant_id: 'ten_default',
      student_id: 'stu_1',
      mess_hall_id: 'mess_1',
      meal_type: selectedMealType,
      token_code: `TOK-${hash}`,
      meal_date: '2026-09-13',
      status: 'GENERATED',
      created_at: new Date().toISOString()
    };

    activeTokens.unshift(newToken);
  }

  function handleScanToken(code: string) {
    scanMessage = null;
    const trimmed = code.trim();
    if (!trimmed) return;

    const token = activeTokens.find((t) => t.token_code === trimmed);
    if (!token) {
      scanMessage = { type: 'error', text: `❌ Invalid Token: "${trimmed}" not found in system.` };
      return;
    }

    if (token.status === 'REDEEMED') {
      scanMessage = {
        type: 'error',
        text: `⚠️ Duplicate Punch Rejected: Token ${trimmed} was already redeemed at ${new Date(
          token.redeemed_at || ''
        ).toLocaleTimeString()}.`
      };
      return;
    }

    // Redeem token
    token.status = 'REDEEMED';
    token.redeemed_at = new Date().toISOString();
    token.device_reader_id = 'READER_MAIN_01';

    scanMessage = {
      type: 'success',
      text: `✓ Punch Accepted: ${token.meal_type} token redeemed for Student ${token.student_id}. Meal serving authorized.`
    };
    simulatedScanCode = '';
  }
</script>

<div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
  <!-- Left: Student Digital Dining Pass & QR Token -->
  <div class="lg:col-span-5 space-y-4">
    <Card class="border shadow-md bg-gradient-to-b from-card to-muted/20">
      <CardHeader class="pb-3 border-b">
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="text-lg">Digital Dining Pass</CardTitle>
            <CardDescription class="text-xs">Single-use cryptographically signed meal token</CardDescription>
          </div>
          <Badge variant="outline" class="bg-primary/10 text-primary font-bold">1-TAP PUNCH</Badge>
        </div>
      </CardHeader>

      <CardContent class="p-6 space-y-5">
        <!-- Meal Slot Selector -->
        <div class="grid grid-cols-4 gap-1 bg-muted p-1 rounded-lg text-center">
          <button
            type="button"
            class="py-1.5 rounded-md text-xs font-semibold transition-colors {selectedMealType === 'BREAKFAST'
              ? 'bg-primary text-primary-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'}"
            onclick={() => (selectedMealType = 'BREAKFAST')}
          >
            Breakfast
          </button>
          <button
            type="button"
            class="py-1.5 rounded-md text-xs font-semibold transition-colors {selectedMealType === 'LUNCH'
              ? 'bg-primary text-primary-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'}"
            onclick={() => (selectedMealType = 'LUNCH')}
          >
            Lunch
          </button>
          <button
            type="button"
            class="py-1.5 rounded-md text-xs font-semibold transition-colors {selectedMealType === 'SNACKS'
              ? 'bg-primary text-primary-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'}"
            onclick={() => (selectedMealType = 'SNACKS')}
          >
            Snacks
          </button>
          <button
            type="button"
            class="py-1.5 rounded-md text-xs font-semibold transition-colors {selectedMealType === 'DINNER'
              ? 'bg-primary text-primary-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'}"
            onclick={() => (selectedMealType = 'DINNER')}
          >
            Dinner
          </button>
        </div>

        <!-- Token Display Card -->
        {#if currentStudentToken}
          <div
            class="p-6 rounded-2xl border text-center space-y-4 {currentStudentToken.status === 'GENERATED'
              ? 'bg-card border-primary/30 shadow-inner'
              : 'bg-muted/50 border-muted opacity-80'}"
          >
            <!-- QR Simulation Code Box -->
            <div class="mx-auto w-44 h-44 bg-white p-3 rounded-xl border-2 border-slate-900 shadow-sm flex flex-col items-center justify-center space-y-2">
              <div class="grid grid-cols-5 gap-1.5 p-2 bg-slate-900 rounded-lg">
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-transparent"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-transparent"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-transparent"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-transparent"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
                <div class="w-4 h-4 bg-white rounded-sm"></div>
              </div>
              <p class="font-mono text-xs font-extrabold text-slate-900 tracking-wider">
                {currentStudentToken.token_code}
              </p>
            </div>

            <div>
              <div class="flex items-center justify-center gap-2 mb-1">
                <span class="text-sm font-bold">{selectedMealType} MEAL TOKEN</span>
                {#if currentStudentToken.status === 'GENERATED'}
                  <Badge class="bg-emerald-100 text-emerald-800 border-emerald-300">ACTIVE · SCAN TO EAT</Badge>
                {:else}
                  <Badge variant="secondary" class="bg-slate-200 text-slate-800">REDEEMED ✓</Badge>
                {/if}
              </div>
              <p class="text-xs text-muted-foreground">Valid for Date: {currentStudentToken.meal_date}</p>
            </div>

            {#if currentStudentToken.status === 'GENERATED'}
              <Button
                variant="outline"
                size="sm"
                class="text-xs w-full"
                onclick={() => handleScanToken(currentStudentToken!.token_code)}
              >
                Simulate Entrance Reader Scan
              </Button>
            {/if}
          </div>
        {:else}
          <div class="p-8 rounded-2xl border border-dashed text-center space-y-3 bg-muted/20">
            <div class="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto text-primary text-xl">
              🍱
            </div>
            <p class="text-sm font-semibold">No Token Generated Yet for {selectedMealType}</p>
            <p class="text-xs text-muted-foreground">
              Click below to generate your single-use dining QR punch.
            </p>
            <Button class="w-full font-semibold mt-2" onclick={handleGenerateToken}>
              Generate {selectedMealType} Token
            </Button>
          </div>
        {/if}
      </CardContent>
    </Card>
  </div>

  <!-- Right: Caterer Entrance Scanner & Access Control Desk -->
  <div class="lg:col-span-7 space-y-4">
    <Card class="border shadow-sm">
      <CardHeader class="pb-3 border-b bg-muted/20">
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="text-base">Entrance Scanner Terminal (Caterer Desk)</CardTitle>
            <CardDescription class="text-xs">Live optical token validator & double punch detection</CardDescription>
          </div>
          <Badge variant="secondary" class="font-mono text-[10px]">READER: ENTRANCE_01</Badge>
        </div>
      </CardHeader>

      <CardContent class="p-4 space-y-4">
        <!-- Input / Barcode Scanner Field -->
        <div class="flex items-center gap-2">
          <Input
            placeholder="Scan QR or Enter Token Code (e.g. TOK-9d7a2f10b4)..."
            bind:value={simulatedScanCode}
            class="font-mono text-sm h-10"
            onkeydown={(e) => e.key === 'Enter' && handleScanToken(simulatedScanCode)}
          />
          <Button class="h-10 px-5 font-semibold" onclick={() => handleScanToken(simulatedScanCode)}>
            Verify Token
          </Button>
        </div>

        <!-- Scan Status Feedback Alert -->
        {#if scanMessage}
          <div
            class="p-4 rounded-xl border text-xs font-semibold animate-in fade-in {scanMessage.type ===
            'success'
              ? 'bg-emerald-50 text-emerald-900 border-emerald-300 dark:bg-emerald-950/30'
              : 'bg-red-50 text-red-900 border-red-300 dark:bg-red-950/30'}"
          >
            {scanMessage.text}
          </div>
        {/if}

        <!-- Live Token Log Table -->
        <div class="border rounded-xl overflow-hidden mt-4">
          <div class="bg-muted/40 px-4 py-2 text-xs font-bold text-muted-foreground border-b flex justify-between">
            <span>Recent Token Redeemed Feed</span>
            <span>Today's Total: {activeTokens.filter((t) => t.status === 'REDEEMED').length} Meals</span>
          </div>
          <div class="divide-y max-h-72 overflow-y-auto">
            {#each activeTokens as tok (tok.id)}
              <div class="p-3 text-xs flex items-center justify-between hover:bg-muted/10">
                <div class="space-y-0.5">
                  <div class="flex items-center gap-2">
                    <span class="font-mono font-bold text-foreground">{tok.token_code}</span>
                    <Badge variant="outline" class="text-[10px]">{tok.meal_type}</Badge>
                  </div>
                  <p class="text-[11px] text-muted-foreground">
                    Student ID: <span class="font-semibold text-foreground">{tok.student_id}</span> · Date: {tok.meal_date}
                  </p>
                </div>

                <div class="text-right">
                  {#if tok.status === 'REDEEMED'}
                    <span class="text-emerald-600 font-bold text-[11px]">✓ Redeemed</span>
                    <p class="text-[10px] text-muted-foreground">
                      {new Date(tok.redeemed_at || '').toLocaleTimeString()}
                    </p>
                  {:else}
                    <span class="text-amber-600 font-semibold text-[11px]">⏳ Generated</span>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</div>
