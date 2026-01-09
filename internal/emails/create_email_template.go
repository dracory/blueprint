package emails

import (
	"project/internal/registry"

	baseEmail "github.com/dracory/base/email"
	"github.com/samber/lo"
)

// CreateEmailTemplate creates an email template using the base email package
// This is a new function to avoid conflicts with the original blankEmailTemplate function
func CreateEmailTemplate(registry registry.RegistryInterface, title string, htmlContent string) string {
	// Create header links
	headerLinks := map[string]string{}

	appName := lo.IfF(registry != nil, func() string {
		if registry.GetConfig() == nil {
			return ""
		}
		return registry.GetConfig().GetAppName()
	}).Else("")

	// Use the base email template
	return baseEmail.DefaultTemplate(baseEmail.TemplateOptions{
		Title:       title,
		Content:     htmlContent,
		AppName:     appName,
		HeaderLinks: headerLinks,
	})
}
