<script lang="ts">
  /**
   * BLOCK_WEB_LIBRARY_CATALOG_DESK_001
   * Subsystem: Rank 9 - E-Library System (library)
   * Purpose:   Bibliographic catalog search, shelf locator, physical copy availability, and digital e-reader viewer.
   */
  import type { LibraryBook, BookCategory } from '$lib/types/library';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let {
    books = [
      {
        id: 'bk_001',
        tenant_id: 'ten_default',
        isbn: '978-0131103627',
        title: 'The C Programming Language',
        author: 'Brian W. Kernighan, Dennis M. Ritchie',
        publisher: 'Prentice Hall',
        edition: '2nd Edition',
        category: 'COMPUTER_SCIENCE',
        total_copies: 8,
        available_copies: 5,
        shelf_location: 'Stack CS-01, Bay A',
        ebook_key: 's3://library/c_programming_ansi.pdf',
        ebook_format: 'PDF',
        ebook_size: 4200000,
        created_at: '2026-08-01T00:00:00Z'
      },
      {
        id: 'bk_002',
        tenant_id: 'ten_default',
        isbn: '978-0262033848',
        title: 'Introduction to Algorithms (CLRS)',
        author: 'Thomas H. Cormen, Charles E. Leiserson',
        publisher: 'MIT Press',
        edition: '3rd Edition',
        category: 'COMPUTER_SCIENCE',
        total_copies: 12,
        available_copies: 3,
        shelf_location: 'Stack CS-02, Bay C',
        ebook_key: 's3://library/clrs_algorithms_3e.pdf',
        ebook_format: 'PDF',
        ebook_size: 18500000,
        created_at: '2026-08-01T00:00:00Z'
      },
      {
        id: 'bk_003',
        tenant_id: 'ten_default',
        isbn: '978-0132350884',
        title: 'Clean Code: A Handbook of Agile Craftsmanship',
        author: 'Robert C. Martin',
        publisher: 'Prentice Hall',
        edition: '1st Edition',
        category: 'COMPUTER_SCIENCE',
        total_copies: 6,
        available_copies: 0,
        shelf_location: 'Stack CS-04, Bay B',
        ebook_key: 's3://library/clean_code.epub',
        ebook_format: 'EPUB',
        ebook_size: 2100000,
        created_at: '2026-08-05T00:00:00Z'
      },
      {
        id: 'bk_004',
        tenant_id: 'ten_default',
        isbn: '978-0198553182',
        title: 'Quantum Physics: A Modern Introduction',
        author: 'Alastair I. M. Rae',
        publisher: 'Oxford University Press',
        edition: '4th Edition',
        category: 'PHYSICS',
        total_copies: 4,
        available_copies: 2,
        shelf_location: 'Stack PHY-01, Bay D',
        created_at: '2026-08-10T00:00:00Z'
      }
    ] as LibraryBook[]
  } = $props();

  let selectedCategory = $state<string>('ALL');
  let searchQuery = $state<string>('');
  let activeEbook = $state<LibraryBook | null>(null);
  let activePage = $state<number>(1);

  let filteredBooks = $derived(
    books.filter((b) => {
      const matchCat = selectedCategory === 'ALL' || b.category === selectedCategory;
      const matchSearch =
        searchQuery.trim() === '' ||
        b.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        b.author.toLowerCase().includes(searchQuery.toLowerCase()) ||
        b.isbn.includes(searchQuery.trim());
      return matchCat && matchSearch;
    })
  );

  const categories = [
    { label: 'All Subjects', value: 'ALL' },
    { label: 'Computer Science', value: 'COMPUTER_SCIENCE' },
    { label: 'Physics', value: 'PHYSICS' },
    { label: 'Mathematics', value: 'MATHEMATICS' },
    { label: 'Electronics', value: 'ELECTRONICS' },
    { label: 'Management', value: 'MANAGEMENT' }
  ];

  function handleOpenEbook(b: LibraryBook) {
    activeEbook = b;
    activePage = 1;
  }

  function handleReserve(b: LibraryBook) {
    alert(`✓ Hold Placed for "${b.title}". You will receive a notification when a copy is checked in.`);
  }
</script>

