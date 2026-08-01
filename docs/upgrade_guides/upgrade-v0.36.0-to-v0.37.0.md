# Upgrade Guide: v0.36.0 to v0.37.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.36.0 to v0.37.0.

## Overview

This release focuses on **Turso connection pool stability**, **security hardening**, and a **Go version bump**. There are no new features, environment variables, or interface changes — the upgrade is low-risk and mostly behavioral.

**Turso Pool Refinement** — Turso (libSQL) no longer shares SQLite's restrictive pool settings. It now gets its own tuned profile (`maxOpenConns=5`, `maxIdleConns=0`, `connMaxLifetime=10s`, `connMaxIdleTime=2s`) suited to Turso's HTTP/2-backed connections. Additionally, the app bootstrap now **explicitly applies** the configured pool settings to the underlying `*sql.DB`, guaranteeing they take effect regardless of whether the `neat` ORM package applies them internally. This fixes intermittent `stream is closed` errors on Turso. The explicit application also defaults SQLite pool settings to 1 when unset (zero), which fixes in-memory SQLite test failures where the DB vanishes between queries.

**Security Hardening** — A pass over gosec findings added justified `#nosec` suppressions, tightened file permissions, switched the contact-form captcha from `math/rand` to `crypto/rand`, added `Secure`/`SameSite` flags to the AI browser auto-login cookie, escaped a user-influenced string in the CDN controller, and replaced manual JSON string concatenation in the shop admin with a safe `writeJSONError` helper. The `gosec` task now runs at `-severity medium`, and a new `govulncheck` task scans dependencies against the Go vulnerability database.

**Key Changes:**
- Go version bumped from `1.26.3` to `1.26.5` in `go.mod`
- Turso gets distinct connection pool settings in `internal/config/database_config.go`
- `internal/app/app_implementation.go` explicitly applies pool settings to `*sql.DB` after `databaseOpen`, with SQLite zero-default fallback
- All dependencies updated to latest versions (`go get -u ./...`)
- `github.com/dracory/llm` bumped from `v1.3.0` to `v1.4.0`
- Contact-form captcha uses `crypto/rand` via new `cryptoRandIntn` helper
- AI browser auto-login cookie gains `Secure: true` and `SameSite: http.SameSiteLaxMode`
- CDN controller escapes the `extension` string with `html.EscapeString`
- Shop admin media upload uses new `writeJSONError` helper (XSS-safe JSON errors)
- File permissions tightened: cache dir `0750`, maintenance state file `0600`
- `taskfile.yml`: `gosec` uses `-severity medium`; new `govulncheck-install` and `govulncheck` tasks

---

## ⚠️ Breaking Changes

---

### 1. Turso Connection Pool Settings Changed

**Change**: In v0.36.0, Turso was treated identically to SQLite for connection pool sizing (`maxOpenConns=1`, `maxIdleConns=1`, `connMaxLifetime=30s`, no special idle time). In v0.37.0, Turso gets its own profile tuned for its HTTP/2-backed remote connections. SQLite is unchanged.

**Old Code**:
```go
// v0.36.0 — internal/config/database_config.go
maxOpenConns := env.GetIntOrDefault(KEY_DB_MAX_OPEN_CONNS, 25)
if driver == driverSQLite || driver == driverTurso {
	maxOpenConns = 1
}

maxIdleConns := env.GetIntOrDefault(KEY_DB_MAX_IDLE_CONNS, 5)
if driver == driverSQLite || driver == driverTurso {
	maxIdleConns = 1
}

connMaxLifetime := time.Duration(env.GetIntOrDefault(KEY_DB_CONN_MAX_LIFETIME_SECONDS, 300)) * time.Second
if driver == driverSQLite || driver == driverTurso {
	connMaxLifetime = 30 * time.Second
}

connMaxIdleTime := time.Duration(env.GetIntOrDefault(KEY_DB_CONN_MAX_IDLE_TIME_SECONDS, 5)) * time.Second
// (no Turso-specific idle time)
```

