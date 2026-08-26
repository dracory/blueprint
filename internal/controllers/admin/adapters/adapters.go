// Package adapters bridges the external blogadmin, shopadmin, and
// useradmin packages to the blueprint's internal services (admin
// layout, LLM config, userstore for customer resolution, vault
// tokenization, flash redirects).
//
// The external packages are decoupled from any specific host project, so
// they accept callback functions and interfaces (FuncLayout, LlmFactory,
// CustomerResolverInterface, VaultTokenizer, FlashRedirectFunc). This
// package provides concrete adapters wired to the blueprint's
// app.AppInterface.
package adapters

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/app"
	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/layouts"

	"github.com/dracory/auth"
	"github.com/dracory/auth/types"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/dracory/llm"
	"github.com/dracory/neat"
	"github.com/dracory/req"
	"github.com/dracory/sessionstore"
	"github.com/dracory/taskstore"
	"github.com/dracory/useradmin"
	"github.com/dracory/userstore"
	"github.com/dromara/carbon/v2"
)

// LayoutOptions is the anonymous struct used by the external blogadmin
// and shopadmin packages for FuncLayout options. Both packages use the
// same struct shape so a single layout adapter works for both.
//
// This is defined as a type alias to the anonymous struct so it matches
// the external packages' FuncLayout signature exactly.
type LayoutOptions = struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}

// NewLayoutFunc returns a FuncLayout adapter that bridges the external
// blogadmin/shopadmin FuncLayout signature to the blueprint's
// layouts.NewAdminLayout. The same function can be passed to both
// blogadmin.AdminOptions.FuncLayout and shopadmin.AdminOptions.FuncLayout.
func NewLayoutFunc(app app.AppInterface) func(w http.ResponseWriter, r *http.Request, title string, body string, options LayoutOptions) string {
	return func(w http.ResponseWriter, r *http.Request, title string, body string, options LayoutOptions) string {
		return layouts.NewAdminLayout(app, r, layouts.Options{
			Title:      title,
			Content:    hb.Raw(body),
			ScriptURLs: options.ScriptURLs,
			Scripts:    options.Scripts,
			StyleURLs:  options.StyleURLs,
			Styles:     options.Styles,
		}).ToHTML()
	}
}

// NewLlmFactory returns an LLM factory function for the external blogadmin
// package's AI controllers. It uses the blueprint's config to select the
// provider (mock in testing, OpenRouter otherwise) and API key.
func NewLlmFactory(app app.AppInterface) func() (llm.LlmInterface, error) {
	return func() (llm.LlmInterface, error) {
		provider := llm.ProviderOpenRouter
		if app.GetConfig().IsEnvTesting() {
			provider = llm.ProviderMock
		}
		return llm.JSONModel(provider, llm.LlmOptions{
			ApiKey: app.GetConfig().GetOpenRouterApiKey(),
			Model:  llm.OPENROUTER_MODEL_GEMINI_2_5_FLASH_LITE,
		})
	}
}

// UserStoreCustomerResolver implements shopadmin.CustomerResolverInterface
// using the blueprint's userstore. It resolves customer display name and
// email for order views, and supports searching by name/email substrings.
type UserStoreCustomerResolver struct {
	userStore userstore.StoreInterface
}

// NewUserStoreCustomerResolver creates a CustomerResolver backed by userstore.
func NewUserStoreCustomerResolver(userStore userstore.StoreInterface) *UserStoreCustomerResolver {
	return &UserStoreCustomerResolver{userStore: userStore}
}

// FindByID returns the customer display name and email for a given user ID.
// Returns empty strings if not found.
func (r *UserStoreCustomerResolver) FindByID(ctx context.Context, customerID string) (name, email string) {
	if r.userStore == nil || customerID == "" {
		return "", ""
	}
	user, err := r.userStore.UserFindByID(ctx, customerID)
	if err != nil || user == nil {
		return "", ""
	}
	name = strings.TrimSpace(user.GetFirstName() + " " + user.GetLastName())
	email = user.GetEmail()
	return name, email
}

