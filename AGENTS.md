# AGENTS.md

## AI Company Structure

### AI Employees
- **Director (MD)** - Strategic decision-making
- **Developer (DEV)** - Platform development and implementation
- **Designer (DES)** - UI/UX design and neo-brutalist branding
- **Quality Assurance Controller (QAC)** - Fact-checking persona biographies (SOP-003)

### AI Company Documentation
- `/ai/` - AI company operational memory and records
- `/ai/index.html` - Mini-site dashboard (browse all AI docs in browser)
- `/ai/styles.css` - Shared neo-brutalist styles for all AI HTML pages
- `/ai/company-profile.html` - Company mission, goals, and constraints

## Project Overview
SwiftAI is a modular AI orchestration platform that enables multi‑provider model management, persona handling, prompt libraries, and MCP tool integration.

## Project Setup
- Use Go 1.25 or later
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

## Git Workflow
- Branch naming: `feature/your-feature-name` or `bugfix/description`
- Write clear, concise commit messages
- Open a pull request for review before merging to main