**New Code**:
```go
// v0.37.0 — internal/config/database_config.go
maxOpenConns := env.GetIntOrDefault(KEY_DB_MAX_OPEN_CONNS, 25)
if driver == driverSQLite {
	maxOpenConns = 1
}
if driver == driverTurso {
	maxOpenConns = 5
}

maxIdleConns := env.GetIntOrDefault(KEY_DB_MAX_IDLE_CONNS, 5)
if driver == driverSQLite {
	maxIdleConns = 1
}
if driver == driverTurso {
	maxIdleConns = 0
}

connMaxLifetime := time.Duration(env.GetIntOrDefault(KEY_DB_CONN_MAX_LIFETIME_SECONDS, 300)) * time.Second
if driver == driverSQLite {
	connMaxLifetime = 30 * time.Second
}
if driver == driverTurso {
	connMaxLifetime = 10 * time.Second
}

connMaxIdleTime := time.Duration(env.GetIntOrDefault(KEY_DB_CONN_MAX_IDLE_TIME_SECONDS, 5)) * time.Second
if driver == driverTurso {
	connMaxIdleTime = 2 * time.Second
}
```

**Action Required**:
- If you use the built-in `databaseConfig()` / `config.NewFromEnv()`, no action needed — the new Turso profile is already applied.
- If you have a **custom** `databaseConfig()` implementation or override pool settings elsewhere, split the `driverSQLite || driverTurso` conditions into separate branches and apply the Turso-specific values above.
- If you previously set `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME_SECONDS`, or `DB_CONN_MAX_IDLE_TIME_SECONDS` in your `.env` to work around Turso's overly-restrictive v0.36.0 defaults, review whether those overrides are still needed — the new defaults may make them unnecessary.

**Migration Command**:
```bash
# Find any code that groups Turso with SQLite for pool sizing
grep -rn "driverTurso" --include="*.go" .
grep -rn "driverSQLite || driver == driverTurso" --include="*.go" .
```

---

### 2. Connection Pool Settings Now Explicitly Applied in `app.New()`

**Change**: `internal/app/app_implementation.go` now calls `db.SetMaxOpenConns`, `db.SetMaxIdleConns`, `db.SetConnMaxLifetime`, and `db.SetConnMaxIdleTime` directly on the `*sql.DB` returned by `neatDB.DB()`, immediately after opening the database. Previously, the app relied on the `neat` package to apply these settings internally. This guarantees the configured pool is always in effect, which is critical for Turso/libSQL where stale HTTP/2 connections cause `stream is closed` errors.

Additionally, when the driver is SQLite and pool settings are unset (0), `app.New()` defaults `MaxOpenConns` and `MaxIdleConns` to 1. This is important for tests: `config.New()` (used in test setups) doesn't run `databaseConfig()` and leaves pool settings at zero. With `MaxIdleConns=0`, every returned connection is closed immediately, causing in-memory SQLite (`mode=memory&cache=shared`) to be destroyed between queries — migrations fail with `no such table: migration_tracker`. The zero-default fallback fixes this automatically for all callers of `app.New()`.

**Old Code**:
```go
// v0.36.0 — internal/app/app_implementation.go
db, err := neatDB.DB()
if err != nil {
	return nil, err
}

// Build app instance
app := &appImplementation{cfg: cfg}
```

