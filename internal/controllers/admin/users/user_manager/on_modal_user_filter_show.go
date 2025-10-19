package admin

import (
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/userstore"
)

const ActionModalUserFilterShow = "modal_user_filter_show"

func (controller *userManagerController) onModalUserFilterShow(data userManagerControllerData) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.Heading5().
		Text("Filters").
		Style(`margin:0px;padding:0px;`)

	buttonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonOk := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Apply").
		Class("btn btn-primary float-end").
		OnClick(`FormFilters.submit();` + modalCloseScript)

	filterForm := form.NewForm(form.FormOptions{
		ID:     "FormFilters",
		Method: http.MethodGet,
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Status",
				Name:  "status",
				Type:  form.FORM_FIELD_TYPE_SELECT,
				Help:  `The status of the user.`,
				Value: data.formStatus,
				Options: []form.FieldOption{
					{
						Value: "",
						Key:   "",
					},
					{
						Value: "Active",
						Key:   userstore.USER_STATUS_ACTIVE,
					},
					{
						Value: "Inactive",
						Key:   userstore.USER_STATUS_INACTIVE,
					},
					{
						Value: "Unverified",
						Key:   userstore.USER_STATUS_UNVERIFIED,
					},
					{
						Value: "Deleted",
						Key:   userstore.USER_STATUS_DELETED,
					},
				},
			}),
			form.NewField(form.FieldOptions{
				Label: "First Name",
				Name:  "first_name",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formFirstName,
				Help:  `Filter by first name.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Last Name",
				Name:  "last_name",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formLastName,
				Help:  `Filter by last name.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Email",
				Name:  "email",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formEmail,
				Help:  `Filter by email.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created From",
				Name:  "created_from",
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedFrom,
				Help:  `Filter by creation date.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created To",
				Name:  "created_to",
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedTo,
				Help:  `Filter by creation date.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "User ID",
				Name:  "user_id",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formUserID,
				Help:  `Find user by reference number (ID).`,
			}),
		},
	}).Build()

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
						Child(filterForm),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonOk),
				}),
			}),
		})

	backdrop := hb.Div().
		ID("ModalBackdrop").
		Class("modal-backdrop fade show").
		Style("display:block;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})

}
