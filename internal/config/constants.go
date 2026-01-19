package config

// ============================================================================
// == START: Types
// ============================================================================

// AuthenticatedUserContextKey is a context key for the authenticated user.
type AuthenticatedUserContextKey struct{}

// AuthenticatedSessionContextKey is a context key for the authenticated session.
type AuthenticatedSessionContextKey struct{}

// APIAuthenticatedUserContextKey is a context key for API authenticated user.
type APIAuthenticatedUserContextKey struct{}

// APIAuthenticatedSessionContextKey is a context key for API authenticated session.
type APIAuthenticatedSessionContextKey struct{}

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
const KEY_DB_SSL_MODE = "DB_SSL_MODE"

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
// == START: Auth Configurations
// ============================================================================

const KEY_AUTH_REGISTRATION_ENABLED = "AUTH_REGISTRATION_ENABLED"

// ============================================================================
// == END: Auth Configurations
// ============================================================================

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
const KEY_ANTHROPIC_API_DEFAULT_MODEL = "ANTHROPIC_DEFAULT_MODEL"

// LLM: Google Gemini
const KEY_GEMINI_API_USED = "GEMINI_API_USED"
const KEY_GEMINI_API_KEY = "GEMINI_API_KEY"
const KEY_GEMINI_API_DEFAULT_MODEL = "GEMINI_API_DEFAULT_MODEL"

// LLM: OpenAI
const KEY_OPENAI_API_USED = "OPENAI_API_USED"
const KEY_OPENAI_API_KEY = "OPENAI_API_KEY"
const KEY_OPENAI_API_DEFAULT_MODEL = "OPENAI_API_DEFAULT_MODEL"

// LLM: OpenRouter
const KEY_OPENROUTER_API_USED = "OPENROUTER_API_USED"
const KEY_OPENROUTER_API_KEY = "OPENROUTER_API_KEY"
const KEY_OPENROUTER_API_DEFAULT_MODEL = "OPENROUTER_API_DEFAULT_MODEL"

// LLM: Vertex AI
const KEY_VERTEX_AI_API_USED = "VERTEX_AI_API_USED"
const KEY_VERTEX_AI_API_MODEL_ID = "VERTEX_AI_API_MODEL_ID"
const KEY_VERTEX_AI_API_PROJECT_ID = "VERTEX_AI_API_PROJECT_ID"
const KEY_VERTEX_AI_API_REGION_ID = "VERTEX_AI_API_REGION_ID"
const KEY_VERTEX_AI_API_DEFAULT_MODEL = "VERTEX_AI_API_DEFAULT_MODEL"

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
const KEY_ENVENC_USED = "ENVENC_USED"

const KEY_ENVENC_KEY_PRIVATE = "ENVENC_KEY_PRIVATE"

/**
 * ENVENC_KEY_PUBLIC: Public Key for Envenc Vault Encryption
 *
 * This constant stores the public key used to encrypt your envenc vault file.
 * * **Important Security Information:**
 * - This public key, combined with a corresponding private key, is used
 *   as input to a secure **one-way** hashing function to derive the final
 *   encryption key.
 * - Both the private and public keys must be at least 32-character strings,
 *   composed of randomly generated characters, numbers, and symbols.
 * - **DO NOT store the actual final key anywhere.** It should be generated dynamically when needed.
 * - **DO NOT directly commit the actual PRIVATE key to version control.** Use environment variables or secure configuration management.
 * - Replace "YOUR_PUBLIC_KEY" with your actual 32-character public key.
 * - The associated private key must be kept extremely secure.
 * - Ensure that the random number generator used to create the keys is cryptographically secure (CSPRNG).
 * - **Ideally, the public key should be obfuscated. See envenc for more details.**
 *
 * Example:
 * const ENVENC_KEY_PUBLIC = "aBcD123$456!eFgH789%iJkL0mNoPqRsTuVwXyZ"; // Replace with your actual key
 */
const ENVENC_KEY_PUBLIC = "YOUR_PUBLIC_KEY"

// ============================================================================
// == END: EnvEnc Configurations
// ============================================================================

// ============================================================================
// == START: Store Configurations
// ============================================================================
//
// This is where you can configure the stores used by the application.
//

// KEY_CMS_STORE_TEMPLATE_ID identifies the CMS template to hydrate when the
// CMS store is enabled.
const KEY_CMS_STORE_TEMPLATE_ID = "CMS_STORE_TEMPLATE_ID"

// KEY_VAULT_STORE_KEY supplies the encryption key required when the vault store
// is enabled.
const KEY_VAULT_STORE_KEY = "VAULT_STORE_KEY"

// ============================================================================
// == END: Store Configurations
// ============================================================================
