package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/shopstore"
)

// ShopStoreInitialize initializes the shop store
func ShopStoreInitialize(db *sql.DB) (*shopstore.Store, error) {
	shopStoreInstance, err := shopstore.NewStore(shopstore.NewStoreOptions{
		DB:                     db,
		CategoryTableName:      "snv_shop_category",
		DiscountTableName:      "snv_shop_discount",
		MediaTableName:         "snv_shop_media",
		OrderTableName:         "snv_shop_order",
		OrderLineItemTableName: "snv_shop_order_line_item",
		ProductTableName:       "snv_shop_product",
	})

	if err != nil {
		return nil, errors.Join(errors.New("shopstore.NewStore"), err)
	}

	if shopStoreInstance == nil {
		return nil, errors.New("ShopStore is nil")
	}

	return shopStoreInstance, nil
}

// ShopStoreAutoMigrate runs migrations for the shop store
func ShopStoreAutoMigrate(ctx context.Context, store *shopstore.Store) error {
	if store == nil {
		return errors.New("shopstore.AutoMigrate: ShopStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("shopstore.AutoMigrate"), err)
	}

	return nil
}
