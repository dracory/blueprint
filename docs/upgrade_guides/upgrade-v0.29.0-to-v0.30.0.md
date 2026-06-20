# Upgrade Guide: v0.29.0 to v0.30.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.29.0 to v0.30.0.

## Overview

This release consolidates database connection management and migrations under the `dracory/neat` package, introduces a Laravel-style multi-connection configuration model, and aligns the taskstore/statsstore APIs with a consistent `Get*` getter convention. The terminology used across the project is also updated from "framework" to "rapid application development (RAD) starter template".

**Key Changes:**
- Database connection layer migrated from `github.com/dracory/database` to `github.com/dracory/neat/database`
- Migration engine migrated from `github.com/dracory/migrate` to `github.com/dracory/neat/database/migrator`
- `AppInterface` gained neat database accessors (`GetNeatDatabase`, `SetNeatDatabase`, `GetDatabaseConnection`)
- New optional `.env` keys: `DB_DEFAULT_CONNECTION`, `DB_DSN`, `DB_PREFIX`
- New `DatabaseConnectionConfigInterface` and extended `DatabaseConfigInterface`
- Migration interface changed: `ID()` → `Signature()`, `Up(ctx, tx)` → `Up()`, `Down(ctx, tx)` → `Down()`, `CreatedAt()` removed
- Custom SQL migrations now use the neat schema builder instead of `sb` + raw SQL
- Subscription store migration switched from `AutoMigrate` to `MigrateUp`/`MigrateDown`
- taskstore deprecated non-`Get` aliases removed from entity interfaces (`ID()`, `Status()`, `Details()`, `TaskID()`, `Attempts()`, `CompletedAt()`, `CreatedAt()`, `Output()`, `Parameters()`, `SoftDeletedAt()`, `StartedAt()`, `UpdatedAt()`, `QueueName()`, `Alias()`, `Description()`, `Memo()`, `IsRecurring()`, `RecurrenceRule()`, `Title()`, `Name()`, `TaskDefinitionID()`); `Data()`, `DataChanged()`, `MarkAsNotDirty()` removed from entity interfaces
- statsstore query API switched from `VisitorQueryOptions` struct to `VisitorQueryInterface` builder (`NewVisitorQuery()` / `VisitorQuery()`)
- statsstore `VisitorInterface` getters renamed to `Get*` convention (`UserAgent()` → `GetUserAgent()`, `IpAddress()` → `GetIpAddress()`, `Country()` → `GetCountry()`, `ID()` → `GetID()`, `Path()` → `GetPath()`, `Fingerprint()` → `GetFingerprint()`, `CreatedAt()` → `GetCreatedAt()`, `UpdatedAt()` → `GetUpdatedAt()`, `DeletedAt()` → `GetSoftDeletedAt()`, and all `User*()` getters); `Data()`/`DataChanged()`/`MarkAsNotDirty()` removed; `IsSoftDeleted()` added
- logstore `LogCount` return type changed from `(int, error)` to `(int64, error)`
- geostore `StateCreate` and `StatesCreate` now require a `context.Context` first parameter
- versionstore `VersionInterface` getters renamed (`CreatedAt()` → `GetCreatedAt()`, `SoftDeletedAt()` → `GetSoftDeletedAt()`); `Data()`/`DataChanged()`/`MarkAsNotDirty()` removed; `IsSoftDeleted()`, `GetCreatedAtCarbon()`, `GetSoftDeletedAtCarbon()` added
- DataObject methods (`Data()`, `DataChanged()`, `MarkAsNotDirty()`) removed from entity interfaces across sessionstore, chatstore, feedstore, settingstore, subscriptionstore, and vaultstore
- cachestore `DriverName(db *sql.DB) string` removed from `StoreInterface`
- Terminology updated from "framework" to "rapid application development (RAD) starter template"
- Removed unused `pkg/blogai/database.go`

---

## ⚠️ Breaking Changes

### 1. Database Connection Layer Migrated to `dracory/neat`

**Change**: The application no longer opens database connections through `github.com/dracory/database`. Connection management is now handled by `github.com/dracory/neat/database`, which provides multi-connection support, pooling, and a Laravel-style configuration model. The `*sql.DB` handle is derived from the neat instance so existing stores continue to receive a standard `*sql.DB`.

`internal/app/database_open.go` now returns `*neatdatabase.Database` instead of `*sql.DB`, and `app.New()` derives the `*sql.DB` via `neatDB.DB()`.

**Old Usage**:
```go
// internal/app/database_open.go (v0.29.0)
import (
    "database/sql"
    "github.com/dracory/database"
)

func databaseOpen(cfg config.ConfigInterface) (*sql.DB, error) {
    options := database.Options().
        SetDatabaseType(cfg.GetDatabaseDriver()).
        SetDatabaseHost(cfg.GetDatabaseHost()).
        // ...
    return database.Open(options)
}
```

**New Usage**:
```go
// internal/app/database_open.go (v0.30.0)
import (
    neatdatabase "github.com/dracory/neat/database"
)

func databaseOpen(cfg config.ConfigInterface) (*neatdatabase.Database, error) {
    neatCfg := config.DatabaseNeatConfig(cfg)
    return neatdatabase.New(neatCfg)
}
```

**Action Required**:
- If you have custom code that calls `databaseOpen()` directly, update it to handle the `*neatdatabase.Database` return type.
- If you imported `github.com/dracory/database` directly, replace it with `github.com/dracory/neat/database` (or remove it if you only used it through `app`).
- Run `go mod tidy` after updating imports to drop `github.com/dracory/database` from your module graph.

**Files to Check**:
- `internal/app/database_open.go` (already updated in template)
- `internal/app/app_implementation.go` (already updated in template)
- Any custom code that imported `github.com/dracory/database`

---

### 2. AppInterface Gained Neat Database Accessors

**Change**: `AppInterface` now exposes the neat database instance and a named-connection accessor. The existing `GetDatabase()` / `SetDatabase()` methods remain unchanged for backward compatibility.

**Old Usage**:
```go
type AppInterface interface {
    // ...
    GetDatabase() *sql.DB
    SetDatabase(db *sql.DB)
    // ...
}
```

**New Usage**:
```go
import neatdatabase "github.com/dracory/neat/database"

type AppInterface interface {
    // ...
    GetDatabase() *sql.DB
    SetDatabase(db *sql.DB)

    GetNeatDatabase() *neatdatabase.Database
    SetNeatDatabase(db *neatdatabase.Database)

    GetDatabaseConnection(name string) *sql.DB
    // ...
}
```

