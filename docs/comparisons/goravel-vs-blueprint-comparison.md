# Goravel vs Blueprint Comparison

## Overview

This document provides a comprehensive comparison between **Goravel** (a Laravel-inspired Go framework) and **Blueprint** (the current MVC web application framework).

## Framework Philosophy

### Goravel
- **Inspiration**: Laravel-like syntax and conventions
- **Approach**: Convention over configuration
- **Target Audience**: Developers familiar with PHP/Laravel transitioning to Go
- **Design Philosophy**: Elegant, expressive syntax with comprehensive built-in features

### Blueprint
- **Inspiration**: Custom MVC architecture with modular design
- **Approach**: Flexible, component-based architecture
- **Target Audience**: Go developers wanting a customizable foundation
- **Design Philosophy**: Modular, extensible with specialized stores/services

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
│   ├── server/
│   └── deploy/
├── internal/
│   ├── controllers/
│   ├── middlewares/
│   ├── routes/
│   ├── tasks/
│   ├── schedules/
│   └── config/
├── pkg/
└── docs/
```

## Core Features Comparison

| Feature | Goravel | Blueprint |
|---------|---------|-----------|
| **Routing** | Chi-based with Laravel-like syntax | Custom router with CLI listing |
| **ORM** | GORM integration | Custom store-based architecture |
| **Authentication** | Built-in auth system | `github.com/dracory/auth` package |
| **Database Migrations** | Built-in migration system | Custom migration tools |
| **Queue System** | Built-in queue management | `github.com/dracory/taskstore` |
| **Scheduler** | Built-in task scheduler | `github.com/go-co-op/gocron` |
| **Cache** | Multiple cache drivers | `github.com/dracory/cachestore` |
| **Validation** | Built-in validation | `github.com/asaskevich/govalidator` |
| **Templating** | Built-in template engine | Pongo2 (`github.com/flosch/pongo2/v6`) |
| **CLI Tools** | Artisan-like commands | Custom CLI with task runner |

## Dependency Management

### Goravel Dependencies
- **Core**: `github.com/goravel/framework`
- **Database**: GORM with multiple drivers
- **Routing**: Chi router
- **Validation**: Built-in validation package
- **Authentication**: Built-in auth package
- **Queue**: Built-in queue system

### Blueprint Dependencies
- **Core**: Custom modular architecture
- **Database**: Multiple specialized stores (`*store` packages)
- **Routing**: Custom router implementation
- **Validation**: `github.com/asaskevich/govalidator`
- **Authentication**: `github.com/dracory/auth`
- **Queue**: `github.com/dracory/taskstore`

## Key Differences

### 1. **Modularity vs Convention**
- **Goravel**: Follows Laravel conventions, less configuration needed
- **Blueprint**: Highly modular, requires more setup but offers greater flexibility

### 2. **Ecosystem**
- **Goravel**: Growing ecosystem with Laravel-like packages
- **Blueprint**: Custom ecosystem with specialized `dracory/*` packages

### 3. **Learning Curve**
- **Goravel**: Easier for Laravel developers, steeper for pure Go developers
- **Blueprint**: More Go-idiomatic, easier for experienced Go developers

### 4. **Database Architecture**
- **Goravel**: Traditional ORM approach with GORM
- **Blueprint**: Store-based architecture with specialized stores for different data types

### 5. **CLI Experience**
- **Goravel**: Artisan-like commands familiar to Laravel developers
- **Blueprint**: Task-based CLI with custom commands

## Performance Considerations

### Goravel
- **Pros**: Optimized for common web patterns, good caching system
- **Cons**: Additional abstraction layers may impact performance slightly

### Blueprint
- **Pros**: Lightweight core, specialized stores for optimal performance
- **Cons**: Performance depends on implementation quality of individual stores

## Development Experience

### Goravel
```go
// Laravel-like syntax
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
// Go-idiomatic approach
func (c *UserController) Show(ctx context.Context, id string) (interface{}, error) {
    user, err := c.userService.FindByID(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

## Migration Considerations

### From Blueprint to Goravel
1. **Pros**: Standardized conventions, larger community
2. **Cons**: Significant refactoring required, loss of custom store architecture

### From Goravel to Blueprint
1. **Pros**: Greater flexibility, custom store architecture
2. **Cons**: More boilerplate code, smaller community

## Recommendations

### Choose Goravel if:
- You're coming from a Laravel background
- You prefer convention over configuration
- You want a framework that "just works" out of the box
- You value Laravel-like ecosystem and patterns

### Choose Blueprint if:
- You want maximum flexibility and customization
- You prefer a modular, component-based architecture
- You need specialized data stores for different use cases
- You want to maintain full control over the architecture

## Conclusion

Both frameworks have their strengths:

- **Goravel** excels at developer productivity and convention-based development
- **Blueprint** offers superior flexibility and a more Go-idiomatic approach

The choice depends on your team's background, project requirements, and long-term architectural goals. Blueprint's modular design with specialized stores provides unique advantages for complex applications with diverse data management needs, while Goravel offers a more standardized, familiar experience for developers transitioning from Laravel.

## Next Steps

Consider creating a proof-of-concept application using both frameworks to evaluate:
1. Development velocity
2. Performance characteristics
3. Team productivity
4. Maintenance overhead

This hands-on comparison will provide the most accurate assessment for your specific use case.
