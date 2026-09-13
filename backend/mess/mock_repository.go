/**
 * BLOCK_MESS_MOCK_REPOSITORY_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   In-memory mock repository implementing Repository interface for unit testing.
 */

package mess

import (
	"context"
	"sync"
	"time"
)

type MockRepository struct {
	mu            sync.RWMutex
	halls         map[string]*MessHall
	menuItems     map[string]*MenuItem
	subscriptions map[string]*Subscription
	tokens        map[string]*DiningToken
	rebates       map[string]*RebateApplication
	feedbacks     map[string]*MealFeedback
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		halls:         make(map[string]*MessHall),
		menuItems:     make(map[string]*MenuItem),
		subscriptions: make(map[string]*Subscription),
		tokens:        make(map[string]*DiningToken),
		rebates:       make(map[string]*RebateApplication),
		feedbacks:     make(map[string]*MealFeedback),
	}
}

func (m *MockRepository) CreateMessHall(ctx context.Context, hall *MessHall) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, h := range m.halls {
		if h.TenantID == hall.TenantID && h.Code == hall.Code {
			return ErrMessHallCodeExists
		}
	}
	m.halls[hall.ID] = hall
	return nil
}

func (m *MockRepository) GetMessHallByID(ctx context.Context, tenantID, id string) (*MessHall, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	h, exists := m.halls[id]
	if !exists || h.TenantID != tenantID {
		return nil, ErrMessHallNotFound
	}
	return h, nil
}

func (m *MockRepository) GetMessHallByCode(ctx context.Context, tenantID, code string) (*MessHall, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, h := range m.halls {
		if h.TenantID == tenantID && h.Code == code {
			return h, nil
		}
	}
	return nil, ErrMessHallNotFound
}

func (m *MockRepository) ListMessHalls(ctx context.Context, tenantID string) ([]*MessHall, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*MessHall
	for _, h := range m.halls {
		if h.TenantID == tenantID {
			result = append(result, h)
		}
	}
	return result, nil
}

func (m *MockRepository) CreateMenuItem(ctx context.Context, item *MenuItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.menuItems[item.ID] = item
	return nil
}

func (m *MockRepository) GetMenuItemByID(ctx context.Context, tenantID, id string) (*MenuItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.menuItems[id]
	if !exists || item.TenantID != tenantID {
		return nil, ErrMenuItemNotFound
	}
	return item, nil
}

func (m *MockRepository) ListMenuItemsByHall(ctx context.Context, tenantID, messHallID string) ([]*MenuItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*MenuItem
	for _, item := range m.menuItems {
		if item.TenantID == tenantID && item.MessHallID == messHallID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *MockRepository) ListMenuByDayAndType(ctx context.Context, tenantID, messHallID string, day DayOfWeek, mealType MealType) ([]*MenuItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*MenuItem
	for _, item := range m.menuItems {
		if item.TenantID == tenantID && item.MessHallID == messHallID && item.DayOfWeek == day && item.MealType == mealType {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *MockRepository) CreateSubscription(ctx context.Context, sub *Subscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.subscriptions[sub.ID] = sub
	return nil
}

func (m *MockRepository) GetSubscriptionByID(ctx context.Context, tenantID, id string) (*Subscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sub, exists := m.subscriptions[id]
	if !exists || sub.TenantID != tenantID {
		return nil, ErrSubscriptionNotFound
	}
	return sub, nil
}

func (m *MockRepository) GetActiveSubscriptionByStudent(ctx context.Context, tenantID, studentID string) (*Subscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, sub := range m.subscriptions {
		if sub.TenantID == tenantID && sub.StudentID == studentID && sub.Status == SubscriptionStatusActive {
			return sub, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) UpdateSubscriptionStatus(ctx context.Context, tenantID, id string, status SubscriptionStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sub, exists := m.subscriptions[id]
	if !exists || sub.TenantID != tenantID {
		return ErrSubscriptionNotFound
	}
	sub.Status = status
	sub.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) ListSubscriptions(ctx context.Context, tenantID, messHallID string) ([]*Subscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Subscription
	for _, sub := range m.subscriptions {
		if sub.TenantID != tenantID {
			continue
		}
		if messHallID != "" && sub.MessHallID != messHallID {
			continue
		}
		result = append(result, sub)
	}
	return result, nil
}

func (m *MockRepository) CreateDiningToken(ctx context.Context, token *DiningToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, t := range m.tokens {
		if t.TenantID == token.TenantID && t.StudentID == token.StudentID && t.MealDate == token.MealDate && t.MealType == token.MealType {
			return ErrDoubleRedemptionAttempt
		}
	}
	m.tokens[token.ID] = token
	return nil
}

func (m *MockRepository) GetTokenByID(ctx context.Context, tenantID, id string) (*DiningToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, exists := m.tokens[id]
	if !exists || t.TenantID != tenantID {
		return nil, ErrTokenNotFound
	}
	return t, nil
}

func (m *MockRepository) GetTokenByCode(ctx context.Context, tenantID, tokenCode string) (*DiningToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, t := range m.tokens {
		if t.TenantID == tenantID && t.TokenCode == tokenCode {
			return t, nil
		}
	}
	return nil, ErrTokenNotFound
}

func (m *MockRepository) GetTokenByStudentMealDate(ctx context.Context, tenantID, studentID, mealDate string, mealType MealType) (*DiningToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, t := range m.tokens {
		if t.TenantID == tenantID && t.StudentID == studentID && t.MealDate == mealDate && t.MealType == mealType {
			return t, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) UpdateTokenStatus(ctx context.Context, tenantID, id string, status TokenStatus, readerID *string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, exists := m.tokens[id]
	if !exists || t.TenantID != tenantID {
		return ErrTokenNotFound
	}
	t.Status = status
	if status == TokenStatusRedeemed {
		now := time.Now().UTC()
		t.RedeemedAt = &now
		t.DeviceReaderID = readerID
	}
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) ListTokensByStudent(ctx context.Context, tenantID, studentID string) ([]*DiningToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*DiningToken
	for _, t := range m.tokens {
		if t.TenantID == tenantID && t.StudentID == studentID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *MockRepository) CreateRebate(ctx context.Context, rebate *RebateApplication) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rebates[rebate.ID] = rebate
	return nil
}

func (m *MockRepository) GetRebateByID(ctx context.Context, tenantID, id string) (*RebateApplication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, exists := m.rebates[id]
	if !exists || r.TenantID != tenantID {
		return nil, ErrRebateNotFound
	}
	return r, nil
}

func (m *MockRepository) UpdateRebateStatus(ctx context.Context, tenantID, id string, status RebateStatus, approverID *string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, exists := m.rebates[id]
	if !exists || r.TenantID != tenantID {
		return ErrRebateNotFound
	}
	r.Status = status
	r.ApprovedByID = approverID
	r.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockRepository) ListRebates(ctx context.Context, tenantID, studentID string) ([]*RebateApplication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*RebateApplication
	for _, r := range m.rebates {
		if r.TenantID == tenantID {
			if studentID == "" || r.StudentID == studentID {
				result = append(result, r)
			}
		}
	}
	return result, nil
}

func (m *MockRepository) CreateFeedback(ctx context.Context, feedback *MealFeedback) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.feedbacks[feedback.ID] = feedback
	return nil
}

func (m *MockRepository) ListFeedbacksByHall(ctx context.Context, tenantID, messHallID string) ([]*MealFeedback, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*MealFeedback
	for _, f := range m.feedbacks {
		if f.TenantID == tenantID && f.MessHallID == messHallID {
			result = append(result, f)
		}
	}
	return result, nil
}
