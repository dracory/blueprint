package aititlegenerator

import (
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func tableGeneratedTitles(data pageData) *hb.Tag {
	if len(data.ExistingPostRecords) == 0 {
		return hb.P().Class("text-muted").HTML("No titles generated yet.")
	}

	activeRecords := []blogai.RecordPost{}
	publishedRecords := []blogai.RecordPost{}

	for _, record := range data.ExistingPostRecords {
		if record.Status == blogai.POST_STATUS_PUBLISHED {
			publishedRecords = append(publishedRecords, record)
			continue
		}

		activeRecords = append(activeRecords, record)
	}

	container := hb.Div().Class("d-flex flex-column gap-4")

	if len(activeRecords) > 0 {
		container = container.Child(renderTitleRecordsSection(
			"Generated Titles",
			"Manage new and in-progress titles below.",
			activeRecords,
		))
	}

	if len(publishedRecords) > 0 {
		container = container.Child(renderTitleRecordsSection(
			"Published Titles (Reference)",
			"Previously published titles are listed here for reference and historical tracking.",
			publishedRecords,
		))
	}

	return container
}

func renderTitleRecordsSection(title string, description string, records []blogai.RecordPost) *hb.Tag {
	section := hb.Div().Class("d-flex flex-column gap-2")

	header := hb.Div().Class("d-flex flex-column gap-1")
	header = header.Child(
		hb.Heading3().
			Class("h5 fw-semibold mb-0").
			Text(title),
	)

	if description != "" {
		header = header.Child(
			hb.Paragraph().
				Class("text-muted mb-0").
				Text(description),
		)
	}

	section = section.
		Child(header).
		Child(buildTitleTable(records))

	return section
}

func buildTitleTable(records []blogai.RecordPost) *hb.Tag {
	tableRows := lo.Map(records, func(recordPost blogai.RecordPost, _ int) hb.TagInterface {
		var actionButtons []hb.TagInterface

		buttonApprove := hb.
			Button().
			Class("btn btn-success btn-sm").
			HTML(`Approve <span class="htmx-indicator spinner-border spinner-border-sm" role="status"></span>`).
			HxPost(shared.NewLinks().AiTitleGenerator(map[string]string{
				"action":         ACTION_APPROVE_TITLE,
				"record_post_id": recordPost.ID,
			})).
			HxTarget("body").
			HxSwap("beforeend").
			Attr("hx-indicator", "this")

		buttonReject := hb.
			Button().
			Class("btn btn-warning btn-sm ms-2").
			HTML(`Reject <span class="htmx-indicator spinner-border spinner-border-sm" role="status"></span>`).
			HxPost(shared.NewLinks().AiTitleGenerator(map[string]string{
				"action":         ACTION_REJECT_TITLE,
				"record_post_id": recordPost.ID,
			})).
			HxTarget("body").
			HxSwap("beforeend").
			Attr("hx-indicator", "this")

		buttonGeneratePost := hb.
			A().
			Class("btn btn-primary btn-sm").
			HTML("Generate Post").
			Attr("href", shared.NewLinks().AiPostGenerator(map[string]string{
				"action":         ACTION_GENERATE_POST,
				"record_post_id": recordPost.ID,
			}))

		makeDeleteButton := func(withMargin bool) hb.TagInterface {
			className := "btn btn-outline-danger btn-sm"
			if withMargin {
				className += " ms-2"
			}

			return hb.
				Button().
				Class(className).
				HTML(`Delete <span class="htmx-indicator spinner-border spinner-border-sm" role="status"></span>`).
				HxPost(shared.NewLinks().AiTitleGenerator(map[string]string{
					"action":         ACTION_DELETE_TITLE,
					"record_post_id": recordPost.ID,
				})).
				HxTarget("body").
				HxSwap("beforeend").
				Attr("hx-indicator", "this")
		}

		switch recordPost.Status {
		case blogai.POST_STATUS_PENDING:
			actionButtons = append(actionButtons, buttonApprove, buttonReject, makeDeleteButton(true))
		case blogai.POST_STATUS_APPROVED:
			actionButtons = append(actionButtons, buttonGeneratePost, makeDeleteButton(true))
		default:
			actionButtons = append(actionButtons, makeDeleteButton(false))
		}

		status := recordPost.Status
		if status == "" {
			status = "N/A"
		}

		statusBadge := hb.Span().
			Class("badge rounded-pill " + getStatusBadgeClass(status) + " px-3").
			Text(status)

		return hb.TR().Children([]hb.TagInterface{
			hb.TD().Text(recordPost.Title),
			hb.TD().Child(statusBadge),
			hb.TD().Children(actionButtons),
		})
	})

	return hb.Table().
		Class("table table-striped table-hover align-middle text-body").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().Class("text-body fw-semibold").Text("Title"),
					hb.TH().Class("text-body fw-semibold").Text("Status"),
					hb.TH().Class("text-body fw-semibold").Text("Actions"),
				}),
			}),
			hb.Tbody().Children(tableRows),
		})
}

func getStatusBadgeClass(status string) string {
	switch status {
	case blogai.POST_STATUS_PENDING:
		return "bg-warning"
	case blogai.POST_STATUS_APPROVED:
		return "bg-success"
	case blogai.POST_STATUS_REJECTED:
		return "bg-danger"
	case blogai.POST_STATUS_DRAFT:
		return "bg-info"
	case blogai.POST_STATUS_PUBLISHED:
		return "bg-primary"
	default:
		return "bg-secondary"
	}
}
