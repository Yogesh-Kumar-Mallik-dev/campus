<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_LIBRARY_001
   * Subsystem: Rank 9 - E-Library Management System (library)
   * Purpose:   Integrated physical book catalog, digital borrow loan management, overdue fine billing, and e-resource access.
   */
  import BookCatalogDesk from '$lib/components/library/BookCatalogDesk.svelte';
  import BorrowLoanDesk from '$lib/components/library/BorrowLoanDesk.svelte';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        E-Library & Resource Hub
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Catalog search, ISBN management, borrow loans with auto overdue fine calculations, and digital DRM assets.
      </p>
    </div>
  </div>

  <Tabs value="catalog" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-sm">
      <TabsTrigger value="catalog">Book Catalog & Inventory</TabsTrigger>
      <TabsTrigger value="loans">Active Borrow Loans & Fines</TabsTrigger>
    </TabsList>

    <TabsContent value="catalog" class="mt-6">
      <BookCatalogDesk />
    </TabsContent>

    <TabsContent value="loans" class="mt-6">
      <BorrowLoanDesk />
    </TabsContent>
  </Tabs>
</div>
