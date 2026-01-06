package app

import (
	"project/internal/types"
)

// dataStoresMigrate performs phase 2 of store setup. Placeholder for upcoming
// two-phase migration once stores are split into create/migrate.
func (r *Registry) dataStoresMigrate() error {
	migrators := []func(app types.RegistryInterface) error{
		auditStoreMigrate,
		blogStoreMigrate,
		blindIndexEmailStoreMigrate,
		blindIndexFirstNameStoreMigrate,
		blindIndexLastNameStoreMigrate,
		cacheStoreMigrate,
		chatStoreMigrate,
		cmsStoreMigrate,
		customStoreMigrate,
		entityStoreMigrate,
		feedStoreMigrate,
		geoStoreMigrate,
		logStoreMigrate,
		metaStoreMigrate,
		sessionStoreMigrate,
		settingStoreMigrate,
		shopStoreMigrate,
		sqlFileStorageMigrate,
		statsStoreMigrate,
		subscriptionStoreMigrate,
		taskStoreMigrate,
		userStoreMigrate,
		vaultStoreMigrate,
	}

	for _, m := range migrators {
		if err := m(r); err != nil {
			return err
		}
	}

	return nil
}
