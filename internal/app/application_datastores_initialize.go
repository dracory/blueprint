package app

import (
	"errors"
	"project/internal/types"
)

// dataStoresInitialize performs phase 1 of store setup. For now it delegates
// to initializeStores to preserve behavior; it will be refactored to create-only.
func (a *Application) dataStoresInitialize() error {
	if a.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	initializers := []func(app types.AppInterface) error{
		blindIndexEmailStoreInitialize,
		blindIndexFirstNameStoreInitialize,
		blindIndexLastNameStoreInitialize,
		blogStoreInitialize,
		cacheStoreInitialize,
		cmsStoreInitialize,
		customStoreInitialize,
		entityStoreInitialize,
		feedStoreInitialize,
		geoStoreInitialize,
		logStoreInitialize,
		metaStoreInitialize,
		sessionStoreInitialize,
		settingStoreInitialize,
		shopStoreInitialize,
		sqlFileStorageInitialize,
		userStoreInitialize,
		vaultStoreInitialize,
		taskStoreInitialize,
		statsStoreInitialize,
	}

	for _, initializer := range initializers {
		if err := initializer(a); err != nil {
			return err
		}
	}

	return nil
}
