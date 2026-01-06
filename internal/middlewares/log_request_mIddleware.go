package middlewares

import (
	"net/http"
	"project/internal/links"
	"project/internal/types"
	"slices"
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
func LogRequestMiddleware(app types.RegistryInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Log Request Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return logRequestHandler(app, next)
		})
}

func logRequestHandler(app types.RegistryInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := req.GetIP(r)

		method := r.Method

		rawPath := r.URL.Path
		if shouldSkipLogForPath(rawPath) {
			next.ServeHTTP(w, r)
			return
		}

		app.GetLogger().Info("["+method+" request by "+ip+"] "+r.RequestURI,
			slog.String("host", r.Host),
			slog.String("path", rawPath),
			slog.String("query", r.URL.RawQuery),
			slog.String("scheme", r.URL.Scheme),
			slog.String("ip", ip),
			slog.String("method", method),
			slog.String("proto", r.Proto),
			slog.String("user_agent", r.Header.Get("User-Agent")),
			slog.String("accept_language", r.Header.Get("Accept-Language")),
			slog.String("referer", r.Header.Get("Referer")),
		)

		next.ServeHTTP(w, r)
	})
}

func shouldSkipLogForPath(rawPath string) bool {
	path := strings.TrimLeft(rawPath, "/")

	skipPrefixes := []string{
		"th/",
	}

	skipSuffixes := []string{
		".css",
		".js",
		".png",
		".jpg",
		".jpeg",
		".ico",
		".svg",
		".woff",
		".woff2",
	}

	skipExact := []string{
		"health",
		links.LIVEFLUX,
		"ping",
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	for _, suffix := range skipSuffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}

	return slices.Contains(skipExact, path)
}
