# Upgrade Guide: v0.43.0 to v0.44.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.43.0 to v0.44.0.

## Overview

This release makes the in-house **email OTP login the default authentication method**, selectable via a new `config.LOGIN_METHOD` constant. The AuthKnight-specific controllers were renamed to make the dependency explicit, the register controller was split into GET page / POST ajax handlers, a shared `SessionLogin` helper was extracted, a new `EmailOTPTask` background task delivers OTP codes, email allowlist enforcement was consolidated into `EmailAllowedRule`, and an XSS in the OTP return URL was fixed.

**Key Changes:**
- New `config.LOGIN_METHOD` constant (`LOGIN_METHOD_OTP` default, `LOGIN_METHOD_AUTHKNIGHT` alternative) in `internal/config/auth_config.go` — a code-level constant, **not** an env var
- New package `internal/controllers/auth/login_otp/` — self-contained email OTP login (page + ajax handlers, OTP send/verify, rate limiting via cache)
- AuthKnight controllers renamed: `internal/controllers/auth/authentication/` → `authentication_authknight/`, `internal/controllers/auth/login/` → `login_authknight/`
- Auth routes now switch on `config.LOGIN_METHOD`; the AuthKnight callback route (`links.AUTH_AUTH`) is only mounted in AuthKnight mode
- Register controller split: `Handler` → `PageHandler` (GET HTML) + `AjaxHandler` (POST JSON); `form_register.go` deleted; register assets extracted to `app.html`/`app.css`/`app.js`
- New `internal/controllers/auth/shared/session_login.go` — shared `SessionLogin()` helper used by OTP login and registration
- New `internal/tasks/email_otp/` task (`EmailOTPTask`, alias `EmailOTPTask`) + `internal/emails/user_email_otp.go` template
- `EmailAllowlistMiddleware` now delegates to `authrules.NewEmailAllowedRule` instead of inline loop; the rule uses `slices.Contains`
- XSS fix: OTP login `return` query param is validated (`/` prefix, `//` rejected) before being embedded
- Dependency bumps: `neat` v0.50→v0.52, `useradmin` v0.3→v0.4, `userstore` v1.19.0→v1.19.1 pseudo-version, plus indirect bumps (`s3`, `gax-go`, `google.golang.org/api`, `genproto`, `modernc.org/libc`, `uid` v1.10 indirect)

---

## ⚠️ Breaking Changes

---

### 1. Login Method Selection via `config.LOGIN_METHOD`

**Change**: The mounted login controller is now selected by a code constant in `internal/config/auth_config.go`. OTP is the default — Blueprint no longer relies on the external AuthKnight service out of the box.

```go
const (
	LOGIN_METHOD_AUTHKNIGHT = "authknight"
	LOGIN_METHOD_OTP        = "otp"
)

const LOGIN_METHOD = LOGIN_METHOD_OTP
```

**Old Usage** (`internal/controllers/auth/routes.go`):
```go
loginRoute := rtr.NewRoute().
	SetName("Auth > Login Controller").
	SetPath(links.AUTH_LOGIN).
	SetHTMLHandler(login.NewLoginController(application).Handler)

authRoute := rtr.NewRoute().
	SetName("Auth > Auth Controller").
	SetPath(links.AUTH_AUTH).
	SetHTMLHandler(authentication.NewAuthenticationController(application).Handler)
```

**New Usage**:
```go
switch config.LOGIN_METHOD {
case config.LOGIN_METHOD_AUTHKNIGHT:
	authRoutes = append(authRoutes,
		rtr.GetHTML(links.AUTH_LOGIN, login_authknight.NewLoginController(application).Handler),
		rtr.NewRoute().
			SetName("Auth > Auth Controller").
			SetPath(links.AUTH_AUTH).
			SetHTMLHandler(authentication_authknight.NewAuthenticationController(application).Handler))
case config.LOGIN_METHOD_OTP:
	otpController := login_otp.NewLoginController(application)
	authRoutes = append(authRoutes,
		rtr.GetHTML(links.AUTH_LOGIN, otpController.PageHandler),
		rtr.PostJSON(links.AUTH_LOGIN, otpController.AjaxHandler))
default:
	panic("invalid config.LOGIN_METHOD: " + config.LOGIN_METHOD)
}
```

