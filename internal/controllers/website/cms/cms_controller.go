package cms

import (
	"net/http"
	"project/internal/registry"
	"project/internal/widgets"
	"project/pkg/webtheme"
	"sync"

	"github.com/dracory/cmsstore"
	cmsFrontend "github.com/dracory/cmsstore/frontend"
	"github.com/dracory/ui"
)

const CMS_ENABLE_CACHE = false

// == CONTROLLER ===============================================================

type cmsController struct {
	frontend cmsFrontend.FrontendInterface
	registry registry.RegistryInterface
}

// == CONSTRUCTOR ==============================================================

func NewCmsController(registry registry.RegistryInterface) *cmsController {
	return &cmsController{registry: registry}
}

// == PUBLIC METHODS ===========================================================

func (controller cmsController) Handler(w http.ResponseWriter, r *http.Request) string {
	instance := GetInstance(controller.registry)
	if instance == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "cms is not configured"
	}
	return instance.StringHandler(w, r)
}

var instance cmsFrontend.FrontendInterface
var once sync.Once

func GetInstance(registry registry.RegistryInterface) cmsFrontend.FrontendInterface {
	once.Do(func() {
		list := widgets.WidgetRegistry(registry)

		shortcodes := []cmsstore.ShortcodeInterface{}
		for _, widget := range list {
			shortcodes = append(shortcodes, widget)
		}

		frontend := cmsFrontend.New(cmsFrontend.Config{
			// BlockEditorDefinitions: webtheme.BlockEditorDefinitions(),
			BlockEditorRenderer: func(blocks []ui.BlockInterface) string {
				return webtheme.New(blocks).ToHtml()
			},
			Store:              registry.GetCmsStore(),
			Shortcodes:         shortcodes,
			Logger:             registry.GetLogger(),
			CacheEnabled:       true,
			CacheExpireSeconds: 1 * 60, // 1 mins
		})

		instance = frontend
	})
	return instance
}
