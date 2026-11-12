# AGENTS.md

## Project Setup
- Use Go 1.27 or later
- Install dependencies: `go mod download`
- Run the application: `go run main.go`
- Run tests: `go test ./...`

## Code Style
- Follow standard Go formatting with `gofmt`
- Use `camelCase` for variable and function names
- Keep lines under 100 characters
- Document all exported functions and types
- Write table‑driven tests for new functionality
- liveflux (`internal/controllers/liveflux`) is deprecated — do not use it for
  new code; use the shopadmin-style folder-per-controller pattern with
  `rtr.GetHTML` / `rtr.PostJSON` instead

## AI Company Layer
- `/ai` – Filesystem-based AI company (employees, decisions, SOPs, tasks, memory). See `/ai/AGENTS.md` for operating rules.

## Project Structure
- `/cmd` – Main entry points (CLI, server)
- `/internal` – Private application code (config, controllers, MCP client, background workers)
- `/pkg` – Public libraries (e.g., `agentstore`)
- `/public` – Web UI assets (dashboard, Swagger UI)
- `/docs` – Documentation files

## Common Tasks
- Add a new package: create a new package directory under `/pkg`
- Run linters: `golangci-lint run`

## AI Browser Sandbox
Agents with browser automation access (Playwright / Chrome DevTools MCP) can run a
sandboxed, auto-authenticated instance of the app for UI testing:

```bash
go run ./cmd/ai-browser    # serves http://127.0.0.1:34756, auto-login enabled
```

It uses a disposable SQLite DB (`tmp/ai-browser.db`), disables all external APIs,
and never touches real data. See `cmd/ai-browser/AGENTS.md` for details.

## Testing
- Run all tests: `go test ./...`
- Run tests with coverage: `go test -coverprofile=coverage.out ./...`
- View coverage: `go tool cover -html=coverage.out`
- Run integration tests: `go test -tags=integration ./...`
- Browser e2e (self-contained in `/e2e`): `task e2e` (Playwright) or
  `task e2e:rod` (go-rod); both auto-start `e2e/browser` on :34757. See
  `e2e/README.md`

## Git Workflow
- Branch naming: `feature/your-feature-name` or `bugfix/description`
- Write clear, concise commit messages
- Open a pull request for review before merging to main
