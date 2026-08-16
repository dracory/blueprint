# Upgrade Guide: v0.37.0 to v0.38.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.37.0 to v0.38.0.

## Overview

This release focuses on **removing the `dracory/sb` dependency**, **environment-aware secure cookies**, a new **Dokploy webhook deployment task**, and **dependency updates**. The main breaking change is the migration of sort-order constants from `github.com/dracory/sb` to `github.com/dracory/neat`, and a signature change to the user-impersonation helper.

**`dracory/sb` Removed** — The `sb.ASC` / `sb.DESC` sort-order constants are now sourced from `github.com/dracory/neat` (`neat.SortAsc` / `neat.SortDesc`), which are aliases for the canonical `orm.SortAsc` / `orm.SortDesc`. The `github.com/dracory/sb` direct dependency has been removed entirely from `go.mod`.

**Environment-Aware Secure Cookies** — Authentication and user-impersonation cookies now respect the application environment. In development (`APP_ENV=development`), the `Secure` flag is disabled so the browser sends the cookie back over plain HTTP. In production/staging, the `Secure` flag is maintained. The `Impersonate()` function signature gained a `secure bool` parameter, and `auth.AuthCookieSet` calls now pass `types.WithSecure(...)` cookie options.

**Dokploy Webhook Deployment** — A new `deploy-dokploy` task was added to `taskfile.yml` for triggering deployments via a Dokploy webhook, with cross-platform support (PowerShell on Windows, curl on Linux/Darwin).

**`AUTH_CSRF_SECRET` Documented** — The `AUTH_CSRF_SECRET` environment variable (whose code logic already existed in v0.37.0) is now documented in `.env.example` with a security warning. The MCP configuration section header was also clarified to spell out the acronym.

**Key Changes:**
- `github.com/dracory/sb` removed from `go.mod`; sort constants migrated to `github.com/dracory/neat`
- `pkg/useradmin/user_impersonate/impersonate.go`: `Impersonate()` signature changed — new `secure bool` parameter
- `internal/controllers/auth/authentication/authentication_controller.go`: environment-aware `Secure` cookie flag
- `pkg/useradmin/user_impersonate/user_impersonate_controller.go`: passes `secure` flag based on `IsEnvDevelopment()`
- `AUTH_CSRF_SECRET` entry added to `.env.example` with security warning
- MCP section header in `.env.example` clarified: `MCP Configuration (Model Context Protocol)`
- New `deploy-dokploy` task in `taskfile.yml` with `DEPLOY_WEBHOOK_URL` variable
- `github.com/dracory/neat` bumped from `v0.33.0` to `v0.39.0`
- `github.com/dracory/rtr` bumped from `v1.7.0` to `v1.9.0`
- `github.com/dracory/websrv` bumped from `v0.1.0` to `v0.4.0`
- `github.com/dracory/api` bumped from `v1.7.0` to `v1.8.0`
- `github.com/dracory/envenc` bumped from `v1.4.1` to `v1.5.0`
- `github.com/dracory/taskstore` bumped from `v1.29.0` to `v1.30.0`
- `modernc.org/sqlite` bumped from `v1.54.0` to `v1.56.0`
- `dromara/carbon/v2` bumped from `v2.6.16` to `v2.6.17`
- Removed unused indirect dependencies: `go-chi/chi/v5`, `go-chi/cors`, `go-chi/httprate`, `zeebo/xxh3`, `mingrammer/cfmt`, `logrusorgru/aurora`, `mattn/go-colorable`, `georgysavva/scany`, `jackc/pgservicefile`, `klauspost/cpuid/v2`

---

## ⚠️ Breaking Changes

---

### 1. `dracory/sb` Removed — Sort Constants Migrated to `dracory/neat`

**Change**: The `github.com/dracory/sb` package has been removed as a direct dependency. All usages of `sb.ASC` and `sb.DESC` must be replaced with `neat.SortAsc` and `neat.SortDesc` from `github.com/dracory/neat`. The `neat` constants are aliases for `orm.SortAsc` / `orm.SortDesc` (the canonical source of truth in `contracts/database/orm`).