**Action Required**:
- No action required for applications using the standard `GetDatabase()` accessor.
- If you implement `AppInterface` with a custom type, add the three new methods (`GetNeatDatabase`, `SetNeatDatabase`, `GetDatabaseConnection`).
- If you embed `appImplementation`, you inherit these methods automatically.

**Close semantics**: `app.Close()` now closes the neat database instance (which closes the underlying `*sql.DB`). Do not call `app.GetDatabase().Close()` separately — this would result in a double close.

---

### 3. New Database Configuration Interface and `.env` Keys

**Change**: The config layer now models database settings as a map of named connections. A new `DatabaseConnectionConfigInterface` and extended `DatabaseConfigInterface` are introduced. The existing single-database getters (`GetDatabaseDriver`, `GetDatabaseHost`, etc.) continue to work and delegate to the default connection.

**New `.env` keys** (all optional):
```env
# Default connection name (default: "default")
DB_DEFAULT_CONNECTION=default

# Direct DSN override (optional)
DB_DSN=

# Table prefix (optional)
DB_PREFIX=
```

**New interface methods**:
```go
type DatabaseConfigInterface interface {
    // ... existing single-db getters/setters ...

    SetDatabaseDSN(string)
    GetDatabaseDSN() string

    SetDatabasePrefix(string)
    GetDatabasePrefix() string

    SetDatabaseDefaultConnection(string)
    GetDatabaseDefaultConnection() string

    GetDatabaseConnections() []DatabaseConnectionConfigInterface
    GetDatabaseConnectionByName(name string) DatabaseConnectionConfigInterface
}

type DatabaseConnectionConfigInterface interface {
    GetName() string
    GetDriver() string
    GetHost() string
    GetPort() string
    GetDatabase() string
    GetUsername() string
    GetPassword() string
    GetSSLMode() string
    GetCharset() string
    GetTimezone() string
    GetDSN() string
    GetPrefix() string
}
```

**Action Required**:
- Add the new optional keys to your `.env` only if you need multi-connection support, DSN overrides, or a table prefix. Existing `.env` files continue to work unchanged.
- If you implement `ConfigInterface` with a custom type, add the new getters/setters. If you embed the default `configImplementation`, you inherit them automatically.
- A new helper `config.DatabaseNeatConfig(cfg)` maps the blueprint config to a `neat/database/db.DBConfig`. Use it if you build a neat database instance outside of `app.New()`.

**Files to Check**:
- `internal/config/z_config_interfaces.go` (already updated in template)
- `internal/config/z_config_constants.go` (already updated in template)
- `internal/config/database_config.go` (already updated in template)
- `internal/config/z_config_implementation.go` (already updated in template)
- `internal/config/database_neat_config.go` (new file in template)
- `.env.example` (already updated in template)

---

### 4. Migration Engine Migrated to `dracory/neat/database/migrator`

**Change**: The migration system no longer uses `github.com/dracory/migrate`. All migrations now implement `migrator.MigrationInterface` from `github.com/dracory/neat/database/migrator`. Store migrations and custom SQL migrations are combined and tracked in a single `migration_tracker` table.

The `database/migrations/registry.go` file has been removed; `getStoreMigrations` and `getSQLMigrations` now live in `database/migrations/migrate.go`.

**Migration interface changes**:

| v0.29.0 (`migrate.MigrationInterface`) | v0.30.0 (`migrator.MigrationInterface`) |
|----------------------------------------|----------------------------------------|
| `ID() string`                          | `Signature() string`                   |
| `Up(ctx context.Context, tx *sql.Tx) error` | `Up() error`                     |
| `Down(ctx context.Context, tx *sql.Tx) error` | `Down() error`                   |
| `CreatedAt() time.Time`                | removed (provided by `BaseMigration`)  |

Migrations must now embed `migrator.BaseMigration` to satisfy the interface.

**Old Usage**:
```go
import (
    "context"
    "database/sql"
    "time"

    "github.com/dracory/migrate"
    "github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreAuditMigrate)(nil)

type StoreAuditMigrate struct {
    app app.AppInterface
}

func (m *StoreAuditMigrate) ID() string {
    return "2026_03_21_0001_store_audit_migrate"
}

func (m *StoreAuditMigrate) Up(ctx context.Context, tx *sql.Tx) error {
    return m.app.GetAuditStore().MigrateUp(ctx)
}

func (m *StoreAuditMigrate) Down(ctx context.Context, tx *sql.Tx) error {
    return m.app.GetAuditStore().MigrateDown(ctx, tx)
}

func (m *StoreAuditMigrate) CreatedAt() time.Time {
    return carbon.Parse("2026-03-21 00:01:00", "UTC").StdTime()
}
```

**New Usage**:
```go
import (
    "context"
    "errors"

    "github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreAuditMigrate)(nil)

type StoreAuditMigrate struct {
    migrator.BaseMigration
    app app.AppInterface
}

func (m *StoreAuditMigrate) Signature() string {
    return "2026_03_21_0001_store_audit_migrate"
}

func (m *StoreAuditMigrate) Up() error {
    if m.app == nil {
        return errors.New("app is nil")
    }
    store := m.app.GetAuditStore()
    if store == nil {
        return errors.New("audit store is not initialized")
    }
    return store.MigrateUp(context.Background())
}

func (m *StoreAuditMigrate) Down() error {
    store := m.app.GetAuditStore()
    if store == nil {
        return errors.New("audit store is not initialized")
    }
    return store.MigrateDown(context.Background())
}
```

**Action Required**:
- For each custom migration file in `database/migrations/`:
  1. Change the import from `github.com/dracory/migrate` to `github.com/dracory/neat/database/migrator`.
  2. Add `migrator.BaseMigration` as an embedded field in the migration struct.
  3. Rename `ID()` to `Signature()`.
  4. Change `Up(ctx context.Context, tx *sql.Tx) error` to `Up() error`.
  5. Change `Down(ctx context.Context, tx *sql.Tx) error` to `Down() error`.
  6. Remove the `CreatedAt() time.Time` method and the `github.com/dromara/carbon/v2` import (no longer needed).
  7. Update the interface assertion from `migrate.MigrationInterface` to `migrator.MigrationInterface`.
