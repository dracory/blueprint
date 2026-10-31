# Upgrade Guide: v0.55.0 to v0.56.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.55.0 to v0.56.0.

## Overview

This is primarily a **dependency maintenance release** plus one performance change: Blueprint now writes log records and visitor stats asynchronously so request handlers no longer pay a synchronous DB INSERT per request.

**Key Changes:**
- **Async log + visitor writes**: `logstore.NewSlogAsyncHandler` replaces the synchronous `NewSlogHandler`, and `statsstore` gains `VisitorCreateAsyncEnabled`. Both are drained in `app.Close()` before the DB shuts down (see "Async writes" below)
- `github.com/dracory/logstore` v1.23.2 → v1.24.0 (adds `SlogAsyncHandler`)
- `github.com/dracory/statsstore` v1.34.2 → v1.35.0 (adds `VisitorCreateAsyncEnabled`, `DrainVisitorWrites`, `VisitorDroppedCount`)
- `github.com/dracory/auth` v0.41.0 → v0.44.0 (local `replace` directive removed)
- `github.com/dracory/rtr` v1.10.1 → v1.11.0
- `github.com/dracory/neat` v0.58.0 → v0.61.0
- `github.com/dracory/shopadmin` v0.8.0 → v0.9.1
- `github.com/dracory/statsadmin` v0.8.0 → v0.9.2
- `github.com/dracory/logadmin` v0.2.0 → v0.2.1
- `github.com/dracory/blogadmin` v0.3.0 → v0.3.2
- `github.com/dracory/cmsstore` v1.43.0 → v1.44.2
- `github.com/dracory/userstore` v1.21.0 → v1.22.2
- All other `dracory/*` store packages bumped to latest patch versions (see table below)
- `github.com/yuin/goldmark` v1.8.6 → v1.8.7
- `modernc.org/sqlite` v1.61.0 → v1.61.1
- Indirect bumps: `cloud.google.com/go/*`, `google.golang.org/api` v0.309.0 → v0.310.0, `google.golang.org/genai` v1.76.0 → v1.77.0, `github.com/aws/aws-sdk-go-v2/credentials` v1.20.12 → v1.20.13

**One small application code change is recommended** — wiring the async handlers in `app.New`/`app.Close` and enabling `VisitorCreateAsyncEnabled` in `NewStatsStore` (see below). If your project has local `replace` directives or pinned older versions, also follow the auth removal step.

---

## ⚠️ Breaking Changes

**None.** All bumps are minor/patch releases that remain backward compatible with the Blueprint codebase. The full test suite passes unchanged.

### 1. Local `replace` directive for dracory/auth removed (conditional)

**Change**: Blueprint no longer carries `replace github.com/dracory/auth => ../auth` (or similar local path) in `go.mod`. The auth package is resolved from the Go module proxy at v0.44.0.

**Action Required** (only if your project copied the replace directive):
- Remove `replace github.com/dracory/auth => ...` from your `go.mod`
- Run `go mod tidy`
- Ensure you are on auth v0.44.0 — earlier code written against v0.41.0 APIs still compiles, but the package's `middlewares.RememberMeMiddleware` and delivery-key contract fixes are only in the newer releases

---

## ⚡ Async Log and Visitor Writes

Every request previously paid two synchronous DB writes: `LogRequestMiddleware` → logstore INSERT, and `StatsMiddleware` → `statsstore.VisitorCreate` INSERT. On busy sites (especially with crawler traffic) this adds real per-request latency. logstore v1.24.0 and statsstore v1.35.0 add opt-in async variants that buffer writes on a background worker and drain them at shutdown.

### 1. logstore: `NewSlogHandler` → `NewSlogAsyncHandler`

**Old** (`internal/app/app_implementation.go`):

```go
if app.IsEnabledLogStore() {
	app.SetLogger(slog.New(logstore.NewSlogHandler(app.GetLogStore())))
}
```

**New**:

```go
// appImplementation struct gains:
asyncLogHandler *logstore.SlogAsyncHandler

// in New():
if app.IsEnabledLogStore() {
	// The logstore handler does a synchronous DB INSERT per record —
	// use the async variant so logging never blocks the request path
	// on a DB write. Drained in Close before the DB shuts down.
	app.asyncLogHandler = logstore.NewSlogAsyncHandler(app.GetLogStore())
	app.SetLogger(slog.New(app.asyncLogHandler))
}
```

The raw handler is stored on the struct because `*slog.Logger` doesn't expose `Drain()` — it is needed at shutdown, nothing else uses it.

### 2. statsstore: enable `VisitorCreateAsyncEnabled`

In `internal/config/store_builders.go` → `NewStatsStore`:

```go
st, err := statsstore.NewStore(statsstore.NewStoreOptions{
	DB:                db,
	SettingsTableName: "snv_stats_settings",
	VisitorTableName:  "snv_stats_visitor",
	// VisitorCreate runs on every request via StatsMiddleware — queue
	// the INSERT instead of paying synchronous DB write latency per
	// request. Flushed by DrainVisitorWrites in app.Close.
	VisitorCreateAsyncEnabled: true,
})
```

### 3. Drain both in `app.Close()` before the DB closes

```go
func (r *appImplementation) Close() error {
	if r == nil {
		return nil
	}

	// Flush buffered log records and visitor writes before the database
	// closes so pending writes don't fail against a closed connection.
	if r.asyncLogHandler != nil {
		r.asyncLogHandler.Drain()
		r.asyncLogHandler = nil
	}
	if r.statsStore != nil {
		r.statsStore.DrainVisitorWrites()
	}

	if r.neatDB == nil {
		return nil
	}
	// ...
}
```

