# Upgrade Guide: v0.30.0 to v0.31.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.30.0 to v0.31.0.

## Overview

This release introduces a business rules engine for authentication logic, centralizes task alias constants into a dedicated package, and adds an AI browser testing entry point with auto-login middleware. The register controller has been refactored to use the new rules engine instead of inline validation checks.

**Key Changes:**
- New `internal/rules/auth/` package with composable business rules (`CanRegisterRule`, `CanUsePasswordAuthRule`, `EmailAllowedRule`, `RegisterFormValidationRule`, `UserActiveRule`)
- New `internal/taskconstants/` package centralizing all task alias string constants
- New `AUTH_PASSWORD_AUTH_ENABLED` environment variable and corresponding `AuthConfigInterface` methods
- New `cmd/ai-browser` entry point for automated browser testing with pre-seeded user/session
- New `routes.AiBrowserRouter()` and `middlewares.AiBrowserAutoLoginMiddleware()` for AI browser integration
- New `github.com/dracory/rule` dependency
- Register controller refactored to use business rules instead of inline checks

---

## ⚠️ Breaking Changes

### 1. `AuthConfigInterface` Extended with Password Auth Methods

**Change**: `config.AuthConfigInterface` now requires two new methods: `SetPasswordAuthEnabled(bool)` and `GetPasswordAuthEnabled() bool`. Any custom implementation of this interface must add these methods.

**Old Usage**:
```go
// Custom config implementing AuthConfigInterface (v0.30.0)
type myConfig struct {
    config.ConfigInterface
}

// No password auth methods needed
```

**New Usage**:
```go
// Custom config implementing AuthConfigInterface (v0.31.0)
type myConfig struct {
    config.ConfigInterface
    passwordAuthEnabled bool
}

func (c *myConfig) SetPasswordAuthEnabled(v bool) {
    c.passwordAuthEnabled = v
}

func (c *myConfig) GetPasswordAuthEnabled() bool {
    return c.passwordAuthEnabled
}
```

**Action Required**:
- If you have a custom type that embeds or implements `config.AuthConfigInterface`, add the `SetPasswordAuthEnabled(bool)` and `GetPasswordAuthEnabled() bool` methods.
- If you use the default `configImplementation`, no changes are needed — the methods are already implemented.
- Add `AUTH_PASSWORD_AUTH_ENABLED` to your `.env` file (defaults to `false` if unset; set to `yes` or `true` to enable password authentication).

**Files to Check**:
- `internal/config/config_interfaces.go` (interface definition)
- `internal/config/config_implementation.go` (default implementation)
- `internal/config/auth_config.go` (config loading)
- `internal/config/constants.go` (new `KEY_AUTH_PASSWORD_AUTH_ENABLED` constant)
- Any custom config types implementing `AuthConfigInterface`

---

### 2. Task Alias Strings Replaced with Constants

**Change**: All task handler `Alias()` methods now return constants from the new `internal/taskconstants` package instead of string literals. While the alias values themselves are unchanged, any code that compares task aliases against string literals should be updated to use the constants for compile-time safety.

**Old Usage**:
```go
// v0.30.0 — string literals
func (h *helloWorldTask) Alias() string {
    return "HelloWorldTask"
}

// Comparing against string literal
if task.Alias() == "HelloWorldTask" {
    // ...
}
```

**New Usage**:
```go
// v0.31.0 — constants from taskconstants package
import "project/internal/taskconstants"

func (h *helloWorldTask) Alias() string {
    return taskconstants.HelloWorldTaskAlias
}

// Comparing against constant
if task.Alias() == taskconstants.HelloWorldTaskAlias {
    // ...
}
```

