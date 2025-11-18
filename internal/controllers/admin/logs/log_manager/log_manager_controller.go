package log_manager

import (
	"net/http"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
)

// == CONTROLLER ==============================================================

type logManagerController struct{ app types.AppInterface }

// == CONSTRUCTOR =============================================================

func NewLogManagerController(app types.AppInterface) *logManagerController {
	return &logManagerController{app: app}
}

func (controller *logManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, errorMessage, links.Admin().Home(), 10)
	}

	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Log Manager | Logs",
		Content: controller.page(data),
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Scripts: []string{
			liveflux.Script().ToHTML(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *logManagerController) page(data logManagerData) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Dashboard",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Logs",
			URL:  links.Admin().Logs(),
		},
	})

	// filterComponent := NewLogFilterComponent(controller.app)
	// filterSSR := liveflux.SSR(filterComponent, map[string]string{
	// 	"level":              data.FilterLevel,
	// 	"search_message":     data.FilterSearchMessage,
	// 	"search_context":     data.FilterSearchContext,
	// 	"search_message_not": data.FilterSearchMessageNot,
	// 	"search_context_not": data.FilterSearchContextNot,
	// 	"from":               data.FilterFrom,
	// 	"to":                 data.FilterTo,
	// })

	// tableComponent := NewLogTableComponent(controller.app)
	// tableSSR := liveflux.SSR(tableComponent, map[string]string{
	// 	"level":              data.FilterLevel,
	// 	"search_message":     data.FilterSearchMessage,
	// 	"search_context":     data.FilterSearchContext,
	// 	"search_message_not": data.FilterSearchMessageNot,
	// 	"search_context_not": data.FilterSearchContextNot,
	// 	"from":               data.FilterFrom,
	// 	"to":                 data.FilterTo,
	// })

	filterComponent := NewLogFilterComponent(controller.app)
	filterSSR := liveflux.Placeholder(filterComponent, map[string]string{
		"level":              data.FilterLevel,
		"search_message":     data.FilterSearchMessage,
		"search_context":     data.FilterSearchContext,
		"search_message_not": data.FilterSearchMessageNot,
		"search_context_not": data.FilterSearchContextNot,
		"from":               data.FilterFrom,
		"to":                 data.FilterTo,
	})

	tableComponent := NewLogTableComponent(controller.app)
	tableSSR := liveflux.Placeholder(tableComponent, map[string]string{
		"level":              data.FilterLevel,
		"search_message":     data.FilterSearchMessage,
		"search_context":     data.FilterSearchContext,
		"search_message_not": data.FilterSearchMessageNot,
		"search_context_not": data.FilterSearchContextNot,
		"from":               data.FilterFrom,
		"to":                 data.FilterTo,
	})

	card := hb.Div().
		Class("card shadow-sm w-100 mb-5").
		Child(
			hb.Div().
				Class("card-body").
				Child(
					hb.Div().
						Class("d-flex justify-content-between align-items-center mb-3").
						Child(hb.Heading1().Class("h3 mb-0").Text("Log Manager")),
				).
				Child(filterSSR).
				Child(tableSSR),
		)

	return hb.Div().
		Class("container min-vh-100 py-4").
		Child(breadcrumbs).
		Child(card)
}

func (controller *logManagerController) prepareData(r *http.Request) (logManagerData, string) {
	if controller.app.GetLogStore() == nil {
		return logManagerData{}, "Log store is not initialized"
	}

	level := req.GetStringTrimmed(r, FILTER_LEVEL)
	searchMessage := req.GetStringTrimmed(r, FILTER_SEARCH_MESSAGE)
	searchContext := req.GetStringTrimmed(r, FILTER_SEARCH_CONTEXT)
	searchMessageNot := req.GetStringTrimmed(r, FILTER_SEARCH_MESSAGE_NOT)
	searchContextNot := req.GetStringTrimmed(r, FILTER_SEARCH_CONTEXT_NOT)
	from := req.GetStringTrimmed(r, FILTER_FROM)
	to := req.GetStringTrimmed(r, FILTER_TO)

	return logManagerData{
		Request:                r,
		FilterLevel:            level,
		FilterSearchMessage:    searchMessage,
		FilterSearchContext:    searchContext,
		FilterSearchMessageNot: searchMessageNot,
		FilterSearchContextNot: searchContextNot,
		FilterFrom:             from,
		FilterTo:               to,
	}, ""
}

type logManagerData struct {
	Request                *http.Request
	FilterLevel            string
	FilterSearchMessage    string
	FilterSearchContext    string
	FilterSearchMessageNot string
	FilterSearchContextNot string
	FilterFrom             string
	FilterTo               string
}
