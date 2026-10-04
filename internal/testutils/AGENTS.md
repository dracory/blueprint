# internal/testutils

Test helpers: `Setup()` builds a full app on a unique in-memory SQLite
DB and runs all migrations — use it instead of hand-rolling env config in
tests.

## Rules

- **Always start from `testutils.Setup(...)`** — not `config.NewFromEnv()`
  or `app.New` directly. Pass `WithXStore(true)` options to enable only
  the stores a test needs (all stores are off by default).
- Migrations run inside `Setup` — do not call `migrations.MigrateAll`
  again in tests.
- Seed fixtures live here too (`seed_user.go`, `seed_session.go`,
  `seed_cms.go`, ...). Reuse them; add new fixtures here rather than
  inlining setup in tests.
- `testutils` is for tests only — never import it from non-test code.
- Test style per repo rules: plain test functions, no test-case tables
  unless the neighboring file does it.