**Old Usage**:
```go
import (
	"github.com/dracory/sb"
)

// ...
countryList, err := c.app.GetGeoStore().CountryList(ctx, geostore.CountryQueryOptions{
	SortOrder: sb.ASC,
	OrderBy:   geostore.COLUMN_NAME,
})

// ...
query := geostore.TimezoneQueryOptions{
	SortOrder: sb.ASC,
	OrderBy:   geostore.COLUMN_TIMEZONE,
}

// ...
postList, err := c.app.GetBlogStore().PostList(context.Background(), blogstore.PostQueryOptions{
	Status:    blogstore.POST_STATUS_PUBLISHED,
	OrderBy:   "title",
	SortOrder: sb.DESC,
	Limit:     1000,
})
```

**New Usage**:
```go
import (
	"github.com/dracory/neat"
)

// ...
countryList, err := c.app.GetGeoStore().CountryList(ctx, geostore.CountryQueryOptions{
	SortOrder: neat.SortAsc,
	OrderBy:   geostore.COLUMN_NAME,
})

// ...
query := geostore.TimezoneQueryOptions{
	SortOrder: neat.SortAsc,
	OrderBy:   geostore.COLUMN_TIMEZONE,
}

// ...
postList, err := c.app.GetBlogStore().PostList(context.Background(), blogstore.PostQueryOptions{
	Status:    blogstore.POST_STATUS_PUBLISHED,
	OrderBy:   "title",
	SortOrder: neat.SortDesc,
	Limit:     1000,
})
```

**Action Required**:
- Search your codebase for all usages of `sb.ASC` and `sb.DESC`:
  ```bash
  grep -rn "sb\.ASC\|sb\.DESC" --include="*.go" .
  ```
- Replace `sb.ASC` with `neat.SortAsc` and `sb.DESC` with `neat.SortDesc`.
- Replace the import `"github.com/dracory/sb"` with `"github.com/dracory/neat"`.
- Remove `github.com/dracory/sb` from `go.mod`:
  ```bash
  go mod tidy
  ```
- **Note**: The string values are identical (`"asc"` / `"desc"`), so this is a source-level rename — no database or runtime behavior changes.

**Migration Command**:
```bash
# Find all files using sb sort constants
grep -rln "dracory/sb" --include="*.go" .

# Find all sb.ASC / sb.DESC usages
grep -rn "sb\.ASC\|sb\.DESC" --include="*.go" .
```

---

### 2. `Impersonate()` Function Signature Changed

**Change**: The `Impersonate()` function in `pkg/useradmin/user_impersonate/impersonate.go` now accepts a `secure bool` parameter. This controls whether the `Secure` flag is set on the authentication cookie. The function now imports `github.com/dracory/auth/types` and passes `types.WithSecure(secure)` to `auth.AuthCookieSet`.

**Old Usage**:
```go
// v0.37.0 — pkg/useradmin/user_impersonate/impersonate.go
import (
	"github.com/dracory/auth"
)

func Impersonate(ss sessionstore.StoreInterface, w http.ResponseWriter, r *http.Request, userID string) error {
	// ...
	auth.AuthCookieSet(w, r, session.GetKey())
	return nil
}

// Caller:
err := Impersonate(c.app.GetSessionStore(), w, r, userID)
```

**New Usage**:
```go
// v0.38.0 — pkg/useradmin/user_impersonate/impersonate.go
import (
	"github.com/dracory/auth"
	"github.com/dracory/auth/types"
)

func Impersonate(ss sessionstore.StoreInterface, w http.ResponseWriter, r *http.Request, userID string, secure bool) error {
	// ...
	auth.AuthCookieSet(w, r, session.GetKey(), types.WithSecure(secure))
	return nil
}

// Caller (environment-aware):
secure := true
if c.app.GetConfig() != nil && c.app.GetConfig().IsEnvDevelopment() {
	secure = false
}
err := Impersonate(c.app.GetSessionStore(), w, r, userID, secure)
```

**Action Required**:
- Update all call sites of `Impersonate()` to pass the new `secure` parameter.
- If you call `Impersonate()` from a controller that has access to `app.AppInterface`, use the environment-aware pattern shown above:
  ```go
  secure := true
  if c.app.GetConfig() != nil && c.app.GetConfig().IsEnvDevelopment() {
      secure = false
  }
  ```
- If you call `Impersonate()` from a context without config access, pass `true` (secure by default) unless you have a specific reason to disable it.
- If you have a **custom** impersonation wrapper, update it to accept and forward the `secure` parameter, and add the `github.com/dracory/auth/types` import.

