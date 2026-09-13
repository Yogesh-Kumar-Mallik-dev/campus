/**
 * BLOCK_LIBRARY_UNIT_TESTS_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   Unit tests covering cataloging, borrowing quota guards, overdue fine calculations, renewals, and reservations.
 */

package library_test

import (
	"context"
	"testing"
	"time"

	"campus/backend/library"
)

func TestLibrary_CatalogBookAndAddCopies(t *testing.T) {
	repo := library.NewMockRepository()
	svc := library.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	// 1. Catalog Book Title
	book, err := svc.CatalogBook(
		ctx,
		tenantID,
		"978-0131103627",
		"The C Programming Language",
		"Brian W. Kernighan, Dennis M. Ritchie",
		"Prentice Hall",
		"2nd Edition",
		library.CategoryComputerScience,
		"Rack CS-01, Shelf B",
		nil,
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("CatalogBook failed: %v", err)
	}
	if book.ISBN != "978-0131103627" {
		t.Errorf("Expected ISBN 978-0131103627, got %s", book.ISBN)
	}

	// 2. Add Copies
	copy1, err := svc.AddBookCopy(ctx, tenantID, book.ID, "LIB-ACC-00101", "BC-9780131103627-01")
	if err != nil {
		t.Fatalf("AddBookCopy 1 failed: %v", err)
	}
	copy2, err := svc.AddBookCopy(ctx, tenantID, book.ID, "LIB-ACC-00102", "BC-9780131103627-02")
	if err != nil {
		t.Fatalf("AddBookCopy 2 failed: %v", err)
	}

	// 3. Verify total and available inventory updated
	updatedBook, _ := repo.GetBookByID(ctx, tenantID, book.ID)
	if updatedBook.TotalCopies != 2 || updatedBook.AvailableCopies != 2 {
		t.Errorf("Expected 2 copies available, got %d total and %d available", updatedBook.TotalCopies, updatedBook.AvailableCopies)
	}

	// 4. Search Books
	books, err := svc.ListBooks(ctx, tenantID, nil, "Kernighan")
	if err != nil {
		t.Fatalf("ListBooks failed: %v", err)
	}
	if len(books) != 1 {
		t.Errorf("Expected 1 match for Kernighan, got %d", len(books))
	}
	_ = copy1
	_ = copy2
}

func TestLibrary_BorrowAndReturnLifecycle_OnTime(t *testing.T) {
	repo := library.NewMockRepository()
	svc := library.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	book, _ := svc.CatalogBook(ctx, tenantID, "978-0262033848", "Introduction to Algorithms", "CLRS", "MIT Press", "3rd", library.CategoryComputerScience, "Rack CS-02", nil, nil, nil)
	copy1, _ := svc.AddBookCopy(ctx, tenantID, book.ID, "LIB-ACC-00201", "BC-CLRS-01")
	studentID := "stu_varun_01"

	// 1. Borrow Book
	borrow, err := svc.BorrowBook(ctx, tenantID, copy1.ID, studentID, "usr_librarian_01", 14)
	if err != nil {
		t.Fatalf("BorrowBook failed: %v", err)
	}
	if borrow.Status != library.BorrowStatusIssued {
		t.Errorf("Expected status ISSUED, got %s", borrow.Status)
	}

	// Verify inventory decremented
	bCheck, _ := repo.GetBookByID(ctx, tenantID, book.ID)
	if bCheck.AvailableCopies != 0 {
		t.Errorf("Expected 0 available copies after borrow, got %d", bCheck.AvailableCopies)
	}

	// 2. Return on time (0 days overdue)
	returned, err := svc.ReturnBook(ctx, tenantID, borrow.ID, 5.0)
	if err != nil {
		t.Fatalf("ReturnBook failed: %v", err)
	}
	if returned.Status != library.BorrowStatusReturned {
		t.Errorf("Expected status RETURNED, got %s", returned.Status)
	}
	if returned.FineAmount != 0.0 {
		t.Errorf("Expected ₹0 fine for on-time return, got %f", returned.FineAmount)
	}

	// Verify inventory restored
	bRestored, _ := repo.GetBookByID(ctx, tenantID, book.ID)
	if bRestored.AvailableCopies != 1 {
		t.Errorf("Expected 1 available copy restored, got %d", bRestored.AvailableCopies)
	}
}