- Update any code that called `migrator.AddMigration(m)` / `migrator.Up(ctx)` directly — the new API uses `m.AddMigrations([]MigrationInterface)` then `m.Up(ctx)` on a `*migrator.Migrator` instance.
- Remove references to the deleted `database/migrations/registry.go` (its helpers moved into `migrate.go`).

**Migration Command** (Unix/macOS):
```bash
# Update imports
find . -type f -name "*.go" -exec sed -i 's|"github.com/dracory/migrate"|"github.com/dracory/neat/database/migrator"|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|migrate\.MigrationInterface|migrator.MigrationInterface|g' {} \;

# Rename ID() to Signature()
find . -type f -name "*.go" -path "*/migrations/*" -exec sed -i 's|func (m \*\([A-Za-z]*\)) ID() string|func (m *\1) Signature() string|g' {} \;

# Rename Up/Down signatures (review carefully before applying)
find . -type f -name "*.go" -path "*/migrations/*" -exec sed -i 's|func (m \*\([A-Za-z]*\)) Up(ctx context.Context, tx \*sql.Tx) error|func (m *\1) Up() error|g' {} \;
find . -type f -name "*.go" -path "*/migrations/*" -exec sed -i 's|func (m \*\([A-Za-z]*\)) Down(ctx context.Context, tx \*sql.Tx) error|func (m *\1) Down() error|g' {} \;
```

> Note: The `sed` commands above are starting points. Review each migration file manually afterward to remove the `CreatedAt()` method, the `carbon` import, and to embed `migrator.BaseMigration`.

**Files to Check**:
- All files under `database/migrations/` (template files already updated)
- Any custom migration files you have added
- `database/migrations/migrate.go` (already updated in template)

---

### 5. Custom SQL Migrations Use the Neat Schema Builder

**Change**: The example custom SQL migration (`2026_03_22_0001_table_custom_create.go`) has been rewritten to use the neat schema builder (`m.GetSchema().Create(...)`) instead of the `github.com/dracory/sb` builder and raw `tx.Exec` SQL. The `database/sql` and `context` imports are no longer needed in custom migrations.

**Old Usage**:
```go
import (
    "context"
    "database/sql"
    "github.com/dracory/database"
    "github.com/dracory/sb"
)

func (m *TableCustomCreate) Up(ctx context.Context, tx *sql.Tx) error {
    dialect := database.DatabaseType(tx)
    tableCreateSql, err := sb.NewBuilder(dialect).
        Table("custom_example").
        Column(sb.Column{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true}).
        // ...
        Create()
    if err != nil {
        return err
    }
    _, err = tx.Exec(tableCreateSql)
    return err
}
```

**New Usage**:
```go
import (
    contractsschema "github.com/dracory/neat/contracts/database/schema"
    "github.com/dracory/neat/database/migrator"
)

func (m *TableCustomCreate) Up() error {
    if m.GetSchema().HasTable("custom_example") {
        return nil
    }
    return m.GetSchema().Create("custom_example", func(blueprint contractsschema.Blueprint) {
        blueprint.ID()
        blueprint.String("name")
        blueprint.String("email")
        blueprint.Unique("email")
        blueprint.String("status")
        blueprint.Timestamps()
    })
}

func (m *TableCustomCreate) Down() error {
    return m.GetSchema().DropIfExists("custom_example")
}
```

**Action Required**:
- If you have custom SQL migrations using `sb` + `tx.Exec`, consider migrating them to the neat schema builder for cross-database compatibility.
- Existing raw-SQL migrations can still be adapted by calling `m.GetSchema().Exec()` or by keeping raw SQL through the neat connection, but the builder approach is now the recommended pattern.
- The `TableCustomCreate` migration is commented out by default in `getSQLMigrations()` — uncomment it to use as a template.

---

### 6. Subscription Store Migration: AutoMigrate → MigrateUp/MigrateDown

**Change**: The subscription store migration no longer calls `store.AutoMigrate(ctx)`. It now uses `store.MigrateUp(context.Background())` and implements a real `Down()` via `store.MigrateDown(context.Background())`.

**Old Usage**:
```go
func (m *StoreSubscriptionMigrate) Up(ctx context.Context, tx *sql.Tx) error {
    return store.AutoMigrate(ctx)
}

func (m *StoreSubscriptionMigrate) Down(ctx context.Context, tx *sql.Tx) error {
    // AutoMigrate doesn't support rollback
    return nil
}
```

**New Usage**:
```go
func (m *StoreSubscriptionMigrate) Up() error {
    return store.MigrateUp(context.Background())
}

func (m *StoreSubscriptionMigrate) Down() error {
    return store.MigrateDown(context.Background())
}
```

**Action Required**:
- If you have custom code calling `subscriptionStore.AutoMigrate(ctx)`, switch to `MigrateUp(ctx)` / `MigrateDown(ctx)`.
- The template migration file has already been updated.

---

### 7. taskstore Deprecated Aliases Removed and DataObject Methods Dropped

**Change**: The `taskstore` package (bumped to v1.25.0) removed the deprecated non-`Get` accessor aliases from the entity interfaces (`TaskQueueInterface`, `TaskDefinitionInterface`, `ScheduleInterface`). In v1.22.0 these were present with `// Deprecated: use Get* instead. Will be removed after 2026-11-30.` comments; they have now been removed ahead of schedule. The `Data()`, `DataChanged()`, and `MarkAsNotDirty()` DataObject methods were also removed from the entity interfaces.

> Note: The **query** interfaces (`TaskQueueQueryInterface`, `TaskDefinitionQueryInterface`, `ScheduleQueryInterface`) keep their non-`Get` accessor names (`ID()`, `Status()`, `TaskID()`, `QueueName()`, etc.) — only the **entity** interfaces changed.

**Removed from `TaskQueueInterface`** (use the `Get*` equivalent):

| Removed alias | Replacement |
|---------------|-------------|
| `ID()` | `GetID()` |
| `Status()` | `GetStatus()` |
| `Details()` | `GetDetails()` |
| `TaskID()` | `GetTaskID()` |
| `Attempts()` | `GetAttempts()` |
| `CompletedAt()` | `GetCompletedAt()` |
| `CreatedAt()` | `GetCreatedAt()` |
| `Output()` | `GetOutput()` |
| `Parameters()` | `GetParameters()` |
| `SoftDeletedAt()` | `GetSoftDeletedAt()` |
| `StartedAt()` | `GetStartedAt()` |
| `UpdatedAt()` | `GetUpdatedAt()` |
| `QueueName()` | `GetQueueName()` |
| `Data()` | removed (no replacement) |
| `DataChanged()` | removed (no replacement) |
| `MarkAsNotDirty()` | removed (no replacement) |

