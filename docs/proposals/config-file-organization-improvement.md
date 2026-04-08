# Proposal: Configuration File Organization Improvement

## Status
**IMPLEMENTED** - 2026-04-08

## Overview
Improve the organization and discoverability of Blueprint's configuration system by restructuring config files into logical, domain-focused modules inspired by Goravel's clean separation pattern. This proposal focuses on practical structural improvements while preserving the existing interface-based architecture.

## Problem Statement

### Current State
Blueprint's configuration is already well-architected with:
- Clean interface-based design (`config_interface.go`)
- Separated loader functions per domain (`app.go`, `database.go`, `llm.go`, etc.)
- Type-safe getter/setter implementation
- Excellent registry/config separation

**However, there are organizational challenges:**

1. **Monolithic Implementation File**: `config_implementation.go` contains 730+ lines with all getters/setters for every config domain mixed together
2. **Interface Fragmentation**: `config_interface.go` has 20+ small interfaces that could be better organized
3. **Discoverability**: Finding config for a specific domain requires scanning through large files
4. **Maintenance Overhead**: Adding new config to a domain requires editing multiple large files
5. **Cognitive Load**: Developers must mentally map between loader files, interface fragments, and implementation methods

### What Works Well (Keep These)
- ✅ Interface-based design with dependency injection
- ✅ Separate loader functions per domain
- ✅ Type-safe getters/setters (no Viper/dynamic config)
- ✅ Clear separation between config and registry
- ✅ Environment variable loading with validation

## Proposed Solution

### Inspiration: Goravel's Organization
Goravel organizes config into focused, self-contained files:
```
config/
├── app.go          # Application settings
├── database.go     # Database configuration
├── cache.go        # Cache settings
├── mail.go         # Email configuration
├── queue.go        # Queue settings
└── ...
```

Each file is complete and self-contained, making it easy to find and modify related configuration.

### Proposed Structure for Blueprint

Reorganize config into domain-focused modules where each domain has its own file containing:
- Interface definition
- Implementation (getters/setters)
- Loader function
- Related types and constants

```
internal/config/
├── config.go                    # Main ConfigInterface composition + New()
├── load.go                      # Main Load() orchestration
├── constants.go                 # All environment variable key constants
├── defaults.go                  # Default value functions
│
├── app_config.go                # App domain (complete)
├── database_config.go           # Database domain (complete)
├── email_config.go              # Email/mail domain (complete)
├── llm_config.go                # LLM providers domain (complete)
├── media_config.go              # Media/storage domain (complete)
├── payment_config.go            # Payment/Stripe domain (complete)
├── auth_config.go               # Authentication domain (complete)
├── i18n_config.go               # Translation/i18n domain (complete)
├── seo_config.go                # SEO domain (complete)
├── encryption_config.go         # Encryption domain (complete)
│
└── stores_config.go             # All data stores (complete)
```

### Example: app_config.go (Complete Domain Module)

```go
package config

import (
	"strings"
	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// ============================================================================
// Interface
// ============================================================================

// AppConfigInterface defines application-level configuration methods.
type AppConfigInterface interface {
	SetAppName(string)
	GetAppName() string

	SetAppType(string)
	GetAppType() string

	SetAppEnv(string)
	GetAppEnv() string

	SetAppHost(string)
	GetAppHost() string

	SetAppPort(string)
	GetAppPort() string

	SetAppUrl(string)
	GetAppUrl() string

	SetAppDebug(bool)
	GetAppDebug() bool

	// Environment helpers
	IsEnvDevelopment() bool
	IsEnvLocal() bool
	IsEnvProduction() bool
	IsEnvStaging() bool
	IsEnvTesting() bool
}

// ============================================================================
// Types
// ============================================================================

// appConfig captures application-level settings.
type appConfig struct {
	name         string // Application name identifier
	url          string // Base URL for the application
	host         string // Host address for the server
	port         string // Port number for the server
	env          string // Environment (development, staging, production)
	debug        bool   // Debug mode flag
	cmsMcpApiKey string // CMS MCP API key for integration
}

// ============================================================================
// Loader
// ============================================================================

// loadAppConfig loads application configuration from environment variables.
func loadAppConfig(acc *baseCfg.LoadAccumulator) appConfig {
	mcpApiKey := strings.TrimSpace(env.GetString(KEY_MCP_API_KEY))

	return appConfig{
		name:         env.GetString(KEY_APP_NAME),
		url:          env.GetString(KEY_APP_URL),
		host:         acc.MustString(KEY_APP_HOST, "set the application host address"),
		port:         acc.MustString(KEY_APP_PORT, "set the application port"),
		env:          acc.MustString(KEY_APP_ENVIRONMENT, "set the application environment"),
		debug:        env.GetBool(KEY_APP_DEBUG),
		cmsMcpApiKey: mcpApiKey,
	}
}

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetAppName(appName string) {
	c.appName = appName
}

func (c *configImplementation) GetAppName() string {
	return c.appName
}

func (c *configImplementation) SetAppType(appType string) {
	c.appType = appType
}

func (c *configImplementation) GetAppType() string {
	return c.appType
}

func (c *configImplementation) SetAppEnv(appEnv string) {
	c.appEnv = appEnv
}

func (c *configImplementation) GetAppEnv() string {
	return c.appEnv
}

func (c *configImplementation) SetAppHost(appHost string) {
	c.appHost = appHost
}

func (c *configImplementation) GetAppHost() string {
	return c.appHost
}

func (c *configImplementation) SetAppPort(appPort string) {
	c.appPort = appPort
}

func (c *configImplementation) GetAppPort() string {
	return c.appPort
}

func (c *configImplementation) SetAppUrl(appUrl string) {
	c.appUrl = appUrl
}

func (c *configImplementation) GetAppUrl() string {
	return c.appUrl
}

func (c *configImplementation) SetAppDebug(appDebug bool) {
	c.appDebug = appDebug
}

func (c *configImplementation) GetAppDebug() bool {
	return c.appDebug
}

// Environment Helpers
func (c *configImplementation) IsEnvDevelopment() bool {
	return c.appEnv == "development"
}

func (c *configImplementation) IsEnvLocal() bool {
	return c.appEnv == "local"
}

func (c *configImplementation) IsEnvProduction() bool {
	return c.appEnv == "production"
}

func (c *configImplementation) IsEnvStaging() bool {
	return c.appEnv == "staging"
}

func (c *configImplementation) IsEnvTesting() bool {
	return c.appEnv == "testing"
}
```

