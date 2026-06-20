package main

import (
	"context"
	"fmt"
	"net/http"

	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/auth"
	"github.com/dracory/base/cfmt"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/sessionstore"
	"github.com/dracory/subscriptionstore"
	"github.com/dracory/userstore"
	"github.com/dromara/carbon/v2"
)

const aiBrowserDefaultEmail = "ai-browser@blueprint.local"

func seedUserAndSession(app app.AppInterface, email string, isAdmin bool) error {
	ctx := context.Background()

	if app.GetUserStore() == nil {
		return fmt.Errorf("user store is not initialized")
	}
	if app.GetSessionStore() == nil {
		return fmt.Errorf("session store is not initialized")
	}

	user, err := findUserByEmail(ctx, app, email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		user, err = createUser(ctx, app, email, isAdmin)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		cfmt.Successln("Created new user:", email)
	} else {
		cfmt.Successln("Found existing user:", email)
	}

	if app.GetSubscriptionStore() != nil {
		if err := seedSubscription(ctx, app, user.GetID()); err != nil {
			return fmt.Errorf("failed to seed subscription: %w", err)
		}
	}

	session, err := findOrCreateSession(ctx, app, user)
	if err != nil {
		return fmt.Errorf("failed to ensure session: %w", err)
	}

	role := user.GetRole()
	if role == "" {
		role = "user"
	}

	appURL := app.GetConfig().GetAppUrl()

	cfmt.Infoln("==================================================")
	cfmt.Infoln("AI Browser Auto-Seed Credentials")
	cfmt.Infoln("==================================================")
	cfmt.Infoln("Email:      ", email)
	cfmt.Infoln("Password:   ", "password123")
	cfmt.Infoln("Role:       ", role)
	cfmt.Infoln("User ID:    ", user.GetID())
	cfmt.Infoln("Session Key:", session.GetKey())
	cfmt.Infoln("Cookie:     ", auth.CookieName+"="+session.GetKey())
	cfmt.Infoln("Login URL:  ", appURL+links.AUTH_LOGIN)
	cfmt.Infoln("Dashboard:  ", appURL+links.USER_HOME)
	cfmt.Infoln("==================================================")

	return nil
}

func findUserByEmail(ctx context.Context, app app.AppInterface, email string) (userstore.UserInterface, error) {
	if app.GetConfig().GetUserStoreVaultEnabled() {
		if app.GetBlindIndexStoreEmail() == nil {
			return nil, fmt.Errorf("blind index store email is not initialized")
		}

		recordsFound, err := app.GetBlindIndexStoreEmail().SearchValueList(ctx, blindindexstore.NewSearchValueQuery().
			SetSearchValue(email).
			SetSearchType(blindindexstore.SEARCH_TYPE_EQUALS))
		if err != nil {
			return nil, err
		}
		if len(recordsFound) == 0 {
			return nil, nil
		}

		userID := recordsFound[0].SourceReferenceID()
		return app.GetUserStore().UserFindByID(ctx, userID)
	}

	return app.GetUserStore().UserFindByEmail(ctx, email)
}

