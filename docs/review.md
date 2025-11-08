# Architectural Review

## Summary

## High Severity
- None currently.

## Medium Severity
- None currently.

## Low Severity / Observations
1. **SQLite is configured with WAL but still reports SSL requirements**  
   `databaseOpen` enforces `SetSSLMode("require")` regardless of driver. While harmless for SQLite, it obscures actual TLS configuration for other drivers and deserves a follow-up to ensure the setting is honored @internal/app/database_open.go#18-49.

2. **Sparse automated testing**  
   Test coverage focuses on initialization helpers; there are no integration or controller tests ensuring routes, middlewares, or stores behave as expected @main_test.go#10-49.

## Recently Mitigated
1. **Scheduler jobs executed synchronously on critical paths (fixed)**  
   `scheduleCleanUpTask` now spins the clean-up handler in its own goroutine, keeping the scheduler loop responsive while still running the task immediately without queue indirection. This prevents long-running clean-up work from blocking future jobs @internal/schedules/start_async.go#33-51 @internal/tasks/clean_up_task.go#68-90.
2. **Background processes lacked shutdown coordination (fixed)**  
   The task queue, cache expiry, session expiry, and scheduler routines now expose context-aware runners. `main.go` orchestrates them through a shared parent context, and each external module respects cancellation signals, preventing dangling goroutines during tests, CLI runs, or graceful shutdown @main.go#78-225 @github.com/dracory/taskstore/Store.go#121-224 @github.com/dracory/cachestore/store.go#109-185 @github.com/dracory/sessionstore/store.go#93-169.

3. **Configuration surface was hard to reason about (fixed)**  
   Configuration loading now delegates to focused helpers that collect validation errors before exiting. Related environment requirements (database, LLM providers, store toggles) are aggregated with clear messaging, providing actionable diagnostics instead of cascading panics @internal/config/load.go#82-299.

4. **Geo lookup lacked timeouts and leaked sensitive logs (fixed)**  
   `StatsVisitorEnhanceTask` now issues outbound requests with a five-second timeout, checks HTTP status codes, and trims the response without logging raw payloads. This prevents scheduler stalls when the upstream hangs and avoids emitting IP-derived data @internal/tasks/stats/stats_visitor_enhance_task.go#1-212.

5. **Project overview doc missing from repository (fixed)**  
   Added `docs/overview.md`, aligning with `.windsurf/rules.yaml` expectations so tooling can load the required overview before sessions.

## Recommendations & Next Steps
1. Keep `docs/overview.md` aligned with project evolution and expand test coverage, especially around routing and critical workflows.
