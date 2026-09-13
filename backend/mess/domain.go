/**
 * BLOCK_MESS_DOMAIN_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   Domain entities, QR dining tokens, diet preferences, and leave-linked rebate invariants.
 */

package mess

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type MealType string

const (
	MealTypeBreakfast MealType = "BREAKFAST"
	MealTypeLunch     MealType = "LUNCH"
	MealTypeSnacks    MealType = "SNACKS"
	MealTypeDinner    MealType = "DINNER"
)

type DietCategory string

const (
	DietCategoryVeg    DietCategory = "VEG"
	DietCategoryNonVeg DietCategory = "NON_VEG"
	DietCategoryJain   DietCategory = "JAIN"
	DietCategoryVegan  DietCategory = "VEGAN"
)

type PlanType string

const (
	PlanTypeSemester PlanType = "SEMESTER"
	PlanTypeMonthly  PlanType = "MONTHLY"
	PlanTypePerMeal  PlanType = "PER_MEAL"
	PlanTypeFlexi    PlanType = "FLEXI"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusPaused    SubscriptionStatus = "PAUSED"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
)

type TokenStatus string

const (
	TokenStatusGenerated TokenStatus = "GENERATED"
	TokenStatusRedeemed  TokenStatus = "REDEEMED"
	TokenStatusExpired   TokenStatus = "EXPIRED"
	TokenStatusCancelled TokenStatus = "CANCELLED"
)

type RebateStatus string

const (
	RebateStatusPending  RebateStatus = "PENDING"
	RebateStatusApproved RebateStatus = "APPROVED"
	RebateStatusRejected RebateStatus = "REJECTED"
	RebateStatusCredited RebateStatus = "CREDITED"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "MONDAY"
	Tuesday   DayOfWeek = "TUESDAY"
	Wednesday DayOfWeek = "WEDNESDAY"
	Thursday  DayOfWeek = "THURSDAY"
	Friday    DayOfWeek = "FRIDAY"
	Saturday  DayOfWeek = "SATURDAY"
	Sunday    DayOfWeek = "SUNDAY"
)

// MessHall represents a residential dining facility
type MessHall struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Capacity    int       `json:"capacity"`
	Location    string    `json:"location"`
	CatererName string    `json:"caterer_name"`
	ManagerID   *string   `json:"manager_id,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MenuItem represents weekly recurring dish items
type MenuItem struct {
	ID           string       `json:"id"`
	TenantID     string       `json:"tenant_id"`
	MessHallID   string       `json:"mess_hall_id"`
	DayOfWeek    DayOfWeek    `json:"day_of_week"`
	MealType     MealType     `json:"meal_type"`
	DietCategory DietCategory `json:"diet_category"`
	Title        string       `json:"title"`
	Description  string       `json:"description,omitempty"`
	Calories     int          `json:"calories"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// Subscription represents an active dining meal plan
type Subscription struct {
	ID             string             `json:"id"`
	TenantID       string             `json:"tenant_id"`
	StudentID      string             `json:"student_id"`
	MessHallID     string             `json:"mess_hall_id"`
	PlanType       PlanType           `json:"plan_type"`
	DietPreference DietCategory       `json:"diet_preference"`
	StartDate      time.Time          `json:"start_date"`
	EndDate        time.Time          `json:"end_date"`
	Status         SubscriptionStatus `json:"status"`
	MonthlyFee     float64            `json:"monthly_fee"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// DiningToken represents a digital QR token for a meal slot
type DiningToken struct {
	ID             string      `json:"id"`
	TenantID       string      `json:"tenant_id"`
	StudentID      string      `json:"student_id"`
	MessHallID     string      `json:"mess_hall_id"`
	MealType       MealType    `json:"meal_type"`
	TokenCode      string      `json:"token_code"`
	MealDate       string      `json:"meal_date"` // YYYY-MM-DD
	Status         TokenStatus `json:"status"`
	RedeemedAt     *time.Time  `json:"redeemed_at,omitempty"`
	DeviceReaderID *string     `json:"device_reader_id,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// RebateApplication represents a dining fee deduction for prolonged absences
type RebateApplication struct {
	ID                string       `json:"id"`
	TenantID          string       `json:"tenant_id"`
	StudentID         string       `json:"student_id"`
	SubscriptionID    string       `json:"subscription_id"`
	FromDate          time.Time    `json:"from_date"`
	ToDate            time.Time    `json:"to_date"`
	TotalDays         int          `json:"total_days"`
	DailyRate         float64      `json:"daily_rate"`
	TotalRebateAmount float64      `json:"total_rebate_amount"`
	Reason            string       `json:"reason"`
	GatePassID        *string      `json:"gate_pass_id,omitempty"`
	Status            RebateStatus `json:"status"`
	ApprovedByID      *string      `json:"approved_by_id,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// MealFeedback represents student satisfaction ratings
type MealFeedback struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	StudentID     string    `json:"student_id"`
	MessHallID    string    `json:"mess_hall_id"`
	MealDate      string    `json:"meal_date"` // YYYY-MM-DD
	MealType      MealType  `json:"meal_type"`
	FoodRating    int       `json:"food_rating"`
	HygieneRating int       `json:"hygiene_rating"`
	Comments      string    `json:"comments,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// ValidateRebateEligibility enforces minimum 3 continuous days for fee rebate
func ValidateRebateEligibility(fromDate, toDate time.Time, dailyRate float64) (int, float64, error) {
	if !toDate.After(fromDate) {
		return 0, 0, ErrInvalidDateRange
	}

	diffHours := toDate.Sub(fromDate).Hours()
	days := int(diffHours / 24)
	if days < 3 {
		return 0, 0, ErrInsufficientRebateDays
	}

	totalRebate := float64(days) * dailyRate
	return days, totalRebate, nil
}

// GenerateTokenCode generates a tamper-proof cryptographic token
func GenerateTokenCode(tenantID, studentID, mealDate string, mealType MealType) string {
	raw := fmt.Sprintf("%s:%s:%s:%s", tenantID, studentID, mealDate, mealType)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("TOK-%x", sum[:8])
}

// CanRedeemToken verifies token validity
func CanRedeemToken(t *DiningToken) error {
	if t.Status == TokenStatusRedeemed {
		return ErrTokenAlreadyRedeemed
	}
	if t.Status == TokenStatusExpired || t.Status == TokenStatusCancelled {
		return ErrTokenExpired
	}
	return nil
}
