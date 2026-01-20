# Upgrade Guide: v0.8.0 to v0.9.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.8.0 to v0.9.0.

## ⚠️ Breaking Changes

### 1. Main.go Architecture Refactor
**Change**: Complete refactor of main.go with graceful shutdown and context management

**Old Usage**:
```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        cfmt.Error("Failed to load config:", slog.Any("error", err))
        return
    }
    
    application, err := app.New(cfg)
    if err != nil {
        cfmt.Error("Failed to initialize app:", slog.Any("error", err))
        return
    }
    
    defer closeResourcesDB(application.GetDB())
    tasks.RegisterTasks(application)
    
    if isCliMode() {
        cli.ExecuteCliCommand(application, os.Args[1:])
        return
    }
    
    startBackgroundProcesses(application)
    
    _, err = websrv.Start(websrv.Options{
        Host:    application.GetConfig().GetAppHost(),
        Port:    application.GetConfig().GetAppPort(),
        URL:     application.GetConfig().GetAppUrl(),
        Handler: routes.Routes(application).ServeHTTP,
    })
    
    if err != nil {
        cfmt.Errorf("Failed to start server: %v", err)
        return
    }
}
```

**New Usage**:
```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        fmt.Printf("Failed to load config: %v\n", err)
        return
    }

    application, err := app.New(cfg)
    if err != nil {
        fmt.Printf("Failed to initialize app: %v\n", err)
        return
    }

    defer closeResourcesDB(application.GetDB())
    tasks.RegisterTasks(application)

    if isCliMode() {
        if len(os.Args) < 2 {
            return
        }
        if err := cli.ExecuteCliCommand(application, os.Args[1:]); err != nil {
            fmt.Printf("Failed to execute CLI command: %v\n", err)
            os.Exit(1)
        }
        return
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    background := newBackgroundGroup(ctx)
    startBackgroundProcesses(ctx, background, application)

    server, err := websrv.Start(websrv.Options{
        Host:    application.GetConfig().GetAppHost(),
        Port:    application.GetConfig().GetAppPort(),
        URL:     application.GetConfig().GetAppUrl(),
        Handler: routes.Routes(application).ServeHTTP,
    })

    if err != nil {
        fmt.Printf("Failed to start server: %v\n", err)
        background.stop()
        return
    }

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    select {
    case <-sigs:
        fmt.Println("Shutdown signal received, draining background workers")
        cancel()
    case <-background.Done():
        cancel()
    }

    background.stop()
    if server != nil {
        shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer shutdownCancel()
        if err := server.Shutdown(shutdownCtx); err != nil {
            slog.Error("Server shutdown failed", "error", err)
        }
    }
}
```

**Action Required**:
- Replace entire main.go file with new version
- Add new imports: `context`, `sync`, `syscall`, `time`
- Add background process management functions
- Implement graceful shutdown handling

### 2. Background Process Management
**Change**: Background processes now use context-aware goroutine management

**Old Usage**:
```go
func startBackgroundProcesses(app types.AppInterface) {
    if app.GetDB() != nil {
        if ts := app.GetTaskStore(); ts != nil {
            go ts.QueueRunGoroutine(10, 2)
        }
        if cs := app.GetCacheStore(); cs != nil {
            go cs.ExpireCacheGoroutine()
        }
        if ss := app.GetSessionStore(); ss != nil {
            go ss.SessionExpiryGoroutine()
        }
    }
    schedules.StartAsync(app)
    emails.InitEmailSender(app)
    middlewares.CmsAddMiddlewares(app)
    widgets.CmsAddShortcodes(app)
}
```

**New Usage**:
```go
func startBackgroundProcesses(ctx context.Context, group *backgroundGroup, app types.AppInterface) {
    if app.GetDB() != nil {
        if ts := app.GetTaskStore(); ts != nil {
            group.Go(func(ctx context.Context) {
                ts.QueueRunGoroutine(ctx, 10, 2)
            })
        }
        if cs := app.GetCacheStore(); cs != nil {
            group.Go(func(ctx context.Context) {
                if err := cs.ExpireCacheGoroutine(ctx); err != nil {
                    slog.Error("Cache expiration goroutine failed", "error", err)
                }
            })
        }
        if ss := app.GetSessionStore(); ss != nil {
            group.Go(func(ctx context.Context) {
                if err := ss.SessionExpiryGoroutine(ctx); err != nil {
                    slog.Error("Session expiry goroutine failed", "error", err)
                }
            })
        }
    }

    group.Go(func(ctx context.Context) {
        schedules.StartAsync(ctx, app)
    })

    emails.InitEmailSender(app)
    middlewares.CmsAddMiddlewares(app)
    widgets.CmsAddShortcodes(app)
}
```