**Migration Command**:
```bash
# Find all call sites of Impersonate
grep -rn "Impersonate(" --include="*.go" .
```

---

### 3. Authentication Cookie Now Environment-Aware (`Secure` Flag)

**Change**: The authentication controller (`internal/controllers/auth/authentication/authentication_controller.go`) now sets the `Secure` cookie flag conditionally based on the application environment. In development, `Secure` is disabled (via `types.WithSecure(false)`) so the browser sends the cookie over plain HTTP. In production/staging, the default `Secure` behavior is maintained.

**Old Usage**:
```go
// v0.37.0 — internal/controllers/auth/authentication/authentication_controller.go
auth.AuthCookieSet(w, r, session.GetKey())
```

**New Usage**:
```go
// v0.38.0 — internal/controllers/auth/authentication/authentication_controller.go
import (
	"github.com/dracory/auth"
	"github.com/dracory/auth/types"
)

// In development (HTTP), the Secure flag must be disabled so the
// browser sends the cookie back over plain HTTP.
cookieOpts := []types.CookieOption{}
if c.app.GetConfig() != nil && c.app.GetConfig().IsEnvDevelopment() {
	cookieOpts = append(cookieOpts, types.WithSecure(false))
}

auth.AuthCookieSet(w, r, session.GetKey(), cookieOpts...)
```

**Action Required**:
- If you use the built-in authentication controller, no action needed — the environment-aware logic is already in place.
- If you have a **custom** authentication controller or any code that calls `auth.AuthCookieSet()` directly, apply the same environment-aware pattern:
  ```go
  cookieOpts := []types.CookieOption{}
  if cfg != nil && cfg.IsEnvDevelopment() {
      cookieOpts = append(cookieOpts, types.WithSecure(false))
  }
  auth.AuthCookieSet(w, r, token, cookieOpts...)
  ```
- Add the `github.com/dracory/auth/types` import where needed.
- **Note**: `auth.AuthCookieSet` uses variadic `...types.CookieOption`, so existing calls without options remain backward compatible at the call site. However, the *behavior* changes: without options, the cookie defaults to `Secure: true` (the `auth` package default). If you need insecure cookies in development, you must explicitly pass `types.WithSecure(false)`.

**Migration Command**:
```bash
# Find all AuthCookieSet call sites
grep -rn "AuthCookieSet(" --include="*.go" .
```

---

### 4. `AUTH_CSRF_SECRET` Added to `.env.example`

**Change**: The `AUTH_CSRF_SECRET` environment variable is now documented in `.env.example` with a security warning. The code logic (the `KEY_AUTH_CSRF_SECRET` constant in `internal/config/constants.go` and the CSRF secret loading/generation in `internal/config/auth_config.go`) already existed in v0.37.0 — this release only adds the example entry and documentation.

**Old `.env.example`** (CSRF section absent — appeared after the `AUTH_EMAILS_ALLOWED_ACCESS` block):
```bash
# AUTH_EMAILS_ALLOWED_ACCESS=""

# ============================================================================
# Payment Configuration
# ============================================================================
```

**New `.env.example`**:
```bash
# AUTH_EMAILS_ALLOWED_ACCESS=""

# CSRF Secret
# Secret key used for CSRF token generation and validation.
# Default: courseflow-auth
# WARNING: Change in production to a strong random string.
AUTH_CSRF_SECRET="YOUR_SECURE_RANDOM_STRING"

# ============================================================================
# Payment Configuration
# ============================================================================
```

**Action Required**:
- Add the `AUTH_CSRF_SECRET` entry to your project's `.env.example` (and `.env` if not already present).
- **In production/staging**: Set `AUTH_CSRF_SECRET` to a strong random string. If unset in production/staging, the application will **panic** at startup (`FATAL: AUTH_CSRF_SECRET must be set in the production environment`).
- **In development/testing**: If `AUTH_CSRF_SECRET` is unset, a random secret is generated automatically for each run and a warning is logged. This is fine for local dev but means CSRF tokens won't survive restarts.
- The constant `KEY_AUTH_CSRF_SECRET` and the `authConfig()` logic were already present in v0.37.0, so no code changes are required — only the `.env.example` documentation.

---

### 5. MCP Section Header Clarified in `.env.example`

**Change**: The MCP configuration section header in `.env.example` was updated to spell out the acronym.

