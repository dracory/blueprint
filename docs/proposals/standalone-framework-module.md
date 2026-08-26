# Proposal: Extracting Blueprint Core into a Standalone Framework Module (`dracory/blueprint-framework`)

## Status

Proposed

## Executive Summary

When developers create a project using **Next.js**, they get a minimal directory structure with a single `package.json`, a few pages (`app/page.tsx`, `app/layout.tsx`), and framework logic completely encapsulated inside `next`. In contrast, a freshly cloned **Blueprint** repository currently exposes over **500 files** across dozens of directories (`cmd/`, `internal/app/`, `internal/config/`, `internal/routes/`, `internal/middlewares/`, `internal/layouts/`, `internal/cli/`, `database/migrations/`, etc.).

While Blueprint's "batteries-included" philosophy is powerful, exposing all low-level infrastructure code directly in user projects creates cognitive overload, makes upstream maintenance/upgrades difficult (requiring constant manual migration steps), and forces users to write or maintain extensive boilerplate.

This proposal outlines a strategy to extract Blueprint's core application runtime, configuration engine, background process manager, middleware pipeline, CLI dispatcher, and base layout engine into a standalone Go module: **`github.com/dracory/blueprint-framework`** (or `github.com/dracory/blueprint/pkg/framework`).

---

## The Next.js vs. Blueprint Project Comparison

### Next.js Project Structure (~10-15 files)
```
my-next-app/
├── app/
│   ├── layout.tsx
│   ├── page.tsx
│   └── api/hello/route.ts
├── public/
├── package.json
├── next.config.js
└── tsconfig.json
```
*Framework core code lives entirely inside `node_modules/next`.*

### Current Blueprint Project Structure (~530+ files)
```
my-blueprint-app/
├── cmd/
│   ├── server/           (main.go, background.go, background_processes.go, cli_mode.go)
│   ├── deploy/           (main.go, config.go, functions.go, types.go, constants.go)
│   ├── envenc/           (main.go)
│   ├── loadtest/         (main.go)
│   ├── snakecase/        (main.go)
│   └── ai-browser/       (main.go, background.go, background_processes.go)
├── database/
│   └── migrations/       (48 migration files)
├── internal/
│   ├── app/              (10 orchestration files - datastores, neatdb, logger, cache)
│   ├── config/           (16 config files - app, db, auth, env, stores, etc.)
│   ├── cli/              (4 cli runner files)
│   ├── routes/           (3 router/middleware binding files)
│   ├── middlewares/      (26 middleware files)
│   ├── layouts/          (28 layout/navigation/UI files)
│   ├── tasks/            (23 background task files)
│   ├── schedules/        (6 scheduled job files)
│   ├── controllers/      (123 controller files across admin, auth, website, user, etc.)
│   └── ...
├── taskfile.yml
├── go.mod
└── go.sum
```

---

## Key Problems to Solve

1. **High Cognitive Load**: New developers opening a Blueprint project face 500+ files before writing their first line of business logic.
2. **Upgrade Friction**: Today, updating Blueprint requires manual upgrade guides (e.g. `docs/upgrade_guides/upgrade-v0.37.0-to-v0.38.0.md`) patching low-level `internal/app/`, `internal/config/`, or `cmd/server/` code directly in user repos.
3. **Boilerplate Duplication**: Every Blueprint project duplicates identical `cmd/server/main.go`, signal handling, background process loops, router wiring, and store initialization logic.

---

## Proposed Architecture: `blueprint-framework`

We propose packaging the core orchestration into `github.com/dracory/blueprint-framework` (or `pkg/framework` inside the mono-repo as a first step).

```
+-------------------------------------------------------------------+
|                        User Application                           |
|  - Custom Controllers & Handlers                                  |
|  - Custom Routes & Domain Models                                  |
|  - Custom Migrations & Seeders                                    |
|  - Simple main.go (~15 lines)                                     |
+-------------------------------------------------------------------+
                                  |
                                  v
+-------------------------------------------------------------------+
|               github.com/dracory/blueprint-framework              |
|                                                                   |
|  +--------------------+  +--------------------+                   |
|  | Config Engine      |  | App Runtime &      |                   |
|  | - Env Loader       |  | Store Factory      |                   |
|  | - Vault/Envenc     |  | - NeatDB / SQL     |                   |
|  +--------------------+  +--------------------+                   |
|  +--------------------+  +--------------------+                   |
|  | Router & Server    |  | Background Manager |                   |
|  | - websrv runner    |  | - Task Queue       |                   |
|  | - Default Middle-  |  | - Scheduler        |                   |
|  |   wares Pipeline   |  | - Expirations      |                   |
|  +--------------------+  +--------------------+                   |
|  +--------------------+                                           |
|  | CLI Engine         |                                           |
|  | - Maintenance      |                                           |
|  | - Task/Job runners |                                           |
|  +--------------------+                                           |
+-------------------------------------------------------------------+
```

---

## Modular Breakdown: What Moves vs. What Stays

