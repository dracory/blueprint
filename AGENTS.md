# AGENTS.md

## Project Setup
- Use Go 1.21 or later
- Install dependencies: `go mod download`
- Run the application: `go run main.go`
- Run tests: `go test ./...`

## Code Style
- Follow standard Go formatting with `gofmt`
- Use `camelCase` for variable and function names
- Keep lines under 100 characters
- Document all exported functions and types
- Write table-driven tests for new functionality

## Project Structure
- `/cmd` - Main application entry points
- `/internal` - Private application code
- `/pkg` - Public libraries that can be used by external applications
- `/docs` - Documentation files

## Common Tasks
- To add a new package: Create a new directory under `/pkg` with appropriate `go.mod`
- To run database migrations: `go run cmd/migrate/main.go`
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
