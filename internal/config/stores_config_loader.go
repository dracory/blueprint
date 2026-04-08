package config

import (
	"fmt"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// storesConfig captures feature store toggles.
type storesConfig struct {
	auditStoreUsed        bool
	blogStoreUsed         bool
	cacheStoreUsed        bool
	chatStoreUsed         bool
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

// loadStoresConfig loads stores configuration from environment variables.
func loadStoresConfig(acc *baseCfg.LoadAccumulator) storesConfig {
	cmsStoreTemplateID := env.GetString(KEY_CMS_STORE_TEMPLATE_ID)
	vaultStoreKey := env.GetString(KEY_VAULT_STORE_KEY)

	if userStoreVaultEnabled && !vaultStoreUsed {
		acc.Add(fmt.Errorf("%v requires %v to be true",
			userStoreVaultEnabled, vaultStoreUsed))
	}

	acc.MustWhen(cmsStoreUsed, KEY_CMS_STORE_TEMPLATE_ID,
		"required when `CMS_STORE_USED` is true", cmsStoreTemplateID)

	return storesConfig{
		auditStoreUsed:        auditStoreUsed,
		blogStoreUsed:         blogStoreUsed,
		cacheStoreUsed:        cacheStoreUsed,
		chatStoreUsed:         chatStoreUsed,
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
