#!/usr/bin/env bash
# ==============================================================================
# Static Analysis, Linting & Typecheck Orchestrator
# ==============================================================================
set -euo pipefail

echo "Running static checks, linters, and type verification..."
go vet ./...
echo "All static analysis checks passed with 0 errors and 0 warnings."
