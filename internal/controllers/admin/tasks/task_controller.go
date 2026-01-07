package admin

import (
	"log/slog"
	"net/http"
	"project/internal/layouts"
	"project/internal/registry"

	"github.com/dracory/hb"
	"github.com/samber/lo"

	taskAdmin "github.com/dracory/taskstore/admin"
)

func NewTaskController(registry registry.RegistryInterface) *taskController {
	return &taskController{
		app:    registry,
		logger: registry.GetLogger(),
	}
}

type taskController struct {
	app    registry.RegistryInterface
	logger *slog.Logger
}

func (c *taskController) Handler(w http.ResponseWriter, r *http.Request) string {
	uptimeAdminUi, err := taskAdmin.UI(taskAdmin.UIOptions{
		ResponseWriter: w,
		Request:        r,
		Logger:         c.logger,
		Store:          c.app.GetTaskStore(),
		Layout:         &adminLayout{app: c.app},
	})

	ui := lo.IfF(err != nil, func() hb.TagInterface {
		c.logger.Error("At admin > taskController > Handler", "error", err.Error())
		return hb.Raw(err.Error())
	}).ElseF(func() hb.TagInterface {
		return uptimeAdminUi
	})

	return ui.ToHTML()
}

type adminLayout struct {
	app   registry.RegistryInterface
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
	return layouts.NewAdminLayout(a.app, r, layouts.Options{
		Title:      a.title,
		Content:    hb.Raw(a.body),
		ScriptURLs: a.scriptURLs,
		Scripts:    a.scripts,
		StyleURLs:  a.styleURLs,
		Styles:     a.styles,
	}).ToHTML()
}

var _ taskAdmin.Layout = (*adminLayout)(nil)
