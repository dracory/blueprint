# Upgrade Guide: v0.11.0 to v0.12.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.11.0 to v0.12.0.

## 🚨 Breaking Changes

### 1. New Security Middlewares Added to Global Stack
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
// New global middleware stack includes HTTPS redirect and security headers
func globalMiddlewares(registry registry.RegistryInterface) []rtr.MiddlewareInterface {
    globalMiddlewares := []rtr.MiddlewareInterface{
        // ... existing middlewares
        middlewares.LogRequestMiddleware(registry),
        middlewares.NewHTTPSRedirectMiddleware(),      // NEW
        middlewares.NewSecurityHeadersMiddleware(),     // NEW
        middlewares.ThemeMiddleware(),
        middlewares.AuthMiddleware(registry),
    }
    return globalMiddlewares
}
```

**Action Required**:
- No action needed - these middlewares are automatically included
- Review your application's behavior if it relies on HTTP responses
- Test HTTPS redirect behavior in production environments

### 2. Email Template Variable Name Fix
**Change**: Fixed variable name from `registryName` to `appName` in email template creation.

**Old Usage**:
```go
// In internal/emails/create_email_template.go
registryName := lo.IfF(registry != nil, func() string {
    if registry.GetConfig() == nil {
        return ""
    }
    return registry.GetConfig().GetAppName()
}).Else("")
```

**New Usage**:
```go
// In internal/emails/create_email_template.go
appName := lo.IfF(registry != nil, func() string {
    if registry.GetConfig() == nil {
        return ""
    }
    return registry.GetConfig().GetAppName()
}).Else("")
```

**Action Required**:
- No action needed - this is an internal fix
- Email templates will now properly display the application name

### 3. Docker Tasks Removed from Taskfile
**Change**: Removed all Docker-related tasks from the taskfile.yml to simplify the project structure.

**Old Usage**:
```yaml
# Previously available tasks
docker:build:
  desc: Build the Docker image
  cmds:
    - docker build -t tap-api .

docker:run:
  desc: Run the Docker container on port 8080
  cmds:
    - docker run --rm -p {{.PORT}}:8080 tap-api

docker:up:
  desc: Run with docker compose
  cmds:
    - docker compose up --build
```

**New Usage**:
```yaml
# Docker tasks removed - use native Docker commands instead
# Build: docker build -t your-app .
# Run: docker run --rm -p 8080:8080 your-app
```

**Action Required**:
- Update CI/CD pipelines to use native Docker commands instead of taskfile tasks
- Update documentation that references the old Docker tasks
- Use `docker build` and `docker run` commands directly

### 4. Config Implementation Reformatted
**Change**: Config getters and setters have been reformatted with consistent spacing and section comments for better code organization.

**Old Usage**:
```go
// Previously compact formatting
func (c *configImplementation) SetAppName(v string) { c.appName = v }
func (c *configImplementation) GetAppName() string  { return c.appName }
```

**New Usage**:
```go
// New formatting with section comments and consistent spacing
// ============================================================================
// START: App Configuration
// ============================================================================

func (c *configImplementation) SetAppName(v string) { 
    c.appName = v 
}

func (c *configImplementation) GetAppName() string  { 
    return c.appName 
}

// ============================================================================
// END: App Configuration
// ============================================================================
```

**Action Required**:
- No action needed - this is a code formatting improvement
- The API remains unchanged, only the source code formatting is different

## 🔄 Migration Steps

### Step 1: Update Dependencies
```bash
# Update to latest dependencies
go mod tidy
go mod download
```

### Step 2: Review Middleware Configuration
```bash
# Check if your application needs custom HTTPS or security header configuration
grep -r "HTTPS" internal/
grep -r "security" internal/
```

### Step 3: Update CI/CD Docker Commands
```bash
# Replace old taskfile Docker commands with native Docker commands
# Old: task docker:build
# New: docker build -t your-app .

# Old: task docker:run
# New: docker run --rm -p 8080:8080 your-app
```

### Step 4: Test Email Functionality
```bash
# Verify email templates still work correctly
go test ./internal/emails/...
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
# Test middleware stack
go test -tags=integration ./internal/middlewares/...

# Test email functionality
go test -tags=integration ./internal/emails/...
```

### 3. Security Tests
```bash
# Test HTTPS redirect middleware
curl -I http://your-domain.com
# Should return 301 redirect to HTTPS

# Test security headers
curl -I https://your-domain.com
# Should include security headers like HSTS, X-Frame-Options, etc.
```

### 4. Manual Testing
- Verify HTTPS redirects work in production
- Check that security headers are present in responses
- Test email template rendering with application name
- Confirm Docker deployment still works with native commands

## 📝 Additional Notes

### New Features
- **Automatic HTTPS Redirect**: Applications now automatically redirect HTTP to HTTPS in production
- **Security Headers**: Built-in security headers including HSTS, X-Frame-Options, and CSP
- **Improved Code Organization**: Better formatting and section comments in config implementation

### Removed Features
- **Docker Taskfile Tasks**: Docker tasks removed from taskfile.yml to reduce complexity
- No breaking API changes - all changes are additive or internal improvements

### Performance Improvements
- No performance impact from the new middlewares
- Code organization improvements may slightly improve build times

## 🆘 Common Issues and Solutions

### Issue 1: HTTPS Redirect in Development
**Problem**: HTTPS redirect interferes with local development
**Solution**: The middleware automatically skips redirects for localhost and 127.0.0.1 addresses

### Issue 2: Security Headers Too Restrictive
**Problem**: Content Security Policy blocks some resources
**Solution**: Modify the CSP policy in `internal/middlewares/security_headers.go` if needed

### Issue 3: Docker Deployment Broken
**Problem**: CI/CD pipeline fails after Docker tasks removed
**Solution**: Update pipeline to use native Docker commands instead of taskfile tasks

### Issue 4: Email Templates Not Showing App Name
**Problem**: Email templates show empty application name
**Solution**: Ensure `APP_NAME` environment variable is set or config is properly initialized

## 📞 Support

For additional support:
- Check the [GitHub repository](https://github.com/dracory/blueprint)
- Review existing issues and discussions
- Create a new issue for upgrade-specific problems
- Consult the main documentation for detailed configuration options

---

**Version Information**:
- From: v0.11.0
- To: v0.12.0 (current main)
- Go Version: 1.25.0
- Release Date: January 2026
