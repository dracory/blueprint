package stats

import (
	"log/slog"
	"net/http"
	"project/app/layouts"
	"project/app/links"
	"project/config"
	"project/internal/helpers"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/router"

	statsAdmin "github.com/gouniverse/statsstore/admin"
	statsAdminShared "github.com/gouniverse/statsstore/admin/shared"
)

func StatsController() router.ControllerInterface {
	return &statsController{
		logger: config.Logger,
	}
}

type statsController struct {
	logger slog.Logger
}

func (c *statsController) Handler(w http.ResponseWriter, r *http.Request) {
	visitorAnalyticsAdmin, err := statsAdmin.New(statsAdmin.Options{
		ResponseWriter: w,
		Request:        r,
		Logger:         &config.Logger,
		Store:          config.StatsStore,
		Layout:         &adminLayout{},
		HomeURL:        links.Admin().Home(),
		WebsiteUrl:     "https://lesichkov.co.uk",
	})

	if err != nil {
		c.logger.Error("At admin > statsController > Handler", "error", err.Error())
		helpers.ToFlashError(w, r, err.Error(), links.Admin().Home(), 30)
		return
	}

	visitorAnalyticsAdmin.ServeHTTP(w, r)
}

type adminLayout struct {
	title string
	body  string

	scriptURLs []string
	scripts    []string

	styleURLs []string
	styles    []string
}

func (a *adminLayout) SetTitle(title string) {
	a.title = title
}

func (a *adminLayout) SetBody(body string) {
	a.body = body
}

func (a *adminLayout) SetScriptURLs(urls []string) {
	a.scriptURLs = urls
}

func (a *adminLayout) SetScripts(scripts []string) {
	a.scripts = scripts
}

func (a *adminLayout) SetStyleURLs(urls []string) {
	a.styleURLs = urls
}

func (a *adminLayout) SetStyles(styles []string) {
	a.styles = styles
}

func (a *adminLayout) Render(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(r, layouts.Options{
		Title:      a.title,
		Content:    hb.Raw(a.body),
		ScriptURLs: a.scriptURLs,
		Scripts:    a.scripts,
		StyleURLs:  a.styleURLs,
		Styles:     a.styles,
	}).ToHTML()
}

var _ statsAdminShared.LayoutInterface = (*adminLayout)(nil)
