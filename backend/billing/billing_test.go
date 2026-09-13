/**
 * BLOCK_BILLING_TEST_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   Exhaustive unit test suite verifying fee templates, invoices, idempotency, double-entry ledger balance, and receipts.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package billing_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"campus/backend/audit"
	"campus/backend/billing"
)

type mockAuditSubscriber struct {
	mu     sync.Mutex
	events []audit.RecordAuditRequest
}

func (m *mockAuditSubscriber) Enqueue(req audit.RecordAuditRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, req)
	return nil
}

func (m *mockAuditSubscriber) Start(ctx context.Context) {}
func (m *mockAuditSubscriber) Stop()                      {}

func setupBillingTestService() (*billing.Service, *billing.MockFeeRepository, *billing.MockInvoiceRepository, *billing.MockPaymentRepository, *billing.MockLedgerRepository, *mockAuditSubscriber) {
	feeRepo := billing.NewMockFeeRepository()
	invRepo := billing.NewMockInvoiceRepository()
	payRepo := billing.NewMockPaymentRepository()
	ledgerRepo := billing.NewMockLedgerRepository()
	seqRepo := billing.NewMockSequenceRepository()
	auditSub := &mockAuditSubscriber{}

	svc := billing.NewService(feeRepo, invRepo, payRepo, ledgerRepo, seqRepo, auditSub)
	return svc, feeRepo, invRepo, payRepo, ledgerRepo, auditSub
}

func TestCreateFeeStructure(t *testing.T) {
	svc, _, _, _, _, _ := setupBillingTestService()
	ctx := context.Background()

	// 1. Success
	fs, err := svc.CreateFeeStructure(ctx, billing.CreateFeeStructureCommand{
		TenantID:     "ten_default",
		ProgramID:    "prog_cs",
		AcademicYear: "2026-2027",
		Semester:     1,
		Name:         "B.Tech CSE - Semester 1 Fee Template",
		DueDate:      "2026-10-15",
		Items: []billing.FeeItemInput{
			{Category: billing.CategoryTuition, Name: "Tuition Fee", Amount: 5000000},
			{Category: billing.CategoryLaboratory, Name: "Computer Lab Fee", Amount: 500000},
			{Category: billing.CategoryLibrary, Name: "Library Fee", Amount: 200000},
		},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fs.TotalAmount != 5700000 {
		t.Errorf("expected total amount 5700000, got %d", fs.TotalAmount)
	}
	if len(fs.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(fs.Items))
	}

	// 2. Missing Tenant
	_, err = svc.CreateFeeStructure(ctx, billing.CreateFeeStructureCommand{
		TenantID: "",
		Name:     "Invalid Fee",
	})
	if err != billing.ErrTenantRequired {
		t.Errorf("expected ErrTenantRequired, got %v", err)
	}

	// 3. Negative Amount
	_, err = svc.CreateFeeStructure(ctx, billing.CreateFeeStructureCommand{
		TenantID:  "ten_default",
		ProgramID: "prog_cs",
		Name:      "Negative Fee",
		Semester:  1,
		Items: []billing.FeeItemInput{
			{Category: billing.CategoryTuition, Name: "Tuition", Amount: -100},
		},
	})
	if err != billing.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestGenerateInvoice(t *testing.T) {
	svc, _, _, _, _, auditSub := setupBillingTestService()
	ctx := context.Background()

	// Create a fee structure first
	fs, err := svc.CreateFeeStructure(ctx, billing.CreateFeeStructureCommand{
		TenantID:     "ten_default",
		ProgramID:    "prog_cs",
		AcademicYear: "2026-2027",
		Semester:     1,
		Name:         "B.Tech CSE - Semester 1",
		DueDate:      "2026-10-15",
		Items: []billing.FeeItemInput{
			{Category: billing.CategoryTuition, Name: "Tuition Fee", Amount: 5000000},
			{Category: billing.CategoryExamination, Name: "Exam Fee", Amount: 500000},
		},
	})
	if err != nil {
		t.Fatalf("CreateFeeStructure failed: %v", err)
	}

	// Generate student invoice with scholarship discount of 1,000,000
	inv, err := svc.GenerateInvoice(ctx, billing.GenerateInvoiceCommand{
		TenantID:       "ten_default",
		StudentID:      "stu_101",
		FeeStructureID: &fs.ID,
		AcademicYear:   "2026-2027",
		Semester:       1,
		DueDate:        "2026-10-15",
		DiscountAmount: 1000000,
		TaxAmount:      0,
	})
	if err != nil {
		t.Fatalf("GenerateInvoice failed: %v", err)
	}

	if inv.SubtotalAmount != 5500000 {
		t.Errorf("expected subtotal 5500000, got %d", inv.SubtotalAmount)
	}
	if inv.TotalAmount != 4500000 {
		t.Errorf("expected total amount 4500000 after discount, got %d", inv.TotalAmount)
	}
	if inv.BalanceAmount != 4500000 {
		t.Errorf("expected initial balance 4500000, got %d", inv.BalanceAmount)
	}
	if inv.Status != billing.InvoiceIssued {
		t.Errorf("expected status ISSUED, got %s", inv.Status)
	}
	if inv.InvoiceNumber != "INV-2026-2027-000001" {
		t.Errorf("expected invoice number INV-2026-2027-000001, got %s", inv.InvoiceNumber)
	}

	if len(auditSub.events) == 0 {
		t.Errorf("expected audit event to be enqueued")
	}
}

func TestPaymentLifecycle_IdempotencyAndDoubleEntry(t *testing.T) {
	svc, _, _, _, ledgerRepo, _ := setupBillingTestService()
	ctx := context.Background()

	// Generate an invoice for 100,000
	inv, err := svc.GenerateInvoice(ctx, billing.GenerateInvoiceCommand{
		TenantID:     "ten_default",
		StudentID:    "stu_202",
		AcademicYear: "2026-2027",
		Semester:     1,
		DueDate:      "2026-10-15",
		CustomItems: []billing.FeeItemInput{
			{Category: billing.CategoryHostel, Name: "Hostel Room Fee", Amount: 100000},
		},
	})
	if err != nil {
		t.Fatalf("GenerateInvoice failed: %v", err)
	}

	// 1. Initiate 1st installment payment of 60,000 with Idempotency Key
	idempKey := "idemp_checkout_xyz_001"
	tx1, err := svc.InitiatePayment(ctx, billing.InitiatePaymentCommand{
		TenantID:       "ten_default",
		InvoiceID:      inv.ID,
		StudentID:      "stu_202",
		Amount:         60000,
		Method:         billing.MethodOnlineGateway,
		GatewayName:    "RAZORPAY",
		IdempotencyKey: &idempKey,
	})
	if err != nil {
		t.Fatalf("InitiatePayment failed: %v", err)
	}
	if tx1.Status != billing.PaymentInitiated {
		t.Errorf("expected status INITIATED, got %s", tx1.Status)
	}

	// 2. Replay same idempotency key returns exact same transaction
	txReplay, err := svc.InitiatePayment(ctx, billing.InitiatePaymentCommand{
		TenantID:       "ten_default",
		InvoiceID:      inv.ID,
		StudentID:      "stu_202",
		Amount:         60000,
		Method:         billing.MethodOnlineGateway,
		IdempotencyKey: &idempKey,
	})
	if err != nil {
		t.Fatalf("Replay InitiatePayment failed: %v", err)
	}
	if txReplay.ID != tx1.ID {
		t.Errorf("expected replayed transaction ID to match %s, got %s", tx1.ID, txReplay.ID)
	}

	// 3. Process payment success for 1st installment
	now := time.Now().UTC()
	paidTx1, err := svc.ProcessPaymentSuccess(ctx, billing.ProcessPaymentSuccessCommand{
		TenantID:      "ten_default",
		TransactionID: tx1.ID,
		PaidAt:        now,
	})
	if err != nil {
		t.Fatalf("ProcessPaymentSuccess failed: %v", err)
	}
	if paidTx1.Status != billing.PaymentSuccess {
		t.Errorf("expected status SUCCESS, got %s", paidTx1.Status)
	}
	if paidTx1.ReceiptNumber == nil || *paidTx1.ReceiptNumber == "" {
		t.Fatalf("expected receipt number to be generated")
	}

	// Verify invoice updated to PARTIALLY_PAID with balance 40,000
	updatedInv, _ := svc.GetInvoice(ctx, "ten_default", inv.ID)
	if updatedInv.Status != billing.InvoicePartiallyPaid {
		t.Errorf("expected status PARTIALLY_PAID, got %s", updatedInv.Status)
	}
	if updatedInv.BalanceAmount != 40000 {
		t.Errorf("expected balance 40000, got %d", updatedInv.BalanceAmount)
	}

	// Verify double-entry ledger entries for transaction 1
	entries1, _ := ledgerRepo.ListEntriesByTransaction(ctx, "ten_default", tx1.ID)
	if len(entries1) != 2 {
		t.Fatalf("expected 2 balanced ledger entries, got %d", len(entries1))
	}
	if err := billing.ValidateBalancedLedger(entries1); err != nil {
		t.Errorf("ledger entries unbalanced: %v", err)
	}

	// 4. Try to pay 50,000 when balance is only 40,000 -> Should fail with ErrOverpaymentNotAllowed
	_, err = svc.InitiatePayment(ctx, billing.InitiatePaymentCommand{
		TenantID:  "ten_default",
		InvoiceID: inv.ID,
		StudentID: "stu_202",
		Amount:    50000,
	})
	if err != billing.ErrOverpaymentNotAllowed {
		t.Errorf("expected ErrOverpaymentNotAllowed, got %v", err)
	}

	// 5. Pay remaining 40,000 to fully clear invoice
	tx2, err := svc.InitiatePayment(ctx, billing.InitiatePaymentCommand{
		TenantID:  "ten_default",
		InvoiceID: inv.ID,
		StudentID: "stu_202",
		Amount:    40000,
		Method:    billing.MethodUPI,
	})
	if err != nil {
		t.Fatalf("Initiate 2nd payment failed: %v", err)
	}

	paidTx2, err := svc.ProcessPaymentSuccess(ctx, billing.ProcessPaymentSuccessCommand{
		TenantID:      "ten_default",
		TransactionID: tx2.ID,
		PaidAt:        time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("Process 2nd payment success failed: %v", err)
	}
	if paidTx2.Status != billing.PaymentSuccess {
		t.Errorf("expected 2nd payment status SUCCESS, got %s", paidTx2.Status)
	}

	// Verify invoice is now fully PAID with balance 0
	clearedInv, _ := svc.GetInvoice(ctx, "ten_default", inv.ID)
	if clearedInv.Status != billing.InvoicePaid {
		t.Errorf("expected status PAID, got %s", clearedInv.Status)
	}
	if clearedInv.BalanceAmount != 0 {
		t.Errorf("expected balance 0, got %d", clearedInv.BalanceAmount)
	}

	// Verify student total outstanding balance is 0
	studentBalance, err := svc.GetStudentBalance(ctx, "ten_default", "stu_202")
	if err != nil {
		t.Fatalf("GetStudentBalance failed: %v", err)
	}
	if studentBalance != 0 {
		t.Errorf("expected student balance 0, got %d", studentBalance)
	}
}

func TestProcessRefund_ReversalsAndLedger(t *testing.T) {
	svc, _, _, _, ledgerRepo, _ := setupBillingTestService()
	ctx := context.Background()

	// Create invoice and pay 50,000
	inv, _ := svc.GenerateInvoice(ctx, billing.GenerateInvoiceCommand{
		TenantID:     "ten_default",
		StudentID:    "stu_refund_test",
		AcademicYear: "2026-2027",
		Semester:     1,
		DueDate:      "2026-10-15",
		CustomItems: []billing.FeeItemInput{
			{Category: billing.CategoryLibrary, Name: "Book Security Deposit", Amount: 50000},
		},
	})

	tx, _ := svc.InitiatePayment(ctx, billing.InitiatePaymentCommand{
		TenantID:  "ten_default",
		InvoiceID: inv.ID,
		StudentID: "stu_refund_test",
		Amount:    50000,
		Method:    billing.MethodOnlineGateway,
	})
	_, _ = svc.ProcessPaymentSuccess(ctx, billing.ProcessPaymentSuccessCommand{
		TenantID:      "ten_default",
		TransactionID: tx.ID,
		PaidAt:        time.Now().UTC(),
	})

	// Check invoice is PAID
	paidInv, _ := svc.GetInvoice(ctx, "ten_default", inv.ID)
	if paidInv.Status != billing.InvoicePaid {
		t.Fatalf("expected invoice to be PAID before refund")
	}

	// Process Refund
	refundedTx, err := svc.ProcessRefund(ctx, "ten_default", tx.ID, "Student withdrew admission")
	if err != nil {
		t.Fatalf("ProcessRefund failed: %v", err)
	}
	if refundedTx.Status != billing.PaymentRefunded {
		t.Errorf("expected status REFUNDED, got %s", refundedTx.Status)
	}

	// Verify invoice balance restored
	restoredInv, _ := svc.GetInvoice(ctx, "ten_default", inv.ID)
	if restoredInv.PaidAmount != 0 || restoredInv.BalanceAmount != 50000 {
		t.Errorf("expected paid=0 balance=50000, got paid=%d balance=%d", restoredInv.PaidAmount, restoredInv.BalanceAmount)
	}
	if restoredInv.Status != billing.InvoiceIssued {
		t.Errorf("expected status ISSUED after full refund, got %s", restoredInv.Status)
	}

	// Verify reversal ledger entries
	entries, _ := ledgerRepo.ListEntriesByTransaction(ctx, "ten_default", tx.ID)
	if len(entries) < 4 { // 2 original + 2 reversal
		t.Errorf("expected at least 4 ledger entries (2 initial + 2 reversal), got %d", len(entries))
	}
}

func TestInvoiceCancel_Guards(t *testing.T) {
	svc, _, _, _, _, _ := setupBillingTestService()
	ctx := context.Background()

	inv, _ := svc.GenerateInvoice(ctx, billing.GenerateInvoiceCommand{
		TenantID:     "ten_default",
		StudentID:    "stu_cancel_test",
		AcademicYear: "2026-2027",
		Semester:     1,
		DueDate:      "2026-10-15",
		CustomItems: []billing.FeeItemInput{
			{Category: billing.CategoryTuition, Name: "Tuition", Amount: 10000},
		},
	})

	// Cancel issued invoice with 0 payments -> succeeds
	if err := inv.Cancel("Administrative adjustment"); err != nil {
		t.Errorf("expected successful cancellation, got %v", err)
	}
	if inv.Status != billing.InvoiceCancelled {
		t.Errorf("expected status CANCELLED, got %s", inv.Status)
	}

	// An invoice with existing payments cannot be cancelled
	invWithPayment := &billing.StudentInvoice{
		Status:        billing.InvoicePartiallyPaid,
		TotalAmount:   20000,
		PaidAmount:    10000,
		BalanceAmount: 10000,
	}
	if err := invWithPayment.Cancel("Illegal attempt"); err == nil {
		t.Errorf("expected error cancelling invoice with payments")
	}
}
