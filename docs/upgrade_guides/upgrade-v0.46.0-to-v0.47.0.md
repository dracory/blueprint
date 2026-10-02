# Upgrade Guide: v0.46.0 to v0.47.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.46.0 to v0.47.0.

## Overview

This release moves login-method selection from a compile-time constant to a runtime environment variable, replaces `panic`-based config failures with collected validation errors, and removes the hardcoded default email allowlist. The auth route tree was refactored so each login method mounts only the routes it needs.

**Key Changes:**
- `config.LOGIN_METHOD` constant removed — replaced by the `AUTH_LOGIN_METHOD` environment variable (`otp`, `magiclink`, `password`, `authknight`; default `otp`)
- `AUTH_PASSWORD_AUTH_ENABLED` / `GetPasswordAuthEnabled` / `SetPasswordAuthEnabled` removed — password auth is now activated by `AUTH_LOGIN_METHOD="password"`
- `internal/rules/auth.CanUsePasswordAuthRule` deleted
- `AUTH_EMAILS_ALLOWED_ACCESS` no longer falls back to a hardcoded two-email allowlist — empty now means "all authenticated emails allowed"
- `config.NewFromEnv()` collects validation errors and returns them instead of panicking
- `github.com/dracory/env` upgraded to v1.4.0; allowlist parsing uses `env.GetArray`
- Magic-link login cache entries are stored as JSON instead of a pipe-delimited string (internal change, no action needed)
- `AuthConfigInterface` gains `SetLoginMethod(string)` / `GetLoginMethod() string`

---

## ⚠️ Breaking Changes

---

### 1. `config.LOGIN_METHOD` removed — use `AUTH_LOGIN_METHOD`

**Change**: The login method was a code constant in `internal/config/auth_config.go`. It is now a runtime environment variable read by `authConfig()` and exposed via `ConfigInterface.GetLoginMethod()`.

**Old Usage**:
```go
// internal/config/auth_config.go
const LOGIN_METHOD = LOGIN_METHOD_OTP

// internal/controllers/auth/routes.go
switch config.LOGIN_METHOD { ... }
```

**New Usage**:
```bash
# .env
AUTH_LOGIN_METHOD="otp"   # otp | magiclink | password | authknight
```

```go
loginMethod := application.GetConfig().GetLoginMethod()
switch loginMethod { ... }
```

**Action Required**:
- If you changed `config.LOGIN_METHOD` in your fork, remove the constant and set `AUTH_LOGIN_METHOD` in your environment instead.
- If you referenced `config.LOGIN_METHOD` in your own code, inject the method via `GetLoginMethod()`:
  ```bash
  grep -rn "config\.LOGIN_METHOD\b" --include="*.go" .
  ```
- The `LOGIN_METHOD_OTP` / `LOGIN_METHOD_MAGICLINK` / `LOGIN_METHOD_PASSWORD` / `LOGIN_METHOD_AUTHKNIGHT` constants still exist — they moved to `internal/config/constants.go` and remain the valid values for `AUTH_LOGIN_METHOD`.
- An invalid `AUTH_LOGIN_METHOD` value is now a startup validation error returned by `NewFromEnv()`, not a panic at route registration.

---

### 2. `AUTH_PASSWORD_AUTH_ENABLED` removed

**Change**: The boolean env var `AUTH_PASSWORD_AUTH_ENABLED` and the config accessors `SetPasswordAuthEnabled(bool)` / `GetPasswordAuthEnabled() bool` are gone. Password authentication is enabled by selecting the password login method.

**Old Usage**:
```bash
AUTH_PASSWORD_AUTH_ENABLED="yes"
```

```go
if cfg.GetPasswordAuthEnabled() { ... }
```

**New Usage**:
```bash
AUTH_LOGIN_METHOD="password"
```

```go
if cfg.GetLoginMethod() == config.LOGIN_METHOD_PASSWORD { ... }
```

**Action Required**:
- Remove `AUTH_PASSWORD_AUTH_ENABLED` from your `.env` files and deployment secrets; set `AUTH_LOGIN_METHOD="password"` if you were using password auth.
- Update any code calling `GetPasswordAuthEnabled`/`SetPasswordAuthEnabled` to `GetLoginMethod`/`SetLoginMethod` (`AuthConfigInterface` in `internal/config/config_interfaces.go`).

---

### 3. `CanUsePasswordAuthRule` deleted

**Change**: `internal/rules/auth/can_use_password_auth_rule.go` was removed — password auth is implied by the selected login method, so a separate rule is no longer needed. The password login/forgot/reset controllers dropped their rule checks.

**Action Required**:
- Remove usages if you referenced it:
  ```bash
  grep -rn "CanUsePasswordAuthRule" --include="*.go" .
  ```
- Replace gating logic with `cfg.GetLoginMethod() == config.LOGIN_METHOD_PASSWORD`.

---

### 4. `AUTH_EMAILS_ALLOWED_ACCESS` default allowlist removed

