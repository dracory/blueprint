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
const KEY_EMAIL_FROM_ADDRESS = "EMAIL_FROM_ADDRESS"
const KEY_EMAIL_FROM_NAME = "EMAIL_FROM_NAME"
const KEY_MAIL_HOST = "MAIL_HOST"
const KEY_MAIL_PASSWORD = "MAIL_PASSWORD"
const KEY_MAIL_PORT = "MAIL_PORT"
const KEY_MAIL_USERNAME = "MAIL_USERNAME"

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
// ============================================================================

const KEY_CACHE_STORE_USED = "CACHE_STORE_USED"
const KEY_CMS_STORE_USED = "CMS_STORE_USED"
const KEY_CMS_STORE_TEMPLATE_ID = "CMS_STORE_TEMPLATE_ID"
const KEY_CUSTOM_STORE_USED = "CUSTOM_STORE_USED"
const KEY_ENTITY_STORE_USED = "ENTITY_STORE_USED"
const KEY_FEED_STORE_USED = "FEED_STORE_USED"
const KEY_GEO_STORE_USED = "GEO_STORE_USED"
const KEY_LOG_STORE_USED = "LOG_STORE_USED"
const KEY_META_STORE_USED = "META_STORE_USED"
const KEY_SESSION_STORE_USED = "SESSION_STORE_USED"
const KEY_SETTING_STORE_USED = "SETTING_STORE_USED"
const KEY_SHOP_STORE_USED = "SHOP_STORE_USED"
const KEY_STATS_STORE_USED = "STATS_STORE_USED"
const KEY_TASK_STORE_USED = "TASK_STORE_USED"
const KEY_TRADING_STORE_USED = "TRADING_STORE_USED"
const KEY_USER_STORE_USED = "USER_STORE_USED"
const KEY_VAULT_STORE_USED = "VAULT_STORE_USED"
const KEY_VAULT_STORE_KEY = "VAULT_STORE_KEY"

// ============================================================================
// == END: Store Configurations
// ============================================================================