**Available Constants**:
| Constant | Value |
|---|---|
| `BlindIndexRebuildTaskAlias` | `"BlindIndexUpdate"` |
| `CleanUpTaskAlias` | `"CleanUpTask"` |
| `EmailTestTaskAlias` | `"EmailTestTask"` |
| `EmailToAdminTaskAlias` | `"EmailToAdminTask"` |
| `EmailToAdminOnNewContactFormSubmittedTaskAlias` | `"email-to-admin-on-new-contact-form-submitted"` |
| `EmailToAdminOnNewUserRegisteredTaskAlias` | `"email-to-admin-on-new-user-registered"` |
| `HelloWorldTaskAlias` | `"HelloWorldTask"` |
| `StatsVisitorEnhanceTaskAlias` | `"StatsVisitorEnhanceTask"` |

**Action Required**:
- If you have custom task handlers, replace string literals in `Alias()` returns with the corresponding `taskconstants` constant.
- If you compare task aliases against string literals elsewhere, update those comparisons to use the constants.
- If you enqueue tasks by alias using string literals (e.g., `TaskDefinitionEnqueueByAlias("HelloWorldTask")`), switch to the constants (e.g., `TaskDefinitionEnqueueByAlias(taskconstants.HelloWorldTaskAlias)`).

**Files to Check**:
- All files in `internal/tasks/` (already updated in template)
- `pkg/useradmin/user_update/handle_user_update_ajax.go` (already updated in template — `TaskDefinitionEnqueueByAlias` call now uses `taskconstants.BlindIndexRebuildTaskAlias`)
- Any custom task handlers in your application
- Any code that calls `TaskDefinitionEnqueueByAlias()` with string literals

---

### 3. Register Controller Uses Business Rules Engine

**Change**: The register controller (`internal/controllers/auth/register/register_controller.go`) now uses `authrules.NewCanRegisterRule()` and `authrules.NewRegisterFormValidationRule()` instead of inline `if` checks. The validation logic is extracted into composable rule objects in the new `internal/rules/auth/` package.

**Old Usage**:
```go
// v0.30.0 — inline checks
if !controller.app.GetConfig().GetRegistrationEnabled() {
    return helpers.ToFlashError(controller.app.GetCacheStore(), w, r,
        `Registrations are currently disabled`, links.Website().Home(), 10)
}

// ...
if data.firstName == "" {
    data.formErrorMessage = "First name is required field"
    return controller.formRegister(ctx, data).ToHTML()
}
if data.lastName == "" {
    data.formErrorMessage = "Last name is required field"
    return controller.formRegister(ctx, data).ToHTML()
}
```

**New Usage**:
```go
// v0.31.0 — business rules
import authrules "project/internal/rules/auth"

canRegister := authrules.NewCanRegisterRule(controller.app, "")
if canRegister.Fails() {
    return helpers.ToFlashError(controller.app.GetCacheStore(), w, r,
        canRegister.FailMessageFirst(), links.Website().Home(), 10)
}

// ...
formValidation := authrules.NewRegisterFormValidationRule(authrules.RegisterFormData{
    FirstName: data.firstName,
    LastName:  data.lastName,
    Email:     data.email,
    Country:   data.country,
    Timezone:  data.timezone,
})
if formValidation.Fails() {
    data.formErrorMessage = formValidation.Message()
    return controller.formRegister(ctx, data).ToHTML()
}
```

**Action Required**:
- If you have a custom register controller or extended the default one, update it to use the new rules from `internal/rules/auth/`.
- The `RegisterFormValidationRule` also validates password fields (min 8 chars, password confirmation) when a password is provided. If your custom controller had separate password validation, consider delegating to the rule.
- The `CanRegisterRule` now accepts an optional email parameter. Pass `""` to check only registration-enabled status, or pass an email to also validate the allowlist.

**Files to Check**:
- `internal/controllers/auth/register/register_controller.go` (already updated in template)
- Any custom register controllers
- Any code that replicated the registration-enabled or email-allowlist checks inline

---

## 🔄 Migration Steps

### Step 1: Add `AUTH_PASSWORD_AUTH_ENABLED` to Your `.env`

Add the new environment variable to your `.env` file. Set it to `yes` or `true` if your application uses email/password authentication:

