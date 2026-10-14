# Upgrade Guide: v0.50.0 to v0.51.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.50.0 to v0.51.0.

## Overview

This release migrates all **short-lived auth credentials** (email OTP codes, magic-link tokens, and password-reset tokens) from the in-process memory cache (`app.GetMemoryCache()`, backed by `ttlcache`) to the **DB-backed cache store** (`app.GetCacheStore()`, table `snv_caches_cache`). Previously these tokens died on every app restart and — more critically — the background email tasks (`EmailOTPTask`, `EmailMagicLinkTask`, `EmailPasswordResetTask`) could never resolve them when the task queue was processed by a separate worker process, because the nonce lookup only worked inside the process that issued it.

The release also adds defensive store-availability checks (blog, geo, subscription) to prevent nil-pointer panics when optional stores are disabled, and updates several dependencies.

**Key Changes:**
- OTP login (`internal/controllers/auth/login_otp/` + `internal/tasks/email_otp/`) stores `otp:<nonce>`, `otp_attempts:<nonce>`, and the `otp_sent:` throttle counter in the cache store
- Magic-link login (`internal/controllers/auth/login_magiclink/` + `internal/tasks/email_magic_link/`) stores `magiclink:<nonce>` and `magiclink:token:<token>` in the cache store
- Password reset (`internal/controllers/auth/forgot_password/`, `internal/controllers/auth/password_reset/`, `internal/tasks/email_password_reset/`) stores `passwordreset:<nonce>` and `passwordreset:token:<token>` in the cache store
- `checkStores` guards in these flows now require `IsEnabledCacheStore()` — the cache store is a hard dependency of all email-token auth flows
- New nil/store checks: blog store in `post_controller`, `blog_post_widget`, `blog_post_list_widget`; geo store in `profile_controller`; country nil check in `stats_controller`; subscription store check in `subscription_middleware`
- Dependency updates: `entitystore` v1.20.1→v1.21.1, `llm` v1.6.0→v1.8.1, `neat` v0.54.0→v0.56.0, `shopadmin` v0.4.0→v0.5.0, `statsstore` v1.28.1→v1.31.1

---

## ⚠️ Breaking Changes

### 1. Cache store is now required for OTP, magic-link, and password-reset flows

**Change**: The auth token flows no longer accept a missing cache store. `checkStores()` returns `500 "Cache store not initialized"` when `app.IsDisabledCacheStore()` is true, and the email tasks fail with "OTP/token not found or expired for the given nonce" when the store is absent.

**Old Usage**:
```go
// v0.50.0 — OTP stored in the in-process ttlcache; dies on restart and is
// invisible to a separate task-worker process.
cache := c.app.GetMemoryCache()
cache.Set("otp:"+nonce, email+":"+otp, otpTTL)
item := cache.Get("otp:" + nonce)
```

**New Usage**:
```go
// v0.51.0 — OTP stored in the DB-backed cache store; survives restarts and
// is resolvable by any process sharing the database.
cache := c.app.GetCacheStore()
err := cache.Set("otp:"+nonce, email+":"+otp, int64(otpTTL.Seconds()))
stored, err := cache.Get("otp:"+nonce, "")
```

**Action Required**:
- Ensure the cache store is enabled: `cacheStoreUsed = true` in `internal/config/stores_config.go` (already `true` by default)
- Ensure the `2026_03_21_0006_store_cache_migrate` migration has run — it creates `snv_caches_cache` via `store.MigrateUp`
- If your project copied or customized `login_otp`, `login_magiclink`, `forgot_password`, `password_reset`, or the three email tasks, port the changes — see Step 3

### 2. `cachestore` API differs from `ttlcache`

**Change**: Any custom code that copied the old pattern must adapt to the `cachestore.StoreInterface` API:

