# Proposal: Blueprint Framework Standalone Engine (`dracory/blueprint-framework`)

## Status

Proposed

## Executive Summary

Developers coming to **Blueprint** appreciate its rapid application development (RAD) capabilities and "batteries-included" infrastructure. However, when comparing the file footprint of a new Blueprint project to frameworks like Next.js or Laravel, a freshly generated Blueprint application currently contains **530+ files** across `cmd/`, `internal/app/`, `internal/config/`, `internal/cli/`, `internal/routes/`, `internal/middlewares/`, `internal/layouts/`, and `database/migrations/`.

Much of this file volume consists of structural boilerplate: database pool initialization, 20+ datastore factory bindings, environment validation schemas, signal handling, CLI dispatchers, and background runner groups.

While Next.js achieves a minimal file footprint by encapsulating runtime logic inside `node_modules/next`, **Go applications require explicit composition, strong typing, and direct configurability**. We do not want to obscure Go's explicit nature or reduce Blueprint's powerful configurability (custom DB drivers, store overrides, custom middlewares, CLI flags, background workers, and email templates).

This proposal outlines how to extract Blueprint's core orchestration into an idiomatic, highly configurable Go module — **`github.com/dracory/blueprint-framework`** (or `pkg/framework`). By leveraging **Functional Options**, **Explicit Configuration Structs**, and **Plugin/Adapter Interfaces**, Blueprint projects can reduce boilerplate by **~350 files** while retaining **100% of Blueprint's configuration flexibility**.

---

## The Footprint Problem & Design Philosophy

### Current Blueprint Project (~530+ Files)
A standard Blueprint repository contains extensive infrastructure code living directly in the user's project:
- `cmd/server/*.go` (~4 files): Graceful shutdown, OS signal channel listening, background process context lifecycle.
- `internal/app/*.go` (~10 files): Direct instantiation and getter/setter bindings for 20+ Dracory stores (`userstore`, `sessionstore`, `cmsstore`, `shopstore`, `logstore`, `vaultstore`, `auditstore`, `geostore`, etc.).
- `internal/config/*.go` (~16 files): Environment parsing, database connection pool defaults, LLM API keys, media storage settings.
- `internal/cli/*.go` (~4 files): Maintenance command handlers and CLI argument dispatching.
- `internal/routes/*.go` (~3 files): Global router builder and middleware binding arrays.

### Design Philosophy for `blueprint-framework`
1. **Idiomatic Go**: Use functional options (`framework.WithConfig()`, `framework.WithStore()`, `framework.WithMiddleware()`) and standard interfaces (`net/http`, `slog.Logger`, `sql.DB`).
2. **Configuration Transparency**: All configuration options currently supported in `internal/config` remain available via environment variables, code options, or custom config builders.
3. **Zero Magic**: No hidden magic, code generation, or obscure reflection. The framework acts as an explicit orchestration facade.
4. **Gradual Adoption**: Allow users to rely on default store factories while seamlessly swapping or customizing any store, driver, or middleware.

---

## Architectural Overview

The proposed `blueprint-framework` module exposes a clean, composable Go API while encapsulating standard setup routines:

```
+---------------------------------------------------------------------------------+
|                               User Application                                  |
|                                                                                 |
|  main.go:                                                                       |
|    fw, err := framework.New(                                                    |
|        framework.WithConfig(cfg),                                               |
|        framework.WithCustomStore("custom", myStore),                            |
|        framework.WithRoutes(myAppRoutes),                                       |
|        framework.WithTasks(myTasks),                                            |
|    )                                                                            |
|    fw.Run()                                                                     |
+---------------------------------------------------------------------------------+
                                         |
                                         v
+---------------------------------------------------------------------------------+
|                     github.com/dracory/blueprint-framework                      |
|                                                                                 |
|  +------------------------+  +------------------------+  +-------------------+  |
|  | Config Engine          |  | Datastore Orchestration|  | Router & Server   |  |
|  | - Env Parsing & Vault  |  | - Automatic Store Init |  | - websrv / rtr    |  |
|  | - Connection Pools     |  | - Store Overrides      |  | - Middleware Chain|  |
|  +------------------------+  +------------------------+  +-------------------+  |
|  +------------------------+  +------------------------+  +-------------------+  |
|  | Background Manager     |  | CLI Engine             |  | Logging & Caches  |  |
|  | - Task Queue & Cron    |  | - Maintenance Mode     |  | - slog Handler    |  |
|  | - Cache Expiration     |  | - Command Dispatcher   |  | - Memory & File   |  |
|  +------------------------+  +------------------------+  +-------------------+  |
+---------------------------------------------------------------------------------+
```

