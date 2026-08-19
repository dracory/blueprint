// Package adapters bridges the external blogadmin and shopadmin packages
// to the blueprint's internal services (admin layout, LLM config, userstore
// for customer resolution).
//
// The external packages are decoupled from any specific host project, so
// they accept callback functions and interfaces (FuncLayout, LlmFactory,
// CustomerResolverInterface). This package provides concrete adapters
// wired to the blueprint's app.AppInterface.
package adapters

import (
	"context"
	"net/http"
	"strings"

	"project/internal/app"
	"project/internal/layouts"

	"github.com/dracory/hb"
	"github.com/dracory/llm"
	"github.com/dracory/userstore"
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