func createUser(ctx context.Context, app app.AppInterface, email string, isAdmin bool) (userstore.UserInterface, error) {
	role := userstore.USER_ROLE_USER
	if isAdmin {
		role = userstore.USER_ROLE_ADMINISTRATOR
	}

	user := userstore.NewUser().
		SetEmail(email).
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetRole(role).
		SetFirstName("AI").
		SetLastName("Browser").
		SetCountry("US").
		SetTimezone("UTC")

	if err := user.SetPasswordAndHash("password123"); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := app.GetUserStore().UserCreate(ctx, user); err != nil {
		return nil, err
	}

	if app.GetConfig().GetUserStoreVaultEnabled() {
		if app.GetVaultStore() == nil {
			return nil, fmt.Errorf("vault store is not initialized")
		}

		emailToken, err := app.GetVaultStore().TokenCreate(ctx, email, app.GetConfig().GetVaultStoreKey(), 20)
		if err != nil {
			return nil, fmt.Errorf("failed to create email token: %w", err)
		}

		firstNameToken, err := app.GetVaultStore().TokenCreate(ctx, "AI", app.GetConfig().GetVaultStoreKey(), 20)
		if err != nil {
			return nil, fmt.Errorf("failed to create first name token: %w", err)
		}

		lastNameToken, err := app.GetVaultStore().TokenCreate(ctx, "Browser", app.GetConfig().GetVaultStoreKey(), 20)
		if err != nil {
			return nil, fmt.Errorf("failed to create last name token: %w", err)
		}

		user.SetEmail(emailToken)
		user.SetFirstName(firstNameToken)
		user.SetLastName(lastNameToken)
		if err := app.GetUserStore().UserUpdate(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update user with tokens: %w", err)
		}

		emailSearchValue := blindindexstore.NewSearchValue().
			SetSourceReferenceID(user.GetID()).
			SetSearchValue(email)
		if err := app.GetBlindIndexStoreEmail().SearchValueCreate(ctx, emailSearchValue); err != nil {
			return nil, fmt.Errorf("failed to create email blind index: %w", err)
		}

		firstNameSearchValue := blindindexstore.NewSearchValue().
			SetSourceReferenceID(user.GetID()).
			SetSearchValue("AI")
		if err := app.GetBlindIndexStoreFirstName().SearchValueCreate(ctx, firstNameSearchValue); err != nil {
			return nil, fmt.Errorf("failed to create first name blind index: %w", err)
		}

		lastNameSearchValue := blindindexstore.NewSearchValue().
			SetSourceReferenceID(user.GetID()).
			SetSearchValue("Browser")
		if err := app.GetBlindIndexStoreLastName().SearchValueCreate(ctx, lastNameSearchValue); err != nil {
			return nil, fmt.Errorf("failed to create last name blind index: %w", err)
		}
	}

	return user, nil
}

func seedSubscription(ctx context.Context, app app.AppInterface, userID string) error {
	if app.GetSubscriptionStore() == nil {
		return nil
	}

	activeSubscriptions, err := app.GetSubscriptionStore().SubscriptionList(ctx,
		subscriptionstore.NewSubscriptionQuery().
			SetSubscriberID(userID).
			SetStatus(subscriptionstore.SUBSCRIPTION_STATUS_ACTIVE))
	if err != nil {
		return fmt.Errorf("failed to list subscriptions: %w", err)
	}

	if len(activeSubscriptions) > 0 {
		return nil
	}

	planID := "ai-browser-seed-plan"
	plan, err := app.GetSubscriptionStore().PlanFindByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("failed to find plan: %w", err)
	}

	if plan == nil {
		plan = subscriptionstore.NewPlan().
			SetID(planID).
			SetTitle("AI Browser Seed Plan").
			SetPrice("0.00").
			SetType(subscriptionstore.PLAN_TYPE_TRIAL).
			SetInterval(subscriptionstore.PLAN_INTERVAL_MONTHLY).
			SetCurrency("USD").
			SetStatus(subscriptionstore.PLAN_STATUS_ACTIVE)

		if err := app.GetSubscriptionStore().PlanCreate(ctx, plan); err != nil {
			return fmt.Errorf("failed to create plan: %w", err)
		}
		cfmt.Successln("Created seed plan:", planID)
	}

	subscription := subscriptionstore.NewSubscription().
		SetStatus(subscriptionstore.SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID(userID).
		SetPlanID(planID)

	if err := app.GetSubscriptionStore().SubscriptionCreate(ctx, subscription); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	cfmt.Successln("Created active subscription for user:", userID)
	return nil
}

func findOrCreateSession(ctx context.Context, app app.AppInterface, user userstore.UserInterface) (sessionstore.SessionInterface, error) {
	sessions, err := app.GetSessionStore().SessionList(ctx, sessionstore.NewSessionQuery().
		SetUserID(user.GetID()).
		SetExpiresAtGte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)))
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	for _, s := range sessions {
		if !s.IsExpired() {
			return s, nil
		}
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	req.Header.Set("User-Agent", "AI-Browser")

	session := sessionstore.NewSession().
		SetUserID(user.GetID()).
		SetUserAgent(req.UserAgent()).
		SetIPAddress("127.0.0.1").
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(24).ToDateTimeString(carbon.UTC))

	if err := app.GetSessionStore().SessionCreate(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}