**Old `.env.example`**:
```bash
# ============================================================================
# MCP Configuration
# ============================================================================
```

**New `.env.example`**:
```bash
# ============================================================================
# MCP Configuration (Model Context Protocol)
# ============================================================================
```

**Action Required**:
- Optionally update your `.env.example` header for clarity. This is a documentation-only change with no functional impact.

---

### 6. Dokploy Webhook Deployment Task Added to `taskfile.yml`

**Change**: A new `deploy-dokploy` task was added to `taskfile.yml` for triggering deployments via a Dokploy webhook. It includes a commented-out `DEPLOY_WEBHOOK_URL` variable and cross-platform commands (PowerShell for Windows, curl for Linux/Darwin).

**New `taskfile.yml`** (vars section):
```yaml
vars:
  APPNAME: The Dracory Blueprint Project
  DATETIME: '{{now | date "20060102_150405"}}'
  LIVEURL: https://dracory.com
  # Uncomment and set your Dokploy webhook URL to use deploy-dokploy task
  # DEPLOY_WEBHOOK_URL: https://your-dokploy-server/api/deploy/YOUR_WEBHOOK_TOKEN
```

**New `taskfile.yml`** (tasks section):
```yaml
  # Alternative deployment via Dokploy webhook (Hetzner).
  # To use: uncomment DEPLOY_WEBHOOK_URL in vars above and set your webhook URL.
  deploy-dokploy:
    desc: Triggers a deployment via Dokploy webhook (Hetzner)
    cmds:
      - echo "Triggering Dokploy deployment webhook..."
      - cmd: echo "Note - webhook acknowledgement only confirms the build was triggered."
      - cmd: echo "Docker image build + container rollout takes ~5 minutes after trigger."
      - cmd: echo ""
      - cmd: echo "-------------------------------------"
      - cmd: echo "--- Dokploy server response below ---"
      - cmd: echo "-------------------------------------"
      - cmd: echo ""
      - cmd: powershell -Command "Invoke-RestMethod -Method Post -Uri '{{.DEPLOY_WEBHOOK_URL}}' -ContentType 'application/json' -Headers @{'X-GitHub-Event'='push'} -Body '{\"ref\":\"refs/heads/main\"}' | Out-String"
        platforms: [windows]
      - cmd: >-
          curl -X POST -fsS '{{.DEPLOY_WEBHOOK_URL}}'
          -H 'Content-Type: application/json'
          -H 'X-GitHub-Event: push'
          -d '{"ref":"refs/heads/main"}'
        platforms: [linux, darwin]
      - cmd: echo ""
      - cmd: echo "-----------------------------------"
      - cmd: echo "--- End Dokploy server response ---"
      - cmd: echo "-----------------------------------"
      - cmd: echo ""
      - cmd: echo "Build triggered. Not yet live - Docker build is now running (~5 min)."
      - echo "Done!"
    silent: true
```

**Action Required**:
- If you use the Blueprint `taskfile.yml` directly, no action needed — the task is already in place.
- If you maintain a **custom** `taskfile.yml` and use Dokploy, add the `deploy-dokploy` task and the `DEPLOY_WEBHOOK_URL` variable (commented by default).
- To use the task: uncomment `DEPLOY_WEBHOOK_URL` in the `vars` section and set it to your Dokploy webhook URL, then run:
  ```bash
  task deploy-dokploy
  ```
- This is an **additive** change — no existing tasks were modified or removed.

---

### 7. Dependency Updates

**Change**: Several `dracory/*` dependencies were bumped, and unused indirect dependencies were removed.

**Old `go.mod`** (key direct dependencies):
```go
github.com/dracory/api v1.7.0
github.com/dracory/envenc v1.4.1
github.com/dracory/neat v0.33.0
github.com/dracory/rtr v1.7.0
github.com/dracory/sb v0.26.0
github.com/dracory/taskstore v1.29.0
github.com/dracory/websrv v0.1.0
github.com/dromara/carbon/v2 v2.6.16
modernc.org/sqlite v1.54.0
```

**New `go.mod`** (key direct dependencies):
```go
github.com/dracory/api v1.8.0
github.com/dracory/envenc v1.5.0
github.com/dracory/neat v0.39.0
github.com/dracory/rtr v1.9.0
github.com/dracory/taskstore v1.30.0
github.com/dracory/websrv v0.4.0
github.com/dromara/carbon/v2 v2.6.17
modernc.org/sqlite v1.56.0
// github.com/dracory/sb — REMOVED
```

