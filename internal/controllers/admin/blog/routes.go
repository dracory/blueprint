package admin

import (
	"net/http"
	"project/internal/links"
	"project/internal/types"

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
	"project/internal/controllers/admin/blog/post_update_v1"
	"project/internal/controllers/admin/blog/shared"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_HOME:
			return post_manager.NewPostManagerController(app).Handler(w, r)
		case shared.CONTROLLER_POST_CREATE:
			return post_create.NewPostCreateController(app).Handler(w, r)
		case shared.CONTROLLER_POST_DELETE:
			return post_delete.NewPostDeleteController(app).Handler(w, r)
		case shared.CONTROLLER_POST_MANAGER:
			return post_manager.NewPostManagerController(app).Handler(w, r)
		case shared.CONTROLLER_POST_UPDATE_V1:
			return post_update_v1.NewPostUpdateController(app).Handler(w, r)
		case shared.CONTROLLER_POST_UPDATE:
			return post_update.NewPostUpdateController(app).Handler(w, r)
		case shared.CONTROLLER_AI_TOOLS:
			return aiTools.NewAiToolsController(app).Handler(w, r)
		case shared.CONTROLLER_BLOG_SETTINGS:
			return blog_settings.NewBlogSettingsController(app).Handler(w, r)
		case shared.CONTROLLER_AI_POST_CONTENT_UPDATE:
			return aiPostContentUpdate.NewController(app).Handler(w, r)
		case shared.CONTROLLER_AI_POST_GENERATOR:
			return aiPostGenerator.NewAiPostGeneratorController(app).Handler(w, r)
		case shared.CONTROLLER_AI_TITLE_GENERATOR:
			return aiTitleGenerator.NewAiTitleGeneratorController(app).Handler(w, r)
		case shared.CONTROLLER_AI_POST_EDITOR:
			return aiPostEditor.NewAiPostEditorController(app).Handler(w, r)
		case shared.CONTROLLER_AI_TEST:
			return aiTest.NewAiTestController(app).Handler(w, r)
		}

		// Default to post manager
		return post_manager.NewPostManagerController(app).Handler(w, r)
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
