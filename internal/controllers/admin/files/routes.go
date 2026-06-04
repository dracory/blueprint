package admin

import (
	"errors"
	"project/internal/links"
	"project/internal/app"

	"github.com/dracory/rtr"
)

func Routes(app app.AppInterface) ([]rtr.RouteInterface, error) {
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}
	fileManager := rtr.NewRoute().
		SetName("Admin > File Manager").
		SetPath(links.ADMIN_FILE_MANAGER).
		SetHandler(NewFileManagerController(app).Handler)

	return []rtr.RouteInterface{
		fileManager,
	}, nil
}
