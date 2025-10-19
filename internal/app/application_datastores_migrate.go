package app

import (
	"project/internal/types"
)

// dataStoresMigrate performs phase 2 of store setup. Placeholder for upcoming
// two-phase migration once stores are split into create/migrate.
func (a *Application) dataStoresMigrate() error {
	migrators := []func(app types.AppInterface) error{
		auditStoreMigrate,
		blogStoreMigrate,
		blindIndexEmailStoreMigrate,
		blindIndexFirstNameStoreMigrate,
		blindIndexLastNameStoreMigrate,
		cacheStoreMigrate,
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
		if err := m(a); err != nil {
			return err
		}
	}

	return nil
}
