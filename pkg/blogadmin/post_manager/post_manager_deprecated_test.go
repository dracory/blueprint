package post_manager

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/test"
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
		if err != "" {
			t.Errorf("expected empty error, got %s", err)
		}
		if data.pageInt != 1 {
			t.Errorf("expected pageInt 1, got %d", data.pageInt)
		}
	})

	t.Run("page", func(t *testing.T) {
		data := postManagerControllerData{
			pageInt: 0,
			perPage: 10,
		}
		html := controller.page(data)
		if html == nil {
			t.Error("expected html to not be nil")
		}
		if !strings.Contains(html.ToHTML(), "Post Manager") {
			t.Error("expected Post Manager in HTML")
		}
	})
}
