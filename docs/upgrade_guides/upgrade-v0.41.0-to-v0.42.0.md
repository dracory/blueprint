# Upgrade Guide: v0.41.0 to v0.42.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.41.0 to v0.42.0.

## Overview

This release is primarily a maintenance release: it removes the `cmd/snakecase` utility, moves the `ads.txt` endpoint from shared routes into a dedicated SEO controller, bumps the Go toolchain to 1.27, fixes auto-login in the `cmd/ai-browser` sandbox, and updates many dependencies.

**Key Changes:**
- `cmd/snakecase` utility removed entirely (`snakecase.go`, `doc.go`, tests)
- `ads.txt` route moved from `internal/controllers/shared/routes.go` to `internal/controllers/website/seo/routes.go`, now served by `seo.NewAdsTxtController()`
- Google AdSense publisher ID updated from the placeholder to `pub-8821108004642146` / `f08c47fec0942fa0`
- Go version bumped from 1.26 to 1.27 (`go.mod` now requires `go 1.27.0`, Dockerfile uses `golang:1.27`)
- `cmd/ai-browser` fixes: `AUTH_EMAILS_ALLOWED_ACCESS` is now set before config load, and the seeded user gets the fixed ID `ai-browser-user` required by `AiBrowserAutoLoginMiddleware`
- New documentation: `cmd/ai-browser/AGENTS.md` and an "AI Browser Sandbox" section in `AGENTS.md`
- `.gitignore` gains `.playwright-mcp/`
- Major dependency bumps: `dracory/auth` v0.35→v0.36, `dracory/feedstore` v1.5→v1.10.2, `dracory/neat` v0.41→v0.48, `dracory/llm` v1.4→v1.5, `dracory/taskstore` v1.30→v1.32, `dracory/useradmin` v0.2→v0.3, `dracory/userstore` v1.18→v1.19, `dracory/cdn` v1.11→v1.13, `dracory/customstore` v1.14→v1.15, `modernc.org/sqlite` v1.57→v1.59, plus AWS SDK, OpenTelemetry, Google API and `golang.org/x/*` indirect bumps
- `github.com/samber/lo` promoted from the trailing indirect require block to a direct dependency

---

## ⚠️ Breaking Changes

---

### 1. `cmd/snakecase` Utility Removed

**Change**: The `cmd/snakecase` command-line utility (string case conversion helper) has been deleted, including `doc.go`, `snakecase.go`, and `snakecase_test.go`.

**Old Usage**:
```bash
go run ./cmd/snakecase "MyVariableName"
```

**New Usage**: None — the utility no longer exists. Use `github.com/dracory/str` or an equivalent library if you need case conversion in code.

**Action Required**:
- Remove any references to `./cmd/snakecase` in scripts, Taskfiles, CI jobs, or documentation:
  ```bash
  grep -rn "cmd/snakecase" --exclude-dir=.git .
  ```
- If your project relied on it, either copy the deleted code into your own `cmd/` tree or replace calls with `str.SnakeCase(...)`-style helpers.
- Remove the `/cmd/snakecase` entry from your project documentation (e.g. `AGENTS.md` "Project Structure" section) if you copied it.

---

### 2. `ads.txt` Route Moved from Shared Routes to SEO Routes

**Change**: The `/ads.txt` endpoint was removed from `internal/controllers/shared/routes.go` and added to `internal/controllers/website/seo/routes.go`, served by a new dedicated controller `seo.NewAdsTxtController()` (`internal/controllers/website/seo/ads_txt_controller.go`). The returned body changed from the placeholder (`pub-YOURNUMBER`) to `google.com, pub-8821108004642146, DIRECT, f08c47fec0942fa0`.

**Old Usage** (`internal/controllers/shared/routes.go`):
```go
adsTxt := rtr.NewRoute().
	SetName("Shared > ads.txt").
	SetPath("/ads.txt").
	SetStringHandler(func(w http.ResponseWriter, r *http.Request) string {
		return "google.com, pub-8821108004642146, DIRECT, f08c47fec0942fa0"
	})
```

**New Usage** (`internal/controllers/website/seo/routes.go`):
```go
adsRoute := rtr.NewRoute().
	SetName("Website > ads.txt").
	SetPath("/ads.txt").
	SetStringHandler(NewAdsTxtController().Handler)
```

