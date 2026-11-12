package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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

// Seeded fixture users. Both are always created so tests can cover the
// admin redirect (/admin) and the regular-user redirect (/user).
const (
	e2eAdminEmail = "e2e-admin@blueprint.local"
	e2eAdminID    = "e2e-admin-user"
	e2eUserEmail  = "e2e-user@blueprint.local"
	e2eUserID     = "e2e-regular-user"
)

const e2ePassword = "password123"

// seedFixtures creates the standard admin + regular user (and an optional
// extra user) and returns the emails that must be on the auth allowlist.
func seedFixtures(app app.AppInterface, withUser string, isAdmin bool) ([]string, error) {
	fixtures := []struct {
		email   string
		id      string
		isAdmin bool
	}{
		{e2eAdminEmail, e2eAdminID, true},
		{e2eUserEmail, e2eUserID, false},
	}
	if withUser != "" {
		fixtures = append(fixtures, struct {
			email   string
			id      string
			isAdmin bool
		}{withUser, "e2e-custom-user", isAdmin})
	}

	emails := []string{}
	for _, f := range fixtures {
		if err := seedUserAndSession(app, f.email, f.id, f.isAdmin); err != nil {
			return nil, err
		}
		emails = append(emails, f.email)
	}
	return emails, nil
}

func seedUserAndSession(app app.AppInterface, email string, userID string, isAdmin bool) error {
	ctx := context.Background()

	if app.IsDisabledUserStore() {
		return fmt.Errorf("user store is not initialized")
	}
	if app.IsDisabledSessionStore() {
		return fmt.Errorf("session store is not initialized")
	}

	user, err := findUserByEmail(ctx, app, email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		user, err = createUser(ctx, app, email, userID, isAdmin)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		cfmt.Successln("Created new user:", email)
	} else {
		cfmt.Successln("Found existing user:", email)
	}

	if app.IsEnabledSubscriptionStore() {
		if err := seedSubscription(ctx, app, user.GetID()); err != nil {
			return fmt.Errorf("failed to seed subscription: %w", err)
		}
	}

	session, err := findOrCreateSession(ctx, app, user)
	if err != nil {
		return fmt.Errorf("failed to ensure session: %w", err)
	}

	roleNames := []string{}
	if roles, err := app.GetUserStore().UserRoles(ctx, user.GetID()); err == nil {
		for _, r := range roles {
			roleNames = append(roleNames, r.GetHandle())
		}
	}
	role := strings.Join(roleNames, ",")
	if role == "" {
		role = "user"
	}

	appURL := app.GetConfig().GetAppUrl()

	cfmt.Infoln("==================================================")
	cfmt.Infoln("E2E Browser Seed Credentials")
	cfmt.Infoln("==================================================")
	cfmt.Infoln("Email:      ", email)
	cfmt.Infoln("Password:   ", e2ePassword)
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
		if app.IsDisabledBlindIndexStoreEmail() {
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

func createUser(ctx context.Context, app app.AppInterface, email string, userID string, isAdmin bool) (userstore.UserInterface, error) {
	role := userstore.USER_ROLE_USER
	if isAdmin {
		role = userstore.USER_ROLE_ADMINISTRATOR
	}

	user := userstore.NewUser().
		SetID(userID).
		SetEmail(email).
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetFirstName("E2E").
		SetLastName("Browser").
		SetCountry("US").
		SetTimezone("UTC")

	if err := user.SetPasswordAndHash(e2ePassword); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := app.GetUserStore().UserCreate(ctx, user); err != nil {
		return nil, err
	}

	roleRecord, err := app.GetUserStore().RoleFindByHandleOrCreate(ctx, role, userstore.ROLE_STATUS_ACTIVE)
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	if _, err := app.GetUserStore().UserRoleFindByUserIDAndRoleIDOrCreate(ctx, user.GetID(), roleRecord.GetID()); err != nil {
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	if app.GetConfig().GetUserStoreVaultEnabled() {
		if app.IsDisabledVaultStore() {
			return nil, fmt.Errorf("vault store is not initialized")
		}

		emailToken, err := app.GetVaultStore().TokenCreate(ctx, email, app.GetConfig().GetVaultStoreKey(), 20)
		if err != nil {
			return nil, fmt.Errorf("failed to create email token: %w", err)
		}

		firstNameToken, err := app.GetVaultStore().TokenCreate(ctx, "E2E", app.GetConfig().GetVaultStoreKey(), 20)
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
			SetSearchValue("E2E")
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
	if app.IsDisabledSubscriptionStore() {
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

	planID := "e2e-browser-seed-plan"
	plan, err := app.GetSubscriptionStore().PlanFindByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("failed to find plan: %w", err)
	}

	if plan == nil {
		plan = subscriptionstore.NewPlan().
			SetID(planID).
			SetTitle("E2E Browser Seed Plan").
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
	sessionStore := app.GetSessionStore()
	if sessionStore == nil {
		return nil, fmt.Errorf("session store is nil")
	}

	sessions, err := sessionStore.SessionList(ctx, sessionstore.NewSessionQuery().
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "E2E-Browser")

	session := sessionstore.NewSession().
		SetUserID(user.GetID()).
		SetUserAgent(req.UserAgent()).
		SetIPAddress("127.0.0.1").
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(24).ToDateTimeString(carbon.UTC))

	if err := sessionStore.SessionCreate(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}
