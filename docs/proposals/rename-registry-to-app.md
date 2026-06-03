# Rename Registry to App

**Status:** Draft

## Problem

The package was originally named `app`, then renamed to `registry`. In practice, `registry` does not feel right for this concept and causes two distinct problems:

**1. Semantic mismatch within the Dracory ecosystem**

Other Dracory modules use `registry` for its classical meaning — a lookup table where things are registered and retrieved by key at runtime:
- `ork` — skill registry (lookup table of runnables)
- `liveflux` — component registry (reflection-based component map)
- `neat` — seeder registry

Blueprint's struct is not a lookup table. It is the application instance, booted once at startup, holding infrastructure dependencies. Calling it `registry` conflates two different patterns.

**2. Poor developer intuition, especially in tests**

In tests, a setup helper spins up the full application context. `testutils.SetupApp()` maps directly to how developers think — "I have an app, I run tests against it." `testutils.SetupRegistry()` is opaque — a registry of what, and why?

Similarly, `app.New(cfg)` and `app.Close()` read as application lifecycle operations. `registry.New(cfg)` and `registry.Close()` do not.

## Proposed Solution

Rename the `registry` package to `app` to better reflect its purpose as the application-level dependency injection container.

**Changes:**
- Rename `internal/registry/` to `internal/app/`
- Update all imports and references throughout the codebase
- Update documentation to reflect the new naming

## Rationale

### Why "app" is better than "registry"

1. **Matches the original intent**: The package was originally named `app` — this rename restores the correct name.
2. **Test readability**: `testutils.SetupApp()` is immediately understood; `testutils.SetupRegistry()` is not.
3. **Lifecycle semantics**: `app.New()`, `app.Close()` read as application lifecycle operations.
4. **Ecosystem clarity**: Frees `registry` to mean only its classical definition (lookup table) across all Dracory modules.
5. **Readability**: `app.GetLogger()`, `app.GetConfig()` read naturally at every call site.

### Why not other alternatives

- **Container**: Standard DI term but more technical/verbose; `container.GetMemoryCache()` is less natural.
- **ServiceLocator**: Too verbose, less Go-idiomatic.
- **ApplicationContext**: Too verbose, Java-centric.
- **services**: Misleading — config and caches are not "services"; `services.GetAppName()` is incoherent.
- **core**: Generic, doesn't convey application instance semantics.
- **deps**: Terse but feels like a shorthand, not a convention.
- **registry**: Conflicts with the classical registry pattern used in `ork`, `liveflux`, and `neat`.

### Addressing naming conflicts

There are existing "app" concepts in the codebase:
- `AppConfigInterface` in `config` package
- `appSettings` struct in `config` package

These are not practical conflicts because:
- Different packages (`config` vs `app`)
- Different purposes (configuration vs runtime services)
- Context makes usage clear in practice
- Different import paths prevent confusion

### Ecosystem distinction

This rename clarifies a meaningful distinction across the Dracory ecosystem:

| Pattern | Name | Meaning |
|---|---|---|
| Dracory modules (`ork`, `liveflux`, `neat`) | `registry` | Lookup table — things registered and retrieved by key |
| Blueprint applications | `app` | Application instance — booted once, holds infrastructure dependencies |

Using `app` for Blueprint and reserving `registry` for lookup-table patterns gives every Dracory contributor a clear, consistent mental model.

## Detailed Changes

### 1. Package Rename

```
internal/registry/
├── doc.go
├── registry_interface.go
└── registry_implementation.go
```

Becomes:

```
internal/app/
├── doc.go
├── app_interface.go
└── app_implementation.go
```

### 2. File Renames

- `registry_interface.go` → `app_interface.go`
- `registry_implementation.go` → `app_implementation.go`
- Update package declaration from `package registry` to `package app`

### 3. Interface Rename

- `RegistryInterface` → `AppInterface`
- `registryImplementation` → `appImplementation`

### 4. Function Renames

- `registry.New()` → `app.New()`
- All getter/setter methods remain the same (no changes needed)

### 5. Import Updates

Update all files that import the registry package:

```go
// Before
import "project/internal/registry"

// After
import "project/internal/app"
```

### 6. Documentation Updates

- Update `docs/proposals/architecture-documentation.md`
- Update `README.md` if it references registry
- Update code comments and docstrings

## Impact Assessment

**Breaking Changes:** Yes - This is a breaking change for external consumers if any exist.

**Risk:** Low - This is purely a mechanical rename with no functional changes.

**Testing:** All existing tests should pass after the rename. No new tests needed.

**Migration Effort:** Medium - Requires updating all imports and references throughout the codebase.

## Migration Steps

1. Create new `internal/app/` directory structure
2. Copy and rename files from `internal/registry/` to `internal/app/`
3. Update package declarations and type names in new files
4. Update all imports throughout the codebase
5. Update documentation
6. Run tests to verify no regressions
7. Delete old `internal/registry/` directory
8. Commit changes

## Benefits

1. **Improved clarity**: Code is more self-documenting
2. **Better onboarding**: New developers understand the purpose faster
3. **Industry alignment**: Follows common naming patterns
4. **Reduced confusion**: Less ambiguity about what the package does
5. **Better readability**: `app.GetLogger()` reads naturally

## Alternatives Considered

### Keep "registry"

**Pros:**
- No migration effort
- Already in use

**Cons:**
- Conflicts semantically with the classical registry pattern used in `ork`, `liveflux`, and `neat`
- Obscures test setup intent (`SetupRegistry` vs `SetupApp`)
- Was already replaced once — the original name was `app`
- Lifecycle methods (`Close()`, `New()`) read awkwardly as registry operations

### Use "Container"

**Pros:**
- Standard DI term
- Clear purpose
- No conflicts

**Cons:**
- More technical/verbose
- Less Go-idiomatic
- `container.GetMemoryCache()` is less natural than `app.GetMemoryCache()`

### Use "ServiceLocator"

**Pros:**
- Explicit about purpose
- No conflicts

**Cons:**
- Too verbose
- Less Go-idiomatic
- Not commonly used in Go

## Related Proposals

- Architecture Documentation (proposals/architecture-documentation.md)
