# Upgrade Guide: v0.47.0 to v0.48.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.47.0 to v0.48.0.

## Overview

This release enables multiple concurrent login methods via the new `AUTH_LOGIN_METHODS` environment variable, splits the ambiguous shared `/auth/auth` callback path into per-method flat paths (`/auth/magiclink-callback`, `/auth/authknight-callback`), and adds "use a different method" alternatives on the login pages.

**Key Changes:**
- `AUTH_LOGIN_METHODS` env var added — comma/semicolon-separated (or JSON array) list of login methods; the first entry is the primary method rendered at `/auth/login`, the rest mount at `/auth/<method>-login`
- `AUTH_LOGIN_METHOD` (singular) deprecated — still works as a fallback when `AUTH_LOGIN_METHODS` is unset; a startup warning is logged if both are set
- `links.AUTH_AUTH` (`/auth/auth`) removed — replaced by `links.AUTH_CALLBACK_MAGICLINK` (`/auth/magiclink-callback`) and `links.AUTH_CALLBACK_AUTHKNIGHT` (`/auth/authknight-callback`)
- `AuthConfigInterface` gains `SetLoginMethods([]string)` / `GetLoginMethods() []string`; `GetLoginMethod()` now returns the first entry of the methods list
- New per-method login page links: `AUTH_LOGIN_OTP`, `AUTH_LOGIN_MAGICLINK`, `AUTH_LOGIN_PASSWORD`, `AUTH_LOGIN_AUTHKNIGHT` plus `links.Auth().LoginOtp()` / `LoginMagiclink()` / `LoginPassword()` / `LoginAuthknight()` helpers
- Login pages render an "OR — use another method" section via `shared.LoginAlternatives()` and a new `{{ alternatives }}` template placeholder
- Login controllers (OTP, magic-link, password) gain `SetBasePath()` so secondary mounts post to their own URL
- Forgot/reset password and the password register form now mount whenever `password` is among the enabled methods — not only when it is the primary method

---

## ⚠️ Breaking Changes

---

### 1. `links.AUTH_AUTH` removed — per-method callback paths

**Change**: The shared `/auth/auth` callback path was used by both the magic-link verifier and the AuthKnight `next_url`. It is replaced by two unambiguous constants; the old path is no longer mounted.

**Old Usage**:
```go
// internal/links/constants.go
const AUTH_AUTH = "/auth/auth"

links.Auth().Auth()        // "/auth/auth?..." — shared magic-link callback
// AuthKnight next_url also pointed at AUTH_AUTH
```

**New Usage**:
```go
const AUTH_CALLBACK_MAGICLINK  = "/auth/magiclink-callback"
const AUTH_CALLBACK_AUTHKNIGHT = "/auth/authknight-callback"

links.Auth().Auth()        // now returns the magic-link callback URL
// AuthKnightLogin() next_url uses AUTH_CALLBACK_AUTHKNIGHT
```

**Action Required**:
- Replace any `links.AUTH_AUTH` references in your code:
  ```bash
  grep -rn "AUTH_AUTH\b" --include="*.go" .
  ```
- Update external configuration that referenced `/auth/auth`:
  - AuthKnight app settings: change the `next_url`/redirect to `https://your-domain/auth/authknight-callback`.
  - Any bookmarks, docs, or mail templates pointing at `/auth/auth` for magic links must use `/auth/magiclink-callback`.
- Magic links issued before the upgrade still point at `/auth/auth` and will 404 — users should request a fresh link.

---

### 2. `AUTH_LOGIN_METHOD` deprecated — use `AUTH_LOGIN_METHODS`

**Change**: Login methods are now a list. `AUTH_LOGIN_METHODS` accepts a comma/semicolon-separated list or JSON array (parsed via `env.GetArrayLower`). The first entry is the primary method at `/auth/login`; additional methods are mounted at `/auth/<method>-login` and offered as alternatives on every login page.

**Old Usage**:
```bash
# .env
AUTH_LOGIN_METHOD="otp"
```

