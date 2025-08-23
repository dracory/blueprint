package app

import (
	"context"
)

// dataStoresMigrate performs phase 2 of store setup. Placeholder for upcoming
// two-phase migration once stores are split into create/migrate.
func (a *Application) dataStoresMigrate() error {
	ctx := context.Background()

	if blindIndexStoreEmail := a.GetBlindIndexStoreEmail(); blindIndexStoreEmail != nil {
		if err := blindIndexStoreEmail.AutoMigrate(); err != nil {
			return err
		}
	}

	if blindIndexStoreFirstName := a.GetBlindIndexStoreFirstName(); blindIndexStoreFirstName != nil {
		if err := blindIndexStoreFirstName.AutoMigrate(); err != nil {
			return err
		}
	}

	if blindIndexStoreLastName := a.GetBlindIndexStoreLastName(); blindIndexStoreLastName != nil {
		if err := blindIndexStoreLastName.AutoMigrate(); err != nil {
			return err
		}
	}

	if blogStore := a.GetBlogStore(); blogStore != nil {
		if err := blogStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if cacheStore := a.GetCacheStore(); cacheStore != nil {
		if err := cacheStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if customStore := a.GetCustomStore(); customStore != nil {
		if err := customStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if geoStore := a.GetGeoStore(); geoStore != nil {
		if err := geoStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if logStore := a.GetLogStore(); logStore != nil {
		if err := logStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if metaStore := a.GetMetaStore(); metaStore != nil {
		if err := metaStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if sessionStore := a.GetSessionStore(); sessionStore != nil {
		ctx := context.Background()
		if err := sessionStore.AutoMigrate(ctx); err != nil {
			return err
		}
	}

	if settingStore := a.GetSettingStore(); settingStore != nil {
		if err := settingStore.AutoMigrate(ctx); err != nil {
			return err
		}
	}

	if statsStore := a.GetStatsStore(); statsStore != nil {
		if err := statsStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if taskStore := a.GetTaskStore(); taskStore != nil {
		if err := taskStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if tradingStore := a.GetTradingStore(); tradingStore != nil {
		ctx := context.Background()
		if err := tradingStore.AutoMigrateInstruments(ctx); err != nil {
			return err
		}
		if err := tradingStore.AutoMigratePrices(ctx); err != nil {
			return err
		}
	}

	if userStore := a.GetUserStore(); userStore != nil {
		if err := userStore.AutoMigrate(); err != nil {
			return err
		}
	}

	if vaultStore := a.GetVaultStore(); vaultStore != nil {
		if err := vaultStore.AutoMigrate(); err != nil {
			return err
		}
	}

	return nil
}
