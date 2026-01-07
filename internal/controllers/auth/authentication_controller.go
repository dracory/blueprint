package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"
	"project/internal/testutils"
	"strings"

	"github.com/dracory/auth"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/req"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

const msgAccountNotFound = `Your account may have been deactivated or deleted. Please contact our support team for assistance.`
const msgAccountNotActive = `Your account is not active. Please contact our support team for assistance.`
const msgUserNotFound = `An unexpected error has occurred trying to find your account. The support team has been notified.`

// == CONTROLLER ==============================================================

// authenticationController handles the authentication of the user,
// once the user has logged in successfully via the AuthKnight service.
type authenticationController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

// NewAuthenticationController creates a new instance with injected app only.
func NewAuthenticationController(application registry.RegistryInterface) *authenticationController {
	return &authenticationController{registry: application}
}

// == PUBLIC METHODS ==========================================================

// Handler handles the authentication.
//
// 1. Checks if there is a once parameter in the request from the AuthKnight service.
// 2. Calls the AuthKnight service with the once parameter.
// 3. Verifies the response from the AuthKnight service.
// 4. Based on the email, it will find or create a user in the database.
// 5. Creates a new session for the user.
// 6. Checks if the user has completed their profile.
// 7. If not, it will redirect the user to the profile page.
// 8. If yes, it will redirect the user to the home page, or the admin panel.
//
// Parameters:
// - w: http.ResponseWriter: the response writer.
// - r: *http.Request: the incoming request.
//
// Return:
// - string: the result of the authentication request.
func (c *authenticationController) Handler(w http.ResponseWriter, r *http.Request) string {
	homeURL := links.Website().Home()
	if c.registry.GetUserStore() == nil {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, `user store is required`, homeURL, 5)
	}

	if c.registry.GetConfig().GetUserStoreVaultEnabled() {
		if c.registry.GetVaultStore() == nil {
			return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, `vault store is required`, homeURL, 5)
		}
	}

	if c.registry.GetConfig().GetUserStoreVaultEnabled() && c.registry.GetBlindIndexStoreEmail() == nil {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, `blind index store is required`, homeURL, 5)
	}

	if c.registry.GetSessionStore() == nil {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, `session store is required`, homeURL, 5)
	}

	email, backUrl, errorMessage := c.emailAndBackUrlFromAuthKnightRequest(r)

	if errorMessage != "" {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, "Authentication Provider Error. "+errorMessage, homeURL, 5)
	}

	user, err := c.userFindByEmailOrCreate(r.Context(), email, userstore.USER_STATUS_ACTIVE)

	if err != nil {
		c.registry.GetLogger().Error("At Auth Controller > AnyIndex > User Create Error", slog.String("error", err.Error()))
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, msgUserNotFound, homeURL, 5)
	}

	if user == nil {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, msgAccountNotFound, homeURL, 5)
	}

	if !user.IsActive() {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, msgAccountNotActive, homeURL, 5)
	}

	session := sessionstore.NewSession().
		SetUserID(user.ID()).
		SetUserAgent(r.UserAgent()).
		SetIPAddress(req.GetIP(r)).
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(2).ToDateTimeString(carbon.UTC))

	if c.registry.GetConfig() != nil && c.registry.GetConfig().IsEnvDevelopment() {
		session.SetExpiresAt(carbon.Now(carbon.UTC).AddHours(4).ToDateTimeString(carbon.UTC))
	}

	err = c.registry.GetSessionStore().SessionCreate(r.Context(), session)

	if err != nil {
		c.registry.GetLogger().Error("At Auth Controller > AnyIndex > Session Store Error", slog.String("error", err.Error()))
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, "Error creating session", homeURL, 5)
	}

	auth.AuthCookieSet(w, r, session.GetKey())

	redirectUrl := c.calculateRedirectURL(user)

	if backUrl != "" {
		redirectUrl = backUrl
	}

	return helpers.ToFlashSuccess(c.registry.GetCacheStore(), w, r, "Login was successful", redirectUrl, 5)
}

// == PRIVATE METHODS =========================================================