### Example: config.go (Main Composition)

```go
package config

// ConfigInterface defines the contract for application configuration.
// It composes all domain-specific configuration interfaces.
type ConfigInterface interface {
	AppConfigInterface
	DatabaseConfigInterface
	EmailConfigInterface
	LLMConfigInterface
	MediaConfigInterface
	PaymentConfigInterface
	AuthConfigInterface
	I18nConfigInterface
	SEOConfigInterface
	EncryptionConfigInterface
	
	// CMS MCP (special case - small interface)
	SetCmsMcpApiKey(string)
	GetCmsMcpApiKey() string

	// Datastores
	AuditStoreConfigInterface
	BlogStoreConfigInterface
	CacheStoreConfigInterface
	ChatStoreConfigInterface
	CmsStoreConfigInterface
	CustomStoreConfigInterface
	EntityStoreConfigInterface
	FeedStoreConfigInterface
	GeoStoreConfigInterface
	LogStoreConfigInterface
	MetaStoreConfigInterface
	SessionStoreConfigInterface
	SettingStoreConfigInterface
	ShopStoreConfigInterface
	SqlFileStoreConfigInterface
	StatsStoreConfigInterface
	SubscriptionStoreConfigInterface
	TaskStoreConfigInterface
	UserStoreConfigInterface
	VaultStoreConfigInterface
}

// configImplementation holds all configuration values.
type configImplementation struct {
	// App configuration
	appName  string
	appType  string
	appEnv   string
	appHost  string
	appPort  string
	appUrl   string
	appDebug bool

	// Database configuration
	databaseDriver   string
	databaseHost     string
	databasePort     string
	databaseName     string
	databaseUsername string
	databasePassword string
	databaseSSLMode  string

	// Email configuration
	emailDriver      string
	emailHost        string
	emailPort        int
	emailUsername    string
	emailPassword    string
	emailFromName    string
	emailFromAddress string

	// LLM configuration
	openRouterApiKey          string
	openRouterApiUsed         bool
	openRouterApiDefaultModel string
	openAiApiKey              string
	openAiApiUsed             bool
	openAiApiDefaultModel     string
	anthropicApiUsed          bool
	anthropicApiKey           string
	anthropicApiDefaultModel  string
	googleGeminiApiUsed       bool
	googleGeminiApiKey        string
	googleGeminiApiDefaultModel string
	vertexAiApiUsed           bool
	vertexAiApiDefaultModel   string
	vertexAiApiProjectID      string
	vertexAiApiRegionID       string
	vertexAiApiModelID        string

	// Media configuration
	mediaBucket   string
	mediaDriver   string
	mediaKey      string
	mediaEndpoint string
	mediaRegion   string
	mediaRoot     string
	mediaSecret   string
	mediaUrl      string

	// Payment configuration
	stripeKeyPrivate string
	stripeKeyPublic  string
	stripeUsed       bool

	// Authentication
	registrationEnabled bool

	// i18n / Translation
	translationLanguageDefault string
	translationLanguageList    map[string]string

	// SEO configuration
	indexNowKey string

	// Encryption
	envEncryptionKey string

	// CMS MCP
	cmsMcpApiKey string

	// Store flags
	auditStoreUsed        bool
	blogStoreUsed         bool
	chatStoreUsed         bool
	cacheStoreUsed        bool
	cmsStoreUsed          bool
	cmsStoreTemplateID    string
	customStoreUsed       bool
	entityStoreUsed       bool
	feedStoreUsed         bool
	geoStoreUsed          bool
	logStoreUsed          bool
	metaStoreUsed         bool
	sessionStoreUsed      bool
	settingStoreUsed      bool
	shopStoreUsed         bool
	sqlFileStoreUsed      bool
	statsStoreUsed        bool
	subscriptionStoreUsed bool
	taskStoreUsed         bool
	userStoreUsed         bool
	userStoreVaultEnabled bool
	vaultStoreUsed        bool
	vaultStoreKey         string
}

// Ensure configImplementation satisfies ConfigInterface
var _ ConfigInterface = (*configImplementation)(nil)

// New constructs a new configuration instance.
func New() ConfigInterface {
	return &configImplementation{}
}

// CMS MCP methods (small interface, kept here)
func (c *configImplementation) SetCmsMcpApiKey(v string) {
	c.cmsMcpApiKey = v
}

func (c *configImplementation) GetCmsMcpApiKey() string {
	return c.cmsMcpApiKey
}
```

