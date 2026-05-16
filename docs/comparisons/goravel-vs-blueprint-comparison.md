# Goravel vs Blueprint Comparison

## Overview

This document provides a comprehensive comparison between **Goravel** (a Laravel-inspired Go framework) and **Blueprint** (a modular MVC web application framework built on specialized packages).

## Framework Philosophy

### Goravel
- **Inspiration**: Laravel-like syntax and conventions
- **Approach**: Convention over configuration
- **Target Audience**: Developers familiar with PHP/Laravel transitioning to Go
- **Design Philosophy**: Elegant, expressive syntax with comprehensive built-in features
- **Architecture**: Monolithic framework with integrated components

### Blueprint
- **Inspiration**: Modular MVC architecture with specialized packages
- **Approach**: Component-based architecture with reusable packages
- **Target Audience**: Go developers wanting a customizable, modular foundation
- **Design Philosophy**: Composable architecture with specialized packages (rtr, database, base, etc.)
- **Architecture**: Modular ecosystem with independent, reusable packages

## Architecture Comparison

### Goravel Architecture
```
├── app/
│   ├── http/
│   │   ├── controllers/
│   │   ├── middleware/
│   │   └── requests/
│   ├── models/
│   ├── jobs/
│   └── providers/
├── config/
├── database/
│   ├── migrations/
│   └── seeders/
├── resources/
│   └── views/
└── routes/
```

### Blueprint Architecture
```
├── cmd/
│   ├── server/          # Main application entry point
│   └── deploy/          # Deployment utilities
├── internal/
│   ├── controllers/     # HTTP controllers
│   ├── middlewares/     # Custom middlewares
│   ├── routes/          # Route definitions using rtr
│   ├── tasks/           # Background tasks
│   ├── schedules/       # Scheduled jobs
│   ├── config/          # Configuration management
│   ├── registry/        # Dependency injection
│   └── widgets/         # CMS widgets
├── pkg/                 # Reusable packages
└── docs/                # Documentation
```

### Blueprint Package Ecosystem
Blueprint is built on a modular ecosystem of specialized packages:

- **github.com/dracory/rtr** - High-performance router with middleware chains, domain routing, route groups
- **github.com/dracory/database** - Database abstraction with SQLite optimizations, connection pooling
- **github.com/dracory/base** - Utilities (htmx, cfmt, files, http, url, types, etc.)
- **github.com/dracory/auth** - Authentication system
- **github.com/dracory/taskstore** - Task queue management
- **github.com/dracory/cachestore** - Caching layer
- **github.com/dracory/cmsstore** - CMS data management
- **github.com/dracory/blogstore** - Blog data management
- **And many more specialized *store packages**

## Core Features Comparison

| Feature | Goravel | Blueprint |
|---------|---------|-----------|
| **Routing** | Chi-based with Laravel-like syntax | `github.com/dracory/rtr` - High-performance router with middleware chains, domain routing, route groups, declarative config |
| **Database** | GORM integration with migrations | `github.com/dracory/database` - Database abstraction with SQLite optimizations, connection pooling, transaction support |
| **Authentication** | Built-in auth system | `github.com/dracory/auth` package - Dedicated authentication system |
| **Data Access** | ORM-based with models | Store-based architecture with specialized packages (blogstore, cmsstore, userstore, etc.) |
| **Queue System** | Built-in queue management | `github.com/dracory/taskstore` - Task queue management |
| **Scheduler** | Built-in task scheduler | `github.com/go-co-op/gocron` - Cron-like job scheduling |
| **Cache** | Multiple cache drivers | `github.com/dracory/cachestore` - Caching layer |
| **Validation** | Built-in validation | `github.com/asaskevich/govalidator` - Validation library |
| **Templating** | Built-in template engine | Pongo2 (`github.com/flosch/pongo2/v6`) - Django-like template engine |
| **CLI Tools** | Artisan-like commands | Custom CLI with task runner using `github.com/dracory/base/cli` |
| **Middleware** | Built-in middleware | `github.com/dracory/rtr/middlewares` - Comprehensive middleware suite (CORS, rate limiting, security headers, jail bots, etc.) |
| **HTMX Support** | Not built-in | `github.com/dracory/base/htmx` - HTMX utilities and helpers |
| **Configuration** | Environment-based config | `github.com/dracory/base/config` with encrypted environment variable support |
| **Email** | Built-in email | `github.com/dracory/email` - Email package with SMTP support |
| **File Handling** | Built-in filesystem | `github.com/dracory/base/files` - File utilities |

## Dependency Management

### Goravel Dependencies
- **Core**: `github.com/goravel/framework` (monolithic framework)
- **Database**: GORM with multiple drivers
- **Routing**: Chi router
- **Validation**: Built-in validation package
- **Authentication**: Built-in auth package
- **Queue**: Built-in queue system

### Blueprint Dependencies
- **Core**: Modular architecture with specialized packages
- **Routing**: `github.com/dracory/rtr` (standalone router package)
- **Database**: `github.com/dracory/database` (standalone database package)
- **Utilities**: `github.com/dracory/base` (utilities: htmx, cfmt, files, http, url, types, etc.)
- **Authentication**: `github.com/dracory/auth` (standalone auth package)
- **Queue**: `github.com/dracory/taskstore` (standalone task queue)
- **Data Stores**: Specialized packages (blogstore, cmsstore, userstore, cachestore, etc.)
- **CLI**: `github.com/dracory/base/cli` (generic CLI dispatcher)
- **Email**: `github.com/dracory/email` (standalone email package)

