// Package shared contains the post-authentication pipeline shared by all
// login methods (AuthKnight, email OTP, ...).
//
// Once a login method has verified the user's email, it calls SessionLogin
// which performs the common steps: find-or-create the user (with optional
// vault encryption and blind indexing), create a session, set the auth
// cookie, and calculate the post-login redirect URL.
package shared

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"
	authrules "project/internal/rules/auth"

	"github.com/dracory/auth"
	"github.com/dracory/auth/types"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/req"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
	"github.com/dromara/carbon/v2"
)

// Authentication error messages
const (
	MsgAccountNotFound  = `Your account may have been deactivated or deleted. Please contact our support team for assistance.`
	MsgAccountNotActive = `Your account is not active. Please contact our support team for assistance.`
	MsgUserNotFound     = `An unexpected error has occurred trying to find your account. The support team has been notified.`
	MsgSessionError     = `Error creating session`
)

// SessionLogin performs the shared post-authentication pipeline:
//
//  1. Finds or creates the user by email (vault + blind index aware).
//  2. Verifies the account exists and is active.
//  3. Creates a session (2h production, 4h development).
//  4. Sets the auth cookie.
//  5. Calculates the redirect URL (admin panel, registration, or home).
//
// Parameters:
//   - application: the app instance
//   - w: the response writer (used to set the auth cookie)
//   - r: the incoming request (used for session metadata)
//   - email: the verified user email
//   - backUrl: optional URL to redirect to instead of the calculated one
//
// Returns:
//   - redirectURL: where to send the user on success ("" on failure)
//   - needsRegistration: true when the user must complete registration
//   - errorMessage: user-safe error message ("" on success)
func SessionLogin(application app.AppInterface, w http.ResponseWriter, r *http.Request, email, backUrl string) (redirectURL string, needsRegistration bool, errorMessage string) {
	// Enforce the email allowlist before find-or-create so a non-allowed
	// email never gets a user record, session, or registration page. The
	// EmailAllowlistMiddleware only guards /user/* and /admin/* — without
	// this check a blocked email could fully authenticate.
	emailAllowed := authrules.NewEmailAllowedRule(application, email)
	if emailAllowed.Fails() {
		return "", false, emailAllowed.FailMessageFirst()
	}

	user, err := userFindByEmailOrCreate(application, r.Context(), email, userstore.USER_STATUS_ACTIVE)

	if err != nil {
		application.GetLogger().Error("At Shared SessionLogin > User Create Error", slog.String("error", err.Error()))
		return "", false, MsgUserNotFound
	}

	if user == nil {
		return "", false, MsgAccountNotFound
	}

	if active := authrules.NewUserActiveRule(user); active.Fails() {
		return "", false, MsgAccountNotActive
	}

	session := sessionstore.NewSession().
		SetUserID(user.GetID()).
		SetUserAgent(r.UserAgent()).
		SetIPAddress(req.GetIP(r)).
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(2).ToDateTimeString(carbon.UTC))

	if application.GetConfig() != nil && application.GetConfig().IsEnvDevelopment() {
		session.SetExpiresAt(carbon.Now(carbon.UTC).AddHours(4).ToDateTimeString(carbon.UTC))
	}

	sessionStore := application.GetSessionStore()
	if sessionStore == nil {
		application.GetLogger().Error("At Shared SessionLogin > Session Store Error", slog.String("error", "session store is nil"))
		return "", false, MsgSessionError
	}

	err = sessionStore.SessionCreate(r.Context(), session)

	if err != nil {
		application.GetLogger().Error("At Shared SessionLogin > Session Store Error", slog.String("error", err.Error()))
		return "", false, MsgSessionError
	}

	// In development (HTTP), the Secure flag must be disabled so the
	// browser sends the cookie back over plain HTTP.
	cookieOpts := []types.CookieOption{}
	if application.GetConfig() != nil && application.GetConfig().IsEnvDevelopment() {
		cookieOpts = append(cookieOpts, types.WithSecure(false))
	}

	auth.AuthCookieSet(w, r, session.GetKey(), cookieOpts...)

	needsRegistration = !user.IsRegistrationCompleted()

	redirectURL = calculateRedirectURL(application, r.Context(), user)

	// Only honour backUrl when it points back at this application (absolute
	// same-origin URL or a relative path). Anything else is ignored to
	// prevent open redirects — even when the URL arrives via a trusted
	// provider such as AuthKnight.
	if isSafeBackURL(application, backUrl) {
		redirectURL = backUrl
	}

	return redirectURL, needsRegistration, ""
}

// isSafeBackURL reports whether backUrl is safe to redirect to: a relative
// path ("/...", but not "//...") or an absolute URL on this application's
// own domain.
func isSafeBackURL(application app.AppInterface, backUrl string) bool {
	if backUrl == "" {
		return false
	}

	if strings.HasPrefix(backUrl, "/") {
		return !strings.HasPrefix(backUrl, "//")
	}

	return strings.HasPrefix(backUrl, links.Website().Home())
}