**Action Required**:
- Update background process function signatures to accept context
- Add error handling for goroutine failures
- Implement backgroundGroup struct for goroutine management

### 3. AppInterface Store Additions
**Change**: New stores added to AppInterface

**Old Usage**:
```go
// In internal/types/app.go
type AppInterface interface {
    // ... existing methods
    GetSqlFileStorage() filesystem.StorageInterface
    SetSqlFileStorage(s filesystem.StorageInterface)
    // ... no audit, chat, or subscription stores
}
```

**New Usage**:
```go
// In internal/types/app_interface.go
type AppInterface interface {
    // ... existing methods
    GetSqlFileStorage() filesystem.StorageInterface
    SetSqlFileStorage(s filesystem.StorageInterface)
    
    // New stores
    GetAuditStore() auditstore.StoreInterface
    SetAuditStore(s auditstore.StoreInterface)
    
    GetChatStore() chatstore.StoreInterface
    SetChatStore(s chatstore.StoreInterface)
    
    GetSubscriptionStore() subscriptionstore.StoreInterface
    SetSubscriptionStore(s subscriptionstore.StoreInterface)
}
```

**Action Required**:
- Rename `internal/types/app.go` to `internal/types/app_interface.go`
- Add new store interfaces and methods
- Update import statements to include new stores

### 4. Dependency Migration from gouniverse to dracory
**Change**: Multiple packages migrated from gouniverse namespace to dracory namespace

**Old Usage**:
```go
// In go.mod
github.com/gouniverse/csrf v0.1.0
github.com/gouniverse/dashboard v1.7.1
github.com/gouniverse/filesystem v0.3.1
github.com/mingrammer/cfmt v1.1.0
```

**New Usage**:
```go
// In go.mod
github.com/dracory/csrf v0.2.0
github.com/dracory/dashboard v1.11.0
github.com/dracory/filesystem v1.0.0
github.com/dracory/base/cfmt (imported from base package)
```

**Action Required**:
- Update all import statements in codebase
- Run `go mod tidy` to resolve dependencies
- Update package references in code

### 5. SQLite Driver Change
**Change**: SQLite driver switched from modernc.org to glebarez

**Old Usage**:
```go
// In go.mod
gorm.io/driver/sqlite v1.6.0
gorm.io/gorm v1.30.5
modernc.org/sqlite v1.38.2
```

**New Usage**:
```go
// In go.mod - SQLite moved to indirect requirements
modernc.org/sqlite v1.40.0 // indirect
// GORM dependencies removed
```

**Action Required**:
- Remove direct GORM dependencies from go.mod
- Update database initialization code if using GORM
- Run `go mod tidy` to clean up dependencies

### 6. Go Version Update
**Change**: Go version updated from 1.24.5 to 1.25.0

**Old Usage**:
```go
// In go.mod
go 1.24.5
```

**New Usage**:
```go
// In go.mod
go 1.25.0
```

**Action Required**:
- Update Go version in go.mod
- Ensure local Go installation is 1.25.0 or later
- Update CI/CD pipeline Go version if applicable

### 7. CLI Error Handling
**Change**: Enhanced error handling for CLI commands

**Old Usage**:
```go
if isCliMode() {
    cli.ExecuteCliCommand(application, os.Args[1:])
    return
}
```

**New Usage**:
```go
if isCliMode() {
    if len(os.Args) < 2 {
        return
    }
    if err := cli.ExecuteCliCommand(application, os.Args[1:]); err != nil {
        fmt.Printf("Failed to execute CLI command: %v\n", err)
        os.Exit(1)
    }
    return
}
```

**Action Required**:
- Add argument count validation
- Add error handling with proper exit codes
- Update CLI command execution flow

### 8. Error Logging Changes
**Change**: Switched from cfmt package to standard fmt for main.go errors

**Old Usage**:
```go
cfmt.Error("Failed to load config:", slog.Any("error", err))
cfmt.Errorf("Failed to start server: %v", err)
```

**New Usage**:
```go
fmt.Printf("Failed to load config: %v\n", err)
fmt.Printf("Failed to start server: %v\n", err)
```

