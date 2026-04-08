package config

import "fmt"

// loadStoresConfig loads stores configuration directly into the config.
func loadStoresConfig(cfg ConfigInterface, v *envValidator) {
	cmsStoreTemplateID := v.GetString(KEY_CMS_STORE_TEMPLATE_ID)
	vaultStoreKey := v.GetString(KEY_VAULT_STORE_KEY)

	if userStoreVaultEnabled && !vaultStoreUsed {
		v.Add(fmt.Errorf("%v requires %v to be true",
			userStoreVaultEnabled, vaultStoreUsed))
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
