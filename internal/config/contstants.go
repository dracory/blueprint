package config

// ============================================================================
// == START: Types
// ============================================================================

// AuthenticatedUserContextKey is a context key for the authenticated user.
type AuthenticatedUserContextKey struct{}

// AuthenticatedSessionContextKey is a context key for the authenticated session.
type AuthenticatedSessionContextKey struct{}

// ============================================================================
// == END: Types
// ============================================================================

// ============================================================================
// == START: AppEnvironment constants
// ============================================================================

const APP_ENVIRONMENT_DEVELOPMENT = "development"
const APP_ENVIRONMENT_LOCAL = "local"
const APP_ENVIRONMENT_PRODUCTION = "production"
const APP_ENVIRONMENT_STAGING = "staging"
const APP_ENVIRONMENT_TESTING = "testing"

const driverSQLite = "sqlite"

// ============================================================================
// == END: AppEnvironment constants
// ============================================================================

// ============================================================================
// == START: App Configurations
// ============================================================================

const KEY_APP_DEBUG = "APP_DEBUG"
const KEY_APP_ENVIRONMENT = "APP_ENV"
const KEY_APP_NAME = "APP_NAME"
const KEY_APP_URL = "APP_URL"
const KEY_APP_HOST = "APP_HOST"
const KEY_APP_PORT = "APP_PORT"

// ============================================================================
// == END: App Configurations
// ============================================================================

// ============================================================================
// == START: Database Configurations
// ============================================================================

const KEY_DB_DRIVER = "DB_DRIVER"
const KEY_DB_HOST = "DB_HOST"
const KEY_DB_PORT = "DB_PORT"
const KEY_DB_DATABASE = "DB_DATABASE"
const KEY_DB_USERNAME = "DB_USERNAME"
const KEY_DB_PASSWORD = "DB_PASSWORD"

// ============================================================================
// == END: Database Configurations
// ============================================================================

// ============================================================================
// == START: Mail Configurations
// ============================================================================

const KEY_MAIL_DRIVER = "MAIL_DRIVER"
const KEY_MAIL_HOST = "MAIL_HOST"
const KEY_MAIL_PASSWORD = "MAIL_PASSWORD"
const KEY_MAIL_PORT = "MAIL_PORT"
const KEY_MAIL_USERNAME = "MAIL_USERNAME"
const KEY_MAIL_FROM_ADDRESS = "MAIL_FROM_ADDRESS"
const KEY_MAIL_FROM_NAME = "MAIL_FROM_NAME"

// ============================================================================
// == END: Mail Configurations
// ============================================================================

// ============================================================================
// == START: Artifical Intelligence Configurations (LLM)
// ============================================================================
//
// This is where you can configure the artificial intelligence configurations.
//
// ============================================================================

// LLM: Anthropic
const KEY_ANTHROPIC_API_USED = "ANTHROPIC_API_USED"
const KEY_ANTHROPIC_API_KEY = "ANTHROPIC_API_KEY"
const KEY_ANTHROPIC_DEFAULT_MODEL = "ANTHROPIC_DEFAULT_MODEL"

// LLM: Google Gemini
const KEY_GEMINI_API_USED = "GEMINI_API_USED"
const KEY_GEMINI_API_KEY = "GEMINI_API_KEY"
const KEY_GEMINI_DEFAULT_MODEL = "GEMINI_DEFAULT_MODEL"

// LLM: OpenAI
const KEY_OPENAI_API_USED = "OPENAI_API_USED"
const KEY_OPENAI_API_KEY = "OPENAI_API_KEY"
const KEY_OPENAI_DEFAULT_MODEL = "OPENAI_DEFAULT_MODEL"

// LLM: OpenRouter
const KEY_OPENROUTER_API_USED = "OPENROUTER_API_USED"
const KEY_OPENROUTER_API_KEY = "OPENROUTER_API_KEY"
const KEY_OPENROUTER_DEFAULT_MODEL = "OPENROUTER_DEFAULT_MODEL"

// LLM: Vertex AI
const KEY_VERTEX_AI_USED = "VERTEX_AI_USED"
const KEY_VERTEX_MODEL_ID = "VERTEX_MODEL_ID"
const KEY_VERTEX_PROJECT_ID = "VERTEX_PROJECT_ID"
const KEY_VERTEX_REGION_ID = "VERTEX_REGION_ID"
const KEY_VERTEX_DEFAULT_MODEL = "VERTEX_DEFAULT_MODEL"

// ============================================================================
// == END: Artifical Intelligence Configurations (LLM)
// ============================================================================

// ============================================================================
// == START: Payment Configurations
// ============================================================================
//
// This is where you can configure the payment configurations.
//
// ============================================================================

const KEY_STRIPE_KEY_PRIVATE = "STRIPE_KEY_PRIVATE"
const KEY_STRIPE_KEY_PUBLIC = "STRIPE_KEY_PUBLIC"

