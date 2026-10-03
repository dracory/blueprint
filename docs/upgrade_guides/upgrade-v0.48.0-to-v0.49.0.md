# Upgrade Guide: v0.48.0 to v0.49.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.48.0 to v0.49.0.

## Overview

This release adds an opt-in "remember me" persistent login: when enabled, login forms show a checkbox that issues a long-lived session (stored as a normal sessionstore row) in a separate HttpOnly `remember_token` cookie. A new `RememberMeMiddleware` runs before `AuthMiddleware` and exchanges a valid remember cookie for a fresh short-lived auth session, rotating the remember session on every use. Logout now deletes the server-side auth and remember sessions in addition to clearing cookies. The stats store initialization gains an explicit `SettingsTableName`, and dependencies were bumped (including `dracory/neat` v0.52.0 → v0.54.0).

**Key Changes:**
- `AUTH_REMEMBER_ME_ENABLED` env var added — shows a "remember me" checkbox on the OTP, magic-link, and password login forms (default: off)
- `AUTH_REMEMBER_ME_DAYS` env var added — lifetime of the remember session and cookie in days (default: 30; must be ≥ 1, otherwise a validation error is collected and the default is used)
- New `config.COOKIE_NAME_REMEMBER_TOKEN` constant (`remember_token`) for the remember cookie
- `AuthConfigInterface` gains `SetRememberMeEnabled(bool)` / `GetRememberMeEnabled() bool` and `SetRememberMeDays(int)` / `GetRememberMeDays() int`
- New `middlewares.RememberMeMiddleware` registered in `globalMiddlewares` before `AuthMiddleware`
- `shared.SessionLogin` / `shared.SessionLoginUser` signatures gained a `rememberMe bool` parameter; session creation was extracted into `shared.CreateSession` and `shared.IssueRememberSession` / `shared.RemoveRememberCookie` were added
- Logout now deletes the server-side auth session and the remember session before clearing cookies
- `statsstore.NewStore` initialization now sets `SettingsTableName: "snv_stats_settings"`
- Dependency upgrades: `dracory/neat` v0.54.0, updated dracory stores (sessionstore, statsstore, userstore, etc.), carbon v2.6.18, sqlite v1.60.1, go-openai v1.43.0

---

## ⚠️ Breaking Changes

---

### 1. `AuthConfigInterface` gained remember-me methods

**Change**: The interface now includes `SetRememberMeEnabled(bool)` / `GetRememberMeEnabled() bool` and `SetRememberMeDays(int)` / `GetRememberMeDays() int`.

**Action Required**:
- If you implement `AuthConfigInterface` yourself (e.g. a mock in tests), add the four new methods. Back them with a `bool` and an `int` field, or delegate to a real `config.New()` instance.

---

### 2. `shared.SessionLogin` / `shared.SessionLoginUser` signatures changed

**Change**: Both functions take a new trailing `rememberMe bool` parameter.

**Old Usage**:
```go
redirectURL, needsRegistration, errorMessage :=
    shared.SessionLogin(application, w, r, email, backUrl)
redirectURL, needsRegistration, errorMessage :=
    shared.SessionLoginUser(application, w, r, user, backUrl)
```

**New Usage**:
```go
redirectURL, needsRegistration, errorMessage :=
    shared.SessionLogin(application, w, r, email, backUrl, rememberMe)
redirectURL, needsRegistration, errorMessage :=
    shared.SessionLoginUser(application, w, r, user, backUrl, rememberMe)
```

**Action Required**:
- Update all call sites — find them with:
  ```bash
  grep -rn "SessionLogin" --include="*.go" .
  ```
- Pass `true` only when the user explicitly opted in (the login controllers read the `remember` form parameter) and `config.GetRememberMeEnabled()` permits it; pass `false` for a standard session-only login.
- A failed remember session does not fail the login — it is logged and the user stays authenticated.

---

### 3. New middleware must run before `AuthMiddleware`

**Change**: `middlewares.RememberMeMiddleware(app)` was added to `globalMiddlewares` immediately before `AuthMiddleware`. It reads the `remember_token` cookie, validates the remember session (`Value == shared.RememberSessionValue`), creates a fresh auth session via `shared.CreateSession`, and rotates the remember session.

**Action Required**:
- If you maintain a forked `internal/routes/global_middlewares.go`, insert `middlewares.RememberMeMiddleware(app)` before `AuthMiddleware` — ordering matters, because the exchange must happen before the auth check runs on the same request.
- The middleware is a no-op when `AUTH_REMEMBER_ME_ENABLED` is unset/off or the session store is disabled, so it is safe to include unconditionally.