## Benefits

### 1. Improved Discoverability
- **Before**: Search through 730-line file to find database config methods
- **After**: Open `database_config.go` - everything is there

### 2. Better Organization
- Each domain is self-contained in one file
- Interface, implementation, loader, and types together
- Easy to understand the complete picture of a domain

### 3. Easier Maintenance
- Adding new config to a domain: edit one file
- All related code in one place
- Reduced chance of missing updates

### 4. Reduced Cognitive Load
- No need to jump between 3-4 files to understand one domain
- Clear file names indicate content
- Smaller, focused files are easier to read

### 5. Better Code Reviews
- Changes to a domain are isolated to one file
- Easier to review and understand impact
- Clear diff boundaries

### 6. Preserved Architecture
- ✅ Interface-based design maintained
- ✅ Type safety preserved
- ✅ Dependency injection unchanged
- ✅ Registry separation intact
- ✅ No breaking changes to consumers

## File Size Comparison

### Before
```
config_interface.go       157 lines (20+ small interfaces)
config_implementation.go  730 lines (all getters/setters)
app.go                     47 lines (loader only)
database.go                52 lines (loader only)
llm.go                    102 lines (loader only)
mail.go                    28 lines (loader only)
... (other loaders)
```

### After
```
config.go                 120 lines (main composition + struct)
load.go                   150 lines (orchestration)
constants.go               80 lines (all keys)
defaults.go                30 lines (defaults)

app_config.go             150 lines (complete domain)
database_config.go        140 lines (complete domain)
email_config.go           130 lines (complete domain)
llm_config.go             280 lines (complete domain)
media_config.go           120 lines (complete domain)
payment_config.go          80 lines (complete domain)
auth_config.go             50 lines (complete domain)
i18n_config.go             70 lines (complete domain)
seo_config.go              40 lines (complete domain)
encryption_config.go       60 lines (complete domain)
stores_config.go          400 lines (all stores complete)
```

**Key Insight**: Same total lines, but organized into logical, discoverable modules.

## Migration Strategy

### Phase 1: Create New Structure (2-3 days)

1. **Create domain config files** with complete implementations
2. **Move interface definitions** from `config_interface.go` to domain files
3. **Move getters/setters** from `config_implementation.go` to domain files
4. **Keep loaders** in domain files (already separate)
5. **Update `config.go`** to compose interfaces
6. **Keep `load.go`** mostly unchanged (orchestration)

### Phase 2: Testing & Validation (1 day)

1. **Run all existing tests** - should pass without changes
2. **Verify no breaking changes** to consumers
3. **Check imports** - may need to update some internal imports
4. **Validate IDE support** - autocomplete should work identically

### Phase 3: Documentation (1 day)

1. **Update config documentation** to reflect new structure
2. **Add examples** showing how to find and modify config
3. **Document conventions** for adding new config domains
4. **Update AGENTS.md** if needed

## Implementation Details

### File Organization Rules

1. **One domain per file**: Each config domain gets its own file
2. **Complete modules**: Interface + implementation + loader + types in same file
3. **Clear naming**: `{domain}_config.go` pattern
4. **Logical grouping**: Related config stays together (e.g., all LLM providers)
5. **Size guideline**: Aim for 50-300 lines per file (readable in one screen)

### Domain Identification

**Core Application Domains:**
- `app_config.go` - Application settings (name, env, host, port, URL, debug)
- `database_config.go` - Database connection settings
- `email_config.go` - Email/mail delivery settings
- `auth_config.go` - Authentication settings (registration, etc.)

