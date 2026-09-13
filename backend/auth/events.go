/**
 * BLOCK_AUTH_EVENTS_001
 * Purpose: Domain event declarations and AuditPublisher contract for IAM events.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"context"
	"time"
)

// AuthEventType defines the specific domain event action.
type AuthEventType string

const (
	EventLoginSuccess         AuthEventType = "auth:login:success"
	EventLoginFailed          AuthEventType = "auth:login:failed"
	EventTokenRefreshed       AuthEventType = "auth:token:refreshed"
	EventTokenBreachDetected  AuthEventType = "auth:token:breach_detected"
	EventLogout               AuthEventType = "auth:logout"
	EventAllSessionsRevoked   AuthEventType = "auth:sessions:revoked_all"
	EventMFAEnrolled          AuthEventType = "auth:mfa:enrolled"
	EventMFADisabled          AuthEventType = "auth:mfa:disabled"
	EventPasswordChanged      AuthEventType = "auth:password:changed"
	EventAccountLocked        AuthEventType = "auth:account:locked"
)

// AuthAuditEvent represents a structured audit event emitted by the auth service.
type AuthAuditEvent struct {
	EventID   string        `json:"event_id"`
	Type      AuthEventType `json:"type"`
	TenantID  string        `json:"tenant_id"`
	UserID    string        `json:"user_id,omitempty"`
	Username  string        `json:"username,omitempty"`
	IPAddress string        `json:"ip_address,omitempty"`
	UserAgent string        `json:"user_agent,omitempty"`
	Severity  string        `json:"severity"` // "INFO", "WARN", "CRITICAL"
	Detail    string        `json:"detail"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// AuditPublisher is the Dependency Injection contract for publishing audit events asynchronously.
type AuditPublisher interface {
	PublishAuthEvent(ctx context.Context, event *AuthAuditEvent) error
}

// NoopAuditPublisher provides a default null implementation for testing.
type NoopAuditPublisher struct{}

func (n *NoopAuditPublisher) PublishAuthEvent(ctx context.Context, event *AuthAuditEvent) error {
	return nil
}
