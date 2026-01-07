package log_manager

import (
	"context"
	"log"
	"net/url"

	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
)

type logFilterComponent struct {
	liveflux.Base

	App              registry.RegistryInterface
	Level            string
	SearchMessage    string
	SearchContext    string
	SearchMessageNot string
	SearchContextNot string
	From             string
	To               string
	IsOpen           bool
	// RedirectURL, if set during Handle, will trigger a client-side redirect in Render
	RedirectURL string
}

func NewLogFilterComponent(registry registry.RegistryInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&logFilterComponent{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if c, ok := inst.(*logFilterComponent); ok {
		c.App = registry
	}

	return inst
}

func (c *logFilterComponent) GetKind() string {
	return "admin_log_manager_filter"
}

func (c *logFilterComponent) Mount(ctx context.Context, params map[string]string) error {
	c.Level = params[FILTER_LEVEL]
	c.SearchMessage = params[FILTER_SEARCH_MESSAGE]
	c.SearchContext = params[FILTER_SEARCH_CONTEXT]
	c.SearchMessageNot = params[FILTER_SEARCH_MESSAGE_NOT]
	c.SearchContextNot = params[FILTER_SEARCH_CONTEXT_NOT]
	c.From = params[FILTER_FROM]
	c.To = params[FILTER_TO]
	c.IsOpen = false

	return nil
}

func (c *logFilterComponent) Handle(ctx context.Context, action string, data url.Values) error {
	switch action {
	case "open":
		c.IsOpen = true
	case "close":
		c.IsOpen = false
	case "submit":
		// LiveFlux form submit: capture filters and redirect to full logs page with query params
		if data == nil {
			data = url.Values{}
		}

		c.Level = data.Get(FILTER_LEVEL)
		c.SearchMessage = data.Get(FILTER_SEARCH_MESSAGE)
		c.SearchContext = data.Get(FILTER_SEARCH_CONTEXT)
		c.SearchMessageNot = data.Get(FILTER_SEARCH_MESSAGE_NOT)
		c.SearchContextNot = data.Get(FILTER_SEARCH_CONTEXT_NOT)
		c.From = data.Get(FILTER_FROM)
		c.To = data.Get(FILTER_TO)

		q := url.Values{}
		if c.Level != "" {
			q.Set(FILTER_LEVEL, c.Level)
		}
		if c.SearchMessage != "" {
			q.Set(FILTER_SEARCH_MESSAGE, c.SearchMessage)
		}
		if c.SearchContext != "" {
			q.Set(FILTER_SEARCH_CONTEXT, c.SearchContext)
		}
		if c.SearchMessageNot != "" {
			q.Set(FILTER_SEARCH_MESSAGE_NOT, c.SearchMessageNot)
		}
		if c.SearchContextNot != "" {
			q.Set(FILTER_SEARCH_CONTEXT_NOT, c.SearchContextNot)
		}
		if c.From != "" {
			q.Set(FILTER_FROM, c.From)
		}
		if c.To != "" {
			q.Set(FILTER_TO, c.To)
		}

		base := links.Admin().Logs()
		if encoded := q.Encode(); encoded != "" {
			c.RedirectURL = base + "?" + encoded
		} else {
			c.RedirectURL = base
		}

		// Close modal on next render; the page will redirect shortly after
		c.IsOpen = false
	}

	return nil
}

func (c *logFilterComponent) Render(ctx context.Context) hb.TagInterface {
	filtersSummary := hb.Div().
		Class("d-flex flex-wrap align-items-center gap-2")

	activeFilters := false

	if c.Level != "" {
		activeFilters = true
		filtersSummary = filtersSummary.Child(
			hb.Span().
				Class("badge rounded-pill text-bg-secondary").
				Text("Level: " + c.Level),
		)
	}

	if c.SearchMessage != "" {
		activeFilters = true
		filtersSummary = filtersSummary.Child(
			hb.Span().
				Class("badge rounded-pill text-bg-secondary").
				Text("Message: " + c.SearchMessage),
		)
	}

	if c.SearchContext != "" {
		activeFilters = true
		filtersSummary = filtersSummary.Child(
			hb.Span().
				Class("badge rounded-pill text-bg-secondary").
				Text("Context: " + c.SearchContext),
		)
	}

	if c.SearchMessageNot != "" {
		activeFilters = true
		filtersSummary = filtersSummary.Child(
			hb.Span().
				Class("badge rounded-pill text-bg-secondary").
				Text("Message ≠ " + c.SearchMessageNot),
		)
	}

	if c.SearchContextNot != "" {
		activeFilters = true
		filtersSummary = filtersSummary.Child(
			hb.Span().
				Class("badge rounded-pill text-bg-secondary").
				Text("Context ≠ " + c.SearchContextNot),
		)
	}

	if c.From != "" {
		activeFilters = true
		filtersSummary = filtersSummary.Child(
			hb.Span().
				Class("badge rounded-pill text-bg-secondary").
				Text("From: " + c.From),
		)
	}

	if c.To != "" {
		activeFilters = true
		filtersSummary = filtersSummary.Child(
			hb.Span().
				Class("badge rounded-pill text-bg-secondary").
				Text("To: " + c.To),
		)
	}

	// Prefix text similar to "Showing users with status: ..."
	prefixText := "Showing all logs"
	if activeFilters {
		prefixText = "Showing logs with:"
	}

	filtersSummary = hb.Div().
		Class("d-flex flex-wrap align-items-center gap-2 ms-2").
		Child(hb.Span().Class("text-muted").Text(prefixText)).
		Child(filtersSummary)

	buttonFilters := hb.Button().
		Type("button").
		Class("btn btn-outline-secondary btn-sm").
		Attr(liveflux.DataFluxAction, "open").
		Attr(liveflux.DataFluxTargetKind, c.GetKind()).
		Attr(liveflux.DataFluxTargetID, c.GetID()).
		Attr(liveflux.DataFluxIndicator, "this").
		Child(hb.I().Class("bi bi-funnel me-1")).
		Child(hb.Span().Text("Filters"))

	filtersBar := hb.Div().
		Class("alert alert-info d-flex align-items-center gap-2 mb-3").
		Child(buttonFilters).
		Child(filtersSummary)

	// If a redirect URL is set (after a LiveFlux submit), emit a script to navigate there.
	if c.RedirectURL != "" {
		redirectScript := hb.Script("window.location.href = '" + c.RedirectURL + "';")
		return c.Root(hb.Div().
			Child(filtersBar).
			Child(redirectScript),
		)
	}

	if !c.IsOpen {
		return c.Root(filtersBar)
	}

	// Level select
	levelOptions := []hb.TagInterface{
		hb.Option().Value("").Selected(c.Level == "").Text("All levels"),
		hb.Option().Value("trace").Selected(c.Level == "trace").Text("Trace"),
		hb.Option().Value("debug").Selected(c.Level == "debug").Text("Debug"),
		hb.Option().Value("info").Selected(c.Level == "info").Text("Info"),
		hb.Option().Value("warn").Selected(c.Level == "warn").Text("Warn"),
		hb.Option().Value("error").Selected(c.Level == "error").Text("Error"),
		hb.Option().Value("fatal").Selected(c.Level == "fatal").Text("Fatal"),
		hb.Option().Value("panic").Selected(c.Level == "panic").Text("Panic"),
	}

	levelSelect := hb.Select().
		Class("form-select").
		ID(FILTER_LEVEL).
		Name(FILTER_LEVEL).
		Children(levelOptions)

	levelRow := hb.Div().
		Class("mb-3").
		Child(hb.Label().
			Class("form-label").
			For(FILTER_LEVEL).
			Text("Level")).
		Child(levelSelect)

	// Search message / context row
	searchMessageInput := hb.Input().
		Type("text").
		Class("form-control").
		ID(FILTER_SEARCH_MESSAGE).
		Name(FILTER_SEARCH_MESSAGE).
		Value(c.SearchMessage).
		Placeholder("Search in message")

	searchContextInput := hb.Input().
		Type("text").
		Class("form-control").
		ID(FILTER_SEARCH_CONTEXT).
		Name(FILTER_SEARCH_CONTEXT).
		Value(c.SearchContext).
		Placeholder("Search in context")

	searchRow := hb.Div().
		Class("row g-3 mb-3").
		Child(hb.Div().
			Class("col-12 col-md-6").
			Child(hb.Label().
				Class("form-label").
				For(FILTER_SEARCH_MESSAGE).
				Text("Search message")).
			Child(searchMessageInput),
		).
		Child(hb.Div().
			Class("col-12 col-md-6").
			Child(hb.Label().
				Class("form-label").
				For(FILTER_SEARCH_CONTEXT).
				Text("Search context")).
			Child(searchContextInput),
		)

	// Message not / Context not row
	messageNotInput := hb.Input().
		Type("text").
		Class("form-control").
		ID(FILTER_SEARCH_MESSAGE_NOT).
		Name(FILTER_SEARCH_MESSAGE_NOT).
		Value(c.SearchMessageNot).
		Placeholder("Exclude messages containing")

	contextNotInput := hb.Input().
		Type("text").
		Class("form-control").
		ID(FILTER_SEARCH_CONTEXT_NOT).
		Name(FILTER_SEARCH_CONTEXT_NOT).
		Value(c.SearchContextNot).
		Placeholder("Exclude contexts containing")

	notRow := hb.Div().
		Class("row g-3 mb-3").
		Child(hb.Div().
			Class("col-12 col-md-6").
			Child(hb.Label().
				Class("form-label").
				For(FILTER_SEARCH_MESSAGE_NOT).
				Text("Message not contains")).
			Child(messageNotInput),
		).
		Child(hb.Div().
			Class("col-12 col-md-6").
			Child(hb.Label().
				Class("form-label").
				For(FILTER_SEARCH_CONTEXT_NOT).
				Text("Context not contains")).
			Child(contextNotInput),
		)

	// From / To row
	fromInput := hb.Input().
		Type("text").
		Class("form-control").
		ID(FILTER_FROM).
		Name(FILTER_FROM).
		Value(c.From).
		Placeholder("e.g. 2025-11-16 21:50:00")

	toInput := hb.Input().
		Type("text").
		Class("form-control").
		ID(FILTER_TO).
		Name(FILTER_TO).
		Value(c.To).
		Placeholder("e.g. 2025-11-16 21:59:59")

	rangeRow := hb.Div().
		Class("row g-3 mb-3").
		Child(hb.Div().
			Class("col-12 col-md-6").
			Child(hb.Label().
				Class("form-label").
				For(FILTER_FROM).
				Text("From (YYYY-MM-DD HH:MM:SS)")).
			Child(fromInput),
		).
		Child(hb.Div().
			Class("col-12 col-md-6").
			Child(hb.Label().
				Class("form-label").
				For(FILTER_TO).
				Text("To (YYYY-MM-DD HH:MM:SS)")).
			Child(toInput),
		)

	footerButtons := hb.Div().
		Class("d-flex justify-content-end gap-2").
		Child(hb.Button().
			Type("button").
			Class("btn btn-outline-secondary me-auto").
			Attr(liveflux.DataFluxAction, "close").
			Attr(liveflux.DataFluxIndicator, "this").
			Child(hb.I().Class("bi bi-chevron-left me-1")).
			Child(hb.Span().Text("Cancel")),
		).
		Child(hb.Button().
			Type("submit").
			Class("btn btn-primary").
			Attr(liveflux.DataFluxAction, "submit").
			Attr(liveflux.DataFluxTargetKind, c.GetKind()).
			Attr(liveflux.DataFluxTargetID, c.GetID()).
			Attr(liveflux.DataFluxIndicator, "this").
			Child(hb.I().Class("bi bi-check2 me-1")).
			Child(hb.Span().Text("Apply")),
		)

	form := hb.Form().
		Method("GET").
		Action(links.Admin().Logs()).
		Child(levelRow).
		Child(searchRow).
		Child(notRow).
		Child(rangeRow).
		Child(footerButtons)

	modalContent := hb.Div().
		Class("modal-content").
		Child(hb.Div().
			Class("modal-header").
			Child(hb.Heading5().Class("modal-title mb-0").Text("Log Filters")).
			Child(hb.Button().
				Type("button").
				Class("btn-close").
				Attr(liveflux.DataFluxAction, "close")),
		).
		Child(hb.Div().
			Class("modal-body").
			Child(form),
		)

	modalDialog := hb.Div().
		Class("modal-dialog modal-xl modal-dialog-centered").
		Child(modalContent)

	modal := hb.Div().
		ID("LogManagerFilterModal").
		Class("modal fade show").
		Attr("tabindex", "-1").
		Style("display:block; background: rgba(0,0,0,0.5);").
		Child(hb.Div().
			Class("modal-dialog-wrapper d-flex align-items-center justify-content-center min-vh-100").
			Child(modalDialog),
		)

	return c.Root(hb.Div().
		Child(filtersBar).
		Child(modal),
	)
}

func init() {
	if err := liveflux.Register(&logFilterComponent{}); err != nil {
		log.Printf("Failed to register logFilterComponent: %v", err)
	}
}
