package middlewares

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

// allowedEmails is a map of allowed emails
// if the email is not in the map, the user is not allowed to access the page
// this is a simple way to restrict access to certain pages to certain emails
// if the map is empty, all emails are allowed
var allowedEmails = map[string]struct{}{
	"info@sinevia.com": {},
}

func NewEmailAllowlistMiddleware(app types.AppInterface) rtr.MiddlewareInterface {
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

				// if the map is empty, all emails are allowed
				if len(allowedEmails) == 0 {
					next.ServeHTTP(w, r)
					return
				}

				email := user.Email()
				_, found := allowedEmails[email]

				// if the email is not in the map, the user is not allowed to access the page
				if found {
					next.ServeHTTP(w, r)
					return
				}

				homeURL := links.Website().Home()
				helpers.ToFlashError(app.GetCacheStore(), w, r, "Access restricted to authorized emails only", homeURL, 15)
			})
		})
}
