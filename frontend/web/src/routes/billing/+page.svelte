<script lang="ts">
  /**
   * BLOCK_WEB_PAGE_BILLING_001
   * Purpose: Institutional Payment & Billing Command Center for Students, Bursar, and Finance Admins.
   */
  import InvoiceDesk from '$lib/components/billing/InvoiceDesk.svelte';
  import PaymentLedger from '$lib/components/billing/PaymentLedger.svelte';
  import type { PaymentMethod } from '$lib/types/billing';
  import { Button } from '$lib/components/ui/button';
  import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';

  async function handlePayInvoice(invoiceID: string, amountPaise: number, method: PaymentMethod) {
    // In production, invokes POST /api/v1/billing/payments/initiate then gateway callback
    console.log('Initiating checkout for invoice:', invoiceID, amountPaise, method);
  }
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight text-foreground">
        Finance & Billing Command Center
      </h1>
      <p class="text-sm text-muted-foreground mt-1">
        Student fee ledger, itemized invoices, split installments, gateway payments, and double-entry accounting.
      </p>
    </div>

    <div class="flex gap-2">
      <Button variant="outline" size="sm">Generate Invoice</Button>
      <Button size="sm">New Fee Schedule</Button>
    </div>
  </div>

  <Tabs value="invoices" class="w-full">
    <TabsList class="grid w-full grid-cols-2 max-w-md">
      <TabsTrigger value="invoices">Invoices & Checkout</TabsTrigger>
      <TabsTrigger value="ledger">Accounting & Receipts</TabsTrigger>
    </TabsList>

    <TabsContent value="invoices" class="mt-6">
      <InvoiceDesk onPayInvoice={handlePayInvoice} />
    </TabsContent>

    <TabsContent value="ledger" class="mt-6">
      <PaymentLedger />
    </TabsContent>
  </Tabs>
</div>
