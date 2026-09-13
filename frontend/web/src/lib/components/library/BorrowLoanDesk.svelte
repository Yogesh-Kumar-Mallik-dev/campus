<script lang="ts">
  /**
   * BLOCK_WEB_LIBRARY_BORROW_DESK_001
   * Subsystem: Rank 9 - E-Library System (library)
   * Purpose:   Circulation desk, student book checkout, dynamic overdue fine calculator, and renewals.
   */
  import type { LibraryBorrowRecord, BorrowStatus } from '$lib/types/library';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let {
    borrows = [
      {
        id: 'brw_001',
        tenant_id: 'ten_default',
        copy_id: 'cpy_101',
        book_title: 'The C Programming Language (2nd Edition)',
        isbn: '978-0131103627',
        accession_number: 'LIB-ACC-00101',
        student_id: 'stu_1',
        student_name: 'Rahul Sharma',
        roll_number: '2026CSE001',
        borrowed_at: '2026-09-01T10:00:00Z',
        due_date: '2026-09-15T18:00:00Z',
        renew_count: 0,
        status: 'ISSUED',
        fine_amount: 0,
        fine_paid: false,
        created_at: '2026-09-01T10:00:00Z'
      },
      {
        id: 'brw_002',
        tenant_id: 'ten_default',
        copy_id: 'cpy_102',
        book_title: 'Introduction to Algorithms (CLRS)',
        isbn: '978-0262033848',
        accession_number: 'LIB-ACC-00201',
        student_id: 'stu_2',
        student_name: 'Aditya Verma',
        roll_number: '2026CSE002',
        borrowed_at: '2026-08-20T11:00:00Z',
        due_date: '2026-09-03T18:00:00Z',
        renew_count: 1,
        status: 'OVERDUE',
        fine_amount: 50,
        fine_paid: false,
        created_at: '2026-08-20T11:00:00Z'
      },
      {
        id: 'brw_003',
        tenant_id: 'ten_default',
        copy_id: 'cpy_103',
        book_title: 'Clean Code: A Handbook of Agile Craftsmanship',
        isbn: '978-0132350884',
        accession_number: 'LIB-ACC-00301',
        student_id: 'stu_3',
        student_name: 'Vikram Mehta',
        roll_number: '2026ECE014',
        borrowed_at: '2026-08-15T09:30:00Z',
        due_date: '2026-08-29T18:00:00Z',
        returned_at: '2026-08-28T16:00:00Z',
        renew_count: 0,
        status: 'RETURNED',
        fine_amount: 0,
        fine_paid: true,
        created_at: '2026-08-15T09:30:00Z'
      }
    ] as LibraryBorrowRecord[]
  } = $props();

  let activeTab = $state<'ALL' | 'ISSUED' | 'OVERDUE' | 'RETURNED'>('ISSUED');
  let searchQuery = $state<string>('');
  let isIssueModalOpen = $state<boolean>(false);

  // New Issue Form
  let issueStudentName = $state<string>('');
  let issueRollNumber = $state<string>('');
  let issueBookTitle = $state<string>('');
  let issueAccession = $state<string>('');

  let filteredBorrows = $derived(
    borrows.filter((b) => {
      const matchTab = activeTab === 'ALL' || b.status === activeTab;
      const matchSearch =
        searchQuery.trim() === '' ||
        b.student_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        b.roll_number?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        b.book_title?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        b.accession_number?.toLowerCase().includes(searchQuery.toLowerCase());
      return matchTab && matchSearch;
    })
  );

  function handleReturnBook(b: LibraryBorrowRecord) {
    const now = new Date();
    const due = new Date(b.due_date);
    let fine = 0;

    if (now > due) {
      const diffDays = Math.ceil((now.getTime() - due.getTime()) / (1000 * 60 * 60 * 24));
      fine = diffDays * 5; // ₹5 per day
    }

    b.status = 'RETURNED';
    b.returned_at = now.toISOString();
    b.fine_amount = fine;

    if (fine > 0) {
      alert(`⚠️ Book Returned Late (${fine / 5} days overdue). Overdue fine of ₹${fine} charged to student billing invoice.`);
    } else {
      alert(`✓ Book Checked In successfully on time. No fines accrued.`);
    }
  }

  function handleRenewBook(b: LibraryBorrowRecord) {
    if (b.renew_count >= 2) {
      alert('Error: Maximum allowed renewals (2) exceeded for this loan.');
      return;
    }
    if (new Date() > new Date(b.due_date)) {
      alert('Error: Overdue books cannot be renewed. Please return the book and settle fines.');
      return;
    }

    b.renew_count++;
    const currentDue = new Date(b.due_date);
    currentDue.setDate(currentDue.getDate() + 14);
    b.due_date = currentDue.toISOString();

    alert(`✓ Book Loan Renewed (+14 days). New due date: ${currentDue.toLocaleDateString()}`);
  }

  function handleIssueNewBook() {
    if (!issueStudentName || !issueRollNumber || !issueBookTitle || !issueAccession) {
      alert('Please fill all required checkout fields');
      return;
    }

    const now = new Date();
    const due = new Date();
    due.setDate(due.getDate() + 14);

    const newBorrow: LibraryBorrowRecord = {
      id: `brw_${Date.now()}`,
      tenant_id: 'ten_default',
      copy_id: `cpy_${Date.now()}`,
      book_title: issueBookTitle,
      isbn: '978-0131103627',
      accession_number: issueAccession,
      student_id: `stu_${Date.now()}`,
      student_name: issueStudentName,
      roll_number: issueRollNumber,
      borrowed_at: now.toISOString(),
      due_date: due.toISOString(),
      renew_count: 0,
      status: 'ISSUED',
      fine_amount: 0,
      fine_paid: false,
      created_at: now.toISOString()
    };

    borrows.unshift(newBorrow);
    isIssueModalOpen = false;
    // reset
    issueStudentName = '';
    issueRollNumber = '';
    issueBookTitle = '';
    issueAccession = '';
    alert('✓ Book Loan Registered: Due date in 14 days.');
  }
