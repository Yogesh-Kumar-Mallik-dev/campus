/**
 * BLOCK_API_SERVER_MAIN_001
 * Purpose: Main Go HTTP server entrypoint initializing gateway, middlewares, and domain routes.
 * Domain:  API Gateway & Transport Layer
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"campus/api"
	authhttp "campus/api/http/auth"
	"campus/api/middleware"
	"campus/backend/auth"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-development-jwt-secret-key-must-be-changed-in-prod-32b!"
	}

	// 1. Initialize Injected Cryptographic Services
	hasher := auth.NewArgon2idHasher(nil)
	signer := auth.NewJWTSigner(jwtSecret, "campus.institution.edu", 15*time.Minute)
	totp := auth.NewStandardTOTPProvider()
	audit := &auth.NoopAuditPublisher{}

	// 2. Initialize Repositories (In-memory mock for dev / PG repository in prod)
	userRepo := auth.NewMockUserRepo()
	sessionRepo := auth.NewMockSessionRepo()
	rolePermRepo := auth.NewMockRolePermRepo()
	mfaRepo := auth.NewMockMFARepo()

	// 3. Initialize Auth Domain Service
	authService := auth.NewAuthService(
		userRepo,
		sessionRepo,
		rolePermRepo,
		mfaRepo,
		hasher,
		signer,
		totp,
		audit,
		nil,
	)

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
		ExposedHeaders:   []string{"Link", "Location", "ETag", "RateLimit-Limit", "RateLimit-Remaining", "RateLimit-Reset", "Retry-After", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Register Health Probes
	api.RegisterHealthProbes(r, func() error {
		return nil // Active DB ping check
	})

	// Register Domain Subsystem Routes
	authhttp.RegisterRoutes(r, authHandler, authMiddleware)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
}