```bash
# .env
AUTH_PASSWORD_AUTH_ENABLED=yes
```

If you rely exclusively on external OAuth providers, set it to `no` or leave it unset (defaults to `false`).

### Step 2: Update Custom `AuthConfigInterface` Implementations

If you have a custom type that implements `config.AuthConfigInterface`, add the two new methods:

```go
func (c *myConfig) SetPasswordAuthEnabled(v bool) {
    c.passwordAuthEnabled = v
}

func (c *myConfig) GetPasswordAuthEnabled() bool {
    return c.passwordAuthEnabled
}
```

### Step 3: Update Task Alias References

Replace string literals with `taskconstants` constants in any custom task handlers or enqueue calls:

```bash
# Find files with task alias string literals
grep -rn '"HelloWorldTask"' --include="*.go" .
grep -rn '"CleanUpTask"' --include="*.go" .
grep -rn '"EmailTestTask"' --include="*.go" .
grep -rn '"EmailToAdminTask"' --include="*.go" .
grep -rn '"BlindIndexUpdate"' --include="*.go" .
grep -rn '"email-to-admin-on-new-contact-form-submitted"' --include="*.go" .
grep -rn '"email-to-admin-on-new-user-registered"' --include="*.go" .
grep -rn '"StatsVisitorEnhanceTask"' --include="*.go" .
```

Replace each match with the corresponding `taskconstants` constant and add the import:

```go
import "project/internal/taskconstants"
```

### Step 4: Update Custom Register Controllers (If Applicable)

If you have custom registration logic, refactor it to use the new business rules:

```go
import authrules "project/internal/rules/auth"

// Check if registration is allowed
canRegister := authrules.NewCanRegisterRule(app, email)
if canRegister.Fails() {
    // handle failure with canRegister.FailMessageFirst()
}

// Validate form data
formValidation := authrules.NewRegisterFormValidationRule(authrules.RegisterFormData{
    FirstName:       firstName,
    LastName:        lastName,
    Email:           email,
    Country:         country,
    Timezone:        timezone,
    Password:        password,
    PasswordConfirm: passwordConfirm,
})
if formValidation.Fails() {
    // handle failure with formValidation.Message()
}
```

### Step 5: Run `go mod tidy`

The new `github.com/dracory/rule v0.8.0` dependency has been added. Run tidy to ensure your module graph is clean:

```bash
go mod tidy
```

If you had a local replace directive for `github.com/dracory/rule`, comment it out so `go mod tidy` resolves the published `v0.8.0` version.

---

## 🧪 Testing After Migration

### 1. Unit Tests

Run the full test suite to verify everything compiles and passes:

```bash
go test ./...
```

### 2. Verify Password Auth Configuration

Test that the new `AUTH_PASSWORD_AUTH_ENABLED` setting is correctly loaded:

```bash
# With password auth enabled
AUTH_PASSWORD_AUTH_ENABLED=yes go test ./internal/config/... -run TestAuth

# With password auth disabled
AUTH_PASSWORD_AUTH_ENABLED=no go test ./internal/config/... -run TestAuth
```

### 3. Verify Task Alias Constants

Ensure all task aliases still match their expected values:

```bash
go test ./internal/tasks/... -run TestMetadata
```

### 4. Verify Business Rules

Run the rules engine tests:

```bash
go test ./internal/rules/...
```

### 5. Verify AI Browser (Optional)

If you plan to use the AI browser testing tool:

```bash
go run ./cmd/ai-browser -with-user test@example.com
```

Then navigate to `http://127.0.0.1:34756` to verify auto-login works.

---

## 📝 Additional Notes

### New Features

- **Business Rules Engine** (`internal/rules/auth/`): Five composable rule types built on `github.com/dracory/rule`:
  - `CanRegisterRule` — checks registration enabled + email allowlist
  - `CanUsePasswordAuthRule` — checks if password auth is enabled
  - `EmailAllowedRule` — checks email against allowlist
  - `RegisterFormValidationRule` — validates registration form fields (first name, last name, country, timezone, optional password)
  - `UserActiveRule` — checks if a user account is active