**Removed from `TaskDefinitionInterface`** (use the `Get*` equivalent):

| Removed alias | Replacement |
|---------------|-------------|
| `ID()` | `GetID()` |
| `Status()` | `GetStatus()` |
| `Alias()` | `GetAlias()` |
| `Description()` | `GetDescription()` |
| `Title()` | `GetTitle()` |
| `Memo()` | `GetMemo()` |
| `IsRecurring()` | `GetIsRecurring()` |
| `RecurrenceRule()` | `GetRecurrenceRule()` |
| `CreatedAt()` | `GetCreatedAt()` |
| `UpdatedAt()` | `GetUpdatedAt()` |
| `SoftDeletedAt()` | `GetSoftDeletedAt()` |
| `Data()` | removed (no replacement) |
| `DataChanged()` | removed (no replacement) |
| `MarkAsNotDirty()` | removed (no replacement) |

**Removed from `ScheduleInterface`** (use the `Get*` equivalent):

| Removed alias | Replacement |
|---------------|-------------|
| `ID()` | `GetID()` |
| `Name()` | `GetName()` |
| `Status()` | `GetStatus()` |
| `QueueName()` | `GetQueueName()` |
| `TaskDefinitionID()` | `GetTaskDefinitionID()` |

**Old Usage**:
```go
if queuedTask.Status() == taskstore.TaskQueueStatusRunning {
    // ...
}
err := t.app.GetTaskStore().TaskQueueDeleteByID(ctx, purgeTask.ID())
details := task.QueuedTask().Details()
```

**New Usage**:
```go
if queuedTask.GetStatus() == taskstore.TaskQueueStatusRunning {
    // ...
}
err := t.app.GetTaskStore().TaskQueueDeleteByID(ctx, purgeTask.GetID())
details := task.QueuedTask().GetDetails()
```

**Action Required**:
- Update all calls to the removed non-`Get` aliases on taskstore entity types to their `Get*` equivalents.
- Remove any calls to `Data()`, `DataChanged()`, or `MarkAsNotDirty()` on taskstore entities. If you relied on `Data()` for introspection, use the typed getters instead.
- This affects task handlers, CLI commands (`internal/cmds/execute_job.go`), scheduled jobs (`internal/schedules/`), and tests.
- If you implement `TaskQueueInterface`, `TaskDefinitionInterface`, or `ScheduleInterface` with a custom type, remove the deprecated alias methods and the DataObject methods from your implementation.

**Migration Command** (Unix/macOS):
```bash
# Rename taskstore entity getters (review results before committing)
find . -type f -name "*.go" -exec sed -i 's|\.ID()|.GetID()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.Status()|.GetStatus()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.Details()|.GetDetails()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.TaskID()|.GetTaskID()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.Attempts()|.GetAttempts()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.CompletedAt()|.GetCompletedAt()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.CreatedAt()|.GetCreatedAt()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.Output()|.GetOutput()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.Parameters()|.GetParameters()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.StartedAt()|.GetStartedAt()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.UpdatedAt()|.GetUpdatedAt()|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|\.QueueName()|.GetQueueName()|g' {} \;
```

> Warning: These broad `sed` replacements may also match unrelated types that legitimately have `ID()`, `Status()`, etc. methods (including taskstore **query** interfaces, which keep their non-`Get` names). Review the diff carefully and revert any incorrect changes (for example, on `app.AppInterface`, store types, query builders, or model types that did not rename their getters).

**Files to Check**:
- `internal/cmds/execute_job.go` (already updated in template)
- `internal/schedules/schedule_queue_clear_job.go` (already updated in template)
- `internal/tasks/clean_up/clean_up_task.go` (already updated in template)
- `internal/tasks/blind_index_rebuild/blind_index_rebuild_task_test.go` (already updated in template)
- Any custom task handlers or schedules
- Any custom implementations of `TaskQueueInterface`, `TaskDefinitionInterface`, or `ScheduleInterface`

---

### 8. statsstore Query API and Visitor Interface Overhaul

**Change**: The `statsstore` package (bumped to v1.2.0) made two breaking changes:

1. **Query type changed**: `VisitorList` and `VisitorCount` now accept `VisitorQueryInterface` (a fluent builder) instead of the `VisitorQueryOptions` struct. Use `statsstore.NewVisitorQuery()` or `statsstore.VisitorQuery()` to build the query.

2. **`VisitorInterface` getters renamed**: All non-`Get` getters were removed and replaced with `Get*` equivalents. `DeletedAt()` was renamed to `GetSoftDeletedAt()`. The `Data()`, `DataChanged()`, and `MarkAsNotDirty()` DataObject methods were removed. `IsSoftDeleted()` was added.

**Removed from `VisitorInterface`** (use the `Get*` equivalent):

| Removed | Replacement |
|---------|-------------|
| `ID()` | `GetID()` |
| `IpAddress()` | `GetIpAddress()` |
| `UserAgent()` | `GetUserAgent()` |
| `Country()` | `GetCountry()` |
| `Path()` | `GetPath()` |
| `Fingerprint()` | `GetFingerprint()` |
| `CreatedAt()` | `GetCreatedAt()` |
| `UpdatedAt()` | `GetUpdatedAt()` |
| `DeletedAt()` | `GetSoftDeletedAt()` |
| `UserAcceptLanguage()` | `GetUserAcceptLanguage()` |
| `UserAcceptEncoding()` | `GetUserAcceptEncoding()` |
| `UserBrowser()` | `GetUserBrowser()` |
| `UserBrowserVersion()` | `GetUserBrowserVersion()` |
| `UserDevice()` | `GetUserDevice()` |
| `UserDeviceType()` | `GetUserDeviceType()` |
| `UserOs()` | `GetUserOs()` |
| `UserOsVersion()` | `GetUserOsVersion()` |
| `UserReferrer()` | `GetUserReferrer()` |
| `Data()` | removed (no replacement) |
| `DataChanged()` | removed (no replacement) |
| `MarkAsNotDirty()` | removed (no replacement) |

