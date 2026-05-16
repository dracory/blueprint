package post_manager

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestPostManagerController_DeprecatedMethods(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	controller := NewPostManagerController(registry)

	t.Run("prepareData", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/posts?page=1", nil)
		req = req.WithContext(context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user))

		data, err := controller.prepareData(req)
		assert.Empty(t, err)
		assert.Equal(t, 1, data.pageInt)
	})

	t.Run("page", func(t *testing.T) {
		data := postManagerControllerData{
			pageInt: 0,
			perPage: 10,
		}
		html := controller.page(data)
		assert.NotNil(t, html)
		assert.Contains(t, html.ToHTML(), "Post Manager")
	})
}
