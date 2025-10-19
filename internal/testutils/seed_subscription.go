package testutils

import (
	"context"

	"github.com/dracory/subscriptionstore"
)

const (
	PLAN_01         = "plan_01"
	SUBSCRIPTION_01 = "subscription_01"
)

func SeedSubscription(store subscriptionstore.StoreInterface, userID string, planID string) (subscriptionstore.SubscriptionInterface, error) {
	subscription := subscriptionstore.NewSubscription().
		SetStatus(subscriptionstore.SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID(userID).
		SetPlanID(planID)

	err := store.SubscriptionCreate(context.Background(), subscription)
	return subscription, err
}