// SearchIDs returns user IDs matching the given name and/or email substrings.
// Empty string means "no filter on that field".
func (r *UserStoreCustomerResolver) SearchIDs(ctx context.Context, name, email string) ([]string, error) {
	if r.userStore == nil {
		return nil, nil
	}
	query := userstore.NewUserQuery()
	if email != "" {
		query.SetEmailLike(email)
	}
	if name != "" {
		query.SetFirstNameLike(name)
	}
	users, err := r.userStore.UserList(ctx, query)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(users))
	for i, u := range users {
		ids[i] = u.GetID()
	}
	return ids, nil
}

// VaultTokenizerAdapter implements useradmin.VaultTokenizer
// using the blueprint's ext.UserTokenize / ext.UserUntokenize helpers.
// It is only active when the blueprint's config has vault store enabled;
// otherwise it returns plain-text passthrough so useradmin treats user
// fields as plain text.
type VaultTokenizerAdapter struct {
	app app.AppInterface
}

// NewVaultTokenizerAdapter creates a VaultTokenizer backed by the
// blueprint's vault store. The adapter checks config.GetVaultStoreUsed()
// on every call so toggling vault at runtime is respected.
func NewVaultTokenizerAdapter(app app.AppInterface) *VaultTokenizerAdapter {
	return &VaultTokenizerAdapter{app: app}
}

// Tokenize upserts vault tokens for the given user fields and returns
// the resulting token strings. When vault is disabled, it returns the
// input values unchanged.
func (a *VaultTokenizerAdapter) Tokenize(
	ctx context.Context,
	user userstore.UserInterface,
	firstName, lastName, email, phone, businessName string,
) (string, string, string, string, string, error) {
	if a.app == nil || a.app.GetConfig() == nil || !a.app.GetConfig().GetVaultStoreUsed() || a.app.GetVaultStore() == nil {
		return firstName, lastName, email, phone, businessName, nil
	}
	return ext.UserTokenize(
		ctx,
		a.app.GetVaultStore(),
		a.app.GetConfig().GetVaultStoreKey(),
		user,
		firstName, lastName, email, phone, businessName,
	)
}

// Untokenize resolves the tokenized fields on the given user back to
// their plain-text values. When vault is disabled, it returns the
// user's stored field values unchanged.
func (a *VaultTokenizerAdapter) Untokenize(
	ctx context.Context,
	user userstore.UserInterface,
) (string, string, string, string, string, error) {
	if a.app == nil || a.app.GetConfig() == nil || !a.app.GetConfig().GetVaultStoreUsed() || a.app.GetVaultStore() == nil {
		if user == nil {
			return "", "", "", "", "", nil
		}
		return user.GetFirstName(), user.GetLastName(), user.GetEmail(), user.GetPhone(), user.GetBusinessName(), nil
	}
	return ext.UserUntokenize(ctx, a.app, a.app.GetConfig().GetVaultStoreKey(), user)
}

// NewFlashRedirectFunc returns a FlashRedirectFunc adapter that bridges
// the useradmin.FlashRedirectFunc signature to the blueprint's
// helpers.ToFlash* helpers. It uses the blueprint's cache store and
// links.Website().Flash() route.
func NewFlashRedirectFunc(app app.AppInterface) useradmin.FlashRedirectFunc {
	return func(w http.ResponseWriter, r *http.Request, messageType, message, redirectURL string, seconds int) string {
		switch messageType {
		case "error":
			return helpers.ToFlashError(app.GetCacheStore(), w, r, message, redirectURL, seconds)
		case "success":
			return helpers.ToFlashSuccess(app.GetCacheStore(), w, r, message, redirectURL, seconds)
		case "info":
			return helpers.ToFlashInfo(app.GetCacheStore(), w, r, message, redirectURL, seconds)
		case "warning":
			return helpers.ToFlashWarning(app.GetCacheStore(), w, r, message, redirectURL, seconds)
		default:
			return helpers.ToFlashError(app.GetCacheStore(), w, r, message, redirectURL, seconds)
		}
	}
}

