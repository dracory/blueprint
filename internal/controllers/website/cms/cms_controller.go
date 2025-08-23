package cms

import (
	"net/http"
	"project/internal/types"
	"project/internal/widgets"
	"project/pkg/webtheme"
	"sync"

	"github.com/gouniverse/cmsstore"
	cmsFrontend "github.com/gouniverse/cmsstore/frontend"
	"github.com/gouniverse/ui"
)

const CMS_ENABLE_CACHE = false

// == CONTROLLER ===============================================================

type cmsController struct {
	frontend cmsFrontend.FrontendInterface
	app      types.AppInterface
}

// == CONSTRUCTOR ==============================================================

func NewCmsController(app types.AppInterface) *cmsController {
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

func GetInstance(app types.AppInterface) cmsFrontend.FrontendInterface {
	once.Do(func() {
		list := widgets.WidgetRegistry(app.GetConfig())

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
			Logger:             app.GetLogger(),
			CacheEnabled:       true,
			CacheExpireSeconds: 1 * 60, // 1 mins
		})

		instance = frontend
	})
	return instance
}
