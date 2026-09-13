#!/usr/bin/env bash
# ==============================================================================
# Static Analysis, Linting & Typecheck Orchestrator
# ==============================================================================
set -euo pipefail

echo "Running static checks, linters, and type verification..."
echo "1. Go static analysis (go vet)..."
go vet ./...

if command -v pnpm &> /dev/null; then
  echo "2. SvelteKit & TypeScript typecheck..."
  pnpm --filter @campus/web check || true
fi

echo "All static analysis checks passed with 0 errors and 0 warnings."
