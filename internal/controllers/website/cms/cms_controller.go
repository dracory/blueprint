package cms

import (
	"net/http"
	"project/internal/app"
	"project/internal/widgets"
	"sync"

	"github.com/dracory/base/webtheme"

	"github.com/dracory/cmsstore"
	cmsFrontend "github.com/dracory/cmsstore/frontend"
	"github.com/dracory/ui"
)

const CMS_ENABLE_CACHE = false

// == CONTROLLER ===============================================================

type cmsController struct {
	frontend cmsFrontend.FrontendInterface
	app app.AppInterface
}

// == CONSTRUCTOR ==============================================================

func NewCmsController(app app.AppInterface) *cmsController {
	return &cmsController{app: app}
}

// == PUBLIC METHODS ===========================================================

func (controller cmsController) Handler(w http.ResponseWriter, r *http.Request) string {
	instance := GetInstance(controller.app)
	if instance == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "cms is not configured"
	}
	return instance.StringHandler(w, r)
}

var instance cmsFrontend.FrontendInterface
var once sync.Once

func GetInstance(app app.AppInterface) cmsFrontend.FrontendInterface {
	once.Do(func() {
		list := widgets.WidgetRegistry(app)

		shortcodes := []cmsstore.ShortcodeInterface{}
		for _, widget := range list {
			shortcodes = append(shortcodes, widget)
		}

		frontend := cmsFrontend.New(cmsFrontend.Config{
			// BlockEditorDefinitions: webtheme.BlockEditorDefinitions(),
			BlockEditorRenderer: func(blocks []ui.BlockInterface) string {
				return webtheme.New(blocks).ToHtml()
			},
			Store:              app.GetCmsStore(),
			Shortcodes:         shortcodes,
			Logger:             app.GetLogger(),
			CacheEnabled:       true,
			CacheExpireSeconds: 1 * 60, // 1 mins
			PageNotFoundHandler: func(w http.ResponseWriter, r *http.Request, alias string) (bool, string) {
				return true, "Not found"
			},
		})

		instance = frontend
	})
	return instance
}
