# cmd/deploy

Deployment utility: builds the executable and uploads/replaces it on a
remote server over SSH (`BuildApp` → `UploadFiles` → `UploadExecutable` →
`ReplaceExecutable`).

## Rules

- **This tool has real side effects.** It SSHs into a remote host and
  swaps the running binary. Never run it (`go run ./cmd/deploy`) without
  explicit user instruction — code-review and test changes only.
- SSH credentials come from `constants.go` placeholders / env — never
  commit real hosts, users, or keys.
- `SSH()` and `PrivateKeyPath()` in `functions.go` shell out to the system
  `ssh` binary — keep commands passing through `validateCommand`.
- Mirror the existing `*_test.go` style: plain test functions, no test
  frameworks.

## Run

```bash
go run ./cmd/deploy   # asks for confirmation? No — it deploys immediately.
```