**Change**: When `AUTH_EMAILS_ALLOWED_ACCESS` was unset, Blueprint previously defaulted to a hardcoded list (`info@sinevia.com`, `lesichkovm@gmail.com`). The default is removed — an empty value now means every authenticated email is allowed. Parsing moved to `env.GetArray` (comma/semicolon separated).

**Action Required**:
- If you relied on the implicit default (i.e., never set the variable but expected only those two emails), explicitly set `AUTH_EMAILS_ALLOWED_ACCESS` in your environment.
- If you parse this variable yourself, note the canonical parser is now `env.GetArray` from `github.com/dracory/env` v1.4.0.

---

### 5. `config.NewFromEnv()` returns validation errors instead of panicking

**Change**: Config sections now receive an `envValidator` and collect problems; `NewFromEnv()` returns a combined `error`. Previously a missing `AUTH_CSRF_SECRET` in production/staging (and a failed random-secret generation) called `panic`.

**Old Usage**:
```go
cfg := config.NewFromEnv() // panicked on invalid config
```

**New Usage**:
```go
cfg, err := config.NewFromEnv()
if err != nil {
	log.Fatal(err)
}
```

**Action Required**:
- Ensure all `config.NewFromEnv()` call sites handle the returned error (stock Blueprint already does in `cmd/server`).
- If you call `authConfig()` or other section loaders directly in tests, note the new `authConfig(env *envValidator)` signature — pass a fresh validator and inspect its collected errors.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.47.0"
```

### Step 2: Update environment variables

- Add `AUTH_LOGIN_METHOD` to `.env` / `.env.example` / deployment secrets (`otp`, `magiclink`, `password`, or `authknight`). Omit it to keep the OTP default.
- Delete `AUTH_PASSWORD_AUTH_ENABLED` — if it was `yes`, set `AUTH_LOGIN_METHOD="password"` instead.
- Review `AUTH_EMAILS_ALLOWED_ACCESS`: set it explicitly if the old hardcoded two-email default was load-bearing for you.

### Step 3: Update code references

- `config.LOGIN_METHOD` → `cfg.GetLoginMethod()`
- `GetPasswordAuthEnabled`/`SetPasswordAuthEnabled` → `GetLoginMethod`/`SetLoginMethod`
- `CanUsePasswordAuthRule` → `cfg.GetLoginMethod() == config.LOGIN_METHOD_PASSWORD`

### Step 4: Clean up dependencies

```bash
go mod tidy
go build ./...
```

`github.com/dracory/env` should resolve to v1.4.0+.

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...` (see `internal/controllers/auth/routes_test.go` for per-method route coverage)
3. **Default login**: with `AUTH_LOGIN_METHOD` unset, `/auth/login` serves the OTP flow
4. **Each method**: set `AUTH_LOGIN_METHOD` to `magiclink` / `password` / `authknight` and verify `/auth/login` mounts the right controller — `password` must also mount `/auth/forgot-password` and `/auth/password-reset`
5. **Invalid method**: `AUTH_LOGIN_METHOD="bogus"` must fail startup with a validation error, not a panic
6. **Allowlist**: with `AUTH_EMAILS_ALLOWED_ACCESS` unset, any authenticated email should be allowed

---

## 📝 Additional Notes

- The route tree in `internal/controllers/auth/routes.go` is now split into `loginRoutes`, `passwordRecoveryRoutes`, and `registerRoutes` helpers — if you forked `Routes`, re-apply your changes on the new structure.
- Magic-link tokens stored in the cache are JSON-encoded now; in-flight links issued before the upgrade will not validate — ask users to request a fresh link.
- `docs/proposals/multiple-login-methods.md` describes a proposed `AUTH_LOGIN_METHODS` (plural, concurrent methods) feature — it is a proposal only and not implemented in this release.
- `config.New()` still defaults `loginMethod` to `LOGIN_METHOD_OTP` for programmatic construction.

---

## 🆘 Common Issues and Solutions

### Issue: `undefined: config.LOGIN_METHOD` after upgrade

**Cause**: The constant was removed in favor of the env var.

**Solution**: Read the method from config: `app.GetConfig().GetLoginMethod()` (Breaking Change #1).

### Issue: `cfg.GetPasswordAuthEnabled undefined`

**Cause**: Accessors removed with `AUTH_PASSWORD_AUTH_ENABLED`.

**Solution**: Use `GetLoginMethod()` and compare against `config.LOGIN_METHOD_PASSWORD` (Breaking Change #2).

### Issue: Users outside the old default allowlist can now log in

**Cause**: The hardcoded two-email fallback was removed — empty `AUTH_EMAILS_ALLOWED_ACCESS` allows everyone.

**Solution**: Set `AUTH_EMAILS_ALLOWED_ACCESS` explicitly (Breaking Change #4).

### Issue: App exits with a validation error mentioning `AUTH_LOGIN_METHOD` or `AUTH_CSRF_SECRET`

**Cause**: Config errors are now returned by `NewFromEnv()` and surfaced at startup instead of panicking later.

**Solution**: Fix the reported env var; the error message lists the valid `AUTH_LOGIN_METHOD` values.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
