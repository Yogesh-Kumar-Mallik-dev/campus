/**
 * BLOCK_API_LIBRARY_HANDLER_001
 * Subsystem: Rank 9 - E-Library System (library)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for cataloging, checkouts, and overdue fines.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package library

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/library"
)

type Handler struct {
	service *library.Service
}

func NewHandler(service *library.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/library", func(r chi.Router) {
		// Books
		r.Post("/books", h.CatalogBook)
		r.Get("/books", h.ListBooks)

		// Copies
		r.Post("/books/{id}/copies", h.AddBookCopy)
		r.Get("/books/{id}/copies", h.ListCopies)

		// Borrows & Returns
		r.Post("/borrows", h.BorrowBook)
		r.Post("/borrows/{id}/return", h.ReturnBook)
		r.Post("/borrows/{id}/renew", h.RenewBook)
		r.Get("/borrows", h.ListBorrows)

		// Reservations & Ebooks
		r.Post("/reservations", h.ReserveBook)
		r.Post("/ebooks/{id}/progress", h.RecordEbookProgress)
	})
}

// CatalogBook handles POST /api/v1/library/books
func (h *Handler) CatalogBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID      string               `json:"tenant_id"`
		ISBN          string               `json:"isbn"`
		Title         string               `json:"title"`
		Author        string               `json:"author"`
		Publisher     string               `json:"publisher"`
		Edition       string               `json:"edition"`
		Category      library.BookCategory `json:"category"`
		ShelfLocation string               `json:"shelf_location"`
		EbookKey      *string              `json:"ebook_key"`
		EbookFormat   *string              `json:"ebook_format"`
		EbookSize     *int                 `json:"ebook_size"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	book, err := h.service.CatalogBook(r.Context(), req.TenantID, req.ISBN, req.Title, req.Author, req.Publisher, req.Edition, req.Category, req.ShelfLocation, req.EbookKey, req.EbookFormat, req.EbookSize)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"book": book})
}

// ListBooks handles GET /api/v1/library/books
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	var catPtr *library.BookCategory
	if c := r.URL.Query().Get("category"); c != "" {
		cat := library.BookCategory(c)
		catPtr = &cat
	}
	query := r.URL.Query().Get("q")

	books, err := h.service.ListBooks(r.Context(), tenantID, catPtr, query)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"books": books})
}

// AddBookCopy handles POST /api/v1/library/books/{id}/copies
func (h *Handler) AddBookCopy(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")
	var req struct {
		TenantID        string `json:"tenant_id"`
		AccessionNumber string `json:"accession_number"`
		Barcode         string `json:"barcode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	copy, err := h.service.AddBookCopy(r.Context(), req.TenantID, bookID, req.AccessionNumber, req.Barcode)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"copy": copy})
}

// ListCopies handles GET /api/v1/library/books/{id}/copies
func (h *Handler) ListCopies(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	bookID := chi.URLParam(r, "id")

	if tenantID == "" || bookID == "" {
		problem.BadRequest(w, r, "tenant_id and book id are required", "MISSING_PARAMS")
		return
	}

	copies, err := h.service.ListCopies(r.Context(), tenantID, bookID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"copies": copies})
}

// BorrowBook handles POST /api/v1/library/borrows
func (h *Handler) BorrowBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID   string `json:"tenant_id"`
		CopyID     string `json:"copy_id"`
		StudentID  string `json:"student_id"`
		IssuedByID string `json:"issued_by_id"`
		LoanDays   int    `json:"loan_days"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	borrow, err := h.service.BorrowBook(r.Context(), req.TenantID, req.CopyID, req.StudentID, req.IssuedByID, req.LoanDays)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"borrow": borrow})
}

// ReturnBook handles POST /api/v1/library/borrows/{id}/return
func (h *Handler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		TenantID  string  `json:"tenant_id"`
		DailyRate float64 `json:"daily_rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	borrow, err := h.service.ReturnBook(r.Context(), req.TenantID, id, req.DailyRate)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"borrow": borrow})
}

// RenewBook handles POST /api/v1/library/borrows/{id}/renew
func (h *Handler) RenewBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		TenantID       string `json:"tenant_id"`
		AdditionalDays int    `json:"additional_days"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	borrow, err := h.service.RenewBook(r.Context(), req.TenantID, id, req.AdditionalDays)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"borrow": borrow})
}

// ListBorrows handles GET /api/v1/library/borrows
func (h *Handler) ListBorrows(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	studentID := r.URL.Query().Get("student_id")

	var statusPtr *library.BorrowStatus
	if s := r.URL.Query().Get("status"); s != "" {
		st := library.BorrowStatus(s)
		statusPtr = &st
	}

	borrows, err := h.service.ListStudentBorrows(r.Context(), tenantID, studentID, statusPtr)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"borrows": borrows})
}

// ReserveBook handles POST /api/v1/library/reservations
func (h *Handler) ReserveBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID  string `json:"tenant_id"`
		BookID    string `json:"book_id"`
		StudentID string `json:"student_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	res, err := h.service.ReserveBook(r.Context(), req.TenantID, req.BookID, req.StudentID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"reservation": res})
}

// RecordEbookProgress handles POST /api/v1/library/ebooks/{id}/progress
func (h *Handler) RecordEbookProgress(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")
	var req struct {
		TenantID     string `json:"tenant_id"`
		StudentID    string `json:"student_id"`
		PageRead     int    `json:"page_read"`
		MinutesSpent int    `json:"minutes_spent"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	access, err := h.service.RecordEbookProgress(r.Context(), req.TenantID, bookID, req.StudentID, req.PageRead, req.MinutesSpent)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ebook_access": access})
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, library.ErrBookNotFound),
		errors.Is(err, library.ErrBookCopyNotFound),
		errors.Is(err, library.ErrBorrowRecordNotFound),
		errors.Is(err, library.ErrReservationNotFound),
		errors.Is(err, library.ErrEbookNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, library.ErrISBNExists),
		errors.Is(err, library.ErrAccessionExists),
		errors.Is(err, library.ErrCopyUnavailable),
		errors.Is(err, library.ErrBorrowLimitExceeded),
		errors.Is(err, library.ErrBorrowAlreadyReturned):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	case errors.Is(err, library.ErrMaxRenewalsExceeded),
		errors.Is(err, library.ErrBorrowOverdueRenewal):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	default:
		var domErr *library.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected library error occurred")
	}
}
