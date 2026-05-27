package aititlegenerator

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

func TestAiTitleGeneratorController_Functional_RenderPage(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-title-generator", nil).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "AI Title Generator") {
		t.Error("expected AI Title Generator in response")
	}
}

func TestAiTitleGeneratorController_Functional_OnAddTitleModal(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-title-generator?action="+ACTION_ADD_TITLE, nil).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "Add Custom Title") {
		t.Error("expected Add Custom Title in response")
	}
}
