<script lang="ts">
  /**
   * BLOCK_WEB_STUDYHUB_MATERIALS_001
   * Subsystem: Rank 10 - Study Hub System (studyhub)
   * Purpose:   Interactive course syllabus viewer, unit-wise notes & lab manuals repository, and faculty publishing desk.
   */
  import type { StudyCourse, StudyMaterial, StudyMaterialType } from '$lib/types/studyhub';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Badge } from '$lib/components/ui/badge';

  let courses = $state<StudyCourse[]>([
    {
      id: 'crs-01',
      tenantId: 'tenant-main',
      code: 'CS-301',
      title: 'Advanced Operating Systems',
      description: 'Microkernel design, virtual memory management, CPU scheduling, and synchronization.',
      department: 'Computer Science',
      credits: 4,
      semester: 5,
      syllabusText: 'Unit 1: Kernel & System Calls\nUnit 2: Memory Paging & Swapping\nUnit 3: File System Drivers\nUnit 4: Concurrency & Lock-Free Structures',
      isActive: true,
      createdAt: '2026-08-01T00:00:00Z',
      updatedAt: '2026-08-01T00:00:00Z'
    },
    {
      id: 'crs-02',
      tenantId: 'tenant-main',
      code: 'CS-302',
      title: 'Database Engine Architecture',
      description: 'Storage layouts, buffer pool management, B-Tree index structures, and ACID transactions.',
      department: 'Computer Science',
      credits: 4,
      semester: 5,
      syllabusText: 'Unit 1: Buffer Pool & Page Formats\nUnit 2: B+ Trees & LSM Trees\nUnit 3: Query Execution Engines\nUnit 4: Two-Phase Locking & MVCC',
      isActive: true,
      createdAt: '2026-08-01T00:00:00Z',
      updatedAt: '2026-08-01T00:00:00Z'
    }
  ]);

  let selectedCourseId = $state<string>('crs-01');
  let selectedUnit = $state<number | null>(null);

  let materials = $state<StudyMaterial[]>([
    {
      id: 'mat-01',
      tenantId: 'tenant-main',
      courseId: 'crs-01',
      title: 'Lecture 01: Microkernel vs Monolithic Architectures',
      description: 'Comparative deep-dive into Mach, Linux, and seL4 microkernel architectures.',
      unitNumber: 1,
      materialType: 'SLIDE',
      fileUrl: 'https://storage.campus.internal/study/cs301-lec01.pdf',
      fileSizeBytes: 2450000,
      publishedAt: '2026-08-10T10:00:00Z',
      createdAt: '2026-08-10T10:00:00Z',
      updatedAt: '2026-08-10T10:00:00Z'
    },
    {
      id: 'mat-02',
      tenantId: 'tenant-main',
      courseId: 'crs-01',
      title: 'Lab 01: Custom Kernel Module & ProcFS Handler in C',
      description: 'Step-by-step kernel programming assignment with debugging via QEMU.',
      unitNumber: 1,
      materialType: 'LAB_MANUAL',
      fileUrl: 'https://storage.campus.internal/study/cs301-lab01.pdf',
      fileSizeBytes: 1850000,
      publishedAt: '2026-08-12T14:30:00Z',
      createdAt: '2026-08-12T14:30:00Z',
      updatedAt: '2026-08-12T14:30:00Z'
    },
    {
      id: 'mat-03',
      tenantId: 'tenant-main',
      courseId: 'crs-01',
      title: 'Lecture 05: Multi-Level Page Tables & TLB Shootdowns',
      description: 'x86_64 4-level and 5-level paging implementation details.',
      unitNumber: 2,
      materialType: 'NOTE',
      fileUrl: 'https://storage.campus.internal/study/cs301-lec05.pdf',
      fileSizeBytes: 3100000,
      publishedAt: '2026-08-20T09:00:00Z',
      createdAt: '2026-08-20T09:00:00Z',
      updatedAt: '2026-08-20T09:00:00Z'
    }
  ]);

  // Publish Material Form
  let isPublishModalOpen = $state(false);
  let newMatTitle = $state('');
  let newMatDesc = $state('');
  let newMatUnit = $state(1);
  let newMatType = $state<StudyMaterialType>('NOTE');
  let newMatUrl = $state('');

  const activeCourse = $derived(courses.find((c) => c.id === selectedCourseId));
  const filteredMaterials = $derived(
    materials.filter(
      (m) => m.courseId === selectedCourseId && (selectedUnit === null || m.unitNumber === selectedUnit)
    )
  );

  function handlePublish() {
    if (!newMatTitle || !newMatUrl) return;
    const newMat: StudyMaterial = {
      id: `mat-${Date.now()}`,
      tenantId: 'tenant-main',
      courseId: selectedCourseId,
      title: newMatTitle,
      description: newMatDesc,
      unitNumber: Number(newMatUnit) || 1,
      materialType: newMatType,
      fileUrl: newMatUrl,
      fileSizeBytes: 1500000,
      publishedAt: new Date().toISOString(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    materials = [newMat, ...materials];
    newMatTitle = '';
    newMatDesc = '';
    newMatUrl = '';
    isPublishModalOpen = false;
  }

  function getBadgeVariant(type: StudyMaterialType) {
    switch (type) {
      case 'NOTE':
        return 'secondary';
      case 'SLIDE':
        return 'default';
      case 'LAB_MANUAL':
        return 'outline';
      case 'RECORDING':
        return 'destructive';
      default:
        return 'secondary';
    }
  }

  function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }
</script>

<div class="space-y-6">
  <!-- Course Selection Bar -->
  <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 bg-card p-4 rounded-xl border">
    <div class="flex items-center gap-3">
      <Label for="courseSelect" class="font-semibold text-sm">Active Course:</Label>
      <select
        id="courseSelect"
        bind:value={selectedCourseId}
        class="bg-background border rounded-lg px-3 py-2 text-sm font-medium focus:ring-2 focus:ring-primary"
      >
        {#each courses as course}
          <option value={course.id}>{course.code} — {course.title}</option>
        {/each}
      </select>
    </div>

    <div class="flex items-center gap-2">
      <Button variant="outline" size="sm" onclick={() => (selectedUnit = null)}>
        All Units
      </Button>
      {#each [1, 2, 3, 4] as u}
        <Button
          variant={selectedUnit === u ? 'default' : 'ghost'}
          size="sm"
          onclick={() => (selectedUnit = u)}
        >
          Unit {u}
        </Button>
      {/each}
      <Button size="sm" onclick={() => (isPublishModalOpen = true)}>
        + Publish Material
      </Button>
    </div>
  </div>

  <!-- Course Overview Card -->
  {#if activeCourse}
    <div class="bg-gradient-to-r from-slate-900 to-slate-800 text-white p-6 rounded-2xl shadow-sm space-y-3">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div>
          <span class="text-xs uppercase tracking-widest text-slate-400 font-bold">
            {activeCourse.department} · Semester {activeCourse.semester} · {activeCourse.credits} Credits
          </span>
          <h2 class="text-2xl font-black mt-0.5">{activeCourse.code}: {activeCourse.title}</h2>
        </div>
        <Badge variant="secondary" class="bg-emerald-500/20 text-emerald-300 border-emerald-500/30">
          ACTIVE SYLLABUS
        </Badge>
      </div>
      <p class="text-xs text-slate-300 leading-relaxed max-w-3xl">
        {activeCourse.description}
      </p>
      {#if activeCourse.syllabusText}
        <div class="mt-3 pt-3 border-t border-slate-700/60 text-xs text-slate-300">
          <span class="font-bold text-slate-100">Syllabus Outline:</span>
          <pre class="font-sans text-xs mt-1 whitespace-pre-line text-slate-300">{activeCourse.syllabusText}</pre>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Materials Grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#each filteredMaterials as mat}
      <div class="bg-card border rounded-2xl p-5 flex flex-col justify-between space-y-4 hover:shadow-md transition-shadow">
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <Badge variant={getBadgeVariant(mat.materialType)}>
              {mat.materialType}
            </Badge>
            <span class="text-[11px] font-semibold text-muted-foreground">
              Unit {mat.unitNumber}
            </span>
          </div>
          <h3 class="text-sm font-bold text-foreground line-clamp-2">{mat.title}</h3>
          {#if mat.description}
            <p class="text-xs text-muted-foreground line-clamp-2">{mat.description}</p>
          {/if}
        </div>

        <div class="pt-3 border-t flex items-center justify-between text-xs">
          <span class="text-muted-foreground font-mono">{formatBytes(mat.fileSizeBytes)}</span>
          <a
            href={mat.fileUrl}
            target="_blank"
            rel="noopener noreferrer"
            class="text-primary font-semibold hover:underline flex items-center gap-1 text-xs"
          >
            📥 Download / View
          </a>
        </div>
      </div>
    {/each}
  </div>

  {#if filteredMaterials.length === 0}
    <div class="p-12 text-center border-2 border-dashed rounded-2xl text-muted-foreground space-y-2">
      <p class="font-semibold text-sm">No materials published for this unit yet.</p>
      <p class="text-xs">Faculty instructors can upload lecture slides and lab manuals above.</p>
    </div>
  {/if}
</div>

<!-- Modal: Publish Material -->
{#if isPublishModalOpen}
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
    <div class="bg-card border max-w-lg w-full rounded-2xl p-6 shadow-2xl space-y-4">
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-lg font-bold text-foreground">Publish Study Resource</h3>
        <button
          onclick={() => (isPublishModalOpen = false)}
          class="text-muted-foreground hover:text-foreground text-sm font-bold"
        >
          ✕
        </button>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <Label for="matTitle" class="text-xs font-semibold">Material Title</Label>
          <Input id="matTitle" bind:value={newMatTitle} placeholder="e.g. Lecture 02: Page Replacement Algorithms" />
        </div>

        <div class="space-y-1">
          <Label for="matDesc" class="text-xs font-semibold">Description (Optional)</Label>
          <Input id="matDesc" bind:value={newMatDesc} placeholder="Key concepts, homework references..." />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <Label for="matUnit" class="text-xs font-semibold">Syllabus Unit #</Label>
            <Input id="matUnit" type="number" min="1" max="10" bind:value={newMatUnit} />
          </div>

          <div class="space-y-1">
            <Label for="matType" class="text-xs font-semibold">Resource Type</Label>
            <select
              id="matType"
              bind:value={newMatType}
              class="w-full bg-background border rounded-lg px-3 py-2 text-xs font-medium"
            >
              <option value="NOTE">Lecture Notes (PDF)</option>
              <option value="SLIDE">Slide Presentation</option>
              <option value="LAB_MANUAL">Lab Manual & Code</option>
              <option value="RECORDING">Session Recording</option>
              <option value="SAMPLE_PAPER">Sample Question Paper</option>
            </select>
          </div>
        </div>

        <div class="space-y-1">
          <Label for="matUrl" class="text-xs font-semibold">CDN / Cloud Storage URL</Label>
          <Input id="matUrl" bind:value={newMatUrl} placeholder="https://storage.campus.internal/study/..." />
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-3 border-t">
        <Button variant="outline" size="sm" onclick={() => (isPublishModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handlePublish}>
          Publish to Students
        </Button>
      </div>
    </div>
  </div>
{/if}
