# Project Overview

## Brand Concept
Blueprint is a ready-to-adapt Go web application starter that gives teams a reliable, opinionated
foundation for MVC projects. It highlights production-ready patterns while leaving room for product-
specific customization.

## Mission
Provide engineering teams with a maintainable baseline that accelerates delivery of secure, observable,
and resilient web applications without locking them into a rigid framework.

## Product Metaphor
Blueprint behaves like an architect's master plan: it lays out structural systems (routing, data access,
background work, configuration) so feature teams can focus on bespoke rooms and finishes.

## Primary Personas
- **Founding engineer** – needs a dependable scaffold to launch customer features rapidly while keeping
  operational risk low.
- **Platform maintainer** – curates shared infrastructure and documentation to ensure consistency across
  multiple deployments and teams.
- **Contributor or contractor** – joins midstream and relies on clear guides, tests, and scripts to become
  productive within a single working session.

## Application Summary
- Root entrypoint `main.go` wires configuration, datastore initialization, routing, background tasks, and
  graceful shutdown handling.
- `/internal/app` manages application services, including database access, caching, stores, and lifecycle
  helpers for background processes.
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

## Linked References
- `README.md` – quickstart commands, environment setup, and deployment pathways.
- `docs/review.md` – current architectural findings, mitigations, and open follow-ups.

## Recommended Focus Areas
Consider expanding automated test coverage around controllers and critical workflows, and keep the
architectural review document updated as mitigations land.