**Old Usage**:
```go
entries, err := store.VisitorList(ctx, statsstore.VisitorQueryOptions{
    Country:      "empty",
    Limit:        10,
    CreatedAtGte: date + " 00:00:00",
    CreatedAtLte: date + " 23:59:59",
    Distinct:     statsstore.COLUMN_IP_ADDRESS,
})

ua := useragent.Parse(visitor.UserAgent())
country := t.findCountryByIp(ctx, visitor.IpAddress())
```

**New Usage**:
```go
query := statsstore.NewVisitorQuery().
    SetCountry("empty").
    SetLimit(10).
    SetCreatedAtGte(date + " 00:00:00").
    SetCreatedAtLte(date + " 23:59:59").
    SetDistinct(statsstore.COLUMN_IP_ADDRESS)
entries, err := store.VisitorList(ctx, query)

ua := useragent.Parse(visitor.GetUserAgent())
country := t.findCountryByIp(ctx, visitor.GetIpAddress())
```

**Action Required**:
- Replace all `statsstore.VisitorQueryOptions{...}` struct literals with `statsstore.NewVisitorQuery().Set...()` (or `statsstore.VisitorQuery().Set...()`) chains.
- Rename every `visitor.X()` getter call to its `visitor.GetX()` equivalent (see the table above). Note that `DeletedAt()` is now `GetSoftDeletedAt()`, not `GetDeletedAt()`.
- Remove any calls to `Data()`, `DataChanged()`, or `MarkAsNotDirty()` on visitor entities.
- If you implement `VisitorInterface` with a custom type, remove the old non-`Get` methods and the DataObject methods, and add the `Get*` methods plus `IsSoftDeleted() bool`.

**Files to Check**:
- `internal/controllers/admin/home_controller.go` (already updated in template)
- `internal/tasks/stats/stats_visitor_enhance_task.go` (already updated in template)
- Any custom code querying the stats store or implementing `VisitorInterface`

---

### 9. logstore `LogCount` Return Type Changed

**Change**: The `logstore` package (bumped to v1.18.0) changed the return type of `LogCount` from `(int, error)` to `(int64, error)`.

**Old Usage**:
```go
total, err := logStore.LogCount(ctx, query)
// total is int
var total int = 0
if n, err := logStore.LogCount(ctx, query); err == nil {
    total = n // n is int
}
```

**New Usage**:
```go
total, err := logStore.LogCount(ctx, query)
// total is int64
var total int64 = 0
if n, err := logStore.LogCount(ctx, query); err == nil {
    total = n // n is int64
}
```

**Action Required**:
- If you assign the result of `LogCount` to an explicitly typed `int` variable, change it to `int64` (or use `:=` type inference).
- If you pass the result to a function expecting `int`, add an explicit cast: `int(n)`.
- The template's `pkg/logadmin/log_manager/log_list_results.go` has been updated (`logListResult.Total` changed from `int` to `int64`).

**Files to Check**:
- `pkg/logadmin/log_manager/log_list_results.go` (already updated in template)
- `pkg/logadmin/log_manager/log_manager_controller.go` (uses `:=` type inference, no change needed)
- Any custom code calling `logStore.LogCount(...)`

---

### 10. geostore `StateCreate` and `StatesCreate` Now Require `context.Context`

**Change**: The `geostore` package (bumped to v1.5.0) added a required `context.Context` first parameter to `StateCreate` and `StatesCreate`.

**Old Usage**:
```go
err := store.StateCreate(state)
err := store.StatesCreate(states)
```

**New Usage**:
```go
err := store.StateCreate(ctx, state)
err := store.StatesCreate(ctx, states)
```

**Action Required**:
- Add a `context.Context` first argument to all calls to `StateCreate` and `StatesCreate`.
- The Blueprint template does not call these methods directly (they are used internally by `store.Seed(ctx)`), so no template files needed updating.
- If you implement the geostore `StoreInterface` with a custom type, update your `StateCreate` and `StatesCreate` method signatures.

**Files to Check**:
- Any custom code calling `geostore` `StateCreate` or `StatesCreate`
- Any custom implementations of the geostore `StoreInterface`

---

### 11. versionstore `VersionInterface` Getters Renamed and DataObject Methods Dropped

**Change**: The `versionstore` package (bumped from v1.1.0 to v1.5.0) renamed `VersionInterface` getters to the `Get*` convention and removed the DataObject methods.

**Removed from `VersionInterface`** (use the `Get*` equivalent):

| Removed | Replacement |
|---------|-------------|
| `CreatedAt()` | `GetCreatedAt()` |
| `SoftDeletedAt()` | `GetSoftDeletedAt()` |
| `Data()` | removed (no replacement) |
| `DataChanged()` | removed (no replacement) |
| `MarkAsNotDirty()` | removed (no replacement) |

**Added to `VersionInterface`**:
- `GetCreatedAtCarbon() *carbon.Carbon`
- `GetSoftDeletedAtCarbon() *carbon.Carbon`
- `IsSoftDeleted() bool`

**Old Usage**:
```go
createdAt := version.CreatedAt()
softDeletedAt := version.SoftDeletedAt()
data := version.Data()
```

**New Usage**:
```go
createdAt := version.GetCreatedAt()
softDeletedAt := version.GetSoftDeletedAt()
// Data() removed — use typed getters instead
```

**Action Required**:
- Rename `version.CreatedAt()` → `version.GetCreatedAt()` and `version.SoftDeletedAt()` → `version.GetSoftDeletedAt()`.
- Remove any calls to `Data()`, `DataChanged()`, or `MarkAsNotDirty()` on version entities.
- If you implement `VersionInterface` with a custom type, remove the old methods and add the new `Get*` methods plus `IsSoftDeleted()`, `GetCreatedAtCarbon()`, and `GetSoftDeletedAtCarbon()`.

**Files to Check**:
- The Blueprint template only uses `versionstore.COLUMN_CREATED_AT` (a constant, unchanged) — no template files needed updating.
- Any custom code using `versionstore.VersionInterface` getters or implementing the interface.

---

### 12. DataObject Methods Removed Across Multiple Store Packages

