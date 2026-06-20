# Proposal: Business Rules Engine

## Status

Implemented

## Context

Blueprint's auth controllers currently inline validation logic with scattered `if` checks. For example, the registration controller checks `GetRegistrationEnabled()`, email allowlist, password length, password confirmation, and required fields — all inline in the handler function. This approach has several problems:

1. **Logic duplication** — The same checks are repeated across controllers (register, login, password reset).
2. **Hard to test** — Validation logic is embedded in HTTP handlers, requiring full request/response setup.
3. **No reusable composition** — Cannot combine or reorder checks without modifying handler code.
4. **No error messages** — Inline checks return generic errors; domain-specific messages are scattered.

CourseThread solved this with a `internal/rules/auth/` package using `github.com/dracory/rule` — a composable rules engine where each rule is a self-contained, testable unit with a condition function and optional message.

## Goal

Introduce a `internal/rules/` package to Blueprint that encapsulates business validation logic into composable, testable rule objects — starting with auth rules and extensible to any domain.

## Key Design Decisions

1. **Use `github.com/dracory/rule`** — Already a Dracory ecosystem package; provides `Rule` struct with `SetContext`, `SetCondition`, and `Evaluate` pattern.
2. **Start with auth rules** — The most immediate value is extracting auth validation (registration, login, password).
3. **Each rule is self-contained** — A rule holds its own context, condition function, and message.
4. **Rules are composable** — Multiple rules can be evaluated in sequence; first failure short-circuits.
5. **Rules return messages** — Each rule produces a user-friendly error message on failure.
6. **Rules are testable in isolation** — No HTTP context needed; just construct the rule with data and call `Evaluate()`.

## Proposed Implementation

### New Files

#### `internal/rules/auth/can_register_rule.go`

Checks if registration is enabled and the email is on the allowlist:

```go
type CanRegisterRule struct {
    rule.Rule
}

func NewCanRegisterRule(app app.AppInterface, email string) *CanRegisterRule

// Condition: registration enabled AND email is allowed (if allowlist is non-empty)
// Returns false with no message if registration is disabled
// Returns false if email is not in allowlist (when allowlist is configured)
```

#### `internal/rules/auth/can_register_rule_test.go`

Tests:
- Registration disabled → rule fails
- Registration enabled, no allowlist → rule passes
- Registration enabled, email in allowlist → rule passes
- Registration enabled, email not in allowlist → rule fails

#### `internal/rules/auth/email_allowed_rule.go`

Checks if an email is in the allowed access list:

```go
type EmailAllowedRule struct {
    rule.Rule
}

func NewEmailAllowedRule(app app.AppInterface, email string) *EmailAllowedRule
```

#### `internal/rules/auth/email_allowed_rule_test.go`

Tests:
- Empty allowlist → passes (open access)
- Email in allowlist → passes
- Email not in allowlist → fails

#### `internal/rules/auth/user_active_rule.go`

Checks if a user account is active:

```go
type UserActiveRule struct {
    rule.Rule
}

func NewUserActiveRule(user userstore.UserInterface) *UserActiveRule
```

#### `internal/rules/auth/user_active_rule_test.go`

Tests:
- Active user → passes
- Inactive user → fails
- Nil user → fails

#### `internal/rules/auth/can_use_password_auth_rule.go`

Checks if password authentication is enabled:

```go
type CanUsePasswordAuthRule struct {
    rule.Rule
}

func NewCanUsePasswordAuthRule(app app.AppInterface) *CanUsePasswordAuthRule
```

#### `internal/rules/auth/can_use_password_auth_rule_test.go`

Tests:
- Password auth enabled → passes
- Password auth disabled → fails

#### `internal/rules/auth/register_form_validation_rule.go`

Validates all registration form fields:

```go
type RegisterFormData struct {
    FirstName       string
    LastName        string
    Email           string
    Country         string
    Timezone        string
    Password        string
    PasswordConfirm string
}

type RegisterFormValidationRule struct {
    rule.Rule
    message string
}

func NewRegisterFormValidationRule(data RegisterFormData) *RegisterFormValidationRule
func (r *RegisterFormValidationRule) Message() string
```

