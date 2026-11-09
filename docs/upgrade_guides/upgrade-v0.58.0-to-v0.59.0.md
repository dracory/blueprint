# Upgrade Guide: v0.58.0 to v0.59.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.58.0 to v0.59.0.

## Overview

This release is a **test-coverage release** — no production code changed. A large suite of unit tests was added across seeders, controllers, admin adapters, rules, tasks, models, widgets-adjacent helpers, and the `ai-browser` tool. There are no breaking changes, no dependency bumps, and no migration steps.

**Key Changes:**
- New tests for seeders: `database/seeders/once/settings` (`installed_at` seed), `database/seeders/once/users` (default-user seeding, idempotency), `database/seeders/sync/settings` (drift correction)
- New tests for `internal/rules/post_image.go` (`PostImageURL` fallback/`PostImageURLs` batch resolver added in v0.58.0)
- New tests for `internal/controllers/admin/adapters` (ID-set helpers, customer/geo resolvers, PII seal/unseal passthrough, impersonate/update hooks, LLM factory, flash redirect)
- New tests for `internal/controllers/shared/theme` (theme cookie + open-redirect guard)
- New tests for `internal/controllers/website/shop/cart` (guest cart, add-to-cart validation, quantity limits, authenticated cart)
- New tests for `internal/tasks/constants` (task-alias uniqueness), `internal/cmsblocks` (block-type registration), `database/models` (`Note`), `cmd/ai-browser` (env setup, background process group, user/session seeding)
- No `go.mod`/`go.sum` dependency changes beyond tidying

---

## ⚠️ Breaking Changes

**None.** No public API, signature, config, or behavior changes in this release.

---

## 🔄 Migration Steps

### Step 1: Update dependencies

```bash
go mod tidy
```

That is all — upgrading only means pulling the new version (or cherry-picking the test files if you track Blueprint's test suite).

### Step 2 (optional): Port the new tests

If your project mirrors Blueprint's test files, copy the new `*_test.go` files listed in the Overview so your coverage stays in step. They exercise existing v0.58.0 code and should pass unmodified.

---

## 🧪 Testing After Migration

1. **Unit Tests**: `go test ./...` — all suites should pass, including the ~70 new test functions
2. **Spot check**: `go test ./internal/rules ./internal/controllers/website/shop/cart ./database/seeders/...` if you ported the tests selectively

---

## 📝 Additional Notes

- Because no production code changed, a project on v0.58.0 can treat this release as a no-op: bump the version constant and move on.
- The new tests document real behavior worth knowing when customizing:
  - `rules.PostImageURLs` returns first media by sequence per post and falls back to `post.GetImageUrl()`
  - The theme controller blocks open redirects (non-relative `return` URLs are rejected)
  - Guest carts enforce quantity limits and reject invalid/missing product payloads
  - Once-seeders are idempotent and no-op on nil app or disabled stores

---

## 🆘 Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| Ported cart/theme test fails | These tests rely on `testutils.Setup()` (v0.58.0+); make sure your `internal/testutils` is current |
| `PostImageURLs` test fails on `EntityIDIn` | Requires `blogstore` v1.39.x or later — already satisfied if you applied the v0.58.0 guide |

## 📞 Support

- Blueprint repository: https://github.com/dracory/blueprint