// GeoResolver implements useradmin.GeoResolverInterface using the
// blueprint's geostore. It adapts the geostore's query-option-based API
// to the simpler GeoResolverInterface shape (no query structs, no column
// constants leaked into useradmin).
type GeoResolver struct {
	geoStore geostore.StoreInterface
}

// Compile-time assertion that GeoResolver satisfies useradmin.GeoResolverInterface.
var _ useradmin.GeoResolverInterface = (*GeoResolver)(nil)

// NewGeoResolver creates a GeoResolver backed by the blueprint's geostore.
func NewGeoResolver(geoStore geostore.StoreInterface) *GeoResolver {
	return &GeoResolver{geoStore: geoStore}
}

// Countries returns all countries sorted by name ascending.
func (r *GeoResolver) Countries(ctx context.Context) ([]useradmin.Country, error) {
	if r.geoStore == nil {
		return nil, nil
	}
	list, err := r.geoStore.CountryList(ctx, geostore.CountryQueryOptions{
		SortOrder: neat.SortAsc,
		OrderBy:   geostore.COLUMN_NAME,
	})
	if err != nil {
		return nil, err
	}
	out := make([]useradmin.Country, 0, len(list))
	for _, c := range list {
		out = append(out, useradmin.Country{
			IsoCode2: c.IsoCode2(),
			Name:     c.Name(),
		})
	}
	return out, nil
}

// Timezones returns timezones for the given country code, sorted by
// timezone ascending. Returns an empty list when no country code is
// provided (no country selected means no timezones to show).
func (r *GeoResolver) Timezones(ctx context.Context, countryCode ...string) ([]useradmin.Timezone, error) {
	if r.geoStore == nil || len(countryCode) == 0 || countryCode[0] == "" {
		return nil, nil
	}
	list, err := r.geoStore.TimezoneList(ctx, geostore.TimezoneQueryOptions{
		SortOrder:   neat.SortAsc,
		OrderBy:     geostore.COLUMN_TIMEZONE,
		CountryCode: countryCode[0],
	})
	if err != nil {
		return nil, err
	}
	out := make([]useradmin.Timezone, 0, len(list))
	for _, tz := range list {
		out = append(out, useradmin.Timezone{
			Code: tz.Timezone(),
		})
	}
	return out, nil
}

// NewOnUserSearchFunc returns an OnUserSearch callback that searches
// the blueprint's blind index stores. It intersects results across
// active filters (AND semantics). When ExactMatch is true, uses
// SEARCH_TYPE_EQUALS; otherwise SEARCH_TYPE_CONTAINS.
func NewOnUserSearchFunc(app app.AppInterface) useradmin.OnUserSearchFunc {
	return func(ctx context.Context, event useradmin.UserSearchEvent) ([]string, error) {
		if app == nil {
			return nil, nil
		}

		searchType := blindindexstore.SEARCH_TYPE_CONTAINS
		if event.ExactMatch {
			searchType = blindindexstore.SEARCH_TYPE_EQUALS
		}

		var idSets [][]string

		if event.FirstName != "" {
			store := app.GetBlindIndexStoreFirstName()
			if store != nil {
				ids, err := store.Search(ctx, event.FirstName, searchType)
				if err != nil {
					return nil, err
				}
				if len(ids) == 0 {
					return nil, nil
				}
				idSets = append(idSets, ids)
			}
		}

		if event.LastName != "" {
			store := app.GetBlindIndexStoreLastName()
			if store != nil {
				ids, err := store.Search(ctx, event.LastName, searchType)
				if err != nil {
					return nil, err
				}
				if len(ids) == 0 {
					return nil, nil
				}
				idSets = append(idSets, ids)
			}
		}

		if event.Email != "" {
			store := app.GetBlindIndexStoreEmail()
			if store != nil {
				ids, err := store.Search(ctx, event.Email, searchType)
				if err != nil {
					return nil, err
				}
				if len(ids) == 0 {
					return nil, nil
				}
				idSets = append(idSets, ids)
			}
		}

		if len(idSets) == 0 {
			return nil, nil
		}
		return intersectIDSets(idSets), nil
	}
}

