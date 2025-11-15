package links

import "github.com/samber/lo"

type userLinks struct{}

// User is a shortcut for NewUserLinks
func User() *userLinks {
	return &userLinks{}
}

// Home URL
func (l *userLinks) Home(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(USER_HOME, p)
}

// Profile URL
func (l *userLinks) Profile(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(USER_PROFILE, p)
}

// SubscriptionsPlanSelect URL
func (l userLinks) SubscriptionsPlanSelect(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(USER_SUBSCRIPTION_PLAN_SELECT, p)
}

// SubscriptionsPlanSelectAjax URL
func (l userLinks) SubscriptionsPlanSelectAjax(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(USER_SUBSCRIPTION_PLAN_SELECT_AJAX, p)
}

// SubscriptionsPaymentSuccess URL
func (l userLinks) SubscriptionsPaymentSuccess(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(USER_SUBSCRIPTION_PAYMENT_SUCCESS, p)
}

// SubscriptionsPaymentCanceled URL
func (l userLinks) SubscriptionsPaymentCanceled(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(USER_SUBSCRIPTION_PAYMENT_CANCELED, p)
}