**Key Difference**: Blueprint uses a modular package ecosystem where each component is a standalone, reusable package that can be used independently or combined. Goravel is a monolithic framework with integrated components.

## Key Differences

### 1. **Modularity vs Convention**
- **Goravel**: Monolithic framework with Laravel conventions, less configuration needed
- **Blueprint**: Highly modular package ecosystem, requires more setup but offers greater flexibility and reusability

### 2. **Ecosystem**
- **Goravel**: Growing ecosystem with Laravel-like packages, framework-specific
- **Blueprint**: Modular ecosystem with standalone `dracory/*` packages that can be used independently in any Go project

### 3. **Learning Curve**
- **Goravel**: Easier for Laravel developers, steeper for pure Go developers due to framework-specific patterns
- **Blueprint**: More Go-idiomatic, easier for experienced Go developers, follows standard Go package patterns

### 4. **Database Architecture**
- **Goravel**: Traditional ORM approach with GORM, model-based
- **Blueprint**: Store-based architecture with specialized packages for different data types, context-aware operations, transaction support

### 5. **CLI Experience**
- **Goravel**: Artisan-like commands familiar to Laravel developers
- **Blueprint**: Generic CLI dispatcher from `github.com/dracory/base/cli` with type-safe command registration

### 6. **Package Reusability**
- **Goravel**: Components are tightly coupled to the framework
- **Blueprint**: Each package (rtr, database, base, auth, etc.) is standalone and can be used in any Go project without the full Blueprint framework

### 7. **Routing Capabilities**
- **Goravel**: Chi-based routing with Laravel-like syntax
- **Blueprint**: Advanced router (`github.com/dracory/rtr`) with middleware chains, domain-based routing, route groups, declarative configuration, multiple handler types

## Performance Considerations

### Goravel
- **Pros**: Optimized for common web patterns, good caching system, Chi router performance
- **Cons**: Additional abstraction layers may impact performance slightly, monolithic framework overhead

### Blueprint
- **Pros**: Lightweight core with modular packages, specialized stores for optimal performance, SQLite optimizations (WAL mode, connection pooling), high-performance router (rtr) with minimal allocations
- **Cons**: Performance depends on implementation quality of individual packages, but packages are independently optimized and tested

## Development Experience

### Goravel
```go
// Laravel-like syntax with framework-specific response types
func (r *UserController) Show(ctx context.Context, id string) response.Response {
    user, err := r.user.FindByID(id)
    if err != nil {
        return response.Json().Error(err)
    }
    return response.Json().Success(user)
}
```

### Blueprint
```go
// Go-idiomatic approach with rtr router
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
    userID := rtr.MustGetParam(r, "id")
    user, err := c.userService.FindByID(userID)
    if err != nil {
        rtr.JSONResponse(w, r, `{"error": "User not found"}`)
        return
    }
    rtr.JSONResponse(w, r, user)
}

// Or using rtr's JSON handler
rtr.GetJSON("/users/:id", func(w http.ResponseWriter, r *http.Request) string {
    userID := rtr.MustGetParam(r, "id")
    user, err := userService.FindByID(userID)
    if err != nil {
        return `{"error": "User not found"}`
    }
    return userJSON
})
```

## Migration Considerations

### From Blueprint to Goravel
1. **Pros**: Standardized conventions, larger community, Laravel-like patterns
2. **Cons**: Significant refactoring required, loss of modular package architecture, loss of specialized stores, framework lock-in

### From Goravel to Blueprint
1. **Pros**: Greater flexibility with modular packages, ability to use packages independently, custom store architecture, no framework lock-in
2. **Cons**: More initial setup required, learning curve for package ecosystem, smaller community

## Recommendations

### Choose Goravel if:
- You're coming from a Laravel background and want familiar patterns
- You prefer convention over configuration with minimal setup
- You want a monolithic framework that "just works" out of the box
- You value Laravel-like ecosystem and don't mind framework lock-in
- You're building a typical web application and don't need custom components

### Choose Blueprint if:
- You want maximum flexibility with a modular package ecosystem
- You prefer a component-based architecture with reusable packages
- You need specialized data stores for different use cases
- You want to maintain full control over the architecture
- You want to use packages independently in other Go projects
- You value Go-idiomatic patterns over framework-specific conventions
- You need advanced routing features (domain-based routing, middleware chains, declarative config)
- You want HTMX support built-in
- You need encrypted environment variable support

## Conclusion

Both frameworks have their strengths:

- **Goravel** excels at developer productivity and convention-based development with a monolithic framework approach
- **Blueprint** offers superior flexibility through a modular package ecosystem with standalone, reusable components

The choice depends on your team's background, project requirements, and long-term architectural goals. Blueprint's modular design with specialized packages (rtr, database, base, auth, etc.) provides unique advantages for complex applications with diverse data management needs and allows for package reuse across multiple projects. Goravel offers a more standardized, familiar experience for developers transitioning from Laravel but comes with framework lock-in.

**Key Takeaway**: Blueprint is not just a framework—it's a modular ecosystem of standalone packages that can be used independently or combined. This gives you the flexibility to use only what you need, mix and match packages, and avoid framework lock-in. Goravel provides a more traditional monolithic framework experience with less setup but less flexibility.

## Next Steps

Consider creating a proof-of-concept application using both frameworks to evaluate:
1. Development velocity
2. Performance characteristics
3. Team productivity
4. Maintenance overhead

This hands-on comparison will provide the most accurate assessment for your specific use case.
