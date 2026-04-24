package middlewares

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

func NewEmailAllowlistMiddleware(registry registry.RegistryInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Email Allowlist Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user := helpers.GetAuthUser(r)

				// if the user is not authenticated, redirect to login
				if user == nil {
					helpers.ToFlashError(registry.GetCacheStore(), w, r, "Only authenticated users can access this page", links.AUTH_LOGIN, 15)
					return
				}

				allowedEmails := registry.GetConfig().GetEmailsAllowedAccess()

				// if the list is empty, all emails are allowed
				if len(allowedEmails) == 0 {
					next.ServeHTTP(w, r)
					return
				}

				email := user.Email()
				// Untokenize email if vault is enabled (with caching to reduce vault hits)
				if registry.GetConfig().GetUserStoreVaultEnabled() && email != "" {
					cacheKey := "email_untokenize:" + email
					// Try to get from cache first
					cachedEmail, err := registry.GetCacheStore().GetJSON(cacheKey, "")
					if err == nil && cachedEmail != "" {
						email = cachedEmail.(string)
					} else {
						// Cache miss - decrypt from vault
						untokenizedEmail, err := registry.GetVaultStore().TokenRead(r.Context(), email, registry.GetConfig().GetVaultStoreKey())
						if err == nil {
							email = untokenizedEmail
							// Cache for 5 minutes
							registry.GetCacheStore().SetJSON(cacheKey, email, 5*60)
						}
					}
				}

				found := false
				for _, allowed := range allowedEmails {
					if allowed == email {
						found = true
						break
					}
				}

				// if the email is not in the list, the user is not allowed to access the page
				if found {
					next.ServeHTTP(w, r)
					return
				}

				homeURL := links.Website().Home()
				helpers.ToFlashError(registry.GetCacheStore(), w, r, "Access restricted to authorized emails only", homeURL, 15)
			})
		})
}
