package cms

import (
	"net/http"
	"project/app/widgets"
	"project/config"
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
}

// == CONSTRUCTOR ==============================================================

func NewCmsController() *cmsController {
	return &cmsController{}
}

// == PUBLIC METHODS ===========================================================

func (controller cmsController) Handler(w http.ResponseWriter, r *http.Request) string {
	instance := GetInstance()
	if instance == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "cms is not configured"
	}
	return instance.StringHandler(w, r)
}

var instance cmsFrontend.FrontendInterface
var once sync.Once

func GetInstance() cmsFrontend.FrontendInterface {
	once.Do(func() {
		list := widgets.WidgetRegistry()

		shortcodes := []cmsstore.ShortcodeInterface{}
		for _, widget := range list {
			shortcodes = append(shortcodes, widget)
		}

		frontend := cmsFrontend.New(cmsFrontend.Config{
			// BlockEditorDefinitions: webtheme.BlockEditorDefinitions(),
			BlockEditorRenderer: func(blocks []ui.BlockInterface) string {
				return webtheme.New(blocks).ToHtml()
			},
			Store:              config.CmsStore,
			Logger:             &config.Logger,
			CacheEnabled:       true,
			CacheExpireSeconds: 1 * 60, // 1 mins
		})

		instance = frontend
	})
	return instance
}
