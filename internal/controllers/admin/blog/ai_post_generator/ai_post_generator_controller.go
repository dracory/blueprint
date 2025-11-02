package aipostgenerator

import (
	"fmt"
	"log/slog"
	"net/http"
	"project/internal/controllers/admin/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"
	"project/pkg/blogai"

	"github.com/dracory/cdn"
	"github.com/dracory/customstore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

const ACTION_GENERATE_POST = "generate_post"

type AiPostGeneratorController struct {
	app types.AppInterface
}

type pageData struct {
	Request             *http.Request
	Action              string
	ApprovedBlogAiPosts []blogai.RecordPost
}

func NewAiPostGeneratorController(app types.AppInterface) *AiPostGeneratorController {
	return &AiPostGeneratorController{app: app}
}

func (c *AiPostGeneratorController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := c.prepareData(r)

	if errorMessage != "" {
		if cache := c.app.GetCacheStore(); cache != nil {
			return helpers.ToFlashError(cache, w, r, errorMessage, shared.NewLinks().Home(), 10)
		}
		return shared.ErrorPopup(errorMessage).ToHTML()
	}

	if r.Method == http.MethodPost && data.Action == ACTION_GENERATE_POST {
		return c.onGeneratePost(r)
	}

	return layouts.NewAdminLayout(c.app, r, layouts.Options{
		Title:   "Post Generator",
		AppName: c.app.GetConfig().GetAppName(),
		Content: c.view(data),
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Sweetalert2_11(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (c *AiPostGeneratorController) view(data pageData) *hb.Tag {
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
			Name: "Post Generator",
			URL:  shared.NewLinks().AiPostGenerator(),
		},
	})

	card := hb.Div().
		Class("card shadow-sm w-100 mb-5")
	card = card.Child(
		hb.Div().Class("card-body text-center p-4").
			Child(
				hb.Heading1().
					HTML("Post Generator").
					Class("h3 mb-3 fw-bold"),
			).
			Child(
				hb.Paragraph().
					HTML("Select an approved title below to generate a blog post.").
					Class("text-muted mb-4"),
			).
			Child(
				hb.Div().Class("text-start").
					Child(c.tableApprovedTitles(data)),
			),
	)

	return hb.Div().
		Class("container").
		Class("min-vh-100 py-4").
		Child(breadcrumbs).
		Child(card)
}

func (c *AiPostGeneratorController) prepareData(r *http.Request) (data pageData, errorMessage string) {
	data.Request = r
	data.Action = req.GetStringTrimmed(r, "action")

	customStore := c.app.GetCustomStore()
	if customStore == nil {
		return data, "custom store not configured"
	}

	approvedTitleRecords, err := customStore.RecordList(customstore.NewRecordQuery().
		SetType(blogai.POST_RECORD_TYPE).
		AddPayloadSearch(`"status":"approved"`).
		AddPayloadSearch(`"status":"draft"`))

	if err != nil {
		return data, fmt.Sprintf("Error fetching approved titles: %s", err.Error())
	}

	approvedBlogAiPosts := []blogai.RecordPost{}
	for _, record := range approvedTitleRecords {
		recordPost, err := blogai.NewRecordPostFromCustomRecord(record)
		if err != nil {
			if logger := c.app.GetLogger(); logger != nil {
				logger.Warn("Failed to parse custom record into RecordPost", slog.String("error", err.Error()))
			}
			continue
		}

		if recordPost.Status == blogai.POST_STATUS_APPROVED || recordPost.Status == blogai.POST_STATUS_DRAFT {
			approvedBlogAiPosts = append(approvedBlogAiPosts, recordPost)
		}
	}

	data.ApprovedBlogAiPosts = approvedBlogAiPosts

	return data, ""
}