- **Task Constants Package** (`internal/taskconstants/`): Centralized task alias constants for compile-time safety and single source of truth.

- **AI Browser Entry Point** (`cmd/ai-browser/`): A standalone server entry point for automated browser testing. It:
  - Seeds a test user and session automatically
  - Uses `routes.AiBrowserRouter()` which prepends `AiBrowserAutoLoginMiddleware`
  - Runs on a separate port (34756) with hardcoded development environment
  - Refuses to start in production environment
  - Supports `-with-user` and `-admin` CLI flags

- **AI Browser Auto-Login Middleware** (`internal/middlewares/ai_browser_auto_login.go`): Automatically authenticates requests as the `ai-browser-user` dev user. **Must never be used in production** — only via `routes.AiBrowserRouter()`.

- **Password Auth Toggle**: New `AUTH_PASSWORD_AUTH_ENABLED` env var allows disabling email/password authentication when using external OAuth providers exclusively.

### Removed Features

- No features were removed in this release.

---

## 🆘 Common Issues and Solutions

### Issue 1: `AuthConfigInterface` not satisfied

**Symptom**: Compile error like `cannot use myConfig (type *myConfig) as type config.AuthConfigInterface: missing method GetPasswordAuthEnabled`.

**Solution**: Add the `SetPasswordAuthEnabled(bool)` and `GetPasswordAuthEnabled() bool` methods to your custom config type. See Breaking Change #1.

### Issue 2: `taskconstants` package not found

**Symptom**: Compile error like `cannot find package "project/internal/taskconstants"`.

**Solution**: The package is new in v0.31.0. Ensure you've merged the release branch completely. Run `go mod tidy` to refresh the module cache.

### Issue 3: `github.com/dracory/rule` not found

**Symptom**: Compile error like `cannot find module providing package github.com/dracory/rule`.

**Solution**: The `github.com/dracory/rule` dependency was added in this release. Run `go mod tidy` to download it. The published version is `v0.8.0` — ensure your `go.mod` requires `github.com/dracory/rule v0.8.0` (not the placeholder `v0.0.0-00010101000000-000000000000` which results from a local replace directive). If you had a local replace directive (`replace github.com/dracory/rule => ../../_modules_dracory/rule`), comment it out and use the published version instead.

### Issue 4: Register form validation behavior changed

**Symptom**: Registration form now validates password length and confirmation when a password is provided, which wasn't checked before.

**Solution**: The `RegisterFormValidationRule` adds password validation (min 8 chars, must match confirmation) when the `Password` field is non-empty. If you need different password rules, modify the rule in `internal/rules/auth/register_form_validation_rule.go` or create a custom rule.

### Issue 5: AI browser server won't start

**Symptom**: `ERROR: AI browser must not run in production environment. Exiting.`

**Solution**: The AI browser entry point refuses to run when `APP_ENV=production`. Ensure `APP_ENV` is set to `development` or `testing`. The `cmd/ai-browser/main.go` hardcodes `APP_ENV=development` via `setHardcodedEnv()`, so this should only happen if you override it.

---

## 📞 Support

For issues or questions about this upgrade:
- Check the [Blueprint repository](https://github.com/dracory/blueprint)
- Review the proposals in `docs/proposals/` for detailed rationale:
  - `docs/proposals/business-rules-engine.md`
  - `docs/proposals/task-constants-file.md`
  - `docs/proposals/ai-browser.md`
  - `docs/proposals/api-layer-cors-token-auth.md`
- Open an issue on GitHub for upgrade-specific problems

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
- [x] File naming follows pattern: `upgrade-vX.Y.Z-to-vX.Y.Z.md`
- [x] Emoji styling used consistently (⚠️, 🔄, 🧪, 📝, 🆘)
- [x] Previous guides reviewed for consistency
