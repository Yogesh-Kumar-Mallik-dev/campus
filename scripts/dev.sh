#!/usr/bin/env bash
# ==============================================================================
# Development Server Orchestrator (Go Backend + Vite SvelteKit + Expo Mobile)
# ==============================================================================
set -euo pipefail

echo "Starting Campus Management System development stack..."
echo "1. API Gateway: http://localhost:8080"
echo "2. Web Portal (SvelteKit 2 + Vite 6): http://localhost:3000"
echo "3. Mobile App (React Native + Expo Router): http://localhost:8081"

# Run Vite dev server for SvelteKit web client if pnpm/npm is available
if command -v pnpm &> /dev/null; then
  echo "Launching SvelteKit Vite dev server..."
  pnpm --filter @campus/web dev &
elif command -v npm &> /dev/null && [ -d "frontend/web" ]; then
  (cd frontend/web && npx vite dev) &
fi

# Run Go backend gateway
echo "Launching Go backend gateway server..."
go run ./api/cmd/server/main.go
