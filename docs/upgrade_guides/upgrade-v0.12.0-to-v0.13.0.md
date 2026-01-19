# Upgrade Guide: v0.12.0 to v0.13.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.12.0 to v0.13.0.

## 🚨 Breaking Changes

### 1. Variable Name Changes: `app` → `registry`
**Change**: All instances of `app` variable have been renamed to `registry` for better semantic clarity and consistency with the registry pattern.

**Old Usage**:
```go
func NewIndexNowController(app registry.RegistryInterface) *indexNowController {
    return &indexNowController{
        app: app,
    }
}

func (controller *indexNowController) Handler(w http.ResponseWriter, r *http.Request) string {
    if controller.app.GetConfig().GetCmsStoreUsed() {
        return layouts.NewCmsLayout(
            controller.app,
            r,
            options).ToHTML()
    }
}
```

**New Usage**:
```go
func NewIndexNowController(registry registry.RegistryInterface) *indexNowController {
    return &indexNowController{
        app: registry,  // Field name remains 'app' for backward compatibility
    }
}

func (controller *indexNowController) Handler(w http.ResponseWriter, r *http.Request) string {
    if controller.app.GetConfig().GetCmsStoreUsed() {
        return layouts.NewCmsLayout(
            controller.app,
            r,
            options).ToHTML()
    }
}
```

**Action Required**:
- Update all variable names from `app` to `registry` in function parameters and local variables
- Update method calls from `app.GetX()` to `registry.GetX()`
- **Note**: Struct field names remain `app` for backward compatibility in existing controllers

### 2. New Security Middlewares Added to Global Stack
**Change**: HTTPS redirect and security headers middlewares are now automatically enabled in the global middleware stack.

**Old Usage**:
```go
// Previous global middleware stack did not include HTTPS redirect or security headers
func globalMiddlewares(registry registry.RegistryInterface) []rtr.MiddlewareInterface {
    globalMiddlewares := []rtr.MiddlewareInterface{
        // ... existing middlewares
        middlewares.LogRequestMiddleware(registry),
        middlewares.ThemeMiddleware(),
        middlewares.AuthMiddleware(registry),
    }
    return globalMiddlewares
}
```

**New Usage**:
```go
// New global middleware stack includes security middlewares
func globalMiddlewares(registry registry.RegistryInterface) []rtr.MiddlewareInterface {
    globalMiddlewares := []rtr.MiddlewareInterface{
        // ... existing middlewares
        rtrMiddleware.TimeoutMiddleware(30 * time.Second),
        rtrMiddleware.RateLimitByIPMiddleware(20, 1),
        rtrMiddleware.RateLimitByIPMiddleware(180, 1*60),
        rtrMiddleware.RateLimitByIPMiddleware(12000, 60*60),
    }

    // Add HTTPS redirect middleware only in production (not in development or testing)
    if registry.GetConfig() != nil &&
        !registry.GetConfig().IsEnvTesting() &&
        !registry.GetConfig().IsEnvDevelopment() {
        globalMiddlewares = append(globalMiddlewares,
            middlewares.NewHTTPSRedirectMiddleware(),
        )
    }

    globalMiddlewares = append(globalMiddlewares,
        middlewares.LogRequestMiddleware(registry),
        middlewares.NewSecurityHeadersMiddleware(),  // NEW
        middlewares.ThemeMiddleware(),
        middlewares.AuthMiddleware(registry),
    )

    return globalMiddlewares
}
```

**Action Required**:
- Review your custom middleware configurations
- Test HTTPS behavior in production environments
- Update CSP policies if you have custom content sources
- **Files to check**: `internal/routes/global_middlewares.go`

### 3. Dependency Updates
**Change**: Several key dependencies have been updated with breaking changes.

**Old Dependencies**:
```go
modernc.org/sqlite v1.42.2
github.com/dracory/blogstore v1.4.0
github.com/dracory/cmsstore v1.4.0
```

**New Dependencies**:
```go
modernc.org/sqlite v1.43.0
github.com/dracory/blogstore v1.4.1
github.com/dracory/cmsstore v1.5.0
```

**Action Required**:
- Run `go mod tidy` to update dependencies
- Test database operations with new SQLite version
- Review blog and CMS store API changes if any

### 4. Task Queue API Changes
**Change**: Task queue initialization has been updated with new configuration options.

**Old Usage**:
```go
// Previous task queue runner initialization
runner := taskstore.NewTaskQueueRunner(ts, 2, 10)
runner.Start(ctx)
```

**New Usage**:
```go
// New task queue runner with structured options
runner := taskstore.NewTaskQueueRunner(ts, taskstore.TaskQueueRunnerOptions{
    IntervalSeconds: 2,
    UnstuckMinutes:  2,
    MaxConcurrency:  10,
    Logger:          log.Default(),
})
runner.Start(ctx)
```

**Action Required**:
- Update task queue runner initialization in main.go
- **File to update**: `cmd/server/main.go` line ~185

### 5. Enhanced Error Handling in Background Processes
**Change**: Added comprehensive nil checks and error handling in background process initialization.

