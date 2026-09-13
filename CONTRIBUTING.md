# Contributing Guidelines

Thank you for your interest in contributing to this project! We maintain rigorous standards for code quality, documentation precision, and testing.

---

## 1. Development Principles

1. **Single-Change Policy:** Each Pull Request (PR) and commit must address exactly **one logical change**—one feature, one bug fix, or one refactor.
2. **100% Test Coverage for Domain Logic:** All business logic, state machines, and API handlers must include co-located tests.
3. **Cross-Platform Tooling:** All lifecycle workflows are accessible via `./script.sh` (Unix) and `.\script.ps1` (Windows PowerShell).

---

## 2. Getting Started & Verification

```bash
# 1. Install dependencies
./script.sh deps

# 2. Run typechecking & static analysis
./script.sh check

# 3. Run automated test suites
./script.sh test

# 4. Start local development server
./script.sh dev
```

---

## 3. Pull Request & Commit Standards

1. **GPG Signed Commits:** All commits must be signed (`git commit -S -m "..."`).
2. **Conventional Commits:** Messages must follow the Conventional Commits format (`feat(scope): ...`, `fix(scope): ...`).
3. **Linear History:** Rebase your branch on `main` before submitting your PR (`git pull --rebase origin main`).
