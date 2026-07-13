---
name: "Developer Guidelines"
description: "Development rules and philosophy adopted by the team"
tools: []
model: ""
maxTurns: 0
skills: []
mcpServers: []
---

KISS — Keep It Simple, Stupid

Keep everything simple. Prefer direct and easy-to-understand solutions over complex abstractions. Simplicity facilitates code review, testing, and long-term maintenance.

Practical Principles

- Keep files small and with a single responsibility, not much larger than 300 lines.
  - Each file must have a clear focus; avoid massive files with multiple responsibilities.
- In Go, prefer small and cohesive functions/methods.
  - Extract functions whenever a piece of code reaches more than a few lines or does more than one thing.
  - Place related helper functions in separate files with descriptive names (e.g., db_migrations.go, workspace_store.go, agent_parser.go).
- Prefer composition over inheritance; keep interfaces small and specific.
- Write automated tests for critical features.
- Document important architectural decisions in the repository (README, docs/) — not just in the code.
- When introducing a new dependency, evaluate the cost/benefit and prefer small, active dependencies.

- **Strict Decoupling & Token Optimization:** Maintain absolute isolation between core sub-packages (e.g., indexer, LLM client, database orchestrator).
  - Sub-packages must never share direct domain models or cross-dependencies unless explicitly designed via lightweight primitives, clean interfaces, or clear API boundaries (JSON/HTTP).
  - When writing code inside a specific micro-utility or sub-package, do not attempt to solve orchestration problems or write business logic meant for other modules. Keep the codebase scope highly local to save LLM context tokens and prevent regression side effects.

Additional Best Practices

- Clear naming: choose descriptive names for packages, functions, and variables in English.
- Do not optimize prematurely: measure before optimizing.
- Make small commits with clear messages — prefer to review frequently.

Example (Go)

- If a service has multiple parts: initialization, migrations, handlers → place each part in a separate file:
  - init.go
  - migrations.go
  - handlers.go
  - store.go

Following these principles, we achieve more readable, testable, and maintainable code.

UI Policy and Fallbacks — Do Not Mask Errors

FAILURES MUST NOT BE MASKED BY AUTOMATIC FALLBACKS

- Rule: Do not implement fallbacks in the UI that hide or mask persistence, migration, or data loading problems. An empty dropdown or an unpopulated field must clearly expose that there is a backend issue (log + user-friendly error message), and the root cause must be fixed.

- Practical Motivation:
  - Fallbacks hide bugs and race conditions. When the UI presents a value via a "fallback," it is not evident that the primary data (e.g., normalized tables in the DB) failed to load correctly.
  - This hinders debugging and promotes technical debt accumulation: short-term fixes inadvertently become permanent solutions.

- Expected Behavior on Missing Data:
  - 1. The UI must clearly indicate an "incomplete" state (e.g., a message or icon stating "Missing data — check logs"), rather than filling the field with data from another invisible source.
  - 2. Record a log on both frontend and backend with sufficient context: endpoint called, timestamp, user action, workspace/ID, and any relevant payload.
  - 3. Provide a clear correction path (e.g., a "Reload data" button, instructions to re-run migrations, or a link to troubleshooting documentation).

- Debugging Procedures (High Priority):
  - 1. Check backend logs at startup — look for migration messages and logs like: "[DB] fixed_model loaded" and "[DB] SaveFixedModelRow".
  - 2. Query the DB table directly (sqlite3) to confirm the contents of fixed_models and fixed_model_tools.
  - 3. Verify if the engine has completed initialization before serving GetAdaConfig (to prevent race conditions).
  - 4. Confirm that provider_models were migrated to provider_models (GetProvidersFull) and that deadaptProviderConfig mapped the models to adaCfg.Providers.

- When a Fallback is Deemed Necessary (Exception):
  - It must be highly visible and transient: clearly display that it is a fallback (UI badge "fallback"), create an automatic ticket/alert, and expire the fallback after X minutes.
  - Always prefer to show the failure and demand a backend fix rather than hiding the problem.

- Secure Implementation Checklist (When adding any behavior involving derived data or models):
  - [ ] Is there sufficient logging in the backend to track data loading/migration?
  - [ ] Does the frontend validate the explicit presence of primary data before rendering (e.g., adaCfg.tiny_brain.provider !== undefined)?
  - [ ] In case of absence, does the UI present an error message and a reload button, instead of a list pre-populated by an invisible asset?
  - [ ] Is there a documented resolution path (in README/docs) for the migration/seed operation that populates the normalized tables?

By following this policy, we avoid masking problems and maintain visibility and correctness of root causes. If a dropdown is empty, we treat it as an error signal to be investigated — not as a reason to silence the bug with an invisible fallback.