### 4. Tests that assert visitor rows need a flush

With async enabled, `VisitorCreate` returns before the row is inserted. Any test that calls `VisitorCreate` (or hits a page under `StatsMiddleware`) and then asserts on visitor rows must call `statsStore.DrainVisitorWrites()` first — otherwise assertions race the background worker.

```go
app.GetStatsStore().DrainVisitorWrites() // flush before asserting rows
```

Both queues are bounded and drop-on-full; `asyncLogHandler.DroppedCount()` and `statsStore.VisitorDroppedCount()` expose drop counts if you want to monitor them.

---

## 📦 Dependency Changes

### dracory/* packages

| Package | Old | New |
|---|---|---|
| dracory/auditstore | v1.12.0 | v1.12.2 |
| dracory/auth | v0.41.0 | v0.44.0 |
| dracory/blindindexstore | v1.19.0 | v1.19.2 |
| dracory/blogadmin | v0.3.0 | v0.3.2 |
| dracory/blogstore | v1.39.0 | v1.39.2 |
| dracory/cachestore | v1.10.0 | v1.10.2 |
| dracory/chatstore | v1.6.0 | v1.6.2 |
| dracory/cmsstore | v1.43.0 | v1.44.2 |
| dracory/customstore | v1.16.0 | v1.16.2 |
| dracory/entitystore | v1.23.0 | v1.23.2 |
| dracory/feedstore | v1.12.0 | v1.12.2 |
| dracory/geostore | v1.14.0 | v1.14.2 |
| dracory/logadmin | v0.2.0 | v0.2.1 |
| dracory/logstore | v1.23.0 | v1.24.0 |
| dracory/metastore | v1.12.0 | v1.12.2 |
| dracory/neat | v0.58.0 | v0.61.0 |
| dracory/rtr | v1.10.1 | v1.11.0 |
| dracory/sessionstore | v1.20.0 | v1.20.2 |
| dracory/settingstore | v1.14.0 | v1.14.2 |
| dracory/shopadmin | v0.8.0 | v0.9.1 |
| dracory/shopstore | v1.27.0 | v1.27.2 |
| dracory/sqlfilestore (indirect) | v1.11.0 | v1.11.2 |
| dracory/statsadmin | v0.8.0 | v0.9.2 |
| dracory/statsstore | v1.33.3 | v1.35.0 |
| dracory/subscriptionstore | v1.6.0 | v1.6.2 |
| dracory/taskstore | v1.33.0 | v1.33.2 |
| dracory/userstore | v1.21.0 | v1.22.2 |
| dracory/vaultstore | v1.6.0 | v1.6.2 |
| dracory/versionstore (indirect) | v1.8.0 | v1.8.2 |

### Third-party packages

| Package | Old | New |
|---|---|---|
| github.com/yuin/goldmark | v1.8.6 | v1.8.7 |
| modernc.org/sqlite | v1.61.0 | v1.61.1 |
| cloud.google.com/go/aiplatform | v1.127.0 | v1.127.1 |
| cloud.google.com/go/auth | v0.24.1 | v0.24.2 |
| cloud.google.com/go/auth/oauth2adapt | v0.3.0 | v0.3.1 |
| cloud.google.com/go/iam | v1.15.0 | v1.15.1 |
| cloud.google.com/go/longrunning | v1.4.0 | v1.4.1 |
| github.com/aws/aws-sdk-go-v2/credentials | v1.20.12 | v1.20.13 |
| google.golang.org/api | v0.309.0 | v0.310.0 |
| google.golang.org/genai | v1.76.0 | v1.77.0 |

---

## 🔄 Migration Steps

### Step 1: Update dependencies

```bash
# Option A: pull the exact versions from this release
go get -u ./...
go mod tidy

# Option B: copy go.mod / go.sum from Blueprint v0.56.0, then
go mod tidy
```

### Step 2: Remove any local auth replace directive

```bash
# Check for leftover replace directives
grep -n "replace" go.mod

# If `github.com/dracory/auth => ...` is listed, delete that line
go mod tidy
```

### Step 3: Wire async log/visitor writes (recommended)

Apply the code changes in the "Async Log and Visitor Writes" section above: swap `NewSlogHandler` for `NewSlogAsyncHandler` in `app.New`, drain it plus `DrainVisitorWrites()` in `app.Close`, and set `VisitorCreateAsyncEnabled: true` in `NewStatsStore`. Add `DrainVisitorWrites()` to any test that asserts on visitor rows after `VisitorCreate`/`StatsMiddleware`.

### Step 4: Update the version constant

If your project tracks the Blueprint version, update `internal/config/version.go`:

```go
const Version = "0.56.0"
```

### Step 5: Verify

```bash
go build ./...
go test ./...
```

---

## ✅ Testing Checklist

- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] Admin areas still render (blogadmin v0.3.2, shopadmin v0.9.1, statsadmin v0.9.2, logadmin v0.2.1)
- [ ] Auth flows work (login, logout, remember-me) on dracory/auth v0.44.0
- [ ] SQLite-backed stores migrate cleanly on modernc.org/sqlite v1.61.1
- [ ] Log records still land in the logstore tables (async handler flushes on its own; wait ~1s or call `Drain()`)
- [ ] Visitor rows appear after page hits — if they seem delayed, that is the async queue working as intended
- [ ] Graceful shutdown flushes pending writes (drain happens before `neatDB.Close()`)

---

## 📝 Notes

- This release exists to keep downstream Blueprint projects on current dependency versions without manual `go get -u` churn.
- Store schema compatibility is unchanged — no new migrations are required by this release.
