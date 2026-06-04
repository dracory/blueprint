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

func TestPostManagerController_DeprecatedMethods_PrepareData(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewPostManagerController(app)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/posts?page=1", nil)
	req = req.WithContext(context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user))

	data, err := controller.prepareData(req)
	if err != "" {
		t.Errorf("expected empty error, got %s", err)
	}
	if data.pageInt != 1 {
		t.Errorf("expected pageInt 1, got %d", data.pageInt)
	}
}

func TestPostManagerController_DeprecatedMethods_Page(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewPostManagerController(app)

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
}