Validates:
- First name required
- Last name required
- Email required
- Country required
- Timezone required
- Password required
- Password minimum 8 characters
- Password confirmation matches

#### `internal/rules/auth/register_form_validation_rule_test.go`

Tests:
- All fields valid → passes
- Each missing field → fails with correct message
- Password too short → fails with correct message
- Password mismatch → fails with correct message

### Modified Files

#### `internal/controllers/auth/register/register_controller.go`

Refactor inline validation to use rules:

```go
// Before (inline):
if !app.GetConfig().GetRegistrationEnabled() {
    return errorResponse
}
if len(allowedEmails) > 0 && !slices.Contains(allowedEmails, email) {
    return errorResponse
}
if password != passwordConfirm {
    return errorResponse
}

// After (rules):
canRegister := authRules.NewCanRegisterRule(app, email)
if !canRegister.Evaluate() {
    return errorResponse
}

formValidation := authRules.NewRegisterFormValidationRule(authRules.RegisterFormData{
    FirstName:       firstName,
    LastName:        lastName,
    Email:           email,
    Country:         country,
    Timezone:        timezone,
    Password:        password,
    PasswordConfirm: passwordConfirm,
})
if !formValidation.Evaluate() {
    return errorResponse(formValidation.Message())
}
```

#### `go.mod`

Add dependency (if not already present):

```
github.com/dracory/rule
```

## Rule Composition Pattern

Rules can be composed for complex validation flows:

```go
func ValidateRegistration(app app.AppInterface, data RegisterFormData) error {
    rules := []rule.Rule{
        *authRules.NewCanRegisterRule(app, data.Email),
        *authRules.NewRegisterFormValidationRule(data),
    }

    for _, r := range rules {
        if !r.Evaluate() {
            if msg, ok := r.(interface{ Message() string }); ok {
                return errors.New(msg.Message())
            }
            return errors.New("validation failed")
        }
    }
    return nil
}
```

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/rules/auth/can_register_rule.go` | Create | Registration enabled + email allowlist rule |
| `internal/rules/auth/can_register_rule_test.go` | Create | Tests |
| `internal/rules/auth/email_allowed_rule.go` | Create | Email allowlist check rule |
| `internal/rules/auth/email_allowed_rule_test.go` | Create | Tests |
| `internal/rules/auth/user_active_rule.go` | Create | User active status rule |
| `internal/rules/auth/user_active_rule_test.go` | Create | Tests |
| `internal/rules/auth/can_use_password_auth_rule.go` | Create | Password auth enabled rule |
| `internal/rules/auth/can_use_password_auth_rule_test.go` | Create | Tests |
| `internal/rules/auth/register_form_validation_rule.go` | Create | Full form validation rule |
| `internal/rules/auth/register_form_validation_rule_test.go` | Create | Tests |
| `internal/controllers/auth/register/register_controller.go` | Modify | Refactor to use rules |
| `go.mod` | Modify | Add `dracory/rule` dependency |

## Testing

Each rule is tested in isolation — no HTTP context, no database, no app initialization required (except rules that check config, which use a mock config interface).

- `go test ./internal/rules/auth/...` — all rule tests pass
- `go test ./internal/controllers/auth/...` — existing controller tests still pass after refactor

## Verification

- `go test ./internal/rules/...` passes
- `go test ./internal/controllers/auth/...` passes
- Registration flow behaves identically before and after refactor
- Error messages are preserved or improved

## Benefits

- **Testable** — Each rule is a unit test with no HTTP or database setup
- **Reusable** — Same rules used across controllers (register, login, password reset)
- **Composable** — Rules can be combined and evaluated in sequence
- **Self-documenting** — Rule names and messages clearly express business intent
- **Consistent** — All validation follows the same pattern
- **Extensible** — New rules added without modifying existing controller code

## Future Extensions

- `internal/rules/blog/` — Blog post validation rules
- `internal/rules/shop/` — Order and product validation rules
- `internal/rules/user/` — User profile validation rules
- Rule chains with short-circuit evaluation
- Rule groups with aggregated error messages