---

### 4. Logout clears remember state

**Change**: `logoutController.AnyIndex` now deletes the auth session and the remember session from the session store (remember sessions identified by `shared.RememberSessionValue`) and calls `shared.RemoveRememberCookie` before removing the auth cookie.

**Action Required**:
- If you forked `internal/controllers/auth/logout/logout_controller.go`, re-apply your changes on the new implementation — a logout that only clears cookies leaves the remember session valid for anyone still holding its key.
- If you have custom logout handlers, replicate the pattern: delete the server-side sessions for both `auth.CookieName` and `config.COOKIE_NAME_REMEMBER_TOKEN`, then expire both cookies.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.49.0"
```

### Step 2: Update environment variables (optional feature)

Remember-me is opt-in and off by default. To enable it, add to `.env` / `.env.example` / deployment secrets:

```bash
AUTH_REMEMBER_ME_ENABLED="yes"
AUTH_REMEMBER_ME_DAYS="30"   # optional, default 30
```

`AUTH_REMEMBER_ME_DAYS` must be a positive integer; invalid values produce a startup validation error and fall back to 30.

### Step 3: Update code references

- Custom `AuthConfigInterface` implementations → add `SetRememberMeEnabled`/`GetRememberMeEnabled`/`SetRememberMeDays`/`GetRememberMeDays`
- `shared.SessionLogin(...)` / `shared.SessionLoginUser(...)` call sites → add the `rememberMe` argument
- Forked `globalMiddlewares` → insert `RememberMeMiddleware` before `AuthMiddleware`
- Forked logout controller → delete the remember session and expire `config.COOKIE_NAME_REMEMBER_TOKEN`
- Custom login `app.html`/`app.js` → add a `remember` checkbox field posting to the login endpoint if you want the option when remember-me is enabled

### Step 4: Clean up dependencies

```bash
go mod tidy
go build ./...
```

`github.com/dracory/neat` should resolve to v0.54.0. The stats store now creates/uses a `snv_stats_settings` table — the store's automigration creates it on first run.

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...` (see `internal/middlewares/remember_me_middleware_test.go` for exchange/rotation coverage and `internal/middlewares/api_auth_middleware_test.go` for the updated auth middleware)
3. **Feature off**: with `AUTH_REMEMBER_ME_ENABLED` unset, no remember checkbox renders and no `remember_token` cookie is ever set
4. **Feature on**: with `AUTH_REMEMBER_ME_ENABLED="yes"`, logging in with the checkbox sets `remember_token`; logging in without it does not
5. **Re-authentication**: delete the auth cookie but keep `remember_token`; the next request is authenticated transparently and the remember cookie value changes (rotation)
6. **Logout**: after logout, both cookies are expired and the session rows are gone from the session store — reusing an old `remember_token` does not log in
7. **Validation**: `AUTH_REMEMBER_ME_DAYS="0"` fails startup validation (collected error) and falls back to 30

---

## 📝 Additional Notes

- Remember sessions are ordinary sessionstore rows marked with `shared.RememberSessionValue` — no separate token table is required, and expired rows are cleaned by the normal session store sweep.
- The remember cookie is HttpOnly and carries only the session key; the Secure flag follows the same development/production logic as the auth cookie.
- `shared.CreateSession` is now the single place that builds the short-lived auth session (2h production, 4h development) — the middleware and the login pipeline share it so behavior stays identical.
- See `docs/proposals/remember-me.md` for the full design rationale.

---

## 🆘 Common Issues and Solutions

### Issue: `ConfigInterface` mock no longer compiles

**Cause**: `AuthConfigInterface` gained four remember-me methods.

**Solution**: Implement `SetRememberMeEnabled`/`GetRememberMeEnabled`/`SetRememberMeDays`/`GetRememberMeDays` (Breaking Change #1).

### Issue: `too few arguments in call to shared.SessionLogin`

**Cause**: The signature gained a `rememberMe bool` parameter.

**Solution**: Pass `false` to preserve the old session-only behavior, or pass the user's checkbox value (Breaking Change #2).

### Issue: Remembered user is not re-authenticated

**Cause**: `RememberMeMiddleware` is missing or runs after `AuthMiddleware` in a forked middleware chain.

**Solution**: Register `middlewares.RememberMeMiddleware(app)` immediately before `AuthMiddleware` (Breaking Change #3).

### Issue: User stays logged in after logout on another device

**Cause**: A custom logout only clears cookies; the remember session row remains valid.

**Solution**: Delete both session rows and expire `config.COOKIE_NAME_REMEMBER_TOKEN` (Breaking Change #4).

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
