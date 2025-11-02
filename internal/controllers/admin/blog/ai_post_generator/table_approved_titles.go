package aipostgenerator

import (
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func (c *AiPostGeneratorController) tableApprovedTitles(data pageData) *hb.Tag {
	tableRows := lo.Map(data.ApprovedBlogAiPosts, func(recordPost blogai.RecordPost, _ int) hb.TagInterface {
		var actionButtons []hb.TagInterface

		if recordPost.Status == blogai.POST_STATUS_APPROVED {
			buttonGenerate := hb.Button().
				Class("btn btn-primary btn-sm").
				HTML(`Generate Post <span class="htmx-indicator spinner-border spinner-border-sm" role="status"></span>`).
				HxPost(shared.NewLinks().AiPostGenerator(map[string]string{
					"action":         ACTION_GENERATE_POST,
					"record_post_id": recordPost.ID,
				})).
				HxTarget("body").
				HxSwap("beforeend").
				HxIndicator("this")

			actionButtons = append(actionButtons, buttonGenerate)
		}

		if recordPost.Status == blogai.POST_STATUS_DRAFT {
			buttonViewDraft := hb.A().
				Class("btn btn-info btn-sm").
				HTML(`View Draft`).
				Href(shared.NewLinks().AiPostEditor(map[string]string{
					"id": recordPost.ID,
				}))

			actionButtons = append(actionButtons, buttonViewDraft)
		}

		status := recordPost.Status
		if status == "" {
			status = "N/A"
		}
		statusBadge := hb.Span().
			Class("badge rounded-pill " + getStatusBadgeClass(status) + " px-3").
			Text(status)

		return hb.TR().
			Child(hb.TD().
				Text(recordPost.Title)).
			Child(hb.TD().
				Child(statusBadge)).
			Child(hb.TD().
				Children(actionButtons))
	})

	tableHead := hb.Thead().
		Child(hb.TR().
			Child(hb.TH().
				Class("fw-semibold text-uppercase small").
				Text("Title")).
			Child(hb.TH().
				Class("fw-semibold text-uppercase small").
				Text("Status")).
			Child(hb.TH().
				Class("fw-semibold text-uppercase small").
				Text("Actions")))

	tableTitles := hb.Table().
		Class("table table-striped table-hover table-bordered align-middle").
		Child(tableHead).
		Child(hb.Tbody().
			Children(tableRows))

	if len(data.ApprovedBlogAiPosts) == 0 {
		return hb.Wrap().Child(
			hb.P().
				Class("text-muted").
				Text("No approved titles found. Please generate and approve some titles first."),
		)
	}

	return tableTitles
}

func getStatusBadgeClass(status string) string {
	switch status {
	case blogai.POST_STATUS_APPROVED:
		return "bg-success"
	case blogai.POST_STATUS_DRAFT:
		return "bg-info"
	default:
		return "bg-secondary"
	}
}
