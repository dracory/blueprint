# Proposal: Multiple Login Methods

## Status

Proposed

## Context

`AUTH_LOGIN_METHOD` is currently a single-value switch. `authConfig()` validates it against four methods (`otp`, `magiclink`, `password`, `authknight`) and `loginRoutes()` mounts exactly one controller at `links.AUTH_LOGIN`.

Some projects need more than one login method simultaneously — for example a password login page with a "Use magic link instead" alternative. That UX is not possible today without forking the auth controllers.

## Goal

Allow projects to configure one or more in-house login methods (`otp`, `magiclink`, `password`). The primary method renders at `/auth/login`; secondary methods are offered as alternatives from the same page. All existing single-method behavior remains unchanged.

## Key Design Decisions

1. **Plural env var `AUTH_LOGIN_METHODS`** — comma-separated list (e.g. `password,magiclink`), parsed with `env.GetArrayLower` (env ≥ v1.4.0 splits on `,`/`;`, trims, drops empty entries, and lowercases; also accepts a JSON array). Every entry is validated against the `LOGIN_METHOD_*` constants; unknown values are collected via the `envValidator` and returned as an aggregated error from `NewFromEnv` (see decision #10).
2. **First entry is the primary method** — it renders the main form at `AUTH_LOGIN`. Secondary methods are rendered as "alternative" links/buttons on the login page, matching the "OR — Use magic link instead" pattern.
3. **`AUTH_LOGIN_METHOD` (singular) stays supported** — treated as a one-element list, so existing deployments and `.env` files continue to work. If both vars are set, the plural var wins and a deprecation warning is logged via `slog.Warn`.
4. **Secondary methods get their own sub-paths** — `/auth/login/magiclink`, `/auth/login/otp`, `/auth/login/password`. The primary method always owns the bare `/auth/login` path. This keeps each method's controller, handler signatures, and tests untouched.
5. **Feature-flag checks become membership checks** — `loginMethod == config.LOGIN_METHOD_PASSWORD` becomes `slices.Contains(loginMethods, config.LOGIN_METHOD_PASSWORD)` for:
   - Mounting forgot-password/password-reset routes (`passwordRecoveryRoutes`)
   - Selecting the register controller (see Registration below)
6. **Each method gets a dedicated callback path under `/auth/callback/`** — `/auth/auth` is currently shared ambiguously by `magiclink` and `authknight`, which is unclear and produces a duplicate-path conflict when both are enabled. Move every callback-bearing method to a self-documenting path: `magiclink` → `links.AUTH_CALLBACK_MAGICLINK` (`/auth/callback/magiclink`), `authknight` → `links.AUTH_CALLBACK_AUTHKNIGHT` (`/auth/callback/authknight`), with `links.Auth().AuthKnightLogin` pointing `next_url` at the latter. `otp` and `password` have no callback — nothing changes for them. `links.AUTH_AUTH` is removed (see Breaking Changes): the shared constant was the ambiguity. This is an intentional breaking change — existing AuthKnight app configurations must be updated to `/auth/callback/authknight`, and magic-link emails sent before upgrade will 404 (tokens are short-lived, so impact is minimal).
7. **Callbacks mount once per method, outside the login-page loop** — `/auth/callback/magiclink` mounts iff `slices.Contains(loginMethods, LOGIN_METHOD_MAGICLINK)`; `/auth/callback/authknight` mounts iff `slices.Contains(loginMethods, LOGIN_METHOD_AUTHKNIGHT)`; independent of whether the method is primary or secondary. This removes the need for an authknight exclusivity guard — combining it with other methods is legal in v1. As a secondary method authknight appears on the login page as a "Sign in with AuthKnight" button linking to `/auth/login/authknight`, which redirects to the external service.
8. **Unique route names** — names must be unique across the whole route table. The primary method's GET route keeps the canonical name `"Auth > Login Controller"` so reverse-URL lookups and existing tests keep working; secondary routes get suffixed names (`"Auth > Login Controller (alt)"`, `"Auth > Login MagicLink Ajax Controller (alt)"`).
9. **Shared session login stays unchanged** — all methods converge on `shared/session_login.go` to create the session; multiple methods only change how the user proves identity, not what happens after.
10. **Config errors are collected, not panicked** — `authConfig()` switches from `panic("FATAL: ...")` to the `envValidator` pattern already used by `appConfig(v)`, `databaseConfig(v)`, `storesConfig(v)`, and `llmConfig(v)`: signature becomes `authConfig(v *envValidator)`, validation failures go through `v.Add(...)`, and `NewFromEnv` returns the aggregated `v.Err()` at the end. This reports *all* config problems in one boot failure instead of forcing a fix-restart-repeat cycle. The existing `authConfig()` panics (invalid method, missing CSRF secret) are converted to `v.Add()` as part of this change since the signature must change anyway. The `default:` panic in `mountMethod` stays — that's a programming error (unreachable after config validation), not a config error.

## Proposed Implementation

### Config

#### `internal/config/auth_config.go`

```go
// authConfig reads authentication configuration from environment variables.
// Validation failures are collected in v and reported together by NewFromEnv.
func authConfig(v *envValidator) authSettings {
    // Login Methods
    //
    // List of in-house login mechanisms (comma/semicolon-separated, or a
    // JSON array). The first entry is the primary method rendered at
    // links.AUTH_LOGIN; the rest are offered as alternatives on the login
    // page at /auth/login/<method>.
    // Valid values: otp (default), magiclink, password, authknight.
    // env.GetArrayLower splits on ','/';', trims, drops empties, and lowercases.
    loginMethods := env.GetArrayLower(KEY_AUTH_LOGIN_METHODS)
    if len(loginMethods) == 0 {
        loginMethods = env.GetArrayLower(KEY_AUTH_LOGIN_METHOD) // legacy singular
    } else if strings.TrimSpace(env.GetString(KEY_AUTH_LOGIN_METHOD)) != "" {
        slog.Warn("Both AUTH_LOGIN_METHODS and AUTH_LOGIN_METHOD are set; " +
            "AUTH_LOGIN_METHODS takes precedence and AUTH_LOGIN_METHOD is deprecated.")
    }
    if len(loginMethods) == 0 {
        loginMethods = []string{LOGIN_METHOD_OTP}
    }
    for _, m := range loginMethods {
        switch m {
        case LOGIN_METHOD_OTP, LOGIN_METHOD_MAGICLINK, LOGIN_METHOD_PASSWORD, LOGIN_METHOD_AUTHKNIGHT:
        default:
            v.Add(fmt.Errorf("invalid %s value %q (expected a comma-separated list of: %s, %s, %s, %s)",
                KEY_AUTH_LOGIN_METHODS, m,
                LOGIN_METHOD_OTP, LOGIN_METHOD_MAGICLINK, LOGIN_METHOD_PASSWORD, LOGIN_METHOD_AUTHKNIGHT))
        }
    }
    loginMethods = lo.Uniq(loginMethods)
    // ... CSRF secret etc. — existing panics converted to v.Add() ...
}
```

`authSettings` stores `loginMethods []string` (replacing `loginMethod`). `NewFromEnv` passes the validator: `cfg.setAuthConfig(authConfig(v))` (`config_implementation.go:180`).

#### `internal/config/config_interfaces.go`

```go
GetLoginMethods() []string       // new: all enabled methods, primary first
SetLoginMethods(methods []string) // new: test/override support
GetLoginMethod() string          // kept: returns GetLoginMethods()[0] for compatibility
```

#### `internal/config/constants.go`

Add `KEY_AUTH_LOGIN_METHODS = "AUTH_LOGIN_METHODS"`; keep `KEY_AUTH_LOGIN_METHOD` as deprecated.

### Links

#### `internal/links/constants.go` / `internal/links/auth_links.go`

Add dedicated per-method callback constants; remove the ambiguous shared one:

```go
const AUTH_CALLBACK_MAGICLINK = "/auth/callback/magiclink"   // replaces AUTH_AUTH = "/auth/auth"
const AUTH_CALLBACK_AUTHKNIGHT = "/auth/callback/authknight"
```

```go
// Auth is now unambiguously the magic-link verification callback.
func (l *authLinks) Auth(params ...map[string]string) string {
    p := lo.FirstOr(params, map[string]string{})
    return URL(AUTH_CALLBACK_MAGICLINK, p)
}

func (l *authLinks) AuthKnightLogin(backUrl string) string {
    params := map[string]string{
        "back_url": backUrl,
        "next_url": URL(AUTH_CALLBACK_AUTHKNIGHT, nil), // was l.Auth() → /auth/auth
    }
    return "https://authknight.com/app/login" + query(params)
}
```

`links.AUTH_AUTH` is removed rather than deprecated — it was the ambiguity; keeping it invites new callers to bind to a method-less path again. All existing callers (`login_magiclink` handler, magic-link email task, authknight `next_url`) are updated to the new constants.

### Routes

#### `internal/controllers/auth/routes.go`

```go
func Routes(application app.AppInterface) []rtr.RouteInterface {
    loginMethods := application.GetConfig().GetLoginMethods()

    authRoutes := loginRoutes(application, loginMethods)

    // Verification callbacks mount once per callback-bearing method,
    // independent of primary/secondary position.
    if slices.Contains(loginMethods, config.LOGIN_METHOD_MAGICLINK) {
        authRoutes = append(authRoutes,
            rtr.GetHTML(links.AUTH_CALLBACK_MAGICLINK, login_magiclink.NewLoginController(application).Handler).
                SetName("Auth > MagicLink Controller"))
    }
    if slices.Contains(loginMethods, config.LOGIN_METHOD_AUTHKNIGHT) {
        authRoutes = append(authRoutes,
            rtr.GetHTML(links.AUTH_CALLBACK_AUTHKNIGHT, authentication_authknight.NewAuthenticationController(application).Handler).
                SetName("Auth > AuthKnight Controller"))
    }

    // Password recovery exists whenever password auth is enabled,
    // not only when it is the sole method.
    if slices.Contains(loginMethods, config.LOGIN_METHOD_PASSWORD) {
        authRoutes = append(authRoutes, passwordRecoveryRoutes(application)...)
    }
    // ... rate limiting, logout, register ...
}
```

`loginRoutes` becomes a `mountMethod` helper parameterized on base path. The primary method mounts at `links.AUTH_LOGIN` and keeps the canonical route name `"Auth > Login Controller"`; secondary methods mount at `links.AUTH_LOGIN + "/" + method` with suffixed names:

```go
func loginRoutes(app app.AppInterface, methods []string) []rtr.RouteInterface {
    routes := mountMethod(app, methods[0], links.AUTH_LOGIN, true)
    for _, m := range methods[1:] {
        routes = append(routes, mountMethod(app, m, links.AUTH_LOGIN+"/"+m, false)...)
    }
    return routes
}

// mountMethod mounts GET page + POST ajax routes for one method at basePath.
// primary=true assigns the canonical "Auth > Login Controller" name.
func mountMethod(app app.AppInterface, method, basePath string, primary bool) []rtr.RouteInterface {
    suffix := ""
    if !primary {
        suffix = " (alt)"
    }
    switch method {
    case config.LOGIN_METHOD_MAGICLINK:
        c := login_magiclink.NewLoginController(app)
        return []rtr.RouteInterface{
            rtr.GetHTML(basePath, c.PageHandler).SetName("Auth > Login Controller" + suffix),
            rtr.PostJSON(basePath, c.AjaxHandler).SetName("Auth > Login MagicLink Ajax Controller" + suffix),
        }
    case config.LOGIN_METHOD_OTP:
        // ... same shape
    case config.LOGIN_METHOD_PASSWORD:
        // ... same shape
    case config.LOGIN_METHOD_AUTHKNIGHT:
        c := login_authknight.NewLoginController(app)
        return []rtr.RouteInterface{
            rtr.GetHTML(basePath, c.Handler).SetName("Auth > Login Controller" + suffix),
        }
    }
    panic("invalid login method: " + method)
}
```

### Login page alternatives

Each login page HTML (`app.html` per controller) gets a `{{ alternatives }}` placeholder. Each controller's `PageHandler` receives the enabled-methods list via its constructor (or an options struct) and replaces the placeholder with the output of a shared helper:

#### `internal/controllers/auth/shared/login_alternatives.go`

```go
// LoginAlternatives renders the "OR — use X instead" section as trusted
// HTML to be injected into the login template's {{ alternatives }}
// placeholder. It returns an empty string when no alternatives exist.
// Method names are fixed constants, so the output contains no
// user-controlled input; callers must not pass raw user input here.
func LoginAlternatives(current string, methods []string) string
```

Display names live in one map in the helper:

```go
var loginMethodLabels = map[string]string{
    config.LOGIN_METHOD_MAGICLINK:  "Use magic link instead",
    config.LOGIN_METHOD_OTP:        "Use a one-time code instead",
    config.LOGIN_METHOD_PASSWORD:   "Use password instead",
    config.LOGIN_METHOD_AUTHKNIGHT: "Sign in with AuthKnight",
}
```

Secondary methods link to `/auth/login/<method>`. Since the generated markup is injected into templates as trusted HTML (via `strings.ReplaceAll` on the placeholder), the helper must escape nothing itself but also must never embed dynamic/user values.

### Registration

`registerRoutes` switches on password *membership* rather than equality:

```go
if slices.Contains(loginMethods, config.LOGIN_METHOD_PASSWORD) {
    // register_password controller (creates email+password account)
} else {
    // register controller (post-auth profile completion)
}
```

**Consequence to document:** `AUTH_LOGIN_METHODS=otp,password` produces password-based registration even though OTP is the primary login. This is intentional for v1 — password registration is the only flow that creates credentials — but must be called out in `authentication-architecture.md`. If a project needs primary-method-driven registration, that is a follow-up (see Open Questions).

## Method Interactions

- **Identity is shared.** All methods authenticate the same user record; a user who registered with a password can log in via magic link (and keeps their password for next time).
- **Magic link callback.** `/auth/callback/magiclink` is mounted exactly once if `magiclink` is anywhere in the list.
- **AuthKnight callback.** `/auth/callback/authknight` is mounted exactly once if `authknight` is anywhere in the list. `links.Auth().AuthKnightLogin` generates the new callback URL automatically. Existing AuthKnight-only deployments that previously pinned the callback to `/auth/auth` must update their AuthKnight app configuration to `/auth/callback/authknight` when upgrading.
- **Password recovery.** `/auth/forgot-password` and `/auth/password-reset` are mounted whenever `password` is enabled.
- **Rate limiting.** All login sub-paths inherit the existing 5 req/min IP rate limit. Note: email security scanners that pre-fetch `/auth/auth` links will consume rate-limit quota; the token itself is safe because the magic-link controller validates IP binding before consuming it. If this proves problematic in practice, exempt `AUTH_AUTH` from the auth rate limiter or give it a more permissive bucket — defer to post-implementation observation.

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/config/auth_config.go` | Modify | Signature `authConfig(v *envValidator)`; parse `AUTH_LOGIN_METHODS` list, legacy fallback + deprecation warning, `v.Add` validation; convert existing panics to `v.Add` |
| `internal/config/constants.go` | Modify | Add `KEY_AUTH_LOGIN_METHODS`; deprecate singular key |
| `internal/config/config_interfaces.go` | Modify | Add `GetLoginMethods()` / `SetLoginMethods()` |
| `internal/config/config_implementation.go` | Modify | Pass `v` to `authConfig`; implement new getters; `GetLoginMethod` returns primary |
| `internal/links/constants.go` | Modify | Add `AUTH_CALLBACK_AUTHKNIGHT = "/auth/callback/authknight"` |
| `internal/links/auth_links.go` | Modify | `AuthKnightLogin` points `next_url` at `AUTH_CALLBACK_AUTHKNIGHT` |
| `internal/links/links_extended_test.go` | Modify | Add `AUTH_CALLBACK_AUTHKNIGHT` constant assertion; update `AuthKnightLogin` assertions for the new `/auth/callback/authknight` `next_url` |
| `internal/controllers/auth/routes.go` | Modify | `mountMethod` helper, sub-path mounting, unique route names, per-method callbacks mounted once, membership checks |
| `internal/controllers/auth/shared/login_alternatives.go` | Create | Trusted-HTML alternatives helper + method label map |
| `internal/controllers/auth/login_password/app.html` | Modify | Add `{{ alternatives }}` placeholder |
| `internal/controllers/auth/login_otp/app.html` | Modify | Add `{{ alternatives }}` placeholder |
| `internal/controllers/auth/login_magiclink/app.html` | Modify | Add `{{ alternatives }}` placeholder |
| `internal/controllers/auth/login_password/login_controller.go` | Modify | Accept methods list; render alternatives |
| `internal/controllers/auth/login_otp/login_controller.go` | Modify | Same |
| `internal/controllers/auth/login_magiclink/login_controller.go` | Modify | Same |
| `internal/controllers/auth/routes_test.go` | Modify | Cover multi-method route tables, sub-paths, unique names, dedupe, callback mounting per method combination |
| `.env.example` | Modify | Document `AUTH_LOGIN_METHODS`; mark singular deprecated |
| `docs/authentication-architecture.md` | Modify | Update single-method description and "Removing a Method" section; document registration consequence |
| `docs/environment-variables.md` | Modify | Document new env var |

## Backward Compatibility

- `AUTH_LOGIN_METHOD` continues to work unchanged (one-element list); deprecation warning only fires when both vars are set.
- `GetLoginMethod()` keeps returning a single string (the primary method), so existing callers don't break.
- With a single configured method, routes, URLs, route names, and page output are identical to today (`{{ alternatives }}` is replaced with empty string).
- The only deliberate single-method behavior change is the AuthKnight callback path: AuthKnight-only deployments previously using `/auth/auth` must update their AuthKnight app configuration to `/auth/callback/authknight`.

### Migration Notes

**AuthKnight callback URL.** AuthKnight-only deployments that configured the AuthKnight service to redirect to `/auth/auth` after login must update the callback URL in the AuthKnight app dashboard to `/auth/callback/authknight`. The old path is now reserved exclusively for the magic-link verification callback and is no longer handled by the AuthKnight flow.

## Open Questions

1. **Registration with multiple methods** — should the register page offer a choice ("create account with password" vs "continue with magic link")? v1 keeps the password-present → password-registration rule. If needed later, add `AUTH_REGISTRATION_METHOD` or make it follow the primary method.
2. **Ordering semantics** — "first entry is primary" vs explicit `AUTH_LOGIN_METHOD_PRIMARY`. First-entry is simpler and covers the common case; revisit if projects need primary to differ from list order.

## Testing

- Config: list parsing, dedupe, order preservation, invalid entries surface in `NewFromEnv`'s aggregated error, both-vars deprecation warning, legacy singular fallback, all existing `authConfig` panics now return errors.
- Routes: each secondary method reachable at `/auth/login/<method>`; `/auth/auth` mounted exactly once when `magiclink` present (primary or secondary); `/auth/callback/authknight` mounted exactly once when `authknight` present; `authknight,magiclink` combination produces no path conflicts; password recovery mounted when `password` is secondary; all route names unique; canonical `"Auth > Login Controller"` name preserved for primary.
- Login pages: alternatives section shown iff >1 method; empty placeholder when single method; links point at correct sub-paths.
- End-to-end: `AUTH_LOGIN_METHODS=password,magiclink` — sign in with password at `/auth/login`, sign in via magic link at `/auth/login/magiclink`, both land on the same session via `shared/session_login.go`.

## Verification

- `go test ./internal/config/...` and `go test ./internal/controllers/auth/...` pass
- `AUTH_LOGIN_METHOD=password` produces identical route table and HTML as before
- `AUTH_LOGIN_METHODS=password,magiclink` serves both methods; login page shows "Use magic link instead"
- `AUTH_LOGIN_METHODS=password,authknight` serves both methods; login page shows "Sign in with AuthKnight"; AuthKnight redirect targets `/auth/callback/authknight`
- `AUTH_LOGIN_METHOD=authknight` verifies the callback at `/auth/callback/authknight` (not `/auth/auth` as before)
