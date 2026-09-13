/**
 * BLOCK_API_HOSTEL_HANDLER_TEST_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   HTTP integration tests verifying REST endpoints, RFC 7807 problem details, and gate pass lifecycles.
 */

package hostel_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	apiHostel "campus/api/http/hostel"
	backendHostel "campus/backend/hostel"
)

func setupTestRouter() (chi.Router, *backendHostel.Service) {
	r := chi.NewRouter()
	repo := backendHostel.NewMockRepository()
	svc := backendHostel.NewService(repo, nil)
	handler := apiHostel.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_HostelBlockAndRoomCreation(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Create Block
	blockPayload := map[string]interface{}{
		"tenant_id":    "tenant_test",
		"name":         "Chanakya Hall of Residence",
		"code":         "BH-C",
		"gender":       "MALE",
		"total_floors": 4,
		"total_rooms":  100,
		"capacity":     200,
	}
	body, _ := json.Marshal(blockPayload)
	req := httptest.NewRequest("POST", "/api/v1/hostel/blocks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var blockResp struct {
		Block backendHostel.Block `json:"block"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &blockResp)
	blockID := blockResp.Block.ID

	// 2. List Blocks
	reqList := httptest.NewRequest("GET", "/api/v1/hostel/blocks?tenant_id=tenant_test", nil)
	recList := httptest.NewRecorder()
	router.ServeHTTP(recList, reqList)

	if recList.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d: %s", recList.Code, recList.Body.String())
	}

	// 3. Create Room in Block
	roomPayload := map[string]interface{}{
		"tenant_id":             "tenant_test",
		"block_id":              blockID,
		"room_number":           "301",
		"floor_number":          3,
		"room_type":             "DOUBLE",
		"is_ac":                 true,
		"base_fee_per_semester": 28000,
		"max_beds":              2,
	}
	roomBody, _ := json.Marshal(roomPayload)
	reqRoom := httptest.NewRequest("POST", "/api/v1/hostel/rooms", bytes.NewReader(roomBody))
	reqRoom.Header.Set("Content-Type", "application/json")
	recRoom := httptest.NewRecorder()
	router.ServeHTTP(recRoom, reqRoom)

	if recRoom.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for room, got %d: %s", recRoom.Code, recRoom.Body.String())
	}
}

func TestHTTP_HostelGatePassAndCurfewBreach(t *testing.T) {
	router, svc := setupTestRouter()

	block, _ := svc.CreateBlock(context.Background(), "tenant_test", "Sarojini Hostel", "GH-S", backendHostel.GenderFemale, 3, 50, 100, nil)

	now := time.Now().UTC()
	outTime := now.Add(1 * time.Hour)
	inTime := now.Add(4 * time.Hour)

	// 1. Apply Gate Pass
	gpPayload := map[string]interface{}{
		"tenant_id":         "tenant_test",
		"student_id":        "stu_simran_01",
		"block_id":          block.ID,
		"reason":            "Medical checkup at Apollo",
		"destination":       "Apollo Hospital",
		"emergency_contact": "+919988776655",
		"expected_out_at":   outTime.Format(time.RFC3339),
		"expected_in_at":    inTime.Format(time.RFC3339),
	}
	gpBody, _ := json.Marshal(gpPayload)
	req := httptest.NewRequest("POST", "/api/v1/hostel/gate-passes", bytes.NewReader(gpBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for GatePass, got %d: %s", rec.Code, rec.Body.String())
	}

	var gpResp struct {
		GatePass backendHostel.GatePass `json:"gate_pass"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &gpResp)
	passID := gpResp.GatePass.ID

	// 2. Review (Approve)
	revPayload := map[string]interface{}{
		"tenant_id": "tenant_test",
		"warden_id": "warden_roy",
		"approve":   true,
	}
	revBody, _ := json.Marshal(revPayload)
	reqRev := httptest.NewRequest("POST", "/api/v1/hostel/gate-passes/"+passID+"/review", bytes.NewReader(revBody))
	recRev := httptest.NewRecorder()
	router.ServeHTTP(recRev, reqRev)

	if recRev.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for Review, got %d: %s", recRev.Code, recRev.Body.String())
	}

	// 3. Exit
	exitPayload := map[string]interface{}{
		"tenant_id":     "tenant_test",
		"actual_out_at": outTime.Format(time.RFC3339),
	}
	exitBody, _ := json.Marshal(exitPayload)
	reqExit := httptest.NewRequest("POST", "/api/v1/hostel/gate-passes/"+passID+"/exit", bytes.NewReader(exitBody))
	recExit := httptest.NewRecorder()
	router.ServeHTTP(recExit, reqExit)

	if recExit.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for Exit, got %d: %s", recExit.Code, recExit.Body.String())
	}

	// 4. Return with Curfew Breach (2 hours late)
	actualIn := inTime.Add(2 * time.Hour)
	retPayload := map[string]interface{}{
		"tenant_id":    "tenant_test",
		"warden_id":    "warden_roy",
		"actual_in_at": actualIn.Format(time.RFC3339),
	}
	retBody, _ := json.Marshal(retPayload)
	reqRet := httptest.NewRequest("POST", "/api/v1/hostel/gate-passes/"+passID+"/return", bytes.NewReader(retBody))
	recRet := httptest.NewRecorder()
	router.ServeHTTP(recRet, reqRet)

	if recRet.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for Return, got %d: %s", recRet.Code, recRet.Body.String())
	}

	var retResp struct {
		GatePass       backendHostel.GatePass    `json:"gate_pass"`
		CurfewIncident *backendHostel.IncidentLog `json:"curfew_incident,omitempty"`
	}
	_ = json.Unmarshal(recRet.Body.Bytes(), &retResp)

	if retResp.CurfewIncident == nil {
		t.Fatalf("Expected CurfewIncident on 2-hour late return, got nil")
	}
	if retResp.CurfewIncident.FineAmount <= 0 {
		t.Errorf("Expected fine amount to be imposed, got %f", retResp.CurfewIncident.FineAmount)
	}
}
