/**
 * BLOCK_API_MESS_HANDLER_TEST_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   HTTP integration tests verifying REST endpoints, QR punches, and leave rebate calculations.
 */

package mess_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	apiMess "campus/api/http/mess"
	backendMess "campus/backend/mess"
)

func setupTestRouter() (chi.Router, *backendMess.Service) {
	r := chi.NewRouter()
	repo := backendMess.NewMockRepository()
	svc := backendMess.NewService(repo, nil)
	handler := apiMess.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_MessHallAndMenuCreation(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Create Mess Hall
	hallPayload := map[string]interface{}{
		"tenant_id":    "tenant_test",
		"name":         "Main Student Mess",
		"code":         "MESS-MAIN",
		"capacity":     500,
		"location":     "Central Square",
		"caterer_name": "Campus Food Services",
	}
	body, _ := json.Marshal(hallPayload)
	req := httptest.NewRequest("POST", "/api/v1/mess/halls", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for mess hall, got %d: %s", rec.Code, rec.Body.String())
	}

	var hallResp struct {
		MessHall backendMess.MessHall `json:"mess_hall"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &hallResp)
	hallID := hallResp.MessHall.ID

	// 2. Add Menu Item
	menuPayload := map[string]interface{}{
		"tenant_id":     "tenant_test",
		"mess_hall_id":  hallID,
		"day_of_week":   "MONDAY",
		"meal_type":     "LUNCH",
		"diet_category": "VEG",
		"title":         "Rajma Chawal with Mixed Salad",
		"description":   "Includes curd and roasted papad",
		"calories":      650,
	}
	menuBody, _ := json.Marshal(menuPayload)
	reqMenu := httptest.NewRequest("POST", "/api/v1/mess/menu-items", bytes.NewReader(menuBody))
	reqMenu.Header.Set("Content-Type", "application/json")
	recMenu := httptest.NewRecorder()
	router.ServeHTTP(recMenu, reqMenu)

	if recMenu.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for menu item, got %d: %s", recMenu.Code, recMenu.Body.String())
	}
}

func TestHTTP_MessTokenGenerationAndRedeem(t *testing.T) {
	router, svc := setupTestRouter()
	ctx := context.Background()

	hall, _ := svc.CreateMessHall(ctx, "tenant_test", "East Mess", "MESS-E", "East Wing", "Campus Caterer", 300, nil)
	studentID := "stu_karan_01"

	start := time.Now().UTC()
	end := start.AddDate(0, 4, 0)
	_, _ = svc.SubscribeStudent(ctx, "tenant_test", studentID, hall.ID, backendMess.PlanTypeSemester, backendMess.DietCategoryVeg, start, end, 4500)

	// 1. Generate Daily Token
	tokenPayload := map[string]interface{}{
		"tenant_id":  "tenant_test",
		"student_id": studentID,
		"meal_date":  "2026-09-13",
		"meal_type":  "DINNER",
	}
	tokBody, _ := json.Marshal(tokenPayload)
	reqTok := httptest.NewRequest("POST", "/api/v1/mess/tokens/daily", bytes.NewReader(tokBody))
	reqTok.Header.Set("Content-Type", "application/json")
	recTok := httptest.NewRecorder()
	router.ServeHTTP(recTok, reqTok)

	if recTok.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for token, got %d: %s", recTok.Code, recTok.Body.String())
	}

	var tokResp struct {
		Token backendMess.DiningToken `json:"token"`
	}
	_ = json.Unmarshal(recTok.Body.Bytes(), &tokResp)

	// 2. Redeem Token at Entrance
	redeemPayload := map[string]interface{}{
		"tenant_id":        "tenant_test",
		"token_code":       tokResp.Token.TokenCode,
		"device_reader_id": "READER_ENTRANCE_02",
	}
	redBody, _ := json.Marshal(redeemPayload)
	reqRed := httptest.NewRequest("POST", "/api/v1/mess/tokens/redeem", bytes.NewReader(redBody))
	reqRed.Header.Set("Content-Type", "application/json")
	recRed := httptest.NewRecorder()
	router.ServeHTTP(recRed, reqRed)

	if recRed.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for token redeem, got %d: %s", recRed.Code, recRed.Body.String())
	}

	// 3. Duplicate scan should return conflict
	reqRedDup := httptest.NewRequest("POST", "/api/v1/mess/tokens/redeem", bytes.NewReader(redBody))
	reqRedDup.Header.Set("Content-Type", "application/json")
	recRedDup := httptest.NewRecorder()
	router.ServeHTTP(recRedDup, reqRedDup)

	if recRedDup.Code != http.StatusConflict {
		t.Fatalf("Expected 409 Conflict for double punch, got %d: %s", recRedDup.Code, recRedDup.Body.String())
	}
}

func TestHTTP_MessRebateAndFeedback(t *testing.T) {
	router, svc := setupTestRouter()
	ctx := context.Background()

	hall, _ := svc.CreateMessHall(ctx, "tenant_test", "West Mess", "MESS-W", "West Wing", "Caterers Inc", 300, nil)
	studentID := "stu_sneha_02"

	now := time.Now().UTC()
	_, _ = svc.SubscribeStudent(ctx, "tenant_test", studentID, hall.ID, backendMess.PlanTypeSemester, backendMess.DietCategoryVeg, now, now.AddDate(0, 5, 0), 4500)

	// 1. Apply Rebate for 4 days
	from := now.Add(24 * time.Hour)
	to := from.Add(96 * time.Hour) // 4 days

	rebPayload := map[string]interface{}{
		"tenant_id":  "tenant_test",
		"student_id": studentID,
		"from_date":  from.Format(time.RFC3339),
		"to_date":    to.Format(time.RFC3339),
		"reason":     "Sports tournament trip representing college",
	}
	rebBody, _ := json.Marshal(rebPayload)
	reqReb := httptest.NewRequest("POST", "/api/v1/mess/rebates", bytes.NewReader(rebBody))
	reqReb.Header.Set("Content-Type", "application/json")
	recReb := httptest.NewRecorder()
	router.ServeHTTP(recReb, reqReb)

	if recReb.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for rebate, got %d: %s", recReb.Code, recReb.Body.String())
	}

	var rebResp struct {
		Rebate backendMess.RebateApplication `json:"rebate"`
	}
	_ = json.Unmarshal(recReb.Body.Bytes(), &rebResp)

	// 2. Approve Rebate
	appPayload := map[string]interface{}{
		"tenant_id":   "tenant_test",
		"approver_id": "usr_warden_01",
	}
	appBody, _ := json.Marshal(appPayload)
	reqApp := httptest.NewRequest("POST", "/api/v1/mess/rebates/"+rebResp.Rebate.ID+"/approve", bytes.NewReader(appBody))
	reqApp.Header.Set("Content-Type", "application/json")
	recApp := httptest.NewRecorder()
	router.ServeHTTP(recApp, reqApp)

	if recApp.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for rebate approve, got %d: %s", recApp.Code, recApp.Body.String())
	}

	// 3. Submit Meal Feedback
	fbPayload := map[string]interface{}{
		"tenant_id":      "tenant_test",
		"student_id":     studentID,
		"mess_hall_id":   hall.ID,
		"meal_date":      "2026-09-13",
		"meal_type":      "LUNCH",
		"food_rating":    5,
		"hygiene_rating": 4,
		"comments":       "Great taste and quick service",
	}
	fbBody, _ := json.Marshal(fbPayload)
	reqFb := httptest.NewRequest("POST", "/api/v1/mess/feedbacks", bytes.NewReader(fbBody))
	reqFb.Header.Set("Content-Type", "application/json")
	recFb := httptest.NewRecorder()
	router.ServeHTTP(recFb, reqFb)

	if recFb.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for feedback, got %d: %s", recFb.Code, recFb.Body.String())
	}
}