| ttlcache (`GetMemoryCache()`) | cachestore (`GetCacheStore()`) |
|---|---|
| `Set(key, val, time.Duration)` | `Set(key, val string, seconds int64) error` |
| `Get(key)` → `*Item` (`item.Value().(string)`) | `Get(key, default)` → `(string, error)`; returns the default for missing/expired keys |
| `Delete(key)` | `Remove(key) error` |
| `Keys()` | not available — query `snv_caches_cache` directly if needed |
| `GetAndDelete(key)` (atomic) | not available — serialize `Get`+`Remove` with a mutex |

**Action Required**:
- Search your project for `GetMemoryCache()` uses that store auth tokens, send throttles, or anything resolvable across processes, and convert them
- `internal/middlewares/api_auth_middleware.go` and `internal/middlewares/auth_middleware.go` still legitimately use `GetMemoryCache()` as a ~1-minute session/user lookup cache — no change needed there

### 3. Password-reset token consumption is no longer a single atomic op

**Change**: `password_reset.consumeToken` previously used ttlcache's atomic `GetAndDelete`. `cachestore` has no atomic consume, so the read-delete is serialized by a package-level `sync.Mutex` (`consumeMu`). In a **multi-instance deployment** (two app processes sharing the database), a tiny race window remains where two concurrent requests could both read the token before either removes it — single-instance deployments are unaffected.

**Action Required**:
- No code change needed unless you copied `consumeToken`; if you run multiple app instances and need strict single-use, wrap the consume in a DB transaction or a `SELECT ... FOR UPDATE`-style lock

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.51.0"
```

### Step 2: Update dependencies

Apply the `go.mod` bumps, then tidy:

```bash
go get github.com/dracory/entitystore@v1.21.1 \
       github.com/dracory/llm@v1.8.1 \
       github.com/dracory/neat@v0.56.0 \
       github.com/dracory/shopadmin@v0.5.0 \
       github.com/dracory/statsstore@v1.31.1
go mod tidy
```

### Step 3: Update the auth token flows

Copy or port these files — the whole change set is contained in them:

```
internal/controllers/auth/login_otp/login_controller.go
internal/controllers/auth/login_otp/login_controller_test.go
internal/controllers/auth/login_magiclink/login_controller.go
internal/controllers/auth/login_magiclink/login_controller_test.go
internal/controllers/auth/forgot_password/forgot_password_controller.go
internal/controllers/auth/password_reset/password_reset_controller.go
internal/tasks/email_otp/email_otp_task.go
internal/tasks/email_otp/email_otp_task_test.go
internal/tasks/email_magic_link/email_magic_link_task.go
internal/tasks/email_magic_link/email_magic_link_task_test.go
internal/tasks/email_password_reset/email_password_reset_task.go
```

If you customized these files, the mechanical translation is:

```go
// Before
cache := c.app.GetMemoryCache()
cache.Set(key, value, ttl)                    // ttl is time.Duration
item := cache.Get(key)                        // *ttlcache.Item, nil on miss
s, ok := item.Value().(string)
cache.Delete(key)

// After
cache := c.app.GetCacheStore()
err := cache.Set(key, value, int64(ttl.Seconds()))   // seconds int64, returns error
s, err := cache.Get(key, "")                          // "" on miss/expired/error
err := cache.Remove(key)
```

And update the store check in each controller's `checkStores`/`PageHandler`/`handleSubmit`:

```go
// Before
if c.app.GetMemoryCache() == nil {
    c.sendErrorResponse(w, "Memory cache not initialized", http.StatusInternalServerError)
    return false
}