---

## Preserving Full Blueprint Configurability

Blueprint is deeply configurable. The standalone module exposes every configuration layer using idiomatic Go constructs:

### 1. Database & Connection Pool Configuration
Users can rely on standard environment variables (`DB_DRIVER`, `DB_HOST`, `DB_MAX_OPEN_CONNS`) or explicitly pass database connections/configs in Go code:

```go
fw, err := framework.New(
    // Option A: Load from environment or custom config
    framework.WithConfig(customCfg),

    // Option B: Provide an existing *sql.DB or custom connection pool
    framework.WithDatabase(customDB),

    // Option C: Multi-database setup (Laravel-style connections supported by neat)
    framework.WithDatabaseConnection("analytics", analyticsDB),
)
```

### 2. Store Overrides & Custom Stores
All 20+ default stores (`UserStore`, `SessionStore`, `VaultStore`, `LogStore`, etc.) are initialized automatically based on configuration, but can be overridden or supplemented:

```go
fw, err := framework.New(
    // Override standard UserStore with a custom implementation
    framework.WithUserStore(myCustomUserStore),

    // Register domain-specific custom store
    framework.WithStore("billing", myBillingStore),
)
```

### 3. Middleware Pipeline Composition
Global and route-level middlewares can be customized, prepended, or appended without modifying framework source code:

```go
fw, err := framework.New(
    // Prepend custom security or logging middleware
    framework.WithGlobalMiddleware(myTelemetryMiddleware),

    // Disable or customize default middlewares (e.g., Maintenance, Security Headers)
    framework.WithMiddlewareOptions(framework.MiddlewareOptions{
        EnableSecurityHeaders: true,
        EnableJailBots:        false,
    }),
)
```

### 4. Background Workers & Job Schedulers
Background task handlers, cron schedules, and maintenance routines remain fully configurable:

```go
fw, err := framework.New(
    framework.WithTaskRegistration(func(app app.AppInterface) {
        tasks.RegisterTasks(app)
        myCustomTasks.Register(app)
    }),
    framework.WithSchedules(myCronScheduleGroup),
    framework.WithBackgroundWorkerCount(10),
)
```

### 5. CLI Command Extension
Custom CLI commands integrate directly into Blueprint's CLI dispatcher (`cmd/server` or custom entry points):

```go
fw, err := framework.New(
    framework.WithCLICommand("seed:demo", "Seeds demo dataset", func(app app.AppInterface, args []string) error {
        return seeders.SeedDemo(app)
    }),
)
```

---

## Code Comparison: Current Blueprint vs. `blueprint-framework`

### `cmd/server/main.go`

#### Current Blueprint (~120 lines of repetitive setup):
```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "project/database/migrations"
    "project/internal/cli"
    "project/internal/app"
    "project/internal/config"
    "project/internal/routes"
    "project/internal/tasks"
    "github.com/dracory/websrv"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    cfg, err := config.NewFromEnv()
    if err != nil {
        fmt.Printf("Failed to load config: %v\n", err)
        return
    }

    app, err := app.New(cfg)
    if err != nil {
        fmt.Printf("Failed to initialize app: %v\n", err)
        return
    }
    defer app.Close()

    if err := migrations.MigrateAll(app); err != nil {
        fmt.Printf("Failed to run migrations: %v\n", err)
        return
    }

    tasks.RegisterTasks(app)

    if isCliMode() {
        if err := cli.ExecuteCliCommand(app, os.Args[1:]); err != nil {
            fmt.Printf("Failed to execute CLI command: %v\n", err)
            os.Exit(1)
        }
        return
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    background := newBackgroundGroup(ctx)
    if err := startBackgroundProcesses(ctx, background, app); err != nil {
        log.Println("Failed to start background processes:", err)
        return
    }

    server, err := websrv.Start(websrv.Options{
        Host:    app.GetConfig().GetAppHost(),
        Port:    app.GetConfig().GetAppPort(),
        URL:     app.GetConfig().GetAppUrl(),
        Handler: routes.Router(app).ServeHTTP,
    })
    if err != nil { ... }

    // Draining background workers and graceful HTTP shutdown logic...
}
```

