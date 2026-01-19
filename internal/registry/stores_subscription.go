package registry

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/subscriptionstore"
)

// subscriptionStoreInitialize initializes the subscription store if enabled in the configuration.
func subscriptionStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetSubscriptionStoreUsed() {
		return nil
	}

	if store, err := newSubscriptionStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetSubscriptionStore(store)
	}

	return nil
}

func subscriptionStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetSubscriptionStoreUsed() {
		return nil
	}

	subscriptionStore := registry.GetSubscriptionStore()
	if subscriptionStore == nil {
		return errors.New("subscription store is not initialized")
	}

	if err := subscriptionStore.AutoMigrate(context.Background()); err != nil {
		return err
	}

	return nil
}

// newSubscriptionStore constructs the Subscription store without running migrations
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
		return nil, errors.New("subscriptionstore.NewStore returned a nil store")
	}
	
	return st, nil
}
