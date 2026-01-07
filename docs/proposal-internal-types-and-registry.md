# Proposal: Improve `internal/types` and `internal/registry`

Date: 2026-01-07

## Goals

- Make ownership and lifecycle of runtime dependencies explicit.
- Reduce coupling caused by wide “service locator” interfaces.
- Improve testability by preferring narrow dependency interfaces.
- Align folder/package naming and docs with the actual code.
- Reduce global state (especially caches) to support isolated tests and multi-instance operation.

## Non-goals

- Replace the project with a full DI framework.
- Perform a large rename/refactor across the whole tree in one PR.
- Redesign external store packages (e.g. `github.com/dracory/*store`).

## Terminology

- This proposal uses “composition root” to mean the place where concrete dependencies are constructed and wired together at startup.
- `Registry` is treated as a runtime service container owned by the composition root.
- Passing `RegistryInterface` broadly into many packages is effectively the “service locator” pattern. The goal is to keep that interface at the edges and prefer constructor parameters that only expose what the component needs.

## Current state (findings)

### `internal/types`

- `internal/types` currently contains:
  - unrelated small app types (`flash_message.go`)

`internal/types` is currently small, but it is historically a common place for catch-all growth. The intent is to prevent it from becoming a dumping ground for runtime wiring interfaces or unrelated DTOs.

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
  - it references an `Application` and `types.AppInterface`, while the code uses `Registry` and `registry.RegistryInterface`.

### What the code currently implies (facts)

- `internal/types` is not a composition-root surface; it currently only contains `FlashMessage`.
- `registry.RegistryInterface` is setter-heavy (service locator + mutation).
- `registry.New(cfg)` sets cache values by mutating `internal/cache` package globals; the registry cache accessors currently delegate back to those globals.
- `main.go` already has partial lifecycle handling:
  - DB closure is deferred explicitly (`closeResourcesDB(registry.GetDatabase())`)
  - background work is started under a context and stopped on signals
  - there is not yet a single `registry.Close()` that owns all shutdown.

## Proposal

### 1) Clarify package boundaries and naming

#### Recommendation

- Standardize on canonical packages:
  - configuration lives under `internal/config` (`config.ConfigInterface`, private implementation)
  - registry lives under `internal/registry` (`registry.RegistryInterface`)

- Keep `internal/types` minimal:
  - current state: it contains `FlashMessage`
  - target state: either move remaining app-specific types to a more specific package, or keep `internal/types` as a small shared-types package with a clear boundary

Note: prefer direct imports from `internal/registry` and `internal/config`; do not reintroduce composition-root interfaces into `internal/types`.

#### Expected impact

- Navigation becomes obvious.
- Future changes are localized.

### 2) Make configuration read-mostly (prefer immutable after build)

#### Recommendation

- Introduce a read-only interface (example name): `ConfigReader`.
  - It should only expose `GetX()` methods.
- Keep existing `ConfigInterface` temporarily for compatibility.
- Gradually update consumers (controllers/tasks/schedules) to accept `ConfigReader`.

#### Critical note

If `ConfigReader` still exposes 100+ getters, it will reduce mutation risk but not coupling. Prefer to also introduce feature-focused config views (example: `EmailConfig`, `DatabaseConfig`) where it reduces dependencies.

#### Follow-up

- Prefer a single construction step that validates invariants:

  - `Load() (ConfigInterface, error)`
  - validation errors returned early

#### Expected impact

- Stronger invariants.
- Less runtime mutation.
- Easier test setup (config builder only in tests).

### 3) Replace “wide Registry passed everywhere” with role-based dependency interfaces

Right now, tasks/controllers commonly accept `registry.RegistryInterface`, which allows pulling any dependency at runtime.

#### Recommendation

- Introduce small, role-based interfaces, but do it in a way that does not simply recreate the full registry via interface embedding.
- Prefer feature-focused dependency sets over one-interface-per-store.

- Update tasks/controllers to accept only what they need:
  - For example, a task that needs logging + task store should accept a local interface that embeds only those methods.

- Keep the full `registry.RegistryInterface` limited to startup wiring/bootstrapping and edge integration points.

#### Suggested style

- Define dependency interfaces close to the consumer (same package) unless reuse is proven.
- Avoid “mega-interfaces” like `HasEverything` that become a renamed `RegistryInterface`.
- Prefer constructor injection over pulling dependencies at runtime.

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

#### Transitional approach

- Keep `internal/cache` temporarily, but stop using its globals as the source of truth.
- If something truly needs process-wide caching, make it explicit (a dedicated package with explicit initialization and lifecycle), rather than accidental globals.

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

#### Current gap

`main.go` currently closes the DB and stops background work, but ownership is split across the application entrypoint and the registry. Consolidating shutdown into `Close()` makes it easier to test and harder to forget resources.

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

Prerequisite (already done in this repo): standardize on canonical packages (`internal/config`, `internal/registry`) and keep `internal/types` minimal.

### Phase 1: Introduce narrow interfaces (no behavior change)

- Start with 1-2 tasks and 1 controller to establish a pattern.
- Add dependency interfaces locally to the consuming package.
- Migrate tasks first (they’re typically simple to adjust).

### Phase 2: Config read-only adoption

- Add `ConfigReader`.
- Migrate constructors to accept it.

Optional: introduce smaller config views where it meaningfully reduces coupling.

### Phase 3: Remove global caches

- Move cache instances into `Registry`.
- Stop using `project/internal/cache` global variables as the registry source of truth.
- Delete or deprecate the globals once no longer referenced.

### Phase 4: Clarify ownership of `internal/types` (optional, last)

- Decide whether `FlashMessage` belongs in a more specific package (e.g. a `flash`/`ui`/controller-adjacent package) or whether `internal/types` should be removed entirely.

## Open decision

Choose one:

- **A)** “Registry is the app container”: standardize constructors on narrow dependency interfaces implemented by `Registry`.
- **B)** “Application is the app container”: keep `Registry` as bootstrap/wiring and add an `Application` type for runtime dependencies.

If you pick A, this proposal naturally evolves the current code with minimal renaming.
If you pick B, update `docs/review.md` and reintroduce `Application` as the canonical DI type.

Recommendation: pick **A** for now, because it matches the current entrypoint (`registry.New(cfg)`) and minimizes churn while still enabling interface segregation.
