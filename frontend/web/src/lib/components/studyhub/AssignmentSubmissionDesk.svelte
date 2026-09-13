<script lang="ts">
  /**
   * BLOCK_WEB_STUDYHUB_ASSIGNMENTS_001
   * Subsystem: Rank 10 - Study Hub System (studyhub)
   * Purpose:   Interactive assignment submission desk, faculty grading matrix with late penalty computations, and peer reviews.
   */
  import type { StudyAssignment, StudySubmission, StudyPeerReview } from '$lib/types/studyhub';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Badge } from '$lib/components/ui/badge';

  let assignments = $state<StudyAssignment[]>([
    {
      id: 'asgn-01',
      tenantId: 'tenant-main',
      courseId: 'crs-01',
      title: 'Lab 01: Kernel Module & ProcFS Debugger',
      description: 'Implement a loadable kernel module in C that exposes memory page telemetry via /proc/mem_telemetry.',
      maxMarks: 100,
      dueDate: '2026-09-18T23:59:59Z',
      allowLateSubmission: true,
      latePenaltyPercentPerDay: 5,
      status: 'PUBLISHED',
      createdAt: '2026-09-01T00:00:00Z',
      updatedAt: '2026-09-01T00:00:00Z'
    },
    {
      id: 'asgn-02',
      tenantId: 'tenant-main',
      courseId: 'crs-01',
      title: 'Assignment 2: Lock-Free Multi-Producer Queue',
      description: 'Design and benchmark a lock-free ring buffer queue using C11 atomic compare-and-swap (CAS).',
      maxMarks: 50,
      dueDate: '2026-09-25T23:59:59Z',
      allowLateSubmission: false,
      latePenaltyPercentPerDay: 0,
      status: 'PUBLISHED',
      createdAt: '2026-09-05T00:00:00Z',
      updatedAt: '2026-09-05T00:00:00Z'
    }
  ]);

  let selectedAsgnId = $state<string>('asgn-01');

  let submissions = $state<StudySubmission[]>([
    {
      id: 'sub-01',
      tenantId: 'tenant-main',
      assignmentId: 'asgn-01',
      studentId: 'stu-001',
      studentName: 'Aarav Sharma',
      rollNumber: 'CS22B001',
      fileUrl: 'https://storage.campus.internal/subs/cs22b001-lab1.tar.gz',
      contentText: 'Kernel module tested on Linux 6.8 with 10,000 page reads.',
      submittedAt: '2026-09-15T18:20:00Z',
      status: 'GRADED',
      marksObtained: 96,
      feedback: 'Excellent clean code with no kernel memory leaks reported by kmemleak.',
      gradedById: 'prof-dev-01',
      gradedAt: '2026-09-16T10:00:00Z',
      createdAt: '2026-09-15T18:20:00Z',
      updatedAt: '2026-09-16T10:00:00Z'
    },
    {
      id: 'sub-02',
      tenantId: 'tenant-main',
      assignmentId: 'asgn-01',
      studentId: 'stu-002',
      studentName: 'Diya Patel',
      rollNumber: 'CS22B014',
      fileUrl: 'https://storage.campus.internal/subs/cs22b014-lab1.zip',
      contentText: 'Includes unit tests in C with QEMU boot scripts.',
      submittedAt: '2026-09-16T09:12:00Z',
      status: 'SUBMITTED',
      createdAt: '2026-09-16T09:12:00Z',
      updatedAt: '2026-09-16T09:12:00Z'
    }
  ]);

  // Peer Reviews
  let peerReviews = $state<StudyPeerReview[]>([
    {
      id: 'prev-01',
      submissionId: 'sub-01',
      reviewerStudentId: 'stu-002',
      score: 95,
      comments: 'Very well documented Makefile and error cleanup paths in module_exit.',
      reviewedAt: '2026-09-16T12:00:00Z'
    }
  ]);

  // Submit Assignment Dialog
  let isSubmitModalOpen = $state(false);
  let myFileUrl = $state('');
  let myContent = $state('');

  // Grading Dialog
  let isGradeModalOpen = $state(false);
  let activeGradingSub = $state<StudySubmission | null>(null);
  let inputRawMarks = $state(90);
  let inputFeedback = $state('');

  const activeAsgn = $derived(assignments.find((a) => a.id === selectedAsgnId));
  const activeSubmissions = $derived(
    submissions.filter((s) => s.assignmentId === selectedAsgnId)
  );

  function handleSubmitAssignment() {
    if (!myFileUrl && !myContent) return;
    const newSub: StudySubmission = {
      id: `sub-${Date.now()}`,
      tenantId: 'tenant-main',
      assignmentId: selectedAsgnId,
      studentId: 'stu-me',
      studentName: 'You (Current Student)',
      rollNumber: 'CS22B099',
      fileUrl: myFileUrl,
      contentText: myContent,
      submittedAt: new Date().toISOString(),
      status: 'SUBMITTED',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    submissions = [newSub, ...submissions];
    myFileUrl = '';
    myContent = '';
    isSubmitModalOpen = false;
  }

  function openGrading(sub: StudySubmission) {
    activeGradingSub = sub;
    inputRawMarks = sub.marksObtained ?? 90;
    inputFeedback = sub.feedback ?? '';
    isGradeModalOpen = true;
  }

  function handleSaveGrade() {
    if (!activeGradingSub) return;
    submissions = submissions.map((s) => {
      if (s.id === activeGradingSub!.id) {
        return {
          ...s,
          status: 'GRADED',
          marksObtained: Number(inputRawMarks),
          feedback: inputFeedback,
          gradedById: 'prof-current',
          gradedAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        };
      }
      return s;
    });
    isGradeModalOpen = false;
    activeGradingSub = null;
  }

  function getSubmissionStatusBadge(status: string) {
    switch (status) {
      case 'GRADED':
        return 'bg-emerald-100 text-emerald-800 border-emerald-200';
      case 'SUBMITTED':
        return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'LATE':
        return 'bg-amber-100 text-amber-800 border-amber-200';
      case 'REJECTED':
        return 'bg-rose-100 text-rose-800 border-rose-200';
      default:
        return 'bg-slate-100 text-slate-800 border-slate-200';
    }
  }
</script>

<div class="space-y-6">
  <!-- Assignment Selector & Actions Bar -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-card p-4 rounded-xl border">
    <div class="flex items-center gap-3">
      <Label for="asgnSelect" class="font-semibold text-sm">Select Assignment:</Label>
      <select
        id="asgnSelect"
        bind:value={selectedAsgnId}
        class="bg-background border rounded-lg px-3 py-2 text-sm font-medium focus:ring-2 focus:ring-primary"
      >
        {#each assignments as asgn}
          <option value={asgn.id}>{asgn.title} (Max {asgn.maxMarks} pts)</option>
        {/each}
      </select>
    </div>

    <div class="flex items-center gap-2">
      <Button size="sm" onclick={() => (isSubmitModalOpen = true)}>
        📤 Turn in Submission
      </Button>
    </div>
  </div>

  <!-- Active Assignment Spec Box -->
  {#if activeAsgn}
    <div class="bg-card border rounded-2xl p-6 shadow-sm space-y-4">
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div>
          <span class="text-xs uppercase tracking-wider font-bold text-muted-foreground">
            Due Date: {new Date(activeAsgn.dueDate).toLocaleDateString()} at {new Date(activeAsgn.dueDate).toLocaleTimeString()}
          </span>
          <h2 class="text-xl font-black text-foreground mt-1">{activeAsgn.title}</h2>
        </div>
        <div class="flex items-center gap-2">
          <Badge variant="outline">Max {activeAsgn.maxMarks} Marks</Badge>
          {#if activeAsgn.allowLateSubmission}
            <Badge variant="secondary">
              Late Allowed (-{activeAsgn.latePenaltyPercentPerDay}%/day)
            </Badge>
          {:else}
            <Badge variant="destructive">Strict Deadline (No Late)</Badge>
          {/if}
        </div>
      </div>

      <p class="text-xs text-muted-foreground leading-relaxed">
        {activeAsgn.description}
      </p>
    </div>
  {/if}

  <!-- Submissions & Grading Matrix -->
  <div class="bg-card border rounded-2xl overflow-hidden shadow-sm">
    <div class="p-4 border-b flex items-center justify-between">
      <div>
        <h3 class="text-sm font-bold text-foreground">Turned In Submissions ({activeSubmissions.length})</h3>
        <p class="text-xs text-muted-foreground">Review student code archives and assign grade evaluations.</p>
      </div>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead class="bg-muted/50 text-muted-foreground font-semibold uppercase tracking-wider border-b">
          <tr>
            <th class="p-3">Student</th>
            <th class="p-3">Turned In</th>
            <th class="p-3">Status</th>
            <th class="p-3">Marks</th>
            <th class="p-3">Artifact</th>
            <th class="p-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          {#each activeSubmissions as sub}
            <tr class="hover:bg-muted/30 transition-colors">
              <td class="p-3">
                <div class="font-bold text-foreground">{sub.studentName ?? sub.studentId}</div>
                <div class="text-[10px] text-muted-foreground font-mono">{sub.rollNumber ?? 'N/A'}</div>
              </td>
              <td class="p-3 text-muted-foreground">
                {new Date(sub.submittedAt).toLocaleDateString()}
                <span class="text-[10px] text-slate-400 block">{new Date(sub.submittedAt).toLocaleTimeString()}</span>
              </td>
              <td class="p-3">
                <span class="px-2 py-0.5 rounded text-[10px] font-bold border {getSubmissionStatusBadge(sub.status)}">
                  {sub.status}
                </span>
              </td>
              <td class="p-3 font-semibold">
                {#if sub.marksObtained !== undefined}
                  <span class="text-emerald-700 font-extrabold">{sub.marksObtained}</span> / {activeAsgn?.maxMarks}
                {:else}
                  <span class="text-muted-foreground italic">Pending</span>
                {/if}
              </td>
              <td class="p-3">
                {#if sub.fileUrl}
                  <a
                    href={sub.fileUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-primary font-semibold hover:underline"
                  >
                    📦 View Code Archive
                  </a>
                {:else}
                  <span class="text-muted-foreground italic">Text only</span>
                {/if}
              </td>
              <td class="p-3 text-right">
                <Button size="sm" variant="outline" onclick={() => openGrading(sub)}>
                  {sub.status === 'GRADED' ? 'Edit Grade' : 'Grade'}
                </Button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>

<!-- Modal: Turn In Submission -->
{#if isSubmitModalOpen}
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
    <div class="bg-card border max-w-lg w-full rounded-2xl p-6 shadow-2xl space-y-4">
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-lg font-bold text-foreground">Turn In Assignment</h3>
        <button
          onclick={() => (isSubmitModalOpen = false)}
          class="text-muted-foreground hover:text-foreground text-sm font-bold"
        >
          ✕
        </button>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <Label for="subFile" class="text-xs font-semibold">Code Archive URL (Zip / Tar.gz)</Label>
          <Input id="subFile" bind:value={myFileUrl} placeholder="https://storage.campus.internal/subs/..." />
        </div>

        <div class="space-y-1">
          <Label for="subText" class="text-xs font-semibold">Submission Notes / README Summary</Label>
          <textarea
            id="subText"
            bind:value={myContent}
            rows="3"
            class="w-full bg-background border rounded-lg p-2.5 text-xs focus:ring-2 focus:ring-primary"
            placeholder="Describe build instructions, test run commands, or system requirements..."
          ></textarea>
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-3 border-t">
        <Button variant="outline" size="sm" onclick={() => (isSubmitModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleSubmitAssignment}>
          Confirm & Turn In
        </Button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal: Grade Submission -->
{#if isGradeModalOpen && activeGradingSub}
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
    <div class="bg-card border max-w-lg w-full rounded-2xl p-6 shadow-2xl space-y-4">
      <div class="flex items-center justify-between border-b pb-3">
        <div>
          <h3 class="text-lg font-bold text-foreground">Grade Submission</h3>
          <p class="text-xs text-muted-foreground">
            Student: {activeGradingSub.studentName} ({activeGradingSub.rollNumber})
          </p>
        </div>
        <button
          onclick={() => (isGradeModalOpen = false)}
          class="text-muted-foreground hover:text-foreground text-sm font-bold"
        >
          ✕
        </button>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <Label for="rawMarks" class="text-xs font-semibold">Marks (Max: {activeAsgn?.maxMarks})</Label>
          <Input id="rawMarks" type="number" min="0" max={activeAsgn?.maxMarks} bind:value={inputRawMarks} />
        </div>

        <div class="space-y-1">
          <Label for="feedText" class="text-xs font-semibold">Evaluator Feedback</Label>
          <textarea
            id="feedText"
            bind:value={inputFeedback}
            rows="3"
            class="w-full bg-background border rounded-lg p-2.5 text-xs focus:ring-2 focus:ring-primary"
            placeholder="Detailed comments on algorithm efficiency, edge cases, and code style..."
          ></textarea>
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-3 border-t">
        <Button variant="outline" size="sm" onclick={() => (isGradeModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleSaveGrade}>
          Post Grade & Notify Student
        </Button>
      </div>
    </div>
  </div>
{/if}
