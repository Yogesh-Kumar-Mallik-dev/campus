/**
 * BLOCK_MESS_SERVICE_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   Business domain service implementing dining plans, QR meal punch redemption, and rebate calculation.
 */

package mess

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"campus/backend/audit"
)

type Service struct {
	repo     Repository
	auditSub audit.Subscriber
}

func NewService(repo Repository, auditSub audit.Subscriber) *Service {
	return &Service{
		repo:     repo,
		auditSub: auditSub,
	}
}

func (s *Service) emitAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
	if s.auditSub == nil {
		return
	}
	var metaRaw json.RawMessage
	if meta != nil {
		if b, err := json.Marshal(meta); err == nil {
			metaRaw = b
		}
	}
	_ = s.auditSub.Enqueue(audit.RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      &actorID,
		ActorType:    audit.ActorType(actorType),
		Action:       action,
		ResourceType: resType,
		ResourceID:   &resID,
		Status:       status,
		Metadata:     metaRaw,
	})
}

// CreateMessHall registers a dining hall facility
func (s *Service) CreateMessHall(ctx context.Context, tenantID, name, code, location, catererName string, capacity int, managerID *string) (*MessHall, error) {
	if tenantID == "" || name == "" || code == "" {
		return nil, NewDomainError(nil, 400, "Bad Request", "tenant_id, name and code are required", "https://campus.internal/errors/invalid-argument")
	}

	hall := &MessHall{
		ID:          fmt.Sprintf("mess_%d", time.Now().UnixNano()),
		TenantID:    tenantID,
		Name:        name,
		Code:        code,
		Capacity:    capacity,
		Location:    location,
		CatererName: catererName,
		ManagerID:   managerID,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.CreateMessHall(ctx, hall); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, "SYSTEM", "SYSTEM", "mess:hall:created", "mess_hall", hall.ID, audit.StatusSuccess, map[string]interface{}{
		"code": code,
		"name": name,
	})

	return hall, nil
}

// ListMessHalls retrieves all mess halls for a tenant
func (s *Service) ListMessHalls(ctx context.Context, tenantID string) ([]*MessHall, error) {
	return s.repo.ListMessHalls(ctx, tenantID)
}

