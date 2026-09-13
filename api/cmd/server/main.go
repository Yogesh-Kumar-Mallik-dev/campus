/**
 * BLOCK_API_SERVER_MAIN_001
 * Purpose: Main Go HTTP server entrypoint initializing gateway, middlewares, and domain routes.
 * Domain:  API Gateway & Transport Layer
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"campus/api"
	audithttp "campus/api/http/audit"
	authhttp "campus/api/http/auth"
	onboardinghttp "campus/api/http/onboarding"
	"campus/api/middleware"
	"campus/backend/audit"
	"campus/backend/auth"
	"campus/backend/onboarding"
)

// AuthAuditBridge adapts auth domain events into the central audit ledger.
type AuthAuditBridge struct {
	auditService audit.Service
}

func (b *AuthAuditBridge) PublishAuthEvent(ctx context.Context, event *auth.AuthAuditEvent) error {
	status := audit.StatusSuccess
	if event.Severity == "WARN" || event.Severity == "CRITICAL" {
		status = audit.StatusFailure
	}

	var meta json.RawMessage
	if len(event.Metadata) > 0 {
		meta, _ = json.Marshal(event.Metadata)
	}

	var userID *string
	if event.UserID != "" {
		userID = &event.UserID
	}
	var ip *string
	if event.IPAddress != "" {
		ip = &event.IPAddress
	}
	var ua *string
	if event.UserAgent != "" {
		ua = &event.UserAgent
	}

	return b.auditService.RecordAsync(audit.RecordAuditRequest{
		TenantID:     event.TenantID,
		ActorID:      userID,
		ActorType:    audit.ActorTypeUser,
		Action:       string(event.Type),
		ResourceType: "auth_session",
		Status:       status,
		IPAddress:    ip,
		UserAgent:    ua,
		Metadata:     meta,
		Timestamp:    &event.Timestamp,
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-development-jwt-secret-key-must-be-changed-in-prod-32b!"
	}

	// 1. Initialize Rank 2: Central Audit & Compliance System
	auditRepo := audit.NewMockRepository()
	auditHasher := audit.NewSHA256Hasher()
	auditSubscriber := audit.NewAsyncSubscriber(auditRepo, auditHasher, audit.BatchConfig{
		QueueCapacity: 10000,
		BatchSize:     100,
		FlushInterval: 50 * time.Millisecond,
	})
	auditSubscriber.Start(context.Background())
	auditService := audit.NewService(auditRepo, auditHasher, auditSubscriber)

	// 2. Initialize Rank 1: Central Auth & Permissions (IAM)
	hasher := auth.NewArgon2idHasher(nil)
	signer := auth.NewJWTSigner(jwtSecret, "campus.institution.edu", 15*time.Minute)
	totp := auth.NewStandardTOTPProvider()
	auditBridge := &AuthAuditBridge{auditService: auditService}

	userRepo := auth.NewMockUserRepo()
	sessionRepo := auth.NewMockSessionRepo()
	rolePermRepo := auth.NewMockRolePermRepo()
	mfaRepo := auth.NewMockMFARepo()

	authService := auth.NewAuthService(
		userRepo,
		sessionRepo,
		rolePermRepo,
		mfaRepo,
		hasher,
		signer,
		totp,
		auditBridge,
		nil,
	)

	// 3. Initialize Rank 3: Student & Staff Registration System (Onboarding)
	onboardingAppRepo := onboarding.NewMockApplicantRepository()
	onboardingDocRepo := onboarding.NewMockDocumentRepository()
	onboardingAcadRepo := onboarding.NewMockAcademicRepository()
	onboardingProfRepo := onboarding.NewMockProfileRepository()
	onboardingSeqRepo := onboarding.NewMockSequenceRepository()
	onboardingSeqEngine := onboarding.NewSequenceEngine(onboardingSeqRepo)

	onboardingService := onboarding.NewService(
		onboardingAppRepo,
		onboardingDocRepo,
		onboardingAcadRepo,
		onboardingProfRepo,
		onboardingSeqEngine,
		auditSubscriber,
	)
	onboardingHandler := onboardinghttp.NewHandler(onboardingService)

	// 4. Initialize HTTP Handlers & Middlewares
	authHandler := authhttp.NewAuthHandler(authService)
	authMiddleware := authhttp.NewAuthMiddleware(signer)

	// 5. Build Chi Router Pipeline
	r := chi.NewRouter()

	// Gateway Hardened Middlewares
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.SecurityHeadersMiddleware)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.BodyLimitMiddleware(250 * 1024)) // 250 KB JSON limit
	r.Use(middleware.EnforceJSONContentType)
	r.Use(middleware.StructuredLoggingMiddleware)

	// Hardened CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "tauri://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "Idempotency-Key", "If-Match"},
		ExposedHeaders:   []string{"Link", "Location", "ETag", "RateLimit-Limit", "RateLimit-Remaining", "RateLimit-Reset", "Retry-After", "X-Request-ID", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Register Health Probes
	api.RegisterHealthProbes(r, func() error {
		return nil // Active DB ping check
	})

	// Register Domain Subsystem Routes
	authhttp.RegisterRoutes(r, authHandler, authMiddleware)
	audithttp.RegisterRoutes(r, auditService)
	onboardingHandler.RegisterRoutes(r)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful Shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Campus Management System Gateway listening on :%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-shutdownChan
	log.Println("Shutting down gateway gracefully...")

	// Drain audit subscriber
	auditSubscriber.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
}