**Action Required**:
- Replace cfmt.Error calls with fmt.Printf
- Update error message formatting
- Remove slog.Any wrapper for simple error messages

## 🔄 Migration Steps

### Step 1: Update Dependencies
```bash
# Backup current go.mod and go.sum
cp go.mod go.mod.backup
cp go.sum go.sum.backup

# Update Go version
sed -i 's/go 1.24.5/go 1.25.0/' go.mod

# Remove old dependencies
go mod edit -droprequire=gorm.io/driver/sqlite
go mod edit -droprequire=gorm.io/gorm

# Add new dependencies
go get github.com/dracory/auditstore@v0.2.0
go get github.com/dracory/chatstore@v0.6.0
go get github.com/dracory/subscriptionstore@v0.5.0
go get github.com/dracory/csrf@v0.2.0
go get github.com/dracory/dashboard@v1.11.0
go get github.com/dracory/filesystem@v1.0.0

# Update all dependencies
go mod tidy
```

### Step 2: Update Import Statements
```bash
# Find and replace old imports
find . -name "*.go" -type f -exec sed -i 's/github.com\/gouniverse\/csrf/github.com\/dracory\/csrf/g' {} \;
find . -name "*.go" -type f -exec sed -i 's/github.com\/gouniverse\/dashboard/github.com\/dracory\/dashboard/g' {} \;
find . -name "*.go" -type f -exec sed -i 's/github.com\/gouniverse\/filesystem/github.com\/dracory\/filesystem/g' {} \;
find . -name "*.go" -type f -exec sed -i 's/github.com\/mingrammer\/cfmt/github.com\/dracory\/base\/cfmt/g' {} \;
```

### Step 3: Update Main.go
```bash
# Replace main.go with new version
# Copy the new main.go content from the template above
```

### Step 4: Update AppInterface
```bash
# Rename the file
mv internal/types/app.go internal/types/app_interface.go

# Add new store interfaces to the interface
# Add the new import statements for auditstore, chatstore, subscriptionstore
```

### Step 5: Update Background Processes
```bash
# Update any custom background processes to use context
# Ensure all goroutines are properly managed with the backgroundGroup
```

### Step 6: Update Database Configuration
```bash
# If using GORM, update database initialization code
# Remove GORM-specific imports and code
# Update to use direct SQL or new database patterns
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
# Run integration tests
go test -tags=integration ./...

# Test CLI commands
go run main.go help
go run main.go version
```

### 3. Application Startup
```bash
# Test application startup
go run main.go

# Test graceful shutdown
# Start the app and send SIGINT/SIGTERM signals
```

### 4. Database Operations
```bash
# Test database connectivity
# Verify all stores are properly initialized
# Test CRUD operations through the application
```

## 📝 Additional Notes

### New Features Added
- **Graceful Shutdown**: Application now properly handles shutdown signals
- **Context Management**: All background processes are context-aware
- **New Stores**: Added audit, chat, and subscription stores
- **Enhanced Error Handling**: Better error reporting and CLI error handling
- **Dependency Cleanup**: Removed unused GORM dependencies

### Removed Features
- **GORM Dependencies**: Direct GORM dependencies removed from main requirements
- **Old Import Paths**: All gouniverse namespace imports removed

### Performance Improvements
- **Better Goroutine Management**: Background processes now properly managed
- **Context Cancellation**: Proper cleanup of resources on shutdown
- **Reduced Dependencies**: Cleaner dependency tree

## 🆘 Common Issues and Solutions

### Issue: Import Path Errors
**Problem**: Build fails due to unknown import paths
**Solution**: 
```bash
go mod tidy
go clean -modcache
go mod download
```

### Issue: Context-related Build Errors
**Problem**: Functions don't accept context parameter
**Solution**: Update function signatures to accept context.Context parameter

### Issue: Store Interface Errors
**Problem**: Missing store interface methods
**Solution**: Add new store methods to your AppInterface implementation

### Issue: Background Process Failures
**Problem**: Goroutines not starting properly
**Solution**: Ensure all background processes use the new backgroundGroup pattern

### Issue: Database Connection Issues
**Problem**: SQLite driver errors
**Solution**: Remove GORM imports and use direct SQL or updated database patterns

## 📞 Support

For additional support:
- Check the repository issues: https://github.com/dracory/blueprint/issues
- Review the documentation: https://github.com/dracory/blueprint/docs
- Compare with the example application in the repository
