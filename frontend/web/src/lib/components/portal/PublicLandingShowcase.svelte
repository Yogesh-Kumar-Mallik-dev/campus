<!--
  BLOCK_WEB_PORTAL_SHOWCASE_001
  Subsystem: Rank 16 - Public Web Portal (portal)
  Purpose:   Public Institutional Landing Showcase, Program Offerings Catalog, and Prospect Inquiry Lead Form.
-->
<script lang="ts">
  import type { PortalLandingPage, PortalProgramCatalog, PortalDegreeType } from '$lib/types/portal';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  let landing = $state<PortalLandingPage>({
    id: 'p-landing-1',
    tenantId: 'tenant-demo',
    heroHeadline: 'Empowering Next-Gen Leaders, Innovators & Engineers',
    heroSubheadline: 'Admissions open for Academic Year 2026-2027. Explore 28+ world-class accredited undergraduate and graduate degree tracks.',
    admissionsOpen: true,
    admissionsCycle: 'Fall 2026',
    contactEmail: 'admissions@campus.edu',
    contactPhone: '+91 80 1234 5678',
    campusAddress: 'Academic City, Innovation Corridor, Bangalore - 560100',
    isPublished: true,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  });

  let programs = $state<PortalProgramCatalog[]>([
    {
      id: 'prog-1',
      tenantId: 'tenant-demo',
      programCode: 'BTECH_CSE_AI',
      programName: 'B.Tech in Computer Science & Artificial Intelligence',
      degreeType: 'UG',
      departmentName: 'Computer Science & Engineering',
      durationYears: 4,
      totalSemesters: 8,
      eligibilityCriteria: '10+2 with Physics & Mathematics (≥75% aggregate score in PCM)',
      annualFee: 250000,
      isFeatured: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'prog-2',
      tenantId: 'tenant-demo',
      programCode: 'BTECH_ECE_VLSI',
      programName: 'B.Tech in Electronics & VLSI Semiconductor Design',
      degreeType: 'UG',
      departmentName: 'Electronics & Communication',
      durationYears: 4,
      totalSemesters: 8,
      eligibilityCriteria: '10+2 with Physics & Mathematics (≥70% aggregate score)',
      annualFee: 220000,
      isFeatured: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'prog-3',
      tenantId: 'tenant-demo',
      programCode: 'MTECH_CYBER_SEC',
      programName: 'M.Tech in Cybersecurity & Cloud Resilience',
      degreeType: 'PG',
      departmentName: 'Computer Science & Engineering',
      durationYears: 2,
      totalSemesters: 4,
      eligibilityCriteria: 'B.E. / B.Tech in CSE/IT/ECE with valid GATE score or Institute CET',
      annualFee: 180000,
      isFeatured: false,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: 'prog-4',
      tenantId: 'tenant-demo',
      programCode: 'PHD_QUANTUM_COMP',
      programName: 'Ph.D. in Quantum Information & High-Performance Computing',
      degreeType: 'PHD',
      departmentName: 'Interdisciplinary Quantum Lab',
      durationYears: 3,
      totalSemesters: 6,
      eligibilityCriteria: 'Master degree in relevant STEM discipline with minimum 65% aggregate',
      annualFee: 60000,
      isFeatured: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ]);

  let selectedDegree = $state<string>('ALL');
  let searchQuery = $state<string>('');

  // Inquiry form modal
  let isInquiryModalOpen = $state<boolean>(false);
  let selectedProgramName = $state<string>('');
  let prospectName = $state<string>('');
  let prospectEmail = $state<string>('');
  let prospectPhone = $state<string>('');
  let prospectMessage = $state<string>('');
  let isSubmittedSuccess = $state<boolean>(false);

  let filteredPrograms = $derived(
    programs.filter((p) => {
      if (selectedDegree !== 'ALL' && p.degreeType !== selectedDegree) return false;
      if (searchQuery.trim() !== '') {
        const q = searchQuery.toLowerCase();
        return (
          p.programName.toLowerCase().includes(q) ||
          p.programCode.toLowerCase().includes(q) ||
          p.departmentName.toLowerCase().includes(q)
        );
      }
      return true;
    })
  );

  function openInquiry(progName: string) {
    selectedProgramName = progName;
    isInquiryModalOpen = true;
    isSubmittedSuccess = false;
  }

  function handleSendInquiry() {
    if (!prospectName.trim() || !prospectEmail.trim()) return;
    isSubmittedSuccess = true;
    setTimeout(() => {
      isInquiryModalOpen = false;
      prospectName = '';
      prospectEmail = '';
      prospectPhone = '';
      prospectMessage = '';
      isSubmittedSuccess = false;
    }, 2000);
  }
</script>

<div class="space-y-8">
  <!-- Hero Section -->
  <div class="p-8 md:p-12 bg-gradient-to-br from-slate-900 via-primary/95 to-slate-950 rounded-3xl text-white shadow-xl relative overflow-hidden">
    <div class="max-w-2xl space-y-4 relative z-10">
      <div class="flex items-center gap-2">
        <span class="px-3 py-1 bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 rounded-full text-xs font-bold uppercase tracking-wider">
          ✦ Admissions Open: {landing.admissionsCycle}
        </span>
      </div>

      <h1 class="text-3xl md:text-5xl font-black tracking-tight leading-tight">
        {landing.heroHeadline}
      </h1>

      <p class="text-sm md:text-base text-slate-200 leading-relaxed">
        {landing.heroSubheadline}
      </p>

      <div class="pt-3 flex flex-wrap items-center gap-4 text-xs text-slate-300">
        <span>📍 {landing.campusAddress}</span>
        <span>✉️ {landing.contactEmail}</span>
        <span>📞 {landing.contactPhone}</span>
      </div>

      <div class="pt-2 flex gap-3">
        <Button size="lg" onclick={() => openInquiry('General Admissions')} class="bg-white text-slate-900 hover:bg-slate-100 font-extrabold shadow-md">
          Apply / Request Info →
        </Button>
      </div>
    </div>
  </div>

  <!-- Program Catalog Explorer Header -->
  <div class="space-y-4">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black text-foreground">Academic Program Catalog</h2>
        <p class="text-xs text-muted-foreground mt-0.5">
          Curated industry-aligned degrees with flexible semesters and global transfer credits.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <Input
          placeholder="Search programs, departments..."
          bind:value={searchQuery}
          class="w-full sm:w-64 text-xs"
        />
      </div>
    </div>

    <!-- Degree Type Tabs -->
    <div class="flex flex-wrap gap-2 pt-1">
      {#each ['ALL', 'UG', 'PG', 'PHD'] as deg}
        <button
          onclick={() => (selectedDegree = deg)}
          class="px-4 py-1.5 rounded-xl text-xs font-bold transition-colors border {selectedDegree === deg ? 'bg-primary text-primary-foreground border-primary' : 'bg-card text-muted-foreground hover:bg-muted'}"
        >
          {deg === 'ALL' ? 'All Degrees' : deg === 'UG' ? 'Undergraduate (UG)' : deg === 'PG' ? 'Postgraduate (PG)' : 'Doctoral (Ph.D.)'}
        </button>
      {/each}
    </div>
  </div>

  <!-- Programs Grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
    {#each filteredPrograms as prog}
      <div class="p-6 bg-card rounded-2xl border shadow-sm flex flex-col justify-between space-y-4 hover:border-primary/40 transition-colors">
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <span class="font-mono text-xs font-bold text-primary bg-primary/10 px-2 py-0.5 rounded">
              {prog.programCode}
            </span>
            {#if prog.isFeatured}
              <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-amber-100 text-amber-800 border border-amber-200">
                ★ FEATURED TRACK
              </span>
            {/if}
          </div>

          <h3 class="text-lg font-bold text-foreground leading-snug">
            {prog.programName}
          </h3>

          <div class="text-xs text-muted-foreground space-y-1">
            <div>🏛️ Department: <strong>{prog.departmentName}</strong></div>
            <div>⏱️ Duration: <strong>{prog.durationYears} Years ({prog.totalSemesters} Semesters)</strong></div>
            <div class="pt-1 text-[11px] leading-relaxed text-muted-foreground">
              🎓 Eligibility: {prog.eligibilityCriteria}
            </div>
          </div>
        </div>

        <div class="pt-3 border-t flex items-center justify-between">
          <div>
            <span class="text-[10px] text-muted-foreground block">Annual Tuition</span>
            <span class="text-base font-black text-foreground">₹{prog.annualFee.toLocaleString()}</span>
          </div>

          <Button size="sm" onclick={() => openInquiry(prog.programName)} class="font-bold">
            Inquire Now
          </Button>
        </div>
      </div>
    {/each}
  </div>

  <!-- Inquiry Lead Capture Modal -->
  {#if isInquiryModalOpen}
    <div class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-card w-full max-w-md rounded-2xl border shadow-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b pb-3">
          <div>
            <h2 class="text-lg font-bold text-foreground">Admissions Information Request</h2>
            <p class="text-xs text-muted-foreground">{selectedProgramName}</p>
          </div>
          <button onclick={() => (isInquiryModalOpen = false)} class="text-muted-foreground hover:text-foreground">✕</button>
        </div>

        {#if !isSubmittedSuccess}
          <div class="space-y-3">
            <div>
              <label for="prospect-name" class="text-xs font-semibold text-muted-foreground block mb-1">Your Full Name</label>
              <Input id="prospect-name" bind:value={prospectName} placeholder="e.g. Rahul Sharma" />
            </div>

            <div>
              <label for="prospect-email" class="text-xs font-semibold text-muted-foreground block mb-1">Email Address</label>
              <Input id="prospect-email" type="email" bind:value={prospectEmail} placeholder="e.g. rahul@example.com" />
            </div>

            <div>
              <label for="prospect-phone" class="text-xs font-semibold text-muted-foreground block mb-1">Phone Number</label>
              <Input id="prospect-phone" bind:value={prospectPhone} placeholder="e.g. +91 98765 43210" />
            </div>

            <div>
              <label for="prospect-msg" class="text-xs font-semibold text-muted-foreground block mb-1">Questions or Specific Inquiries</label>
              <textarea
                id="prospect-msg"
                bind:value={prospectMessage}
                rows={3}
                placeholder="Inquire about scholarships, hostel facilities, entrance criteria..."
                class="w-full bg-background text-foreground text-xs rounded-xl border p-3"
              ></textarea>
            </div>
          </div>

          <div class="flex items-center justify-end gap-3 pt-3 border-t">
            <Button variant="outline" onclick={() => (isInquiryModalOpen = false)}>Cancel</Button>
            <Button onclick={handleSendInquiry} class="font-bold">Submit Request</Button>
          </div>
        {:else}
          <div class="p-6 text-center space-y-2">
            <div class="text-3xl text-emerald-600 font-bold">✓</div>
            <h3 class="text-base font-bold text-foreground">Inquiry Received!</h3>
            <p class="text-xs text-muted-foreground">
              Our admissions counseling desk has received your request and will contact you within 24 hours.
            </p>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>
