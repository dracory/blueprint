package middlewares

import (
	"net/http"
	"slices"

	"github.com/dracory/rtr"
)

// NewHTTPSRedirectMiddleware creates middleware that redirects HTTP requests to HTTPS
func NewHTTPSRedirectMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("HTTPS Redirect Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Skip redirection for localhost and in development
				if slices.Contains([]string{
					"localhost",
					"127.0.0.",
					".local",
				}, r.Host) || r.TLS != nil {
					next.ServeHTTP(w, r)
					return
				}

				// Redirect to HTTPS version of same URL
				httpsURL := "https://" + r.Host + r.URL.Path
				if r.URL.RawQuery != "" {
					httpsURL += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			})
		})
}
