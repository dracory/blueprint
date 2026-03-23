# Proposal: Registry Base & Store Manager Migration

## Overview

This proposal outlines the migration of registry orchestration logic and store management patterns from the Blueprint project to a shared base package. This will create a reusable foundation for all Dracory projects while maintaining flexibility for application-specific configurations.

## Current State Analysis

### Blueprint Registry Implementation
The current Blueprint registry contains approximately 400+ lines of boilerplate code that handles:

1. **Registry Orchestration** (`registry_implementation.go`):
   - Cache directory detection and setup
   - Logger initialization (console + database)
   - Database connection management
   - Memory and file cache setup
   - Store initialization and migration orchestration

2. **Store Management Pattern** (24 stores):
   - Each store follows identical `Initialize` → `Migrate` → `NewStore` pattern
   - Configuration-based conditional initialization
   - Repetitive accessor methods (Get/Set for each store)
   - Manual registration in initialization/migration arrays

### Identified Issues
- **Code Duplication**: Every Dracory project replicates this orchestration logic
- **Maintenance Overhead**: Adding new stores requires updates in multiple files
- **Inconsistency Risk**: Manual registration leads to missed stores or wrong ordering
- **Testing Complexity**: Large monolithic registry implementation is hard to test

## Proposed Solution

### 1. Generic Registry Base (`github.com/dracory/base/registry`)

Create a generic, type-safe registry base that handles all common orchestration:

```go
// Generic registry base for any configuration type
type Registry[T ConfigInterface] struct {
    cfg T
    db  *sql.DB
    
    // Common infrastructure
    databaseLogger *slog.Logger
    consoleLogger  *slog.Logger
    memoryCache    *ttlcache.Cache[string, any]
    fileCache      cachego.Cache
    
    // Store manager
    storeManager *StoreManager
}

// Constructor with automatic orchestration
func New[T ConfigInterface](cfg T) (*Registry[T], error)

// Graceful cleanup
func (r *Registry[T]) Close() error
```

### 2. Store Manager (`github.com/dracory/base/registry/manager`)

Implement a generic store manager that handles initialization and migration:

```go
// Store configuration and factory
type StoreConfig struct {
    Name      string
    Enabled   func(cfg ConfigInterface) bool
    Factory   func(db *sql.DB) (StoreInterface, error)
    Migrator  func(store StoreInterface) error
}

// Generic store manager
type StoreManager struct {
    stores map[string]StoreInterface
    configs map[string]StoreConfig
}

// Registration and lifecycle management
func (sm *StoreManager) Register(config StoreConfig) error
func (sm *StoreManager) Initialize(cfg ConfigInterface, db *sql.DB) error
func (sm *StoreManager) Migrate() error
func (sm *StoreManager) Get(name string) StoreInterface
```

### 3. Standard Store Configurations

Pre-defined configurations for common Dracory stores:

```go
// Standard store configurations registry
var StandardStores = map[string]StoreConfig{
    "audit": {
        Name:    "audit",
        Enabled: func(cfg ConfigInterface) bool { return cfg.GetAuditStoreUsed() },
        Factory: func(db *sql.DB) (StoreInterface, error) { return auditstore.New(db) },
        Migrator: func(store StoreInterface) error { return store.AutoMigrate() },
    },
    "blog": {
        Name:    "blog", 
        Enabled: func(cfg ConfigInterface) bool { return cfg.GetBlogStoreUsed() },
        Factory: func(db *sql.DB) (StoreInterface, error) { return blogstore.New(db) },
        Migrator: func(store StoreInterface) error { return store.AutoMigrate() },
    },
    // ... 22 more standard stores
}
```

## Migration Plan

### Phase 1: Create Base Registry Package

1. **Create Base Registry Structure**:
   ```
   github.com/dracory/base/registry/
   ├── registry.go          # Generic registry base
   ├── manager/
   │   ├── manager.go       # Store manager implementation
   │   ├── config.go        # Store configuration types
   │   └── standard.go      # Standard store definitions
   ├── cache/
   │   └── directory.go     # Cache directory logic
   └── logger/
       └── setup.go         # Logger initialization
   ```

2. **Implement Core Components**:
   - Generic registry base with configuration type parameter
   - Store manager with registration and lifecycle management
   - Cache directory detection and setup
   - Logger initialization patterns
   - Standard store configurations for all 24 stores

### Phase 2: Blueprint Integration

1. **Blueprint Registry Simplification**:
   ```go
   // Before: 400+ lines of boilerplate
   type registryImplementation struct { /* 24 store fields + infrastructure */ }
   
   // After: Application-specific configuration only
   type BlueprintRegistry struct {
       registry.BaseRegistry[config.Config]
   }
   
   func New(cfg config.ConfigInterface) (RegistryInterface, error) {
       base, err := registry.New(cfg)
       if err != nil {
           return nil, err
       }
       return &BlueprintRegistry{BaseRegistry: base}, nil
   }
   ```

2. **Remove Duplicate Code**:
   - Delete `registry_implementation.go` (400+ lines)
   - Delete `registry_datastores_initialize.go` (48 lines)
   - Delete `registry_datastores_migrate.go` (40 lines)
   - Simplify all 24 store files to remove Initialize/Migrate boilerplate

