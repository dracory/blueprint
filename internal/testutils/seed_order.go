package testutils

import (
	"context"
	"errors"

	"github.com/dracory/shopstore"
)

func SeedOrder(shopStore shopstore.StoreInterface, orderID string, customerID string) (shopstore.OrderInterface, error) {
	if shopStore == nil {
		return nil, errors.New("shopstore is nil")
	}

	order, err := shopStore.OrderFindByID(context.Background(), orderID)

	if err != nil {
		return nil, err
	}

	if order != nil {
		return order, nil
	}

	order = shopstore.NewOrder()
	order.SetID(orderID)
	order.SetCustomerID(customerID)

	if err := shopStore.OrderCreate(context.Background(), order); err != nil {
		return nil, err
	}

	return order, nil
}
