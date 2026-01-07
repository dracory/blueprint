# Project Overview

## Brand Concept
Blueprint is a ready-to-adapt Go web application starter that gives teams
a reliable, opinionated foundation for MVC projects. It highlights
production-ready patterns while leaving room for product-specific
customization.

The main aim is to ship with "everything" by default: it is easier to remove
what you do not need than to add missing infrastructure later. End users should
delete what is not needed for their deployment.

## Mission
Provide engineering teams with a maintainable baseline that accelerates
delivery of secure, observable, and resilient web applications without
locking them into a rigid framework.

## Product Metaphor
Blueprint behaves like an architect's master plan: it lays out structural
systems (routing, data access, background work, configuration) so feature
teams can focus on bespoke rooms and finishes.

## Primary Personas
- **Founding engineer** – needs a dependable scaffold to launch customer
features rapidly while keeping operational risk low.
- **Platform maintainer** – curates shared infrastructure and documentation
 to ensure consistency across multiple deployments and teams.
- **Contributor or contractor** – joins midstream and relies on clear guides,
 tests, and scripts to become productive within a single working session.

## Application Summary
- Root entrypoint `cmd/server/main.go` wires configuration, datastore initialization, routing, background tasks, and
  graceful shutdown handling.
- `/internal/registry` owns runtime dependencies, including database access, caching, and store wiring.
- `/internal/routes` (and feature subpackages) collect HTTP handlers, middlewares, and view composition for
  the MVC layer.
- `/internal/tasks` and `/internal/schedules` coordinate asynchronous work, cron-like jobs, and task queue
  runners with context-aware shutdown.
- `/pkg` exposes reusable libraries (blog, testimonials, theme helpers) that can be imported by other
  applications or extensions.

## Operational Practices
- Configuration is centralized in `internal/config` with validation and environment-specific overrides.
- Testing emphasizes realistic integrations using SQLite-backed stores and utilities shared with the base
  project. Run `go test ./...` or invoke Taskfile targets for coverage and specialized suites.
- Deployment commands and development workflows are scripted via `taskfile.yml`, `cmd/` utilities, and CI
  pipelines defined under `.github/workflows`.

## Architecture Contracts

- **Composition root:** `cmd/server/main.go` constructs the `registry` via `registry.New(cfg)` and wires the HTTP router.
- **Registry ownership:** `registry` owns the database handle and store instances.
- **Caches:** caches are instance-scoped on the registry (do not rely on package-level cache globals).
- **Lifecycle:** `registry.Close()` is responsible for shutting down registry-owned resources (for example, the
  database). Background goroutine lifecycle is coordinated by `cmd/server/main.go` via context cancellation.

## Linked References
- `README.md` – quickstart commands, environment setup, and deployment pathways.
- `docs/review.md` – current architectural findings, mitigations, and open follow-ups.

## Recommended Focus Areas
Consider expanding automated test coverage around controllers and critical workflows, and keep the
architectural review document updated as mitigations land.
