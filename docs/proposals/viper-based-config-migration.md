# Proposal: Migrate to Viper-Based Configuration System

## Status
**REJECTED** - 2026-03-23

## Rejection Reason

**This proposal is NOT recommended for implementation.** While Viper offers brevity, Blueprint's current getter/setter pattern provides superior benefits:

- **Type Safety**: Compile-time errors vs runtime string typo failures
- **IDE Support**: Full autocomplete and refactoring support
- **Discoverability**: Interface clearly shows all available config
- **Explicit Contracts**: Clear API vs implicit string keys
- **Maintainability**: Easier to track config usage across codebase

The current 730 lines of getters/setters are **not a problem** - they're explicit, type-safe, and maintainable. The verbosity is a feature, not a bug.

**Keep this proposal for reference only.** If config management becomes a real pain point (not just "lots of lines"), revisit with specific problems to solve.

---

## Overview
Migrate Blueprint's configuration system from a verbose getter/setter pattern to a Viper-based dynamic configuration system, inspired by Goravel's elegant config architecture.

## Problem Statement

### Current Issues
1. **Excessive Boilerplate**: 730+ lines of getter/setter methods in `config_implementation.go`
2. **Maintenance Overhead**: Every new config key requires 2-4 new methods (getter, setter, interface definition)
3. **Code Duplication**: Similar patterns repeated across all config categories
4. **Limited Flexibility**: Cannot dynamically access nested config without predefined methods
5. **Testing Complexity**: Mocking requires implementing entire interface with 100+ methods

### Current Architecture
```go
// Adding a new config requires:
// 1. Add field to struct
type configImplementation struct {
    appName string
    appUrl  string
    // ... 50+ more fields
}

// 2. Add interface method
type appConfigInterface interface {
    SetAppName(string)
    GetAppName() string
    // ... repeated for every field
}

// 3. Implement getter
func (c *configImplementation) GetAppName() string {
    return c.appName
}

// 4. Implement setter
func (c *configImplementation) SetAppName(v string) {
    c.appName = v
}

// Usage:
cfg.SetAppName(env.GetString("APP_NAME"))
name := cfg.GetAppName()
```

## Proposed Solution

### Viper-Based Configuration
Adopt a Viper-based configuration system similar to Goravel, providing dynamic access with type safety.

### New Architecture
```go
// base/config/config.go
package config

import (
    "fmt"
    "github.com/spf13/viper"
)

type Config struct {
    v *viper.Viper
}

func New() *Config {
    v := viper.New()
    v.AutomaticEnv()
    return &Config{v: v}
}

// Type-safe accessors with optional defaults
func (c *Config) String(key string, defaultValue ...string) string {
    if !c.v.IsSet(key) && len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return c.v.GetString(key)
}

func (c *Config) Bool(key string, defaultValue ...bool) bool {
    if !c.v.IsSet(key) && len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return c.v.GetBool(key)
}

func (c *Config) Int(key string, defaultValue ...int) int {
    if !c.v.IsSet(key) && len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return c.v.GetInt(key)
}

// Validation helper
func (c *Config) MustString(key string) (string, error) {
    if !c.v.IsSet(key) {
        return "", fmt.Errorf("required config key missing: %s", key)
    }
    return c.v.GetString(key), nil
}

// Set values programmatically
func (c *Config) Set(key string, value any) {
    c.v.Set(key, value)
}

// Usage:
cfg := config.New()
name := cfg.String("APP_NAME")
debug := cfg.Bool("APP_DEBUG", false)
port := cfg.Int("APP_PORT", 8080)

// Validation
host, err := cfg.MustString("APP_HOST")
```

## Benefits

### 1. **Massive Code Reduction**
- **Before**: 730 lines in `config_implementation.go`
- **After**: ~100 lines for core config logic
- **Savings**: ~630 lines eliminated (86% reduction)

### 2. **Zero Boilerplate for New Config**
```go
// Before: 4 code changes required
// After: Just use it
apiKey := cfg.String("NEW_API_KEY")
```

### 3. **Built-in Features**
- Automatic environment variable binding
- Dot notation support: `cfg.String("database.host")`
- Type conversion with defaults
- Config file support (JSON, YAML, TOML, ENV)
- Watch for config changes (optional)

### 4. **Better Testing**
```go
// Before: Mock entire interface
type mockConfig struct {
    appName string
    appUrl  string
    // ... 50+ fields
}

// After: Simple setup
cfg := config.New()
cfg.Set("APP_NAME", "test-app")
```

### 5. **Consistency with Ecosystem**
- Viper is the de facto standard in Go (used by Hugo, Kubernetes, etc.)
- Familiar to Go developers
- Well-documented and maintained

## Migration Strategy

### Phase 1: Add Viper to Base Package
**Timeline**: 1-2 days

