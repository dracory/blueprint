package middlewares

import (
	"net/http"

	"github.com/dracory/rtr"
)

// NewHTTPSRedirectMiddleware creates middleware that redirects HTTP requests to HTTPS
func NewHTTPSRedirectMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("HTTPS Redirect Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Skip redirection for localhost and in development
				host := r.Host
				isLocal := host == "localhost" ||
					host == "127.0.0.1" ||
					host == "0.0.0.0" ||
					len(host) > 6 && host[len(host)-6:] == ".local" || // Ends with .local
					len(host) > 9 && host[:9] == "127.0.0." || // Starts with 127.0.0.
					len(host) > 10 && host[:10] == "192.168." || // Starts with 192.168.
					len(host) > 7 && host[:7] == "10.0.0." || // Starts with 10.0.0.
					r.TLS != nil

				if isLocal {
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