**New Code**:
```go
// v0.37.0 — internal/app/app_implementation.go
db, err := neatDB.DB()
if err != nil {
	return nil, err
}

// Explicitly apply connection pool settings to the *sql.DB.
// This guarantees the pool config is applied regardless of whether
// the neat package applies it internally. Critical for Turso/libsql
// where stale HTTP/2 connections cause "stream is closed" errors.
//
// SQLite defaults: when pool settings are unset (0), default to a single
// connection. In-memory SQLite (mode=memory&cache=shared) is destroyed
// when the last connection closes, so MaxIdleConns=0 (no retained idle
// connections) causes the DB to vanish between queries. This is
// especially relevant for tests using config.New() which doesn't run
// databaseConfig() and leaves pool settings at zero.
maxOpen := cfg.GetDatabaseMaxOpenConns()
maxIdle := cfg.GetDatabaseMaxIdleConns()
if strings.EqualFold(cfg.GetDatabaseDriver(), "sqlite") {
	if maxOpen == 0 {
		maxOpen = 1
	}
	if maxIdle == 0 {
		maxIdle = 1
	}
}
db.SetMaxOpenConns(maxOpen)
db.SetMaxIdleConns(maxIdle)
db.SetConnMaxLifetime(time.Duration(cfg.GetDatabaseConnMaxLifetimeSeconds()) * time.Second)
db.SetConnMaxIdleTime(time.Duration(cfg.GetDatabaseConnMaxIdleTimeSeconds()) * time.Second)

// Build app instance
app := &appImplementation{cfg: cfg}
```

**Action Required**:
- If you use the built-in `app.New()`, no action needed — the explicit pool application and SQLite zero-default fallback are already in place.
- If you have a **custom** app constructor or bypass `app.New()` (e.g., you open the database and build the app instance manually), add the pool-setting logic above after obtaining the `*sql.DB`. This requires importing the `time` and `strings` packages if not already imported.
- The config interface methods used (`GetDatabaseMaxOpenConns`, `GetDatabaseMaxIdleConns`, `GetDatabaseConnMaxLifetimeSeconds`, `GetDatabaseConnMaxIdleTimeSeconds`) already existed in v0.36.0, so no interface changes are required.
- **Test setups**: If you have test helpers that call `app.New()` with SQLite and previously pinned pool settings manually (e.g., `cfg.SetDatabaseMaxOpenConns(1)` / `cfg.SetDatabaseMaxIdleConns(1)`), you can remove those manual pins — `app.New()` now handles it automatically.

---

### 3. `github.com/dracory/llm` Bumped to v1.4.0

**Change**: The `github.com/dracory/llm` dependency was bumped from `v1.3.0` to `v1.4.0`. The indirect dependency `github.com/googleapis/enterprise-certificate-proxy` was adjusted from `v0.3.19` to `v0.3.18` as a side effect of the resolution.

**Old `go.mod`**:
```go
github.com/dracory/llm v1.3.0
```

**New `go.mod`**:
```go
github.com/dracory/llm v1.4.0
```

**Action Required**:
- Update the dependency:
  ```bash
  go get github.com/dracory/llm@v1.4.0
  go mod tidy
  ```
