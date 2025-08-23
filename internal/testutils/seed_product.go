package testutils

import (
	"context"
	"errors"

	"github.com/dracory/shopstore"
)

func SeedProduct(shopStore shopstore.StoreInterface, productID string, price float64) (shopstore.ProductInterface, error) {
	if shopStore == nil {
		return nil, errors.New("shopstore is nil")
	}

	product, err := shopStore.ProductFindByID(context.Background(), productID)

	if err != nil {
		return nil, err
	}

	if product != nil {
		return product, nil
	}

	product = shopstore.NewProduct()
	product.SetID(productID)
	product.SetTitle("Test Product")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)
	product.SetPriceFloat(price)
	product.SetQuantityInt(10)

	if err := shopStore.ProductCreate(context.Background(), product); err != nil {
		return nil, err
	}

	return product, nil
}
