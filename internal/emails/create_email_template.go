package emails

import (
	"project/internal/types"

	baseEmail "github.com/dracory/base/email"
	"github.com/samber/lo"
)

// CreateEmailTemplate creates an email template using the base email package
// This is a new function to avoid conflicts with the original blankEmailTemplate function
func CreateEmailTemplate(app types.AppInterface, title string, htmlContent string) string {
	// Create header links
	headerLinks := map[string]string{}

	appName := lo.IfF(app != nil, func() string {
		if app.GetConfig() == nil {
			return ""
		}
		return app.GetConfig().GetAppName()
	}).Else("")

	// Use the base email template
	return baseEmail.DefaultTemplate(baseEmail.TemplateOptions{
		Title:       title,
		Content:     htmlContent,
		AppName:     appName,
		HeaderLinks: headerLinks,
	})
}
