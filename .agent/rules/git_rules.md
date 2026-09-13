# Rule: Git Governance, Commit Conventions & Branch Lifecycle

## 1. Core Commit Invariants

1. **Strict Single-Change Policy (One Task, One Commit):**
   - Perform only **one logical change at a time**—one feature, one bug fix, or one refactor per task/commit.
   - If a prompt or request contains multiple unrelated changes, reject the combined execution, explain the single-change policy, and execute the changes sequentially.

2. **Cryptographically Signed Commits (GPG / SSH):**
   - Every commit must be signed using `git commit -S -m "..."`.
   - Never commit unsigned code to protected branches.

3. **Conventional Commit Message Standard:**
   - Commit messages must strictly follow the Conventional Commits specification:
     ```text
     <type>(<scope>): <concise description in imperative mood>

     [optional body explaining rationale, trade-offs, and non-obvious context]

     [optional footer(s): BREAKING CHANGE, Closes #123]
     ```

### Standard Commit Types

| Type | Purpose | Example |
| :--- | :--- | :--- |
| `feat` | New user-facing feature or capability | `feat(auth): add multi-factor TOTP authentication` |
| `fix` | Bug fix or error resolution | `fix(api): return 200 OK with empty array on clean query` |
| `refactor` | Code restructuring with zero behavioral change | `refactor(db): extract connection pool factory` |
| `perf` | Performance optimization | `perf(cache): add single-flight mutex to prevent thundering herd` |
| `test` | Adding or updating unit/integration tests | `test(engine): add AST formula cycle detection test` |
| `docs` | Documentation additions or updates | `docs(api): document RFC 7807 problem detail error models` |
| `chore` | Maintenance, dependencies, toolchain updates | `chore(deps): upgrade Go toolchain to 1.23.4` |
| `ci` | CI/CD pipeline and workflow script updates | `ci(actions): add automated test matrix` |

---

## 2. Pre-Commit Verification Gate

Before staging and committing any code:

1. **Run Linter & Typechecker:** All type-checkers and linters must pass with **0 errors and 0 warnings**.
2. **Run Full Test Suite:** 100% of unit and integration tests must pass.
3. **Format Code:** Run repository formatters (e.g. Prettier, `gofmt`, `rustfmt`, `black`) across all modified files.
4. **Clean Git Status:** Ensure no unwanted debug files, scratch logs, or temporary artifacts are staged.

---

## 3. Branching & Merge Strategy

1. **Branch Naming Conventions:**
   - Feature: `feat/<ticket-or-short-name>` (e.g. `feat/oauth-providers`)
   - Bugfix: `fix/<issue-name>` (e.g. `fix/session-leak`)
   - Refactor: `refactor/<scope>` (e.g. `refactor/grpc-transport`)
   - Documentation: `docs/<topic>` (e.g. `docs/api-guide`)

2. **Linear History Preference:**
   - Rebase feature branches on top of `main` before merging (`git pull --rebase origin main`).
   - Squash-and-merge or rebase-and-merge pull requests to preserve a clean, bisectable history.
   - Avoid empty merge commits (`Merge branch 'main' into ...`) on feature branches.

3. **Atomic Rollback Guarantee:**
   - Every commit must leave the repository in a fully buildable and testable state.
   - If a commit is reverted (`git revert <commit-sha>`), the rest of the application must continue to function without crashing.