// ============================================================================
// == END: Payment Configurations
// ============================================================================

// ============================================================================
// == START: Daily Analysis Configurations
// ============================================================================

const KEY_DAILY_ANALYSIS_SYMBOLS = "DAILY_ANALYSIS_SYMBOLS"
const KEY_DAILY_ANALYSIS_TIME_UTC = "DAILY_ANALYSIS_TIME_UTC"
const KEY_DAILY_ANALYSIS_CADENCE_HOURS = "DAILY_ANALYSIS_CADENCE_HOURS"

// ============================================================================
// == END: Daily Analysis Configurations
// ============================================================================

// ============================================================================
// == START: i18n Configurations
// ============================================================================

const KEY_TRANSLATION_LANGUAGE_DEFAULT = "TRANSLATION_LANGUAGE_DEFAULT"

// ============================================================================
// == END: i18n Configurations
// ============================================================================

// ============================================================================
// == START: EnvEnc Configurations
// ============================================================================
//
// This is where you can configure the EnvEnc encryption key.
//
// ============================================================================
const KEY_ENV_ENCRYPTION_KEY = "ENV_ENCRYPTION_KEY"
const KEY_ENVENC_KEY_PRIVATE = "ENVENC_KEY_PRIVATE"

// ============================================================================
// == END: EnvEnc Configurations
// ============================================================================

// ============================================================================
// == START: Store Configurations
// ============================================================================
//
// This is where you can configure the stores used by the application.
//

// KEY_CACHE_STORE_USED toggles cache store bootstrapping. When true, the
// application expects backing cache infrastructure to be reachable.
const KEY_CACHE_STORE_USED = "CACHE_STORE_USED"

// KEY_CMS_STORE_USED enables the CMS store and requires related templates and
// backing tables so the CMS module can respond to requests.
const KEY_CMS_STORE_USED = "CMS_STORE_USED"

// KEY_CMS_STORE_TEMPLATE_ID identifies the CMS template to hydrate when the
// CMS store is enabled.
const KEY_CMS_STORE_TEMPLATE_ID = "CMS_STORE_TEMPLATE_ID"

// KEY_CUSTOM_STORE_USED gates initialization of custom store resources and any
// dependent background jobs.
const KEY_CUSTOM_STORE_USED = "CUSTOM_STORE_USED"

// KEY_ENTITY_STORE_USED toggles domain entity persistence. When enabled, entity
// migrations must be applied before startup.
const KEY_ENTITY_STORE_USED = "ENTITY_STORE_USED"

// KEY_FEED_STORE_USED enables feed processing pipelines and the database
// structures that back them.
const KEY_FEED_STORE_USED = "FEED_STORE_USED"

// KEY_GEO_STORE_USED activates geographic data hydration and requires region
// lookup tables to exist.
const KEY_GEO_STORE_USED = "GEO_STORE_USED"

// KEY_LOG_STORE_USED controls persistence of structured logs. When true, the
// log store tables are expected to be present for ingestion.
const KEY_LOG_STORE_USED = "LOG_STORE_USED"

// KEY_META_STORE_USED toggles metadata storage. Enabling it means metadata
// tables will be touched during initialization.
const KEY_META_STORE_USED = "META_STORE_USED"

// KEY_SESSION_STORE_USED activates the session store. When enabled, session
// tables must be migrated to avoid authentication failures.
const KEY_SESSION_STORE_USED = "SESSION_STORE_USED"

// KEY_SETTING_STORE_USED toggles application setting synchronization and
// expects settings tables to be available.
const KEY_SETTING_STORE_USED = "SETTING_STORE_USED"

// KEY_SHOP_STORE_USED enables commerce-related database entities and services.
// Ensure shop migrations run before enabling it in production.
const KEY_SHOP_STORE_USED = "SHOP_STORE_USED"

// KEY_STATS_STORE_USED controls analytics/statistics aggregation stores. When
// enabled, reporting jobs will read/write supporting tables.
const KEY_STATS_STORE_USED = "STATS_STORE_USED"

// KEY_TASK_STORE_USED toggles the task orchestration store and requires task
// queues to be reachable.
const KEY_TASK_STORE_USED = "TASK_STORE_USED"

// KEY_USER_STORE_USED activates the user store. User authentication and profile
// management will fail if the necessary tables are missing.
const KEY_USER_STORE_USED = "USER_STORE_USED"

// KEY_VAULT_STORE_USED toggles secret vault storage. When true, vault keys and
// encrypted records must be provisioned.
const KEY_VAULT_STORE_USED = "VAULT_STORE_USED"

// KEY_VAULT_STORE_KEY supplies the encryption key required when the vault store
// is enabled.
const KEY_VAULT_STORE_KEY = "VAULT_STORE_KEY"

// ============================================================================
// == END: Store Configurations
// ============================================================================
