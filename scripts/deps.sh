#!/usr/bin/env bash
# ==============================================================================
# Dependency Installation Orchestrator (Go Modules + pnpm Workspace)
# ==============================================================================
set -euo pipefail

echo "Installing project dependencies..."

# Download and tidy Go modules
echo "Downloading Go modules..."
go mod tidy
go mod download

# Install pnpm workspace dependencies
if command -v pnpm &> /dev/null; then
  echo "Installing pnpm workspace dependencies..."
  pnpm install || true
fi

echo "Dependencies up to date."
