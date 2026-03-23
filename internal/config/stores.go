package config

import (
	"fmt"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// storesConfig captures feature store toggles.
// It manages the enablement/disablement of various data stores and services
// in the application, allowing for modular feature activation.
type storesConfig struct {
	auditStoreUsed        bool   // Enable audit logging store
	blogStoreUsed         bool   // Enable blog content store
	cacheStoreUsed        bool   // Enable caching store
	cmsStoreUsed          bool   // Enable CMS content store
	cmsStoreTemplateID    string // Default template ID for CMS store
	customStoreUsed       bool   // Enable custom data store
	entityStoreUsed       bool   // Enable entity management store
	feedStoreUsed         bool   // Enable RSS/Atom feed store
	geoStoreUsed          bool   // Enable geolocation data store
	logStoreUsed          bool   // Enable application logging store
	metaStoreUsed         bool   // Enable metadata store
	sessionStoreUsed      bool   // Enable user session store
	settingStoreUsed      bool   // Enable application settings store
	shopStoreUsed         bool   // Enable e-commerce store
	sqlFileStoreUsed      bool   // Enable SQL file storage store
	statsStoreUsed        bool   // Enable analytics/statistics store
	subscriptionStoreUsed bool   // Enable subscription management store
	taskStoreUsed         bool   // Enable background task store
	userStoreUsed         bool   // Enable user management store
	userStoreVaultEnabled bool   // Enable vault encryption for user store
	vaultStoreUsed        bool   // Enable secure vault store
	vaultStoreKey         string // Encryption key for vault store
}

// loadStoresConfig loads stores configuration from environment variables.
// It validates store dependencies and requirements, ensuring that dependent
// stores are properly configured. For example, CMS store requires a template ID,
// and user vault encryption requires the vault store to be enabled.
//
// Parameters:
//   - acc: LoadAccumulator for collecting validation errors and dependency checks
//
// Returns:
//   - storesConfig: Populated configuration struct with store enablement flags
func loadStoresConfig(acc *baseCfg.LoadAccumulator) storesConfig {
	cmsStoreTemplateID := env.GetString(KEY_CMS_STORE_TEMPLATE_ID)
	vaultStoreKey := env.GetString(KEY_VAULT_STORE_KEY)

	if userStoreVaultEnabled && !vaultStoreUsed {
		acc.Add(fmt.Errorf("%v requires %v to be true", userStoreVaultEnabled, vaultStoreUsed))
	}

	acc.MustWhen(cmsStoreUsed, KEY_CMS_STORE_TEMPLATE_ID, "required when `CMS_STORE_USED` is true", cmsStoreTemplateID)

	return storesConfig{
		auditStoreUsed:        auditStoreUsed,
		blogStoreUsed:         blogStoreUsed,
		cacheStoreUsed:        cacheStoreUsed,
		cmsStoreUsed:          cmsStoreUsed,
		cmsStoreTemplateID:    cmsStoreTemplateID,
		customStoreUsed:       customStoreUsed,
		entityStoreUsed:       entityStoreUsed,
		feedStoreUsed:         feedStoreUsed,
		geoStoreUsed:          geoStoreUsed,
		logStoreUsed:          logStoreUsed,
		metaStoreUsed:         metaStoreUsed,
		sessionStoreUsed:      sessionStoreUsed,
		settingStoreUsed:      settingStoreUsed,
		shopStoreUsed:         shopStoreUsed,
		sqlFileStoreUsed:      sqlFileStoreUsed,
		statsStoreUsed:        statsStoreUsed,
		subscriptionStoreUsed: subscriptionStoreUsed,
		taskStoreUsed:         taskStoreUsed,
		userStoreUsed:         userStoreUsed,
		userStoreVaultEnabled: userStoreVaultEnabled,
		vaultStoreUsed:        vaultStoreUsed,
		vaultStoreKey:         vaultStoreKey,
	}
}