func TestLibrary_BorrowLimitGuard(t *testing.T) {
	repo := library.NewMockRepository()
	svc := library.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"
	studentID := "stu_quota_tester"

	// Create 5 copies
	book, _ := svc.CatalogBook(ctx, tenantID, "978-1111111111", "Clean Code Series", "Uncle Bob", "Pearson", "1st", library.CategoryComputerScience, "Rack CS-03", nil, nil, nil)
	var copies []*library.BookCopy
	for i := 1; i <= 5; i++ {
		cpy, _ := svc.AddBookCopy(ctx, tenantID, book.ID, "LIB-ACC-0030"+string(rune('0'+i)), "BC-CC-"+string(rune('0'+i)))
		copies = append(copies, cpy)
	}

	// Borrow 4 books (Allowed)
	for i := 0; i < 4; i++ {
		_, err := svc.BorrowBook(ctx, tenantID, copies[i].ID, studentID, "usr_lib", 14)
		if err != nil {
			t.Fatalf("Borrowing book %d failed: %v", i+1, err)
		}
	}

	// Attempt to borrow 5th book (Should fail due to quota of 4)
	_, err := svc.BorrowBook(ctx, tenantID, copies[4].ID, studentID, "usr_lib", 14)
	if err == nil {
		t.Fatalf("Expected ErrBorrowLimitExceeded for 5th book, got nil")
	}
}

func TestLibrary_OverdueFineCalculation(t *testing.T) {
	repo := library.NewMockRepository()
	svc := library.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	book, _ := svc.CatalogBook(ctx, tenantID, "978-0132350884", "Clean Code", "Robert C. Martin", "Prentice Hall", "1st", library.CategoryComputerScience, "Rack CS-04", nil, nil, nil)
	copy1, _ := svc.AddBookCopy(ctx, tenantID, book.ID, "LIB-ACC-00401", "BC-CLEAN-01")

	borrow, _ := svc.BorrowBook(ctx, tenantID, copy1.ID, "stu_late_01", "usr_lib", 14)

	// Simulate due date was 6 days ago
	borrow.DueDate = time.Now().UTC().AddDate(0, 0, -6)
	_ = repo.UpdateBorrowRecord(ctx, borrow)

	// Return late
	dailyFine := 5.0
	returned, err := svc.ReturnBook(ctx, tenantID, borrow.ID, dailyFine)
	if err != nil {
		t.Fatalf("ReturnBook late failed: %v", err)
	}

	// 6 days * ₹5 = ₹30
	if returned.FineAmount < 30.0 {
		t.Errorf("Expected fine at least ₹30.0 for 6 days late, got %f", returned.FineAmount)
	}
}

func TestLibrary_RenewalLifecycleAndGuards(t *testing.T) {
	repo := library.NewMockRepository()
	svc := library.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	book, _ := svc.CatalogBook(ctx, tenantID, "978-0596517748", "JavaScript: The Good Parts", "Douglas Crockford", "O'Reilly", "1st", library.CategoryComputerScience, "Rack CS-05", nil, nil, nil)
	copy1, _ := svc.AddBookCopy(ctx, tenantID, book.ID, "LIB-ACC-00501", "BC-JS-01")

	borrow, _ := svc.BorrowBook(ctx, tenantID, copy1.ID, "stu_renew_01", "usr_lib", 14)

	// 1. First renewal (allowed)
	r1, err := svc.RenewBook(ctx, tenantID, borrow.ID, 14)
	if err != nil {
		t.Fatalf("Renew 1 failed: %v", err)
	}
	if r1.RenewCount != 1 {
		t.Errorf("Expected renew count 1, got %d", r1.RenewCount)
	}

	// 2. Second renewal (allowed)
	r2, err := svc.RenewBook(ctx, tenantID, borrow.ID, 14)
	if err != nil {
		t.Fatalf("Renew 2 failed: %v", err)
	}
	if r2.RenewCount != 2 {
		t.Errorf("Expected renew count 2, got %d", r2.RenewCount)
	}

	// 3. Third renewal (should be rejected: max renewals = 2)
	_, err = svc.RenewBook(ctx, tenantID, borrow.ID, 14)
	if err == nil {
		t.Fatalf("Expected ErrMaxRenewalsExceeded on 3rd renewal, got nil")
	}
}

func TestLibrary_ReservationAndEbookTelemetry(t *testing.T) {
	repo := library.NewMockRepository()
	svc := library.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	book, _ := svc.CatalogBook(ctx, tenantID, "978-0321125217", "Domain-Driven Design", "Eric Evans", "Addison-Wesley", "1st", library.CategoryComputerScience, "Rack CS-06", nil, nil, nil)

	// 1. Reserve Book
	res, err := svc.ReserveBook(ctx, tenantID, book.ID, "stu_eric_01")
	if err != nil {
		t.Fatalf("ReserveBook failed: %v", err)
	}
	if res.Status != library.ReservationStatusPending {
		t.Errorf("Expected status PENDING, got %s", res.Status)
	}

	// 2. Track Ebook progress
	access, err := svc.RecordEbookProgress(ctx, tenantID, book.ID, "stu_eric_01", 142, 45)
	if err != nil {
		t.Fatalf("RecordEbookProgress failed: %v", err)
	}
	if access.LastPageRead != 142 || access.TotalReadingMinutes != 45 {
		t.Errorf("Expected page 142 and 45 mins, got page %d and %d mins", access.LastPageRead, access.TotalReadingMinutes)
	}
}
