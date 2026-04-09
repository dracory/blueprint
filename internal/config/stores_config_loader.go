package config

import "fmt"

// loadStoresConfig loads datastore feature flags directly into the config.
// Each store is opt-in via configuration_stores.go - set the flag to true to enable it.
func loadStoresConfig(cfg ConfigInterface, v *envValidator) {
	// CMS Store Template ID
	//
	// The template ID used by the CMS store for rendering content.
	// Required when CMS_STORE_USED is true.
	cmsStoreTemplateID := v.GetString(KEY_CMS_STORE_TEMPLATE_ID)

	// Vault Store Key
	//
	// Encryption key used by the vault store to encrypt sensitive data.
	// Required when VAULT_STORE_USED is true.
	vaultStoreKey := v.GetString(KEY_VAULT_STORE_KEY)

	if userStoreVaultEnabled && !vaultStoreUsed {
		v.Add(fmt.Errorf("userStoreVaultEnabled requires vaultStoreUsed to be true"))
	}

	v.MustWhen(cmsStoreUsed, KEY_CMS_STORE_TEMPLATE_ID,
		"required when `CMS_STORE_USED` is true", cmsStoreTemplateID)

	cfg.SetAuditStoreUsed(auditStoreUsed)
	cfg.SetBlogStoreUsed(blogStoreUsed)
	cfg.SetCacheStoreUsed(cacheStoreUsed)
	cfg.SetChatStoreUsed(chatStoreUsed)
	cfg.SetCmsStoreUsed(cmsStoreUsed)
	cfg.SetCmsStoreTemplateID(cmsStoreTemplateID)
	cfg.SetCustomStoreUsed(customStoreUsed)
	cfg.SetEntityStoreUsed(entityStoreUsed)
	cfg.SetFeedStoreUsed(feedStoreUsed)
	cfg.SetGeoStoreUsed(geoStoreUsed)
	cfg.SetLogStoreUsed(logStoreUsed)
	cfg.SetMetaStoreUsed(metaStoreUsed)
	cfg.SetSessionStoreUsed(sessionStoreUsed)
	cfg.SetSettingStoreUsed(settingStoreUsed)
	cfg.SetShopStoreUsed(shopStoreUsed)
	cfg.SetSqlFileStoreUsed(sqlFileStoreUsed)
	cfg.SetStatsStoreUsed(statsStoreUsed)
	cfg.SetSubscriptionStoreUsed(subscriptionStoreUsed)
	cfg.SetTaskStoreUsed(taskStoreUsed)
	cfg.SetUserStoreUsed(userStoreUsed)
	cfg.SetUserStoreVaultEnabled(userStoreVaultEnabled)
	cfg.SetVaultStoreUsed(vaultStoreUsed)
	cfg.SetVaultStoreKey(vaultStoreKey)
}
