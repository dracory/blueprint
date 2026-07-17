package app

import (
	"errors"

	"github.com/samber/lo"

	"project/internal/config"
)

// ============================================================================
// == START: Data Stores Initialization
// ============================================================================
//
// This is where all data stores are initialized based on the enabled flags
// configured in stores_config.go. Each enabled store's setup function (below)
// is called to create the store via config builders and wire it into the
// AppInterface.
//
// ============================================================================

// dataStoresInitialize creates and wires all data stores based on configuration.
func (app *appImplementation) dataStoresInitialize() error {
	cfg := app.GetConfig()
	if cfg == nil {
		return errors.New("config is not initialized")
	}

	if app.GetDatabase() == nil {
		return errors.New("database is not initialized")
	}

	blindIndexEnabled := cfg.GetUserStoreUsed() && cfg.GetVaultStoreUsed()

	stores := []struct {
		enabled bool
		init    func(AppInterface) error
	}{
		{cfg.GetAuditStoreUsed(), setupAuditStore},
		{cfg.GetBlogStoreUsed(), setupBlogStore},
		{cfg.GetCacheStoreUsed(), setupCacheStore},
		{cfg.GetChatStoreUsed(), setupChatStore},
		{cfg.GetCmsStoreUsed(), setupCmsStore},
		{cfg.GetCustomStoreUsed(), setupCustomStore},
		{cfg.GetEntityStoreUsed(), setupEntityStore},
		{cfg.GetFeedStoreUsed(), setupFeedStore},
		{cfg.GetGeoStoreUsed(), setupGeoStore},
		{cfg.GetLogStoreUsed(), setupLogStore},
		{cfg.GetMetaStoreUsed(), setupMetaStore},
		{cfg.GetSessionStoreUsed(), setupSessionStore},
		{cfg.GetSettingStoreUsed(), setupSettingStore},
		{cfg.GetShopStoreUsed(), setupShopStore},
		{cfg.GetSqlFileStoreUsed(), setupSqlFileStorage},
		{cfg.GetStatsStoreUsed(), setupStatsStore},
		{cfg.GetSubscriptionStoreUsed(), setupSubscriptionStore},
		{cfg.GetTaskStoreUsed(), setupTaskStore},
		{cfg.GetUserStoreUsed(), setupUserStore},
		{cfg.GetVaultStoreUsed(), setupVaultStore},
		{blindIndexEnabled, setupBlindIndexEmailStore},
		{blindIndexEnabled, setupBlindIndexFirstNameStore},
		{blindIndexEnabled, setupBlindIndexLastNameStore},
	}

	enabledStores := lo.Filter(stores, func(s struct {
		enabled bool
		init    func(AppInterface) error
	}, _ int) bool {
		return s.enabled
	})

	for _, s := range enabledStores {
		if err := s.init(app); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// == START: Store Setup Functions
// ============================================================================
//
// Each function here wires a store created by config.NewXxxStore into the
// AppInterface. These are called by dataStoresInitialize above based on the
// enabled flags in stores_config.go.
//
// ============================================================================

func setupAuditStore(app AppInterface) error {
	st, err := config.NewAuditStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetAuditStore(st)
	return nil
}

func setupBlogStore(app AppInterface) error {
	st, err := config.NewBlogStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetBlogStore(st)
	return nil
}

func setupBlindIndexEmailStore(app AppInterface) error {
	st, err := config.NewBlindIndexEmailStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetBlindIndexStoreEmail(st)
	return nil
}

func setupBlindIndexFirstNameStore(app AppInterface) error {
	st, err := config.NewBlindIndexFirstNameStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetBlindIndexStoreFirstName(st)
	return nil
}

func setupBlindIndexLastNameStore(app AppInterface) error {
	st, err := config.NewBlindIndexLastNameStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetBlindIndexStoreLastName(st)
	return nil
}

func setupCacheStore(app AppInterface) error {
	st, err := config.NewCacheStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetCacheStore(st)
	return nil
}

func setupChatStore(app AppInterface) error {
	st, err := config.NewChatStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetChatStore(st)
	return nil
}

func setupCmsStore(app AppInterface) error {
	st, err := config.NewCmsStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetCmsStore(st)
	return nil
}

func setupCustomStore(app AppInterface) error {
	st, err := config.NewCustomStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetCustomStore(st)
	return nil
}

func setupEntityStore(app AppInterface) error {
	st, err := config.NewEntityStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetEntityStore(st)
	return nil
}

func setupFeedStore(app AppInterface) error {
	st, err := config.NewFeedStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetFeedStore(st)
	return nil
}

func setupGeoStore(app AppInterface) error {
	st, err := config.NewGeoStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetGeoStore(st)
	return nil
}

func setupLogStore(app AppInterface) error {
	st, err := config.NewLogStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetLogStore(st)
	return nil
}

func setupMetaStore(app AppInterface) error {
	st, err := config.NewMetaStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetMetaStore(st)
	return nil
}

func setupSessionStore(app AppInterface) error {
	st, err := config.NewSessionStore(app.GetDatabase(), app.GetConfig().GetAppDebug(), app.GetConfig().IsEnvDevelopment())
	if err != nil {
		return err
	}
	app.SetSessionStore(st)
	return nil
}

func setupSettingStore(app AppInterface) error {
	st, err := config.NewSettingStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetSettingStore(st)
	return nil
}

func setupShopStore(app AppInterface) error {
	st, err := config.NewShopStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetShopStore(st)
	return nil
}

func setupSqlFileStorage(app AppInterface) error {
	st, err := config.NewSqlFileStorage(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetSqlFileStorage(st)
	return nil
}

func setupStatsStore(app AppInterface) error {
	st, err := config.NewStatsStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetStatsStore(st)
	return nil
}

func setupSubscriptionStore(app AppInterface) error {
	st, err := config.NewSubscriptionStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetSubscriptionStore(st)
	return nil
}

func setupTaskStore(app AppInterface) error {
	st, err := config.NewTaskStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetTaskStore(st)
	return nil
}

func setupUserStore(app AppInterface) error {
	st, err := config.NewUserStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetUserStore(st)
	return nil
}

func setupVaultStore(app AppInterface) error {
	st, err := config.NewVaultStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
	if err != nil {
		return err
	}
	app.SetVaultStore(st)
	return nil
}

// ============================================================================
// == END: Data Stores Initialization
// ============================================================================
