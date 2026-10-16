# Upgrade Guide: v0.52.0 to v0.53.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.52.0 to v0.53.0.

## Overview

This is a maintenance release. Blueprint's "insecure cookie" behavior — disabling the `Secure` flag on auth cookies and extending the session lifetime — now applies to the `local` environment in addition to `development`. This makes `APP_ENV=local` behave like `development` for plain-HTTP local testing, including the remember-me cookie flow.

**Key Changes:**
- `internal/controllers/auth/shared/session_login.go` — `CreateSession`, `IssueRememberSession`, and `RemoveRememberCookie` now check `IsEnvDevelopment() || IsEnvLocal()` when deciding to drop `Secure` and extend the session lifetime.
- `internal/middlewares/remember_me_middleware.go` — `expireCookie` applies the same `IsEnvDevelopment() || IsEnvLocal()` check.

---

## ⚠️ Breaking Changes

### 1. Auth cookies are no longer `Secure` in the `local` environment

**Change**: When `APP_ENV=local`, Blueprint now sets session, remember-me, and expired cookies **without** the `Secure` flag, matching the existing `development` behavior. Session lifetime in `local` also extends from 2 to 4 hours.

**Old Behavior**: `Secure` was only disabled when `APP_ENV=development`; `local` still emitted `Secure` cookies, which browsers drop over plain HTTP.

**New Behavior**: `local` and `development` both emit non-`Secure` cookies over HTTP.

**Action Required**:
- Do **not** use `APP_ENV=local` for any deployment reachable over the network or serving real users — cookies will be sent over plain HTTP.
- If project-local code relies on `Secure` cookies under `APP_ENV=local`, override the cookie options or use a different environment name.
- No code changes or schema migrations are required.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.53.0"
```

### Step 2: Update dependencies

```bash
go get -u ./...
go mod tidy
```

### Step 3: Update copied insecure-cookie checks

If your project copied Blueprint's auth code, mirror the expanded environment check:

```go
// Before
if application.GetConfig() != nil && application.GetConfig().IsEnvDevelopment() {
    cookieOpts = append(cookieOpts, types.WithSecure(false))
}

// After
if application.GetConfig() != nil &&
    (application.GetConfig().IsEnvDevelopment() || application.GetConfig().IsEnvLocal()) {
    cookieOpts = append(cookieOpts, types.WithSecure(false))
}
```

Apply to any project-local equivalents of:

- `internal/controllers/auth/shared/session_login.go` — `CreateSession` (cookie options and session expiry), `IssueRememberSession`, `RemoveRememberCookie`
- `internal/middlewares/remember_me_middleware.go` — `expireCookie`

### Step 4: Verify

```bash
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

1. **Unit tests**: `go test ./internal/controllers/auth/... ./internal/middlewares/...`
2. **Local login**: with `APP_ENV=local`, log in over `http://localhost` and confirm the session cookie is set (no `Secure` attribute) and persists.
3. **Remember me**: enable remember-me login under `APP_ENV=local` and confirm the remember cookie is stored and auto-login works after the session expires.

---

## 📝 Additional Notes

- **Environment values**: `IsEnvLocal()` checks `APP_ENV=local`; `IsEnvDevelopment()` checks `APP_ENV=development`. Both are defined on `config.ConfigInterface`.
- **Production unaffected**: `production`, `staging`, and `testing` environments still emit `Secure` cookies with the 2-hour session lifetime.
- **No data migration**: only cookie attributes and session expiry behavior changed.

---

## 🆘 Common Issues and Solutions

### Issue: Session cookie still has `Secure` under `APP_ENV=local`

**Cause**: The app was not rebuilt after upgrading, or project-local cookie code does not include the `IsEnvLocal()` check.

**Solution**: Rebuild (`go build ./...`) and update any copied cookie options to include the `IsEnvLocal()` condition.

### Issue: Remember-me auto-login works in `development` but not `local`

**Cause**: `expireCookie` or `IssueRememberSession` in a copied middleware still only checks `IsEnvDevelopment()`.

**Solution**: Add `|| cfg.IsEnvLocal()` to the same condition, matching Blueprint's updated code.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