1. Create `base/config/config.go` with Viper wrapper
2. Add type-safe accessor methods
3. Add validation helpers
4. Write comprehensive tests
5. Document usage patterns

**Deliverables**:
- `base/config/config.go` - Core Viper wrapper
- `base/config/config_test.go` - Test coverage
- `base/config/README.md` - Usage documentation

### Phase 2: Parallel Implementation in Blueprint
**Timeline**: 2-3 days

1. Keep existing config system working
2. Add new Viper-based config alongside
3. Update `Load()` to populate both systems
4. Gradually migrate controllers to use new config

**Example**:
```go
// blueprint/internal/config/load.go
func Load() (ConfigInterface, error) {
    // Load env
    env.Load(".env")
    
    // Create new Viper config
    viperCfg := baseConfig.New()
    
    // Populate from env (automatic via AutomaticEnv)
    // Or manually set values
    viperCfg.Set("app.name", env.GetString("APP_NAME"))
    
    // Keep old config for backward compatibility
    oldCfg := New()
    oldCfg.SetAppName(env.GetString("APP_NAME"))
    
    return oldCfg, nil // Initially return old config
}
```

### Phase 3: Incremental Migration
**Timeline**: 1-2 weeks

1. Update controllers one by one
2. Replace `cfg.GetAppName()` with `cfg.String("APP_NAME")`
3. Run tests after each migration
4. Keep both systems until all code migrated

### Phase 4: Cleanup
**Timeline**: 1 day

1. Remove old config implementation
2. Remove old interfaces
3. Remove getter/setter methods
4. Update all documentation

## Breaking Changes

### API Changes
```go
// Before
name := cfg.GetAppName()
cfg.SetAppName("myapp")

// After
name := cfg.String("APP_NAME")
cfg.Set("APP_NAME", "myapp")
```

### Interface Changes
- `ConfigInterface` will be simplified or removed
- Dependency injection will use `*config.Config` directly

### Mitigation
- Gradual migration allows testing at each step
- Both systems can coexist during transition
- No external API changes (only internal refactoring)

## Risks & Mitigation

### Risk 1: Viper Dependency
**Impact**: Medium  
**Mitigation**: Viper is stable, widely-used, and actively maintained. It's already a transitive dependency in many Go projects.

### Risk 2: Type Safety Loss
**Impact**: Low  
**Mitigation**: Provide type-safe accessor methods (`String()`, `Bool()`, `Int()`) with compile-time checking.

### Risk 3: Migration Effort
**Impact**: Medium  
**Mitigation**: Incremental migration allows spreading work over time. Both systems can coexist.

### Risk 4: Breaking Existing Code
**Impact**: High  
**Mitigation**: Thorough testing at each phase. Keep old system until 100% migrated.

## Performance Considerations

### Viper Performance
- Viper uses sync.RWMutex for thread-safe access
- Negligible overhead for config reads (microseconds)
- Config is typically read at startup, not in hot paths

### Memory Impact
- Viper stores config in memory (map-based)
- Minimal memory overhead vs current struct-based approach
- No performance degradation expected

## Alternative Approaches Considered

### 1. Keep Current System
**Pros**: No migration effort, familiar to team  
**Cons**: Continues maintenance burden, boilerplate growth

### 2. Custom Dynamic Config
**Pros**: Full control, no external dependency  
**Cons**: Reinventing the wheel, maintenance burden

### 3. Struct Tags + Reflection
**Pros**: Type-safe, less boilerplate  
**Cons**: Complex, less flexible than Viper

## Success Criteria

1. ✅ Reduce config code by >80%
2. ✅ Zero boilerplate for new config keys
3. ✅ All existing tests pass
4. ✅ No performance degradation
5. ✅ Complete documentation
6. ✅ Team approval and buy-in

## Timeline

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| Phase 1: Base Package | 1-2 days | Viper wrapper in base |
| Phase 2: Parallel Implementation | 2-3 days | Both systems working |
| Phase 3: Incremental Migration | 1-2 weeks | All code migrated |
| Phase 4: Cleanup | 1 day | Old system removed |
| **Total** | **2-3 weeks** | **Production-ready** |

## Open Questions

1. Should we support config file watching for hot-reload?
2. Should we support nested config structures via dot notation?
3. Should we migrate all Dracory projects to this pattern?
4. Should validation be centralized in base or app-specific?

## References

- [Goravel Config Implementation](https://github.com/goravel/framework/blob/master/config/application.go)
- [Viper Documentation](https://github.com/spf13/viper)
- [Go Config Best Practices](https://github.com/golang-standards/project-layout)

## Approval

- [ ] Technical Lead Review
- [ ] Architecture Review
- [ ] Security Review
- [ ] Performance Review
- [ ] Team Consensus

---

**Author**: Cascade AI  
**Date**: 2026-03-23  
**Version**: 1.0
