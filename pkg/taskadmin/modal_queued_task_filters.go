package taskadmin

import (
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func (p *pageQueueManager) modalQueuedTaskFilters(data pageQueueManagerData) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.NewHeading5().
		Text("Queued Task Filters").
		Style(`margin:0px;padding:0px;`)

	buttonModalClose := hb.NewButton().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.NewButton().
		Child(hb.NewI().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonOk := hb.NewButton().
		Child(hb.NewI().Class("bi bi-check me-2")).
		HTML("Apply").
		Class("btn btn-primary float-end").
		OnClick(`FormFilters.submit();` + modalCloseScript)

	statusList := []map[string]string{
		{
			"id":   "all",
			"name": "All",
		},
		{
			"id":   "pending",
			"name": "Pending",
		},
		{
			"id":   "running",
			"name": "Running",
		},
		{
			"id":   "completed",
			"name": "Completed",
		},
		{
			"id":   "failed",
			"name": "Failed",
		},
	}

	groupSearch := bs.FormGroup().
		Child(
			hb.NewDiv().
				HTML("Search:").
				Style(`font-size:18px;color:black;font-weight:bold;`),
		).
		Child(
			hb.NewInput().
				Type("search").
				Name("search").
				Class("form-control"),
		)

	groupStatus := bs.FormGroup().
		Child(
			hb.NewDiv().
				HTML("Status:").
				Style(`font-size:18px;color:black;font-weight:bold;`),
		).
		Child(
			hb.NewSelect().
				Name("status").
				Class("form-select").
				Children(lo.Map(statusList, func(status map[string]string, index int) hb.TagInterface {
					return hb.NewOption().
						Value(status["id"]).
						HTML(status["name"]).
						AttrIf(status["id"] == data.status, "selected", "selected")
				})),
		)

	groupPeriod := bs.FormGroup().
		Child(
			hb.NewDiv().
				HTML("Period:").
				Style(`font-size:18px;color:black;font-weight:bold;`),
		).
		Child(
			hb.NewInput().
				Type(hb.TYPE_DATE).
				Name("date_from").
				Value(data.dateFrom).
				Class("form-control"),
		).Child(
		hb.NewSpan().
			HTML(" : "),
	).Child(
		hb.NewInput().
			Type(hb.TYPE_DATE).
			Name("date_to").
			Value(data.dateTo).
			Class("form-control"),
	)

	modal := bs.Modal().
		ID("ModalMessage").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						title,
						buttonModalClose,
					}),

					bs.ModalBody().
						Child(hb.NewForm().
							ID("FormFilters").
							Child(groupSearch).
							Child(groupStatus).
							Child(groupPeriod)),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonOk),
				}),
			}),
		})

	backdrop := hb.NewDiv().
		ID("ModalBackdrop").
		Class("modal-backdrop fade show").
		Style("display:block;")

	return hb.NewWrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})

}
