package middlewares

import (
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	"github.com/dracory/rtr"
)

func NewEmailAllowlistMiddleware(app app.AppInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Email Allowlist Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user := helpers.GetAuthUser(r)

				// if the user is not authenticated, redirect to login
				if user == nil {
					helpers.ToFlashError(app.GetCacheStore(), w, r, "Only authenticated users can access this page", links.AUTH_LOGIN, 15)
					return
				}

				allowedEmails := app.GetConfig().GetEmailsAllowedAccess()

				// if the list is empty, all emails are allowed
				if len(allowedEmails) == 0 {
					next.ServeHTTP(w, r)
					return
				}

				email := user.GetEmail()
				// Untokenize email if vault is enabled (with caching to reduce vault hits)
				if app.GetConfig().GetUserStoreVaultEnabled() && email != "" {
					cacheKey := "email_untokenize:" + email
					// Try to get from cache first
					cachedEmail, err := app.GetCacheStore().GetJSON(cacheKey, "")
					if err == nil && cachedEmail != "" {
						email = cachedEmail.(string)
					} else {
						// Cache miss - decrypt from vault
						untokenizedEmail, err := app.GetVaultStore().TokenRead(r.Context(), email, app.GetConfig().GetVaultStoreKey())
						if err == nil {
							email = untokenizedEmail
							// Cache for 5 minutes
							app.GetCacheStore().SetJSON(cacheKey, email, 5*60)
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
				helpers.ToFlashError(app.GetCacheStore(), w, r, "Access restricted to authorized emails only", homeURL, 15)
			})
		})
}
