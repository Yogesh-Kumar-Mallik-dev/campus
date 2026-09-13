/**
 * BLOCK_API_EVENTS_HANDLER_TEST_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   HTTP integration tests verifying REST endpoints for events, venue conflicts, ticketing, and QR check-in.
 */

package events_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	apiEvents "campus/api/http/events"
	backendEvents "campus/backend/events"
)

func setupTestRouter() (chi.Router, backendEvents.Service) {
	r := chi.NewRouter()
	repo := backendEvents.NewMockRepository()
	svc := backendEvents.NewService(repo, nil)
	handler := apiEvents.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_EventsCreationAndTicketing(t *testing.T) {
	router, _ := setupTestRouter()

	startTime := time.Now().UTC().Add(48 * time.Hour).Format(time.RFC3339)
	endTime := time.Now().UTC().Add(52 * time.Hour).Format(time.RFC3339)
	regDeadline := time.Now().UTC().Add(40 * time.Hour).Format(time.RFC3339)

	// 1. Create Event
	eventPayload := map[string]interface{}{
		"tenant_id":             "tenant-test",
		"title":                 "Campus Hackathon 2026",
		"description":           "36-hour continuous software sprint with cloud sponsors.",
		"category":              "TECHNICAL",
		"venue_name":            "Innovation Hub Floor 2",
		"venue_capacity":        100,
		"start_time":            startTime,
		"end_time":              endTime,
		"registration_deadline": regDeadline,
		"organizer_id":          "staff-cs-coordinator",
		"is_ticketed":           true,
		"ticket_price":          0.0,
		"max_tickets":           100,
	}
	body, _ := json.Marshal(eventPayload)
	req := httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for event, got %d: %s", rec.Code, rec.Body.String())
	}

	var event backendEvents.CampusEvent
	_ = json.Unmarshal(rec.Body.Bytes(), &event)

	// 2. Book Ticket
	bookPayload := map[string]interface{}{
		"tenant_id":  "tenant-test",
		"student_id": "stu-hack-01",
	}
	bookBody, _ := json.Marshal(bookPayload)
	reqBook := httptest.NewRequest("POST", "/api/v1/events/"+event.ID+"/tickets", bytes.NewReader(bookBody))
	reqBook.Header.Set("Content-Type", "application/json")
	recBook := httptest.NewRecorder()
	router.ServeHTTP(recBook, reqBook)

	if recBook.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for ticket booking, got %d: %s", recBook.Code, recBook.Body.String())
	}

	var ticket backendEvents.EventTicket
	_ = json.Unmarshal(recBook.Body.Bytes(), &ticket)

	// 3. Gate Check-In by QR Ticket Code
	checkInPayload := map[string]interface{}{
		"tenant_id":   "tenant-test",
		"ticket_code": ticket.TicketCode,
		"staff_id":    "staff-gate-scanner",
	}
	checkInBody, _ := json.Marshal(checkInPayload)
	reqCheckIn := httptest.NewRequest("POST", "/api/v1/events/tickets/checkin", bytes.NewReader(checkInBody))
	reqCheckIn.Header.Set("Content-Type", "application/json")
	recCheckIn := httptest.NewRecorder()
	router.ServeHTTP(recCheckIn, reqCheckIn)

	if recCheckIn.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for gate check-in, got %d: %s", recCheckIn.Code, recCheckIn.Body.String())
	}
}
