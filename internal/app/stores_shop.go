package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/shopstore"
)

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
