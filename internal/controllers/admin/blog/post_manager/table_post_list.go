package post_manager

import (
	"net/http"
	"project/internal/controllers/admin/blog/shared"
	"project/internal/links"
	"strings"

	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// tablePostList returns a table of posts
func tablePostList(data postManagerControllerData) hb.TagInterface {
	table := hb.Table().
		Class("table table-striped table-hover table-bordered align-middle").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().
						Child(sortableColumnLabel(data, "Post", "title")).
						Text(", ").
						Child(sortableColumnLabel(data, "Reference", "title")).
						Style(`cursor: pointer;`),
					hb.TH().
						Child(sortableColumnLabel(data, "Status", "status")).
						Style("width: 200px;cursor: pointer;"),
					hb.TH().
						Child(sortableColumnLabel(data, "Featured", "featured")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(sortableColumnLabel(data, "Published", "published_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(sortableColumnLabel(data, "Created", "created_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(sortableColumnLabel(data, "Modified", "updated_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.blogList, func(blog blogstore.Post, _ int) hb.TagInterface {
				blogLink := hb.Hyperlink().
					HTML(blog.Title()).
					Href(links.Website().BlogPost(blog.ID(), blog.Slug())).
					Target("_blank")

				status := hb.Span().
					Class("fw-semibold text-capitalize").
					ClassIf(blog.IsPublished(), "text-success").
					ClassIf(blog.IsTrashed(), "text-secondary").
					ClassIf(blog.IsDraft(), "text-primary").
					ClassIf(blog.IsUnpublished(), "text-danger").
					HTML(blog.Status())

				buttonAiContent := hb.Hyperlink().
					Class("btn btn-warning me-2").
					Child(hb.I().Class("bi bi-magic")).
					Title("AI Content Editor").
					Href(shared.NewLinks().AiPostContentUpdate(map[string]string{"post_id": blog.ID()})).
					Target("_blank")

				buttonEdit := hb.Hyperlink().
					Class("btn btn-primary me-2").
					Child(hb.I().Class("bi bi-pencil-square")).
					Title("Edit").
					Href(shared.NewLinks().PostUpdate(map[string]string{"post_id": blog.ID()})).
					Target("_blank")

				buttonDelete := hb.Hyperlink().
					Class("btn btn-danger").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete").
					HxGet(shared.NewLinks().PostDelete(map[string]string{"post_id": blog.ID()})).
					HxTarget("body").
					HxSwap("beforeend")

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().
						Child(hb.Div().Child(blogLink)).
						Child(hb.Div().
							Class("text-muted small").
							HTML("Ref: ").
							HTML(blog.ID())),
					hb.TD().
						Child(status),
					hb.TD().
						Child(hb.Span().Class("text-muted text-uppercase").HTML(blog.Featured())),
					hb.TD().
						Child(hb.Div().
							Class("small text-muted text-nowrap").
							HTML(blog.PublishedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(hb.Div().
							Class("small text-muted text-nowrap").
							HTML(blog.CreatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(hb.Div().
							Class("small text-muted text-nowrap").
							HTML(blog.UpdatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(buttonAiContent).
						Child(buttonEdit).
						Child(buttonDelete),
				})
			})),
		})

	// cfmt.Successln("Table: ", table)

	return hb.Wrap().Children([]hb.TagInterface{
		tableFilter(data),
		table,
		tablePagination(data, int(data.blogCount), data.pageInt, data.perPage),
	})
}

func sortableColumnLabel(data postManagerControllerData, tableLabel string, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == "asc", "desc").Else("asc")

	if !isSelected {
		direction = "asc"
	}

	link := shared.NewLinks().PostManager(map[string]string{
		"page":        "0",
		"by":          columnName,
		"sort":        direction,
		"date_from":   data.dateFrom,
		"date_to":     data.dateTo,
		"status":      data.status,
		"search":      data.search,
		"customer_id": data.customerID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(sortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func sortingIndicator(columnName string, sortByColumnName string, sortOrder string) hb.TagInterface {
	isSelected := strings.EqualFold(sortByColumnName, columnName)

	direction := lo.If(isSelected && sortOrder == "asc", "up").
		ElseIf(isSelected && sortOrder == "desc", "down").
		Else("none")

	sortingIndicator := hb.Span().
		Class("sorting").
		HTMLIf(direction == "up", "&#8595;").
		HTMLIf(direction == "down", "&#8593;").
		HTMLIf(direction != "down" && direction != "up", "")

	return sortingIndicator
}

func tableFilter(data postManagerControllerData) hb.TagInterface {
	statusList := []map[string]string{
		{"id": "", "name": "All Statuses"},
		{"id": blogstore.POST_STATUS_DRAFT, "name": "Draft"},
		{"id": blogstore.POST_STATUS_UNPUBLISHED, "name": "Unpublished"},
		{"id": blogstore.POST_STATUS_PUBLISHED, "name": "Published"},
		{"id": blogstore.POST_STATUS_TRASH, "name": "Deleted"},
	}

	searchButton := hb.Button().
		Type("submit").
		Child(hb.I().Class("bi bi-search")).
		Class("btn btn-primary w-100 h-100")

	period := hb.Div().Class("form-group").Children([]hb.TagInterface{
		hb.Label().
			HTML("Period").
			Style("margin-bottom: 0px;"),
		hb.Div().Class("input-group").Children([]hb.TagInterface{
			hb.Input().
				Type(hb.TYPE_DATE).
				Name("date_from").
				Value(data.dateFrom).
				OnChange("FORM_TRANSACTIONS.submit()").
				Class("form-control"),
			hb.Span().
				HTML(" : ").
				Class("input-group-text"),
			hb.Input().
				Type(hb.TYPE_DATE).
				Name("date_to").
				Value(data.dateTo).
				OnChange("FORM_TRANSACTIONS.submit()").
				Class("form-control"),
		}),
	})

	search := hb.Div().Class("form-group").Children([]hb.TagInterface{
		hb.Label().
			HTML("Search").
			Style("margin-bottom: 0px;"),
		hb.Input().
			Type("search").
			Name("search").
			Value(data.search).
			Class("form-control").
			Placeholder("reference, title, content ..."),
	})

	status := hb.Div().Class("form-group").Children([]hb.TagInterface{
		hb.Label().
			HTML("Status").
			Style("margin-bottom: 0px;"),
		hb.Select().
			Name("status").
			Class("form-select").
			OnChange("FORM_TRANSACTIONS.submit()").
			Children(lo.Map(statusList, func(status map[string]string, index int) hb.TagInterface {
				return hb.Option().
					Value(status["id"]).
					HTML(status["name"]).
					AttrIf(data.status == status["id"], "selected", "selected")
			})),
	})

	form := hb.Form().
		ID("FORM_TRANSACTIONS").
		Style("display: block").
		Method(http.MethodGet).
		Children([]hb.TagInterface{
			hb.Div().Class("row").Children([]hb.TagInterface{
				hb.Div().Class("col-md-2").Children([]hb.TagInterface{
					search,
				}),
				hb.Div().Class("col-md-4").Children([]hb.TagInterface{
					period,
				}),
				hb.Div().Class("col-md-2").Children([]hb.TagInterface{
					status,
				}),
				hb.Div().Class("col-md-1").Children([]hb.TagInterface{
					searchButton,
				}),
			}),
		})

	return hb.Div().
		Class("card bg-light mb-3").
		Style("").
		Children([]hb.TagInterface{
			hb.Div().Class("card-body").Children([]hb.TagInterface{
				form,
			}),
		})
}

func tablePagination(data postManagerControllerData, count int, page int, perPage int) hb.TagInterface {
	url := shared.NewLinks().PostManager(map[string]string{
		"search":    data.search,
		"status":    data.status,
		"date_from": data.dateFrom,
		"date_to":   data.dateTo,
		"by":        data.sortBy,
		"order":     data.sortOrder,
	})

	url = lo.Ternary(strings.Contains(url, "?"), url+"&page=", url+"?page=") // page must be last

	pagination := bs.Pagination(bs.PaginationOptions{
		NumberItems:       count,
		CurrentPageNumber: page,
		PagesToShow:       5,
		PerPage:           perPage,
		URL:               url,
	})

	return hb.Div().
		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
		HTML(pagination)
}
