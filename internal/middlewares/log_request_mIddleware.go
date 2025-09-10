package middlewares

import (
	"net/http"
	"project/internal/types"
	"strings"

	"log/slog"

	"github.com/dracory/req"
	"github.com/dracory/rtr"
)

// LogRequestMiddleware logs every request to the database using the LogStore logger
// ==================================================================
// This is userful so that we can identify where all the visits
// come from and keep the application protected - i.e. bots,
// malicious spiders, DDOS, etc
// ==================================================================
// it is useful to detect spamming bots
func LogRequestMiddleware(app types.AppInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Log Request Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// uri := r.RequestURI

				ip := req.GetIP(r)

				method := r.Method

				app.GetLogger().Info("request",
					slog.String("host", r.Host),
					slog.String("path", strings.TrimLeft(r.URL.Path, "/")),
					slog.String("ip", ip),
					slog.String("method", method),
					slog.String("useragent", r.Header.Get("User-Agent")),
					slog.String("acceptlanguage", r.Header.Get("Accept-Language")),
					slog.String("referer", r.Header.Get("Referer")),
				)

				next.ServeHTTP(w, r)
			})
		},
		)
}