With the controller in `internal/controllers/website/seo/ads_txt_controller.go`:
```go
func (c adsTxtController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "text/plain")
	return "google.com, pub-8821108004642146, DIRECT, f08c47fec0942fa0"
}
```

**Action Required**:
- If your project has its own `ads.txt` route (in shared routes or elsewhere), remove the old copy and keep exactly one `/ads.txt` route — duplicate paths may conflict or shadow each other.
- Move the route to your SEO routes file and extract the handler into an `adsTxtController` if you copied the blueprint structure.
- Update the publisher ID to your own AdSense publisher ID in `ads_txt_controller.go` — the blueprint now ships a real ID, not a placeholder.
- If your project counted shared routes in a test (the blueprint test expects 10 routes now, down from 11), update the expected count.
- If you do NOT use AdSense, you can delete the route and controller entirely.

---

### 3. Go Version Requirement Bumped to 1.27

**Change**: `go.mod` now declares `go 1.27.0` (was `go 1.26.5`), and the Dockerfile uses `golang:1.27` (was `golang:1.26`).

**Old Usage**:
```
go 1.26.5
```
```dockerfile
FROM golang:1.26 AS builder
```

**New Usage**:
```
go 1.27.0
```
```dockerfile
FROM golang:1.27 AS builder
```

**Action Required**:
- Install Go 1.27 or later locally (`go version` to verify).
- Update `go` directives in `go.mod` files of any sub-modules under `/pkg` if present.
- Update CI pipelines, Dockerfiles, and dev-container images that pin a Go version.
- Run `go mod tidy` after switching toolchains.

---

### 4. `cmd/ai-browser` Seeding and Env Ordering

**Change**: Two fixes in the AI browser sandbox:

1. `os.Setenv("AUTH_EMAILS_ALLOWED_ACCESS", email)` is now called **before** `config.NewFromEnv()` in `cmd/ai-browser/main.go`. Previously the env var was never set, so `EmailAllowlistMiddleware` fell back to the hardcoded allowlist (`info@sinevia.com`, `lesichkovm@gmail.com`) and blocked the seeded user.
2. The seeded user is now created with `SetID("ai-browser-user")` (`aiBrowserUserID` constant in `cmd/ai-browser/seed_user.go`), matching the fixed ID looked up by `middlewares.AiBrowserAutoLoginMiddleware`. Without the fixed ID, auto-login silently does nothing.

**Old Usage**:
```go
// env var never set; allowlist fell back to hardcoded emails
user := userstore.NewUser().
	SetEmail(email).
	...
```

**New Usage**:
```go
email := withUser
if email == "" {
	email = aiBrowserDefaultEmail
}
os.Setenv("AUTH_EMAILS_ALLOWED_ACCESS", email)
cfg, err := config.NewFromEnv()
// ...
user := userstore.NewUser().
	SetID(aiBrowserUserID). // "ai-browser-user"
	SetEmail(email).
	...
```

**Action Required**:
- Only applies if your project has a `cmd/ai-browser` (or similar sandbox entry point). If you copied it, apply both fixes.
- Delete any stale `tmp/ai-browser.db` so the user is reseeded with the correct ID.
- Note the known gotcha documented in `cmd/ai-browser/AGENTS.md`: the auto-login cookie is `Secure: true`, so browsers won't persist it over `http://`; the middleware re-mints a session per request so it still works.

---

### 5. Dependency Updates

**Change**: Many dependencies were bumped. The notable direct-dependency changes:

| Package | Old | New |
|---|---|---|
| `github.com/dracory/auth` | v0.35.0 | v0.36.0 |
| `github.com/dracory/cdn` | v1.11.0 | v1.13.0 |
| `github.com/dracory/customstore` | v1.14.0 | v1.15.0 |
| `github.com/dracory/feedstore` | v1.5.0 | v1.10.2 |
| `github.com/dracory/llm` | v1.4.0 | v1.5.0 |
| `github.com/dracory/neat` | v0.41.0 | v0.48.0 |
| `github.com/dracory/taskstore` | v1.30.0 | v1.32.0 |
| `github.com/dracory/useradmin` | v0.2.0 | v0.3.0 |
| `github.com/dracory/userstore` | v1.18.0 | v1.19.0 |
| `github.com/go-sql-driver/mysql` | v1.10.0 | v1.10.1 |
| `github.com/yuin/goldmark` | v1.8.5 | v1.8.6 |
| `modernc.org/sqlite` | v1.57.0 | v1.59.0 |