func (c *authenticationController) findUserIDInBlindIndex(ctx context.Context, email string) (userID string, err error) {
	recordsFound, err := c.registry.GetBlindIndexStoreEmail().SearchValueList(ctx, blindindexstore.NewSearchValueQuery().
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

func (c *authenticationController) emailAndBackUrlFromAuthKnightRequest(r *http.Request) (email, backUrl, errorMessage string) {
	once := strings.TrimSpace(req.GetStringTrimmed(r, "once"))

	if once == "" {
		return "", "", "Once is required field"
	}

	response, err := c.callAuthKnight(once)

	if err != nil {
		c.registry.GetLogger().Error("At Auth Controller > emailFromAuthKnightRequest > Call Auth Knight Error", slog.String("error", err.Error()))
		return "", "", "No response from authentication provider"
	}

	c.registry.GetLogger().Info("At Auth Controller > emailFromAuthKnightRequest > Call Auth Knight Response", slog.Any("response", response))

	status := lo.ValueOr(response, "status", "")
	message := lo.ValueOr(response, "message", "")
	data := lo.ValueOr(response, "data", "")

	if status == "" {
		return "", "", "No status found"
	}

	if message == "" {
		return "", "", "No message found"
	}

	if data == "" {
		return "", "", "No data found"
	}

	if status != "success" {
		c.registry.GetLogger().Error("At Auth Controller > AnyIndex > Response Status", slog.String("error", message.(string)))
		return "", "", "Invalid authentication response status"
	}

	mapData := data.(map[string]any)

	// Required
	email = strings.TrimSpace(lo.ValueOr(mapData, "email", "").(string))

	// Optional
	backUrl = strings.TrimSpace(lo.ValueOr(mapData, "back_url", "").(string))

	return email, backUrl, ""
}

// callAuthKnight makes a request to the authentication server
// to verify the provided "once" token. The "once" token is provided
// by the AuthKnight service.
//
// Note! If the environment is "testing", it will return a predefined response
// which is used only for testing purposes. In the case of a successful response,
// the email is "test@test.com".
//
// Parameters:
//   - once: The once token to be verified.
//
// Returns:
//   - response: A map containing the response data from the authentication server.
//   - error: An error object if an error occurred during the request.
func (c *authenticationController) callAuthKnight(once string) (map[string]interface{}, error) {
	var response map[string]interface{}

	if c.registry.GetConfig() != nil && c.registry.GetConfig().IsEnvTesting() {
		var testResponseJSONString = ""
		if once == testutils.TestKey(c.registry.GetConfig()) {
			testResponseJSONString = `{"status":"success","message":"success","data":{"email":"test@test.com"}}`
		} else {
			testResponseJSONString = `{"status":"error","message":"once data is invalid:test","data":{}}`
		}
		err := json.NewDecoder(bytes.NewReader([]byte(testResponseJSONString))).Decode(&response)
		if err != nil {
			return nil, fmt.Errorf("failed to decode test response: %v", err)
		}
		return response, nil
	}

	req, err := http.PostForm("https://authknight.com/api/who?once="+once, url.Values{
		"once": {once},
	})

	if err != nil {
		return nil, err
	}

	if req == nil {
		return nil, errors.New("no response")
	}

	defer func() {
		if err := req.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if err := json.NewDecoder(req.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

// calculateRedirectURL calculates the redirect URL based on the user's role and profile completeness.
//
// 1. By default all users redirect to home
// 2. If user is manager or admin, redirect to admin panel
// 3. If user does not have any names, redirect to profile
//
// Parameters:
// - user (models.User): The user object.
//
// Returns:
// - string: The redirect URL.
func (c *authenticationController) calculateRedirectURL(user userstore.UserInterface) string {
	// 1. By default all users redirect to home
	redirectUrl := links.User().Home()

	// 2. If user is manager or admin, redirect to admin panel
	if user.IsManager() || user.IsAdministrator() || user.IsSuperuser() {
		redirectUrl = links.Admin().Home()
	}

	// 3. If user does not have any names, redirect to profile
	if !user.IsRegistrationCompleted() {
		redirectUrl = links.Auth().Register()
		redirectUrl = helpers.ToFlashInfoURL(c.registry.GetCacheStore(), "Thank you for logging in. Please complete your data to finish your registration", redirectUrl, 5)
	}

	return redirectUrl
}

// userCreate creates a new user and returns the user object.
//
// Business Logic:
// 1. Create a new user object.
// 2. If the vault store is not used, return the user object.
// 3. Create a new email token and store it in the vault store.
// 4. Replace the email in user object with the email token.
// 5. Update the user object in the database.
// 6. Insert the email token in the blind index.
// 7. Return the user object.
//
// Parameters:
//   - ctx: The context for the request.
//   - email: The email address of the user.
//   - status: The status of the user.
//
// Returns:
//   - userstore.UserInterface: The user object.
//   - error: An error object if an error occurred during the operation.
func (c *authenticationController) userCreate(ctx context.Context, email string, status string) (userstore.UserInterface, error) {
	user := userstore.NewUser().
		SetStatus(status).
		SetEmail(email)

	if c.registry.GetUserStore() == nil {
		return nil, errors.New("user store is nil")
	}

	if c.registry.GetConfig().GetUserStoreVaultEnabled() && c.registry.GetVaultStore() == nil {
		return nil, errors.New(`vault store is nil`)
	}

	err := c.registry.GetUserStore().UserCreate(ctx, user)

	if err != nil {
		return nil, err
	}

	if !c.registry.GetConfig().GetUserStoreVaultEnabled() {
		return user, nil
	}

	if c.registry.GetVaultStore() == nil {
		return nil, errors.New(`vault store is nil`)
	}

	emailToken, err := c.registry.GetVaultStore().TokenCreate(ctx, email, c.registry.GetConfig().GetVaultStoreKey(), 20)

	if err != nil {
		return nil, err
	}

	user.SetEmail(emailToken)

	err = c.registry.GetUserStore().UserUpdate(context.Background(), user)

	if err != nil {
		return nil, err
	}

	searchValue := blindindexstore.NewSearchValue().
		SetSourceReferenceID(user.ID()).
		SetSearchValue(email)

	err = c.registry.GetBlindIndexStoreEmail().SearchValueCreate(ctx, searchValue)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// userFindByEmailOrCreate finds or creates a user based on the provided email.
//
// Business Logic:
//  1. If VultStore is used:
//     a. Check if the email is in the blind index, and get the user ID.
//     b. If the user ID is not found, create a new user.
//     c. Find the user by ID.
//  2. If VultStore is not used:
//     a. Find the user by email.
//     b. If the user is not found, create a new user.
//
// Parameters:
//   - ctx: The context for the request.
//   - email: The email address of the user.
//   - status: The status of the user.
//
// Returns:
//   - userstore.UserInterface: The user object.
//   - error: An error object if an error occurred during the operation.
func (c *authenticationController) userFindByEmailOrCreate(ctx context.Context, email string, status string) (userstore.UserInterface, error) {
	if c.registry.GetUserStore() == nil {
		return nil, errors.New("user store is nil")
	}

	if c.registry.GetConfig().GetUserStoreVaultEnabled() {
		if c.registry.GetVaultStore() == nil {
			return nil, errors.New(`vault store is nil`)
		}

		userID, err := c.findUserIDInBlindIndex(ctx, email)
		if err != nil {
			return nil, err
		}

		if userID == "" {
			return c.userCreate(ctx, email, status)
		}

		user, err := c.registry.GetUserStore().UserFindByID(context.Background(), userID)

		if err != nil {
			return nil, err
		}

		if user == nil {
			c.registry.GetLogger().Error("At Auth Controller > userFindByEmailOrCreate",
				slog.String("error", "User not found, even though email was found in the blind index, and user ID returned successfully"),
				slog.String("user", userID))
			return nil, nil
		}

		return user, nil
	}

	user, err := c.registry.GetUserStore().UserFindByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return c.userCreate(ctx, email, status)
	}

	return user, nil
}