**External Service Domains:**
- `llm_config.go` - All LLM providers (OpenAI, Anthropic, Gemini, Vertex, OpenRouter)
- `media_config.go` - Media storage settings (S3, local, etc.)
- `payment_config.go` - Payment provider settings (Stripe)
- `seo_config.go` - SEO tools (IndexNow)

**Feature Domains:**
- `i18n_config.go` - Translation/internationalization
- `encryption_config.go` - Encryption keys and settings
- `stores_config.go` - All data store enablement flags

### Backward Compatibility

**No Breaking Changes:**
- `ConfigInterface` remains the same (just composed differently)
- All getter/setter methods unchanged
- `Load()` function signature unchanged
- Registry integration unchanged
- Consumer code unchanged

**Internal Changes Only:**
- File organization changes
- Import paths may change within config package
- Implementation details reorganized

## Comparison with Goravel

### What We're Adopting from Goravel
✅ **File-per-domain organization** - Clear, focused files  
✅ **Self-contained modules** - Everything for a domain in one place  
✅ **Logical naming** - Clear, predictable file names  
✅ **Discoverability** - Easy to find what you need

### What We're Keeping from Blueprint
✅ **Interface-based design** - Type-safe, explicit contracts  
✅ **Getter/setter pattern** - Clear, predictable API  
✅ **Separate loaders** - Clean separation of concerns  
✅ **Registry separation** - Config vs. runtime dependencies  
✅ **Validation** - LoadAccumulator pattern for errors

### Why Not Full Goravel Approach?

Goravel uses a simpler, more dynamic approach:
```go
// Goravel style
config.Get("app.name")
config.GetString("database.host", "localhost")
```

Blueprint's interface-based approach provides:
- **Compile-time type safety** - Typos caught at compile time
- **IDE autocomplete** - Full method discovery
- **Explicit contracts** - Clear API surface
- **Refactoring support** - Rename methods safely
- **Better for large teams** - Explicit is better than implicit

**Decision**: Keep Blueprint's type-safe approach, adopt Goravel's organization.

## Risks and Mitigations

### Risk 1: Import Path Changes
**Impact**: Low  
**Mitigation**: All imports are internal to config package. External consumers use `config.ConfigInterface` which doesn't change.

### Risk 2: Merge Conflicts
**Impact**: Medium (during migration)  
**Mitigation**: Complete migration in one PR. Coordinate with team to avoid concurrent config changes.

### Risk 3: File Size Imbalance
**Impact**: Low  
**Mitigation**: Some domains (LLM, stores) will be larger. This is acceptable - they're still more focused than current monolithic files.

### Risk 4: Learning Curve
**Impact**: Very Low  
**Mitigation**: New structure is more intuitive. Documentation will guide developers.

## Success Criteria

1. ✅ All existing tests pass without modification
2. ✅ No breaking changes to consumers
3. ✅ Each domain config in single, focused file
4. ✅ Improved discoverability (measured by developer feedback)
5. ✅ Easier maintenance (measured by time to add new config)
6. ✅ Complete documentation
7. ✅ Team approval

## Timeline

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| Phase 1: Restructure | 2-3 days | New file organization |
| Phase 2: Testing | 1 day | All tests passing |
| Phase 3: Documentation | 1 day | Updated docs |
| **Total** | **4-5 days** | **Production-ready** |

## Alternative Approaches Considered

### Alternative 1: Keep Current Structure
**Pros**: No migration effort, familiar to team  
**Cons**: Continues discoverability and maintenance challenges

### Alternative 2: Viper-Based Dynamic Config
**Pros**: Less boilerplate, more flexible  
**Cons**: Loses type safety, IDE support, compile-time checking  
**Status**: Already rejected in previous proposal

### Alternative 3: Partial Reorganization
**Pros**: Lower effort, incremental improvement  
**Cons**: Inconsistent structure, doesn't fully solve problem

**Decision**: Full reorganization (Alternative 4) provides best balance of benefits vs. effort.

## Open Questions

1. Should we group all store configs together or separate by category?
   - **Recommendation**: Keep together in `stores_config.go` - they're all similar
2. Should constants stay in separate file or move to domain files?
   - **Recommendation**: Keep in `constants.go` - easier to see all env vars at once
3. Should we create subdirectories for major domains?
   - **Recommendation**: No - flat structure is simpler for this size

## References

- [Goravel Config Structure](https://github.com/goravel/goravel/tree/v1.17.x/config)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Blueprint Config Package](internal/config/)

## Approval

- [ ] Technical Lead Review
- [ ] Architecture Review
- [ ] Team Consensus
- [ ] Implementation Plan Approved

---

**Author**: Kiro AI  
**Date**: 2026-04-08  
**Version**: 1.0
