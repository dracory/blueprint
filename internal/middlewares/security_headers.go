package middlewares

import (
	"net/http"

	"github.com/dracory/rtr"
)

// NewSecurityHeadersMiddleware creates middleware that sets security headers
func NewSecurityHeadersMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Security Headers Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set security headers
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
				w.Header().Set("X-Frame-Options", "DENY")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
				w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:")

				next.ServeHTTP(w, r)
			})
		})
}
