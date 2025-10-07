package emails

import baseEmail "github.com/dracory/base/email"

// CreateEmailTemplate creates an email template using the base email package
// This is a new function to avoid conflicts with the original blankEmailTemplate function
func CreateEmailTemplate(title string, htmlContent string) string {
	// Create header links
	headerLinks := map[string]string{}

	// Use the base email template
	return baseEmail.DefaultTemplate(baseEmail.TemplateOptions{
		Title:   title,
		Content: htmlContent,
		AppName: func() string {
			if cfg != nil {
				return cfg.GetAppName()
			}
			return ""
		}(),
		HeaderLinks: headerLinks,
	})
}
