# Blog Admin Package

A self-contained admin interface for blog management, following the same pattern as `cmsstore/admin`. This package is designed to be moved to a separate external package (`github.com/dracory/blogstore/admin`) in the future.

## Features

- **Dashboard**: Overview with post, category, and tag counts
- **Post Management**: List, create, update, and delete blog posts
- **Category Management**: Manage blog categories (if taxonomy enabled)
- **Tag Management**: Manage blog tags (if taxonomy enabled)
- **Pagination**: Full pagination support for lists
- **Search**: Search functionality for posts and terms
- **Responsive UI**: Bootstrap 5 based responsive interface

## Installation

```bash
go get project/pkg/blogadmin
```

## Usage

### Basic Setup

```go
package main

import (
    "log/slog"
    "net/http"
    
    "github.com/dracory/blogstore"
    "project/pkg/blogadmin"
)

func main() {
    // Initialize your blog store
    store, _ := blogstore.NewStore(blogstore.NewStoreOptions{
        DB: db,
    })
    
    // Create the blog admin
    admin, err := blogadmin.New(blogadmin.AdminOptions{
        Store:        store,
        Logger:       slog.Default(),
        BlogAdminURL: "/admin/blog",
        AdminHomeURL: "/admin",
        AuthUserID: func(r *http.Request) string {
            // Return authenticated user ID or empty string
            return getUserIDFromSession(r)
        },
        FuncLayout: func(pageTitle string, pageContent string, options struct {
            Styles     []string
            StyleURLs  []string
            Scripts    []string
            ScriptURLs []string
        }) string {
            // Your layout function that wraps the content
            return renderLayout(pageTitle, pageContent, options)
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    // Register the handler
    http.HandleFunc("/admin/blog/", admin.Handle)
    http.ListenAndServe(":8080", nil)
}
```

### Integration with Registry Pattern

```go
// In your controller
func NewBlogAdminController(registry registry.RegistryInterface) http.HandlerFunc {
    admin, err := blogadmin.New(blogadmin.AdminOptions{
        Store:        registry.GetBlogStore(),
        Logger:       registry.GetLogger(),
        BlogAdminURL: "/admin/blog",
        AdminHomeURL: "/admin",
        AuthUserID: func(r *http.Request) string {
            user := helpers.GetAuthUser(r)
            if user == nil {
                return ""
            }
            return user.GetID()
        },
        FuncLayout: func(pageTitle string, pageContent string, options struct {
            Styles     []string
            StyleURLs  []string
            Scripts    []string
            ScriptURLs []string
        }) string {
            return layouts.NewAdminLayout(registry, nil, layouts.Options{
                Title:      pageTitle,
                Content:    hb.Raw(pageContent),
                StyleURLs:  options.StyleURLs,
                ScriptURLs: options.ScriptURLs,
                Styles:     options.Styles,
                Scripts:    options.Scripts,
            }).ToHTML()
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
| `Store` | `blogstore.StoreInterface` | Yes | The blog store instance |
| `Logger` | `*slog.Logger` | No | Logger for errors and info |
| `FuncLayout` | `func(string, string, struct) string` | Yes | Layout wrapper function |
| `AdminHomeURL` | `string` | No | URL for admin home (default: "/admin") |
| `BlogAdminURL` | `string` | No | Base URL for blog admin (default: "/admin/blog") |
| `AuthUserID` | `func(*http.Request) string` | No | Function to get authenticated user ID |

## Routes

The package handles the following routes under the configured `BlogAdminURL`:

| Path | Action |
|------|--------|
| `/` or `/dashboard` | Dashboard with stats |
| `/posts` | Post list/manager |
| `/post-create` | Create new post |
| `/post-update?id={id}` | Edit post |
| `/post-delete?id={id}` | Delete post confirmation |
| `/categories` | Category manager (if taxonomy enabled) |
| `/tags` | Tag manager (if taxonomy enabled) |

## Architecture

This package follows the same pattern as `cmsstore/admin`:

1. **Single Entry Point**: The `New()` function creates an admin instance
2. **Self-Contained**: All handlers are internal to the package
3. **Layout Delegation**: The parent application provides the layout function
4. **Dependency Injection**: Store and logger are passed via options
5. **Authentication**: Optional auth check via callback function

## File Structure

```
pkg/blogadmin/
├── blogadmin.go    # Main admin struct, New(), and Handle()
├── dashboard.go    # Dashboard handler and UI
├── posts.go        # Post CRUD handlers
├── taxonomy.go     # Category and tag handlers
├── errors.go       # Error definitions
└── README.md       # Documentation
```

## Future Migration

When moving to an external package (`github.com/dracory/blogstore/admin`):

1. Update imports from `project/pkg/blogadmin` to `github.com/dracory/blogstore/admin`
2. No other code changes needed - the API remains the same
3. The package has no dependencies on internal project code

## Testing

```bash
go test ./pkg/blogadmin/...
```

## Dependencies

- `github.com/dracory/blogstore` - Blog store interface
- `github.com/dracory/hb` - HTML builder
- `github.com/dracory/uid` - Unique ID generation
