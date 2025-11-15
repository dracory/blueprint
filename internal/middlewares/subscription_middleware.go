package middlewares

import (
	"context"
	"log"
	"net/http"
	"project/internal/config"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/subscriptionstore"

	"github.com/dracory/api"
	"github.com/dracory/rtr"
	"github.com/dracory/userstore"
)

// SubscriptionOnlyMiddleware checks the user has an active subscription
// ==========================================================
// Business Logic:
// 1. Checks an active subscription exists for the specified user
// ==========================================================
func NewSubscriptionOnlyMiddleware(app types.AppInterface) rtr.MiddlewareInterface {
	m := rtr.NewMiddleware().
		SetName("Subscription Only Middleware").
		SetHandler(subscriptionOnlyMiddlewareHandler(app))

	return m
}

func subscriptionOnlyMiddlewareHandler(app types.AppInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUser := helpers.GetAuthUser(r)

			if authUser == nil {
				api.Respond(w, r, api.Unauthenticated("user id empty"))
				return
			}

			authenticatedUser := r.Context().Value(config.AuthenticatedUserContextKey{}).(userstore.UserInterface)

			if authenticatedUser == nil {
				api.Respond(w, r, api.Unauthenticated("user not found"))
				return
			}

			// Check if user is an admin? Yes => Allow
			if authUser.IsAdministrator() || authUser.IsSuperuser() {
				next.ServeHTTP(w, r)
				return
			}

			activeSubscriptions, errSubscriptions := app.GetSubscriptionStore().
				SubscriptionList(context.Background(), subscriptionstore.NewSubscriptionQuery().
					SetSubscriberID(authenticatedUser.ID()).
					SetStatus(subscriptionstore.SUBSCRIPTION_STATUS_ACTIVE))

			if errSubscriptions != nil {
				log.Println(errSubscriptions.Error())
				api.Respond(w, r, api.Error("error listing subscriptions"))
				return
			}

			if len(activeSubscriptions) > 0 {
				next.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, links.User().SubscriptionsPlanSelect(), http.StatusTemporaryRedirect)
		})
	}
}
