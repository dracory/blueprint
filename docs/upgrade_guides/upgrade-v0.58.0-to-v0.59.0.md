# Upgrade Guide: v0.58.0 to v0.59.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.58.0 to v0.59.0.

## Overview

This release adds a large suite of unit tests across seeders, controllers, admin adapters, rules, tasks, models, widgets-adjacent helpers, and the `ai-browser` tool — and **replaces the hand-rolled admin media manager with the external `github.com/dracory/fileadmin` package**.

**Key Changes:**
- **Media manager now uses `dracory/fileadmin`**: `internal/controllers/admin/media/media_manager_controller.go` (~1,026 LOC) and its test (~601 LOC) were deleted and replaced by a ~90-line `MediaManagerController` that wraps `fileadmin` with the S3 storage built from `cfg.GetMediaUrl/Region/Key/Secret/Bucket`. `internal/controllers/admin/media/routes.go` now registers a catch-all route (`/admin/media/*`) like the file manager does. The URL `/admin/media` is unchanged, and the UI gains fileadmin's extra features (bulk ops, move, clone, file tree).
- New tests for seeders: `database/seeders/once/settings` (`installed_at` seed), `database/seeders/once/users` (default-user seeding, idempotency), `database/seeders/sync/settings` (drift correction)
- New tests for `internal/rules/post_image.go` (`PostImageURL` fallback/`PostImageURLs` batch resolver added in v0.58.0)
- New tests for `internal/controllers/admin/adapters` (ID-set helpers, customer/geo resolvers, PII seal/unseal passthrough, impersonate/update hooks, LLM factory, flash redirect)
- New tests for `internal/controllers/shared/theme` (theme cookie + open-redirect guard)
- New tests for `internal/controllers/website/shop/cart` (guest cart, add-to-cart validation, quantity limits, authenticated cart)
- New tests for `internal/tasks/constants` (task-alias uniqueness), `internal/cmsblocks` (block-type registration), `database/models` (`Note`), `cmd/ai-browser` (env setup, background process group, user/session seeding)
- No `go.mod`/`go.sum` dependency changes beyond tidying (`dracory/fileadmin` was already a dependency)

---

## ⚠️ Breaking Changes

### 1. `internal/controllers/admin/media` internals removed

The old `mediaManagerController` (unexported) and its AJAX contract (`action=file_rename|file_delete|file_upload|directory_create|directory_delete`, `FileEntry`, `uiLayout`, `MAX_UPLOAD_SIZE`, `JSON_ACTION_*` constants) are gone. `NewMediaManagerController(app)` now returns `*MediaManagerController` with a `Handler(w, r)` method instead of `AnyIndex`.

**Impact**: only if your project customized files inside `internal/controllers/admin/media/` or called the old AJAX endpoint directly (e.g., custom scripts posting `action=...`). The media UI is served entirely by `fileadmin` and is self-contained, so normal admin usage is unaffected.

**Migration**: delete any custom media-controller code and port customizations into the `fileadmin.AdminOptions` of the new wrapper (or fork `fileadmin`).

### 2. Media manager AJAX endpoints changed

`fileadmin` serves its own endpoints under `/admin/media/*` (upload, rename, delete, directory ops, bulk ops, move, clone, `load_files`). If you have external tooling posting to `/admin/media`, update it to `fileadmin`'s contract — see `github.com/dracory/fileadmin` `file_manager/` for the request shapes.

---

## 🔄 Migration Steps

### Step 1: Update dependencies

```bash
go mod tidy
```

### Step 2: Replace the media manager (if you carry Blueprint internals)

If your project mirrors Blueprint's `internal/controllers/admin/media/` directory:

1. Delete `internal/controllers/admin/media/media_manager_controller.go` and `media_manager_controller_test.go`.
2. Copy the new `media_manager_controller.go` and `routes.go` from Blueprint v0.59.0 — the new controller is a ~90-line `fileadmin` wrapper, and routes now include the `/admin/media/*` catch-all.
3. No config or link changes needed — `/admin/media` and the media S3 env vars (`GetMediaUrl/Region/Key/Secret/Bucket/Root`) are used exactly as before.

If your project only mounts the admin routes (no custom media code), no action is needed.

### Step 3 (optional): Port the new tests

If your project mirrors Blueprint's test files, copy the new `*_test.go` files listed in the Overview so your coverage stays in step. They exercise existing v0.58.0 code and should pass unmodified.

---

## 🧪 Testing After Migration

1. **Unit Tests**: `go test ./...` — all suites should pass, including the ~70 new test functions
2. **Media manager**: open `/admin/media` as an admin — the fileadmin UI should load, list your S3 bucket contents, and allow upload/rename/delete/directory ops
3. **Spot check**: `go test ./internal/rules ./internal/controllers/website/shop/cart ./internal/controllers/admin/media ./database/seeders/...` if you ported the tests selectively

---

## 📝 Additional Notes

- If your project has no custom media-manager code, this release is nearly a no-op: bump the version constant and move on — the media page just gets a richer UI.
- The media manager now requires valid S3 media config at request time (same as before — the old controller also built the S3 disk per request and returned the error to the page).
- The new tests document real behavior worth knowing when customizing:
  - `rules.PostImageURLs` returns first media by sequence per post and falls back to `post.GetImageUrl()`
  - The theme controller blocks open redirects (non-relative `return` URLs are rejected)
  - Guest carts enforce quantity limits and reject invalid/missing product payloads
  - Once-seeders are idempotent and no-op on nil app or disabled stores

---

## 🆘 Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| `/admin/media` 404s on sub-actions | The catch-all route is missing — ensure `routes.go` registers both `ADMIN_MEDIA` and `ADMIN_MEDIA + links.CATCHALL` |
| `/admin/media` shows a storage error | `fileadmin` requires a valid `Storage` — check `MEDIA_URL/REGION/KEY/SECRET/BUCKET` env config (same vars the old controller used) |
| Old `action=file_upload` scripts break | The AJAX contract changed to `fileadmin`'s — update callers or drop them (the UI is self-contained) |
| Ported cart/theme test fails | These tests rely on `testutils.Setup()` (v0.58.0+); make sure your `internal/testutils` is current |
| `PostImageURLs` test fails on `EntityIDIn` | Requires `blogstore` v1.39.x or later — already satisfied if you applied the v0.58.0 guide |

## 📞 Support

- Blueprint repository: https://github.com/dracory/blueprint
