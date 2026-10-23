package stats

import (
	"errors"
	"log/slog"
	"net/http"
	"project/internal/app"
	"project/internal/controllers/admin/adapters"
	"project/internal/helpers"
	"project/internal/links"

	basesession "github.com/dracory/base/session"

	statsadmin "github.com/dracory/statsadmin"
)

// statsController wraps the external github.com/dracory/statsadmin package
// for integration with the blueprint admin interface.
type statsController struct {
	logger *slog.Logger
	app    app.AppInterface
}

// NewStatsController creates a new stats admin controller
func NewStatsController(app app.AppInterface) *statsController {
	var logger *slog.Logger
	if app != nil {
		logger = app.GetLogger()
	}
	return &statsController{
		logger: logger,
		app:    app,
	}
}

// Handler processes stats admin requests
func (c *statsController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := statsadmin.New(statsadmin.AdminOptions{
		Store:             c.app.GetStatsStore(),
		Logger:            c.app.GetLogger(),
		FuncLayout:        adapters.NewLayoutFunc(c.app),
		HomeURL:           links.Admin().Home(),
		WebsiteUrl:        links.Website().Home(),
		BaseURL:           links.Admin().Stats(),
		CountryNameByIso2: c.countryNameByIso2(r),
		AuthUserID: func(r *http.Request) string {
			user := basesession.GetAuthUser(r)
			if user == nil {
				return ""
			}
			return user.GetID()
		},
		FlashError: func(w http.ResponseWriter, r *http.Request, message string, redirectURL string, delaySeconds int) string {
			return helpers.ToFlashError(c.app.GetCacheStore(), w, r, message, redirectURL, delaySeconds)
		},
	})

	if err != nil {
		if c.logger != nil {
			c.logger.Error("At admin > statsController > Handler", "error", err.Error())
		}
		helpers.ToFlashError(c.app.GetCacheStore(), w, r, err.Error(), links.Admin().Home(), 30)
		return
	}

	admin.Handle(w, r)
}

// countryNameByIso2 resolves an ISO2 country code to a human-readable
// country name using the blueprint's geostore.
func (c *statsController) countryNameByIso2(r *http.Request) func(iso2Code string) (string, error) {
	return func(iso2Code string) (string, error) {
		geostore := c.app.GetGeoStore()
		if geostore == nil {
			return "N/A", errors.New("geostore is nil")
		}
		country, err := geostore.CountryFindByIso2(r.Context(), iso2Code)
		if err != nil {
			return "", err
		}
		if country == nil {
			return "", errors.New("country not found for iso2: " + iso2Code)
		}
		return country.Name(), nil
	}
}