**New Usage**:
```bash
# .env — single method (equivalent to before)
AUTH_LOGIN_METHODS="otp"

# .env — multiple concurrent methods (primary first)
AUTH_LOGIN_METHODS="password,magiclink"
```

**Action Required**:
- Rename `AUTH_LOGIN_METHOD` to `AUTH_LOGIN_METHODS` in `.env`, `.env.example`, and deployment secrets. The value stays valid — a single-entry list behaves exactly like the old singular variable.
- The singular variable still works as a fallback when `AUTH_LOGIN_METHODS` is unset, but it is deprecated. If both are set, `AUTH_LOGIN_METHODS` wins and a warning is logged at startup.
- Invalid entries are collected as `NewFromEnv()` validation errors listing the valid values (`otp`, `magiclink`, `password`, `authknight`).

---

### 3. `AuthConfigInterface` — login method is now a list

**Change**: Config stores `loginMethods []string` instead of `loginMethod string`. `AuthConfigInterface` gains `SetLoginMethods([]string)` / `GetLoginMethods() []string`. `SetLoginMethod`/`GetLoginMethod` remain for compatibility — `SetLoginMethod(v)` is equivalent to `SetLoginMethods([]string{v})`, and `GetLoginMethod()` returns the first entry or `""`.

**Old Usage**:
```go
loginMethod := application.GetConfig().GetLoginMethod()
if loginMethod == config.LOGIN_METHOD_PASSWORD { ... }
```

**New Usage**:
```go
loginMethods := application.GetConfig().GetLoginMethods()
if slices.Contains(loginMethods, config.LOGIN_METHOD_PASSWORD) { ... }
// primary method only:
primary := application.GetConfig().GetLoginMethod() // == loginMethods[0]
```

**Action Required**:
- If you implement `AuthConfigInterface` yourself (e.g. a mock in tests), add `SetLoginMethods([]string)` and `GetLoginMethods() []string` — the interface grew.
- Review code comparing `GetLoginMethod()` to a method name: it now only reflects the *primary* method. Use `slices.Contains(cfg.GetLoginMethods(), ...)` when checking whether a method is *enabled*.
- `config.New()` still defaults to OTP (`loginMethods: []string{LOGIN_METHOD_OTP}`).

---

### 4. `{{ alternatives }}` placeholder in login templates

**Change**: The stock `app.html` templates for OTP, magic-link, and password login now include a `{{ alternatives }}` placeholder filled by `shared.LoginAlternatives(current, methods)`. It renders an "OR" divider plus a button per secondary method; empty when fewer than two methods are enabled.

**Action Required**:
- If you maintain custom copies of `internal/controllers/auth/login_*/app.html`, add `{{ alternatives }}` inside the card body where you want the alternative-method buttons to appear.
- If you render the login templates yourself, call `shared.LoginAlternatives(config.LOGIN_METHOD_*, cfg.GetLoginMethods())` and inject the result — it is trusted HTML (fixed labels and constant paths only), safe to inject unescaped.

---

### 5. Login controllers are mount-path aware (`SetBasePath`)

**Change**: `login_otp`, `login_magiclink`, and `login_password` controllers gained a `SetBasePath(string)` method (default `links.AUTH_LOGIN`). The ajax URLs rendered into the page are built from the base path so secondary mounts at `/auth/<method>-login` post back to the correct URL.