// AddMenuItem registers a meal dish item
func (s *Service) AddMenuItem(ctx context.Context, tenantID, messHallID string, day DayOfWeek, mealType MealType, diet DietCategory, title, description string, calories int) (*MenuItem, error) {
	if _, err := s.repo.GetMessHallByID(ctx, tenantID, messHallID); err != nil {
		return nil, err
	}

	item := &MenuItem{
		ID:           fmt.Sprintf("item_%d", time.Now().UnixNano()),
		TenantID:     tenantID,
		MessHallID:   messHallID,
		DayOfWeek:    day,
		MealType:     mealType,
		DietCategory: diet,
		Title:        title,
		Description:  description,
		Calories:     calories,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := s.repo.CreateMenuItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// ListMenu retrieves dishes for a specific day and meal slot
func (s *Service) ListMenu(ctx context.Context, tenantID, messHallID string, day DayOfWeek, mealType MealType) ([]*MenuItem, error) {
	if day != "" && mealType != "" {
		return s.repo.ListMenuByDayAndType(ctx, tenantID, messHallID, day, mealType)
	}
	return s.repo.ListMenuItemsByHall(ctx, tenantID, messHallID)
}

// SubscribeStudent assigns a student to a dining facility
func (s *Service) SubscribeStudent(ctx context.Context, tenantID, studentID, messHallID string, planType PlanType, diet DietCategory, startDate, endDate time.Time, monthlyFee float64) (*Subscription, error) {
	if _, err := s.repo.GetMessHallByID(ctx, tenantID, messHallID); err != nil {
		return nil, err
	}

	activeSub, err := s.repo.GetActiveSubscriptionByStudent(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if activeSub != nil {
		return nil, ErrActiveSubscriptionExists
	}

	sub := &Subscription{
		ID:             fmt.Sprintf("sub_%d", time.Now().UnixNano()),
		TenantID:       tenantID,
		StudentID:      studentID,
		MessHallID:     messHallID,
		PlanType:       planType,
		DietPreference: diet,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         SubscriptionStatusActive,
		MonthlyFee:     monthlyFee,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	if err := s.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, "SYSTEM", "SYSTEM", "mess:subscription:created", "mess_subscription", sub.ID, audit.StatusSuccess, map[string]interface{}{
		"student_id":   studentID,
		"mess_hall_id": messHallID,
		"diet":         string(diet),
	})

	return sub, nil
}

// GenerateDailyToken generates a single-use QR dining token for a meal slot
func (s *Service) GenerateDailyToken(ctx context.Context, tenantID, studentID, mealDate string, mealType MealType) (*DiningToken, error) {
	activeSub, err := s.repo.GetActiveSubscriptionByStudent(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if activeSub == nil {
		return nil, ErrSubscriptionNotFound
	}

	// Verify no double-generation for this slot
	existing, err := s.repo.GetTokenByStudentMealDate(ctx, tenantID, studentID, mealDate, mealType)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	tokenCode := GenerateTokenCode(tenantID, studentID, mealDate, mealType)

	token := &DiningToken{
		ID:         fmt.Sprintf("tok_%d", time.Now().UnixNano()),
		TenantID:   tenantID,
		StudentID:  studentID,
		MessHallID: activeSub.MessHallID,
		MealType:   mealType,
		TokenCode:  tokenCode,
		MealDate:   mealDate,
		Status:     TokenStatusGenerated,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	if err := s.repo.CreateDiningToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

// RedeemDiningToken processes QR punch at the mess entrance scanner
func (s *Service) RedeemDiningToken(ctx context.Context, tenantID, tokenCode, deviceReaderID string) (*DiningToken, error) {
	token, err := s.repo.GetTokenByCode(ctx, tenantID, tokenCode)
	if err != nil {
		return nil, err
	}

	if err := CanRedeemToken(token); err != nil {
		return nil, err
	}

	readerPtr := &deviceReaderID
	if err := s.repo.UpdateTokenStatus(ctx, tenantID, token.ID, TokenStatusRedeemed, readerPtr); err != nil {
		return nil, err
	}

	token.Status = TokenStatusRedeemed
	now := time.Now().UTC()
	token.RedeemedAt = &now
	token.DeviceReaderID = readerPtr

	s.emitAudit(tenantID, token.StudentID, "USER", "mess:token:redeemed", "mess_dining_token", token.ID, audit.StatusSuccess, map[string]interface{}{
		"meal_type": token.MealType,
		"meal_date": token.MealDate,
		"reader_id": deviceReaderID,
	})

	return token, nil
}

// ApplyRebate submits a dining fee deduction for prolonged absences (>= 3 days)
func (s *Service) ApplyRebate(ctx context.Context, tenantID, studentID string, fromDate, toDate time.Time, reason, gatePassID string) (*RebateApplication, error) {
	activeSub, err := s.repo.GetActiveSubscriptionByStudent(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if activeSub == nil {
		return nil, ErrSubscriptionNotFound
	}

	dailyRate := activeSub.MonthlyFee / 30.0
	days, totalRebate, err := ValidateRebateEligibility(fromDate, toDate, dailyRate)
	if err != nil {
		return nil, err
	}

	var gatePassPtr *string
	if gatePassID != "" {
		gatePassPtr = &gatePassID
	}

	rebate := &RebateApplication{
		ID:                fmt.Sprintf("reb_%d", time.Now().UnixNano()),
		TenantID:          tenantID,
		StudentID:         studentID,
		SubscriptionID:    activeSub.ID,
		FromDate:          fromDate,
		ToDate:            toDate,
		TotalDays:         days,
		DailyRate:         dailyRate,
		TotalRebateAmount: totalRebate,
		Reason:            reason,
		GatePassID:        gatePassPtr,
		Status:            RebateStatusPending,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := s.repo.CreateRebate(ctx, rebate); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, studentID, "USER", "mess:rebate:applied", "mess_rebate_application", rebate.ID, audit.StatusSuccess, map[string]interface{}{
		"total_days":   days,
		"rebate_total": totalRebate,
	})

	return rebate, nil
}

// ApproveRebate authorizes fee credit to student account
func (s *Service) ApproveRebate(ctx context.Context, tenantID, rebateID, approverID string) (*RebateApplication, error) {
	rebate, err := s.repo.GetRebateByID(ctx, tenantID, rebateID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateRebateStatus(ctx, tenantID, rebateID, RebateStatusApproved, &approverID); err != nil {
		return nil, err
	}

	rebate.Status = RebateStatusApproved
	rebate.ApprovedByID = &approverID

	s.emitAudit(tenantID, approverID, "USER", "mess:rebate:approved", "mess_rebate_application", rebate.ID, audit.StatusSuccess, map[string]interface{}{
		"student_id":   rebate.StudentID,
		"rebate_total": rebate.TotalRebateAmount,
	})

	return rebate, nil
}

// SubmitFeedback logs student meal ratings and feedback
func (s *Service) SubmitFeedback(ctx context.Context, tenantID, studentID, messHallID, mealDate string, mealType MealType, foodRating, hygieneRating int, comments string) (*MealFeedback, error) {
	fb := &MealFeedback{
		ID:            fmt.Sprintf("fb_%d", time.Now().UnixNano()),
		TenantID:      tenantID,
		StudentID:     studentID,
		MessHallID:    messHallID,
		MealDate:      mealDate,
		MealType:      mealType,
		FoodRating:    foodRating,
		HygieneRating: hygieneRating,
		Comments:      comments,
		CreatedAt:     time.Now().UTC(),
	}

	if err := s.repo.CreateFeedback(ctx, fb); err != nil {
		return nil, err
	}

	return fb, nil
}

// ListStudentTokens retrieves tokens for a student
func (s *Service) ListStudentTokens(ctx context.Context, tenantID, studentID string) ([]*DiningToken, error) {
	return s.repo.ListTokensByStudent(ctx, tenantID, studentID)
}

// ListStudentRebates retrieves rebate history for a student
func (s *Service) ListStudentRebates(ctx context.Context, tenantID, studentID string) ([]*RebateApplication, error) {
	return s.repo.ListRebates(ctx, tenantID, studentID)
}

// ListFeedbacks retrieves feedbacks for a mess hall
func (s *Service) ListFeedbacks(ctx context.Context, tenantID, messHallID string) ([]*MealFeedback, error) {
	return s.repo.ListFeedbacksByHall(ctx, tenantID, messHallID)
}