**Change**: Several store packages removed the DataObject-style methods (`Data()`, `DataChanged()`, `MarkAsNotDirty()`) from their entity interfaces. This is part of a codebase-wide cleanup to replace generic data introspection with typed getters. The taskstore (#7), statsstore (#8), and versionstore (#11) removals are already documented above. The following additional stores are also affected:

| Package | Interface(s) | Removed Methods |
|---------|-------------|-----------------|
| `sessionstore` v1.11.0 → v1.15.0 | `SessionInterface` | `Data()`, `DataChanged()`, `MarkAsNotDirty(...string)` |
| `chatstore` v0.10.0 → v1.1.0 | `ChatInterface` | `Data()`, `DataChanged()` (kept `MarkAsNotDirty()`) |
| `feedstore` v0.9.0 → v1.1.0 | `FeedInterface`, `LinkInterface` | `DataChanged()` (kept `Data()` and `MarkAsNotDirty()`) |
| `settingstore` v1.7.0 → v1.9.0 | `SettingInterface` | `Data()`, `DataChanged()`, `MarkAsNotDirty(...string)` |
| `subscriptionstore` v1.0.0 → v1.2.0 | `PlanInterface`, `SubscriptionInterface` | `Data()`, `DataChanged()`, `MarkAsNotDirty(columns ...string)` |
| `vaultstore` v1.0.0 → v1.2.0 | `RecordInterface`, `MetaInterface` | `Data()`, `DataChanged()` |

> Note: `chatstore` `MessageInterface` also **added** `ID()`, `SetID()`, and `MarkAsNotDirty()` as new required methods. If you implement `MessageInterface` with a custom type, you must add these methods.

**Action Required**:
- Remove any calls to `Data()`, `DataChanged()`, or `MarkAsNotDirty()` on entities from the listed stores. If you relied on `Data()` for introspection, use the typed getters (`GetID()`, `GetStatus()`, `GetCreatedAt()`, etc.) instead.
- If you implement any of these interfaces with a custom type, remove the DataObject methods from your implementation.
- The Blueprint template does not call these methods on the affected stores, so no template files needed updating.

**Files to Check**:
- Any custom code using `Data()`, `DataChanged()`, or `MarkAsNotDirty()` on entities from sessionstore, chatstore, feedstore, settingstore, subscriptionstore, or vaultstore
- Any custom implementations of these store interfaces

---

### 13. cachestore `DriverName` Removed from `StoreInterface`

**Change**: The `cachestore` package (bumped from v1.2.0 to v1.6.0) removed the `DriverName(db *sql.DB) string` method from `StoreInterface`.

**Old Usage**:
```go
driver := cacheStore.DriverName(db)
```

**New Usage**:
```go
// DriverName removed — use database/sql or dracory/neat utilities to determine the driver
```

**Action Required**:
- Remove any calls to `cacheStore.DriverName(db)`. Use standard `database/sql` driver detection or `dracory/neat` utilities instead.
- If you implement `cachestore.StoreInterface` with a custom type, remove the `DriverName` method from your implementation.
- The Blueprint template does not call this method — no template files needed updating.

---

### 14. Terminology: "Framework" → "RAD Starter Template"

**Change**: Documentation and code comments now refer to Blueprint as a "rapid application development (RAD) starter template" instead of a "framework". The version constant comment in `internal/config/version.go` and the README have been updated.

**Action Required**:
- No code changes required. If you maintain documentation that references Blueprint as a "framework", consider updating the wording for consistency.

---

### 15. Removed `pkg/blogai/database.go`

**Change**: The file `pkg/blogai/database.go` has been removed. It was already fully commented out and unused.

**Action Required**:
- None, unless you had uncommented or referenced this file locally. Remove any local references.

---

## 🔄 Migration Steps

### Step 1: Update Dependencies

Update your `go.mod` to the new dependency versions. The key changes are:

- Add `github.com/dracory/neat v0.23.0`
- Remove `github.com/dracory/database` (direct) and `github.com/dracory/migrate`
- Bump store packages: `auditstore v1.5.0`, `blindindexstore v1.12.0`, `blogstore v1.25.0`, `cachestore v1.6.0`, `chatstore v1.1.0`, `cmsstore v1.33.0`, `customstore v1.10.0`, `entitystore v1.10.0`, `feedstore v1.1.0`, `geostore v1.5.0`, `logstore v1.18.0`, `metastore v1.7.0`, `sessionstore v1.15.0`, `settingstore v1.9.0`, `shopstore v1.17.0`, `statsstore v1.2.0`, `subscriptionstore v1.2.0`, `taskstore v1.25.0`, `userstore v1.15.0`, `vaultstore v1.2.0`, `versionstore v1.5.0`, `email v0.2.0`

```bash
go get github.com/dracory/neat@v0.23.0
go get github.com/dracory/taskstore@v1.25.0
go get github.com/dracory/statsstore@v1.2.0
go get github.com/dracory/sessionstore@v1.15.0
go get github.com/dracory/vaultstore@v1.2.0
# ... and the other bumped stores listed above
go mod tidy
```

### Step 2: Update Database Connection Layer

If you have custom code that imported `github.com/dracory/database`, replace it with `github.com/dracory/neat/database`. The template's `internal/app/database_open.go` and `internal/app/app_implementation.go` have already been updated.

```bash
# Update direct database imports (if any)
find . -type f -name "*.go" -exec sed -i 's|"github.com/dracory/database"|neatdatabase "github.com/dracory/neat/database"|g' {} \;
```

Review the results manually — the API differs (`database.Open(options)` vs `neatdatabase.New(dbConfig)`), so calls will need rewriting rather than a simple import swap.

### Step 3: Update Migration Files

For each migration file under `database/migrations/`:

1. Swap the import from `github.com/dracory/migrate` to `github.com/dracory/neat/database/migrator`.
2. Embed `migrator.BaseMigration` in the struct.
3. Rename `ID()` → `Signature()`.
4. Change `Up(ctx, tx)` → `Up()` and `Down(ctx, tx)` → `Down()`.
5. Remove `CreatedAt()` and the `carbon` import.
6. Update the interface assertion to `migrator.MigrationInterface`.

The template migration files have already been updated — copy the pattern from `2026_03_21_0001_store_audit_migrate.go`.

### Step 4: Update Custom SQL Migrations (Optional)

If you have custom SQL migrations using `sb` + `tx.Exec`, migrate them to the neat schema builder. See `2026_03_22_0001_table_custom_create.go` for the new pattern (`m.GetSchema().Create(...)` / `m.GetSchema().DropIfExists(...)`).

### Step 5: Update taskstore Getter Calls

Rename all taskstore entity getter calls to the `Get*` convention, and remove any `Data()`/`DataChanged()`/`MarkAsNotDirty()` calls. See Breaking Change #7 for the full list. Do **not** rename getters on taskstore **query** interfaces — those keep their non-`Get` names.

### Step 6: Update statsstore Query and Visitor Calls

Replace `VisitorQueryOptions{...}` struct literals with `NewVisitorQuery().Set...()` (or `VisitorQuery().Set...()`) builders, and rename all `visitor.X()` getter calls to `visitor.GetX()` (see Breaking Change #8 for the full list, including `DeletedAt()` → `GetSoftDeletedAt()`).

### Step 7: Update logstore `LogCount` Consumers

If you assign the result of `logStore.LogCount(...)` to an explicitly typed `int` variable, change it to `int64`. See Breaking Change #9.

### Step 8: Update geostore `StateCreate`/`StatesCreate` Calls

If you call `store.StateCreate(state)` or `store.StatesCreate(states)` directly, add a `context.Context` first argument. See Breaking Change #10.

### Step 9: Update versionstore `VersionInterface` Calls

Rename `version.CreatedAt()` → `version.GetCreatedAt()` and `version.SoftDeletedAt()` → `version.GetSoftDeletedAt()`. Remove any `Data()`/`DataChanged()`/`MarkAsNotDirty()` calls. See Breaking Change #11.

### Step 10: Remove DataObject Method Calls from Other Stores

Remove any calls to `Data()`, `DataChanged()`, or `MarkAsNotDirty()` on entities from sessionstore, chatstore, feedstore, settingstore, subscriptionstore, and vaultstore. See Breaking Change #12 for the full list.

### Step 11: Remove cachestore `DriverName` Calls

Remove any calls to `cacheStore.DriverName(db)`. See Breaking Change #13.

### Step 12: Add Optional `.env` Keys (Optional)

If you need multi-connection support, a DSN override, or a table prefix, add the new keys to your `.env`:

```env
DB_DEFAULT_CONNECTION=default
DB_DSN=
DB_PREFIX=
```

Existing `.env` files continue to work without these keys.

### Step 13: Build and Test

Build the application to ensure all changes are compatible:

```bash
go build ./...
```

Run tests:

```bash
go test ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests

Run all unit tests to ensure no regressions:

```bash
go test ./...
```

### 2. Integration Tests

Run integration tests if applicable:

```bash
go test -tags=integration ./...
```

### 3. Build Verification

Verify the application builds successfully:

```bash
go build -o ./bin/server ./cmd/server
```

### 4. Database Migration Verification

Verify migrations run successfully with the new neat migrator:

```bash
# Start the server (migrations run automatically on boot)
go run ./cmd/server
```

Check that:
- The `migration_tracker` table is created and populated.
- All enabled store tables are created.
- No duplicate migration entries from the old `migrate` framework remain (the tracker table name may differ — if you see stale records from the previous tracker, clean them up manually).

### 5. Manual Testing

Test the application manually:

```bash
# Start the server
go run ./cmd/server

# Test key functionality:
# - Application starts successfully
# - Database connections work (SQLite, PostgreSQL, MySQL)
# - Migrations complete without errors
# - Admin interfaces load
# - Task queue processes jobs
# - Stats visitor enhancement task runs
```

---

## 📝 Additional Notes

### Rationale for the neat Migration

The project previously used two separate packages for database concerns: `github.com/dracory/database` for opening connections and `github.com/dracory/neat` for migrations and ORM operations. The `neat` package already contains a full database connection layer with multi-connection support, read/write replicas, pooling, and a Laravel-style configuration model. Consolidating under `neat` removes a redundant dependency and unlocks neat's advanced features. See `docs/proposals/change-to-neat.md` for the full proposal.

### Backward Compatibility

- All existing `.env` keys remain unchanged and continue to work.
- Existing single-database getters (`GetDatabaseDriver`, `GetDatabaseHost`, etc.) delegate to the default connection.
- Stores continue to receive `*sql.DB`, so no store constructors needed to change.
- The default connection name is `"default"` when `DB_DEFAULT_CONNECTION` is unset.

### Close Semantics

`app.Close()` now closes the neat database instance, which in turn closes the underlying `*sql.DB`. Do not call `app.GetDatabase().Close()` separately. The `appImplementation.Close()` method guards against double-close by niling out `neatDB` after closing.

### Dependency Cleanup

This release removes several transitive dependencies that were no longer needed after dropping `dracory/database`, `dracory/migrate`, and the GORM stack, including `gorm.io/*`, `jackc/pgx`, `doug-martin/goqu`, `jinzhu/*`, and `mattn/go-sqlite3`. Duplicate OpenTelemetry dependencies were also cleaned up (now consistently at v1.44.0).

### New Tests

- `internal/config/database_neat_config_test.go` — tests for the config-to-neat mapper.
- `internal/app/app_close_test.go` — new `TestClose_ClosesNeatDatabase` test verifying neat DB lifecycle on close.
- `pkg/useradmin/user_update/handle_user_update_ajax_test.go`, `handle_user_fetch_ajax_test.go`, `handle_timezones_fetch_ajax_test.go`, `user_update_page_test.go` — expanded AJAX handler tests for the user update controller (replaced `t.Skip()` placeholders with real test implementations using `testutils.Setup`).

---

## 🆘 Common Issues and Solutions

### Issue 1: "undefined: database" after migration

**Symptom**: Compilation errors about undefined `database` package.

**Solution**: Replace `github.com/dracory/database` imports with `github.com/dracory/neat/database`. Note the API differs — `database.Open(options)` becomes `neatdatabase.New(dbConfig)`. Use `config.DatabaseNeatConfig(cfg)` to build the config.

### Issue 2: "undefined: migrate.MigrationInterface"

**Symptom**: Compilation errors in migration files.

**Solution**: Update the import to `github.com/dracory/neat/database/migrator` and the type to `migrator.MigrationInterface`. Embed `migrator.BaseMigration` in your struct.

### Issue 3: Migration method signature mismatch

**Symptom**: Errors like `cannot use X (type *Foo) as type migrator.MigrationInterface in assignment: wrong type for Up method`.

**Solution**: Update method signatures to `Up() error` and `Down() error` (no `ctx`/`tx` parameters), rename `ID()` to `Signature()`, and remove `CreatedAt()`.

### Issue 4: "undefined: registry" in migrations

**Symptom**: Compilation errors referencing the deleted `database/migrations/registry.go`.

**Solution**: The `getStoreMigrations` and `getSQLMigrations` helpers moved into `database/migrations/migrate.go`. Remove any local references to the deleted file. The `validateMigrationID` helper was removed with the old `migrate` package.

### Issue 5: taskstore method not found (`ID`, `Status`, `Details`, etc.)

**Symptom**: Errors like `queuedTask.ID undefined (type *TaskQueue has no field or method ID)`, or `taskDefinition.Alias undefined`, or `schedule.Name undefined`.

**Solution**: The deprecated non-`Get` aliases were removed from the entity interfaces. Rename to the `Get*` getters: `GetID()`, `GetStatus()`, `GetDetails()`, `GetTaskID()`, `GetAttempts()`, `GetCompletedAt()`, `GetCreatedAt()`, `GetOutput()`, `GetParameters()`, `GetSoftDeletedAt()`, `GetStartedAt()`, `GetUpdatedAt()`, `GetQueueName()`, `GetAlias()`, `GetDescription()`, `GetTitle()`, `GetMemo()`, `GetIsRecurring()`, `GetRecurrenceRule()`, `GetName()`, `GetTaskDefinitionID()`. See Breaking Change #7 for the full table. Note that taskstore **query** interfaces keep their non-`Get` names — do not rename those.

### Issue 6: taskstore `Data` / `DataChanged` / `MarkAsNotDirty` undefined

**Symptom**: Errors like `queuedTask.Data undefined (type *TaskQueue has no field or method Data)`.

**Solution**: The DataObject methods were removed from the taskstore entity interfaces. Remove these calls from your code. If you relied on `Data()` for introspection, use the typed getters (`GetID()`, `GetStatus()`, etc.) instead.

### Issue 7: statsstore `VisitorQueryOptions` undefined

**Symptom**: Errors about `statsstore.VisitorQueryOptions` not being a valid type or wrong argument type for `VisitorList`.

**Solution**: Switch to the builder: `statsstore.NewVisitorQuery().SetCountry(...).SetLimit(...)` (or `statsstore.VisitorQuery().Set...()`) and pass the builder result to `VisitorList` / `VisitorCount`.

### Issue 8: statsstore visitor getters undefined (`UserAgent`, `IpAddress`, `Country`, etc.)

**Symptom**: Errors like `visitor.UserAgent undefined` or `visitor.DeletedAt undefined`.

**Solution**: All `VisitorInterface` getters were renamed to the `Get*` convention. Use `GetUserAgent()`, `GetIpAddress()`, `GetCountry()`, `GetID()`, `GetPath()`, `GetFingerprint()`, `GetCreatedAt()`, `GetUpdatedAt()`, `GetSoftDeletedAt()` (note: `DeletedAt()` became `GetSoftDeletedAt()`, not `GetDeletedAt()`), and the `GetUser*()` variants. See Breaking Change #8 for the full table.

### Issue 9: logstore `LogCount` type mismatch

**Symptom**: Errors like `cannot use n (type int64) as type int in assignment` or `cannot assign int64 to total (type int)`.

**Solution**: `logstore.LogCount` now returns `(int64, error)`. Change the receiving variable from `int` to `int64`, or add an explicit cast: `int(n)`. See Breaking Change #9.

### Issue 10: geostore `StateCreate` / `StatesCreate` argument mismatch

**Symptom**: Errors like `not enough arguments in call to store.StateCreate` or `cannot use state (type *State) as type context.Context`.

**Solution**: Both methods now require a `context.Context` first argument: `store.StateCreate(ctx, state)` and `store.StatesCreate(ctx, states)`. See Breaking Change #10.

### Issue 11: versionstore `VersionInterface` getters undefined (`CreatedAt`, `SoftDeletedAt`)

**Symptom**: Errors like `version.CreatedAt undefined` or `version.SoftDeletedAt undefined`.

**Solution**: The getters were renamed: `version.CreatedAt()` → `version.GetCreatedAt()` and `version.SoftDeletedAt()` → `version.GetSoftDeletedAt()`. The `Data()`/`DataChanged()`/`MarkAsNotDirty()` methods were also removed. See Breaking Change #11.

### Issue 12: `Data` / `DataChanged` / `MarkAsNotDirty` undefined on store entities

**Symptom**: Errors like `session.Data undefined` or `setting.DataChanged undefined` or `subscription.MarkAsNotDirty undefined` on entities from sessionstore, chatstore, feedstore, settingstore, subscriptionstore, or vaultstore.

**Solution**: The DataObject methods were removed from entity interfaces across these stores. Remove these calls from your code. If you relied on `Data()` for introspection, use the typed getters (`GetID()`, `GetCreatedAt()`, etc.) instead. See Breaking Change #12 for the full list of affected stores and interfaces.

### Issue 13: cachestore `DriverName` undefined

**Symptom**: Errors like `cacheStore.DriverName undefined (type *Store has no field or method DriverName)`.

**Solution**: The `DriverName(db *sql.DB) string` method was removed from `cachestore.StoreInterface`. Remove calls to this method and use standard `database/sql` driver detection or `dracory/neat` utilities instead. See Breaking Change #13.

### Issue 14: Double close panic on shutdown

**Symptom**: `sql: database is closed` or a panic during `app.Close()`.

**Solution**: Do not call `app.GetDatabase().Close()` in your shutdown code. `app.Close()` now closes the neat database instance which closes the underlying `*sql.DB`.

### Issue 15: Stale migration tracker entries

**Symptom**: Migrations appear to not run after upgrading because they were already recorded.

**Solution**: The new migrator uses a `migration_tracker` table. If the previous `migrate` framework used a different tracker table name, old entries will not be reused. Either drop the old tracker table or manually mark migrations as pending in the new tracker. Since signatures are unchanged (`ID()` values became `Signature()` values with the same string), you can copy over the recorded migration IDs if needed.

---

## 📞 Support

For issues or questions about this upgrade:
- Check the [Blueprint repository](https://github.com/dracory/blueprint)
- Review the [neat migration proposal](docs/proposals/change-to-neat.md) for detailed rationale
- Open an issue on GitHub for upgrade-specific problems

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
- [x] File naming follows pattern: `upgrade-vX.Y.Z-to-vX.Y.Z.md`
- [x] Emoji styling used consistently (⚠️, 🔄, 🧪, 📝, 🆘)
- [x] Git tag verified for previous version (v0.29.0)
- [x] Previous guides reviewed for consistency
