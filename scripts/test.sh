#!/usr/bin/env bash
# ==============================================================================
# Test Suite Orchestrator
# ==============================================================================
set -euo pipefail

echo "Executing test suites with coverage tracking..."
go test ./... -v -cover
echo "All tests passed with 100% success."
