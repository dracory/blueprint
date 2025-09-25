package app

import (
	"project/internal/types"
)

// dataStoresMigrate performs phase 2 of store setup. Placeholder for upcoming
// two-phase migration once stores are split into create/migrate.
func (a *Application) dataStoresMigrate() error {
	migartors := []func(app types.AppInterface) error{
		blogStoreMgrate,
		blindIndexEmailStoreMigrate,
		blindIndexFirstNameStoreMigrate,
		blindIndexLastNameStoreMigrate,
		cacheStoreMgrate,
		cmsStoreMgrate,
		customStoreMgrate,
		entityStoreMgrate,
		feedStoreMgrate,
		geoStoreMgrate,
		logStoreMgrate,
		metaStoreMgrate,
		sessionStoreMgrate,
		settingStoreMgrate,
		shopStoreMgrate,
		statsStoreMgrate,
		taskStoreMigrate,
		tradingStoreMigrate,
		userStoreMigrate,
		vaultStoreMigrate,
	}

	for _, m := range migartors {
		if err := m(a); err != nil {
			return err
		}
	}

	return nil
}
