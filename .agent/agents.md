# Agent Instructions

Welcome, Agents! When working in this repository, you must adhere to the following rules:

1. **Context Tracking:**
   Before ending your task, you must always record your current context, the last step you completed, and any ongoing plans or immediate next steps into `.agent/current_context.md`.
2. **Reviewing Context:**
   At the start of your task, always review `.agent/current_context.md` to understand where the previous agent left off.
3. **Single-Change Policy:**
   Only perform **one logical change at a time**—one feature, one refactor, or one fix per task/commit. Reject multi-change requests.
4. **Mandatory Hand-in-Hand Testing:**
   Whenever a new feature or helper is added, corresponding unit tests must be created in the exact same change.
5. **Signed Commits:**
   Always perform a `git commit -S -m "..."` immediately after completing any functional change or documentation update.
6. **Documentation Standards:**
   Adhere strictly to `rules/documentation_standards.md` and eliminate all 10 AI anti-patterns.
7. **Mobile-First Responsive Layouts:**
   Ensure all UI components and templates are fluid and resilient across viewports (280px to 4K).
