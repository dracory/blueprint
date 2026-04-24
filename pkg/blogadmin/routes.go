package blogadmin

import (
	"errors"
	"net/http"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	aiPostContentUpdate "project/pkg/blogadmin/ai_post_content_update"
	aiPostEditor "project/pkg/blogadmin/ai_post_editor"
	aiPostGenerator "project/pkg/blogadmin/ai_post_generator"
	aiTest "project/pkg/blogadmin/ai_test"
	aiTitleGenerator "project/pkg/blogadmin/ai_title_generator"
	aiTools "project/pkg/blogadmin/ai_tools"
	blogSettings "project/pkg/blogadmin/blog_settings"
	"project/pkg/blogadmin/category_manager"
	"project/pkg/blogadmin/dashboard"
	"project/pkg/blogadmin/post_create"
	"project/pkg/blogadmin/post_delete"
	"project/pkg/blogadmin/post_manager"
	"project/pkg/blogadmin/post_update"
	"project/pkg/blogadmin/shared"
	"project/pkg/blogadmin/tag_manager"
)

func Routes(registry registry.RegistryInterface) ([]rtr.RouteInterface, error) {
	if registry == nil {
		return nil, errors.New("registry cannot be nil")
	}
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_HOME:
			return post_manager.NewPostManagerController(registry).Handler(w, r)
		case shared.CONTROLLER_POST_CREATE:
			return post_create.NewPostCreateController(registry).Handler(w, r)
		case shared.CONTROLLER_POST_DELETE:
			return post_delete.NewPostDeleteController(registry).Handler(w, r)
		case shared.CONTROLLER_POST_MANAGER:
			return post_manager.NewPostManagerController(registry).Handler(w, r)
		case shared.CONTROLLER_POST_UPDATE:
			return post_update.NewPostUpdateController(registry).Handler(w, r)
		case shared.CONTROLLER_AI_TOOLS:
			return aiTools.NewAiToolsController(registry).Handler(w, r)
		case shared.CONTROLLER_BLOG_SETTINGS:
			return blogSettings.NewBlogSettingsController(registry).Handler(w, r)
		case shared.CONTROLLER_AI_POST_CONTENT_UPDATE:
			return aiPostContentUpdate.NewController(registry).Handler(w, r)
		case shared.CONTROLLER_AI_POST_GENERATOR:
			return aiPostGenerator.NewAiPostGeneratorController(registry).Handler(w, r)
		case shared.CONTROLLER_AI_TITLE_GENERATOR:
			return aiTitleGenerator.NewAiTitleGeneratorController(registry).Handler(w, r)
		case shared.CONTROLLER_AI_POST_EDITOR:
			return aiPostEditor.NewAiPostEditorController(registry).Handler(w, r)
		case shared.CONTROLLER_AI_TEST:
			return aiTest.NewAiTestController(registry).Handler(w, r)
		case shared.CONTROLLER_DASHBOARD:
			return dashboard.NewDashboardController(registry).Handler(w, r)
		case shared.CONTROLLER_CATEGORY_MANAGER:
			return category_manager.NewCategoryManagerController(registry).Handler(w, r)
		case shared.CONTROLLER_TAG_MANAGER:
			return tag_manager.NewTagManagerController(registry).Handler(w, r)
		}

		// Default to dashboard
		return dashboard.NewDashboardController(registry).Handler(w, r)
	}

	blog := rtr.NewRoute().
		SetName("Admin > Blog").
		SetPath(links.ADMIN_BLOG).
		SetHTMLHandler(handler)

	blogCatchAll := rtr.NewRoute().
		SetName("Admin > Blog > Catchall").
		SetPath(links.ADMIN_BLOG + links.CATCHALL).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		blog,
		blogCatchAll,
	}, nil
}
