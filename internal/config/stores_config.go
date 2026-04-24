package config

import "fmt"

// ============================================================================
// == START: Enabled Database Stores
// ============================================================================
//
// This is where you can configure which database stores will be enabled
// an available through the registry.
//
// ============================================================================

// auditStoreUsed enables / disables the audit store responsible for
// recording change history.
// Enable it when audit logging should be persisted.
const auditStoreUsed = false

// blogStoreUsed enables / disables the blog store.
// Enable it when the application should expose blog content backed by blogstore tables.
const blogStoreUsed = false

// cacheStoreUsed enables / disables cache store bootstrapping.
// Enable it when you need database-based caching.
const cacheStoreUsed = true

// chatStoreUsed enables / disables chat store.
const chatStoreUsed = false

// cmsStoreUsed enables / disables the CMS store and requires related templates and
// backing tables so the CMS module can respond to requests.
const cmsStoreUsed = false

// customStoreUsed enables / disables initialization of custom store resources and any
// dependent background jobs.
const customStoreUsed = false

// entityStoreUsed enables / disables domain entity persistence. When enabled, entity
// migrations must be applied before startup.
const entityStoreUsed = false

// feedStoreUsed enables / disables feed processing pipelines and the database
// structures that back them.
const feedStoreUsed = false

// geoStoreUsed enables / disables geographic data hydration and requires region
// lookup tables to exist.
const geoStoreUsed = true

// logStoreUsed enables / disables persistence of structured logs. When true, the
// log store tables are expected to be present for ingestion.
const logStoreUsed = true

// metaStoreUsed enables / disables metadata storage. Enabling it means metadata
// tables will be touched during initialization.
const metaStoreUsed = false

// sessionStoreUsed enables / disables the session store. When enabled, session
// tables must be migrated to avoid authentication failures.
const sessionStoreUsed = true

// settingStoreUsed enables / disables application setting synchronization and
// expects settings tables to be available.
const settingStoreUsed = false

// shopStoreUsed enables / disables commerce-related database entities and services.
// Ensure shop migrations run before enabling it in production.
const shopStoreUsed = false

// sqlFileStoreUsed enables / disables the SQL-backed file storage. Enable it when
// uploads should be persisted via `filesystem.DRIVER_SQL` tables.
const sqlFileStoreUsed = false

// statsStoreUsed enables / disables analytics/statistics aggregation stores. When
// enabled, reporting jobs will read/write supporting tables.
const statsStoreUsed = false

// subscriptionStoreUsed enables / disables subscription store bootstrapping. Enable
// it when subscription plans and billing data should be managed.
const subscriptionStoreUsed = false

// taskStoreUsed enables / disables the task orchestration store and requires task
// queues to be reachable.
const taskStoreUsed = true

// userStoreUsed enables / disables the user store. User authentication and profile
// management will fail if the necessary tables are missing.
const userStoreUsed = true

// userStoreVaultEnabled enables / disables the user store vault
const userStoreVaultEnabled = false

// vaultStoreUsed enables / disables secret vault storage. When true, vault keys and
// encrypted records must be provisioned.
const vaultStoreUsed = false

// ============================================================================
// == END: Enabled Database Stores
// ============================================================================

// storesConfig reads datastore feature flags from environment variables.
// Each store is opt-in via configuration_stores.go - set the flag to true to enable it.
func storesConfig(env *envValidator) storesSettings {
	// CMS Store Template ID
	//
	// The template ID used by the CMS store for rendering content.
	// Required when CMS_STORE_USED is true.
	cmsStoreTemplateID := env.GetString(KEY_CMS_STORE_TEMPLATE_ID)

	// Vault Store Key
	//
	// Encryption key used by the vault store to encrypt sensitive data.
	// Required when VAULT_STORE_USED is true.
	vaultStoreKey := env.GetString(KEY_VAULT_STORE_KEY)

	if userStoreVaultEnabled && !vaultStoreUsed {
		env.Add(fmt.Errorf("userStoreVaultEnabled requires vaultStoreUsed to be true"))
	}

	env.RequireWhen(cmsStoreUsed, KEY_CMS_STORE_TEMPLATE_ID,
		"required when `CMS_STORE_USED` is true", cmsStoreTemplateID)

	env.RequireWhen(vaultStoreUsed, KEY_VAULT_STORE_KEY,
		"required when `VAULT_STORE_USED` is true", vaultStoreKey)

	return storesSettings{cmsStoreTemplateID: cmsStoreTemplateID, vaultStoreKey: vaultStoreKey}
}

type storesSettings struct {
	cmsStoreTemplateID string
	vaultStoreKey      string
}
