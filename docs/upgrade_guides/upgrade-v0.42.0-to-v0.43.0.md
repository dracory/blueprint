# Upgrade Guide: v0.42.0 to v0.43.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.42.0 to v0.43.0.

## Overview

This is primarily a maintenance and hardening release. The `github.com/dracory/uid` package was dropped as a direct dependency in favor of `neat.GenerateID()`, error handling was improved across the codebase (unchecked `w.Write`, `os.Setenv`, `r.ParseForm`, `SetMeta`, and similar calls now handle or explicitly discard errors), nil guards were added to app/store accessors and task registration, and several dependencies were bumped.

**Key Changes:**
- `github.com/dracory/uid` removed as a direct dependency — use `neat.GenerateID()` instead of `uid.HumanUid()`
- Nil-safety guards added: `RegisterTasks` returns early when `app == nil`, store setup helpers tolerate a nil config, CMS block `Render` returns an error on nil block
- errcheck cleanup: `_, _ =` discards added to `w.Write`, `os.Setenv`, `r.Body.Close`, `SetJSON`, `SetMeta`, etc.
- Dependency bumps: `cmsstore` v1.41→v1.42, `entitystore` v1.18→v1.20, `neat` v0.48→v0.50, `shopadmin` v0.3→v0.4, `vaultstore` v1.4.1→v1.5, plus indirect bumps (`smithy-go`, `genproto`, `modernc.org/libc`)

---

## ⚠️ Breaking Changes

---

### 1. `github.com/dracory/uid` Replaced by `neat.GenerateID()`

**Change**: The blueprint no longer requires `github.com/dracory/uid` directly (it remains only as an indirect dependency). All internal uses of `uid.HumanUid()` were replaced with `neat.GenerateID()` from `github.com/dracory/neat`, which was already a direct dependency.

**Old Usage**:
```go
import "github.com/dracory/uid"

id := uid.HumanUid()
```

**New Usage**:
```go
import "github.com/dracory/neat"

id := neat.GenerateID()
```

Files changed in the blueprint:
- `internal/helpers/flash.go` — flash message IDs
- `internal/controllers/user/partials/page_header.go` — page header element IDs

**Action Required**:
- Search your project for `dracory/uid` imports:
  ```bash
  grep -rn "dracory/uid" --include="*.go" .
  ```
- Replace `uid.HumanUid()` calls with `neat.GenerateID()` and swap the import.
- Run `go mod tidy` — `dracory/uid` should drop out of your direct requires.
- If you depend on uid-specific formatting (e.g. `uid.HumanUid()` vs other uid generators), verify `neat.GenerateID()` output is acceptable for your use case; otherwise keep `dracory/uid` as a direct dependency in your own `go.mod`.

---

### 2. `RegisterTasks` Now Nil-Safe; Store Setup Tolerates Nil Config

**Change**: `internal/tasks.RegisterTasks(app)` now returns early when `app == nil` in addition to when the task store is disabled. New internal helpers `appDebugEnabled` / `appEnvDevelopment` in `internal/app/datastores.go` guard `app.GetConfig()` before calling `GetAppDebug()` / `IsEnvDevelopment()`, so store constructors no longer panic on a nil config.

**Old Usage**:
```go
func RegisterTasks(app app.AppInterface) {
	if app.IsDisabledTaskStore() { // panics if app is nil
		return
	}
	...
}

st, err := config.NewBlogStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
```

**New Usage**:
```go
func RegisterTasks(app app.AppInterface) {
	if app == nil || app.IsDisabledTaskStore() {
		return
	}
	...
}

st, err := config.NewBlogStore(app.GetDatabase(), appDebugEnabled(app))
```

**Action Required**:
- If your project copied `RegisterTasks` or the `setup*Store` helpers, apply the same nil guards — especially if you construct the app lazily or in tests with partial fixtures.
- No action needed if you call `RegisterTasks` with a fully built app.

---

### 3. CMS Block `Render` Returns Error on Nil Block

**Change**: `BlogPostBlockType.Render`, `BlogPostListBlockType.Render`, and `SearchBlockType.Render` now return `("", fmt.Errorf("block is nil"))` when called with a nil `cmsstore.BlockInterface` instead of panicking. `SaveAdminFields` now propagates `r.ParseForm()` errors and explicitly discards `SetMeta` errors.

