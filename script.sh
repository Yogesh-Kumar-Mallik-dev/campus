#!/usr/bin/env bash
# ==============================================================================
# Universal Project Lifecycle Orchestrator (POSIX Entrypoint)
# Usage: ./script.sh [dev|build|check|test|deps|envi|uenvi|flush]
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMMAND="${1:-dev}"
shift || true

function show_help() {
  echo "Usage: ./script.sh <command> [args...]"
  echo ""
  echo "Available commands:"
  echo "  dev      Start local development servers and live watchers"
  echo "  build    Build all production assets and binaries"
  echo "  check    Run typecheckers, linters, and static analysis"
  echo "  test     Run unit, integration, and contract test suites"
  echo "  deps     Install repository dependencies"
  echo "  envi     Initialize local environment and start Docker services"
  echo "  uenvi    Tear down local Docker containers and cleanup"
  echo "  flush    Flush database and caches for clean-slate testing"
  echo "  help     Display this help message"
}

case "$COMMAND" in
  dev)
    if [ -f "$SCRIPT_DIR/scripts/dev.sh" ]; then
      "$SCRIPT_DIR/scripts/dev.sh" "$@"
    else
      echo "Starting development mode..."
    fi
    ;;
  build)
    if [ -f "$SCRIPT_DIR/scripts/build.sh" ]; then
      "$SCRIPT_DIR/scripts/build.sh" "$@"
    else
      echo "Building production bundles..."
    fi
    ;;
  check)
    if [ -f "$SCRIPT_DIR/scripts/check.sh" ]; then
      "$SCRIPT_DIR/scripts/check.sh" "$@"
    else
      echo "Running typechecks and linters..."
    fi
    ;;
  test)
    if [ -f "$SCRIPT_DIR/scripts/test.sh" ]; then
      "$SCRIPT_DIR/scripts/test.sh" "$@"
    else
      echo "Running test suites..."
    fi
    ;;
  deps)
    if [ -f "$SCRIPT_DIR/scripts/deps.sh" ]; then
      "$SCRIPT_DIR/scripts/deps.sh" "$@"
    else
      echo "Installing dependencies..."
    fi
    ;;
  envi)
    if [ -f "$SCRIPT_DIR/scripts/envi.sh" ]; then
      "$SCRIPT_DIR/scripts/envi.sh" "$@"
    else
      echo "Initializing environment..."
    fi
    ;;
  uenvi)
    if [ -f "$SCRIPT_DIR/scripts/uenvi.sh" ]; then
      "$SCRIPT_DIR/scripts/uenvi.sh" "$@"
    else
      echo "Tearing down environment..."
    fi
    ;;
  flush)
    if [ -f "$SCRIPT_DIR/scripts/flush_db.sh" ]; then
      "$SCRIPT_DIR/scripts/flush_db.sh" "$@"
    else
      echo "Flushing database..."
    fi
    ;;
  help|--help|-h)
    show_help
    ;;
  *)
    echo "Unknown command: $COMMAND"
    show_help
    exit 1
    ;;
esac
