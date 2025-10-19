package admin

import (
	"context"
	"log/slog"
	"strings"

	"project/internal/ext"
	"project/internal/links"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/userstore"
	"github.com/samber/lo"
)

func (controller *userManagerController) tableUsers(data userManagerControllerData) hb.TagInterface {
	table := hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().
						Child(controller.sortableColumnLabel(data, "First Name", "first_name")).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Last Name", "last_name")).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Reference", "id")).
						Style(`cursor: pointer;`),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Status", "status")).
						Style("width: 200px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "E-mail", "email")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Created", "created_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Modified", "updated_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.userList, func(user userstore.UserInterface, _ int) hb.TagInterface {
				var err error
				firstName := user.FirstName()
				lastName := user.LastName()
				email := user.Email()

				if controller.app.GetConfig().GetVaultStoreUsed() {
					firstName, lastName, email, _, _, err = ext.UserUntokenize(context.Background(), controller.app, controller.app.GetConfig().GetVaultStoreKey(), user)

					if err != nil {
						controller.app.GetLogger().Error("At userManagerController > tableUsers", slog.String("error", err.Error()))
						firstName = "n/a"
						lastName = "n/a"
						email = "n/a"
					}
				}

				userLink := hb.Hyperlink().
					Text(firstName).
					Text(` `).
					Text(lastName).
					Href(links.Admin().UsersUserUpdate(map[string]string{"user_id": user.ID()}))

				status := hb.Span().
					Style(`font-weight: bold;`).
					StyleIf(user.IsActive(), `color:green;`).
					StyleIf(user.IsSoftDeleted(), `color:silver;`).
					StyleIf(user.IsUnverified(), `color:blue;`).
					StyleIf(user.IsInactive(), `color:red;`).
					HTML(user.Status())

				buttonEdit := hb.Hyperlink().
					Class("btn btn-primary me-2").
					Child(hb.I().Class("bi bi-pencil-square")).
					Title("Edit").
					Href(links.Admin().UsersUserUpdate(map[string]string{"user_id": user.ID()})).
					Target("_blank")

				buttonDelete := hb.Hyperlink().
					Class("btn btn-danger").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete").
					HxGet(links.Admin().UsersUserDelete(map[string]string{"user_id": user.ID()})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonImpersonate := hb.Hyperlink().
					Class("btn btn-warning me-2").
					Child(hb.I().Class("bi bi-shuffle")).
					Title("Impersonate").
					Href(links.Admin().UsersUserImpersonate(map[string]string{"user_id": user.ID()}))

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().
						Child(hb.Div().Child(userLink)).
						Child(hb.Div().
							Style("font-size: 11px;").
							HTML("Ref: ").
							HTML(user.ID())),
					hb.TD().
						Child(status),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(email)),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(user.CreatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(user.UpdatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(buttonEdit).
						Child(buttonImpersonate).
						Child(buttonDelete),
				})
			})),
		})

	return hb.Wrap().Children([]hb.TagInterface{
		controller.tableFilter(data),
		table,
		controller.tablePagination(data, int(data.userCount), data.pageInt, data.perPage),
	})
}

func (controller *userManagerController) sortableColumnLabel(data userManagerControllerData, tableLabel string, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == "asc", "desc").Else("asc")

	if !isSelected {
		direction = "asc"
	}

	link := links.Admin().UsersUserManager(map[string]string{
		"page":      "0",
		"by":        columnName,
		"sort":      direction,
		"date_from": data.formCreatedFrom,
		"date_to":   data.formCreatedTo,
		"status":    data.formStatus,
		"user_id":   data.formUserID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(controller.sortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func (controller *userManagerController) sortingIndicator(columnName string, sortByColumnName string, sortOrder string) hb.TagInterface {
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

func (controller *userManagerController) tableFilter(data userManagerControllerData) hb.TagInterface {
	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(links.Admin().UsersUserManager(map[string]string{
			"action":       ActionModalUserFilterShow,
			"first_name":   data.formFirstName,
			"last_name":    data.formLastName,
			"email":        data.formEmail,
			"status":       data.formStatus,
			"user_id":      data.formUserID,
			"created_from": data.formCreatedFrom,
			"created_to":   data.formCreatedTo,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing users").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formEmail != "" {
		description = append(description, hb.Span().Text("and email: "+data.formEmail).ToHTML())
	}

	if data.formUserID != "" {
		description = append(description, hb.Span().Text("and ID: "+data.formUserID).ToHTML())
	}

	if data.formFirstName != "" {
		description = append(description, hb.Span().Text("and first name: "+data.formFirstName).ToHTML())
	}

	if data.formLastName != "" {
		description = append(description, hb.Span().Text("and last name: "+data.formLastName).ToHTML())
	}

	if data.formCreatedFrom != "" && data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created between: "+data.formCreatedFrom+" and "+data.formCreatedTo).ToHTML())
	} else if data.formCreatedFrom != "" {
		description = append(description, hb.Span().Text("and created after: "+data.formCreatedFrom).ToHTML())
	} else if data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created before: "+data.formCreatedTo).ToHTML())
	}

	return hb.Div().
		Class("card bg-light mb-3").
		Style("").
		Children([]hb.TagInterface{
			hb.Div().Class("card-body").
				Child(buttonFilter).
				Child(hb.Span().
					HTML(strings.Join(description, " "))),
		})
}

func (controller *userManagerController) tablePagination(data userManagerControllerData, count int, page int, perPage int) hb.TagInterface {
	url := links.Admin().UsersUserManager(map[string]string{
		"status":       data.formStatus,
		"first_name":   data.formFirstName,
		"last_name":    data.formLastName,
		"email":        data.formEmail,
		"created_from": data.formCreatedFrom,
		"created_to":   data.formCreatedTo,
		"by":           data.sortBy,
		"order":        data.sortOrder,
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
