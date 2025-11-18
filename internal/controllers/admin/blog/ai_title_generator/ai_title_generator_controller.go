package aititlegenerator

import (
	"fmt"
	"net/http"
	"project/internal/controllers/admin/blog/shared"
	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"
	"project/pkg/blogai"
	"strings"

	"github.com/dracory/cdn"
	"github.com/dracory/customstore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
)

const (
	ACTION_ADD_TITLE       = "add_title"
	ACTION_GENERATE_TITLES = "generate_titles"
	ACTION_APPROVE_TITLE   = "approve_title"
	ACTION_REJECT_TITLE    = "reject_title"
	ACTION_GENERATE_POST   = "generate_post"
	ACTION_DELETE_TITLE    = "delete_title"
)

type AiTitleGeneratorController struct {
	app types.AppInterface
}

type pageData struct {
	Request             *http.Request
	Action              string
	ExistingPostRecords []blogai.RecordPost
	HasSystemPrompt     bool
}

func NewAiTitleGeneratorController(app types.AppInterface) *AiTitleGeneratorController {
	return &AiTitleGeneratorController{app: app}
}

func (c *AiTitleGeneratorController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := c.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, errorMessage, shared.NewLinks().Home(), 10)
	}

	if r.Method == http.MethodGet && data.Action == ACTION_ADD_TITLE {
		return c.onAddTitleModal(r)
	}

	if r.Method == http.MethodPost {
		switch data.Action {
		case ACTION_ADD_TITLE:
			return c.onAddTitle(r)
		case ACTION_GENERATE_TITLES:
			return c.onGenerateTitles(r)
		case ACTION_APPROVE_TITLE:
			return c.onApproveTitle(r)
		case ACTION_REJECT_TITLE:
			return c.onRejectTitle(r)
		case ACTION_DELETE_TITLE:
			return c.onDeleteTitle(r)
		}
	}

	return layouts.NewAdminLayout(c.app, r, layouts.Options{
		Title:   "AI Title Generator",
		Content: c.view(data),
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Sweetalert2_11(),
		},
		Scripts: []string{
			liveflux.Script().ToHTML(),
		},
		Styles: []string{
			ext.HxHideIndicatorCSS(),
		},
	}).ToHTML()
}

func (c *AiTitleGeneratorController) view(data pageData) *hb.Tag {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Dashboard",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Blog",
			URL:  links.Admin().Blog(),
		},
		{
			Name: "AI Tools",
			URL:  shared.NewLinks().AiTools(),
		},
		{
			Name: "Title Generator",
			URL:  shared.NewLinks().AiTitleGenerator(),
		},
	})

	settingsComponent := NewTitleGeneratorSettingsModal(c.app)
	settingsSSR := liveflux.SSR(settingsComponent, map[string]string{
		"return_url": shared.NewLinks().AiTitleGenerator(),
	})

	settingsButtonSpinner := hb.Span().
		Class("spinner spinner-border spinner-border-sm align-middle ms-2").
		Attr("role", "status").
		Style(`display: none;`).
		Child(hb.Span().Class("visually-hidden").Text("Loading"))

	settingsButton := hb.Button().
		Class("btn btn-outline-secondary btn-sm").
		Attr(liveflux.DataFluxAction, "open").
		Attr(liveflux.DataFluxTargetKind, settingsComponent.GetKind()).
		Attr(liveflux.DataFluxTargetID, settingsComponent.GetID()).
		Attr(liveflux.DataFluxIndicator, "this, .spinner").
		Child(hb.I().Class("bi bi-gear me-1")).
		Child(hb.Span().Text("Settings")).
		Child(settingsButtonSpinner)

	card := hb.Div().
		Class("card shadow-sm w-100 mb-5")
	card = card.
		Child(
			hb.Div().Class("card-body text-center p-4").
				Child(hb.Div().
					Class("d-flex justify-content-between align-items-center mb-3").
					Child(hb.Heading1().
						HTML("Title Generator").
						Class("h3 mb-0 fw-bold text-dark")).
					Child(settingsButton),
				).
				Child(
					hb.Paragraph().
						HTML("Create up to 10 fresh AI titles per run—existing titles are skipped automatically.").
						Class("text-muted mb-4"),
				).
				Child(
					func() hb.TagInterface {
						if !data.HasSystemPrompt {
							return hb.Div().
								Class("col-8 mx-auto mb-4").
								Child(hb.Div().
									Class("alert alert-info d-flex align-items-center gap-2 mb-0").
									Attr("role", "alert").
									Child(hb.I().Class("bi bi-info-circle-fill")).
									Child(hb.Span().Text("Set the Title Generator settings first, then you can generate new titles.")))
						}

						return hb.Div().
							Class("d-grid gap-3 col-8 mx-auto mb-4").
							Children([]hb.TagInterface{
								hb.Button().
									Class("btn btn-primary btn-lg fw-semibold").
									HTML(`Generate New Titles <span class="htmx-indicator spinner-border spinner-border-sm" role="status"></span>`).
									HxPost(shared.NewLinks().AiTitleGenerator(map[string]string{"action": ACTION_GENERATE_TITLES})).
									HxTarget("body").
									HxSwap("beforeend").
									Attr("hx-indicator", "this"),
								hb.Button().
									Class("btn btn-outline-primary btn-lg fw-semibold").
									HTML(`Add Custom Title <span class="htmx-indicator spinner-border spinner-border-sm" role="status"></span>`).
									HxGet(shared.NewLinks().AiTitleGenerator(map[string]string{"action": ACTION_ADD_TITLE})).
									HxTarget("body").
									HxSwap("beforeend").
									Attr("hx-indicator", "this"),
							})
					}(),
				).
				Child(
					hb.Div().Class("text-start").
						Child(tableGeneratedTitles(data)),
				).
				Child(settingsSSR),
		)

	return hb.Div().
		Class("container").
		Class("min-vh-100 py-4").
		Child(breadcrumbs).
		Child(card)
}

func (c *AiTitleGeneratorController) prepareData(r *http.Request) (data pageData, errorMessage string) {
	data = pageData{
		Request: r,
		Action:  req.GetStringTrimmed(r, "action"),
	}

	if c.app.GetCustomStore() == nil {
		return data, "Custom store is not initialized"
	}

	records, err := c.app.GetCustomStore().RecordList(customstore.RecordQuery().
		SetType(blogai.POST_RECORD_TYPE))
	if err != nil {
		return data, fmt.Sprintf("Failed to fetch titles: %s", err.Error())
	}

	recordPosts := []blogai.RecordPost{}
	for _, record := range records {
		recordPost, err := blogai.NewRecordPostFromCustomRecord(record)
		if err != nil {
			c.app.GetLogger().Warn("Failed to parse custom record into RecordPost: " + err.Error())
			continue
		}
		recordPosts = append(recordPosts, recordPost)
	}

	data.ExistingPostRecords = recordPosts

	// Determine if the system prompt setting is configured
	if c.app.GetSettingStore() != nil {
		value, err := c.app.GetSettingStore().Get(r.Context(), SETTING_KEY_BLOG_TOPIC, "")
		if err == nil && strings.TrimSpace(value) != "" {
			data.HasSystemPrompt = true
		}
	}

	return data, ""
}
