<!--
  BLOCK_WEB_HUB_CONFIG_001
  Subsystem: Rank 17 - The Hub Root Super-App (hub)
  Purpose:   Dashboard Widget Customizer and Layout Theme Configurator for personalizing user cockpits.
-->
<script lang="ts">
  import type { HubWidgetConfig, HubWidgetSize, HubWidgetCategory } from '$lib/types/hub';
  import { Button } from '$lib/components/ui/button';

  let selectedTheme = $state<string>('default');

  let widgets = $state<HubWidgetConfig[]>([
    {
      id: 'w-1',
      tenantId: 'tenant-demo',
      dashboardId: 'dash-1',
      widgetKey: 'kpi_overview',
      title: 'Vital KPIs Overview (Attendance, Dues, Borrows)',
      category: 'METRICS',
      size: 'FULL_WIDTH',
      orderIndex: 1,
      isEnabled: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'w-2',
      tenantId: 'tenant-demo',
      dashboardId: 'dash-1',
      widgetKey: 'quick_actions',
      title: '1-Tap Quick Action Shortcuts',
      category: 'ACTIONS',
      size: 'LARGE',
      orderIndex: 2,
      isEnabled: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'w-3',
      tenantId: 'tenant-demo',
      dashboardId: 'dash-1',
      widgetKey: 'safety_emergency',
      title: 'Campus Safety & SOS Live Radar',
      category: 'SAFETY',
      size: 'MEDIUM',
      orderIndex: 3,
      isEnabled: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'w-4',
      tenantId: 'tenant-demo',
      dashboardId: 'dash-1',
      widgetKey: 'telemetry_status',
      title: 'Subsystems Telemetry Grid',
      category: 'METRICS',
      size: 'MEDIUM',
      orderIndex: 4,
      isEnabled: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'w-5',
      tenantId: 'tenant-demo',
      dashboardId: 'dash-1',
      widgetKey: 'compliance_vault',
      title: 'UGC Vigilance & Audit Checkpoints',
      category: 'COMMUNICATION',
      size: 'SMALL',
      orderIndex: 5,
      isEnabled: false,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ]);

  let isSaved = $state<boolean>(false);

  function toggleWidget(w: HubWidgetConfig) {
    w.isEnabled = !w.isEnabled;
  }

  function setSize(w: HubWidgetConfig, s: HubWidgetSize) {
    w.size = s;
  }

  function moveWidget(idx: number, dir: 'UP' | 'DOWN') {
    if (dir === 'UP' && idx > 0) {
      const temp = widgets[idx];
      widgets[idx] = widgets[idx - 1];
      widgets[idx - 1] = temp;
    } else if (dir === 'DOWN' && idx < widgets.length - 1) {
      const temp = widgets[idx];
      widgets[idx] = widgets[idx + 1];
      widgets[idx + 1] = temp;
    }
  }

  function handleSave() {
    isSaved = true;
    setTimeout(() => {
      isSaved = false;
    }, 2500);
  }
</script>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h2 class="text-xl font-black text-foreground">Cockpit Widget Customizer</h2>
      <p class="text-xs text-muted-foreground mt-0.5">
        Rearrange dashboard cards, toggle subsystem modules, and adjust grid layout proportions.
      </p>
    </div>

    <Button onclick={handleSave} class="font-bold">
      {isSaved ? '✓ Layout Saved!' : 'Save Dashboard Preferences'}
    </Button>
  </div>

  <!-- Theme selector -->
  <div class="p-4 bg-card rounded-2xl border space-y-2">
    <span class="text-xs font-bold text-foreground">Cockpit Visual Theme</span>
    <div class="flex flex-wrap gap-2">
      {#each [
        { id: 'default', label: 'Default Enterprise' },
        { id: 'compact_dark', label: 'High-Density Compact' },
        { id: 'grid_expanded', label: 'Expanded Modular Grid' },
      ] as th}
        <button
          onclick={() => (selectedTheme = th.id)}
          class="px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all border {selectedTheme === th.id ? 'bg-primary text-primary-foreground border-primary' : 'bg-muted/50 text-muted-foreground border-border hover:bg-muted'}"
        >
          {th.label}
        </button>
      {/each}
    </div>
  </div>

  <!-- Widget Card Reorder List -->
  <div class="space-y-3">
    {#each widgets as w, idx}
      <div class="p-4 bg-card rounded-2xl border shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4 {w.isEnabled ? '' : 'opacity-50'}">
        <div class="space-y-1">
          <div class="flex items-center gap-2">
            <span class="text-xs font-mono font-bold text-primary px-2 py-0.5 bg-primary/10 rounded">
              #{idx + 1}
            </span>
            <h3 class="text-sm font-bold text-foreground">{w.title}</h3>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-muted text-muted-foreground border">
              {w.category}
            </span>
          </div>
          <p class="text-xs text-muted-foreground font-mono">
            Key: {w.widgetKey} · Current Span: <strong>{w.size}</strong>
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <!-- Size Selector -->
          <div class="flex gap-1 bg-muted p-1 rounded-xl">
            {#each ['SMALL', 'MEDIUM', 'LARGE', 'FULL_WIDTH'] as sz}
              <button
                onclick={() => setSize(w, sz as HubWidgetSize)}
                class="px-2 py-0.5 rounded-lg text-[10px] font-bold transition-all {w.size === sz ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'}"
              >
                {sz === 'FULL_WIDTH' ? 'FULL' : sz}
              </button>
            {/each}
          </div>

          <!-- Reorder buttons -->
          <div class="flex gap-1">
            <Button
              variant="outline"
              size="sm"
              disabled={idx === 0}
              onclick={() => moveWidget(idx, 'UP')}
              class="h-8 px-2 text-xs"
            >
              ▲
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={idx === widgets.length - 1}
              onclick={() => moveWidget(idx, 'DOWN')}
              class="h-8 px-2 text-xs"
            >
              ▼
            </Button>
          </div>

          <!-- Enable Toggle -->
          <Button
            variant={w.isEnabled ? 'default' : 'secondary'}
            size="sm"
            onclick={() => toggleWidget(w)}
            class="h-8 text-xs font-bold"
          >
            {w.isEnabled ? 'Enabled' : 'Disabled'}
          </Button>
        </div>
      </div>
    {/each}
  </div>
</div>