**Action Required**:
- If you mount these controllers at a custom path in forked routes, call `controller.SetBasePath(yourPath)` after `NewLoginController(application)`.
- If you forked `internal/controllers/auth/routes.go`, re-apply your changes on the new structure: `loginRoutes` now iterates the methods list and delegates to `mountMethod(application, method, primary)`, and callback routes are mounted once per callback-bearing method via `slices.Contains`.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.48.0"
```

### Step 2: Update environment variables

- Rename `AUTH_LOGIN_METHOD` → `AUTH_LOGIN_METHODS` in `.env` / `.env.example` / deployment secrets. To enable several methods at once, list them comma-separated with the primary first (e.g. `AUTH_LOGIN_METHODS="password,magiclink"`).
- Update your AuthKnight app configuration to redirect to `/auth/authknight-callback` instead of `/auth/auth`.

### Step 3: Update code references

- `links.AUTH_AUTH` → `links.AUTH_CALLBACK_MAGICLINK` or `links.AUTH_CALLBACK_AUTHKNIGHT`
- `cfg.GetLoginMethod() == config.LOGIN_METHOD_X` → `slices.Contains(cfg.GetLoginMethods(), config.LOGIN_METHOD_X)` when testing whether a method is enabled (keep `GetLoginMethod()` for the primary only)
- Custom `AuthConfigInterface` implementations → add `SetLoginMethods`/`GetLoginMethods`
- Custom login `app.html` templates → add the `{{ alternatives }}` placeholder

### Step 4: Clean up dependencies

```bash
go mod tidy
go build ./...
```

`github.com/dracory/env` should resolve to v1.4.0+ (`GetArrayLower` is used for the methods list). `github.com/samber/lo` is used for de-duplication via `lo.Uniq`.

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...` (see `internal/controllers/auth/routes_test.go` for per-method mount coverage and `internal/links/links_extended_test.go` for link helpers)
3. **Single method**: with `AUTH_LOGIN_METHODS="otp"`, `/auth/login` serves OTP and no `/auth/*-login` paths are mounted
4. **Multiple methods**: with `AUTH_LOGIN_METHODS="password,magiclink"`, `/auth/login` serves password login with a "Use magic link instead" button linking to `/auth/magiclink-login`; `/auth/forgot-password` and `/auth/password-reset` are mounted because password is enabled
5. **Callbacks**: `/auth/magiclink-callback?token=...` verifies tokens; `/auth/authknight-callback` handles the AuthKnight return; `/auth/auth` returns 404
6. **Legacy var**: with only `AUTH_LOGIN_METHOD="magiclink"` set, the app still starts and magiclink is the primary method
7. **Invalid entry**: `AUTH_LOGIN_METHODS="otp,bogus"` fails startup with a validation error, not a panic

---

## 📝 Additional Notes

- Secondary mounts use ` (alt)`-suffixed route names (e.g. `Auth > Login Controller (alt)`) so route names stay unique — if you match routes by name in tests or middleware, account for the suffix.
- The primary method is always served at `links.AUTH_LOGIN` (`/auth/login`) regardless of list order beyond position 0; `loginMethodPath()` maps each method to its page path.
- `docs/authentication-architecture.md` and `docs/proposals/multiple-login-methods.md` document the flat `/auth/<method>-login` / `/auth/<method>-callback` convention.
- Duplicates in `AUTH_LOGIN_METHODS` are removed via `lo.Uniq`; values are lowercased before validation.

---

## 🆘 Common Issues and Solutions

### Issue: `undefined: links.AUTH_AUTH` after upgrade

**Cause**: The shared callback constant was split per method.

**Solution**: Use `links.AUTH_CALLBACK_MAGICLINK` or `links.AUTH_CALLBACK_AUTHKNIGHT` (Breaking Change #1).

### Issue: Magic-link emails sent before the upgrade 404

**Cause**: They point at the removed `/auth/auth` path.

**Solution**: Expected — ask users to request a fresh link, which now targets `/auth/magiclink-callback`.

### Issue: AuthKnight login loop or redirect failure

**Cause**: The AuthKnight app still returns users to `/auth/auth`.

**Solution**: Update the AuthKnight `next_url`/redirect configuration to `/auth/authknight-callback` (Breaking Change #1).

### Issue: Mock config no longer satisfies `AuthConfigInterface`

**Cause**: The interface gained `SetLoginMethods`/`GetLoginMethods`.

**Solution**: Implement the two new methods; `SetLoginMethod(v)` can delegate to `SetLoginMethods([]string{v})` (Breaking Change #3).

### Issue: Secondary login page posts to `/auth/login` instead of its own path

**Cause**: The controller was mounted without `SetBasePath`.

**Solution**: Call `controller.SetBasePath(links.AUTH_LOGIN_*)` after construction (Breaking Change #5).

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
