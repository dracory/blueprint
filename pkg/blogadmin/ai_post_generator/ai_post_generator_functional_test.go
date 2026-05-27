package aipostgenerator

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

func TestAiPostGeneratorController_Functional(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithCustomStore(true),
	)

	user, _ := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	controller := NewAiPostGeneratorController(registry)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	t.Run("renderPage", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-generator", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "Post Generator") {
			t.Error("expected Post Generator in response")
		}
	})
}
