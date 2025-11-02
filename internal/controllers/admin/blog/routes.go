package admin

import (
	"net/http"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	aiPostEditor "project/internal/controllers/admin/blog/ai_post_editor"
	aiPostGenerator "project/internal/controllers/admin/blog/ai_post_generator"
	aiTest "project/internal/controllers/admin/blog/ai_test"
	aiTitleGenerator "project/internal/controllers/admin/blog/ai_title_generator"
	aiTools "project/internal/controllers/admin/blog/ai_tools"
	blogSettings "project/internal/controllers/admin/blog/blog_settings"
	postCreate "project/internal/controllers/admin/blog/post_create"
	postDelete "project/internal/controllers/admin/blog/post_delete"
	postManager "project/internal/controllers/admin/blog/post_manager"
	postUpdate "project/internal/controllers/admin/blog/post_update"
	"project/internal/controllers/admin/blog/shared"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_HOME:
			return postManager.NewPostManagerController(app).Handler(w, r)
		case shared.CONTROLLER_POST_CREATE:
			return postCreate.NewPostCreateController(app).Handler(w, r)
		case shared.CONTROLLER_POST_DELETE:
			return postDelete.NewPostDeleteController(app).Handler(w, r)
		case shared.CONTROLLER_POST_MANAGER:
			return postManager.NewPostManagerController(app).Handler(w, r)
		case shared.CONTROLLER_POST_UPDATE:
			return postUpdate.NewPostUpdateController(app).Handler(w, r)
		case shared.CONTROLLER_AI_TOOLS:
			return aiTools.NewAiToolsController(app).Handler(w, r)
		case shared.CONTROLLER_BLOG_SETTINGS:
			return blogSettings.NewBlogSettingsController(app).Handler(w, r)
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
		return postManager.NewPostManagerController(app).Handler(w, r)
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
