/**
 * BLOCK_MESS_UNIT_TESTS_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   Unit tests covering mess halls, weekly menus, QR dining tokens, double punch prevention, and leave rebates.
 */

package mess_test

import (
	"context"
	"testing"
	"time"

	"campus/backend/mess"
)

func TestMess_HallAndMenuManagement(t *testing.T) {
	repo := mess.NewMockRepository()
	svc := mess.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	// 1. Create Mess Hall
	hall, err := svc.CreateMessHall(ctx, tenantID, "North Dining Hall", "MESS-N1", "North Quad", "Sodexo Campus Dining", 500, nil)
	if err != nil {
		t.Fatalf("CreateMessHall failed: %v", err)
	}
	if hall.Code != "MESS-N1" {
		t.Errorf("Expected code MESS-N1, got %s", hall.Code)
	}

	// 2. Add Menu Item
	item, err := svc.AddMenuItem(
		ctx,
		tenantID,
		hall.ID,
		mess.Monday,
		mess.MealTypeLunch,
		mess.DietCategoryVeg,
		"Paneer Butter Masala & Butter Naan",
		"Includes Jeera Rice and Gulab Jamun",
		750,
	)
	if err != nil {
		t.Fatalf("AddMenuItem failed: %v", err)
	}
	if item.Calories != 750 {
		t.Errorf("Expected calories 750, got %d", item.Calories)
	}

	// 3. List Menu for Monday Lunch
	menu, err := svc.ListMenu(ctx, tenantID, hall.ID, mess.Monday, mess.MealTypeLunch)
	if err != nil {
		t.Fatalf("ListMenu failed: %v", err)
	}
	if len(menu) != 1 {
		t.Errorf("Expected 1 menu item, got %d", len(menu))
	}
}

func TestMess_StudentSubscriptionLifecycle(t *testing.T) {
	repo := mess.NewMockRepository()
	svc := mess.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	hall, _ := svc.CreateMessHall(ctx, tenantID, "Central Mess", "MESS-C", "Central Square", "Campus Caterers", 600, nil)
	studentID := "stu_aarav_01"

	start := time.Now().UTC()
	end := start.AddDate(0, 4, 0) // 4 months semester

	// 1. Subscribe student
	sub, err := svc.SubscribeStudent(ctx, tenantID, studentID, hall.ID, mess.PlanTypeSemester, mess.DietCategoryVeg, start, end, 4500)
	if err != nil {
		t.Fatalf("SubscribeStudent failed: %v", err)
	}
	if sub.Status != mess.SubscriptionStatusActive {
		t.Errorf("Expected status ACTIVE, got %s", sub.Status)
	}

	// 2. Duplicate subscription attempt should fail
	_, err = svc.SubscribeStudent(ctx, tenantID, studentID, hall.ID, mess.PlanTypeMonthly, mess.DietCategoryNonVeg, start, end, 5000)
	if err == nil {
		t.Fatalf("Expected ErrActiveSubscriptionExists on duplicate subscription, got nil")
	}
}

func TestMess_QRTokenGenerationAndRedemption(t *testing.T) {
	repo := mess.NewMockRepository()
	svc := mess.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	hall, _ := svc.CreateMessHall(ctx, tenantID, "South Mess", "MESS-S", "South Block", "Elite Foods", 400, nil)
	studentID := "stu_kavya_02"

	start := time.Now().UTC()
	end := start.AddDate(0, 6, 0)
	_, _ = svc.SubscribeStudent(ctx, tenantID, studentID, hall.ID, mess.PlanTypeSemester, mess.DietCategoryVeg, start, end, 4500)

	mealDate := "2026-09-13"

	// 1. Generate Breakfast Token
	tok, err := svc.GenerateDailyToken(ctx, tenantID, studentID, mealDate, mess.MealTypeBreakfast)
	if err != nil {
		t.Fatalf("GenerateDailyToken failed: %v", err)
	}
	if tok.Status != mess.TokenStatusGenerated {
		t.Errorf("Expected status GENERATED, got %s", tok.Status)
	}

	// 2. Scanning / Redeeming token at entrance
	readerID := "READER_GATE_01"
	redeemedTok, err := svc.RedeemDiningToken(ctx, tenantID, tok.TokenCode, readerID)
	if err != nil {
		t.Fatalf("RedeemDiningToken failed: %v", err)
	}
	if redeemedTok.Status != mess.TokenStatusRedeemed {
		t.Errorf("Expected status REDEEMED, got %s", redeemedTok.Status)
	}
	if redeemedTok.RedeemedAt == nil {
		t.Errorf("Expected RedeemedAt to be set")
	}

	// 3. Invariant: Double punch / re-redemption must be rejected
	_, err = svc.RedeemDiningToken(ctx, tenantID, tok.TokenCode, readerID)
	if err == nil {
		t.Fatalf("Expected ErrTokenAlreadyRedeemed on second scan, got nil")
	}
}

