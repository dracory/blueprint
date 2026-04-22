# Log Admin Package

A self-contained admin interface for log management, following the folder-per-controller pattern. This package is designed to be moved to a separate external package in the future.

## Features

- **Log Manager**: View and filter application logs
- **Filter Options**: Filter by log level, message, context, and date range
- **Live Components**: Uses LiveFlux for dynamic log viewing
- **Pagination**: Full pagination support for large log sets
- **Bulk Actions**: Delete multiple logs at once
- **Context Viewer**: View detailed context for individual log entries

## Installation

```bash
go get project/pkg/logadmin
```

## Usage

### Basic Setup

```go
package main

import (
    "log/slog"
    "net/http"

    "project/pkg/logadmin"
)

func main() {
    // Create the log admin
    admin, err := logadmin.New(logadmin.AdminOptions{
        Registry:     registry,
        LogAdminURL:  "/admin/logs",
        AdminHomeURL: "/admin",
        AuthUserID: func(r *http.Request) string {
            // Return authenticated user ID or empty string
            return getUserIDFromSession(r)
        },
    })

    if err != nil {
        panic(err)
    }

    // Register the handler
    http.HandleFunc("/admin/logs/", admin.Handle)
    http.ListenAndServe(":8080", nil)
}
```

### Integration with Registry Pattern

```go
// In your controller
func NewLogsAdminController(registry registry.RegistryInterface) http.HandlerFunc {
    admin, err := logadmin.New(logadmin.AdminOptions{
        Registry:     registry,
        LogAdminURL:  "/admin/logs",
        AdminHomeURL: "/admin",
        AuthUserID: func(r *http.Request) string {
            user := helpers.GetAuthUser(r)
            if user == nil {
                return ""
            }
            return user.ID()
        },
    })

    if err != nil {
        return func(w http.ResponseWriter, r *http.Request) {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }

    return admin.Handle
}
```

## Configuration

### AdminOptions

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `Registry` | `registry.RegistryInterface` | Yes | The registry instance for accessing stores and services |
| `AdminHomeURL` | `string` | No | URL for admin home (default: "/admin") |
| `LogAdminURL` | `string` | No | Base URL for log admin (default: "/admin/logs") |
| `AuthUserID` | `func(*http.Request) string` | No | Function to get authenticated user ID |

## Routes

The package handles the following routes under the configured `LogAdminURL`:

| Path | Action |
|------|--------|
| `/` or `/log-manager` | Log manager with filtering and viewing |

## Architecture

This package follows the folder-per-controller pattern:

1. **Single Entry Point**: The `New()` function creates an admin instance
2. **Self-Contained**: All handlers are internal to the package
3. **Layout Delegation**: The parent application provides the layout function
4. **Dependency Injection**: Registry is passed via options
5. **Authentication**: Optional auth check via callback function

## File Structure

```
pkg/logadmin/
├── logadmin.go           # Main admin struct, New(), and Handle()
├── routes.go              # Subcontroller routing/dispatch
├── errors.go              # Error definitions
├── README.md              # This file
├── shared/                # Shared constants and utilities
│   ├── constants.go       # Controller names
│   └── links.go           # URL generation helpers
└── log_manager/           # Log manager subcontroller
    ├── log_manager_controller.go
    ├── log_manager_controller_test.go
    ├── constants.go
    ├── log_filter_component.go
    ├── log_table_component.go
    └── log_list_results.go
```

## Subcontrollers

### Log Manager

The log manager subcontroller provides:
- Log filtering by level, message, context, and date range
- Paginated log table display
- Bulk delete operations
- Context viewer for detailed log inspection
- LiveFlux components for dynamic updates

## Future Migration

When moving to an external package:

1. Update imports from `project/pkg/logadmin` to the external package name
2. No other code changes needed - the API remains the same
3. The package has no dependencies on internal project code

## Testing

```bash
go test ./pkg/logadmin/...
```

## Dependencies

- `project/internal/registry` - Registry interface
- `project/internal/layouts` - Layout rendering
- `project/internal/links` - URL generation
- `project/internal/helpers` - Helper functions
- `github.com/dracory/rtr` - Routing
- `github.com/dracory/hb` - HTML builder
- `github.com/dracory/liveflux` - Live components
- `github.com/dracory/cdn` - CDN resources

## See Also

- Folder-per-controller pattern documentation: `docs/folder-per-controller-pattern.md`
- Example implementation: `pkg/blogadmin/`
