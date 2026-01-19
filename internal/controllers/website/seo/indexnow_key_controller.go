package seo

import (
	"net/http"
	"project/internal/registry"
)

type indexNowKeyController struct {
	registry registry.RegistryInterface
}

func NewIndexNowKeyController(registry registry.RegistryInterface) *indexNowKeyController {
	return &indexNowKeyController{
		registry: registry,
	}
}

func (c indexNowKeyController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "text/plain")
	return c.registry.GetConfig().GetIndexNowKey()
}
