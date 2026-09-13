# ==============================================================================
# Universal Project Lifecycle Orchestrator (Windows PowerShell Entrypoint)
# Usage: .\script.ps1 [dev|build|check|test|deps|envi|uenvi|flush]
# ==============================================================================
[CmdletBinding()]
param (
    [Parameter(Position = 0)]
    [string]$Command = "dev",

    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$RemainingArgs
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

function Show-Help {
    Write-Host "Usage: .\script.ps1 <command> [args...]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Available commands:"
    Write-Host "  dev      Start local development servers and live watchers"
    Write-Host "  build    Build all production assets and binaries"
    Write-Host "  check    Run typecheckers, linters, and static analysis"
    Write-Host "  test     Run unit, integration, and contract test suites"
    Write-Host "  deps     Install repository dependencies"
    Write-Host "  envi     Initialize local environment and start Docker services"
    Write-Host "  uenvi    Tear down local Docker containers and cleanup"
    Write-Host "  flush    Flush database and caches for clean-slate testing"
    Write-Host "  help     Display this help message"
}

switch ($Command.ToLower()) {
    "dev" {
        $target = Join-Path $ScriptDir "scripts\dev.ps1"
        if (Test-Path $target) { & $target @RemainingArgs } else { Write-Host "Starting development mode..." }
    }
    "build" {
        $target = Join-Path $ScriptDir "scripts\build.ps1"
        if (Test-Path $target) { & $target @RemainingArgs } else { Write-Host "Building production bundles..." }
    }
    "check" {
        $target = Join-Path $ScriptDir "scripts\check.ps1"
        if (Test-Path $target) { & $target @RemainingArgs } else { Write-Host "Running typechecks and linters..." }
    }
    "test" {
        $target = Join-Path $ScriptDir "scripts\test.ps1"
        if (Test-Path $target) { & $target @RemainingArgs } else { Write-Host "Running test suites..." }
    }
    "deps" {
        $target = Join-Path $ScriptDir "scripts\deps.ps1"
        if (Test-Path $target) { & $target @RemainingArgs } else { Write-Host "Installing dependencies..." }
    }
    "envi" {
        $target = Join-Path $ScriptDir "scripts\envi.ps1"
        if (Test-Path $target) { & $target @RemainingArgs } else { Write-Host "Initializing environment..." }
    }
    "uenvi" {
        $target = Join-Path $ScriptDir "scripts\uenvi.ps1"
        if (Test-Path $target) { & $target @RemainingArgs } else { Write-Host "Tearing down environment..." }
    }
    "flush" {
        $target = Join-Path $ScriptDir "scripts\flush_db.ps1"
        if (Test-Path $target) { & $target @RemainingArgs } else { Write-Host "Flushing database..." }
    }
    { $_ -in "help", "--help", "-h" } {
        Show-Help
    }
    Default {
        Write-Warning "Unknown command: $Command"
        Show-Help
        exit 1
    }
}
