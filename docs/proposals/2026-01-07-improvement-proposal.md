# System Review & Improvement Proposal

**Date:** January 7, 2026
**Author:** Golang Principal Engineer Agent (GoGod)
**Status:** Draft

## 1. Executive Summary

This document outlines a systematic review of the `dracory.com/blueprint` project. The current codebase demonstrates a mature, opinionated architecture with a strong focus on modularity and dependency injection via a central Registry. It leverages modern Go (1.25) and a robust internal ecosystem (`dracory/*`).

The primary goal of this proposal is to refine the existing architecture to improve testability, reduce coupling, and enhance operational excellence without rewriting the core logic. Key recommendations include narrowing dependency interfaces, standardizing project layout, and hardening the deployment pipeline.

## 2. Architecture & Design Review

### 2.1. The Registry Pattern ("God Object")
**Current State:** The `RegistryInterface` serves as a central service locator, exposing over 25 different stores and services. Controllers and tasks currently depend on this entire interface.
**Analysis:** While this simplifies wiring in `main.go`, it violates the **Interface Segregation Principle**. A controller that only needs `UserStore` has access to (and technically depends on) `ShopStore`, `BlogStore`, etc. This makes unit testing difficult as mocks become complex.
**Recommendation:**
*   **Refactor Consumers:** Modify controllers and services to accept specific interfaces (e.g., `NewHomeController(userStore userstore.StoreInterface)` instead of `NewHomeController(reg registry.RegistryInterface)`).
*   **Retain Registry for Startup:** Keep the Registry for the composition root (`main.go`) to manage lifecycles, but do not pass it deep into the application.

### 2.2. Web Server Abstraction
**Current State:** Controllers like `homeController` return `string` directly. The `websrv` package likely handles the `net/http` wrapping.
**Analysis:** This abstraction simplifies handlers but may limit flexibility regarding HTTP status codes, headers, and content types (e.g., returning JSON vs HTML).
**Recommendation:** Ensure `websrv` supports a `Result` or `Response` type that allows handlers to control status codes and headers while maintaining the convenience of the current abstraction.

## 3. Code Structure & Organization

### 3.1. Project Layout
**Current State:** `main.go` resides in the project root.
**Analysis:** Standard Go practice (per `golang-standards/project-layout`) encourages placing entry points in `cmd/`. This keeps the root directory clean and supports multiple binaries easily.
**Recommendation:** Move `main.go` to `cmd/server/main.go`. Update `Dockerfile` and `taskfile.yml` accordingly.

### 3.2. Configuration Management
**Current State:** `ConfigInterface` is extremely comprehensive and granular.
**Analysis:** The granular getters/setters are good for mocking but create a lot of boilerplate.
**Recommendation:** Consider grouping related configs (e.g., `DatabaseConfig`, `AuthConfig`) into structs or smaller interfaces to reduce verbosity while maintaining strict typing.

## 4. Operational Excellence

### 4.1. Docker Optimization
**Current State:**
```dockerfile
COPY . ./
RUN ... go build ...
```
**Analysis:** Copying the entire source before building invalidates the build cache whenever *any* file changes (including docs or README).
**Recommendation:**
1.  `COPY go.mod go.sum ./` -> `RUN go mod download` (Already done - Good!).
2.  `COPY . ./` -> Use `.dockerignore` effectively (already present).
3.  Consider building specific packages if possible, but the current approach is generally acceptable if `.dockerignore` is strict.

### 4.2. Observability
**Current State:** `slog` is used.
**Recommendation:**
*   **Tracing:** Integrate OpenTelemetry (OTEL) to trace requests across the seemingly micro-service-like store architecture.
*   **Metrics:** Ensure key metrics (HTTP latency, error rates, store operation times) are exposed (e.g., via Prometheus).

## 5. Testing Strategy

**Current State:** Tests exist (`main_test.go`, `cli_test.go`).
**Recommendation:**
*   **Table-Driven Tests:** Enforce table-driven tests for all business logic in `internal/`.
*   **Interface Mocks:** With the architectural shift proposed in 2.1, generating mocks for specific stores will become easier, leading to better unit test coverage.
*   **Integration Tests:** Create a dedicated integration test suite (e.g., `tests/integration`) that spins up the full registry with a test database.

## 6. Action Plan

### Phase 1: Foundation (Immediate)
1.  **Refactor Directory:** Move `main.go` to `cmd/server/main.go`.
2.  **Linting:** Ensure `golangci-lint` is configured with strict rules (including `cyclop`, `gocognit`, `funlen`) to maintain code quality.

### Phase 2: Decoupling (Medium Term)
3.  **Refactor Controllers:** Pick one controller (e.g., `HomeController`) and refactor it to accept only its direct dependencies. Measure the impact on testability.
4.  **Apply to All:** Roll out the dependency injection refactor across `internal/controllers` and `internal/tasks`.

### Phase 3: Observability & Scale (Long Term)
5.  **Instrumentation:** Add OpenTelemetry traces to `internal/routes` and `internal/registry`.
6.  **Load Testing:** Implement k6 or similar load tests to benchmark the system limits.

---
**Signed:**
*GoGod (Golang Principal Engineer Agent)*