- If your application code directly imports and uses `github.com/dracory/llm`, review the [llm changelog](https://github.com/dracory/llm) for any API changes between v1.3.0 and v1.4.0. Blueprint itself only consumes llm indirectly and required no code changes for this bump.
- If you have a vendored copy of `go.sum`, regenerate it with `go mod tidy`.

---

### 4. AI Browser Auto-Login Cookie Now Sets `Secure` and `SameSite`

**Change**: The cookie set by `internal/middlewares/ai_browser_auto_login.go` now includes `Secure: true` and `SameSite: http.SameSiteLaxMode`. This is a security improvement but is **behavioral**: the `Secure` flag means the cookie is only sent over HTTPS.

**Old Code**:
```go
// v0.36.0 — internal/middlewares/ai_browser_auto_login.go
http.SetCookie(w, &http.Cookie{
	Name:     auth.CookieName,
	Value:    session.GetKey(),
	Path:     "/",
	Expires:  time.Now().Add(24 * time.Hour),
	HttpOnly: true,
})
```

**New Code**:
```go
// v0.37.0 — internal/middlewares/ai_browser_auto_login.go
http.SetCookie(w, &http.Cookie{
	Name:     auth.CookieName,
	Value:    session.GetKey(),
	Path:     "/",
	Expires:  time.Now().Add(24 * time.Hour),
	HttpOnly: true,
	Secure:   true,
	SameSite: http.SameSiteLaxMode,
})
```

**Action Required**:
- If you use the built-in middleware, no code changes needed.
- **Local development impact**: If you run the AI browser auto-login flow over plain `http://` (no TLS) in local dev, the browser will **not** store the cookie and auto-login will silently fail. Run local dev behind HTTPS (e.g., via the built-in HTTPS redirect with a self-signed cert, or a tunnel like `mkcert`/`ngrok`) for this flow to work.
- If you have a **custom** AI browser auto-login middleware, add the two new cookie fields to match the security posture.

---

### 5. `gosec` Task Severity Raised and `govulncheck` Task Added

**Change**: In `taskfile.yml`, the `gosec` task now runs with `-severity medium` (previously it ran with no severity filter, which reports low-severity findings). Two new tasks were added: `govulncheck-install` and `govulncheck`.

**Old `taskfile.yml`**:
```yaml
  gosec:
    desc: Tests for security
    cmds:
      - echo "Checking for security..."
      - gosec ./...
      - echo "Done!"
    silent: true
```

**New `taskfile.yml`**:
```yaml
  gosec:
    desc: Tests for security
    cmds:
      - echo "Checking for security..."
      - gosec -severity medium ./...
      - echo "Done!"
    silent: true

  govulncheck-install:
    desc: Install govulncheck or update to latest
    cmds:
      - echo "Installing/updating govulncheck..."
      - go install golang.org/x/vuln/cmd/govulncheck@latest
      - echo "Done!"
    silent: true

  govulncheck:
    desc: Scans dependencies against the Go vulnerability database (blocking)
    cmds:
      - echo "Scanning for known vulnerabilities in dependencies..."
      - govulncheck ./...
      - echo "Done!"
    silent: true
```

**Action Required**:
- If you use the Blueprint `taskfile.yml` directly, no action needed — the changes are already in place.
- If you maintain a **custom** `taskfile.yml` or CI pipeline that invokes `gosec`, decide whether to adopt the `-severity medium` filter. Running without it will still surface low-severity findings; the Blueprint project chose `medium` to reduce noise after adding targeted `#nosec` suppressions.
- To use the new `govulncheck` task, install it first:
  ```bash
  task govulncheck-install
  task govulncheck
  ```
  Or install directly:
  ```bash
  go install golang.org/x/vuln/cmd/govulncheck@latest
  ```

---

### 6. Go Version Bumped to 1.26.5

**Change**: The `go` directive in `go.mod` was bumped from `1.26.3` to `1.26.5`. This aligns with the Go toolchain used to build and test this release.

**Old `go.mod`**:
```go
go 1.26.3
```

**New `go.mod`**:
```go
go 1.26.5
```

**Action Required**:
- Ensure you have Go 1.26.5 or later installed (`go version`).
- Update your `go.mod` directive:
  ```bash
  go mod edit -go=1.26.5
  ```
  Or manually edit `go.mod` and change the `go` line.
- If your CI pipeline pins a specific Go version, update it to 1.26.5 or later.

---

### 7. All Dependencies Updated to Latest

**Change**: All dependencies in `go.mod` were updated to their latest versions via `go get -u ./...` followed by `go mod tidy`. This includes the `github.com/dracory/llm` bump (v1.3.0 → v1.4.0) documented in Breaking Change #3, plus updates to many `dracory/*` store packages, `modernc.org/sqlite`, Google Cloud libraries, and other indirect dependencies.

**Action Required**:
- Update all dependencies:
  ```bash
  go get -u ./...
  go mod tidy
  ```
- Review the resulting `go.mod` / `go.sum` diff for any unexpected version jumps.
- If your application code directly imports any of the updated packages, verify it still compiles and tests pass.
- If you maintain a vendored copy of dependencies, regenerate it after updating.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.37.0"
```

### Step 2: Update Go version and dependencies

Update the `go` directive in `go.mod`:
```bash
go mod edit -go=1.26.5
```

Update all dependencies to latest:
```bash
go get -u ./...
go mod tidy
```

### Step 3: Apply Turso pool setting changes (if customized)

If you have a custom `databaseConfig()` or any code that groups `driverTurso` with `driverSQLite` for pool sizing, split them into separate branches and apply the Turso-specific values:

| Setting | SQLite | Turso (new) | postgres/mysql default |
|---------|--------|-------------|------------------------|
| `maxOpenConns` | 1 | 5 | 25 |
| `maxIdleConns` | 1 | 0 | 5 |
| `connMaxLifetime` | 30s | 10s | 300s |
| `connMaxIdleTime` | 5s (default) | 2s | 5s |

### Step 4: Apply explicit pool settings in your app constructor (if customized)

If you bypass `app.New()`, add after obtaining the `*sql.DB`:

```go
import (
	"strings"
	"time"
)

maxOpen := cfg.GetDatabaseMaxOpenConns()
maxIdle := cfg.GetDatabaseMaxIdleConns()
if strings.EqualFold(cfg.GetDatabaseDriver(), "sqlite") {
	if maxOpen == 0 {
		maxOpen = 1
	}
	if maxIdle == 0 {
		maxIdle = 1
	}
}
db.SetMaxOpenConns(maxOpen)
db.SetMaxIdleConns(maxIdle)
db.SetConnMaxLifetime(time.Duration(cfg.GetDatabaseConnMaxLifetimeSeconds()) * time.Second)
db.SetConnMaxIdleTime(time.Duration(cfg.GetDatabaseConnMaxIdleTimeSeconds()) * time.Second)
```

### Step 5: Apply security hardening (optional but recommended)

These are non-breaking refactors. Apply them if you want to match the Blueprint security posture. They are **not required** for the upgrade to function, but they resolve gosec findings and close minor XSS/permission gaps.

**5a. Contact-form captcha — switch to `crypto/rand`** (`internal/controllers/website/contact/form_contact.go`):

Replace `math/rand` + `time` imports with `crypto/rand` and `encoding/binary`, and add the `cryptoRandIntn` helper:

```go
import (
	"crypto/rand"
	"encoding/binary"
	// ... remove "math/rand" and "time" if no longer used elsewhere
)

// In Mount():
a := cryptoRandIntn(5) + 1 // 1-5
b := cryptoRandIntn(5) + 1 // 1-5

// New helper:
// cryptoRandIntn returns a non-negative pseudo-random int in [0, n) using crypto/rand.
// It falls back to 0 if an error occurs.
func cryptoRandIntn(n int) int {
	if n <= 0 {
		return 0
	}
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return 0
	}
	v := binary.BigEndian.Uint64(buf[:])
	return int(v % uint64(n)) // #nosec G115 -- n is bounded to 5, no overflow risk
}
```

**5b. AI browser auto-login cookie** (`internal/middlewares/ai_browser_auto_login.go`): add `Secure: true` and `SameSite: http.SameSiteLaxMode` to the `http.SetCookie` call. See Breaking Change #4.

**5c. CDN controller extension escaping** (`internal/controllers/shared/cdn/cdn_controller.go`): escape the `extension` string before writing it to the response:

```go
import "html"