// After
if c.app.IsDisabledCacheStore() {
    c.sendErrorResponse(w, "Cache store not initialized", http.StatusInternalServerError)
    return false
}
```

### Step 4: Copy the defensive store checks (recommended)

These prevent nil-pointer panics when optional stores are disabled:

- `internal/controllers/website/blog/post/post_controller.go` — `IsDisabledBlogStore()` flash-redirect guard
- `internal/widgets/blog_post_widget.go`, `internal/widgets/blog_post_list_widget.go` — `IsDisabledBlogStore()` guards
- `internal/controllers/user/account/profile_controller.go` — `IsDisabledGeoStore()` guard
- `internal/controllers/admin/stats/stats_controller.go` — nil `country` check
- `internal/middlewares/subscription_middleware.go` — `IsDisabledSubscriptionStore()` guard

### Step 5: Verify the build and migrations

```bash
go mod tidy
go build ./...
go test ./internal/controllers/auth/... ./internal/tasks/...
```

Then boot once so `MigrateAll` applies `2026_03_21_0006_store_cache_migrate` if it hasn't run yet — confirm the `snv_caches_cache` table exists.

---

## 🧪 Testing After Migration

1. **Unit tests**: `go test ./internal/controllers/auth/... ./internal/tasks/email_otp ./internal/tasks/email_magic_link` — all OTP/magic-link tests exercise the cache store end-to-end
2. **OTP flow**: request a code, restart the app, then verify — the code must still be accepted (within its 15-minute TTL)
3. **Restart resilience**: request an OTP/magic link, restart the app, then submit — verify succeeds; under v0.50.0 this always failed
4. **Separate worker**: run the task queue in a second process (`go run ./cmd/server task <Alias>` or a dedicated worker) and confirm the OTP/magic-link/reset email is sent — under v0.50.0 the nonce could not be resolved cross-process
5. **Single-use**: replay a consumed magic link or reset token — must fail with "Invalid or expired"

---

## 📝 Additional Notes

- **In-flight tokens are dropped on upgrade.** OTPs/links issued under v0.50.0 lived in process memory and cannot be migrated — affected users simply request a new code/link. No action needed.
- **Throttle counters also moved.** The `otp_sent:` / `magiclink_sent:` / `passwordreset_sent:` counters now survive restarts, so the per-email send quota can no longer be reset by restarting the app — this is a security improvement, not a regression.
- **`internal/cache` singletons are untouched.** `cache.Memory` (ttlcache) and `cache.File` still exist for genuinely process-local caching; only cross-process/shared auth state moved.
- **Expired rows are filtered at read time** (`expires_at > now` in `FindByKey`) and cleaned by the store's `ExpireCacheGoroutine` — no manual cleanup needed.
- **Email task nonce cleanup**: after a successful send, `EmailMagicLinkTask` and `EmailPasswordResetTask` now `Remove` the nonce entry from the cache store; the OTP task leaves it to expire naturally since the verify handler consumes it.

---

## 🆘 Common Issues and Solutions

### Issue: `500 "Cache store not initialized"` on login/forgot-password

**Cause**: `cacheStoreUsed = false` in `internal/config/stores_config.go`, or the store failed to initialize.

**Solution**: Set `cacheStoreUsed = true` and confirm `MigrateAll` runs `2026_03_21_0006_store_cache_migrate` — check for the `snv_caches_cache` table.

### Issue: Task logs "OTP not found or expired for the given nonce" even though the OTP was just sent

**Cause**: The nonce row was written to a different database, or the cache store write failed silently in a copied controller.

**Solution**: Verify both the web process and the worker point at the same database, and that `handleOtpSend` checks the `cache.Set` error return (the new code fails the request with 500 on write errors instead of proceeding).

### Issue: Compile error `cache.GetAndDelete undefined` after copying old `password_reset` code

**Cause**: `cachestore.StoreInterface` has no atomic get-and-delete.

**Solution**: Copy the new `consumeToken` from v0.51.0 — it serializes `Get`+`Remove` behind the package-level `consumeMu` mutex.

### Issue: Test helpers that listed cache keys fail — `cache.Keys undefined`

**Cause**: `ttlcache.Keys()` has no cachestore equivalent.

**Solution**: Query the table directly, e.g. `SELECT cache_key, cache_value FROM snv_caches_cache WHERE cache_key LIKE 'magiclink:token:%'` (see `login_magiclink/login_controller_test.go` `sendLink` for the pattern).

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