Indirect bumps include AWS SDK for Go v2, OpenTelemetry (v0.70→v0.71 contrib, v1.45→v1.46), `google.golang.org/api` v0.293→v0.298, `google.golang.org/grpc` v1.83→v1.84, `google.golang.org/genai` v1.69→v1.71, and `golang.org/x/*` packages. `github.com/samber/lo` moved to the direct require block, and the commented-out `replace` directive block in `go.mod` was trimmed to a single `dracory/base` line.

**Action Required**:
- Update your `go.mod`/`go.sum` to match, then run:
  ```bash
  go mod tidy
  go build ./...
  go test ./...
  ```
- The `feedstore` jump (v1.5 → v1.10.2) is the largest — if you call feedstore APIs directly, check its changelog for signature changes between those versions.
- If your `go.mod` carried the old commented-out `replace` block for local dracory module development, note the blueprint now keeps only the `dracory/base` example; re-add your own replace lines locally as needed.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.42.0"
```

### Step 2: Upgrade the Go toolchain

Install Go 1.27+, then update `go.mod` (`go 1.27.0`), your Dockerfile (`golang:1.27`), and CI images.

```bash
go version   # must be >= 1.27
go mod tidy
```

### Step 3: Update dependencies

Bump the dependency versions listed in Breaking Change #5, then:

```bash
go mod tidy
go build ./...
```

### Step 4: Remove `cmd/snakecase` references

```bash
grep -rn "cmd/snakecase" --exclude-dir=.git .
```

Delete the directory if present in your project and clean up references in docs/scripts.

### Step 5: Consolidate the `ads.txt` route

Remove the `ads.txt` route from `internal/controllers/shared/routes.go`, add `internal/controllers/website/seo/ads_txt_controller.go`, and register the route in `internal/controllers/website/seo/routes.go` with your own publisher ID. Update the shared routes count test if you have one (11 → 10).

### Step 6: Apply the ai-browser fixes (if applicable)

If your project includes `cmd/ai-browser`, set `AUTH_EMAILS_ALLOWED_ACCESS` before `config.NewFromEnv()` and seed the user with `SetID("ai-browser-user")`. Delete `tmp/ai-browser.db` to reseed.

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...`
3. **Route check**: `curl -s http://localhost:8080/ads.txt` should return your AdSense line with `Content-Type: text/plain`.
4. **AI browser sandbox**: `go run ./cmd/ai-browser` — it should print a ready-made `authtoken=...` cookie and auto-login should work at `http://127.0.0.1:34756`.

---

## 📝 Additional Notes

- **`.gitignore`**: `.playwright-mcp/` was added for Playwright MCP artifacts.
- **Docs**: `cmd/ai-browser/AGENTS.md` documents the sandbox flags, fixed user ID, and known gotchas; root `AGENTS.md` gained an "AI Browser Sandbox" section.
- **No API changes to `app.AppInterface`**: The `IsEnabled*Store()` / `IsDisabled*Store()` methods introduced in v0.41.0 are unchanged.
- **No config/env var changes** for the main app — only `cmd/ai-browser` sets `AUTH_EMAILS_ALLOWED_ACCESS` itself before config load.

---

## 🆘 Common Issues and Solutions

### Issue: `go.mod requires go >= 1.27.0` build error

**Cause**: Local toolchain is older than 1.27.

**Solution**: Install Go 1.27 or later. With toolchain management enabled, `go` may auto-download it; otherwise install it manually.

### Issue: `/ads.txt` 404 or duplicate route

**Cause**: Route removed from shared routes but not re-added to SEO routes, or defined twice during migration.

**Solution**: Ensure exactly one `/ads.txt` route exists, registered under the website/seo routes. Check that the SEO route group is actually mounted by your website router.

### Issue: ai-browser auto-login not working

**Cause**: Seeded user lacks the fixed ID `ai-browser-user`, or `AUTH_EMAILS_ALLOWED_ACCESS` was not set before config load.

**Solution**: Apply both fixes from Breaking Change #4 and delete `tmp/ai-browser.db` to reseed.

### Issue: "unable to open database file" from ai-browser

**Cause**: `tmp/` was deleted by `task dev` cleanup.

**Solution**: `mkdir -p tmp` and restart.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
