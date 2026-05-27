# AGENTS.md

## Project Setup
- Use Go 1.26 or later
- Install dependencies: `go mod download`
- Run the application: `go run ./cmd/server`
- Run tests: `go test ./...`

## Code Style
- Follow standard Go formatting with `gofmt`
- Use `camelCase` for variable and function names
- Keep lines under 100 characters
- Document all exported functions and types
- Write tests for new functionality
- Test cases are difficult to maintain, use only test functions 

## Project Structure
- `/cmd` - Main application entry points
- **cmd/server/** - Main application entry point (`main.go`), background processes, CLI mode
- **cmd/deploy/** - Deployment utility
- **cmd/envenc/** - Environment encryption utility
- **cmd/loadtest/** - Load testing utility
- **cmd/snakecase/** - String case conversion utility
- `/docs` - Documentation files
- `/internal` - Private application code
- **internal/cmsblocks/** - CMS block components
- **internal/config/** - Configuration system with Blueprint-style environment loading
- **internal/registry/** - Registry pattern for service management (database, cache, session, logger, geo, task, asset, user stores)
- **internal/cli/** - CLI dispatcher using github.com/dracory/base/cli
- **internal/routes/** - Router setup using github.com/dracory/rtr
- **internal/controllers/** - HTTP controllers (admin, auth, liveflux, user, webhook, website)
- **internal/middlewares/** - HTTP middlewares
- **internal/emails/** - Email templates and sending logic
- **internal/tasks/** - Background task definitions
- **internal/schedules/** - Scheduled job configuration
- **internal/helpers/** - Utility helpers
- **internal/testutils/** - Test utilities and fixtures
- **internal/resources/** - Static resource handling
- **internal/widgets/** - UI widget components
- **internal/layouts/** - Layout templates
- **internal/links/** - Link management
- **internal/rules/** - Business rule definitions
- **internal/ext/** - External integrations
- `/pkg` - Public libraries that can be used by external applications

## Common Tasks
- To add a new package: Create a new directory under `/pkg` with appropriate `go.mod`
- To run linters: `golangci-lint run`

## Testing
- Run all tests: `go test ./...`
- Run tests with coverage: `go test -coverprofile=coverage.out ./...`
- View coverage: `go tool cover -html=coverage.out`
- Run integration tests: `go test -tags=integration ./...`

## Git Workflow
- Branch naming: `feature/your-feature-name` or `bugfix/description`
- Write clear, concise commit messages
- Open a pull request for review before merging to main
