package register

import (
	"context"
	"project/internal/links"

	"github.com/dracory/bs"
	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func (controller *registerController) formRegister(ctx context.Context, data registerControllerData) hb.TagInterface {
	required := hb.Sup().
		Text("required").
		Style("margin-left:5px;color:lightcoral;")

	buttonSave := bs.Button().
		Class("btn-primary mb-0 w-100 py-3 fs-5").
		Attr("type", "button").
		Child(hb.I().Class("bi bi-check-circle me-2")).
		Text("Save changes").
		HxInclude("#FormRegister").
		HxTarget("#CardRegister").
		HxTrigger("click").
		HxSwap("outerHTML").
		HxPost(links.Auth().Register())

	firstNameGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("First name").
				Child(required),
			bs.FormInput().
				Name(controller.formFirstName).
				Value(data.firstName),
		})

	lastNameGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Last name").
				Child(required),
			bs.FormInput().
				Name(controller.formLastName).
				Value(data.lastName),
		})

	businessNameGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Company / buiness name"),
			bs.FormInput().
				Name("business_name").
				Value(data.buinessName),
		})

	phoneGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Phone"),
			bs.FormInput().
				Name("phone").
				Value(data.phone),
		})

	emailGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Email").
				Child(required),
			bs.FormInput().
				Name("email").
				Value(data.email).
				Attr("readonly", "readonly").
				Style("background-color:#F8F8F8;"),
		})

	selectCountries := bs.FormSelect().
		ID("SelectCountries").
		Name(controller.formCountry).
		Child(bs.FormSelectOption("", "")).
		Children(lo.Map(data.countryList, func(country geostore.Country, _ int) hb.TagInterface {
			return bs.FormSelectOption(country.IsoCode2(), country.Name()).
				AttrIf(data.country == country.IsoCode2(), "selected", "selected")
		})).
		Hx("post", links.Auth().Register(map[string]string{
			"action": "on-country-selected-timezone-options",
		})).
		Hx("target", "#SelectTimezones").
		Hx("swap", "outerHTML")

	countryGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Country").
				Child(required),
			selectCountries,
		})

	timezoneGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Timezone").
				Child(required),
			controller.selectTimezoneByCountry(ctx, data.country, data.timezone),
		})

	formProfile := hb.Div().
		ID("FormRegister").
		Child(
			bs.Row().
				Class("g-3").
				Children([]hb.TagInterface{
					bs.Column(12).
						Child(emailGroup),
					bs.Column(6).
						Child(firstNameGroup),
					bs.Column(6).
						Child(lastNameGroup),
					bs.Column(6).
						Child(businessNameGroup),
					bs.Column(6).
						Child(phoneGroup),
					bs.Column(6).
						Child(countryGroup),
					bs.Column(6).
						Child(timezoneGroup),
				}),
		).
		Child(
			bs.Row().Class("mt-4").Children([]hb.TagInterface{
				bs.Column(12).Class("d-sm-flex justify-content-end").
					Children([]hb.TagInterface{
						buttonSave,
					}),
			}),
		)

	return hb.Div().ID("CardRegister").
		Class("card bg-white border rounded-3").
		Style("text-align:left;").
		Children([]hb.TagInterface{
			hb.Div().Class("card-header  bg-transparent").Children([]hb.TagInterface{
				hb.Heading3().
					Text("Your Details").
					Style("text-align:left;font-size:12px;color:#333;margin:0px;"),
			}),
			hb.Div().Class("card-body").Children([]hb.TagInterface{
				formProfile,
			}),
		}).
		ChildIf(data.formErrorMessage != "", hb.Swal(hb.SwalOptions{
			Icon: "error",
			// Title:             "Oops...",
			Text:              data.formErrorMessage,
			ShowCancelButton:  false,
			ShowConfirmButton: false,
			Timer:             5000,
			TimerProgressBar:  true,
			Position:          "top-end",
		})).
		ChildIf(data.formSuccessMessage != "", hb.Swal(hb.SwalOptions{
			Icon:              "success",
			Title:             "Saved",
			Text:              data.formSuccessMessage,
			ShowCancelButton:  false,
			ShowConfirmButton: false,
			ConfirmCallback:   "window.location.href = window.location.href",
			Timer:             5000,
			TimerProgressBar:  true,
			Position:          "top-end",
		})).
		ChildIf(data.formRedirectURL != "", hb.Script(`window.location.href = '`+data.formRedirectURL+`'`))
}
