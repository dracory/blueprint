package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/rtr"
)

// NewSecurityHeadersMiddleware creates middleware that sets security headers
func NewSecurityHeadersMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Security Headers Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Allowed domains for scripts, styles, fonts, and images
				scriptDomains := []string{
					"'self'",
					"'unsafe-inline'",
					"https://cdn.jsdelivr.net",
					"https://unpkg.com",
					"https://code.jquery.com",
					"https://cdnjs.cloudflare.com",
					"https://www.googletagmanager.com",
					"https://www.statcounter.com",
				}
				styleDomains := []string{
					"'self'",
					"'unsafe-inline'",
					"https://cdn.jsdelivr.net",
					"https://maxcdn.bootstrapcdn.com",
					"https://cdnjs.cloudflare.com",
					"https://fonts.googleapis.com",
				}
				fontDomains := []string{
					"'self'",
					"https://cdn.jsdelivr.net",
					"https://fonts.googleapis.com",
					"https://fonts.gstatic.com",
					"https://cdnjs.cloudflare.com",
					"https://maxcdn.bootstrapcdn.com",
				}
				imgDomains := []string{
					"'self'",
					"data:",
					"https://sfs.ams3.digitaloceanspaces.com",
					"https://lesichkov.ams3.digitaloceanspaces.com",
				}

				// Join arrays into CSP strings
				scriptSrc := strings.Join(scriptDomains, " ")
				styleSrc := strings.Join(styleDomains, " ")
				fontSrc := strings.Join(fontDomains, " ")
				imgSrc := strings.Join(imgDomains, " ")

				// Set security headers
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
				w.Header().Set("X-Frame-Options", "DENY")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
				w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'self'; script-src %s; style-src %s; font-src %s; img-src %s", scriptSrc, styleSrc, fontSrc, imgSrc))

				next.ServeHTTP(w, r)
			})
		})
}
