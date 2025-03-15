package config

import (
	"context"
	"database/sql"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/geostore"
	"github.com/gouniverse/metastore"
	"github.com/gouniverse/sessionstore"
	"github.com/gouniverse/shopstore"
	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/taskstore"
	"github.com/gouniverse/vaultstore"

	"project/app/config/stores"
)

// InitializeStores initializes all stores
func (c *Config) InitializeStores(db *sql.DB) error {
	if c.BlindIndexStoreUsed {
		// Initialize email blind index store
		emailStore, err := stores.BlindIndexStoreInitialize(db, "snv_bindx_email")
		if err != nil {
			return err
		}
		c.BlindIndexStoreEmail = emailStore

		// Initialize first name blind index store
		firstNameStore, err := stores.BlindIndexStoreInitialize(db, "snv_bindx_first_name")
		if err != nil {
			return err
		}
		c.BlindIndexStoreFirstName = firstNameStore

		// Initialize last name blind index store
		lastNameStore, err := stores.BlindIndexStoreInitialize(db, "snv_bindx_last_name")
		if err != nil {
			return err
		}
		c.BlindIndexStoreLastName = lastNameStore

		// Add migrations for all blind index stores
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if err := stores.MigrateBlindIndexStore(ctx, c.BlindIndexStoreEmail); err != nil {
				return err
			}
			if err := stores.MigrateBlindIndexStore(ctx, c.BlindIndexStoreFirstName); err != nil {
				return err
			}
			if err := stores.MigrateBlindIndexStore(ctx, c.BlindIndexStoreLastName); err != nil {
				return err
			}
			return nil
		})
	}

	if c.CacheStoreUsed {
		// Initialize cache store
		cacheStore, err := stores.CacheStoreInitialize(db)
		if err != nil {
			return err
		}
		c.CacheStore = cacheStore

		// Add migration for cache store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if c.CacheStore != nil {
				return stores.CacheStoreAutoMigrate(ctx, c.CacheStore)
			}
			return nil
		})
	}

	if c.CmsStoreUsed {
		// Initialize CMS store
		cmsStore, err := stores.CmsStoreInitialize(db)
		if err != nil {
			return err
		}
		c.CmsStore = cmsStore

		// Add migration for CMS store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if cmsStore, ok := c.CmsStore.(cmsstore.StoreInterface); ok {
				return stores.CmsStoreAutoMigrate(ctx, cmsStore)
			}
			return nil
		})
	}

	if c.CustomStoreUsed {
		// Initialize custom store
		customStore, err := stores.CustomStoreInitialize(db)
		if err != nil {
			return err
		}
		c.CustomStore = customStore

		// Add migration for custom store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if c.CustomStore != nil {
				return stores.CustomStoreAutoMigrate(ctx, c.CustomStore)
			}
			return nil
		})
	}

	if c.EntityStoreUsed {
		// Initialize entity store
		entityStore, err := stores.EntityStoreInitialize(db)
		if err != nil {
			return err
		}
		c.EntityStore = entityStore

		// Add migration for entity store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if c.EntityStore != nil {
				return stores.EntityStoreAutoMigrate(ctx, c.EntityStore)
			}
			return nil
		})
	}

	if c.GeoStoreUsed {
		// Initialize geo store
		geoStore, err := stores.GeoStoreInitialize(db)
		if err != nil {
			return err
		}
		c.GeoStore = geoStore

		// Add migration for geo store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if geoStore, ok := c.GeoStore.(*geostore.Store); ok {
				return stores.GeoStoreAutoMigrate(ctx, geoStore)
			}
			return nil
		})
	}

	if c.LogStoreUsed {
		// Initialize log store
		logStore, err := stores.LogStoreInitialize(db)
		if err != nil {
			return err
		}
		c.LogStore = logStore

		// Add migration for log store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if c.LogStore != nil {
				return stores.LogStoreAutoMigrate(ctx, c.LogStore)
			}
			return nil
		})
	}

	if c.MetaStoreUsed {
		// Initialize meta store
		metaStore, err := stores.MetaStoreInitialize(db)
		if err != nil {
			return err
		}
		c.MetaStore = metaStore

		// Add migration for meta store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if metaStore, ok := c.MetaStore.(metastore.StoreInterface); ok {
				return stores.MetaStoreAutoMigrate(ctx, metaStore)
			}
			return nil
		})
	}

	if c.SessionStoreUsed {
		// Initialize session store
		sessionStore, err := stores.SessionStoreInitialize(db)
		if err != nil {
			return err
		}
		c.SessionStore = sessionStore

		// Add migration for session store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if sessionStore, ok := c.SessionStore.(sessionstore.StoreInterface); ok {
				return stores.SessionStoreAutoMigrate(ctx, sessionStore)
			}
			return nil
		})
	}

	if c.ShopStoreUsed {
		// Initialize shop store
		shopStore, err := stores.ShopStoreInitialize(db)
		if err != nil {
			return err
		}
		c.ShopStore = shopStore

		// Add migration for shop store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if shopStore, ok := c.ShopStore.(*shopstore.Store); ok {
				return stores.ShopStoreAutoMigrate(ctx, shopStore)
			}
			return nil
		})
	}

	if c.StatsStoreUsed {
		// Initialize stats store
		statsStore, err := stores.StatsStoreInitialize(db)
		if err != nil {
			return err
		}
		c.StatsStore = statsStore

		// Add migration for stats store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if statsStore, ok := c.StatsStore.(*statsstore.Store); ok {
				return stores.StatsStoreAutoMigrate(ctx, statsStore)
			}
			return nil
		})
	}

	if c.TaskStoreUsed {
		// Initialize task store
		taskStore, err := stores.TaskStoreInitialize(db)
		if err != nil {
			return err
		}
		c.TaskStore = taskStore

		// Add migration for task store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if taskStore, ok := c.TaskStore.(*taskstore.Store); ok {
				return stores.TaskStoreAutoMigrate(ctx, taskStore)
			}
			return nil
		})
	}

	if c.VaultStoreUsed {
		// Initialize vault store
		vaultStore, err := stores.VaultStoreInitialize(db)
		if err != nil {
			return err
		}
		c.VaultStore = vaultStore

		// Add migration for vault store
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if vaultStore, ok := c.VaultStore.(*vaultstore.Store); ok {
				return stores.VaultStoreAutoMigrate(ctx, vaultStore)
			}
			return nil
		})
	}

	return nil
}
