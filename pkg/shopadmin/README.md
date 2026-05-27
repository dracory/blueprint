# Shop Admin Package

A self-contained admin interface for shop management, following the folder-per-controller pattern with Vue.js SPA frontend. This package is designed to be moved to a separate external package in the future.

## Features

- **Dashboard**: Overview with product, category, and order counts
- **Product Management**: List, create, update, and delete products
- **Category Management**: Manage product categories
- **Discount Management**: Manage discounts
- **Order Management**: View and manage orders
- **Vue.js SPA**: Modern single-page application interface
- **AJAX Actions**: Dynamic content loading without page refreshes
- **Responsive UI**: Bootstrap 5 based responsive interface

## Installation

```bash
go get project/pkg/shopadmin
```

## Usage

### Basic Setup

```go
package main

import (
    "log/slog"
    "net/http"
    
    "project/pkg/shopadmin"
)

func main() {
    // Create the shop admin
    admin, err := shopadmin.New(shopadmin.AdminOptions{
        Registry:     registry,
        ShopAdminURL: "/admin/shop",
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
    http.HandleFunc("/admin/shop/", admin.Handle)
    http.ListenAndServe(":8080", nil)
}
```

### Integration with Registry Pattern

```go
// In your controller
func NewShopAdminController(registry registry.RegistryInterface) http.HandlerFunc {
    admin, err := shopadmin.New(shopadmin.AdminOptions{
        Registry:     registry,
        ShopAdminURL:  "/admin/shop",
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
| `ShopAdminURL` | `string` | No | Base URL for shop admin (default: "/admin/shop") |
| `AuthUserID` | `func(*http.Request) string` | No | Function to get authenticated user ID |

## Routes

The package handles the following routes under the configured `ShopAdminURL`:

| Path | Action |
|------|--------|
| `/` or `/home` | Dashboard with stats |
| `/products` | Product list/manager |
| `/categories` | Category manager |
| `/discounts` | Discount manager |
| `/orders` | Order manager |

## Architecture

This package follows the folder-per-controller pattern with Vue.js SPA:

1. **Single Entry Point**: The `New()` function creates an admin instance
2. **Self-Contained**: All handlers are internal to the package
3. **Vue.js SPA**: Embedded HTML/JS files for frontend logic
4. **AJAX Actions**: Dynamic content loading via JSON responses
5. **Dependency Injection**: Registry is passed via options
6. **Authentication**: Optional auth check via callback function

## File Structure

```
pkg/shopadmin/
├── shopadmin.go           # Main admin struct, New(), and Handle()
├── routes.go              # Subcontroller routing/dispatch
├── errors.go              # Error definitions
├── README.md              # This file
├── shared/                # Shared constants and utilities
│   ├── constants.go       # Controller names
│   └── links.go           # URL generation helpers
├── home/                  # Dashboard subcontroller
│   ├── home_controller.go
│   ├── home.html          # Vue.js template
│   └── home.js            # Vue.js component logic
└── product_manager/       # Product manager subcontroller
    └── product_manager_controller.go
```

## Subcontrollers

### Home (Dashboard)

The dashboard subcontroller provides:
- Overview stats (products, categories, orders)
- Navigation tiles to other sections
- Vue.js component for dynamic updates

### Product Manager

The product manager subcontroller provides:
- Product listing with pagination
- Product deletion (single and bulk)
- Filtering and sorting
- AJAX-based operations

## Future Migration

When moving to an external package:

1. Update imports from `project/pkg/shopadmin` to the external package name
2. No other code changes needed - the API remains the same
3. The package has no dependencies on internal project code

## Testing

```bash
go test ./pkg/shopadmin/...
```

## Dependencies

- `project/internal/registry` - Registry interface
- `project/internal/layouts` - Layout rendering
- `project/internal/links` - URL generation
- `project/internal/helpers` - Helper functions
- `github.com/dracory/rtr` - Routing
- `github.com/dracory/hb` - HTML builder
- `github.com/dracory/api` - API responses
- `github.com/dracory/cdn` - CDN resources
- `github.com/dracory/shopstore` - Shop store interface

## See Also

- Folder-per-controller pattern documentation
- Vue.js pattern example: `pkg/logadmin/`