</script>

<div class="space-y-6">
  <!-- Controls Bar -->
  <Card class="p-4">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div class="flex items-center gap-1.5 bg-muted p-1 rounded-lg">
        <Button
          size="sm"
          variant={activeTab === 'ISSUED' ? 'default' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'ISSUED')}
        >
          Active Loans ({borrows.filter((b) => b.status === 'ISSUED').length})
        </Button>
        <Button
          size="sm"
          variant={activeTab === 'OVERDUE' ? 'destructive' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'OVERDUE')}
        >
          Overdue Fines ({borrows.filter((b) => b.status === 'OVERDUE').length})
        </Button>
        <Button
          size="sm"
          variant={activeTab === 'RETURNED' ? 'default' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'RETURNED')}
        >
          Returned History ({borrows.filter((b) => b.status === 'RETURNED').length})
        </Button>
        <Button
          size="sm"
          variant={activeTab === 'ALL' ? 'default' : 'ghost'}
          class="h-8 text-xs font-semibold"
          onclick={() => (activeTab = 'ALL')}
        >
          All Records ({borrows.length})
        </Button>
      </div>

      <div class="flex items-center gap-3">
        <Input
          type="search"
          placeholder="Search by student, book, accession..."
          bind:value={searchQuery}
          class="w-64 h-9"
        />
        <Button size="sm" class="h-9 font-semibold" onclick={() => (isIssueModalOpen = true)}>
          + Issue Book Loan
        </Button>
      </div>
    </div>
  </Card>

  <!-- Loans Circulation List -->
  <div class="space-y-4">
    {#each filteredBorrows as borrow (borrow.id)}
      <Card class="shadow-sm border hover:border-primary/40 transition-colors">
        <CardHeader class="pb-3 border-b bg-muted/10">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="font-bold text-base text-foreground">{borrow.book_title}</span>
              <Badge variant="outline" class="font-mono text-xs">{borrow.accession_number}</Badge>
            </div>

            <div>
              {#if borrow.status === 'ISSUED'}
                <Badge class="bg-blue-100 text-blue-800 border-blue-300">ACTIVE ON LOAN</Badge>
              {:else if borrow.status === 'OVERDUE'}
                <Badge variant="destructive">OVERDUE · FINE CHARGED</Badge>
              {:else if borrow.status === 'RETURNED'}
                <Badge class="bg-emerald-100 text-emerald-800 border-emerald-300">RETURNED ✓</Badge>
              {/if}
            </div>
          </div>
        </CardHeader>

        <CardContent class="p-4 space-y-3">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
            <div>
              <span class="text-muted-foreground">Borrower Student:</span>
              <p class="font-bold text-foreground text-sm">{borrow.student_name}</p>
              <p class="font-mono text-muted-foreground">Roll No: {borrow.roll_number}</p>
            </div>

            <div class="space-y-1">
              <p><strong class="text-muted-foreground">Issue Date:</strong> {new Date(borrow.borrowed_at).toLocaleDateString()}</p>
              <p><strong class="text-muted-foreground">Due Date:</strong> {new Date(borrow.due_date).toLocaleDateString()}</p>
              <p><strong class="text-muted-foreground">Renewals:</strong> {borrow.renew_count}/2 Used</p>
            </div>

            <div class="space-y-1">
              {#if borrow.fine_amount > 0}
                <div class="p-2 rounded bg-destructive/10 border border-destructive/20 text-destructive">
                  <p class="font-bold">Overdue Penalty: ₹{borrow.fine_amount}</p>
                  <p class="text-[10px] text-muted-foreground">Calculated at ₹5/day overdue rate</p>
                </div>
              {:else}
                <p class="text-emerald-600 font-semibold">✓ No Overdue Fines</p>
              {/if}
            </div>
          </div>
        </CardContent>

        <CardFooter class="p-3 bg-muted/20 border-t flex items-center justify-between">
          <span class="text-[11px] text-muted-foreground">
            Loan Window: 14 Days · Auto-fine on overdue
          </span>

          <div class="flex items-center gap-2">
            {#if borrow.status === 'ISSUED'}
              <Button
                size="sm"
                variant="outline"
                class="h-8 text-xs font-semibold"
                onclick={() => handleRenewBook(borrow)}
              >
                Renew (+14 Days)
              </Button>
              <Button
                size="sm"
                class="h-8 text-xs bg-slate-900 text-white font-semibold"
                onclick={() => handleReturnBook(borrow)}
              >
                Check In / Return
              </Button>
            {:else if borrow.status === 'OVERDUE'}
              <Button
                size="sm"
                class="h-8 text-xs bg-destructive text-white font-semibold"
                onclick={() => handleReturnBook(borrow)}
              >
                Return & Settle Fine (₹{borrow.fine_amount})
              </Button>
            {/if}
          </div>
        </CardFooter>
      </Card>
    {/each}
  </div>
</div>

<!-- Issue Book Loan Modal -->
{#if isIssueModalOpen}
  <div class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
    <Card class="w-full max-w-lg shadow-2xl bg-card animate-in fade-in zoom-in-95">
      <CardHeader>
        <CardTitle class="text-lg">Issue Book Checkout Loan</CardTitle>
        <CardDescription>
          Record physical copy checkout to enrolled student (14-day circulation period).
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-3">
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="stu_borrow_name">Student Name</label>
            <Input id="stu_borrow_name" placeholder="Full name" bind:value={issueStudentName} />
          </div>
          <div class="space-y-1">
            <label class="text-xs font-semibold" for="stu_borrow_roll">Roll Number</label>
            <Input id="stu_borrow_roll" placeholder="Roll No" bind:value={issueRollNumber} />
          </div>
        </div>

        <div class="space-y-1">
          <label class="text-xs font-semibold" for="bk_borrow_title">Book Title</label>
          <Input id="bk_borrow_title" placeholder="e.g. The C Programming Language" bind:value={issueBookTitle} />
        </div>

        <div class="space-y-1">
          <label class="text-xs font-semibold" for="bk_borrow_acc">Accession Barcode Number</label>
          <Input id="bk_borrow_acc" placeholder="e.g. LIB-ACC-001042" bind:value={issueAccession} />
        </div>
      </CardContent>
      <CardFooter class="flex justify-end gap-2 border-t pt-4">
        <Button variant="ghost" size="sm" onclick={() => (isIssueModalOpen = false)}>
          Cancel
        </Button>
        <Button size="sm" onclick={handleIssueNewBook}>
          Confirm Book Checkout
        </Button>
      </CardFooter>
    </Card>
  </div>
{/if}
