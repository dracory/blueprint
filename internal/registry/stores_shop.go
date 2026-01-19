package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/shopstore"
)

// shopStoreInitialize initializes the shop store if enabled in the configuration.
func shopStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetShopStoreUsed() {
		return nil
	}

	if store, err := newShopStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetShopStore(store)
	}

	return nil
}

func shopStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetShopStoreUsed() {
		return nil
	}

	shopStore := registry.GetShopStore()
	if shopStore == nil {
		return errors.New("shop store is not initialized")
	}

	err := shopStore.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

// newShopStore constructs the Shop store without running migrations
func newShopStore(db *sql.DB) (shopstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := shopstore.NewStore(shopstore.NewStoreOptions{
		DB:                     db,
		CategoryTableName:      "snv_shop_category",
		DiscountTableName:      "snv_shop_discount",
		MediaTableName:         "snv_shop_media",
		OrderTableName:         "snv_shop_order",
		OrderLineItemTableName: "snv_shop_order_line_item",
		ProductTableName:       "snv_shop_product",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("shopstore.NewStore returned a nil store")
	}

	return st, nil
}
