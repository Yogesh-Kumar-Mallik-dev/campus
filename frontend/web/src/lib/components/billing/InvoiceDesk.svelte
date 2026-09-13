<script lang="ts">
  /**
   * BLOCK_WEB_BILLING_INVOICE_DESK_001
   * Subsystem: Rank 5 - Central Payment & Billing System (billing)
   * Purpose:   Student fee invoicing desk with itemized breakdown, scholarship deductions, and split installment checkout.
   */
  import type { StudentInvoice, PaymentMethod } from '$lib/types/billing';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';

  let {
    invoices = [
      {
        id: 'inv_101',
        tenant_id: 'ten_default',
        student_id: 'stu_101',
        student_name: 'Aarav Sharma',
        invoice_number: 'INV-2026-2027-000101',
        academic_year: '2026-2027',
        semester: 1,
        subtotal_amount: 5500000,
        discount_amount: 500000,
        tax_amount: 0,
        total_amount: 5000000,
        paid_amount: 2500000,
        balance_amount: 2500000,
        due_date: '2026-10-15',
        status: 'PARTIALLY_PAID',
        items: [
          { id: 'it_1', tenant_id: 'ten_default', invoice_id: 'inv_101', category: 'TUITION', name: 'Semester Tuition', amount: 4500000, created_at: '', updated_at: '' },
          { id: 'it_2', tenant_id: 'ten_default', invoice_id: 'inv_101', category: 'EXAMINATION', name: 'University Examination', amount: 500000, created_at: '', updated_at: '' },
          { id: 'it_3', tenant_id: 'ten_default', invoice_id: 'inv_101', category: 'LABORATORY', name: 'Advanced Computing Lab', amount: 500000, created_at: '', updated_at: '' }
        ],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      },
      {
        id: 'inv_102',
        tenant_id: 'ten_default',
        student_id: 'stu_101',
        student_name: 'Aarav Sharma',
        invoice_number: 'INV-2026-2027-000102',
        academic_year: '2026-2027',
        semester: 1,
        subtotal_amount: 4500000,
        discount_amount: 0,
        tax_amount: 0,
        total_amount: 4500000,
        paid_amount: 0,
        balance_amount: 4500000,
        due_date: '2026-10-31',
        status: 'ISSUED',
        items: [
          { id: 'it_4', tenant_id: 'ten_default', invoice_id: 'inv_102', category: 'HOSTEL', name: 'Hostel Room (AC Double)', amount: 3000000, created_at: '', updated_at: '' },
          { id: 'it_5', tenant_id: 'ten_default', invoice_id: 'inv_102', category: 'MESS', name: 'Mess Dining Plan (Semester)', amount: 1500000, created_at: '', updated_at: '' }
        ],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
    ],
    onPayInvoice = async (_invoiceID: string, _amount: number, _method: PaymentMethod) => {}
  }: {
    invoices?: StudentInvoice[];
    onPayInvoice?: (invoiceID: string, amount: number, method: PaymentMethod) => Promise<void>;
  } = $props();

  let selectedInvoice = $state<StudentInvoice | null>(invoices[0] || null);
  let customPayAmount = $state<number>(invoices[0]?.balance_amount ? invoices[0].balance_amount / 100 : 0);
  let selectedMethod = $state<PaymentMethod>('ONLINE_GATEWAY');
  let isProcessing = $state(false);
  let paymentSuccessReceipt = $state<string | null>(null);

  function selectInvoice(inv: StudentInvoice) {
    selectedInvoice = inv;
    customPayAmount = inv.balance_amount / 100;
    paymentSuccessReceipt = null;
  }

  function formatCurrency(amountCents: number): string {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
      maximumFractionDigits: 0
    }).format(amountCents / 100);
  }

  async function handleCheckout() {
    if (!selectedInvoice || customPayAmount <= 0) return;
    isProcessing = true;
    try {
      const amountPaise = Math.round(customPayAmount * 100);
      await onPayInvoice(selectedInvoice.id, amountPaise, selectedMethod);
      paymentSuccessReceipt = `REC-${new Date().getFullYear()}-${Math.floor(100000 + Math.random() * 900000)}`;
      selectedInvoice.paid_amount += amountPaise;
      selectedInvoice.balance_amount = selectedInvoice.total_amount - selectedInvoice.paid_amount;
      if (selectedInvoice.balance_amount === 0) {
        selectedInvoice.status = 'PAID';
      } else {
        selectedInvoice.status = 'PARTIALLY_PAID';
      }
    } finally {
      isProcessing = false;
    }
  }
</script>