// ...
if _, err := w.Write([]byte("Extension " + html.EscapeString(extension) + " not supported")); err != nil { // #nosec G705 -- extension is escaped
```

**5d. Shop admin JSON error helper** (`pkg/shopadmin/product_update/product_update_controller.go`): replace manual `[]byte(`{"status":"error","message":"...` + err.Error() + `"}`) concatenation with the new `writeJSONError` helper to prevent XSS via error messages:

```go
// writeJSONError writes a JSON error response with proper escaping to prevent XSS.
func writeJSONError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	payload := map[string]string{"status": "error", "message": message}
	if data, err := json.Marshal(payload); err == nil {
		w.Write(data)
	}
}

// Usage (replaces manual concatenation):
writeJSONError(w, "Failed to parse upload: "+err.Error())
```

**5e. File permissions** (optional): tighten where applicable:
- `internal/app/app_implementation.go`: `os.MkdirAll(cacheDir, os.ModePerm)` → `os.MkdirAll(cacheDir, 0750)`
- `internal/cli/maintenance_handler.go`: `os.WriteFile(filePath, data, 0644)` → `os.WriteFile(filePath, data, 0600)`

**5f. gosec `#nosec` suppressions**: A large number of justified `#nosec` comments were added across controllers, middlewares, and CLI utilities. These are not required for functionality. If you run `gosec` and want a clean report at `-severity medium`, port the suppressions from the Blueprint reference for the equivalent files in your project. The suppressions cover findings like G204 (hardcoded subprocess commands), G304 (file reads from validated paths), G404 (cosmetic math/rand), G705 (trusted response writes), G710 (validated redirects), G101 (env var key names), G115 (bounded integer conversion), G120 (bounded multipart), G122 (local file rename), G124 (in-request cookie injection), G203 (intentional template URLs).

### Step 6: Update `taskfile.yml` (if customized)

If you maintain a custom `taskfile.yml`, update the `gosec` task to use `-severity medium` and add the `govulncheck-install` / `govulncheck` tasks. See Breaking Change #5.

### Step 7: Run `go mod tidy` and verify

```bash
go mod tidy
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests

```bash
go test ./...
```

All existing tests should pass without modification — this release contains no interface or API signature changes.

### 2. Test SQLite In-Memory Database (Tests)

The SQLite zero-default fallback in `app.New()` ensures in-memory SQLite databases stay alive between queries in tests. If you previously had manual `cfg.SetDatabaseMaxOpenConns(1)` / `cfg.SetDatabaseMaxIdleConns(1)` calls in your test helpers, you can remove them — `app.New()` now handles this automatically.

```bash
go test ./...
```

Verify no `no such table: migration_tracker` or `no such table: snv_files_file` errors appear.

### 3. Test Turso Connection Pool (If Using Turso)

1. Set `DB_DRIVER="turso"` and `DB_DSN="libsql://your-database-url"` in your `.env`.
2. Start the server: `go run ./cmd/server`.
3. Run a workload that opens multiple concurrent queries.
4. Verify no `stream is closed` or `connection reset` errors appear in logs.
5. Confirm the pool is using the new settings (maxOpenConns=5, maxIdleConns=0, connMaxLifetime=10s, connMaxIdleTime=2s). You can log them temporarily:
   ```go
   db, _ := app.GetDatabase()
   fmt.Println("maxOpen:", db.Stats().MaxOpenConnections)
   ```

### 4. Test AI Browser Auto-Login Over HTTPS

1. Run the server behind HTTPS (e.g., `mkcert localhost` + TLS config, or `ngrok http 8080`).
2. Trigger the AI browser auto-login flow.
3. Verify the session cookie is set in the browser (check DevTools → Application → Cookies; `Secure` and `SameSite=Lax` should be present).
4. Confirm the cookie is **not** set when accessing over plain `http://` (expected behavior with `Secure: true`).

### 5. Test Contact Form Captcha

1. Load the contact form page.
2. Verify the math captcha renders (e.g., "3 + 5 =").
3. Submit with the correct answer — should succeed.
4. Submit with the wrong answer — should fail with a captcha error.
5. Submit multiple times — verify the captcha values vary (crypto/rand produces different values each render).

### 6. Test CDN Error Escaping

1. Request a CDN URL with an unsupported extension containing HTML characters, e.g. `/cdn/file.<script>`.
2. Verify the response body contains the escaped string `Extension &lt;script&gt; not supported` and **not** raw `<script>`.

### 7. Run Security Scans

```bash
# Install govulncheck if not present
task govulncheck-install   # or: go install golang.org/x/vuln/cmd/govulncheck@latest

# Run gosec at the new severity threshold
task gosec

# Run govulncheck
task govulncheck
```

---

## 📝 Additional Notes

### New Features

- None. This is a refinement and security-hardening release.

### Removed Features

- None.

### Configuration Changes

No new environment variables. The existing Turso-related variables (`DB_DRIVER`, `DB_DSN`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME_SECONDS`, `DB_CONN_MAX_IDLE_TIME_SECONDS`) now use different **defaults** for Turso, but the variables themselves are unchanged.

| Variable | SQLite default | Turso default (v0.37.0) | Turso default (v0.36.0) |
|----------|----------------|--------------------------|--------------------------|
| `DB_MAX_OPEN_CONNS` | 1 | 5 | 1 |
| `DB_MAX_IDLE_CONNS` | 1 | 0 | 1 |
| `DB_CONN_MAX_LIFETIME_SECONDS` | 30 | 10 | 30 |
| `DB_CONN_MAX_IDLE_TIME_SECONDS` | 5 | 2 | 5 |

Explicit env var overrides still take precedence over these defaults.

### Why Turso Needed Different Pool Settings

Turso connects over HTTP/2 to a remote edge database. SQLite's `maxOpenConns=1` (designed to avoid local file lock contention) is far too restrictive for Turso and caused serialization of all queries. Conversely, keeping idle connections alive (`maxIdleConns=1`) and a long lifetime (`30s`) allowed stale HTTP/2 streams to be reused, producing `stream is closed` errors. The new profile (`maxOpenConns=5`, `maxIdleConns=0`, `connMaxLifetime=10s`, `connMaxIdleTime=2s`) allows reasonable concurrency while aggressively retiring connections before they go stale.

### Why Pool Settings Are Now Applied Explicitly

The `neat` ORM wrapper may or may not forward pool settings to the underlying `*sql.DB` depending on its version and configuration. By calling `db.Set*` directly in `app.New()`, Blueprint removes that uncertainty. This is especially important for Turso, where the wrong pool settings manifest as intermittent runtime errors rather than compile-time failures.

The SQLite zero-default fallback (defaulting `MaxOpenConns` and `MaxIdleConns` to 1 when unset) was added as part of this explicit application. It addresses a test-only issue: `config.New()` (used in test setups) doesn't run `databaseConfig()` and leaves pool settings at zero, which causes in-memory SQLite databases to be destroyed between queries. Centralizing the fix in `app.New()` means all test setups — including future ones — automatically get the correct behavior without needing to remember to pin pool settings manually.

---

## 🆘 Common Issues and Solutions

### Issue 1: Turso `stream is closed` Errors Persist After Upgrade

**Symptom**: After upgrading to v0.37.0, you still see intermittent `stream is closed` or `connection reset by peer` errors on Turso.

**Solution**:
1. Confirm you applied **both** Breaking Change #1 (Turso pool defaults) **and** Breaking Change #2 (explicit `db.Set*` calls). The explicit application is what actually enforces the settings.
2. Check your `.env` for explicit overrides that may be forcing the old restrictive values:
   ```bash
   grep -E "DB_MAX_OPEN_CONNS|DB_MAX_IDLE_CONNS|DB_CONN_MAX_LIFETIME|DB_CONN_MAX_IDLE" .env
   ```
   Remove or adjust any that conflict with the new Turso defaults.
3. Verify the settings are applied at runtime by logging `db.Stats()`.

### Issue 2: AI Browser Auto-Login Stops Working in Local Dev

**Symptom**: After upgrading, the AI browser auto-login flow no longer logs you in when running locally over `http://localhost`.

**Solution**: The `Secure: true` cookie flag prevents the browser from storing the cookie over plain HTTP. Run local dev over HTTPS:
- Use `mkcert localhost` and configure TLS, or
- Use a tunnel (`ngrok http 8080`) to get an HTTPS endpoint, or
- Temporarily comment out `Secure: true` in your local dev branch (do **not** commit this).

### Issue 3: `gosec` Reports Many Findings After Adopting `-severity medium`

**Symptom**: Running `task gosec` now reports findings that were previously hidden.

**Solution**: The `-severity medium` filter actually **reduces** noise by suppressing low-severity findings. If you see more findings than before, you likely were not running gosec previously, or your project has not yet ported the `#nosec` suppressions from Blueprint. Either:
- Port the relevant `#nosec` suppressions from the Blueprint reference for equivalent files (see Step 5f), or
- Audit each finding and add a justified `#nosec` comment where the concern is mitigated by context.

### Issue 4: `govulncheck` Reports Vulnerabilities in Dependencies

**Symptom**: `task govulncheck` exits non-zero reporting vulnerabilities.

**Solution**:
1. Review the reported vulnerabilities — `govulncheck` only reports vulnerabilities in code paths your project actually calls (call-reachable analysis).
2. Update the affected dependency to a fixed version if available.
3. If no fix is available yet, document the accepted risk and consider excluding the vulnerability if your CI supports it.
4. Re-run `govulncheck` to confirm resolution.

### Issue 5: Compile Error After `go get github.com/dracory/llm@v1.4.0`

**Symptom**: After bumping llm, you get compile errors in code that directly uses the `github.com/dracory/llm` package.

**Solution**: Review the [llm changelog](https://github.com/dracory/llm) for breaking changes between v1.3.0 and v1.4.0 and update your call sites accordingly. Blueprint itself does not directly use llm APIs in a way that broke, so if you only consume llm through Blueprint-provided wrappers, no changes should be needed. Run `go mod tidy` to ensure the dependency graph is consistent.

### Issue 6: Tests Fail with `no such table: migration_tracker` on SQLite

**Symptom**: After upgrading, tests that use in-memory SQLite fail with errors like `no such table: migration_tracker` or `no such table: snv_files_file`.

**Solution**: This happens when a test setup bypasses `app.New()` or uses a custom app constructor that doesn't include the SQLite zero-default fallback. The in-memory SQLite database (`mode=memory&cache=shared`) is destroyed when the last connection closes, and with `MaxIdleConns=0` (the default from `config.New()`), every returned connection is closed immediately.

1. If your test calls `app.New()`, ensure you're using the updated `app.New()` which includes the SQLite zero-default fallback.
2. If you have a **custom** app constructor, add the SQLite zero-default logic (see Breaking Change #2 / Step 4):
   ```go
   if strings.EqualFold(cfg.GetDatabaseDriver(), "sqlite") {
       if maxOpen == 0 { maxOpen = 1 }
       if maxIdle == 0 { maxIdle = 1 }
   }
   ```
3. Alternatively, pin the pool settings manually in your test config:
   ```go
   cfg.SetDatabaseMaxOpenConns(1)
   cfg.SetDatabaseMaxIdleConns(1)
   ```
   (This is what `app.New()` does automatically — only needed if you bypass it.)

---

## Support

For issues or questions about this upgrade:
- Check the Turso pool config: `internal/config/database_config.go`
- Check the explicit pool application with SQLite defaults: `internal/app/app_implementation.go`
- Check the contact-form captcha helper: `internal/controllers/website/contact/form_contact.go`
- Check the AI browser auto-login cookie: `internal/middlewares/ai_browser_auto_login.go`
- Check the shop admin JSON error helper: `pkg/shopadmin/product_update/product_update_controller.go`
- Check the security task definitions: `taskfile.yml`
- Check the Go version and dependencies: `go.mod`
- Open an issue on the [Blueprint repository](https://github.com/dracory/blueprint)
