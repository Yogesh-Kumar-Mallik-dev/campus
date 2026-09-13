# Developer Workstation Onboarding Guide

## 1. Prerequisites & Version Invariants

Ensure the following tools are installed on your workstation:

| Tool | Required Version | Verification Command | Installation Guide |
| :--- | :--- | :--- | :--- |
| **Node.js** | `22.x LTS` | `node -v` | [nodejs.org](https://nodejs.org) |
| **pnpm** | `9.x+` | `pnpm -v` | `npm install -g pnpm` |
| **Go** | `1.23+` | `go version` | [golang.org](https://golang.org/dl) |
| **Docker** | `24+` | `docker --version` | [docker.com](https://docker.com) |
| **Git** | `2.40+` | `git --version` | [git-scm.com](https://git-scm.com) |

---

## 2. 5-Minute Bootstrap

### Unix (Linux / macOS / WSL)
```bash
# Clone the repository
git clone git@github.com:your-org/your-repo.git
cd your-repo

# Bootstrap dependencies and run environment
./script.sh deps
./script.sh envi
./script.sh dev
```

### Windows (PowerShell)
```powershell
# Clone the repository
git clone git@github.com:your-org/your-repo.git
cd your-repo

# Bootstrap dependencies and run environment
.\script.ps1 deps
.\script.ps1 envi
.\script.ps1 dev
```

---

## 3. Verification & Diagnostic Commands

```bash
# Run all linters and typecheckers (Must pass with 0 errors / 0 warnings)
./script.sh check

# Run automated test suites (Must achieve 100% pass rate)
./script.sh test
```