<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
  <!-- Invoice List -->
  <div class="lg:col-span-1 space-y-3">
    <h3 class="text-sm font-semibold uppercase tracking-wider text-muted-foreground px-1">
      Institutional Invoices
    </h3>

    {#each invoices as inv (inv.id)}
      <div
        role="button"
        tabindex="0"
        class="p-4 rounded-xl border transition-all text-left cursor-pointer {selectedInvoice?.id === inv.id ? 'border-primary bg-primary/5 shadow-sm' : 'border-border bg-card hover:bg-muted/50'}"
        onclick={() => selectInvoice(inv)}
        onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') selectInvoice(inv); }}
      >
        <div class="flex items-center justify-between mb-2">
          <span class="font-mono text-xs font-bold text-foreground">{inv.invoice_number}</span>
          <Badge variant={inv.status === 'PAID' ? 'secondary' : inv.status === 'PARTIALLY_PAID' ? 'default' : 'outline'}>
            {inv.status}
          </Badge>
        </div>

        <div class="flex justify-between items-baseline">
          <div>
            <p class="text-xs text-muted-foreground">Due: {inv.due_date}</p>
            <p class="text-xs text-muted-foreground">Semester {inv.semester} · {inv.academic_year}</p>
          </div>
          <div class="text-right">
            <p class="text-xs text-muted-foreground">Balance</p>
            <p class="font-bold text-sm text-foreground">{formatCurrency(inv.balance_amount)}</p>
          </div>
        </div>
      </div>
    {/each}
  </div>

  <!-- Invoice Detail & Checkout Panel -->
  <div class="lg:col-span-2">
    {#if selectedInvoice}
      <Card class="border shadow-sm">
        <CardHeader class="pb-4 border-b bg-muted/20">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <div>
              <div class="flex items-center gap-2">
                <CardTitle class="text-xl font-bold">{selectedInvoice.invoice_number}</CardTitle>
                <Badge variant={selectedInvoice.status === 'PAID' ? 'secondary' : 'default'}>
                  {selectedInvoice.status}
                </Badge>
              </div>
              <CardDescription class="mt-0.5">
                Billed to: <span class="font-medium text-foreground">{selectedInvoice.student_name || selectedInvoice.student_id}</span> |
                Due Date: <span class="font-medium text-foreground">{selectedInvoice.due_date}</span>
              </CardDescription>
            </div>
            <div class="text-right">
              <span class="text-xs text-muted-foreground">Outstanding Balance</span>
              <p class="text-2xl font-extrabold text-foreground">{formatCurrency(selectedInvoice.balance_amount)}</p>
            </div>
          </div>
        </CardHeader>

        <CardContent class="p-6 space-y-6">
          <!-- Fee Heads Breakdown -->
          <div>
            <h4 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3">
              Itemized Fee Heads
            </h4>
            <div class="rounded-lg border divide-y overflow-hidden text-sm">
              {#each selectedInvoice.items || [] as item}
                <div class="flex justify-between items-center p-3 bg-card hover:bg-muted/30">
                  <div>
                    <span class="font-medium text-foreground">{item.name}</span>
                    <span class="text-xs text-muted-foreground ml-2">({item.category})</span>
                  </div>
                  <span class="font-mono font-medium">{formatCurrency(item.amount)}</span>
                </div>
              {/each}

              {#if selectedInvoice.discount_amount > 0}
                <div class="flex justify-between items-center p-3 bg-emerald-50/50 dark:bg-emerald-950/20 text-emerald-700 dark:text-emerald-300">
                  <span>Scholarship / Institutional Waiver</span>
                  <span class="font-mono font-medium">- {formatCurrency(selectedInvoice.discount_amount)}</span>
                </div>
              {/if}

              <div class="flex justify-between items-center p-3 bg-muted/40 font-bold">
                <span>Total Invoice Value</span>
                <span class="font-mono">{formatCurrency(selectedInvoice.total_amount)}</span>
              </div>
            </div>
          </div>

          <!-- Checkout Box (if not paid) -->
          {#if selectedInvoice.status !== 'PAID'}
            <div class="p-5 rounded-xl border bg-card space-y-4">
              <h4 class="text-sm font-bold text-foreground">Initiate Payment Checkout</h4>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-1.5">
                  <Label for="payAmount">Amount to Pay (INR)</Label>
                  <Input
                    id="payAmount"
                    type="number"
                    bind:value={customPayAmount}
                    min="1"
                    max={selectedInvoice.balance_amount / 100}
                  />
                  <span class="text-[11px] text-muted-foreground">
                    Max payable balance: {formatCurrency(selectedInvoice.balance_amount)}
                  </span>
                </div>

                <div class="space-y-1.5">
                  <Label for="payMethod">Payment Channel</Label>
                  <select
                    id="payMethod"
                    bind:value={selectedMethod}
                    class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  >
                    <option value="ONLINE_GATEWAY">Instant Payment Gateway (Razorpay/Stripe)</option>
                    <option value="UPI">UPI Fast QR / VPA</option>
                    <option value="BANK_TRANSFER">NEFT / RTGS Bank Transfer</option>
                    <option value="CHEQUE">Banker's Demand Draft / Cheque</option>
                  </select>
                </div>
              </div>

              {#if paymentSuccessReceipt}
                <div class="p-4 rounded-lg bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 space-y-1">
                  <p class="text-sm font-bold text-emerald-900 dark:text-emerald-200">✓ Payment Settled Successfully!</p>
                  <p class="text-xs text-emerald-700 dark:text-emerald-400">
                    Official Institutional Receipt Issued: <span class="font-mono font-bold">{paymentSuccessReceipt}</span>
                  </p>
                </div>
              {/if}

              <Button
                class="w-full"
                disabled={isProcessing || customPayAmount <= 0 || customPayAmount > (selectedInvoice.balance_amount / 100)}
                onclick={handleCheckout}
              >
                {isProcessing ? 'Processing Gateway Settlement...' : `Pay ${new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(customPayAmount || 0)}`}
              </Button>
            </div>
          {:else}
            <div class="p-4 rounded-xl bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-200 text-center space-y-1">
              <p class="text-sm font-bold text-emerald-800 dark:text-emerald-300">✓ Invoice Fully Paid & Settled</p>
              <p class="text-xs text-emerald-600 dark:text-emerald-400">No outstanding balance remains on this schedule.</p>
            </div>
          {/if}
        </CardContent>
      </Card>
    {/if}
  </div>
</div>
