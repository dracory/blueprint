package flash

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONTROLLER ==============================================================

type flashController struct {
	app types.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewFlashController(app types.AppInterface) *flashController {
	return &flashController{app: app}
}

// == PUBLIC METHODS ==========================================================

func (controller flashController) Handler(w http.ResponseWriter, r *http.Request) string {
	authUser := helpers.GetAuthUser(r)

	title := "System Message"
	html := controller.pageHTML(r)

	if authUser != nil && authUser.IsRegistrationCompleted() {
		return layouts.NewUserLayout(controller.app, r, layouts.Options{
			Title:      title,
			Content:    html,
			ScriptURLs: []string{},
			Styles:     []string{`.Center > div{padding:0px !important;margin:0px !important;}`},
			StyleURLs: []string{
				cdn.BootstrapIconsCss_1_13_1(),
			},
		}).ToHTML()
	}

	if controller.app.GetCmsStore() != nil && controller.app.GetConfig() != nil && controller.app.GetConfig().GetCMSTemplateID() != "" {
		return layouts.NewCmsLayout(
			controller.app,
			layouts.Options{
				Request: r,
				Title:   title,
				Content: html,
				Styles:  []string{`.Center > div{padding:0px !important;margin:0px !important;}`},
				StyleURLs: []string{
					cdn.BootstrapIconsCss_1_13_1(),
				},
			},
		).ToHTML()
	} else {
		return layouts.NewUserLayout(controller.app, r, layouts.Options{
			Title:   title,
			Content: html,
			StyleURLs: []string{
				cdn.BootstrapIconsCss_1_13_1(),
			},
		}).ToHTML()
	}
}

func (c flashController) pageHTML(r *http.Request) hb.TagInterface {
	messageID := req.GetStringTrimmed(r, "message_id")
	msgData, err := c.app.GetCacheStore().GetJSON(messageID+"_flash_message", "")

	msgType := "error"
	message := "The message is no longer available"
	url := links.Website().Home()
	time := "5"

	if err != nil {
		message = "The message is no longer available"
	}

	if msgData == "" {
		message = "The message is no longer available"
	}

	if msgData != "" {
		msgDataAny := msgData.(map[string]interface{})
		msgType = cast.ToString(msgDataAny["type"])
		message = cast.ToString(msgDataAny["message"])
		url = cast.ToString(msgDataAny["url"])
		time = cast.ToString(msgDataAny["time"])
	}

	alertIcon := hb.I().
		Class("me-2").
		Class("bi").
		Class(lo.If(msgType == "error", "bi-exclamation-octagon-fill").
			ElseIf(msgType == "success", "bi-check-circle-fill").
			ElseIf(msgType == "warning", "bi-exclamation-triangle-fill").
			Else("bi-info-circle-fill"))

	alert := hb.Div().
		Class("alert").
		Class(lo.
			If(msgType == "error", "alert-danger").
			ElseIf(msgType == "success", "alert-success").
			ElseIf(msgType == "warning", "alert-warning").
			Else("alert-info")).
		Child(alertIcon).
		HTML(message)

	linkRedirect := hb.Hyperlink().Href(url).HTML("Click here to continue")
	scriptRedirect := hb.Script("setTimeout(()=>{location.href=\"" + url + "\"}, " + time + "*1000)")

	container := hb.Div().
		Class("container").
		Style("padding:0px 0px 20px 0px;text-align:left;").
		Child(alert).
		ChildIf(url != "", hb.Div().
			Child(linkRedirect).
			Style("padding:20px 0px 20px 0px;")).
		ChildIf(url != "" && time != "", scriptRedirect)

	return hb.Section().
		Child(container).
		Style("padding: 80px 0px 40px 0px;")
}
