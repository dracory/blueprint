package log_manager

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strconv"

	livefluxctl "project/internal/controllers/liveflux"
	"project/internal/types"

	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/logstore"
)

type logTableComponent struct {
	liveflux.Base

	App              types.AppInterface
	Level            string
	SearchMessage    string
	SearchContext    string
	SearchMessageNot string
	SearchContextNot string
	From             string
	To               string
	SortBy           string
	SortDirection    string
	Logs             []logstore.LogInterface
	Page             int
	PerPage          int
	HasMore          bool
	Total            int
	// Modal state for viewing a single log's context
	IsContextOpen bool
	ContextLog    logstore.LogInterface
}

func NewLogTableComponent(app types.AppInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&logTableComponent{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if c, ok := inst.(*logTableComponent); ok {
		c.App = app
	}

	return inst
}

func (c *logTableComponent) GetKind() string {
	return "admin_log_manager_table"
}

func (c *logTableComponent) Mount(ctx context.Context, params map[string]string) error {
	// Ensure App is set when component is instantiated via Liveflux placeholder
	if c.App == nil {
		if app, ok := ctx.Value(livefluxctl.AppContextKey).(types.AppInterface); ok {
			c.App = app
		}
	}

	c.Level = params[FILTER_LEVEL]
	c.SearchMessage = params[FILTER_SEARCH_MESSAGE]
	c.SearchContext = params[FILTER_SEARCH_CONTEXT]
	c.SearchMessageNot = params[FILTER_SEARCH_MESSAGE_NOT]
	c.SearchContextNot = params[FILTER_SEARCH_CONTEXT_NOT]
	c.From = params[FILTER_FROM]
	c.To = params[FILTER_TO]
	c.SortBy = params[SORT_BY]
	c.SortDirection = params[SORT_DIRECTION]

	// Pagination defaults
	c.Page = 0
	c.PerPage = 100

	if c.SortBy == "" {
		c.SortBy = logstore.COLUMN_TIME
	}

	if c.SortDirection == "" {
		c.SortDirection = "desc"
	}

	return c.loadLogs()
}

func (c *logTableComponent) Handle(ctx context.Context, action string, data url.Values) error {
	switch action {
	case "sort":
		if data == nil {
			data = url.Values{}
		}

		by := data.Get(SORT_BY)
		if by == "" {
			by = logstore.COLUMN_TIME
		}

		if c.SortBy == by {
			if c.SortDirection == "asc" {
				c.SortDirection = "desc"
			} else {
				c.SortDirection = "asc"
			}
		} else {
			c.SortBy = by
			c.SortDirection = "asc"
		}

		return c.loadLogs()

	case "page":
		if data == nil {
			data = url.Values{}
		}

		pageStr := data.Get(PAGE)
		if pageStr == "" {
			return nil
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			return nil
		}

		c.Page = page

		return c.loadLogs()

	case "per_page":
		if data == nil {
			data = url.Values{}
		}

		perPageStr := data.Get(PER_PAGE)
		if perPageStr == "" {
			return nil
		}

		perPage, err := strconv.Atoi(perPageStr)
		if err != nil || perPage <= 0 {
			return nil
		}

		c.PerPage = perPage
		c.Page = 0

		return c.loadLogs()

	case ACTION_DELETE:
		if data == nil {
			data = url.Values{}
		}

		logID := data.Get(PARAM_LOG_ID)
		if logID == "" {
			return nil
		}

		if c.App == nil || c.App.GetLogStore() == nil {
			return nil
		}

		if err := c.App.GetLogStore().LogDeleteByID(logID); err != nil {
			return nil
		}

		return c.loadLogs()

	case ACTION_SHOW_CONTEXT:
		if data == nil {
			data = url.Values{}
		}

		logID := data.Get(PARAM_LOG_ID)
		if logID == "" {
			return nil
		}

		if c.App == nil || c.App.GetLogStore() == nil {
			return nil
		}

		logEntry, err := c.App.GetLogStore().LogFindByID(logID)
		if err != nil || logEntry == nil {
			return nil
		}

		c.ContextLog = logEntry
		c.IsContextOpen = true

		return nil

	case ACTION_CLOSE_CONTEXT:
		c.ContextLog = nil
		c.IsContextOpen = false
		return nil

	case ACTION_DELETE_SELECTED:
		if data == nil {
			data = url.Values{}
		}

		if c.App == nil || c.App.GetLogStore() == nil {
			return nil
		}

		ids := data[PARAM_BULK_LOG_IDS]
		if len(ids) == 0 {
			return nil
		}

		for _, id := range ids {
			if id == "" {
				continue
			}
			_ = c.App.GetLogStore().LogDeleteByID(id)
		}

		return c.loadLogs()
	}

	return nil
}

func (c *logTableComponent) loadLogs() error {
	res, err := listLogs(c.App, logListFilters{
		FilterLevel:            c.Level,
		FilterSearchMessage:    c.SearchMessage,
		FilterSearchContext:    c.SearchContext,
		FilterSearchMessageNot: c.SearchMessageNot,
		FilterSearchContextNot: c.SearchContextNot,
		FilterFrom:             c.From,
		FilterTo:               c.To,
		FilterSortBy:           c.SortBy,
		FilterSortDirection:    c.SortDirection,
		FilterPage:             c.Page,
		FilterPerPage:          c.PerPage,
	})
	if err != nil {
		return nil
	}

	c.Logs = res.Logs
	c.Total = res.Total
	c.HasMore = res.HasMore

	return nil
}

func (c *logTableComponent) Render(ctx context.Context) hb.TagInterface {
	table := hb.Table().Class("table table-striped table-hover align-middle mb-0")

	thead := c.renderTableHeader()
	tbody := c.renderTableBody()
	table = table.Child(thead).Child(tbody)

	footer := c.renderTableFooter()
	modal := c.renderContextModal()
	bulkActions := c.renderBulkActionsBar()

	content := hb.Div().
		Child(bulkActions).
		Child(table).
		Child(footer).
		Child(modal)

	return c.Root(content)
}

func (c *logTableComponent) renderTableHeader() hb.TagInterface {
	arrow := func(column string) string {
		if c.SortBy != column {
			return ""
		}

		if c.SortDirection == "asc" {
			return " \u2191"
		}

		return " \u2193"
	}

	return hb.Thead().
		Child(hb.Tr().
			Child(hb.Th().
				Class("text-center").
				Child(
					// Select-all checkbox
					hb.Input().
						Type("checkbox").
						Class("form-check-input").
						Attr("data-role", "log-select-all").
						OnChange(`const checked = this.checked; const table = this.closest('table'); if (!table) return; table.querySelectorAll("tbody input[data-role='log-select']").forEach(function(cb){ cb.checked = checked; }); const root = this.closest('[data-liveflux-root]') || document; const checkedBoxes = root.querySelectorAll("input[data-role='log-select']:checked"); const count = checkedBoxes.length; const btn = root.querySelector("[data-role='bulk-delete-selected-button']"); if (btn) { const labelSpan = btn.querySelector('span'); if (labelSpan) { labelSpan.textContent = count > 0 ? 'Delete selected (' + count + ')' : 'Delete selected'; } if (count > 0) { btn.classList.remove('d-none'); btn.removeAttribute('disabled'); } else { btn.classList.add('d-none'); btn.setAttribute('disabled','disabled'); } }`),
				),
			).
			Child(hb.Th().
				Child(hb.Button().
					Type("submit").
					Class("btn btn-link p-0 text-decoration-none").
					Attr(liveflux.DataFluxAction, "sort").
					Attr(liveflux.DataFluxTargetKind, c.GetKind()).
					Attr(liveflux.DataFluxTargetID, c.GetID()).
					Name(SORT_BY).
					Value(logstore.COLUMN_TIME).
					Text("Time" + arrow(logstore.COLUMN_TIME)),
				),
			).
			Child(hb.Th().
				Child(hb.Button().
					Type("submit").
					Class("btn btn-link p-0 text-decoration-none").
					Attr(liveflux.DataFluxAction, "sort").
					Attr(liveflux.DataFluxTargetKind, c.GetKind()).
					Attr(liveflux.DataFluxTargetID, c.GetID()).
					Name(SORT_BY).
					Value(logstore.COLUMN_LEVEL).
					Text("Level" + arrow(logstore.COLUMN_LEVEL)),
				),
			).
			Child(hb.Th().
				Child(hb.Button().
					Type("submit").
					Class("btn btn-link p-0 text-decoration-none").
					Attr(liveflux.DataFluxAction, "sort").
					Attr(liveflux.DataFluxTargetKind, c.GetKind()).
					Attr(liveflux.DataFluxTargetID, c.GetID()).
					Name(SORT_BY).
					Value(logstore.COLUMN_MESSAGE).
					Text("Message" + arrow(logstore.COLUMN_MESSAGE)),
				),
			).
			Child(hb.Th().Class("text-end").Text("Actions")),
		)
}

func (c *logTableComponent) renderTableBody() hb.TagInterface {
	tbody := hb.Tbody()

	for _, l := range c.Logs {
		selectCheckbox := hb.Input().
			Type("checkbox").
			Class("form-check-input").
			Attr("data-role", "log-select").
			Name(PARAM_BULK_LOG_IDS).
			Value(l.GetID()).
			OnChange(`const root = this.closest('[data-liveflux-root]') || document; const checkedBoxes = root.querySelectorAll("input[data-role='log-select']:checked"); const count = checkedBoxes.length; const btn = root.querySelector("[data-role='bulk-delete-selected-button']"); if (btn) { const labelSpan = btn.querySelector('span'); if (labelSpan) { labelSpan.textContent = count > 0 ? 'Delete selected (' + count + ')' : 'Delete selected'; } if (count > 0) { btn.classList.remove('d-none'); btn.removeAttribute('disabled'); } else { btn.classList.add('d-none'); btn.setAttribute('disabled','disabled'); } }`)

		contextButton := hb.Button().
			Type("submit").
			Class("btn btn-sm btn-outline-secondary me-1").
			Attr(liveflux.DataFluxAction, ACTION_SHOW_CONTEXT).
			Attr(liveflux.DataFluxTargetKind, c.GetKind()).
			Attr(liveflux.DataFluxTargetID, c.GetID()).
			Name(PARAM_LOG_ID).
			Value(l.GetID()).
			Child(hb.I().Class("bi bi-eye"))

		deleteButton := hb.Button().
			Type("button").
			Class("btn btn-sm btn-outline-danger").
			OnClick("Swal.fire({title: 'Delete log?', text: 'This action cannot be undone.', icon: 'warning', showCancelButton: true, confirmButtonText: 'Delete', cancelButtonText: 'Cancel'}).then((result) => { if (result.isConfirmed) { const hiddenBtn = this.closest('td').querySelector('[data-role=\\'delete-submit\\']'); if (hiddenBtn) { hiddenBtn.click(); } } });").
			Child(hb.I().Class("bi bi-trash"))

		hiddenSubmitButton := hb.Button().
			Type("submit").
			Class("d-none").
			Attr("data-role", "delete-submit").
			Attr(liveflux.DataFluxAction, ACTION_DELETE).
			Attr(liveflux.DataFluxTargetKind, c.GetKind()).
			Attr(liveflux.DataFluxTargetID, c.GetID()).
			Name(PARAM_LOG_ID).
			Value(l.GetID())

		actionsCell := hb.Td().Class("text-end").
			Child(contextButton).
			Child(deleteButton).
			Child(hiddenSubmitButton)

		row := hb.Tr().
			Child(hb.Td().Class("text-center").Child(selectCheckbox)).
			Child(hb.Td().Text(l.GetTimeCarbon().ToDateTimeString())).
			Child(hb.Td().Text(l.GetLevel())).
			Child(hb.Td().Text(l.GetMessage())).
			Child(actionsCell)

		tbody = tbody.Child(row)
	}

	return tbody
}

func (c *logTableComponent) renderBulkActionsBar() hb.TagInterface {
	spinner := hb.Span().
		Class("logs-bulk-spinner spinner-border spinner-border-sm align-middle ms-2").
		Style("display: none;").
		Attr("role", "status").
		Child(hb.Span().Class("visually-hidden").Text("Loading"))

	deleteSelectedBtn := hb.Button().
		Type("submit").
		Class("btn btn-sm btn-outline-danger d-none").
		Attr("data-role", "bulk-delete-selected-button").
		Attr(liveflux.DataFluxAction, ACTION_DELETE_SELECTED).
		Attr(liveflux.DataFluxTargetKind, c.GetKind()).
		Attr(liveflux.DataFluxTargetID, c.GetID()).
		Attr("disabled", "disabled").
		Attr(liveflux.DataFluxIndicator, "this, .logs-bulk-spinner").
		Attr(liveflux.DataFluxInclude, "input[data-role='log-select']:checked").
		Child(hb.I().Class("bi bi-trash me-1")).
		Child(hb.Span().Text("Delete selected")).
		Child(spinner)

	bar := hb.Div().
		Class("d-flex justify-content-between align-items-center mb-2 gap-2")

	left := hb.Div().
		Class("d-flex align-items-center gap-2").
		Child(deleteSelectedBtn)

	right := hb.Div()

	bar = bar.
		Child(left).
		Child(right)

	return bar
}

func (c *logTableComponent) renderContextModal() hb.TagInterface {
	if !c.IsContextOpen || c.ContextLog == nil {
		return hb.Div()
	}

	id := c.ContextLog.GetID()
	timeStr := ""
	if c.ContextLog.GetTimeCarbon() != nil {
		timeStr = c.ContextLog.GetTimeCarbon().ToDateTimeString()
	}
	message := c.ContextLog.GetMessage()
	contextRaw := c.ContextLog.GetContext()
	if contextRaw == "" {
		contextRaw = "(no context)"
	}

	contextPretty := contextRaw
	var anyJSON any
	if err := json.Unmarshal([]byte(contextRaw), &anyJSON); err == nil {
		if b, err := json.MarshalIndent(anyJSON, "", "  "); err == nil {
			contextPretty = string(b)
		}
	}

	metaRow := func(label, value string, striped bool) hb.TagInterface {
		rowClass := "d-flex justify-content-between px-2 py-1 small"
		if striped {
			rowClass += " bg-light"
		}
		return hb.Div().
			Class(rowClass).
			Child(hb.Span().Class("fw-bold me-2").Text(label)).
			Child(hb.Span().Class("text-monospace text-break").Text(value))
	}

	body := hb.Div().
		Class("modal-body").
		Child(
			hb.Div().
				Class("border rounded mb-3").
				Child(metaRow("ID", id, false)).
				Child(metaRow("Time", timeStr, true)).
				Child(metaRow("Message", message, false)),
		).
		Child(
			hb.PRE().
				Class("mb-0 small").
				Text(contextPretty),
		)

	footer := hb.Div().
		Class("modal-footer d-flex justify-content-end").
		Child(
			hb.Button().
				Type("button").
				Class("btn btn-secondary").
				Attr(liveflux.DataFluxAction, ACTION_CLOSE_CONTEXT).
				Attr(liveflux.DataFluxTargetKind, c.GetKind()).
				Attr(liveflux.DataFluxTargetID, c.GetID()).
				Child(hb.I().Class("bi bi-x-lg me-1")).
				Child(hb.Span().Text("Close")),
		)

	modalContent := hb.Div().
		Class("modal-content").
		Child(
			hb.Div().
				Class("modal-header").
				Child(hb.Heading5().Class("modal-title mb-0").Text("Log context")).
				Child(hb.Button().
					Type("button").
					Class("btn-close").
					Attr(liveflux.DataFluxAction, ACTION_CLOSE_CONTEXT).
					Attr(liveflux.DataFluxTargetKind, c.GetKind()).
					Attr(liveflux.DataFluxTargetID, c.GetID())),
		).
		Child(body).
		Child(footer)

	modalDialog := hb.Div().
		Class("modal-dialog modal-lg modal-dialog-centered").
		Child(modalContent)

	return hb.Div().
		Class("modal fade show").
		Attr("tabindex", "-1").
		Style("display:block; background: rgba(0,0,0,0.5);").
		Child(
			hb.Div().
				Class("modal-dialog-wrapper d-flex align-items-center justify-content-center min-vh-100").
				Child(modalDialog),
		)
}

func (c *logTableComponent) renderTableFooter() hb.TagInterface {
	footer := hb.Div().Class("d-flex flex-wrap justify-content-between align-items-center mt-3 gap-2")

	left := hb.Div()
	if len(c.Logs) == 0 {
		left = left.Child(hb.Span().Class("text-muted").Text("No logs to display"))
	} else {
		start := c.Page*c.PerPage + 1
		end := c.Page*c.PerPage + len(c.Logs)

		text := "Showing " + strconv.Itoa(start) + "–" + strconv.Itoa(end)
		if c.Total > 0 {
			text += " of " + strconv.Itoa(c.Total)
		}
		text += " logs"

		left = left.Child(
			hb.Span().Class("text-muted").Text(text),
		)
	}

	// Pagination controls: Bootstrap btn-group with nested page dropdown
	controls := hb.Div().Class("d-flex align-items-center gap-2 flex-wrap")

	totalPages := 0
	if c.PerPage > 0 && c.Total > 0 {
		totalPages = (c.Total + c.PerPage - 1) / c.PerPage
	}
	if totalPages < 1 {
		totalPages = 1
	}

	pageForm := hb.Form().
		Class("btn-group btn-group-sm align-items-stretch").
		Attr("role", "group").
		Attr("aria-label", "Page navigation")

	// Prev button
	prevBtn := hb.Button().
		Type("submit").
		Class("btn btn-outline-secondary btn-sm").
		Attr(liveflux.DataFluxAction, "page").
		Attr(liveflux.DataFluxTargetKind, c.GetKind()).
		Attr(liveflux.DataFluxTargetID, c.GetID()).
		Name(PAGE).
		Value(strconv.Itoa(c.Page - 1)).
		Text("«")
	if c.Page <= 0 {
		prevBtn = prevBtn.Attr("disabled", "disabled")
	}
	pageForm = pageForm.Child(prevBtn)

	// Dropdown in the middle
	dropdownID := "AdminLogPageDropdown"
	currentPageLabel := strconv.Itoa(c.Page + 1)
	dropdownToggle := hb.Button().
		Type("button").
		Class("btn btn-primary btn-sm dropdown-toggle").
		Attr("id", dropdownID).
		Attr("data-bs-toggle", "dropdown").
		Attr("aria-expanded", "false").
		Text(currentPageLabel)

	menu := hb.NewUL().
		Class("dropdown-menu").
		Attr("aria-labelledby", dropdownID)

	for i := 0; i < totalPages; i++ {
		pageNum := i + 1
		itemBtn := hb.Button().
			Type("submit").
			Class("dropdown-item btn-sm").
			Attr(liveflux.DataFluxAction, "page").
			Attr(liveflux.DataFluxTargetKind, c.GetKind()).
			Attr(liveflux.DataFluxTargetID, c.GetID()).
			Name(PAGE).
			Value(strconv.Itoa(i)).
			Text(strconv.Itoa(pageNum))

		menu = menu.Child(
			hb.NewLI().Child(itemBtn),
		)
	}

	dropdownGroup := hb.Div().
		Class("btn-group").
		Attr("role", "group").
		Child(dropdownToggle).
		Child(menu)

	pageForm = pageForm.Child(dropdownGroup)

	// Next button
	hasNext := (c.Page+1)*c.PerPage < c.Total
	nextBtn := hb.Button().
		Type("submit").
		Class("btn btn-outline-secondary btn-sm").
		Attr(liveflux.DataFluxAction, "page").
		Attr(liveflux.DataFluxTargetKind, c.GetKind()).
		Attr(liveflux.DataFluxTargetID, c.GetID()).
		Name(PAGE).
		Value(strconv.Itoa(c.Page + 1)).
		Text("»")
	if !hasNext {
		nextBtn = nextBtn.Attr("disabled", "disabled")
	}
	pageForm = pageForm.Child(nextBtn)

	controls = controls.Child(
		hb.Div().
			Class("d-flex align-items-center gap-2").
			Child(hb.Span().Class("text-muted small mb-0 text-nowrap").Text("Page")).
			Child(pageForm),
	)

	perPageForm := hb.Form().Class("d-flex align-items-center gap-2 flex-nowrap")
	perPageForm = perPageForm.Child(hb.Span().Class("text-muted small mb-0 text-nowrap").Text("Per page"))

	perPageSelect := hb.Select().
		Class("form-select form-select-sm").
		Name(PER_PAGE).
		OnChange("this.form.querySelector('[" + liveflux.DataFluxAction + "=\\'per_page\\']').click();")

	perPageValues := []int{100, 200, 300, 500, 1000, 3000, 5000, 10000}
	for _, v := range perPageValues {
		perPageSelect = perPageSelect.Child(
			hb.Option().
				Value(strconv.Itoa(v)).
				Text(strconv.Itoa(v)).
				AttrIf(c.PerPage == v, "selected", "selected"),
		)
	}

	perPageForm = perPageForm.Child(perPageSelect)
	perPageForm = perPageForm.Child(
		hb.Button().
			Type("submit").
			Class("d-none").
			Attr(liveflux.DataFluxAction, "per_page").
			Attr(liveflux.DataFluxTargetKind, c.GetKind()).
			Attr(liveflux.DataFluxTargetID, c.GetID()).
			Attr(liveflux.DataFluxIndicator, "this"),
	)

	footer = footer.
		Child(left).
		Child(controls).
		Child(perPageForm)

	return footer
}

func init() {
	if err := liveflux.Register(&logTableComponent{}); err != nil {
		log.Printf("Failed to register logTableComponent: %v", err)
	}
}