// calculateRedirectURL calculates the redirect URL based on the user's role and profile completeness.
//
// 1. By default all users redirect to home
// 2. If user is manager or admin, redirect to admin panel
// 3. If user does not have any names, redirect to profile
func calculateRedirectURL(application app.AppInterface, ctx context.Context, user userstore.UserInterface) string {
	// 1. By default all users redirect to home
	redirectUrl := links.User().Home()

	// 2. If user is manager or admin, redirect to admin panel
	if helpers.UserHasAnyActiveRole(ctx, application, user,
		userstore.USER_ROLE_MANAGER,
		userstore.USER_ROLE_ADMINISTRATOR,
		userstore.USER_ROLE_SUPERUSER) {
		redirectUrl = links.Admin().Home()
	}

	// 3. If user does not have any names, redirect to profile
	if !user.IsRegistrationCompleted() {
		redirectUrl = links.Auth().Register()
		redirectUrl = helpers.ToFlashInfoURL(application.GetCacheStore(), "Thank you for logging in. Please complete your data to finish your registration", redirectUrl, 5)
	}

	return redirectUrl
}

// userCreate creates a new user with privacy-first email encryption and blind indexing.
//
// When the vault store is enabled the plaintext email is replaced with an
// encrypted token and a blind index record is created for future lookups.
func userCreate(application app.AppInterface, ctx context.Context, email string, status string) (userstore.UserInterface, error) {
	user := userstore.NewUser().
		SetStatus(status).
		SetEmail(email)

	if application.IsDisabledUserStore() {
		return nil, errors.New("user store is nil")
	}

	if application.GetConfig().GetUserStoreVaultEnabled() && application.IsDisabledVaultStore() {
		return nil, errors.New(`vault store is nil`)
	}

	err := application.GetUserStore().UserCreate(ctx, user)

	if err != nil {
		return nil, err
	}

	if err := helpers.UserActiveRoleAssign(ctx, application, user.GetID(), userstore.USER_ROLE_USER); err != nil {
		return nil, err
	}

	if !application.GetConfig().GetUserStoreVaultEnabled() {
		return user, nil
	}

	if application.IsDisabledVaultStore() {
		return nil, errors.New(`vault store is nil`)
	}

	emailToken, err := application.GetVaultStore().TokenCreate(ctx, email, application.GetConfig().GetVaultStoreKey(), 20)

	if err != nil {
		return nil, err
	}

	user.SetEmail(emailToken)

	err = application.GetUserStore().UserUpdate(ctx, user)

	if err != nil {
		return nil, err
	}

	searchValue := blindindexstore.NewSearchValue().
		SetSourceReferenceID(user.GetID()).
		SetSearchValue(email)

	err = application.GetBlindIndexStoreEmail().SearchValueCreate(ctx, searchValue)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// userFindByEmailOrCreate finds or creates a user based on the provided email.
//
// Business Logic:
//  1. If VaultStore is used:
//     a. Check if the email is in the blind index, and get the user ID.
//     b. If the user ID is not found, create a new user.
//     c. Find the user by ID.
//  2. If VaultStore is not used:
//     a. Find the user by email.
//     b. If the user is not found, create a new user.
func userFindByEmailOrCreate(application app.AppInterface, ctx context.Context, email string, status string) (userstore.UserInterface, error) {
	if application.IsDisabledUserStore() {
		return nil, errors.New("user store is nil")
	}

	if application.GetConfig().GetUserStoreVaultEnabled() {
		if application.IsDisabledVaultStore() {
			return nil, errors.New(`vault store is nil`)
		}

		userID, err := findUserIDInBlindIndex(application, ctx, email)
		if err != nil {
			return nil, err
		}

		if userID == "" {
			return userCreate(application, ctx, email, status)
		}

		user, err := application.GetUserStore().UserFindByID(ctx, userID)

		if err != nil {
			return nil, err
		}

		if user == nil {
			application.GetLogger().Warn("At Shared SessionLogin > userFindByEmailOrCreate",
				slog.String("error", "User not found, even though email was found in the blind index, and user ID returned successfully"),
				slog.String("user", userID))
			return nil, nil
		}

		return user, nil
	}

	user, err := application.GetUserStore().UserFindByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return userCreate(application, ctx, email, status)
	}

	return user, nil
}

func findUserIDInBlindIndex(application app.AppInterface, ctx context.Context, email string) (userID string, err error) {
	recordsFound, err := application.GetBlindIndexStoreEmail().SearchValueList(ctx, blindindexstore.NewSearchValueQuery().
		SetSearchValue(email).
		SetSearchType(blindindexstore.SEARCH_TYPE_EQUALS))

	if err != nil {
		return "", err
	}

	if len(recordsFound) < 1 {
		return "", nil
	}

	return recordsFound[0].SourceReferenceID(), nil
}
