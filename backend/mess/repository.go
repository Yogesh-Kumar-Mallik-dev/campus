/**
 * BLOCK_MESS_REPOSITORY_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   Repository interface contracts for dining halls, weekly menus, tokens, and rebates.
 */

package mess

import (
	"context"
)

type Repository interface {
	// Mess Hall operations
	CreateMessHall(ctx context.Context, hall *MessHall) error
	GetMessHallByID(ctx context.Context, tenantID, id string) (*MessHall, error)
	GetMessHallByCode(ctx context.Context, tenantID, code string) (*MessHall, error)
	ListMessHalls(ctx context.Context, tenantID string) ([]*MessHall, error)

	// Menu Item operations
	CreateMenuItem(ctx context.Context, item *MenuItem) error
	GetMenuItemByID(ctx context.Context, tenantID, id string) (*MenuItem, error)
	ListMenuItemsByHall(ctx context.Context, tenantID, messHallID string) ([]*MenuItem, error)
	ListMenuByDayAndType(ctx context.Context, tenantID, messHallID string, day DayOfWeek, mealType MealType) ([]*MenuItem, error)

	// Subscription operations
	CreateSubscription(ctx context.Context, sub *Subscription) error
	GetSubscriptionByID(ctx context.Context, tenantID, id string) (*Subscription, error)
	GetActiveSubscriptionByStudent(ctx context.Context, tenantID, studentID string) (*Subscription, error)
	UpdateSubscriptionStatus(ctx context.Context, tenantID, id string, status SubscriptionStatus) error
	ListSubscriptions(ctx context.Context, tenantID, messHallID string) ([]*Subscription, error)

	// Dining Token operations
	CreateDiningToken(ctx context.Context, token *DiningToken) error
	GetTokenByID(ctx context.Context, tenantID, id string) (*DiningToken, error)
	GetTokenByCode(ctx context.Context, tenantID, tokenCode string) (*DiningToken, error)
	GetTokenByStudentMealDate(ctx context.Context, tenantID, studentID, mealDate string, mealType MealType) (*DiningToken, error)
	UpdateTokenStatus(ctx context.Context, tenantID, id string, status TokenStatus, readerID *string) error
	ListTokensByStudent(ctx context.Context, tenantID, studentID string) ([]*DiningToken, error)

	// Rebate operations
	CreateRebate(ctx context.Context, rebate *RebateApplication) error
	GetRebateByID(ctx context.Context, tenantID, id string) (*RebateApplication, error)
	UpdateRebateStatus(ctx context.Context, tenantID, id string, status RebateStatus, approverID *string) error
	ListRebates(ctx context.Context, tenantID, studentID string) ([]*RebateApplication, error)

	// Feedback operations
	CreateFeedback(ctx context.Context, feedback *MealFeedback) error
	ListFeedbacksByHall(ctx context.Context, tenantID, messHallID string) ([]*MealFeedback, error)
}