**Action Required**:
- Decide which login method your project uses and set `config.LOGIN_METHOD` accordingly. If your project depends on AuthKnight, set it to `LOGIN_METHOD_AUTHKNIGHT`.
- Update imports in `internal/controllers/auth/routes.go` (see Breaking Change #2 for the renamed packages).
- In OTP mode the `links.AUTH_AUTH` callback route is not mounted — remove any code that depends on it, or switch the constant to AuthKnight.
- `LOGIN_METHOD` is deliberately a constant, not an env var — do not convert it to config-from-env.

---

### 2. AuthKnight Controller Packages Renamed

**Change**: The AuthKnight-specific controllers were renamed to make the external dependency explicit:

| Old package | New package |
|---|---|
| `internal/controllers/auth/authentication` | `internal/controllers/auth/authentication_authknight` |
| `internal/controllers/auth/login` | `internal/controllers/auth/login_authknight` |

**Old Usage**:
```go
import (
	"project/internal/controllers/auth/authentication"
	"project/internal/controllers/auth/login"
)
```

**New Usage**:
```go
import (
	"project/internal/controllers/auth/authentication_authknight"
	"project/internal/controllers/auth/login_authknight"
)
```

**Action Required**:
- Update imports in any file referencing the old paths:
  ```bash
  grep -rn "controllers/auth/authentication\"\|controllers/auth/login\"" --include="*.go" .
  ```
- If your project uses OTP login exclusively, you can delete both AuthKnight packages entirely (in OTP mode they are unreachable).

---

### 3. Register Controller Split into Page + Ajax Handlers

**Change**: `registerController.Handler` was split into `PageHandler` (GET, renders HTML) and `AjaxHandler` (POST, returns JSON). `form_register.go` was removed and the register page assets moved to `app.html`/`app.css`/`app.js` in the same package. Two routes are now registered for `links.AUTH_REGISTER`, both behind the 10-req/min rate limiter.

**Old Usage**:
```go
registerRoute := rtr.NewRoute().
	SetName("Auth > Register Controller").
	SetPath(links.AUTH_REGISTER).
	SetHTMLHandler(register.NewRegisterController(application).Handler)
```

**New Usage**:
```go
registerController := register.NewRegisterController(application)
registerRoute := rtr.GetHTML(links.AUTH_REGISTER, registerController.PageHandler).
	SetName("Auth > Register Controller")
registerAjaxRoute := rtr.PostJSON(links.AUTH_REGISTER, registerController.AjaxHandler).
	SetName("Auth > Register Ajax Controller")

middlewares := []rtr.MiddlewareInterface{rtrMiddleware.RateLimitByIPMiddleware(10, 60)}
registerRoute.AddBeforeMiddlewares(middlewares)
registerAjaxRoute.AddBeforeMiddlewares(middlewares)
routes = append(routes, registerRoute, registerAjaxRoute)
```

**Action Required**:
- If you copied the register controller, apply the same split and register both routes.
- Form submissions must now hit the POST JSON endpoint (`AjaxHandler`) — update any custom registration forms/templates that posted to the old combined handler.
- If you reference `formRegister` helpers from the deleted `form_register.go`, port them into your own code.

---

### 4. New `shared.SessionLogin` Helper

**Change**: Post-authentication session creation, user lookup/creation (including blind-index email lookup), back-URL safety validation, and redirect-URL calculation were extracted into `internal/controllers/auth/shared/session_login.go`:

```go
func SessionLogin(application app.AppInterface, w http.ResponseWriter, r *http.Request,
	email, backUrl string) (redirectURL string, needsRegistration bool, errorMessage string)
```

It is used by both the OTP login controller and the register controller. `isSafeBackURL` rejects non-relative and `//` scheme-relative URLs to prevent open-redirect/XSS.

**Action Required**:
- If your project has its own login/registration flows that duplicate this logic, switch them to `shared.SessionLogin` to inherit the back-URL validation.
- No action needed if you use the built-in controllers.

---

### 5. `EmailAllowlistMiddleware` Delegates to `EmailAllowedRule`

**Change**: The inline allowlist loop in `internal/middlewares/email_allowlist_middleware.go` was replaced with `authrules.NewEmailAllowedRule(app, email)`, so the allowlist check lives in exactly one place. The rule itself now uses `slices.Contains`.

**Old Usage**:
```go
found := false
for _, allowed := range allowedEmails {
	if allowed == email {
		found = true
		break
	}
}
if !found {
	helpers.ToFlashError(...)
	return
}
next.ServeHTTP(w, r)
```

**New Usage**:
```go
if rule := authrules.NewEmailAllowedRule(app, email); rule.Fails() {
	helpers.ToFlashError(app.GetCacheStore(), w, r,
		"Access restricted to authorized emails only", links.Website().Home(), 15)
	return
}
next.ServeHTTP(w, r)
```

**Action Required**:
- If you copied the middleware, adopt the rule-based check (import `authrules "project/internal/rules/auth"`).
- Behavior is unchanged; the OTP login also enforces the same rule before sending codes.

---

### 6. New `EmailOTPTask` Background Task

**Change**: OTP codes are delivered via a queued task. `internal/tasks/constants` gained `EmailOTPTaskAlias = "EmailOTPTask"`, `internal/tasks/email_otp/email_otp_task.go` implements the handler, and `RegisterTasks` registers it. The email template is `internal/emails/user_email_otp.go`.

**Action Required**:
- If you maintain your own `RegisterTasks`, add `email_otp.NewEmailOTPTask(app)` to the task list.
- OTP login requires a working task queue and mailer — verify `taskstore` and mail settings are configured, or switch `LOGIN_METHOD` to AuthKnight.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.44.0"
```

### Step 2: Choose your login method

Set `config.LOGIN_METHOD` in `internal/config/auth_config.go` to `LOGIN_METHOD_OTP` (default, self-contained) or `LOGIN_METHOD_AUTHKNIGHT` (external service).

### Step 3: Update auth routes

Update `internal/controllers/auth/routes.go` to the `switch config.LOGIN_METHOD` structure shown in Breaking Change #1, including the split register routes from Breaking Change #3.

### Step 4: Update renamed imports

```bash
grep -rln "controllers/auth/authentication\"" --include="*.go" . | xargs sed -i 's|controllers/auth/authentication"|controllers/auth/authentication_authknight"|'
grep -rln "controllers/auth/login\"" --include="*.go" . | xargs sed -i 's|controllers/auth/login"|controllers/auth/login_authknight"|'
```

(Review each match — quoted package paths only; adjust package references `authentication.` → `authentication_authknight.` and `login.` → `login_authknight.` accordingly.)

### Step 5: Register the OTP task

Add `email_otp.NewEmailOTPTask(app)` to your `RegisterTasks` list.

### Step 6: Update dependencies

```bash
go mod tidy
go build ./...
```

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...`
3. **OTP flow smoke test**: `task dev`, open the login page, request an OTP code, verify the email is queued/sent and the code logs you in.
4. **AuthKnight mode** (if used): set `config.LOGIN_METHOD = LOGIN_METHOD_AUTHKNIGHT`, rebuild, and confirm the login page redirects through AuthKnight and the `links.AUTH_AUTH` callback still works.
5. **Allowlist**: log in with a non-allowlisted email (when `AUTH_EMAILS_ALLOWED_ACCESS` is set) and confirm the flash error redirect.

---

## 📝 Additional Notes

- **Security fix**: the OTP login `return` query parameter is now validated (`strings.HasPrefix(returnURL, "/") && !strings.HasPrefix(returnURL, "//")`) before being appended to the verify URL — apply the same check to any custom back/return URL handling.
- **Docs**: `docs/authentication-architecture.md` was expanded to document the OTP flow and the login-method decision.
- **No new env vars**: `LOGIN_METHOD` is a code constant by design.
- **Registration flow unchanged externally**: same `links.AUTH_REGISTER` path, now served by two method-specific routes.

---

## 🆘 Common Issues and Solutions

### Issue: `cannot find package .../controllers/auth/login` or `.../authentication`

**Cause**: The AuthKnight packages were renamed in v0.44.0.

**Solution**: Update imports to `login_authknight` / `authentication_authknight` (see Breaking Change #2), or delete references if you use OTP login.

### Issue: OTP code email never arrives

**Cause**: `EmailOTPTask` not registered, task queue not running, or mailer not configured.

**Solution**: Ensure `email_otp.NewEmailOTPTask(app)` is in `RegisterTasks`, the task queue worker is started, and mail settings (SMTP/mail driver env vars) are valid.

### Issue: Login page 404s after upgrade

**Cause**: `links.AUTH_LOGIN` is mounted via `rtr.GetHTML`/`rtr.PostJSON` now; an old catch-all route may be shadowing or the switch has no matching case.

**Solution**: Verify `config.LOGIN_METHOD` is one of the two valid constants — an invalid value panics at startup; ensure the login routes are registered inside the switch.

### Issue: Registration form submits but nothing happens / 405

**Cause**: Form still posts to the old combined handler instead of the POST JSON `AjaxHandler`.

**Solution**: Use the provided `app.js` registration assets or post to `links.AUTH_REGISTER` as JSON and handle the JSON response.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