#### Proposed `blueprint-framework` (~25 lines of idiomatic Go):
```go
package main

import (
    "log"

    "project/database/migrations"
    "project/internal/routes"
    "project/internal/tasks"

    "github.com/dracory/blueprint-framework"
)

func main() {
    engine, err := framework.New(
        framework.WithMigrations(migrations.MigrateAll),
        framework.WithTaskRegistration(tasks.RegisterTasks),
        framework.WithRoutes(routes.AppRoutes),
    )
    if err != nil {
        log.Fatalf("Failed to initialize framework: %v", err)
    }

    if err := engine.Run(); err != nil {
        log.Fatalf("Application runtime error: %v", err)
    }
}
```

---

## Directory Footprint Comparison

| Area | Current Blueprint File Count | Post-Framework File Count | Primary Location After Extraction |
| :--- | :--- | :--- | :--- |
| `cmd/server/` | 7 files | 1 file | Encapsulated in `framework.Engine` |
| `internal/app/` | 10 files | 0 files (or custom extensions) | `framework/app` package |
| `internal/config/` | 16 files | 0 files (or custom additions) | `framework/config` package |
| `internal/cli/` | 4 files | 0 files (or custom commands) | `framework/cli` package |
| `internal/routes/` | 3 files | 1 file (`routes.go`) | Custom app route definitions |
| `internal/middlewares/` | 26 files | ~3 files (app specific) | Standard middlewares move to `framework/middleware` |
| **Total Project Files** | **~530+ files** | **~180 files** | **~350 boilerplate files eliminated** |

*(Note: The remaining ~180 files in a full-stack Blueprint app consist of domain controllers, views, layouts, email templates, and database migrations — pure application logic.)*

---

## Implementation & Transition Strategy

### Phase 1: Internal Framework Extraction (`pkg/framework`)
1. Create `pkg/framework` inside the `dracory/blueprint` repository.
2. Port `internal/app`, `internal/config`, `internal/cli`, and background worker orchestrators into `pkg/framework`.
3. Implement Functional Options (`framework.Option`) for store overrides, DB connections, route registration, and middleware pipeline adjustments.
4. Update `cmd/server/main.go` to use `pkg/framework`.
5. Verify complete test suite (`go test ./...`) passes.

### Phase 2: Standalone Module (`github.com/dracory/blueprint-framework`)
1. Extract `pkg/framework` into its own dedicated repository: `github.com/dracory/blueprint-framework`.
2. Set up automated CI/CD and unit testing for the framework module.
3. Publish version `v1.0.0`.

### Phase 3: Project Template Update & Migration Path
1. Update Blueprint starter template (`github.com/dracory/blueprint`) to use `github.com/dracory/blueprint-framework`.
2. Provide a clear migration guide for existing Blueprint applications to transition to `blueprint-framework` without breaking database models or custom controllers.

---

## Verification & Safety Plan

1. **Type Safety & Go Idioms**: Guarantee all configuration options remain fully typed, documented, and discoverable via standard Go IDE auto-completion.
2. **Backward Compatibility**: Ensure `app.AppInterface` and `config.ConfigInterface` remain untouched so existing controllers, stores, and widgets operate seamlessly.
3. **Test Suite Verification**: Run all unit and integration tests (`go test ./...`) after refactoring `cmd/server/main.go` to ensure 100% test suite pass rate.
