package app

import (
	"context"
	"database/sql"
	"errors"
	"project/internal/types"

	"github.com/dracory/subscriptionstore"
)

func subscriptionStoreInitialize(app types.AppInterface) error {
	if !app.GetConfig().GetSubscriptionStoreUsed() {
		return nil
	}

	if store, err := newSubscriptionStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetSubscriptionStore(store)
	}

	return nil
}

func subscriptionStoreMigrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetSubscriptionStoreUsed() {
		return nil
	}

	if app.GetSubscriptionStore() == nil {
		return errors.New("subscription store is not initialized")
	}

	if err := app.GetSubscriptionStore().AutoMigrate(context.Background()); err != nil {
		return err
	}

	return nil
}

func newSubscriptionStore(db *sql.DB) (subscriptionstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := subscriptionstore.NewStore(subscriptionstore.NewStoreOptions{
		DB:                    db,
		PlanTableName:         "snv_subscriptions_plan",
		SubscriptionTableName: "snv_subscriptions_subscription",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("sessionstore.NewStore returned a nil store")
	}

	return st, nil
}