// intersectIDSets computes the intersection of multiple ID slices.
// Returns IDs that appear in ALL input slices. If any slice is empty,
// the result is empty. Order follows the first slice.
func intersectIDSets(sets [][]string) []string {
	if len(sets) == 0 {
		return nil
	}
	if len(sets) == 1 {
		return sets[0]
	}

	counts := make(map[string]int, len(sets[0]))
	order := make([]string, 0, len(sets[0]))
	for _, id := range sets[0] {
		if counts[id] == 0 {
			order = append(order, id)
		}
		counts[id]++
	}
	for _, s := range sets[1:] {
		seen := make(map[string]bool, len(s))
		for _, id := range s {
			if !seen[id] {
				seen[id] = true
				counts[id]++
			}
		}
	}

	result := make([]string, 0, len(order))
	for _, id := range order {
		if counts[id] == len(sets) {
			result = append(result, id)
		}
	}
	return result
}

// SessionResolver implements useradmin.SessionResolverInterface using
// the blueprint's sessionstore and auth cookie helpers. It absorbs the
// session creation, expiry policy, and cookie setting logic that
// previously lived in useradmin/user_impersonate/impersonate.go.
type SessionResolver struct {
	sessionStore sessionstore.StoreInterface
}

// NewSessionResolver creates a SessionResolver backed by the
// blueprint's sessionstore.
func NewSessionResolver(sessionStore sessionstore.StoreInterface) *SessionResolver {
	return &SessionResolver{sessionStore: sessionStore}
}

// Compile-time assertion that SessionResolver satisfies useradmin.SessionResolverInterface.
var _ useradmin.SessionResolverInterface = (*SessionResolver)(nil)

// Create creates a new session for the given user ID and sets the auth
// cookie on the response. The session expires after 2 hours.
func (r *SessionResolver) Create(w http.ResponseWriter, httpReq *http.Request, userID string, secure bool) error {
	if r.sessionStore == nil {
		return errors.New("session store is nil")
	}

	session := sessionstore.NewSession().
		SetUserID(userID).
		SetUserAgent(httpReq.UserAgent()).
		SetIPAddress(req.GetIP(httpReq)).
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(2).ToDateTimeString(carbon.UTC))

	if err := r.sessionStore.SessionCreate(httpReq.Context(), session); err != nil {
		return err
	}

	auth.AuthCookieSet(w, httpReq, session.GetKey(), types.WithSecure(secure))
	return nil
}

// NewOnUserUpdateFunc returns an OnUserUpdate callback that enqueues a
// blind index rebuild task when a user's email changes. It bridges
// useradmin's event-based hook to the blueprint's taskstore.
func NewOnUserUpdateFunc(app app.AppInterface, taskAlias string) useradmin.OnUserUpdateFunc {
	return func(ctx context.Context, event useradmin.UserUpdateEvent) {
		if app == nil || app.GetTaskStore() == nil || taskAlias == "" {
			return
		}
		_, err := app.GetTaskStore().TaskDefinitionEnqueueByAlias(
			ctx,
			taskstore.DefaultQueueName,
			taskAlias,
			map[string]any{
				"index":    "email",
				"truncate": "no",
			},
		)
		if err != nil && app.GetLogger() != nil {
			app.GetLogger().Error("adapters.NewOnUserUpdateFunc enqueue failed", slog.String("error", err.Error()))
		}
	}
}
