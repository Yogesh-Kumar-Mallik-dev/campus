/**
 * BLOCK_API_HELPDESK_HANDLER_TEST_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   HTTP integration tests verifying REST endpoints for ticket lifecycle, category management, live chat messages, and rating desk.
 */

package helpdesk_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apiHelpdesk "campus/api/http/helpdesk"
	backendHelpdesk "campus/backend/helpdesk"
)

func setupTestRouter() (chi.Router, backendHelpdesk.Service) {
	r := chi.NewRouter()
	repo := backendHelpdesk.NewMockRepository()
	svc := backendHelpdesk.NewService(repo, nil)
	handler := apiHelpdesk.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_HelpdeskCategoryAndTicketLifecycle(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Create Category
	catPayload := map[string]interface{}{
		"tenant_id":            "tenant-http-1",
		"code":                 "ACADEMIC_EXAM",
		"name":                 "Academic & Exam Cell",
		"description":          "Inquiries regarding semester examinations and grade transcripts",
		"default_priority":     "HIGH",
		"sla_response_hours":   12,
		"sla_resolution_hours": 48,
		"actor_id":             "admin-1",
	}
	body, _ := json.Marshal(catPayload)
	req := httptest.NewRequest("POST", "/api/v1/helpdesk/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for category, got %d: %s", rec.Code, rec.Body.String())
	}

	var cat backendHelpdesk.HelpdeskCategory
	_ = json.Unmarshal(rec.Body.Bytes(), &cat)

	// 2. Create Ticket
	ticketPayload := map[string]interface{}{
		"tenant_id":    "tenant-http-1",
		"category_id":  cat.ID,
		"requester_id": "student-100",
		"title":        "Transcript grade discrepancy in CS301",
		"description":  "Grade shown as B instead of A in published marksheet",
	}
	ticketBody, _ := json.Marshal(ticketPayload)
	reqTicket := httptest.NewRequest("POST", "/api/v1/helpdesk/tickets", bytes.NewReader(ticketBody))
	reqTicket.Header.Set("Content-Type", "application/json")
	recTicket := httptest.NewRecorder()
	router.ServeHTTP(recTicket, reqTicket)

	if recTicket.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for ticket, got %d: %s", recTicket.Code, recTicket.Body.String())
	}

	var ticket backendHelpdesk.HelpdeskTicket
	_ = json.Unmarshal(recTicket.Body.Bytes(), &ticket)

	// 3. Assign Ticket
	assignPayload := map[string]interface{}{
		"tenant_id": "tenant-http-1",
		"staff_id":  "staff-exam-officer",
		"actor_id":  "admin-1",
	}
	assignBody, _ := json.Marshal(assignPayload)
	reqAssign := httptest.NewRequest("POST", "/api/v1/helpdesk/tickets/"+ticket.ID+"/assign", bytes.NewReader(assignBody))
	reqAssign.Header.Set("Content-Type", "application/json")
	recAssign := httptest.NewRecorder()
	router.ServeHTTP(recAssign, reqAssign)

	if recAssign.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for assign ticket, got %d: %s", recAssign.Code, recAssign.Body.String())
	}

	// 4. Add Message from Staff
	msgPayload := map[string]interface{}{
		"tenant_id":        "tenant-http-1",
		"sender_id":        "staff-exam-officer",
		"is_staff_reply":   true,
		"is_internal_note": false,
		"message":          "We are reviewing your evaluation copy with the instructor.",
	}
	msgBody, _ := json.Marshal(msgPayload)
	reqMsg := httptest.NewRequest("POST", "/api/v1/helpdesk/tickets/"+ticket.ID+"/messages", bytes.NewReader(msgBody))
	reqMsg.Header.Set("Content-Type", "application/json")
	recMsg := httptest.NewRecorder()
	router.ServeHTTP(recMsg, reqMsg)

	if recMsg.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for message, got %d: %s", recMsg.Code, recMsg.Body.String())
	}

	// 5. Update Status to RESOLVED
	statusPayload := map[string]interface{}{
		"tenant_id": "tenant-http-1",
		"status":    "RESOLVED",
		"actor_id":  "staff-exam-officer",
	}
	statusBody, _ := json.Marshal(statusPayload)
	reqStatus := httptest.NewRequest("POST", "/api/v1/helpdesk/tickets/"+ticket.ID+"/status", bytes.NewReader(statusBody))
	reqStatus.Header.Set("Content-Type", "application/json")
	recStatus := httptest.NewRecorder()
	router.ServeHTTP(recStatus, reqStatus)

	if recStatus.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for status update, got %d: %s", recStatus.Code, recStatus.Body.String())
	}

	// 6. Rate Ticket
	ratePayload := map[string]interface{}{
		"tenant_id": "tenant-http-1",
		"user_id":   "student-100",
		"rating":    5,
		"feedback":  "Very fast turnaround by exam wing. Thank you!",
	}
	rateBody, _ := json.Marshal(ratePayload)
	reqRate := httptest.NewRequest("POST", "/api/v1/helpdesk/tickets/"+ticket.ID+"/rate", bytes.NewReader(rateBody))
	reqRate.Header.Set("Content-Type", "application/json")
	recRate := httptest.NewRecorder()
	router.ServeHTTP(recRate, reqRate)

	if recRate.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for rating ticket, got %d: %s", recRate.Code, recRate.Body.String())
	}
}
