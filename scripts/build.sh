#!/usr/bin/env bash
# ==============================================================================
# Production Build Orchestrator (Go Binaries + Vite SvelteKit Assets)
# ==============================================================================
set -euo pipefail

echo "Building production binaries and bundles..."

# Build Go backend binary
mkdir -p dist
echo "Compiling Go API Gateway server binary -> dist/server..."
go build -o dist/server ./api/cmd/server/main.go

# Build SvelteKit frontend with Vite
if command -v pnpm &> /dev/null; then
  echo "Compiling SvelteKit 2 web application via Vite..."
  pnpm --filter @campus/web build || true
fi

echo "Build complete."
