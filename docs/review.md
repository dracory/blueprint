# Architectural Review

## Summary
- Background workers (task queue, cache expiry, sessions, scheduler) now inherit the application context and shut down gracefully alongside the HTTP server @main.go#78-225 @internal/app/stores_task.go#52-69 @github.com/dracory/taskstore/Store.go#121-224.
- Geo lookups in background tasks now enforce request timeouts and suppress sensitive payload logging, reducing the blast radius of third-party outages @internal/tasks/stats/stats_visitor_enhance_task.go#1-212.

## High Severity

## Medium Severity
1. **Scheduler jobs execute synchronously on critical paths**  
   `scheduleCleanUpTask` executes heavy clean-up synchronously within the scheduler handler, sharing the same goroutine as timing logic. Failures or long-running clean-up will block subsequent jobs. Offload into task queue when possible @internal/schedules/start_async.go#27-53 @internal/tasks/clean_up_task.go#1-70.

## Recently Mitigated
1. **Background processes lacked shutdown coordination (fixed)**  
   The task queue, cache expiry, session expiry, and scheduler routines now expose context-aware runners. `main.go` orchestrates them through a shared parent context, and each external module respects cancellation signals, preventing dangling goroutines during tests, CLI runs, or graceful shutdown @main.go#78-225 @github.com/dracory/taskstore/Store.go#121-224 @github.com/dracory/cachestore/store.go#109-185 @github.com/dracory/sessionstore/store.go#93-169.

2. **Configuration surface was hard to reason about (fixed)**  
   Configuration loading now delegates to focused helpers that collect validation errors before exiting. Related environment requirements (database, LLM providers, store toggles) are aggregated with clear messaging, providing actionable diagnostics instead of cascading panics @internal/config/load.go#82-299.

3. **Geo lookup lacked timeouts and leaked sensitive logs (fixed)**  
   `StatsVisitorEnhanceTask` now issues outbound requests with a five-second timeout, checks HTTP status codes, and trims the response without logging raw payloads. This prevents scheduler stalls when the upstream hangs and avoids emitting IP-derived data @internal/tasks/stats/stats_visitor_enhance_task.go#1-212.

## Low Severity / Observations
1. **Missing `docs/overview.md` referenced by tooling**  
   `.windsurf/rules.yaml` mandates reading `docs/overview.md`, but the file/dir is absent. Automations depending on it currently fail @.windsurf/rules.yaml#15-24.

2. **SQLite is configured with WAL but still reports SSL requirements**  
   `databaseOpen` enforces `SetSSLMode("require")` regardless of driver. While harmless for SQLite, it obscures actual TLS configuration for other drivers and deserves a follow-up to ensure the setting is honored @internal/app/database_open.go#18-49.

3. **Sparse automated testing**  
   Test coverage focuses on initialization helpers; there are no integration or controller tests ensuring routes, middlewares, or stores behave as expected @main_test.go#10-49.

## Recommendations & Next Steps
1. Wrap background and scheduled tasks in contexts with timeouts; use an HTTP client with deadlines and avoid logging raw geo responses.
2. Introduce an application lifecycle manager (context-aware runner or errgroup) to coordinate goroutines and graceful shutdown.
3. Break the monolithic config loader into domain-specific modules with clearer error reporting; document required env vars.
4. Restore the expected `docs/overview.md` (or update tooling) and expand test coverage, especially around routing and critical workflows.