**Removed indirect dependencies**:
- `github.com/go-chi/chi/v5`
- `github.com/go-chi/cors`
- `github.com/go-chi/httprate`
- `github.com/zeebo/xxh3`
- `github.com/mingrammer/cfmt`
- `github.com/logrusorgru/aurora`
- `github.com/mattn/go-colorable`
- `github.com/georgysavva/scany`
- `github.com/jackc/pgservicefile`
- `github.com/klauspost/cpuid/v2`

**Action Required**:
- Update dependencies:
  ```bash
  go get github.com/dracory/api@v1.8.0
  go get github.com/dracory/envenc@v1.5.0
  go get github.com/dracory/neat@v0.39.0
  go get github.com/dracory/rtr@v1.9.0
  go get github.com/dracory/taskstore@v1.30.0
  go get github.com/dracory/websrv@v0.4.0
  go get github.com/dromara/carbon/v2@v2.6.17
  go get modernc.org/sqlite@v1.56.0
  go mod tidy
  ```
- The `dracory/sb` removal is handled by Breaking Change #1 (migrate sort constants first, then `go mod tidy` will remove it).
- The removed indirect dependencies (chi, xxh3, cfmt, etc.) were previously pulled in by `dracory/websrv` and `dracory/sb`. After bumping `websrv` to `v0.4.0` and removing `sb`, `go mod tidy` will clean them up automatically.
- If your application code directly imports any of the removed indirect dependencies, you will need to find an alternative or vendor the code.
- Review the resulting `go.mod` / `go.sum` diff for any unexpected version jumps.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.38.0"
```

### Step 2: Migrate sort constants from `sb` to `neat`

Search for all usages of `sb.ASC` and `sb.DESC`:

```bash
grep -rn "sb\.ASC\|sb\.DESC" --include="*.go" .
```

For each file found:
1. Replace `"github.com/dracory/sb"` import with `"github.com/dracory/neat"`.
2. Replace `sb.ASC` with `neat.SortAsc`.
3. Replace `sb.DESC` with `neat.SortDesc`.

If `sb` is imported but only used for sort constants, remove the `sb` import entirely. If `sb` is used for other purposes, you will need to find alternatives or keep a vendored copy (though `sb` is being phased out).

### Step 3: Update `Impersonate()` call sites

Search for all call sites:

```bash
grep -rn "Impersonate(" --include="*.go" .
```

For each call site, add the `secure` parameter using the environment-aware pattern:

```go
secure := true
if c.app.GetConfig() != nil && c.app.GetConfig().IsEnvDevelopment() {
	secure = false
}
err := Impersonate(c.app.GetSessionStore(), w, r, userID, secure)
```

### Step 4: Update authentication cookie calls (if customized)

If you have a custom authentication controller or call `auth.AuthCookieSet()` directly, apply the environment-aware `Secure` flag pattern (see Breaking Change #3). Add the `github.com/dracory/auth/types` import.

### Step 5: Update `.env.example`

Add the `AUTH_CSRF_SECRET` entry after the `AUTH_EMAILS_ALLOWED_ACCESS` block:

```bash
# CSRF Secret
# Secret key used for CSRF token generation and validation.
# Default: courseflow-auth
# WARNING: Change in production to a strong random string.
AUTH_CSRF_SECRET="YOUR_SECURE_RANDOM_STRING"
```

Optionally update the MCP section header:

```bash
# ============================================================================
# MCP Configuration (Model Context Protocol)
# ============================================================================
```

### Step 6: Add Dokploy deployment task (optional)

If you use Dokploy for deployments, add the `deploy-dokploy` task and `DEPLOY_WEBHOOK_URL` variable to your `taskfile.yml` (see Breaking Change #6). Uncomment and set the webhook URL to use it.

### Step 7: Update dependencies

```bash
go get github.com/dracory/api@v1.8.0
go get github.com/dracory/envenc@v1.5.0
go get github.com/dracory/neat@v0.39.0
go get github.com/dracory/rtr@v1.9.0
go get github.com/dracory/taskstore@v1.30.0
go get github.com/dracory/websrv@v0.4.0
go get github.com/dromara/carbon/v2@v2.6.17
go get modernc.org/sqlite@v1.56.0
go mod tidy
```

### Step 8: Run `go mod tidy` and verify

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

All existing tests should pass without modification — the sort constant migration is a source-level rename with identical string values, and the cookie changes are behavioral (environment-aware).

### 2. Test Authentication Over HTTP (Development)

1. Set `APP_ENV=development` in your `.env`.
2. Start the server over plain HTTP: `go run ./cmd/server`.
3. Log in via the authentication flow.
4. Verify the session cookie is set in the browser (check DevTools → Application → Cookies). The `Secure` flag should **not** be present (because `APP_ENV=development`).
5. Confirm the cookie is sent back on subsequent HTTP requests and you remain logged in.

### 3. Test Authentication Over HTTPS (Production)

1. Set `APP_ENV=production` and ensure `AUTH_CSRF_SECRET` is set.
2. Start the server behind HTTPS.
3. Log in via the authentication flow.
4. Verify the session cookie has the `Secure` flag set in DevTools.
5. Confirm the cookie is **not** sent over plain HTTP (expected behavior with `Secure: true`).

### 4. Test User Impersonation

1. As an admin, trigger user impersonation.
2. In development (`APP_ENV=development`): verify the impersonation cookie works over HTTP (no `Secure` flag).
3. In production (`APP_ENV=production`): verify the impersonation cookie has the `Secure` flag and only works over HTTPS.

### 5. Test CSRF Secret Behavior

1. **Production without `AUTH_CSRF_SECRET`**: Start the server with `APP_ENV=production` and `AUTH_CSRF_SECRET` unset. Verify the application panics with `FATAL: AUTH_CSRF_SECRET must be set in the production environment`.
2. **Development without `AUTH_CSRF_SECRET`**: Start the server with `APP_ENV=development` and `AUTH_CSRF_SECRET` unset. Verify a warning is logged and a random secret is generated.
3. **With `AUTH_CSRF_SECRET` set**: Verify CSRF tokens are generated and validated correctly across restarts (persistent secret).

### 6. Test Dokploy Deployment (If Using)

1. Uncomment `DEPLOY_WEBHOOK_URL` in `taskfile.yml` and set your webhook URL.
2. Run `task deploy-dokploy`.
3. Verify the webhook is triggered and the Dokploy server responds.
4. Confirm the deployment proceeds (Docker build + container rollout, ~5 minutes).

### 7. Verify `sb` Is Fully Removed

```bash
# Should return no results
grep -rn "dracory/sb" --include="*.go" .
grep -rn "sb\.ASC\|sb\.DESC" --include="*.go" .

