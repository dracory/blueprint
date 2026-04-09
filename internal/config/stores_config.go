package config

import "fmt"

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

	return storesSettings{cmsStoreTemplateID: cmsStoreTemplateID, vaultStoreKey: vaultStoreKey}
}

type storesSettings struct {
	cmsStoreTemplateID string
	vaultStoreKey      string
}
