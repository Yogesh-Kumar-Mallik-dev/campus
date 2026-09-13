/**
 * BLOCK_AUTH_TEST_MOCKS_001
 * Purpose: In-memory mock test doubles for fast, deterministic unit testing.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"context"
	"strings"
	"sync"
	"time"
)

type MockUserRepo struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users: make(map[string]*User),
	}
}

func (m *MockUserRepo) GetByID(ctx context.Context, tenantID, userID string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (m *MockUserRepo) GetByEmailNormalized(ctx context.Context, tenantID, emailNormalized string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, u := range m.users {
		if (tenantID == "" || u.TenantID == tenantID) && u.EmailNormalized == emailNormalized {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepo) GetByUsername(ctx context.Context, tenantID, username string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, u := range m.users {
		if (tenantID == "" || u.TenantID == tenantID) && u.Username == username {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepo) Create(ctx context.Context, user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepo) Update(ctx context.Context, user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepo) RecordFailedLogin(ctx context.Context, tenantID, userID string, lockoutDuration time.Duration, maxAttempts int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[userID]
	if !ok {
		return ErrUserNotFound
	}
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= maxAttempts {
		until := time.Now().UTC().Add(lockoutDuration)
		u.LockoutUntil = &until
	}
	return nil
}

func (m *MockUserRepo) ResetFailedLogins(ctx context.Context, tenantID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[userID]
	if !ok {
		return ErrUserNotFound
	}
	u.FailedLoginAttempts = 0
	u.LockoutUntil = nil
	return nil
}

func (m *MockUserRepo) UpdateLastLogin(ctx context.Context, tenantID, userID string, loginTime time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[userID]
	if !ok {
		return ErrUserNotFound
	}
	u.LastLoginAt = &loginTime
	return nil
}

type MockSessionRepo struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	families map[string]*RefreshTokenFamily
}

func NewMockSessionRepo() *MockSessionRepo {
	return &MockSessionRepo{
		sessions: make(map[string]*Session),
		families: make(map[string]*RefreshTokenFamily),
	}
}

func (m *MockSessionRepo) CreateSession(ctx context.Context, session *Session, family *RefreshTokenFamily) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[session.ID] = session
	m.families[family.FamilyID] = family
	return nil
}

func (m *MockSessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return s, nil
}

func (m *MockSessionRepo) GetTokenFamily(ctx context.Context, familyID string) (*RefreshTokenFamily, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.families[familyID]
	if !ok {
		return nil, ErrUnauthorized
	}
	return f, nil
}

func (m *MockSessionRepo) RotateTokenFamily(ctx context.Context, familyID string, oldHash, newHash string, newRotationCount int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.families[familyID]
	if !ok {
		return ErrUnauthorized
	}
	f.ActiveTokenHash = newHash
	f.RotationCount = newRotationCount
	f.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockSessionRepo) RevokeTokenFamily(ctx context.Context, familyID string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.families[familyID]
	if !ok {
		return ErrUnauthorized
	}
	now := time.Now().UTC()
	f.IsCompromised = true
	f.RevokedAt = &now
	f.RevokedReason = reason
	return nil
}

func (m *MockSessionRepo) RevokeSession(ctx context.Context, sessionID string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}
	now := time.Now().UTC()
	s.RevokedAt = &now
	s.RevokedReason = reason
	return nil
}

func (m *MockSessionRepo) RevokeAllUserSessions(ctx context.Context, tenantID, userID string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	for _, s := range m.sessions {
		if s.UserID == userID {
			s.RevokedAt = &now
			s.RevokedReason = reason
		}
	}
	for _, f := range m.families {
		if f.UserID == userID {
			f.RevokedAt = &now
			f.RevokedReason = reason
		}
	}
	return nil
}

func (m *MockSessionRepo) ListUserSessions(ctx context.Context, tenantID, userID string) ([]*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]*Session, 0)
	now := time.Now().UTC()
	for _, s := range m.sessions {
		if s.UserID == userID && s.IsActive(now) {
			res = append(res, s)
		}
	}
	return res, nil
}

type MockRolePermRepo struct {
	mu          sync.RWMutex
	userRoles   map[string][]*UserRole
	permissions map[string][]string
}

func NewMockRolePermRepo() *MockRolePermRepo {
	return &MockRolePermRepo{
		userRoles:   make(map[string][]*UserRole),
		permissions: make(map[string][]string),
	}
}

func (m *MockRolePermRepo) GetUserRoles(ctx context.Context, userID string) ([]*UserRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.userRoles[userID], nil
}

func (m *MockRolePermRepo) GetEffectivePermissions(ctx context.Context, userID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.permissions[userID], nil
}

func (m *MockRolePermRepo) AssignRoleToUser(ctx context.Context, userRole *UserRole) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userRoles[userRole.UserID] = append(m.userRoles[userRole.UserID], userRole)
	return nil
}

func (m *MockRolePermRepo) RemoveRoleFromUser(ctx context.Context, userRoleID string) error {
	return nil
}

func (m *MockRolePermRepo) GetRoleByKey(ctx context.Context, tenantID, roleKey string) (*Role, error) {
	return &Role{
		Key:      roleKey,
		TenantID: tenantID,
		Name:     strings.Title(roleKey),
	}, nil
}

type MockMFARepo struct {
	mu          sync.RWMutex
	enrollments map[string]*MFAEnrollment
}

func NewMockMFARepo() *MockMFARepo {
	return &MockMFARepo{
		enrollments: make(map[string]*MFAEnrollment),
	}
}

func (m *MockMFARepo) GetEnrollment(ctx context.Context, userID string) (*MFAEnrollment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enrollments[userID], nil
}

func (m *MockMFARepo) SaveEnrollment(ctx context.Context, enrollment *MFAEnrollment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enrollments[enrollment.UserID] = enrollment
	return nil
}

func (m *MockMFARepo) DeleteEnrollment(ctx context.Context, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.enrollments, userID)
	return nil
}
