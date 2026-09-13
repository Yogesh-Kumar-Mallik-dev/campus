<script lang="ts">
  /**
   * BLOCK_WEB_BILLING_LEDGER_001
   * Subsystem: Rank 5 - Central Payment & Billing System (billing)
   * Purpose:   Double-entry accounting ledger view, receipt history, and payment transactions audit monitor.
   */
  import type { PaymentTransaction, LedgerEntry } from '$lib/types/billing';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';

  let {
    transactions = [
      {
        id: 'txn_1',
        tenant_id: 'ten_default',
        invoice_id: 'inv_101',
        student_id: 'stu_101',
        transaction_ref: 'TXN-2026-000001',
        gateway_name: 'RAZORPAY',
        amount: 2500000,
        currency: 'INR',
        method: 'ONLINE_GATEWAY',
        status: 'SUCCESS',
        receipt_number: 'REC-2026-000001',
        paid_at: '2026-09-10T10:30:00Z',
        created_at: '2026-09-10T10:30:00Z',
        updated_at: '2026-09-10T10:30:00Z'
      }
    ],
    ledgerEntries = [
      {
        id: 'lent_1',
        tenant_id: 'ten_default',
        transaction_id: 'txn_1',
        account_id: 'acc_bank',
        account_name: '1100-BANK (HDFC Institutional)',
        entry_type: 'DEBIT',
        amount: 2500000,
        description: 'Payment for INV-2026-2027-000101',
        posted_at: '2026-09-10T10:30:00Z',
        created_at: '2026-09-10T10:30:00Z'
      },
      {
        id: 'lent_2',
        tenant_id: 'ten_default',
        transaction_id: 'txn_1',
        account_id: 'acc_ar',
        account_name: '2000-AR (Accounts Receivable)',
        entry_type: 'CREDIT',
        amount: 2500000,
        description: 'AR Clearance for Aarav Sharma',
        posted_at: '2026-09-10T10:30:00Z',
        created_at: '2026-09-10T10:30:00Z'
      }
    ]
  }: {
    transactions?: PaymentTransaction[];
    ledgerEntries?: LedgerEntry[];
  } = $props();

  function formatCurrency(amountCents: number): string {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
      maximumFractionDigits: 0
    }).format(amountCents / 100);
  }

  let totalDebit = $derived(ledgerEntries.filter(e => e.entry_type === 'DEBIT').reduce((acc, e) => acc + e.amount, 0));
  let totalCredit = $derived(ledgerEntries.filter(e => e.entry_type === 'CREDIT').reduce((acc, e) => acc + e.amount, 0));
  let isBalanced = $derived(totalDebit === totalCredit);
</script>

<div class="space-y-6">
  <!-- Double Entry Integrity Pill -->
  <div class="p-4 rounded-xl border bg-card flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-sm">
    <div class="flex items-center gap-3">
      <div class="h-10 w-10 rounded-full flex items-center justify-center {isBalanced ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300' : 'bg-rose-100 text-rose-700'}">
        {isBalanced ? '✓' : '✗'}
      </div>
      <div>
        <h4 class="text-sm font-bold text-foreground">Double-Entry Accounting Ledger Invariant</h4>
        <p class="text-xs text-muted-foreground">
          Total Debits: <span class="font-mono font-semibold">{formatCurrency(totalDebit)}</span> |
          Total Credits: <span class="font-mono font-semibold">{formatCurrency(totalCredit)}</span>
        </p>
      </div>
    </div>

    <Badge variant={isBalanced ? 'secondary' : 'destructive'} class="self-start sm:self-auto">
      {isBalanced ? 'LEDGER BALANCED (DEBITS = CREDITS)' : 'UNBALANCED WARNING'}
    </Badge>
  </div>

  <!-- Transactions Table -->
  <Card class="border shadow-sm">
    <CardHeader class="pb-4 border-b bg-muted/30">
      <CardTitle class="text-xl font-bold">Settlement & Receipt Ledger</CardTitle>
      <CardDescription>
        Immutable audit trail of gateway collections, offline clearances, and sequential receipt numbers.
      </CardDescription>
    </CardHeader>

    <CardContent class="p-0">
      <div class="overflow-x-auto">
        <table class="w-full text-sm text-left">
          <thead class="bg-muted/40 text-muted-foreground uppercase text-xs">
            <tr>
              <th class="px-6 py-3">Txn Ref</th>
              <th class="px-6 py-3">Receipt No</th>
              <th class="px-6 py-3">Method / Gateway</th>
              <th class="px-6 py-3 text-right">Amount</th>
              <th class="px-6 py-3 text-center">Status</th>
              <th class="px-6 py-3">Date</th>
            </tr>
          </thead>
          <tbody class="divide-y">
            {#each transactions as tx (tx.id)}
              <tr class="hover:bg-muted/20 transition-colors">
                <td class="px-6 py-4 font-mono font-medium text-xs">
                  {tx.transaction_ref}
                </td>
                <td class="px-6 py-4 font-mono font-bold text-xs text-emerald-600 dark:text-emerald-400">
                  {tx.receipt_number || '—'}
                </td>
                <td class="px-6 py-4 text-xs">
                  <span class="font-medium text-foreground">{tx.method}</span>
                  <span class="text-muted-foreground ml-1">({tx.gateway_name})</span>
                </td>
                <td class="px-6 py-4 text-right font-mono font-bold">
                  {formatCurrency(tx.amount)}
                </td>
                <td class="px-6 py-4 text-center">
                  <Badge variant={tx.status === 'SUCCESS' ? 'secondary' : tx.status === 'REFUNDED' ? 'outline' : 'default'}>
                    {tx.status}
                  </Badge>
                </td>
                <td class="px-6 py-4 text-xs text-muted-foreground">
                  {tx.paid_at ? new Date(tx.paid_at).toLocaleString() : 'Pending'}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </CardContent>
  </Card>
</div>
