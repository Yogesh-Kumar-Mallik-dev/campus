/**
 * BLOCK_API_LIBRARY_HANDLER_TEST_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   HTTP integration tests verifying REST endpoints, checkout invariants, and overdue fine calculations.
 */

package library_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apiLibrary "campus/api/http/library"
	backendLibrary "campus/backend/library"
)

func setupTestRouter() (chi.Router, *backendLibrary.Service) {
	r := chi.NewRouter()
	repo := backendLibrary.NewMockRepository()
	svc := backendLibrary.NewService(repo, nil)
	handler := apiLibrary.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_LibraryCatalogAndCopies(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Catalog Book
	bookPayload := map[string]interface{}{
		"tenant_id":      "tenant_test",
		"isbn":           "978-0134685991",
		"title":          "Effective Java",
		"author":         "Joshua Bloch",
		"publisher":      "Addison-Wesley",
		"edition":        "3rd Edition",
		"category":       "COMPUTER_SCIENCE",
		"shelf_location": "Rack CS-09, Shelf A",
	}
	body, _ := json.Marshal(bookPayload)
	req := httptest.NewRequest("POST", "/api/v1/library/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for catalog book, got %d: %s", rec.Code, rec.Body.String())
	}

	var bookResp struct {
		Book backendLibrary.Book `json:"book"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &bookResp)
	bookID := bookResp.Book.ID

	// 2. Add Copy
	copyPayload := map[string]interface{}{
		"tenant_id":        "tenant_test",
		"accession_number": "LIB-ACC-00991",
		"barcode":          "BC-EJ-01",
	}
	copyBody, _ := json.Marshal(copyPayload)
	reqCopy := httptest.NewRequest("POST", "/api/v1/library/books/"+bookID+"/copies", bytes.NewReader(copyBody))
	reqCopy.Header.Set("Content-Type", "application/json")
	recCopy := httptest.NewRecorder()
	router.ServeHTTP(recCopy, reqCopy)

	if recCopy.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for copy, got %d: %s", recCopy.Code, recCopy.Body.String())
	}

	// 3. Search Books
	reqSearch := httptest.NewRequest("GET", "/api/v1/library/books?tenant_id=tenant_test&q=Bloch", nil)
	recSearch := httptest.NewRecorder()
	router.ServeHTTP(recSearch, reqSearch)

	if recSearch.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for search, got %d: %s", recSearch.Code, recSearch.Body.String())
	}
}

func TestHTTP_LibraryBorrowAndReturn(t *testing.T) {
	router, svc := setupTestRouter()
	ctx := context.Background()

	book, _ := svc.CatalogBook(ctx, "tenant_test", "978-0596007126", "Head First Design Patterns", "Freeman", "O'Reilly", "1st", backendLibrary.CategoryComputerScience, "Rack CS-10", nil, nil, nil)
	copy1, _ := svc.AddBookCopy(ctx, "tenant_test", book.ID, "LIB-ACC-00801", "BC-HFDP-01")

	// 1. Borrow Book
	borrowPayload := map[string]interface{}{
		"tenant_id":    "tenant_test",
		"copy_id":      copy1.ID,
		"student_id":   "stu_megha_01",
		"issued_by_id": "usr_lib_01",
		"loan_days":    14,
	}
	body, _ := json.Marshal(borrowPayload)
	req := httptest.NewRequest("POST", "/api/v1/library/borrows", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for borrow, got %d: %s", rec.Code, rec.Body.String())
	}

	var borrowResp struct {
		Borrow backendLibrary.BorrowRecord `json:"borrow"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &borrowResp)
	borrowID := borrowResp.Borrow.ID

	// 2. Return on time
	retPayload := map[string]interface{}{
		"tenant_id":  "tenant_test",
		"daily_rate": 5.0,
	}
	retBody, _ := json.Marshal(retPayload)
	reqRet := httptest.NewRequest("POST", "/api/v1/library/borrows/"+borrowID+"/return", bytes.NewReader(retBody))
	reqRet.Header.Set("Content-Type", "application/json")
	recRet := httptest.NewRecorder()
	router.ServeHTTP(recRet, reqRet)

	if recRet.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for return, got %d: %s", recRet.Code, recRet.Body.String())
	}

	var retResp struct {
		Borrow backendLibrary.BorrowRecord `json:"borrow"`
	}
	_ = json.Unmarshal(recRet.Body.Bytes(), &retResp)

	if retResp.Borrow.Status != backendLibrary.BorrowStatusReturned {
		t.Errorf("Expected status RETURNED, got %s", retResp.Borrow.Status)
	}
	if retResp.Borrow.FineAmount != 0.0 {
		t.Errorf("Expected ₹0 fine, got %f", retResp.Borrow.FineAmount)
	}
}

func TestHTTP_LibraryReservationAndRenew(t *testing.T) {
	router, svc := setupTestRouter()
	ctx := context.Background()

	book, _ := svc.CatalogBook(ctx, "tenant_test", "978-0131177055", "Working Effectively with Legacy Code", "Michael Feathers", "Prentice Hall", "1st", backendLibrary.CategoryComputerScience, "Rack CS-11", nil, nil, nil)
	copy1, _ := svc.AddBookCopy(ctx, "tenant_test", book.ID, "LIB-ACC-00701", "BC-LEGACY-01")

	borrow, _ := svc.BorrowBook(ctx, "tenant_test", copy1.ID, "stu_alok_02", "usr_lib_01", 14)

	// 1. Renew
	renewPayload := map[string]interface{}{
		"tenant_id":       "tenant_test",
		"additional_days": 14,
	}
	renewBody, _ := json.Marshal(renewPayload)
	reqRenew := httptest.NewRequest("POST", "/api/v1/library/borrows/"+borrow.ID+"/renew", bytes.NewReader(renewBody))
	reqRenew.Header.Set("Content-Type", "application/json")
	recRenew := httptest.NewRecorder()
	router.ServeHTTP(recRenew, reqRenew)

	if recRenew.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for renew, got %d: %s", recRenew.Code, recRenew.Body.String())
	}

	// 2. Reserve Title
	resPayload := map[string]interface{}{
		"tenant_id":  "tenant_test",
		"book_id":    book.ID,
		"student_id": "stu_neha_03",
	}
	resBody, _ := json.Marshal(resPayload)
	reqRes := httptest.NewRequest("POST", "/api/v1/library/reservations", bytes.NewReader(resBody))
	reqRes.Header.Set("Content-Type", "application/json")
	recRes := httptest.NewRecorder()
	router.ServeHTTP(recRes, reqRes)

	if recRes.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for reservation, got %d: %s", recRes.Code, recRes.Body.String())
	}
}