func TestMess_RebateCalculationAndGuards(t *testing.T) {
	repo := mess.NewMockRepository()
	svc := mess.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	hall, _ := svc.CreateMessHall(ctx, tenantID, "East Dining Hall", "MESS-E", "East Wing", "Campus Caterers", 400, nil)
	studentID := "stu_priya_03"

	now := time.Now().UTC()
	_, _ = svc.SubscribeStudent(ctx, tenantID, studentID, hall.ID, mess.PlanTypeSemester, mess.DietCategoryVeg, now, now.AddDate(0, 5, 0), 4500)

	// 1. Invalid Rebate: only 2 days (< 3 continuous days required)
	fromDate2 := now.Add(24 * time.Hour)
	toDate2 := fromDate2.Add(48 * time.Hour) // 2 days
	_, err := svc.ApplyRebate(ctx, tenantID, studentID, fromDate2, toDate2, "Weekend short trip", "")
	if err == nil {
		t.Fatalf("Expected ErrInsufficientRebateDays for 2 days absence, got nil")
	}

	// 2. Valid Rebate: 5 days continuous absence
	fromDate5 := now.Add(24 * time.Hour)
	toDate5 := fromDate5.Add(120 * time.Hour) // 5 days
	rebate, err := svc.ApplyRebate(ctx, tenantID, studentID, fromDate5, toDate5, "Festival leave", "gp_fest_01")
	if err != nil {
		t.Fatalf("ApplyRebate for 5 days failed: %v", err)
	}
	if rebate.TotalDays != 5 {
		t.Errorf("Expected 5 total days, got %d", rebate.TotalDays)
	}
	// Monthly 4500 -> dailyRate = 150 -> 5 * 150 = 750
	expectedRebate := 5.0 * 150.0
	if rebate.TotalRebateAmount != expectedRebate {
		t.Errorf("Expected total rebate %f, got %f", expectedRebate, rebate.TotalRebateAmount)
	}

	// 3. Warden/Manager Approves Rebate
	approvedRebate, err := svc.ApproveRebate(ctx, tenantID, rebate.ID, "usr_manager_01")
	if err != nil {
		t.Fatalf("ApproveRebate failed: %v", err)
	}
	if approvedRebate.Status != mess.RebateStatusApproved {
		t.Errorf("Expected status APPROVED, got %s", approvedRebate.Status)
	}
}

func TestMess_FeedbackRatings(t *testing.T) {
	repo := mess.NewMockRepository()
	svc := mess.NewService(repo, nil)
	ctx := context.Background()
	tenantID := "tenant_main"

	hall, _ := svc.CreateMessHall(ctx, tenantID, "West Hall", "MESS-W", "West Wing", "Daily Fresh Foods", 350, nil)

	fb, err := svc.SubmitFeedback(
		ctx,
		tenantID,
		"stu_rohan_04",
		hall.ID,
		"2026-09-13",
		mess.MealTypeDinner,
		4,
		5,
		"Excellent food and hygienic serving stations",
	)
	if err != nil {
		t.Fatalf("SubmitFeedback failed: %v", err)
	}
	if fb.FoodRating != 4 || fb.HygieneRating != 5 {
		t.Errorf("Expected ratings (4, 5), got (%d, %d)", fb.FoodRating, fb.HygieneRating)
	}

	feedbacks, err := svc.ListFeedbacks(ctx, tenantID, hall.ID)
	if err != nil {
		t.Fatalf("ListFeedbacks failed: %v", err)
	}
	if len(feedbacks) != 1 {
		t.Errorf("Expected 1 feedback, got %d", len(feedbacks))
	}
}