# Should not appear in go.mod
grep "dracory/sb" go.mod
```

---

## 📝 Additional Notes

### New Features

- **Dokploy webhook deployment task** (`taskfile.yml`): Cross-platform task for triggering Dokploy deployments via webhook.
- **Environment-aware secure cookies**: Authentication and impersonation cookies now adapt to the application environment, improving local development experience over plain HTTP.

### Removed Features

- **`github.com/dracory/sb` dependency**: Removed entirely. Sort constants migrated to `github.com/dracory/neat`.
- **Unused indirect dependencies**: `go-chi/chi/v5`, `go-chi/cors`, `go-chi/httprate`, `zeebo/xxh3`, `mingrammer/cfmt`, `logrusorgru/aurora`, `mattn/go-colorable`, `georgysavva/scany`, `jackc/pgservicefile`, `klauspost/cpuid/v2` — these were pulled in by `websrv` (pre-v0.4.0) and `sb`, and are no longer needed.

### Configuration Changes

| Variable | Status | Notes |
|----------|--------|-------|
| `AUTH_CSRF_SECRET` | Documented in `.env.example` | Code logic existed in v0.37.0; now documented with security warning. Required in production/staging (panics if unset). Auto-generated in development. |

### Why `sb` Was Removed

The `github.com/dracory/sb` package provided sort-order constants (`ASC` / `DESC`) among other utilities. These constants are now canonically defined in `contracts/database/orm` and re-exported via `github.com/dracory/neat` (`neat.SortAsc` / `neat.SortDesc`). Removing `sb` eliminates a redundant dependency and consolidates sort-order definitions in the ORM layer. The string values (`"asc"` / `"desc"`) are identical, so there is no runtime behavior change.

### Why Cookies Are Now Environment-Aware

In v0.37.0, the AI browser auto-login cookie gained `Secure: true` (see v0.36.0→v0.37.0 guide). This broke local development over plain HTTP because browsers refuse to store `Secure` cookies without TLS. v0.38.0 extends environment-awareness to the main authentication and impersonation cookies: in development, `Secure` is disabled so cookies work over HTTP; in production, `Secure` is maintained for security. This aligns the authentication flow with the same pattern used for the AI browser auto-login cookie.

---

## 🆘 Common Issues and Solutions

### Issue 1: Compile Error — `undefined: sb.ASC` or `undefined: sb.DESC`

**Symptom**: After upgrading, you get compile errors like `undefined: sb.ASC` or `undefined: sb.DESC`.

**Solution**: You have code that still references `sb.ASC` / `sb.DESC` but the `dracory/sb` import was removed. Apply Breaking Change #1:
1. Replace `"github.com/dracory/sb"` import with `"github.com/dracory/neat"`.
2. Replace `sb.ASC` with `neat.SortAsc`.
3. Replace `sb.DESC` with `neat.SortDesc`.
4. Run `go mod tidy`.

```bash
grep -rn "sb\.ASC\|sb\.DESC" --include="*.go" .
```

### Issue 2: Compile Error — `not enough arguments in call to Impersonate`

**Symptom**: After upgrading, you get a compile error like `not enough arguments in call to Impersonate`.

**Solution**: The `Impersonate()` function now requires a `secure bool` parameter. Apply Breaking Change #2 — update all call sites to pass the `secure` flag:

```go
secure := true
if c.app.GetConfig() != nil && c.app.GetConfig().IsEnvDevelopment() {
	secure = false
}
err := Impersonate(c.app.GetSessionStore(), w, r, userID, secure)
```

### Issue 3: Login Stops Working Over HTTP in Development

**Symptom**: After upgrading, logging in over `http://localhost` no longer works — the browser doesn't store the session cookie.

