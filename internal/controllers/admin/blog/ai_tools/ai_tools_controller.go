package ai_tools

import (
	"net/http"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

// == CONTROLLER ==============================================================

type aiToolsController struct {
	app types.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewAiToolsController(app types.AppInterface) *aiToolsController {
	return &aiToolsController{app: app}
}

func (c *aiToolsController) Handler(w http.ResponseWriter, r *http.Request) string {
	if r.Method == http.MethodPost && r.FormValue("action") == "testai" {
		w.Header().Set("Content-Type", "application/json")
		model, err := shared.LlmEngine(c.app)
		if err != nil {
			w.Write([]byte(hb.Swal(hb.SwalOptions{
				Title:            "Error",
				Text:             err.Error(),
				Icon:             "error",
				Timer:            15000,
				TimerProgressBar: true,
			}).ToHTML()))
			return ""
		}
		if model == nil {
			w.Write([]byte(hb.Swal(hb.SwalOptions{
				Title:            "Error",
				Text:             "model is nil",
				Icon:             "error",
				Timer:            15000,
				TimerProgressBar: true,
			}).ToHTML()))
			return ""
		}

		response, err := model.GenerateText("You are a helpful assistant.", "Tell me shortly about blogs.")
		if err != nil {
			w.Write([]byte(hb.Swal(hb.SwalOptions{
				Title:            "Error",
				Text:             err.Error(),
				Icon:             "error",
				Timer:            15000,
				TimerProgressBar: true,
			}).ToHTML()))
			return ""
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hb.Swal(hb.SwalOptions{
			Title:            "Success",
			Text:             response,
			Icon:             "success",
			Timer:            15000,
			TimerProgressBar: true,
		}).ToHTML()))
		return ""
	}

	return layouts.NewAdminLayout(c.app, r, layouts.Options{
		Title:      "BlogAI",
		AppName:    c.app.GetConfig().GetAppName(),
		Content:    c.view(),
		ScriptURLs: []string{cdn.Sweetalert2_11()},
		Styles:     []string{},
	}).ToHTML()
}

func (c *aiToolsController) view() *hb.Tag {
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
	})

	card := hb.Div().
		Class("card shadow-sm w-100 mb-5")
	card = card.Child(
		hb.Div().Class("card-body text-center p-4").
			Child(
				hb.Heading1().
					HTML("BlogAI").
					Class("h3 mb-3 fw-bold"),
			).
			Child(
				hb.Paragraph().
					HTML("Welcome to your AI-powered blog tools.").
					Class("text-muted mb-4"),
			).
			Child(
				hb.Div().Class("d-grid gap-3 col-8 mx-auto").
					Child(
						hb.Hyperlink().
							HTML("📝 Post Generator").
							Href(shared.NewLinks().AiPostGenerator()).
							Class("btn btn-primary btn-lg fw-semibold"),
					).
					Child(
						hb.Hyperlink().
							HTML("💡 Title Generator").
							Href(shared.NewLinks().AiTitleGenerator()).
							Class("btn btn-success btn-lg fw-semibold"),
					).
					Child(
						hb.Hyperlink().
							Class("btn btn-warning btn-lg fw-semibold").
							HTML("Test AI is working").
							Href(shared.NewLinks().AiTest()),
					).
					Child(
						hb.Hyperlink().
							Class("btn btn-outline-secondary btn-lg fw-semibold d-inline-flex align-items-center justify-content-center").
							Child(hb.I().Class("bi bi-arrow-left-circle me-2")).
							HTML("Back to Blog Home").
							Href(shared.NewLinks().Home()),
					),
			),
	)

	return hb.Div().
		Class("container").
		Class("min-vh-100 py-4").
		Child(breadcrumbs).
		Child(card)
}
