# Upgrade Guide: v0.26.0 to v0.27.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.26.0 to v0.27.0.

## ⚠️ Breaking Changes

### 1. migrate.New() Error Handling
**Change**: The `migrate.New()` function now returns an error that must be handled. This is a breaking change in the `github.com/dracory/migrate` package (v0.2.0 → v0.3.0).

**Old Usage**:
```go
migrator := migrate.New(db, nil)
```

**New Usage**:
```go
migrator, err := migrate.New(db, nil)
if err != nil {
    return fmt.Errorf("failed to create migrator: %w", err)
}
```

**Action Required**:
- Update any custom code that calls `migrate.New()` directly
- Add error handling for the new error return value
- The Blueprint framework's `database/migrations/migrate.go` has been updated automatically
- If you have custom migration code outside of the framework, update it to handle the error

**Files to Check**:
- `database/migrations/migrate.go` (already updated in framework)
- Any custom migration files in your project that use `migrate.New()`

### 2. Go Version Update
**Change**: Minimum Go version updated from 1.26.1 to 1.26.3.

**Action Required**:
- Update your Go installation to version 1.26.3 or later
- Update the `go` directive in your `go.mod` file:
  ```bash
  go mod edit -go=1.26.3
  ```

### 3. Dependency Updates
**Change**: Multiple dependency packages have been updated to newer versions:
- `github.com/dracory/blogstore`: v1.21.0 → v1.22.0
- `github.com/dracory/liveflux`: v0.25.0 → v0.26.0
- `github.com/dracory/migrate`: v0.2.0 → v0.3.0
- `github.com/dracory/sb`: v0.24.0 → v0.26.0
- `github.com/dracory/shopstore`: v1.13.0 → v1.14.0
- Various indirect dependencies (golang.org/x/*, google.golang.org/*, etc.)

**Action Required**:
- Run `go mod tidy` to update your dependencies
- Review the dependency changelogs if you use any of these packages directly in custom code

## 🔄 Migration Steps

### Step 1: Update Go Version
Ensure you have Go 1.26.3 or later installed:
```bash
go version
```

Update the go.mod file:
```bash
go mod edit -go=1.26.3
```

### Step 2: Update Dependencies
Update all dependencies to their latest versions:
```bash
go get -u ./...
go mod tidy
```

### Step 3: Update migrate.New() Calls
Search for any custom usage of `migrate.New()` in your project:
```bash
grep -r "migrate.New" --include="*.go" .
```

Update any found instances to handle the error:
```go
// Old
migrator := migrate.New(db, nil)

// New
migrator, err := migrate.New(db, nil)
if err != nil {
    return fmt.Errorf("failed to create migrator: %w", err)
}
```

### Step 4: Update Version Constant
Update the version constant in `internal/config/version.go`:
```go
const Version = "0.27.0"
```

### Step 5: Build and Test
Build the application to ensure all changes are compatible:
```bash
go build ./...
```

Run tests:
```bash
go test ./...
```

## 🧪 Testing After Migration

### 1. Unit Tests
Run all unit tests to ensure no regressions:
```bash
go test ./...
```

### 2. Integration Tests
Run integration tests if applicable:
```bash
go test -tags=integration ./...
```

### 3. Migration Tests
Test database migrations specifically:
```bash
go test ./database/migrations/...
```

### 4. Build Verification
Verify the application builds successfully:
```bash
go build -o ./bin/server ./cmd/server
```

### 5. Runtime Verification
Start the application and verify:
- Database migrations run successfully
- All services initialize correctly
- No errors related to migrate.New() in logs

## 📝 Additional Notes

### New Features
- No new features introduced in this release
- This is primarily a maintenance release with dependency updates

### Removed Features
- No features removed in this release

### Dependency Security Updates
- Multiple indirect dependencies updated including golang.org/x/* packages
- These updates include security and bug fixes

## 🆘 Common Issues and Solutions

### Issue: migrate.New() compilation error
**Symptom**: Compilation error about too many arguments or missing error handling
**Solution**: Update all `migrate.New()` calls to handle the error return value as shown in Breaking Change #1

### Issue: Go version mismatch
**Symptom**: Error about Go version requirement
**Solution**: Install Go 1.26.3 or later and update the go.mod directive

### Issue: Dependency conflicts
**Symptom**: `go mod tidy` fails with dependency conflicts
**Solution**: Run `go clean -modcache` then retry `go mod tidy`

## 📞 Support

For issues or questions about this upgrade:
- Check the [Blueprint repository](https://github.com/dracory/blueprint)
- Review the changelog for detailed changes
- Open an issue on GitHub for upgrade-specific problems
