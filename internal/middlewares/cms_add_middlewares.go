package middlewares

import (
	"net/http"
	"project/internal/registry"

	"github.com/dracory/base/cfmt"
	"github.com/dracory/cmsstore"
)

func CmsAddMiddlewares(registry registry.RegistryInterface) {
	if !registry.GetConfig().GetCmsStoreUsed() {
		return
	}

	if registry.GetCmsStore() == nil {
		return
	}

	helloMiddleware := cmsstore.Middleware().
		SetIdentifier("HelloMiddleware").
		SetName("HelloMiddleware").
		SetType(cmsstore.MIDDLEWARE_TYPE_BEFORE).
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cfmt.Infoln("Hello from Middleware")
				next.ServeHTTP(w, r)
			})
		})
	afterMiddleware := cmsstore.Middleware().
		SetIdentifier("CmsLayoutMiddleware").
		SetName("Cms Layout Middleware").
		SetType(cmsstore.MIDDLEWARE_TYPE_AFTER).
		SetHandler(NewCmsLayoutMiddleware(registry).GetHandler())

	registry.GetCmsStore().AddMiddleware(helloMiddleware)
	registry.GetCmsStore().AddMiddleware(afterMiddleware)
}
