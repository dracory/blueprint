package admin

import (
	"net/http"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	aiPostContentUpdate "project/internal/controllers/admin/blog/ai_post_content_update"
	aiPostEditor "project/internal/controllers/admin/blog/ai_post_editor"
	aiPostGenerator "project/internal/controllers/admin/blog/ai_post_generator"
	aiTest "project/internal/controllers/admin/blog/ai_test"
	aiTitleGenerator "project/internal/controllers/admin/blog/ai_title_generator"
	aiTools "project/internal/controllers/admin/blog/ai_tools"
	blog_settings "project/internal/controllers/admin/blog/blog_settings"
	"project/internal/controllers/admin/blog/post_create"
	"project/internal/controllers/admin/blog/post_delete"
	"project/internal/controllers/admin/blog/post_manager"
	"project/internal/controllers/admin/blog/post_update"
	"project/internal/controllers/admin/blog/shared"
)

func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
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
			return blog_settings.NewBlogSettingsController(registry).Handler(w, r)
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
		}

		// Default to post manager
		return post_manager.NewPostManagerController(registry).Handler(w, r)
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
	}
}