**Old Usage**:
```go
func startBackgroundProcesses(group *backgroundGroup, registry registry.RegistryInterface) error {
    // Direct usage without nil checks
    ts := registry.GetTaskStore()
    if ts != nil {
        // ... task queue setup
    }
}
```

**New Usage**:
```go
func startBackgroundProcesses(ctx context.Context, group *backgroundGroup, registry registry.RegistryInterface) error {
    if registry == nil {
        return errors.New("startBackgroundProcesses called with nil registry")
    }
    
    if registry.GetConfig() == nil {
        return errors.New("startBackgroundProcesses called with nil config")
    }

    if registry.GetDatabase() == nil {
        return errors.New("startBackgroundProcesses called with nil db")
    }

    // Additional nil checks for each store...
    if registry.GetConfig().GetTaskStoreUsed() && registry.GetTaskStore() == nil {
        return errors.New("startBackgroundProcesses task store is enabled but not initialized")
    }
    // ... rest of initialization
}
```

**Action Required**:
- Review error handling in your custom background processes
- **File to update**: `cmd/server/main.go` lines ~151-177

## 🔄 Migration Steps

### Step 1: Update Dependencies
```bash
# Update Go modules
go mod tidy
go mod download

# Verify dependency versions
go list -m all | grep -E "(sqlite|blogstore|cmsstore)"
```

### Step 2: Update Variable Names
```bash
# Find all instances of 'app' variable usage
find . -name "*.go" -not -path "./vendor/*" -exec grep -l "app\." {} \;

# Manual replacement required for semantic clarity
# Replace function parameters and local variables:
# app -> registry
# app.GetX() -> registry.GetX()
```

### Step 3: Update Task Queue Initialization
```bash
# Update cmd/server/main.go task queue initialization
# Replace old TaskQueueRunner call with new options structure
```

### Step 4: Review Security Middleware Configuration
```bash
# Check if your application needs custom CSP policies
# Review HTTPS redirect behavior in development vs production
# Test security headers in browser dev tools
```

### Step 5: Update Test Files
```bash
# Update test files that use 'app' variable
find . -name "*_test.go" -exec sed -i 's/app\./registry\./g' {} \;
# Manual review required for function parameters
```

## 🧪 Testing After Migration

### 1. Unit Tests
```bash
# Run all unit tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 2. Integration Tests
```bash
# Run integration tests with tags
go test -tags=integration ./...

# Test middleware stack
go test ./internal/middlewares/...
```

### 3. Manual Testing
- **Security Headers**: Check browser dev tools Network tab for security headers
- **HTTPS Redirect**: Test HTTP to HTTPS redirection in production
- **Task Queue**: Verify background tasks are processed correctly
- **Database Operations**: Test all database CRUD operations

### 4. Performance Testing
```bash
# Test application startup time
time go run ./cmd/server

# Test middleware performance
go test -bench=. ./internal/middlewares/...
```

## 📝 Additional Notes

### New Features Added
- **IndexNow Integration**: Automatic search engine pinging for faster indexing
- **MCP Support**: Model Context Protocol routes and API key configuration
- **Blog Post Versioning**: Enable versioning for blog posts
- **Enhanced Security**: Automatic HTTPS redirect and security headers

### Removed Features
- **Gitpod Configuration**: Removed `.gitpod.yml` file
- **Docker Tasks**: Removed obsolete Docker tasks from taskfile
- **CMS Old Routes**: Removed deprecated CMS controller routes

### Configuration Changes
- **IndexNow Key**: New configuration option `GetIndexNowKey()`
- **MCP API Key**: New MCP API key configuration support
- **Environment Detection**: Improved development environment detection

## 🆘 Common Issues and Solutions

### Issue 1: Variable Name Confusion
**Problem**: Mixing `app` and `registry` variable names
**Solution**: Use `registry` for all variables, keep `app` only for struct fields

### Issue 2: HTTPS Redirect in Development
**Problem**: HTTPS redirect interfering with local development
**Solution**: The middleware automatically excludes localhost, 127.0.0.1, and private IP ranges

### Issue 3: CSP Policy Too Restrictive
**Problem**: Content Security Policy blocking external resources
**Solution**: Update CSP policy in `security_headers.go` or configure custom middleware

### Issue 4: Task Queue Not Starting
**Problem**: Background processes failing due to nil checks
**Solution**: Ensure all required stores are properly initialized before starting

### Issue 5: Database Connection Issues
**Problem**: SQLite version compatibility issues
**Solution**: Run `go mod tidy` and test database operations thoroughly

## 📞 Support Information

- **Documentation**: Check `docs/` directory for detailed guides
- **Examples**: Review test files for usage patterns
- **Issues**: Report bugs to the repository issue tracker
- **Community**: Join community discussions for migration help

---

**Migration Checklist**:
- [ ] Dependencies updated with `go mod tidy`
- [ ] Variable names changed from `app` to `registry`
- [ ] Task queue initialization updated
- [ ] Security middlewares reviewed and tested
- [ ] All tests passing
- [ ] Manual testing completed
- [ ] Performance verified
- [ ] Documentation updated
