package admin

import (
	"project/internal/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	fileManager := rtr.NewRoute().
		SetName("Admin > File Manager").
		SetPath(links.ADMIN_FILE_MANAGER).
		SetHTMLHandler(NewFileManagerController().Handler)

	return []rtr.RouteInterface{
		fileManager,
	}
}
