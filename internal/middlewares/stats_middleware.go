package middlewares

import (
	"log/slog"
	"net/http"
	"project/internal/config"

	"github.com/dracory/base/req"
	"github.com/dracory/rtr"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/statsstore"
)

func NewStatsMiddleware() rtr.MiddlewareInterface {
	stats := new(statsMiddleware)

	return rtr.NewMiddleware().
		SetName(stats.Name()).
		SetHandler(stats.Handler)
}

type statsMiddleware struct{}

func (m statsMiddleware) Name() string {
	return "Stats Middleware"
}

func (m statsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.StatsStoreUsed {
			next.ServeHTTP(w, r)
			return
		}

		if config.StatsStore == nil {
			config.Logger.Error("stats_middleware", "error", "stats store is marked as used but is nil")
			next.ServeHTTP(w, r)
			return
		}

		ip := req.IP(r)
		userAgent := r.UserAgent()
		userAcceptLanguage := r.Header.Get("Accept-Language")
		country := "" // empty by default (will be filled in later in the backend)
		userReferer := r.Header.Get("Referer")
		userAcceptEncoding := r.Header.Get("Accept-Encoding")
		// userRequestedWith := r.Header.Get("X-Requested-With")
		// userIsBot := r.Header.Get("X-Bot")

		if config.AppEnvironment == config.APP_ENVIRONMENT_TESTING {
			ip = "127.0.0.1"
			userAcceptLanguage = "us"
			userAgent = "testing"
			country = "us"
			userReferer = "testing"
			userAcceptEncoding = "testing"
		}

		visitor := statsstore.NewVisitor()
		visitor.SetCountry(country)
		visitor.SetIpAddress(ip)
		visitor.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
		visitor.SetUserAgent(userAgent)
		visitor.SetUserAcceptLanguage(userAcceptLanguage)
		visitor.SetUserAcceptEncoding(userAcceptEncoding)
		visitor.SetUserReferrer(userReferer)
		visitor.SetPath("[" + r.Method + "] " + r.RequestURI)

		err := config.StatsStore.VisitorCreate(r.Context(), visitor)

		if err != nil {
			config.Logger.Error("Error at statsMiddleware", slog.String("error", err.Error()))
		}

		next.ServeHTTP(w, r)
	})
}