**Old Usage**:
```go
html, err := blockType.Render(ctx, block) // panics if block == nil
```

**New Usage**:
```go
html, err := blockType.Render(ctx, block)
if err != nil {
	// err may now be "block is nil" — handle it
}
```

**Action Required**:
- If you call these block types directly with a possibly-nil block, handle the new error return.
- If you implemented your own `cmsstore` block types copied from these, consider adding the same guard for consistency.

---

### 4. Dependency Updates

**Change**: Direct dependency bumps:

| Package | Old | New |
|---|---|---|
| `github.com/dracory/cmsstore` | v1.41.0 | v1.42.0 |
| `github.com/dracory/entitystore` | v1.18.0 | v1.20.0 |
| `github.com/dracory/neat` | v0.48.0 | v0.50.0 |
| `github.com/dracory/shopadmin` | v0.3.0 | v0.4.0 |
| `github.com/dracory/vaultstore` | v1.4.1 | v1.5.0 |
| `github.com/dracory/uid` | v1.9.0 (direct) | removed (indirect only) |

Indirect bumps include `github.com/aws/smithy-go` v1.28.1→v1.28.2, `google.golang.org/genproto` (newer pseudo-version), and `modernc.org/libc` v1.76→v1.77.

**Action Required**:
- Update your `go.mod`/`go.sum` to match, then run:
  ```bash
  go mod tidy
  go build ./...
  go test ./...
  ```
- `entitystore` jumped two minor versions and `shopadmin` one — if you call those APIs directly, check their changelogs for signature changes.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.43.0"
```

### Step 2: Update dependencies

```bash
go mod tidy
go build ./...
```

### Step 3: Replace `dracory/uid` usage

```bash
grep -rn "dracory/uid\|uid\.HumanUid" --include="*.go" .
```

Swap `uid.HumanUid()` → `neat.GenerateID()` and update imports, then `go mod tidy`.

### Step 4: Apply nil-safety fixes (if applicable)

If your project has its own copies of `RegisterTasks`, store `setup*` helpers, or CMS block types, apply the nil guards described in Breaking Changes #2 and #3.

### Step 5: errcheck cleanup (optional but recommended)

The blueprint now explicitly discards errors from `w.Write`, `os.Setenv`, `r.Body.Close`, and similar calls (`_, _ = ...`). If you run `errcheck` or `gocritic` in your project, apply the same pattern to silence warnings.

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...`
3. **Lint** (if configured): `task errcheck`, `task nilaway`, `golangci-lint run`
4. **Smoke test**: `task dev` — verify the server starts and flash messages, cart, and CMS blocks still render.

---

## 📝 Additional Notes

- **No new env vars or config changes** in this release.
- **No API changes to `app.AppInterface`** — only internal hardening.
- **Test-only changes**: many `*_test.go` files gained explicit error discards/nil checks; no behavior change.
- **`cmd/ai-browser`**: `os.Setenv` calls now discard errors explicitly; `findOrCreateSession` guards against a nil session store and checks `http.NewRequestWithContext` errors.

---

## 🆘 Common Issues and Solutions

### Issue: `undefined: uid.HumanUid` or missing `dracory/uid` import after `go mod tidy`

**Cause**: `dracory/uid` is no longer a direct dependency; tidy may have removed it while code still references it.

**Solution**: Replace with `neat.GenerateID()` (see Breaking Change #1), or re-add `github.com/dracory/uid` to your `go.mod` explicitly if you need it.

### Issue: Panic `nil pointer dereference` in store setup during tests

**Cause**: Calling store constructors with an app whose config is nil.

**Solution**: Adopt the `appDebugEnabled`/`appEnvDevelopment` guard pattern from `internal/app/datastores.go`, or ensure test fixtures provide a config.

### Issue: `Render` returns "block is nil" error where it previously "worked"

**Cause**: A nil `cmsstore.BlockInterface` was being passed; it would have panicked before.

**Solution**: Fix the caller to pass a valid block — the error surfaced a latent bug.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
