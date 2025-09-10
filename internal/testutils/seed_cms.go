package testutils

import (
	"context"
	"errors"

	"github.com/dracory/cmsstore"
)

func SeedTemplate(cmsStore cmsstore.StoreInterface, siteID, templateID string) (err error) {
	if cmsStore == nil {
		return errors.New("cmsstore.seed template: cmsstore is nil")
	}

	templateContent := `
	<html>
	    <head>
			<title>[[PageTitle]]</title>
		</head>
		<body>
			[[PageContent]]
		</body>
	</html>
	`

	template := cmsstore.NewTemplate().
		SetID(templateID).
		SetSiteID(siteID).
		SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE).
		SetName(templateID).
		SetContent(templateContent)

	return cmsStore.TemplateCreate(context.Background(), template)
}