<div class="space-y-6">
  <!-- Filters and Search Bar -->
  <Card class="p-4">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div class="flex flex-wrap items-center gap-1 bg-muted p-1 rounded-lg">
        {#each categories as cat}
          <Button
            size="sm"
            variant={selectedCategory === cat.value ? 'default' : 'ghost'}
            class="h-8 text-xs font-semibold px-3"
            onclick={() => (selectedCategory = cat.value)}
          >
            {cat.label}
          </Button>
        {/each}
      </div>

      <div class="flex items-center gap-3">
        <Input
          type="search"
          placeholder="Search by Title, Author, ISBN..."
          bind:value={searchQuery}
          class="w-72 h-9"
        />
      </div>
    </div>
  </Card>

  <!-- Books Catalog Grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
    {#each filteredBooks as book (book.id)}
      <Card class="shadow-sm border hover:border-primary/40 transition-colors flex flex-col justify-between">
        <CardHeader class="pb-3 border-b bg-muted/10">
          <div class="flex items-start justify-between gap-2">
            <div>
              <CardTitle class="text-base font-bold text-foreground leading-snug line-clamp-2">
                {book.title}
              </CardTitle>
              <CardDescription class="text-xs text-muted-foreground mt-0.5">
                By <span class="font-semibold text-foreground">{book.author}</span>
              </CardDescription>
            </div>
            {#if book.available_copies > 0}
              <Badge class="bg-emerald-100 text-emerald-800 border-emerald-300 text-[10px] shrink-0">
                {book.available_copies}/{book.total_copies} AVAIL
              </Badge>
            {:else}
              <Badge variant="destructive" class="text-[10px] shrink-0">ALL ISSUED</Badge>
            {/if}
          </div>
        </CardHeader>

        <CardContent class="p-4 space-y-3 text-xs flex-1">
          <div class="grid grid-cols-2 gap-2 text-muted-foreground">
            <div>
              <span>Publisher:</span>
              <p class="font-semibold text-foreground">{book.publisher}</p>
            </div>
            <div>
              <span>Edition:</span>
              <p class="font-semibold text-foreground">{book.edition || 'Standard'}</p>
            </div>
          </div>

          <div class="p-2.5 rounded-lg bg-muted/40 border space-y-1">
            <div class="flex items-center justify-between">
              <span class="text-muted-foreground">Physical Shelf:</span>
              <span class="font-mono font-bold text-primary">{book.shelf_location}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-muted-foreground">ISBN:</span>
              <span class="font-mono text-muted-foreground">{book.isbn}</span>
            </div>
          </div>

          {#if book.ebook_key}
            <div class="p-2 rounded-lg bg-blue-50/50 border border-blue-200 text-blue-900 flex items-center justify-between">
              <div class="flex items-center gap-1.5">
                <span>📖 Digital E-Book</span>
                <Badge variant="outline" class="text-[9px] bg-white text-blue-800 font-mono">
                  {book.ebook_format}
                </Badge>
              </div>
              <span class="text-[10px] text-blue-700 font-mono">
                {((book.ebook_size || 0) / (1024 * 1024)).toFixed(1)} MB
              </span>
            </div>
          {/if}
        </CardContent>

        <CardFooter class="p-3 bg-muted/20 border-t flex items-center justify-between gap-2">
          {#if book.ebook_key}
            <Button
              size="sm"
              variant="outline"
              class="h-8 text-xs font-semibold flex-1 text-blue-700 hover:bg-blue-50"
              onclick={() => handleOpenEbook(book)}
            >
              Read Digital E-Book
            </Button>
          {/if}

          {#if book.available_copies === 0}
            <Button
              size="sm"
              class="h-8 text-xs font-semibold flex-1 bg-amber-600 hover:bg-amber-700 text-white"
              onclick={() => handleReserve(book)}
            >
              Place Reservation Hold
            </Button>
          {:else}
            <Button
              size="sm"
              class="h-8 text-xs font-semibold flex-1 bg-slate-900 text-white"
              onclick={() => alert(`Locate physical copy at: ${book.shelf_location}`)}
            >
              Locate in Stacks
            </Button>
          {/if}
        </CardFooter>
      </Card>
    {/each}
  </div>
</div>

<!-- Digital E-Reader Modal -->
{#if activeEbook}
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
    <Card class="w-full max-w-2xl h-[550px] shadow-2xl bg-card flex flex-col justify-between animate-in fade-in zoom-in-95">
      <CardHeader class="border-b pb-3">
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="text-base font-bold">{activeEbook.title}</CardTitle>
            <CardDescription class="text-xs">{activeEbook.author} · {activeEbook.edition}</CardDescription>
          </div>
          <Button variant="ghost" size="sm" onclick={() => (activeEbook = null)}>✕</Button>
        </div>
      </CardHeader>

      <CardContent class="flex-1 p-6 flex flex-col items-center justify-center space-y-4 bg-muted/10 text-center">
        <div class="w-full max-w-lg p-8 rounded-2xl bg-background border shadow-inner space-y-4">
          <Badge variant="outline" class="font-mono text-xs">PAGE {activePage} OF 480</Badge>
          <h3 class="text-lg font-bold text-foreground">Chapter {activePage}: Core Architecture & Fundamentals</h3>
          <p class="text-xs text-muted-foreground leading-relaxed">
            "In computing, modular design and strict boundaries ensure decoupling and long-term maintainability across enterprise architectures..."
          </p>
          <div class="pt-4 flex items-center justify-center gap-3">
            <Button
              size="sm"
              variant="outline"
              class="h-8 px-4 text-xs"
              disabled={activePage <= 1}
              onclick={() => activePage--}
            >
              ← Previous Page
            </Button>
            <Button
              size="sm"
              variant="outline"
              class="h-8 px-4 text-xs"
              onclick={() => activePage++}
            >
              Next Page →
            </Button>
          </div>
        </div>
      </CardContent>

      <CardFooter class="border-t p-3 bg-muted/20 flex items-center justify-between text-xs text-muted-foreground">
        <span>Session Reading Time: 12 minutes</span>
        <Button size="sm" class="h-7 text-xs" onclick={() => (activeEbook = null)}>
          Save Progress & Close
        </Button>
      </CardFooter>
    </Card>
  </div>
{/if}