### Phase 3: Enhanced Features

1. **Store Discovery**:
   ```go
   // Automatic store discovery for custom stores
   func (sm *StoreManager) DiscoverStores(pkgPath string) error
   
   // Runtime store registration
   func (sm *StoreManager) RegisterCustom(name string, config StoreConfig) error
   ```

2. **Health Monitoring**:
   ```go
   // Store health checks
   func (sm *StoreManager) HealthCheck() map[string]error
   
   // Connection status monitoring
   func (r *Registry[T]) GetStoreStatus() StoreStatusMap
   ```

3. **Configuration Validation**:
   ```go
   // Validate store configurations before initialization
   func (sm *StoreManager) ValidateConfig(cfg ConfigInterface) error
   
   // Dependency checking between stores
   func (sm *StoreManager) ValidateDependencies() error
   ```

## Benefits

### 1. Code Reduction
- **Blueprint**: Remove ~500 lines of boilerplate code
- **All Projects**: Eliminate registry duplication across Dracory ecosystem
- **Maintenance**: Single source of truth for store management patterns

### 2. Developer Experience
- **Faster Onboarding**: New projects get full store management instantly
- **Type Safety**: Generic registry prevents configuration type mismatches
- **Consistency**: All projects use identical store initialization patterns

### 3. Extensibility
- **Custom Stores**: Easy registration of application-specific stores
- **Plugin Architecture**: Store discovery enables plugin-based extensions
- **Configuration Flexibility**: Projects can enable/disable stores via configuration

### 4. Reliability
- **Testing**: Smaller, focused components are easier to test comprehensively
- **Error Handling**: Centralized error handling for store lifecycle
- **Health Monitoring**: Built-in store health checks and monitoring

## Implementation Details

### Configuration Interface Requirements

Projects must implement a common configuration interface:

```go
type ConfigInterface interface {
    // Database configuration
    GetDatabaseType() string
    GetDatabaseConnectionString() string
    
    // Store enablement flags (generated from standard stores)
    GetAuditStoreUsed() bool
    GetBlogStoreUsed() bool
    GetCacheStoreUsed() bool
    // ... etc for all standard stores
    
    // Cache configuration
    GetCacheDirectory() string
    
    // Logging configuration
    GetLogLevel() string
    GetLogFormat() string
}
```

### Migration Compatibility

1. **Backward Compatibility**: Existing Blueprint code continues to work unchanged
2. **Gradual Migration**: Projects can migrate incrementally
3. **Configuration Bridge**: Adapter pattern for existing configuration formats

### Testing Strategy

1. **Unit Tests**: Each component tested in isolation
2. **Integration Tests**: Full registry lifecycle testing
3. **Mock Stores**: Test store manager with mock store implementations
4. **Configuration Tests**: Validate all configuration combinations

## File Structure

### Before (Blueprint)
```
internal/registry/
├── registry_implementation.go           (431 lines)
├── registry_datastores_initialize.go    (48 lines)
├── registry_datastores_migrate.go       (40 lines)
├── stores_audit.go                      (73 lines)
├── stores_blog.go                       (78 lines)
├── stores_*.go                          (22 more files)
└── [24 store files with repetitive patterns]
```

### After (Base Package)
```
github.com/dracory/base/registry/
├── registry.go                          (120 lines)
├── manager/
│   ├── manager.go                       (180 lines)
│   ├── config.go                        (65 lines)
│   └── standard.go                      (200 lines)
├── cache/
│   └── directory.go                     (45 lines)
└── logger/
    └── setup.go                         (35 lines)

# Blueprint (simplified)
internal/registry/
├── registry.go                          (25 lines)
└── [application-specific configuration only]
```

## Success Metrics

1. **Code Reduction**: 80% reduction in registry-related code across projects
2. **Adoption Rate**: All new Dracory projects use base registry within 6 months
3. **Bug Reduction**: 50% reduction in store initialization bugs
4. **Developer Satisfaction**: Faster project setup time and fewer configuration issues

## Timeline

- **Phase 1**: 2-3 weeks (Base package creation)
- **Phase 2**: 1-2 weeks (Blueprint integration)
- **Phase 3**: 2-3 weeks (Enhanced features)
- **Testing & Documentation**: 1-2 weeks

**Total Estimated Time**: 6-10 weeks

## Risks and Mitigations

### Risk: Configuration Interface Changes
- **Mitigation**: Provide adapter pattern for existing configurations
- **Fallback**: Maintain backward compatibility for existing projects

### Risk: Store Dependency Complexity
- **Mitigation**: Implement dependency validation and clear error messages
- **Testing**: Comprehensive dependency testing in base package

### Risk: Performance Overhead
- **Mitigation**: Use generics for compile-time type safety
- **Benchmarking**: Performance testing against current implementation

## Conclusion

This migration will significantly reduce code duplication, improve maintainability, and provide a solid foundation for all Dracory projects. The generic registry base with store manager will enable faster project bootstrapping while maintaining the flexibility needed for application-specific requirements.

The proposed solution balances code reuse with extensibility, ensuring that all Dracory projects benefit from standardized store management while retaining the ability to customize and extend as needed.