**Solution**: Ensure `APP_ENV=development` is set in your `.env`. The environment-aware cookie logic checks `IsEnvDevelopment()` to disable the `Secure` flag. If `APP_ENV` is unset or set to `production`, the cookie will have `Secure: true` and won't be stored over HTTP.

### Issue 4: Application Panics on Startup in Production

**Symptom**: After deploying to production, the application panics with `FATAL: AUTH_CSRF_SECRET must be set in the production environment`.

**Solution**: Set `AUTH_CSRF_SECRET` in your production environment to a strong random string:
```bash
# Generate a random 32-byte hex string
openssl rand -hex 32
```
Add it to your `.env` or environment configuration. This requirement existed in v0.37.0 (the code logic was already present), but the `.env.example` documentation is new in v0.38.0.

### Issue 5: `go mod tidy` Re-adds `dracory/sb`

**Symptom**: After running `go mod tidy`, `github.com/dracory/sb` reappears in `go.mod`.

**Solution**: Some code still imports `github.com/dracory/sb`. Search for remaining imports:
```bash
grep -rn "dracory/sb" --include="*.go" .
```
Migrate any remaining usages to `neat` (or another appropriate package), then run `go mod tidy` again.

### Issue 6: Compile Error After Bumping `dracory/websrv` to v0.4.0

**Symptom**: After bumping `websrv`, you get compile errors in code that directly uses the `websrv` package.

**Solution**: Review the `websrv` changelog for breaking changes between v0.1.0 and v0.4.0. The `chi` router dependencies were removed from `websrv` in this range, so if your code relied on `chi`-related types exposed through `websrv`, you may need to import `chi` directly or migrate to the `websrv`-native routing API. Blueprint itself required no code changes for this bump.

---

## Support

For issues or questions about this upgrade:
- Check the sort constant migration: `internal/controllers/user/account/form_profile_update.go`, `internal/controllers/website/seo/sitemap_xml_controller.go`
- Check the impersonation signature change: `pkg/useradmin/user_impersonate/impersonate.go`, `pkg/useradmin/user_impersonate/user_impersonate_controller.go`
- Check the environment-aware cookie logic: `internal/controllers/auth/authentication/authentication_controller.go`
- Check the CSRF secret configuration: `internal/config/auth_config.go`, `internal/config/constants.go`
- Check the Dokploy deployment task: `taskfile.yml`
- Check the Go version and dependencies: `go.mod`
- Open an issue on the [Blueprint repository](https://github.com/dracory/blueprint)
