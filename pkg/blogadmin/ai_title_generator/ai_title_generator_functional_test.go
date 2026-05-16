package aititlegenerator

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

func TestAiTitleGeneratorController_Functional(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithCustomStore(true),
	)

	user, _ := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	controller := NewAiTitleGeneratorController(registry)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	t.Run("renderPage", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-title-generator", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		assert.Contains(t, resp, "AI Title Generator")
	})

	t.Run("onAddTitleModal", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-title-generator?action="+ACTION_ADD_TITLE, nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		assert.Contains(t, resp, "Add Custom Title")
	})
}