| Component | Current Location | Framework Module (`blueprint-framework`) | User Application (`my-app`) |
| :--- | :--- | :--- | :--- |
| **App Orchestration** | `internal/app/*` | `framework/app` (Automated initialization of stores, logger, caches) | Imports framework app interface |
| **Config Loader** | `internal/config/*` | `framework/config` (Env parsing, constants, store builders) | Optional app config extensions |
| **Server Engine** | `cmd/server/main.go`, `cmd/server/background*.go` | `framework.Server` / `framework.Run()` | Single entry point in `cmd/server/main.go` |
| **CLI Runner** | `internal/cli/*`, `cmd/server/cli_mode.go` | `framework/cli` (CLI commands dispatcher, maintenance mode) | Registers custom CLI flags/cmds |
| **Middlewares** | `internal/middlewares/*` | `framework/middleware` (Security, Maintenance, API Auth, Logger) | App custom middlewares |
| **Default Router** | `internal/routes/*` | `framework/router` (Base router builder & middleware appender) | Custom route registrations |
| **Admin Adapters** | `internal/controllers/admin/*` | Move base admin interfaces/adapters to `framework/admin` | Custom admin views / options |
| **Domain Logic** | `internal/controllers/*`, `database/*` | Keep out of framework | User controllers, views, DB migrations |

---

## Before & After Comparison

### 1. `cmd/server/main.go`

#### Before (Current Blueprint - ~120 lines):
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
    if err != nil { ... }

    app, err := app.New(cfg)
    if err != nil { ... }
    defer app.Close()

    if err := migrations.MigrateAll(app); err != nil { ... }
    tasks.RegisterTasks(app)

    if isCliMode() {
        if err := cli.ExecuteCliCommand(app, os.Args[1:]); err != nil { ... }
        return
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    background := newBackgroundGroup(ctx)
    if err := startBackgroundProcesses(ctx, background, app); err != nil { ... }

    server, err := websrv.Start(websrv.Options{
        Host:    app.GetConfig().GetAppHost(),
        Port:    app.GetConfig().GetAppPort(),
        URL:     app.GetConfig().GetAppUrl(),
        Handler: routes.Router(app).ServeHTTP,
    })
    // + 30 lines of graceful shutdown and signal handling...
}
```

#### After (With `blueprint-framework` - ~15 lines):
```go
package main

import (
    "log"
    "github.com/dracory/blueprint-framework"
    "project/database/migrations"
    "project/routes"
    "project/tasks"
)

func main() {
    app := framework.NewApp()

    app.OnBoot(func(a framework.App) error {
        if err := migrations.MigrateAll(a); err != nil {
            return err
        }
        tasks.RegisterTasks(a)
        return nil
    })

    app.RegisterRoutes(routes.AppRoutes)

    if err := app.Run(); err != nil {
        log.Fatalf("App error: %v", err)
    }
}
```

---

### 2. User Directory Footprint

After extracting framework boilerplate into `dracory/blueprint-framework`, a target Blueprint project looks as clean as Next.js:

```
my-blueprint-app/
├── cmd/
│   └── server/
│       └── main.go           (~15 lines)
├── controllers/              (User HTTP controllers)
├── database/
│   └── migrations/           (User migrations)
├── models/                   (Domain entities)
├── routes/
│   └── routes.go             (User route registrations)
├── views/                    (Templates / Layouts)
├── .env.example
├── go.mod
└── taskfile.yml
```

**Result**: Redirection of ~350 boilerplate framework files into a versioned dependency module (`go.mod`). Total project file count reduced from **530+ files** to **< 50 files**.

---

## Key Benefits

1. **Next.js-like Ergonomics**: Developers can start writing code instantly in standard directories (`controllers/`, `routes/`, `views/`) without wading through application setup, database pool configuration, signal handlers, or store creation code.
2. **Seamless Upgrades**: Upgrading Blueprint core becomes a single command:
   ```bash
   go get github.com/dracory/blueprint-framework@latest
   ```
   No more copying manual migration guides for `cmd/server/main.go` or `internal/app/app_implementation.go`.
3. **Pluggable & Extensible Hooks**: Provide `OnBoot`, `OnShutdown`, `RegisterMiddlewares`, and `RegisterCLICommands` lifecycle hooks for custom application behavior.
4. **Backward Compatibility**: `blueprint-framework` retains full compatibility with existing `AppInterface`, `ConfigInterface`, and Dracory store packages (`userstore`, `sessionstore`, `rtr`, etc.).

---

## Migration & Implementation Plan

### Phase 1: Package Framework Internals inside Blueprint (`pkg/framework`)
- Create `pkg/framework` in the main Blueprint repository.
- Move `internal/app`, `internal/config`, `internal/cli`, and default middleware setups into `pkg/framework`.
- Refactor `cmd/server/main.go` to use `pkg/framework`.
- Verify full test suite passes.

### Phase 2: Standalone Repository Split (`github.com/dracory/blueprint-framework`)
- Move `pkg/framework` to its own repo: `github.com/dracory/blueprint-framework`.
- Add unit tests & documentation for framework initialization lifecycle.
- Tag initial release `v1.0.0`.

### Phase 3: Update Blueprint Starter Template
- Update the Blueprint starter template repository to depend on `github.com/dracory/blueprint-framework`.
- Simplify directory structure to match the minimal project blueprint.
- Provide upgrade documentation for existing applications wishing to migrate to `blueprint-framework`.

---

## Verification & Testing Strategy

1. **Unit & Integration Tests**: Ensure all store builders, configuration loaders, background runner loops, and route handling pass existing unit tests within `blueprint-framework`.
2. **Benchmark & Boot Time**: Verify that wrapping app boot in `framework.Run()` introduces zero performance or startup overhead.
3. **Template Smoke Test**: Build a sample app with the new minimal structure and test standard flows (auth, migrations, tasks, CLI commands, HTTP routes).
