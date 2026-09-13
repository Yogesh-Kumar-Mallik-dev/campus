# Rule: Repository Documentation Standards & AI Anti-Pattern Prevention

## 1. Core Documentation Philosophy

Every document in any repository is an authoritative, first-class engineering artifact. Documentation must represent actual codebase behavior with 100% precision.

### Strict Prohibitions (AI Anti-Patterns)

Do NOT:

- ❌ **Invent or Hallucinate Commands/Flags:** Never document CLI arguments, flags, environment variables, or script names that do not exist in workspace configuration files (`package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml`, or shell scripts).
- ❌ **Use Lazy Placeholders:** Never write `// TODO: add logic here`, `// ... rest of code`, or generic hand-waving stubs in examples. Every code example must be complete, syntactically valid, and typed.
- ❌ **Use Robotic Marketing Fluff:** Avoid hyperbolic adjectives (_"revolutionary"_, _"seamless"_, _"cutting-edge"_, _"game-changing"_, _"empowers users with unmatched synergy"_). Write crisp, factual engineering prose.
- ❌ **Allow Code-Doc Drift:** Every documented endpoint route, query parameter, database column, or schema model must match the backend implementation verbatim.
- ❌ **Invert Technical Depth:** Do not spend paragraphs on generic programming concepts while glossing over complex algorithms, state machines, consensus protocols, or DAG cycle detection.
- ❌ **Write Broken/Hypothetical Paths:** Never output fake paths like `file:///path/to/...` or links to non-existent files. Always use workspace-relative paths or verified markdown links.
- ❌ **Ignore Failure Modes:** Never document only the happy path (`200 OK`). Always document RFC 7807 problem details, error codes, validation failures (`422`), authentication errors (`401`), authorization denials (`403`), and edge case recoveries.
- ❌ **Produce Monolithic Walls of Text:** Always use structured markdown with headings, comparison tables, GitHub alert callouts, and syntax-highlighted code blocks.
- ❌ **Use ASCII Art Box Drawings:** Never draw diagrams or tables using ASCII text boxes (`┌─┐`, `│ │`, `└─┘`, `+---+`). Always use proper tools: **Mermaid.js** for visual diagrams/DAGs/flows and **Markdown Tables** for structured tabular mappings.

---

## 2. Standard Documentation Requirements

1. **5-Minute Quickstart Guarantee:** Fast, copy-pasteable bootstrap commands with explicit, verified version numbers for all required toolchains.
2. **Cross-Platform Parity:** Document Unix (`.sh` / Bash / Zsh) and Windows (`.ps1` / PowerShell) commands equally.
3. **Atomic Synchronization:** Update documentation in the exact same commit as code changes.
4. **Diagnostic & Test Verification:** Before committing documentation, all repository diagnostics, type-checks, and test suites must pass with 0 errors and 0 warnings.
5. **GPG Signed Commits:** All documentation commits must be cryptographically signed (`git commit -S`).

---

## 3. Standard Repository Documentation Hierarchy

Every production-grade repository should adhere to this standardized file layout:

```text
Repository Root/
├── README.md                      # Primary entry point: Value proposition, architecture overview, quickstart
├── CONTRIBUTING.md                # Contributor workflows: Branching, PRs, signed commits, test execution
├── SECURITY.md                    # Vulnerability disclosure policy, response SLAs, defense-in-depth model
├── LICENSE                        # Open-source or proprietary licensing terms
├── script.sh                      # Single unified Unix entrypoint (dev, build, check, test, flush)
├── script.ps1                     # Single unified Windows PowerShell entrypoint
└── docs/
    ├── DOCUMENTATION_STANDARDS.md # Authoritative doc quality & anti-pattern rules
    ├── ARCHITECTURE.md            # System topology, subsystem breakdown, data flow
    ├── API_GUIDE.md               # API endpoints, request/response models, problem details
    ├── ONBOARDING.md              # Workstation setup for Linux, macOS, and Windows
    ├── changes.md                 # Chronological version changelog with migration notes
    └── design_decisions.md        # Architecture Decision Records (ADRs) with rationale & trade-offs
```
