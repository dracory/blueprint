package fileadmin

import (
	"errors"
	"net/http"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	"project/pkg/fileadmin/file_manager"
	"project/pkg/fileadmin/shared"
)

func Routes(app app.AppInterface) ([]rtr.RouteInterface, error) {
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}

	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_FILE_MANAGER:
			return file_manager.NewFileManagerController(app).Handler(w, r)
		}

		// Default to file manager
		return file_manager.NewFileManagerController(app).Handler(w, r)
	}

	fileManager := rtr.NewRoute().
		SetName("Admin > File Manager").
		SetPath(links.ADMIN_FILE_MANAGER).
		SetHTMLHandler(handler)

	fileManagerCatchAll := rtr.NewRoute().
		SetName("Admin > File Manager > Catchall").
		SetPath(links.ADMIN_FILE_MANAGER + links.CATCHALL).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		fileManager,
		fileManagerCatchAll,
	}, nil
}
