# Proposal: Improve `internal/types` and `internal/registry`

Date: 2026-01-07

## Goals

- Make ownership and lifecycle of runtime dependencies explicit.
- Reduce coupling caused by wide “service locator” interfaces.
- Improve testability by preferring narrow dependency interfaces.
- Align folder/package naming and docs with the actual code.
- Reduce global state (especially caches) to support isolated tests and multi-instance operation.

## Terminology

- This proposal uses “composition root” to mean the place where concrete dependencies are constructed and wired together at startup.
- `Registry` is treated as a runtime service container owned by the composition root.
- Passing `RegistryInterface` broadly into many packages is effectively the “service locator” pattern. The goal is to keep that interface at the edges and prefer constructor parameters that only expose what the component needs.

## Current state (findings)

### `internal/types`

- `internal/types` currently contains:
  - composition root surface (`RegistryInterface`)
  - unrelated small app types (`flash_message.go`)

This makes the package a catch-all. Over time it becomes harder to reason about responsibilities and changes tend to cause broad rebuild ripple.

#### Configuration is mutable and very wide

- `ConfigInterface` is extremely large and setter-heavy (`SetX`/`GetX` for nearly everything).
- Since configuration can be modified at runtime, it is difficult to enforce invariants and understand who is allowed to mutate what.

### `internal/registry`

#### Registry mixes concerns and depends on globals

`New(cfg)` currently:

- initializes caches via package-level variables (`project/internal/cache`)
- ensures a cache directory exists (but ignores `os.MkdirAll` errors)
- opens the database
- initializes stores
- migrates stores
- switches logger implementation if log store exists

Risks:

- cache globals complicate tests (cross-test contamination)
- error swallowing can hide operational failures
- registry mutability (setter-based interface) makes lifecycle ownership unclear

#### Two-phase store setup is a good direction but not fully enforced

`dataStoresInitialize` and `dataStoresMigrate` are a solid foundation, but the design should ensure:

- initializers do creation/wiring only
- migrators do schema migrations only
- no initializer relies on another initializer’s hidden side effects

### Documentation

- `docs/overview.md` is a good introduction, but it does not define runtime ownership contracts (who closes DB? who owns task runners? are caches global?).
- `docs/review.md` appears partially out of sync with the current code:
  - it references an `Application` and `types.AppInterface`, while the code uses `Registry` and `types.RegistryInterface`.

## Proposal

### 1) Clarify package boundaries and naming

#### Recommendation

- Split responsibilities currently living in `internal/types`:
  - current state: configuration types live under `internal/config` (`config.ConfigInterface`, private implementation)
  - current state: `RegistryInterface` exists under `internal/registry` (`registry.RegistryInterface`)
  - keep `internal/types` only for truly shared domain types (or delete it over time)

Note: `internal/types` no longer provides compatibility aliases to config; consumers should import `internal/config`.

#### Expected impact

- Navigation becomes obvious.
- Future changes are localized.

### 2) Make configuration read-mostly (prefer immutable after build)

#### Recommendation

- Introduce a read-only interface (example name): `ConfigReader`.
  - It should only expose `GetX()` methods.
- Keep existing `ConfigInterface` temporarily for compatibility.
- Gradually update consumers (controllers/tasks/schedules) to accept `ConfigReader`.

#### Follow-up

- Prefer a single construction step that validates invariants:

  - `Load() (ConfigInterface, error)`
  - validation errors returned early

#### Expected impact

- Stronger invariants.
- Less runtime mutation.
- Easier test setup (config builder only in tests).

### 3) Replace “wide Registry passed everywhere” with role-based dependency interfaces

Right now, tasks/controllers commonly accept `types.RegistryInterface`, which allows pulling any dependency at runtime.

#### Recommendation

- Introduce small, role-based interfaces:
  - `HasLogger`
  - `HasConfig`
  - `HasDatabase`
  - `HasUserStore`, `HasTaskStore`, etc.

- Update tasks/controllers to accept only what they need:
  - For example, a task that needs logging + task store should accept an interface that embeds only those.

- Keep the full `registry.RegistryInterface` limited to startup wiring/bootstrapping and edge integration points.

#### Expected impact

- Better modularity.
- Easier unit testing.
- Adding/removing stores does not force pervasive constructor changes.

### 4) Remove global cache singletons (instance-scope cache objects)

#### Current issue

- Caches are stored in `project/internal/cache` package-level variables.

#### Recommendation

- Store cache instances on `Registry` (fields) instead of in a global package.
- `Registry` can still expose `GetMemoryCache()` / `GetFileCache()`, but they should return the registry’s fields.

- Handle filesystem errors:
  - do not ignore `os.MkdirAll`; bubble up errors from `New()`.

#### Expected impact

- Isolated tests.
- Supports running multiple registries in-process.

### 5) Define lifecycle ownership (shutdown/cleanup)

#### Recommendation

- Add an explicit lifecycle method (example): `Close() error` on the registry.
  - closes DB
  - stops background workers/schedulers/task runners (or returns closers owned by those systems)
  - flushes/cleans up other resources as needed

- Document lifecycle contract:
  - what is safe to call after shutdown
  - order of shutdown

#### Expected impact

- Prevents leaks.
- Improves operational predictability.

### 6) Update docs to match code

#### Recommendation

- Update `docs/review.md` to reflect current architecture:
  - change references from `Application`/`types.AppInterface` to `Registry`/`RegistryInterface` if that is the current direction
  - or reintroduce `Application` only if you intend to make it the canonical service container again

- Extend `docs/overview.md` with a short “Architecture Contracts” section:
  - registry owns DB, stores, caches
  - no global singletons (target state)
  - components should depend on narrow interfaces

## Migration plan (incremental)

### Phase 1: Naming and docs alignment (low risk)

- Align folder/package naming (`internal/registry` package name).
- Update docs to match the current code shape.

### Phase 2: Introduce narrow interfaces (no behavior change)

- Add role-based dependency interfaces.
- Migrate tasks first (they’re typically simple to adjust).

### Phase 3: Config read-only adoption

- Add `ConfigReader`.
- Migrate constructors to accept it.

### Phase 4: Remove global caches

- Move cache instances into `Registry`.
- Delete or deprecate `project/internal/cache` global variables.

### Phase 5: Lifecycle closure

- Add `Close()` on registry.
- Ensure `main` calls it during shutdown.

## Open decision

Choose one:

- **A)** “Registry is the app container”: standardize constructors on narrow dependency interfaces implemented by `Registry`.
- **B)** “Application is the app container”: keep `Registry` as bootstrap/wiring and add an `Application` type for runtime dependencies.

If you pick A, this proposal naturally evolves the current code with minimal renaming.
If you pick B, update `docs/review.md` and reintroduce `Application` as the canonical DI type.
